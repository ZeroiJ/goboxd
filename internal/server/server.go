package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thesouldev/goboxd/internal/types"
)

type Runner interface {
	Run(r *types.RunRequest) (*types.RunResponse, error)
}

func New(runner Runner) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, &types.APIError{Status: http.StatusMethodNotAllowed, Code: "method_not_allowed", Message: "POST required"})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB limit
		var req types.RunRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) || err.Error() == "http: request body too large" {
				writeError(w, &types.APIError{Status: http.StatusRequestEntityTooLarge, Code: "request_too_large", Message: "request body exceeds 1MB limit"})
				return
			}
			writeError(w, &types.APIError{Status: http.StatusBadRequest, Code: "invalid_json", Message: "invalid JSON body"})
			return
		}

		if err := types.ValidateRunRequest(&req); err != nil {
			writeError(w, err)
			return
		}

		resp, err := runner.Run(&req)
		if err != nil {
			var ae *types.APIError
			if errors.As(err, &ae) {
				writeError(w, ae)
				return
			}
			writeError(w, &types.APIError{Status: http.StatusInternalServerError, Code: "internal_error", Message: "internal server error"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})

	return mux
}

func writeError(w http.ResponseWriter, err *types.APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	_ = json.NewEncoder(w).Encode(types.ErrorResponse{
		Error: types.ErrorBody{Code: err.Code, Message: err.Message},
	})
}
