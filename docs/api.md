# goboxd API

## POST /run
Executes untrusted code in a sandboxed environment.
- **Request Body:** JSON payload including `language`, `source`, `tests`, and optional `build`/`run` configurations.
- **Limits:** Request body capped at 1 MiB. Max tests capped at 50.
- **Response:** 200 OK with `status`, `build`, and `tests` containing execution outcomes.

## GET /healthz
Liveness probe.
- **Response:** 200 OK with `{"status": "ok"}`.

## GET /readyz
Readiness probe. Verifies that `nsjail` is executable and all configured languages pass their smoke probes.
- **Response:** 200 OK if fully healthy, or 503 Service Unavailable if any probe fails.

## GET /info
System information and statistics.
- **Response:** JSON payload detailing build information, `nsjail` path and version, registered languages, configured limits, and internal concurrency stats.
