package types

const (
	MaxTests = 50
)

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
	if len(r.Tests) > MaxTests {
		return &APIError{Status: 400, Code: "too_many_tests", Message: "tests exceeds max limit"}
	}
	if r.SourceFilename != "" {
		if err := validateFilename(r.SourceFilename); err != nil {
			return err
		}
	}
	if r.ArtifactFilename != "" {
		if err := validateFilename(r.ArtifactFilename); err != nil {
			return err
		}
	}
	return nil
}

func validateFilename(name string) *APIError {
	if name == "." || name == ".." {
		return &APIError{Status: 400, Code: "invalid_filename", Message: "filename must be a single path component"}
	}
	if name == "" {
		return &APIError{Status: 400, Code: "invalid_filename", Message: "filename must be a single path component"}
	}
	if name[0] == '.' {
		return &APIError{Status: 400, Code: "invalid_filename", Message: "filename must be a single path component"}
	}
	for _, ch := range name {
		if ch == '/' || ch == '\\' {
			return &APIError{Status: 400, Code: "invalid_filename", Message: "filename must be a single path component"}
		}
	}
	return nil
}
