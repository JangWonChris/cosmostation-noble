package db

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/go-pg/pg"
)

// SaveExportedData saves exported blockchain data
// if function returns an error transaction is rollbacked, otherwise transaction is committed.
func (db *Database) SaveExportedData(blockInfo []*schema.BlockInfo, evidenceInfo []*schema.EvidenceInfo, genesisValidatorsInfo []*schema.ValidatorSetInfo,
	missInfo []*schema.MissInfo, accumMissInfo []*schema.MissInfo, missDetailInfo []*schema.MissDetailInfo, transactionInfo []*schema.TransactionInfo,
	voteInfo []*schema.VoteInfo, depositInfo []*schema.DepositInfo, proposalInfo []*schema.ProposalInfo, validatorSetInfo []*schema.ValidatorSetInfo) error {

	err := db.RunInTransaction(func(tx *pg.Tx) error {
		if len(blockInfo) > 0 {
			err := tx.Insert(&blockInfo)
			if err != nil {
				return err
			}
		}

		if len(genesisValidatorsInfo) > 0 {
			err := tx.Insert(&genesisValidatorsInfo)
			if err != nil {
				return err
			}
		}

		if len(validatorSetInfo) > 0 {
			err := tx.Insert(&validatorSetInfo)
			if err != nil {
				return err
			}
		}

		if len(evidenceInfo) > 0 {
			err := tx.Insert(&evidenceInfo)
			if err != nil {
				return err
			}
		}

		if len(missInfo) > 0 {
			err := tx.Insert(&missInfo)
			if err != nil {
				return err
			}
		}

		if len(missDetailInfo) > 0 {
			err := tx.Insert(&missDetailInfo)
			if err != nil {
				return err
			}
		}

		if len(transactionInfo) > 0 {
			err := tx.Insert(&transactionInfo)
			if err != nil {
				return err
			}
		}

		if len(depositInfo) > 0 {
			err := tx.Insert(&depositInfo)
			if err != nil {
				return err
			}
		}

		var tempMissInfo schema.MissInfo
		if len(accumMissInfo) > 0 {
			for i := 0; i < len(accumMissInfo); i++ {
				_, err := tx.Model(&tempMissInfo).
					Set("address = ?", accumMissInfo[i].Address).
					Set("start_height = ?", accumMissInfo[i].StartHeight).
					Set("end_height = ?", accumMissInfo[i].EndHeight).
					Set("missing_count = ?", accumMissInfo[i].MissingCount).
					Set("start_time = ?", accumMissInfo[i].StartTime).
					Set("end_time = ?", blockInfo[0].Time).
					Where("end_height = ? AND address = ?", accumMissInfo[i].EndHeight-int64(1), accumMissInfo[i].Address).
					Update()
				if err != nil {
					return err
				}
			}
		}

		if len(voteInfo) > 0 {
			var tempVoteInfo schema.VoteInfo
			for i := 0; i < len(voteInfo); i++ {
				// Check if a validator already voted
				count, _ := tx.Model(&tempVoteInfo).
					Where("proposal_id = ? AND voter = ?", voteInfo[i].ProposalID, voteInfo[i].Voter).
					Count()
				if count > 0 {
					_, err := tx.Model(&tempVoteInfo).
						Set("height = ?", voteInfo[i].Height).
						Set("option = ?", voteInfo[i].Option).
						Set("tx_hash = ?", voteInfo[i].TxHash).
						Set("gas_wanted = ?", voteInfo[i].GasWanted).
						Set("gas_used = ?", voteInfo[i].GasUsed).
						Set("time = ?", voteInfo[i].Time).
						Where("proposal_id = ? AND voter = ?", voteInfo[i].ProposalID, voteInfo[i].Voter).
						Update()
					if err != nil {
						return err
					}
				} else {
					err := tx.Insert(&voteInfo)
					if err != nil {
						return err
					}
				}
			}
		}

		if len(proposalInfo) > 0 {
			var tempProposalInfo schema.ProposalInfo
			for i := 0; i < len(proposalInfo); i++ {
				// check if a validator already voted
				count, _ := tx.Model(&tempProposalInfo).
					Where("id = ?", proposalInfo[i].ID).
					Count()

				if count > 0 {
					// save and update proposalInfo
					_, err := tx.Model(&tempProposalInfo).
						Set("tx_hash = ?", proposalInfo[i].TxHash).
						Set("proposer = ?", proposalInfo[i].Proposer).
						Set("initial_deposit_amount = ?", proposalInfo[i].InitialDepositAmount).
						Set("initial_deposit_denom = ?", proposalInfo[i].InitialDepositDenom).
						Where("id = ?", proposalInfo[i].ID).
						Update()
					if err != nil {
						return err
					}
				} else {
					err := tx.Insert(&proposalInfo)
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
