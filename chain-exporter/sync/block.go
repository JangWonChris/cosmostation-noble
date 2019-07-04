package sync

// import (
// 	"context"
// 	"crypto/tls"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/app/config"
// 	ctypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/app/types"
// 	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/app/utils"

// 	"github.com/cosmos/cosmos-sdk/codec"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/go-pg/pg"
// 	"github.com/olivere/elastic"
// 	"github.com/tendermint/tendermint/crypto"
// 	"github.com/tendermint/tendermint/rpc/client"
// 	resty "gopkg.in/resty.v1"
// )

// // Sync syncs the blockchain and missed blocks from a node
// func Sync(ctx context.Context, client *client.HTTP, db *pg.DB, cdc *codec.Codec, config *config.Config, elasticAPI *elastic.Client) error {
// 	// Check current height in db
// 	var blocks []ctypes.BlockInfo
// 	err := db.Model(&blocks).Order("height DESC").Limit(1).Select()
// 	if err != nil {
// 		return err
// 	}
// 	currentHeight := int64(1)
// 	if len(blocks) > 0 {
// 		currentHeight = blocks[0].Height
// 	}

// 	// Query the node for its height
// 	status, err := client.Status()
// 	if err != nil {
// 		return err
// 	}
// 	maxHeight := status.SyncInfo.LatestBlockHeight

// 	if currentHeight == 1 {
// 		currentHeight = 0
// 	}

// 	// Ingest all blocks up to the best height
// 	for i := currentHeight + 1; i <= maxHeight; i++ {
// 		err = IngestPrevBlock(ctx, i, client, db, cdc, config, elasticAPI)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Printf("synced block %d/%d \n", i, maxHeight)
// 	}
// 	return nil
// }

// // IngestPrevBlock queries the block at the given height-1 from the node and ingests its metadata (blockinfo,evidence)
// // into the database. It also queries the next block to access the commits and stores the missed signatures.
// func IngestPrevBlock(ctx context.Context, height int64, client *client.HTTP, db *pg.DB, cdc *codec.Codec, config *config.Config, elasticAPI *elastic.Client) error {

// 	// This is for syncing at the right height
// 	nextHeight := height + 1

// 	// Query the current block
// 	block, err := client.Block(&height)
// 	if err != nil {
// 		return err
// 	}

// 	// Query the next block to access the commits
// 	nextBlock, err := client.Block(&nextHeight)
// 	if err != nil {
// 		return err
// 	}

// 	// Query the block result for the block
// 	// blockResult, err := client.BlockResults(&height)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// Get validator set for the block
// 	validators, err := client.Validators(&height)
// 	if err != nil {
// 		return err
// 	}

// 	// Parse blockinfo & height needs to be previous height for the first block
// 	blockInfo := new(ctypes.BlockInfo)
// 	blockInfo.BlockHash = block.BlockMeta.BlockID.Hash.String()
// 	blockInfo.Proposer = block.Block.ProposerAddress.String()
// 	blockInfo.Height = block.Block.Height
// 	blockInfo.TotalTxs = block.Block.TotalTxs
// 	blockInfo.NumTxs = block.Block.NumTxs
// 	blockInfo.Time = block.BlockMeta.Header.Time

// 	// Validator information & Missing block, detail, consecutive information
// 	validatorSetInfo := make([]*ctypes.ValidatorSetInfo, 0)
// 	misseInfo := make([]*ctypes.MissInfo, 0)
// 	misseDetailInfo := make([]*ctypes.MissDetailInfo, 0)
// 	accumMissInfo := make([]*ctypes.MissInfo, 0)
// 	for i, validator := range validators.Validators {
// 		// At height 1, insert validators as an event_type of create_validator
// 		if validators.BlockHeight == 1 {
// 			tempValidatorSetInfo := &ctypes.ValidatorSetInfo{
// 				IDValidator:          i + 1,
// 				Height:               blockInfo.Height,
// 				Proposer:             validator.Address.String(),
// 				VotingPower:          float64(validator.VotingPower),
// 				NewVotingPowerAmount: float64(validator.VotingPower),
// 				NewVotingPowerDenom:  "uatom",
// 				EventType:            "create_validator",
// 				Time:                 block.BlockMeta.Header.Time,
// 			}
// 			validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)
// 		}

// 		// Missing information
// 		if nextBlock.Block.LastCommit.Precommits[i] == nil {
// 			// Missing Detail information (save every single height)
// 			tempMissDetailInfo := &ctypes.MissDetailInfo{
// 				Height:   block.BlockMeta.Header.Height,
// 				Address:  validator.Address.String(),
// 				Proposer: block.BlockMeta.Header.ProposerAddress.String(),
// 				Alerted:  false,
// 				Time:     block.BlockMeta.Header.Time,
// 			}
// 			misseDetailInfo = append(misseDetailInfo, tempMissDetailInfo)

// 			// Initial variables
// 			startHeight := block.BlockMeta.Header.Height
// 			endHeight := block.BlockMeta.Header.Height
// 			missingCount := int64(1)

// 			// Query to check if a validator missed previous block
// 			var missInfo ctypes.MissInfo
// 			_ = db.Model(&missInfo).
// 				Where("end_height = ? AND address = ?", endHeight-int64(1), validator.Address.String()).
// 				Order("end_height DESC").
// 				Select()

// 			if missInfo.Address == "" {
// 				tempMissInfo := &ctypes.MissInfo{
// 					Address:      validator.Address.String(),
// 					StartHeight:  startHeight,
// 					EndHeight:    endHeight,
// 					MissingCount: missingCount,
// 					StartTime:    block.BlockMeta.Header.Time,
// 					EndTime:      block.BlockMeta.Header.Time,
// 					Alerted:      false,
// 				}
// 				misseInfo = append(misseInfo, tempMissInfo)
// 			} else {
// 				tempMissInfo := &ctypes.MissInfo{
// 					Address:      missInfo.Address,
// 					StartHeight:  missInfo.StartHeight,
// 					EndHeight:    missInfo.EndHeight + int64(1),
// 					MissingCount: missInfo.MissingCount + int64(1),
// 					StartTime:    missInfo.StartTime,
// 					EndTime:      block.BlockMeta.Header.Time,
// 					Alerted:      false,
// 				}
// 				accumMissInfo = append(accumMissInfo, tempMissInfo)
// 			}
// 			continue
// 		}
// 	}

// 	// Evidence information
// 	evidenceInfo := make([]*ctypes.EvidenceInfo, 0)
// 	for _, evidence := range nextBlock.Block.Evidence.Evidence {
// 		tempEvidenceInfo := &ctypes.EvidenceInfo{
// 			Address: strings.ToUpper(string(hex.EncodeToString(evidence.Address()))),
// 			Height:  evidence.Height(),
// 			Hash:    strings.ToUpper(string(hex.EncodeToString(evidence.Hash()))),
// 			Time:    block.BlockMeta.Header.Time,
// 		}
// 		evidenceInfo = append(evidenceInfo, tempEvidenceInfo)
// 	}

// 	// Transaction information & validatorSetInfo for voting power history
// 	transactionInfo := make([]*ctypes.TransactionInfo, 0)
// 	voteInfo := make([]*ctypes.VoteInfo, 0)
// 	depositInfo := make([]*ctypes.DepositInfo, 0)
// 	proposalInfo := make([]*ctypes.ProposalInfo, 0)
// 	for _, tx := range block.Block.Data.Txs {
// 		// Use tx codec to unmarshal binary length prefix
// 		var sdkTx sdk.Tx
// 		err := cdc.UnmarshalBinaryLengthPrefixed([]byte(tx), &sdkTx)
// 		if err != nil {
// 			return err
// 		}

// 		// Transaction hash
// 		txByte := crypto.Sha256(tx)
// 		txHash := hex.EncodeToString(txByte)

// 		// Insert data for PostgreSQL
// 		tempTransactionInfo := &ctypes.TransactionInfo{
// 			Height: block.Block.Height,
// 			TxHash: txHash,
// 			Time:   block.BlockMeta.Header.Time,
// 		}
// 		transactionInfo = append(transactionInfo, tempTransactionInfo)

// 		// Query LCD
// 		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
// 		resp, err := resty.R().Get(config.Node.LCDURL + "/txs/" + txHash)
// 		if err != nil {
// 			fmt.Printf("Transaction LCD resty - %v\n", err)
// 		}

// 		// Unmarshal general transaction format
// 		var generalTx ctypes.GeneralTx
// 		_ = json.Unmarshal(resp.Body(), &generalTx)

// 		// Check log to see if tx is success
// 		for j, log := range generalTx.Logs {
// 			// Check log to see if tx is success
// 			if log.Success {
// 				switch generalTx.Tx.Value.Msg[j].Type {
// 				case "cosmos-sdk/MsgCreateValidator":
// 					var createValidatorTx ctypes.CreateValidatorMsgValueTx
// 					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &createValidatorTx)

// 					// 이렇게 넣는것도 문제가 발생!
// 					// 동일한 블록안에 create_validator 가 있을 경우 id_validator 를 체크하기가 힘들다
// 					// Query the highest height of id_validator
// 					var lastValidatorSetInfo ctypes.ValidatorSetInfo
// 					_ = db.Model(&lastValidatorSetInfo).
// 						Column("id_validator").
// 						Order("id_validator DESC").
// 						Limit(1).
// 						Select()

// 					// Conversion
// 					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
// 					newVotingPowerAmount, _ := strconv.ParseFloat(createValidatorTx.Value.Amount.String(), 64) // parseFloat from sdk.Dec.String()
// 					newVotingPowerAmount = float64(newVotingPowerAmount) / 1000000

// 					// Insert data
// 					tempValidatorSetInfo := &ctypes.ValidatorSetInfo{
// 						IDValidator:          lastValidatorSetInfo.IDValidator + 1,
// 						Height:               height,
// 						Proposer:             utils.ConsensusPubkeyToProposer(createValidatorTx.Pubkey),
// 						VotingPower:          newVotingPowerAmount,
// 						NewVotingPowerAmount: newVotingPowerAmount,
// 						NewVotingPowerDenom:  createValidatorTx.Value.Denom,
// 						EventType:            "create_validator",
// 						TxHash:               generalTx.TxHash,
// 						Time:                 block.BlockMeta.Header.Time,
// 					}
// 					validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

// 				case "cosmos-sdk/MsgDelegate":
// 					var delegateTx ctypes.DelegateMsgValueTx
// 					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &delegateTx)

// 					// Transaction Messeage
// 					var tempValidatorInfo ctypes.ValidatorInfo
// 					_ = db.Model(&tempValidatorInfo).
// 						Column("proposer").
// 						Where("operator_address = ?", delegateTx.ValidatorAddress).
// 						Limit(1).
// 						Select()

// 					// Query last id_validator
// 					var lastValidatorSetInfo ctypes.ValidatorSetInfo
// 					_ = db.Model(&lastValidatorSetInfo).
// 						Column("id_validator").
// 						Where("proposer = ?", tempValidatorInfo.Proposer).
// 						Order("id DESC").
// 						Limit(1).
// 						Select()

// 					// Conversion
// 					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
// 					newVotingPowerAmount, _ := strconv.ParseFloat(delegateTx.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
// 					newVotingPowerAmount = newVotingPowerAmount / 1000000

// 					// Current Voting Power
// 					var votingPower float64
// 					validators, _ := client.Validators(&height)
// 					for _, validator := range validators.Validators {
// 						if validator.Address.String() == tempValidatorInfo.Proposer {
// 							votingPower = float64(validator.VotingPower)
// 						}
// 					}

// 					// 동일한 블록에서 서로 다른 주소에서 동일한 검증인에게 위임한 트랜잭션이 있을 경우 현재 VotingPower는 같다. 기술적 한계 (Certus One 17번째 블록에 두번)
// 					// Insert data
// 					tempValidatorSetInfo := &ctypes.ValidatorSetInfo{
// 						IDValidator:          lastValidatorSetInfo.IDValidator,
// 						Height:               height,
// 						Proposer:             tempValidatorInfo.Proposer,
// 						VotingPower:          votingPower + newVotingPowerAmount,
// 						EventType:            "delegate",
// 						NewVotingPowerAmount: newVotingPowerAmount,
// 						NewVotingPowerDenom:  delegateTx.Amount.Denom,
// 						TxHash:               generalTx.TxHash,
// 						Time:                 block.BlockMeta.Header.Time,
// 					}
// 					validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

// 				case "cosmos-sdk/MsgUndelegate":
// 					var undelegateTx ctypes.UndelegateMsgValueTx
// 					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &undelegateTx)

// 					// Transaction Messeage
// 					var tempValidatorInfo ctypes.ValidatorInfo
// 					_ = db.Model(&tempValidatorInfo).
// 						Column("proposer").
// 						Where("operator_address = ?", undelegateTx.ValidatorAddress).
// 						Limit(1).
// 						Select()

// 					// Query last id_validator
// 					var lastValidatorSetInfo ctypes.ValidatorSetInfo
// 					_ = db.Model(&lastValidatorSetInfo).
// 						Column("id_validator", "voting_power").
// 						Where("proposer = ?", tempValidatorInfo.Proposer).
// 						Order("id DESC").
// 						Limit(1).
// 						Select()

// 					// Conversion
// 					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
// 					newVotingPowerAmount, _ := strconv.ParseFloat(undelegateTx.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
// 					newVotingPowerAmount = -newVotingPowerAmount / 1000000

// 					// Current Voting Power
// 					var votingPower float64
// 					validators, _ := client.Validators(&height)
// 					for _, validator := range validators.Validators {
// 						if validator.Address.String() == tempValidatorInfo.Proposer {
// 							votingPower = float64(validator.VotingPower)
// 						}
// 					}

// 					// Insert data
// 					tempValidatorSetInfo := &ctypes.ValidatorSetInfo{
// 						IDValidator:          lastValidatorSetInfo.IDValidator,
// 						Height:               height,
// 						Proposer:             tempValidatorInfo.Proposer,
// 						VotingPower:          votingPower + newVotingPowerAmount,
// 						EventType:            "begin_unbonding",
// 						NewVotingPowerAmount: newVotingPowerAmount,
// 						NewVotingPowerDenom:  undelegateTx.Amount.Denom,
// 						TxHash:               generalTx.TxHash,
// 						Time:                 block.BlockMeta.Header.Time,
// 					}
// 					validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

// 				case "cosmos-sdk/MsgBeginRedelegate":
// 					var redelegateTx ctypes.RedelegateMsgValueTx
// 					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &redelegateTx)

// 					/*
// 						Redelegate 당한 검증인은 -
// 						Redelegate 한 검증인은 +
// 					*/

// 					// Query validator_dst_address's proposer address
// 					var tempDstValidatorInfo ctypes.ValidatorInfo
// 					_ = db.Model(&tempDstValidatorInfo).
// 						Column("proposer").
// 						Where("operator_address = ?", redelegateTx.ValidatorDstAddress).
// 						Limit(1).
// 						Select()

// 					// Query validator_src_address's proposer address
// 					var tempSrcValidatorInfo ctypes.ValidatorInfo
// 					_ = db.Model(&tempSrcValidatorInfo).
// 						Column("proposer").
// 						Where("operator_address = ?", redelegateTx.ValidatorSrcAddress).
// 						Limit(1).
// 						Select()

// 					// Query last id_validator
// 					var lastDstValidatorSetInfo ctypes.ValidatorSetInfo
// 					_ = db.Model(&lastDstValidatorSetInfo).
// 						Column("id_validator", "voting_power").
// 						Where("proposer = ?", tempDstValidatorInfo.Proposer).
// 						Order("id DESC").
// 						Limit(1).
// 						Select()

// 					// Query last id_validator
// 					var lastSrcValidatorSetInfo ctypes.ValidatorSetInfo
// 					_ = db.Model(&lastSrcValidatorSetInfo).
// 						Column("id_validator", "voting_power").
// 						Where("proposer = ?", tempSrcValidatorInfo.Proposer).
// 						Order("id DESC").
// 						Limit(1).
// 						Select()

// 					// Conversion
// 					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
// 					newVotingPowerAmount, _ := strconv.ParseFloat(redelegateTx.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
// 					newVotingPowerAmount = newVotingPowerAmount / 1000000

// 					// Current Destination Validator's Voting Power
// 					var dstValidatorVotingPower float64
// 					validators, _ := client.Validators(&height)
// 					for _, validator := range validators.Validators {
// 						if validator.Address.String() == tempDstValidatorInfo.Proposer {
// 							dstValidatorVotingPower = float64(validator.VotingPower)
// 						}
// 					}

// 					// Insert ValidatorDstAddress data
// 					tempDstValidatorSetInfo := &ctypes.ValidatorSetInfo{
// 						IDValidator:          lastDstValidatorSetInfo.IDValidator,
// 						Height:               height,
// 						Proposer:             tempDstValidatorInfo.Proposer,
// 						VotingPower:          dstValidatorVotingPower + newVotingPowerAmount,
// 						EventType:            "begin_redelegate",
// 						NewVotingPowerAmount: newVotingPowerAmount,
// 						NewVotingPowerDenom:  redelegateTx.Amount.Denom,
// 						TxHash:               generalTx.TxHash,
// 						Time:                 block.BlockMeta.Header.Time,
// 					}
// 					validatorSetInfo = append(validatorSetInfo, tempDstValidatorSetInfo)

// 					// Current Source Validator's Voting Power
// 					var srcValidatorVotingPower float64
// 					validators, _ = client.Validators(&height)
// 					for _, validator := range validators.Validators {
// 						if validator.Address.String() == tempSrcValidatorInfo.Proposer {
// 							srcValidatorVotingPower = float64(validator.VotingPower)
// 						}
// 					}

// 					// Insert ValidatorSrcAddress data
// 					tempSrcValidatorSetInfo := &ctypes.ValidatorSetInfo{
// 						IDValidator:          lastSrcValidatorSetInfo.IDValidator,
// 						Height:               height,
// 						Proposer:             tempSrcValidatorInfo.Proposer,
// 						VotingPower:          srcValidatorVotingPower - newVotingPowerAmount,
// 						EventType:            "begin_redelegate",
// 						NewVotingPowerAmount: -newVotingPowerAmount,
// 						NewVotingPowerDenom:  redelegateTx.Amount.Denom,
// 						TxHash:               generalTx.TxHash,
// 						Time:                 block.BlockMeta.Header.Time,
// 					}
// 					validatorSetInfo = append(validatorSetInfo, tempSrcValidatorSetInfo)

// 				case "cosmos-sdk/MsgSubmitProposal":
// 					var submitTx ctypes.SubmitProposalMsgValueTx
// 					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &submitTx)

// 					// 멀티메시지 일 경우 Tags 번호가 달라져서 아래와 같이 key 값을 찾고 value를 넣어줘야 된다
// 					// 141050 블록높이: 7de25c478cf26eb6843c6a1b7a1cb550c8ab77ba9563a252677c059572bea6c3
// 					var proposalID int64
// 					for _, tag := range generalTx.Tags {
// 						if tag.Key == "proposal-id" {
// 							proposalID, _ = strconv.ParseInt(tag.Value, 10, 64)
// 						}
// 					}

// 					initialDepositAmount, _ := strconv.ParseInt(submitTx.InitialDeposit[0].Amount, 10, 64)
// 					initialDepositDenom := submitTx.InitialDeposit[0].Denom

// 					// Insert data
// 					tempProposalInfo := &ctypes.ProposalInfo{
// 						ID:                   proposalID,
// 						TxHash:               generalTx.TxHash,
// 						Proposer:             submitTx.Proposer,
// 						InitialDepositAmount: string(initialDepositAmount),
// 						InitialDepositDenom:  initialDepositDenom,
// 					}
// 					proposalInfo = append(proposalInfo, tempProposalInfo)

// 					// Conversion
// 					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
// 					gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
// 					gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)

// 					// Insert data
// 					tempDepositInfo := &ctypes.DepositInfo{
// 						Height:     height,
// 						ProposalID: proposalID,
// 						Depositor:  submitTx.Proposer,
// 						Amount:     initialDepositAmount,
// 						Denom:      initialDepositDenom,
// 						TxHash:     generalTx.TxHash,
// 						GasWanted:  gasWanted,
// 						GasUsed:    gasUsed,
// 						Time:       block.BlockMeta.Header.Time,
// 					}
// 					depositInfo = append(depositInfo, tempDepositInfo)

// 				case "cosmos-sdk/MsgVote":
// 					var voteTx ctypes.VoteMsgValueTx
// 					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &voteTx)

// 					// Transaction Messeage
// 					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
// 					proposalID, _ := strconv.ParseInt(voteTx.ProposalID, 10, 64)
// 					gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
// 					gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)
// 					tempVoteInfo := &ctypes.VoteInfo{
// 						Height:     height,
// 						ProposalID: proposalID,
// 						Voter:      voteTx.Voter,
// 						Option:     voteTx.Option,
// 						TxHash:     generalTx.TxHash,
// 						GasWanted:  gasWanted,
// 						GasUsed:    gasUsed,
// 						Time:       block.BlockMeta.Header.Time,
// 					}
// 					voteInfo = append(voteInfo, tempVoteInfo)

// 				case "cosmos-sdk/MsgDeposit":
// 					var depositTx ctypes.DepositMsgValueTx
// 					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &depositTx)

// 					// Transaction Messeage
// 					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
// 					proposalID, _ := strconv.ParseInt(depositTx.ProposalID, 10, 64)
// 					amount, _ := strconv.ParseInt(depositTx.Amount[0].Amount, 10, 64)
// 					gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
// 					gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)
// 					tempDepositInfo := &ctypes.DepositInfo{
// 						Height:     height,
// 						ProposalID: proposalID,
// 						Depositor:  depositTx.Depositor,
// 						Amount:     amount,
// 						Denom:      depositTx.Amount[j].Denom,
// 						TxHash:     generalTx.TxHash,
// 						GasWanted:  gasWanted,
// 						GasUsed:    gasUsed,
// 						Time:       block.BlockMeta.Header.Time,
// 					}
// 					depositInfo = append(depositInfo, tempDepositInfo)
// 				default:
// 					continue
// 				}
// 			}
// 		}

// 	}

// 	// Insert in DB - err to rollback tx
// 	err = db.RunInTransaction(func(tx *pg.Tx) error {
// 		// Insert block info
// 		err = tx.Insert(blockInfo)
// 		if err != nil {
// 			return err
// 		}

// 		// Insert evidence info
// 		if len(evidenceInfo) > 0 {
// 			err = tx.Insert(&evidenceInfo)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// Insert missing block info
// 		if len(misseInfo) > 0 {
// 			err = tx.Insert(&misseInfo)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// Insert missing block info (every height)
// 		if len(misseDetailInfo) > 0 {
// 			err = tx.Insert(&misseDetailInfo)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// Insert transaction info
// 		if len(transactionInfo) > 0 {
// 			err = tx.Insert(&transactionInfo)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// Insert validator info
// 		if len(validatorSetInfo) > 0 {
// 			err = tx.Insert(&validatorSetInfo)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// Insert vote tx info
// 		if len(voteInfo) > 0 {
// 			var tempVoteInfo ctypes.VoteInfo
// 			for i := 0; i < len(voteInfo); i++ {
// 				// Check if a validator already voted
// 				count, _ := tx.Model(&tempVoteInfo).
// 					Where("proposal_id = ? AND voter = ?", voteInfo[i].ProposalID, voteInfo[i].Voter).
// 					Count()
// 				if count > 0 {
// 					_, err = tx.Model(&tempVoteInfo).
// 						Set("height = ?", voteInfo[i].Height).
// 						Set("option = ?", voteInfo[i].Option).
// 						Set("tx_hash = ?", voteInfo[i].TxHash).
// 						Set("gas_wanted = ?", voteInfo[i].GasWanted).
// 						Set("gas_used = ?", voteInfo[i].GasUsed).
// 						Set("time = ?", voteInfo[i].Time).
// 						Where("proposal_id = ? AND voter = ?", voteInfo[i].ProposalID, voteInfo[i].Voter).
// 						Update()
// 					if err != nil {
// 						return err
// 					}
// 				} else {
// 					err = tx.Insert(&voteInfo)
// 					if err != nil {
// 						return err
// 					}
// 				}
// 			}
// 		}

// 		// Insert deposit tx info
// 		if len(depositInfo) > 0 {
// 			err = tx.Insert(&depositInfo)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// Exist and update proposerInfo
// 		if len(proposalInfo) > 0 {
// 			var tempProposalInfo ctypes.ProposalInfo
// 			for i := 0; i < len(proposalInfo); i++ {
// 				// Check if a validator already voted
// 				count, _ := tx.Model(&tempProposalInfo).
// 					Where("id = ?", proposalInfo[i].ID).
// 					Count()

// 				if count > 0 {
// 					// Save and update proposalInfo
// 					_, err = tx.Model(&tempProposalInfo).
// 						Set("tx_hash = ?", proposalInfo[i].TxHash).
// 						Set("proposer = ?", proposalInfo[i].Proposer).
// 						Set("initial_deposit_amount = ?", proposalInfo[i].InitialDepositAmount).
// 						Set("initial_deposit_denom = ?", proposalInfo[i].InitialDepositDenom).
// 						Where("id = ?", proposalInfo[i].ID).
// 						Update()
// 					if err != nil {
// 						return err
// 					}
// 				} else {
// 					err = tx.Insert(&proposalInfo)
// 					if err != nil {
// 						return err
// 					}
// 				}
// 			}
// 		}

// 		// Update accumulative missing block info
// 		var tempMissInfo ctypes.MissInfo
// 		if len(accumMissInfo) > 0 {
// 			for i := 0; i < len(accumMissInfo); i++ {
// 				_, err = tx.Model(&tempMissInfo).
// 					Set("address = ?", accumMissInfo[i].Address).
// 					Set("start_height = ?", accumMissInfo[i].StartHeight).
// 					Set("end_height = ?", accumMissInfo[i].EndHeight).
// 					Set("missing_count = ?", accumMissInfo[i].MissingCount).
// 					Set("start_time = ?", accumMissInfo[i].StartTime).
// 					Set("end_time = ?", block.BlockMeta.Header.Time).
// 					Where("end_height = ? AND address = ?", accumMissInfo[i].EndHeight-int64(1), accumMissInfo[i].Address).
// 					Update()
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
