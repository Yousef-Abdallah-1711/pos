//go:build integration

package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"pos-saas/tests/testdb"
)

// TestFinancialAcceptanceLiveAPIOrderToLedger is enabled only by the dedicated
// acceptance workflow, which starts the real API against its disposable stack.
func TestFinancialAcceptanceLiveAPIOrderToLedger(t *testing.T) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("POS_FINANCIAL_ACCEPTANCE_API_URL")), "/")
	if baseURL == "" {
		t.Skip("requires the dedicated disposable financial acceptance API process")
	}

	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable acceptance database: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(4)
	ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("connect acceptance database: %v", err)
	}

	fixture := seedLiveFinancialOrderFixture(t, db)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("create isolated cookie jar: %v", err)
	}
	client := &http.Client{Jar: jar, Timeout: 10 * time.Second}
	loginBody := map[string]string{"email": fixture.email, "password": fixture.password}
	loginStatus, loginResponse := acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/auth/login", loginBody)
	if loginStatus != http.StatusOK {
		t.Fatalf("synthetic owner login returned %d: %s", loginStatus, loginResponse)
	}

	orderBody := map[string]any{
		"customer_name":  "Financial Acceptance Guest",
		"customer_phone": "01000000000",
		"payment_method": "cash",
		"order_source":   "pos",
		"items":          []map[string]any{{"product_id": fixture.productID, "quantity": 2}},
	}
	orderStatus, orderResponse := acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/pos/orders", orderBody)
	if orderStatus != http.StatusCreated {
		t.Fatalf("authenticated order creation returned %d: %s", orderStatus, orderResponse)
	}
	var created struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(orderResponse, &created); err != nil || created.Data.ID <= 0 {
		t.Fatalf("order response has no persisted order id: %s", orderResponse)
	}

	paidBody := map[string]string{"payment_status": "paid"}
	paidPath := fmt.Sprintf("%s/api/v1/admin/orders/%d/payment-status", baseURL, created.Data.ID)
	paidStatus, paidResponse := acceptanceJSONRequest(t, client, http.MethodPut, paidPath, paidBody)
	if paidStatus != http.StatusOK {
		t.Fatalf("authenticated payment confirmation returned %d: %s", paidStatus, paidResponse)
	}
	repeatedPaidStatus, repeatedPaidResponse := acceptanceJSONRequest(t, client, http.MethodPut, paidPath, paidBody)
	if repeatedPaidStatus != http.StatusOK {
		t.Fatalf("idempotent payment confirmation returned %d: %s", repeatedPaidStatus, repeatedPaidResponse)
	}

	deadline := time.Now().Add(30 * time.Second)
	var eventID, eventStatus string
	var eventError sql.NullString
	var eventAttempts int
	for time.Now().Before(deadline) {
		err := db.QueryRowContext(ctx, `
			SELECT event_id, status, last_error, attempts FROM accounting_event_outbox
			WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3
			ORDER BY created_at DESC LIMIT 1`, fixture.tenantID, fixture.restaurantID, created.Data.ID).Scan(&eventID, &eventStatus, &eventError, &eventAttempts)
		if err == nil && eventStatus == "completed" {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
	if eventStatus != "completed" {
		t.Fatalf("accounting worker did not complete synthetic order event; last status=%q attempts=%d last_error=%q", eventStatus, eventAttempts, eventError.String)
	}
	var eventCount int
	if err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM accounting_event_outbox
		WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3 AND event_type = 'paid_order'`,
		fixture.tenantID, fixture.restaurantID, created.Data.ID).Scan(&eventCount); err != nil {
		t.Fatalf("count persisted paid-order accounting intents: %v", err)
	}
	if eventCount != 1 {
		t.Fatalf("repeated authenticated payment confirmation created %d accounting intents, want exactly one", eventCount)
	}

	var expectedSubtotalCents, expectedTaxCents, expectedDiscountCents, expectedDeliveryCents, expectedTotalCents int64
	err = db.QueryRowContext(ctx, `
		SELECT ROUND(COALESCE(SUM(oi.total_price), 0) * 100)::bigint,
		       ROUND(o.tax_amount * 100)::bigint,
		       ROUND(o.discount_amount * 100)::bigint,
		       ROUND(o.delivery_fee * 100)::bigint,
		       ROUND((COALESCE(SUM(oi.total_price), 0) + o.tax_amount - o.discount_amount + o.delivery_fee) * 100)::bigint
		FROM orders o JOIN order_items oi ON oi.order_id = o.id
		WHERE o.id = $1 AND o.tenant_id = $2 AND o.restaurant_id = $3
		GROUP BY o.id`, created.Data.ID, fixture.tenantID, fixture.restaurantID).Scan(
		&expectedSubtotalCents, &expectedTaxCents, &expectedDiscountCents, &expectedDeliveryCents, &expectedTotalCents)
	if err != nil {
		t.Fatalf("independently calculate source total: %v", err)
	}
	var orderTotalCents int64
	if err := db.QueryRowContext(ctx, `SELECT ROUND(total_amount * 100)::bigint FROM orders WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, created.Data.ID, fixture.tenantID, fixture.restaurantID).Scan(&orderTotalCents); err != nil {
		t.Fatalf("read persisted order total: %v", err)
	}
	if orderTotalCents != expectedTotalCents {
		t.Fatalf("source/order total cents=%d/%d, difference=%d", expectedTotalCents, orderTotalCents, orderTotalCents-expectedTotalCents)
	}

	var journalID int64
	var journalCount int
	if err := db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(MIN(id), 0)
		FROM journal_entries
		WHERE tenant_id = $1 AND restaurant_id = $2 AND reference_type = 'order'
		  AND reference_id = $3 AND entry_type = 'auto_sales' AND status = 'posted'`,
		fixture.tenantID, fixture.restaurantID, created.Data.ID).Scan(&journalCount, &journalID); err != nil {
		t.Fatalf("read posted sales journal: %v", err)
	}
	if journalCount != 1 || journalID <= 0 {
		t.Fatalf("posted sales journals=%d id=%d, want exactly one", journalCount, journalID)
	}
	var glCashCents, glRevenueCents, glTaxPayableCents, glRows int64
	if err := db.QueryRowContext(ctx, `
		SELECT COALESCE(ROUND(SUM(gl.debit_amount - gl.credit_amount) FILTER (WHERE a.account_code = '1110') * 100), 0)::bigint,
		       COALESCE(ROUND(SUM(gl.credit_amount - gl.debit_amount) FILTER (WHERE a.account_code = '4100') * 100), 0)::bigint,
		       COALESCE(ROUND(SUM(gl.credit_amount - gl.debit_amount) FILTER (WHERE a.account_code = '2120') * 100), 0)::bigint,
	       COUNT(*)
		FROM general_ledger gl JOIN chart_of_accounts a
		  ON a.id = gl.account_id AND a.tenant_id = gl.tenant_id AND a.restaurant_id = gl.restaurant_id
		WHERE gl.tenant_id = $1 AND gl.restaurant_id = $2 AND gl.journal_entry_id = $3`,
		fixture.tenantID, fixture.restaurantID, journalID).Scan(&glCashCents, &glRevenueCents, &glTaxPayableCents, &glRows); err != nil {
		t.Fatalf("reconcile general ledger: %v", err)
	}
	if glCashCents != expectedTotalCents || glRevenueCents != expectedSubtotalCents || glTaxPayableCents != expectedTaxCents || glRows != 3 {
		t.Fatalf("journal=%d source subtotal/tax/discount/delivery/total=%d/%d/%d/%d/%d actual cash/revenue/tax=%d/%d/%d GL rows=%d; differences=%d/%d/%d",
			journalID, expectedSubtotalCents, expectedTaxCents, expectedDiscountCents, expectedDeliveryCents, expectedTotalCents,
			glCashCents, glRevenueCents, glTaxPayableCents, glRows,
			glCashCents-expectedTotalCents, glRevenueCents-expectedSubtotalCents, glTaxPayableCents-expectedTaxCents)
	}
	verifyLiveFinancialReportsAndExports(t, liveFinancialReportAcceptance{
		baseURL: baseURL, client: client, db: db, ctx: ctx, fixture: fixture,
		orderID: created.Data.ID, sourceRevenueCents: expectedSubtotalCents,
		sourceGrossSalesCents: expectedTotalCents,
	})

	otherTenantID, otherRestaurantID, otherEmail, otherPassword := seedLiveFinancialOtherTenant(t, db)
	otherJar, _ := cookiejar.New(nil)
	otherClient := &http.Client{Jar: otherJar, Timeout: 10 * time.Second}
	otherLoginStatus, otherLoginResponse := acceptanceJSONRequest(t, otherClient, http.MethodPost, baseURL+"/api/v1/auth/login", map[string]string{"email": otherEmail, "password": otherPassword})
	if otherLoginStatus != http.StatusOK {
		t.Fatalf("other-tenant synthetic login returned %d: %s", otherLoginStatus, otherLoginResponse)
	}
	foreignReadPath := fmt.Sprintf("%s/api/v1/admin/orders/%d", baseURL, created.Data.ID)
	foreignStatus, foreignResponse := acceptanceJSONRequest(t, otherClient, http.MethodGet, foreignReadPath, nil)
	if foreignStatus != http.StatusNotFound {
		t.Fatalf("cross-tenant order read returned %d, want 404: %s", foreignStatus, foreignResponse)
	}
	if otherRestaurantID <= 0 {
		t.Fatal("other synthetic restaurant was not created")
	}

	t.Logf("financial_acceptance scenario=authenticated_order_to_ledger tenant_id=%d restaurant_id=%d order_id=%d source_subtotal_cents=%d source_tax_cents=%d source_discount_cents=%d source_delivery_cents=%d source_total_cents=%d actual_order_cents=%d order_difference_cents=%d event_id=%s event_count=%d repeated_payment_http_status=%d journal_id=%d cash_gl_cents=%d cash_difference_cents=%d revenue_gl_cents=%d revenue_difference_cents=%d tax_payable_gl_cents=%d tax_difference_cents=%d gl_rows=%d event_status=%s foreign_tenant_id=%d cross_tenant_read_status=%d",
		fixture.tenantID, fixture.restaurantID, created.Data.ID,
		expectedSubtotalCents, expectedTaxCents, expectedDiscountCents, expectedDeliveryCents, expectedTotalCents,
		orderTotalCents, orderTotalCents-expectedTotalCents, eventID, eventCount, repeatedPaidStatus, journalID,
		glCashCents, glCashCents-expectedTotalCents, glRevenueCents, glRevenueCents-expectedSubtotalCents,
		glTaxPayableCents, glTaxPayableCents-expectedTaxCents, glRows, eventStatus, otherTenantID, foreignStatus)
}

type liveFinancialFixture struct {
	tenantID, restaurantID, productID int64
	email, password                   string
}

func seedLiveFinancialOrderFixture(t *testing.T, db *sql.DB) liveFinancialFixture {
	t.Helper()
	runID := strconv.FormatInt(time.Now().UnixNano(), 10)
	name := "financial-live-" + runID
	var fixture liveFinancialFixture
	fixture.email = "owner-" + runID + "@example.test"
	fixture.password = "synthetic-financial-password-" + runID
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $2, $3) RETURNING id`, name, name, name+".example.test").Scan(&fixture.tenantID); err != nil {
		t.Fatalf("seed synthetic tenant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $3, 'Synthetic test address') RETURNING id`, fixture.tenantID, name, name).Scan(&fixture.restaurantID); err != nil {
		t.Fatalf("seed synthetic restaurant: %v", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(fixture.password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash synthetic owner password: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, role, status) VALUES ($1, $2, $3, $4, 'Synthetic Finance Owner', 'owner', 'active')`, fixture.tenantID, fixture.restaurantID, fixture.email, string(hash)); err != nil {
		t.Fatalf("seed synthetic owner: %v", err)
	}
	var categoryID int64
	if err := db.QueryRow(`INSERT INTO categories (tenant_id, restaurant_id, name, name_en, is_active) VALUES ($1, $2, 'Acceptance', 'Acceptance', TRUE) RETURNING id`, fixture.tenantID, fixture.restaurantID).Scan(&categoryID); err != nil {
		t.Fatalf("seed synthetic category: %v", err)
	}
	productName := `=HYPERLINK("https://example.invalid","Synthetic item")`
	if err := db.QueryRow(`INSERT INTO products (tenant_id, restaurant_id, category_id, name, name_en, price, cost, sku, status, is_available, track_inventory) VALUES ($1, $2, $3, $4, $4, 25, 0, $5, 'active', TRUE, FALSE) RETURNING id`, fixture.tenantID, fixture.restaurantID, categoryID, productName, "FA-"+runID).Scan(&fixture.productID); err != nil {
		t.Fatalf("seed synthetic product: %v", err)
	}
	for _, account := range []struct{ code, name, kind, normal string }{
		{"1110", "Synthetic Cash", "asset", "debit"},
		{"1120", "Synthetic Accounts Receivable", "asset", "debit"},
		{"1150", "Synthetic Input VAT", "asset", "debit"},
		{"1130", "Synthetic Inventory", "asset", "debit"},
		{"2101", "Synthetic Accounts Payable", "liability", "credit"},
		{"4100", "Synthetic Sales", "revenue", "credit"},
		{"2120", "Synthetic Tax Payable", "liability", "credit"},
		{"5110", "Synthetic Food COGS", "expense", "debit"},
	} {
		if _, err := db.Exec(`INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance) VALUES ($1, $2, $3, $4, $5, $6)`, fixture.tenantID, fixture.restaurantID, account.code, account.name, account.kind, account.normal); err != nil {
			t.Fatalf("seed account %s: %v", account.code, err)
		}
	}
	defaults := `{"cash":"1110","bank":"1110","card":"1110","online":"1110","wallet":"1110","receivable":"1120","sales":"4100","delivery_sales":"4120","tax":"2120","discount":"4150"}`
	if _, err := db.Exec(`INSERT INTO auto_posting_rules (tenant_id, restaurant_id, rule_code, rule_name, is_active, default_accounts) VALUES ($1, $2, 'ORDER_PAID', 'Acceptance order sales', TRUE, $3::jsonb)`, fixture.tenantID, fixture.restaurantID, defaults); err != nil {
		t.Fatalf("seed synthetic sales posting rule: %v", err)
	}
	return fixture
}

func seedLiveFinancialOtherTenant(t *testing.T, db *sql.DB) (int64, int64, string, string) {
	t.Helper()
	runID := strconv.FormatInt(time.Now().UnixNano(), 10)
	name := "financial-foreign-" + runID
	var tenantID, restaurantID int64
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $2, $3) RETURNING id`, name, name, name+".example.test").Scan(&tenantID); err != nil {
		t.Fatalf("seed second synthetic tenant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $3, 'Synthetic test address') RETURNING id`, tenantID, name, name).Scan(&restaurantID); err != nil {
		t.Fatalf("seed second synthetic restaurant: %v", err)
	}
	password := "synthetic-foreign-password-" + runID
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash second synthetic password: %v", err)
	}
	email := "owner-" + runID + "@example.test"
	if _, err := db.Exec(`INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, role, status) VALUES ($1, $2, $3, $4, 'Synthetic Foreign Owner', 'owner', 'active')`, tenantID, restaurantID, email, string(hash)); err != nil {
		t.Fatalf("seed second synthetic owner: %v", err)
	}
	return tenantID, restaurantID, email, password
}

func acceptanceJSONRequest(t *testing.T, client *http.Client, method, url string, body any) (int, []byte) {
	t.Helper()
	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("encode acceptance request: %v", err)
		}
		payload = bytes.NewReader(encoded)
	}
	request, err := http.NewRequestWithContext(t.Context(), method, url, payload)
	if err != nil {
		t.Fatalf("create acceptance request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "http://localhost:3002")
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("send acceptance request: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		t.Fatalf("read acceptance response: %v", err)
	}
	return response.StatusCode, responseBody
}
