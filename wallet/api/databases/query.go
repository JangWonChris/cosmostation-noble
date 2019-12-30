package databases

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/wallet/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/models"
	"github.com/go-pg/pg"
)

// InsertAccount inserts new account information
func InsertAccount(w http.ResponseWriter, db *pg.DB, account models.Account) (uint64, models.Account, error) {
	err := db.Insert(&account)
	if err != nil {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
		return 0, account, err
	}
	return 1, account, nil
}

// InsertAppVersion inserts new app version
func InsertAppVersion(w http.ResponseWriter, db *pg.DB, version models.AppVersion) (uint64, models.AppVersion, error) {
	err := db.Insert(&version)
	if err != nil {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
		return 0, version, err
	}
	return 1, version, nil
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

// QueryAccount queries account information
func QueryAccount(w http.ResponseWriter, db *pg.DB, address string) (models.Account, error) {
	var account models.Account
	err := db.Model(&account).
		Where("address = ?", address).
		Select()
	if err != nil {
		errors.ErrNotFound(w, http.StatusNotFound)
	}
	return account, nil
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

// UpdateAccount updates the account
func UpdateAccount(w http.ResponseWriter, db *pg.DB, account models.Account) (bool, error) {
	_, err := db.Model(&account).
		Set("alarm_status = ?", account.AlarmStatus).
		Where("device_type = ? AND alarm_token = ? AND address = ?", account.DeviceType, account.AlarmToken, account.Address).
		Update()
	if err != nil {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
		return false, err
	}
	return true, nil
}

// DeleteAccount deletes the account
func DeleteAccount(w http.ResponseWriter, db *pg.DB, account models.Account) (bool, error) {
	_, err := db.Model(&account).
		Where("device_type = ? AND alarm_token = ? AND address = ?", account.DeviceType, account.AlarmToken, account.Address).
		Delete()
	if err != nil {
		errors.ErrInternalServer(w, http.StatusInternalServerError)
	}
	return true, nil
}
