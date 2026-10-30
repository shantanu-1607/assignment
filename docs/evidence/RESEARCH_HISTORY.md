# Historical blocked research (superseded)

The following ledger predates environment publication. Public-site access now
works; these failures are not current blockers or findings about Cynbit.

# Evidence ledger

## Scope and collection rules

Assessment: Cynbit, Backend Developer Intern. Collection attempted 2026-10-31
00:28–00:35 IST (2026-10-30 18:58–19:05 UTC). These times describe the research
session, not a target-server observation. Only public navigation and a single
homepage/header request were attempted. No scanning, enumeration, submissions,
authenticated requests, load generation, or crafted attack inputs were used.

| ID | Source and method | Result | Supported conclusion | Limits |
|---|---|---|---|---|
| E-001 | Candidate `https://cynbit.com/`, curl HEAD through inherited proxy | CONNECT rejected with HTTP 403 by environment proxy | Research access was blocked before origin TLS/HTTP | Does **not** establish that this domain is official, reachable, insecure, or unavailable |
| E-002 | Public search `https://www.google.com/search?q=Cynbit+official+website`, ordinary GET | CONNECT rejected with HTTP 403 | Independent company/domain corroboration unavailable | No search results were received |
| E-003 | Candidate homepage ordinary GET after saving draft allowlist additions | CONNECT rejected with HTTP 403 | Draft saving had not enabled runtime access | No HTML, cookies, API payloads or origin headers collected |

Sanitized proxy excerpt (E-001):

```text
HTTP/1.1 403 Forbidden
server: envoy
curl: (56) CONNECT tunnel failed, response 403
```

This is intermediary evidence, **not a Cynbit response**. No raw cookie, token,
credential or personal data was captured. The official site, company products,
backend technology, APIs and authentication flows remain **unverified**. This
report is an explicitly bounded assessment package, not a completed website audit.

## Candidate domain supplied by requester

On October 31, 2026 IST, the requester clarified that the assessment email
contained no website and suggested `https://www.cynbit.in/`. This is a better
research lead than the initial guessed `.com` domain, but it is not independent
confirmation of company identity. The `.com` requests above are retained only as
historical unsuccessful research attempts, not as a claimed company association.

| ID | Source and method | Result | Supported conclusion | Limits |
|---|---|---|---|---|
| E-004 | User-suggested `https://www.cynbit.in/`; one ordinary HTTPS GET through the configured proxy, October 31, 2026 at 00:41:20 IST (October 30 at 19:11:20 UTC) | CONNECT rejected with HTTP 403; curl exit 56 | Active environment policy still blocks candidate-site research | No origin response, company identity, products, API, cookies or headers verified |

The response timestamp is the local completion time, not an origin timestamp.
`cynbit.in` and `www.cynbit.in` were added to the environment configuration draft,
preserving its existing domains. Saving that draft does not activate the running
network policy. Review/save the settings and publish the environment through the
product before retrying; do not bypass the proxy or substitute third-party fetch
services. If normal navigation redirects to another domain, inspect the redirect
and request only the necessary supported policy addition.

The current public-review scope and evidence are in [EVIDENCE](../EVIDENCE.md).
