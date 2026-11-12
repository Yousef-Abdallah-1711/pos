# Error Handling Lifecycle

## DB -> repo
Repositories wrap errors with `fmt.Errorf("failed to <op>: %w", err)` (e.g. repository/order_repo.go lines 113, 166, 227, 247, 378), preserving the driver message. No typed/sentinel errors for pq codes were found (grep for pq.Error / 23505 / SQLSTATE in handler/http returned nothing).

## usecase -> HTTP mapping
No central error mapper. Each handler decides ad hoc:
- `respondError(w, code, err.Error())` in auth_handler.go:181 returns `{"error": message}`.
- `writeJSON(w, 400, {"success": false, "error": err.Error()})` in accounting_advanced_handler.go (20+ sites) and account_handler.go:150.
- String comparison on messages: account_handler.go:237 `err.Error() == "account not found"`.
- 29 sites map a server failure to 500 with the raw `err.Error()` (e.g. admin_order_handler.go:177 ListOrders; admin_theme_handler_v2.go:80,119,153,184,329,361,394). Since repo errors wrap driver text, raw SQL/driver error strings can reach clients.
- Account creation failures return 400 with raw error (account_handler.go:150) even for internal DB errors.
- `http.Error` is not used in handler/http (0), while writeJSON ignores encoder errors; respondJSON ignores Marshal errors (auth_handler.go:175).
- Two further 401/403 helpers with `{success,error,code}` shape exist in utils.go (ForbiddenResponseWithCode).

## Panic/recovery
No recover middleware found (grep `recover()` in middleware and cmd/api). Server uses plain `http.ListenAndServe` (main.go ~1887) with no timeouts or graceful shutdown.

## API client -> UI
dashboard/src/lib/api.ts response interceptor: on 401, unless message contains "tenant" (substring match, `errorMessage.includes('tenant')`) or an impersonation_token exists in sessionStorage, it removes token/user/auth-storage from localStorage and redirects to landing `/auth/login`. Concerns: any 401 whose message includes "tenant" never logs out; no refresh-token flow; themeApi.ts has a separate 401 handler (line 76) duplicating logic. error-logger.ts supports an optional remoteEndpoint (fetch at line 381); no Sentry usage found.
