package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"sync"
	"testing"

	"pos-saas/internal/domain"
)

type stockFixture struct {
	db                          *sql.DB
	repo                        *OrderRepository
	tenant, restaurant, product int
}

func newStockFixture(t *testing.T, stock int) stockFixture {
	t.Helper()
	db, scope := openProductScopeFixture(t)
	category := insertScopeCategory(t, db, scope.tenantA, scope.restaurantA1, "stock")
	product := insertScopeProduct(t, db, scope.tenantA, scope.restaurantA1, category, "stock")
	if _, err := db.Exec(`UPDATE products SET track_inventory = true, quantity_in_stock = $1 WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4`, stock, product, scope.tenantA, scope.restaurantA1); err != nil {
		t.Fatal(err)
	}
	return stockFixture{db, NewOrderRepository(db), scope.tenantA, scope.restaurantA1, product}
}

func (f stockFixture) order(t *testing.T, suffix string, quantities ...int) int64 {
	t.Helper()
	id := insertScopeOrder(t, f.db, f.tenant, f.restaurant, suffix)
	for _, quantity := range quantities {
		_, err := f.repo.CreateOrderItem(int64(f.tenant), int64(f.restaurant), &domain.OrderItem{OrderID: id, ProductID: int64(f.product), ProductName: "Stock test", Quantity: quantity, UnitPrice: 10, TotalPrice: float64(quantity) * 10})
		if err != nil {
			t.Fatal(err)
		}
	}
	return id
}

func (f stockFixture) assertStock(t *testing.T, want int) {
	t.Helper()
	var got int
	if err := f.db.QueryRow(`SELECT quantity_in_stock FROM products WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, f.product, f.tenant, f.restaurant).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("stock = %d; want %d", got, want)
	}
}

func (f stockFixture) transition(id int64, status string) error {
	return f.repo.UpdateOrderStatus(int64(f.tenant), int64(f.restaurant), id, status, nil, "stock acceptance")
}

func TestOrderStockAggregatesDuplicateLinesAndDeductsOnce(t *testing.T) {
	f := newStockFixture(t, 20)
	id := f.order(t, "stock-once", 2, 3)
	if err := f.transition(id, "confirmed"); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 15)
	if err := f.transition(id, "confirmed"); err == nil {
		t.Fatal("duplicate transition accepted")
	}
	for _, status := range []string{"preparing", "ready", "out_for_delivery", "delivered"} {
		if err := f.transition(id, status); err != nil {
			t.Fatal(err)
		}
	}
	for _, payment := range []string{"paid", "paid", "refunded", "refunded"} {
		if err := f.repo.UpdatePaymentStatus(int64(f.tenant), int64(f.restaurant), id, payment); err != nil {
			t.Fatal(err)
		}
	}
	// A monetary refund does not imply physical inventory was returned.
	f.assertStock(t, 15)
}

func TestOrderStockCancellationReleasesExactlyOnce(t *testing.T) {
	f := newStockFixture(t, 20)
	id := f.order(t, "stock-cancel", 2)
	if err := f.transition(id, "confirmed"); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 18)
	if err := f.transition(id, "cancelled"); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 20)
	if err := f.transition(id, "cancelled"); err == nil {
		t.Fatal("duplicate cancellation accepted")
	}
	f.assertStock(t, 20)
	pending := f.order(t, "stock-pending-cancel", 2)
	if err := f.transition(pending, "cancelled"); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 20)
}

func TestOrderStockConcurrentConfirmationCannotOversell(t *testing.T) {
	f := newStockFixture(t, 10)
	ids := []int64{f.order(t, "stock-race-a", 8), f.order(t, "stock-race-b", 8)}
	start := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for _, id := range ids {
		workers.Add(1)
		go func(id int64) { defer workers.Done(); <-start; results <- f.transition(id, "confirmed") }(id)
	}
	close(start)
	workers.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		}
	}
	if success != 1 {
		t.Fatalf("confirmed %d orders, want exactly one", success)
	}
	f.assertStock(t, 2)
	var pending, allocations int
	if err := f.db.QueryRow(`SELECT count(*) FROM orders WHERE id = ANY($1) AND tenant_id = $2 AND restaurant_id = $3 AND status = 'pending'`, pq.Array(ids), f.tenant, f.restaurant).Scan(&pending); err != nil {
		t.Fatal(err)
	}
	if err := f.db.QueryRow(`SELECT count(*) FROM order_stock_reservations WHERE tenant_id = $1 AND restaurant_id = $2`, f.tenant, f.restaurant).Scan(&allocations); err != nil {
		t.Fatal(err)
	}
	if pending != 1 || allocations != 1 {
		t.Fatalf("pending=%d allocations=%d, want 1/1", pending, allocations)
	}
}

func TestOrderStockAuditFailureRollsBackReservation(t *testing.T) {
	f := newStockFixture(t, 20)
	id := f.order(t, "stock-audit-rollback", 2)
	missingActor := int64(9223372036854775807)
	if err := f.repo.UpdateOrderStatus(int64(f.tenant), int64(f.restaurant), id, "confirmed", &missingActor, "rollback"); err == nil {
		t.Fatal("invalid audit actor accepted")
	}
	f.assertStock(t, 20)
	assertScopeOrderState(t, f.db, id, f.tenant, f.restaurant, "pending", "pending")
	if err := f.transition(id, "confirmed"); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 18)
}

func TestOrderStockInsufficientQuantityAndScopeDenyChanges(t *testing.T) {
	f := newStockFixture(t, 3)
	id := f.order(t, "stock-insufficient", 2, 2)
	if err := f.transition(id, "confirmed"); err == nil {
		t.Fatal("aggregate insufficient stock accepted")
	}
	f.assertStock(t, 3)
	assertScopeOrderState(t, f.db, id, f.tenant, f.restaurant, "pending", "pending")
	if err := f.repo.UpdateOrderStatus(int64(f.tenant+1), int64(f.restaurant), id, "confirmed", nil, ""); err == nil {
		t.Fatal("foreign tenant confirmed order")
	}
	f.assertStock(t, 3)
}

func TestOrderStockRejectsItemEditsAfterConfirmation(t *testing.T) {
	f := newStockFixture(t, 20)
	id := f.order(t, "stock-immutable", 2)
	if err := f.transition(id, "confirmed"); err != nil {
		t.Fatal(err)
	}
	if _, err := f.db.Exec(`UPDATE order_items SET quantity = 5 WHERE order_id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, f.tenant, f.restaurant); err == nil {
		t.Fatal("confirmed order item edited")
	}
	if _, err := f.db.Exec(`DELETE FROM orders WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, f.tenant, f.restaurant); err == nil {
		t.Fatal("inventory history deleted")
	}
	f.assertStock(t, 18)
}

func TestOrderCreationRollsBackAllItemsWhenAnyItemIsInvalid(t *testing.T) {
	f := newStockFixture(t, 20)
	order := &domain.Order{TenantID: int64(f.tenant), RestaurantID: int64(f.restaurant), OrderNumber: "invalid-items-transaction", CustomerName: "Synthetic", CustomerPhone: "000", PaymentMethod: "cash", PaymentStatus: "pending", Status: "pending", Items: []domain.OrderItem{{ProductID: int64(f.product), ProductName: "Valid", Quantity: 1, UnitPrice: 10, TotalPrice: 10}, {ProductID: 9223372036854775807, ProductName: "Invalid", Quantity: 1, UnitPrice: 10, TotalPrice: 10}}}
	if _, err := f.repo.CreateOrderWithItems(order); err == nil {
		t.Fatal("invalid child item accepted")
	}
	var count int
	if err := f.db.QueryRow(`SELECT count(*) FROM orders WHERE tenant_id=$1 AND restaurant_id=$2`, f.tenant, f.restaurant).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("partial order committed; count=%d", count)
	}
	f.assertStock(t, 20)
}
