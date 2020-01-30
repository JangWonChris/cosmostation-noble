package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	u "github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	resty "gopkg.in/resty.v1"
)

// GetMintingInflation returns minting inflation rate
func GetMintingInflation(config *config.Config, db *db.Database, w http.ResponseWriter, r *http.Request) error {
	resp, _ := resty.R().Get(config.Node.LCDEndpoint + "/minting/inflation")

	var tempInflation string
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &tempInflation)
	if err != nil {
		fmt.Printf("failed to unmarshal tempInflation: %t\n", err)
	}

	inflationRate, _ := strconv.ParseFloat(tempInflation, 64)

	result := &models.ResultInflation{
		Inflation: inflationRate,
	}

	u.Respond(w, result)
	return nil
}
