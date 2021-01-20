package handler

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/client"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
)

// Sessions is shorten for s will be used throughout this handler pakcage.
var s *Session

// Session is struct for wrapping both client and db structs.
type Session struct {
	Client *client.Client
	DB     *db.Database
}

// SetSession set Session object.
func SetSession(client *client.Client, db *db.Database) *Session {
	return &Session{client, db}
}
