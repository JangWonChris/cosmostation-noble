package db

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/go-pg/pg"
)

// InsertProposal saves on-chain proposals getting from /gov/proposals REST API
func (db *Database) InsertProposal(data *schema.Proposal) (bool, error) {
	err := db.Insert(data)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// InsertOrUpdateValidators updates the given validator set
func (db *Database) InsertOrUpdateValidators(data []*schema.Validator) (bool, error) {
	_, err := db.Model(&data).
		OnConflict("(operator_address) DO UPDATE").
		Set("rank = EXCLUDED.rank").
		Set("consensus_pubkey = EXCLUDED.consensus_pubkey").
		Set("proposer = EXCLUDED.proposer").
		Set("jailed = EXCLUDED.jailed").
		Set("status = EXCLUDED.status").
		Set("tokens = EXCLUDED.tokens").
		Set("delegator_shares = EXCLUDED.delegator_shares").
		Set("moniker = EXCLUDED.moniker").
		Set("identity = EXCLUDED.identity").
		Set("website = EXCLUDED.website").
		Set("details = EXCLUDED.details").
		Set("unbonding_height = EXCLUDED.unbonding_height").
		Set("unbonding_time = EXCLUDED.unbonding_time").
		Set("commission_rate = EXCLUDED.commission_rate").
		Set("commission_max_rate = EXCLUDED.commission_max_rate").
		Set("commission_change_rate = EXCLUDED.commission_change_rate").
		Set("update_time = EXCLUDED.update_time").
		Set("min_self_delegation = EXCLUDED.min_self_delegation").
		Insert()
	if err != nil {
		return false, nil
	}
	return true, nil
}

// InsertExportedData saves exported blockchain data
// if function returns an error transaction is rollbacked, otherwise transaction is committed.
func (db *Database) InsertExportedData(block []*schema.BlockCosmoshub3, evidence []*schema.Evidence, genesisValSet []*schema.PowerEventHistory,
	missingBlocks []*schema.Miss, accumMissingBlocks []*schema.Miss, missingBlocksDetail []*schema.MissDetail, txs []*schema.TxCosmoshub3,
	txIndex []*schema.TxIndex, votes []*schema.Vote, deposits []*schema.Deposit, proposals []*schema.Proposal, powerEventHistory []*schema.PowerEventHistory) error {

	err := db.RunInTransaction(func(tx *pg.Tx) error {
		if len(block) > 0 {
			err := tx.Insert(&block)
			if err != nil {
				return err
			}
		}

		if len(evidence) > 0 {
			err := tx.Insert(&evidence)
			if err != nil {
				return err
			}
		}

		if len(genesisValSet) > 0 {
			err := tx.Insert(&genesisValSet)
			if err != nil {
				return err
			}
		}

		if len(powerEventHistory) > 0 {
			err := tx.Insert(&powerEventHistory)
			if err != nil {
				return err
			}
		}

		if len(missingBlocks) > 0 {
			err := tx.Insert(&missingBlocks)
			if err != nil {
				return err
			}
		}

		if len(missingBlocksDetail) > 0 {
			err := tx.Insert(&missingBlocksDetail)
			if err != nil {
				return err
			}
		}

		if len(txs) > 0 {
			err := tx.Insert(&txs)
			if err != nil {
				return err
			}
		}

		if len(txIndex) > 0 {
			err := tx.Insert(&txIndex)
			if err != nil {
				return err
			}
		}

		if len(deposits) > 0 {
			err := tx.Insert(&deposits)
			if err != nil {
				return err
			}
		}

		var tempMiss schema.Miss
		if len(accumMissingBlocks) > 0 {
			for i := 0; i < len(accumMissingBlocks); i++ {
				_, err := tx.Model(&tempMiss).
					Set("address = ?", accumMissingBlocks[i].Address).
					Set("start_height = ?", accumMissingBlocks[i].StartHeight).
					Set("end_height = ?", accumMissingBlocks[i].EndHeight).
					Set("missing_count = ?", accumMissingBlocks[i].MissingCount).
					Set("start_time = ?", accumMissingBlocks[i].StartTime).
					Set("end_time = ?", block[0].Timestamp).
					Where("end_height = ? AND address = ?", accumMissingBlocks[i].EndHeight-int64(1), accumMissingBlocks[i].Address).
					Update()
				if err != nil {
					return err
				}
			}
		}

		if len(votes) > 0 {
			var tempVotes schema.Vote
			for _, vote := range votes {
				// Check if a validator already voted
				count, _ := tx.Model(&tempVotes).
					Where("proposal_id = ? AND voter = ?", vote.ProposalID, vote.Voter).
					Count()
				if count > 0 {
					_, err := tx.Model(&tempVotes).
						Set("height = ?", vote.Height).
						Set("option = ?", vote.Option).
						Set("tx_hash = ?", vote.TxHash).
						Set("gas_wanted = ?", vote.GasWanted).
						Set("gas_used = ?", vote.GasUsed).
						Set("timestamp = ?", vote.Timestamp).
						Where("proposal_id = ? AND voter = ?", vote.ProposalID, vote.Voter).
						Update()
					if err != nil {
						return err
					}
				} else {
					err := tx.Insert(&votes)
					if err != nil {
						return err
					}
				}
			}
		}

		if len(proposals) > 0 {
			var tempProposal schema.Proposal
			for i := 0; i < len(proposals); i++ {
				// check if a validator already voted
				count, _ := tx.Model(&tempProposal).
					Where("id = ?", proposals[i].ID).
					Count()

				if count > 0 {
					// save and update proposal
					_, err := tx.Model(&tempProposal).
						Set("tx_hash = ?", proposals[i].TxHash).
						Set("proposer = ?", proposals[i].Proposer).
						Set("initial_deposit_amount = ?", proposals[i].InitialDepositAmount).
						Set("initial_deposit_denom = ?", proposals[i].InitialDepositDenom).
						Where("id = ?", proposals[i].ID).
						Update()
					if err != nil {
						return err
					}
				} else {
					err := tx.Insert(&proposals)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	// roll back
	if err != nil {
		return err
	}

	return nil
}
