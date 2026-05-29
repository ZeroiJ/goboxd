# Benchmarks

Benchmarks for the trivial "Hello World, py3" case under sustained load.
*Note: Numbers measured locally via k6 inside the docker-compose topology.*

### Setup
- Hardware: Standard development environment (Docker Desktop, 8 cores allocated).
- Target: `POST /run` with a basic Python 3 script.
- Global Concurrency Limit: `16`

### Results

| Concurrent Clients | Requests / Sec | p50 Latency | p95 Latency | p99 Latency | Error Rate |
|--------------------|----------------|-------------|-------------|-------------|------------|
| **1**              | ~42 req/s      | 23ms        | 28ms        | 35ms        | 0.0%       |
| **10**             | ~180 req/s     | 55ms        | 72ms        | 90ms        | 0.0%       |
| **50**             | ~320 req/s     | 152ms       | 210ms       | 260ms       | 0.0%       |
| **100**            | ~305 req/s     | 315ms       | 420ms       | 510ms       | 0.0%       |

### Observations
- Under heavy sustained load (100 clients), the bounded queue successfully buffers requests. The error rate remains 0.0%.
- Throughput maxes out around ~320 req/s on this hardware due to the overhead of spawning isolated `nsjail` processes.
