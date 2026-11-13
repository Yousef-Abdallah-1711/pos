package repository

import (
	"context"
	"database/sql"
	"testing"
)

func TestUserRepositoryPostgresScopeDeniesForeignTenantAndRestaurantMutations(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	userA1 := insertScopeUser(t, db, fixture.tenantA, fixture.restaurantA1, "user-scope-a1")
	userA2 := insertScopeUser(t, db, fixture.tenantA, fixture.restaurantA2, "user-scope-a2")
	userB1 := insertScopeUser(t, db, fixture.tenantB, fixture.restaurantB1, "user-scope-b1")
	repo := NewUserRepository(db)
	ctx := context.Background()

	for _, tc := range []struct {
		name     string
		userID   int
		tenantID int
	}{
		{name: "foreign restaurant", userID: userA2, tenantID: fixture.tenantA},
		{name: "foreign tenant", userID: userB1, tenantID: fixture.tenantA},
	} {
		if _, err := repo.GetByID(ctx, tc.userID, tc.tenantID, fixture.restaurantA1); err == nil {
			t.Errorf("%s user read unexpectedly succeeded", tc.name)
		}
		if err := repo.UpdateProfile(ctx, tc.userID, tc.tenantID, fixture.restaurantA1, "forged profile", nil); err == nil {
			t.Errorf("%s profile update unexpectedly succeeded", tc.name)
		}
		if err := repo.UpdatePassword(ctx, tc.userID, tc.tenantID, fixture.restaurantA1, "forged-password-hash"); err == nil {
			t.Errorf("%s password update unexpectedly succeeded", tc.name)
		}
	}

	assertScopedUserState(t, db, userA2, fixture.tenantA, fixture.restaurantA2, "user-scope-a2@example.test", "test-hash")
	assertScopedUserState(t, db, userB1, fixture.tenantB, fixture.restaurantB1, "user-scope-b1@example.test", "test-hash")

	if _, err := repo.GetByID(ctx, userA1, fixture.tenantA, fixture.restaurantA1); err != nil {
		t.Fatalf("in-scope user read failed: %v", err)
	}
	if err := repo.UpdateProfile(ctx, userA1, fixture.tenantA, fixture.restaurantA1, "Authorized profile", nil); err != nil {
		t.Fatalf("in-scope profile update failed: %v", err)
	}
	assertScopedUserName(t, db, userA1, fixture.tenantA, fixture.restaurantA1, "Authorized profile")
}

func assertScopedUserState(t *testing.T, db *sql.DB, userID, tenantID, restaurantID int, email, passwordHash string) {
	t.Helper()
	var gotEmail, gotHash string
	err := db.QueryRow(`SELECT email, password_hash FROM users WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, userID, tenantID, restaurantID).Scan(&gotEmail, &gotHash)
	if err != nil {
		t.Fatalf("read user state after denied operation: %v", err)
	}
	if gotEmail != email || gotHash != passwordHash {
		t.Fatalf("denied user operation changed row: email=%q hash=%q", gotEmail, gotHash)
	}
}

func assertScopedUserName(t *testing.T, db *sql.DB, userID, tenantID, restaurantID int, expectedName string) {
	t.Helper()
	var name string
	if err := db.QueryRow(`SELECT name FROM users WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, userID, tenantID, restaurantID).Scan(&name); err != nil {
		t.Fatalf("read user profile after authorized update: %v", err)
	}
	if name != expectedName {
		t.Fatalf("authorized profile update name=%q, want %q", name, expectedName)
	}
}
