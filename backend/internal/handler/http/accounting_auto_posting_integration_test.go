package http

import (
	"bytes"
	"context"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	pkgjwt "pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

type recordingAutoPostingService struct {
	salesOrders     []*domain.Order
	salesPostedBy   []int
	cogsOrders      []*domain.Order
	cogsPostedBy    []int
	payrollSalaries []*domain.Salary
	payrollPostedBy []int
}

func (s *recordingAutoPostingService) PostSalesEntry(order *domain.Order, postedBy *int) error {
	s.salesOrders = append(s.salesOrders, cloneOrderForTest(order))
	s.salesPostedBy = append(s.salesPostedBy, derefPostedBy(postedBy))
	return nil
}

func (s *recordingAutoPostingService) PostCOGSEntry(order *domain.Order, postedBy *int) error {
	s.cogsOrders = append(s.cogsOrders, cloneOrderForTest(order))
	s.cogsPostedBy = append(s.cogsPostedBy, derefPostedBy(postedBy))
	return nil
}

func (s *recordingAutoPostingService) PostPayrollEntry(salary *domain.Salary, postedBy *int) error {
	s.payrollSalaries = append(s.payrollSalaries, cloneSalaryForTest(salary))
	s.payrollPostedBy = append(s.payrollPostedBy, derefPostedBy(postedBy))
	return nil
}

func TestAdminOrderHandler_UpdatePaymentStatusPersistsDurableAccountingIntent(t *testing.T) {
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
	mock.ExpectExec(`UPDATE orders`).WithArgs("paid", int64(501), int64(1), int64(2)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO accounting_event_outbox`).WithArgs("paid-order:1:2:501", int64(1), int64(2), int64(501), 77).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := bytes.NewBufferString(`{"payment_status":"paid"}`)
	req := httptest.NewRequest(stdhttp.MethodPut, "/api/v1/admin/orders/501/payment-status", body)
	req.SetPathValue("id", "501")
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
func TestHRSalaryHandler_MarkAsPaidTriggersPayrollAutoPosting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	salaryRepo := repository.NewSalaryRepository(db)
	handler := NewHRSalaryHandler(salaryRepo, db)
	autoPosting := &recordingAutoPostingService{}
	handler.SetAutoPostingEngine(autoPosting)

	paymentDate := time.Date(2026, time.February, 28, 16, 45, 0, 0, time.UTC)
	periodStart := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC)
	salaryID := 44
	userID := 88

	mock.ExpectExec(`UPDATE salaries`).
		WithArgs(userID, "BANK-REF-44", salaryID, 1, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(`FROM salaries`).
		WithArgs(salaryID, 1, 2).
		WillReturnRows(
			sqlmock.NewRows(salaryColumnsForAccountingTests()).
				AddRow(
					salaryID, 1, 2, 9, periodStart,
					periodEnd, paymentDate, 2, 2026, 3000.0,
					"USD", 0.0, 0.0, 0.0,
					200.0, 0.0, 100.0, 0.0, 0.0,
					"{}", 150.0, 50.0, 25.0,
					0.0, 0.0, 0.0, 0.0,
					"{}", 3300.0, 225.0, 3075.0,
					22, 0, 176.0, 0.0,
					"bank_transfer", "BANK-REF-44", nil,
					nil, "paid", true, nil, paymentDate, paymentDate,
				),
		)

	body, err := json.Marshal(map[string]string{"payment_reference": "BANK-REF-44"})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(stdhttp.MethodPost, "/api/v1/hr/salaries/44/mark-paid", bytes.NewReader(body))
	req.SetPathValue("id", "44")
	req = withAccountingAuthContext(req, 1, 2, userID)

	rec := httptest.NewRecorder()
	handler.MarkAsPaid(rec, req)

	if rec.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	assertJSONSuccess(t, rec)

	if len(autoPosting.payrollSalaries) != 1 {
		t.Fatalf("expected 1 payroll auto-post, got %d", len(autoPosting.payrollSalaries))
	}
	if autoPosting.payrollSalaries[0].ID != salaryID {
		t.Fatalf("expected payroll salary id %d, got %d", salaryID, autoPosting.payrollSalaries[0].ID)
	}
	if autoPosting.payrollSalaries[0].Status != "paid" || !autoPosting.payrollSalaries[0].IsPaid {
		t.Fatalf("expected paid salary payload, got status=%s is_paid=%v", autoPosting.payrollSalaries[0].Status, autoPosting.payrollSalaries[0].IsPaid)
	}
	if autoPosting.payrollPostedBy[0] != userID {
		t.Fatalf("expected payroll posted_by=%d, got %d", userID, autoPosting.payrollPostedBy[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func withAccountingAuthContext(req *stdhttp.Request, tenantID, restaurantID int64, userID int) *stdhttp.Request {
	ctx := context.WithValue(req.Context(), middleware.TenantContextKey, tenantID)
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, restaurantID)
	ctx = context.WithValue(ctx, middleware.UserContextKey, &pkgjwt.Claims{
		UserID:       userID,
		TenantID:     int(tenantID),
		RestaurantID: int(restaurantID),
	})
	return req.WithContext(ctx)
}

func assertJSONSuccess(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	var payload map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	success, ok := payload["success"].(bool)
	if !ok || !success {
		t.Fatalf("expected success=true response, got %+v", payload)
	}
}

func derefPostedBy(postedBy *int) int {
	if postedBy == nil {
		return 0
	}
	return *postedBy
}

func cloneOrderForTest(order *domain.Order) *domain.Order {
	if order == nil {
		return nil
	}

	clone := *order
	if len(order.Items) > 0 {
		clone.Items = append([]domain.OrderItem(nil), order.Items...)
	}
	return &clone
}

func cloneSalaryForTest(salary *domain.Salary) *domain.Salary {
	if salary == nil {
		return nil
	}

	clone := *salary
	return &clone
}

func orderDetailColumnsForAccountingTests() []string {
	return []string{
		"id", "tenant_id", "restaurant_id", "order_number", "customer_name", "customer_email",
		"customer_phone", "delivery_address", "delivery_city", "delivery_area",
		"delivery_zip_code", "delivery_latitude", "delivery_longitude",
		"delivery_instructions", "subtotal", "tax_amount", "discount_amount",
		"delivery_fee", "total_amount", "payment_method", "payment_status",
		"status", "estimated_delivery_time", "actual_delivery_time", "notes", "order_source",
		"table_id", "table_number", "table_name", "table_location", "capacity",
		"created_at", "updated_at",
	}
}

func orderItemColumnsForAccountingTests() []string {
	return []string{
		"id", "order_id", "product_id", "variant_id", "product_name", "variant_name",
		"quantity", "unit_price", "discount_amount", "total_price",
		"special_instructions", "addons", "created_at", "updated_at",
	}
}

func salaryColumnsForAccountingTests() []string {
	return []string{
		"id", "tenant_id", "restaurant_id", "employee_id", "pay_period_start",
		"pay_period_end", "payment_date", "month", "year", "base_salary",
		"currency", "overtime_hours", "overtime_rate", "overtime_amount",
		"bonus", "commission", "allowances", "tips", "other_earnings",
		"earnings_details", "tax", "social_insurance", "health_insurance",
		"pension", "loan_deduction", "advance_deduction", "other_deductions",
		"deductions_details", "gross_salary", "total_deductions", "net_salary",
		"days_worked", "days_absent", "total_hours_worked", "total_overtime_hours",
		"payment_method", "payment_reference", "bank_account_number",
		"bank_name", "status", "is_paid", "notes", "created_at", "updated_at",
	}
}
