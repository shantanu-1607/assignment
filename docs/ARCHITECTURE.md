# Architecture

## Observed research path

```mermaid
flowchart LR
  A[Research client] --> B[Managed outbound proxy]
  B --> C[CONNECT rejected: 403]
  C -. no origin observation .-> D[Candidate domain: official identity unverified]
```

No target backend, API, database, CDN or authentication provider can be inferred
from this path.

## Proposed standalone design

```mermaid
flowchart LR
  A[Local caller] --> B[Go HTTP server]
  B --> C[Response security policy / request ID]
  C --> D[Peer quota and bounded admission]
  D --> E[Context deadline / buffered response]
  E --> F[Versioned routes / strict input decoding]
  F --> G[Context-aware calculator]
  F --> H[Immutable public catalog / ETag]
  B --> I[JSON logs: ID and duration]
  J[Signal] --> K[Readiness false / drain then close]
  K --> B
```

For an actual public deployment, place an independently configured HTTPS edge
before B and configure trusted proxy behavior. This diagram claims no deployed
infrastructure.

## Error and middleware flow

```mermaid
flowchart TD
  R[Request] --> ID[Generate request ID and response policy]
  ID --> P{GET health/readiness?}
  P -->|yes| Probe[Return probe result]
  P -->|no| Ready{Accepting work?}
  Ready -->|no| Drain[503 no-store]
  Ready -->|yes| Rate{Peer quota?}
  Rate -->|no| Q[429 plus Retry-After]
  Rate -->|yes| Slot{Worker slot?}
  Slot -->|no| Busy[503 busy]
  Slot -->|yes| Time[Run handler under deadline]
  Time --> Decode[Method, media type, JSON, domain validation]
  Decode --> Result[Small buffered JSON response]
  Time -->|deadline| Deadline[503 timeout and cancel context]
  Result --> Release[Release slot after worker exits]
  Deadline -. worker exit .-> Release
```

## Package boundaries

- `cmd/server`: process lifecycle and HTTP server limits.
- `internal/config`: validated non-secret environment configuration.
- `internal/api`: HTTP contract, response policy, limiter and admission.
- `internal/service`: HTTP-independent deterministic domain logic.

Tests inject an estimator to demonstrate cancellation, failure hygiene and
concurrency admission without calling production systems.
