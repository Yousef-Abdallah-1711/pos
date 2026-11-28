# RUN-004 designs, batch 5 (F-072, F-073: bank reconciliation).
R4 = {}

R4['F-072'] = dict(
    design='Proposed: rebuild reconciliation on ledger data. (1) Session creation reads the book balance from the GL for the bank account (general_ledger rows of bank_accounts_ledger.gl_account_id up to the statement date) and stores it; the client no longer supplies book_balance; (2) AddBankReconciliationLine (backend/internal/repository/accounting_advanced_repo.go:113-166) takes the journal line id, verifies it is a line on the bank GL account in the same tenant/restaurant, takes statement/book amounts from that line (the client supplies only the statement amount to match) and rejects a line already matched in another session; (3) refreshReconciliationSummaryTx (:2126-2149) computes difference = statement balance - (ledger balance adjusted by outstanding items) from those matched lines; (4) CompleteBankReconciliationSession (:169-230) completes only when the difference is zero (or an explicit adjustment journal is recorded) and otherwise returns 409 without setting completed_at; last_reconciled_* are written only on completion; (5) the dashboard shows the ledger balance next to the statement balance.',
    alternatives='A: ledger-derived book balance with line validation (recommended). B: keep manual entry but add a visible "unverified" flag (does not remove the fabricated-completion risk). C: automatic statement import and matching (larger, later phase).',
    acceptance=['A session with a line whose journal has no bank-account line is rejected (400).', 'Completing a session whose difference is not zero returns 409, leaves status in_progress and completed_at null.', 'A session matched to real ledger lines completes and sets last_reconciled_balance equal to the ledger balance at the statement date.', 'The book balance shown equals the GL balance of the bank account (see F-029).'],
    proposed_files=['backend/migrations/<next>_bank_reconciliation_line_constraints.sql (new)'],
    api=['POST /api/v1/accounting/bank-reconciliation/sessions: book_balance ignored (derived)', 'POST .../sessions/{id}/lines: journal line id required; 400 for invalid lines', 'POST .../sessions/{id}/complete: 409 when difference is not zero'],
    policy='Reconciliation tolerance and adjustment-journal rules are an accounting decision (proposed tolerance 0.01).',
    migration='Yes: unique matched journal-line per active session and a CHECK/foreign-key set tying lines to the session tenant; back-fill review of existing sessions.',
    historical='Existing completed sessions must be re-verified against the ledger; list sessions whose book balance differs from the GL at the statement date.',
    authz='Reconciliation completion should require a finance permission distinct from data entry.',
)

R4['F-073'] = dict(
    design='Proposed: in AddBankReconciliationLine (backend/internal/repository/accounting_advanced_repo.go:113-166) validate that the referenced journal entry belongs to the same tenant and restaurant (and, with F-072, touches the bank account); return 400/404 with a generic message; map foreign-key violations through the shared error mapper (F-025); enable row-level security on bank_reconciliation_sessions and bank_reconciliation_lines (pattern of migrations 164-171, see F-019).',
    alternatives='A: validate in the repository plus RLS (recommended). B: composite foreign key (tenant_id, journal_entry_id) to journal_entries (also strong, needs a unique index on (tenant_id, id)).',
    acceptance=['A journal id from another tenant returns 4xx and stores nothing.', 'A nonexistent id returns 4xx with no constraint name in the body.', 'RLS prevents a tenant GUC for tenant 4 from reading tenant 2 reconciliation rows.'],
    proposed_files=['backend/migrations/<next>_rls_bank_reconciliation.sql (new)'],
    api=['POST .../sessions/{id}/lines: 4xx for foreign or unknown journals'],
    policy='None.',
    migration='Yes: RLS policies (and optional composite foreign key).',
    historical='List existing lines whose journal belongs to a different tenant.',
    authz='Negative cross-tenant tests are mandatory.',
)
