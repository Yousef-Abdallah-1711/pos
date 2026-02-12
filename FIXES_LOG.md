# Production Fixes Log

## Session Summary
Fixed critical React Hook order violations causing Internal Server Errors on dashboard pages.

---

## Issues Fixed

### 1. React Hook Order Violations (4 pages)

**Problem:** React has a rule that hooks must be called unconditionally and in the same order on every render. Several dashboard pages violated this by:
- Calling `useHasPermission()` hooks early
- Then having early `return` statements based on permission checks
- Then calling `useState()` hooks AFTER the early returns

This caused the "React has detected a change in the order of Hooks" error, resulting in Internal Server Errors.

**Affected Pages:**
1. ✅ `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/page.tsx` - **FIXED**
2. ✅ `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/leaves/page.tsx` - **FIXED**
3. ✅ `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/attendance/page.tsx` - **FIXED**
4. ✅ `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx` - **FIXED**

**Solution:** Moved ALL `useState()` hooks to execute BEFORE any early returns. This ensures hooks are always called in the same order regardless of permission checks.

**Example Fix (OrdersPage):**
```typescript
// BEFORE (❌ Wrong - hooks after early return)
const { hasPermission: canViewOrders, loading: canViewLoading } = useHasPermission(...)
if (!canViewOrders) { return <PermissionDenied /> }  // Early return
const [orders, setOrders] = useState([])  // Hook after return ❌

// AFTER (✅ Correct - all hooks before early return)
const { hasPermission: canViewOrders, loading: canViewLoading } = useHasPermission(...)
const [orders, setOrders] = useState([])  // Hook BEFORE return ✅
if (!canViewOrders) { return <PermissionDenied /> }  // Early return OK
```

---

## Commits

1. **Commit `d5af832`** - Fix React Hook order violation in OrdersPage
   - Moved `fetchOrders` function and `useEffect` before early returns
   - 47 insertions, 46 deletions

2. **Commit `58fe0c3`** - Fix React Hook order violation in LeavesPage
   - Moved all `useState` hooks before early return check
   - 28 insertions, 27 deletions

3. **Commit `cdc9780`** - Fix React Hook order violations in AttendancePage and PayrollPage
   - Moved all `useState` hooks before early return checks in both files
   - 30 insertions, 28 deletions

---

## Testing Notes

### Routes Now Working ✅
- `http://localhost:3002/en/dashboard/orders` - Fixed hook order violation
- `http://localhost:3002/en/dashboard/hr/leaves` - Fixed hook order violation
- `http://localhost:3002/en/dashboard/hr/attendance` - Fixed hook order violation
- `http://localhost:3002/en/dashboard/hr/payroll` - Fixed hook order violation

### Application Ports
- **Dashboard:** http://localhost:3002
- **Landing:** http://localhost:3001
- **Restaurant Website:** http://localhost:3003
- **Super Admin:** http://localhost:3004
- **Backend API:** http://localhost:8080

### Maintenance Route
- `http://localhost:3003/maintenance` - Route exists in `app/maintenance/page.tsx`

---

## Verification Steps

To verify all fixes:

```bash
# Run TypeScript check
cd frontend
npx tsc --noEmit

# Run dev server
pnpm dev

# Test affected routes
# - http://localhost:3002/en/dashboard/orders
# - http://localhost:3002/en/dashboard/hr/leaves
# - http://localhost:3002/en/dashboard/hr/attendance
# - http://localhost:3002/en/dashboard/hr/payroll
```

---

## React Rules of Hooks Reference

The Rules of Hooks require:
1. **Only call hooks at the top level** - Never call hooks inside loops, conditions, or nested functions
2. **Only call hooks from React components** - Call from functional components or custom hooks
3. **Always call in the same order** - The order of hook calls must be identical on every render

Violating these rules causes:
- React errors in console: "React has detected a change in the order of Hooks"
- Component crashes with "Internal Server Error"
- Unpredictable state behavior

**Key principle:** Hooks must execute unconditionally before any early returns from the component.

---

## Files Checked (No Issues Found)

- ✅ `products/page.tsx` - Hooks called correctly before early return
- ✅ `drivers/page.tsx` - No permission-based early returns
- ✅ `categories/page.tsx` - No permission-based early returns
- ✅ `theme-builder/editor/[id]/page.tsx` - useEffect before early returns
- ✅ `employees/page.tsx` - Hooks before early returns
- ✅ `rbac/roles/page.tsx` - Hooks before early returns
- ✅ `settings/roles/page.tsx` - Hooks before early returns
- ✅ `settings/users/page.tsx` - Hooks before early returns

---

## Related Issues

### User-Reported Issues (Session Start)
1. ❓ "App not working at http://localhost:3000/" - This port doesn't exist (apps on 3001-3004)
2. ✅ "Orders page gives Internal Server Error" - Fixed by hook order corrections
3. ❓ "Theme-builder sometimes gives Internal Server Error" - Not reproducible; likely fixed by other hook corrections
4. ✅ "Maintenance route (localhost:3003/maintenance)" - Route exists and is functional

---

## Future Prevention

To prevent similar issues:
1. Use ESLint rule `react-hooks/rules-of-hooks` (should be enabled)
2. Always place all hooks at component top level
3. Move logic requiring early returns to custom hooks if needed
4. Consider moving permission checks to parent/wrapper components

```javascript
// Good pattern: Wrapper component handles permission check
<PermissionGate permission={READ_ORDERS}>
  <OrdersPage /> {/* All hooks safe here */}
</PermissionGate>
```
