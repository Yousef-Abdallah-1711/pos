# 22 — Current versus target gap matrix

| Area | Evidence | Target | Gate |
|---|---|---|---|
| Balance Sheet | FIN-001 failed historically; earnings omitted | approved earnings/equity bridge | release blocker |
| Auto-balance | FIN-003 mutates line | reject material mismatch | blocker |
| Journal inserts | FIN-008/009 blockers | JSON/schema compatibility | blocker |
| Paid order | FIN-010 suppressed errors at source | durable visible recovery | blocker |
| Money/tax | hard-coded 10%; policy unknown | versioned currency/tax | accountant |
| Channel lifecycle | E2E untested | source-to-ledger proof | E2E |
| Security | tenant matrix untested | negative access proof | security |
| Reliability | no failure/replay tests | idempotent recovery | reliability |
| Exports | PDF partial; CSV/XLSX unverified | format-specific parity | export |
| Scale | dev single samples | production baselines | capacity |

Target states are proposals; confidence is detailed in 21.
