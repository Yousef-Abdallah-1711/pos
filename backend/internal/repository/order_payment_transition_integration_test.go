package repository

import (
	"fmt"
	"sync"
	"testing"
)

func TestOrderPaymentStatusRejectsUnbackedRefundAndPaidDowngrade(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	orderID := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA1, "payment-guard")
	repo := NewOrderRepository(db)
	tenantID, restaurantID := int64(fixture.tenantA), int64(fixture.restaurantA1)

	if err := repo.UpdatePaymentStatus(tenantID, restaurantID, orderID, "paid"); err != nil {
		t.Fatalf("mark order paid: %v", err)
	}
	for _, requested := range []string{"refunded", "failed", "pending"} {
		if err := repo.UpdatePaymentStatus(tenantID, restaurantID, orderID, requested); err == nil {
			t.Errorf("transition paid → %s unexpectedly succeeded", requested)
		}
	}
	assertScopeOrderState(t, db, orderID, fixture.tenantA, fixture.restaurantA1, "pending", "paid")

	var count int
	var eventID string
	if err := db.QueryRow(`
		SELECT COUNT(*), MIN(event_id)
		FROM accounting_event_outbox
		WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`, tenantID, restaurantID, orderID,
	).Scan(&count, &eventID); err != nil {
		t.Fatalf("read paid-order intent: %v", err)
	}
	wantEventID := fmt.Sprintf("paid-order:%d:%d:%d", tenantID, restaurantID, orderID)
	if count != 1 || eventID != wantEventID {
		t.Fatalf("paid event count/id=%d/%q, want 1/%q", count, eventID, wantEventID)
	}
}

func TestCancelledOrderCannotBecomePaid(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	orderID := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA1, "cancel-payment-guard")
	repo := NewOrderRepository(db)
	tenantID, restaurantID := int64(fixture.tenantA), int64(fixture.restaurantA1)
	if err := repo.CancelOrder(tenantID, restaurantID, orderID, "customer cancelled"); err != nil {
		t.Fatalf("cancel unpaid order: %v", err)
	}
	if err := repo.UpdatePaymentStatus(tenantID, restaurantID, orderID, "paid"); err == nil {
		t.Fatal("cancelled order was marked paid")
	}
	assertScopeOrderState(t, db, orderID, fixture.tenantA, fixture.restaurantA1, "cancelled", "pending")
	var intents int
	if err := db.QueryRow(`SELECT COUNT(*) FROM accounting_event_outbox WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`, tenantID, restaurantID, orderID).Scan(&intents); err != nil {
		t.Fatal(err)
	}
	if intents != 0 {
		t.Fatalf("cancelled order has %d paid-order intents, want none", intents)
	}
}

func TestOrderPaymentTransitionRequiresTenantAndRestaurantScope(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	orderID := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA1, "payment-scope")
	repo := NewOrderRepository(db)

	for _, scope := range []struct {
		name         string
		tenantID     int64
		restaurantID int64
	}{
		{name: "foreign tenant", tenantID: int64(fixture.tenantB), restaurantID: int64(fixture.restaurantB1)},
		{name: "sibling restaurant", tenantID: int64(fixture.tenantA), restaurantID: int64(fixture.restaurantA2)},
	} {
		t.Run(scope.name, func(t *testing.T) {
			if err := repo.UpdatePaymentStatus(scope.tenantID, scope.restaurantID, orderID, "paid"); err == nil {
				t.Fatal("foreign scope updated payment status")
			}
		})
	}

	assertScopeOrderState(t, db, orderID, fixture.tenantA, fixture.restaurantA1, "pending", "pending")
	var intents int
	if err := db.QueryRow(`SELECT COUNT(*) FROM accounting_event_outbox WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`, fixture.tenantA, fixture.restaurantA1, orderID).Scan(&intents); err != nil {
		t.Fatal(err)
	}
	if intents != 0 {
		t.Fatalf("foreign transition created %d paid-order intents", intents)
	}
}

func TestConcurrentPaidTransitionAndCancellationHaveSingleSafeWinner(t *testing.T) {
	db, fixture := openProductScopeFixture(t)
	repo := NewOrderRepository(db)
	tenantID, restaurantID := int64(fixture.tenantA), int64(fixture.restaurantA1)

	for attempt := range 8 {
		orderID := insertScopeOrder(t, db, fixture.tenantA, fixture.restaurantA1, fmt.Sprintf("payment-race-%d", attempt))
		start := make(chan struct{})
		results := make(chan error, 2)
		var group sync.WaitGroup
		group.Add(2)
		go func() {
			defer group.Done()
			<-start
			results <- repo.UpdatePaymentStatus(tenantID, restaurantID, orderID, "paid")
		}()
		go func() {
			defer group.Done()
			<-start
			results <- repo.CancelOrder(tenantID, restaurantID, orderID, "customer cancellation")
		}()
		close(start)
		group.Wait()
		close(results)

		successes := 0
		for err := range results {
			if err == nil {
				successes++
			}
		}
		if successes != 1 {
			t.Fatalf("race %d successful transitions=%d, want exactly one", attempt, successes)
		}

		var orderStatus, paymentStatus string
		if err := db.QueryRow(`SELECT status, payment_status FROM orders WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, orderID, tenantID, restaurantID).Scan(&orderStatus, &paymentStatus); err != nil {
			t.Fatal(err)
		}
		var intents int
		if err := db.QueryRow(`SELECT COUNT(*) FROM accounting_event_outbox WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`, tenantID, restaurantID, orderID).Scan(&intents); err != nil {
			t.Fatal(err)
		}
		switch {
		case orderStatus == "cancelled" && paymentStatus == "pending" && intents == 0:
		case orderStatus == "pending" && paymentStatus == "paid" && intents == 1:
		default:
			t.Fatalf("race %d ended in order/payment/intents=%s/%s/%d", attempt, orderStatus, paymentStatus, intents)
		}
	}
}
