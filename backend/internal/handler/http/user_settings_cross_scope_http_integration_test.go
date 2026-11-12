package http

import (
	"database/sql"
	"os"
	"testing"

	"pos-saas/tests/testdb"
)

func TestUserSettingsHTTPRejectsSiblingAndForeignTenantUsers(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	actorID := insertSettingsHTTPUser(t, db, fixture.tenantA, fixture.restaurantA1, "settings-cross-scope-actor")
	targets := []struct {
		name, suffix           string
		tenantID, restaurantID int
	}{
		{name: "sibling restaurant", suffix: "settings-sibling-target", tenantID: fixture.tenantA, restaurantID: fixture.restaurantA2},
		{name: "foreign tenant", suffix: "settings-foreign-target", tenantID: fixture.tenantB, restaurantID: fixture.restaurantB1},
	}
	for _, target := range targets {
		t.Run(target.name, func(t *testing.T) {
			verifyCrossScopeSettingsTarget(t, db, fixture, actorID, target.tenantID, target.restaurantID, target.suffix)
		})
	}
}

func verifyCrossScopeSettingsTarget(t *testing.T, db *sql.DB, fixture productHTTPFixture, actorID int64, tenantID, restaurantID int, suffix string) {
	t.Helper()
	targetID := insertSettingsHTTPUser(t, db, tenantID, restaurantID, suffix)
	if _, err := db.Exec(`INSERT INTO user_settings (user_id, language) VALUES ($1, 'en')`, targetID); err != nil {
		t.Fatalf("insert target settings: %v", err)
	}
	assertCrossUserSettingsRoutesDenied(t, db, fixture, actorID, targetID)
	assertSettingsUserStateUnchanged(t, db, targetID, suffix+"@example.test")
}
