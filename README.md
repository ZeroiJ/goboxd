# goboxd

goboxd is a Go HTTP service that compiles and runs untrusted code inside an nsjail sandbox and returns per-test results.

## Run it

Prereqs: Docker with Compose v2.

  make build
  make run

## Test it

  make test
  make integration
  make lint

## API

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

## Docs

- docs/ for API, languages, security, benchmarks, architecture
- Spec: https://intern-iitm.github.io/goboxd-hackathon/spec.html
- Discussions: https://github.com/intern-iitm/goboxd-hackathon/discussions
- Submission repo: https://github.com/thesouldev/goboxd
