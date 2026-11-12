"""One-off codemod (Phase 2A C007): startup guard, auth rate limits, protected translations upload and metrics."""
P = "backend/cmd/api/main.go"
s = open(P, encoding="utf-8").read()

old = """		startDBHealthMonitor(db, 15*time.Second)
		defer db.Close()
	}
"""
new = """		startDBHealthMonitor(db, 15*time.Second)
		defer db.Close()
	}
	if config.IsProductionMode() {
		if err := productionStartupGuards(context.Background(), db); err != nil {
			log.Fatalf("Unsafe production configuration: %v", err)
		}
	}
"""
assert old in s
s = s.replace(old, new, 1)

rep = {
    'mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)':
        'mux.Handle("POST /api/v1/auth/register", middleware.RateLimitByIP(10)(http.HandlerFunc(authHandler.Register)))',
    'mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)':
        'mux.Handle("POST /api/v1/auth/login", middleware.RateLimitByIP(20)(http.HandlerFunc(authHandler.Login)))',
    'mux.HandleFunc("POST /api/v1/auth/login/confirm", authHandler.LoginConfirm)':
        'mux.Handle("POST /api/v1/auth/login/confirm", middleware.RateLimitByIP(20)(http.HandlerFunc(authHandler.LoginConfirm)))',
    'mux.HandleFunc("POST /api/v1/auth/token-exchange", authHandler.TokenExchange)':
        'mux.Handle("POST /api/v1/auth/token-exchange", middleware.RateLimitByIP(20)(http.HandlerFunc(authHandler.TokenExchange)))',
    'mux.HandleFunc("POST /api/v1/auth/forgot-password", authHandler.ForgotPassword)':
        'mux.Handle("POST /api/v1/auth/forgot-password", middleware.RateLimitByIP(5)(http.HandlerFunc(authHandler.ForgotPassword)))',
    'mux.HandleFunc("POST /api/v1/super-admin/auth/login", superAdminHandler.Login)':
        'mux.Handle("POST /api/v1/super-admin/auth/login", middleware.RateLimitByIP(10)(http.HandlerFunc(superAdminHandler.Login)))',
    'mux.HandleFunc("POST /api/v1/translations/{language}/upload", translationHandler.UploadTranslations)':
        'mux.Handle("POST /api/v1/translations/{language}/upload", middleware.SuperAdminAuthMiddleware(superAdminTokenService)(http.HandlerFunc(translationHandler.UploadTranslations)))',
    'mux.HandleFunc("GET /api/v1/health/metrics", func(w http.ResponseWriter, r *http.Request) {':
        'mux.Handle("GET /api/v1/health/metrics", middleware.RequireInternalToken("METRICS_TOKEN", config.IsProductionMode())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {',
}
for a, b in rep.items():
    assert s.count(a) == 1, (a, s.count(a))
    s = s.replace(a, b, 1)

# close the metrics handler: `})` -> `})))`
i = s.index('mux.Handle("GET /api/v1/health/metrics"')
j = s.index("\n\t})\n", i)
s = s[:j] + "\n\t})))\n" + s[j + len("\n\t})\n"):]

open(P, "w", encoding="utf-8", newline="").write(s)
print("ok")
