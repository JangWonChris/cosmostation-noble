package exporter

import (
	"strconv"
	"sync"
	"time"

	mdschema "github.com/cosmostation/mintscan-database/schema"

	//cosmos-sdk
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	tmcTypes "github.com/tendermint/tendermint/rpc/core/types"

	"go.uber.org/zap"
)

var muProp sync.RWMutex
var propList = make(map[uint64]struct{})

func (ex *Exporter) watchLiveProposals() {
	for {
		p, err := ex.DB.GetLiveProposalIDs()
		if err != nil {
			zap.S().Info("failed to get live proposals")
			time.Sleep(2 * time.Second)
			continue
		}
		muProp.Lock()
		for i := range p {
			_, ok := propList[p[i].ID]
			if !ok {
				propList[p[i].ID] = struct{}{}
			}
		}
		muProp.Unlock()
		zap.S().Info("proposal list updated")
		time.Sleep(6 * time.Second)
	}
}

func (ex *Exporter) updateProposals() {
	for {
		if !ex.App.CatchingUp {
			zap.S().Info("start updating proposals : ", propList)
			muProp.RLock()
			for id := range propList {
				if err := ex.updateProposal(id); err != nil {
					continue
				}
				delete(propList, id)
			}
			muProp.RUnlock()
			zap.S().Info("finish update proposals : ", propList)
		} else {
			zap.S().Info("pending update proposals, app is catching up")
		}
		time.Sleep(10 * time.Second)
	}
}

// getGovernance returns governance by decoding governance related transactions in a block.
func (ex *Exporter) getGovernance(block *tmcTypes.ResultBlock, txResp []*sdkTypes.TxResponse) ([]mdschema.Proposal, []mdschema.Deposit, []mdschema.Vote, error) {
	proposals := make([]mdschema.Proposal, 0)
	deposits := make([]mdschema.Deposit, 0)
	votes := make([]mdschema.Vote, 0)

	if len(txResp) <= 0 {
		return proposals, deposits, votes, nil
	}

	for _, tx := range txResp {
		// code == 0 이면, 오류 트랜잭션이다.
		if tx.Code != 0 {
			continue
		}

		msgs := tx.GetTx().GetMsgs()

		for _, msg := range msgs {

			switch m := msg.(type) {
			case *govtypes.MsgSubmitProposal:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				// Get proposal id for this proposal.
				// Handle case of multiple messages which has multiple events and attributes.
				var proposalID uint64
				for _, log := range tx.Logs {
					for _, event := range log.Events {
						if event.Type == "submit_proposal" {
							for _, attribute := range event.Attributes {
								if attribute.Key == "proposal_id" {
									proposalID, _ = strconv.ParseUint(attribute.Value, 10, 64)
								}
							}
						}
					}
				}

				var initialDepositAmount string
				var initialDepositDenom string

				if len(m.InitialDeposit) > 0 {
					initialDepositAmount = m.InitialDeposit[0].Amount.String()
					initialDepositDenom = m.InitialDeposit[0].Denom
				}

				p := mdschema.Proposal{
					ID:                   proposalID,
					TxHash:               tx.TxHash,
					Proposer:             m.Proposer,
					InitialDepositAmount: initialDepositAmount,
					InitialDepositDenom:  initialDepositDenom,
				}

				proposals = append(proposals, p)

				d := mdschema.Deposit{
					Height:     tx.Height,
					ProposalID: proposalID,
					Depositor:  m.Proposer,
					Amount:     initialDepositAmount,
					Denom:      initialDepositDenom,
					TxHash:     tx.TxHash,
					GasWanted:  tx.GasWanted,
					GasUsed:    tx.GasUsed,
					Timestamp:  block.Block.Header.Time,
				}

				deposits = append(deposits, d)

				go ex.ProposalNotificationToSlack(p.ID)

			case *govtypes.MsgDeposit:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				var amount string
				var denom string

				if len(m.Amount) > 0 {
					amount = m.Amount[0].Amount.String()
					denom = m.Amount[0].Denom
				}

				d := mdschema.Deposit{
					Height:     tx.Height,
					ProposalID: m.ProposalId,
					Depositor:  m.Depositor,
					Amount:     amount,
					Denom:      denom,
					TxHash:     tx.TxHash,
					GasWanted:  tx.GasWanted,
					GasUsed:    tx.GasUsed,
					Timestamp:  block.Block.Header.Time,
				}

				deposits = append(deposits, d)

				go ex.ProposalNotificationToSlack(d.ProposalID)

			case *govtypes.MsgVote:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				v := mdschema.Vote{
					Height:     tx.Height,
					ProposalID: m.ProposalId,
					Voter:      m.Voter,
					Option:     m.Option.String(),
					TxHash:     tx.TxHash,
					GasWanted:  tx.GasWanted,
					GasUsed:    tx.GasUsed,
					Timestamp:  block.Block.Header.Time,
				}

				votes = append(votes, v)

			default:
				continue
			}
		}
	}

	return proposals, deposits, votes, nil
}

// saveProposals saves all governance proposals
func (ex *Exporter) saveAllProposals() {
	NodePropCount, err := ex.Client.GetNumberofProposals()
	if err != nil {
		zap.S().Errorf("failed to get number of proposal from DB: %s", err)
		return
	}
	DBPropCount, err := ex.DB.GetNumberofValidProposal()
	if err != nil {
		zap.S().Errorf("failed to get number of proposal from Node: %s", err)
		return
	}
	// database에 저장된 프로포절의 수와 노드의 수가 같으면 업데이트 하지 않는다.
	if NodePropCount == uint64(DBPropCount) {
		zap.S().Info("skip saveAllProposals, all proposals have already been stored in database, count : ", NodePropCount)
		return
	}

	proposals, err := ex.Client.GetAllProposals()
	if err != nil {
		zap.S().Errorf("failed to get proposals: %s", err)
		return
	}

	if len(proposals) <= 0 {
		zap.S().Info("found empty proposals")
		return
	}

	err = ex.DB.InsertOrUpdateProposals(proposals)
	if err != nil {
		zap.S().Errorf("failed to insert or update proposal: %s", err)
		return
	}
}

// saveLiveProposals saves live governance proposals
func (ex *Exporter) saveLiveProposals() {
	vp, err := ex.Client.GetProposalsByStatus(govtypes.StatusVotingPeriod)
	if err != nil {
		zap.S().Errorf("failed to get proposals: %s", err)
		return
	}
	// 2022.01.17 deposit period 인 프로포절을 가져오는데 행이 걸림
	dp, err := ex.Client.GetProposalsByStatus(govtypes.StatusDepositPeriod)
	if err != nil {
		zap.S().Errorf("failed to get proposals: %s", err)
		return
	}
	proposals := make([]mdschema.Proposal, 0)
	proposals = append(proposals, vp...)
	proposals = append(proposals, dp...)

	if len(proposals) <= 0 {
		zap.S().Info("found empty proposals")
		return
	}

	err = ex.DB.InsertOrUpdateProposals(proposals)
	if err != nil {
		zap.S().Errorf("failed to insert or update proposal: %s", err)
		return
	}
}

// updateProposal update proposal which is passed voting end time
func (ex *Exporter) updateProposal(id uint64) error {
	p, err := ex.Client.GetProposal(id)
	if err != nil {
		zap.S().Errorf("failed to get proposal: %s", err)
		return err
	}

	return ex.DB.InsertOrUpdateProposal(p)
}
