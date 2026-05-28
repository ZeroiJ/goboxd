# AI usage log (prompts)

## 2026-05-24 · Initial HTTP server + validation

**Prompt:**
Asked an AI agent for a step-by-step plan and initial skeleton for the Go HTTP server, request/response types, and validation rules based on the hackathon spec.

**Response summary:**
Suggested a small net/http server with /healthz and /run, request/response structs matching the spec, and validation for required fields and filename safety.

**What we used / didn't use:**
Used the server/handler structure and validation ideas; did not copy code verbatim and adjusted details to match our repo layout.

## 2026-05-24 · README cleanup

**Prompt:**
Asked an AI agent to rewrite the README to be short and strict: what it is, how to run, where docs are.

**Response summary:**
Provided a minimal README format with run/test commands and a curl example.

**What we used / didn't use:**
Used the structure and examples; removed any extra descriptive text.
