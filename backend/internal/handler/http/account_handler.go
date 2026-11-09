package http

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/usecase"
)

type AccountHandler struct {
	accountingUC *usecase.AccountingUseCase
}

const accountListTimeout = 10 * time.Second

func NewAccountHandler(accountingUC *usecase.AccountingUseCase) *AccountHandler {
	return &AccountHandler{accountingUC: accountingUC}
}

func (h *AccountHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	tree := GetQueryParamBool(r, "tree")
	includeInactive := GetQueryParamBool(r, "include_inactive")
	page := GetQueryParamInt(r, "page", 1)
	pageSize := GetQueryParamInt(r, "page_size", 50)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}

	ctx, cancel := context.WithTimeout(r.Context(), accountListTimeout)
	defer cancel()

	var (
		accounts []domain.Account
		total    int
		err      error
	)
	if tree {
		accounts, err = h.accountingUC.ListAccounts(tenantID, restaurantID, true, includeInactive)
		total = len(accounts)
	} else {
		accounts, total, err = h.accountingUC.ListAccountsPaged(ctx, tenantID, restaurantID, includeInactive, page, pageSize)
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch accounts",
		})
		return
	}

	pagedAccounts := accounts
	if tree {
		offset := (page - 1) * pageSize
		if offset > total {
			offset = total
		}
		end := offset + pageSize
		if end > total {
			end = total
		}
		pagedAccounts = accounts[offset:end]
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    pagedAccounts,
		"pagination": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid account ID",
		})
		return
	}

	account, err := h.accountingUC.GetAccountByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch account",
		})
		return
	}
	if account == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Account not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    account,
	})
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req domain.CreateAccountRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	account := domain.Account{
		TenantID:       tenantID,
		RestaurantID:   restaurantID,
		AccountCode:    strings.TrimSpace(req.AccountCode),
		AccountName:    strings.TrimSpace(req.AccountName),
		AccountNameAr:  strings.TrimSpace(req.AccountNameAr),
		Description:    strings.TrimSpace(req.Description),
		DescriptionAr:  strings.TrimSpace(req.DescriptionAr),
		AccountType:    strings.TrimSpace(req.AccountType),
		AccountSubtype: strings.TrimSpace(req.AccountSubtype),
		NormalBalance:  strings.TrimSpace(req.NormalBalance),
		ParentID:       req.ParentID,
		IsHeader:       req.IsHeader,
		Currency:       strings.TrimSpace(req.Currency),
		IsActive:       isActive,
		IsReconcilable: req.IsReconcilable,
		TaxRate:        req.TaxRate,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	id, err := h.accountingUC.CreateAccount(&account)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	created, err := h.accountingUC.GetAccountByID(tenantID, restaurantID, id)
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

func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid account ID",
		})
		return
	}

	existing, err := h.accountingUC.GetAccountByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch account",
		})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Account not found",
		})
		return
	}

	var req domain.UpdateAccountRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	account := domain.Account{
		ID:             id,
		TenantID:       tenantID,
		RestaurantID:   restaurantID,
		AccountCode:    strings.TrimSpace(req.AccountCode),
		AccountName:    strings.TrimSpace(req.AccountName),
		AccountNameAr:  strings.TrimSpace(req.AccountNameAr),
		Description:    strings.TrimSpace(req.Description),
		DescriptionAr:  strings.TrimSpace(req.DescriptionAr),
		AccountType:    strings.TrimSpace(req.AccountType),
		AccountSubtype: strings.TrimSpace(req.AccountSubtype),
		NormalBalance:  strings.TrimSpace(req.NormalBalance),
		ParentID:       req.ParentID,
		IsHeader:       req.IsHeader,
		Currency:       strings.TrimSpace(req.Currency),
		IsActive:       isActive,
		IsReconcilable: req.IsReconcilable,
		TaxRate:        req.TaxRate,
		UpdatedBy:      &userID,
	}

	if err := h.accountingUC.UpdateAccount(&account); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "account not found" {
			statusCode = http.StatusNotFound
		}
		writeJSON(w, statusCode, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, err := h.accountingUC.GetAccountByID(tenantID, restaurantID, id)
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

func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid account ID",
		})
		return
	}

	if err := h.accountingUC.DeleteAccount(tenantID, restaurantID, id); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "account not found" {
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
