package repository

import (
	"database/sql"
	"sync"
	"testing"

	"pos-saas/internal/domain"
)

func TestOrderRepositoryPostgresScopeDeniesForeignReadsAndMutations(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	orderA1 := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA1, "order-a1")
	orderA2 := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA2, "order-a2")
	orderB1 := insertScopeOrder(t, db, fixture.tenantB, fixture.restaurantB1, "order-b1")
	repo := NewOrderRepository(db)

	for _, tc := range []struct {
		name       string
		tenantID   int64
		restaurant int64
		orderID    int64
	}{
		{name: "foreign restaurant", tenantID: int64(fixture.tenantA), restaurant: int64(fixture.restaurantA1), orderID: orderA2},
		{name: "foreign tenant", tenantID: int64(fixture.tenantA), restaurant: int64(fixture.restaurantA1), orderID: orderB1},
	} {
		if _, err := repo.GetOrderByID(tc.tenantID, tc.restaurant, tc.orderID); err == nil {
			t.Errorf("%s order read unexpectedly succeeded", tc.name)
		}
		if err := repo.UpdatePaymentStatus(tc.tenantID, tc.restaurant, tc.orderID, "paid"); err == nil {
			t.Errorf("%s order payment status update unexpectedly succeeded", tc.name)
		}
		if err := repo.DeleteOrder(tc.tenantID, tc.restaurant, tc.orderID); err == nil {
			t.Errorf("%s order delete unexpectedly succeeded", tc.name)
		}
	}

	assertScopeOrderState(t, db, orderA2, fixture.tenantA, fixture.restaurantA2, "pending", "pending")
	assertScopeOrderState(t, db, orderB1, fixture.tenantB, fixture.restaurantB1, "pending", "pending")

	if _, err := repo.GetOrderByID(int64(fixture.tenantA), int64(fixture.restaurantA1), orderA1); err != nil {
		t.Fatalf("in-scope order read failed: %v", err)
	}
}

func TestOrderStatusTransitionRollsBackWhenAuditInsertFails(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	orderID := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA1, "status-rollback")
	repo := NewOrderRepository(db)
	invalidActorID := int64(9223372036854775807)

	if err := repo.UpdateOrderStatus(int64(fixture.tenantA), int64(fixture.restaurantA1), orderID, "confirmed", &invalidActorID, "test rollback"); err == nil {
		t.Fatal("status transition unexpectedly succeeded with a nonexistent audit actor")
	}
	assertScopeOrderState(t, db, orderID, fixture.tenantA, fixture.restaurantA1, "pending", "pending")
	assertOrderHistoryCount(t, db, orderID, 0)
}

func TestConcurrentOrderStatusTransitionsSerializeOnCurrentStatus(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	orderID := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA1, "status-concurrency")
	repo := NewOrderRepository(db)
	start := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup

	for _, status := range []string{"confirmed", "preparing"} {
		workers.Add(1)
		go func(nextStatus string) {
			defer workers.Done()
			<-start
			results <- repo.UpdateOrderStatus(int64(fixture.tenantA), int64(fixture.restaurantA1), orderID, nextStatus, nil, "concurrent transition")
		}(status)
	}
	close(start)
	workers.Wait()
	close(results)

	successes := 0
	failures := 0
	for err := range results {
		if err == nil {
			successes++
		} else {
			failures++
		}
	}
	if successes != 1 || failures != 1 {
		t.Fatalf("got %d successful and %d failed transitions, want one of each", successes, failures)
	}

	var status string
	if err := db.QueryRow(`SELECT status FROM orders WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, orderID, fixture.tenantA, fixture.restaurantA1).Scan(&status); err != nil {
		t.Fatalf("read final order status: %v", err)
	}
	if status != "confirmed" {
		t.Fatalf("final order status = %q, want confirmed", status)
	}
	assertOrderHistoryCount(t, db, orderID, 1)
}

func assertOrderHistoryCount(t *testing.T, db *sql.DB, orderID int64, want int) {
	t.Helper()
	var got int
	if err := db.QueryRow(`SELECT COUNT(*) FROM order_status_history WHERE order_id = $1`, orderID).Scan(&got); err != nil {
		t.Fatalf("count order status history: %v", err)
	}
	if got != want {
		t.Fatalf("order history count = %d, want %d", got, want)
	}
}

func TestOrderItemRepositoryPostgresCreateRejectsForeignOrderAndProductReferences(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	categoryA1 := insertScopeCategory(t, db, fixture.tenantA, fixture.restaurantA1, "order-item-a1")
	categoryB1 := insertScopeCategory(t, db, fixture.tenantB, fixture.restaurantB1, "order-item-b1")
	productA1 := insertScopeProduct(t, db, fixture.tenantA, fixture.restaurantA1, categoryA1, "order-item-a1")
	productB1 := insertScopeProduct(t, db, fixture.tenantB, fixture.restaurantB1, categoryB1, "order-item-b1")
	orderA1 := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA1, "order-item-a1")
	orderB1 := insertScopeOrder(t, db, fixture.tenantB, fixture.restaurantB1, "order-item-b1")
	repo := NewOrderRepository(db)

	for _, tc := range []struct {
		name      string
		orderID   int64
		productID int
	}{
		{name: "foreign order", orderID: orderB1, productID: productA1},
		{name: "foreign product", orderID: orderA1, productID: productB1},
	} {
		_, err := repo.CreateOrderItem(int64(fixture.tenantA), int64(fixture.restaurantA1), &domain.OrderItem{
			OrderID: tc.orderID, ProductID: int64(tc.productID), ProductName: "Synthetic item",
			Quantity: 1, UnitPrice: 10, TotalPrice: 10,
		})
		if err == nil {
			t.Errorf("order item with %s reference unexpectedly succeeded", tc.name)
		}
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM order_items WHERE tenant_id = $1 AND restaurant_id = $2`, fixture.tenantA, fixture.restaurantA1).Scan(&count); err != nil {
		t.Fatalf("verify denied order item side effects: %v", err)
	}
	if count != 0 {
		t.Fatalf("denied foreign order item creates left %d rows, want none", count)
	}
}

func TestOrderRepositoryPostgresHistoryAndItemsRemainInScope(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	categoryB1 := insertScopeCategory(t, db, fixture.tenantB, fixture.restaurantB1, "order-history-b1")
	productB1 := insertScopeProduct(t, db, fixture.tenantB, fixture.restaurantB1, categoryB1, "order-history-b1")
	orderA1 := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA1, "order-history-a1")
	orderB1 := insertScopeOrder(t, db, fixture.tenantB, fixture.restaurantB1, "order-history-b1")
	if _, err := db.Exec(`INSERT INTO order_status_history (order_id, old_status, new_status) VALUES ($1, 'pending', 'confirmed')`, orderB1); err != nil {
		t.Fatalf("insert synthetic order history: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO order_items (order_id, tenant_id, restaurant_id, product_id, product_name, quantity, unit_price, total_price) VALUES ($1, $2, $3, $4, 'Foreign item', 1, 10, 10)`, orderB1, fixture.tenantB, fixture.restaurantB1, productB1); err != nil {
		t.Fatalf("insert synthetic order item: %v", err)
	}
	repo := NewOrderRepository(db)

	items, err := repo.GetOrderItems(int64(fixture.tenantA), int64(fixture.restaurantA1), orderB1)
	if err != nil || len(items) != 0 {
		t.Fatalf("foreign order items returned %#v, err=%v", items, err)
	}
	history, err := repo.GetOrderStatusHistory(int64(fixture.tenantA), int64(fixture.restaurantA1), orderB1)
	if err != nil || len(history) != 0 {
		t.Fatalf("foreign order history returned %#v, err=%v", history, err)
	}
	if _, err := repo.GetOrderByID(int64(fixture.tenantA), int64(fixture.restaurantA1), orderA1); err != nil {
		t.Fatalf("in-scope order read failed: %v", err)
	}
}

func assertScopeOrderState(t *testing.T, db *sql.DB, id int64, tenantID, restaurantID int, status, paymentStatus string) {
	t.Helper()
	var gotStatus, gotPayment string
	err := db.QueryRow(`SELECT status, payment_status FROM orders WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, id, tenantID, restaurantID).Scan(&gotStatus, &gotPayment)
	if err != nil {
		t.Fatalf("read order state after denied operation: %v", err)
	}
	if gotStatus != status || gotPayment != paymentStatus {
		t.Fatalf("denied order operation changed row: status=%q payment=%q", gotStatus, gotPayment)
	}
}
