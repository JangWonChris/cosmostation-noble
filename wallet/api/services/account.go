package services

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/wallet/api/databases"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/models"
	u "github.com/cosmostation/cosmostation-cosmos/wallet/api/utils"
	"github.com/go-pg/pg"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Regsiter registers an account for our mobile users
func Register(db *pg.DB, w http.ResponseWriter, r *http.Request) {
	var account models.Account

	// get post data from request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		errors.ErrBadRequest(w, http.StatusBadRequest)
		return
	}

	// lower case
	account.Address = strings.ToLower(account.Address)
	account.DeviceType = strings.ToLower(account.DeviceType)

	// check device type
	if account.DeviceType != models.Android && account.DeviceType != models.IOS {
		errors.ErrInvalidDeviceType(w, http.StatusBadRequest)
		return
	}

	// check chain id
	if account.ChainID != models.CosmosHub && account.ChainID != models.IrisHub && account.ChainID != models.Kava {
		errors.ErrInvalidChainID(w, http.StatusBadRequest)
		return
	}

	// check validity of an address
	if !strings.Contains(account.Address, sdk.Bech32PrefixAccAddr) || len(account.Address) != 45 {
		errors.ErrInvalidFormat(w, http.StatusBadRequest)
		return
	}

	// check if there is the same account
	// alarm_token with same address
	exist, _ := databases.QueryExistsAccount(w, db, account)
	if exist {
		errors.ErrDuplicateAccount(w, http.StatusConflict)
		return
	}

	// insert account
	databases.InsertAccount(w, db, account)

	u.Result(w, true, "successfully saved")
	return
}

// Update updates the account information
func Update(db *pg.DB, w http.ResponseWriter, r *http.Request) {
	var account models.Account

	// get post data from request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		errors.ErrBadRequest(w, http.StatusBadRequest)
		return
	}

	u.Result(w, true, "successfully updated")
	return
}
