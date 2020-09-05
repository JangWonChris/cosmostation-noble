package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/wallet/errors"
	"github.com/cosmostation/cosmostation-cosmos/wallet/model"
	"github.com/cosmostation/cosmostation-cosmos/wallet/schema"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

// GetAppVersion returns version number of an app.
func GetAppVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceType := vars["device_type"]

	deviceType = strings.ToLower(deviceType)

	version, err := s.db.QueryAppVersion(deviceType)
	if err != nil {
		zap.L().Error("failed to get app version information", zap.Error(err))
		errors.ErrInternalServer(w, http.StatusInternalServerError)
		return
	}

	// Handle when data is empty
	if version.Version == 0 {
		model.Respond(w, schema.AppVersion{})
		return
	}

	version.IdfVersion = 0

	model.Respond(w, version)
	return
}

// SetAppVersion sets version number of an app.
func SetAppVersion(w http.ResponseWriter, r *http.Request) {
	var version schema.AppVersion

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&version)
	if err != nil {
		errors.ErrBadRequest(w, http.StatusBadRequest)
		return
	}

	appName := strings.ToLower(version.AppName)
	deviceType := strings.ToLower(version.DeviceType)

	// Insert or udpate version information
	exist, _ := s.db.ExistAppVersion(appName, deviceType)
	if !exist {
		s.db.InsertAppVersion(version)
	} else {
		s.db.UpdateAppVersion(version)
	}

	model.Respond(w, version)
	return
}
