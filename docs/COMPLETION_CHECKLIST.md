# Assessment completion checklist

This maps the uploaded assessment prompt to concrete deliverables. The technical
work is complete within the requested passive scope. Repository access must be
arranged before submitting a private repository URL to an external reviewer.

| Prompt requirement | Status | Evidence / deliverable |
|---|---|---|
| Inspect existing project and preserve files | Done | Initially empty checkout; incremental history; no pre-existing project files removed |
| Research company and public website | Done within public scope | `EVIDENCE.md` E-005–E-007: matching brand, services, products and Backend Developer Intern page; assessment-email domain was not supplied, so independent ownership authentication is limited |
| Passive technical observation | Done | `EVIDENCE.md` E-008–E-017: browser navigation, HTTP headers, empty form schemas, loaded client code, cookies without values and blog pagination |
| Timestamped, sanitized evidence | Done | `evidence/public-observations.json`, empty-form screenshot and cited source URLs |
| Backend audit with executive summary and architecture | Done | `AUDIT_REPORT.md`, `ARCHITECTURE.md` |
| Categorized security/performance findings | Done | `SECURITY_FINDINGS.md`, `PERFORMANCE_ANALYSIS.md`; observed, inferred and proposed claims distinguished |
| API/payload analysis | Done within observable scope | DOM input contracts, RSC/page behavior and Server Action references; no public business-JSON response occurred in captured navigation; no mutation was invoked to manufacture evidence |
| Concrete recommendations and snippets | Done | `RECOMMENDATIONS.md`: priorities, TypeScript exception-recovery example and links to executable Go controls |
| Functional Go reference backend | Done | `cmd/server`, `internal/api`, `internal/config`, `internal/service`; stateless and independent of production |
| HTTP contracts, validation and safe errors | Done | Versioned routes, exact JSON schema, size/range limits, request IDs, status codes and sanitized messages |
| Security/resource controls | Done | Response policy, peer quota, bounded key table, body/header budgets and active worker limit |
| Reliability/performance patterns | Done | Cancellation, application/server deadlines, caching, probes, structured logging and graceful shutdown |
| Meaningful automated tests | Done | 16 top-level tests, 39 pass events including subtests; validation/cache/rate/concurrency/deadline/cancellation/error tests |
| Architecture diagrams | Done | Mermaid observed/proposed and middleware flows in `ARCHITECTURE.md`; compact overview in README |
| Evidence-to-code/test traceability | Done | F-001–F-007 matrix in `RECOMMENDATIONS.md`; explicitly marks recommendations outside the stateless implementation |
| Engineering decisions and roadmap | Done | Five ADRs and impact/confidence/cost-based roadmap in `RECOMMENDATIONS.md` |
| Polished setup/review documentation | Done | README: five-minute path, setup, requests/responses, checks, configuration and limitations |
| Formatting, vet, tests and race checks | Passed | Local verification recorded in `VERIFICATION.md`; hosted CI also passed |
| Real startup, HTTP and shutdown checks | Passed | `examples/smoke.py`; local and hosted CI execution |
| Clean, meaningful history and publication | Done | Milestone commits on GitHub `main`; no co-author trailers added |
| Submission text, 100–150 words | Done | `SUBMISSION.md` |
| External reviewer access | Owner action required | Owner approved public visibility; GitHub rejected the integration’s PATCH with 403 (`Resource not accessible by integration`). Repository owner must change visibility in Settings → Danger Zone → Change repository visibility → Make public |

## Deliberate scope decisions

- The cloud checkout is `/workspace/assignment`. The prompt's local Desktop folder
  does not exist in this environment; the requester explicitly selected this
  GitHub repository. No duplicate project or worktree was created.
- Live deployment is optional and not needed for the repository-based submission.
  No public deployment is claimed. Local curl commands are demonstrations, not
  publicly reachable links.
- Benchmarks were optional. None were added because arithmetic throughput would
  not establish capacity for Cynbit's unknown workload.
- Authenticated product tests, internal validation, query plans, retention and
  delivery guarantees require a different authorized scope. They are limitations
  of a passive audit, not unfinished implementations in this assignment.
- No contact/application/newsletter form was submitted and no real applicant or
  private response data was collected. Recommendations do not assert exploited
  vulnerabilities or legal violations.

## Before portal submission

1. Make the repository public in GitHub Settings → Danger Zone → Change repository
   visibility → Make public. The owner already approved this; the platform
   integration lacks the administration permission needed to do it. A private
   URL alone is insufficient.
2. Submit `https://github.com/shantanu-1607/assignment` with `SUBMISSION.md` text.
3. If the portal accepts files, a source ZIP can also be generated with
   `git archive --format=zip --output=/path/outside/checkout/assignment.zip HEAD`.
4. The portal submission itself remains the requester's action; no assessment
   application, external message or form has been sent by this project.
