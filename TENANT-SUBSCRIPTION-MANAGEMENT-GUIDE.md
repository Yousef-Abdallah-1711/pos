# Tenant Subscription Management - Complete Guide

**Date**: 2026-02-08
**Status**: ✅ IMPLEMENTED, BUILD-VERIFIED AND READY TO USE
**Location**: Super Admin Dashboard → Tenants

---

## 📍 WHERE TO FIND IT

### Access Path:
1. Open Super Admin Dashboard: **http://localhost:3004/dashboard**
2. Navigate to: **Tenants** (left sidebar)
3. Click on any tenant name to view details
4. Scroll down to see **"Subscription & Trial Management"** section

### Ports:
- **Super Admin Dashboard**: http://localhost:3004
- **Regular Dashboard**: http://localhost:3000 (or 3002)
- **Restaurant Website**: http://localhost:3001
- **Landing Page**: http://localhost:3003
- **Backend API**: http://localhost:8080

---

## 🎮 WHAT YOU CAN DO

### 1. **Extend Trial Period**
**Button**: "Extend Trial" (Blue button)

- Add days to tenant's trial period
- Works for active trials only
- Cannot extend expired trials
- Example: Add 7 days to a tenant with 3 days left

**How to use**:
1. Click "Extend Trial" button
2. Enter number of days (1-365)
3. Enter reason (optional but recommended)
4. Click "Extend"
5. See confirmation message

### 2. **Disable Tenant Access**
**Button**: "Disable Tenant" (Red button - only if active)

- Immediately blocks tenant access
- Shows maintenance page on restaurant website
- Shows blocked overlay in dashboard
- Requires reason for audit trail

**How to use**:
1. Click "Disable Tenant" button
2. Enter reason (required): "Non-payment", "Trial expired", "Violates ToS", etc.
3. Click "Confirm"
4. Access immediately denied for that tenant

**Result**:
- Tenant's dashboard shows "Access Denied" overlay
- Restaurant website shows "Maintenance Page"
- All API requests return 403 Forbidden
- Change recorded in subscription history

### 3. **Enable Tenant Access**
**Button**: "Enable Tenant" (Green button - only if disabled)

- Re-enables access for disabled tenants
- Removes block without deleting data
- Immediate effect

**How to use**:
1. Click "Enable Tenant" button
2. Click "Confirm"
3. Tenant regains full access

### 4. **Delete Tenant**
**Button**: "Delete Tenant" (Gray button)

- Schedules tenant for deletion
- Data is soft-deleted (kept for 30 days)
- Can be canceled within grace period
- Requires confirmation

**How to use**:
1. Click "Delete Tenant" button
2. Enter reason (required)
3. Click "Schedule Deletion"
4. Tenant marked for deletion but data preserved for 30 days

---

## 📊 WHAT YOU CAN SEE

### Subscription Status Card

The page displays 4 key metrics:

#### 1. **Status Badge**
- **trial**: Active trial period (blue)
- **active**: Paid subscription (green)
- **expired**: Trial or subscription ended (red)
- **suspended**: Payment issues (yellow)
- **disabled**: Manually disabled (red)

#### 2. **Trial Days Remaining**
- Shows days left in trial
- Green if > 5 days
- Yellow if 2-5 days
- Red if < 2 days or expired

#### 3. **Trial End Date**
- When the trial period expires
- Format: MM/DD/YYYY
- Critical for monitoring

#### 4. **Account Status**
- **Active** (Green): Tenant can access
- **Disabled** (Red): Tenant is blocked

### Recent Subscription Changes
Shows last 5 actions taken:
- When action was performed
- Who performed it (admin ID)
- What changed (status, trial dates, days added)
- Reason provided

---

## 🧪 TESTING THE FEATURE

### Test Case 1: Extend Trial
**Objective**: Add days to a tenant's trial

```
1. Go to a tenant with trial status
2. Current days remaining: Note the number
3. Click "Extend Trial"
4. Enter 7 days
5. Enter reason: "Early adopter"
6. Click "Extend"
7. Expected: Days remaining increased by 7
8. Verify: Entry appears in "Recent Changes" section
```

### Test Case 2: Disable Tenant
**Objective**: Block a tenant's access

```
1. Go to an active tenant
2. Click "Disable Tenant"
3. Enter reason: "Non-payment"
4. Click "Confirm"
5. Expected: Button changes to "Enable Tenant"
6. Expected: Subscription status remains same but "Account Status" = "Disabled"
7. Verify: Entry appears in history with "TENANT_DISABLED" action
```

### Test Case 3: Dashboard Blocking
**Objective**: Verify disabled tenant sees blocked message

```
1. Disable a tenant (from Test Case 2)
2. Impersonate that tenant (click "Impersonate" in users list)
3. Open tenant's dashboard in new tab
4. Expected: Blocked overlay appears
5. Expected: Sidebar is disabled
6. Expected: "Contact Support" message shown
7. Note: Actual dashboard still loads but overlay blocks interaction
```

### Test Case 4: Restaurant Website Blocking
**Objective**: Verify maintenance page shows

```
1. From same disabled tenant, go to restaurant website
2. Visit: http://localhost:3001
3. Expected: Maintenance page appears
4. Expected: Professional "Under Maintenance" message shown
5. Expected: Expected return date might be displayed (if showing)
```

### Test Case 5: Re-enable Tenant
**Objective**: Restore access to disabled tenant

```
1. Go to the disabled tenant
2. Click "Enable Tenant"
3. Click "Confirm"
4. Expected: Button changes to "Disable Tenant"
5. Expected: Status shows "Active"
6. Expected: Entry appears in history with "TENANT_ENABLED" action
```

### Test Case 6: Subscription History
**Objective**: Verify all changes are logged

```
1. Complete Test Cases 1-5 on same tenant
2. Scroll to "Recent Subscription Changes"
3. Expected: See entries in reverse chronological order
4. Expected: Each entry shows:
   - Action type (TRIAL_EXTENDED, TENANT_DISABLED, etc.)
   - Old and new values
   - Reason provided
   - Timestamp
```

---

## 🔍 DEBUGGING TIPS

### Issue: "Cannot load subscription data"
**Solution**:
1. Check backend is running: `curl http://localhost:8080/api/v1/health`
2. Check database migrations: `go run cmd/migrate/main.go`
3. Check browser console for error messages
4. Try refresh button in subscription card

### Issue: Buttons are disabled
**Possible reasons**:
- Trial already expired (can't extend expired trials)
- Subscription data hasn't loaded yet
- Tenant has no subscription record

**Solution**: Click refresh button to reload data

### Issue: Changes don't appear in history
**Solution**:
1. Click refresh button
2. Wait 2-3 seconds
3. Check browser console for errors
4. Verify backend API is responding

### Issue: Tenant still has access after disabling
**Solution**:
1. This is a caching issue (30-second cache)
2. Wait 30 seconds or restart browser
3. Check that user is using same tenant context
4. Verify tenant ID matches

---

## 📱 API ENDPOINTS BEING USED

All operations are backed by these API endpoints:

| Action | Endpoint | Method |
|--------|----------|--------|
| Get subscription info | `/api/v1/super-admin/tenants/{id}/subscription` | GET |
| Extend trial | `/api/v1/super-admin/tenants/{id}/extend-trial` | POST |
| Disable tenant | `/api/v1/super-admin/tenants/{id}/disable` | POST |
| Enable tenant | `/api/v1/super-admin/tenants/{id}/enable` | POST |
| Delete tenant | `/api/v1/super-admin/tenants/{id}` | DELETE |
| Get history | `/api/v1/super-admin/tenants/{id}/subscription-history` | GET |

---

## 🎯 KEY FEATURES

### Automatic Data Refresh
- Data refreshes after each action
- Subscribe history updates automatically
- UI stays in sync with backend

### Real-Time Status
- Status badge shows current state
- Days counter updates immediately
- Account status reflects changes instantly

### Comprehensive Audit Trail
- Every action is logged
- Shows who, what, when, why
- Immutable history for compliance

### User-Friendly Modals
- Confirmation dialogs for destructive actions
- Clear reasons required
- Success/error notifications

### Responsive Design
- Works on desktop and tablet
- Mobile-friendly layout
- Touch-friendly buttons

---

## 📝 COMMON SCENARIOS

### Scenario 1: Free Trial Expiring Soon
**What you see**:
- Days remaining: 2
- Status badge: "trial" (blue)
- Trial end date shows tomorrow

**What you do**:
1. Review customer's usage
2. Click "Extend Trial"
3. Add 7-14 days
4. Reason: "Good customer, extending trial"
5. Customer gets more time

### Scenario 2: Customer Non-Payment
**What you see**:
- Customer hasn't paid
- Trial has expired
- Status: "expired"

**What you do**:
1. Click "Disable Tenant"
2. Reason: "Trial expired - no payment received"
3. Customer sees "Access Denied"
4. After payment received, click "Enable Tenant"

### Scenario 3: Offensive Content Violation
**What you see**:
- Customer created inappropriate content
- Need to block immediately

**What you do**:
1. Click "Disable Tenant"
2. Reason: "Violates community guidelines - offensive content"
3. Tenant access blocked immediately
4. Schedule deletion after review period

---

## 🔐 SECURITY FEATURES

✅ **Only Super Admin Can**:
- See subscription management
- Modify trial status
- Disable/enable tenants
- Delete tenants
- View audit trail

✅ **Every Action Is Logged**:
- Who made the change
- When it was made
- Why (reason provided)
- Old and new values

✅ **Soft Deletion**:
- Data kept for 30 days
- Can restore within grace period
- Audit trail permanent

✅ **Immediate Effect**:
- No waiting period for disabling
- Tenant blocked instantly
- No data loss (can re-enable)

---

## 🚀 PRODUCTION CHECKLIST

Before using in production:

- [ ] Test all 6 test cases above
- [ ] Verify backend is running
- [ ] Verify database migrations applied
- [ ] Test with real tenant data
- [ ] Test blocking works on both dashboard and restaurant
- [ ] Verify audit trail is comprehensive
- [ ] Train team on how to use
- [ ] Set up monitoring for subscription events
- [ ] Create runbook for common issues
- [ ] Test rollback procedure (enable after disable)

---

## 📚 ADDITIONAL RESOURCES

- **Complete Reference**: `/docs-organized/SAAS-TRIAL-SUBSCRIPTION-COMPLETE-REFERENCE.md`
- **Quick Reference**: `/docs-organized/SAAS-TRIAL-SUBSCRIPTION-QUICK-REFERENCE.md`
- **API Documentation**: `/docs-organized/02-API-BACKEND/API_REFERENCE_V2.0.0.md`

---

## ✨ WHAT CHANGED

**New Components Added**:
- TenantStatusBadge - Shows subscription status with color-coded badges
- TrialDaysCounter - Shows days remaining with progress bar
- ExtendTrialModal - Dialog for extending trial (with days + reason input)
- DisableTenantModal - Dialog for disabling tenant (with reason input)
- EnableTenantModal - Dialog for enabling tenant (shows disable context)
- DeleteTenantModal - Dialog for deleting tenant (with confirmation)
- SubscriptionHistoryTable - Self-contained history table with pagination

**Updated Pages**:
- `/dashboard/tenants/[id]/page.tsx` - Full subscription management UI integrated

**New Sections in Tenant Detail**:
- "Subscription & Trial Management" card with 4 status indicators + 4 action buttons
- "Recent Subscription Changes" table with paginated history
- Disable reason display when tenant is disabled

**Bug Fixes Applied (2026-02-08)**:
- Fixed JSX syntax error (mismatched `</p>` tag → `</div>`)
- Fixed TypeScript type errors with API response handling (removed incorrect `.data` unwrapping)
- Fixed `SubscriptionStatus` type mismatch (string → proper union type)
- Fixed `TrialDaysCounter` props (now uses `trialEndDate`/`trialDays` instead of `daysRemaining`)
- Fixed `SubscriptionHistoryTable` integration (uses `tenantId` prop, self-fetching)
- Fixed all modal props (added `isOpen`, `tenantName`, `currentTrialEndDate`)
- Removed duplicate `disableTenant`/`enableTenant` API methods
- Cleaned up debug console.logs across all pages
- Fixed ESLint unescaped entities errors
- Fixed `StatCard` icon prop type (`ReactNode` → `ElementType`)
- Fixed API constructor env var type safety

**Updated API Layer**:
- `api.ts` - Removed duplicate methods, fixed env var default, clean type exports

---

**Status**: ✅ READY TO USE
**Build**: ✅ PASSES (verified 2026-02-08)
**Documentation**: ✅ COMPLETE

Start managing tenant subscriptions now!
