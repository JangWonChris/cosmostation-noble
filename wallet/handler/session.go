package handler

import (
	"github.com/cosmostation/cosmostation-cosmos/wallet/config"
	"github.com/cosmostation/cosmostation-cosmos/wallet/db"
)

// Sessions is shorten for s will be used throughout this handler pakcage.
var s *Session

// Session is struct for wrapping both client and db structs.
type Session struct {
	config *config.Config
	db     *db.Database
}

// SetSession set Session object.
func SetSession(config *config.Config, db *db.Database) {
	s = &Session{config, db}
}
