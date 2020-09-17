package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/wallet/errors"
	"github.com/cosmostation/cosmostation-cosmos/wallet/model"
)

// SignSignature signs signature with the moonpay's secret key and return output.
func SignSignature(w http.ResponseWriter, r *http.Request) {
	var mp model.MoonPay
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		fmt.Printf("failed to decode MoonPay: %t\n", err)
		errors.ErrBadRequest(w, http.StatusBadRequest)
		return
	}

	// Create a new HMAC by defining the hash type and the key as byte array and encode as hexadecimal string.
	hash := hmac.New(sha256.New, []byte(s.config.Payment.MoonPaySecretKey))
	hash.Write([]byte(mp.APIKey))

	// Convert to base64 encoding format.
	sig := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	result := &model.ResultMoonPay{
		Signature: sig,
	}

	model.Respond(w, result)
	return
}
