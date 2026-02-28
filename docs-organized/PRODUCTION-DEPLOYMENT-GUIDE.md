# Production Deployment Guide

**Last Updated:** February 22, 2026

---

## Current Status (What Is Already Done In Repo)

The following deployment-readiness work is already implemented in the codebase:

- URL/domain runtime hardening across backend + all frontend apps (`restaurant-website`, `dashboard`, `landing`, `super-admin`)
- Development-only localhost fallbacks with production fail-fast behavior for missing required URL/domain env vars
- Clean local vs production env templates (`.env.example` and `.env.production.example`)
- CI guardrail for hardcoded localhost regressions in runtime code:
  - `.github/workflows/runtime-localhost-guardrail.yml`
  - `scripts/check_runtime_localhost.py`
- CI smoke checks for URL/domain configuration:
  - `.github/workflows/url-domain-config-smoke.yml`
  - `scripts/smoke_check_url_domain_config.py`
- Unit tests for URL/config resolvers (frontend apps) and backend production URL validation
- Post-deploy endpoint verification helper:
  - `scripts/verify_deployment_endpoints.py`

What is still manual/external (not done from repo code):

- Real DNS records and wildcard DNS
- Real staging/production secrets and env values
- Reverse proxy deployment (Nginx/Caddy) and SSL certificate setup
- Paymob webhook configuration in Paymob dashboard
- Staging deployment execution, QA signoff, rollback rehearsal, production cutover

---

## Architecture Overview

```
                    ┌──────────────┐
                    │   DNS/CDN    │
                    │ (Cloudflare) │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
     ┌────────▼──┐  ┌──────▼───┐  ┌────▼────────┐
     │  Landing  │  │ Dashboard│  │ Super Admin │
     │  :3001    │  │  :3002   │  │   :3004     │
     │ Next.js   │  │ Next.js  │  │  Next.js    │
     └───────────┘  └──────────┘  └─────────────┘
                           │
              ┌────────────▼────────────┐
              │   Restaurant Websites   │
              │  *.yourdomain.com:3003  │
              │      Next.js (SSR)      │
              └────────────┬────────────┘
                           │
              ┌────────────▼────────────┐
              │    Backend API :8080    │
              │     Go (net/http)       │
              └────────┬───────┬────────┘
                       │       │
              ┌────────▼──┐ ┌──▼────────┐
              │ PostgreSQL│ │   Redis   │
              │   :5432   │ │   :6379   │
              └───────────┘ └───────────┘
```

---

## Domain Configuration

### Required DNS Records

| Subdomain | Service | Record Type | Value |
|-----------|---------|-------------|-------|
| `www.yourdomain.com` | Landing Page | A/CNAME | Server IP / CDN |
| `dashboard.yourdomain.com` | Dashboard | A/CNAME | Server IP / CDN |
| `admin.yourdomain.com` | Super Admin | A/CNAME | Server IP / CDN |
| `api.yourdomain.com` | Backend API | A | Server IP |
| `*.yourdomain.com` | Restaurant Sites | A/CNAME | Server IP / CDN |

### Wildcard Subdomain for Restaurants

Each tenant gets a subdomain: `{restaurant-slug}.yourdomain.com`

Requirements:
- Wildcard DNS: `*.yourdomain.com → Server IP`
- Wildcard SSL certificate (Let's Encrypt supports this)
- Nginx/Caddy reverse proxy to route subdomains to restaurant-website app

---

## Environment Variables Reference

### Backend (`backend/.env.production.example`)

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `PORT` | Yes | API server port | `8080` |
| `ENV` | Yes | Environment mode | `production` |
| `DATABASE_URL` | Yes | PostgreSQL connection | `postgresql://user:pass@host:5432/db` |
| `REDIS_URL` | Yes | Redis connection | `redis://host:6379` |
| `JWT_SECRET` | Yes | JWT signing secret | `openssl rand -base64 64` |
| `CORS_ORIGINS` | Yes | Allowed CORS origins | `https://www.yourdomain.com,...` |
| `APP_BASE_URL` | Yes | Public backend URL | `https://api.yourdomain.com` |
| `DASHBOARD_URL` | Yes | Dashboard app URL | `https://dashboard.yourdomain.com` |
| `LANDING_URL` | Yes | Landing page URL | `https://www.yourdomain.com` |
| `RESTAURANT_DOMAIN` | Yes | Base domain for tenants | `yourdomain.com` |
| `UPLOAD_BASE_URL` | Yes | File upload URL | `https://api.yourdomain.com/uploads` |
| `WEBHOOK_BASE_URL` | Yes | Webhook callback URL | `https://api.yourdomain.com` |
| `PAYMOB_MODE` | Yes | Payment mode | `live` |
| `PAYMOB_API_KEY` | Yes | Paymob API key | From Paymob dashboard |
| `PAYMOB_INTEGRATION_ID` | Yes | Paymob integration ID | From Paymob dashboard |
| `PAYMOB_IFRAME_ID` | Yes | Paymob iframe ID | From Paymob dashboard |
| `PAYMOB_HMAC_SECRET` | Yes | Webhook HMAC secret | From Paymob dashboard |
| `SENTRY_DSN` | Recommended | Error tracking | From Sentry |

### Frontend Apps

| Variable | Apps | Description |
|----------|------|-------------|
| `NEXT_PUBLIC_API_URL` | Dashboard, Landing, Super Admin | Backend API base URL |
| `NEXT_PUBLIC_API_BASE_URL` | Restaurant Website | Backend API base URL |
| `NEXT_PUBLIC_LANDING_URL` | Dashboard, Landing | Landing page URL |
| `NEXT_PUBLIC_DASHBOARD_URL` | Landing, Super Admin | Dashboard URL |
| `NEXT_PUBLIC_RESTAURANT_DOMAIN` | Dashboard, Landing, Restaurant | Base domain for tenants |

---

## Deployment Steps

### Step 0: CI + Smoke Validation (Before Touching Staging/Production)

Run these checks from your local machine (or CI) before any deployment:

```bash
# URL/domain hardening smoke checks (env templates + resolver tests + backend config validation)
python scripts/smoke_check_url_domain_config.py
```

Required CI workflows to enable as branch-protection checks:
- `Runtime Localhost Guardrail` (`.github/workflows/runtime-localhost-guardrail.yml`)
- `URL/Domain Config Smoke` (`.github/workflows/url-domain-config-smoke.yml`)

### Recommended Execution Order (Do These 8 In Order)

1. Set real production/staging domains and secrets (DNS + env vars + Paymob credentials).
2. Enable the CI workflows above as required checks.
3. Run `python scripts/smoke_check_url_domain_config.py` and fix failures.
4. Perform a full staging deployment (Steps 1-6 below).
5. Complete the Pre-Launch QA Checklist on staging.
6. Verify wildcard tenant routing and a themed storefront page on staging.
7. Rehearse the Rollback Procedure on staging.
8. Perform production deployment and post-deploy verification.

### Staging Deployment First (Recommended)

Use this guide for **staging first**, then repeat for production after QA signoff.

Minimum staging flow:

1. Create staging DNS records (or temporary hosts) for:
   - API
   - landing
   - dashboard
   - admin
   - one sample tenant subdomain
2. Prepare staging secrets/env values from:
   - `backend/.env.production.example`
   - `frontend/apps/*/.env.production.example`
3. Run local/CI preflight checks:
   - `python scripts/smoke_check_url_domain_config.py`
4. Execute Steps 1-6 below against staging infrastructure
5. Complete the Pre-Launch QA Checklist (on staging)
6. Rehearse rollback on staging
7. Promote same configuration pattern to production

### Step 1: Database Setup

```bash
# Create PostgreSQL database
createdb -U postgres pos_saas

# Run all migrations
cd backend
go run cmd/migrate/main.go

# Verify tables
psql -U postgres -d pos_saas -c "\dt"
```

### Step 2: Backend Deployment

```bash
# Build
cd backend
CGO_ENABLED=0 go build -o api ./cmd/api

# Set environment variables from secrets manager
export $(cat .env.production | xargs)

# Start with process manager (systemd, PM2, etc.)
./api
```

### Step 3: Frontend Deployment

```bash
cd frontend

# Install dependencies
pnpm install

# Build each app
cd apps/landing && pnpm build
cd ../dashboard && pnpm build
cd ../restaurant-website && pnpm build
cd ../super-admin && pnpm build

# Deploy to hosting (Vercel, AWS, etc.)
```

### Step 4: Reverse Proxy (Nginx Example)

`proxy_pass http://localhost:PORT` below refers to internal upstreams on the same server/VM.
In containers or multi-host deployments, replace those with internal service names/private IPs.
Public URLs should still be the production domains declared in each `server_name`.

```nginx
# Landing Page
server {
    server_name www.yourdomain.com;
    location / { proxy_pass http://localhost:3001; }
}

# Dashboard
server {
    server_name dashboard.yourdomain.com;
    location / { proxy_pass http://localhost:3002; }
}

# Restaurant Websites (wildcard)
server {
    server_name ~^(?<subdomain>.+)\.yourdomain\.com$;
    location / { proxy_pass http://localhost:3003; }
}

# Super Admin
server {
    server_name admin.yourdomain.com;
    location / { proxy_pass http://localhost:3004; }
}

# Backend API
server {
    server_name api.yourdomain.com;
    location / { proxy_pass http://localhost:8080; }
}
```

### Step 5: SSL/TLS Setup

```bash
# Using Certbot with wildcard certificate
certbot certonly --dns-cloudflare \
  -d yourdomain.com \
  -d *.yourdomain.com
```

### Step 6: Configure Paymob Webhook

1. Login to https://dashboard.paymob.com
2. Navigate to Developers → Webhooks
3. Set webhook URL: `https://api.yourdomain.com/api/v1/webhook/paymob`
4. Method: POST, Format: JSON

---

## Docker Deployment

```bash
# Build and start all services
docker-compose up -d

# Check health
curl https://api.yourdomain.com/api/v1/health

# View logs
docker-compose logs -f backend
```

---

## Pre-Launch QA Checklist (Complete On Staging Before Production)

### Backend
- [ ] `go build ./cmd/api` passes
- [ ] `go vet ./cmd/api` passes
- [ ] Health endpoint returns 200
- [ ] Plans endpoint returns valid JSON
- [ ] All migrations applied successfully
- [ ] CORS configured for production domains

### Frontend
- [ ] All 4 apps build without errors
- [ ] Landing page loads and redirects to `/en`
- [ ] Dashboard login works
- [ ] Restaurant website loads with theme data
- [ ] Super Admin dashboard accessible
- [ ] i18n works (EN + AR)
- [ ] Mobile responsive

### Payment System
- [ ] Paymob credentials configured
- [ ] Webhook URL set in Paymob dashboard
- [ ] Test transaction successful
- [ ] Subscription activation works
- [ ] Refund flow tested
- [ ] HMAC verification working

### Security
- [ ] HTTPS enabled on all domains
- [ ] JWT secret is strong and unique
- [ ] No hardcoded credentials in code
- [ ] CORS restricted to production domains
- [ ] Rate limiting enabled
- [ ] Database uses SSL connections

### Monitoring
- [ ] Sentry configured and receiving errors
- [ ] Webhook health dashboard operational
- [ ] Database connection monitoring active
- [ ] Log aggregation configured

### Post-Deploy Verification Command (Staging/Production)

Use this helper script after staging or production deploys:

```bash
python scripts/verify_deployment_endpoints.py \
  --api-url https://api.yourdomain.com \
  --landing-url https://www.yourdomain.com \
  --dashboard-url https://dashboard.yourdomain.com \
  --admin-url https://admin.yourdomain.com \
  --restaurant-url https://demo.yourdomain.com
```

Notes:
- `--restaurant-url` should be a real tenant subdomain in the target environment.
- Add `--insecure` only for temporary staging checks with self-signed certificates.
- Run this after **staging deploy** and again after **production deploy**.

---

## Rollback Procedure

```bash
# 1. Stop current version
docker-compose down

# 2. Restore database from backup
pg_restore -U postgres -d pos_saas backup.dump

# 3. Deploy previous version
git checkout previous-release-tag
docker-compose up -d

# 4. Verify
curl https://api.yourdomain.com/api/v1/health
```

---

## Support

- Payment Issues: See `docs-organized/11-PAYMOB-PAYMENT-INTEGRATION/`
- Theme System: See `docs-organized/01-THEME-SYSTEM/`
- API Reference: See `docs-organized/02-API-BACKEND/`
- Multi-Tenancy: See `docs-organized/05-INFRASTRUCTURE/MULTI_TENANCY/`
