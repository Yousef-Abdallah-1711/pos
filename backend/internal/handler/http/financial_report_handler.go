package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
)

type FinancialReportHandler struct {
	financialReportUC *usecase.FinancialReportUseCase
	cache             service.KVCache
	cacheTTL          time.Duration
}

const defaultFinancialReportCacheTTL = 5 * time.Minute

func NewFinancialReportHandler(financialReportUC *usecase.FinancialReportUseCase) *FinancialReportHandler {
	return NewFinancialReportHandlerWithCache(financialReportUC, service.NewCacheFromEnv(), defaultFinancialReportCacheTTL)
}

func NewFinancialReportHandlerWithCache(financialReportUC *usecase.FinancialReportUseCase, cache service.KVCache, cacheTTL time.Duration) *FinancialReportHandler {
	if cache == nil {
		cache = service.NewInMemoryCache()
	}
	if cacheTTL <= 0 {
		cacheTTL = defaultFinancialReportCacheTTL
	}

	return &FinancialReportHandler{
		financialReportUC: financialReportUC,
		cache:             cache,
		cacheTTL:          cacheTTL,
	}
}

func (h *FinancialReportHandler) GetProfitLoss(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	from, to, err := parseFinancialDateRange(r.URL.Query().Get("date_from"), r.URL.Query().Get("date_to"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	compareFromRaw := r.URL.Query().Get("compare_date_from")
	compareToRaw := r.URL.Query().Get("compare_date_to")
	useComparison := false
	var compareFrom time.Time
	var compareTo time.Time
	if compareFromRaw != "" && compareToRaw != "" {
		useComparison = true
		parsedFrom, parsedTo, compareErr := parseFinancialDateRange(compareFromRaw, compareToRaw)
		if compareErr != nil {
			respondError(w, http.StatusBadRequest, compareErr.Error())
			return
		}
		compareFrom = parsedFrom
		compareTo = parsedTo
	}

	cacheKey := h.financialReportCacheKey(tenantID, restaurantID, r.URL.Path, r.URL.RawQuery)
	h.respondFinancialReportWithCache(w, r, cacheKey, func() (interface{}, error) {
		var report usecase.ProfitLossReport
		if useComparison {
			report, err = h.financialReportUC.GenerateProfitLossWithComparison(tenantID, restaurantID, from, to, compareFrom, compareTo)
		} else {
			report, err = h.financialReportUC.GenerateProfitLoss(tenantID, restaurantID, from, to)
		}
		if err != nil {
			return nil, err
		}
		return report, nil
	}, "Failed to generate profit and loss report")
}

func (h *FinancialReportHandler) GetBalanceSheet(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	asOf, err := parseAsOfDate(r.URL.Query().Get("as_of"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	cacheKey := h.financialReportCacheKey(tenantID, restaurantID, r.URL.Path, r.URL.RawQuery)
	h.respondFinancialReportWithCache(w, r, cacheKey, func() (interface{}, error) {
		return h.financialReportUC.GenerateBalanceSheet(tenantID, restaurantID, asOf)
	}, "Failed to generate balance sheet")
}

func (h *FinancialReportHandler) GetCashFlow(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	from, to, err := parseFinancialDateRange(r.URL.Query().Get("date_from"), r.URL.Query().Get("date_to"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	cacheKey := h.financialReportCacheKey(tenantID, restaurantID, r.URL.Path, r.URL.RawQuery)
	h.respondFinancialReportWithCache(w, r, cacheKey, func() (interface{}, error) {
		return h.financialReportUC.GenerateCashFlow(tenantID, restaurantID, from, to)
	}, "Failed to generate cash flow report")
}

func (h *FinancialReportHandler) GetReceivablesAging(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	asOf, err := parseAsOfDate(r.URL.Query().Get("as_of"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	cacheKey := h.financialReportCacheKey(tenantID, restaurantID, r.URL.Path, r.URL.RawQuery)
	h.respondFinancialReportWithCache(w, r, cacheKey, func() (interface{}, error) {
		return h.financialReportUC.GenerateReceivablesAging(tenantID, restaurantID, asOf)
	}, "Failed to generate receivables aging report")
}

func (h *FinancialReportHandler) GetPayablesAging(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	asOf, err := parseAsOfDate(r.URL.Query().Get("as_of"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	cacheKey := h.financialReportCacheKey(tenantID, restaurantID, r.URL.Path, r.URL.RawQuery)
	h.respondFinancialReportWithCache(w, r, cacheKey, func() (interface{}, error) {
		return h.financialReportUC.GeneratePayablesAging(tenantID, restaurantID, asOf)
	}, "Failed to generate payables aging report")
}

func (h *FinancialReportHandler) GetInventoryValuation(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	asOf, err := parseAsOfDate(r.URL.Query().Get("as_of"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	cacheKey := h.financialReportCacheKey(tenantID, restaurantID, r.URL.Path, r.URL.RawQuery)
	h.respondFinancialReportWithCache(w, r, cacheKey, func() (interface{}, error) {
		return h.financialReportUC.GenerateInventoryValuation(tenantID, restaurantID, asOf)
	}, "Failed to generate inventory valuation report")
}

func (h *FinancialReportHandler) GetTrialBalance(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	asOf, err := parseAsOfDate(r.URL.Query().Get("as_of"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	cacheKey := h.financialReportCacheKey(tenantID, restaurantID, r.URL.Path, r.URL.RawQuery)
	h.respondFinancialReportWithCache(w, r, cacheKey, func() (interface{}, error) {
		return h.financialReportUC.GenerateTrialBalance(tenantID, restaurantID, asOf)
	}, "Failed to generate trial balance report")
}

func (h *FinancialReportHandler) financialReportCacheKey(tenantID, restaurantID int, path, rawQuery string) string {
	return fmt.Sprintf("financial_reports:%d:%d:%s:%s", tenantID, restaurantID, strings.TrimSpace(path), strings.TrimSpace(rawQuery))
}

func (h *FinancialReportHandler) respondFinancialReportWithCache(
	w http.ResponseWriter,
	r *http.Request,
	cacheKey string,
	build func() (interface{}, error),
	errorMessage string,
) {
	if cached := h.getCachedFinancialReport(r, cacheKey); cached != "" {
		writeFinancialReportJSON(w, cached)
		return
	}

	report, err := build()
	if err != nil {
		respondError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	payloadBytes, err := json.Marshal(map[string]interface{}{
		"success": true,
		"data":    report,
	})
	if err != nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    report,
		})
		return
	}

	payload := string(payloadBytes)
	if h.cache != nil {
		_ = h.cache.Set(r.Context(), cacheKey, payload, h.cacheTTL)
	}
	writeFinancialReportJSON(w, payload)
}

func (h *FinancialReportHandler) getCachedFinancialReport(r *http.Request, cacheKey string) string {
	if h.cache == nil {
		return ""
	}
	cached, err := h.cache.Get(r.Context(), cacheKey)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(cached)
}

func writeFinancialReportJSON(w http.ResponseWriter, payload string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(payload))
}

func parseFinancialDateRange(fromRaw, toRaw string) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	from := to.AddDate(0, 0, -29)

	if fromRaw != "" {
		parsed, err := time.Parse("2006-01-02", fromRaw)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		from = parsed.UTC()
	}
	if toRaw != "" {
		parsed, err := time.Parse("2006-01-02", toRaw)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		to = parsed.UTC()
	}
	return from, to, nil
}

func parseAsOfDate(raw string) (time.Time, error) {
	if raw == "" {
		now := time.Now().UTC()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}
