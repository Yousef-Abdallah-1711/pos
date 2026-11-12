package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

// The hash below identifies accounts that were not rotated after the original migration 109.
const legacySeededSuperAdminHash = "$2a$10$6tgEoGE0dusUMtJ0RUd4VO3IACTDD2vb0676WfMyPvbVQU9oo2p06"

func TestIsSeededDefaultPassword(t *testing.T) {
	if !isSeededDefaultPassword(legacySeededSuperAdminHash) {
		t.Fatal("the legacy migration 109 hash must be rejected")
	}
	rotated, err := bcrypt.GenerateFromPassword([]byte("a-new-strong-password-1"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	if isSeededDefaultPassword(string(rotated)) {
		t.Fatal("a rotated password must not be flagged")
	}
	if isSeededDefaultPassword("") || isSeededDefaultPassword("not-a-hash") {
		t.Fatal("empty or invalid hashes must not be flagged")
	}
}

func TestProductionStartupGuards_NoDatabaseIsRejected(t *testing.T) {
	if err := productionStartupGuards(context.Background(), nil); err == nil {
		t.Fatal("production without a database (stub mode) must be refused")
	}
}

func TestProductionStartupGuards_RejectsUnsafeRuntimeRoles(t *testing.T) {
	tests := []struct {
		name      string
		role      string
		superuser bool
		bypass    bool
	}{
		{name: "superuser", role: "pos_app", superuser: true},
		{name: "bypass RLS", role: "pos_app", bypass: true},
		{name: "unexpected role", role: "postgres"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			mock.ExpectQuery(`SELECT rolname, rolsuper, rolbypassrls FROM pg_roles WHERE rolname = current_user`).
				WillReturnRows(sqlmock.NewRows([]string{"rolname", "rolsuper", "rolbypassrls"}).AddRow(test.role, test.superuser, test.bypass))
			if err := productionStartupGuards(context.Background(), db); err == nil {
				t.Fatal("unsafe runtime role must prevent production startup")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestProductionStartupGuards_AcceptsRestrictedRuntimeRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectQuery(`SELECT rolname, rolsuper, rolbypassrls FROM pg_roles WHERE rolname = current_user`).
		WillReturnRows(sqlmock.NewRows([]string{"rolname", "rolsuper", "rolbypassrls"}).AddRow("pos_app", false, false))
	mock.ExpectQuery(`SELECT password_hash FROM users WHERE email = \$1`).
		WithArgs(seededSuperAdminEmail).WillReturnError(sql.ErrNoRows)
	if err := productionStartupGuards(context.Background(), db); err != nil {
		t.Fatalf("restricted runtime role should pass startup guard: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
