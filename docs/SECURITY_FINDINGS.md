# Security findings and control review

**Confirmed Cynbit vulnerabilities: none. Target security posture: unverified.**
E-001–E-004 are environment access failures. They are not security findings about
Cynbit. The following opportunities apply only if matching functionality exists.

| ID | Category / applicability | Evidence and confidence | Recommendation / reference coverage |
|---|---|---|---|
| F-001 | Validation and error contracts; if JSON APIs accept input | No target payload available; low applicability confidence, high confidence in general practice | Strict schema, bounded body, stable codes and correlation IDs. `decodeInput`, `route`, `fail`; `TestEstimateValidation`, `TestBodyLimitBoundary`, `TestFailureHygiene` |
| F-002 | Abuse/resource limits; if endpoints perform costly work | No target abuse control observed or tested | Bound bodies, time, active work and peer counters. `limiter.allow`, `ServeHTTP`; `TestLimiterConcurrent`, `TestAdmissionBound`, `TestTimeoutRetainsSlotUntilWorkerStops` |
| F-003 | Data/cache isolation; if public and private responses coexist | No target cache headers collected | Cache only explicitly public representations; use no-store on errors and calculations. `route`, `fail`; `TestCatalogCaching`, error assertions |
| F-004 | Reliability and diagnostic hygiene; if services have dependencies | No origin error/log evidence | Propagate cancellation, sanitize failures, generated request IDs and structured logs. `ServeHTTP`, `Calculator.Estimate`, server lifecycle; `TestDeadlineCancelsWorker`, `TestClientCancellation`, `TestFailureHygiene` |

## Authentication, authorization and browser boundaries

No public login, cookies, session handling or authorization behavior could be
observed. Do not infer an authentication flaw from the lack of evidence. This
reference API is intentionally public and stateless: it stores no user data and
has no privileged operations. Adding protected data requires an identity boundary
and object-level authorization tests before launch, not simply a token parser.

For a future cookie-backed product, use Secure, HttpOnly and appropriate SameSite
attributes; rotate session identifiers, enforce expiry and authorize every object
access server-side. Evaluate CSRF for state-changing cookie-authenticated requests;
SameSite alone is not a complete CSRF design. None of those controls is claimed to
exist or fail on Cynbit. They are outside the public calculator's scope.

The reference permits no cross-origin browser access (no CORS headers); requiring
JSON does not authenticate callers or stop non-browser clients. If browser clients
need CORS, use an explicit origin allowlist with Vary: Origin, and separately
authorize requests. No reflected Origin or credentialed wildcard policy.

## Implemented controls and deployment boundaries

- CSP `default-src 'none'; frame-ancestors 'none'`, nosniff, no-referrer and DENY
  frame policy. JSON API headers do not replace HTML app security review.
- TLS and HSTS belong to an explicitly configured HTTPS reverse proxy. Local
  server binds loopback by default; never advertise this plain HTTP process as a
  secure public deployment. Do not set HSTS on local HTTP demonstrations.
- Never trust client-supplied request IDs or X-Forwarded-For. Direct socket peers
  define quotas. Behind a proxy this groups all users into one quota until a
  trusted-proxy boundary is deliberately implemented.
- Logs include generated IDs and duration, not request bodies, IPs, query strings,
  authorization headers or internal error strings. Configure retention and access
  controls before deployment; no metrics/tracing platform is claimed here.
- No credentials, external SDKs or runtime dependencies. If secrets become needed,
  use managed environment bindings and rotation; never commit `.env` values.

## Standards informing recommendations

These are technical references, not evidence about Cynbit:

- [OWASP API Security Top 10 2023](https://owasp.org/API-Security/editions/2023/en/0x11-t10/): object-level authorization, resource consumption, security configuration.
- [OWASP REST Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/REST_Security_Cheat_Sheet.html): HTTPS, validation, error handling and access control.
- [OWASP CSRF Prevention](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html): conditional guidance for cookie-authenticated state changes.
- [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110): HTTP semantics and conditional requests.
- [RFC 9111](https://www.rfc-editor.org/rfc/rfc9111): cache control semantics.

Reference URLs identify established standards; this environment did not retrieve
those documents during research. No assertion of a standards compliance audit.
