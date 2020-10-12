package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
)

func TestParseConfig(t *testing.T) {
	config := ParseConfig()

	require.NotEmpty(t, config.Node.RPCNode, "Node RPCNode should not be empty")
	require.NotEmpty(t, config.Node.LCDEndpoint, "Node LCDEndpoint should not be empty")
	require.NotEmpty(t, config.DB.Host, "Database Host is empty")
	require.NotEmpty(t, config.DB.Port, "Database Port is empty")
	require.NotEmpty(t, config.DB.User, "Database User is empty")
	require.NotEmpty(t, config.DB.Password, "Database Password is empty")
	require.NotEmpty(t, config.DB.Table, "Database Table is empty")
	require.NotEmpty(t, config.Web.Port, "Web Port is empty")
	require.NotEmpty(t, config.Market.CoinGeckoEndpoint, "CoinGecko Endpoint is empty")
	require.NotEmpty(t, config.Market.CoinmarketCapEndpoint, "Coinmarketcap Endpoint is empty")
	require.NotEmpty(t, config.Market.CoinmarketCapCoinID, "Coinmarketcap CoinID is empty")
	require.NotEmpty(t, config.Market.CoinmarketCapAPIKey, "Coinmarketcap APIKey is empty")
	assert.IsType(t, false, Common.Maintenance, "Maintenance")
}
