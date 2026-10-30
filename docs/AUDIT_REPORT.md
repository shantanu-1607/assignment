# Cynbit: backend technical audit

## Executive summary

Cynbit's public site at [www.cynbit.in](https://www.cynbit.in/) presents AI automation,
custom software/cloud services and three products: Cynbit Studio, TenantEase and
Revio. Its careers page lists the relevant Backend Developer Intern role. Brand,
location, domain-linked social profile and matching role support the selection;
the employer did not provide a URL, and independent ownership confirmation remains
limited. See E-005–E-007 in the [evidence ledger](EVIDENCE.md).

The public marketing frontend exhibits Next.js/RSC and Vercel caching indicators.
Static marketing pages already benefit from cache hits and validators; careers
and blog documents use no-store. Forms integrate with Next.js Server Actions in
public client bundles. There was no ordinary-navigation REST JSON API payload to
audit, and no submission was made to manufacture one.

The strongest improvement opportunities are client/server contract alignment for
form lengths, local recovery around action exceptions, sampled response-header
hardening, applicant data minimization, and analytics inventory/consent review.
No exploit, authentication bypass, production bottleneck or data leak is claimed.

A tested standalone Go reference API demonstrates strict input contracts,
consistent safe errors, bounded work, deadline propagation and public-only
conditional caching. It is a reference design, not a clone or integration with
Cynbit. The examples remain deliberately illustrative and collect no personal data.

## Scope, evidence and limits

Observation: October 31, 2026 IST, ordinary public browser navigation plus public
headers and browser-loaded assets. No accounts, forms, authenticated product areas,
private APIs, crafted failures, load generation or vulnerability scanners. Cookies
were inspected by name/attributes only. Public contact details and applicant data
are not included in evidence artifacts. [EVIDENCE](EVIDENCE.md) provides exact
sources, UTC timestamps, sanitized excerpts and limitations.

- **Observed:** directly visible behavior, DOM, headers or public client code.
- **Inferred:** interpretation with stated assumptions, not runtime proof.
- **Proposed:** engineering recommendation or standalone reference design.

## Backend architecture

**Observed (E-008–E-010):** the browser receives HTML, versioned Next.js assets and
RSC responses. Link prefetching requests normally linked pages. Header evidence
supports prerendering and Vercel edge caching for sampled pages. Careers/blog are
dynamically delivered or otherwise explicitly uncacheable at the observed layer.

**Disclosed (E-017):** privacy text names Resend and Prisma. This is not proof of a
specific email pipeline, schema, database engine or live credential configuration.
The services page's broad technology list is a capability claim, not a stack bill
of materials.

**Inferred:** inquiry/demo/newsletter actions likely hand data to server-side
handlers. Their implementation, validation, storage and delivery guarantees are
not publicly demonstrated. Revio's advertised role hierarchy and TenantEase's
financial features imply important future authorization boundaries, but those
products were not exercised and their controls cannot be graded here.

**Proposed:** retain public-page caching, validate untrusted input at every action
boundary, use explicit response/error contracts, and place irreversible delivery
behind a durable handoff only if real delivery requirements justify it. No reason
was established to split this website into microservices. See [architecture
comparison](ARCHITECTURE.md).

## API, payload and error review

| Area | Observed | Assessment / recommendation |
|---|---|---|
| Endpoints and methods | Normal page GETs, RSC GETs and client Server Action references | No public REST/version convention established; do not relabel RSC as a business JSON API |
| Payloads | Contact fields fullName/email/company/industry/message; demo fields plus client-added serviceType | Analyze as visible input contracts; actual wire action encoding and server validation unverified |
| Codes | 200 on inspected pages; no production errors deliberately caused | Reference has explicit 4xx/5xx; no claim about actual action error status |
| Validation | HTML required/email constraints; no HTML maxLength | Share length/range rules with server; HTML validation alone never protects a backend (F-001) |
| Error handling | Action callers branch on success and show generic messages | Good error-detail hygiene; add exception recovery and guaranteed pending cleanup (F-004) |
| Pagination | Blog has page/category links; second page changes articles | Bounded public pagination exists; internal query shape and maximum limits unknown |
| Caching | Public HIT + ETag on marketing; no-store MISS on careers/blog | Preserve working cache; review freshness needs before changing dynamic pages (F-003) |
| Retry/idempotency | Delivery semantics not tested | Avoid blind retries of inquiry/email mutations; add idempotency only if retry/replay requirements warrant it |
| Exposure | Ordinary HTML/RSC and public client code | No private applicant payloads, session tokens or authenticated API data inspected |

Illustrative reference contracts, **not Cynbit endpoints**:

```http
GET /v1/services
POST /v1/estimates
Content-Type: application/json

{"service_id":"api-review","integrations":2}
```

A valid calculation returns 200 and effort_units 7; errors have stable code/message
and request ID. The cacheable catalog puts request ID only in its header so
representation bytes and ETag remain stable. Full examples are in the README.

## Priorities and confidence

1. **F-007 / data minimization:** mandatory applicant DOB/address/dual contact
   fields justify a purpose and collection-stage review. Strong UI evidence;
   retention/access controls remain unknown.
2. **F-001 and F-004 / submission reliability:** align client/server bounds and
   recover rejected action calls. Strong DOM/client-code evidence; no server
   rejection defect or live failure proven.
3. **F-005 / response policy:** add a reviewed CSP, nosniff, referrer and framing
   policy on sampled HTML routes. Header absence confirmed, exploitability unknown.
4. **F-006 / analytics:** reconcile pre-choice analytics cookies with consent,
   documented categories and applicable regional requirements. Not a legal finding.
5. **F-002 and F-003 / capacity and caching:** validate existing private controls
   internally and measure before optimizing. No demonstrated production bottleneck.

Detailed confidence, scope and tradeoffs are in [security findings](SECURITY_FINDINGS.md),
[performance](PERFORMANCE_ANALYSIS.md), and [recommendations](RECOMMENDATIONS.md).
