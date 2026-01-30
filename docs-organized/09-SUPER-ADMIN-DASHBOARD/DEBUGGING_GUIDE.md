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

**Last Updated:** 2024-01-30
**Status:** Production Ready
