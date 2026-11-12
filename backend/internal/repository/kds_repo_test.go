package repository

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"pos-saas/internal/domain"
)

var kdsQueueColumns = []string{
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
}

func newKDSRepositoryWithMock(t *testing.T) (*KDSRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}

	cleanup := func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sql expectations: %v", err)
		}
		_ = db.Close()
	}

	return NewKDSRepository(db), mock, cleanup
}

func TestNextBumpStatusTransitions(t *testing.T) {
	t.Parallel()

	cases := []struct {
		current string
		want    string
		wantErr bool
	}{
		{current: "new", want: "in_progress"},
		{current: "in_progress", want: "ready"},
		{current: "ready", want: "completed"},
		{current: "recalled", want: "ready"},
		{current: "completed", wantErr: true},
		{current: "unknown", wantErr: true},
	}

	for _, tc := range cases {
		got, err := nextBumpStatus(tc.current)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("expected error for status %q", tc.current)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for status %q: %v", tc.current, err)
		}
		if got != tc.want {
			t.Fatalf("nextBumpStatus(%q)=%q, want %q", tc.current, got, tc.want)
		}
	}
}

func TestScanKDSOrderSuccess(t *testing.T) {
	repo, mock, cleanup := newKDSRepositoryWithMock(t)
	defer cleanup()

	now := time.Now().UTC().Add(-5 * time.Minute)
	started := now.Add(1 * time.Minute)
	completed := now.Add(4 * time.Minute)
	itemsJSON := `[{"product_id":5,"product_name":"Burger","quantity":2}]`

	rows := sqlmock.NewRows(kdsQueueColumns).AddRow(
		int64(1),
		int64(10),
		int64(20),
		int64(30),
		"ORD-1001",
		"pos",
		"Alice",
		"completed",
		1,
		"expo",
		[]byte(itemsJSON),
		"extra spicy",
		now,
		started,
		completed,
		int64(7),
		now,
		now,
		int64(4),
		int32(12),
		"Hall 12",
	)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	row := repo.db.QueryRowContext(context.Background(), "SELECT 1")
	order, err := scanKDSOrder(row)
	if err != nil {
		t.Fatalf("scanKDSOrder returned error: %v", err)
	}

	if order.ID != 1 || order.OrderNumber != "ORD-1001" {
		t.Fatalf("unexpected mapped order payload: %+v", order)
	}
	if order.Status != domain.KDSStatusCompleted {
		t.Fatalf("expected completed status, got %s", order.Status)
	}
	if order.ElapsedSeconds <= 0 {
		t.Fatalf("expected elapsed seconds > 0, got %d", order.ElapsedSeconds)
	}
	if order.BumpedBy == nil || *order.BumpedBy != 7 {
		t.Fatalf("expected bumped_by=7, got %+v", order.BumpedBy)
	}
	if len(order.Items) != 1 || order.Items[0].ProductName != "Burger" {
		t.Fatalf("unexpected items payload: %+v", order.Items)
	}
}

func TestScanKDSOrderInvalidItemsPayload(t *testing.T) {
	repo, mock, cleanup := newKDSRepositoryWithMock(t)
	defer cleanup()

	now := time.Now().UTC().Add(-1 * time.Minute)

	rows := sqlmock.NewRows(kdsQueueColumns).AddRow(
		int64(1),
		int64(10),
		int64(20),
		int64(30),
		"ORD-1001",
		"website",
		"Alice",
		"new",
		0,
		"expo",
		[]byte("{invalid-json"),
		"",
		now,
		nil,
		nil,
		nil,
		now,
		now,
		nil,
		nil,
		nil,
	)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	row := repo.db.QueryRowContext(context.Background(), "SELECT 1")
	_, err := scanKDSOrder(row)
	if err == nil {
		t.Fatal("expected error for invalid items payload")
	}
	if !strings.Contains(err.Error(), "unmarshal") {
		t.Fatalf("expected unmarshal error, got %v", err)
	}
}

func TestGetQueueRejectsInvalidStatusFilter(t *testing.T) {
	repo, _, cleanup := newKDSRepositoryWithMock(t)
	defer cleanup()

	_, err := repo.GetQueue(context.Background(), 1, 2, []string{"new", "not_a_status"})
	if err == nil {
		t.Fatal("expected invalid status filter error")
	}
	if !strings.Contains(err.Error(), "invalid kds status filter") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetQueueAppliesStatusFilter(t *testing.T) {
	repo, mock, cleanup := newKDSRepositoryWithMock(t)
	defer cleanup()

	now := time.Now().UTC().Add(-2 * time.Minute)
	rows := sqlmock.NewRows(kdsQueueColumns).AddRow(
		int64(5),
		int64(10),
		int64(20),
		int64(50),
		"ORD-2001",
		"website",
		"Bob",
		"ready",
		0,
		"expo",
		[]byte("[]"),
		"",
		now,
		nil,
		nil,
		nil,
		now,
		now,
		nil,
		nil,
		nil,
	)

	mock.ExpectQuery("(?s)FROM kds_queue\\s+k.*WHERE k\\.tenant_id = \\$1 AND k\\.restaurant_id = \\$2\\s+AND k\\.status IN \\(\\$3,\\$4\\)\\s+ORDER BY k\\.priority DESC, k\\.received_at ASC").
		WithArgs(int64(10), int64(20), "new", "ready").
		WillReturnRows(rows)

	orders, err := repo.GetQueue(context.Background(), 10, 20, []string{"new", "ready"})
	if err != nil {
		t.Fatalf("GetQueue returned error: %v", err)
	}
	if len(orders) != 1 {
		t.Fatalf("expected one order, got %d", len(orders))
	}
	if orders[0].Status != domain.KDSStatusReady {
		t.Fatalf("expected ready status, got %s", orders[0].Status)
	}
}

func TestGetSettingsReturnsDefaultsWhenMissing(t *testing.T) {
	repo, mock, cleanup := newKDSRepositoryWithMock(t)
	defer cleanup()

	mock.ExpectQuery("FROM kds_settings").
		WithArgs(int64(10), int64(20)).
		WillReturnError(sql.ErrNoRows)

	settings, err := repo.GetSettings(context.Background(), 10, 20)
	if err != nil {
		t.Fatalf("GetSettings returned error: %v", err)
	}
	if settings.DisplayDensity != "comfortable" {
		t.Fatalf("expected default display density, got %s", settings.DisplayDensity)
	}
	if settings.ShowAnalyticsPanel {
		t.Fatal("expected show_analytics_panel default to false")
	}
	if settings.ColorWarningMinutes != 10 || settings.ColorDangerMinutes != 20 {
		t.Fatalf("unexpected default thresholds: %+v", settings)
	}
}

func TestUpdateSettingsUpsertsValues(t *testing.T) {
	repo, mock, cleanup := newKDSRepositoryWithMock(t)
	defer cleanup()

	settings := &domain.KDSSettings{
		TenantID:            10,
		RestaurantID:        20,
		AutoBumpMinutes:     15,
		AlertSoundEnabled:   false,
		ShowAnalyticsPanel:  true,
		DisplayDensity:      "compact",
		ColorWarningMinutes: 6,
		ColorDangerMinutes:  12,
	}

	mock.ExpectExec("INSERT INTO kds_settings").
		WithArgs(
			int64(10),
			int64(20),
			15,
			false,
			true,
			"compact",
			6,
			12,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.UpdateSettings(context.Background(), settings); err != nil {
		t.Fatalf("UpdateSettings returned error: %v", err)
	}
}

func TestGetAnalyticsReturnsOverviewTrendAndHours(t *testing.T) {
	repo, mock, cleanup := newKDSRepositoryWithMock(t)
	defer cleanup()

	mock.ExpectQuery("average_prep_minutes").
		WithArgs(int64(10), int64(20)).
		WillReturnRows(
			sqlmock.NewRows([]string{"average_prep_minutes", "total_completed_orders"}).
				AddRow(8.5, int64(42)),
		)

	mock.ExpectQuery("WITH days AS").
		WithArgs(int64(10), int64(20)).
		WillReturnRows(
			sqlmock.NewRows([]string{"date", "avg_prep_minutes", "completed_orders"}).
				AddRow("2026-02-18", 7.2, int64(5)).
				AddRow("2026-02-19", 9.1, int64(6)),
		)

	mock.ExpectQuery("WITH hours AS").
		WithArgs(int64(10), int64(20)).
		WillReturnRows(
			sqlmock.NewRows([]string{"hour", "order_count"}).
				AddRow(11, int64(9)).
				AddRow(12, int64(12)),
		)

	analytics, err := repo.GetAnalytics(context.Background(), 10, 20)
	if err != nil {
		t.Fatalf("GetAnalytics returned error: %v", err)
	}

	if analytics.AveragePrepMinutes != 8.5 {
		t.Fatalf("expected average prep 8.5, got %f", analytics.AveragePrepMinutes)
	}
	if analytics.TotalCompletedOrders != 42 {
		t.Fatalf("expected total completed 42, got %d", analytics.TotalCompletedOrders)
	}
	if len(analytics.PrepTimeTrend) != 2 {
		t.Fatalf("expected 2 trend points, got %d", len(analytics.PrepTimeTrend))
	}
	if len(analytics.BusiestHours) != 2 {
		t.Fatalf("expected 2 hour buckets, got %d", len(analytics.BusiestHours))
	}
}
