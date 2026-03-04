package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type PaymentAccountingHandler struct {
	paymentRepo *repository.PaymentAccountingRepository
}

func NewPaymentAccountingHandler(paymentRepo *repository.PaymentAccountingRepository) *PaymentAccountingHandler {
	return &PaymentAccountingHandler{paymentRepo: paymentRepo}
}

func (h *PaymentAccountingHandler) ListPaymentsReceived(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	payments, err := h.paymentRepo.ListPaymentsReceived(tenantID, restaurantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch received payments",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    payments,
	})
}

func (h *PaymentAccountingHandler) CreatePaymentReceived(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req domain.CreatePaymentReceivedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	payment, err := buildPaymentReceivedFromRequest(req, tenantID, restaurantID, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	id, err := h.paymentRepo.CreatePaymentReceived(&payment)
	if err != nil {
		writeAccountingMutationError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    map[string]int{"id": id},
	})
}

func (h *PaymentAccountingHandler) ListPaymentsMade(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	payments, err := h.paymentRepo.ListPaymentsMade(tenantID, restaurantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch supplier payments",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    payments,
	})
}

func (h *PaymentAccountingHandler) CreatePaymentMade(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req domain.CreatePaymentMadeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	payment, err := buildPaymentMadeFromRequest(req, tenantID, restaurantID, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	id, err := h.paymentRepo.CreatePaymentMade(&payment)
	if err != nil {
		writeAccountingMutationError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    map[string]int{"id": id},
	})
}

func buildPaymentReceivedFromRequest(req domain.CreatePaymentReceivedRequest, tenantID, restaurantID, userID int) (domain.PaymentReceived, error) {
	if req.Amount <= 0 {
		return domain.PaymentReceived{}, fmt.Errorf("amount must be greater than zero")
	}
	if req.CustomerID <= 0 && req.InvoiceID == nil {
		return domain.PaymentReceived{}, fmt.Errorf("customer_id or invoice_id is required")
	}
	if err := requireAccountingString(req.PaymentMethod, "payment_method"); err != nil {
		return domain.PaymentReceived{}, err
	}

	paymentDate, err := parseAccountingInputDate(req.PaymentDate, "payment_date")
	if err != nil {
		return domain.PaymentReceived{}, err
	}

	return domain.PaymentReceived{
		TenantID:      tenantID,
		RestaurantID:  restaurantID,
		PaymentDate:   paymentDate,
		CustomerID:    req.CustomerID,
		InvoiceID:     req.InvoiceID,
		Amount:        req.Amount,
		PaymentMethod: strings.TrimSpace(req.PaymentMethod),
		Reference:     strings.TrimSpace(req.Reference),
		BankAccountID: req.BankAccountID,
		Notes:         strings.TrimSpace(req.Notes),
		CreatedBy:     &userID,
	}, nil
}

func buildPaymentMadeFromRequest(req domain.CreatePaymentMadeRequest, tenantID, restaurantID, userID int) (domain.PaymentMade, error) {
	if req.Amount <= 0 {
		return domain.PaymentMade{}, fmt.Errorf("amount must be greater than zero")
	}
	if req.SupplierID <= 0 && req.InvoiceID == nil {
		return domain.PaymentMade{}, fmt.Errorf("supplier_id or invoice_id is required")
	}
	if err := requireAccountingString(req.PaymentMethod, "payment_method"); err != nil {
		return domain.PaymentMade{}, err
	}

	paymentDate, err := parseAccountingInputDate(req.PaymentDate, "payment_date")
	if err != nil {
		return domain.PaymentMade{}, err
	}

	return domain.PaymentMade{
		TenantID:      tenantID,
		RestaurantID:  restaurantID,
		PaymentDate:   paymentDate,
		SupplierID:    req.SupplierID,
		InvoiceID:     req.InvoiceID,
		Amount:        req.Amount,
		PaymentMethod: strings.TrimSpace(req.PaymentMethod),
		Reference:     strings.TrimSpace(req.Reference),
		BankAccountID: req.BankAccountID,
		Notes:         strings.TrimSpace(req.Notes),
		CreatedBy:     &userID,
	}, nil
}
