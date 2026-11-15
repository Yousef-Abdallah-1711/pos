package repository

import (
	"database/sql"
	"math"
	"testing"
	"time"

	"pos-saas/internal/domain"
)

func TestCustomerInvoiceLifecycleKeepsSubledgerAndReceivableGLAligned(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	receivableAccountID := seedCustomerInvoiceBalanceAccounts(t, db, fixture.tenantA, fixture.restaurantA1)
	customerRepo := NewCustomerAccountRepository(db)
	customer := &domain.CustomerAccount{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		CustomerName: "Synthetic balance acceptance", CustomerType: "individual",
		Currency: "EGP", IsActive: true, ReceivableAccountID: &receivableAccountID,
	}
	customerID, err := customerRepo.CreateCustomer(customer)
	if err != nil {
		t.Fatalf("create synthetic customer: %v", err)
	}

	invoices := NewCustomerInvoiceRepository(db, NewJournalEntryRepository(db), NewAccountRepository(db))
	firstInvoiceID := createBalanceAcceptanceInvoice(t, invoices, fixture.tenantA, fixture.restaurantA1, customerID, 100)
	assertCustomerBalance(t, db, fixture.tenantA, fixture.restaurantA1, customerID, 100)
	settleCustomerInvoicePartially(t, db, fixture.tenantA, fixture.restaurantA1, customerID, firstInvoiceID)
	assertCustomerBalance(t, db, fixture.tenantA, fixture.restaurantA1, customerID, 60)

	secondInvoiceID := createBalanceAcceptanceInvoice(t, invoices, fixture.tenantA, fixture.restaurantA1, customerID, 50)
	assertCustomerBalance(t, db, fixture.tenantA, fixture.restaurantA1, customerID, 110)
	updateAndVoidCustomerInvoice(t, invoices, fixture.tenantA, fixture.restaurantA1, secondInvoiceID, customerID, db)
	assertCustomerBalance(t, db, fixture.tenantA, fixture.restaurantA1, customerID, 60)
	assertCustomerReceivableReconciliation(t, db, fixture.tenantA, fixture.restaurantA1, customerID, receivableAccountID)
}

func seedCustomerInvoiceBalanceAccounts(t *testing.T, db *sql.DB, tenantID, restaurantID int) int {
	t.Helper()
	receivable := insertCustomerInvoiceAccount(t, db, tenantID, restaurantID, "1120", "Accounts Receivable", "asset", "debit")
	insertCustomerInvoiceAccount(t, db, tenantID, restaurantID, "4100", "Sales Revenue", "revenue", "credit")
	insertCustomerInvoiceAccount(t, db, tenantID, restaurantID, "2120", "Output VAT", "liability", "credit")
	insertCustomerInvoiceAccount(t, db, tenantID, restaurantID, "1110", "Cash", "asset", "debit")
	return receivable
}

func insertCustomerInvoiceAccount(t *testing.T, db *sql.DB, tenantID, restaurantID int, code, name, accountType, normalBalance string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name,
			account_type, normal_balance, level, currency)
		VALUES ($1,$2,$3,$4,$5,$6,1,'EGP') RETURNING id`,
		tenantID, restaurantID, code, name, accountType, normalBalance).Scan(&id)
	if err != nil {
		t.Fatalf("create synthetic %s account: %v", code, err)
	}
	return id
}

func createBalanceAcceptanceInvoice(t *testing.T, repo *CustomerInvoiceRepository, tenantID, restaurantID, customerID int, amount float64) int {
	t.Helper()
	now := time.Now()
	invoice := &domain.CustomerInvoice{
		TenantID: tenantID, RestaurantID: restaurantID, CustomerID: customerID,
		InvoiceDate: now, DueDate: now.AddDate(0, 0, 15), Currency: "EGP", Status: domain.InvoiceStatusDraft,
		Lines: []domain.InvoiceLine{{Description: "Synthetic receivable", Quantity: 1, UnitPrice: amount, LineTotal: amount}},
	}
	id, err := repo.CreateInvoice(invoice)
	if err != nil {
		t.Fatalf("create synthetic customer invoice for %.2f: %v", amount, err)
	}
	return id
}

func settleCustomerInvoicePartially(t *testing.T, db *sql.DB, tenantID, restaurantID, customerID, invoiceID int) {
	t.Helper()
	invoiceIDCopy := invoiceID
	_, err := NewPaymentAccountingRepository(db, NewJournalEntryRepository(db), NewAccountRepository(db), nil).CreatePaymentReceived(&domain.PaymentReceived{
		TenantID: tenantID, RestaurantID: restaurantID, CustomerID: customerID,
		InvoiceID: &invoiceIDCopy, PaymentDate: time.Now(), Amount: 40, PaymentMethod: "cash",
	})
	if err != nil {
		t.Fatalf("record synthetic partial customer payment: %v", err)
	}
}

func updateAndVoidCustomerInvoice(t *testing.T, repo *CustomerInvoiceRepository, tenantID, restaurantID, invoiceID, customerID int, db *sql.DB) {
	t.Helper()
	invoice, err := repo.GetInvoiceByID(tenantID, restaurantID, invoiceID)
	if err != nil || invoice == nil {
		t.Fatalf("read synthetic invoice before edit: invoice=%v err=%v", invoice, err)
	}
	invoice.Lines[0].UnitPrice = 70
	invoice.Lines[0].LineTotal = 70
	if err := repo.UpdateInvoice(invoice); err != nil {
		t.Fatalf("update unpaid invoice amount: %v", err)
	}
	assertCustomerBalance(t, db, tenantID, restaurantID, customerID, 130)
	if err := repo.VoidInvoice(tenantID, restaurantID, invoiceID, 0, "synthetic acceptance cleanup"); err != nil {
		t.Fatalf("void unpaid invoice: %v", err)
	}
}

func assertCustomerBalance(t *testing.T, db *sql.DB, tenantID, restaurantID, customerID int, expected float64) {
	t.Helper()
	var actual float64
	err := db.QueryRow(`SELECT current_balance FROM customers_accounting WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`,
		customerID, tenantID, restaurantID).Scan(&actual)
	if err != nil {
		t.Fatalf("read synthetic customer balance: %v", err)
	}
	if math.Abs(actual-expected) > 0.001 {
		t.Fatalf("customer current_balance=%.2f, expected %.2f", actual, expected)
	}
}

func assertCustomerReceivableReconciliation(t *testing.T, db *sql.DB, tenantID, restaurantID, customerID, receivableAccountID int) {
	t.Helper()
	var customerBalance, openInvoiceBalance, receivableGL float64
	err := db.QueryRow(`
		SELECT c.current_balance,
		       (SELECT COALESCE(SUM(i.balance_due),0) FROM invoices_customer i
		        WHERE i.tenant_id=c.tenant_id AND i.restaurant_id=c.restaurant_id
	          AND i.customer_id=c.id AND i.status <> 'void'),
		       (SELECT COALESCE(SUM(g.debit_amount-g.credit_amount),0) FROM general_ledger g
		        WHERE g.tenant_id=c.tenant_id AND g.restaurant_id=c.restaurant_id AND g.account_id=$4)
		FROM customers_accounting c WHERE c.id=$1 AND c.tenant_id=$2 AND c.restaurant_id=$3`,
		customerID, tenantID, restaurantID, receivableAccountID).Scan(&customerBalance, &openInvoiceBalance, &receivableGL)
	if err != nil {
		t.Fatalf("read independent customer, invoice, and receivable balances: %v", err)
	}
	if math.Abs(customerBalance-openInvoiceBalance) > 0.001 || math.Abs(customerBalance-receivableGL) > 0.001 {
		t.Fatalf("customer/AP reconciliation mismatch: customer=%.2f open_invoices=%.2f receivable_GL=%.2f", customerBalance, openInvoiceBalance, receivableGL)
	}
}
