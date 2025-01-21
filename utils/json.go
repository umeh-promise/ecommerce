package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "applications/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func ParseJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1_048_576 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteJSONError(w http.ResponseWriter, status int, error []string, message string) error {
	type envelop struct {
		Status  string   `json:"status"`
		Error   []string `json:"error"`
		Message string   `json:"message"`
	}

	return WriteJSON(w, status, &envelop{
		Status:  "error",
		Error:   error,
		Message: message,
	})
}

func JSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	type envelop struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}

	return WriteJSON(w, status, &envelop{
		Status: "success",
		Data:   data,
	})
}
