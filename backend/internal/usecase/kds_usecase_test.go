package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

type capturedKDSBroadcast struct {
	messageTypes []string
}

func (c *capturedKDSBroadcast) BroadcastKDSUpdate(
	_ context.Context,
	_ int64,
	_ int64,
	messageType string,
	_ interface{},
) error {
	c.messageTypes = append(c.messageTypes, messageType)
	return nil
}

func newKDSUseCaseWithMock(t *testing.T) (*KDSUseCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}

	repo := repository.NewKDSRepository(db)
	uc := NewKDSUseCase(repo)

	cleanup := func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sql expectations: %v", err)
		}
		_ = db.Close()
	}

	return uc, mock, cleanup
}

func TestKDSUseCaseGetQueueAppliesDefaultStatuses(t *testing.T) {
	uc, mock, cleanup := newKDSUseCaseWithMock(t)
	defer cleanup()

	receivedAt := time.Now().UTC().Add(-3 * time.Minute)
	completedAt := receivedAt.Add(2 * time.Minute)

	rows := sqlmock.NewRows([]string{
		"id",
		"tenant_id",
		"restaurant_id",
		"order_id",
		"order_number",
		"order_source",
		"customer_name",
		"status",
		"priority",
		"station_id",
		"items",
		"notes",
		"received_at",
		"started_at",
		"completed_at",
		"bumped_by",
		"created_at",
		"updated_at",
		"table_id",
		"table_number",
		"table_name",
	}).
		AddRow(
			int64(1), int64(10), int64(20), int64(101), "ORD-1", "website", "Alice", "new", 0, "expo",
			[]byte("[]"), "", receivedAt, nil, nil, nil, receivedAt, receivedAt, nil, nil, "",
		).
		AddRow(
			int64(2), int64(10), int64(20), int64(102), "ORD-2", "website", "Bob", "completed", 1, "expo",
			[]byte("[]"), "", receivedAt, receivedAt, completedAt, int64(9), receivedAt, completedAt, int64(44), int32(7), "Patio",
		)

	mock.ExpectQuery("(?s)FROM kds_queue\\s+k.*WHERE k\\.tenant_id = \\$1 AND k\\.restaurant_id = \\$2\\s+AND k\\.status IN \\(\\$3,\\$4,\\$5,\\$6\\)\\s+ORDER BY k\\.priority DESC, k\\.received_at ASC").
		WithArgs(int64(10), int64(20), "new", "in_progress", "ready", "recalled").
		WillReturnRows(rows)

	orders, err := uc.GetQueue(context.Background(), 10, 20, domain.KDSFilters{})
	if err != nil {
		t.Fatalf("GetQueue returned error: %v", err)
	}

	if len(orders) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(orders))
	}
	if orders[0].ElapsedSeconds <= 0 {
		t.Fatalf("expected positive elapsed seconds for active order, got %d", orders[0].ElapsedSeconds)
	}
	if orders[1].ElapsedSeconds != int64(completedAt.Sub(receivedAt).Seconds()) {
		t.Fatalf("unexpected completed elapsed seconds: got %d", orders[1].ElapsedSeconds)
	}
}

func TestKDSUseCaseGetQueueRejectsInvalidStatus(t *testing.T) {
	uc := NewKDSUseCase(repository.NewKDSRepository(nil))

	_, err := uc.GetQueue(context.Background(), 1, 2, domain.KDSFilters{
		Statuses: []domain.KDSStatus{"not_valid"},
	})
	if err == nil {
		t.Fatal("expected invalid status error")
	}
	if !strings.Contains(err.Error(), "invalid kds status") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKDSUseCaseBumpOrderBroadcastsUpdates(t *testing.T) {
	uc, mock, cleanup := newKDSUseCaseWithMock(t)
	defer cleanup()

	receivedAt := time.Now().UTC().Add(-4 * time.Minute)
	startedAt := receivedAt.Add(30 * time.Second)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT status, order_id\\s+FROM kds_queue").
		WithArgs(int64(99), int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"status", "order_id"}).AddRow("new", int64(501)))
	mock.ExpectExec("UPDATE kds_queue\\s+SET").
		WithArgs("in_progress", int64(9), int64(99), int64(1), int64(2), "new").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE orders\\s+SET status = \\$1").
		WithArgs("preparing", int64(501), int64(1), int64(2), "confirmed").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO order_status_history").
		WithArgs(int64(501), "confirmed", "preparing", int64(9), "Synced from KDS bump").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("(?s)FROM kds_queue\\s+k.*WHERE k\\.id = \\$1 AND k\\.tenant_id = \\$2 AND k\\.restaurant_id = \\$3").
		WithArgs(int64(99), int64(1), int64(2)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"tenant_id",
				"restaurant_id",
				"order_id",
				"order_number",
				"order_source",
				"customer_name",
				"status",
				"priority",
				"station_id",
				"items",
				"notes",
				"received_at",
				"started_at",
				"completed_at",
				"bumped_by",
				"created_at",
				"updated_at",
				"table_id",
				"table_number",
				"table_name",
			}).AddRow(
				int64(99), int64(1), int64(2), int64(501), "ORD-501", "website", "Alice", "in_progress", 0, "expo",
				[]byte("[]"), "", receivedAt, startedAt, nil, int64(9), receivedAt, startedAt, int64(77), int32(12), "Main Hall",
			),
		)

	captured := &capturedKDSBroadcast{}
	uc.SetBroadcaster(captured)

	order, err := uc.BumpOrder(context.Background(), 1, 2, 99, 9)
	if err != nil {
		t.Fatalf("BumpOrder returned error: %v", err)
	}

	if order.Status != domain.KDSStatusInProgress {
		t.Fatalf("expected in_progress, got %s", order.Status)
	}
	if len(captured.messageTypes) != 2 {
		t.Fatalf("expected 2 broadcasts, got %d", len(captured.messageTypes))
	}
	if captured.messageTypes[0] != "order_bumped" || captured.messageTypes[1] != "order_updated" {
		t.Fatalf("unexpected broadcast sequence: %+v", captured.messageTypes)
	}
}

func TestKDSUseCaseRecallOrderBroadcastsUpdates(t *testing.T) {
	uc, mock, cleanup := newKDSUseCaseWithMock(t)
	defer cleanup()

	receivedAt := time.Now().UTC().Add(-7 * time.Minute)
	startedAt := receivedAt.Add(1 * time.Minute)

	mock.ExpectExec("UPDATE kds_queue\\s+SET\\s+status = 'ready'").
		WithArgs(int64(88), int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery("(?s)FROM kds_queue\\s+k.*WHERE k\\.id = \\$1 AND k\\.tenant_id = \\$2 AND k\\.restaurant_id = \\$3").
		WithArgs(int64(88), int64(1), int64(2)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"tenant_id",
				"restaurant_id",
				"order_id",
				"order_number",
				"order_source",
				"customer_name",
				"status",
				"priority",
				"station_id",
				"items",
				"notes",
				"received_at",
				"started_at",
				"completed_at",
				"bumped_by",
				"created_at",
				"updated_at",
				"table_id",
				"table_number",
				"table_name",
			}).AddRow(
				int64(88), int64(1), int64(2), int64(500), "ORD-500", "website", "Bob", "ready", 0, "expo",
				[]byte("[]"), "", receivedAt, startedAt, nil, int64(7), receivedAt, startedAt, nil, nil, "",
			),
		)

	captured := &capturedKDSBroadcast{}
	uc.SetBroadcaster(captured)

	order, err := uc.RecallOrder(context.Background(), 1, 2, 88)
	if err != nil {
		t.Fatalf("RecallOrder returned error: %v", err)
	}

	if order.Status != domain.KDSStatusReady {
		t.Fatalf("expected ready status, got %s", order.Status)
	}
	if len(captured.messageTypes) != 2 {
		t.Fatalf("expected 2 broadcasts, got %d", len(captured.messageTypes))
	}
	if captured.messageTypes[0] != "order_recalled" || captured.messageTypes[1] != "order_updated" {
		t.Fatalf("unexpected broadcast sequence: %+v", captured.messageTypes)
	}
}

func TestKDSUseCaseUpdateSettingsValidationAndPersist(t *testing.T) {
	uc, mock, cleanup := newKDSUseCaseWithMock(t)
	defer cleanup()

	if err := uc.UpdateSettings(context.Background(), nil); err == nil {
		t.Fatal("expected nil settings validation error")
	}

	invalidDensity := &domain.KDSSettings{
		TenantID:            1,
		RestaurantID:        2,
		AutoBumpMinutes:     5,
		DisplayDensity:      "dense",
		ColorWarningMinutes: 5,
		ColorDangerMinutes:  10,
	}
	if err := uc.UpdateSettings(context.Background(), invalidDensity); err == nil {
		t.Fatal("expected invalid display density error")
	}

	invalidThresholds := &domain.KDSSettings{
		TenantID:            1,
		RestaurantID:        2,
		AutoBumpMinutes:     5,
		DisplayDensity:      "compact",
		ColorWarningMinutes: 10,
		ColorDangerMinutes:  10,
	}
	if err := uc.UpdateSettings(context.Background(), invalidThresholds); err == nil {
		t.Fatal("expected invalid threshold error")
	}

	valid := &domain.KDSSettings{
		TenantID:            1,
		RestaurantID:        2,
		AutoBumpMinutes:     15,
		AlertSoundEnabled:   true,
		ShowAnalyticsPanel:  true,
		DisplayDensity:      "spacious",
		ColorWarningMinutes: 8,
		ColorDangerMinutes:  16,
	}

	mock.ExpectExec("INSERT INTO kds_settings").
		WithArgs(int64(1), int64(2), 15, true, true, "spacious", 8, 16).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := uc.UpdateSettings(context.Background(), valid); err != nil {
		t.Fatalf("expected valid settings to persist, got %v", err)
	}
}
