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

type TrainingHandler struct {
	repo *repository.TrainingRepository
}

const trainingHandlerTimeout = 10 * time.Second

func NewTrainingHandler(repo *repository.TrainingRepository) *TrainingHandler {
	return &TrainingHandler{repo: repo}
}

func (h *TrainingHandler) List(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel := context.WithTimeout(r.Context(), trainingHandlerTimeout)
	defer cancel()
	data, total, err := h.repo.ListPaged(ctx, tenantID, restaurantID, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load training records")
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

func (h *TrainingHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid training id")
		return
	}
	item, err := h.repo.GetByID(tenantID, restaurantID, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load training record")
		return
	}
	if item == nil {
		respondError(w, http.StatusNotFound, "Training record not found")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

func (h *TrainingHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	item, err := decodeTrainingPayload(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.TenantID = tenantID
	item.RestaurantID = restaurantID
	item.CreatedBy = &userID
	if item.Status == "" {
		item.Status = "scheduled"
	}
	id, err := h.repo.Create(&item)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ID = id
	respondJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": item})
}

func (h *TrainingHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid training id")
		return
	}
	item, err := decodeTrainingPayload(r)
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

func (h *TrainingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid training id")
		return
	}
	if err := h.repo.Delete(tenantID, restaurantID, id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *TrainingHandler) Expiring(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	data, err := h.repo.GetExpiringCertificates(tenantID, restaurantID, days)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load expiring certificates")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": data})
}

func decodeTrainingPayload(r *http.Request) (domain.TrainingRecord, error) {
	var body struct {
		EmployeeID        int    `json:"employee_id"`
		TrainingTitle     string `json:"training_title"`
		TrainingProvider  string `json:"training_provider"`
		TrainingType      string `json:"training_type"`
		StartDate         string `json:"start_date"`
		EndDate           string `json:"end_date"`
		Status            string `json:"status"`
		CertificateNumber string `json:"certificate_number"`
		CertificateExpiry string `json:"certificate_expiry"`
		Notes             string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return domain.TrainingRecord{}, err
	}
	if body.EmployeeID <= 0 || body.TrainingTitle == "" || body.TrainingType == "" || body.StartDate == "" {
		return domain.TrainingRecord{}, errValidation("employee_id, training_title, training_type, start_date are required")
	}
	startDate, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		return domain.TrainingRecord{}, errValidation("start_date must be YYYY-MM-DD")
	}
	var endDate *time.Time
	if body.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", body.EndDate)
		if err != nil {
			return domain.TrainingRecord{}, errValidation("end_date must be YYYY-MM-DD")
		}
		v := parsed.UTC()
		endDate = &v
	}
	var certificateExpiry *time.Time
	if body.CertificateExpiry != "" {
		parsed, err := time.Parse("2006-01-02", body.CertificateExpiry)
		if err != nil {
			return domain.TrainingRecord{}, errValidation("certificate_expiry must be YYYY-MM-DD")
		}
		v := parsed.UTC()
		certificateExpiry = &v
	}
	return domain.TrainingRecord{
		EmployeeID:        body.EmployeeID,
		TrainingTitle:     body.TrainingTitle,
		TrainingProvider:  body.TrainingProvider,
		TrainingType:      body.TrainingType,
		StartDate:         startDate.UTC(),
		EndDate:           endDate,
		Status:            body.Status,
		CertificateNumber: body.CertificateNumber,
		CertificateExpiry: certificateExpiry,
		Notes:             body.Notes,
	}, nil
}
