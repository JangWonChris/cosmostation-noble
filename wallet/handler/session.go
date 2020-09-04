package handler

import (
	"github.com/cosmostation/cosmostation-cosmos/wallet/db"
)

// Sessions is shorten for s will be used throughout this handler pakcage.
var s *Session

// Session is struct for wrapping both client and db structs.
type Session struct {
	db *db.Database
}

// SetSession set Session object.
func SetSession(db *db.Database) {
	s = &Session{db}
}
