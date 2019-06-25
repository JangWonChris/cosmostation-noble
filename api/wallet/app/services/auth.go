package services

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app/exception"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app/models"
	u "github.com/cosmostation/cosmostation-cosmos/api/wallet/app/utils"

	"github.com/go-pg/pg"
)

func Test(DB *pg.DB, w http.ResponseWriter, r *http.Request) error {
	status := false
	if status {
		exception.ErrDuplicateAccount(w, http.StatusUnauthorized)
		return nil
	}

	var account []models.Account
	tempAccount := models.Account{
		AlarmToken: "AAABBBCCC",
		DeviceType: "android",
		CoinType:   "ATOM",
		Status:     true,
	}
	account = append(account, tempAccount)

	resp := models.AccountResponse{
		Result: "success",
		Data:   account,
	}
	u.Respond(w, resp)
	return nil
}
