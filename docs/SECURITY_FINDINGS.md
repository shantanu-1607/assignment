# Security, validation and privacy findings

No exploited vulnerability or breached authorization boundary was established.
Priorities below reflect useful engineering work, not invented CVSS scores.
Each confirmed observation is bounded to the public sample in [EVIDENCE](EVIDENCE.md).

| ID | Category / classification | Evidence | Significance and recommendation |
|---|---|---|---|
| F-001 | Validation: confirmed client gap; backend unknown | E-012: required/email fields, no HTML maxLength on contact/demo | Declare shared lengths and types and enforce server-side limits independently. Reference strict decoder/body/domain limits demonstrate boundary enforcement; no claim that Cynbit's server accepts oversized data |
| F-002 | Abuse control: conditional recommendation | E-006/E-013: lead/subscription workflows have potential email/processing side effects; limits not tested | Verify per-action rate limits, body caps, deadlines and bounded upstream work internally. Reference bounded quota/admission/deadline shows the control pattern |
| F-003 | Cache/data isolation: observed good baseline + opportunity | E-009/E-010: public marketing cache vs no-store careers/blog | Preserve this separation. Assess dynamic-page freshness and personalized response isolation before enabling more caching; no cached private data found |
| F-004 | Reliability/error handling: source-backed inference | E-013: visible action callers lack local rejection recovery; contact resets pending after await | Add try/catch/finally and accessible retryable errors. A rejected await could skip local cleanup; this was not induced in production. Generic existing messages avoid exposing internals |
| F-005 | Response security: confirmed sampled header omissions | E-011: no CSP, X-Frame-Options, nosniff or Referrer-Policy response header; HSTS present | Add nosniff and reviewed referrer/framing policy; stage CSP in report-only before enforcement. Review Next.js inline scripts, nonces and third-party tags. No XSS/clickjacking exploit attempted or established |
| F-006 | Analytics/privacy: confirmed cookies, policy review needed | E-015: three analytics/measurement cookie names before any consent choice; SameSite=Lax, Secure=false | Inventory tags, purpose/lifetime and regional consent requirements; gate optional tags where required and set Secure where supported. Analytics cookies are not evidence of insecure auth sessions; no legal noncompliance conclusion |
| F-007 | Applicant data minimization: confirmed UI scope | E-014: required DOB, gender choice, address, phone and WhatsApp; privacy summary is broad | Evaluate whether each field is necessary at initial application; defer excess fields to later hiring stages, explain purposes and retention, restrict reviewer access and audit deletion. Data collection is not itself a leak or proof of unlawful processing |

## Authentication and authorization

No login/session/private product workflow was exercised. Marketing pages can be
public without authentication. Revio advertises Admin/Manager/Advisor roles and
TenantEase advertises tenant/payment data (E-006); these make object-level and
tenant-level authorization important **review targets**, not verified flaws.
Require authorization tests for cross-tenant IDs and role changes in an authorized
staging environment. Do not attempt those tests against production as part of
this passive assessment.

For any cookie-backed product: review Secure/HttpOnly/SameSite, session rotation,
expiry and CSRF protection. The observed analytics cookies do not establish the
properties of future authentication cookies. Client-visible consent/required
attributes never establish server-side authorization or validation.

Some public HTML/RSC responses permit `Access-Control-Allow-Origin: *`. No
credentialed/private response was captured, so a CORS data-exposure finding would
be unsupported. The reference API deliberately has no CORS headers or cookies.

## Reference controls and deployment limits

The Go reference implements response headers, strict body/type/range checks,
server-generated IDs, no-store errors, and sanitized failure messages. Its logs
exclude bodies, IPs, query strings and credentials. It is public/stateless and
cannot demonstrate real tenant authorization, consent storage or applicant data
retention. No TLS termination or identity provider is included; local binding is
loopback. Configure HTTPS/HSTS at a trusted edge for a real deployment.

Forwarded headers are ignored to prevent quota spoofing. A reverse proxy would
collapse clients into one quota until trusted identity mapping is implemented.
The demo limiter is per replica, bounded and restart-volatile, not a DDoS defense.

## Technical references

Standards inform recommendations, not claims about Cynbit compliance:

- [OWASP API Security Top 10](https://owasp.org/API-Security/editions/2023/en/0x11-t10/): authorization and resource consumption.
- [OWASP REST Security](https://cheatsheetseries.owasp.org/cheatsheets/REST_Security_Cheat_Sheet.html): validation, access control and safe errors.
- [OWASP HTTP Headers](https://cheatsheetseries.owasp.org/cheatsheets/HTTP_Headers_Cheat_Sheet.html): response-policy guidance.
- [OWASP CSRF Prevention](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html): conditional guidance for cookie-authenticated mutations.
- [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110) and [RFC 9111](https://www.rfc-editor.org/rfc/rfc9111): HTTP and cache semantics.
