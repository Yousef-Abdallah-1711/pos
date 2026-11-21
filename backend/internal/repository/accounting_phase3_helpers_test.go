package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLookupAccountByCodeTxUsesTransactionConnection(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	defer tx.Rollback()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id
		FROM chart_of_accounts
		WHERE tenant_id = $1 AND restaurant_id = $2 AND account_code = $3
	`)).WithArgs(7, 11, "2101").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	accountID, err := lookupAccountByCodeTx(tx, 7, 11, "2101")
	if err != nil || accountID != 42 {
		t.Fatalf("lookup account through transaction: id=%d err=%v", accountID, err)
	}
	mock.ExpectCommit()
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit transaction: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL expectations: %v", err)
	}
}
