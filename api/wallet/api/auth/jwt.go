package auth

// import (
// 	"net/http"
// 	"strings"
// 	"time"

// 	u "github.com/cosmostation/cosmostation-cosmos/api/wallet/api/utils"
// 	"github.com/cosmostation/cosmostation-cosmos/api/wallet/config"

// 	jwt "github.com/dgrijalva/jwt-go"
// )

// // Authenticates JSON Web Token
// var JwtAuthentication = func(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Requested endpoint
// 		requestedPath := r.URL.Path

// 		// Endpoints that are required for JWT authentication
// 		requiredAuth := []string{"/v1/auth/account/testv"}

// 		// Check if the endpoint needs an authentication
// 		for _, value := range requiredAuth {
// 			if value == requestedPath {
// 				w.Header().Set("Content-Type", "application/json")

// 				resp := make(map[string]interface{})

// 				// X-Auth-Token is missing, returns with error code 403 Unauthorized
// 				tokenHeader := r.Header.Get("X-Auth-Token") // Grab the token from the header
// 				if tokenHeader == "" {
// 					resp = u.Message(false, "missing auth token")
// 					w.WriteHeader(http.StatusForbidden)
// 					u.Respond(w, resp)
// 					return
// 				}

// 				// Check Json Web Token format
// 				splitted := strings.Split(tokenHeader, ".")
// 				if len(splitted) != 3 {
// 					resp = u.Message(false, "invalid or malformed auth token")
// 					w.WriteHeader(http.StatusForbidden)
// 					u.Respond(w, resp)
// 					return
// 				}

// 				// Malformed token, returns with http code 403
// 				tk := &u.Token{}
// 				token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
// 					return []byte(config.GetMainnetDevConfig().JWT.Token), nil
// 				})
// 				if err != nil {
// 					resp = u.Message(false, "malformed authentication token")
// 					w.WriteHeader(http.StatusForbidden)
// 					u.Respond(w, resp)
// 					return
// 				}

// 				// Token is invalid, maybe not signed on this server
// 				if !token.Valid {
// 					resp = u.Message(false, "token is invalid")
// 					w.WriteHeader(http.StatusForbidden)
// 					u.Respond(w, resp)
// 					return
// 				}

// 				// Check expiration time
// 				currentTime := time.Now().UTC()
// 				expiredTime := tk.ExpiredAt.UTC()
// 				diff := expiredTime.Sub(currentTime)
// 				if diff < 0 {
// 					w.WriteHeader(http.StatusForbidden)
// 					resp = u.Message(false, "token is expired")
// 					u.Respond(w, resp)
// 					return
// 				}

// 				// Proceed with the request and set the caller to the user retrieved from the parsed token
// 				next.ServeHTTP(w, r) //proceed in the middleware chain!
// 				return
// 			}
// 		}
// 		// Process if requestedPath is not required to authorize
// 		next.ServeHTTP(w, r)
// 		return
// 	})
// }
