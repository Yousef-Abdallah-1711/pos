package usecase

import (
	"errors"
	"math"
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/domain"
)

func TestValidateAutoPostingLinesRejectsBothDirectionsOfImbalance(t *testing.T) {
	tests := []struct {
		name  string
		lines []domain.JournalEntryLine
	}{
		{
			name: "debit heavy",
			lines: []domain.JournalEntryLine{
				{AccountID: 1, DebitAmount: 100},
				{AccountID: 2, CreditAmount: 90},
			},
		},
		{
			name: "credit heavy",
			lines: []domain.JournalEntryLine{
				{AccountID: 1, DebitAmount: 90},
				{AccountID: 2, CreditAmount: 100},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := validateAutoPostingLines(test.lines)
			var validationErr *AutoPostingValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected typed validation error, got %v", err)
			}
			if validationErr.Code != "unbalanced_journal" {
				t.Fatalf("expected unbalanced_journal code, got %q", validationErr.Code)
			}
			if validationErr.Error() == "" {
				t.Fatal("expected actionable validation message")
			}
		})
	}
}

func TestValidateAutoPostingLinesAcceptsBalancedAmountsWithoutMutatingInput(t *testing.T) {
	lines := []domain.JournalEntryLine{
		{AccountID: 1, DebitAmount: 100},
		{AccountID: 2, CreditAmount: 100},
	}
	original := append([]domain.JournalEntryLine(nil), lines...)

	validated, err := validateAutoPostingLines(lines)
	if err != nil {
		t.Fatalf("expected balanced lines to validate, got %v", err)
	}
	if !reflect.DeepEqual(lines, original) {
		t.Fatalf("validation mutated input lines: got %#v, want %#v", lines, original)
	}
	if len(validated) != 2 || validated[0].DebitAmount != 100 || validated[1].CreditAmount != 100 {
		t.Fatalf("unexpected validated lines: %#v", validated)
	}
	if validated[0].LineOrder != 1 || validated[1].LineOrder != 2 {
		t.Fatalf("expected missing line orders to be assigned on copies, got %#v", validated)
	}
}

func TestAutoPostingRejectsUnbalancedEntriesBeforeJournalOrLedgerWrites(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	engine := newAutoPostingEngineForTest(db)
	tests := []struct {
		name  string
		lines []domain.JournalEntryLine
	}{
		{
			name: "debit heavy",
			lines: []domain.JournalEntryLine{
				{AccountID: 1, DebitAmount: 100},
				{AccountID: 2, CreditAmount: 90},
			},
		},
		{
			name: "credit heavy",
			lines: []domain.JournalEntryLine{
				{AccountID: 1, DebitAmount: 90},
				{AccountID: 2, CreditAmount: 100},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event := &domain.PostingEvent{
				TenantID:     10,
				RestaurantID: 20,
				Description:  "validation regression",
			}
			err := engine.PostSupplierInvoiceEntry(event, test.lines)
			var validationErr *AutoPostingValidationError
			if !errors.As(err, &validationErr) || validationErr.Code != "unbalanced_journal" {
				t.Fatalf("expected unbalanced_journal validation error, got %v", err)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("invalid entries reached a journal or ledger repository: %v", err)
	}
}

func TestValidateAutoPostingLinesRejectsMalformedAmounts(t *testing.T) {
	validCredit := domain.JournalEntryLine{AccountID: 2, CreditAmount: 1}
	tests := []struct {
		name  string
		lines []domain.JournalEntryLine
	}{
		{name: "negative debit", lines: []domain.JournalEntryLine{{AccountID: 1, DebitAmount: -1}, validCredit}},
		{name: "not a number", lines: []domain.JournalEntryLine{{AccountID: 1, DebitAmount: math.NaN()}, validCredit}},
		{name: "positive infinity", lines: []domain.JournalEntryLine{{AccountID: 1, DebitAmount: math.Inf(1)}, validCredit}},
		{name: "negative infinity", lines: []domain.JournalEntryLine{{AccountID: 1, DebitAmount: math.Inf(-1)}, validCredit}},
		{name: "sub-cent precision", lines: []domain.JournalEntryLine{{AccountID: 1, DebitAmount: 1.001}, validCredit}},
		{name: "both sides on one line", lines: []domain.JournalEntryLine{{AccountID: 1, DebitAmount: 1, CreditAmount: 1}, validCredit}},
		{name: "zero line", lines: []domain.JournalEntryLine{{AccountID: 1}, validCredit}},
		{name: "amount beyond database precision", lines: []domain.JournalEntryLine{{AccountID: 1, DebitAmount: 10_000_000_000_000}, validCredit}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := validateAutoPostingLines(test.lines)
			var validationErr *AutoPostingValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected typed validation error, got %v", err)
			}
		})
	}
}
