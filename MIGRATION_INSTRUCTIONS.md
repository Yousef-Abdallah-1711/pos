# Database Migration Setup Guide

**Status**: Backend is running, migrations need to be applied

---

## Option 1: Using pgAdmin (Easiest - GUI)

### Step 1: Open pgAdmin
1. Open pgAdmin (comes with PostgreSQL installation)
2. Right-click on **Databases** → **Create** → **Database**
3. Name: `pos_saas` (if not already created)

### Step 2: Open Query Tool
1. Expand **Databases** → **pos_saas**
2. Right-click on **pos_saas** → **Query Tool**

### Step 3: Copy and Paste SQL

**Migration 110** - Add subscription fields:
```sql
-- Copy the entire contents of: backend/migrations/110_add_tenant_subscription_fields.sql
-- Paste it here and execute
```

**Migration 111** - Create subscription history table:
```sql
-- Copy the entire contents of: backend/migrations/111_create_subscription_history_table.sql
-- Paste it here and execute
```

**Migration 112** - Seed existing tenants:
```sql
-- Copy the entire contents of: backend/migrations/112_seed_existing_tenants_trial.sql
-- Paste it here and execute
```

**Test Data** - Populate test records:
```sql
-- Copy the entire contents of: docs-organized/10-SAAS-TRIAL-SUBSCRIPTION/test-data-setup.sql
-- Paste it here and execute
```

---

## Option 2: Using Command Line (Windows PowerShell)

### Step 1: Open PowerShell as Administrator

### Step 2: Find PostgreSQL Installation
```powershell
# Find psql executable
Get-Command psql
# OR
dir "C:\Program Files\PostgreSQL\*\bin\psql.exe"
```

### Step 3: Run Migrations

```powershell
# Set PostgreSQL path
$psqlPath = "C:\Program Files\PostgreSQL\15\bin\psql.exe"

# Run each migration in order
& $psqlPath -U postgres -d pos_saas -f "C:\Users\OPT\Desktop\POS\backend\migrations\110_add_tenant_subscription_fields.sql"
& $psqlPath -U postgres -d pos_saas -f "C:\Users\OPT\Desktop\POS\backend\migrations\111_create_subscription_history_table.sql"
& $psqlPath -U postgres -d pos_saas -f "C:\Users\OPT\Desktop\POS\backend\migrations\112_seed_existing_tenants_trial.sql"
& $psqlPath -U postgres -d pos_saas -f "C:\Users\OPT\Desktop\POS\docs-organized\10-SAAS-TRIAL-SUBSCRIPTION\test-data-setup.sql"

Write-Host "✅ All migrations completed!"
```

---

## Option 3: Using DBeaver or DataGrip

### Step 1: Create New Connection
1. File → New → Database Connection
2. Select PostgreSQL
3. Hostname: `localhost`
4. Port: `5432`
5. Database: `pos_saas`
6. Username: `postgres`
7. Password: (your password)
8. Test Connection → Finish

### Step 2: Create SQL Script
1. Right-click on connection → New → SQL Script
2. Copy entire contents of first migration file
3. Execute (Ctrl+Enter)
4. Repeat for each migration file

---

## SQL Files Location

All migration files are in your project:

```
C:\Users\OPT\Desktop\POS\
├── backend\migrations\
│   ├── 110_add_tenant_subscription_fields.sql
│   ├── 111_create_subscription_history_table.sql
│   └── 112_seed_existing_tenants_trial.sql
└── docs-organized\10-SAAS-TRIAL-SUBSCRIPTION\
    └── test-data-setup.sql
```

---

## Verification

After running all migrations, verify in pgAdmin or psql:

```sql
-- Check if subscription fields exist
SELECT column_name FROM information_schema.columns
WHERE table_name='tenants' AND column_name LIKE '%trial%';

-- Check if history table exists
SELECT table_name FROM information_schema.tables
WHERE table_name='tenant_subscription_history';

-- Check test data
SELECT id, slug, subscription_status, is_disabled FROM tenants LIMIT 5;
```

---

## After Migrations: Restart Backend

Once migrations are complete:

```bash
# Kill the old backend process
taskkill /F /IM go.exe

# Or use this bash command
kill -9 $(pgrep -f "go run cmd/api")

# Start fresh
cd C:\Users\OPT\Desktop\POS\backend
go run cmd/api/main.go
```

---

## Test After Migrations

Once backend is restarted with DB migrations applied:

```bash
# Test with real data
curl http://localhost:8080/api/v1/public/tenant-status/pizza-palace
```

Expected response:
```json
{
  "can_access": true,
  "block_reason": null,
  "is_maintenance": false
}
```

---

## Troubleshooting

### "database does not exist"
- Create the database first: `CREATE DATABASE pos_saas;`

### "relation already exists"
- Migrations may have been already applied
- Check with the verification SQL above

### "password authentication failed"
- Verify PostgreSQL is running
- Check username/password are correct
- Default username is usually `postgres`

### "permission denied"
- Make sure you're connected as `postgres` user
- May need superuser rights to create tables

---

## Quick Checklist

- [ ] PostgreSQL is running
- [ ] `pos_saas` database exists
- [ ] Migration 110 applied (subscription fields)
- [ ] Migration 111 applied (history table)
- [ ] Migration 112 applied (seed data)
- [ ] Test data applied
- [ ] Backend restarted
- [ ] API returns real data: `curl http://localhost:8080/api/v1/public/tenant-status/pizza-palace`

---

## Next Steps After Migrations

1. ✅ Apply migrations (THIS STEP)
2. ✅ Restart backend server
3. ✅ Test endpoints with real data
4. Test frontend (Dashboard & Restaurant Website)
5. Run full testing suite

---

**Timeline**: 5-10 minutes to apply migrations + restart
**Status**: Backend waiting for migrations to enable full functionality
