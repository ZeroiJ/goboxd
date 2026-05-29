# Changelog

This file tracks work done by each teammate.

## Unreleased

### Sujal
- Added HTTP server skeleton with /healthz and /run handlers
- Added request/response types and validation errors
- Added stub runner and initial unit tests
- Fixed variable shadowing bug in /run handler
- Added validation limits (max tests) and filename safety checks
- Added tests for validation limits and filename checks
- Overhauled README with stage-1 guidance and a curl example
- Stub runner now returns fake accepted output
- Added test for stub runner 200 response
- Aligned API status strings and top-level response logic with the exact spec
- Enforced HTTP request size limits (`http.MaxBytesReader`) to prevent payload OOM
- Added `/readyz` and `/info` endpoints returning system stats and smoke probes
- Authored required documentation files (`api.md`, `architecture.md`, `benchmarks.md`, `languages.md`, `security.md`)

### Archi
- Added nsjail runner for executing code in sandbox
- Added language registry (YAML/config driven)
- Added Python + C++ execution support
- Added integration tests for end-to-end runs
- Updated Dockerfile/main wiring for runner
- Implemented bounded concurrency queue (worker pool) to handle sustained load
- Capped child stdout/stderr streams to prevent OOM
- Migrated language registry entirely to `languages.yaml`
- Implemented Compiler Flag Allowlists for secure build/run templates
- Added toolchains for Java, Node.js, and Verilog in Dockerfile

### Shared
- 
