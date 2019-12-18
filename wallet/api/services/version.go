package services

import (
	"fmt"
	"net/http"

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
	case "android":
		_ = DB.Model(&version).
			Where("device_type = ?", deviceType).
			Select()
	case "ios":
		_ = DB.Model(&version).
			Where("device_type = ?", deviceType).
			Select()
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if version.Latest == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	u.Respond(w, version)
	return
}

// SetVersion sets version number of an app
func SetVersion(DB *pg.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("SetVersion")
}
