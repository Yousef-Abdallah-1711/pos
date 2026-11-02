package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/tests/testdb"
)

func TestDriverRepositoryPostgresScopeDeniesForeignTenantWrites(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	repo := NewDriverRepository(db)
	ctx := context.Background()

	if _, err := repo.GetByID(ctx, fixture.tenantBDriver, fixture.tenantA, fixture.restaurantA1); err == nil {
		t.Fatal("cross-tenant driver read unexpectedly succeeded")
	}
	if err := repo.UpdateStatus(ctx, fixture.tenantBDriver, fixture.tenantA, fixture.restaurantA1, "inactive"); err == nil {
		t.Fatal("cross-tenant status update unexpectedly succeeded")
	}
	if err := repo.Delete(ctx, fixture.tenantBDriver, fixture.tenantA, fixture.restaurantA1); err == nil {
		t.Fatal("cross-tenant delete unexpectedly succeeded")
	}
	assertDriverState(t, db, fixture.tenantBDriver, fixture.tenantB, fixture.restaurantB1, "active")
}

func TestDriverRepositoryPostgresScopeDeniesForeignRestaurantWrites(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	repo := NewDriverRepository(db)
	ctx := context.Background()

	if _, err := repo.GetByID(ctx, fixture.restaurantA2Driver, fixture.tenantA, fixture.restaurantA1); err == nil {
		t.Fatal("cross-restaurant driver read unexpectedly succeeded")
	}
	if err := repo.UpdateStatus(ctx, fixture.restaurantA2Driver, fixture.tenantA, fixture.restaurantA1, "inactive"); err == nil {
		t.Fatal("cross-restaurant status update unexpectedly succeeded")
	}
	if err := repo.Delete(ctx, fixture.restaurantA2Driver, fixture.tenantA, fixture.restaurantA1); err == nil {
		t.Fatal("cross-restaurant delete unexpectedly succeeded")
	}
	assertDriverState(t, db, fixture.restaurantA2Driver, fixture.tenantA, fixture.restaurantA2, "active")
}

func TestDriverRepositoryPostgresCreateRequiresRestaurantTenantMatch(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	driver := &domain.Driver{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantB1,
		FirstName: "Forged", LastName: "Scope", Email: "forged-scope@example.test",
		PhoneNumber: "000", LicenseNumber: "FORGED-SCOPE", Status: "active", AvailabilityStatus: "available",
	}
	if err := NewDriverRepository(db).Create(context.Background(), driver); err == nil {
		t.Fatal("driver creation with a foreign restaurant unexpectedly succeeded")
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM drivers WHERE email = $1`, driver.Email).Scan(&count); err != nil {
		t.Fatalf("count forged driver rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("forged driver creation left %d rows", count)
	}
}

type driverScopeFixture struct {
	tenantA, tenantB                  int
	restaurantA1, restaurantA2        int
	restaurantB1, restaurantB2        int
	tenantBDriver, restaurantA2Driver int
}

func openDriverScopeFixture(t *testing.T) (*sql.DB, driverScopeFixture) {
	t.Helper()
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}

	fixture := seedDriverScopeFixture(t, ctx, db)
	t.Cleanup(func() {
		_, _ = db.Exec(`DELETE FROM tenants WHERE id IN ($1, $2)`, fixture.tenantA, fixture.tenantB)
		_ = db.Close()
	})
	return db, fixture
}

func seedDriverScopeFixture(t *testing.T, ctx context.Context, db *sql.DB) driverScopeFixture {
	t.Helper()
	var fixture driverScopeFixture
	fixture.tenantA = insertTenant(t, ctx, db, "closure-driver-tenant-a")
	fixture.tenantB = insertTenant(t, ctx, db, "closure-driver-tenant-b")
	fixture.restaurantA1 = insertRestaurant(t, ctx, db, fixture.tenantA, "closure-driver-a1")
	fixture.restaurantA2 = insertRestaurant(t, ctx, db, fixture.tenantA, "closure-driver-a2")
	fixture.restaurantB1 = insertRestaurant(t, ctx, db, fixture.tenantB, "closure-driver-b1")
	fixture.restaurantB2 = insertRestaurant(t, ctx, db, fixture.tenantB, "closure-driver-b2")
	insertDriver(t, ctx, db, fixture.tenantA, fixture.restaurantA1, "a1")
	fixture.restaurantA2Driver = insertDriver(t, ctx, db, fixture.tenantA, fixture.restaurantA2, "a2")
	fixture.tenantBDriver = insertDriver(t, ctx, db, fixture.tenantB, fixture.restaurantB1, "b1")
	insertDriver(t, ctx, db, fixture.tenantB, fixture.restaurantB2, "b2")
	return fixture
}

func insertTenant(t *testing.T, ctx context.Context, db *sql.DB, key string) int {
	t.Helper()
	var id int
	err := db.QueryRowContext(ctx, `INSERT INTO tenants (name, slug, domain) VALUES ($1, $2, $3) RETURNING id`, key, key, key+".example.test").Scan(&id)
	if err != nil {
		t.Fatalf("insert %s tenant: %v", key, err)
	}
	return id
}

func insertRestaurant(t *testing.T, ctx context.Context, db *sql.DB, tenantID int, key string) int {
	t.Helper()
	var id int
	err := db.QueryRowContext(ctx, `INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $3, $4) RETURNING id`, tenantID, key, key, "synthetic address").Scan(&id)
	if err != nil {
		t.Fatalf("insert %s restaurant: %v", key, err)
	}
	return id
}

func insertDriver(t *testing.T, ctx context.Context, db *sql.DB, tenantID, restaurantID int, suffix string) int {
	t.Helper()
	var id int
	err := db.QueryRowContext(ctx, `INSERT INTO drivers (tenant_id, restaurant_id, first_name, last_name, email, phone_number, license_number, status, availability_status) VALUES ($1, $2, $3, 'Test', $4, '000', $5, 'active', 'available') RETURNING id`, tenantID, restaurantID, "Driver "+suffix, fmt.Sprintf("driver-%s@example.test", suffix), "LICENSE-"+suffix).Scan(&id)
	if err != nil {
		t.Fatalf("insert %s driver: %v", suffix, err)
	}
	return id
}

func assertDriverState(t *testing.T, db *sql.DB, id, tenantID, restaurantID int, expectedStatus string) {
	t.Helper()
	var count int
	var status string
	err := db.QueryRow(`SELECT COUNT(*), COALESCE(MAX(status), '') FROM drivers WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, tenantID, restaurantID).Scan(&count, &status)
	if err != nil {
		t.Fatalf("read driver state after denied operation: %v", err)
	}
	if count != 1 || status != expectedStatus {
		t.Fatalf("denied operation changed driver state: count=%d status=%q, want count=1 status=%q", count, status, expectedStatus)
	}
}
