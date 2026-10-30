# Observed and proposed architecture

## Public request flow (E-008–E-017)

```mermaid
flowchart LR
  Browser[Public browser] --> Edge[Vercel cache indicators]
  Edge --> HTML[HTML and prerendered marketing content]
  Edge --> Dynamic[Careers and blog: no-store responses]
  Browser --> Assets[Next.js assets and RSC prefetch]
  Assets --> Forms[Contact / demo / newsletter client code]
  Forms -. referenced actions only .-> Actions[Server Actions boundary]
  Actions -. internals unverified .-> Unknown[Validation / storage / delivery unknown]
  Policy[Privacy-policy disclosure] -. disclosure only .-> Providers[Resend and Prisma]
```

Solid edges describe captured navigation/response behavior; dotted edges describe
client-code references or policy disclosures without backend execution proof.
No database engine, service topology or product authentication stack is inferred.

## Proposed application boundary

```mermaid
flowchart LR
  Caller --> Edge[HTTPS edge / reviewed browser policy]
  Edge --> Boundary[Size / quota / deadline / correlation]
  Boundary --> Validate[Schema and authorization where required]
  Validate --> Domain[Application operation]
  Domain --> Public[Explicitly public cacheable reads]
  Domain -. only for real delivery needs .-> Queue[Durable job acceptance]
  Queue -. proposed bounded worker .-> Email[Email provider]
  Boundary --> Telemetry[Sanitized logs / status and duration metrics]
```

This is a recommendation, not deployed Cynbit infrastructure. A durable queue is
conditional on real delivery requirements and is not added to the stateless demo.

## Implemented Go reference

```mermaid
flowchart TD
  R[Request] --> ID[Generate request ID / JSON security policy]
  ID --> Probe{GET health/readiness?}
  Probe -->|yes| Health[Health or readiness response]
  Probe -->|no| Ready{Accepting work?}
  Ready -->|no| Drain[503 draining]
  Ready -->|yes| Rate{Peer quota?}
  Rate -->|no| Quota[429 / Retry-After]
  Rate -->|yes| Slot{Worker capacity?}
  Slot -->|no| Busy[503 busy]
  Slot -->|yes| Budget[Context deadline / buffered small response]
  Budget --> Route[Versioned route / strict JSON / domain rules]
  Route --> Read[Public catalog / ETag]
  Route --> Compute[Context-aware calculator]
  Budget -->|deadline| Timeout[503 and cancel]
  Route --> Exit[Worker exit releases slot]
  Timeout -. eventual worker exit .-> Exit
```

Packages: `cmd/server` owns process/server lifecycle; `internal/config` validates
non-secret configuration; `internal/api` owns contracts and HTTP resource policy;
`internal/service` owns deterministic logic. Tests inject estimators to simulate
failure/cancellation locally. No test touches Cynbit.
