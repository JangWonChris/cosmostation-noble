package db

import (
	"fmt"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/wallet/config"
	"github.com/cosmostation/cosmostation-cosmos/wallet/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	// columnLength is the column length of varchar type in every table.
	// This needs to be considered again to set it to what specific length is needed, but right now set it to 99999.
	columnLength = 99999
)

const (
	indexAccountAddress    = "`CREATE INDEX account_address_idx ON account USING btree(address);`"
	indexAccountAlarmToken = "`CREATE INDEX account_alarm_token_idx ON account USING btree(alarm_token);`"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*pg.DB
}

// Connect opens a database connections with the given database connection info from config.
func Connect(cfg config.DBConfig) *Database {
	db := pg.Connect(&pg.Options{
		Addr:         cfg.Host + ":" + cfg.Port,
		User:         cfg.User,
		Password:     cfg.Password,
		Database:     cfg.Table,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// Disable pluralization
	orm.SetTableNameInflector(func(s string) string {
		return s
	})

	return &Database{db}
}

// Ping returns a database connection handle or an error if the connection fails.
func (db *Database) Ping() error {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

// CreateTables creates database tables using object relational mapper (ORM)
func (db *Database) CreateTables() error {
	for _, model := range []interface{}{(*schema.AppAccount)(nil), (*schema.AppVersion)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     columnLength,
		})

		if err != nil {
			return err
		}
	}

	// runs a function in a transaction.
	// if function returns an error transaction is rollbacked, otherwise transaction is committed.
	// create indexes to reduce the cost of lookup queries in case of server traffic jams (B-Tree Index)
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		_, err := db.Model(schema.AppAccount{}).Exec(indexAccountAddress)
		if err != nil {
			return fmt.Errorf("failed to create account address index: %s", err)
		}
		_, err = db.Model(schema.AppAccount{}).Exec(indexAccountAlarmToken)
		if err != nil {
			return fmt.Errorf("failed to create account alarm token index: %s", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// --------------------
// Query
// --------------------

// QueryAppVersion returns mobile app version with given device type.
func (db *Database) QueryAppVersion(deviceType string) (mv schema.AppVersion, err error) {
	err = db.Model(&mv).
		Where("device_type = ?", deviceType).
		Select()

	if err == pg.ErrNoRows {
		return schema.AppVersion{}, nil
	}

	if err != nil {
		return schema.AppVersion{}, err
	}

	return mv, nil
}

// QueryAccount returns user information.
func (db *Database) QueryAccount(address string) (account schema.AppAccount, err error) {
	err = db.Model(&account).
		Where("address = ?", address).
		Select()

	if err == pg.ErrNoRows {
		return schema.AppAccount{}, nil
	}

	if err != nil {
		return schema.AppAccount{}, err
	}

	return account, nil
}

// --------------------
// Insert
// --------------------

// InsertAccount inserts new account information.
func (db *Database) InsertAccount(account schema.AppAccount) error {
	err := db.Insert(&account)
	if err != nil {
		return err
	}

	return nil
}

// InsertVersion inserts new app version
func (db *Database) InsertVersion(version schema.AppVersion) error {
	err := db.Insert(&version)
	if err != nil {
		return err
	}

	return nil
}

// --------------------
// Exist
// --------------------

// ExistsAccount queries to check if the account data exists.
func (db *Database) ExistsAccount(account schema.Account) (bool, error) {
	exist, err := db.Model(&account).
		Where("alarm_token = ? AND address = ?", account.AlarmToken, account.Address).
		Exists()

	if err != nil {
		return exist, err
	}

	return exist, nil
}

// ExistsVersion queries to check if the app data exists
func (db *Database) ExistsVersion(version schema.AppVersion) (bool, error) {
	exist, err := db.Model(&version).
		Where("app_name = ? AND device_type = ?", version.AppName, version.DeviceType).
		Exists()

	if err != nil {
		return exist, err
	}

	return exist, nil
}

// --------------------
// Update
// --------------------

// UpdateVersion updates the app version.
func (db *Database) UpdateVersion(version schema.AppVersion) (schema.AppVersion, error) {
	_, err := db.Model(&version).
		Set("version = ?", version.Version).
		Set("enable = ?", version.Enable).
		Where("app_name = ? AND device_type = ?", version.AppName, version.DeviceType).
		Update()

	if err != nil {
		return schema.AppVersion{}, err
	}

	return version, nil
}

// UpdateAccount updates the account information.
func (db *Database) UpdateAccount(account schema.Account) error {
	_, err := db.Model(&account).
		Set("alarm_status = ?", account.AlarmStatus).
		Where("device_type = ? AND alarm_token = ? AND address = ?", account.DeviceType, account.AlarmToken, account.Address).
		Update()

	if err != nil {
		return err
	}
	return nil
}

// --------------------
// Delete
// --------------------

// DeleteAccount deletes the account
func (db *Database) DeleteAccount(account schema.Account) error {
	_, err := db.Model(&account).
		Where("device_type = ? AND alarm_token = ? AND address = ?", account.DeviceType, account.AlarmToken, account.Address).
		Delete()

	if err != nil {
		return err
	}

	return nil
}
