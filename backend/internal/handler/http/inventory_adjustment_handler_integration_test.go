package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/tests/testdb"

	_ "github.com/lib/pq"
)

func TestInventoryAdjustmentPostgresRollsBackWhenMovementInsertFails(t *testing.T) {
	db := openInventoryRollbackDB(t)
	marker := fmt.Sprintf("ws10-inventory-rollback-%d", time.Now().UnixNano())
	tenantID, restaurantID, userID, productID := createInventoryRollbackFixture(t, db, marker)
	removeMovementFailureTrigger := installMovementFailureTrigger(t, db, marker)
	t.Cleanup(removeMovementFailureTrigger)

	claims := &jwt.Claims{UserID: userID, TenantID: tenantID, RestaurantID: restaurantID}
	body, err := json.Marshal(inventoryAdjustmentRequest{
		ProductID: productID, QuantityChange: 3, Reason: "adjustment", Notes: marker,
	})
	if err != nil {
		t.Fatalf("encode adjustment request: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/inventory/adjust", bytes.NewReader(body))
	request = request.WithContext(context.WithValue(request.Context(), middleware.UserContextKey, claims))
	response := httptest.NewRecorder()

	NewInventoryAdjustmentHandler(db).Adjust(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("adjustment status = %d, want %d; body=%s", response.Code, http.StatusInternalServerError, response.Body.String())
	}
	assertInventoryRollbackState(t, db, productID, tenantID, restaurantID, marker)
}

func openInventoryRollbackDB(t *testing.T) *sql.DB {
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
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close disposable PostgreSQL: %v", err)
		}
	})
	return db
}

func createInventoryRollbackFixture(t *testing.T, db *sql.DB, marker string) (int, int, int, int64) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var tenantID, restaurantID, userID int
	var productID int64
	if err := db.QueryRowContext(ctx, `INSERT INTO tenants (name, slug, domain) VALUES ($1, $2, $3) RETURNING id`, marker, marker, marker+".example.test").Scan(&tenantID); err != nil {
		t.Fatalf("insert synthetic tenant: %v", err)
	}
	if err := db.QueryRowContext(ctx, `INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $3, $4) RETURNING id`, tenantID, marker, marker, "synthetic test address").Scan(&restaurantID); err != nil {
		t.Fatalf("insert synthetic restaurant: %v", err)
	}
	if err := db.QueryRowContext(ctx, `INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, role) VALUES ($1, $2, $3, 'test-hash', 'WS10 rollback fixture', 'owner') RETURNING id`, tenantID, restaurantID, marker+"@example.test").Scan(&userID); err != nil {
		t.Fatalf("insert synthetic user: %v", err)
	}
	if err := db.QueryRowContext(ctx, `INSERT INTO products (tenant_id, restaurant_id, name, price, quantity_in_stock, track_inventory) VALUES ($1, $2, $3, 1, 10, TRUE) RETURNING id`, tenantID, restaurantID, marker).Scan(&productID); err != nil {
		t.Fatalf("insert synthetic product: %v", err)
	}
	return tenantID, restaurantID, userID, productID
}

func installMovementFailureTrigger(t *testing.T, db *sql.DB, marker string) func() {
	t.Helper()
	identifier := "ws10_movement_failure_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	functionName := identifier + "_fn"
	triggerName := identifier + "_trigger"
	createFunction := fmt.Sprintf(`CREATE FUNCTION %s() RETURNS trigger AS $$
BEGIN
	IF NEW.notes = '%s' THEN
		RAISE EXCEPTION 'injected WS10 inventory movement failure';
	END IF;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql`, functionName, marker)
	if _, err := db.Exec(createFunction); err != nil {
		t.Fatalf("create movement failure function: %v", err)
	}
	if _, err := db.Exec(fmt.Sprintf(`CREATE TRIGGER %s BEFORE INSERT ON inventory_movements FOR EACH ROW EXECUTE FUNCTION %s()`, triggerName, functionName)); err != nil {
		db.Exec(fmt.Sprintf(`DROP FUNCTION %s()`, functionName))
		t.Fatalf("create movement failure trigger: %v", err)
	}
	return func() {
		if _, err := db.Exec(fmt.Sprintf(`DROP TRIGGER IF EXISTS %s ON inventory_movements`, triggerName)); err != nil {
			t.Errorf("drop movement failure trigger: %v", err)
		}
		if _, err := db.Exec(fmt.Sprintf(`DROP FUNCTION IF EXISTS %s()`, functionName)); err != nil {
			t.Errorf("drop movement failure function: %v", err)
		}
	}
}

func assertInventoryRollbackState(t *testing.T, db *sql.DB, productID int64, tenantID, restaurantID int, marker string) {
	t.Helper()
	var stock int
	if err := db.QueryRow(`SELECT quantity_in_stock FROM products WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, productID, tenantID, restaurantID).Scan(&stock); err != nil {
		t.Fatalf("read synthetic stock after rejected adjustment: %v", err)
	}
	if stock != 10 {
		t.Fatalf("rejected adjustment left stock at %d, want original 10", stock)
	}
	var movements int
	if err := db.QueryRow(`SELECT COUNT(*) FROM inventory_movements WHERE product_id = $1 AND tenant_id = $2 AND restaurant_id = $3 AND notes = $4`, productID, tenantID, restaurantID, marker).Scan(&movements); err != nil {
		t.Fatalf("count movements after rejected adjustment: %v", err)
	}
	if movements != 0 {
		t.Fatalf("rejected adjustment left %d movement rows, want none", movements)
	}
}
