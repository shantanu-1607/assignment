# Assessment submission

I reviewed Cynbit's public website using ordinary browser navigation, HTTP
responses, empty form schemas and publicly loaded client code. The audit records
Next.js/Vercel delivery behavior and identifies opportunities in form validation,
Server Action error recovery, response security headers, applicant data
minimization and analytics-cookie handling. Each finding includes evidence, its
limits and a practical recommendation. I also built a standalone Go API to
demonstrate strict JSON contracts, safe errors, request correlation, conditional
caching, bounded rate and concurrency limits, cancellation and graceful shutdown.
Formatting, vet, automated tests, race detection and real HTTP/shutdown checks
passed locally, with hosted CI results documented. The repository includes
sanitized evidence, architecture diagrams, implementation tradeoffs and a
finding-to-test matrix. The review remained passive: no forms were submitted or
private resources accessed. Internal authorization, retention and delivery
controls require a separate authorized review.
