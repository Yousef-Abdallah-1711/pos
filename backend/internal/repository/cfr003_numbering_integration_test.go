//go:build integration

package repository

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCFR003PostgresBaselineMaxPlusOneCollidesUnderInterleaving(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	entryDate := time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC)
	firstTx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin first legacy allocation: %v", err)
	}
	defer firstTx.Rollback()
	secondTx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin second legacy allocation: %v", err)
	}
	defer secondTx.Rollback()

	firstNumber := cfr003LegacyNextJournalNumber(t, firstTx, fixture, entryDate)
	secondNumber := cfr003LegacyNextJournalNumber(t, secondTx, fixture, entryDate)
	if firstNumber != secondNumber {
		t.Fatalf("baseline allocators returned %q and %q; expected the uncommitted MAX()+1 collision", firstNumber, secondNumber)
	}

	insertCfr003Draft(t, firstTx, fixture, firstNumber, entryDate)
	if err := firstTx.Commit(); err != nil {
		t.Fatalf("commit first legacy journal: %v", err)
	}

	if err := insertCfr003DraftError(secondTx, fixture, secondNumber, entryDate); err == nil {
		t.Fatal("second legacy journal unexpectedly committed a duplicate document number")
	}
	if err := secondTx.Rollback(); err != nil && err != sql.ErrTxDone {
		t.Fatalf("rollback rejected duplicate allocation: %v", err)
	}
}

func TestCFR003PostgresConcurrentDocumentNumbersAreUniqueBySeries(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	var customerID, supplierID int
	if err := db.QueryRow(`
		INSERT INTO customers_accounting (tenant_id, restaurant_id, customer_code, customer_name)
		VALUES ($1, $2, 'CFR003-C', 'CFR-003 numbering customer') RETURNING id
	`, fixture.tenantID, fixture.restaurantID).Scan(&customerID); err != nil {
		t.Fatalf("create customer fixture: %v", err)
	}
	if err := db.QueryRow(`
		INSERT INTO suppliers (tenant_id, restaurant_id, supplier_code, company_name)
		VALUES ($1, $2, 'CFR003-S', 'CFR-003 numbering supplier') RETURNING id
	`, fixture.tenantID, fixture.restaurantID).Scan(&supplierID); err != nil {
		t.Fatalf("create supplier fixture: %v", err)
	}

	entryDate := time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC)
	date := entryDate.Format("2006-01-02")
	period := entryDate.Format("2006-01")
	tests := []struct {
		name      string
		table     string
		column    string
		prefix    string
		allocate  func(*sql.Tx) (string, error)
		insertDoc func(*sql.Tx, string) error
	}{
		{
			name:   "journal entries",
			table:  "journal_entries",
			column: "entry_number",
			prefix: "JE-2026-%",
			allocate: func(tx *sql.Tx) (string, error) {
				return NewJournalEntryRepository(db).generateEntryNumberTx(tx, fixture.tenantID, fixture.restaurantID, entryDate)
			},
			insertDoc: func(tx *sql.Tx, number string) error {
				_, err := tx.Exec(`INSERT INTO journal_entries (tenant_id, restaurant_id, entry_number, entry_date, fiscal_period, description, entry_type, status) VALUES ($1, $2, $3, $4, $5, 'CFR-003 concurrent number', 'manual', 'draft')`, fixture.tenantID, fixture.restaurantID, number, date, period)
				return err
			},
		},
		{
			name:   "customer invoices",
			table:  "invoices_customer",
			column: "invoice_number",
			prefix: "INV-2026-%",
			allocate: func(tx *sql.Tx) (string, error) {
				return NewCustomerInvoiceRepository(db, nil, nil).generateInvoiceNumberTx(tx, fixture.tenantID, fixture.restaurantID, entryDate)
			},
			insertDoc: func(tx *sql.Tx, number string) error {
				_, err := tx.Exec(`INSERT INTO invoices_customer (tenant_id, restaurant_id, invoice_number, invoice_date, due_date, customer_id, status) VALUES ($1, $2, $3, $4, $4, $5, 'draft')`, fixture.tenantID, fixture.restaurantID, number, date, customerID)
				return err
			},
		},
		{
			name:   "supplier invoices",
			table:  "invoices_supplier",
			column: "invoice_number",
			prefix: "PINV-2026-%",
			allocate: func(tx *sql.Tx) (string, error) {
				return NewSupplierInvoiceRepository(db, nil, nil).generateInvoiceNumberTx(tx, fixture.tenantID, fixture.restaurantID, entryDate)
			},
			insertDoc: func(tx *sql.Tx, number string) error {
				_, err := tx.Exec(`INSERT INTO invoices_supplier (tenant_id, restaurant_id, invoice_number, invoice_date, due_date, supplier_id, status) VALUES ($1, $2, $3, $4, $4, $5, 'draft')`, fixture.tenantID, fixture.restaurantID, number, date, supplierID)
				return err
			},
		},
		{
			name:   "received payments",
			table:  "payments_received",
			column: "payment_number",
			prefix: "PMT-R-2026-%",
			allocate: func(tx *sql.Tx) (string, error) {
				return NewPaymentAccountingRepository(db, nil, nil, nil).generatePaymentNumberTx(tx, "received", fixture.tenantID, fixture.restaurantID, entryDate)
			},
			insertDoc: func(tx *sql.Tx, number string) error {
				_, err := tx.Exec(`INSERT INTO payments_received (tenant_id, restaurant_id, payment_number, payment_date, customer_id, amount, payment_method) VALUES ($1, $2, $3, $4, $5, 1, 'cash')`, fixture.tenantID, fixture.restaurantID, number, date, customerID)
				return err
			},
		},
		{
			name:   "made payments",
			table:  "payments_made",
			column: "payment_number",
			prefix: "PMT-M-2026-%",
			allocate: func(tx *sql.Tx) (string, error) {
				return NewPaymentAccountingRepository(db, nil, nil, nil).generatePaymentNumberTx(tx, "made", fixture.tenantID, fixture.restaurantID, entryDate)
			},
			insertDoc: func(tx *sql.Tx, number string) error {
				_, err := tx.Exec(`INSERT INTO payments_made (tenant_id, restaurant_id, payment_number, payment_date, supplier_id, amount, payment_method) VALUES ($1, $2, $3, $4, $5, 1, 'cash')`, fixture.tenantID, fixture.restaurantID, number, date, supplierID)
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			const count = 10
			start := make(chan struct{})
			numbers := make(chan string, count)
			errs := make(chan error, count)
			var workers sync.WaitGroup
			for range count {
				workers.Add(1)
				go func() {
					defer workers.Done()
					<-start
					tx, err := db.Begin()
					if err != nil {
						errs <- err
						return
					}
					defer tx.Rollback()
					number, err := test.allocate(tx)
					if err == nil {
						err = test.insertDoc(tx, number)
					}
					if err == nil {
						err = tx.Commit()
					}
					if err != nil {
						errs <- err
						return
					}
					numbers <- number
				}()
			}
			close(start)
			workers.Wait()
			close(numbers)
			close(errs)
			for err := range errs {
				t.Errorf("concurrent %s creation failed: %v", test.name, err)
			}
			seen := make(map[string]struct{}, count)
			for number := range numbers {
				if _, exists := seen[number]; exists {
					t.Errorf("duplicate %s number allocated: %s", test.name, number)
				}
				seen[number] = struct{}{}
			}
			if len(seen) != count {
				t.Fatalf("unique %s numbers = %d, want %d", test.name, len(seen), count)
			}
			query := fmt.Sprintf("SELECT COUNT(DISTINCT %s) FROM %s WHERE tenant_id = $1 AND restaurant_id = $2 AND %s LIKE $3", test.column, test.table, test.column)
			var persisted int
			if err := db.QueryRow(query, fixture.tenantID, fixture.restaurantID, test.prefix).Scan(&persisted); err != nil {
				t.Fatalf("count persisted %s numbers: %v", test.name, err)
			}
			if persisted != count {
				t.Fatalf("persisted unique %s numbers = %d, want %d", test.name, persisted, count)
			}
		})
	}
}

func TestCFR003PostgresJournalNumberAllocationRollsBackWithDocument(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	journalRepo := NewJournalEntryRepository(db)
	entry := fixture.entry()
	entry.ReferenceID = nil
	entry.ReferenceType = ""
	entry.EntryNumber = ""

	rolledBack, err := db.Begin()
	if err != nil {
		t.Fatalf("begin rollback fixture: %v", err)
	}
	if _, err := journalRepo.CreateJournalEntryTx(rolledBack, entry); err != nil {
		t.Fatalf("allocate number and create journal in rollback transaction: %v", err)
	}
	if err := rolledBack.Rollback(); err != nil {
		t.Fatalf("rollback number allocation transaction: %v", err)
	}

	retry := fixture.entry()
	retry.ReferenceID = nil
	retry.ReferenceType = ""
	retry.EntryNumber = ""
	id, err := journalRepo.CreateJournalEntry(retry)
	if err != nil {
		t.Fatalf("retry journal creation after rollback: %v", err)
	}
	if id <= 0 || retry.EntryNumber == "" {
		t.Fatalf("retry did not persist a journal number: id=%d number=%q", id, retry.EntryNumber)
	}
}

func TestCFR003PostgresConcurrentJournalPostingPreservesRunningBalance(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	const count = 12
	start := make(chan struct{})
	errs := make(chan error, count)
	numbers := make(chan string, count)
	var workers sync.WaitGroup
	for index := range count {
		workers.Add(1)
		go func(referenceID int) {
			defer workers.Done()
			<-start
			entry := fixture.entry()
			entry.ReferenceID = &referenceID
			entry.ReferenceNumber = fmt.Sprintf("CFR003-concurrent-%d", referenceID)
			id, err := fixture.ledger.CreateAndPostAutoEntry(entry, 0)
			if err != nil {
				errs <- err
				return
			}
			var number string
			if err := db.QueryRow(`SELECT entry_number FROM journal_entries WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, fixture.tenantID, fixture.restaurantID).Scan(&number); err != nil {
				errs <- err
				return
			}
			numbers <- number
		}(41000 + index)
	}
	close(start)
	workers.Wait()
	close(errs)
	close(numbers)
	for err := range errs {
		t.Errorf("concurrent journal posting failed: %v", err)
	}
	seen := make(map[string]struct{}, count)
	for number := range numbers {
		if _, exists := seen[number]; exists {
			t.Errorf("duplicate concurrent journal number %q", number)
		}
		seen[number] = struct{}{}
	}
	if len(seen) != count {
		t.Fatalf("unique concurrent journal numbers = %d, want %d", len(seen), count)
	}

	for _, accountID := range []int{fixture.debitAccountID, fixture.creditAccountID} {
		var accountBalance, latestRunningBalance float64
		if err := db.QueryRow(`SELECT current_balance FROM chart_of_accounts WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, accountID, fixture.tenantID, fixture.restaurantID).Scan(&accountBalance); err != nil {
			t.Fatalf("read account %d balance: %v", accountID, err)
		}
		if err := db.QueryRow(`SELECT running_balance FROM general_ledger WHERE account_id = $1 AND tenant_id = $2 AND restaurant_id = $3 ORDER BY entry_date DESC, id DESC LIMIT 1`, accountID, fixture.tenantID, fixture.restaurantID).Scan(&latestRunningBalance); err != nil {
			t.Fatalf("read account %d latest running balance: %v", accountID, err)
		}
		if accountBalance != count*100 || latestRunningBalance != count*100 {
			t.Fatalf("account %d current/latest running balance = %.2f/%.2f, want %.2f", accountID, accountBalance, latestRunningBalance, float64(count*100))
		}
	}
}

func TestCFR003PostgresConcurrentJournalNumbersKeepTenantRestaurantScopes(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	otherTenant := openAtomicPostingTenantFixture(t, db)
	siblingRestaurant := openAtomicPostingRestaurantFixture(t, db, fixture.tenantID)
	fixtures := []*atomicPostingFixture{fixture, otherTenant, siblingRestaurant}
	start := make(chan struct{})
	results := make(chan string, len(fixtures))
	errs := make(chan error, len(fixtures))
	var workers sync.WaitGroup
	for _, scopedFixture := range fixtures {
		workers.Add(1)
		go func(scopedFixture *atomicPostingFixture) {
			defer workers.Done()
			<-start
			entry := scopedFixture.entry()
			id, err := scopedFixture.ledger.CreateAndPostAutoEntry(entry, 0)
			if err != nil {
				errs <- err
				return
			}
			var number string
			if err := db.QueryRow(`SELECT entry_number FROM journal_entries WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, scopedFixture.tenantID, scopedFixture.restaurantID).Scan(&number); err != nil {
				errs <- err
				return
			}
			results <- number
		}(scopedFixture)
	}
	close(start)
	workers.Wait()
	close(errs)
	close(results)
	for err := range errs {
		t.Errorf("scoped concurrent journal creation failed: %v", err)
	}
	for number := range results {
		if number != "JE-2026-00001" {
			t.Errorf("first number in independent tenant/restaurant scope = %q, want JE-2026-00001", number)
		}
	}
}

func cfr003LegacyNextJournalNumber(t *testing.T, tx *sql.Tx, fixture *atomicPostingFixture, entryDate time.Time) string {
	t.Helper()
	prefix := fmt.Sprintf("JE-%s-", entryDate.Format("2006"))
	var next int
	err := tx.QueryRow(`
		SELECT COALESCE(MAX(CAST(RIGHT(entry_number, 5) AS INT)), 0) + 1
		FROM journal_entries
		WHERE tenant_id = $1 AND restaurant_id = $2 AND entry_number LIKE $3
	`, fixture.tenantID, fixture.restaurantID, prefix+"%").Scan(&next)
	if err != nil {
		t.Fatalf("run baseline MAX()+1 allocator: %v", err)
	}
	return fmt.Sprintf("%s%05d", prefix, next)
}

func insertCfr003Draft(t *testing.T, tx *sql.Tx, fixture *atomicPostingFixture, number string, entryDate time.Time) {
	t.Helper()
	if err := insertCfr003DraftError(tx, fixture, number, entryDate); err != nil {
		t.Fatalf("insert baseline journal %s: %v", number, err)
	}
}

func insertCfr003DraftError(tx *sql.Tx, fixture *atomicPostingFixture, number string, entryDate time.Time) error {
	_, err := tx.Exec(`
		INSERT INTO journal_entries (tenant_id, restaurant_id, entry_number, entry_date, fiscal_period, description, entry_type, status)
		VALUES ($1, $2, $3, $4, $5, 'CFR-003 baseline reproducer', 'manual', 'draft')
	`, fixture.tenantID, fixture.restaurantID, number, entryDate, entryDate.Format("2006-01"))
	return err
}
