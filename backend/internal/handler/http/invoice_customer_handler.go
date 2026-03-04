package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type CustomerInvoiceHandler struct {
	invoiceRepo *repository.CustomerInvoiceRepository
}

func NewCustomerInvoiceHandler(invoiceRepo *repository.CustomerInvoiceRepository) *CustomerInvoiceHandler {
	return &CustomerInvoiceHandler{invoiceRepo: invoiceRepo}
}

func (h *CustomerInvoiceHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	invoices, err := h.invoiceRepo.ListInvoices(tenantID, restaurantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch customer invoices",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    invoices,
	})
}

func (h *CustomerInvoiceHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid customer invoice ID",
		})
		return
	}

	invoice, err := h.invoiceRepo.GetInvoiceByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch customer invoice",
		})
		return
	}
	if invoice == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Customer invoice not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    invoice,
	})
}

func (h *CustomerInvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req domain.CreateCustomerInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	invoice, err := buildCustomerInvoiceFromRequest(req, tenantID, restaurantID, 0, nil, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	id, err := h.invoiceRepo.CreateInvoice(&invoice)
	if err != nil {
		writeAccountingMutationError(w, err)
		return
	}

	created, err := h.invoiceRepo.GetInvoiceByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"data":    map[string]int{"id": id},
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    created,
	})
}

func (h *CustomerInvoiceHandler) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid customer invoice ID",
		})
		return
	}

	existing, err := h.invoiceRepo.GetInvoiceByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch customer invoice",
		})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Customer invoice not found",
		})
		return
	}

	var req domain.UpdateCustomerInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	invoice, err := buildCustomerInvoiceFromRequest(req, tenantID, restaurantID, id, existing.CreatedBy, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := h.invoiceRepo.UpdateInvoice(&invoice); err != nil {
		writeAccountingMutationError(w, err)
		return
	}

	updated, err := h.invoiceRepo.GetInvoiceByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    updated,
	})
}

func (h *CustomerInvoiceHandler) MarkAsSent(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid customer invoice ID",
		})
		return
	}

	if err := h.invoiceRepo.MarkAsSent(tenantID, restaurantID, id, &userID); err != nil {
		writeAccountingMutationError(w, err)
		return
	}

	updated, err := h.invoiceRepo.GetInvoiceByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    updated,
	})
}

func (h *CustomerInvoiceHandler) VoidInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid customer invoice ID",
		})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := decodeAccountingActionBody(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = "Voided by user"
	}

	if err := h.invoiceRepo.VoidInvoice(tenantID, restaurantID, id, userID, reason); err != nil {
		writeAccountingMutationError(w, err)
		return
	}

	updated, err := h.invoiceRepo.GetInvoiceByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    updated,
	})
}

func (h *CustomerInvoiceHandler) ListOverdueInvoices(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	invoices, err := h.invoiceRepo.GetOverdueInvoices(tenantID, restaurantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch overdue customer invoices",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    invoices,
	})
}

func buildCustomerInvoiceFromRequest(req domain.CreateCustomerInvoiceRequest, tenantID, restaurantID, id int, createdBy *int, userID int) (domain.CustomerInvoice, error) {
	if req.CustomerID <= 0 {
		return domain.CustomerInvoice{}, fmt.Errorf("customer_id is required")
	}
	if len(req.Lines) == 0 {
		return domain.CustomerInvoice{}, fmt.Errorf("at least one invoice line is required")
	}

	invoiceDate, err := parseAccountingInputDate(req.InvoiceDate, "invoice_date")
	if err != nil {
		return domain.CustomerInvoice{}, err
	}
	dueDate, err := parseAccountingInputDate(req.DueDate, "due_date")
	if err != nil {
		return domain.CustomerInvoice{}, err
	}

	invoice := domain.CustomerInvoice{
		ID:            id,
		TenantID:      tenantID,
		RestaurantID:  restaurantID,
		InvoiceDate:   invoiceDate,
		DueDate:       dueDate,
		CustomerID:    req.CustomerID,
		OrderID:       req.OrderID,
		Currency:      defaultAccountingString(req.Currency, "EGP"),
		Status:        defaultAccountingString(req.Status, domain.InvoiceStatusDraft),
		Notes:         strings.TrimSpace(req.Notes),
		Terms:         strings.TrimSpace(req.Terms),
		InternalNotes: strings.TrimSpace(req.InternalNotes),
		CreatedBy:     createdBy,
		UpdatedBy:     &userID,
		Lines:         mapAccountingInvoiceLines(req.Lines),
	}
	if createdBy == nil {
		invoice.CreatedBy = &userID
	}
	return invoice, nil
}
