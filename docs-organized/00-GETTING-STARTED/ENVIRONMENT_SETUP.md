# Environment Setup Guide

## Quick Start (Development)

### Prerequisites
- Go 1.24+
- Node.js 18+
- pnpm 8+
- PostgreSQL 15+
- Redis (optional, has in-memory fallback)

### 1. Backend

```bash
cd backend
cp .env.example .env
# Edit .env with your local database credentials

go build ./cmd/api
./api.exe  # or ./api on Linux/Mac
```

### 2. Frontend

```bash
cd frontend
pnpm install

# Start individual apps
cd apps/landing && pnpm dev      # http://localhost:3001
cd apps/dashboard && pnpm dev    # http://localhost:3002
cd apps/restaurant-website && pnpm dev  # http://demo.localhost:3003
cd apps/super-admin && pnpm dev  # http://localhost:3004
```

### 3. Database

```bash
# Create database
createdb -U postgres pos_saas

# Run migrations
cd backend && go run cmd/migrate/main.go
```

---

## Environment Variables

### Development Defaults

All apps have localhost defaults for development. No additional configuration needed.

### Production Setup

See `.env.production.example` files in each app:
- `backend/.env.production.example`
- `frontend/apps/dashboard/.env.production.example`
- `frontend/apps/landing/.env.production.example`
- `frontend/apps/restaurant-website/.env.production.example`
- `frontend/apps/super-admin/.env.production.example`

For full production deployment, see [PRODUCTION-DEPLOYMENT-GUIDE.md](../PRODUCTION-DEPLOYMENT-GUIDE.md).

---

## Port Assignments

| Port | Service | Description |
|------|---------|-------------|
| 3001 | Landing Page | Public marketing site + auth pages |
| 3002 | Dashboard | Tenant admin panel |
| 3003 | Restaurant Website | Customer-facing restaurant sites |
| 3004 | Super Admin | Platform administration |
| 8080 | Backend API | Go REST API server |
| 5432 | PostgreSQL | Database |
| 6379 | Redis | Cache (optional) |
