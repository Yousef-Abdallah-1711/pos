package repository

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestPhase6TenantIsolation_AccountRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewAccountRepository(db)
	mock.ExpectQuery(`FROM chart_of_accounts`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	accounts, err := repo.ListAccounts(7, 11, false)
	if err != nil {
		t.Fatalf("ListAccounts returned error: %v", err)
	}
	if len(accounts) != 0 {
		t.Fatalf("expected no accounts, got %d", len(accounts))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_TrainingRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewTrainingRepository(db)
	mock.ExpectQuery(`FROM training_records`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.List(7, 11)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_ContractRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewContractRepository(db)
	mock.ExpectQuery(`FROM employment_contracts`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.List(7, 11)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_ResignationRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewResignationRepository(db)
	mock.ExpectQuery(`FROM resignations`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.List(7, 11)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_SupplierRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewSupplierRepository(db)
	mock.ExpectQuery(`FROM suppliers`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.ListSuppliers(7, 11)
	if err != nil {
		t.Fatalf("ListSuppliers returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_CustomerAccountRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewCustomerAccountRepository(db)
	mock.ExpectQuery(`FROM customers_accounting`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.ListCustomers(7, 11)
	if err != nil {
		t.Fatalf("ListCustomers returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_BankAccountRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewBankAccountRepository(db)
	mock.ExpectQuery(`FROM bank_accounts_ledger`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.ListBankAccounts(7, 11)
	if err != nil {
		t.Fatalf("ListBankAccounts returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_CustomerInvoiceRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewCustomerInvoiceRepository(db, nil, nil)
	mock.ExpectQuery(`FROM invoices_customer`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.ListInvoices(7, 11)
	if err != nil {
		t.Fatalf("ListInvoices returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_SupplierInvoiceRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewSupplierInvoiceRepository(db, nil, nil)
	mock.ExpectQuery(`FROM invoices_supplier`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.ListInvoices(7, 11)
	if err != nil {
		t.Fatalf("ListInvoices returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_PaymentAccountingRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewPaymentAccountingRepository(db, nil, nil, nil)

	mock.ExpectQuery(`FROM payments_received`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	received, err := repo.ListPaymentsReceived(7, 11)
	if err != nil {
		t.Fatalf("ListPaymentsReceived returned error: %v", err)
	}
	if len(received) != 0 {
		t.Fatalf("expected no records, got %d", len(received))
	}

	mock.ExpectQuery(`FROM payments_made`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	made, err := repo.ListPaymentsMade(7, 11)
	if err != nil {
		t.Fatalf("ListPaymentsMade returned error: %v", err)
	}
	if len(made) != 0 {
		t.Fatalf("expected no records, got %d", len(made))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_PerformanceReviewRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewPerformanceReviewRepository(db)
	mock.ExpectQuery(`FROM performance_reviews`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.List(7, 11)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_InsuranceRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewInsuranceRepository(db)
	mock.ExpectQuery(`FROM insurance_tracking`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.List(7, 11)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6TenantIsolation_TaxRepositoryListUsesTenantAndRestaurantScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewTaxRepository(db)
	mock.ExpectQuery(`FROM tax_periods`).
		WithArgs(7, 11).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	records, err := repo.ListTaxPeriods(7, 11)
	if err != nil {
		t.Fatalf("ListTaxPeriods returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected no records, got %d", len(records))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
