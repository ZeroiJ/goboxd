# Architecture

goboxd is a robust Go HTTP daemon that safely executes untrusted code using `nsjail`.

## Core Components
- **`cmd/goboxd`:** The application entrypoint. Initializes the language registry and starts the HTTP server.
- **`internal/server`:** Contains the HTTP transport layer. It implements a **bounded concurrency queue** utilizing a buffered Go channel (`s.Sem`). This ensures the system queues requests up to the defined `MAX_CONCURRENT_JOBS` (defaulting to `runtime.NumCPU()`) instead of overloading the host or dropping requests.
- **`internal/runner`:** The sandboxing engine. It evaluates `RunRequests`, formats the command templates, and spawns `nsjail` processes. It manages unique temporary directories and capped memory buffers to safely capture output.
- **`internal/types`:** Contains request/response structures and standard validations.

## Concurrency Model
The server utilizes `sync/atomic` counters to track active, completed, and failed jobs. Requests exceeding the concurrent job limit block on the `s.Sem` channel until an execution slot frees up, ensuring sustained throughput under heavy load without degrading host stability.
