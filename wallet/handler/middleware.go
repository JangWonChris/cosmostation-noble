package handler

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/wallet/client"
	"github.com/cosmostation/cosmostation-cosmos/wallet/db"

	"github.com/tomasen/realip"

	"go.uber.org/zap"
)

// Middleware logs incoming requests and calls next handler.
func Middleware(next http.Handler, c *client.Client, db *db.Database) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		clientIP := realip.FromRequest(r)
		zap.S().Infof("%s %s [%s]", r.Method, r.URL, clientIP)

		// Session will wrap both client and database and be used throughout all handlers.
		SetSession(c, db)

		next.ServeHTTP(rw, r)
	})
}
