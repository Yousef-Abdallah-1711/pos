package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type ContractHandler struct {
	repo *repository.ContractRepository
}

const contractHandlerTimeout = 10 * time.Second

func NewContractHandler(repo *repository.ContractRepository) *ContractHandler {
	return &ContractHandler{repo: repo}
}

func (h *ContractHandler) List(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel := context.WithTimeout(r.Context(), contractHandlerTimeout)
	defer cancel()
	data, total, err := h.repo.ListPaged(ctx, tenantID, restaurantID, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load contracts")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
		"pagination": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *ContractHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid contract id")
		return
	}
	item, err := h.repo.GetByID(tenantID, restaurantID, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load contract")
		return
	}
	if item == nil {
		respondError(w, http.StatusNotFound, "Contract not found")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

func (h *ContractHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	item, err := decodeContractPayload(r)
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

func (h *ContractHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid contract id")
		return
	}
	item, err := decodeContractPayload(r)
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

func (h *ContractHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid contract id")
		return
	}
	if err := h.repo.Delete(tenantID, restaurantID, id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *ContractHandler) Renew(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid contract id")
		return
	}
	var body struct {
		StartDate    string  `json:"start_date"`
		EndDate      *string `json:"end_date"`
		SalaryAmount float64 `json:"salary_amount"`
		Currency     string  `json:"currency"`
	}
	if err := decodeJSONStrict(r, &body); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	startDate, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "start_date must be YYYY-MM-DD")
		return
	}
	var endDate *time.Time
	if body.EndDate != nil && *body.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", *body.EndDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "end_date must be YYYY-MM-DD")
			return
		}
		v := parsed.UTC()
		endDate = &v
	}
	if body.Currency == "" {
		body.Currency = "EGP"
	}

	newID, err := h.repo.Renew(tenantID, restaurantID, id, startDate.UTC(), endDate, body.SalaryAmount, body.Currency, userID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]int{"id": newID}})
}

func (h *ContractHandler) Expiring(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	data, err := h.repo.GetExpiring(tenantID, restaurantID, days)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load expiring contracts")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": data})
}

func decodeContractPayload(r *http.Request) (domain.EmploymentContract, error) {
	var body struct {
		EmployeeID     int     `json:"employee_id"`
		ContractType   string  `json:"contract_type"`
		ContractNumber string  `json:"contract_number"`
		StartDate      string  `json:"start_date"`
		EndDate        *string `json:"end_date"`
		SalaryAmount   float64 `json:"salary_amount"`
		Currency       string  `json:"currency"`
		Status         string  `json:"status"`
		Notes          string  `json:"notes"`
	}
	if err := decodeJSONStrict(r, &body); err != nil {
		return domain.EmploymentContract{}, err
	}
	if body.EmployeeID <= 0 || body.ContractType == "" || body.ContractNumber == "" || body.StartDate == "" {
		return domain.EmploymentContract{}, errValidation("employee_id, contract_type, contract_number, start_date are required")
	}
	startDate, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		return domain.EmploymentContract{}, errValidation("start_date must be YYYY-MM-DD")
	}
	var endDate *time.Time
	if body.EndDate != nil && *body.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", *body.EndDate)
		if err != nil {
			return domain.EmploymentContract{}, errValidation("end_date must be YYYY-MM-DD")
		}
		v := parsed.UTC()
		endDate = &v
	}
	if body.Currency == "" {
		body.Currency = "EGP"
	}
	return domain.EmploymentContract{
		EmployeeID:     body.EmployeeID,
		ContractType:   body.ContractType,
		ContractNumber: body.ContractNumber,
		StartDate:      startDate.UTC(),
		EndDate:        endDate,
		SalaryAmount:   body.SalaryAmount,
		Currency:       body.Currency,
		Status:         body.Status,
		Notes:          body.Notes,
	}, nil
}
