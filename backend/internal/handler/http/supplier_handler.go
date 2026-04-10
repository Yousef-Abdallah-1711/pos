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

type SupplierHandler struct {
	supplierRepo *repository.SupplierRepository
}

const supplierHandlerTimeout = 10 * time.Second

func NewSupplierHandler(supplierRepo *repository.SupplierRepository) *SupplierHandler {
	return &SupplierHandler{supplierRepo: supplierRepo}
}

func (h *SupplierHandler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
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

	ctx, cancel := context.WithTimeout(r.Context(), supplierHandlerTimeout)
	defer cancel()
	suppliers, total, err := h.supplierRepo.ListSuppliersPaged(ctx, tenantID, restaurantID, page, pageSize)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch suppliers",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    suppliers,
		"pagination": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *SupplierHandler) GetSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid supplier ID",
		})
		return
	}

	supplier, err := h.supplierRepo.GetSupplierByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch supplier",
		})
		return
	}
	if supplier == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Supplier not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    supplier,
	})
}

func (h *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req domain.CreateSupplierRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	if strings.TrimSpace(req.CompanyName) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "company_name is required",
		})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	supplier := domain.Supplier{
		TenantID:         tenantID,
		RestaurantID:     restaurantID,
		CompanyName:      strings.TrimSpace(req.CompanyName),
		CompanyNameAr:    strings.TrimSpace(req.CompanyNameAr),
		ContactPerson:    strings.TrimSpace(req.ContactPerson),
		Email:            strings.TrimSpace(req.Email),
		Phone:            strings.TrimSpace(req.Phone),
		Mobile:           strings.TrimSpace(req.Mobile),
		Website:          strings.TrimSpace(req.Website),
		Address:          strings.TrimSpace(req.Address),
		AddressAr:        strings.TrimSpace(req.AddressAr),
		City:             strings.TrimSpace(req.City),
		State:            strings.TrimSpace(req.State),
		PostalCode:       strings.TrimSpace(req.PostalCode),
		Country:          defaultAccountingString(req.Country, "Egypt"),
		TaxNumber:        strings.TrimSpace(req.TaxNumber),
		PaymentTerms:     defaultAccountingString(req.PaymentTerms, "net_30"),
		CreditLimit:      req.CreditLimit,
		Currency:         defaultAccountingString(req.Currency, "EGP"),
		PayableAccountID: req.PayableAccountID,
		Category:         strings.TrimSpace(req.Category),
		BankName:         strings.TrimSpace(req.BankName),
		BankAccount:      strings.TrimSpace(req.BankAccount),
		IBAN:             strings.TrimSpace(req.IBAN),
		Notes:            strings.TrimSpace(req.Notes),
		IsActive:         isActive,
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}

	id, err := h.supplierRepo.CreateSupplier(&supplier)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	created, err := h.supplierRepo.GetSupplierByID(tenantID, restaurantID, id)
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

func (h *SupplierHandler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid supplier ID",
		})
		return
	}

	existing, err := h.supplierRepo.GetSupplierByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch supplier",
		})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Supplier not found",
		})
		return
	}

	var req domain.UpdateSupplierRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	if strings.TrimSpace(req.CompanyName) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "company_name is required",
		})
		return
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	supplier := domain.Supplier{
		ID:               id,
		TenantID:         tenantID,
		RestaurantID:     restaurantID,
		SupplierCode:     existing.SupplierCode,
		CompanyName:      strings.TrimSpace(req.CompanyName),
		CompanyNameAr:    strings.TrimSpace(req.CompanyNameAr),
		ContactPerson:    strings.TrimSpace(req.ContactPerson),
		Email:            strings.TrimSpace(req.Email),
		Phone:            strings.TrimSpace(req.Phone),
		Mobile:           strings.TrimSpace(req.Mobile),
		Website:          strings.TrimSpace(req.Website),
		Address:          strings.TrimSpace(req.Address),
		AddressAr:        strings.TrimSpace(req.AddressAr),
		City:             strings.TrimSpace(req.City),
		State:            strings.TrimSpace(req.State),
		PostalCode:       strings.TrimSpace(req.PostalCode),
		Country:          defaultAccountingString(req.Country, "Egypt"),
		TaxNumber:        strings.TrimSpace(req.TaxNumber),
		PaymentTerms:     defaultAccountingString(req.PaymentTerms, "net_30"),
		CreditLimit:      req.CreditLimit,
		CurrentBalance:   existing.CurrentBalance,
		Currency:         defaultAccountingString(req.Currency, "EGP"),
		PayableAccountID: req.PayableAccountID,
		Category:         strings.TrimSpace(req.Category),
		BankName:         strings.TrimSpace(req.BankName),
		BankAccount:      strings.TrimSpace(req.BankAccount),
		IBAN:             strings.TrimSpace(req.IBAN),
		Notes:            strings.TrimSpace(req.Notes),
		IsActive:         isActive,
		CreatedBy:        existing.CreatedBy,
		UpdatedBy:        &userID,
	}

	if err := h.supplierRepo.UpdateSupplier(&supplier); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "supplier not found" {
			statusCode = http.StatusNotFound
		}
		writeJSON(w, statusCode, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, err := h.supplierRepo.GetSupplierByID(tenantID, restaurantID, id)
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

func (h *SupplierHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid supplier ID",
		})
		return
	}

	if err := h.supplierRepo.DeleteSupplier(tenantID, restaurantID, id); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "supplier not found" {
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

func defaultAccountingString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
