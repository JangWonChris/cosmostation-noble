package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdkCodec "github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkUtils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmcTypes "github.com/tendermint/tendermint/rpc/core/types"

	resty "github.com/go-resty/resty/v2"
)

// Client implements a wrapper around both Tendermint RPC HTTP client and
// Cosmos SDK REST client that allow for essential data queries.
type Client struct {
	cliCtx        context.CLIContext
	cdc           *sdkCodec.Codec
	rpcClient     rpcclient.Client
	apiClient     *resty.Client
	keyBaseClient *resty.Client
}

// NewClient creates a new client with the given configuration and
// return Client struct. An error is returned if it fails.
func NewClient(nodeCfg config.Node, keyBaseURL string) (*Client, error) {
	cliCtx := context.NewCLIContext().
		WithCodec(codec.Codec).
		WithNodeURI(nodeCfg.RPCNode).
		WithTrustNode(true)

	rpcClient := rpcclient.NewHTTP(nodeCfg.RPCNode, "/websocket")
	if err := rpcClient.Start(); err != nil {
		return &Client{}, err
	}
	apiClient := resty.New().
		SetHostURL(nodeCfg.LCDEndpoint).
		SetTimeout(time.Duration(10 * time.Second))

	keyBaseClient := resty.New().
		SetHostURL(keyBaseURL).
		SetTimeout(time.Duration(5 * time.Second))

	return &Client{cliCtx, codec.Codec, rpcClient, apiClient, keyBaseClient}, nil
}

// --------------------
// RPC APIs
// --------------------

// GetNetworkChainID returns network chain id.
func (c *Client) GetNetworkChainID() (string, error) {
	status, err := c.rpcClient.Status()
	if err != nil {
		return "", err
	}

	return status.NodeInfo.Network, nil
}

// GetBondDenom returns bond denomination for the network.
func (c *Client) GetBondDenom() (string, error) {
	route := fmt.Sprintf("custom/%s/%s", stakingTypes.StoreKey, stakingTypes.QueryParameters)
	bz, _, err := c.cliCtx.QueryWithData(route, nil)
	if err != nil {
		return "", err
	}

	var params stakingTypes.Params
	c.cdc.MustUnmarshalJSON(bz, &params)

	return params.BondDenom, nil
}

// GetStatus queries for status on the active chain.
func (c *Client) GetStatus() (*tmcTypes.ResultStatus, error) {
	return c.rpcClient.Status()
}

// GetBlock queries for a block with height.
func (c *Client) GetBlock(height int64) (*tmcTypes.ResultBlock, error) {
	return c.rpcClient.Block(&height)
}

// GetLatestBlockHeight returns the latest block height on the active network.
func (c *Client) GetLatestBlockHeight() (int64, error) {
	status, err := c.rpcClient.Status()
	if err != nil {
		return -1, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

// GetValidators returns all the known Tendermint validators for a given block
// height. An error is returned if the query fails.
func (c *Client) GetValidators(height int64, page int, perPage int) (*tmcTypes.ResultValidators, error) {
	// This method will change to c.rpcClient.Validators(&height, page, perPage) in above v0.38.+
	return c.rpcClient.Validators(&height)
}

// GetTendermintTx queries for a transaction by hash.
func (c *Client) GetTendermintTx(hash string) (*tmcTypes.ResultTx, error) {
	hashRaw, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	return c.rpcClient.Tx(hashRaw, false)
}

// GetTxs queries for all the transactions in a block.
// Transactions are returned in the sdkTypes.TxResponse format which internally contains an sdkTypes.Tx.
func (c *Client) GetTxs(block *tmcTypes.ResultBlock) ([]*sdkTypes.TxResponse, error) {
	txResponses := make([]*sdkTypes.TxResponse, len(block.Block.Txs), len(block.Block.Txs))

	if len(block.Block.Txs) <= 0 {
		return txResponses, nil
	}

	for i, tx := range block.Block.Txs {
		txResponse, err := c.GetTx(fmt.Sprintf("%X", tx.Hash()))
		if err != nil {
			return nil, err
		}

		txResponses[i] = &txResponse
	}

	return txResponses, nil
}

// GetTx queries for a single transaction by a hash string in hex format. An
// error is returned if the transaction does not exist or cannot be queried.
func (c *Client) GetTx(hash string) (sdkTypes.TxResponse, error) {
	txResponse, err := sdkUtils.QueryTx(c.cliCtx, hash) // use RPC under the hood
	if err != nil {
		return sdkTypes.TxResponse{}, fmt.Errorf("failed to query tx hash: %s", err)
	}

	if txResponse.Empty() {
		return sdkTypes.TxResponse{}, fmt.Errorf("tx hash has empty tx response: %s", err)
	}

	return txResponse, nil
}

// GetAccount checks account type and returns account interface.
func (c *Client) GetAccount(address string) (exported.Account, error) {
	accAddr, err := sdkTypes.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	acc, err := auth.NewAccountRetriever(c.cliCtx).GetAccount(accAddr)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

// GetValidatorCommission queries validator's commission and returns the coins with truncated decimals and the change.
func (c *Client) GetValidatorCommission(address string) (sdkTypes.Coins, error) {
	valAddr, err := sdkTypes.ValAddressFromBech32(address)
	if err != nil {
		return sdkTypes.Coins{}, err
	}

	res, err := common.QueryValidatorCommission(c.cliCtx, distr.QuerierRoute, valAddr)
	if err != nil {
		return sdkTypes.Coins{}, err
	}

	var valCom distr.ValidatorAccumulatedCommission
	c.cliCtx.Codec.MustUnmarshalJSON(res, &valCom)

	truncatedCoins, _ := valCom.TruncateDecimal()

	return truncatedCoins, nil
}

// GetDelegatorDelegations returns a list of delegations made by a certain delegator address
func (c *Client) GetDelegatorDelegations(address string) (staking.DelegationResponses, error) {
	delAddr, err := sdkTypes.AccAddressFromBech32(address)
	if err != nil {
		return staking.DelegationResponses{}, err
	}

	bz, err := c.cdc.MarshalJSON(staking.NewQueryDelegatorParams(delAddr))
	if err != nil {
		return staking.DelegationResponses{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", staking.QuerierRoute, staking.QueryDelegatorDelegations)
	res, _, err := c.cliCtx.QueryWithData(route, bz)
	if err != nil {
		return staking.DelegationResponses{}, err
	}

	var delegations staking.DelegationResponses
	if err := c.cdc.UnmarshalJSON(res, &delegations); err != nil {
		return staking.DelegationResponses{}, err
	}

	return delegations, nil
}

// --------------------
// REST SERVER APIs
// --------------------

// GetTxAPIClient queries for a transaction from the REST client and decodes it into a sdkTypes.Tx [Another way to query a transaction.]
// if the transaction exists. An error is returned if the tx doesn't exist or
// decoding fails.
func (c *Client) GetTxAPIClient(hash string) (sdkTypes.TxResponse, error) {
	resp, err := c.apiClient.R().Get("/txs/" + hash)
	if err != nil {
		return sdkTypes.TxResponse{}, fmt.Errorf("failed to request tx hash: %s", err)
	}

	var txResponse sdkTypes.TxResponse
	if err := c.cdc.UnmarshalJSON(resp.Body(), &txResponse); err != nil {
		return sdkTypes.TxResponse{}, fmt.Errorf("failed to unmarshal tx hash: %s", err)
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
		// accAddr := types.AccAddressFromOperatorAddress(val.OperatorAddress)
		// consAddr := types.ConsAddrFromConsPubkey(val.ConsensusPubkey)
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
		// accAddr := types.AccAddressFromOperatorAddress(val.OperatorAddress)
		// consAddr := types.ConsAddrFromConsPubkey(val.ConsensusPubkey)
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
