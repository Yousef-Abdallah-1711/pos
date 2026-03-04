package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type ResignationHandler struct {
	repo *repository.ResignationRepository
}

const resignationHandlerTimeout = 10 * time.Second

func NewResignationHandler(repo *repository.ResignationRepository) *ResignationHandler {
	return &ResignationHandler{repo: repo}
}

func (h *ResignationHandler) List(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel := context.WithTimeout(r.Context(), resignationHandlerTimeout)
	defer cancel()
	data, total, err := h.repo.ListPaged(ctx, tenantID, restaurantID, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load resignations")
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

func (h *ResignationHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid resignation id")
		return
	}
	item, err := h.repo.GetByID(tenantID, restaurantID, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load resignation")
		return
	}
	if item == nil {
		respondError(w, http.StatusNotFound, "Resignation not found")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

func (h *ResignationHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	item, err := decodeResignationPayload(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.TenantID = tenantID
	item.RestaurantID = restaurantID
	item.CreatedBy = &userID
	item.UpdatedBy = &userID
	if item.NoticePeriodDays == 0 {
		item.NoticePeriodDays = int(item.LastWorkingDate.Sub(item.ResignationDate).Hours() / 24)
	}
	item.SettlementAmount = h.repo.CalculateSettlement(&item)
	id, err := h.repo.Create(&item)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ID = id
	respondJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": item})
}

func (h *ResignationHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid resignation id")
		return
	}
	item, err := decodeResignationPayload(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ID = id
	item.TenantID = tenantID
	item.RestaurantID = restaurantID
	item.UpdatedBy = &userID
	if item.NoticePeriodDays == 0 {
		item.NoticePeriodDays = int(item.LastWorkingDate.Sub(item.ResignationDate).Hours() / 24)
	}
	item.SettlementAmount = h.repo.CalculateSettlement(&item)
	if err := h.repo.Update(tenantID, restaurantID, &item); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *ResignationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid resignation id")
		return
	}
	if err := h.repo.Delete(tenantID, restaurantID, id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *ResignationHandler) Acknowledge(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid resignation id")
		return
	}
	if err := h.repo.Acknowledge(tenantID, restaurantID, id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]string{"status": "acknowledged"}})
}

func (h *ResignationHandler) Complete(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid resignation id")
		return
	}
	item, err := h.repo.GetByID(tenantID, restaurantID, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load resignation")
		return
	}
	if item == nil {
		respondError(w, http.StatusNotFound, "Resignation not found")
		return
	}
	settlement := h.repo.CalculateSettlement(item)
	if err := h.repo.Complete(tenantID, restaurantID, id, settlement); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]interface{}{"status": "completed", "settlement_amount": settlement}})
}

func (h *ResignationHandler) CalculateSettlement(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid resignation id")
		return
	}
	item, err := h.repo.GetByID(tenantID, restaurantID, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load resignation")
		return
	}
	if item == nil {
		respondError(w, http.StatusNotFound, "Resignation not found")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]float64{
			"settlement_amount": h.repo.CalculateSettlement(item),
		},
	})
}

func decodeResignationPayload(r *http.Request) (domain.Resignation, error) {
	var body struct {
		EmployeeID         int     `json:"employee_id"`
		ResignationDate    string  `json:"resignation_date"`
		LastWorkingDate    string  `json:"last_working_date"`
		Reason             string  `json:"reason"`
		Status             string  `json:"status"`
		NoticePeriodDays   int     `json:"notice_period_days"`
		PendingSalary      float64 `json:"pending_salary"`
		UnusedLeaveBalance float64 `json:"unused_leave_balance"`
		SettlementAmount   float64 `json:"settlement_amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return domain.Resignation{}, err
	}
	if body.EmployeeID <= 0 || body.ResignationDate == "" || body.LastWorkingDate == "" || body.Reason == "" {
		return domain.Resignation{}, errValidation("employee_id, resignation_date, last_working_date, reason are required")
	}
	resignationDate, err := time.Parse("2006-01-02", body.ResignationDate)
	if err != nil {
		return domain.Resignation{}, errValidation("resignation_date must be YYYY-MM-DD")
	}
	lastWorkingDate, err := time.Parse("2006-01-02", body.LastWorkingDate)
	if err != nil {
		return domain.Resignation{}, errValidation("last_working_date must be YYYY-MM-DD")
	}
	if body.Status == "" {
		body.Status = "submitted"
	}
	return domain.Resignation{
		EmployeeID:         body.EmployeeID,
		ResignationDate:    resignationDate.UTC(),
		LastWorkingDate:    lastWorkingDate.UTC(),
		Reason:             body.Reason,
		Status:             body.Status,
		NoticePeriodDays:   body.NoticePeriodDays,
		PendingSalary:      body.PendingSalary,
		UnusedLeaveBalance: body.UnusedLeaveBalance,
		SettlementAmount:   body.SettlementAmount,
	}, nil
}
