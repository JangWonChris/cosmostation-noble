package client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	// sdkUtils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdkcodec "github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	resty "github.com/go-resty/resty/v2"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpc "github.com/tendermint/tendermint/rpc/client/http"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"google.golang.org/grpc"
)

// Client implements a wrapper around both Tendermint RPC HTTP client and
// Cosmos SDK REST client that allow for essential data queries.
type Client struct {
	cliCtx        client.Context
	grpcClient    *grpc.ClientConn
	rpcClient     rpcclient.Client
	apiClient     *resty.Client
	keyBaseClient *resty.Client
}

// NewClient creates a new client with the given configuration and
// return Client struct. An error is returned if it fails.
func NewClient(nodeCfg config.Node, keyBaseURL string) (*Client, error) {
	cliCtx := client.Context{}.
		WithNodeURI(nodeCfg.RPCNode).
		WithJSONMarshaler(codec.EncodingConfig.Marshaler).
		WithLegacyAmino(codec.EncodingConfig.Amino).
		WithTxConfig(codec.EncodingConfig.TxConfig).
		WithInterfaceRegistry(codec.EncodingConfig.InterfaceRegistry).
		WithAccountRetriever(authtypes.AccountRetriever{})

	grpcClient, err := grpc.Dial(nodeCfg.GRPCEndpoint,
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*10),
		grpc.WithInsecure())
	if err != nil {
		return &Client{}, err
	}

	rpcClient, err := rpc.NewWithTimeout(nodeCfg.RPCNode, "/websocket", 10)
	if err != nil {
		return &Client{}, err
	}

	apiClient := resty.New().
		SetHostURL(nodeCfg.LCDEndpoint).
		SetTimeout(time.Duration(10 * time.Second))

	keyBaseClient := resty.New().
		SetHostURL(keyBaseURL).
		SetTimeout(time.Duration(5 * time.Second))

	return &Client{cliCtx, grpcClient, rpcClient, apiClient, keyBaseClient}, nil
}

// Close close the connection
func (c *Client) Close() {
	c.grpcClient.Close()
}

// --------------------
// RPC APIs
// --------------------

// GetNetworkChainID returns network chain id.
func (c *Client) GetNetworkChainID() (string, error) {
	status, err := c.rpcClient.Status(context.Background())
	if err != nil {
		return "", err
	}

	return status.NodeInfo.Network, nil
}

// GetBondDenom returns bond denomination for the network.
func (c *Client) GetBondDenom() (string, error) {
	route := fmt.Sprintf("custom/%s/%s", stakingtypes.StoreKey, stakingtypes.QueryParameters)
	bz, _, err := c.cliCtx.QueryWithData(route, nil)
	if err != nil {
		return "", err
	}

	var params stakingtypes.Params
	c.cliCtx.LegacyAmino.Amino.MustUnmarshalJSON(bz, &params)

	return params.BondDenom, nil
}

// GetStatus queries for status on the active chain.
func (c *Client) GetStatus() (*tmctypes.ResultStatus, error) {
	return c.rpcClient.Status(context.Background())
}

// GetBlock queries for a block with height.
func (c *Client) GetBlock(height int64) (*tmctypes.ResultBlock, error) {
	return c.rpcClient.Block(context.Background(), &height)
}

// GetLatestBlockHeight returns the latest block height on the active network.
func (c *Client) GetLatestBlockHeight() (int64, error) {
	status, err := c.rpcClient.Status(context.Background())
	if err != nil {
		return -1, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

// GetValidators returns all the known Tendermint validators for a given block
// height. An error is returned if the query fails.
func (c *Client) GetValidators(height int64, page int, perPage int) (*tmctypes.ResultValidators, error) {
	return c.rpcClient.Validators(context.Background(), &height, &page, &perPage)
}

// GetGenesisAccounts extracts all genesis accounts from genesis file and return them.
func (c *Client) GetGenesisAccounts() (authtypes.GenesisAccounts, error) {
	gen, err := c.rpcClient.Genesis(context.Background())
	if err != nil {
		return authtypes.GenesisAccounts{}, err
	}

	// jeonghwan : LegacyAmino 로 풀어도 되는지 확인이 필요
	// GetGenesisStateFromAppState() 함수에서 JSONMarshaler 타입만 받음
	appState := make(map[string]json.RawMessage)
	err = c.cliCtx.LegacyAmino.UnmarshalJSON(gen.Genesis.AppState, &appState)
	if err != nil {
		return authtypes.GenesisAccounts{}, err
	}

	genesisState := authtypes.GetGenesisStateFromAppState(c.cliCtx.JSONMarshaler.(sdkcodec.Marshaler), appState)
	accs, err := authtypes.UnpackAccounts(genesisState.Accounts)
	if err != nil {
		return nil, err
	}
	genesisAccts := authtypes.SanitizeGenesisAccounts(accs)

	return genesisAccts, nil
}

// GetTendermintTx queries for a transaction by hash.
func (c *Client) GetTendermintTx(hash string) (*tmctypes.ResultTx, error) {
	hashRaw, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	return c.rpcClient.Tx(context.Background(), hashRaw, false)
}

// GetTendermintTxSearch queries for a transaction search by condition.
// TODO: need more tests. ex:) query := "tx.height=75960",prove := true, page := 1, perPage := 30, orderBy := "asc"
// If this is not needed for this project, let's just remove.
func (c *Client) GetTendermintTxSearch(query string, prove bool, page, perPage int, orderBy string) (*tmctypes.ResultTxSearch, error) {
	txResp, err := c.rpcClient.TxSearch(context.Background(), query, prove, &page, &perPage, orderBy)
	if err != nil {
		return nil, err
	}

	return txResp, nil
}

// GetTxs queries for all the transactions in a block.
// Transactions are returned in the sdktypes.TxResponse format which internally contains an sdktypes.Tx.
func (c *Client) GetTxs(block *tmctypes.ResultBlock) ([]*sdktypes.TxResponse, error) {
	txResponses := make([]*sdktypes.TxResponse, len(block.Block.Txs), len(block.Block.Txs))

	if len(block.Block.Txs) <= 0 {
		return txResponses, nil
	}

	for i, tx := range block.Block.Txs {
		txResponse, err := c.GetTx(fmt.Sprintf("%X", tx.Hash()))
		if err != nil {
			return nil, err
		}

		txResponses[i] = txResponse
	}

	return txResponses, nil
}

// GetTx queries for a single transaction by a hash string in hex format. An
// error is returned if the transaction does not exist or cannot be queried.
func (c *Client) GetTx(hash string) (*sdktypes.TxResponse, error) {
	txResponse, err := authclient.QueryTx(c.cliCtx, hash) // use RPC under the hood
	if err != nil {
		return &sdktypes.TxResponse{}, fmt.Errorf("failed to query tx hash: %s", err)
	}

	if txResponse.Empty() {
		return &sdktypes.TxResponse{}, fmt.Errorf("tx hash has empty tx response: %s", err)
	}

	return txResponse, nil
}

// GetAccount checks account type and returns account interface.
func (c *Client) GetAccount(address string) (authtypes.AccountI, error) {
	accAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	ar := authtypes.AccountRetriever{}
	acc, err := ar.GetAccount(c.cliCtx, accAddr)
	if err != nil {
		log.Println(err)
	}
	// acc, err := auth.NewAccountRetriever(c.cliCtx).GetAccount(accAddr)
	// if err != nil {
	// 	return nil, err
	// }
	// jeonghwan : need grpc
	// queryClient := authtypes.NewQueryClient(c.cliCtx)
	// res, err := queryClient.Account(context.Background(), &authtypes.QueryAccountRequest{Address: accAddr.String()})
	// if err != nil {
	// 	return nil, err
	// }

	// out, err := c.cliCtx.LegacyAmino.MarshalJSON(res.Account)
	// if err != nil {
	// 	return nil, err
	// }
	// var acc authtypes.AccountI
	// c.cliCtx.LegacyAmino.UnmarshalJSON(out, &acc)

	return acc, nil
}

// GetValidatorCommission queries validator's commission and returns the coins with truncated decimals and the change.
func (c *Client) GetValidatorCommission(address string) (sdktypes.Coins, error) {
	valAddr, err := sdktypes.ValAddressFromBech32(address)
	if err != nil {
		return sdktypes.Coins{}, err
	}

	res, err := common.QueryValidatorCommission(c.cliCtx, valAddr)
	if err != nil {
		return sdktypes.Coins{}, err
	}

	var valCom distrtypes.ValidatorAccumulatedCommission
	c.cliCtx.LegacyAmino.MustUnmarshalJSON(res, &valCom)

	truncatedCoins, _ := valCom.Commission.TruncateDecimal()

	return truncatedCoins, nil
}

// GetDelegatorDelegations returns a list of delegations made by a certain delegator address
func (c *Client) GetDelegatorDelegations(address string) (stakingtypes.DelegationResponses, error) {
	delAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return stakingtypes.DelegationResponses{}, err
	}

	bz, err := c.cliCtx.LegacyAmino.MarshalJSON(stakingtypes.NewQueryDelegatorParams(delAddr))
	if err != nil {
		return stakingtypes.DelegationResponses{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", stakingtypes.QuerierRoute, stakingtypes.QueryDelegatorDelegations)
	res, _, err := c.cliCtx.QueryWithData(route, bz)
	if err != nil {
		return stakingtypes.DelegationResponses{}, err
	}

	var delegations stakingtypes.DelegationResponses
	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(res, &delegations); err != nil {
		return stakingtypes.DelegationResponses{}, err
	}

	return delegations, nil
}

// GetDelegatorUndelegations returns a list of undelegations made by a certain delegator address
func (c *Client) GetDelegatorUndelegations(address string) (stakingtypes.UnbondingDelegations, error) {
	delAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return stakingtypes.UnbondingDelegations{}, err
	}

	bz, err := c.cliCtx.LegacyAmino.MarshalJSON(stakingtypes.NewQueryDelegatorParams(delAddr))
	if err != nil {
		return stakingtypes.UnbondingDelegations{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", stakingtypes.QuerierRoute, stakingtypes.QueryDelegatorUnbondingDelegations)
	res, _, err := c.cliCtx.QueryWithData(route, bz)
	if err != nil {
		return stakingtypes.UnbondingDelegations{}, err
	}

	var undelegations stakingtypes.UnbondingDelegations
	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(res, &undelegations); err != nil {
		return stakingtypes.UnbondingDelegations{}, err
	}

	return undelegations, nil
}

// GetDelegatorTotalRewards returns the total rewards balance from all delegations by a delegator
func (c *Client) GetDelegatorTotalRewards(address string) (distrtypes.QueryDelegatorTotalRewardsResponse, error) {
	delAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return distrtypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	bz, err := c.cliCtx.LegacyAmino.MarshalJSON(distrtypes.NewQueryDelegatorParams(delAddr))
	if err != nil {
		return distrtypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", distrtypes.QuerierRoute, distrtypes.QueryDelegatorTotalRewards)
	res, _, err := c.cliCtx.QueryWithData(route, bz)
	if err != nil {
		return distrtypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	var totalRewards distrtypes.QueryDelegatorTotalRewardsResponse
	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(res, &totalRewards); err != nil {
		return distrtypes.QueryDelegatorTotalRewardsResponse{}, err
	}

	return totalRewards, nil
}

// GetBaseAccountTotalAsset returns total available, rewards, commission, delegations, and undelegations from a delegator.
// func (c *Client) GetBaseAccountTotalAsset(address string) (sdktypes.Coin, sdktypes.Coin, sdktypes.Coin, sdktypes.Coin, sdktypes.Coin, error) {
// 	account, err := c.GetAccount(address)
// 	if err != nil {
// 		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
// 	}

// 	denom, err := c.GetBondDenom()
// 	if err != nil {
// 		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
// 	}

// 	spendable := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
// 	delegated := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
// 	undelegated := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
// 	rewards := sdktypes.NewCoin(denom, sdktypes.NewInt(0))
// 	commission := sdktypes.NewCoin(denom, sdktypes.NewInt(0))

// 	// Get total spendable coins.
// 	if len(account.GetCoins()) > 0 {
// 		for _, coin := range account.GetCoins() {
// 			if coin.Denom == denom {
// 				spendable = spendable.Add(coin)
// 			}
// 		}
// 	}

// 	// Get total delegated coins.
// 	delegations, err := c.GetDelegatorDelegations(address)
// 	if err != nil {
// 		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
// 	}

// 	if len(delegations) > 0 {
// 		for _, delegation := range delegations {
// 			delegated = delegated.Add(delegation.Balance)
// 		}
// 	}

// 	// Get total undelegated coins.
// 	undelegations, err := c.GetDelegatorUndelegations(address)
// 	if err != nil {
// 		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
// 	}

// 	if len(undelegations) > 0 {
// 		for _, undelegation := range undelegations {
// 			for _, e := range undelegation.Entries {
// 				undelegated = undelegated.Add(sdktypes.NewCoin(denom, e.Balance))
// 			}
// 		}
// 	}

// 	// Get total rewarded coins.
// 	totalRewards, err := c.GetDelegatorTotalRewards(address)
// 	if err != nil {
// 		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
// 	}

// 	if len(totalRewards.Rewards) > 0 {
// 		for _, tr := range totalRewards.Rewards {
// 			for _, reward := range tr.Reward {
// 				if reward.Denom == denom {
// 					truncatedRewards, _ := reward.TruncateDecimal()
// 					rewards = rewards.Add(truncatedRewards)
// 				}
// 			}
// 		}
// 	}

// 	valAddr, err := types.ConvertValAddrFromAccAddr(address)
// 	if err != nil {
// 		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
// 	}

// 	// Get commission
// 	commissions, err := c.GetValidatorCommission(valAddr)
// 	if err != nil {
// 		return sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, sdktypes.Coin{}, err
// 	}

// 	if len(commissions) > 0 {
// 		for _, c := range commissions {
// 			commission = commission.Add(c)
// 		}
// 	}

// 	return spendable, delegated, undelegated, rewards, commission, nil
// }

// --------------------
// REST SERVER APIs
// --------------------

// GetTxAPIClient queries for a transaction from the REST client and decodes it into a sdktypes.Tx [Another way to query a transaction.]
// if the transaction exists. An error is returned if the tx doesn't exist or
// decoding fails.
func (c *Client) GetTxAPIClient(hash string) (sdktypes.TxResponse, error) {
	resp, err := c.apiClient.R().Get("/txs/" + hash)
	if err != nil {
		return sdktypes.TxResponse{}, fmt.Errorf("failed to request tx hash: %s", err)
	}

	var txResponse sdktypes.TxResponse
	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(resp.Body(), &txResponse); err != nil {
		return sdktypes.TxResponse{}, fmt.Errorf("failed to unmarshal tx hash: %s", err)
	}

	return txResponse, nil
}

// GetProposals returns all governance proposals
func (c *Client) GetProposals() (result []schema.Proposal, err error) {
	resp, err := c.apiClient.R().Get("/gov/proposals")
	if err != nil {
		return []schema.Proposal{}, fmt.Errorf("failed to request gov proposals: %s", err)
	}

	var proposals []types.Proposal
	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &proposals)
	if err != nil {
		return []schema.Proposal{}, fmt.Errorf("failed to unmarshal gov proposals: %s", err)
	}

	if len(proposals) <= 0 {
		return []schema.Proposal{}, nil
	}

	for _, proposal := range proposals {
		proposalID, _ := strconv.ParseInt(proposal.ID, 10, 64)

		var totalDepositAmount string
		var totalDepositDenom string
		if proposal.TotalDeposit != nil {
			totalDepositAmount = proposal.TotalDeposit[0].Amount
			totalDepositDenom = proposal.TotalDeposit[0].Denom
		}

		resp, err := c.apiClient.R().Get("/gov/proposals/" + proposal.ID + "/tally")
		if err != nil {
			return []schema.Proposal{}, fmt.Errorf("failed to request gov tally: %s", err)
		}

		var tally types.Tally
		err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &tally)
		if err != nil {
			return []schema.Proposal{}, fmt.Errorf("failed to unmarshal gov tally: %s", err)
		}

		p := schema.NewProposal(schema.Proposal{
			ID:                 proposalID,
			Title:              proposal.Content.Value.Title,
			Description:        proposal.Content.Value.Description,
			ProposalType:       proposal.Content.Type,
			ProposalStatus:     proposal.ProposalStatus,
			Yes:                tally.Yes,
			Abstain:            tally.Abstain,
			No:                 tally.No,
			NoWithVeto:         tally.NoWithVeto,
			SubmitTime:         proposal.SubmitTime,
			DepositEndtime:     proposal.DepositEndTime,
			TotalDepositAmount: totalDepositAmount,
			TotalDepositDenom:  totalDepositDenom,
			VotingStartTime:    proposal.VotingStartTime,
			VotingEndTime:      proposal.VotingEndTime,
		})

		result = append(result, *p)
	}

	return result, nil
}

// GetBondedValidators returns all bonded validators
func (c *Client) GetBondedValidators() (validators []schema.Validator, err error) {
	resp, err := c.apiClient.R().Get("/staking/validators?status=bonded")
	if err != nil {
		return []schema.Validator{}, fmt.Errorf("failed to request bonded vals: %s", err)
	}

	var bondedVals []types.Validator
	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &bondedVals)
	if err != nil {
		return []schema.Validator{}, fmt.Errorf("failed to unmarshal bonded vals: %s", err)
	}

	// Sort bondedVals by highest token amount
	sort.Slice(bondedVals[:], func(i, j int) bool {
		tempTk1, _ := strconv.Atoi(bondedVals[i].Tokens)
		tempTk2, _ := strconv.Atoi(bondedVals[j].Tokens)
		return tempTk1 > tempTk2
	})

	if len(bondedVals) <= 0 {
		return []schema.Validator{}, nil
	}

	for i, val := range bondedVals {
		accAddr, _ := types.ConvertAccAddrFromValAddr(val.OperatorAddress)
		consAddr, _ := types.ConvertConsAddrFromConsPubkey(val.ConsensusPubkey)

		v := schema.NewValidator(schema.Validator{
			Rank:                 i + 1,
			OperatorAddress:      val.OperatorAddress,
			Address:              accAddr,
			ConsensusPubkey:      val.ConsensusPubkey,
			Proposer:             consAddr,
			Jailed:               val.Jailed,
			Status:               val.Status,
			Tokens:               val.Tokens,
			DelegatorShares:      val.DelegatorShares,
			Moniker:              val.Description.Moniker,
			Identity:             val.Description.Identity,
			Website:              val.Description.Website,
			Details:              val.Description.Details,
			UnbondingHeight:      val.UnbondingHeight,
			UnbondingTime:        val.UnbondingTime,
			CommissionRate:       val.Commission.CommissionRates.Rate,
			CommissionMaxRate:    val.Commission.CommissionRates.MaxRate,
			CommissionChangeRate: val.Commission.CommissionRates.MaxChangeRate,
			MinSelfDelegation:    val.MinSelfDelegation,
			UpdateTime:           val.Commission.UpdateTime,
		})

		validators = append(validators, *v)
	}

	return validators, nil
}

// GetUnbondingValidators returns unbonding validators
func (c *Client) GetUnbondingValidators() (validators []schema.Validator, err error) {
	resp, err := c.apiClient.R().Get("/staking/validators?status=unbonding")
	if err != nil {
		return []schema.Validator{}, fmt.Errorf("failed to request unbonding vals: %s", err)
	}

	var unbondingVals []*types.Validator
	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &unbondingVals)
	if err != nil {
		return []schema.Validator{}, fmt.Errorf("failed to unmarshal unbonding vals: %s", err)
	}

	// Sort bondedValidators by highest token amount
	sort.Slice(unbondingVals[:], func(i, j int) bool {
		tempTk1, _ := strconv.Atoi(unbondingVals[i].Tokens)
		tempTk2, _ := strconv.Atoi(unbondingVals[j].Tokens)
		return tempTk1 > tempTk2
	})

	if len(unbondingVals) <= 0 {
		return []schema.Validator{}, nil
	}

	for _, val := range unbondingVals {
		accAddr, _ := types.ConvertAccAddrFromValAddr(val.OperatorAddress)
		consAddr, _ := types.ConvertConsAddrFromConsPubkey(val.ConsensusPubkey)

		v := schema.NewValidator(schema.Validator{
			OperatorAddress:      val.OperatorAddress,
			Address:              accAddr,
			ConsensusPubkey:      val.ConsensusPubkey,
			Proposer:             consAddr,
			Jailed:               val.Jailed,
			Status:               val.Status,
			Tokens:               val.Tokens,
			DelegatorShares:      val.DelegatorShares,
			Moniker:              val.Description.Moniker,
			Identity:             val.Description.Identity,
			Website:              val.Description.Website,
			Details:              val.Description.Details,
			UnbondingHeight:      val.UnbondingHeight,
			UnbondingTime:        val.UnbondingTime,
			CommissionRate:       val.Commission.CommissionRates.Rate,
			CommissionMaxRate:    val.Commission.CommissionRates.MaxRate,
			CommissionChangeRate: val.Commission.CommissionRates.MaxChangeRate,
			MinSelfDelegation:    val.MinSelfDelegation,
			UpdateTime:           val.Commission.UpdateTime,
		})

		validators = append(validators, *v)
	}

	return validators, nil
}

// GetUnbondedValidators returns unbonded validators
func (c *Client) GetUnbondedValidators() (validators []schema.Validator, err error) {
	resp, err := c.apiClient.R().Get("/staking/validators?status=unbonded")
	if err != nil {
		return []schema.Validator{}, fmt.Errorf("failed to request unbonded vals: %s", err)
	}

	var unbondedVals []types.Validator
	err = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &unbondedVals)
	if err != nil {
		return []schema.Validator{}, fmt.Errorf("failed to unmarshal unbonded vals: %s", err)
	}

	// Sort bondedValidators by highest token amount
	sort.Slice(unbondedVals[:], func(i, j int) bool {
		tempTk1, _ := strconv.Atoi(unbondedVals[i].Tokens)
		tempTk2, _ := strconv.Atoi(unbondedVals[j].Tokens)
		return tempTk1 > tempTk2
	})

	if len(unbondedVals) <= 0 {
		return []schema.Validator{}, nil
	}

	for _, val := range unbondedVals {
		accAddr, _ := types.ConvertAccAddrFromValAddr(val.OperatorAddress)
		consAddr, _ := types.ConvertConsAddrFromConsPubkey(val.ConsensusPubkey)

		v := schema.NewValidator(schema.Validator{
			OperatorAddress:      val.OperatorAddress,
			Address:              accAddr,
			ConsensusPubkey:      val.ConsensusPubkey,
			Proposer:             consAddr,
			Jailed:               val.Jailed,
			Status:               val.Status,
			Tokens:               val.Tokens,
			DelegatorShares:      val.DelegatorShares,
			Moniker:              val.Description.Moniker,
			Identity:             val.Description.Identity,
			Website:              val.Description.Website,
			Details:              val.Description.Details,
			UnbondingHeight:      val.UnbondingHeight,
			UnbondingTime:        val.UnbondingTime,
			CommissionRate:       val.Commission.CommissionRates.Rate,
			CommissionMaxRate:    val.Commission.CommissionRates.MaxRate,
			CommissionChangeRate: val.Commission.CommissionRates.MaxChangeRate,
			MinSelfDelegation:    val.MinSelfDelegation,
			UpdateTime:           val.Commission.UpdateTime,
		})

		validators = append(validators, *v)
	}

	return validators, nil
}

// GetValidatorsIdentities returns identities of all validators in the active chain.
func (c *Client) GetValidatorsIdentities(vals []schema.Validator) (result []schema.Validator, err error) {
	for _, val := range vals {
		if val.Identity != "" {
			resp, err := c.keyBaseClient.R().Get("_/api/1.0/user/lookup.json?fields=pictures&key_suffix=" + val.Identity)
			if err != nil {
				return []schema.Validator{}, fmt.Errorf("failed to request identity: %s", err)
			}

			var keyBase types.KeyBase
			err = json.Unmarshal(resp.Body(), &keyBase)
			if err != nil {
				return []schema.Validator{}, fmt.Errorf("failed to unmarshal keybase: %s", err)
			}

			var url string
			if len(keyBase.Them) > 0 {
				for _, k := range keyBase.Them {
					url = k.Pictures.Primary.URL
				}
			}

			v := schema.NewValidator(schema.Validator{
				ID:         val.ID,
				KeybaseURL: url,
			})

			result = append(result, *v)
		}
	}

	return result, nil
}
