# Super Admin Dashboard

The Super Admin Dashboard is a secure, privileged interface for SaaS administrators to manage the entire platform including all tenants, users, and perform advanced operations like impersonation.

## Features

- **Secure Authentication**: JWT-based authentication with separate super admin token signing
- **Tenant Management**: View and manage all platform tenants
- **User Management**: View all platform users across all tenants
- **Impersonation**: One-click access to any tenant dashboard without credentials
- **Audit Logging**: Comprehensive audit trail of all super admin actions
- **IP Whitelisting**: Restrict super admin access to specific IP addresses
- **Rate Limiting**: Built-in rate limiting for security
- **Session Tracking**: Monitor active super admin sessions

## Project Structure

```
src/
├── app/                      # Next.js app directory
│   ├── layout.tsx            # Root layout
│   ├── page.tsx              # Home page (redirects to login/dashboard)
│   ├── login/                # Authentication pages (Phase 7)
│   ├── dashboard/            # Dashboard pages (Phase 8)
│   └── globals.css           # Global styles
├── components/
│   ├── auth/                 # Auth-related components (Phase 7)
│   ├── layout/               # Layout components
│   ├── dashboard/            # Dashboard components (Phase 8)
│   └── impersonation/        # Impersonation UI (Phase 9)
├── hooks/                    # Custom React hooks
├── lib/                      # Utilities and API client
├── stores/                   # Zustand stores
└── types/                    # TypeScript types
```

## Getting Started

### Prerequisites

- Node.js 18+
- pnpm (or npm/yarn)

### Installation

From the root directory, run:

```bash
pnpm install
```

### Development

Run the super admin dashboard in development mode:

```bash
pnpm --filter @pos-saas/super-admin dev
```

The app will be available at `http://localhost:3004`

### Backend Required (No Docker Fallback)

If Docker Desktop is not running, start backend + embedded PostgreSQL directly:

```bash
cd backend
powershell -ExecutionPolicy Bypass -File .\scripts\start_local_no_docker.ps1
```

Then run super-admin normally.

### Build

Build for production:

```bash
pnpm --filter @pos-saas/super-admin build
```

### Start Production Server

```bash
pnpm --filter @pos-saas/super-admin start
```

## Configuration

Environment variables are configured in `.env.local`:

- `NEXT_PUBLIC_API_URL`: Backend API base URL
- `NEXT_PUBLIC_SUPER_ADMIN_API_URL`: Super admin endpoint URL
- `NEXT_PUBLIC_APP_NAME`: Application display name

## API Integration

The app uses the `superAdminAPI` client (in `src/lib/api.ts`) which automatically:
- Adds JWT tokens to all requests
- Handles 401 unauthorized responses
- Redirects to login on authentication failure

## Security

- All tokens are stored in localStorage with the super-admin prefix
- Tokens are automatically included in Authorization headers
- IP whitelist validation on the backend
- Rate limiting: 30 requests per minute
- Audit logging of all super admin actions
- Session tracking and management

## Testing

Run tests with:

```bash
pnpm --filter @pos-saas/super-admin test
```

Run tests once:

```bash
pnpm --filter @pos-saas/super-admin test:run
```

## Development Phases

- **Phase 1-5**: Backend implementation ✅
- **Phase 6**: Next.js app structure (Current)
- **Phase 7**: Authentication pages
- **Phase 8**: Dashboard and management pages
- **Phase 9**: Impersonation UI components

## Default Credentials

For development/testing:
- Email: `superadmin@pos-saas.com`
- Password: `SuperAdmin@123`

## Further Reading

See `/docs-organized/09-SUPER-ADMIN-DASHBOARD/` for detailed implementation documentation.
