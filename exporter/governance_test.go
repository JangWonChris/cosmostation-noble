package exporter

import (
	"testing"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/stretchr/testify/require"
)

func TestSaveAllProposals(t *testing.T) {
	ex.saveAllProposals()
}

func TestGetProposalVyStatus(t *testing.T) {
	ex.saveAllProposals()
	dp, err := ex.Client.GetProposalsByStatus(govtypes.StatusDepositPeriod)
	require.NoError(t, err)
	t.Log(dp)
}
