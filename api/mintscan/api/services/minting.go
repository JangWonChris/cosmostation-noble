package services

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	u "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	"github.com/go-pg/pg"
	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

// GetMinting returns minting parameters
func GetMintingInflation(RPCClient *client.HTTP, DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Query inflation
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, _ := resty.R().Get(Config.Node.LCDURL + "/minting/inflation")

	var tempInflation string
	_ = json.Unmarshal(resp.Body(), &tempInflation)

	// Conversion
	inflation, _ := strconv.ParseFloat(tempInflation, 64)

	resultMinting := &models.ResultMinting{
		Inflation: inflation,
	}

	u.Respond(w, resultMinting)
	return nil
}
