# API Reference v2.0.0

**Version**: 2.0.0
**Base URL**: `https://api.example.com/api/v1`
**Auth**: Bearer Token (JWT)

---

## Quick Start

```bash
# Get authentication token
curl -X POST https://api.example.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass123"}'

# Response
{
  "token": "eyJhbGc...",
  "refreshToken": "eyJhbGc...",
  "expiresIn": 3600
}

# Use token in requests
curl -H "Authorization: Bearer eyJhbGc..." \
  https://api.example.com/api/v1/products
```

---

## Authentication

### Login
```
POST /auth/login
Body: { email, password }
Response: { token, refreshToken, expiresIn }
```

### Refresh Token
```
POST /auth/refresh
Body: { refreshToken }
Response: { token, expiresIn }
```

### Logout
```
POST /auth/logout
Headers: Authorization: Bearer {token}
Response: { success: true }
```

---

## Products API

### List Products
```
GET /products
Query: ?page=1&limit=20&category=electronics&sort=name
Response: { data: [], total, page, limit }
```

**Example**:
```bash
curl -H "Authorization: Bearer TOKEN" \
  'https://api.example.com/api/v1/products?limit=10'
```

### Get Product Detail
```
GET /products/:id
Response: { id, name, price, description, category, inStock, rating }
```

### Search Products
```
GET /products/search
Query: ?q=laptop&category=electronics&priceMin=500&priceMax=2000
Response: { data: [], total }
```

### Create Product (Admin)
```
POST /products
Body: { name, price, description, category, stock }
Response: { id, name, ... }
```

### Update Product (Admin)
```
PATCH /products/:id
Body: { name?, price?, description?, ... }
Response: { id, name, ... }
```

### Delete Product (Admin)
```
DELETE /products/:id
Response: { success: true }
```

---

## Orders API

### Current POS Order Endpoints (v1)
```
Public
GET /api/v1/public/orders/number/{number}/lookup
GET /api/v1/public/orders/{id}
GET /api/v1/public/orders/{id}/status
GET /api/v1/public/orders/{id}/track
POST /api/v1/public/orders
DELETE /api/v1/public/orders/{id}

Admin
GET /api/v1/admin/orders
GET /api/v1/admin/orders/{id}
PUT /api/v1/admin/orders/{id}/status
POST /api/v1/admin/orders/{id}/assign
DELETE /api/v1/admin/orders/{id}
GET /api/v1/admin/orders/stats
GET /api/v1/admin/orders/export
```

### List User Orders
```
GET /orders
Query: ?page=1&status=completed
Response: { data: [], total }
```

### Get Order Details
```
GET /orders/:id
Response: { id, userId, items, total, status, createdAt }
```

### Create Order
```
POST /orders
Body: { items: [{ productId, quantity }], shippingAddress }
Response: { id, total, status: "pending" }
```

### Update Order Status (Admin)
```
PATCH /orders/:id
Body: { status: "shipped" | "delivered" | "cancelled" }
Response: { id, status }
```

---

## Drivers API (Admin)

### Driver Management
```
GET /api/v1/admin/drivers
GET /api/v1/admin/drivers/available
GET /api/v1/admin/drivers/{id}
POST /api/v1/admin/drivers
PUT /api/v1/admin/drivers/{id}
PUT /api/v1/admin/drivers/{id}/status
DELETE /api/v1/admin/drivers/{id}
```

## Users API

### Get Current User
```
GET /users/me
Headers: Authorization: Bearer {token}
Response: { id, email, name, role, createdAt }
```

### Update Profile
```
PATCH /users/me
Body: { name?, email? }
Response: { id, email, name }
```

### Change Password
```
POST /users/me/password
Body: { currentPassword, newPassword }
Response: { success: true }
```

### Get User Preferences
```
GET /users/me/preferences
Response: { theme, notifications, language }
```

### Update Preferences
```
PATCH /users/me/preferences
Body: { theme?, notifications?, language? }
Response: { theme, notifications, language }
```

---

## Themes API

### List Themes
```
GET /themes
Query: ?page=1
Response: { data: [], total }
```

### Get Theme
```
GET /themes/:id
Response: { id, name, colors, typography, spacing }
```

### Create Theme
```
POST /themes
Body: { name, colors: {}, typography: {}, spacing: {} }
Response: { id, name, ... }
```

### Update Theme
```
PATCH /themes/:id
Body: { name?, colors?, typography?, ... }
Response: { id, name, ... }
```

### Delete Theme
```
DELETE /themes/:id
Response: { success: true }
```

### Export Theme
```
GET /themes/:id/export
Response: JSON file
```

### Import Theme
```
POST /themes/import
Body: FormData with theme.json file
Response: { id, name }
```

---

## Analytics API

### Dashboard Stats
```
GET /analytics/dashboard
Query: ?startDate=2024-01-01&endDate=2024-01-31
Response: {
  pageViews: 1500,
  uniqueVisitors: 450,
  bounceRate: 0.42,
  avgSessionDuration: 180
}
```

### Revenue Report
```
GET /analytics/revenue
Query: ?period=monthly
Response: { data: [], total }
```

### User Activity
```
GET /analytics/users
Query: ?metric=active
Response: { active: 1250, inactive: 300 }
```

---

## Error Responses

### Standard Error Format
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email is required",
    "details": {}
  }
}
```

### Common Error Codes

| Code | Status | Meaning |
|------|--------|---------|
| `AUTH_REQUIRED` | 401 | Login required |
| `INVALID_TOKEN` | 401 | Token expired/invalid |
| `PERMISSION_DENIED` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Invalid input data |
| `CONFLICT` | 409 | Resource already exists |
| `RATE_LIMITED` | 429 | Too many requests |
| `SERVER_ERROR` | 500 | Internal error |

---

## Rate Limiting

```
Rate Limit: 100 requests per minute per user
Rate Limit Window: 1 minute

Headers:
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1234567890
```

When limited:
```
Status: 429 Too Many Requests
Retry-After: 60
```

---

## Pagination

```
GET /products?page=1&limit=20
Response: {
  data: [],
  total: 150,
  page: 1,
  limit: 20,
  pages: 8,
  hasNext: true,
  hasPrev: false
}
```

---

## Filtering & Sorting

### Filtering
```
GET /products?category=electronics&inStock=true&priceMin=100&priceMax=2000
```

### Sorting
```
GET /products?sort=name                    # ASC
GET /products?sort=-price                  # DESC
GET /products?sort=rating,-createdAt       # Multiple
```

---

## Status Codes

| Code | Meaning |
|------|---------|
| `200` | OK - Success |
| `201` | Created - New resource created |
| `204` | No Content - Success, no body |
| `400` | Bad Request - Invalid input |
| `401` | Unauthorized - Auth required |
| `403` | Forbidden - No permission |
| `404` | Not Found - Resource missing |
| `429` | Too Many Requests - Rate limited |
| `500` | Server Error - Internal error |

---

## Caching Headers

**v2.0.0 includes intelligent caching**:

```
GET /products
Response Headers:
  Cache-Control: public, max-age=300     # 5-min cache
  ETag: "v1-abc123"
  Last-Modified: Mon, 15 Jan 2024 10:00:00 GMT
```

For cached responses:
```
GET /products
Request Headers:
  If-None-Match: "v1-abc123"

Response:
  304 Not Modified (no body sent, saves bandwidth!)
```

---

## Example: Complete Workflow

```bash
# 1. Login
TOKEN=$(curl -s -X POST https://api.example.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass"}' \
  | jq -r '.token')

# 2. Get products
curl -s -H "Authorization: Bearer $TOKEN" \
  https://api.example.com/api/v1/products?limit=5

# 3. Create order
PRODUCT_ID="prod-123"
curl -s -X POST https://api.example.com/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"items\":[{\"productId\":\"$PRODUCT_ID\",\"quantity\":2}]}"

# 4. Get order status
ORDER_ID="order-456"
curl -s -H "Authorization: Bearer $TOKEN" \
  https://api.example.com/api/v1/orders/$ORDER_ID
```

---

## SDK/Libraries

Coming soon in v2.0.0:

- JavaScript SDK
- Python SDK
- Go SDK
- Ruby SDK

---

## POS SaaS Phase 6 Endpoints (Accounting + Advanced HR)

All endpoints below are protected (`Authorization: Bearer <token>`) and tenant-scoped.

### Accounting (RBAC module 15)

```
GET/POST     /api/v1/accounting/accounts
GET/PUT/DEL  /api/v1/accounting/accounts/{id}
GET/POST     /api/v1/accounting/journal-entries
GET/PUT      /api/v1/accounting/journal-entries/{id}
POST         /api/v1/accounting/journal-entries/{id}/post
POST         /api/v1/accounting/journal-entries/{id}/void
GET          /api/v1/accounting/ledger
GET          /api/v1/accounting/trial-balance
GET/POST     /api/v1/accounting/bank-accounts
GET/PUT/DEL  /api/v1/accounting/bank-accounts/{id}
GET          /api/v1/accounting/bank-accounts/{id}/transactions
GET/POST     /api/v1/accounting/customers
GET/PUT/DEL  /api/v1/accounting/customers/{id}
GET          /api/v1/accounting/customers/{id}/statement
GET/POST     /api/v1/accounting/customer-invoices
GET/PUT      /api/v1/accounting/customer-invoices/{id}
POST         /api/v1/accounting/customer-invoices/{id}/send
POST         /api/v1/accounting/customer-invoices/{id}/void
GET          /api/v1/accounting/customer-invoices/overdue
GET/POST     /api/v1/accounting/supplier-invoices
GET/PUT      /api/v1/accounting/supplier-invoices/{id}
POST         /api/v1/accounting/supplier-invoices/{id}/approve
POST         /api/v1/accounting/supplier-invoices/{id}/void
GET/POST     /api/v1/accounting/payments/received
GET/POST     /api/v1/accounting/payments/made
```

### Suppliers (RBAC module 16)

```
GET/POST     /api/v1/suppliers
GET/PUT/DEL  /api/v1/suppliers/{id}
```

### Financial Reports + Tax (RBAC module 17)

```
GET/POST     /api/v1/accounting/tax/periods
GET/PUT/DEL  /api/v1/accounting/tax/periods/{id}
POST         /api/v1/accounting/tax/periods/{id}/close
GET          /api/v1/accounting/tax/vat-report
```

### HR Performance (RBAC module 18)

```
GET/POST     /api/v1/hr/performance
GET/PUT/DEL  /api/v1/hr/performance/{id}
POST         /api/v1/hr/performance/{id}/acknowledge
POST         /api/v1/hr/performance/{id}/complete
GET          /api/v1/hr/performance/employee/{employeeId}
```

### HR Training (RBAC module 19)

```
GET/POST     /api/v1/hr/training
GET/PUT/DEL  /api/v1/hr/training/{id}
GET          /api/v1/hr/training/expiring
```

### Contracts/Insurance/Resignations (RBAC module 20)

```
GET/POST     /api/v1/hr/insurance
GET/PUT/DEL  /api/v1/hr/insurance/{id}
GET          /api/v1/hr/insurance/expiring
GET/POST     /api/v1/hr/contracts
GET/PUT/DEL  /api/v1/hr/contracts/{id}
POST         /api/v1/hr/contracts/{id}/renew
GET          /api/v1/hr/contracts/expiring
GET/POST     /api/v1/hr/resignations
GET/PUT/DEL  /api/v1/hr/resignations/{id}
POST         /api/v1/hr/resignations/{id}/acknowledge
POST         /api/v1/hr/resignations/{id}/complete
GET          /api/v1/hr/resignations/{id}/settlement
```

### HR Reports Dashboard

```
GET          /api/v1/reports/hr/headcount
GET          /api/v1/reports/hr/turnover
GET          /api/v1/reports/hr/absenteeism
GET          /api/v1/reports/hr/overtime
GET          /api/v1/reports/hr/leave-balances
GET          /api/v1/reports/hr/training-compliance
```

---

## Support

- 📖 **Full Docs**: https://docs.example.com
- 🐛 **Issues**: https://github.com/example/issues
- 💬 **Chat**: https://slack.example.com
- 📧 **Email**: api-support@example.com
