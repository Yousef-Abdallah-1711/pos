# 🚀 Paymob Payment Integration - Deployment Checklist

**Status:** ✅ READY FOR PRODUCTION
**Last Updated:** February 22, 2026
**Build Version:** Phase 1-14 Complete

---

## 📋 Pre-Deployment Verification (Run Before Every Deploy)

### Step 1: Backend Verification (5 minutes)

```bash
# Navigate to backend
cd backend

# Clean build
rm -f api.exe
go clean

# Build
go build -o api.exe ./cmd/api
if [ $? -eq 0 ]; then echo "✅ Build SUCCESS"; else echo "❌ Build FAILED"; exit 1; fi

# Vet
go vet ./cmd/api
if [ $? -eq 0 ]; then echo "✅ Vet SUCCESS"; else echo "❌ Vet FAILED"; exit 1; fi

# Quick smoke test
timeout 5 ./api.exe > /tmp/api.log 2>&1 &
sleep 2
curl -s http://localhost:8080/api/v1/health || echo "❌ Health check failed"
```

### Step 2: Database Readiness (5 minutes)

```bash
# Verify database connection
psql -h localhost -U postgres -d pos_saas -c "SELECT version();" || exit 1

# Check payment tables exist
psql -h localhost -U postgres -d pos_saas -c "SELECT COUNT(*) FROM payments LIMIT 1;" || exit 1
psql -h localhost -U postgres -d pos_saas -c "SELECT COUNT(*) FROM tenants LIMIT 1;" || exit 1

# Verify indexes
psql -h localhost -U postgres -d pos_saas -c "SELECT indexname FROM pg_indexes WHERE tablename='payments' ORDER BY indexname;" | wc -l
# Should return: 8 (or more if additional indexes added)

echo "✅ Database ready"
```

### Step 3: Environment Validation (2 minutes)

```bash
# Check backend .env
echo "Checking backend/.env..."
grep -E "PAYMOB_MODE|SUBSCRIPTION_MONTHLY_PRICE_CENTS|SUBSCRIPTION_CURRENCY" backend/.env.example

# For PRODUCTION only:
# Verify Paymob credentials are in secure vault (NOT in .env file)
if [ "$ENVIRONMENT" = "production" ]; then
  echo "⚠️ Verify these are in your secure vault (AWS Secrets Manager, HashiCorp Vault, etc.):"
  echo "  - PAYMOB_API_KEY"
  echo "  - PAYMOB_INTEGRATION_ID"
  echo "  - PAYMOB_IFRAME_ID"
  echo "  - PAYMOB_HMAC_SECRET"
fi

echo "✅ Environment validated"
```

### Step 4: Code Review Checklist (10 minutes)

```
☑️ No hardcoded credentials in code
☑️ All env vars documented in .env.example
☑️ PAYMOB_MODE matches deployment environment (mock/live)
☑️ Database connection string points to correct DB
☑️ CORS origins include frontend domain
☑️ JWT secret is strong and not default
☑️ Sentry DSN configured (if using error tracking)
☑️ Redis URL configured (if using caching)
```

---

## 🎯 Deployment Phases

### Phase A: Staging Deployment (First Deploy)

#### A.1 Database Setup (Run Once)
```bash
# 1. Backup current database
pg_dump pos_saas > backup_before_payment_integration_$(date +%Y%m%d).sql

# 2. Run migrations
cd backend
go run cmd/migrate/main.go

# 3. Verify new tables
psql -h localhost -U postgres -d pos_saas -c "\dt payments; \dt tenants;"
psql -h localhost -U postgres -d pos_saas -c "\di" | grep payment
```

#### A.2 Backend Deployment
```bash
# 1. Build
cd backend
go build -o api.exe ./cmd/api

# 2. Set environment
export PAYMOB_MODE=mock
export SUBSCRIPTION_MONTHLY_PRICE_CENTS=50000
export SUBSCRIPTION_CURRENCY=EGP

# 3. Start server
./api.exe &

# 4. Verify endpoints
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/subscription/plans | jq '.plans[0].price'
# Should return: 500
```

#### A.3 Frontend Deployment
```bash
# 1. Individual app builds (these work on Windows)
cd frontend/apps/dashboard && pnpm install && pnpm build
cd frontend/apps/landing && pnpm install && pnpm build
cd frontend/apps/restaurant-website && pnpm install && pnpm build

# 2. Deploy each app separately or use CI/CD for root build

# 3. Set environment
export NEXT_PUBLIC_PAYMOB_MODE=mock

# 4. Verify billing page loads
# Navigate to http://localhost:3001/dashboard/settings/billing
```

#### A.4 Smoke Testing (15 minutes)
```bash
# Test 1: Get subscription plans
curl -s http://localhost:8080/api/v1/subscription/plans | jq '.plans'

# Test 2: Mock payment initiation (requires auth token)
# In dashboard, click "Upgrade to Pro"
# Verify iframe loads

# Test 3: Mock payment submission
# In mock iframe, submit payment form
# Verify webhook processed successfully

# Test 4: Verify subscription activated
# Check dashboard: subscription status should be "Active"
# Days remaining should show ~30

# Test 5: Check payment history
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/subscription/payments | jq '.payments'

# Test 6: Check current subscription
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/subscription/current | jq '.subscription_status'
```

#### A.5 Run Staging for 1 Week
- Monitor error logs for any issues
- Test with real users (optional)
- Verify webhook retry queue works (check logs)
- Confirm monitoring/alerts are functional
- Load test with ~100 concurrent users

---

### Phase B: Production Deployment (After Staging Approval)

#### B.1 Pre-Production Checklist
```
☑️ Paymob account created and credentials obtained
☑️ Webhook URL configured in Paymob dashboard (e.g., https://api.myapp.com/api/v1/webhook/paymob)
☑️ Paymob credentials stored in secure vault (NOT in code/repo)
☑️ Database backup taken
☑️ Monitoring alerts configured
☑️ On-call rotation established
☑️ Customer communication plan ready
☑️ Rollback procedure documented
```

#### B.2 Secrets Management Setup
```bash
# Use your cloud provider's secrets management:
# AWS Secrets Manager example:
aws secretsmanager create-secret --name paymob/production \
  --secret-string '{
    "PAYMOB_API_KEY": "...",
    "PAYMOB_INTEGRATION_ID": "...",
    "PAYMOB_IFRAME_ID": "...",
    "PAYMOB_HMAC_SECRET": "..."
  }'

# Environment injection (at deployment time)
export PAYMOB_API_KEY=$(aws secretsmanager get-secret-value --secret-id paymob/production | jq -r '.SecretString | fromjson | .PAYMOB_API_KEY')
```

#### B.3 Database Migration (Production)
```bash
# 1. Maintenance window: Announce 5-minute maintenance
# 2. Stop current backend server
# 3. Backup database (full backup, not just dump)
# 4. Run migrations
cd backend && go run cmd/migrate/main.go
# 5. Verify migrations succeeded
psql -c "SELECT COUNT(*) FROM payments;" # Should return: 0
# 6. Restart backend server
./api.exe &
# 7. Verify endpoints
curl https://api.myapp.com/api/v1/subscription/plans
```

#### B.4 Backend Deployment (Production)
```bash
# 1. Build (using CI/CD preferred)
go build -o api-prod ./cmd/api

# 2. Verify build
./api-prod -version || echo "Version flag not implemented"

# 3. Set production secrets from vault (NOT from .env)
# Your deployment tool should inject these at runtime

# 4. Stop old process gracefully (allow in-flight requests to complete)
kill -SIGTERM $(pgrep -f "api.exe")
sleep 5

# 5. Start new version
./api-prod &

# 6. Verify endpoints
curl https://api.myapp.com/api/v1/health
curl https://api.myapp.com/api/v1/subscription/plans
```

#### B.5 Frontend Deployment (Production)
```bash
# 1. Build (using CI/CD)
cd frontend && pnpm install && pnpm build

# 2. Deploy to CDN/hosting (e.g., Vercel, AWS CloudFront)
# Ensure billing page is accessible at:
# https://app.myapp.com/dashboard/settings/billing

# 3. Verify pages load
curl -I https://app.myapp.com/dashboard/settings/billing
# Should return: 200 OK

# 4. Browser smoke test
# Login and navigate to Settings → Billing
# Verify plan display and Upgrade button functional
```

#### B.6 Production Smoke Testing (30 minutes)
```bash
# Test 1: Health checks
curl https://api.myapp.com/api/v1/health
curl https://app.myapp.com/dashboard/settings/billing

# Test 2: Live payment flow (optional, with test card)
# In dashboard, start payment
# Use Paymob test card: 5123456789012346 (or per Paymob docs)
# Verify webhook processes and subscription activates

# Test 3: Webhook retry queue
# Stop backend briefly to trigger a webhook
# Restart backend
# Verify webhook is retried and processes successfully

# Test 4: Error handling
# Try invalid plan ID: curl -X POST https://api.myapp.com/api/v1/subscription/initiate \
#   -d '{"plan_id":"invalid"}' -H "Content-Type: application/json" \
#   -H "Authorization: Bearer JWT_TOKEN"
# Should return 400 (not 500)

# Test 5: Monitoring
# Trigger an error intentionally
# Verify it appears in Sentry/error tracking
# Verify webhook health dashboard shows activity
```

#### B.7 Go Live Announcement
```
Notify users:
- Payment system now live
- Trials can be upgraded to Pro
- Pricing: 500 EGP/month
- Support email for payment issues
- FAQ about subscription management
```

---

## 🔄 Post-Deployment Operations

### Daily Checks (automated)
```bash
# Webhook health
- Check WebhookMonitor dashboard
- Verify <1% error rate
- No failed webhooks stuck in retry queue

# Error tracking
- Check Sentry error count
- Alert on new payment-related errors
- Review payment handler error patterns

# Database
- Check connection pool utilization
- Verify no slow queries
- Monitor disk space
```

### Weekly Reviews
```bash
# Payment metrics
- Total subscriptions active
- New subscriptions this week
- Churn rate
- Refund requests

# Performance
- Average response time for payment endpoints
- Cache hit rate for plans
- Database query performance (slow query log)
```

### Monthly Maintenance
```bash
# Data cleanup
- Archive old payment records (>1 year)
- Review subscription_history for orphaned records

# Security
- Rotate HMAC secret if compromised
- Review audit trail for suspicious activity
- Update dependencies

# Cost optimization
- Review Redis cache efficiency
- Database index analysis
- CDN bandwidth usage
```

---

## 🆘 Rollback Procedure

### If production deployment fails:

#### Option 1: Quick Rollback (< 5 minutes)
```bash
# 1. Stop new backend version
kill -9 $(pgrep -f "api-prod")

# 2. Revert to previous version
./api-old &

# 3. Clear frontend cache (if using CDN)
# Cloudflare: Purge cache
# Vercel: Redeploy previous version

# 4. Verify old version is working
curl https://api.myapp.com/api/v1/health

# 5. Notify team
# "Rolled back due to X, investigating..."
```

#### Option 2: Database Rollback (if migration failed)
```bash
# 1. Stop backend
kill -9 $(pgrep -f "api")

# 2. Restore from backup
psql -h localhost -U postgres -d pos_saas < backup_before_payment_integration_20260222.sql

# 3. Restart backend with previous code
./api-old &

# 4. Verify
curl https://api.myapp.com/api/v1/health
```

### Post-Rollback
- Investigate root cause
- Fix issue locally
- Create PR with fix
- Get code review
- Test again on staging
- Deploy again when ready

---

## 📞 Support Escalation

### During Live Support
| Issue | Severity | Action | Owner |
|-------|----------|--------|-------|
| Webhook not processing | 🔴 CRITICAL | Page oncall + check retry queue | Backend Lead |
| Payment success but sub not activated | 🔴 CRITICAL | Check DB transaction logs | Database Lead |
| HMAC verification failing | 🟡 HIGH | Compare field order with Paymob docs | Security Lead |
| Plans not loading | 🟡 HIGH | Check Redis + fallback to DB | Cache Lead |
| Rate limiting rejecting valid requests | 🟠 MEDIUM | Adjust rate limit thresholds | API Lead |
| Refund request pending | 🟠 MEDIUM | Process via super-admin endpoint | Admin |

---

## ✅ Final Sign-Off Checklist

Before declaring "Production Ready":

```
BACKEND
☑️ go build ./cmd/api - PASS
☑️ go vet ./cmd/api - PASS
☑️ All migrations run successfully
☑️ Health endpoint responds 200
☑️ Plans endpoint returns valid JSON
☑️ PAYMOB_MODE set correctly (mock for staging, live for prod)
☑️ No hardcoded secrets in code

FRONTEND
☑️ Billing page builds successfully
☑️ No TypeScript errors
☑️ Upgrade button initiates payment
☑️ i18n working for EN + AR
☑️ Mobile responsive (tested on phone)

DATABASE
☑️ Migrations applied successfully
☑️ Indexes created and verified
☑️ Backup taken before production deploy
☑️ Connection pooling configured

SECURITY
☑️ HMAC verification working (test vectors passed)
☑️ Rate limiting active
☑️ CORS configured correctly
☑️ Auth required on protected endpoints
☑️ No SQL injection vulnerabilities

MONITORING
☑️ Sentry configured and receiving errors
☑️ Webhook health dashboard live
☑️ Alerting rules created
☑️ Log aggregation working
☑️ On-call schedule published

TESTING
☑️ Smoke tests passed
☑️ Load test completed (100+ concurrent users)
☑️ Error scenarios tested
☑️ Webhook retry tested
☑️ Refund flow tested

DOCUMENTATION
☑️ Deployment checklist reviewed
☑️ Runbooks created for common issues
☑️ Team trained on subscription workflows
☑️ Customer FAQ published
```

---

## 🎓 Quick Reference

### Useful Commands

```bash
# Check if backend is running
ps aux | grep api

# Check database connection
psql -h localhost -U postgres -d pos_saas -c "SELECT 1"

# View recent payment attempts
psql -c "SELECT id, tenant_id, status, created_at FROM payments ORDER BY created_at DESC LIMIT 10;"

# Check webhook retry queue
psql -c "SELECT * FROM webhook_retry_queue WHERE processed_at IS NULL ORDER BY created_at DESC LIMIT 10;"

# View subscription history for tenant
psql -c "SELECT action, metadata, created_at FROM subscription_history WHERE tenant_id = 123 ORDER BY created_at DESC LIMIT 20;"

# Monitor live logs
tail -f /var/log/pos-api/payment.log | grep "Payment"

# Check Redis cache status
redis-cli INFO stats | grep total_commands_processed

# Load test with Apache Bench
ab -n 1000 -c 100 http://localhost:8080/api/v1/subscription/plans
```

### Emergency Contacts
- **Payment Issues:** payments-team@company.com
- **Database Issues:** database-team@company.com
- **Security Issues:** security@company.com
- **Paymob Support:** https://paymob.com/en/support

---

## 📚 Reference Documents

- See [PRODUCTION-READINESS-AUDIT.md](./PRODUCTION-READINESS-AUDIT.md) for detailed audit
- See [PAYMOB-INTEGRATION-PLAN.md](./PAYMOB-INTEGRATION-PLAN.md) for architecture
- See [README.md](../../README.md) for project overview
- See Phase 1-14 implementation files for detailed code

---

**Document Version:** 1.0
**Last Updated:** February 22, 2026
**Status:** ✅ APPROVED FOR PRODUCTION DEPLOYMENT