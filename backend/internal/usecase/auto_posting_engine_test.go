package usecase

import (
	"database/sql"
	"sort"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

func TestAutoPostingEnginePostSalesEntryCreatesJournalAndLedger(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	engine := newAutoPostingEngineForTest(db)
	entryDate := time.Date(2026, 2, 28, 12, 0, 0, 0, time.UTC)
	order := &domain.Order{
		ID:            501,
		TenantID:      10,
		RestaurantID:  20,
		OrderNumber:   "ORD-501",
		Subtotal:      100,
		TaxAmount:     10,
		TotalAmount:   110,
		PaymentMethod: "cash",
		PaymentStatus: "paid",
		CreatedAt:     entryDate,
	}
	postedBy := 7

	mockAutoPostingDefaults(mock, domain.AutoPostingEventOrderPaid)
	mockAccountByCode(mock, 10, 20, "1110", 1101, "1110", "Cash on Hand", domain.NormalBalanceDebit)
	mockAccountByCode(mock, 10, 20, "4100", 4100, "4100", "Food Sales", domain.NormalBalanceCredit)
	mockAccountByCode(mock, 10, 20, "2120", 2120, "2120", "VAT Payable", domain.NormalBalanceCredit)

	mockCreateJournalEntry(mock, entryDate, 7001, []postingExpectation{
		{accountID: 1101, accountCode: "1110", normalBalance: domain.NormalBalanceDebit, debit: 110, credit: 0},
		{accountID: 4100, accountCode: "4100", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{accountID: 2120, accountCode: "2120", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 10},
	})
	mockPostJournalEntry(mock, 7001, entryDate, []postingExpectation{
		{lineID: 11, accountID: 1101, normalBalance: domain.NormalBalanceDebit, debit: 110, credit: 0},
		{lineID: 12, accountID: 4100, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{lineID: 13, accountID: 2120, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 10},
	}, postedBy)

	if err := engine.PostSalesEntry(order, &postedBy); err != nil {
		t.Fatalf("PostSalesEntry returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestAutoPostingEnginePostPayrollEntryCreatesJournalAndLedger(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	engine := newAutoPostingEngineForTest(db)
	entryDate := time.Date(2026, 2, 28, 15, 30, 0, 0, time.UTC)
	salary := &domain.Salary{
		ID:              88,
		TenantID:        10,
		RestaurantID:    20,
		GrossSalary:     1000,
		NetSalary:       800,
		Tax:             100,
		SocialInsurance: 100,
		PaymentMethod:   "bank_transfer",
		Status:          "paid",
		IsPaid:          true,
		PaymentDate:     &entryDate,
	}
	postedBy := 9

	mockAutoPostingDefaults(mock, domain.AutoPostingEventSalaryPaid)
	mockAccountByCode(mock, 10, 20, "1111", 1111, "1111", "Main Bank Account", domain.NormalBalanceDebit)
	mockAccountByCode(mock, 10, 20, "5210", 5210, "5210", "Salaries Expense", domain.NormalBalanceDebit)
	mockAccountByCode(mock, 10, 20, "2110", 2110, "2110", "Accrued Expenses", domain.NormalBalanceCredit)
	mockAccountByCode(mock, 10, 20, "2110", 2110, "2110", "Accrued Expenses", domain.NormalBalanceCredit)

	mockCreateJournalEntry(mock, entryDate, 7002, []postingExpectation{
		{accountID: 5210, accountCode: "5210", normalBalance: domain.NormalBalanceDebit, debit: 1000, credit: 0},
		{accountID: 2110, accountCode: "2110", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{accountID: 2110, accountCode: "2110", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{accountID: 1111, accountCode: "1111", normalBalance: domain.NormalBalanceDebit, debit: 0, credit: 800},
	})
	mockPostJournalEntry(mock, 7002, entryDate, []postingExpectation{
		{lineID: 21, accountID: 5210, normalBalance: domain.NormalBalanceDebit, debit: 1000, credit: 0},
		{lineID: 22, accountID: 2110, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{lineID: 23, accountID: 2110, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{lineID: 24, accountID: 1111, normalBalance: domain.NormalBalanceDebit, debit: 0, credit: 800},
	}, postedBy)

	if err := engine.PostPayrollEntry(salary, &postedBy); err != nil {
		t.Fatalf("PostPayrollEntry returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestAutoPostingEnginePostSupplierEntriesRequireTenantContext(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	engine := newAutoPostingEngineForTest(db)
	event := &domain.PostingEvent{
		ReferenceType: domain.AutoPostingReferenceTypeOrder,
		ReferenceID:   1,
		Description:   "supplier posting",
	}
	lines := []domain.JournalEntryLine{
		{AccountID: 1, DebitAmount: 100, LineOrder: 1},
		{AccountID: 2, CreditAmount: 100, LineOrder: 2},
	}

	if err := engine.PostSupplierInvoiceEntry(event, lines); err == nil {
		t.Fatalf("expected tenant/restaurant validation error for supplier invoice posting")
	}
	if err := engine.PostSupplierPaymentEntry(event, lines); err == nil {
		t.Fatalf("expected tenant/restaurant validation error for supplier payment posting")
	}
}

type postingExpectation struct {
	lineID        int
	accountID     int
	accountCode   string
	normalBalance string
	debit         float64
	credit        float64
}

func newAutoPostingEngineForTest(db *sql.DB) *AutoPostingEngine {
	journalRepo := repository.NewJournalEntryRepository(db)
	return NewAutoPostingEngine(
		repository.NewAutoPostingRepository(db),
		repository.NewAccountRepository(db),
		repository.NewInventoryMovementRepository(db),
		repository.NewGeneralLedgerRepository(db, journalRepo),
		repository.NewProductRepository(db),
	)
}

func mockAutoPostingDefaults(mock sqlmock.Sqlmock, eventCode string) {
	mock.ExpectQuery(`FROM auto_posting_rules`).
		WithArgs(10, 20, eventCode).
		WillReturnError(sql.ErrNoRows)
}

func mockAccountByCode(mock sqlmock.Sqlmock, tenantID, restaurantID int, code string, id int, accountCode, accountName, normalBalance string) {
	mock.ExpectQuery(`FROM chart_of_accounts`).
		WithArgs(tenantID, restaurantID, code).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tenant_id", "restaurant_id", "account_code", "account_name",
			"account_name_ar", "description", "description_ar", "account_type",
			"account_subtype", "normal_balance", "parent_id", "level", "is_header",
			"current_balance", "currency", "is_active", "is_system", "is_reconcilable",
			"tax_rate", "created_at", "updated_at", "created_by", "updated_by",
		}).AddRow(
			id, tenantID, restaurantID, accountCode, accountName,
			"", "", "", "asset",
			"", normalBalance, nil, 3, false,
			0, "EGP", true, true, false,
			nil, time.Now(), time.Now(), nil, nil,
		))
}

func mockCreateJournalEntry(mock sqlmock.Sqlmock, entryDate time.Time, journalEntryID int, lines []postingExpectation) {
	mock.ExpectBegin()
	mock.ExpectQuery(`pg_advisory_xact_lock`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"?column?"}).AddRow(1))
	mock.ExpectQuery(`SELECT id, status\s+FROM journal_entries`).
		WithArgs(10, 20, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}))
	mock.ExpectQuery(`pg_advisory_xact_lock`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"?column?"}).AddRow(1))
	mock.ExpectQuery(`SELECT COALESCE\(MAX\(CAST\(RIGHT\(entry_number, 5\) AS INT\)\), 0\) \+ 1`).
		WithArgs(10, 20, "JE-2026-%").
		WillReturnRows(sqlmock.NewRows([]string{"next_number"}).AddRow(1))
	mock.ExpectQuery(`INSERT INTO journal_entries`).
		WithArgs(
			10, 20, "JE-2026-00001", entryDate, "2026-02",
			sqlmock.AnyArg(), "", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), domain.JournalStatusDraft, "", sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(journalEntryID))

	for _, line := range lines {
		mock.ExpectQuery(`SELECT account_code, is_header`).
			WithArgs(line.accountID, 10, 20).
			WillReturnRows(sqlmock.NewRows([]string{"account_code", "is_header"}).AddRow(line.accountCode, false))
		mock.ExpectExec(`INSERT INTO journal_entry_lines`).
			WithArgs(
				journalEntryID,
				10,
				line.accountID,
				line.accountCode,
				line.debit,
				line.credit,
				sqlmock.AnyArg(),
				"",
				sqlmock.AnyArg(),
				"",
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func mockPostJournalEntry(mock sqlmock.Sqlmock, journalEntryID int, entryDate time.Time, lines []postingExpectation, postedBy int) {
	mock.ExpectQuery(`SELECT status, entry_date`).
		WithArgs(journalEntryID, 10, 20).
		WillReturnRows(sqlmock.NewRows([]string{"status", "entry_date", "fiscal_period"}).AddRow(domain.JournalStatusDraft, entryDate, "2026-02"))
	mock.ExpectQuery(`SELECT status\s+FROM fiscal_periods`).
		WithArgs(10, 20, "2026-02", entryDate).
		WillReturnError(sql.ErrNoRows)

	rows := sqlmock.NewRows([]string{"id", "account_id", "debit_amount", "credit_amount", "normal_balance"})
	for _, line := range lines {
		rows.AddRow(line.lineID, line.accountID, line.debit, line.credit, line.normalBalance)
	}
	mock.ExpectQuery(`FROM journal_entry_lines jel`).
		WithArgs(journalEntryID).
		WillReturnRows(rows)

	mock.ExpectExec(`UPDATE journal_entries`).
		WithArgs(domain.JournalStatusPosted, postedBy, journalEntryID, 10, 20).
		WillReturnResult(sqlmock.NewResult(1, 1))
	accountIDs := make([]int, 0, len(lines))
	seenAccounts := make(map[int]struct{}, len(lines))
	for _, line := range lines {
		if _, exists := seenAccounts[line.accountID]; exists {
			continue
		}
		seenAccounts[line.accountID] = struct{}{}
		accountIDs = append(accountIDs, line.accountID)
	}
	sort.Ints(accountIDs)
	for _, accountID := range accountIDs {
		mock.ExpectQuery(`SELECT id\s+FROM chart_of_accounts`).
			WithArgs(accountID, 10, 20).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(accountID))
	}

	for _, line := range lines {
		mock.ExpectQuery(`FROM general_ledger`).
			WithArgs(10, 20, line.accountID).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectExec(`INSERT INTO general_ledger`).
			WithArgs(
				10, 20, line.accountID, journalEntryID, line.lineID,
				entryDate, "2026-02", line.debit, line.credit, sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`UPDATE chart_of_accounts`).
			WithArgs(sqlmock.AnyArg(), line.accountID, 10, 20).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()
}
