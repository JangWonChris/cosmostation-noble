package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorCode type
type ErrorCode uint32

// ErrorMsg type
type ErrorMsg string

// Parses the error into an object-like struct for exporting
type WrapError struct {
	ErrorCode ErrorCode `json:"error_code"`
	ErrorMsg  ErrorMsg  `json:"error_msg"`
}

// Error code numbers
const (
	InternalServer ErrorCode = 101

	DuplicateAccount ErrorCode = 201
	InvalidFormat    ErrorCode = 202
	NotExist         ErrorCode = 203
)

// ErrorCodeToErrorMsg returns error message from error code
func ErrorCodeToErrorMsg(code ErrorCode) ErrorMsg {
	switch code {
	case InternalServer:
		return "Internal server error"
	case DuplicateAccount:
		return "Duplicate account"
	case InvalidFormat:
		return "Invalid format"
	case NotExist:
		return "NotExist"
	default:
		return "Unknown"
	}
}

/*
	----------------------------------------------- Error Types
*/

func ErrInternalServer(w http.ResponseWriter, statusCode int) {
	wrapError := WrapError{
		ErrorCode: InternalServer,
		ErrorMsg:  ErrorCodeToErrorMsg(InternalServer),
	}
	PrintException(w, statusCode, wrapError)
}

func ErrDuplicateAccount(w http.ResponseWriter, statusCode int) {
	wrapError := WrapError{
		ErrorCode: DuplicateAccount,
		ErrorMsg:  ErrorCodeToErrorMsg(DuplicateAccount),
	}
	PrintException(w, statusCode, wrapError)
}

func ErrInvalidFormat(w http.ResponseWriter, statusCode int) {
	wrapError := WrapError{
		ErrorCode: InvalidFormat,
		ErrorMsg:  ErrorCodeToErrorMsg(InvalidFormat),
	}
	PrintException(w, statusCode, wrapError)
}

func ErrNotExist(w http.ResponseWriter, statusCode int) {
	wrapError := WrapError{
		ErrorCode: NotExist,
		ErrorMsg:  ErrorCodeToErrorMsg(NotExist),
	}
	PrintException(w, statusCode, wrapError)
}

/*
	----------------------------------------------- PrintException
*/

// PrintException prints out the exception result
func PrintException(w http.ResponseWriter, statusCode int, err WrapError) {
	w.Header().Add("Content-Type", "application/json")

	// Write HTTP status code
	w.WriteHeader(statusCode)

	result, _ := json.Marshal(err)

	fmt.Fprint(w, string(result))
}
