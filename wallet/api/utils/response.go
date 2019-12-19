package utils

import (
	"encoding/json"
	"net/http"
)

// type ResponseMsg struct {
// 	Code   uint16      `json:"code"`
// 	Result string      `json:"result"`
// 	Msg    interface{} `json:"msg"`
// }

// SuccessResult responds map string format with any data type
func SuccessResult(w http.ResponseWriter, msg string) {
	w.Header().Add("Content-Type", "application/json")
	result := make(map[string]string)
	result["result"] = msg
	json.NewEncoder(w).Encode(result)
}

// Respond responds json format with any data type
func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// func Respond(w http.ResponseWriter, data interface{}) {
// 	w.Header().Add("Content-Type", "application/json")
// 	resp := &ResponseMsg{
// 		Code:   101,
// 		Result: Success,
// 		Msg:    data,
// 	}
// 	json.NewEncoder(w).Encode(resp)
// }
