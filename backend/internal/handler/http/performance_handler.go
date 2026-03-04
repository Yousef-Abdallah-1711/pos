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

type PerformanceHandler struct {
	repo *repository.PerformanceReviewRepository
}

func NewPerformanceHandler(repo *repository.PerformanceReviewRepository) *PerformanceHandler {
	return &PerformanceHandler{repo: repo}
}

func (h *PerformanceHandler) List(w http.ResponseWriter, r *http.Request) {
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
	offset := (page - 1) * pageSize

	items, total, err := h.repo.ListPaginated(tenantID, restaurantID, pageSize, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load performance reviews")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    items,
		"pagination": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *PerformanceHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid review id")
		return
	}

	item, err := h.repo.GetByID(tenantID, restaurantID, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load performance review")
		return
	}
	if item == nil {
		respondError(w, http.StatusNotFound, "Performance review not found")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

func (h *PerformanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	req, err := decodePerformanceReviewPayload(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	req.TenantID = tenantID
	req.RestaurantID = restaurantID
	req.CreatedBy = &userID
	if req.Status == "" {
		req.Status = "draft"
	}

	id, err := h.repo.Create(&req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create performance review")
		return
	}
	req.ID = id
	respondJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": req})
}

func (h *PerformanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid review id")
		return
	}

	req, err := decodePerformanceReviewPayload(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = id
	req.TenantID = tenantID
	req.RestaurantID = restaurantID
	req.UpdatedBy = &userID

	if err := h.repo.Update(tenantID, restaurantID, &req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.repo.GetByID(tenantID, restaurantID, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load updated performance review")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": updated})
}

func (h *PerformanceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid review id")
		return
	}

	if err := h.repo.Delete(tenantID, restaurantID, id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]int{"id": id}})
}

func (h *PerformanceHandler) Acknowledge(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid review id")
		return
	}
	if err := h.repo.Acknowledge(tenantID, restaurantID, id, userID); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]string{"status": "acknowledged"}})
}

func (h *PerformanceHandler) Complete(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid review id")
		return
	}
	if err := h.repo.Complete(tenantID, restaurantID, id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]string{"status": "completed"}})
}

func (h *PerformanceHandler) EmployeeHistory(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	employeeID, err := strconv.Atoi(r.PathValue("employeeId"))
	if err != nil || employeeID <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid employee id")
		return
	}
	items, average, err := h.repo.GetEmployeeHistory(tenantID, restaurantID, employeeID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load employee performance history")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"rows":           items,
			"average_rating": average,
			"employee_id":    employeeID,
			"review_count":   len(items),
		},
	})
}

func decodePerformanceReviewPayload(r *http.Request) (domain.PerformanceReview, error) {
	var body struct {
		EmployeeID        int     `json:"employee_id"`
		ReviewerID        int     `json:"reviewer_id"`
		ReviewPeriodStart string  `json:"review_period_start"`
		ReviewPeriodEnd   string  `json:"review_period_end"`
		ReviewDate        string  `json:"review_date"`
		GoalsAchievement  float64 `json:"goals_achievement"`
		TechnicalSkills   float64 `json:"technical_skills"`
		Communication     float64 `json:"communication"`
		Leadership        float64 `json:"leadership"`
		Strengths         string  `json:"strengths"`
		AreasToImprove    string  `json:"areas_to_improve"`
		ActionPlan        string  `json:"action_plan"`
		Status            string  `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return domain.PerformanceReview{}, err
	}
	if body.EmployeeID <= 0 || body.ReviewerID <= 0 {
		return domain.PerformanceReview{}, errValidation("employee_id and reviewer_id are required")
	}
	start, err := time.Parse("2006-01-02", body.ReviewPeriodStart)
	if err != nil {
		return domain.PerformanceReview{}, errValidation("review_period_start must be YYYY-MM-DD")
	}
	end, err := time.Parse("2006-01-02", body.ReviewPeriodEnd)
	if err != nil {
		return domain.PerformanceReview{}, errValidation("review_period_end must be YYYY-MM-DD")
	}
	reviewDate, err := time.Parse("2006-01-02", body.ReviewDate)
	if err != nil {
		return domain.PerformanceReview{}, errValidation("review_date must be YYYY-MM-DD")
	}

	return domain.PerformanceReview{
		EmployeeID:        body.EmployeeID,
		ReviewerID:        body.ReviewerID,
		ReviewPeriodStart: start.UTC(),
		ReviewPeriodEnd:   end.UTC(),
		ReviewDate:        reviewDate.UTC(),
		GoalsAchievement:  body.GoalsAchievement,
		TechnicalSkills:   body.TechnicalSkills,
		Communication:     body.Communication,
		Leadership:        body.Leadership,
		Strengths:         body.Strengths,
		AreasToImprove:    body.AreasToImprove,
		ActionPlan:        body.ActionPlan,
		Status:            body.Status,
	}, nil
}
