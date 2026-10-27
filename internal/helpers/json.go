package helpers

import (
	"encoding/json/v2"
	"net/http"
)

func Encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.MarshalWrite(w, v); err != nil {
		return err
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T

	if err := json.UnmarshalRead(r.Body, &v); err != nil {
		return v, err
	}

	return v, nil
}
