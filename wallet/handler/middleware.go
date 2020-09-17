package handler

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/wallet/config"
	"github.com/cosmostation/cosmostation-cosmos/wallet/db"

	"github.com/tomasen/realip"

	"go.uber.org/zap"
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

// Middleware logs incoming requests and calls next handler.
func Middleware(next http.Handler, config *config.Config, db *db.Database) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		clientIP := realip.FromRequest(r)
		zap.S().Infof("%s %s [%s]", r.Method, r.URL, clientIP)

		// Session will wrap both client and database and be used throughout all handlers.
		SetSession(config, db)

		next.ServeHTTP(rw, r)
	})
}
