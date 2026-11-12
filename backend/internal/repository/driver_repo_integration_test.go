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

func TestDriverRepositoryPostgresAssignmentRejectsForeignScopeWithoutSideEffects(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	repo := NewDriverRepository(db)
	ownOrder := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA1, "driver-assign-own")
	siblingOrder := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA2, "driver-assign-sibling")
	foreignOrder := insertScopeOrder(t, db, fixture.tenantB, fixture.restaurantB1, "driver-assign-foreign")

	attempts := []struct {
		name     string
		orderID  int64
		driverID int
	}{
		{name: "foreign tenant driver", orderID: ownOrder, driverID: fixture.tenantBDriver},
		{name: "sibling restaurant driver", orderID: ownOrder, driverID: fixture.restaurantA2Driver},
		{name: "foreign tenant order", orderID: foreignOrder, driverID: fixture.restaurantA1Driver},
		{name: "sibling restaurant order", orderID: siblingOrder, driverID: fixture.restaurantA1Driver},
	}

	for _, attempt := range attempts {
		t.Run(attempt.name, func(t *testing.T) {
			assignment := &domain.DriverAssignment{OrderID: int(attempt.orderID), DriverID: attempt.driverID, Status: "pending"}
			if err := repo.AssignOrder(context.Background(), assignment, fixture.tenantA, fixture.restaurantA1); err == nil {
				t.Fatal("foreign assignment unexpectedly succeeded")
			}
			var assignmentCount, activeOrders int
			var availability string
			if err := db.QueryRow(`SELECT COUNT(*) FROM driver_assignments WHERE order_id = $1`, attempt.orderID).Scan(&assignmentCount); err != nil {
				t.Fatalf("count assignment side effects: %v", err)
			}
			if err := db.QueryRow(`SELECT active_orders_count, availability_status FROM drivers WHERE id = $1`, attempt.driverID).Scan(&activeOrders, &availability); err != nil {
				t.Fatalf("read driver state: %v", err)
			}
			if assignmentCount != 0 || activeOrders != 0 || availability != "available" {
				t.Fatalf("rejected assignment changed state: assignments=%d active_orders=%d availability=%q", assignmentCount, activeOrders, availability)
			}
		})
	}
}

func TestDriverRepositoryPostgresCreateRequiresRestaurantTenantMatch(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	driver := &domain.Driver{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantB1,
		FirstName: "Forged", LastName: "Scope", Email: fmt.Sprintf("forged-scope-%d@example.test", fixture.tenantA),
		PhoneNumber: "000", LicenseNumber: fmt.Sprintf("FORGED-SCOPE-%d", fixture.tenantA), Status: "active", AvailabilityStatus: "available",
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
	tenantA, tenantB                                      int
	restaurantA1, restaurantA2                            int
	restaurantB1, restaurantB2                            int
	tenantBDriver, restaurantA2Driver, restaurantA1Driver int
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
		// Audit triggers reject tenant cascades after their parent row is removed.
		// Fixtures therefore live only in this disposable database, whose tmpfs
		// container is discarded after the integration run.
		if err := db.Close(); err != nil {
			t.Errorf("close disposable PostgreSQL connection: %v", err)
		}
	})
	return db, fixture
}

func seedDriverScopeFixture(t *testing.T, ctx context.Context, db *sql.DB) driverScopeFixture {
	t.Helper()
	var fixture driverScopeFixture
	runID := fmt.Sprintf("%d", time.Now().UnixNano())
	fixture.tenantA = insertTenant(t, ctx, db, "closure-driver-tenant-a-"+runID)
	fixture.tenantB = insertTenant(t, ctx, db, "closure-driver-tenant-b-"+runID)
	fixture.restaurantA1 = insertRestaurant(t, ctx, db, fixture.tenantA, "closure-driver-a1-"+runID)
	fixture.restaurantA2 = insertRestaurant(t, ctx, db, fixture.tenantA, "closure-driver-a2-"+runID)
	fixture.restaurantB1 = insertRestaurant(t, ctx, db, fixture.tenantB, "closure-driver-b1-"+runID)
	fixture.restaurantB2 = insertRestaurant(t, ctx, db, fixture.tenantB, "closure-driver-b2-"+runID)
	fixture.restaurantA1Driver = insertDriver(t, ctx, db, fixture.tenantA, fixture.restaurantA1, runID+"-a1")
	fixture.restaurantA2Driver = insertDriver(t, ctx, db, fixture.tenantA, fixture.restaurantA2, runID+"-a2")
	fixture.tenantBDriver = insertDriver(t, ctx, db, fixture.tenantB, fixture.restaurantB1, runID+"-b1")
	insertDriver(t, ctx, db, fixture.tenantB, fixture.restaurantB2, runID+"-b2")
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
