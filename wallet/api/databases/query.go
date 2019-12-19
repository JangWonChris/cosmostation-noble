package databases

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/wallet/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/models"
	"github.com/go-pg/pg"
)

// InsertAccount inserts new account information
func InsertAccount(w http.ResponseWriter, db *pg.DB, account models.Account) (models.Account, error) {
	err := db.Insert(&account)
	if err != nil {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
	}
	return account, nil
}

// InsertAppVersion inserts new app version
func InsertAppVersion(w http.ResponseWriter, db *pg.DB, version models.AppVersion) (models.AppVersion, error) {
	err := db.Insert(&version)
	if err != nil {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
	}
	return version, nil
}

// QueryAppVersion queries mobile app version
func QueryAppVersion(w http.ResponseWriter, db *pg.DB, deviceType string) (models.AppVersion, error) {
	var version models.AppVersion
	err := db.Model(&version).
		Where("device_type = ?", deviceType).
		Select()
	if err != nil {
		errors.ErrNotFound(w, http.StatusNotFound)
	}
	return version, nil
}

// QueryExistsAppVersion queries to check if the app data already exists
func QueryExistsAppVersion(w http.ResponseWriter, db *pg.DB, version models.AppVersion) (bool, error) {
	exist, err := db.Model(&version).
		Where("app_name = ? AND device_type = ?", version.AppName, version.DeviceType).
		Exists()
	if err != nil {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
	}
	return exist, nil
}

// QueryExistsAccount queries to check if the account data already exists
func QueryExistsAccount(w http.ResponseWriter, db *pg.DB, account models.Account) (bool, error) {
	exist, err := db.Model(&account).
		Where("alarm_token = ? AND address = ?", account.AlarmToken, account.Address).
		Exists()
	if err != nil {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
	}
	return exist, nil
}

// UpdateAppVersion updates the app version
func UpdateAppVersion(w http.ResponseWriter, db *pg.DB, version models.AppVersion) (models.AppVersion, error) {
	_, err := db.Model(&version).
		Set("acceptable = ?", version.Acceptable).
		Set("latest = ?", version.Latest).
		Where("app_name = ? AND device_type = ?", version.AppName, version.DeviceType).
		Update()
	if err != nil {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
	}
	return version, nil
}

// UpdateAccount updates the account information
func UpdateAccount(w http.ResponseWriter, db *pg.DB, account models.Account) (models.Account, error) {
	_, err := db.Model(&account).
		Set("alarm_status = ?", account.AlarmStatus).
		Where("device_type = ? alarm_token = ? AND address = ?", account.DeviceType, account.AlarmToken, account.Address).
		Update()
	if err != nil {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
	}
	return account, nil
}
