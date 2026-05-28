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

### Archi
- Added nsjail runner for executing code in sandbox
- Added language registry (YAML/config driven)
- Added Python + C++ execution support
- Added integration tests for end-to-end runs
- Updated Dockerfile/main wiring for runner

### Shared
- 
