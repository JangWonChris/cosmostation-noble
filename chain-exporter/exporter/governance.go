package exporter

import (
	"fmt"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"go.uber.org/zap"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/gov"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getGovernance returns governance by decoding governance related transactions in a block.
func (ex *Exporter) getGovernance(block *tmctypes.ResultBlock, txs []*sdk.TxResponse) ([]schema.Proposal, []schema.Deposit, []schema.Vote, error) {
	proposals := make([]schema.Proposal, 0)
	deposits := make([]schema.Deposit, 0)
	votes := make([]schema.Vote, 0)

	if len(txs) <= 0 {
		return proposals, deposits, votes, nil
	}

	for _, tx := range txs {
		// Other than code equals to 0, it is failed transaction.
		if tx.Code != 0 {
			return proposals, deposits, votes, nil
		}

		stdTx, ok := tx.Tx.(auth.StdTx)
		if !ok {
			return proposals, deposits, votes, fmt.Errorf("unsupported tx type: %s", tx.Tx)
		}

		switch stdTx.Msgs[0].(type) {
		case gov.MsgSubmitProposal:
			zap.S().Info(zap.Any("MsgType ", stdTx.Msgs[0].Type()), zap.Any("Hash ", tx.TxHash))

			msgSubmitProposal := stdTx.Msgs[0].(gov.MsgSubmitProposal)

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

			if len(msgSubmitProposal.InitialDeposit) > 0 {
				initialDepositAmount = msgSubmitProposal.InitialDeposit[0].Amount.String()
				initialDepositDenom = msgSubmitProposal.InitialDeposit[0].Denom
			}

			p := schema.NewProposal(schema.Proposal{
				ID:                   int64(proposalID),
				TxHash:               tx.TxHash,
				Proposer:             msgSubmitProposal.Proposer.String(),
				InitialDepositAmount: initialDepositAmount,
				InitialDepositDenom:  initialDepositDenom,
			})

			proposals = append(proposals, *p)

			d := schema.NewDeposit(schema.Deposit{
				Height:     tx.Height,
				ProposalID: proposalID,
				Depositor:  msgSubmitProposal.Proposer.String(),
				Amount:     initialDepositAmount,
				Denom:      initialDepositDenom,
				TxHash:     tx.TxHash,
				GasWanted:  tx.GasWanted,
				GasUsed:    tx.GasUsed,
				Timestamp:  block.Block.Header.Time,
			})

			deposits = append(deposits, *d)

		case gov.MsgDeposit:
			zap.S().Info(zap.Any("MsgType ", stdTx.Msgs[0].Type()), zap.Any("Hash ", tx.TxHash))

			msgDeposit := stdTx.Msgs[0].(gov.MsgDeposit)

			var amount string
			var denom string

			if len(msgDeposit.Amount) > 0 {
				amount = msgDeposit.Amount[0].Amount.String()
				denom = msgDeposit.Amount[0].Denom
			}

			d := schema.NewDeposit(schema.Deposit{
				Height:     tx.Height,
				ProposalID: msgDeposit.ProposalID,
				Depositor:  msgDeposit.Depositor.String(),
				Amount:     amount,
				Denom:      denom,
				TxHash:     tx.TxHash,
				GasWanted:  tx.GasWanted,
				GasUsed:    tx.GasUsed,
				Timestamp:  block.Block.Header.Time,
			})

			deposits = append(deposits, *d)

		case gov.MsgVote:
			zap.S().Info(zap.Any("MsgType ", stdTx.Msgs[0].Type()), zap.Any("Hash ", tx.TxHash))

			msgVote := stdTx.Msgs[0].(gov.MsgVote)

			v := schema.NewVote(schema.Vote{
				Height:     tx.Height,
				ProposalID: msgVote.ProposalID,
				Voter:      msgVote.Voter.String(),
				Option:     msgVote.Option.String(),
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
