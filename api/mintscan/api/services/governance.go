package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	ctypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/sync"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/libs/bech32"
	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

// GetProposals returns all existing proposals
func GetProposals(DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Query all proposals
	proposalInfo := make([]*ctypes.ProposalInfo, 0)
	_ = DB.Model(&proposalInfo).Select()

	// Check if any proposal exists
	if len(proposalInfo) <= 0 {
		return json.NewEncoder(w).Encode(proposalInfo)
	}

	resultProposal := make([]*models.ResultProposal, 0)
	for _, proposal := range proposalInfo {
		// Convert Cosmos Address to Opeartor Address
		_, decoded, _ := bech32.DecodeAndConvert(proposal.Proposer)
		cosmosOperAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)

		// Check if the address matches any moniker in our DB
		var validatorInfo ctypes.ValidatorInfo
		_ = DB.Model(&validatorInfo).
			Column("moniker").
			Where("operator_address = ?", cosmosOperAddress).
			Limit(1).
			Select()

		// Insert proposal data
		tempProposal := &models.ResultProposal{
			ProposalID:           proposal.ID,
			TxHash:               proposal.TxHash,
			Proposer:             proposal.Proposer,
			Moniker:              validatorInfo.Moniker,
			Title:                proposal.Title,
			Description:          proposal.Description,
			ProposalType:         proposal.ProposalType,
			ProposalStatus:       proposal.ProposalStatus,
			Yes:                  proposal.Yes,
			Abstain:              proposal.Abstain,
			No:                   proposal.No,
			NoWithVeto:           proposal.NoWithVeto,
			InitialDepositAmount: proposal.InitialDepositAmount,
			InitialDepositDenom:  proposal.InitialDepositDenom,
			TotalDepositAmount:   proposal.TotalDepositAmount,
			TotalDepositDenom:    proposal.TotalDepositDenom,
			SubmitTime:           proposal.SubmitTime,
			DepositEndtime:       proposal.DepositEndtime,
			VotingStartTime:      proposal.VotingStartTime,
			VotingEndTime:        proposal.VotingEndTime,
		}
		resultProposal = append(resultProposal, tempProposal)
	}
	return json.NewEncoder(w).Encode(resultProposal)
}

// GetProposal receives proposal id and returns particular proposal
func GetProposal(DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Receive proposal id
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// Query particular proposal
	var proposalInfo ctypes.ProposalInfo
	err := DB.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Convert Cosmos Address to Opeartor Address
	_, decoded, _ := bech32.DecodeAndConvert(proposalInfo.Proposer)
	cosmosOperAddress, _ := bech32.ConvertAndEncode(sdk.Bech32PrefixValAddr, decoded)

	// Check if the address matches any moniker in our DB
	var validatorInfo ctypes.ValidatorInfo
	_ = DB.Model(&validatorInfo).
		Column("moniker").
		Where("operator_address = ?", cosmosOperAddress).
		Limit(1).
		Select()

	return json.NewEncoder(w).Encode(&models.ResultProposal{
		ProposalID:           proposalInfo.ID,
		TxHash:               proposalInfo.TxHash,
		Proposer:             proposalInfo.Proposer,
		Moniker:              validatorInfo.Moniker,
		Title:                proposalInfo.Title,
		Description:          proposalInfo.Description,
		ProposalType:         proposalInfo.ProposalType,
		ProposalStatus:       proposalInfo.ProposalStatus,
		Yes:                  proposalInfo.Yes,
		Abstain:              proposalInfo.Abstain,
		No:                   proposalInfo.No,
		NoWithVeto:           proposalInfo.NoWithVeto,
		InitialDepositAmount: proposalInfo.InitialDepositAmount,
		InitialDepositDenom:  proposalInfo.InitialDepositDenom,
		TotalDepositAmount:   proposalInfo.TotalDepositAmount,
		TotalDepositDenom:    proposalInfo.TotalDepositDenom,
		SubmitTime:           proposalInfo.SubmitTime,
		DepositEndtime:       proposalInfo.DepositEndtime,
		VotingStartTime:      proposalInfo.VotingStartTime,
		VotingEndTime:        proposalInfo.VotingEndTime,
	})
}

// GetProposalVotes receives proposal id and returns voting information
func GetProposalVotes(DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Receive proposal id
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// Check if proposal id exists
	var proposalInfo ctypes.ProposalInfo
	err := DB.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query all votes
	voteInfo := make([]*ctypes.VoteInfo, 0)
	_ = DB.Model(&voteInfo).
		Where("proposal_id = ?", proposalID).
		Order("id DESC").
		Select()

	// Check if votes exists
	if len(voteInfo) <= 0 {
		return json.NewEncoder(w).Encode(&models.ResultVoteInfo{
			Tally: &models.Tally{},
			Votes: []*models.Votes{},
		})
	}

	// Query count for respective votes
	yesCnt, _ := DB.Model(&voteInfo).
		Where("proposal_id = ? AND option = ?", proposalID, "Yes").
		Count()
	abstainCnt, _ := DB.Model(&voteInfo).
		Where("proposal_id = ? AND option = ?", proposalID, "Abstain").
		Count()
	noCnt, _ := DB.Model(&voteInfo).
		Where("proposal_id = ? AND option = ?", proposalID, "No").
		Count()
	noWithVetoCnt, _ := DB.Model(&voteInfo).
		Where("proposal_id = ? AND option = ?", proposalID, "NoWithVeto").
		Count()

	// Votes
	votes := make([]*models.Votes, 0)
	for _, vote := range voteInfo {
		moniker := utils.ConvertCosmosAddressToMoniker(vote.Voter, DB)

		tempVoteInfo := &models.Votes{
			Voter:   vote.Voter,
			Moniker: moniker,
			Option:  vote.Option,
			TxHash:  vote.TxHash,
			Time:    vote.Time,
		}
		votes = append(votes, tempVoteInfo)
	}

	// Query LCD
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err := resty.R().Get(Config.Node.LCDURL + "/gov/proposals/" + proposalID + "/tally")
	if err != nil {
		fmt.Printf("Proposal LCD resty - %v\n", err)
	}

	// Parse Tally struct
	var tallyInfo models.TallyInfo
	err = json.Unmarshal(resp.Body(), &tallyInfo)
	if err != nil {
		fmt.Printf("Proposal unmarshal error - %v\n", err)
	}

	// Tally
	tempTally := &models.Tally{
		YesAmount:        tallyInfo.Yes,
		NoAmount:         tallyInfo.No,
		AbstainAmount:    tallyInfo.Abstain,
		NoWithVetoAmount: tallyInfo.NoWithVeto,
		YesNum:           yesCnt,
		AbstainNum:       abstainCnt,
		NoNum:            noCnt,
		NoWithVetoNum:    noWithVetoCnt,
	}

	return json.NewEncoder(w).Encode(&models.ResultVoteInfo{
		Tally: tempTally,
		Votes: votes,
	})
}

// GetProposalDeposits receives proposal id and returns deposit information
func GetProposalDeposits(db *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Receive proposal id
	vars := mux.Vars(r)
	proposalID := vars["proposalId"]

	// Check if proposal id exists
	var proposalInfo ctypes.ProposalInfo
	err := db.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Select()
	if err != nil {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Result Response
	resultDepositInfo := make([]*models.DepositInfo, 0)

	// Query all deposit info
	depositInfo := make([]*ctypes.DepositInfo, 0)
	_ = db.Model(&depositInfo).
		Where("proposal_id = ?", proposalID).
		Order("id DESC").
		Select()

	// Check if the deposit exists
	if len(depositInfo) <= 0 {
		return json.NewEncoder(w).Encode(resultDepositInfo)
	}

	for _, deposit := range depositInfo {
		// Convert Cosmos Address to Opeartor Address
		moniker := utils.ConvertCosmosAddressToMoniker(deposit.Depositor, db)

		// Insert deposits
		tempDepositInfo := &models.DepositInfo{
			Depositor:     deposit.Depositor,
			Moniker:       moniker,
			DepositAmount: deposit.Amount,
			DepositDenom:  deposit.Denom,
			Height:        deposit.Height,
			TxHash:        deposit.TxHash,
			Time:          deposit.Time,
		}
		resultDepositInfo = append(resultDepositInfo, tempDepositInfo)
	}

	return json.NewEncoder(w).Encode(resultDepositInfo)
}

func Test(RPCClient *client.HTTP, DB *pg.DB, w http.ResponseWriter, r *http.Request) {
	// q, _ := tmquery.New("account.owner='Cosmostation'")
	// fmt.Println(q.String())
	// fmt.Println(reflect.TypeOf(q))

	// tx, _ := rpcClient.TxSearch(q.String(), false, 1, 1)

	// fmt.Println(tx)

	type ValidatorDelegations struct {
		DelegatorAddress string  `json:"delegator_address"`
		ValidatorAddress string  `json:"validator_address"`
		Shares           sdk.Dec `json:"shares"`
	}

	type ValidatorDetailInfo struct {
		OperatorAddress string  `json:"operator_address"`
		ConsensusPubkey string  `json:"consensus_pubkey"`
		Jailed          bool    `json:"jailed"`
		Status          int     `json:"status"`
		Tokens          sdk.Dec `json:"tokens"`
		DelegatorShares sdk.Dec `json:"delegator_shares"`
		Description     struct {
			Moniker  string `json:"moniker"`
			Identity string `json:"identity"`
			Website  string `json:"website"`
			Details  string `json:"details"`
		} `json:"description"`
		UnbondingHeight string    `json:"unbonding_height"`
		UnbondingTime   time.Time `json:"unbonding_time"`
		Commission      struct {
			Rate          sdk.Dec   `json:"rate"`
			MaxRate       sdk.Dec   `json:"max_rate"`
			MaxChangeRate sdk.Dec   `json:"max_change_rate"`
			UpdateTime    time.Time `json:"update_time"`
		} `json:"commission"`
		MinSelfDelegation string `json:"min_self_delegation"`
	}

	// Query all validators' operating addresses
	var validatorInfo []ctypes.ValidatorInfo
	_ = DB.Model(&validatorInfo).
		Column("cosmos_address", "operator_address").
		Order("id ASC").
		Select()

	/*
		시도 3 : LCD 서버가 요청을 느리게 받아주는 건지, RPC Full Node가 느리게 받아주는 건지 확인 (lcd-do-not-abuse 로 테스트)
	*/

	// // 여기서부터 test
	// validatorAddr, err := sdktypes.ValAddressFromBech32(operatorAddress)
	// if err != nil {
	// 	fmt.Println("err", err)
	// }
	// params := staking.NewQueryValidatorParams(validatorAddr)

	// bz, err := cdc.MarshalJSON(params)
	// if err != nil {
	// 	fmt.Println("MarshalJSON", err)
	// }

	// opts := rpcclient.ABCIQueryOptions{
	// 	// Height: height,
	// 	// Prove: true,
	// }

	// result, err := client.ABCIQueryWithOptions(fmt.Sprintf("custom/%s/%s", staking.QuerierRoute, staking.QueryValidatorDelegations), bz, opts)
	// if err != nil {
	// 	fmt.Println("ABCIQueryWithOptions", err)
	// }

	// resp := result.Response
	// if !resp.IsOK() {
	// 	fmt.Println("err", err)
	// }

	// var validatorDelegations []*models.ValidatorDelegations
	// err = json.Unmarshal(resp.Value, &validatorDelegations)
	// if err != nil {
	// 	fmt.Printf("staking/validators/{address}/delegations unmarshal error - %v\n", err)
	// }

	/*
		시도 2 : 느림
	*/
	// Query each validator's delegations
	// for _, validator := range validatorInfo {
	// 	validatorResp, _ := resty.R().Get("https://lcd-do-not-abuse.cosmostation.io/staking/validators/" + validator.OperatorAddress)
	// 	validatorDelegationsResp, _ := resty.R().Get("https://lcd-do-not-abuse.cosmostation.io/staking/validators/" + validator.OperatorAddress + "/delegations")

	// 	var totalDelegatorShares float64
	// 	var selfDelegatedShares float64
	// 	var othersShares float64

	// 	// Parse ValidatorDelegations struct
	// 	var validatorDetailInfo ValidatorDetailInfo
	// 	_ = json.Unmarshal(validatorResp.Body(), &validatorDetailInfo)

	// 	// Parse ValidatorDelegations struct
	// 	var validatorDelegations []ValidatorDelegations
	// 	_ = json.Unmarshal(validatorDelegationsResp.Body(), &validatorDelegations)

	// 	validatorCosmosAddress := utils.OperatorAddressToCosmosAddress(validatorDetailInfo.OperatorAddress)
	// 	for _, validatorDelegation := range validatorDelegations {
	// 		// Calculate Self-Delegated Shares
	// 		if validatorDelegation.DelegatorAddress == validatorCosmosAddress {
	// 			selfDelegatedShares, _ = strconv.ParseFloat(validatorDelegation.Shares.String(), 64)
	// 		}
	// 	}

	// 	othersShares, _ = strconv.ParseFloat(validatorDetailInfo.DelegatorShares.String(), 64)
	// 	totalDelegatorShares = selfDelegatedShares + othersShares

	// 	fmt.Println("validator.OperatorAddress: ", validatorDetailInfo.OperatorAddress)
	// 	fmt.Println("totalDelegatorShares: ", totalDelegatorShares)
	// 	fmt.Println("selfDelegatedShares: ", selfDelegatedShares)
	// 	fmt.Println("othersShares: ", othersShares-selfDelegatedShares)
	// 	fmt.Println("")
	// }

	/*
		시도 1 : 정석대로 요청한 결과 느림
	*/
	// Query each validator's delegations
	// for _, validator := range validatorInfo {
	// 	resp, _ := resty.R().Get("https://lcd-do-not-abuse.cosmostation.io/staking/validators/" + validator.OperatorAddress + "/delegations")

	// 	// Parse ValidatorDelegations struct
	// 	var validatorDelegations []ValidatorDelegations
	// 	_ = json.Unmarshal(resp.Body(), &validatorDelegations)

	// 	fmt.Println("validatorDelegations: ", validatorDelegations)
	// 	fmt.Println("")

	// 	var totalDelegatorShares float64
	// 	var selfDelegatedShares float64
	// 	var othersShares float64

	// 	for _, validatorDelegation := range validatorDelegations {
	// 		// Calculate self-delegated and others shares
	// 		if validatorDelegation.DelegatorAddress == validator.CosmosAddress {
	// 			selfDelegatedShares, _ = strconv.ParseFloat(validatorDelegation.Shares.String(), 64)
	// 		} else {
	// 			tempOthersShares, _ := strconv.ParseFloat(validatorDelegation.Shares.String(), 64)
	// 			othersShares += tempOthersShares
	// 		}
	// 	}

	// 	totalDelegatorShares = selfDelegatedShares + othersShares

	// 	fmt.Println("validator.OperatorAddress: ", validator.OperatorAddress)
	// 	fmt.Println("totalDelegatorShares: ", totalDelegatorShares)
	// 	fmt.Println("selfDelegatedShares: ", selfDelegatedShares)
	// 	fmt.Println("othersShares: ", othersShares)
	// 	fmt.Println("")
	// }

}
