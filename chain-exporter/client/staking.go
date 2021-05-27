package client

import (
	"context"
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	"github.com/cosmostation/mintscan-backend-library/types"
	"github.com/cosmostation/mintscan-database/schema"

	//cosmos-sdk
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetValidatorsByStatus 는 MBL의 GetValidatorsByStatus을 wrap한 함수 (codec 사용을 분리하기 위해)
// 필요한 함수를 우선 모듈에 맞게 정의한 후, 나중에 코어로 이전
// 코어로 분리가 가능할 것 같다.
func (c *Client) GetValidatorsByStatus(ctx context.Context, status stakingtypes.BondStatus) (validators []schema.Validator, err error) {
	res, err := c.GRPC.GetValidatorsByStatus(ctx, status)
	if err != nil {
		return []schema.Validator{}, nil
	}

	if res == nil {
		return []schema.Validator{}, nil
	}

	for i, val := range res.Validators {
		accAddr, err := types.ConvertAccAddrFromValAddr(val.OperatorAddress)
		if err != nil {
			return []schema.Validator{}, fmt.Errorf("failed to convert address from validator Address : %s", err)
		}

		var conspubkey cryptotypes.PubKey
		custom.AppCodec.UnpackAny(val.ConsensusPubkey, &conspubkey)

		valconspub, err := sdktypes.Bech32ifyPubKey(sdktypes.Bech32PubKeyTypeConsPub, conspubkey)
		if err != nil {
			return []schema.Validator{}, fmt.Errorf("failed to get consesnsus pubkey : %s", err)
		}

		// log.Println("conspubkey get cached value : ", val.ConsensusPubkey.GetCachedValue())
		// conspubkey, err := val.TmConsPubKey()
		// if err != nil {
		// 	return []schema.Validator{}, fmt.Errorf("failed to get consesnsus pubkey : %s", err)
		// }

		v := schema.Validator{
			Rank:                 i + 1,
			OperatorAddress:      val.OperatorAddress,
			Address:              accAddr,
			ConsensusPubkey:      valconspub,
			Proposer:             conspubkey.Address().String(),
			Jailed:               val.Jailed,
			Status:               int(val.Status),
			Tokens:               val.Tokens.String(),
			DelegatorShares:      val.DelegatorShares.String(),
			Moniker:              val.Description.Moniker,
			Identity:             val.Description.Identity,
			Website:              val.Description.Website,
			Details:              val.Description.Details,
			UnbondingHeight:      val.UnbondingHeight,
			UnbondingTime:        val.UnbondingTime,
			CommissionRate:       val.Commission.CommissionRates.Rate.String(),
			CommissionMaxRate:    val.Commission.CommissionRates.MaxRate.String(),
			CommissionChangeRate: val.Commission.CommissionRates.MaxChangeRate.String(),
			MinSelfDelegation:    val.MinSelfDelegation.String(),
			UpdateTime:           val.Commission.UpdateTime,
		}

		validators = append(validators, v)
	}

	return validators, nil
}
