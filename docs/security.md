# Security Fixes

The following vulnerabilities present in the reference implementation have been closed in goboxd:

1. **Path Traversal via Filename:** Fixed in `internal/types/validate.go`. `validateFilename` strictly rejects `.` , `..`, and any slashes.
2. **Shell-style Directory Commands:** Fixed in `internal/runner/runner.go`. We use `os.MkdirTemp` and direct file I/O `os.WriteFile`, completely avoiding `sh -c`.
3. **Compiler-Flag Injection:** Fixed in `internal/runner/runner.go`. The `validateFlags` function ensures any user-supplied flags match the `flag_allowlist` defined in `languages.yaml`. Unsafe flags return HTTP 400.
4. **No Request Size Limits:** Fixed in `internal/server/server.go`. We wrap `r.Body` in `http.MaxBytesReader` (capped at 1 MiB) and validate `MaxTests` in `types.ValidateRunRequest`.
5. **UID Collisions Under Load:** Fixed in `internal/runner/runner.go`. `os.MkdirTemp` guarantees unique directory names automatically using process-unique random suffixes.
6. **Unbounded Child Output (OOM):** Fixed in `internal/runner/runner.go`. Both `cmd.Stdout` and `cmd.Stderr` are assigned to a custom `cappedWriter` which stops accumulating memory after 1 MiB and appends a `[output truncated]` marker.
7. **Stale Jail Directories:** Fixed in `internal/runner/runner.go`. A `defer os.RemoveAll(workDir)` is invoked immediately after temporary directory creation, ensuring cleanup on every exit path including panics.
