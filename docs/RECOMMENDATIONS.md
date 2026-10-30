# Recommendations and decisions

## Evidence-to-code traceability

All applicability is **proposed/unverified** because E-001–E-003 prevented target
observation. The chain is: documented access limitation → explicit assumption of
a possible public JSON workflow → general engineering control → runnable example
→ behavior test. It is not a chain of discovered vulnerabilities.

| ID | Proposed improvement | Code | Principal verification |
|---|---|---|---|
| F-001 | Strict contracts and consistent errors | `internal/api/api.go`: `decodeInput`, `route`, `fail` | `TestEstimateValidation`, `TestBodyLimitBoundary`, `TestRoutingAndHeaders` |
| F-002 | Bound time, memory, active work and quotas | `internal/api/limiter.go`, `API.ServeHTTP` | `TestRateLimitIgnoresForwardedHeaders`, `TestLimiterConcurrent`, `TestLimiterBoundAndExpiry`, `TestAdmissionBound` |
| F-003 | Public-only conditional caching | Catalog initialization and GET/HEAD route | `TestCatalogCaching` |
| F-004 | Cancellation, correlation and draining | `API.ServeHTTP`, `cmd/server/main.go` | `TestDeadlineCancelsWorker`, `TestClientCancellation`, `TestTimeoutRetainsSlotUntilWorkerStops`, manual SIGTERM smoke |

## Roadmap

| Priority | Action | Impact / confidence / cost |
|---|---|---|
| Gate | Corroborate official domain and collect public navigation evidence | Essential to company-specific credibility; confidence high; low effort once access granted |
| Quick win, conditional | Review real API schemas, error policy, cache directives and request correlation | Broad benefit if APIs exist; applicability unknown; low/medium effort |
| Medium term, conditional | Validate identity/object authorization; trusted edge quotas; operational metrics | High value for authenticated/expensive workflows; applicability unknown; medium cost |
| Longer term, evidence-driven | Optimize actual query plans, capacity, distributed caching and failure isolation | Impact depends on measured bottlenecks; highest uncertainty/cost |

Do not treat missing public evidence as a reason to replace an existing working
system. Confirm current controls with maintainers before proposing migrations.

## ADR-001: standard library and stateless example

Accepted. Go's `net/http`, `context`, JSON and structured logging suffice. A small
catalog and deterministic estimator demonstrate contracts without collecting
contact details or pretending to implement an unverified Cynbit workflow. No
third-party module, database or framework increases review/setup cost. Loss:
this cannot demonstrate transactions, real authorization or persistence.

## ADR-002: in-memory fixed-window admission

Accepted for demonstration. Mutex-protected quota state and bounded keys keep
memory predictable. Table saturation fails closed; existing peers continue within
their quotas. Availability for new peers can suffer when the table fills. Quotas
reset on restart and are per replica; window-boundary bursts are possible. Adopt
edge/global quotas if deployment requires shared guarantees. Do not silently trust
forwarded headers to make per-user quotas appear to work behind a proxy.

## ADR-003: caching without correlation in the representation

Accepted. Stable public catalog bytes yield stable ETags; the request ID stays in
headers on cacheable responses. Errors and calculations include it in JSON and
use no-store. A shared cache can replay the origin's correlation header; an edge
should issue its own delivery ID if needed. A 60s freshness window trades immediate
updates for fewer origin requests. No sensitive representation is cacheable.

## ADR-004: deadlines and admission cooperate

Accepted. `http.TimeoutHandler` bounds client wait for these tiny responses and
passes a cancelable context to the service. Worker-held slots preserve the active
work cap even when cancellation is ignored. A stuck service can exhaust slots;
Go has no safe goroutine kill. Future dependencies need bounded, context-aware
operations. Do not wrap unbounded/streaming responses in this buffered handler.

## ADR-005: no speculative auth, retries or infrastructure

Accepted. There is no user state, side effect or upstream dependency. Adding fake
JWT authentication, a Redis instance or retry loops would obscure the actual
scope. Introduce identity and object-level authorization before private data,
idempotency storage before retryable mutations, and controlled backoff only when
real upstream semantics justify it.
