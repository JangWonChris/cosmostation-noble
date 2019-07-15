package utils

// import (
// 	"fmt"
// 	"time"

// 	"github.com/cosmostation/cosmostation-cosmos/api/wallet/config"

// 	jwt "github.com/dgrijalva/jwt-go"
// )

// var (
// 	expiration = 1000000
// )

// // JWT claims struct
// type Token struct {
// 	ID        uint
// 	IssuedAt  time.Time `json:"issued_at"`
// 	ExpiredAt time.Time `json:"expired_at"`
// 	jwt.StandardClaims
// }

// // Creates expirable token
// func CreateToken(id uint) string {
// 	currentTime := time.Now().UTC()
// 	expiredAt := currentTime.Add(time.Hour * time.Duration(expiration))

// 	// Insert token information
// 	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{
// 		ID:        id,
// 		IssuedAt:  currentTime,
// 		ExpiredAt: expiredAt,
// 	})

// 	tokenString, err := token.SignedString([]byte(config.GetMainnetDevConfig().JWT.Token))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return tokenString
// }
