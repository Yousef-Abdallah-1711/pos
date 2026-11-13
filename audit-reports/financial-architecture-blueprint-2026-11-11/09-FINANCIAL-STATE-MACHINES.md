# 09 — Financial state machines (PROPOSED)

| Aggregate | Allowed progression | Illegal or correction behavior |
|---|---|---|
| Order | draft → accepted → fulfilling → fulfilled | cancellation follows policy; paid does not mean posted |
| Payment | initiated → authorized → captured → settled | repeated capture rejected; failure terminal absent new attempt |
| Refund | requested → approved → provider_pending → succeeded/failed | cumulative amount <= capture less prior refunds |
| Journal | draft → posted or void-before-post | posted immutable; correction via linked reversal |
| Customer/supplier invoice | draft → approved → posted → partial/paid | correction note/void by policy, no posted edit |
| Settlement | imported → matched/exception → reconciled | discrepancy requires explicit adjustment |
| Fiscal period | open → closing → closed | reopen only authorized, reasoned, audited |
| Inventory adjustment | proposed → approved → posted | correction is opposite movement |

Each transition validates tenant/restaurant, role, version, amounts and period; uses idempotency and audit actor/time/reason; fails visibly and leaves retryable or terminal state. Provider actions are external. Concurrency uses version checks/locks and has not been runtime-tested.
