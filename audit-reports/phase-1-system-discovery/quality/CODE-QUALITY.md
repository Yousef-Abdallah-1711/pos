# Code Quality

Verified metrics (static):
- Go ~116.7k lines tracked; TS/TSX ~125.5k lines tracked.
- cmd/api/main.go is 1894 lines with 361 route registrations inline; god-file wiring.
- Oversized files: report_handler.go 2630, accounting_advanced_repo.go 2133, accounting_advanced_handler.go 1346, payment_usecase.go 1110 (see inventory).
- Backend logging: 1325 `log.Printf/Println/Fatal` calls in non-test Go code; 0 uses of slog/zap/logrus/zerolog.
- 254 `err.Error()` occurrences in handler/http non-test code; two response shapes (`{"error": msg}` via respondError, `{"success":false,"error":msg}` via writeJSON); 139 of the `"error"` literals sit beside `success`.
- Frontend: 374 `: any`/`as any` in non-test app code; 45 `console.log/error` in non-test app code (dashboard has a debugLog wrapper).
- Dashboard has two overlapping API layers: lib/api.ts (axios + interceptors) and lib/api/themeApi.ts with its own 401 handling.
- Dead/unwired code: backend/internal/handler/*.go (gin-based, 8 files incl. blockchain, personalization, predictive analytics, white_label), backend/internal/handlers/ (7 files), internal/models.disabled/, internal/scaling and internal/sharding largely tested but not referenced from main.go (not exhaustively verified; handler/handlers are verified not imported by cmd/api/main.go).
- Linter config: CODING_STANDARDS.md Appendix C recommends golangci-lint, but no .golangci.yml was found tracked; frontend `lint` scripts exist but nothing runs them in CI.
- Mixed framework styles: net/http 1.22 mux in main code versus gin in legacy code.
