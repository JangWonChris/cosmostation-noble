package utils

import (
	"encoding/json"
	"net/http"
)

/*
	어떠한 타입의 struct가 와도 json 포맷으로 리턴할 수 있게끔 구현
*/

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
