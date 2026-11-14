package middleware

import (
	"net/http"
	"pos-saas/internal/commerce"
)

// CommercialCheckoutGate rejects new charges before any payment record or provider call.
// Enabling checkout requires a separately reviewed commercial policy change.
func CommercialCheckoutGate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !commerce.Current().CheckoutAvailable {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":"Paid checkout is not available.","code":"checkout_unavailable"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
