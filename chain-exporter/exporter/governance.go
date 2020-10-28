package exporter

import (
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	tmcTypes "github.com/tendermint/tendermint/rpc/core/types"

	"go.uber.org/zap"
)

// getGovernance returns governance by decoding governance related transactions in a block.
func (ex *Exporter) getGovernance(block *tmcTypes.ResultBlock, txResp []*sdkTypes.TxResponse) ([]schema.Proposal, []schema.Deposit, []schema.Vote, error) {
	proposals := make([]schema.Proposal, 0)
	deposits := make([]schema.Deposit, 0)
	votes := make([]schema.Vote, 0)

	if len(txResp) <= 0 {
		return proposals, deposits, votes, nil
	}

	for _, tx := range txResp {
		// Other than code equals to 0, it is failed transaction.
		if tx.Code != 0 {
			return proposals, deposits, votes, nil
		}

		// stdTx, ok := tx.Tx.(auth.StdTx)
		// if !ok {
		// 	return proposals, deposits, votes, fmt.Errorf("unsupported tx type: %s", tx.Tx)
		// }

		msgs := tx.GetTx().GetMsgs()

		for _, msg := range msgs {

			switch m := msg.(type) {
			case *govtypes.MsgSubmitProposal:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				// msgSubmitProposal := m.(gov.MsgSubmitProposal)

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

				p := schema.NewProposal(schema.Proposal{
					ID:                   int64(proposalID),
					TxHash:               tx.TxHash,
					Proposer:             m.Proposer,
					InitialDepositAmount: initialDepositAmount,
					InitialDepositDenom:  initialDepositDenom,
				})

				proposals = append(proposals, *p)

				d := schema.NewDeposit(schema.Deposit{
					Height:     tx.Height,
					ProposalID: proposalID,
					Depositor:  m.Proposer,
					Amount:     initialDepositAmount,
					Denom:      initialDepositDenom,
					TxHash:     tx.TxHash,
					GasWanted:  tx.GasWanted,
					GasUsed:    tx.GasUsed,
					Timestamp:  block.Block.Header.Time,
				})

				deposits = append(deposits, *d)

			case *govtypes.MsgDeposit:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				// msgDeposit := m.(gov.MsgDeposit)

				var amount string
				var denom string

				if len(m.Amount) > 0 {
					amount = m.Amount[0].Amount.String()
					denom = m.Amount[0].Denom
				}

				d := schema.NewDeposit(schema.Deposit{
					Height:     tx.Height,
					ProposalID: m.ProposalId,
					Depositor:  m.Depositor,
					Amount:     amount,
					Denom:      denom,
					TxHash:     tx.TxHash,
					GasWanted:  tx.GasWanted,
					GasUsed:    tx.GasUsed,
					Timestamp:  block.Block.Header.Time,
				})

				deposits = append(deposits, *d)

			case *govtypes.MsgVote:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				// msgVote := m.(gov.MsgVote)

				v := schema.NewVote(schema.Vote{
					Height:     tx.Height,
					ProposalID: m.ProposalId,
					Voter:      m.Voter,
					Option:     m.Option.String(),
					TxHash:     tx.TxHash,
					GasWanted:  tx.GasWanted,
					GasUsed:    tx.GasUsed,
					Timestamp:  block.Block.Header.Time,
				})

				votes = append(votes, *v)

			default:
				continue
			}
		}
	}

	return proposals, deposits, votes, nil
}

// saveProposals saves all governance proposals
func (ex *Exporter) saveProposals() {
	proposals, err := ex.client.GetProposals()
	if err != nil {
		zap.S().Errorf("failed to get proposals: %s", err)
		return
	}

	if len(proposals) <= 0 {
		zap.S().Info("found empty proposals")
		return
	}

	err = ex.db.InsertOrUpdateProposals(proposals)
	if err != nil {
		zap.S().Errorf("failed to insert or update proposal: %s", err)
		return
	}
}
