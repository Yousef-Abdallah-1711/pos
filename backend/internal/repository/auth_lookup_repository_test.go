package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAuthLookupRepositoryFindsVerifiedIdentityByNormalizedEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	wantTime := time.Date(2026, 11, 22, 20, 0, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta("FROM auth_private.lookup_verified_identity_by_email($1)")).
		WithArgs("owner@example.test").
		WillReturnRows(sqlmock.NewRows([]string{"identity_id", "normalized_email", "email_verified_at", "password_hash", "security_version"}).
			AddRow("00000000-0000-4000-8000-000000000001", "owner@example.test", wantTime, "$2a$synthetic-hash", int64(4)))

	identity, err := NewAuthLookupRepository(db).FindVerifiedIdentityByEmail(context.Background(), "  Owner@Example.Test ")
	if err != nil {
		t.Fatal(err)
	}
	if identity.ID != "00000000-0000-4000-8000-000000000001" || identity.NormalizedEmail != "owner@example.test" ||
		identity.EmailVerifiedAt != wantTime || identity.PasswordHash != "$2a$synthetic-hash" || identity.SecurityVersion != 4 {
		t.Fatalf("unexpected identity lookup result: %#v", identity)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAuthLookupRepositoryListsOnlyResolvedMembershipChoices(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta("FROM auth_private.lookup_eligible_memberships($1::uuid)")).
		WithArgs("00000000-0000-4000-8000-000000000001").
		WillReturnRows(sqlmock.NewRows([]string{"membership_id", "tenant_id", "legacy_user_id", "tenant_name", "restaurant_id", "restaurant_name", "role_id", "role_code"}).
			AddRow("00000000-0000-4000-8000-000000000002", 7, int64(42), "Tenant 7", 11, "Restaurant 11", 13, "CASHIER"))

	memberships, err := NewAuthLookupRepository(db).FindEligibleMemberships(context.Background(), "00000000-0000-4000-8000-000000000001")
	if err != nil {
		t.Fatal(err)
	}
	if len(memberships) != 1 || memberships[0].TenantID != 7 || memberships[0].LegacyUserID != (sql.NullInt64{Int64: 42, Valid: true}) ||
		memberships[0].RestaurantID != 11 || memberships[0].RoleCode != "CASHIER" {
		t.Fatalf("unexpected membership choices: %#v", memberships)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAuthLookupRepositoryRejectsUnavailablePool(t *testing.T) {
	if _, err := NewAuthLookupRepository(nil).FindVerifiedIdentityByEmail(context.Background(), "owner@example.test"); err != sql.ErrConnDone {
		t.Fatalf("nil pool identity lookup error = %v", err)
	}
	if _, err := NewAuthLookupRepository(nil).FindEligibleMemberships(context.Background(), "00000000-0000-4000-8000-000000000001"); err != sql.ErrConnDone {
		t.Fatalf("nil pool membership lookup error = %v", err)
	}
}
