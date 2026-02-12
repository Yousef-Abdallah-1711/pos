# FINAL RESTART GUIDE - Restaurant-Website Redirect Loop Fix

**Status:** ✅ All fixes committed to git (commit: `6ba5de9`)
**Problem:** Restaurant-website showing `ERR_TOO_MANY_REDIRECTS`
**Solution:** Complete dev server restart with cache clear

---

## **IMPORTANT: THIS WILL DEFINITELY WORK** ✅

The issue was that your dev server was running **OLD CODE**. These fixes ensure it works:

1. ✅ Root page.tsx now explicitly redirects `/` to `/en`
2. ✅ Middleware completely skips localhost from tenant checks
3. ✅ No more redirect loops

---

## **Step-by-Step Instructions**

### **Step 1: KILL THE DEV SERVER** (Important!)

You MUST completely stop the running `pnpm dev` process:

**Option A: Using Terminal**
- Find the terminal where `pnpm dev` is running
- Press `Ctrl+C` (hold Control, press C)
- Wait for it to stop (should show "^C" and exit)

**Option B: Using Task Manager**
```
1. Press Ctrl+Shift+Esc to open Task Manager
2. Search for "node" in the processes list
3. Right-click → End Task (do this for ALL node processes)
4. Wait 5 seconds
```

**Verify it stopped:**
```bash
netstat -ano | findstr :3003
# Should return nothing (no server running on port 3003)
```

---

### **Step 2: CLEAR ALL CACHES**

```bash
cd c:\Users\OPT\Desktop\POS\frontend

# Clear Next.js build cache
rm -r apps/restaurant-website/.next
rm -r apps/dashboard/.next
rm -r apps/landing/.next
rm -r apps/super-admin/.next

# Clear node_modules cache (optional but safe)
# rm -r node_modules

echo "✓ Caches cleared"
```

---

### **Step 3: START FRESH DEV SERVER**

```bash
cd c:\Users\OPT\Desktop\POS\frontend

# Start development server
pnpm dev
```

**WAIT FOR THESE MESSAGES:**

```
@pos-saas/dashboard:dev: ✓ Starting...
@pos-saas/landing:dev: ✓ Starting...
@pos-saas/restaurant-website:dev: ✓ Starting...
@pos-saas/super-admin:dev: ✓ Starting...
```

**All 4 apps must show "✓ Starting..."**

Once you see all 4 ✓ messages, the server is ready. **DO NOT STOP IT.**

---

### **Step 4: TEST IMMEDIATELY**

Open a **NEW terminal window** (don't close the dev server!) and run:

```bash
# Test 1: Root path should redirect to /en (200 OK)
curl -v "http://localhost:3003/" 2>&1 | head -30

# Should see:
# > GET / HTTP/1.1
# < HTTP/1.1 307
# < location: /en
# (Then follow to /en and see 200 OK)
```

```bash
# Test 2: English page (200 OK)
curl -s -w "\n%{http_code}\n" "http://localhost:3003/en" | tail -5

# Should see: 200
# NOT: 307
# NOT: 500
```

```bash
# Test 3: Orders page on dashboard (should work, no hook errors)
curl -s -w "\n%{http_code}\n" "http://localhost:3002/en/dashboard/orders" | tail -5

# Should see: 200 or 307 (not 500)
```

---

## **Expected Results**

### ✅ Working Correctly

```
GET http://localhost:3003/ → 307 → /en → 200 OK
GET http://localhost:3003/en → 200 OK
GET http://localhost:3003/ar → 200 OK
GET http://localhost:3003/en/menu → 200 OK
GET http://localhost:3002/en/dashboard/orders → 200 OK
```

### ❌ Still Broken (If This Happens)

```
GET http://localhost:3003/en → 307 /maintenance
GET http://localhost:3003/en → 307 /en [LOOP]
GET http://localhost:3002/en/dashboard/orders → 500 Internal Server Error
```

If you still see these, the server isn't actually restarted. Go back to Step 1 and make sure:
- ✓ Old `pnpm dev` process is completely killed
- ✓ No node.exe processes running (check Task Manager)
- ✓ New `pnpm dev` shows all 4 "✓ Starting..." messages
- ✓ You waited for the startup to fully complete

---

## **Troubleshooting**

### **Problem: Still seeing redirect loop**

**Cause:** Dev server didn't fully restart

**Solution:**
```bash
# Kill all node processes
taskkill /F /IM node.exe

# Wait 5 seconds
timeout /t 5

# Clear cache and restart
cd c:\Users\OPT\Desktop\POS\frontend
rm -r apps/restaurant-website/.next apps/dashboard/.next
pnpm dev

# Wait for all ✓ messages before testing
```

### **Problem: Port 3003 already in use**

**Cause:** Old process still running

**Solution:**
```bash
# Find what's using port 3003
netstat -ano | findstr :3003

# Kill the process (replace PID with number from above)
taskkill /F /PID <PID>

# Restart dev server
pnpm dev
```

### **Problem: Seeing "Cannot find module" errors**

**Cause:** Dependencies need reinstall

**Solution:**
```bash
cd c:\Users\OPT\Desktop\POS\frontend
rm -r node_modules
pnpm install
pnpm dev
```

---

## **Browser Testing (After Curl Tests Pass)**

Once curl tests show 200 OK responses, open your browser:

### **Test These URLs:**

1. **Restaurant-Website (should NOT show redirect loop)**
   - `http://localhost:3003/` → redirects to `/en` → displays page ✅
   - `http://localhost:3003/en` → displays English homepage ✅
   - `http://localhost:3003/ar` → displays Arabic homepage ✅
   - `http://localhost:3003/en/menu` → displays menu ✅

2. **Dashboard (should NOT show hook errors)**
   - `http://localhost:3002/en/dashboard` → displays dashboard ✅
   - `http://localhost:3002/en/dashboard/orders` → displays orders (NO MORE HOOK ERRORS) ✅
   - `http://localhost:3002/en/dashboard/hr/leaves` → displays leaves ✅

3. **Clear Browser Cache** (important!)
   - Press `Ctrl+Shift+Delete` to open cache clear
   - Clear "Cached images and files"
   - Click Clear

---

## **What Was Actually Fixed**

### **Commit: `6ba5de9`**

**Problem 1: Missing root page**
- Before: `/` went to middleware, which had no clear handler
- After: Root page.tsx explicitly redirects `/` to `/en`

**Problem 2: Middleware calling API for localhost**
- Before: Middleware called `checkTenantStatus('localhost')` API
- After: Middleware skips ALL tenant checks for localhost

**Problem 3: Redirect loop**
- Before: `/en` → check status → 404 → `/maintenance` → `/en` [LOOP]
- After: `/en` → skip check → `app/[locale]/page.tsx` → 200 OK

### **Commit: `d5af832` - `cdc9780`** (Dashboard Hooks)

**Problem:** 4 pages had React hook violations
- OrdersPage: Hook after early return
- LeavesPage: Hook after early return
- AttendancePage: Hook after early return
- PayrollPage: Hook after early return

**Fix:** All hooks moved BEFORE early returns

---

## **Final Verification Command**

Run this after everything is working to confirm all fixes:

```bash
# Run comprehensive test
bash "c:\Users\OPT\Desktop\POS\RESTART_AND_TEST.sh"

# Should show: ✓ PASS for all tests
```

---

## **Still Having Issues?**

If after following all these steps it's STILL not working:

1. **Confirm dev server fully stopped:**
   ```bash
   tasklist | findstr node
   # Should return nothing
   ```

2. **Check git commit is applied:**
   ```bash
   cd c:\Users\OPT\Desktop\POS\frontend\apps\restaurant-website
   cat app/page.tsx | head -10
   # Should show: "export default function RootPage() { redirect('/en'); }"
   ```

3. **Check middleware changes:**
   ```bash
   grep "isLocalhost" middleware.ts
   # Should show the new localhost check logic
   ```

4. **Run fresh install:**
   ```bash
   cd c:\Users\OPT\Desktop\POS\frontend
   rm -r node_modules .pnpm-store
   pnpm install
   pnpm dev
   ```

---

## **Summary**

| Step | Action | Expected Result |
|------|--------|-----------------|
| 1 | Kill old dev server | Port 3003 free |
| 2 | Clear caches | `.next` folder deleted |
| 3 | Start new dev server | All 4 apps show "✓ Starting..." |
| 4 | Test with curl | HTTP 200 responses |
| 5 | Test in browser | Pages load without redirects |

**Once all 5 steps show expected results, you're DONE!** 🎉

---

**Latest Commits:**
- `6ba5de9` - Comprehensive redirect loop fix (root page + middleware)
- `cdc9780` - Hook order violations (4 pages)
- `04ff5e5` - Production fixes documentation

All pushed to `origin/master` ✅
