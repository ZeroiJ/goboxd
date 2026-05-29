//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/thesouldev/goboxd/internal/types"
)

const baseURL = "http://goboxd:8080"

func resolveBaseURL() string {
	if os.Getenv("GOBXD_BASE_URL") != "" {
		return os.Getenv("GOBXD_BASE_URL")
	}
	return baseURL
}

func postRun(t *testing.T, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	url := resolveBaseURL() + "/run"
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func TestIntegrationPythonAccepted(t *testing.T) {
	req := types.RunRequest{
		Language: "py3",
		Source:   `print("hello world")`,
		Tests: []types.TestCase{
			{Stdin: "", ExpectedStdout: "hello world"},
		},
	}
	resp := postRun(t, req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result types.RunResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "accepted" {
		t.Fatalf("expected accepted, got %s", result.Status)
	}
	if len(result.Tests) != 1 {
		t.Fatalf("expected 1 test, got %d", len(result.Tests))
	}
	if result.Tests[0].Status != "accepted" {
		t.Fatalf("expected test accepted, got %s", result.Tests[0].Status)
	}
}

func TestIntegrationCppAccepted(t *testing.T) {
	req := types.RunRequest{
		Language: "cpp",
		Source: `#include <iostream>
int main() {
    std::cout << "hello world" << std::endl;
    return 0;
}`,
		Tests: []types.TestCase{
			{Stdin: "", ExpectedStdout: "hello world"},
		},
	}
	resp := postRun(t, req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result types.RunResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "accepted" {
		t.Fatalf("expected accepted, got %s", result.Status)
	}
	if result.Build.Status != "ok" {
		t.Fatalf("expected build success, got %s", result.Build.Status)
	}
}

func TestIntegrationPythonWrongAnswer(t *testing.T) {
	req := types.RunRequest{
		Language: "py3",
		Source:   `print("hello world")`,
		Tests: []types.TestCase{
			{Stdin: "", ExpectedStdout: "goodbye"},
		},
	}
	resp := postRun(t, req)
	defer resp.Body.Close()

	var result types.RunResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "wrong_output" {
		t.Fatalf("expected wrong_answer, got %s", result.Status)
	}
}
