# Super Admin Dashboard - Debugging Guide

## Console Debugging Setup

This guide helps you debug and fix issues with the Super Admin Dashboard development server.

---

## Common Issues & Solutions

### Issue 1: Port 3004 Already in Use

**Error:**
```
Error: listen EADDRINUSE: address already in use :::3004
```

**Cause:** A previous dev server process is still running on port 3004.

**Solutions:**

#### Windows (Command Prompt as Admin):
```cmd
REM Find process using port 3004
netstat -ano | findstr :3004

REM Kill the process (replace PID with actual process ID)
taskkill /PID <PID> /F

REM Example:
taskkill /PID 33928 /F
```

#### Windows (PowerShell as Admin):
```powershell
# Find and kill process on port 3004
Get-NetTCPConnection -LocalPort 3004 | ForEach-Object {Stop-Process -Id $_.OwningProcess -Force}
```

#### Linux/Mac:
```bash
# Find process using port 3004
lsof -i :3004

# Kill the process
kill -9 <PID>
```

**Then restart:**
```bash
cd frontend/apps/super-admin
pnpm dev
```

---

### Issue 2: Tailwind Config Build Error

**Error:**
```
Error: Package subpath './tailwind.config' is not defined by "exports"
in @pos-saas/ui-themes/package.json
```

**Cause:** The tailwind.config.ts file is trying to import from a package that doesn't export the subpath.

**Solution:**

The file should NOT import from ui-themes. Verify `tailwind.config.ts` contains:

```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: 'rgb(var(--primary) / <alpha-value>)',
        secondary: 'rgb(var(--secondary) / <alpha-value>)',
        accent: 'rgb(var(--accent) / <alpha-value>)',
      },
    },
  },
  plugins: [],
};
```

**If still broken:**
```bash
cd frontend/apps/super-admin
rm -rf .next node_modules
cd ..
pnpm install --filter @pos-saas/super-admin
cd apps/super-admin
pnpm dev
```

---

### Issue 3: Cache-Related Errors

**Symptoms:** Changes don't take effect, old errors persist

**Solutions:**

```bash
cd frontend/apps/super-admin

# Clear all caches
rm -rf .next
rm -rf node_modules
rm -rf dist

# Reinstall and rebuild
cd ..
pnpm install --filter @pos-saas/super-admin
cd apps/super-admin
pnpm dev
```

---

### Issue 4: Module Not Found Errors

**Error:**
```
Cannot find module '@pos-saas/ui-components' or its corresponding type declarations
```

**Solution:**

```bash
# Ensure all workspace dependencies are installed
cd frontend
pnpm install

# Then restart super-admin
cd apps/super-admin
pnpm dev
```

---

## Console Debugging Tips

### 1. Monitor Build Output

```bash
cd frontend/apps/super-admin
pnpm dev 2>&1 | tee debug.log

# In another terminal, watch the log file
tail -f debug.log
```

### 2. Check Network Requests

Open browser DevTools (F12) → Network tab to see:
- Failed API calls (look for 500/404 errors)
- CSS loading issues
- JavaScript module errors

### 3. Check Browser Console

Open browser DevTools (F12) → Console tab to see:
- JavaScript errors
- Module resolution issues
- React/Next.js errors
- Network errors

### 4. Verify API Connection

```bash
# Test backend connectivity
curl -X POST http://localhost:8080/api/v1/super-admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"superadmin@pos-saas.com","password":"SuperAdmin@123"}'
```

Expected response:
```json
{
  "success": true,
  "data": {
    "token": "...",
    "super_admin": {...}
  }
}
```

---

## Step-by-Step Debugging Process

### Step 1: Kill Existing Processes

```bash
# Windows
taskkill /IM node.exe /F

# Linux/Mac
killall -9 node
```

### Step 2: Clean Installation

```bash
cd C:\Users\OPT\Desktop\POS\frontend

# Remove all caches
rm -rf node_modules
rm -rf apps/super-admin/.next
rm -rf apps/super-admin/node_modules
rm -rf pnpm-lock.yaml

# Reinstall everything
pnpm install

# Verify installation
pnpm --filter @pos-saas/super-admin list
```

### Step 3: Start Single Dev Server

```bash
# Start only super-admin, not all apps
cd apps/super-admin
pnpm dev

# You should see:
# ✓ Starting...
# ▲ Next.js 15.5.9
# - Local:        http://localhost:3004
```

### Step 4: Test in Browser

1. Open http://localhost:3004 in browser
2. Open DevTools (F12)
3. Go to Console tab
4. Look for any errors
5. Try logging in with credentials:
   - Email: superadmin@pos-saas.com
   - Password: SuperAdmin@123

### Step 5: Check Build Output

Look for these success messages:
```
✓ Starting...
▲ Next.js 15.5.9
- Local:        http://localhost:3004
- Network:      http://192.168.x.x:3004
```

---

## Debugging Checklist

- [ ] Port 3004 is free (no other process using it)
- [ ] No `.next` cache files interfering
- [ ] `node_modules` is properly installed
- [ ] `tailwind.config.ts` doesn't import from ui-themes
- [ ] Backend is running on port 8080
- [ ] Environment variables in `.env.local` are correct
- [ ] Browser console has no errors
- [ ] Network tab shows successful API responses
- [ ] Both frontend and backend are on the same machine

---

## Advanced Debugging

### Enable Verbose Logging

```bash
cd frontend/apps/super-admin
DEBUG=* pnpm dev
```

### Check Next.js Logs

```bash
cd frontend/apps/super-admin
pnpm dev 2>&1 | grep -E "error|Error|ERROR"
```

### Verify File Structure

```bash
cd frontend/apps/super-admin

# Check critical files exist
ls -la src/app/layout.tsx
ls -la src/app/login/page.tsx
ls -la tailwind.config.ts
ls -la next.config.js
ls -la package.json
```

### Reset Everything (Nuclear Option)

```bash
cd C:\Users\OPT\Desktop\POS

# Kill all Node processes
taskkill /IM node.exe /F

# Clean entire frontend
cd frontend
rm -rf node_modules
rm -rf apps/*/node_modules
rm -rf apps/*/.next
rm -rf .pnpm-store
rm -rf pnpm-lock.yaml

# Reinstall from scratch
pnpm install

# Start super-admin
cd apps/super-admin
pnpm dev
```

---

## Performance Monitoring

Monitor these metrics during development:

```bash
# Watch build time (should be < 5 seconds)
time pnpm dev

# Monitor memory usage
# Windows: Open Task Manager → Processes → Look for "node"
# Linux: top | grep node
# Mac: top | grep node
```

---

## Getting Help

If you encounter issues:

1. **Check this guide** - Most issues are listed above
2. **Check console output** - Copy the full error message
3. **Look at browser DevTools** - F12 → Console & Network tabs
4. **Check backend** - Ensure `http://localhost:8080` is responding
5. **Review git changes** - Ensure no unexpected modifications

---

## Quick Reference Commands

```bash
# Start super-admin alone
cd frontend/apps/super-admin && pnpm dev

# Start all frontend apps
cd frontend && pnpm dev

# Kill port 3004 process (Windows)
FOR /F "tokens=5" %A IN ('netstat -ano ^| findstr :3004') DO taskkill /PID %A /F

# Kill port 3004 process (Mac/Linux)
kill -9 $(lsof -t -i :3004)

# Clear all caches
rm -rf ~/.next node_modules .next dist pnpm-lock.yaml

# Test API
curl http://localhost:8080/api/v1/super-admin/dashboard/stats

# Verify backend running
curl http://localhost:8080
```

---

## CURRENT ISSUE: Dashboard Access and Debug Mode

### Problem Statement

When a super admin tries to access the tenant dashboard via the auto-login flow:

**Expected Flow:**
1. Super admin clicks "Enter Dashboard" button on tenant detail page (port 3004)
2. Generate one-time token on backend
3. Redirect to landing app token exchange page (port 3001)
4. Token exchange returns JWT and user data
5. Redirect to dashboard app with auth params in URL: `/dashboard?auth_token=...&auth_user=...&debug=true`
6. **Debug panel should appear** showing all authentication steps
7. User can click "Try Dashboard" to access the actual dashboard

**Actual Behavior:**
- When visiting `/dashboard?auth_token=...&auth_user=...&debug=true`
- Page immediately redirects to login page: `http://localhost:3001/auth/login`
- Debug panel never appears
- Auth parameters in URL are ignored

### Root Cause Analysis

**Discovery 1:** Even the unprotected test page (`/en/test-url`) redirects to login
- This proves the redirect is NOT coming from DashboardAuthProvider component
- The redirect happens at a **higher level** (layout, middleware, or Next.js router)

**Current Finding:**
- The dashboard app has a `[locale]` routing pattern
- All routes must match `/[locale]/path` pattern
- When visiting `/dashboard`, Next.js may not match the route correctly
- When visiting `/en/dashboard`, it should work, but still redirects

### Testing Procedure

To debug this issue:

**Test 1: Simple Test Page (No Auth)**
```
http://localhost:3002/en/test-url?auth_token=...&auth_user=...&debug=true
```

Expected: Blue test page appears showing URL and query parameters
Actual: Redirects to login instead

**Test 2: Dashboard with Debug Mode**
```
http://localhost:3002/en/dashboard?auth_token=...&auth_user=...&debug=true
```

Expected: Dark debug panel with all logs
Actual: Redirects to login instead

### Next Steps to Resolve

1. **Investigate Route Protection**
   - Check if there's a global middleware redirecting unauthenticated users
   - Check the [locale] layout for any auth checks
   - Look for Next.js route protection in next.config.js

2. **Verify URL Parameters**
   - Confirm `?auth_token=` and `?auth_user=` reach the component
   - Debug if they're being stripped by Next.js router

3. **Disable Auth Check for Debug Mode**
   - Once debug=true is confirmed in URL, skip all auth checks
   - Render the debug panel regardless of authentication status

4. **Files to Investigate**
   - `frontend/apps/dashboard/src/app/[locale]/layout.tsx` - Auth checks?
   - `frontend/apps/dashboard/src/app/[locale]/dashboard/layout.tsx` - DashboardAuthProvider
   - `frontend/apps/dashboard/next.config.js` - Route redirects?
   - `frontend/apps/dashboard/src/middleware.ts` - Global middleware?

### Debug Logs Location

When debug mode works, logs are stored in localStorage:
```javascript
localStorage.getItem('debug_logs') // Array of {time, msg} objects
```

### Related Files

- Landing app token exchange: `frontend/apps/landing/src/app/auth/super-admin-login/page.tsx`
- Dashboard auth provider: `frontend/apps/dashboard/src/components/providers/DashboardAuthProvider.tsx`
- Dashboard layouts: `frontend/apps/dashboard/src/app/[locale]/dashboard/layout.tsx`

---

**Last Updated:** 2026-01-31
**Status:** In Progress - Debugging Global Redirect Issue
**Priority:** High - Blocks super admin dashboard access feature
