package utils

import (
	"encoding/json"
	"net/http"
)

// Respond responds json format with any data type
func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RespondSuccessMessage(message string) map[string]interface{} {
	return map[string]interface{}{
		"code":   101,
		"result": "success",
		"msg":    message,
	}
}

// func Respond(w http.ResponseWriter, data map[string]interface{}) {
// 	w.Header().Add("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(data)
// }
