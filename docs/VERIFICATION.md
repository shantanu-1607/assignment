# Verification record

Reference-backend checks completed October 31, 2026 (IST), using **Go 1.27.3 on
Linux amd64**. These results describe the local implementation, not Cynbit's
production performance or security.

## Results

| Check | Result |
|---|---|
| Formatting (`gofmt` gate) | Passed |
| `go vet ./...` | Passed |
| `go test ./...` | 16 top-level tests; 39 pass events including subtests; 0 failures/skips |
| `go test -race ./...` | Passed for all three tested packages; server command has no unit tests |
| `make build` | Passed |
| Real HTTP smoke check | Passed readiness, catalog, conditional 304, result 7, validation 422 and 404 |
| SIGTERM shutdown | Passed with exit status 0 |
| Documentation review | Local links, whitespace, finding/test references and submission word count checked |
| Source/evidence review | Credential-pattern checks passed; cookie values and private data excluded |

The review is not an exhaustive secret scan or a security certification. No tests
were disabled to obtain passing results. No benchmark, production load test or
authenticated product test was performed.

## Hosted CI

GitHub [run 41115659636](https://github.com/shantanu-1607/assignment/actions/runs/41115659636)
completed successfully for commit `38aae4b1180d49c381c4dd2f3484b8b529a0a8c7`.
The workflow checked formatting, vet and race tests, then built and exercised the
server. Earlier code/evidence milestone [run 41115164304](https://github.com/shantanu-1607/assignment/actions/runs/41115164304)
also passed. These are recorded results for those commits; the workflow page shows
subsequent runs.

The workflow now uses the same `make smoke` entry point documented for reviewers.
It builds the binary and runs `examples/smoke.py`, which owns its local process on
port 18081 and stops that process after checking responses.

## Reproduce

```sh
make check
make smoke
```

Requires Go, a C compiler for race detection, Make and Python 3. The tests and
smoke script contact only local servers. The repository has no third-party Go
modules, external service dependency or secret requirement.

Public-site evidence is recorded separately in [EVIDENCE.md](EVIDENCE.md). The
browser review used Chromium 154 with the supplied proxy CA trusted in NSS;
TLS verification remained enabled. That review did not submit any forms.

## Cloud workspace notes

The cloud toolchain is installed outside the checkout. To activate it:

```sh
export PATH=/workspace/.tools/go/bin:$PATH
export GOCACHE=/workspace/.cache/go-build
export GOPATH=/workspace/.cache/gopath
cd /workspace/assignment
make check
make smoke
```

The official Go archive from `dl.google.com` was checked against HTTPS release
metadata from `go.dev/dl/?mode=json`. SHA-256:

```text
65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40
```

Initial environment failures involved a read-only default cache and restricted
local sockets. Selecting writable caches and permitting local socket execution
resolved them. After environment restoration, build and smoke checks passed again;
Go reused valid test-cache entries where reported. Live processes are restarted
when needed rather than assumed to survive restoration.

See the [submission checklist](COMPLETION_CHECKLIST.md#before-portal-submission)
for reviewer-access status.
