package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestWithTenantTx_SetsTransactionLocalContextThenCommits(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`SELECT set_config\(\$1, \$2, true\), set_config\(\$3, \$4, true\)`).
		WithArgs(TenantSetting, "7", RestaurantSetting, "11").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	called := false
	if err := WithTenantTx(context.Background(), db, 7, 11, func(tx *sql.Tx) error {
		called = true
		return nil
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("fn must run")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestWithTenantTx_TenantWideWorkHasEmptyRestaurant(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec(`set_config`).WithArgs(TenantSetting, "7", RestaurantSetting, "").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	if err := WithTenantTx(context.Background(), db, 7, 0, func(*sql.Tx) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestWithTenantTx_RollsBackWhenFnFails(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec(`set_config`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	boom := errors.New("boom")
	if err := WithTenantTx(context.Background(), db, 1, 1, func(*sql.Tx) error { return boom }); !errors.Is(err, boom) {
		t.Fatalf("want boom, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestWithTenantTx_RollsBackWhenContextCannotBeSet(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec(`set_config`).WillReturnError(errors.New("no such setting"))
	mock.ExpectRollback()

	called := false
	if err := WithTenantTx(context.Background(), db, 1, 1, func(*sql.Tx) error { called = true; return nil }); err == nil {
		t.Fatal("must fail")
	}
	if called {
		t.Fatal("fn must not run without a tenant context")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestWithTenantTx_RollsBackAndRepanics(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec(`set_config`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	defer func() {
		if recover() == nil {
			t.Fatal("panic must propagate")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	}()
	_ = WithTenantTx(context.Background(), db, 1, 1, func(*sql.Tx) error { panic("handler bug") })
}

func TestWithTenantTx_RejectsInvalidIdentityWithoutTouchingTheDatabase(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	for _, c := range [][2]int64{{0, 1}, {-3, 1}, {1, -1}} {
		if err := WithTenantTx(context.Background(), db, c[0], c[1], func(*sql.Tx) error { return nil }); err == nil {
			t.Errorf("tenant=%d restaurant=%d must be rejected", c[0], c[1])
		}
	}
	if err := WithTenantTx(context.Background(), nil, 1, 1, func(*sql.Tx) error { return nil }); err == nil {
		t.Error("nil db must be rejected")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err) // no expectations were set: any DB call would have failed above
	}
}
