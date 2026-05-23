package types

import "testing"

func TestValidateRunRequest(t *testing.T) {
	bad := &RunRequest{}
	if ValidateRunRequest(bad) == nil {
		t.Fatal("expected validation error")
	}

	ok := &RunRequest{Language: "py3", Source: "print(1)", Tests: []TestCase{{Stdin: ""}}}
	if ValidateRunRequest(ok) != nil {
		t.Fatal("expected no validation error")
	}
}
