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

func TestValidateRunRequestTooManyTests(t *testing.T) {
	r := &RunRequest{Language: "py3", Source: "print(1)"}
	for i := 0; i < MaxTests+1; i++ {
		r.Tests = append(r.Tests, TestCase{Stdin: ""})
	}
	if ValidateRunRequest(r) == nil {
		t.Fatal("expected too_many_tests error")
	}
}

func TestValidateFilename(t *testing.T) {
	good := &RunRequest{Language: "py3", Source: "print(1)", SourceFilename: "main.py", Tests: []TestCase{{Stdin: ""}}}
	if ValidateRunRequest(good) != nil {
		t.Fatal("expected no validation error")
	}

	bad := &RunRequest{Language: "py3", Source: "print(1)", SourceFilename: "../x.py", Tests: []TestCase{{Stdin: ""}}}
	if ValidateRunRequest(bad) == nil {
		t.Fatal("expected invalid_filename error")
	}
}
