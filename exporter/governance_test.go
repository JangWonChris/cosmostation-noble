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

func TestSaveAllProposalsWithoutCondition(t *testing.T) {
	proposals, err := ex.Client.GetAllProposals()
	require.NoError(t, err)

	if len(proposals) <= 0 {
		t.Log("found empty proposals")
		return
	}

	err = ex.DB.InsertOrUpdateProposals(proposals)
	require.NoError(t, err)
}
