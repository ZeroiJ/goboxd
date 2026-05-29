package server

import "github.com/thesouldev/goboxd/internal/types"

type stubRunner struct{}

func NewStubRunner() Runner {
	return &stubRunner{}
}

func (s *stubRunner) Run(r *types.RunRequest) (*types.RunResponse, error) {
	resp := &types.RunResponse{
		Status: "accepted",
		Build: types.BuildResult{
			Status:     "ok",
			Stdout:     "",
			Stderr:     "",
			DurationMS: 0,
		},
	}

	resp.Tests = make([]types.TestResult, 0, len(r.Tests))
	for _, tc := range r.Tests {
		resp.Tests = append(resp.Tests, types.TestResult{
			Status:       "accepted",
			Stdout:       tc.ExpectedStdout,
			Stderr:       "",
			DurationMS:   0,
			MemoryPeakKB: 0,
		})
	}

	return resp, nil
}

func (s *stubRunner) Probe() types.Readiness {
	return types.Readiness{Status: "ok"}
}
func (s *stubRunner) Info() types.RunnerInfo {
	return types.RunnerInfo{}
}
