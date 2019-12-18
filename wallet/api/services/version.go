package services

import (
	"encoding/json"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/wallet/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/models"
	u "github.com/cosmostation/cosmostation-cosmos/wallet/api/utils"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// GetVersion returns version number of an app
func GetVersion(DB *pg.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceType := vars["deviceType"]

	var version models.AppVersion

	switch deviceType {
	case models.Android:
		_ = DB.Model(&version).
			Where("device_type = ?", deviceType).
			Select()
	case models.IOS:
		_ = DB.Model(&version).
			Where("device_type = ?", deviceType).
			Select()
	default:
		errors.ErrInvalidDeviceType(w, http.StatusBadRequest)
		return
	}

	if version.Latest == 0 {
		errors.ErrNotFound(w, http.StatusNotFound)
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

	exist, err := db.Model(&version).
		Where("app_name = ? AND device_type = ?", version.AppName, version.DeviceType).
		Count()

	// update version info if it exists. Otherwise, insert version info
	if exist > 0 {
		_, err = db.Model(&version).
			Set("acceptable = ?", version.Acceptable).
			Set("latest = ?", version.Latest).
			Where("app_name = ? AND device_type = ?", version.AppName, version.DeviceType).
			Update()
		if err != nil {
			errors.ErrInternalServer(w, http.StatusInternalServerError)
			return
		}
	} else {
		err = db.Insert(&version)
		if err != nil {
			errors.ErrInternalServer(w, http.StatusInternalServerError)
			return
		}
	}

}
