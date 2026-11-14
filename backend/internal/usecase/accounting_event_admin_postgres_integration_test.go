//go:build integration

package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"pos-saas/internal/repository"
)

type failedEventFixture struct {
	db                  *sql.DB
	events              *repository.AccountingEventRepository
	scope, foreign      repository.AccountingEventScope
	order, foreignOrder int64
}

func newFailedEventFixture(t *testing.T) failedEventFixture {
	t.Helper()
	db := openAccountingEventTestDB(t)
	tenant, restaurant, order := createPaidOrderEventFixture(t, db)
	otherTenant, otherRestaurant, otherOrder := createPaidOrderEventFixture(t, db)
	fixture := failedEventFixture{db, repository.NewAccountingEventRepository(db),
		repository.AccountingEventScope{TenantID: tenant, RestaurantID: restaurant},
		repository.AccountingEventScope{TenantID: otherTenant, RestaurantID: otherRestaurant}, order, otherOrder}
	for _, item := range []struct {
		scope repository.AccountingEventScope
		order int64
	}{{fixture.scope, order}, {fixture.foreign, otherOrder}} {
		if err := repository.NewOrderRepository(db).UpdatePaymentStatus(item.scope.TenantID, item.scope.RestaurantID, item.order, "paid"); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`UPDATE accounting_event_outbox SET status='failed',attempts=max_attempts,last_error=repeat('x',3000)
   WHERE tenant_id=$1 AND restaurant_id=$2 AND order_id=$3`, item.scope.TenantID, item.scope.RestaurantID, item.order); err != nil {
			t.Fatal(err)
		}
	}
	return fixture
}

func failedEventInScope(t *testing.T, fixture failedEventFixture) repository.FailedAccountingEvent {
	t.Helper()
	visible, total, err := fixture.events.ListFailed(context.Background(), fixture.scope, 25, 0)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(visible) != 1 || visible[0].OrderID != fixture.order || visible[0].RedriveCount != 0 {
		t.Fatalf("wrong scoped events: total=%d items=%+v", total, visible)
	}
	if len(visible[0].LastError) != 2000 {
		t.Fatalf("error was not bounded: %d", len(visible[0].LastError))
	}
	return visible[0]
}

func redriveRequest(scope repository.AccountingEventScope, eventID int64) repository.AccountingEventRedriveRequest {
	return repository.AccountingEventRedriveRequest{Scope: scope, EventID: eventID, ActorUserID: 91, Reason: "corrected the missing account mapping"}
}

func TestFailedAccountingEventVisibilityAndAuditAreScoped(t *testing.T) {
	fixture := newFailedEventFixture(t)
	event := failedEventInScope(t, fixture)
	request := redriveRequest(fixture.scope, event.ID)
	eventID, err := fixture.events.RedriveFailed(context.Background(), request)
	if err != nil || eventID != event.EventID {
		t.Fatalf("redrive=%q err=%v", eventID, err)
	}
	assertAccountingEventState(t, fixture.db, fixture.scope.TenantID, fixture.scope.RestaurantID, fixture.order, repository.AccountingEventRetryable, 0, 1)
	var actor, attempts int64
	var reason, lastError string
	err = fixture.db.QueryRow(`SELECT a.actor_user_id,a.previous_attempts,a.previous_error,a.reason
  FROM accounting_event_redrive_audit a JOIN accounting_event_outbox e ON e.id=a.outbox_id
  WHERE a.outbox_id=$1 AND e.tenant_id=$2 AND e.restaurant_id=$3`, event.ID, fixture.scope.TenantID, fixture.scope.RestaurantID).
		Scan(&actor, &attempts, &lastError, &reason)
	if err != nil || actor != 91 || attempts != 8 || len(lastError) != 2000 || reason != request.Reason {
		t.Fatalf("audit not preserved: actor=%d attempts=%d error length=%d reason=%q err=%v", actor, attempts, len(lastError), reason, err)
	}
	if _, err := fixture.events.RedriveFailed(context.Background(), request); !errors.Is(err, repository.ErrAccountingEventNotRedrivable) {
		t.Fatalf("duplicate redrive: %v", err)
	}
}

func TestAccountingEventRedriveRejectsForeignOrNoLongerPaidOrders(t *testing.T) {
	fixture := newFailedEventFixture(t)
	event := failedEventInScope(t, fixture)
	for _, scope := range []repository.AccountingEventScope{fixture.foreign, {TenantID: fixture.scope.TenantID, RestaurantID: fixture.foreign.RestaurantID}} {
		if _, err := fixture.events.RedriveFailed(context.Background(), redriveRequest(scope, event.ID)); !errors.Is(err, repository.ErrAccountingEventNotFound) {
			t.Fatalf("cross-scope redrive: %v", err)
		}
		rows, total, err := fixture.events.ListFailed(context.Background(), scope, 25, 0)
		if scope.TenantID == fixture.scope.TenantID && (err != nil || total != 0 || len(rows) != 0) {
			t.Fatalf("sibling scope leaked failed events: total=%d err=%v", total, err)
		}
	}
	request := redriveRequest(fixture.foreign, event.ID)
	// Simulate a provider-side state change in this disposable fixture. The public
	// payment API deliberately rejects paid-to-refunded transitions without policy.
	_, err := fixture.db.Exec(`UPDATE orders SET payment_status='refunded'
  WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3 AND payment_status='paid'`, fixture.foreignOrder, fixture.foreign.TenantID, fixture.foreign.RestaurantID)
	if err != nil {
		t.Fatal(err)
	}
	var foreignEvent int64
	err = fixture.db.QueryRow(`SELECT id FROM accounting_event_outbox WHERE tenant_id=$1 AND restaurant_id=$2 AND order_id=$3`, fixture.foreign.TenantID, fixture.foreign.RestaurantID, fixture.foreignOrder).Scan(&foreignEvent)
	if err != nil {
		t.Fatal(err)
	}
	request.EventID = foreignEvent
	if _, err := fixture.events.RedriveFailed(context.Background(), request); !errors.Is(err, repository.ErrAccountingEventNotRedrivable) {
		t.Fatalf("non-paid redrive: %v", err)
	}
}

func TestAccountingEventAuditRuntimePrivilegesAreAppendOnly(t *testing.T) {
	db := openAccountingEventTestDB(t)
	var readable, insertable, sequenceUsable, updatable, deletable, truncatable bool
	err := db.QueryRow(`SELECT
  has_table_privilege('pos_app','accounting_event_redrive_audit','SELECT'),
  has_table_privilege('pos_app','accounting_event_redrive_audit','INSERT'),
  has_sequence_privilege('pos_app','accounting_event_redrive_audit_id_seq','USAGE'),
  has_table_privilege('pos_app','accounting_event_redrive_audit','UPDATE'),
  has_table_privilege('pos_app','accounting_event_redrive_audit','DELETE'),
  has_table_privilege('pos_app','accounting_event_redrive_audit','TRUNCATE')`).
		Scan(&readable, &insertable, &sequenceUsable, &updatable, &deletable, &truncatable)
	if err != nil {
		t.Fatal(err)
	}
	if !readable || !insertable || !sequenceUsable || updatable || deletable || truncatable {
		t.Fatalf("runtime audit privileges read=%t insert=%t seq=%t update=%t delete=%t truncate=%t", readable, insertable, sequenceUsable, updatable, deletable, truncatable)
	}
}
