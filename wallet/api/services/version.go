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
	"github.com/gorilla/mux"
)

// GetVersion returns version number of an app
func GetVersion(db *pg.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceType := vars["deviceType"]

	// lower case
	deviceType = strings.ToLower(deviceType)

	var version models.AppVersion

	switch deviceType {
	case models.Android:
		version, _ = databases.QueryAppVersion(w, db, deviceType)
	case models.IOS:
		version, _ = databases.QueryAppVersion(w, db, deviceType)
	default:
		errors.ErrInvalidDeviceType(w, http.StatusBadRequest)
		return
	}

	// in case when data is empty
	if version.Latest == 0 {
		return
	}

	u.Respond(w, version)
	return
}

// SetVersion sets version number of an app
func SetVersion(db *pg.DB, w http.ResponseWriter, r *http.Request) {
	var version models.AppVersion

	// get post data from request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&version)
	if err != nil {
		errors.ErrBadRequest(w, http.StatusBadRequest)
		return
	}

	// lower case
	version.AppName = strings.ToLower(version.AppName)
	version.DeviceType = strings.ToLower(version.DeviceType)

	exist, _ := databases.QueryExistsAppVersion(w, db, version)

	// update version info if it exists. Otherwise, insert version info
	if exist {
		databases.UpdateAppVersion(w, db, version)
	} else {
		databases.InsertAppVersion(w, db, version)
	}

	u.Respond(w, version)
	return
}
