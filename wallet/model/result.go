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

type ResultMessage struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
}

// Result returns result of the message
func Result(w http.ResponseWriter, result bool, msg string) {
	w.Header().Add("Content-Type", "application/json")
	data := &ResultMessage{
		Result: result,
		Msg:    msg,
	}
	json.NewEncoder(w).Encode(data)
}
