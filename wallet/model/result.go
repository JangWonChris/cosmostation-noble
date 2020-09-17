package model

import (
	"encoding/json"
	"net/http"
)

// Respond responds json format with any data type
func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ResultMessage defines the structure for result message.
type ResultMessage struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
}

// Result returns ResultMessage.
func Result(w http.ResponseWriter, result bool, msg string) {
	w.Header().Add("Content-Type", "application/json")
	data := &ResultMessage{
		Result: result,
		Msg:    msg,
	}
	json.NewEncoder(w).Encode(data)
}
