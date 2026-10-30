# Improvement roadmap and engineering decisions

## Evidence-to-implementation traceability

This matrix distinguishes implemented examples from recommendations that need
product-specific work. No reference endpoint is represented as a discovered
Cynbit API, and no test is represented as a production-site check.

| Finding | Evidence and implication | Proposed solution | Reference implementation | Verification |
|---|---|---|---|---|
| F-001 | E-012: no HTML length caps; backend checks unknown | Shared contract and independent server bounds | `internal/api/api.go`: `decodeInput`, `route`, `fail` | `TestEstimateValidation`, `TestBodyLimitBoundary`, `TestRoutingAndHeaders` |
| F-002 | E-013: submission actions imply resource/provider work; protection unknown | Quotas, body/time limits, bounded active work | `limiter.go`, `API.ServeHTTP` | `TestLimiterConcurrent`, `TestLimiterBoundAndExpiry`, `TestRateLimitIgnoresForwardedHeaders`, `TestAdmissionBound` |
| F-003 | E-009/E-010: real public-cache/no-store separation | Preserve public-only caching and evaluate freshness | Catalog initialization and GET/HEAD route | `TestCatalogCaching`; error helpers assert no-store |
| F-004 | E-013: no local rejection recovery in inspected callers | Try/catch/finally client adapter; safe server errors and cancellation | API error/deadline handling; illustrative adapter below | `TestFailureHygiene`, `TestDeadlineCancelsWorker`, `TestClientCancellation`, `TestTimeoutRetainsSlotUntilWorkerStops`, SIGTERM smoke |
| F-005 | E-011: sampled HTML header omissions | Reviewed HTML policy; CSP rollout in report-only | API JSON headers only; not a ready-made HTML CSP | `TestRoutingAndHeaders`; target HTML policy needs separate staging tests |
| F-006 | E-015: analytics cookies before consent choice | Inventory/region-appropriate optional-tag gating | Not implemented: stateless API has no tracking | Browser consent/tag tests required in authorized staging |
| F-007 | E-014: mandatory applicant data beyond basic contact | Defer fields; purpose/retention/access review | Not implemented: demo collects no applicant data | Hiring-policy and deletion/access tests require internal scope |

## Concrete client reliability pattern (F-004)

Illustrative TypeScript adapter; it does not import, invoke or recreate Cynbit's
Server Actions. Wire it to a documented application contract in an authorized
codebase. This is a proposed snippet, not executed site code.

```ts
type SubmitResult = { success: boolean };
async function submitSafely(
  action: () => Promise<SubmitResult>,
  setPending: (pending: boolean) => void,
  setError: (message: string | null) => void,
  onSuccess: () => void,
): Promise<void> {
  setPending(true);
  setError(null);
  try {
    const result = await action();
    if (result.success) onSuccess();
    else setError("Unable to submit. Please try again or contact support.");
  } catch {
    setError("Connection interrupted. Please check before submitting again.");
  } finally {
    setPending(false);
  }
}
```

Do not automatically retry a mutation on ambiguous transport failure: it may
already have succeeded. Use a documented idempotency contract or a user-visible
reconciliation flow if duplicate leads/emails matter. Add a bounded timeout at the
server/provider boundary; a browser Promise.race alone does not stop server work.

The reference's actual Go boundary uses `http.MaxBytesReader(w, r.Body, 4096)`
before decoding, rejects duplicate/unknown/null/missing keys and trailing JSON,
then checks catalog IDs and a 0..20 integer range. Those checks are executable
examples of independent server validation. Contact/email/applicant schemas would
need their own domain limits, Unicode rules and privacy decisions.

## Roadmap

| Horizon | Work | Impact / confidence / cost |
|---|---|---|
| Quick win | Add local exception cleanup and accessible errors to action callers | Clear source-backed reliability opportunity; low cost; reproduce rejected calls only in staging |
| Quick win | Define input lengths; align browser and server rules | Confirmed UI omission; server enforcement unknown; low/medium cost |
| Quick win | Review response policy and stage CSP reporting | Confirmed sampled header gap; low effort for nosniff/referrer, medium for Next.js-compatible CSP |
| Near term | Review mandatory applicant fields, purposes, access and retention | Directly observed collection scope; potentially high privacy value; needs hiring/product owners |
| Near term | Inventory analytics cookies and consent behavior | Direct observation; regional requirements unknown; product/legal review before conclusions |
| Medium term | Inspect actual action quotas, provider budgets and delivery guarantees | Potentially high abuse/reliability value; existing controls unknown; medium cost |
| Evidence-driven | Measure dynamic-page caching, prefetch costs and query performance | Avoid speculative optimization; depends on traffic, personalization and freshness objectives |
| Longer term | Cross-tenant/role authorization tests for advertised products | High value for sensitive product workflows; requires separate explicit internal/staging scope |

## ADR-001: standard-library stateless reference

Accepted. A small catalog and deterministic calculator isolate HTTP engineering
without storing applicant/contact information or connecting to production. No
framework, DB or external module is needed. Tradeoff: no transaction, durable
mutation, email-delivery, identity or tenant isolation demonstration. Real observed
forms justify the validation/reliability patterns without requiring a fake CRM.

## ADR-002: bounded in-memory fixed windows

Accepted for the demo. Mutex-protected state is capped at 1,024 peers; full tables
fail closed for new peers, preserving current counters. This sacrifices availability
for new peers under saturation. Quotas are per process and reset on restart;
window-boundary bursts are possible. Distributed/edge quotas are conditional on
actual deployment. Forwarded headers remain untrusted until an explicit trusted
proxy boundary exists.

## ADR-003: stable public representations

Accepted. Serialize catalog once and compute an ETag over those exact bytes.
Per-request IDs stay in headers on cacheable reads; private/error/calculation
responses use no-store. Shared caches can replay an origin ID; a real edge should
issue its own delivery ID where needed. Max-age=60 is the demo's chosen freshness
budget, not a claim about Cynbit's needs. Keep private representations separate.

## ADR-004: cooperative deadlines plus worker-held slots

Accepted. `http.TimeoutHandler` bounds response wait and cancels context for small
buffered outputs. Slots release when workers actually stop, so timed-out work
cannot bypass the concurrency cap. Stuck workers can still exhaust capacity;
future dependencies must honor context. The wrapper is inappropriate for unbounded
streaming responses. Graceful shutdown drains before forced connection close.

## ADR-005: no speculative infrastructure

Accepted. No auth, Redis, queue or retries were added to a calculator without
private state or upstream effects. If implementing the real lead pipeline, first
establish existing validation/delivery behavior; add durable handoff, identity,
idempotency or retries only against concrete requirements. No evidence supports
rewriting the existing website as microservices.
