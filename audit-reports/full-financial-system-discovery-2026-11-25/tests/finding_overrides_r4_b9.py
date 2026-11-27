R4 = {}
R4['F-079'] = dict(
    design='Proposed: in frontend/apps/dashboard/src/components/hr/LeaveForm.tsx:114-115 send data.start_date and data.end_date exactly as the date input returns them (YYYY-MM-DD) instead of new Date(...).toISOString(); additionally make backend/internal/handler/http/hr_leave_handler.go:186 accept both YYYY-MM-DD and RFC3339 (normalise to a date in the tenant time zone, see F-070) so other clients are not affected; add a contract test that posts the exact form payload.',
    alternatives='A: fix the client and widen the API parser (recommended). B: client only (other API consumers still fail on RFC3339).',
    acceptance=['Submitting the leave form with valid dates returns 201 and the request appears in the list with the same dates.', 'The API accepts YYYY-MM-DD and RFC3339 and rejects other formats with a field error.', 'A form-contract test fails if the payload format changes.'],
    proposed_files=[], api=['POST/PUT /api/v1/hr/leaves: accepts both date formats (additive)'], policy='None.', migration='No.', historical='None.')
