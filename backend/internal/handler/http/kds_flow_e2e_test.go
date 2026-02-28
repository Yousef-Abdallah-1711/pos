package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

func TestKDSFlowDisplayBumpAndStatusUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sql expectations: %v", err)
		}
		_ = db.Close()
	}()

	repo := repository.NewKDSRepository(db)
	uc := usecase.NewKDSUseCase(repo)
	handler := NewKDSHandler(uc)

	mux := http.NewServeMux()
	mux.Handle("GET /api/v1/kds/queue", withFixedKDSContext(http.HandlerFunc(handler.GetQueue)))
	mux.Handle("PUT /api/v1/kds/orders/{id}/bump", withFixedKDSContext(http.HandlerFunc(handler.BumpOrder)))

	receivedAt := time.Now().UTC().Add(-3 * time.Minute)
	startedAt := receivedAt.Add(20 * time.Second)

	// Step 1: KDS display fetches active queue (new order appears).
	mock.ExpectQuery("(?s)FROM kds_queue\\s+k.*WHERE k\\.tenant_id = \\$1 AND k\\.restaurant_id = \\$2\\s+AND k\\.status IN \\(\\$3,\\$4,\\$5,\\$6\\)\\s+ORDER BY k\\.priority DESC, k\\.received_at ASC").
		WithArgs(int64(1), int64(2), "new", "in_progress", "ready", "recalled").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"tenant_id",
				"restaurant_id",
				"order_id",
				"order_number",
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
				int64(501), int64(1), int64(2), int64(9001), "ORD-9001", "Kitchen Ticket", "new", 1, "expo",
				[]byte(`[{"product_name":"Burger","quantity":1}]`), "no onion", receivedAt, nil, nil, nil, receivedAt, receivedAt, int64(88), int32(8), "VIP",
			),
		)

	// Step 2: Kitchen bumps order (new -> in_progress), backend syncs order status.
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT status, order_id\\s+FROM kds_queue").
		WithArgs(int64(501), int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"status", "order_id"}).AddRow("new", int64(9001)))
	mock.ExpectExec("UPDATE kds_queue\\s+SET").
		WithArgs("in_progress", int64(7), int64(501), int64(1), int64(2), "new").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE orders\\s+SET status = \\$1").
		WithArgs("preparing", int64(9001), int64(1), int64(2), "confirmed").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO order_status_history").
		WithArgs(int64(9001), "confirmed", "preparing", int64(7), "Synced from KDS bump").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("(?s)FROM kds_queue\\s+k.*WHERE k\\.id = \\$1 AND k\\.tenant_id = \\$2 AND k\\.restaurant_id = \\$3").
		WithArgs(int64(501), int64(1), int64(2)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"tenant_id",
				"restaurant_id",
				"order_id",
				"order_number",
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
				int64(501), int64(1), int64(2), int64(9001), "ORD-9001", "Kitchen Ticket", "in_progress", 1, "expo",
				[]byte(`[{"product_name":"Burger","quantity":1}]`), "no onion", receivedAt, startedAt, nil, int64(7), receivedAt, startedAt, int64(88), int32(8), "VIP",
			),
		)

	// Step 3: KDS display refresh shows bumped order in in_progress.
	mock.ExpectQuery("(?s)FROM kds_queue\\s+k.*WHERE k\\.tenant_id = \\$1 AND k\\.restaurant_id = \\$2\\s+AND k\\.status IN \\(\\$3\\)\\s+ORDER BY k\\.priority DESC, k\\.received_at ASC").
		WithArgs(int64(1), int64(2), "in_progress").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"tenant_id",
				"restaurant_id",
				"order_id",
				"order_number",
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
				int64(501), int64(1), int64(2), int64(9001), "ORD-9001", "Kitchen Ticket", "in_progress", 1, "expo",
				[]byte(`[{"product_name":"Burger","quantity":1}]`), "no onion", receivedAt, startedAt, nil, int64(7), receivedAt, startedAt, int64(88), int32(8), "VIP",
			),
		)

	// Initial display check.
	firstQueueReq := httptest.NewRequest(http.MethodGet, "/api/v1/kds/queue", nil)
	firstQueueRec := httptest.NewRecorder()
	mux.ServeHTTP(firstQueueRec, firstQueueReq)
	if firstQueueRec.Code != http.StatusOK {
		t.Fatalf("queue step failed: %d body=%s", firstQueueRec.Code, firstQueueRec.Body.String())
	}

	var firstQueue struct {
		Success bool `json:"success"`
		Data    []struct {
			ID     int64  `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(firstQueueRec.Body).Decode(&firstQueue); err != nil {
		t.Fatalf("failed to decode first queue response: %v", err)
	}
	if !firstQueue.Success || len(firstQueue.Data) != 1 || firstQueue.Data[0].Status != "new" {
		t.Fatalf("unexpected first queue payload: %+v", firstQueue)
	}

	// Bump check.
	bumpReq := httptest.NewRequest(http.MethodPut, "/api/v1/kds/orders/501/bump", nil)
	bumpRec := httptest.NewRecorder()
	mux.ServeHTTP(bumpRec, bumpReq)
	if bumpRec.Code != http.StatusOK {
		t.Fatalf("bump step failed: %d body=%s", bumpRec.Code, bumpRec.Body.String())
	}

	var bumpPayload struct {
		Success bool `json:"success"`
		Data    struct {
			ID     int64  `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(bumpRec.Body).Decode(&bumpPayload); err != nil {
		t.Fatalf("failed to decode bump response: %v", err)
	}
	if !bumpPayload.Success || bumpPayload.Data.ID != 501 || bumpPayload.Data.Status != "in_progress" {
		t.Fatalf("unexpected bump payload: %+v", bumpPayload)
	}

	// Updated display check.
	secondQueueReq := httptest.NewRequest(http.MethodGet, "/api/v1/kds/queue?status=in_progress", nil)
	secondQueueRec := httptest.NewRecorder()
	mux.ServeHTTP(secondQueueRec, secondQueueReq)
	if secondQueueRec.Code != http.StatusOK {
		t.Fatalf("post-bump queue step failed: %d body=%s", secondQueueRec.Code, secondQueueRec.Body.String())
	}

	var secondQueue struct {
		Success bool `json:"success"`
		Data    []struct {
			ID     int64  `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(secondQueueRec.Body).Decode(&secondQueue); err != nil {
		t.Fatalf("failed to decode second queue response: %v", err)
	}
	if !secondQueue.Success || len(secondQueue.Data) != 1 || secondQueue.Data[0].Status != "in_progress" {
		t.Fatalf("unexpected second queue payload: %+v", secondQueue)
	}
}
