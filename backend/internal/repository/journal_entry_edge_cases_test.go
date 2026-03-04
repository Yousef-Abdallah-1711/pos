package repository

import (
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/domain"
)

func TestValidateJournalLines_EdgeCases(t *testing.T) {
	validLines := []domain.JournalEntryLine{
		{AccountID: 101, DebitAmount: 100},
		{AccountID: 202, CreditAmount: 100},
	}

	tests := []struct {
		name       string
		lines      []domain.JournalEntryLine
		wantErrMsg string
	}{
		{
			name:       "requires at least two lines",
			lines:      []domain.JournalEntryLine{{AccountID: 101, DebitAmount: 100}},
			wantErrMsg: "at least two lines",
		},
		{
			name: "rejects zero amount lines",
			lines: []domain.JournalEntryLine{
				{AccountID: 101, DebitAmount: 0, CreditAmount: 0},
				{AccountID: 202, CreditAmount: 0, DebitAmount: 0},
			},
			wantErrMsg: "must have a debit or credit amount",
		},
		{
			name: "rejects negative amount lines",
			lines: []domain.JournalEntryLine{
				{AccountID: 101, DebitAmount: -1},
				{AccountID: 202, CreditAmount: 1},
			},
			wantErrMsg: "cannot be negative",
		},
		{
			name: "rejects debit and credit in same line",
			lines: []domain.JournalEntryLine{
				{AccountID: 101, DebitAmount: 50, CreditAmount: 50},
				{AccountID: 202, CreditAmount: 100},
			},
			wantErrMsg: "cannot be both debit and credit",
		},
		{
			name: "rejects unbalanced lines",
			lines: []domain.JournalEntryLine{
				{AccountID: 101, DebitAmount: 120},
				{AccountID: 202, CreditAmount: 100},
			},
			wantErrMsg: "not balanced",
		},
		{
			name:       "accepts balanced lines",
			lines:      validLines,
			wantErrMsg: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := validateJournalLines(tc.lines)
			if tc.wantErrMsg == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.wantErrMsg)
			}
			if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tc.wantErrMsg)) {
				t.Fatalf("expected error containing %q, got %q", tc.wantErrMsg, err.Error())
			}
		})
	}
}

func TestJournalEntryRepository_UpdateDraftEntryRejectsPostedEntries(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewJournalEntryRepository(db)
	entry := &domain.JournalEntry{
		ID:           55,
		TenantID:     1,
		RestaurantID: 2,
		EntryDate:    time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC),
		Lines: []domain.JournalEntryLine{
			{AccountID: 101, DebitAmount: 100},
			{AccountID: 202, CreditAmount: 100},
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT status FROM journal_entries`).
		WithArgs(entry.ID, entry.TenantID, entry.RestaurantID).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(domain.JournalStatusPosted))
	mock.ExpectRollback()

	err = repo.UpdateDraftEntry(entry)
	if err == nil || !strings.Contains(err.Error(), "only draft journal entries can be updated") {
		t.Fatalf("expected draft-only error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestJournalEntryRepository_PostEntryRejectsNonDraft(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewJournalEntryRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT status, entry_date`).
		WithArgs(88, 1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"status", "entry_date", "fiscal_period"}).
			AddRow(domain.JournalStatusPosted, time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC), "2026-03"))
	mock.ExpectRollback()

	err = repo.PostEntry(1, 2, 88, 7)
	if err == nil || !strings.Contains(err.Error(), "only draft journal entries can be posted") {
		t.Fatalf("expected draft-only post error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestJournalEntryRepository_PostEntryRejectsClosedFiscalPeriod(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewJournalEntryRepository(db)
	entryDate := time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT status, entry_date`).
		WithArgs(90, 1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"status", "entry_date", "fiscal_period"}).
			AddRow(domain.JournalStatusDraft, entryDate, "2026-03"))
	mock.ExpectQuery(`SELECT status\s+FROM fiscal_periods`).
		WithArgs(1, 2, "2026-03", entryDate).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("closed"))
	mock.ExpectRollback()

	err = repo.PostEntry(1, 2, 90, 7)
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "fiscal period is closed") {
		t.Fatalf("expected closed fiscal period error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestJournalEntryRepository_VoidEntryRejectsAlreadyVoid(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewJournalEntryRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT status FROM journal_entries`).
		WithArgs(89, 1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(domain.JournalStatusVoid))
	mock.ExpectRollback()

	err = repo.VoidEntry(1, 2, 89, 7, "already void")
	if err == nil || !strings.Contains(err.Error(), "already void") {
		t.Fatalf("expected already-void error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
