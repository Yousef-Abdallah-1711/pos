package http

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type CustomerAccountHandler struct {
	customerRepo *repository.CustomerAccountRepository
}

const customerAcctHandlerTimeout = 10 * time.Second

func NewCustomerAccountHandler(customerRepo *repository.CustomerAccountRepository) *CustomerAccountHandler {
	return &CustomerAccountHandler{customerRepo: customerRepo}
}

func (h *CustomerAccountHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	page := GetQueryParamInt(r, "page", 1)
	pageSize := GetQueryParamInt(r, "page_size", 20)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	ctx, cancel := context.WithTimeout(r.Context(), customerAcctHandlerTimeout)
	defer cancel()
	customers, total, err := h.customerRepo.ListCustomersPaged(ctx, tenantID, restaurantID, page, pageSize)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch customer accounts",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    customers,
		"pagination": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *CustomerAccountHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid customer account ID",
		})
		return
	}

	customer, err := h.customerRepo.GetCustomerByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch customer account",
		})
		return
	}
	if customer == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Customer account not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    customer,
	})
}

func (h *CustomerAccountHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req domain.CreateCustomerAccountRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	if strings.TrimSpace(req.CustomerName) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "customer_name is required",
		})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	customer := domain.CustomerAccount{
		TenantID:            tenantID,
		RestaurantID:        restaurantID,
		CustomerName:        strings.TrimSpace(req.CustomerName),
		CustomerNameAr:      strings.TrimSpace(req.CustomerNameAr),
		CustomerType:        defaultAccountingString(req.CustomerType, "individual"),
		ContactPerson:       strings.TrimSpace(req.ContactPerson),
		Email:               strings.TrimSpace(req.Email),
		Phone:               strings.TrimSpace(req.Phone),
		Mobile:              strings.TrimSpace(req.Mobile),
		Address:             strings.TrimSpace(req.Address),
		AddressAr:           strings.TrimSpace(req.AddressAr),
		City:                strings.TrimSpace(req.City),
		Country:             defaultAccountingString(req.Country, "Egypt"),
		TaxNumber:           strings.TrimSpace(req.TaxNumber),
		PaymentTerms:        defaultAccountingString(req.PaymentTerms, "immediate"),
		CreditLimit:         req.CreditLimit,
		Currency:            defaultAccountingString(req.Currency, "EGP"),
		ReceivableAccountID: req.ReceivableAccountID,
		Notes:               strings.TrimSpace(req.Notes),
		IsActive:            isActive,
		CreatedBy:           &userID,
		UpdatedBy:           &userID,
	}

	id, err := h.customerRepo.CreateCustomer(&customer)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	created, err := h.customerRepo.GetCustomerByID(tenantID, restaurantID, id)
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

func (h *CustomerAccountHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid customer account ID",
		})
		return
	}

	existing, err := h.customerRepo.GetCustomerByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch customer account",
		})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Customer account not found",
		})
		return
	}

	var req domain.UpdateCustomerAccountRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	if strings.TrimSpace(req.CustomerName) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "customer_name is required",
		})
		return
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	customer := domain.CustomerAccount{
		ID:                  id,
		TenantID:            tenantID,
		RestaurantID:        restaurantID,
		CustomerCode:        existing.CustomerCode,
		CustomerName:        strings.TrimSpace(req.CustomerName),
		CustomerNameAr:      strings.TrimSpace(req.CustomerNameAr),
		CustomerType:        defaultAccountingString(req.CustomerType, "individual"),
		ContactPerson:       strings.TrimSpace(req.ContactPerson),
		Email:               strings.TrimSpace(req.Email),
		Phone:               strings.TrimSpace(req.Phone),
		Mobile:              strings.TrimSpace(req.Mobile),
		Address:             strings.TrimSpace(req.Address),
		AddressAr:           strings.TrimSpace(req.AddressAr),
		City:                strings.TrimSpace(req.City),
		Country:             defaultAccountingString(req.Country, "Egypt"),
		TaxNumber:           strings.TrimSpace(req.TaxNumber),
		PaymentTerms:        defaultAccountingString(req.PaymentTerms, "immediate"),
		CreditLimit:         req.CreditLimit,
		CurrentBalance:      existing.CurrentBalance,
		Currency:            defaultAccountingString(req.Currency, "EGP"),
		ReceivableAccountID: req.ReceivableAccountID,
		Notes:               strings.TrimSpace(req.Notes),
		IsActive:            isActive,
		CreatedBy:           existing.CreatedBy,
		UpdatedBy:           &userID,
	}

	if err := h.customerRepo.UpdateCustomer(&customer); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "customer account not found" {
			statusCode = http.StatusNotFound
		}
		writeJSON(w, statusCode, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, err := h.customerRepo.GetCustomerByID(tenantID, restaurantID, id)
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

func (h *CustomerAccountHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid customer account ID",
		})
		return
	}

	if err := h.customerRepo.DeleteCustomer(tenantID, restaurantID, id); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "customer account not found" {
			statusCode = http.StatusNotFound
		}
		writeJSON(w, statusCode, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func (h *CustomerAccountHandler) GetCustomerStatement(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid customer account ID",
		})
		return
	}

	statement, err := h.customerRepo.GetCustomerStatement(tenantID, restaurantID, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "customer account not found" {
			statusCode = http.StatusNotFound
		}
		writeJSON(w, statusCode, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    statement,
	})
}
