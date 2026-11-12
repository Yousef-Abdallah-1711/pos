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

type BankAccountHandler struct {
	bankRepo *repository.BankAccountRepository
}

const bankAccountHandlerTimeout = 10 * time.Second

func NewBankAccountHandler(bankRepo *repository.BankAccountRepository) *BankAccountHandler {
	return &BankAccountHandler{bankRepo: bankRepo}
}

func (h *BankAccountHandler) ListBankAccounts(w http.ResponseWriter, r *http.Request) {
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

	ctx, cancel := context.WithTimeout(r.Context(), bankAccountHandlerTimeout)
	defer cancel()
	accounts, total, err := h.bankRepo.ListBankAccountsPaged(ctx, tenantID, restaurantID, page, pageSize)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch bank accounts",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    accounts,
		"pagination": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *BankAccountHandler) GetBankAccount(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid bank account ID",
		})
		return
	}

	account, err := h.bankRepo.GetBankAccountByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch bank account",
		})
		return
	}
	if account == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Bank account not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    account,
	})
}

func (h *BankAccountHandler) CreateBankAccount(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req domain.CreateBankAccountRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	account, err := buildBankAccountFromRequest(req, tenantID, restaurantID, nil, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	account.CurrentBalance = account.OpeningBalance

	id, err := h.bankRepo.CreateBankAccount(&account)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	created, err := h.bankRepo.GetBankAccountByID(tenantID, restaurantID, id)
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

func (h *BankAccountHandler) UpdateBankAccount(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid bank account ID",
		})
		return
	}

	existing, err := h.bankRepo.GetBankAccountByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch bank account",
		})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Bank account not found",
		})
		return
	}

	var req domain.UpdateBankAccountRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	account, err := buildBankAccountFromRequest(req, tenantID, restaurantID, existing, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	account.ID = id
	account.CurrentBalance = existing.CurrentBalance
	account.CreatedAt = existing.CreatedAt
	account.UpdatedAt = time.Now()
	account.CreatedBy = existing.CreatedBy

	if err := h.bankRepo.UpdateBankAccount(&account); err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		writeJSON(w, statusCode, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, err := h.bankRepo.GetBankAccountByID(tenantID, restaurantID, id)
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

func (h *BankAccountHandler) DeleteBankAccount(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid bank account ID",
		})
		return
	}

	if err := h.bankRepo.DeleteBankAccount(tenantID, restaurantID, id); err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "not found") {
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

func (h *BankAccountHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid bank account ID",
		})
		return
	}

	transactions, err := h.bankRepo.GetBankTransactions(tenantID, restaurantID, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
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
		"data":    transactions,
	})
}

func buildBankAccountFromRequest(req domain.CreateBankAccountRequest, tenantID, restaurantID int, existing *domain.BankAccount, userID int) (domain.BankAccount, error) {
	if err := requireAccountingString(req.AccountName, "account_name"); err != nil {
		return domain.BankAccount{}, err
	}
	if err := requireAccountingString(req.BankName, "bank_name"); err != nil {
		return domain.BankAccount{}, err
	}
	if err := requireAccountingString(req.AccountNumber, "account_number"); err != nil {
		return domain.BankAccount{}, err
	}

	openingDate, err := parseAccountingInputDate(req.OpeningDate, "opening_date")
	if err != nil {
		return domain.BankAccount{}, err
	}
	if openingDate.IsZero() {
		if existing != nil {
			openingDate = existing.OpeningDate
		} else {
			openingDate = time.Now()
		}
	}

	isActive := true
	if existing != nil {
		isActive = existing.IsActive
	}
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	account := domain.BankAccount{
		TenantID:       tenantID,
		RestaurantID:   restaurantID,
		AccountName:    strings.TrimSpace(req.AccountName),
		AccountNameAr:  strings.TrimSpace(req.AccountNameAr),
		BankName:       strings.TrimSpace(req.BankName),
		BankBranch:     strings.TrimSpace(req.BankBranch),
		AccountNumber:  strings.TrimSpace(req.AccountNumber),
		IBAN:           strings.TrimSpace(req.IBAN),
		SwiftCode:      strings.TrimSpace(req.SwiftCode),
		Currency:       defaultAccountingString(req.Currency, "EGP"),
		AccountType:    defaultAccountingString(req.AccountType, domain.BankAccountTypeChecking),
		OpeningBalance: req.OpeningBalance,
		OpeningDate:    openingDate,
		GLAccountID:    req.GLAccountID,
		IsActive:       isActive,
		IsPrimary:      req.IsPrimary,
		UpdatedBy:      &userID,
	}
	if existing == nil {
		account.CreatedBy = &userID
	}
	return account, nil
}
