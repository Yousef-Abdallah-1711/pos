//go:build integration

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
)

func TestEInvoiceHTTPRejectsForeignInvoiceAndListsScopedSubmissions(t *testing.T) {
	admin := openReportRLSDB(t, "POS_RLS_ADMIN_URL")
	appURL := os.Getenv("POS_RLS_APP_URL")
	if appURL == "" {
		t.Skip("POS_RLS_APP_URL is not set")
	}
	schema := fmt.Sprintf("einvoice_http_rls_%d", time.Now().UnixNano())
	if !regexp.MustCompile(`^einvoice_http_rls_[0-9]+$`).MatchString(schema) {
		t.Fatal("generated unsafe schema name")
	}
	ctx := context.Background()
	admin.SetMaxOpenConns(1)
	if _, err := admin.ExecContext(ctx, `CREATE SCHEMA `+schema); err != nil {
		t.Fatalf("create isolated schema: %v", err)
	}
	t.Cleanup(func() {
		if _, err := admin.ExecContext(context.Background(), `DROP SCHEMA `+schema+` CASCADE`); err != nil {
			t.Errorf("drop isolated schema: %v", err)
		}
	})
	if _, err := admin.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set admin search path: %v", err)
	}
	createReportRLSFixture(t, ctx, admin)
	for _, statement := range []string{
		`CREATE TABLE invoices_customer (id INT PRIMARY KEY, tenant_id BIGINT NOT NULL, restaurant_id BIGINT NOT NULL)`,
		`CREATE TABLE invoices_supplier (id INT PRIMARY KEY, tenant_id BIGINT NOT NULL, restaurant_id BIGINT NOT NULL)`,
		`CREATE TABLE e_invoice_submissions (id SERIAL PRIMARY KEY, tenant_id INT NOT NULL, restaurant_id INT NOT NULL, invoice_type VARCHAR(20) NOT NULL, invoice_id INT NOT NULL, eta_uuid VARCHAR(100), status VARCHAR(20) NOT NULL, submitted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, response_payload JSONB NOT NULL DEFAULT '{}'::jsonb, error_message TEXT, created_by INT NULL)`,
	} {
		if _, err := admin.ExecContext(ctx, statement); err != nil {
			t.Fatalf("create e-invoice HTTP fixture: %v", err)
		}
	}
	if _, err := admin.ExecContext(ctx, `GRANT USAGE ON SCHEMA `+schema+` TO pos_app; GRANT SELECT ON tenants, restaurants, users, invoices_customer, invoices_supplier TO pos_app; GRANT SELECT, INSERT, UPDATE, DELETE ON e_invoice_submissions TO pos_app; GRANT USAGE, SELECT ON SEQUENCE e_invoice_submissions_id_seq TO pos_app`); err != nil {
		t.Fatalf("grant e-invoice runtime access: %v", err)
	}
	contents, err := os.ReadFile("../../../migrations/171_rls_e_invoice_submissions.sql")
	if err != nil {
		t.Fatalf("read e-invoice migration: %v", err)
	}
	if _, err := admin.ExecContext(ctx, string(contents)); err != nil {
		t.Fatalf("apply e-invoice migration: %v", err)
	}
	appURL = reportURLWithSearchPath(t, appURL, schema)
	app := openReportRLSURL(t, appURL)
	app.SetMaxOpenConns(1)
	if _, err := app.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set runtime search path: %v", err)
	}
	assertReportRLSAppRole(t, ctx, app)
	for _, statement := range []string{
		`INSERT INTO tenants (id) VALUES (1), (2)`,
		`INSERT INTO restaurants (id, tenant_id) VALUES (11, 1), (12, 1), (21, 2)`,
		`INSERT INTO users (id, tenant_id, restaurant_id) VALUES (101, 1, 11), (103, 1, 12), (202, 2, 21)`,
		`INSERT INTO invoices_customer (id, tenant_id, restaurant_id) VALUES (1001, 1, 11), (1002, 1, 12), (2001, 2, 21)`,
		`INSERT INTO invoices_supplier (id, tenant_id, restaurant_id) VALUES (3001, 1, 11), (3002, 1, 12), (4001, 2, 21)`,
	} {
		if _, err := admin.ExecContext(ctx, statement); err != nil {
			t.Fatalf("seed e-invoice HTTP fixture: %v", err)
		}
	}
	accountingRepo := repository.NewAccountingAdvancedRepository(app, nil)
	handler := NewAccountingAdvancedHandler(app, accountingRepo, nil)
	if got := submitEInvoiceForScope(handler, 1, 11, 101, 1001); got != http.StatusOK {
		t.Fatalf("owned customer invoice submit status = %d, want %d", got, http.StatusOK)
	}
	if got := submitEInvoiceForScope(handler, 1, 11, 101, 2001); got != http.StatusNotFound {
		t.Fatalf("foreign invoice submit status = %d, want %d", got, http.StatusNotFound)
	}
	if got := submitEInvoiceForScope(handler, 1, 12, 103, 1002); got != http.StatusOK {
		t.Fatalf("sibling restaurant's own invoice submit status = %d, want %d", got, http.StatusOK)
	}
	items, err := accountingRepo.ListEInvoiceSubmissions(ctx, 1, 11)
	if err != nil || len(items) != 1 {
		t.Fatalf("foreign submit must leave Tenant A restaurant 11 with one row, got %d (err=%v)", len(items), err)
	}
	listRequest := eInvoiceRequest(http.MethodGet, "/api/v1/accounting/einvoice/submissions", 1, 11, 101)
	listResponse := httptest.NewRecorder()
	handler.ListEInvoiceSubmissions(listResponse, listRequest)
	var response struct {
		Data []map[string]interface{} `json:"data"`
	}
	if listResponse.Code != http.StatusOK || json.Unmarshal(listResponse.Body.Bytes(), &response) != nil || len(response.Data) != 1 {
		t.Fatalf("scoped e-invoice list response status=%d body=%s", listResponse.Code, listResponse.Body.String())
	}
}

func submitEInvoiceForScope(handler *AccountingAdvancedHandler, tenantID, restaurantID, userID, invoiceID int64) int {
	request := eInvoiceRequest(http.MethodPost, "/api/v1/accounting/einvoice/"+fmt.Sprint(invoiceID)+"/submit", tenantID, restaurantID, userID)
	request.SetPathValue("id", fmt.Sprint(invoiceID))
	response := httptest.NewRecorder()
	handler.SubmitEInvoice(response, request)
	return response.Code
}

func eInvoiceRequest(method, path string, tenantID, restaurantID, userID int64) *http.Request {
	request := httptest.NewRequest(method, path, nil)
	ctx := context.WithValue(request.Context(), middleware.TenantContextKey, tenantID)
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, restaurantID)
	ctx = context.WithValue(ctx, middleware.UserContextKey, &jwt.Claims{UserID: int(userID), TenantID: int(tenantID), RestaurantID: int(restaurantID)})
	return request.WithContext(ctx)
}
