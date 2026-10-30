# Performance and reliability

## Observed public behavior

E-008–E-010 show a useful existing baseline: Next.js assets/RSC delivery and Vercel
HITs for sampled marketing pages. Home, services, products, contact and privacy
returned ETags and `public, max-age=0, must-revalidate`. Careers/blog returned
`private, no-cache, no-store, max-age=0, must-revalidate`, age 0 and MISS.

Browser freshness at max-age=0 requires revalidation; it does not mean the origin
is recomputing every request. A homepage Age of 108453 describes shared response
age, not latency or stale database records. `x-nextjs-stale-time` and
`x-nextjs-prerender` are framework metadata, not universal HTTP cache guarantees.

The initial curl homepage body was 100,498 bytes. This is one retrieved HTML
representation, not compressed wire size or a performance budget measurement.
The seven-page browser session captured 307 same-site responses including assets
and automatic linked-page prefetches. That count is a research-session artifact,
not 307 API calls per user or a load benchmark. Ordinary fetch/XHR responses in
this capture were RSC, not business JSON. Blog page/category links and different
second-page articles are visible (E-016); SQL pagination and query plans are unknown.

## Prioritized recommendations

| ID | Recommendation | What must be measured / verified |
|---|---|---|
| F-003 | Preserve public cache validators; review dynamic careers/blog render policy | Content update frequency, personalization and acceptable freshness; no-store may be intentional |
| F-003 | Profile RSC prefetch and HTML/JS transfer on representative mobile navigation | Transferred compressed bytes, cache reuse, route conversion and user latency; do not remove prefetch just because requests are visible |
| F-002 | Bound mutation work and protect email/provider budgets | Actual action quotas, request size/time limits, concurrency and provider failure behavior |
| F-004 | Recover action exceptions and correlate failures | Rejected transport/provider calls in authorized staging; no production failure was triggered |
| F-004 | Introduce status/duration metrics and trace IDs if missing internally | Existing telemetry and retention policy; lack of public metrics is not proof they are absent |

## Reference implementation: explicit bounds

- 4 KiB body, 16 KiB configured header budget, 32 active worker slots.
- 60 requests/minute per socket peer, at most 1,024 quota keys. Fixed windows allow
  boundary bursts; full tables reject new peers and expired keys are reclaimed.
- 2s application context deadline; server header/read/idle limits 5s/10s/60s;
  write budget exceeds the application deadline by 10s.
- Worker slots release on actual worker exit. A timeout does not create unlimited
  detached work. Go cannot forcibly kill a goroutine; future dependencies must
  honor cancellation or stuck workers can exhaust slots.
- Immutable pre-serialized catalog, SHA-256 ETag, conditional GET/HEAD, max-age=60
  and must-revalidate. Code changes update content/validator on restart. Errors
  and calculations use no-store; cacheable JSON excludes per-request IDs.
- Readiness becomes unavailable during shutdown; liveness stays available. SIGTERM
  allows 10s to drain before forced close. Tests and real smoke checks cover it.

## Conditional internal architecture work

No database or pool was observed. Review capped connection reuse, result sizes,
indexes and query plans only against actual slow paths. Privacy text's reference
to Prisma is not evidence of N+1 queries. If reliable email handoff is required,
use durable job acceptance and provider delivery state; do not promise delivery
on an uncommitted background goroutine. Introduce idempotency and bounded retries
only with clear mutation/provider semantics; blind retries can duplicate inquiries.

No benchmark or authorized production load test was run. Calculator throughput
would not establish Cynbit capacity. Measure real staging workloads and service
objectives before adding distributed caches, queues or microservices.
