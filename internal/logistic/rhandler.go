// Package ragger defines server work and functions to configure server.
package logistic

import (
	"encoding/json"
	"net/http"
)

type raggerHandler func(w http.ResponseWriter, r *http.Request) error

func handle(rh raggerHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := rh(w, r); err != nil {
			switch err.(type) {
			case Error:
				RespondError(w, err.(Error))
			default:
				httpErr := Error{
					Code: http.StatusInternalServerError,
					Err:  err,
				}

				RespondError(w, httpErr)
			}
		}
	})
}

// Error struct for return api error
type Error struct {
	Code int
	Err  error
}

func newError(code int, err error) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func (e Error) Error() string {
	return e.Err.Error()
}

// MarshalJSON implements Marshaller interface.
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Code int    `json:"code"`
		Err  string `json:"error"`
	}{
		Code: e.Code,
		Err:  e.Err.Error(),
	})
}

// Message returns map with data for answer on API request
func Message(message interface{}) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

// Respond does answer with data on API request
func Respond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RespondError does answer with error on API request
func RespondError(w http.ResponseWriter, err Error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	if err := json.NewEncoder(w).Encode(err); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
