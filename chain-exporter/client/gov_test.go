package client

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetProposals(t *testing.T) {
	resp, err := cli.GetProposals()
	require.NoError(t, err)

	log.Println(resp)

}
