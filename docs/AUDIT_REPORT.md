# Backend technical assessment: Cynbit

## Executive summary

This submission pairs a deliberately bounded audit with a runnable Go reference
backend. **No Cynbit-origin behavior was verified.** The environment rejected
public research requests at its outbound proxy; even the candidate official domain
could not be independently corroborated. E-001–E-003 in [EVIDENCE](EVIDENCE.md)
document the limitation. There are no confirmed target vulnerabilities, performance
measurements, API endpoints, products or stack identifications in this report.

The highest-value deliverable currently supported by evidence is a small, tested
reference design: strict JSON contracts, sanitized errors, bounded resource use,
context cancellation and conditional HTTP caching. Findings F-001–F-004 are
**conditional engineering opportunities**, not defects attributed to Cynbit.

## Scope and confidence

| Label | Meaning in this report |
|---|---|
| Observed | Directly collected evidence; here, research-proxy rejection only |
| Inferred | Interpretation with explicit assumptions; no target stack inference made |
| Proposed | Standalone design that could apply if corresponding product needs exist |
| Unverified | Requires permitted origin/browser or internal evidence |

The employer's assessment context is user-provided. Independent confirmation of
company identity, offerings and publicly visible flows remains outstanding. The
reference service names are invented examples, explicitly labeled in code and
responses. They must never be quoted as discovered production routes or payloads.

## Backend architecture and lifecycle

Observed: a local client attempted to access a candidate public domain through the
managed proxy and was denied before origin communication. This tells us nothing
about Cynbit's CDN, backend language, hosting or data stores.

Proposed: HTTP edge termination feeds a single Go process; API policy surrounds
small handlers that call a context-aware service. No database, background workers,
authentication provider or external network calls are needed for this stateless
example. Its response lifecycle is documented in [ARCHITECTURE](ARCHITECTURE.md).
Adding services or a database without a demonstrated need would introduce failure
modes without improving this assessment.

## API and payload analysis

All target-specific entries below remain unverified. The rightmost column is a
reference design, not an observation.

| Area | Needed evidence | Illustrative implementation |
|---|---|---|
| Naming/versioning | Actual ordinary-navigation requests | `/v1/services`, `/v1/estimates` |
| Methods and status | Request/response pairs | GET/HEAD catalog; POST stateless calculation; 200, 304, 400, 404, 405, 413, 415, 422, 429, 500, 503 |
| JSON and validation | Public payloads, documented input schema | Exact keys and types, required fields, duplicate/null rejection, bounded ranges |
| Error hygiene | Naturally occurring errors | Stable error code/message and server-generated request ID; internal details excluded |
| Pagination | Actual collection size and access pattern | Omitted for two immutable entries; add bounded cursor pagination for real large collections |
| Data exposure | Public response field inventory | No personal/contact data; only illustrative service IDs and effort units |
| Caching | Origin cache headers and content semantics | Public catalog uses ETag, 60-second freshness and conditional GET; errors and estimates use no-store |
| Timeout/retry | Actual behavior and upstream contracts | Two-second application budget; 503 on deadline; no external retries |
| Idempotency | Side effects and failure/retry semantics | Calculation has no persisted effect; repeated POST yields same domain result, different request ID |

The catalog returns `{ "data": [...] }` with request ID in `X-Request-ID` only.
Putting per-request IDs in a cacheable representation would change its ETag every
request. Uncacheable responses include `request_id` in JSON as well.

## Security and performance conclusions

See [SECURITY_FINDINGS](SECURITY_FINDINGS.md) and
[PERFORMANCE_ANALYSIS](PERFORMANCE_ANALYSIS.md). No severity rating is assigned to
unobserved target behavior. Priorities express the value of validating a control,
not a confirmed breach risk. The reference backend demonstrates selected controls;
it is not a security certification or a drop-in production application.

## Prioritized next work

1. Restore permitted research access and verify company identity. This gates all
   target-specific conclusions and is more valuable than expanding the demo.
2. Capture representative, sanitized browser evidence and map actual contracts.
3. Validate applicable controls against owners' internal documentation and existing
   tests, with authorization. Public absence of a header or endpoint cannot prove
   absence of server-side validation, authorization or observability.
4. Reprioritize conditional recommendations with impact, confidence and cost.

No real site latency, throughput, database efficiency or scalability result is
claimed. Required follow-up is explicit rather than hidden behind generic findings.
