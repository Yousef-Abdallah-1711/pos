# Security and data integrity

- Read-only PostgreSQL catalog query found row security enabled on 15 of 92 public ordinary tables (16.3%). Core tables `users`, `tenants`, `restaurants`, `orders`, `order_items`, `products`, `categories`, `inventory`, and `employees` have RLS disabled.
- The two database roles are non-superusers and cannot create databases or roles.
- The bounded API checks returned 401 for unauthenticated session/user settings/product/admin-order reads. The second disposable tenant received 403 for user 2 settings and 200 for its own settings.
- The API-level checks support those specific paths only; they do not prove every tenant-scoped handler or database query is isolated. RLS enforcement under `pos_app` for representative core tables, role-based mutation matrix, SQL injection, XSS, CSRF, rate limits, webhook signatures/idempotency, secret scanning, backups, and recovery were NOT RUN.
- No confirmed exploit was established by these checks. Database-level RLS coverage is limited and needs an authorized security review before asserting defense in depth.

