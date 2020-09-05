package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/wallet/errors"
	"github.com/cosmostation/cosmostation-cosmos/wallet/model"
	"github.com/cosmostation/cosmostation-cosmos/wallet/schema"
	"go.uber.org/zap"
)

// RegisterOrUpdateAccount registers an account if it doesn't exist or
// updates the account information if it exists.
func RegisterOrUpdateAccount(w http.ResponseWriter, r *http.Request) {
	var account schema.AppAccount

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		errors.ErrBadRequest(w, http.StatusBadRequest)
		return
	}

	accAddr := strings.ToLower(account.Address)
	deviceType := strings.ToLower(account.DeviceType)

	if deviceType != model.Android && deviceType != model.IOS {
		errors.ErrInvalidDeviceType(w, http.StatusBadRequest)
		return
	}

	if account.ChainID != model.Cosmos &&
		account.ChainID != model.Iris &&
		account.ChainID != model.Kava {
		errors.ErrInvalidChainID(w, http.StatusBadRequest)
		return
	}

	err = model.VerifyBech32AccAddr(accAddr)
	if err != nil {
		zap.L().Debug("failed to verify account address", zap.Error(err))
		errors.ErrInvalidParam(w, http.StatusBadRequest, "account address is invalid")
		return
	}

	// Insert account information if it doesn't exist
	exist, _ := s.db.ExistAppAccount(account.AlarmToken, account.Address)
	if !exist {
		err := s.db.InsertAppAccount(account)
		if err != nil {
			zap.L().Error("failed to insert app account information", zap.Error(err))
			errors.ErrInternalServer(w, http.StatusInternalServerError)
			return
		}

		model.Result(w, true, "successfully inserted")
		return
	}

	err = s.db.UpdateAppAccount(account)
	if err != nil {
		zap.L().Error("failed to insert app account information", zap.Error(err))
		errors.ErrInternalServer(w, http.StatusInternalServerError)
		return
	}

	model.Result(w, true, "successfully updated")
	return
}

// DeleteAccount deletes the account information.
// [NOT USED]
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var account schema.AppAccount

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		errors.ErrBadRequest(w, http.StatusBadRequest)
		return
	}

	exist, _ := s.db.ExistAppAccount(account.AlarmToken, account.Address)
	if !exist {
		errors.ErrNotFound(w, http.StatusNotFound)
		return
	}

	err = s.db.DeleteAppAccount(account)
	if err != nil {
		zap.L().Error("failed to delete app account information", zap.Error(err))
		errors.ErrInternalServer(w, http.StatusInternalServerError)
		return
	}

	model.Result(w, true, "successfully deleted")
	return
}
