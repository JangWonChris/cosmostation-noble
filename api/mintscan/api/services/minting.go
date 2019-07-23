package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	u "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	"github.com/go-pg/pg"
	resty "gopkg.in/resty.v1"
)

// GetMinting returns minting parameters
func GetMintingInflation(config *config.Config, db *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Query inflation
	resp, _ := resty.R().Get(config.Node.LCDURL + "/minting/inflation")

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
