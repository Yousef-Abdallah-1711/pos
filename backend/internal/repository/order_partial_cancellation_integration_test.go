package repository

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/google/uuid"
	"pos-saas/internal/domain"
)

func cancellationFixture(t *testing.T, quantity int) (stockFixture, domain.OrderItemCancellation) {
	t.Helper()
	f := newStockFixture(t, 20)
	order := f.order(t, "partial-"+uuid.NewString()[:8], quantity)
	var item, actor int64
	if err := f.db.QueryRow(`SELECT id FROM order_items WHERE order_id=$1 AND tenant_id=$2 AND restaurant_id=$3`, order, f.tenant, f.restaurant).Scan(&item); err != nil {
		t.Fatal(err)
	}
	if err := f.db.QueryRow(`INSERT INTO users(tenant_id,restaurant_id,email,name,password_hash,role,status) VALUES($1,$2,$3,'Synthetic actor','!test-only','owner','active') RETURNING id`, f.tenant, f.restaurant, uuid.NewString()+"@example.invalid").Scan(&actor); err != nil {
		t.Fatal(err)
	}
	if _, err := f.db.Exec(`UPDATE orders SET subtotal=$1,tax_amount=$2,total_amount=$3 WHERE id=$4 AND tenant_id=$5 AND restaurant_id=$6`, quantity*10, quantity, quantity*11, order, f.tenant, f.restaurant); err != nil {
		t.Fatal(err)
	}
	return f, domain.OrderItemCancellation{TenantID: int64(f.tenant), RestaurantID: int64(f.restaurant), OrderID: order, ItemID: item, Quantity: 2, ActorID: actor, RequestKey: uuid.NewString(), Reason: "Synthetic partial cancellation"}
}

func TestPartialCancellationReplayAndWholeCancellationReleaseOnlyRemainder(t *testing.T) {
	f, request := cancellationFixture(t, 5)
	if err := f.transition(request.OrderID, "confirmed"); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 15)
	for _, replay := range []bool{false, true} {
		raw, err := f.repo.CancelOrderItem(context.Background(), request)
		if err != nil {
			t.Fatal(err)
		}
		var body struct {
			Replayed bool `json:"replayed"`
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatal(err)
		}
		if body.Replayed != replay {
			t.Fatalf("replayed=%v want %v", body.Replayed, replay)
		}
		f.assertStock(t, 17)
	}
	items, err := f.repo.GetOrderItems(request.TenantID, request.RestaurantID, request.OrderID)
	if err != nil || len(items) != 1 || items[0].Quantity != 3 || items[0].TotalPrice != 30 {
		t.Fatalf("effective items=%+v err=%v", items, err)
	}
	if err := f.transition(request.OrderID, "cancelled"); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 20)
}

func TestPendingPartialCancellationReservesOnlyRemainingQuantity(t *testing.T) {
	f, request := cancellationFixture(t, 5)
	if _, err := f.repo.CancelOrderItem(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 20)
	if err := f.transition(request.OrderID, "confirmed"); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 17)
	var total float64
	if err := f.db.QueryRow(`SELECT total_amount FROM orders WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, request.OrderID, request.TenantID, request.RestaurantID).Scan(&total); err != nil {
		t.Fatal(err)
	}
	if total != 33 {
		t.Fatalf("remaining order total=%v want 33", total)
	}
}

func TestPartialCancellationRejectsForeignScopeExcessAndReusedPayload(t *testing.T) {
	f, request := cancellationFixture(t, 5)
	if err := f.transition(request.OrderID, "confirmed"); err != nil {
		t.Fatal(err)
	}
	invalid := request
	invalid.TenantID++
	if _, err := f.repo.CancelOrderItem(context.Background(), invalid); err == nil {
		t.Fatal("foreign cancellation accepted")
	}
	invalid = request
	invalid.Quantity = 6
	if _, err := f.repo.CancelOrderItem(context.Background(), invalid); err == nil {
		t.Fatal("excess cancellation accepted")
	}
	f.assertStock(t, 15)
	if _, err := f.repo.CancelOrderItem(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	invalid = request
	invalid.Quantity = 1
	if _, err := f.repo.CancelOrderItem(context.Background(), invalid); err == nil {
		t.Fatal("same key accepted different payload")
	}
	f.assertStock(t, 17)
}

func TestPartialCancellationCannotRestockPaidOrFulfilledOrders(t *testing.T) {
	for _, state := range []string{"paid", "preparing"} {
		t.Run(state, func(t *testing.T) {
			f, request := cancellationFixture(t, 5)
			if err := f.transition(request.OrderID, "confirmed"); err != nil {
				t.Fatal(err)
			}
			if state == "paid" {
				if err := f.repo.UpdatePaymentStatus(request.TenantID, request.RestaurantID, request.OrderID, "paid"); err != nil {
					t.Fatal(err)
				}
			} else if err := f.transition(request.OrderID, state); err != nil {
				t.Fatal(err)
			}
			if _, err := f.repo.CancelOrderItem(context.Background(), request); err == nil {
				t.Fatal("paid/fulfilled cancellation accepted")
			}
			f.assertStock(t, 15)
		})
	}
}

func TestPartialCancellationConcurrentReplayDoesNotDoubleRelease(t *testing.T) {
	f, request := cancellationFixture(t, 5)
	if err := f.transition(request.OrderID, "confirmed"); err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for i := 0; i < 2; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			_, err := f.repo.CancelOrderItem(context.Background(), request)
			results <- err
		}()
	}
	close(start)
	workers.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatal(err)
		}
	}
	f.assertStock(t, 17)
	var receipts int
	if err := f.db.QueryRow(`SELECT count(*) FROM order_item_cancellations WHERE tenant_id=$1 AND restaurant_id=$2 AND order_id=$3`, request.TenantID, request.RestaurantID, request.OrderID).Scan(&receipts); err != nil {
		t.Fatal(err)
	}
	if receipts != 1 {
		t.Fatalf("receipts=%d want one", receipts)
	}
}

func TestPartialCancellationAllUnitsClosesOrderWithoutDoubleRelease(t *testing.T) {
	f, request := cancellationFixture(t, 2)
	if err := f.transition(request.OrderID, "confirmed"); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 18)
	if _, err := f.repo.CancelOrderItem(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 20)
	assertScopeOrderState(t, f.db, request.OrderID, f.tenant, f.restaurant, "cancelled", "pending")
	if _, err := f.repo.CancelOrderItem(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	f.assertStock(t, 20)
}
