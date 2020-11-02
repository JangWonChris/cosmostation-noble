package handler

import (
	"net/http"

	"github.com/tomasen/realip"

	"go.uber.org/zap"
)

// TODO: Response Status Code needs to be logged. Find out how others do to resolve this, maybe third party library?
// Should we use key:value pairs for better json formatted output?

// Middleware logs incoming requests and calls next handler.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		clientIP := realip.FromRequest(r)
		zap.S().Infof("%s %s [%s]", r.Method, r.URL, clientIP)

		next.ServeHTTP(rw, r)
	})
}
