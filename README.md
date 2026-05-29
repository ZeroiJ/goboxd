# goboxd

goboxd is a Go HTTP service that compiles and runs untrusted code inside an nsjail sandbox. It evaluates submissions against test cases and returns structured results.

We use the standard library `net/http` for our routing framework. The API surface is small enough that an external framework is unnecessary, and the standard `ServeMux` keeps the dependency tree minimal.

## Running

Requires Docker and Docker Compose v2. No local nsjail installation is needed; it compiles from source during the container build.

    make build
    make run

## Testing

    make test         # Unit tests
    make integration  # End-to-end tests inside Docker
    make lint         # golangci-lint

## Documentation

Extended documentation is located in the `docs/` directory:
- `docs/api.md`: Endpoint specifications
- `docs/languages.md`: Language registry and YAML configuration
- `docs/security.md`: Hardening and vulnerability fixes
- `docs/architecture.md`: System design and concurrency queue
- `docs/benchmarks.md`: Load testing results

## Submission Info

- Spec: https://intern-iitm.github.io/goboxd-hackathon/spec.html
- Submission repo: https://github.com/thesouldev/goboxd
