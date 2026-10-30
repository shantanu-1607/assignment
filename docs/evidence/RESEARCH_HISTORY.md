# Initial research access failures

This historical record preserves evidence IDs from before the environment's
network settings were activated. It has been superseded by the successful public
review in [EVIDENCE.md](../EVIDENCE.md). The failures below came from the local
outbound proxy and say nothing about Cynbit's availability or security.

| ID | Attempt | Result and limit |
|---|---|---|
| E-001 | HEAD to the initially guessed `https://cynbit.com/`, October 31, 2026, approximately 00:28 IST | Proxy CONNECT 403; no origin response. The `.com` domain was never established as associated with Cynbit |
| E-002 | Public Google search for Cynbit's official website in the same initial session | Proxy CONNECT 403; no search results received |
| E-003 | Ordinary GET to the candidate homepage after saving a network draft | Proxy CONNECT 403; draft persistence had not activated network access |
| E-004 | GET to the subsequently suggested `https://www.cynbit.in/`, October 31, 2026, 00:41:20 IST (October 30, 19:11:20 UTC) | Proxy CONNECT 403, curl exit 56; no origin data received |

Sanitized intermediary excerpt:

```text
HTTP/1.1 403 Forbidden
server: envoy
curl: (56) CONNECT tunnel failed, response 403
```

After environment publication, verified-TLS HTTP and browser access succeeded.
The current report is based on that successful review, E-005–E-017. No cookies,
credentials or personal data were retained from these failed attempts.
