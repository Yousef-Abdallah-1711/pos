# RUN-004 designs, batch 4 (F-053..F-071). Policy IDs refer to 56 (POL-01..POL-15); POL-16+ are PROPOSED additions. F-052/F-040/F-041 keep their earlier hand-written designs.
R4 = {}

R4['F-053'] = dict(
    design='Proposed: strict clock-time handling in backend/internal/handler/http/hr_attendance_handler.go (CreateAttendance, parse block at :235-243): accept RFC3339 and HH:MM (combined with attendance_date in the tenant time zone), return 400 with a field error when a supplied value cannot be parsed (never ignore the error), require clock_out > clock_in, compute and return total hours and overtime server-side; add a CHECK clock_out > clock_in on the table.',
    alternatives='A: strict parse + 400 + computed hours (recommended). B: keep lenient parsing and warn (data stays wrong).',
    acceptance=['POST with clock_in "09:00" and clock_out "19:30" stores both times and 10.5 hours.', 'POST with clock_in "abc" returns 400 and stores nothing.', 'clock_out before clock_in returns 400.'],
    proposed_files=['backend/migrations/<next>_attendance_clock_check.sql (new)'],
    api=['POST /api/v1/hr/attendance: 400 field errors; response includes total_hours'],
    policy='Overtime rules (thresholds) are HR policy; not changed here.',
    migration='CHECK constraint after listing existing rows with NULL times and zero hours.',
    historical='Rows with NULL times and 0 hours (evidence) are listed for HR correction.',
)

R4['F-054'] = dict(
    design='Proposed: derive and enforce leave entitlement. In backend/internal/handler/http/hr_leave_handler.go (:207-210 create, :328-329 update, approve near :412) ignore client total_days and compute working days from start/end (excluding configured weekends/holidays); at approval compare with the remaining annual allowance and either reject, or approve the excess as unpaid leave according to POL-14, recording the unpaid portion; the leave balances report reads the same ledger (it currently clamps remaining to 0 when used exceeds the allowance, hiding the overdraw).',
    alternatives='A: server-derived days + balance check at approval (recommended). B: keep client days and only warn.',
    acceptance=['Requesting 3 calendar days stores total_days 3 regardless of the client value.', 'Approving 30 days against a remaining allowance of 21 follows POL-14 (reject, or 21 paid + 9 unpaid recorded).', 'The balances report shows the true remaining balance including negative values or the unpaid portion.'],
    proposed_files=[],
    api=['POST/PUT /api/v1/hr/leaves: total_days ignored on input (returned computed)'],
    policy='POL-14 (leave entitlement) in 56.',
    migration='Optional unpaid_days column.',
    historical='Existing over-allowance approvals (30 vs 21 in the evidence) are listed for HR.',
)

R4['F-055'] = dict(
    design='Proposed: bound and record trial extensions. In backend/internal/handler/http/subscription_handler.go ExtendTenantTrial (:128-138) enforce 1 <= days <= 90 (configurable), require a reason, and write a subscription history row (actor, days, previous end, new end); show the history in Super Admin.',
    alternatives='A: cap + reason + history (recommended). B: cap only.',
    acceptance=['days 100000 returns 400.', 'A valid extension records actor, reason and dates and appears in the tenant history.'],
    proposed_files=[],
    api=['POST /api/v1/super-admin/tenants/{id}/extend-trial: reason required; days capped'],
    policy='Maximum extension is a business decision (proposed 90).',
    migration='Possibly a subscription_history table if none exists.',
    historical='Tenants with trial ends beyond the cap are listed.',
)

R4['F-056'] = dict(
    design='Proposed: validate ordered date ranges. In backend/internal/handler/http/contract_handler.go (create :165-170, update :221-226) and the resignation create/update (backend/internal/repository/resignation_repo.go:117-143 and its handler) require end_date >= start_date and last_working_date >= resignation_date, returning 400; add CHECK constraints so direct writers are covered.',
    alternatives='A: handler validation + CHECK constraints (recommended). B: handler only.',
    acceptance=['A contract ending before its start returns 400; a valid one returns 201.', 'A resignation with last working day before the resignation date returns 400.', 'The CHECK constraints reject the same data written directly.'],
    proposed_files=['backend/migrations/<next>_hr_date_range_checks.sql (new)'],
    api=['contracts/resignations create and update: 400 for inverted ranges'],
    policy='None.',
    migration='CHECK constraints after listing and correcting existing inverted rows (contract 1 and resignation 1 in the evidence).',
    historical='Existing inverted rows are listed for HR correction.',
)

R4['F-057'] = dict(
    design='Proposed: cash-on-delivery settlement. Extend OrderUseCase.CompleteDelivery (backend/internal/usecase/order_usecase.go:433-447) so completing a COD delivery records the amount the driver collected (field, default order total) and creates a driver cash liability (Dr driver cash receivable / Cr revenue-and-VAT through the normal sales posting); add a driver settlement operation that records the cash handed in (Dr cash / Cr driver receivable) and reports variances per driver and day.',
    alternatives='A: collected amount at delivery + driver settlement document (recommended). B: mark paid at delivery without tracking the driver (cash risk invisible).',
    acceptance=['Completing a COD order records the collected amount and posts the sales journal once.', 'Settling a driver for a day moves exactly the handed-in cash and shows any shortage.', 'Unsettled driver cash is visible in a report.'],
    proposed_files=['backend/internal/handler/http/driver_settlement_handler.go (new)', 'backend/migrations/<next>_driver_cash_settlements.sql (new)'],
    api=['POST /api/v1/admin/orders/{id}/complete-delivery: accepts collected_amount (additive)', 'POST /api/v1/admin/drivers/{id}/settlements (new)'],
    policy='Proposed POL-18 (driver cash handling and shortage treatment) to be added to 56.',
    migration='Yes: driver settlement tables and a driver receivable account mapping.',
    historical='Existing delivered-unpaid COD orders are listed for reconciliation with drivers.',
)

R4['F-058'] = dict(
    design='Proposed: validate driver phone numbers in backend/internal/handler/http/admin_driver_handler.go (:146-150) with the E.164 validator used by registration, plus email format and licence number length; return 400 with field errors.',
    alternatives='A: shared validator (recommended). B: client-side only.',
    acceptance=['Phone "abc" returns 400.', 'A valid E.164 number is stored normalised.'],
    proposed_files=[],
    api=['POST /api/v1/admin/drivers: 400 for invalid phone'],
    policy='None.',
    migration='No.',
    historical='Existing invalid phones are listed.',
)

R4['F-059'] = dict(
    design='Proposed: schema-validate theme JSON. In the theme save path (backend/internal/handler/http/admin_theme_handler_v2.go:293 and the theme service) validate every colour against a hex/rgb(a)/hsl pattern, fonts against an allow-list, URLs against http/https and the upload domain only (reject javascript:, data: and url() with other schemes), sizes within ranges, and reject unknown keys; sanitise again when rendering on the storefront.',
    alternatives='A: strict schema on write + defensive render (recommended). B: render-side escaping only (bad data remains).',
    acceptance=['A colour "not-a-color" returns 400.', 'A value containing url(javascript:...) returns 400 and never reaches the public theme endpoint.', 'Valid themes round-trip unchanged.'],
    proposed_files=['backend/internal/service/theme_validator.go (new)'],
    api=['POST/PUT /api/v1/admin/themes: 400 with field errors'],
    policy='None.',
    migration='No; existing themes are scanned and reported.',
    historical='Existing themes with invalid values (evidence) are listed and fixed by the tenant.',
)

R4['F-060'] = dict(
    design='Proposed: implement theme history and protect the active theme. Replace the 501 stub at backend/internal/handler/http/admin_theme_handler.go:68 with a version list (theme_versions table: id, theme_id, config, saved_by, saved_at) and a restore action; in DeleteTheme (admin_theme_handler_v2.go:371) refuse to delete the active theme (409) or require choosing a replacement.',
    alternatives='A: versions + guard (recommended). B: guard only.',
    acceptance=['Each save creates a version; restoring version N reproduces its config.', 'Deleting the active theme returns 409.'],
    proposed_files=['backend/migrations/<next>_theme_versions.sql (new)'],
    api=['GET /api/v1/admin/themes/{id}/history (implemented); POST .../restore (new)'],
    policy='None.',
    migration='Yes: theme_versions table.',
    historical='None.',
)

R4['F-061'] = dict(
    design='Proposed: complete zero-value events. In backend/internal/usecase/auto_posting_engine.go (the two-line check at :535-537 and PostSalesEntry/PostCOGSEntry) return a distinct "nothing to post" result when the total is zero and mark the outbox event completed with a note, instead of building an empty journal that can never post and retries forever; keep rejecting genuinely unbalanced journals.',
    alternatives='A: short-circuit zero-total events (recommended). B: forbid zero-total paid orders (policy decision on free items, see F-062 price validation).',
    acceptance=['Paying a zero-total order completes its outbox event with status completed and no journal.', 'The outbox shows no retryable entry for it.', 'A non-zero order still posts exactly one journal.'],
    proposed_files=[],
    api=['unchanged'],
    policy='Whether free items are allowed at all belongs with the product validation decision (F-062).',
    migration='No; stuck events are completed by a one-off script after review.',
    historical='Stuck zero-total events (retryable:1 in tenant 4) are listed and closed.',
)

R4['F-062'] = dict(
    design='Proposed: restore and extend input validation. In backend/internal/usecase/product_usecase.go validateProductRequest (:294-306, all checks commented out) re-enable and extend: name required and <= 100 chars, price > 0 for sellable products and <= a sane maximum with 2 decimals, cost >= 0, stock >= 0, discount <= price, SKU unique per restaurant (return 409, not 500); apply equivalent validation to categories (unique name), customers/suppliers (email/phone format, credit limit >= 0). Add table-driven tests for every rule.',
    alternatives='A: handler/usecase validation + DB CHECK constraints for the numeric rules (recommended). B: usecase validation only.',
    acceptance=['Each boundary from file 54 (negative price, zero price, cost above price allowed or not per decision, negative stock, empty or over-long name, duplicate SKU) returns the documented 4xx and creates nothing.', 'Valid payloads behave as before.'],
    proposed_files=['backend/migrations/<next>_product_check_constraints.sql (new)'],
    api=['POST/PUT /api/v1/products, categories, customers, suppliers: 400/409 with field errors'],
    policy='Whether a zero price (free item) is allowed and whether cost may exceed price are product decisions.',
    migration='CHECK constraints after listing violating rows.',
    historical='Existing invalid products (negative or zero price, empty name) are listed for correction.',
)

R4['F-063'] = dict(
    design='Proposed: make journal CSV import safe. In backend/internal/handler/http/accounting_advanced_handler.go ImportJournalEntriesCSV (:1067-1140) (1) require a client batch id (or hash of the file) with a unique constraint so a re-upload is rejected as a duplicate; (2) validate the whole file first and return 422 (or 207 with per-row errors) when any entry fails instead of 200 success:true; (3) process large files as a background job returning 202 plus a status endpoint, because imports of 1,500 journals take about 60 s and 6,000 outlive the HTTP connection; (4) wrap each file in one transaction when all-or-nothing is chosen.',
    alternatives='A: batch id + validation response + background jobs (recommended). B: only fix the response code.',
    acceptance=['Re-uploading the same file creates no duplicate journals.', 'A file with one bad entry returns a per-row error report and the documented all-or-nothing behaviour.', 'A 27,000-journal file returns 202 within 2 s and completes in the background with a progress endpoint.'],
    proposed_files=['backend/internal/jobs/journal_import_job.go (new)', 'backend/migrations/<next>_journal_import_batches.sql (new)'],
    api=['POST /api/v1/accounting/csv/journal-entries/import: 202 + job id; GET .../import/{job}'],
    policy='All-or-nothing versus partial import is an accounting decision.',
    migration='Yes: import batches table.',
    historical='Duplicate journals created by repeated imports are listed.',
)

R4['F-064'] = dict(
    design='Proposed (scope reduced in RUN-004: uniqueness is already enforced by uq_restaurant_tables_number and mapped to 409): validate the table name length (<= 100) and location in backend/internal/handler/http/table_handler.go before the insert and map pq 22001 (value too long) to 400 via the shared error mapper (F-025).',
    alternatives='A: handler validation + shared mapper (recommended).',
    acceptance=['A 1,800-character name returns 400 with a field error.', 'Duplicate table numbers still return 409.'],
    proposed_files=[],
    api=['POST /api/v1/tables: 400 for over-long names'],
    policy='None.',
    migration='No.',
    historical='None.',
)

R4['F-065'] = dict(
    design='Proposed: convert numeric form fields before submit. In frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx (:27 state, :129 payload) convert display_order with the same toInt helper the product form uses (components/products/ProductForm.tsx:100), move toInt to a shared lib, and show the server field error message instead of a generic toast; add a component test that submits a category with a display order and asserts the request body type.',
    alternatives='A: shared converter + field errors (recommended). B: make the backend accept numeric strings (loosens the API contract).',
    acceptance=['Saving a category with display order 7 returns 201 and stores 7.', 'A server validation error is shown next to the field.'],
    proposed_files=['frontend/apps/dashboard/src/lib/formNumbers.ts (new)'],
    api=['unchanged'],
    policy='None.',
    migration='No.',
    historical='None.',
)

R4['F-066'] = dict(
    design='Proposed: read the product total from the API contract. In frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx (extractTotalFromCollection :173-178, call :329) read pagination.total (fall back to total, then items.length only when no total exists) and add a unit test with the real /products response shape; consider a lightweight GET /products/count endpoint so the home page does not fetch rows.',
    alternatives='A: fix the extractor (recommended, smallest). B: add a count endpoint (cleaner long term).',
    acceptance=['With 12 non-deleted products the tile shows 12.', 'The unit test fails if the extractor reads the wrong key.'],
    proposed_files=[],
    api=['optional GET /api/v1/products/count (new)'],
    policy='None.',
    migration='No.',
    historical='None.',
)

R4['F-067'] = dict(
    design='Proposed: sort the lowest-margin list by margin. In backend/internal/handler/http/report_handler.go replace the backwards iteration at :2470-2476 with a query (or an explicit sort) by gross_margin_pct ASC then revenue DESC, keeping the top list ordered by gross profit (:2432); label both lists with their sort key in the UI and exports.',
    alternatives='A: separate sort for the lowest list (recommended). B: rename the section "Lowest profit" (changes the business meaning).',
    acceptance=['For items with margins 0%, 60%, 60%, 60%, 100% the lowest list is ordered 0%, 60%, 60%, 60%, 100%.', 'The top list stays ordered by gross profit.', 'CSV and PDF exports match the screen.'],
    proposed_files=[],
    api=['response shape unchanged; ordering of lowest_margin_items changes'],
    policy='None.',
    migration='No.',
    historical='None.',
)

R4['F-068'] = dict(
    design='Proposed: one row per salary record in the social-insurance export. In backend/internal/repository/accounting_advanced_repo.go GetSocialInsuranceReportRows (:1421-1446) replace the LEFT JOIN insurance_tracking at :1437 with a LATERAL subquery that returns one active policy per employee (latest start) or aggregates policy numbers with string_agg, so salary and deductions appear once per employee; add a test with two active policies.',
    alternatives='A: LATERAL one-row-per-employee (recommended). B: DISTINCT ON (employee) in the outer query.',
    acceptance=['With two active policies for one employee the CSV has exactly one data row for that employee.', 'Totals of base_salary and deductions over the file equal the salaries table.'],
    proposed_files=[],
    api=['CSV columns unchanged'],
    policy='Which policy is reported when several are active is an HR/legal decision (proposed: list all numbers in one cell).',
    migration='No.',
    historical='Previously exported files may contain duplicates; list affected months.',
)

R4['F-069'] = dict(
    design='Proposed: terminate the employee when a resignation completes. In backend/internal/repository/resignation_repo.go Complete (:201-215) update employees in the same transaction: employment_status = terminated, termination_date = last_working_date, is_active = false from that date; make DeleteEmployee (employee_repo.go:370-380) use the contract/resignation date when one exists; exclude terminated employees from future payroll and attendance creation; the turnover report then counts the departure.',
    alternatives='A: complete resignation terminates employee (recommended). B: keep manual termination and warn when a completed resignation exists.',
    acceptance=['Completing a resignation with last working day 2026-11-25 sets termination_date 2026-11-25 and status terminated.', 'The turnover report shows departures 1 and the matching rate for the period.', 'No salary or attendance can be created for the employee after the termination date.'],
    proposed_files=[],
    api=['POST /api/v1/hr/resignations/{id}/complete: also updates the employee (documented side effect)'],
    policy='Notice period and final-settlement rules are HR policy (proposed POL-19).',
    migration='No schema change; data-correction for completed resignations whose employee is still active.',
    historical='Employees with completed resignations still active are listed for HR review.',
)

R4['F-070'] = dict(
    design='Proposed: introduce a restaurant time zone. Add a tenant/restaurant timezone setting (default Africa/Cairo) and use it for KDS "today" and busiest-hours buckets (backend/internal/repository/kds_repo.go:314-317 and :418-430: convert with AT TIME ZONE before CURRENT_DATE / EXTRACT(HOUR ...)); store timestamps as timestamptz going forward; audit other date-bounded reports (sales, dashboard, orders stats) for the same UTC cut-off (currently a CANDIDATE, not tested).',
    alternatives='A: time-zone setting + conversion in queries (recommended). B: convert in the frontend only (wrong boundaries in aggregates).',
    acceptance=['An order received at 23:30 local appears in the previous local day and in hour 23 (not hour 21).', 'KDS completed_today uses local midnight.', 'The audit of other reports lists every query that uses CURRENT_DATE or date_trunc on a UTC timestamp.'],
    proposed_files=['backend/migrations/<next>_restaurant_timezone.sql (new)'],
    api=['settings gain timezone (additive)'],
    policy='Default time zone per tenant.',
    migration='Yes: setting column; timestamps conversion plan.',
    historical='Historical aggregates are recomputed on read; no stored data changes.',
)

R4['F-071'] = dict(
    design='Proposed: invert the environment check. Replace config.IsProductionMode() (backend/internal/config/runtime_env.go:20-23) usage with an explicit IsDevelopmentMode() allow-list (development, dev, local, test) so that staging and any unknown value behave like production for: RBAC default and refusing RBAC_ENFORCEMENT=audit (middleware/permission.go:36-52), mock payment refusal (config/paymob_mode.go:29), Redis/session and rate-limit store requirements (cmd/api/main.go:148-160), secure cookies (handler/http/auth_handler.go:223) and the metrics token (cmd/api/health_routes.go:52); fail startup on an unrecognised ENV value; keep the independent-secret check (cmd/api/signing_secrets.go:9).',
    alternatives='A: development allow-list (recommended; unknown values become safe). B: add staging to IsProductionMode (fixes only the named value).',
    acceptance=['Starting with ENV=staging and no RBAC variable enforces RBAC: a cashier gets 403 on finance, HR and user endpoints (re-run tests/runtime/ph17_sec013.mjs).', 'ENV=staging with PAYMOB_MODE unset or mock fails to start.', '/health/metrics returns 401 without a token outside development.', 'An unrecognised ENV value stops the API with a clear message.'],
    proposed_files=[],
    api=['/health/metrics requires a token outside development'],
    policy='List of environments considered development (owner decision).',
    migration='No.',
    historical='Review staging data for exposure while RBAC was in audit mode.',
    authz='Core of the fix: enforcement mode and role matrix tests.',
)
