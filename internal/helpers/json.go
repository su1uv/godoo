package helpers

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	const maxBytes = 1 << 20

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	err := json.UnmarshalRead(r.Body, dst, json.RejectUnknownMembers(true))
	if err != nil {
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not exceed %d bytes", maxBytes)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		default:
			return fmt.Errorf("invalid JSON request body: %w", err)
		}
	}

	return nil
}

func RespondWithJSON(w http.ResponseWriter, code int, data any) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(res)
	return err
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
}
