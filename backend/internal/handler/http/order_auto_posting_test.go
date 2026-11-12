package http

import (
	"bytes"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

func TestOrderAutoPostingIntegration_OrderPaidCreatesSalesAndCOGS(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderUC := usecase.NewOrderUseCase(orderRepo, productRepo)
	autoPosting := &recordingAutoPostingService{}
	orderUC.SetAutoPostingEngine(autoPosting)
	handler := NewAdminOrderHandler(orderUC, nil, nil)

	now := time.Date(2026, time.March, 1, 11, 0, 0, 0, time.UTC)
	orderID := int64(601)
	userID := 99

	mock.ExpectExec(`UPDATE orders`).
		WithArgs("paid", orderID, int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(`FROM orders o`).
		WithArgs(orderID, int64(1), int64(2)).
		WillReturnRows(
			sqlmock.NewRows(orderDetailColumnsForAccountingTests()).
				AddRow(
					orderID, int64(1), int64(2), "ORD-2026-000601", "Walk-in Guest", nil,
					"01000000000", nil, nil, nil,
					nil, nil, nil,
					nil, 100.0, 14.0, 0.0,
					0.0, 114.0, "cash", "paid",
					"confirmed", nil, nil, "", "pos",
					nil, nil, "", "", nil,
					now, now,
				),
		)

	mock.ExpectQuery(`FROM order_items`).
		WithArgs(int64(1), int64(2), orderID).
		WillReturnRows(
			sqlmock.NewRows(orderItemColumnsForAccountingTests()).
				AddRow(
					int64(1), orderID, int64(201), nil, "Meal Combo", nil,
					2, 50.0, 0.0, 100.0,
					"", []byte(`[]`), now, now,
				),
		)

	body, err := json.Marshal(map[string]string{"payment_status": "paid"})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(stdhttp.MethodPut, "/api/v1/admin/orders/601/payment-status", bytes.NewReader(body))
	req.SetPathValue("id", "601")
	req = withAccountingAuthContext(req, 1, 2, userID)

	rec := httptest.NewRecorder()
	handler.UpdatePaymentStatus(rec, req)

	if rec.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	assertJSONSuccess(t, rec)

	if len(autoPosting.salesOrders) != 1 || len(autoPosting.cogsOrders) != 1 {
		t.Fatalf("expected 1 sales + 1 COGS posting, got sales=%d cogs=%d", len(autoPosting.salesOrders), len(autoPosting.cogsOrders))
	}
	if autoPosting.salesOrders[0].ID != orderID || autoPosting.cogsOrders[0].ID != orderID {
		t.Fatalf("expected order id %d in both postings", orderID)
	}
	if autoPosting.salesPostedBy[0] != userID || autoPosting.cogsPostedBy[0] != userID {
		t.Fatalf("expected posted_by=%d for both postings", userID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

