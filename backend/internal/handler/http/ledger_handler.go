package http

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type LedgerHandler struct {
	generalLedgerRepo *repository.GeneralLedgerRepository
}

func NewLedgerHandler(generalLedgerRepo *repository.GeneralLedgerRepository) *LedgerHandler {
	return &LedgerHandler{generalLedgerRepo: generalLedgerRepo}
}

func (h *LedgerHandler) GetLedger(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	accountID, err := strconv.Atoi(r.URL.Query().Get("account_id"))
	if err != nil || accountID <= 0 {
		respondError(w, http.StatusBadRequest, "Valid account_id is required")
		return
	}

	startDate, err := parseLedgerDate(r.URL.Query().Get("start"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid start date")
		return
	}
	endDate, err := parseLedgerDate(r.URL.Query().Get("end"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid end date")
		return
	}

	limit, err := parseLedgerLimit(r.URL.Query().Get("limit"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid limit")
		return
	}
	cursor, err := parseLedgerCursor(r.URL.Query().Get("cursor"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid cursor")
		return
	}

	lines, nextCursor, err := h.generalLedgerRepo.GetLedgerByAccountCursor(tenantID, restaurantID, accountID, startDate, endDate, limit, cursor)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load ledger")
		return
	}
	periodSummaries, err := h.generalLedgerRepo.GetLedgerPeriodSummaries(tenantID, restaurantID, accountID, startDate, endDate)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load ledger period summaries")
		return
	}

	balance, err := h.generalLedgerRepo.GetAccountBalance(tenantID, restaurantID, accountID, endDate)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to calculate account balance")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"account_id":       accountID,
			"lines":            lines,
			"balance":          balance,
			"next_cursor":      nextCursor,
			"limit":            limit,
			"period_summaries": periodSummaries,
		},
	})
}

func (h *LedgerHandler) GetTrialBalance(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	asOf, err := parseLedgerDate(r.URL.Query().Get("as_of"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid as_of date")
		return
	}

	asOfDate := time.Now()
	if asOf != nil {
		asOfDate = *asOf
	}

	rows, err := h.generalLedgerRepo.GetTrialBalance(tenantID, restaurantID, asOfDate)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load trial balance")
		return
	}

	totalDebit := 0.0
	totalCredit := 0.0
	for _, row := range rows {
		totalDebit += row.DebitBalance
		totalCredit += row.CreditBalance
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"as_of":        asOfDate.Format("2006-01-02"),
			"rows":         rows,
			"total_debit":  totalDebit,
			"total_credit": totalCredit,
			"is_balanced":  roundLedgerAmount(totalDebit-totalCredit) == 0,
		},
	})
}

func parseLedgerDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseLedgerLimit(value string) (int, error) {
	if strings.TrimSpace(value) == "" {
		return 200, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 || parsed > 500 {
		return 0, strconv.ErrSyntax
	}
	return parsed, nil
}

func parseLedgerCursor(value string) (*int, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return nil, strconv.ErrSyntax
	}
	return &parsed, nil
}

func roundLedgerAmount(value float64) float64 {
	return math.Round(value*100) / 100
}
