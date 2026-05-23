# goboxd

A Go HTTP service for executing untrusted code in isolated sandboxes.

## Important (read this first)

- Stage 1 target: a Go HTTP server that runs in Docker, with /healthz and POST /run working for 2 languages (one interpreted + one compiled).
- Use the branch: team/ZeroTrust_Hustlers.
- The submission repo is thesouldev/goboxd; work on your fork and open a PR back to master.
- Everything should run inside Docker (use make build / make run / make test).

## What this service does

- Accepts code + test cases over HTTP
- Runs the code inside a sandbox (nsjail)
- Returns per-test results as JSON

## Project structure

.
├── cmd/goboxd/   binary entry point
├── internal/     private application packages
├── docs/         api, languages, security, benchmarks, architecture
└── tests/        integration tests

## Getting started

Prerequisites:
- Docker with Compose v2

Clone and build:

  git clone https://github.com/thesouldev/goboxd.git
  cd goboxd
  make build

Run locally:

  make run

Tests:

  make test
  make integration
  make lint

## API quickstart

Health check:

  curl -s http://localhost:8080/healthz

Run code (example):

  curl -s http://localhost:8080/run \
    -H 'Content-Type: application/json' \
    -d '{
      "language": "py3",
      "source": "print(\"hi\")",
      "tests": [
        { "stdin": "", "expected_stdout": "hi" }
      ]
    }'

## Notes

- See the hackathon spec for the full API contract and stage requirements.
- Stage 1 only needs two languages working end-to-end.
