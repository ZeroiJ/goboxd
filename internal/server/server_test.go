package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	h := New(NewStubRunner())
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRunInvalidJSON(t *testing.T) {
	h := New(NewStubRunner())
	req := httptest.NewRequest(http.MethodPost, "/run", bytes.NewBufferString("{bad"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRunUnknownLanguage(t *testing.T) {
	h := New(NewStubRunner())
	body := `{"language":"","source":"print(1)","tests":[{"stdin":""}]}`
	req := httptest.NewRequest(http.MethodPost, "/run", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRunStubAccepted(t *testing.T) {
	h := New(NewStubRunner())
	body := `{"language":"py3","source":"print(1)","tests":[{"stdin":"","expected_stdout":"hi"}]}`
	req := httptest.NewRequest(http.MethodPost, "/run", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
