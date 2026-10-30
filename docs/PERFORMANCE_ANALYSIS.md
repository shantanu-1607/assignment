# Performance and reliability analysis

No origin measurements were possible. Proxy rejections are excluded from latency
analysis. No Cynbit bottleneck, query pattern, N+1 problem, missing pool or slow
endpoint is asserted. Local tests validate behavior, not production capacity.

## Proposed controls

**F-002: bounded work.** A 4 KiB request body, 16 KiB header budget, 32 active
handler slots and a per-peer fixed window of 60 requests/minute bound representative
resources. The limiter stores at most 1,024 peer records. Expired records are
reclaimed on admission when full, under one mutex. Saturated tables reject new
peers; fixed windows can permit twice the limit across a window boundary. These
are deliberate demo tradeoffs, not a DDoS defense. A trusted edge, global quotas
and per-identity policies are needed at scale.

**F-003: immutable public caching.** The catalog is serialized once, with SHA-256
ETag computed over those bytes. HEAD shares GET headers and sends no body. Weak
or strong If-None-Match validators and wildcard yield empty 304 responses on a
match. Max-age=60 allows bounded staleness; must-revalidate prevents stale reuse
once freshness expires without validation. Updating code updates the representation
and validator on restart. No cache invalidation service or database exists.

**F-004: lifecycle budgets.** The server has 5s header, 10s read and 60s idle
limits; write timeout exceeds application deadline by 10s. Application work gets
a 2s context budget by default. Go's timeout wrapper buffers the small response,
returns 503 on deadline and cancels work. Slots release when work exits, not when
the HTTP deadline fires: noncooperative dependencies cannot create unbounded
background execution. Go cannot kill stuck goroutines; a stuck worker can occupy
a slot indefinitely. Domain and future outbound calls must honor context.
Readiness becomes unavailable during shutdown, liveness stays available; the
server allows 10s to drain then closes active connections. Open connections and
processes must be restarted after environment restoration.

## Conditional scaling recommendations

| Trigger to measure first | Recommendation | Tradeoff |
|---|---|---|
| Repeated unchanged public data | Shared edge cache with explicit invalidation/freshness | Lower origin load versus staleness |
| Multiple API replicas / proxy ingress | Edge/global limiter and trusted identity mapping | Consistent quotas versus infrastructure/failure dependency |
| Actual database-backed reads | Bound result sizes; measure query plans; reuse capped connection pools | Throughput versus memory, connections and queueing |
| External dependency failures | Per-call deadlines; bounded retries only for safe/idempotent operations; jitter | Recovery versus amplified load and tail latency |
| Latency or availability objectives | Status/duration metrics, traces, error budgets and load tests on authorized staging | Diagnostic value versus operational cost and data retention |

No retries or connection pools were added to a calculator with no outbound calls.
No benchmark is included: nanosecond arithmetic throughput would not establish
readiness for Cynbit's unknown workload. Use representative authorized staging
load only after actual architecture, targets and traffic are known.
