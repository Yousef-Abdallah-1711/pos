//go:build integration

package repository

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
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
		var accountBalance float64
		if err := db.QueryRow(`SELECT current_balance FROM chart_of_accounts WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, accountID, fixture.tenantID, fixture.restaurantID).Scan(&accountBalance); err != nil {
			t.Fatalf("read account %d balance: %v", accountID, err)
		}
		ledgerBalance, err := fixture.ledger.GetAccountBalance(fixture.tenantID, fixture.restaurantID, accountID, nil)
		if err != nil {
			t.Fatalf("read account %d movement balance: %v", accountID, err)
		}
		assertFinancialBalanceCents(t, fmt.Sprintf("account %d stored current balance", accountID), int64(count*10000), accountBalance)
		assertFinancialBalanceCents(t, fmt.Sprintf("account %d independent movement balance", accountID), int64(count*10000), ledgerBalance)
	}
}

func TestCFR003PostgresBackdatedEntryKeepsAsOfBalancesChronological(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	lateDate := time.Date(2026, time.November, 20, 0, 0, 0, 0, time.UTC)
	earlyDate := time.Date(2026, time.November, 10, 0, 0, 0, 0, time.UTC)
	for _, posting := range []struct {
		referenceID int
		entryDate   time.Time
		amount      float64
	}{{43020, lateDate, 100}, {43010, earlyDate, 50}, {43011, earlyDate, 25}} {
		entry := fixture.entry()
		entry.ReferenceID = &posting.referenceID
		entry.EntryDate = posting.entryDate
		entry.Lines[0].DebitAmount = posting.amount
		entry.Lines[1].CreditAmount = posting.amount
		if _, err := fixture.ledger.CreateAndPostAutoEntry(entry, 0); err != nil {
			t.Fatalf("post journal %d: %v", posting.referenceID, err)
		}
	}

	debitBalance, err := fixture.ledger.GetAccountBalance(fixture.tenantID, fixture.restaurantID, fixture.debitAccountID, &earlyDate)
	if err != nil {
		t.Fatalf("get debit balance as of early date: %v", err)
	}
	assertFinancialBalanceCents(t, "debit balance as of backdated date", 7500, debitBalance)

	trialBalance, err := fixture.ledger.GetTrialBalance(fixture.tenantID, fixture.restaurantID, earlyDate)
	if err != nil {
		t.Fatalf("get trial balance as of early date: %v", err)
	}
	var debitFound, creditFound bool
	for _, row := range trialBalance {
		switch row.AccountID {
		case fixture.debitAccountID:
			debitFound = row.DebitBalance == 75 && row.CreditBalance == 0
		case fixture.creditAccountID:
			creditFound = row.DebitBalance == 0 && row.CreditBalance == 75
		}
	}
	if !debitFound || !creditFound {
		t.Fatalf("as-of trial balance does not reflect only the early posting: %+v", trialBalance)
	}
	var independentAsOfBalance float64
	if err := db.QueryRow(`SELECT COALESCE(SUM(debit_amount - credit_amount), 0) FROM general_ledger WHERE tenant_id = $1 AND restaurant_id = $2 AND account_id = $3 AND entry_date <= $4`, fixture.tenantID, fixture.restaurantID, fixture.debitAccountID, earlyDate).Scan(&independentAsOfBalance); err != nil {
		t.Fatalf("independently reconcile backdated as-of balance: %v", err)
	}
	assertFinancialBalanceCents(t, "independent SQL debit balance as of backdated date", 7500, independentAsOfBalance)

	ledgerLines, err := fixture.ledger.GetLedgerByAccount(fixture.tenantID, fixture.restaurantID, fixture.debitAccountID, nil, nil)
	if err != nil {
		t.Fatalf("list debit ledger: %v", err)
	}
	if len(ledgerLines) != 3 {
		t.Fatalf("ledger row count = %d, want 3 rows including equal-date movements", len(ledgerLines))
	}
	assertFinancialBalanceCents(t, "latest ledger running balance", 17500, ledgerLines[0].RunningBalance)
	assertFinancialBalanceCents(t, "second ledger running balance", 7500, ledgerLines[1].RunningBalance)
	assertFinancialBalanceCents(t, "oldest ledger running balance", 5000, ledgerLines[2].RunningBalance)
	if !ledgerLines[1].EntryDate.Equal(earlyDate) || !ledgerLines[2].EntryDate.Equal(earlyDate) || ledgerLines[1].JournalEntryID <= ledgerLines[2].JournalEntryID {
		t.Fatalf("equal-date ledger rows are not deterministically ordered by descending id: %+v", ledgerLines)
	}

	firstPage, nextCursor, err := fixture.ledger.GetLedgerByAccountCursor(fixture.tenantID, fixture.restaurantID, fixture.debitAccountID, nil, nil, 1, nil)
	if err != nil {
		t.Fatalf("list first chronological ledger page: %v", err)
	}
	if len(firstPage) != 1 || nextCursor == nil || math.Round(firstPage[0].RunningBalance*100) != 17500 {
		t.Fatalf("first ledger page/cursor = %+v/%v, want one newest row and a continuation cursor", firstPage, nextCursor)
	}

	secondPage, finalCursor, err := fixture.ledger.GetLedgerByAccountCursor(fixture.tenantID, fixture.restaurantID, fixture.debitAccountID, nil, nil, 1, nextCursor)
	if err != nil {
		t.Fatalf("continue chronological ledger page: %v", err)
	}
	if len(secondPage) != 1 || finalCursor == nil || !secondPage[0].EntryDate.Equal(earlyDate) || math.Round(secondPage[0].RunningBalance*100) != 7500 {
		t.Fatalf("second ledger page/cursor = %+v/%v, want newest equal-date row and another cursor", secondPage, finalCursor)
	}

	thirdPage, finalCursor, err := fixture.ledger.GetLedgerByAccountCursor(fixture.tenantID, fixture.restaurantID, fixture.debitAccountID, nil, nil, 1, finalCursor)
	if err != nil {
		t.Fatalf("continue ledger pagination across equal-date rows: %v", err)
	}
	if len(thirdPage) != 1 || finalCursor != nil || !thirdPage[0].EntryDate.Equal(earlyDate) || math.Round(thirdPage[0].RunningBalance*100) != 5000 {
		t.Fatalf("third ledger page/cursor = %+v/%v, want oldest equal-date row and no further cursor", thirdPage, finalCursor)
	}

	var persistedLateSnapshot float64
	if err := db.QueryRow(`SELECT running_balance FROM general_ledger WHERE tenant_id = $1 AND restaurant_id = $2 AND account_id = $3 AND entry_date = $4 ORDER BY id LIMIT 1`, fixture.tenantID, fixture.restaurantID, fixture.debitAccountID, lateDate).Scan(&persistedLateSnapshot); err != nil {
		t.Fatalf("read immutable posting-time snapshot: %v", err)
	}
	assertFinancialBalanceCents(t, "unchanged original posting-time snapshot", 10000, persistedLateSnapshot)

	var storedCurrentBalance float64
	if err := db.QueryRow(`SELECT current_balance FROM chart_of_accounts WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, fixture.debitAccountID, fixture.tenantID, fixture.restaurantID).Scan(&storedCurrentBalance); err != nil {
		t.Fatalf("read account current balance: %v", err)
	}
	assertFinancialBalanceCents(t, "stored current balance after backdated and equal-date postings", 17500, storedCurrentBalance)
	var independentCurrentBalance float64
	if err := db.QueryRow(`SELECT COALESCE(SUM(debit_amount - credit_amount), 0) FROM general_ledger WHERE tenant_id = $1 AND restaurant_id = $2 AND account_id = $3`, fixture.tenantID, fixture.restaurantID, fixture.debitAccountID).Scan(&independentCurrentBalance); err != nil {
		t.Fatalf("independently reconcile current account balance: %v", err)
	}
	assertFinancialBalanceCents(t, "independent SQL current balance", 17500, independentCurrentBalance)
}

func TestCFR003PostgresOperationalRepairReconcilesExistingBalanceCaches(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	entry := fixture.entry()
	entry.Lines[0].DebitAmount = 40.25
	entry.Lines[1].CreditAmount = 40.25
	if _, err := fixture.ledger.CreateAndPostAutoEntry(entry, 0); err != nil {
		t.Fatalf("post source journal: %v", err)
	}

	var zeroMovementAccountID int
	if err := db.QueryRow(`
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance, current_balance)
		VALUES ($1, $2, 'CFR004-Z', 'CFR-004 zero-movement fixture', 'asset', 'debit', 88.88)
		RETURNING id
	`, fixture.tenantID, fixture.restaurantID).Scan(&zeroMovementAccountID); err != nil {
		t.Fatalf("create zero-movement account with stale cache: %v", err)
	}

	if _, err := db.Exec(`
		UPDATE chart_of_accounts SET current_balance = CASE id
			WHEN $1 THEN 999.99 WHEN $2 THEN -999.99 ELSE current_balance END
		WHERE tenant_id = $3 AND restaurant_id = $4 AND id IN ($1, $2)
	`, fixture.debitAccountID, fixture.creditAccountID, fixture.tenantID, fixture.restaurantID); err != nil {
		t.Fatalf("seed stale derived balance caches: %v", err)
	}
	var journalCountBefore, ledgerCountBefore int
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entries WHERE tenant_id = $1 AND restaurant_id = $2`, fixture.tenantID, fixture.restaurantID).Scan(&journalCountBefore); err != nil {
		t.Fatalf("count posted journal history before repair: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM general_ledger WHERE tenant_id = $1 AND restaurant_id = $2`, fixture.tenantID, fixture.restaurantID).Scan(&ledgerCountBefore); err != nil {
		t.Fatalf("count GL history before repair: %v", err)
	}

	repairSQL, err := os.ReadFile(filepath.Join("..", "..", "scripts", "reconcile_chart_of_accounts_current_balance.sql"))
	if err != nil {
		t.Fatalf("read operational balance repair SQL: %v", err)
	}
	if _, err := db.Exec(string(repairSQL)); err != nil {
		t.Fatalf("run operational balance repair transaction: %v", err)
	}

	for _, check := range []struct {
		name      string
		accountID int
		wantCents int64
	}{
		{name: "debit-normal account", accountID: fixture.debitAccountID, wantCents: 4025},
		{name: "credit-normal account", accountID: fixture.creditAccountID, wantCents: 4025},
		{name: "zero-movement account", accountID: zeroMovementAccountID, wantCents: 0},
	} {
		var observed float64
		if err := db.QueryRow(`
			SELECT current_balance FROM chart_of_accounts
			WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
		`, check.accountID, fixture.tenantID, fixture.restaurantID).Scan(&observed); err != nil {
			t.Fatalf("read repaired %s cache: %v", check.name, err)
		}
		assertFinancialBalanceCents(t, "repaired "+check.name, check.wantCents, observed)
	}

	var journalCountAfter, ledgerCountAfter int
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entries WHERE tenant_id = $1 AND restaurant_id = $2`, fixture.tenantID, fixture.restaurantID).Scan(&journalCountAfter); err != nil {
		t.Fatalf("count posted journal history after repair: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM general_ledger WHERE tenant_id = $1 AND restaurant_id = $2`, fixture.tenantID, fixture.restaurantID).Scan(&ledgerCountAfter); err != nil {
		t.Fatalf("count GL history after repair: %v", err)
	}
	if journalCountAfter != journalCountBefore || ledgerCountAfter != ledgerCountBefore {
		t.Fatalf("repair changed posted history: journals %d -> %d, GL rows %d -> %d", journalCountBefore, journalCountAfter, ledgerCountBefore, ledgerCountAfter)
	}
}

func TestCFR003PostgresNormalBalanceCanChangeBeforeLedgerHistory(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	accountRepo := NewAccountRepository(db)
	account, err := accountRepo.GetAccountByID(fixture.tenantID, fixture.restaurantID, fixture.debitAccountID)
	if err != nil || account == nil {
		t.Fatalf("read unposted account: account=%+v err=%v", account, err)
	}
	account.NormalBalance = "credit"
	if err := accountRepo.UpdateAccount(account); err != nil {
		t.Fatalf("change normal balance before ledger history: %v", err)
	}
	account, err = accountRepo.GetAccountByID(fixture.tenantID, fixture.restaurantID, fixture.debitAccountID)
	if err != nil || account == nil || account.NormalBalance != "credit" {
		t.Fatalf("unposted account normal balance = %+v err=%v, want credit", account, err)
	}
	advancedRepo := NewAccountingAdvancedRepository(db, fixture.ledger.journalRepo)
	if err := advancedRepo.UpsertAccountFromImport(
		fixture.tenantID, fixture.restaurantID, "CFR004-D", "Imported debit fixture", "", "asset", "credit", "EGP", "", "", true,
	); err != nil {
		t.Fatalf("import normal balance before ledger history: %v", err)
	}
}

func TestCFR003PostgresNormalBalanceChangeRejectedAfterPosting(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	entry := fixture.entry()
	if _, err := fixture.ledger.CreateAndPostAutoEntry(entry, 0); err != nil {
		t.Fatalf("post historical journal: %v", err)
	}
	accountRepo := NewAccountRepository(db)
	account, err := accountRepo.GetAccountByID(fixture.tenantID, fixture.restaurantID, fixture.debitAccountID)
	if err != nil || account == nil {
		t.Fatalf("read posted account: account=%+v err=%v", account, err)
	}
	account.NormalBalance = "credit"
	if err := accountRepo.UpdateAccount(account); err == nil {
		t.Fatal("normal balance change after posting unexpectedly succeeded")
	}
	if _, err := db.Exec(`
		UPDATE chart_of_accounts
		SET normal_balance = 'credit'
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, fixture.debitAccountID, fixture.tenantID, fixture.restaurantID); err == nil {
		t.Fatal("direct SQL normal balance change after posting unexpectedly succeeded")
	}
	advancedRepo := NewAccountingAdvancedRepository(db, fixture.ledger.journalRepo)
	if err := advancedRepo.UpsertAccountFromImport(
		fixture.tenantID, fixture.restaurantID, "CFR004-D", "Imported debit fixture", "", "asset", "credit", "EGP", "", "", true,
	); err == nil {
		t.Fatal("account import changed normal balance after posting unexpectedly succeeded")
	}
	account, err = accountRepo.GetAccountByID(fixture.tenantID, fixture.restaurantID, fixture.debitAccountID)
	if err != nil || account == nil || account.NormalBalance != "debit" {
		t.Fatalf("rejected mutation changed historical normal balance: account=%+v err=%v", account, err)
	}
}

func assertFinancialBalanceCents(t *testing.T, label string, expected int64, observedAmount float64) {
	t.Helper()
	observed := int64(math.Round(observedAmount * 100))
	difference := observed - expected
	t.Logf("financial_reconciliation scenario=%s expected_cents=%d observed_cents=%d difference_cents=%d", label, expected, observed, difference)
	if difference != 0 {
		t.Fatalf("%s reconciliation: expected=%d cents observed=%d cents difference=%d cents", label, expected, observed, difference)
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
