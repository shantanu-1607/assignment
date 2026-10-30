# Verification record

Executed in the managed Linux amd64 environment on October 31, 2026 IST
(October 30 UTC). These are reference-backend results, not Cynbit measurements.

## Toolchain and setup

- Go 1.27.3 linux/amd64, official archive from `dl.google.com`.
- SHA-256 checked against HTTPS release metadata from `go.dev/dl/?mode=json`:
  `65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40`.
- No third-party Go modules. `go.mod` requires Go 1.24 or newer.
- Toolchain installed outside the checkout at `/workspace/.tools/go`.
- System `/usr/bin/go` is a board-game executable, not the compiler.
- Default home cache is read-only; GOCACHE/GOPATH use `/workspace/.cache`.

## Actual checks

| Check | Outcome |
|---|---|
| `gofmt` / Make formatting gate | Passed; no formatting differences |
| `go vet ./...` | Passed |
| `go test ./...` (JSON event capture) | Passed: 16 top-level tests; 39 test/subtest pass events; 0 failures/skips |
| `go test -race ./...` via `make check` | Passed; all three tested packages; server command package has no unit tests |
| `make build` | Passed; built `bin/server` (ignored) |
| `python3 examples/smoke.py` | Passed real HTTP readiness, catalog, conditional 304, estimate result 7, validation 422, JSON 404 |
| SIGTERM built server | Passed, graceful exit status 0 in smoke script |
| Saved installation script's retained-toolchain path | Executed successfully; version gate and build passed |

Initial failures were environmental: default Go cache was read-only, and the
restricted command sandbox denied the live HTTP test's listening socket. After
selecting writable caches and obtaining local socket execution permission, the
ordinary and race suites passed. No tests were disabled or skipped to achieve
this result. The initial run validated that instance; the later restored-machine check is
recorded separately below.

The GitHub workflow runs formatting, vet, race tests, build and smoke checks on
push/PR. Its hosted result is **unverified locally**; adding a workflow is not proof
that GitHub Actions ran it. No benchmark, load test, production authentication
check was performed. Public-site navigation was subsequently completed; see the
separate evidence ledger for its scope and limits.

## Review boundaries

Tracked files were checked for accidental credential files, private-key/token
patterns, whitespace errors and broken local documentation links. This manual and
pattern-based review is not an exhaustive secret scanner or security audit.
Generated binaries, environment files and logs are ignored. No co-author trailer
was added. See [EVIDENCE](EVIDENCE.md) for the completed passive review and its limits.

## Published-environment revalidation

After publication restored a new machine, `make check`, `make build`, and
`python3 examples/smoke.py` were rerun successfully on October 31, 2026 IST.
The check command reused valid cached Go test results where reported; the built
server smoke check executed again. Public browser inspection used Chromium 154
and the supplied proxy CA trusted in NSS, with certificate verification enabled.
No prior network changes were reapplied and no environment publication was
performed by the agent. Documentation links and whitespace were checked again.
The hosted GitHub Actions outcome remains unverified.
