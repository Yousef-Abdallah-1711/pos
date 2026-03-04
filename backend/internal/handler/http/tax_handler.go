package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type TaxHandler struct {
	taxRepo *repository.TaxRepository
}

func NewTaxHandler(taxRepo *repository.TaxRepository) *TaxHandler {
	return &TaxHandler{taxRepo: taxRepo}
}

func (h *TaxHandler) ListTaxPeriods(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	periods, err := h.taxRepo.ListTaxPeriods(tenantID, restaurantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load tax periods")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    periods,
	})
}

func (h *TaxHandler) GetTaxPeriod(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid period id")
		return
	}

	period, err := h.taxRepo.GetTaxPeriodByID(tenantID, restaurantID, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Tax period not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    period,
	})
}

func (h *TaxHandler) CreateTaxPeriod(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	input, err := parseTaxPeriodInput(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.taxRepo.CreateTaxPeriod(tenantID, restaurantID, input)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    created,
	})
}

func (h *TaxHandler) UpdateTaxPeriod(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid period id")
		return
	}

	input, err := parseTaxPeriodInput(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.taxRepo.UpdateTaxPeriod(tenantID, restaurantID, id, input)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    updated,
	})
}

func (h *TaxHandler) DeleteTaxPeriod(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid period id")
		return
	}

	if err := h.taxRepo.DeleteTaxPeriod(tenantID, restaurantID, id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id": id,
		},
	})
}

func (h *TaxHandler) GetVATReport(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	from, to, err := parseTaxDateRange(r.URL.Query().Get("date_from"), r.URL.Query().Get("date_to"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid date range")
		return
	}

	summary, err := h.taxRepo.GetVATSummary(tenantID, restaurantID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load VAT report")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    summary,
	})
}

func (h *TaxHandler) CloseTaxPeriod(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	periodID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || periodID <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid period id")
		return
	}

	var body struct {
		Notes string `json:"notes"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	if err := h.taxRepo.CloseTaxPeriod(tenantID, restaurantID, periodID, userID); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":     periodID,
			"status": "closed",
			"notes":  body.Notes,
		},
	})
}

func parseTaxDateRange(fromRaw, toRaw string) (time.Time, time.Time, error) {
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

func parseTaxPeriodInput(r *http.Request) (repository.TaxPeriodInput, error) {
	var body struct {
		PeriodName  string `json:"period_name"`
		PeriodStart string `json:"period_start"`
		PeriodEnd   string `json:"period_end"`
		Notes       string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return repository.TaxPeriodInput{}, err
	}
	if body.PeriodName == "" || body.PeriodStart == "" || body.PeriodEnd == "" {
		return repository.TaxPeriodInput{}, errValidation("period_name, period_start and period_end are required")
	}

	start, err := time.Parse("2006-01-02", body.PeriodStart)
	if err != nil {
		return repository.TaxPeriodInput{}, errValidation("period_start must be YYYY-MM-DD")
	}
	end, err := time.Parse("2006-01-02", body.PeriodEnd)
	if err != nil {
		return repository.TaxPeriodInput{}, errValidation("period_end must be YYYY-MM-DD")
	}
	if end.Before(start) {
		return repository.TaxPeriodInput{}, errValidation("period_end must be >= period_start")
	}

	return repository.TaxPeriodInput{
		PeriodName:  body.PeriodName,
		PeriodStart: start.UTC(),
		PeriodEnd:   end.UTC(),
		Notes:       body.Notes,
	}, nil
}

type errValidation string

func (e errValidation) Error() string { return string(e) }
