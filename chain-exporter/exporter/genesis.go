package exporter

import (
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getGenesisValidatorsSet returns validator set in genesis.
func (ex *Exporter) getGenesisValidatorsSet(block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) ([]schema.PowerEventHistory, error) {
	genesisValsSet := make([]schema.PowerEventHistory, 0)

	if block.Block.Height != 1 {
		return []schema.PowerEventHistory{}, nil
	}

	denom, err := ex.client.GetBondDenom()
	if err != nil {
		return []schema.PowerEventHistory{}, err
	}

	// Get genesis validator set (block height 1).
	for i, val := range vals.Validators {
		gvs := schema.NewPowerEventHistoryForGenesisValidatorSet(schema.PowerEventHistory{
			IDValidator:          i + 1,
			Height:               block.Block.Height,
			Moniker:              "",
			OperatorAddress:      "",
			Proposer:             val.Address.String(),
			VotingPower:          float64(val.VotingPower),
			MsgType:              types.TypeMsgCreateValidator,
			NewVotingPowerAmount: float64(val.VotingPower),
			NewVotingPowerDenom:  denom,
			TxHash:               "",
			Timestamp:            block.Block.Header.Time,
		})

		genesisValsSet = append(genesisValsSet, *gvs)
	}

	return genesisValsSet, nil
}
