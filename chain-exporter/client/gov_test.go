package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetProposals(t *testing.T) {
	resp, err := cli.GetProposals()
	require.NoError(t, err)
	t.Log(resp)
}
