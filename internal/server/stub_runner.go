package server

import "github.com/thesouldev/goboxd/internal/types"

type stubRunner struct{}

func NewStubRunner() Runner {
	return &stubRunner{}
}

func (s *stubRunner) Run(r *types.RunRequest) (*types.RunResponse, error) {
	return nil, &types.APIError{Status: 501, Code: "not_implemented", Message: "runner not implemented yet"}
}
