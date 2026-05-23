package types

func ValidateRunRequest(r *RunRequest) *APIError {
	if r.Language == "" {
		return &APIError{Status: 400, Code: "missing_language", Message: "language is required"}
	}
	if r.Source == "" {
		return &APIError{Status: 400, Code: "missing_source", Message: "source is required"}
	}
	if len(r.Tests) == 0 {
		return &APIError{Status: 400, Code: "missing_tests", Message: "tests must have at least one entry"}
	}
	return nil
}
