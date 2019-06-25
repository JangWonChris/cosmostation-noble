package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app/exception"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app/models"
)

func UpdateAlarmStatus(DB *pg.DB, w http.ResponseWriter, r *http.Request) {
	// Get post data from request
	var account models.Account
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check the validity of cosmos address
	if !strings.Contains(account.Address, sdk.Bech32PrefixAccAddr) || len(account.Address) != 45 {
		exception.ErrInvalidFormat(w, http.StatusBadRequest)
		return
	}

	// Check if same account already exists (alarm_token with same address)
	exist, err := DB.Model(&account).
		Where("alarm_token = ? AND address = ?", account.AlarmToken, account.Address).
		Exists()
	if !exist {
		exception.ErrNotExist(w, http.StatusBadRequest)
		return
	}

	// Update alarm status
	err = DB.Insert(&account)
	fmt.Println(err)
	if err != nil {
		exception.ErrInternalServer(w, http.StatusInternalServerError)
		return
	}

	// resp := models.Message(101, "UpdateAlarm")
	// models.Respond(w, resp)
}
