# Phase 16 — Client Requirements

**Created:** 2026-03-13
**Status:** PLANNING

---

## Scope

This phase covers seven distinct client-requested features and fixes grouped into four sprints.

### Feature 1 — Reseller Commission & Payment System

**Problem:** No way for resellers to see how much they owe their manager, and no way for managers to record payments received from resellers.

**Solution:** Two new pages:

| Role | Page | Purpose |
|------|------|---------|
| Reseller | `/reseller/payment-status` | Read-only view of sales, commission owed, amount paid, outstanding balance |
| Manager / Manager-Parent | `/manager/reseller-payments` | List all resellers with owed/paid amounts, edit payment records, payment history |

**New database tables:**
- `reseller_commissions` — monthly commission snapshots per reseller
- `reseller_payments` — individual payment records (amount, date, method, notes)

**Backend endpoints:**
- `GET /api/reseller/payment-status` — reseller sees own status
- `GET /api/manager/reseller-payments` — manager sees all resellers
- `POST /api/manager/reseller-payments` — record a payment
- `PUT /api/manager/reseller-payments/{id}` — update a payment
- `GET /api/manager/reseller-payments/{resellerId}/history` — payment history

---

### Feature 2 — BIOS ID Change Request Workflow

**Problem:** When a customer needs to change their BIOS ID (hardware swap, etc.), there is no formal workflow. It happens outside the system.

**Solution:**

1. **Reseller** sees a "Request BIOS ID Change" action on the customer detail page
2. Opens a dialog: New BIOS ID + Reason for change
3. Request is submitted and appears in manager/manager-parent queue
4. **Manager/Manager-Parent** sees pending requests on a new page or notification badge
5. Can **Approve** (triggers real external API call to change BIOS in the software) or **Reject** (with rejection reason)
6. Result logged in customer activity/history

**New database table:**
- `bios_change_requests` — request_id, license_id, old_bios, new_bios, reason, status (pending/approved/rejected), reviewer_id, reviewer_notes, created_at

**Backend endpoints:**
- `POST /api/reseller/bios-change-requests` — submit request
- `GET /api/manager/bios-change-requests` — list pending requests
- `PUT /api/manager/bios-change-requests/{id}/approve` — approve and call external API
- `PUT /api/manager/bios-change-requests/{id}/reject` — reject with reason

---

### Feature 3 — Customer License History (Reseller Timeline)

**Problem:** On the customer detail page (`/manager/customers/44`), there is no way to see the full history of which resellers sold to this customer, what periods each reseller provided, and the full activation/renewal timeline.

**Solution:** Add a "License History" tab/section to the manager and manager-parent customer detail pages showing:

- Chronological list of every license activation, renewal, and expiry
- Which reseller performed each action
- Duration of each period (start → end)
- Visual timeline of all periods grouped by reseller

**Files to edit:**
- `frontend/src/pages/manager/CustomerDetail.tsx`
- `frontend/src/pages/manager-parent/CustomerDetail.tsx`
- `backend/app/Http/Controllers/Manager/CustomerController.php` — add `licenseHistory()` endpoint
- `backend/app/Http/Controllers/ManagerParent/CustomerController.php` — same

---

### Feature 4 — Team Page Label Fixes

**Problem:** On `/manager/team`, the buttons say "Invite Manager" / "Invite Reseller" but the client wants different labels.

**Solution:** Change button text to:
- "Invite Reseller" → "Create Reseller"
- Keep "Invite Manager" as-is, or change per client direction

**Files to edit:**
- `frontend/src/locales/en.json` — update translation keys
- `frontend/src/locales/ar.json` — update Arabic translations

---

### Feature 5 — Software Page Fixes

**Problem 5a:** "Licenses Sold" on the software page shows the global count, not the count for the logged-in user's scope.

**Problem 5b:** "Active Licenses" field should be removed entirely.

**Solution:**
- Backend: Scope `licenses_sold` to the current user (reseller sees own count, manager sees team count, etc.)
- Frontend: Remove the "Active Licenses" display from `ProgramCatalogPage.tsx`

**Files to edit:**
- `frontend/src/components/shared/ProgramCatalogPage.tsx` — remove Active Licenses card
- `backend/app/Http/Controllers/ManagerParent/ProgramController.php` — scope licenses_sold
- `backend/app/Http/Controllers/Manager/SoftwareController.php` — scope licenses_sold
- `backend/app/Http/Controllers/Reseller/SoftwareController.php` — scope licenses_sold

---

### Feature 6 — Customer Page Improvements

**Problem 6a:** Manager customers page has no manager filter. Manager-Parent has no manager/manager-parent filter. The table headers scroll out of view.

**Problem 6b:** Status cards show wrong counts when a reseller filter is active (cards show global, table shows filtered).

**Problem 6c:** Reports page label should say "Current Active Customers" not "Active Customers".

**Solution:**
- Add manager filter dropdown to manager customer page
- Add manager-parent and manager filter dropdowns to manager-parent customer page
- Make table headers sticky (`position: sticky; top: 0`)
- Fix status card counts to respect all active filters (reseller, program, etc.)
- Rename label in reports

**Files to edit:**
- `frontend/src/pages/manager/Customers.tsx` — add manager filter, sticky headers
- `frontend/src/pages/manager-parent/Customers.tsx` — add manager/manager-parent filter, sticky headers
- `frontend/src/locales/en.json` — rename "Active Customers" → "Current Active Customers"
- `frontend/src/locales/ar.json` — same in Arabic
- Backend customer controllers — accept `manager_id` filter param

---

### Feature 7 — BIOS Pages Cleanup

**Problem 7a:** BIOS Details page is too crowded and hard to understand.

**Problem 7b:** BIOS Conflicts page: remove "Mark as Resolved" action; add a "Reseller" column showing which reseller caused the conflict.

**Problem 7c:** BIOS Details reseller tab shows "Activations: 1" for every reseller instead of the real count.

**Solution:**
- BIOS Details: Simplify layout, use clearer section headers, reduce visual noise, group related info
- BIOS Conflicts: Remove resolve action, add reseller column from conflict metadata
- Fix activation count to aggregate properly per reseller per BIOS

**Files to edit:**
- `frontend/src/pages/manager-parent/BiosDetails.tsx` — redesign layout
- `frontend/src/pages/manager/BiosDetails.tsx` — same
- `frontend/src/pages/manager-parent/BiosConflicts.tsx` — remove resolve, add reseller column
- `frontend/src/pages/manager/BiosConflicts.tsx` — same
- Backend BIOS controllers — fix activation count aggregation, include reseller in conflict response

---

## Architecture Decisions

| Decision | Choice | Reason |
|----------|--------|--------|
| Commission calculation | Manual entry by manager, not auto-calculated | Client wants manager to control the numbers directly |
| BIOS change API | Real external API call on approve | Must actually update the software, not just the local DB |
| Payment history | Separate table, not inline | Allows multiple partial payments per period |
| Customer license history | Query from licenses + activity_logs | No new table needed, data already exists |
| Status card counts | Must include all active filters | Bug fix — cards currently ignore reseller/program filters |

---

## Risk Assessment

| Risk | Severity | Mitigation |
|------|----------|------------|
| External BIOS API failure on approve | HIGH | Queue the change, retry on failure, show status to user |
| Commission data integrity | MEDIUM | Audit log on every payment edit |
| Sticky headers browser compat | LOW | Standard CSS, works in all modern browsers |
| Migration on production | LOW | All new tables, no schema changes to existing tables |

