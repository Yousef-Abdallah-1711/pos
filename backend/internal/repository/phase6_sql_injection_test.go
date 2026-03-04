package repository

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/domain"
)

func TestPhase6SQLInjection_ListJournalEntriesTreatsInjectedFiltersAsParameters(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewJournalEntryRepository(db)
	injectedStatus := "posted' OR 1=1 --"
	injectedType := "manual'; DROP TABLE journal_entries; --"

	mock.ExpectQuery(`FROM journal_entries`).
		WithArgs(7, 11, injectedStatus, injectedType).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tenant_id", "restaurant_id", "entry_number", "entry_date", "fiscal_period",
			"description", "description_ar", "reference_type", "reference_id", "reference_number",
			"entry_type", "total_debit", "total_credit", "status", "posted_at", "posted_by",
			"voided_at", "voided_by", "void_reason", "is_approved", "approved_by", "approved_at",
			"notes", "attachments", "created_at", "updated_at", "created_by", "updated_by",
		}))

	entries, err := repo.ListJournalEntries(7, 11, domain.JournalEntryFilters{
		Status:    injectedStatus,
		EntryType: injectedType,
	})
	if err != nil {
		t.Fatalf("ListJournalEntries returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected zero entries, got %d", len(entries))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6SQLInjection_GetAccountByCodeTreatsCodeAsParameter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewAccountRepository(db)
	injectedCode := "1110' UNION SELECT * FROM users --"

	mock.ExpectQuery(`FROM chart_of_accounts`).
		WithArgs(7, 11, injectedCode).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tenant_id", "restaurant_id", "account_code", "account_name", "account_name_ar",
			"description", "description_ar", "account_type", "account_subtype", "normal_balance",
			"parent_id", "level", "is_header", "current_balance", "currency", "is_active",
			"is_system", "is_reconcilable", "tax_rate", "created_at", "updated_at",
		}))

	account, err := repo.GetAccountByCode(7, 11, injectedCode)
	if err != nil {
		t.Fatalf("GetAccountByCode returned error: %v", err)
	}
	if account != nil {
		t.Fatalf("expected nil account for injected code lookup, got %+v", account)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestPhase6SQLInjection_UpdateRuleUsesParameterizedJSONPayload(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewAutoPostingRepository(db)
	rule := &domain.AutoPostingRule{
		ID:           15,
		TenantID:     7,
		RestaurantID: 11,
		RuleName:     "Supplier Payment Rule",
		Description:  "Parameterized payload test",
		IsActive:     true,
		DefaultAccounts: map[string]string{
			"sales": "4100\" || (SELECT password_hash FROM users LIMIT 1) || \"",
		},
	}

	mock.ExpectExec(`UPDATE auto_posting_rules`).
		WithArgs(rule.RuleName, rule.Description, true, sqlmock.AnyArg(), 15, 7, 11).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateRule(rule)
	if err != nil {
		t.Fatalf("UpdateRule returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}

}
