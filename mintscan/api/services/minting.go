package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models/types"
	u "github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/go-pg/pg"
	resty "gopkg.in/resty.v1"
)

// GetMintingInflation returns minting inflation rate
func GetMintingInflation(config *config.Config, db *pg.DB, w http.ResponseWriter, r *http.Request) error {
	resp, _ := resty.R().Get(config.Node.LCDURL + "/minting/inflation")

	var tempInflation string
	_ = json.Unmarshal(types.ReadRespWithHeight(resp).Result, &tempInflation)

	inflation, _ := strconv.ParseFloat(tempInflation, 64)

	resultInflation := &models.ResultInflation{
		Inflation: inflation,
	}

	u.Respond(w, resultInflation)
	return nil
}
