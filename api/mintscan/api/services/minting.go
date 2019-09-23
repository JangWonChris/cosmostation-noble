package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
	u "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	"github.com/go-pg/pg"
	resty "gopkg.in/resty.v1"
)

// GetMintingInflation returns minting parameters
func GetMintingInflation(config *config.Config, db *pg.DB, w http.ResponseWriter, r *http.Request) error {
	// Query inflation
	resp, _ := resty.R().Get(config.Node.LCDURL + "/minting/inflation")

	var responseWithHeight types.ResponseWithHeight
	err := json.Unmarshal(resp.Body(), &responseWithHeight)
	if err != nil {
		fmt.Printf("unmarshal responseWithHeight error - %v\n", err)
	}

	var tempInflation string
	_ = json.Unmarshal(responseWithHeight.Result, &tempInflation)

	inflation, _ := strconv.ParseFloat(tempInflation, 64)

	resultInflation := &models.ResultInflation{
		Inflation: inflation,
	}

	u.Respond(w, resultInflation)
	return nil
}
