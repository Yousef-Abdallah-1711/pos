# RUN-004 designs, batch 6 (F-074, F-075).
R4 = {}

R4['F-074'] = dict(
    design='Proposed (root cause PROBABLE: first trace the submit handlers of the four forms): add one shared hook (proposed frontend/apps/dashboard/src/hooks/useFormSubmit.ts) that wraps the submit promise in try/catch, maps a server {error, field_errors} body to react-hook-form setError calls and shows a generic toast otherwise; replace the ad-hoc handlers in the customer, supplier, chart-of-accounts and contract forms (and CategoryFormPage.tsx:129) with it; combine with the backend error mapper (F-025) so messages are safe to display.',
    alternatives='A: shared hook + backend error contract (recommended). B: global axios interceptor toast (hides field context).',
    acceptance=['A 400 on each listed form shows a readable inline or toast message and no framework error overlay.', 'The form remains editable and the entered values are kept.', 'A unit test simulates a 400 and asserts the message is rendered.'],
    proposed_files=['frontend/apps/dashboard/src/hooks/useFormSubmit.ts (new)'],
    api=['error bodies gain field_errors (additive, see F-025)'],
    policy='None.',
    migration='No.',
    historical='None.',
)

R4['F-075'] = dict(
    design='Proposed (root cause PROBABLE: confirm the handler path first): validate account_code against the chart pattern (digits only, 3-10 characters, prefix must match an existing parent when parent_id is set), trim names and limit them to the column size (return 400 not 500), and set maxlength on the account form inputs; escape account fields in CSV/PDF exports; reject codes that would change report classification prefixes unexpectedly.',
    alternatives='A: server validation + input limits (recommended). B: UI-only limits.',
    acceptance=['A code containing letters or HTML characters returns 400 with a field error.', 'A 300-character name returns 400, not a database error.', 'Valid numeric codes are accepted and appear in reports.'],
    proposed_files=[],
    api=['POST/PUT /api/v1/accounting/accounts: 400 field errors'],
    policy='Account code format of the chart (owner/accountant decision).',
    migration='Optional CHECK on account_code format after listing non-conforming rows (tenant 4 ids 389-391 in the evidence).',
    historical='Existing accounts with non-numeric codes are listed.',
)
