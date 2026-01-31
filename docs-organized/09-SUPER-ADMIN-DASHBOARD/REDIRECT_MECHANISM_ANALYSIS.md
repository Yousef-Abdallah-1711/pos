# Dashboard Redirect Mechanisms - Complete Analysis

**Date:** 2026-01-31
**Analysis Type:** Code audit of all redirect logic

---

## All Redirect Points in Dashboard App

### 🔴 Most Likely Culprit for Our Issue

#### **DashboardAuthProvider** (PRIMARY)
- **File:** `frontend/apps/dashboard/src/components/providers/DashboardAuthProvider.tsx`
- **Line:** 199-221
- **Trigger:** User not authenticated AND `?debug=true` NOT in URL
- **Code:**
  ```typescript
  if (!finalIsAuthenticated) {
    // Double-check debug mode before redirecting
    const debugParam = urlParams.get('debug')
    if (debugParam === 'true') {
      return  // ✅ SKIP REDIRECT IF DEBUG MODE
    }
    // Redirect to login after 100ms
    setTimeout(() => {
      window.location.href = `${landingUrl}/auth/login`
    }, 100)
  }
  ```
- **Status:** ✅ Has safeguard for debug mode
- **Why It Might Not Work:** The debug parameter might not be reaching this component

---

### 🟡 Secondary Redirect Points

#### **Dashboard Page Guard**
- **File:** `frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx`
- **Line:** 21-24
- **Trigger:** Page mount if `!isAuthenticated`
- **Code:**
  ```typescript
  useEffect(() => {
    if (!isAuthenticated) {
      window.location.href = `${getConfig().landingUrl}/auth/login`
    }
  }, [isAuthenticated, user])
  ```
- **Status:** ❌ **NO debug mode check** - REDIRECTS REGARDLESS
- **Problem:** This might redirect BEFORE DashboardAuthProvider can process debug mode

---

#### **Locale Home Page** (`/[locale]/`)
- **File:** `frontend/apps/dashboard/src/app/[locale]/page.tsx`
- **Line:** 20-31
- **Trigger:** Accessing `/en/` path
- **Code:**
  ```typescript
  useEffect(() => {
    if (isAuthenticated) {
      router.replace(`/${locale}/dashboard`)  // Go to dashboard
    } else {
      window.location.href = `${getConfig().landingUrl}/auth/login`  // Go to login
    }
  }, [isAuthenticated, user, router, locale])
  ```
- **Status:** ❌ **NO debug mode check**
- **Why It Might Help:** User directly accesses `/en/` instead of `/en/dashboard`

---

#### **Auth Callback**
- **File:** `frontend/apps/dashboard/src/app/auth/callback/page.tsx` and `[locale]/auth/callback/page.tsx`
- **Trigger:** Exchange token and redirect
- **Code:**
  ```typescript
  if (token && userParam) {
    setAuth(user, token)
    window.location.href = `${getConfig().dashboardUrl}/${locale}/dashboard`
  } else {
    window.location.href = `${getConfig().landingUrl}/auth/login`
  }
  ```
- **Status:** ❌ **NO debug mode check**
- **Note:** This is used by auth callback, not our direct URL approach

---

### 🟢 Other Redirects (Not Related)

- Error boundary logout → `/login`
- Notification actions → Dynamic URL
- Root `/` → `/en` (server-side redirect)

---

## Critical Discovery: The Problem Path

### When user visits:
```
http://localhost:3002/en/dashboard?auth_token=JWT&auth_user=JSON&debug=true
```

### Current Route Execution:
```
1. Next.js routes to /[locale]/dashboard
   └─ Matches [locale]="en", path="dashboard"

2. Renders /app/[locale]/dashboard/layout.tsx
   └─ Wraps with DashboardAuthProvider
   └─ DashboardAuthProvider mounts

3. DashboardAuthProvider component function runs
   ├─ Checks for ?debug=true ✅
   ├─ If found, should return debug panel
   └─ If not found, check auth in useEffect

4. First useEffect (MOUNT EFFECT)
   ├─ Reads ?auth_token and ?auth_user from URL ✅
   ├─ Stores in localStorage ✅
   ├─ Cleans URL preserving ?debug=true ✅
   └─ Sets urlAuthProcessed=true

5. Second useEffect (AUTH CHECK)
   ├─ Checks debug mode again ✅
   ├─ If debug=true, return early ✅
   ├─ Otherwise check auth
   └─ Redirect if not authenticated

6. Page renders
   ├─ If debug=true → Debug panel ✅
   ├─ Else if !auth → Loading spinner then redirect
   └─ Else → Dashboard
```

### The Issue:
**The redirect is happening, which means:**
- Either `?debug=true` is NOT being detected
- Or `?auth_token` and `?auth_user` are NOT being received
- Or the component is not rendering the return statement

---

## Root Cause Hypotheses

### Hypothesis 1: URL Parameters Lost
The `?debug=true`, `?auth_token=...`, `?auth_user=...` parameters might be:
- Stripped by Next.js router
- Lost during redirect from port 3001 to 3002
- URL-encoded incorrectly
- Not matching URLSearchParams parsing

**Test:** Check what `window.location.search` actually contains when component mounts

### Hypothesis 2: Component Not Rendering Debug Panel
Even if parameters are present, the return statement might not execute if:
- Component crashes before reaching that code
- React error boundary catches and unmounts component
- Some state initialization prevents render

**Test:** Add early console.error logs at component top

### Hypothesis 3: Layout-Level Redirect
A parent layout might be redirecting BEFORE DashboardAuthProvider renders:
- `/[locale]/layout.tsx` - AuthProvider (only preferences, not auth)
- Some other layout or middleware

**Test:** Check if test page at `/[locale]/test-url` still redirects (it does!)

### Hypothesis 4: Zustand Initialization Causes Early Redirect
Zustand store initialization might trigger auth checks before our code runs:
- Store hydration reads `auth-storage` from localStorage
- If `auth-storage` is empty, `isAuthenticated = false`
- Some component checks `isAuthenticated` and redirects early

**Test:** Check if localStorage `auth-storage` exists before component runs

---

## The Dashboard Page Guard Issue

**CRITICAL FINDING:** The `dashboard/page.tsx` file has its own auth guard:

```typescript
useEffect(() => {
  if (!isAuthenticated) {
    window.location.href = `${getConfig().landingUrl}/auth/login`
  }
}, [isAuthenticated, user])
```

This has **NO debug mode check** and will redirect even if `?debug=true` is in the URL!

**This could be the actual culprit.** If this effect runs and `isAuthenticated` is false, it redirects immediately, **before the debug panel can even render.**

---

## How to Fix

### Option 1: Add Debug Mode Check to Dashboard Page
Add check for `?debug=true` before redirecting in dashboard/page.tsx:

```typescript
useEffect(() => {
  const debugMode = new URLSearchParams(window.location.search).get('debug') === 'true'

  if (!isAuthenticated && !debugMode) {
    window.location.href = `${getConfig().landingUrl}/auth/login`
  }
}, [isAuthenticated, user])
```

### Option 2: Remove Auth Check from Dashboard Page
Let DashboardAuthProvider handle all auth logic, remove the duplicate check

### Option 3: Use Query Parameter in All Redirects
Pass `?debug=true` through the redirect chain so it's preserved

---

## Files Needing Audit

1. ✅ `DashboardAuthProvider.tsx` - Has debug mode safeguard
2. ❌ `dashboard/page.tsx` - **MISSING debug mode check**
3. ✅ `dashboard/layout.tsx` - Wraps with DashboardAuthProvider
4. ✅ `[locale]/layout.tsx` - Only preferences, not auth
5. ⚠️ `[locale]/page.tsx` - Has auth redirect, needs debug check
6. ✅ `middleware.ts.disabled` - Disabled in Phase 1
7. ✅ `lib/config.ts` - Just returns config URLs

---

## Next Steps

1. **Check if dashboard/page.tsx is the culprit**
   - Add debug console log at top
   - Check if isAuthenticated is false
   - Confirm redirect happens there

2. **Add debug mode bypass to dashboard/page.tsx**
   - Check for `?debug=true` before redirecting

3. **Test the fix**
   - Visit: `http://localhost:3002/en/dashboard?auth_token=...&auth_user=...&debug=true`
   - Should see debug panel instead of redirect

4. **Verify end-to-end flow**
   - Button on super admin page
   - Token exchange on landing page
   - Redirect to dashboard
   - Debug panel appears
   - Can proceed to actual dashboard

---

## Summary

The **most likely culprit is `dashboard/page.tsx`** which redirects unauthenticated users WITHOUT checking for debug mode. This redirect happens at the page level, bypassing the DashboardAuthProvider's debug mode safeguards.

**The fix:** Add a single debug mode check to the dashboard page guard.
