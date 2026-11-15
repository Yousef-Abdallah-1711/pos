//go:build integration

package usecase

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"pos-saas/tests/testdb"
)

func TestReceivablesAgingReconcilesPositiveARAndExcludesVoids(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(8)

	fixture := createReceivablesAgingFixture(t, db, "primary")
	asOf := time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC)
	insertReceivablesAgingInvoice(t, db, fixture, asOf.AddDate(0, 0, -45), 85, "sent")
	insertReceivablesAgingInvoice(t, db, fixture, asOf.AddDate(0, 0, -45), 700, "void")
	foreign := createReceivablesAgingFixture(t, db, "foreign")
	insertReceivablesAgingInvoice(t, db, foreign, asOf.AddDate(0, 0, -45), 900, "sent")

	report, err := NewFinancialReportUseCase(db, nil).GenerateReceivablesAging(fixture.tenantID, fixture.restaurantID, asOf)
	if err != nil {
		t.Fatalf("generate synthetic receivables aging report: %v", err)
	}
	if report.AsOfDate != "2026-11-15" || len(report.Rows) != 1 {
		t.Fatalf("aging report scope/date = %s rows=%d; want one scoped row", report.AsOfDate, len(report.Rows))
	}
	bucket := report.Rows[0]
	if bucket.EntityID != fixture.customerID || bucket.Days60 != 85 || bucket.Total != 85 || report.Total != 85 {
		t.Fatalf("aging bucket = %+v total=%.2f; want one 85.00 balance in 60-day bucket", bucket, report.Total)
	}
	assertReceivablesAgingSQLBalance(t, db, fixture, asOf, report.Total)
}

type receivablesAgingFixture struct {
	tenantID, restaurantID, customerID int
}

func createReceivablesAgingFixture(t *testing.T, db *sql.DB, label string) receivablesAgingFixture {
	t.Helper()
	name := fmt.Sprintf("ar-aging-%s-%d", label, time.Now().UnixNano())
	fixture := receivablesAgingFixture{}
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1,$1,$2) RETURNING id`, name, name+".example.test").Scan(&fixture.tenantID); err != nil {
		t.Fatalf("create %s tenant: %v", label, err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id,name,slug,address) VALUES ($1,$2,$2,'test') RETURNING id`, fixture.tenantID, name).Scan(&fixture.restaurantID); err != nil {
		t.Fatalf("create %s restaurant: %v", label, err)
	}
	if err := db.QueryRow(`INSERT INTO customers_accounting (tenant_id,restaurant_id,customer_code,customer_name,current_balance) VALUES ($1,$2,$3,$4,85) RETURNING id`, fixture.tenantID, fixture.restaurantID, name, name).Scan(&fixture.customerID); err != nil {
		t.Fatalf("create %s customer: %v", label, err)
	}
	return fixture
}

func insertReceivablesAgingInvoice(t *testing.T, db *sql.DB, fixture receivablesAgingFixture, dueDate time.Time, balance float64, status string) {
	t.Helper()
	invoiceNumber := fmt.Sprintf("AR-%s-%d", status, time.Now().UnixNano())
	_, err := db.Exec(`
		INSERT INTO invoices_customer (tenant_id,restaurant_id,invoice_number,invoice_date,due_date,customer_id,total_amount,balance_due,currency,status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$7,'EGP',$8)
	`, fixture.tenantID, fixture.restaurantID, invoiceNumber, dueDate.AddDate(0, 0, -1), dueDate, fixture.customerID, balance, status)
	if err != nil {
		t.Fatalf("create synthetic %s invoice: %v", status, err)
	}
}

func assertReceivablesAgingSQLBalance(t *testing.T, db *sql.DB, fixture receivablesAgingFixture, asOf time.Time, actual float64) {
	t.Helper()
	var expected float64
	err := db.QueryRow(`
		SELECT COALESCE(SUM(balance_due),0) FROM invoices_customer
		WHERE tenant_id=$1 AND restaurant_id=$2 AND status <> 'void' AND balance_due > 0
		  AND due_date < $3 AND due_date >= $3 - INTERVAL '60 day'
	`, fixture.tenantID, fixture.restaurantID, asOf).Scan(&expected)
	if err != nil {
		t.Fatalf("independently calculate scoped receivable aging: %v", err)
	}
	if expected != 85 || actual != expected {
		t.Fatalf("SQL aging balance=%.2f report total=%.2f; expected 85.00 and zero difference", expected, actual)
	}
}
