# Phase 16 — TODO List & Sprint Plan

**Created:** 2026-03-13
**Total Tasks:** 40
**Sprints:** 4

---

## Sprint 1 — Quick Fixes & Label Changes (10 tasks)

These are small, independent fixes that can ship fast.

### S1-01: Remove "Active Licenses" from Software pages
- [ ] Edit `frontend/src/components/shared/ProgramCatalogPage.tsx` — remove the Active Licenses card (line ~114)
- [ ] Verify the card is gone on `/reseller/software`, `/manager/software`, `/manager-parent/software`, `/super-admin/software`

### S1-02: Scope "Licenses Sold" to current user
- [ ] Edit `backend/app/Http/Controllers/Reseller/SoftwareController.php` — filter `licenses_sold` to `WHERE reseller_id = current_reseller`
- [ ] Edit `backend/app/Http/Controllers/Manager/SoftwareController.php` — filter `licenses_sold` to `WHERE reseller_id IN (team_seller_ids)`
- [ ] Edit `backend/app/Http/Controllers/ManagerParent/ProgramController.php` — filter `licenses_sold` to `WHERE tenant_id = current_tenant`
- [ ] If the query comes from a shared controller or the Program model, check `backend/app/Models/Program.php` for `licenses_sold` accessor
- [ ] Verify each role sees their own count only

### S1-03: Rename team page buttons
- [ ] Edit `frontend/src/locales/en.json` — change `managerParent.pages.teamManagement.inviteReseller` to "Create Reseller"
- [ ] Edit `frontend/src/locales/ar.json` — change Arabic equivalent to "إنشاء موزع"
- [ ] Verify on `/manager/team` and `/manager-parent/team-management`

### S1-04: Rename "Active Customers" label in reports
- [ ] Edit `frontend/src/locales/en.json` — change `managerParent.pages.financialReports.activeCustomers` from "Active Customers" to "Current Active Customers"
- [ ] Edit `frontend/src/locales/ar.json` — Arabic: "العملاء النشطين حالياً"
- [ ] Verify on `/manager/reports` and `/manager-parent/reports`

### S1-05: Verify "Total Activations" count is correct
- [ ] Read `backend/app/Http/Controllers/ManagerParent/FinancialReportController.php` — confirm `total_activations` counts license rows, not customers
- [ ] Read `backend/app/Http/Controllers/Manager/ReportController.php` — same check
- [ ] Cross-check: compare the number shown on reports page with the actual license count in the database for the same date range

### S1-06: Remove "Mark as Resolved" from BIOS Conflicts
- [ ] Edit `frontend/src/pages/manager-parent/BiosConflicts.tsx` — remove the Resolve button and resolve dialog
- [ ] Edit `frontend/src/pages/manager/BiosConflicts.tsx` — same (if exists)
- [ ] Keep the data, just remove the UI action

### S1-07: Add "Reseller" column to BIOS Conflicts table
- [ ] Edit `backend/app/Http/Controllers/ManagerParent/BiosConflictController.php` — include reseller name in conflict response (from the second activation's reseller)
- [ ] Edit `backend/app/Http/Controllers/Manager/BiosConflictController.php` — same
- [ ] Edit `frontend/src/pages/manager-parent/BiosConflicts.tsx` — add Reseller column after "Affected Customers"
- [ ] Edit `frontend/src/pages/manager/BiosConflicts.tsx` — same

### S1-08: Fix BIOS Details reseller activation count
- [ ] Edit backend BIOS details controller — change reseller activation count from hardcoded/wrong value to `COUNT(*) WHERE bios_id = X AND reseller_id = Y`
- [ ] Check `backend/app/Http/Controllers/ManagerParent/BiosDetailsController.php`
- [ ] Check `backend/app/Http/Controllers/Manager/BiosDetailsController.php`
- [ ] Verify the count is correct on the BIOS Details Resellers tab

### S1-09: Make customer table headers sticky
- [ ] Edit `frontend/src/components/shared/DataTable.tsx` — add `sticky top-0 z-10 bg-white dark:bg-slate-950` to the `<thead>` or `<tr>` in the header
- [ ] Verify on `/manager/customers`, `/manager-parent/customers`, `/reseller/customers`
- [ ] Test with long customer lists (scroll down, headers should stay visible)

### S1-10: Fix status card counts with reseller filter
- [ ] Read `frontend/src/pages/manager-parent/Customers.tsx` — check if the count queries include `reseller_id` in their queryKey and queryFn
- [ ] Read `frontend/src/pages/manager/Customers.tsx` — same check
- [ ] If counts ignore the reseller filter, add `reseller_id` (and `program_id`) to the count query parameters
- [ ] The count queries should use the same filters as the main table query (except status)
- [ ] Verify: filter by reseller → cards should show counts for that reseller only

---

## Sprint 2 — Customer Page Filters & BIOS Details Redesign (10 tasks)

### S2-01: Add manager filter to Manager customer page
- [ ] Edit `frontend/src/pages/manager/Customers.tsx` — add a dropdown for filtering by manager (team member)
- [ ] Query team members using existing `managerService.getTeam()` or similar
- [ ] Pass `manager_id` param to the customer list API
- [ ] Edit `backend/app/Http/Controllers/Manager/CustomerController.php` — accept and apply `manager_id` filter

### S2-02: Add manager/manager-parent filter to Manager-Parent customer page
- [ ] Edit `frontend/src/pages/manager-parent/Customers.tsx` — add dropdown for filtering by manager or manager-parent
- [ ] Query team members (managers + resellers) using existing service
- [ ] Pass filter param to customer list API
- [ ] Edit `backend/app/Http/Controllers/ManagerParent/CustomerController.php` — accept and apply manager filter

### S2-03: Redesign BIOS Details page — Simplify Overview tab
- [ ] Edit `frontend/src/pages/manager-parent/BiosDetails.tsx` — restructure Overview:
  - Group 1: BIOS Info (BIOS ID, status, first/last activity)
  - Group 2: Current License (program, duration, price, dates)
  - Group 3: Customer & Reseller (names, emails)
  - Group 4: Sale Summary (total activations, avg days, revenue)
- [ ] Use clear section headers with dividers
- [ ] Remove duplicate/redundant fields
- [ ] Edit `frontend/src/pages/manager/BiosDetails.tsx` — same changes

### S2-04: Redesign BIOS Details page — Simplify Activity tab
- [ ] Edit `frontend/src/pages/manager-parent/BiosDetails.tsx` — Activity tab:
  - Group duplicate entries (e.g., `bios.activate` + `license.activated` for same event → show as one entry)
  - Add date/time column
  - Add reseller name column
  - Use color-coded action badges
- [ ] Edit `frontend/src/pages/manager/BiosDetails.tsx` — same

### S2-05: Redesign BIOS Details page — Fix Resellers tab
- [ ] Already fixed activation count in S1-08
- [ ] Add "Last Activity" date per reseller
- [ ] Add "Programs Sold" per reseller
- [ ] Make reseller name clickable (link to team member detail if available)

### S2-06: Add customer license history to Manager customer detail
- [ ] Create backend endpoint `GET /api/manager/customers/{id}/license-history`
- [ ] Edit `backend/app/Http/Controllers/Manager/CustomerController.php` — add `licenseHistory()` method
- [ ] Query: all licenses for this customer, ordered by `activated_at DESC`, include reseller name, program name, status, duration
- [ ] Edit `frontend/src/pages/manager/CustomerDetail.tsx` — add "License History" section/tab
- [ ] Show timeline: each entry = reseller name, program, start date, end date, status, duration

### S2-07: Add customer license history to Manager-Parent customer detail
- [ ] Same as S2-06 but for manager-parent:
- [ ] Edit `backend/app/Http/Controllers/ManagerParent/CustomerController.php` — add `licenseHistory()`
- [ ] Edit `frontend/src/pages/manager-parent/CustomerDetail.tsx` — add "License History" tab
- [ ] Add route in `backend/routes/api.php`

### S2-08: Add reseller timeline view to customer detail
- [ ] In the license history section (S2-06/S2-07), group entries by reseller
- [ ] Show: "Reseller: Ahmed Ali — 3 activations (Jan 2026 – Mar 2026)"
- [ ] Under each reseller, list individual activation/renewal entries

### S2-09: Add backend route for license history
- [ ] Edit `backend/routes/api.php` — add:
  - `GET /api/manager/customers/{customer}/license-history`
  - `GET /api/manager-parent/customers/{customer}/license-history`
- [ ] Verify routes are protected by correct role middleware

### S2-10: Add i18n translations for new labels
- [ ] Add English translations for:
  - `customerDetail.licenseHistory` = "License History"
  - `customerDetail.resellerTimeline` = "Reseller Timeline"
  - `customerDetail.period` = "Period"
  - `customerDetail.soldBy` = "Sold by"
  - `biosDetails.overview.biosInfo` = "BIOS Information"
  - `biosDetails.overview.currentLicense` = "Current License"
  - `biosDetails.overview.saleSummary` = "Sale Summary"
- [ ] Add Arabic translations for the same

---

## Sprint 3 — BIOS Change Request Workflow (10 tasks)

### S3-01: Create `bios_change_requests` migration
- [ ] Create `backend/database/migrations/2026_03_14_000000_create_bios_change_requests_table.php`
- [ ] Columns: id, license_id, reseller_id, old_bios_id, new_bios_id, reason, status (enum: pending/approved/rejected), reviewer_id (nullable), reviewer_notes (nullable), reviewed_at (nullable), timestamps
- [ ] Add foreign keys: license_id → licenses, reseller_id → users, reviewer_id → users
- [ ] Run migration

### S3-02: Create BiosChangeRequest model
- [ ] Create `backend/app/Models/BiosChangeRequest.php`
- [ ] Define relationships: `license()`, `reseller()`, `reviewer()`
- [ ] Add fillable fields, casts for status enum

### S3-03: Create reseller BIOS change request endpoint
- [ ] Create or edit `backend/app/Http/Controllers/Reseller/BiosChangeRequestController.php`
- [ ] `store()` — validate new_bios_id + reason, create request with status=pending
- [ ] `index()` — list reseller's own requests with status
- [ ] Add routes in `backend/routes/api.php`:
  - `POST /api/reseller/bios-change-requests`
  - `GET /api/reseller/bios-change-requests`

### S3-04: Create manager BIOS change request endpoints
- [ ] Create `backend/app/Http/Controllers/Manager/BiosChangeRequestController.php`
- [ ] `index()` — list all pending requests for the manager's tenant
- [ ] `approve()` — set status=approved, call external API to change BIOS, log activity
- [ ] `reject()` — set status=rejected, save reviewer_notes, log activity
- [ ] Add routes:
  - `GET /api/manager/bios-change-requests`
  - `PUT /api/manager/bios-change-requests/{id}/approve`
  - `PUT /api/manager/bios-change-requests/{id}/reject`
- [ ] Same for manager-parent routes

### S3-05: Implement external API call for BIOS change
- [ ] In the approve handler, call the external software API to change the BIOS ID
- [ ] Handle API failure: set status to `approved_pending_sync` or similar, retry later
- [ ] Update the license record: set `bios_id = new_bios_id`
- [ ] Log the change in activity_logs

### S3-06: Add "Request BIOS ID Change" action to reseller customer detail
- [ ] Edit `frontend/src/pages/reseller/CustomerDetail.tsx` — add action button in dropdown
- [ ] Create dialog component with fields: New BIOS ID (text input), Reason (textarea)
- [ ] On submit, call `POST /api/reseller/bios-change-requests`
- [ ] Show success/error toast

### S3-07: Create BIOS change requests page for manager
- [ ] Create `frontend/src/pages/manager/BiosChangeRequests.tsx`
- [ ] Table: Customer Name, Old BIOS, New BIOS, Reason, Reseller, Date, Status, Actions
- [ ] Filter: Status (pending/approved/rejected), date range
- [ ] Actions: Approve button, Reject button (opens dialog for rejection reason)
- [ ] Add route in `frontend/src/router/routes.ts`

### S3-08: Create BIOS change requests page for manager-parent
- [ ] Create `frontend/src/pages/manager-parent/BiosChangeRequests.tsx` — same as S3-07
- [ ] Add route and navigation link

### S3-09: Add navigation links and notification badges
- [ ] Edit manager sidebar — add "BIOS Change Requests" link under BIOS section
- [ ] Edit manager-parent sidebar — same
- [ ] Optional: Show badge with pending count on the nav item
- [ ] Edit `frontend/src/components/layout/ManagerLayout.tsx` and `ManagerParentLayout.tsx`

### S3-10: Add activity log entries for BIOS changes
- [ ] On request submit: log `bios.change_requested` with old/new BIOS in metadata
- [ ] On approve: log `bios.change_approved` with reviewer info
- [ ] On reject: log `bios.change_rejected` with rejection reason
- [ ] These should appear in customer activity tab and reseller logs

---

## Sprint 4 — Reseller Commission & Payment System (10 tasks)

### S4-01: Create commission database tables
- [ ] Create `backend/database/migrations/2026_03_15_000000_create_reseller_commissions_table.php`
  - Columns: id, reseller_id, manager_id, tenant_id, period (varchar: YYYY-MM), total_sales, commission_rate, commission_owed, amount_paid, outstanding, status (unpaid/partial/paid), notes, timestamps
- [ ] Create `backend/database/migrations/2026_03_15_000001_create_reseller_payments_table.php`
  - Columns: id, commission_id, reseller_id, manager_id, amount, payment_date, payment_method (enum: bank_transfer/cash/other), reference, notes, timestamps
- [ ] Run migrations

### S4-02: Create Commission and Payment models
- [ ] Create `backend/app/Models/ResellerCommission.php`
  - Relationships: `reseller()`, `manager()`, `payments()`
  - Scopes: `unpaid()`, `forPeriod()`, `forReseller()`
- [ ] Create `backend/app/Models/ResellerPayment.php`
  - Relationships: `commission()`, `reseller()`, `manager()`

### S4-03: Create reseller payment status endpoint
- [ ] Create `backend/app/Http/Controllers/Reseller/PaymentStatusController.php`
- [ ] `index()` — return current reseller's commission summary:
  - Total sales (from reports data)
  - Pending contracts count + value
  - Commission rate
  - Total owed
  - Total paid
  - Outstanding balance
  - Payment history (last 12 months)
- [ ] Add route: `GET /api/reseller/payment-status`

### S4-04: Create manager reseller payments endpoints
- [ ] Create `backend/app/Http/Controllers/Manager/ResellerPaymentController.php`
- [ ] `index()` — list all resellers with commission summary
- [ ] `show($resellerId)` — detailed payment history for one reseller
- [ ] `storePayment()` — record a payment (amount, date, method, notes)
- [ ] `updatePayment($paymentId)` — edit a payment record
- [ ] `updateCommission($commissionId)` — edit commission rate or owed amount
- [ ] Add routes:
  - `GET /api/manager/reseller-payments`
  - `GET /api/manager/reseller-payments/{resellerId}`
  - `POST /api/manager/reseller-payments`
  - `PUT /api/manager/reseller-payments/{paymentId}`
  - `PUT /api/manager/reseller-commissions/{commissionId}`
- [ ] Same for manager-parent routes

### S4-05: Create reseller payment status page (frontend)
- [ ] Create `frontend/src/pages/reseller/PaymentStatus.tsx`
- [ ] Layout:
  - Summary card: Total Sales, Commission Rate, Owed, Paid, Outstanding
  - Payment history table: Date, Amount, Method, Reference
  - Monthly breakdown chart (optional)
- [ ] Add route in `frontend/src/router/routes.ts`: `/reseller/payment-status`
- [ ] Add navigation link in reseller sidebar

### S4-06: Create manager reseller payments page (frontend)
- [ ] Create `frontend/src/pages/manager/ResellerPayments.tsx`
- [ ] Layout:
  - Summary cards: Total Owed, Total Paid, Total Outstanding
  - Reseller table: Name, Sales, Commission, Paid, Outstanding, Status, Actions
  - Filters: Period (month picker), Status (all/unpaid/partial/paid)
- [ ] Add route: `/manager/reseller-payments`
- [ ] Add navigation link in manager sidebar

### S4-07: Create manager reseller payment detail page
- [ ] Create `frontend/src/pages/manager/ResellerPaymentDetail.tsx`
- [ ] Layout:
  - Reseller info header
  - Monthly commission breakdown table
  - Payment history table
  - "Record Payment" button → dialog with amount, date, method, reference, notes
- [ ] Add route: `/manager/reseller-payments/:resellerId`

### S4-08: Create manager-parent payment pages
- [ ] Create `frontend/src/pages/manager-parent/ResellerPayments.tsx` — same as S4-06
- [ ] Create `frontend/src/pages/manager-parent/ResellerPaymentDetail.tsx` — same as S4-07
- [ ] Add routes and navigation links

### S4-09: Add commission service to frontend
- [ ] Add methods to `frontend/src/services/manager.service.ts`:
  - `getResellerPayments(filters)`
  - `getResellerPaymentDetail(resellerId)`
  - `recordPayment(payload)`
  - `updatePayment(paymentId, payload)`
  - `updateCommission(commissionId, payload)`
- [ ] Add methods to `frontend/src/services/reseller.service.ts`:
  - `getPaymentStatus()`
- [ ] Same for `manager-parent.service.ts`

### S4-10: Add types and translations
- [ ] Edit `frontend/src/types/manager-reseller.types.ts` — add:
  - `ResellerCommission` type
  - `ResellerPayment` type
  - `ResellerPaymentSummary` type
- [ ] Edit `frontend/src/locales/en.json` — add all commission/payment labels
- [ ] Edit `frontend/src/locales/ar.json` — Arabic translations

---

## File Reference — All Files to Edit or Create

### New Files to Create

| File | Sprint | Purpose |
|------|--------|---------|
| `backend/database/migrations/2026_03_14_000000_create_bios_change_requests_table.php` | S3 | BIOS change request table |
| `backend/database/migrations/2026_03_15_000000_create_reseller_commissions_table.php` | S4 | Commission tracking table |
| `backend/database/migrations/2026_03_15_000001_create_reseller_payments_table.php` | S4 | Payment records table |
| `backend/app/Models/BiosChangeRequest.php` | S3 | Model |
| `backend/app/Models/ResellerCommission.php` | S4 | Model |
| `backend/app/Models/ResellerPayment.php` | S4 | Model |
| `backend/app/Http/Controllers/Reseller/BiosChangeRequestController.php` | S3 | Reseller submits requests |
| `backend/app/Http/Controllers/Reseller/PaymentStatusController.php` | S4 | Reseller views payment status |
| `backend/app/Http/Controllers/Manager/BiosChangeRequestController.php` | S3 | Manager reviews requests |
| `backend/app/Http/Controllers/Manager/ResellerPaymentController.php` | S4 | Manager manages payments |
| `backend/app/Http/Controllers/ManagerParent/BiosChangeRequestController.php` | S3 | Manager-parent reviews requests |
| `backend/app/Http/Controllers/ManagerParent/ResellerPaymentController.php` | S4 | Manager-parent manages payments |
| `frontend/src/pages/manager/BiosChangeRequests.tsx` | S3 | BIOS change request queue |
| `frontend/src/pages/manager-parent/BiosChangeRequests.tsx` | S3 | Same for manager-parent |
| `frontend/src/pages/reseller/PaymentStatus.tsx` | S4 | Reseller payment view |
| `frontend/src/pages/manager/ResellerPayments.tsx` | S4 | Manager payment admin |
| `frontend/src/pages/manager/ResellerPaymentDetail.tsx` | S4 | Per-reseller payment detail |
| `frontend/src/pages/manager-parent/ResellerPayments.tsx` | S4 | Same for manager-parent |
| `frontend/src/pages/manager-parent/ResellerPaymentDetail.tsx` | S4 | Same for manager-parent |

### Existing Files to Edit

| File | Sprint | Change |
|------|--------|--------|
| `frontend/src/components/shared/ProgramCatalogPage.tsx` | S1 | Remove Active Licenses card |
| `frontend/src/components/shared/DataTable.tsx` | S1 | Add sticky headers |
| `frontend/src/pages/manager-parent/BiosConflicts.tsx` | S1 | Remove resolve action, add reseller column |
| `frontend/src/pages/manager/BiosConflicts.tsx` | S1 | Same |
| `frontend/src/pages/manager-parent/BiosDetails.tsx` | S2 | Redesign layout |
| `frontend/src/pages/manager/BiosDetails.tsx` | S2 | Same |
| `frontend/src/pages/manager/Customers.tsx` | S1+S2 | Fix card counts, add manager filter, sticky headers |
| `frontend/src/pages/manager-parent/Customers.tsx` | S1+S2 | Fix card counts, add manager filter, sticky headers |
| `frontend/src/pages/manager/CustomerDetail.tsx` | S2 | Add license history section |
| `frontend/src/pages/manager-parent/CustomerDetail.tsx` | S2 | Add license history section |
| `frontend/src/pages/reseller/CustomerDetail.tsx` | S3 | Add "Request BIOS ID Change" action |
| `frontend/src/locales/en.json` | S1+S2+S3+S4 | All new English labels |
| `frontend/src/locales/ar.json` | S1+S2+S3+S4 | All new Arabic labels |
| `frontend/src/router/routes.ts` | S3+S4 | New page routes |
| `frontend/src/services/manager.service.ts` | S3+S4 | New API methods |
| `frontend/src/services/manager-parent.service.ts` | S3+S4 | New API methods |
| `frontend/src/services/reseller.service.ts` | S3+S4 | New API methods |
| `frontend/src/types/manager-reseller.types.ts` | S3+S4 | New types |
| `frontend/src/types/manager-parent.types.ts` | S3+S4 | New types |
| `backend/routes/api.php` | S2+S3+S4 | New routes |
| `backend/app/Http/Controllers/Manager/CustomerController.php` | S2 | Add license history, accept manager filter |
| `backend/app/Http/Controllers/ManagerParent/CustomerController.php` | S2 | Same |
| `backend/app/Http/Controllers/ManagerParent/BiosConflictController.php` | S1 | Add reseller to response |
| `backend/app/Http/Controllers/Manager/BiosConflictController.php` | S1 | Same |
| `backend/app/Http/Controllers/ManagerParent/BiosDetailsController.php` | S1 | Fix activation count |
| `backend/app/Http/Controllers/Manager/BiosDetailsController.php` | S1 | Same |
| `backend/app/Http/Controllers/Reseller/SoftwareController.php` | S1 | Scope licenses_sold |
| `backend/app/Http/Controllers/Manager/SoftwareController.php` | S1 | Scope licenses_sold |
| `backend/app/Http/Controllers/ManagerParent/ProgramController.php` | S1 | Scope licenses_sold |
| Layout files (sidebar navigation) | S3+S4 | Add new nav links |

---

## Sprint Execution Order

| Sprint | Tasks | Estimated Time | Dependencies |
|--------|-------|---------------|--------------|
| **Sprint 1** | S1-01 to S1-10 | 4-6 hours | None |
| **Sprint 2** | S2-01 to S2-10 | 6-8 hours | S1 complete |
| **Sprint 3** | S3-01 to S3-10 | 8-10 hours | S1 complete |
| **Sprint 4** | S4-01 to S4-10 | 10-12 hours | S1 complete |

**Sprints 2, 3, and 4 can run in parallel after Sprint 1 is done.**

**Total estimated time: 28-36 hours**

