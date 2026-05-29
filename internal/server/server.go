package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/thesouldev/goboxd/internal/types"
)

type Runner interface {
	Run(r *types.RunRequest) (*types.RunResponse, error)
}

// Server encapsulates the HTTP handlers, runner, and global state
type Server struct {
	Runner        Runner
	Sem           chan struct{}
	MaxJobs       int
	InFlight      atomic.Int64
	JobsTotal     atomic.Int64
	JobsFailed    atomic.Int64
	LastErrorTime atomic.Pointer[time.Time]
}

func New(runner Runner) *Server {
	maxJobs := runtime.NumCPU()
	if env := os.Getenv("MAX_CONCURRENT_JOBS"); env != "" {
		if v, err := strconv.Atoi(env); err == nil && v > 0 {
			maxJobs = v
		}
	}
	return &Server{
		Runner:  runner,
		Sem:     make(chan struct{}, maxJobs),
		MaxJobs: maxJobs,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/run", s.handleRun)
	return mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
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

	s.Sem <- struct{}{} // wait for concurrency slot
	s.InFlight.Add(1)
	s.JobsTotal.Add(1)
	defer func() {
		s.InFlight.Add(-1)
		<-s.Sem // release concurrency slot
	}()

	resp, err := s.Runner.Run(&req)
	if err != nil {
		s.JobsFailed.Add(1)
		now := time.Now()
		s.LastErrorTime.Store(&now)

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
}

func writeError(w http.ResponseWriter, err *types.APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	_ = json.NewEncoder(w).Encode(types.ErrorResponse{
		Error: types.ErrorBody{Code: err.Code, Message: err.Message},
	})
}
