package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/stats"
	dbtypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
	u "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	resty "gopkg.in/resty.v1"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"

	// resty "gopkg.in/resty.v1"

	"github.com/cosmos/cosmos-sdk/codec"
)

// GetValidators returns all existing validators
func GetValidators(RPCClient *client.HTTP, DB *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Query all validators order by their tokens
	validatorInfo := make([]*dbtypes.ValidatorInfo, 0)

	// Check status param
	statusParam := r.FormValue("status")
	switch statusParam {
	case "":
		_ = DB.Model(&validatorInfo).
			Order("id ASC").
			Select()
	case "active":
		_ = DB.Model(&validatorInfo).
			Where("status = ?", 2).
			Order("id ASC").
			Select()
	case "inactive":
		_ = DB.Model(&validatorInfo).
			Where("status = ? OR status = ?", 0, 1).
			Order("id ASC").
			Select()
	default:
		return json.NewEncoder(w).Encode(validatorInfo)
	}

	if len(validatorInfo) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query status for the current height
	status, _ := RPCClient.Status()
	currentHeight := status.SyncInfo.LatestBlockHeight

	// Sort validatorInfo by highest tokens
	sort.Slice(validatorInfo[:], func(i, j int) bool {
		tempToken1, _ := strconv.Atoi(validatorInfo[i].Tokens)
		tempToken2, _ := strconv.Atoi(validatorInfo[j].Tokens)
		return tempToken1 > tempToken2
	})

	resultValidator := make([]*models.ResultValidator, 0)
	for _, validator := range validatorInfo {
		// Convert to proposer address
		proposerAddress := u.ConsensusPubkeyToProposer(validator.ConsensusPubkey)

		// Query how many missed blocks each validator has for the last 100 blocks
		var missDetailInfo []dbtypes.MissDetailInfo
		missBlockCount, _ := DB.Model(&missDetailInfo).
			Where("address = ? AND height BETWEEN ? AND ?", proposerAddress, currentHeight-100, currentHeight).
			Count()

		// If a validator is jailed, missing block is 100
		if validator.Jailed {
			missBlockCount = 100
		}

		// Insert uptime data
		tempUptime := &models.Uptime{
			Address:      proposerAddress,
			MissedBlocks: int64(missBlockCount),
			OverBlocks:   100,
		}

		// Insert validator data
		tempResultValidator := &models.ResultValidator{
			Rank:                 validator.Rank,
			OperatorAddress:      validator.OperatorAddress,
			ConsensusPubkey:      validator.ConsensusPubkey,
			Jailed:               validator.Jailed,
			Status:               validator.Status,
			Tokens:               validator.Tokens,
			DelegatorShares:      validator.DelegatorShares,
			Moniker:              validator.Moniker,
			Identity:             validator.Identity,
			Website:              validator.Website,
			Details:              validator.Details,
			UnbondingHeight:      validator.UnbondingHeight,
			UnbondingTime:        validator.UnbondingTime,
			CommissionRate:       validator.CommissionRate,
			CommissionMaxRate:    validator.CommissionMaxRate,
			CommissionChangeRate: validator.CommissionChangeRate,
			UpdateTime:           validator.UpdateTime,
			Uptime:               *tempUptime,
			MinSelfDelegation:    validator.MinSelfDelegation,
			KeybaseURL:           validator.KeybaseURL,
		}
		resultValidator = append(resultValidator, tempResultValidator)
	}

	u.Respond(w, resultValidator)
	return nil
}

// GetValidator receives validator address and returns that validator
func GetValidator(RPCClient *client.HTTP, DB *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Receive address
	vars := mux.Vars(r)
	address := vars["address"]

	// Convert to proposer address format
	validatorInfo, _ := u.ConvertToProposer(address, DB)

	// Check if the input validator address exists
	if validatorInfo.Proposer == "" {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Proposer Address
	address = validatorInfo.Proposer

	// Query status for the current height
	status, _ := RPCClient.Status()
	currentHeight := status.SyncInfo.LatestBlockHeight

	// Query how many missed blocks each validator has for the last 100 blocks
	proposerAddress := u.ConsensusPubkeyToProposer(validatorInfo.ConsensusPubkey)
	var missDetailInfo []dbtypes.MissDetailInfo
	missBlockCount, _ := DB.Model(&missDetailInfo).
		Where("address = ? AND height BETWEEN ? AND ?", proposerAddress, currentHeight-100, currentHeight).
		Count()

	// If a validator is jailed, missing block is 100
	if validatorInfo.Jailed {
		missBlockCount = 100
	}

	tempUptime := &models.Uptime{
		Address:      proposerAddress,
		MissedBlocks: int64(missBlockCount),
		OverBlocks:   100,
	}

	// Validator's bonded height and its timestamp
	var validatorSetInfo dbtypes.ValidatorSetInfo
	_ = DB.Model(&validatorSetInfo).
		Where("proposer = ? AND event_type = ?", proposerAddress, "create_validator").
		Select()

	resultValidatorDetail := &models.ResultValidatorDetail{
		Rank:                 validatorInfo.Rank,
		OperatorAddress:      validatorInfo.OperatorAddress,
		ConsensusPubkey:      validatorInfo.ConsensusPubkey,
		BondedHeight:         validatorSetInfo.Height,
		BondedTime:           validatorSetInfo.Time,
		Jailed:               validatorInfo.Jailed,
		Status:               validatorInfo.Status,
		Tokens:               validatorInfo.Tokens,
		DelegatorShares:      validatorInfo.DelegatorShares,
		Moniker:              validatorInfo.Moniker,
		Identity:             validatorInfo.Identity,
		Website:              validatorInfo.Website,
		Details:              validatorInfo.Details,
		UnbondingHeight:      validatorInfo.UnbondingHeight,
		UnbondingTime:        validatorInfo.UnbondingTime,
		CommissionRate:       validatorInfo.CommissionRate,
		CommissionMaxRate:    validatorInfo.CommissionMaxRate,
		CommissionChangeRate: validatorInfo.CommissionChangeRate,
		UpdateTime:           validatorInfo.UpdateTime,
		Uptime:               *tempUptime,
		MinSelfDelegation:    validatorInfo.MinSelfDelegation,
		KeybaseURL:           validatorInfo.KeybaseURL,
	}

	u.Respond(w, resultValidatorDetail)
	return nil
}

// GetValidatorBlockMisses receives validator address and returns the validator's block consencutive misses
func GetValidatorBlockMisses(RPCClient *client.HTTP, DB *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Receive address
	vars := mux.Vars(r)
	address := vars["address"]

	// Change to proposer address format
	validatorInfo, _ := u.ConvertToProposerSlice(address, DB)

	// Check if the validator address exists
	if len(validatorInfo) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Validator's proposer address
	address = validatorInfo[0].Proposer

	// Query a validator's missing blocks
	var missInfos []dbtypes.MissInfo
	err := DB.Model(&missInfos).
		Where("address = ?", address).
		Limit(50).
		Order("start_height DESC").
		Select()
	if err != nil {
		return err
	}

	resultMisses := make([]*models.ResultMisses, 0)
	for _, missInfo := range missInfos {
		tempResultMisses := &models.ResultMisses{
			StartHeight:  missInfo.StartHeight,
			EndHeight:    missInfo.EndHeight,
			MissingCount: missInfo.MissingCount,
			StartTime:    missInfo.StartTime,
			EndTime:      missInfo.EndTime,
		}
		resultMisses = append(resultMisses, tempResultMisses)
	}

	u.Respond(w, resultMisses)	
	return nil
}

// GetValidatorBlockMissesDetail receives validator address and returns the validator's block misses (uptime)
func GetValidatorBlockMissesDetail(RPCClient *client.HTTP, DB *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Receive address
	vars := mux.Vars(r)
	address := vars["address"]

	// Change to proposer address format
	validatorInfo, _ := u.ConvertToProposerSlice(address, DB)

	// Check if the input validator address exists
	if len(validatorInfo) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Validator's proposer address
	address = validatorInfo[0].Proposer

	// Query the current block
	// * Second highest block number saved in a database due to client's handling
	var blockInfo []dbtypes.BlockInfo
	_ = DB.Model(&blockInfo).
		Column("height").
		Order("id DESC").
		Limit(2).
		Select()

	// Query a validator's missing blocks
	var missDetailInfos []dbtypes.MissDetailInfo
	_ = DB.Model(&missDetailInfos).
		Where("address = ? AND height >= ?", address, blockInfo[1].Height-100).
		Limit(100).
		Order("height DESC").
		Select()

	resultMissesDetail := make([]*models.ResultMissesDetail, 0)
	if len(missDetailInfos) <= 0 {
		return json.NewEncoder(w).Encode(resultMissesDetail)
	}

	for _, missDetailInfo := range missDetailInfos {
		tempResultMissesDetail := &models.ResultMissesDetail{
			Height: missDetailInfo.Height,
			Time:   missDetailInfo.Time,
		}
		resultMissesDetail = append(resultMissesDetail, tempResultMissesDetail)
	}

	u.Respond(w, resultMissesDetail)	
	return nil
}

// GetValidatorEvents receives validator address and returns the validator's events
func GetValidatorEvents(DB *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Receive Address
	vars := mux.Vars(r)
	address := vars["address"]

	// Define default variables
	limit := int(50)
	from := int(1)

	// Check limit param
	tempLimit := r.URL.Query()["limit"]
	if len(tempLimit) > 0 {
		tempLimit, _ := strconv.Atoi(tempLimit[0])
		limit = tempLimit
	}

	// Max Limit
	if limit > 50 {
		errors.ErrOverMaxLimit(w, http.StatusRequestedRangeNotSatisfiable)
		return nil
	}

	// Check from param
	tempFrom := r.URL.Query()["from"]
	if len(tempFrom) > 0 {
		tempFrom, _ := strconv.Atoi(tempFrom[0])
		from = tempFrom
	} else {
		// Get current height in DB
		var blocks dbtypes.BlockInfo
		_ = DB.Model(&blocks).
			Order("id DESC").
			Limit(1).
			Select()
		from = int(blocks.Height)
	}

	// Change to proposer address format
	validatorInfo, _ := u.ConvertToProposerSlice(address, DB)

	// Check if the input validator address exists
	if len(validatorInfo) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Address
	address = validatorInfo[0].Proposer

	// Query id_validator
	var idValidatorSetInfo dbtypes.ValidatorSetInfo
	_ = DB.Model(&idValidatorSetInfo).
		Column("id_validator").
		Where("proposer = ?", address).
		Limit(1).
		Select()

	resultVotingPowerHistory := make([]*models.ResultVotingPowerHistory, 0)
	if idValidatorSetInfo.IDValidator != 0 {
		var validatorSetInfo []dbtypes.ValidatorSetInfo
		_ = DB.Model(&validatorSetInfo).
			Where("id_validator = ? AND height <= ?", idValidatorSetInfo.IDValidator, from).
			Limit(limit).
			Order("id DESC").
			Select()

		for _, validatorSet := range validatorSetInfo {
			tempResultValidatorSet := &models.ResultVotingPowerHistory{
				Height:         validatorSet.Height,
				EventType:      validatorSet.EventType,
				VotingPower:    validatorSet.VotingPower * 1000000,
				NewVotingPower: validatorSet.NewVotingPowerAmount * 1000000,
				TxHash:         validatorSet.TxHash,
				Timestamp:      validatorSet.Time,
			}
			resultVotingPowerHistory = append(resultVotingPowerHistory, tempResultValidatorSet)
		}
	}

	u.Respond(w, resultVotingPowerHistory)	
	return nil
}

// GetRedelegations receives delegator, srcvalidator, dstvalidator address and returns redelegations information
func GetRedelegations(DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	endpoint := "/staking/redelegations?"
	if len(r.URL.Query()["delegator"]) > 0 {
		endpoint += fmt.Sprintf("delegator=%s&", r.URL.Query()["delegator"][0])
	}

	if len(r.URL.Query()["validator_from"]) > 0 {
		endpoint += fmt.Sprintf("validator_from=%s&", r.URL.Query()["validator_from"][0])
	}

	if len(r.URL.Query()["validator_to"]) > 0 {
		endpoint += fmt.Sprintf("validator_to=%s&", r.URL.Query()["validator_to"][0])
	}

	// Query LCD
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, _ := resty.R().Get(Config.Node.LCDURL + endpoint)

	var redelegations []models.Redelegations
	err := json.Unmarshal(resp.Body(), &redelegations)
	if err != nil {
		fmt.Printf("staking/redelegations? unmarshal error - %v\n", err)
	}

	u.Respond(w, redelegations)	
	return nil
}

// GetValidatorDelegations receives validator address and returns all existing delegations that are delegated to the validator
func GetValidatorDelegations(Codec *codec.Codec, RPCClient *client.HTTP, DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Receive address
	vars := mux.Vars(r)
	operatorAddress := vars["address"]

	// Change to proposer address format
	validatorInfo, _ := u.ConvertToProposerSlice(operatorAddress, DB)

	// Check if the input validator address exists
	if len(validatorInfo) <= 0 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// OperatorAddress & proposer (to decrease database overload, use shorter string)
	operatorAddress = validatorInfo[0].OperatorAddress
	proposer := validatorInfo[0].Proposer

	// Query a validator's information
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // Local 환경에서 테스트를 위해
	validatorResp, _ := resty.R().Get(Config.Node.LCDURL + "/staking/validators/" + operatorAddress)

	var validator models.Validator
	err := json.Unmarshal(validatorResp.Body(), &validator)
	if err != nil {
		fmt.Printf("staking/validators/ unmarshal error - %v\n", err)
	}

	// Query delegations of a validator
	resp, _ := resty.R().Get(Config.Node.LCDURL + "/staking/validators/" + operatorAddress + "/delegations")

	var validatorDelegations []*models.ValidatorDelegations
	err = json.Unmarshal(resp.Body(), &validatorDelegations)
	if err != nil {
		fmt.Printf("staking/validators/{address}/delegations unmarshal error - %v\n", err)
	}

	// Validator's token divide by delegator_shares equals amount of uatom
	tokens, _ := strconv.ParseFloat(validator.Tokens.String(), 64)
	delegatorShares, _ := strconv.ParseFloat(validator.DelegatorShares.String(), 64)
	uatom := tokens / delegatorShares

	// Validator's all delegations
	var resultValidatorDelegations []*models.ValidatorDelegations
	for _, validatorDelegation := range validatorDelegations {
		shares, _ := strconv.ParseFloat(validatorDelegation.Shares.String(), 64)
		amount := fmt.Sprintf("%f", shares*uatom)

		tempValidatorDelegations := &models.ValidatorDelegations{
			DelegatorAddress: validatorDelegation.DelegatorAddress,
			ValidatorAddress: validatorDelegation.ValidatorAddress,
			Shares:           validatorDelegation.Shares,
			Amount:           amount,
		}
		resultValidatorDelegations = append(resultValidatorDelegations, tempValidatorDelegations)
	}

	// Query delegation change rate in 24 hours by 24 rows order by descending id
	latestDelegatorNum := make([]*stats.ValidatorStats, 0)
	_ = DB.Model(&latestDelegatorNum).
		Where("proposer_address = ?", proposer).
		Order("id DESC").
		Limit(24).
		Select()

	// 24 hrs
	count := len(latestDelegatorNum) - 1

	// Initial variables and current Delegator Num
	before24HBondedTokens := int(0)
	delegatorNumChange24H := int(0)
	currentDelegatorNum := len(validatorDelegations)

	// Get change delegator num in 24 hours
	if len(latestDelegatorNum) > 0 {
		before24HBondedTokens = latestDelegatorNum[count].DelegatorNum1H // Validator's bonded tokens 24 hours ago
		delegatorNumChange24H = currentDelegatorNum - before24HBondedTokens
	}

	// Result response
	return json.NewEncoder(w).Encode(&models.ResultValidatorDelegations{
		TotalDelegatorNum:     currentDelegatorNum,
		DelegatorNumChange24H: delegatorNumChange24H,
		ValidatorDelegations:  validatorDelegations,
	})

	// u.Respond(w, redelegations)	
	// return nil
}
