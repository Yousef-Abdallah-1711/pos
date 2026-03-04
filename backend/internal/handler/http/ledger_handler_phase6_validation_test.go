package http

import (
	"context"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"pos-saas/internal/middleware"
)

func TestLedgerHandler_GetLedgerRejectsInvalidInputs(t *testing.T) {
	handler := NewLedgerHandler(nil)

	tests := []struct {
		name   string
		target string
	}{
		{
			name:   "missing account_id",
			target: "/api/v1/accounting/ledger",
		},
		{
			name:   "invalid start date",
			target: "/api/v1/accounting/ledger?account_id=1101&start=2026-31-01",
		},
		{
			name:   "invalid limit",
			target: "/api/v1/accounting/ledger?account_id=1101&limit=0",
		},
		{
			name:   "invalid cursor",
			target: "/api/v1/accounting/ledger?account_id=1101&cursor=-3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(stdhttp.MethodGet, tc.target, nil)
			req = withLedgerTenantContext(req)

			rec := httptest.NewRecorder()
			handler.GetLedger(rec, req)

			if rec.Code != stdhttp.StatusBadRequest {
				t.Fatalf("expected status 400, got %d for target %s", rec.Code, tc.target)
			}
		})
	}
}

func TestLedgerHandler_GetTrialBalanceRejectsInvalidAsOfDate(t *testing.T) {
	handler := NewLedgerHandler(nil)
	req := httptest.NewRequest(stdhttp.MethodGet, "/api/v1/accounting/trial-balance?as_of=invalid", nil)
	req = withLedgerTenantContext(req)

	rec := httptest.NewRecorder()
	handler.GetTrialBalance(rec, req)

	if rec.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func withLedgerTenantContext(req *stdhttp.Request) *stdhttp.Request {
	ctx := context.WithValue(req.Context(), middleware.TenantContextKey, int64(1))
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(2))
	return req.WithContext(ctx)
}
