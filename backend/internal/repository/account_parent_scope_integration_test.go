package repository

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"pos-saas/internal/domain"
)

func TestAccountRepositoryPostgresRejectsForeignScopeParentReferences(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	foreignParent := insertScopedAccount(t, db, fixture.tenantA, fixture.restaurantA2, "HRP-FOREIGN", "Foreign parent")
	localParent := insertScopedAccount(t, db, fixture.tenantA, fixture.restaurantA1, "HRP-LOCAL", "Local parent")
	targetID := insertScopedAccount(t, db, fixture.tenantA, fixture.restaurantA1, "HRP-TARGET", "Original account")
	repo := NewAccountRepository(db)

	for _, foreignID := range []int{foreignParent, insertScopedAccount(t, db, fixture.tenantB, fixture.restaurantB1, "HRP-TENANT", "Tenant B parent")} {
		newAccount := scopedAccount(fixture.tenantA, fixture.restaurantA1, "HRP-NEW", "Rejected child", &foreignID)
		if _, err := repo.CreateAccount(newAccount); err == nil {
			t.Fatalf("account creation accepted foreign parent %d", foreignID)
		}
		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM chart_of_accounts WHERE tenant_id = $1 AND restaurant_id = $2 AND account_code = $3`, fixture.tenantA, fixture.restaurantA1, "HRP-NEW").Scan(&count); err != nil {
			t.Fatalf("check rejected child creation: %v", err)
		}
		if count != 0 {
			t.Fatalf("foreign parent creation left %d child rows", count)
		}

		newName := "Forged parent update"
		if err := repo.UpdateAccount(&domain.Account{
			ID: targetID, TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
			AccountCode: "HRP-TARGET", AccountName: newName, AccountType: "asset",
			NormalBalance: "debit", Level: 1, Currency: "EGP", IsActive: true,
			ParentID: &foreignID,
		}); err == nil {
			t.Fatalf("account update accepted foreign parent %d", foreignID)
		}
		var name string
		var parent sql.NullInt64
		if err := db.QueryRow(`SELECT account_name, parent_id FROM chart_of_accounts WHERE id = $1`, targetID).Scan(&name, &parent); err != nil {
			t.Fatalf("read account after rejected parent update: %v", err)
		}
		if name != "Original account" || parent.Valid {
			t.Fatalf("rejected parent update changed target: name=%q parent=%v", name, parent)
		}
	}

	// A same-scope parent remains valid so the guard preserves the intended hierarchy.
	validChild := scopedAccount(fixture.tenantA, fixture.restaurantA1, "HRP-VALID", "Valid child", &localParent)
	if _, err := repo.CreateAccount(validChild); err != nil {
		t.Fatalf("create child under same-scope parent: %v", err)
	}
}

func TestAccountRepositoryPostgresRejectsMismatchedTenantRestaurant(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	account := scopedAccount(fixture.tenantA, fixture.restaurantB1, "HRP-MISMATCH", "Mismatched restaurant", nil)
	if _, err := NewAccountRepository(db).CreateAccount(account); err == nil {
		t.Fatal("account creation accepted a restaurant owned by another tenant")
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM chart_of_accounts WHERE tenant_id = $1 AND restaurant_id = $2 AND account_code = $3`, fixture.tenantA, fixture.restaurantB1, "HRP-MISMATCH").Scan(&count); err != nil {
		t.Fatalf("verify mismatched account creation: %v", err)
	}
	if count != 0 {
		t.Fatalf("mismatched tenant/restaurant account creation left %d rows", count)
	}
}

func TestJournalEntryRepositoryPostgresRejectsForeignScopeAccountLineAtomically(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	localAccount := insertScopedAccount(t, db, fixture.tenantA, fixture.restaurantA1, "JNL-LOCAL", "Local account")
	foreignAccount := insertScopedAccount(t, db, fixture.tenantA, fixture.restaurantA2, "JNL-FOREIGN", "Sibling account")
	unique := "scope-journal-" + fmt.Sprint(time.Now().UnixNano())
	entry := &domain.JournalEntry{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		EntryDate: time.Now(), Description: unique, EntryType: "manual", Status: domain.JournalStatusDraft,
		Lines: []domain.JournalEntryLine{
			{AccountID: localAccount, DebitAmount: 10, LineOrder: 1},
			{AccountID: foreignAccount, CreditAmount: 10, LineOrder: 2},
		},
	}
	if _, err := NewJournalEntryRepository(db).CreateJournalEntry(entry); err == nil {
		t.Fatal("journal entry accepted a line referencing a sibling restaurant account")
	}
	var entries, lines int
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entries WHERE tenant_id = $1 AND restaurant_id = $2 AND description = $3`, fixture.tenantA, fixture.restaurantA1, unique).Scan(&entries); err != nil {
		t.Fatalf("count rejected journal entry: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entry_lines jel JOIN journal_entries je ON je.id = jel.journal_entry_id WHERE je.tenant_id = $1 AND je.restaurant_id = $2 AND je.description = $3`, fixture.tenantA, fixture.restaurantA1, unique).Scan(&lines); err != nil {
		t.Fatalf("count lines from rejected journal entry: %v", err)
	}
	if entries != 0 || lines != 0 {
		t.Fatalf("rejected journal entry left side effects: entries=%d lines=%d", entries, lines)
	}
}

func TestAccountingPartyRepositoriesPostgresRejectForeignScopeAccounts(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	localAccount := insertScopedAccount(t, db, fixture.tenantA, fixture.restaurantA1, "ACCT-PARTY-LOCAL", "Local party account")
	foreignAccount := insertScopedAccount(t, db, fixture.tenantA, fixture.restaurantA2, "ACCT-PARTY-FOREIGN", "Sibling party account")
	localID := localAccount
	foreignID := foreignAccount

	supplierRepo := NewSupplierRepository(db)
	supplier := &domain.Supplier{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		CompanyName: "Scope supplier", Currency: "EGP", IsActive: true,
		PayableAccountID: &localID,
	}
	supplierID, err := supplierRepo.CreateSupplier(supplier)
	if err != nil {
		t.Fatalf("create in-scope supplier: %v", err)
	}
	supplier.PayableAccountID = &foreignID
	if err := supplierRepo.UpdateSupplier(supplier); err == nil {
		t.Fatal("supplier update accepted sibling-restaurant payable account")
	}
	assertAccountingPartyAccount(t, db, "suppliers", supplierID, localAccount)
	_, err = supplierRepo.CreateSupplier(&domain.Supplier{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		CompanyName: "Rejected scope supplier", Currency: "EGP", IsActive: true,
		PayableAccountID: &foreignID,
	})
	if err == nil {
		t.Fatal("supplier creation accepted sibling-restaurant payable account")
	}

	customerRepo := NewCustomerAccountRepository(db)
	customer := &domain.CustomerAccount{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		CustomerName: "Scope customer", CustomerType: "individual",
		Currency: "EGP", IsActive: true, ReceivableAccountID: &localID,
	}
	customerID, err := customerRepo.CreateCustomer(customer)
	if err != nil {
		t.Fatalf("create in-scope accounting customer: %v", err)
	}
	customer.ReceivableAccountID = &foreignID
	if err := customerRepo.UpdateCustomer(customer); err == nil {
		t.Fatal("customer update accepted sibling-restaurant receivable account")
	}
	assertAccountingPartyAccount(t, db, "customers_accounting", customerID, localAccount)
	_, err = customerRepo.CreateCustomer(&domain.CustomerAccount{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		CustomerName: "Rejected scope customer", CustomerType: "individual",
		Currency: "EGP", IsActive: true, ReceivableAccountID: &foreignID,
	})
	if err == nil {
		t.Fatal("customer creation accepted sibling-restaurant receivable account")
	}
}

func TestAccountingInvoiceRepositoriesPostgresRejectForeignPartiesWithoutSideEffects(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	supplierRepo := NewSupplierRepository(db)
	siblingSupplier, err := supplierRepo.CreateSupplier(&domain.Supplier{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA2,
		CompanyName: "Sibling invoice supplier", Currency: "EGP", IsActive: true,
	})
	if err != nil {
		t.Fatalf("create sibling supplier fixture: %v", err)
	}
	customerRepo := NewCustomerAccountRepository(db)
	siblingCustomer, err := customerRepo.CreateCustomer(&domain.CustomerAccount{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA2,
		CustomerName: "Sibling invoice customer", CustomerType: "individual", Currency: "EGP", IsActive: true,
	})
	if err != nil {
		t.Fatalf("create sibling customer fixture: %v", err)
	}
	unique := fmt.Sprintf("%d", time.Now().UnixNano())
	line := domain.InvoiceLine{Description: "Scope test line", Quantity: 1, UnitPrice: 10, LineTotal: 10}
	supplierInvoices := NewSupplierInvoiceRepository(db, NewJournalEntryRepository(db), NewAccountRepository(db))
	if _, err := supplierInvoices.CreateInvoice(&domain.SupplierInvoice{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		SupplierID: siblingSupplier, InvoiceDate: time.Now(), DueDate: time.Now(),
		InvoiceNumber: "ignored", Currency: "EGP", Status: domain.InvoiceStatusDraft,
		SupplierRef: "ws02-" + unique, Lines: []domain.InvoiceLine{line},
	}); err == nil {
		t.Fatal("supplier invoice accepted sibling-restaurant supplier")
	}
	customerInvoices := NewCustomerInvoiceRepository(db, NewJournalEntryRepository(db), NewAccountRepository(db))
	if _, err := customerInvoices.CreateInvoice(&domain.CustomerInvoice{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		CustomerID: siblingCustomer, InvoiceDate: time.Now(), DueDate: time.Now(),
		Currency: "EGP", Status: domain.InvoiceStatusDraft,
		Lines: []domain.InvoiceLine{line},
	}); err == nil {
		t.Fatal("customer invoice accepted sibling-restaurant customer")
	}
	for table := range map[string]struct{}{"invoices_supplier": {}, "invoices_customer": {}} {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE tenant_id = $1 AND restaurant_id = $2", table)
		if err := db.QueryRow(query, fixture.tenantA, fixture.restaurantA1).Scan(&count); err != nil {
			t.Fatalf("count %s after denied invoice creation: %v", table, err)
		}
		if count != 0 {
			t.Errorf("denied foreign-party request left %d rows in %s", count, table)
		}
	}
}

func TestAccountingInvoiceLinesPostgresRejectForeignScopeAccounts(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	foreignAccount := insertScopedAccount(t, db, fixture.tenantA, fixture.restaurantA2, "ACCT-LINE-FOREIGN", "Foreign invoice line account")
	foreignAccountID := foreignAccount
	supplierID, err := NewSupplierRepository(db).CreateSupplier(&domain.Supplier{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		CompanyName: "Local invoice supplier", Currency: "EGP", IsActive: true,
	})
	if err != nil {
		t.Fatalf("create local supplier: %v", err)
	}
	unique := fmt.Sprintf("%d", time.Now().UnixNano())
	line := domain.InvoiceLine{
		Description: "Foreign account line", Quantity: 1, UnitPrice: 10,
		LineTotal: 10, AccountID: &foreignAccountID,
	}
	if _, err := NewSupplierInvoiceRepository(db, NewJournalEntryRepository(db), NewAccountRepository(db)).CreateInvoice(&domain.SupplierInvoice{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		SupplierID: supplierID, InvoiceDate: time.Now(), DueDate: time.Now(),
		Currency: "EGP", Status: domain.InvoiceStatusDraft,
		SupplierRef: "line-" + unique, Lines: []domain.InvoiceLine{line},
	}); err == nil {
		t.Fatal("supplier invoice accepted a sibling-restaurant line account")
	}

	var customerID int
	if err := db.QueryRow(`
		INSERT INTO customers_accounting (tenant_id, restaurant_id, customer_code, customer_name)
		VALUES ($1, $2, $3, 'Line test customer') RETURNING id
	`, fixture.tenantA, fixture.restaurantA1, "LC"+unique[:12]).Scan(&customerID); err != nil {
		t.Fatalf("create customer line fixture: %v", err)
	}
	var invoiceID int
	if err := db.QueryRow(`
		INSERT INTO invoices_customer (tenant_id, restaurant_id, invoice_number, invoice_date, due_date, customer_id)
		VALUES ($1, $2, $3, CURRENT_DATE, CURRENT_DATE, $4) RETURNING id
	`, fixture.tenantA, fixture.restaurantA1, "LINE-INV-"+unique[:10], customerID).Scan(&invoiceID); err != nil {
		t.Fatalf("create customer invoice-line fixture: %v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin customer line transaction: %v", err)
	}
	lineErr := NewCustomerInvoiceRepository(db, NewJournalEntryRepository(db), NewAccountRepository(db)).insertCustomerInvoiceLinesTx(tx, invoiceID, fixture.tenantA, fixture.restaurantA1, []domain.InvoiceLine{line})
	if err := tx.Rollback(); err != nil {
		t.Fatalf("rollback customer invoice-line probe: %v", err)
	}
	if lineErr == nil {
		t.Fatal("customer invoice line accepted a sibling-restaurant account")
	}
	var lineCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM invoice_customer_lines WHERE invoice_id = $1`, invoiceID).Scan(&lineCount); err != nil {
		t.Fatalf("count customer invoice lines after denied write: %v", err)
	}
	if lineCount != 0 {
		t.Fatalf("denied customer line write left %d rows", lineCount)
	}
}

func assertAccountingPartyAccount(t *testing.T, db *sql.DB, table string, id, wantAccountID int) {
	t.Helper()
	var accountID sql.NullInt64
	query := fmt.Sprintf("SELECT payable_account_id FROM %s WHERE id = $1", table)
	if table == "customers_accounting" {
		query = "SELECT receivable_account_id FROM customers_accounting WHERE id = $1"
	}
	if err := db.QueryRow(query, id).Scan(&accountID); err != nil {
		t.Fatalf("read accounting party reference after rejected update: %v", err)
	}
	if !accountID.Valid || int(accountID.Int64) != wantAccountID {
		t.Fatalf("rejected update changed %s account reference: got=%v want=%d", table, accountID, wantAccountID)
	}
}

func insertScopedAccount(t *testing.T, db *sql.DB, tenantID, restaurantID int, code, name string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`
		INSERT INTO chart_of_accounts (
			tenant_id, restaurant_id, account_code, account_name, account_type,
			normal_balance, level, currency
		) VALUES ($1, $2, $3, $4, 'asset', 'debit', 1, 'EGP')
		RETURNING id
	`, tenantID, restaurantID, code, name).Scan(&id)
	if err != nil {
		t.Fatalf("insert test account: %v", err)
	}
	return id
}

func scopedAccount(tenantID, restaurantID int, code, name string, parentID *int) *domain.Account {
	return &domain.Account{
		TenantID: tenantID, RestaurantID: restaurantID, AccountCode: code,
		AccountName: name, AccountType: "asset", NormalBalance: "debit",
		ParentID: parentID, Level: 2, Currency: "EGP", IsActive: true,
	}
}
