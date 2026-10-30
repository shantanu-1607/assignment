# Public evidence ledger

## Method and scope

Public navigation completed October 31, 2026, approximately **00:47–00:52 IST**
(October 30, 19:17–19:22 UTC). Sources: verified-TLS curl and Chromium 154 through
the environment's configured proxy. The supplied proxy CA was imported into the
browser's NSS trust store; certificate verification was not disabled. Browser
prefetching generated ordinary Next.js RSC requests automatically.

No forms were submitted, accounts created, private resources accessed, endpoint
lists guessed, attacks crafted, or load tests run. No production errors were
intentionally triggered. Navigation included homepage, services, products,
contact, careers, privacy, blog, cookie policy, the linked backend internship,
and the blog's linked second page. Passive inspection included the exact public
JavaScript assets loaded for contact/products/blog.

[Structured observations](evidence/public-observations.json) preserve collection
times, selected response headers, empty form schemas, cookie attributes **without
values**, and client-asset hashes. Query strings, cookie values, authentication
headers and full HARs are deliberately excluded. This is a sanitized evidence
selection, not a complete traffic archive. Historic network failures are retained
separately in [research history](evidence/RESEARCH_HISTORY.md); they no longer block
this review and are not findings about the company.

## Company and domain selection

| ID | Source / method | Observation | Conclusion and limit |
|---|---|---|---|
| E-005 | [Homepage](https://www.cynbit.in/), GET + rendered page; 19:17:54 and 19:19:49 UTC | Title identifies Cynbit Technologies, AI automation/software solutions in Jaipur; links to `/services`, `/products`, `/careers` and LinkedIn `/company/cynbit-in/` | Strong public brand/domain match to the user-suggested site; the assessment email did not supply a URL, so employer ownership is not independently authenticated |
| E-006 | [Services](https://www.cynbit.in/services) and [products](https://www.cynbit.in/products), rendered navigation, 19:19:51–53 UTC | Advertises AI agents, backend/cloud development; names Cynbit Studio, TenantEase and Revio | Product context is observed marketing content, not verification of deployed features, uptime or claimed scale |
| E-007 | [Careers](https://www.cynbit.in/careers), 19:19:57 UTC; [linked internship](https://www.cynbit.in/careers/backend-developer-intern), 19:21–22 UTC | Backend Developer Intern role matches the assignment; role mentions REST/GraphQL, auth, databases and Go/Node/Python | Corroborates assessment relevance; does not prove those technologies are currently deployed |

A normal request to the homepage-linked LinkedIn profile returned status 999.
No attempt was made to bypass it. Public Google search did not yield usable
independent corroboration. The report reviews the matching publicly identified
site, with this identity limitation stated rather than hidden.

## Backend and control observations

| ID | Resource / method and timestamp UTC | What was observed | Reasonable conclusion / limit |
|---|---|---|---|
| E-008 | Seven main document responses and ordinary navigation, 19:19:49–19:20:02 | `/_next/static/...` assets; `text/x-component` fetches; `x-nextjs-prerender: 1` and `x-vercel-cache` on sampled public pages | Next.js/RSC frontend and Vercel caching indicators; no backend language, database or product-service topology established |
| E-009 | Home/services/products/contact/privacy document headers, same interval | 200, `public, max-age=0, must-revalidate`; ETags; Vercel HIT; homepage age 108453, contact age 3601 | Public caching is already present; Age describes cached response age, not database freshness or response latency |
| E-010 | Careers/blog document headers, 19:19:57 and 19:20:02 | 200, `private, no-cache, no-store, max-age=0, must-revalidate`, Vercel MISS, age 0 | Different cache treatment is observed; whether deliberate freshness or unnecessary dynamic rendering requires internal evidence |
| E-011 | Same seven document responses | HSTS `max-age=63072000`; no CSP, X-Frame-Options, X-Content-Type-Options or Referrer-Policy response header captured | Hardening opportunity scoped to sampled documents; no exploit or global absence established. HTML meta policies, browser defaults and product subdomains need separate review |
| E-012 | Contact and product-demo DOM, 19:19:53–55; [empty form screenshot](evidence/contact-form.png) | Required text/email/message fields; `maxLength: -1`; products also requires company; contact company/industry optional | Client has no HTML length cap; server-side size/schema validation remains unknown |
| E-013 | Exact browser-loaded contact/products/blog route bundles, fetched around 19:20 UTC; URLs and SHA-256 in structured evidence | References `submitInquiry`, `submitServiceLead`, `subscribeToNewsletter`; public client functions await these actions, branching on `success`; no local try/catch/finally around those awaits in inspected functions | Server Action integration, not a discovered REST API. Exception recovery gap is inferred from client code; no failed production submission was performed |
| E-014 | Linked internship page, ordinary GET 200, 19:21–22 UTC; privacy page | Application labels require DOB, gender (with “Prefer not to say”), permanent address, phone and WhatsApp in addition to career details. Privacy page broadly describes minimum name/email/professional-context collection | Specific data-minimization and purpose/retention review opportunity; no actual applicant data viewed, retention measured, or legal violation asserted |
| E-015 | Fresh browser context, after ordinary navigation; cookie policy read 19:20–21 UTC | `_ga`, `_ga_LK6CQWCWX0`, `_gcl_au`; Secure=false, HttpOnly=false, SameSite=Lax. No consent choice was made. Policy describes analytics and browser controls | Review analytics loading/consent and inventory; these are not proven authentication cookies. JS analytics cookies lacking HttpOnly is not itself an auth defect; no jurisdictional compliance conclusion |
| E-016 | Blog first page and linked `?page=2`, 200, 19:20–22 UTC | Page/category query links; second page shows a different article set. Captured same-site fetches in initial seven-page review were RSC, with zero non-RSC XHR/fetch JSON responses | Visible URL pagination exists; cannot infer SQL offset, database query efficiency, private API envelopes or maximum page enforcement |
| E-017 | [Privacy policy](https://www.cynbit.in/privacy), rendered text, 19:19:59 | Names Resend for email delivery and Prisma for database management | Company disclosure, not independently verified deployment architecture; Prisma is tooling, not identification of the database engine |

Selected header excerpt (homepage browser response):

```http
HTTP 200
content-type: text/html; charset=utf-8
cache-control: public, max-age=0, must-revalidate
etag: W/"babdce9795e4ef3d84461090987852d3"
age: 108453
strict-transport-security: max-age=63072000
vary: rsc, next-router-state-tree, next-router-prefetch, next-router-segment-prefetch
x-vercel-cache: HIT
x-nextjs-prerender: 1
```

Some public prerendered documents also returned `Access-Control-Allow-Origin: *`.
This was not observed on an authenticated API response; do not label public HTML
CORS as a credential-exposure vulnerability. `server: envoy` can be intermediary
metadata and is not used to infer the origin web server.

Client behavior excerpt, normalized from E-013 (action IDs deliberately omitted):

```js
// Contact: simplified control flow, not source formatting.
preventDefault(); setPending(true); clearError();
const result = await submitInquiry(new FormData(form));
setPending(false);
if (result.success) showSuccess(); else showGenericError();
```

An exception would bypass the normal reset in this function unless another layer
recovers it. That is a reason to add local recovery, not a claim that real users
have experienced a stuck form.
