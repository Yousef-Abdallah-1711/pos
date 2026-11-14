package http

import (
	"bytes"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

func TestOrderPaymentStatusPersistsDurableAccountingIntent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	defer db.Close()

	orderUC := usecase.NewOrderUseCase(repository.NewOrderRepository(db), repository.NewProductRepository(db))
	autoPosting := &recordingAutoPostingService{}
	orderUC.SetAutoPostingEngine(autoPosting)
	handler := NewAdminOrderHandler(orderUC, nil, nil)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT status, payment_status`).
		WithArgs(int64(1), int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"status", "payment_status"}).AddRow("pending", "pending"))
	mock.ExpectExec(`UPDATE orders`).WithArgs("paid", int64(1), int64(1), int64(2)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO accounting_event_outbox`).WithArgs("paid-order:1:2:1", int64(1), int64(2), int64(1), 77).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := bytes.NewBufferString(`{"payment_status":"paid"}`)
	req := httptest.NewRequest(stdhttp.MethodPut, "/api/v1/admin/orders/1/payment-status", body)
	req.SetPathValue("id", "1")
	req = withAccountingAuthContext(req, 1, 2, 77)
	rec := httptest.NewRecorder()
	handler.UpdatePaymentStatus(rec, req)
	if rec.Code != stdhttp.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}

	var response map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response["accounting_intent"] != "persisted" {
		t.Fatalf("accounting_intent=%v, want persisted", response["accounting_intent"])
	}
	if len(autoPosting.salesOrders) != 0 || len(autoPosting.cogsOrders) != 0 {
		t.Fatal("handler must leave posting to durable worker")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
