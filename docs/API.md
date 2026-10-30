# Reference API

Base URL for local development: `http://127.0.0.1:8080`. These routes belong to the
standalone demonstration, not Cynbit's website. All operations are public and
stateless. The API stores no personal data and makes no external requests.

## Routes

| Method | Path | Successful response | Cache policy |
|---|---|---|---|
| GET | `/healthz` | 200, `data.status: "ok"` | no-store |
| GET | `/readyz` | 200 while accepting work; 503 while draining | no-store |
| GET, HEAD | `/v1/services` | 200 catalog; 304 on matching validator; HEAD has no body | public, max-age=60, must-revalidate |
| POST | `/v1/estimates` | 200 calculation result | no-store |

GET health/readiness bypass quotas and admission limits. Other requests share the
per-peer quota. Catalog and estimate routes return 405 with `Allow` for unsupported
methods; unmatched routes (including unsupported probe methods) return 404.

Every normal response includes a generated `X-Request-ID`. Uncacheable JSON also
includes `request_id`; catalog JSON omits it to keep representation bytes stable.
This ID is correlation metadata, not a credential. Client-supplied IDs are not reused.

## Estimate contract

Send `Content-Type: application/json` and exactly one JSON object:

```json
{"service_id":"api-review","integrations":2}
```

| Field | Required | Constraint |
|---|---|---|
| `service_id` | Yes | `api-review` or `integration-review` |
| `integrations` | Yes | Integer from 0 through 20 |

Unknown, duplicate, null and missing fields are rejected. Trailing JSON is rejected.
The body limit is 4,096 bytes, including whitespace. No input coercion is applied.
The illustrative formula is base units (3 or 5) plus two units per integration.
A calculation creates no resource, so success returns 200 rather than 201.

```sh
curl -i -H 'Content-Type: application/json' \
  -d '{"service_id":"api-review","integrations":21}' \
  http://127.0.0.1:8080/v1/estimates
```

This returns HTTP 422:

```json
{
  "error": {
    "code": "validation_failed",
    "message": "service_id must be a catalog ID; integrations must be 0..20"
  },
  "request_id": "..."
}
```

## Errors

| HTTP | Code | Meaning |
|---|---|---|
| 400 | `invalid_json` | Invalid JSON shape, types, fields or trailing input |
| 404 | `not_found` | No matching route |
| 405 | `method_not_allowed` | Wrong method on catalog/estimate route |
| 413 | `body_too_large` | Body exceeds 4,096 bytes |
| 415 | `unsupported_media_type` | Missing/unsupported JSON content type |
| 422 | `validation_failed` | Unknown service or integration count outside 0..20 |
| 429 | `rate_limited` | Peer quota/key capacity exceeded; `Retry-After: 60` |
| 500 | `internal_error` | Sanitized unexpected service error or panic |
| 503 | `busy` | Active worker capacity reached |
| 503 | `not_ready` | Server draining |
| 503 | `timeout` | Application deadline exceeded or service reports cancellation |
| 503 | `unavailable` | Request ID generation unavailable |

Errors use `Cache-Control: no-store`. Earlier middleware decisions take precedence:
a quota-exhausted caller can receive 429 before route or body validation. A client
that disconnects/cancels may receive no JSON body; that is distinct from a normal
application error response. Server-level malformed HTTP/timeouts may be handled
by Go's HTTP server before API middleware runs.

## Conditional reads

```sh
curl -i http://127.0.0.1:8080/v1/services
curl -I http://127.0.0.1:8080/v1/services
```

Copy the actual quoted `ETag` into an `If-None-Match` header on a later GET. Matching
strong or weak validators, and `*` for the existing catalog, return 304 with no body.
The wildcard example is directly runnable:

```sh
curl -i -H 'If-None-Match: *' http://127.0.0.1:8080/v1/services
```

Catalog content changes when code changes and the process restarts. A 60-second
freshness window is a demo decision; a real product needs its own freshness budget.

## Configuration and resource budgets

| Variable | Default | Accepted values |
|---|---|---|
| `ADDR` | `127.0.0.1:8080` | host:port; choose an available address |
| `REQUEST_TIMEOUT` | `2s` | Go duration from 1ms through 30s |
| `RATE_PER_MINUTE` | `60` | Integer from 1 through 10000 |

```sh
ADDR=127.0.0.1:8090 REQUEST_TIMEOUT=1s go run ./cmd/server
```

No `.env` loader is included. Other reference limits are code defaults: 32 active
workers, 1,024 peer keys, 4 KiB body cap, 16 KiB configured header budget, 5s header
read timeout, 10s read timeout, and 60s idle timeout. Write timeout is the request
budget plus 10s; shutdown allows 10s to drain. The quota uses direct socket peers
and ignores forwarded headers. See [decisions](RECOMMENDATIONS.md) for deployment
tradeoffs and [tests](../internal/api/api_test.go) for the executable contract.
