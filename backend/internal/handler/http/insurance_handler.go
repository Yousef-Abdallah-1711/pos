package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type InsuranceHandler struct {
	repo *repository.InsuranceRepository
}

func NewInsuranceHandler(repo *repository.InsuranceRepository) *InsuranceHandler {
	return &InsuranceHandler{repo: repo}
}

func (h *InsuranceHandler) List(w http.ResponseWriter, r *http.Request) {
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
	data, err := h.repo.List(tenantID, restaurantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load insurance records")
		return
	}
	total := len(data)
	offset := (page - 1) * pageSize
	if offset > total {
		offset = total
	}
	end := offset + pageSize
	if end > total {
		end = total
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data[offset:end],
		"pagination": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *InsuranceHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid insurance id")
		return
	}
	item, err := h.repo.GetByID(tenantID, restaurantID, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load insurance record")
		return
	}
	if item == nil {
		respondError(w, http.StatusNotFound, "Insurance record not found")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

func (h *InsuranceHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	item, err := decodeInsurancePayload(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.TenantID = tenantID
	item.RestaurantID = restaurantID
	item.CreatedBy = &userID
	if item.Status == "" {
		item.Status = "active"
	}
	id, err := h.repo.Create(&item)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ID = id
	respondJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": item})
}

func (h *InsuranceHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid insurance id")
		return
	}
	item, err := decodeInsurancePayload(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ID = id
	item.TenantID = tenantID
	item.RestaurantID = restaurantID
	item.UpdatedBy = &userID
	if err := h.repo.Update(tenantID, restaurantID, &item); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *InsuranceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid insurance id")
		return
	}
	if err := h.repo.Delete(tenantID, restaurantID, id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *InsuranceHandler) Expiring(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	data, err := h.repo.GetExpiring(tenantID, restaurantID, days)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load expiring policies")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": data})
}

func decodeInsurancePayload(r *http.Request) (domain.InsuranceTracking, error) {
	var body struct {
		EmployeeID     int     `json:"employee_id"`
		PolicyType     string  `json:"policy_type"`
		ProviderName   string  `json:"provider_name"`
		PolicyNumber   string  `json:"policy_number"`
		CoverageStart  string  `json:"coverage_start"`
		CoverageEnd    string  `json:"coverage_end"`
		MonthlyPremium float64 `json:"monthly_premium"`
		Status         string  `json:"status"`
		Notes          string  `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return domain.InsuranceTracking{}, err
	}
	if body.EmployeeID <= 0 || body.PolicyType == "" || body.ProviderName == "" || body.PolicyNumber == "" || body.CoverageStart == "" || body.CoverageEnd == "" {
		return domain.InsuranceTracking{}, errValidation("employee_id, policy_type, provider_name, policy_number, coverage_start, coverage_end are required")
	}
	startDate, err := time.Parse("2006-01-02", body.CoverageStart)
	if err != nil {
		return domain.InsuranceTracking{}, errValidation("coverage_start must be YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", body.CoverageEnd)
	if err != nil {
		return domain.InsuranceTracking{}, errValidation("coverage_end must be YYYY-MM-DD")
	}
	return domain.InsuranceTracking{
		EmployeeID:     body.EmployeeID,
		PolicyType:     body.PolicyType,
		ProviderName:   body.ProviderName,
		PolicyNumber:   body.PolicyNumber,
		CoverageStart:  startDate.UTC(),
		CoverageEnd:    endDate.UTC(),
		MonthlyPremium: body.MonthlyPremium,
		Status:         body.Status,
		Notes:          body.Notes,
	}, nil
}
