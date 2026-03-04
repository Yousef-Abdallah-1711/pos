package usecase

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

func TestPhase6FullAccountingCycle_OrderSalaryAndTrialBalanceBalanced(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	engine := newAutoPostingEngineForTest(db)
	journalRepo := repository.NewJournalEntryRepository(db)
	generalLedgerRepo := repository.NewGeneralLedgerRepository(db, journalRepo)
	reportsUC := NewFinancialReportUseCase(db, generalLedgerRepo)

	postedBy := 31
	tenantID := 10
	restaurantID := 20

	orderDate := time.Date(2026, time.February, 1, 12, 0, 0, 0, time.UTC)
	order := &domain.Order{
		ID:            701,
		TenantID:      int64(tenantID),
		RestaurantID:  int64(restaurantID),
		OrderNumber:   "ORD-701",
		Subtotal:      100,
		TaxAmount:     14,
		TotalAmount:   114,
		PaymentMethod: "cash",
		PaymentStatus: "paid",
		CreatedAt:     orderDate,
	}

	mockAutoPostingDefaults(mock, domain.AutoPostingEventOrderPaid, int(order.ID))
	mockAccountByCode(mock, tenantID, restaurantID, "1110", 1101, "1110", "Cash on Hand", domain.NormalBalanceDebit)
	mockAccountByCode(mock, tenantID, restaurantID, "4100", 4100, "4100", "Food Sales", domain.NormalBalanceCredit)
	mockAccountByCode(mock, tenantID, restaurantID, "2120", 2120, "2120", "VAT Payable", domain.NormalBalanceCredit)
	mockCreateJournalEntry(mock, orderDate, 9101, []postingExpectation{
		{accountID: 1101, accountCode: "1110", normalBalance: domain.NormalBalanceDebit, debit: 114, credit: 0},
		{accountID: 4100, accountCode: "4100", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{accountID: 2120, accountCode: "2120", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 14},
	})
	mockPostJournalEntry(mock, 9101, orderDate, []postingExpectation{
		{lineID: 111, accountID: 1101, normalBalance: domain.NormalBalanceDebit, debit: 114, credit: 0},
		{lineID: 112, accountID: 4100, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{lineID: 113, accountID: 2120, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 14},
	}, postedBy)

	if err := engine.PostSalesEntry(order, &postedBy); err != nil {
		t.Fatalf("PostSalesEntry failed: %v", err)
	}

	paymentDate := time.Date(2026, time.February, 2, 16, 0, 0, 0, time.UTC)
	salary := &domain.Salary{
		ID:              801,
		TenantID:        tenantID,
		RestaurantID:    restaurantID,
		GrossSalary:     1000,
		NetSalary:       800,
		Tax:             100,
		SocialInsurance: 100,
		PaymentMethod:   "bank_transfer",
		Status:          "paid",
		IsPaid:          true,
		PaymentDate:     &paymentDate,
	}

	mockAutoPostingDefaults(mock, domain.AutoPostingEventSalaryPaid, salary.ID)
	mockAccountByCode(mock, tenantID, restaurantID, "1111", 1111, "1111", "Main Bank Account", domain.NormalBalanceDebit)
	mockAccountByCode(mock, tenantID, restaurantID, "5210", 5210, "5210", "Salaries Expense", domain.NormalBalanceDebit)
	mockAccountByCode(mock, tenantID, restaurantID, "2110", 2110, "2110", "Accrued Expenses", domain.NormalBalanceCredit)
	mockAccountByCode(mock, tenantID, restaurantID, "2110", 2110, "2110", "Accrued Expenses", domain.NormalBalanceCredit)
	mockCreateJournalEntry(mock, paymentDate, 9102, []postingExpectation{
		{accountID: 5210, accountCode: "5210", normalBalance: domain.NormalBalanceDebit, debit: 1000, credit: 0},
		{accountID: 2110, accountCode: "2110", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{accountID: 2110, accountCode: "2110", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{accountID: 1111, accountCode: "1111", normalBalance: domain.NormalBalanceDebit, debit: 0, credit: 800},
	})
	mockPostJournalEntry(mock, 9102, paymentDate, []postingExpectation{
		{lineID: 211, accountID: 5210, normalBalance: domain.NormalBalanceDebit, debit: 1000, credit: 0},
		{lineID: 212, accountID: 2110, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{lineID: 213, accountID: 2110, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 100},
		{lineID: 214, accountID: 1111, normalBalance: domain.NormalBalanceDebit, debit: 0, credit: 800},
	}, postedBy)

	if err := engine.PostPayrollEntry(salary, &postedBy); err != nil {
		t.Fatalf("PostPayrollEntry failed: %v", err)
	}

	asOf := time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`WITH latest_balances AS`).
		WithArgs(tenantID, restaurantID, asOf).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "account_code", "account_name", "account_type", "normal_balance", "net_balance",
		}).
			AddRow(1101, "1110", "Cash on Hand", "asset", "debit", 114.0).
			AddRow(1111, "1111", "Main Bank Account", "asset", "debit", -800.0).
			AddRow(4100, "4100", "Food Sales", "revenue", "credit", 100.0).
			AddRow(2120, "2120", "VAT Payable", "liability", "credit", 14.0).
			AddRow(5210, "5210", "Salaries Expense", "expense", "debit", 1000.0).
			AddRow(2110, "2110", "Accrued Expenses", "liability", "credit", 200.0))

	trialBalance, err := reportsUC.GenerateTrialBalance(tenantID, restaurantID, asOf)
	if err != nil {
		t.Fatalf("GenerateTrialBalance failed: %v", err)
	}
	if !trialBalance.IsBalanced {
		t.Fatalf("expected balanced trial balance, got debits=%f credits=%f", trialBalance.TotalDebits, trialBalance.TotalCredits)
	}
	if roundedAmount(trialBalance.TotalDebits) != 1114 || roundedAmount(trialBalance.TotalCredits) != 1114 {
		t.Fatalf("expected trial balance totals 1114/1114, got %f/%f", trialBalance.TotalDebits, trialBalance.TotalCredits)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6AutoPostingAccuracy_HighVolumeOrdersAndSalaries(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	engine := newAutoPostingEngineForTest(db)
	postedBy := 19

	totalOrderDebits := 0.0
	totalOrderCredits := 0.0

	for i := 1; i <= 50; i++ {
		subtotal := float64(100 + i)
		tax := roundedAmount(subtotal * 0.14)
		total := roundedAmount(subtotal + tax)
		entryDate := time.Date(2026, time.February, (i%28)+1, 10, 0, 0, 0, time.UTC)
		orderID := int64(3000 + i)
		journalID := 10000 + i
		lineBase := i * 10

		order := &domain.Order{
			ID:            orderID,
			TenantID:      10,
			RestaurantID:  20,
			OrderNumber:   "ORD-VOLUME",
			Subtotal:      subtotal,
			TaxAmount:     tax,
			TotalAmount:   total,
			PaymentMethod: "cash",
			PaymentStatus: "paid",
			CreatedAt:     entryDate,
		}

		mockAutoPostingDefaults(mock, domain.AutoPostingEventOrderPaid, int(orderID))
		mockAccountByCode(mock, 10, 20, "1110", 1101, "1110", "Cash on Hand", domain.NormalBalanceDebit)
		mockAccountByCode(mock, 10, 20, "4100", 4100, "4100", "Food Sales", domain.NormalBalanceCredit)
		mockAccountByCode(mock, 10, 20, "2120", 2120, "2120", "VAT Payable", domain.NormalBalanceCredit)
		mockCreateJournalEntry(mock, entryDate, journalID, []postingExpectation{
			{accountID: 1101, accountCode: "1110", normalBalance: domain.NormalBalanceDebit, debit: total, credit: 0},
			{accountID: 4100, accountCode: "4100", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: subtotal},
			{accountID: 2120, accountCode: "2120", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: tax},
		})
		mockPostJournalEntry(mock, journalID, entryDate, []postingExpectation{
			{lineID: lineBase + 1, accountID: 1101, normalBalance: domain.NormalBalanceDebit, debit: total, credit: 0},
			{lineID: lineBase + 2, accountID: 4100, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: subtotal},
			{lineID: lineBase + 3, accountID: 2120, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: tax},
		}, postedBy)

		if err := engine.PostSalesEntry(order, &postedBy); err != nil {
			t.Fatalf("PostSalesEntry failed on order %d: %v", i, err)
		}

		totalOrderDebits += total
		totalOrderCredits += subtotal + tax
	}

	totalSalaryDebits := 0.0
	totalSalaryCredits := 0.0

	for i := 1; i <= 20; i++ {
		gross := float64(1200 + (i * 5))
		tax := roundedAmount(gross * 0.10)
		social := roundedAmount(gross * 0.05)
		net := roundedAmount(gross - tax - social)
		paymentDate := time.Date(2026, time.February, (i%28)+1, 17, 0, 0, 0, time.UTC)
		salaryID := 4000 + i
		journalID := 20000 + i
		lineBase := 5000 + (i * 10)

		salary := &domain.Salary{
			ID:              salaryID,
			TenantID:        10,
			RestaurantID:    20,
			GrossSalary:     gross,
			NetSalary:       net,
			Tax:             tax,
			SocialInsurance: social,
			PaymentMethod:   "bank_transfer",
			Status:          "paid",
			IsPaid:          true,
			PaymentDate:     &paymentDate,
		}

		mockAutoPostingDefaults(mock, domain.AutoPostingEventSalaryPaid, salaryID)
		mockAccountByCode(mock, 10, 20, "1111", 1111, "1111", "Main Bank Account", domain.NormalBalanceDebit)
		mockAccountByCode(mock, 10, 20, "5210", 5210, "5210", "Salaries Expense", domain.NormalBalanceDebit)
		mockAccountByCode(mock, 10, 20, "2110", 2110, "2110", "Accrued Expenses", domain.NormalBalanceCredit)
		mockAccountByCode(mock, 10, 20, "2110", 2110, "2110", "Accrued Expenses", domain.NormalBalanceCredit)
		mockCreateJournalEntry(mock, paymentDate, journalID, []postingExpectation{
			{accountID: 5210, accountCode: "5210", normalBalance: domain.NormalBalanceDebit, debit: gross, credit: 0},
			{accountID: 2110, accountCode: "2110", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: tax},
			{accountID: 2110, accountCode: "2110", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: social},
			{accountID: 1111, accountCode: "1111", normalBalance: domain.NormalBalanceDebit, debit: 0, credit: net},
		})
		mockPostJournalEntry(mock, journalID, paymentDate, []postingExpectation{
			{lineID: lineBase + 1, accountID: 5210, normalBalance: domain.NormalBalanceDebit, debit: gross, credit: 0},
			{lineID: lineBase + 2, accountID: 2110, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: tax},
			{lineID: lineBase + 3, accountID: 2110, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: social},
			{lineID: lineBase + 4, accountID: 1111, normalBalance: domain.NormalBalanceDebit, debit: 0, credit: net},
		}, postedBy)

		if err := engine.PostPayrollEntry(salary, &postedBy); err != nil {
			t.Fatalf("PostPayrollEntry failed on salary %d: %v", i, err)
		}

		totalSalaryDebits += gross
		totalSalaryCredits += tax + social + net
	}

	if roundedAmount(totalOrderDebits) != roundedAmount(totalOrderCredits) {
		t.Fatalf("orders are not balanced: debits=%f credits=%f", totalOrderDebits, totalOrderCredits)
	}
	if roundedAmount(totalSalaryDebits) != roundedAmount(totalSalaryCredits) {
		t.Fatalf("salaries are not balanced: debits=%f credits=%f", totalSalaryDebits, totalSalaryCredits)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6FullAccountingCycleE2E_OrderSalarySupplierFlowAndProfitLoss(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	engine := newAutoPostingEngineForTest(db)
	uc := NewFinancialReportUseCase(db, nil)
	postedBy := 55
	tenantID := 10
	restaurantID := 20

	// Step 1: paid order -> auto sales posting
	orderDate := time.Date(2026, time.February, 5, 11, 0, 0, 0, time.UTC)
	order := &domain.Order{
		ID:            9901,
		TenantID:      int64(tenantID),
		RestaurantID:  int64(restaurantID),
		OrderNumber:   "ORD-E2E-9901",
		Subtotal:      1000,
		TaxAmount:     140,
		TotalAmount:   1140,
		PaymentMethod: "cash",
		PaymentStatus: "paid",
		CreatedAt:     orderDate,
	}
	mockAutoPostingDefaults(mock, domain.AutoPostingEventOrderPaid, int(order.ID))
	mockAccountByCode(mock, tenantID, restaurantID, "1110", 1101, "1110", "Cash on Hand", domain.NormalBalanceDebit)
	mockAccountByCode(mock, tenantID, restaurantID, "4100", 4100, "4100", "Food Sales", domain.NormalBalanceCredit)
	mockAccountByCode(mock, tenantID, restaurantID, "2120", 2120, "2120", "VAT Payable", domain.NormalBalanceCredit)
	mockCreateJournalEntry(mock, orderDate, 9301, []postingExpectation{
		{accountID: 1101, accountCode: "1110", normalBalance: domain.NormalBalanceDebit, debit: 1140, credit: 0},
		{accountID: 4100, accountCode: "4100", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 1000},
		{accountID: 2120, accountCode: "2120", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 140},
	})
	mockPostJournalEntry(mock, 9301, orderDate, []postingExpectation{
		{lineID: 901, accountID: 1101, normalBalance: domain.NormalBalanceDebit, debit: 1140, credit: 0},
		{lineID: 902, accountID: 4100, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 1000},
		{lineID: 903, accountID: 2120, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 140},
	}, postedBy)
	if err := engine.PostSalesEntry(order, &postedBy); err != nil {
		t.Fatalf("PostSalesEntry failed: %v", err)
	}

	// Step 2: paid salary -> auto payroll posting
	paymentDate := time.Date(2026, time.February, 10, 15, 0, 0, 0, time.UTC)
	salary := &domain.Salary{
		ID:              9902,
		TenantID:        tenantID,
		RestaurantID:    restaurantID,
		GrossSalary:     600,
		NetSalary:       500,
		Tax:             60,
		SocialInsurance: 40,
		PaymentMethod:   "bank_transfer",
		Status:          "paid",
		IsPaid:          true,
		PaymentDate:     &paymentDate,
	}
	mockAutoPostingDefaults(mock, domain.AutoPostingEventSalaryPaid, salary.ID)
	mockAccountByCode(mock, tenantID, restaurantID, "1111", 1111, "1111", "Main Bank Account", domain.NormalBalanceDebit)
	mockAccountByCode(mock, tenantID, restaurantID, "5210", 5210, "5210", "Salaries Expense", domain.NormalBalanceDebit)
	mockAccountByCode(mock, tenantID, restaurantID, "2110", 2110, "2110", "Accrued Expenses", domain.NormalBalanceCredit)
	mockAccountByCode(mock, tenantID, restaurantID, "2110", 2110, "2110", "Accrued Expenses", domain.NormalBalanceCredit)
	mockCreateJournalEntry(mock, paymentDate, 9302, []postingExpectation{
		{accountID: 5210, accountCode: "5210", normalBalance: domain.NormalBalanceDebit, debit: 600, credit: 0},
		{accountID: 2110, accountCode: "2110", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 60},
		{accountID: 2110, accountCode: "2110", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 40},
		{accountID: 1111, accountCode: "1111", normalBalance: domain.NormalBalanceDebit, debit: 0, credit: 500},
	})
	mockPostJournalEntry(mock, 9302, paymentDate, []postingExpectation{
		{lineID: 911, accountID: 5210, normalBalance: domain.NormalBalanceDebit, debit: 600, credit: 0},
		{lineID: 912, accountID: 2110, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 60},
		{lineID: 913, accountID: 2110, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 40},
		{lineID: 914, accountID: 1111, normalBalance: domain.NormalBalanceDebit, debit: 0, credit: 500},
	}, postedBy)
	if err := engine.PostPayrollEntry(salary, &postedBy); err != nil {
		t.Fatalf("PostPayrollEntry failed: %v", err)
	}

	// Step 3: supplier invoice approved -> payable journal
	invoiceDate := time.Date(2026, time.February, 12, 12, 0, 0, 0, time.UTC)
	mockCreateJournalEntry(mock, invoiceDate, 9303, []postingExpectation{
		{accountID: 5110, accountCode: "5110", normalBalance: domain.NormalBalanceDebit, debit: 250, credit: 0},
		{accountID: 1150, accountCode: "1150", normalBalance: domain.NormalBalanceDebit, debit: 35, credit: 0},
		{accountID: 2101, accountCode: "2101", normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 285},
	})
	mockPostJournalEntry(mock, 9303, invoiceDate, []postingExpectation{
		{lineID: 921, accountID: 5110, normalBalance: domain.NormalBalanceDebit, debit: 250, credit: 0},
		{lineID: 922, accountID: 1150, normalBalance: domain.NormalBalanceDebit, debit: 35, credit: 0},
		{lineID: 923, accountID: 2101, normalBalance: domain.NormalBalanceCredit, debit: 0, credit: 285},
	}, postedBy)
	_, err = engine.createAndPostEntry(
		tenantID,
		restaurantID,
		domain.PostingEvent{
			Code:            domain.AutoPostingEventSupplierInvoice,
			EntryType:       domain.AutoPostingEntryTypeSupplierInvoice,
			ReferenceType:   "supplier_invoice",
			ReferenceID:     9903,
			ReferenceNumber: "PINV-9903",
			Description:     "Auto-posted supplier invoice PINV-9903",
			PostedBy:        &postedBy,
		},
		invoiceDate,
		[]domain.JournalEntryLine{
			{AccountID: 5110, DebitAmount: 250, Description: "Expense", LineOrder: 1},
			{AccountID: 1150, DebitAmount: 35, Description: "Input VAT", LineOrder: 2},
			{AccountID: 2101, CreditAmount: 285, Description: "Accounts payable", LineOrder: 3},
		},
	)
	if err != nil {
		t.Fatalf("supplier invoice posting failed: %v", err)
	}

	// Step 4: supplier payment -> clear AP against bank
	supplierPaymentDate := time.Date(2026, time.February, 18, 14, 0, 0, 0, time.UTC)
	mockCreateJournalEntry(mock, supplierPaymentDate, 9304, []postingExpectation{
		{accountID: 2101, accountCode: "2101", normalBalance: domain.NormalBalanceCredit, debit: 285, credit: 0},
		{accountID: 1111, accountCode: "1111", normalBalance: domain.NormalBalanceDebit, debit: 0, credit: 285},
	})
	mockPostJournalEntry(mock, 9304, supplierPaymentDate, []postingExpectation{
		{lineID: 931, accountID: 2101, normalBalance: domain.NormalBalanceCredit, debit: 285, credit: 0},
		{lineID: 932, accountID: 1111, normalBalance: domain.NormalBalanceDebit, debit: 0, credit: 285},
	}, postedBy)
	_, err = engine.createAndPostEntry(
		tenantID,
		restaurantID,
		domain.PostingEvent{
			Code:            domain.AutoPostingEventSupplierPayment,
			EntryType:       domain.AutoPostingEntryTypeSupplierPayment,
			ReferenceType:   "supplier_payment",
			ReferenceID:     9904,
			ReferenceNumber: "PMT-M-9904",
			Description:     "Auto-posted supplier payment PMT-M-9904",
			PostedBy:        &postedBy,
		},
		supplierPaymentDate,
		[]domain.JournalEntryLine{
			{AccountID: 2101, DebitAmount: 285, Description: "Settle AP", LineOrder: 1},
			{AccountID: 1111, CreditAmount: 285, Description: "Bank payment", LineOrder: 2},
		},
	)
	if err != nil {
		t.Fatalf("supplier payment posting failed: %v", err)
	}

	// Step 5: generate P&L and validate profit
	from := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC)
	compareFrom := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	compareTo := time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC)

	columns := []string{"account_code", "account_name", "account_type", "debit_total", "credit_total"}
	mock.ExpectQuery(`FROM journal_entry_lines jel`).
		WithArgs(tenantID, restaurantID, from, to).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("4100", "Food Sales", "revenue", 0.0, 1000.0).
			AddRow("5210", "Salaries Expense", "expense", 600.0, 0.0).
			AddRow("5110", "Operating Expense", "expense", 250.0, 0.0))
	mock.ExpectQuery(`FROM journal_entry_lines jel`).
		WithArgs(tenantID, restaurantID, compareFrom, compareTo).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("4100", "Food Sales", "revenue", 0.0, 800.0).
			AddRow("5210", "Salaries Expense", "expense", 500.0, 0.0))

	report, err := uc.GenerateProfitLossWithComparison(tenantID, restaurantID, from, to, compareFrom, compareTo)
	if err != nil {
		t.Fatalf("GenerateProfitLossWithComparison failed: %v", err)
	}
	if roundedAmount(report.TotalRevenue) != 1000 {
		t.Fatalf("expected revenue 1000, got %f", report.TotalRevenue)
	}
	if roundedAmount(report.TotalExpenses) != 850 {
		t.Fatalf("expected expenses 850, got %f", report.TotalExpenses)
	}
	if roundedAmount(report.NetProfit) != 150 {
		t.Fatalf("expected net profit 150, got %f", report.NetProfit)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
