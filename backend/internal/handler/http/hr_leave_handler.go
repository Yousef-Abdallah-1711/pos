package http

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

// HRLeaveHandler handles HR leave management endpoints
type HRLeaveHandler struct {
	leaveRepo *repository.LeaveRepository
	db        *sql.DB
}

// NewHRLeaveHandler creates a new HR leave handler
func NewHRLeaveHandler(leaveRepo *repository.LeaveRepository, db *sql.DB) *HRLeaveHandler {
	return &HRLeaveHandler{leaveRepo: leaveRepo, db: db}
}

// ListLeaves handles GET /hr/leaves?start_date=...&end_date=...
func (h *HRLeaveHandler) ListLeaves(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	log.Printf("[HR-Leave] ListLeaves: tenantID=%d, restaurantID=%d, start=%s, end=%s\n",
		tenantID, restaurantID, startDateStr, endDateStr)

	// Default to wide range if no dates provided
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0) // 1 year ago
	endDate := now.AddDate(1, 0, 0)    // 1 year from now

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed
		}
	}

	leaves, err := h.leaveRepo.ListLeaves(tenantID, restaurantID, startDate, endDate)
	if err != nil {
		log.Printf("[HR-Leave] ListLeaves ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch leave records",
		})
		return
	}

	log.Printf("[HR-Leave] ListLeaves: found %d records\n", len(leaves))
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": leaves,
	})
}

// GetLeave handles GET /hr/leaves/{id}
func (h *HRLeaveHandler) GetLeave(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid leave ID",
		})
		return
	}

	leave, err := h.leaveRepo.GetLeaveByID(tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[HR-Leave] GetLeave ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch leave record",
		})
		return
	}

	if leave == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "error": "Leave record not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": leave,
	})
}

// CreateLeave handles POST /hr/leaves
func (h *HRLeaveHandler) CreateLeave(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req struct {
		EmployeeID       int    `json:"employee_id"`
		StartDate        string `json:"start_date"`
		EndDate          string `json:"end_date"`
		TotalDays        int    `json:"total_days"`
		IsHalfDay        bool   `json:"is_half_day"`
		HalfDayPeriod    string `json:"half_day_period"`
		LeaveType        string `json:"leave_type"`
		LeaveCategory    string `json:"leave_category"`
		Reason           string `json:"reason"`
		ContactNumber    string `json:"contact_number"`
		ContactAddress   string `json:"contact_address"`
		EmergencyContact string `json:"emergency_contact"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid request body",
		})
		return
	}

	log.Printf("[HR-Leave] CreateLeave: tenantID=%d, restaurantID=%d, employeeID=%d, type=%s\n",
		tenantID, restaurantID, req.EmployeeID, req.LeaveType)

	if req.EmployeeID == 0 || req.StartDate == "" || req.EndDate == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "employee_id, start_date, and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid start_date format (expected YYYY-MM-DD)",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid end_date format (expected YYYY-MM-DD)",
		})
		return
	}

	if req.LeaveType == "" {
		req.LeaveType = "annual"
	}
	if req.LeaveCategory == "" {
		req.LeaveCategory = "paid"
	}

	// Calculate total days if not provided
	if req.TotalDays == 0 {
		req.TotalDays = int(endDate.Sub(startDate).Hours()/24) + 1
		if req.IsHalfDay {
			req.TotalDays = 1
		}
	}

	leave := &domain.Leave{
		TenantID:         tenantID,
		RestaurantID:     restaurantID,
		EmployeeID:       req.EmployeeID,
		StartDate:        startDate,
		EndDate:          endDate,
		TotalDays:        req.TotalDays,
		IsHalfDay:        req.IsHalfDay,
		HalfDayPeriod:    req.HalfDayPeriod,
		LeaveType:        req.LeaveType,
		LeaveCategory:    req.LeaveCategory,
		Reason:           req.Reason,
		ContactNumber:    req.ContactNumber,
		ContactAddress:   req.ContactAddress,
		EmergencyContact: req.EmergencyContact,
		Status:           "pending",
		CreatedBy:        &userID,
	}

	id, err := h.leaveRepo.CreateLeave(leave)
	if err != nil {
		log.Printf("[HR-Leave] CreateLeave ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to create leave request: " + err.Error(),
		})
		return
	}

	leave.ID = id
	log.Printf("[HR-Leave] CreateLeave SUCCESS: id=%d\n", id)

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true, "data": leave,
	})
}

// UpdateLeave handles PUT /hr/leaves/{id}
func (h *HRLeaveHandler) UpdateLeave(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid leave ID",
		})
		return
	}

	// Get existing
	existing, err := h.leaveRepo.GetLeaveByID(tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[HR-Leave] UpdateLeave ERROR fetching: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch leave record",
		})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "error": "Leave record not found",
		})
		return
	}

	var req struct {
		StartDate        string `json:"start_date"`
		EndDate          string `json:"end_date"`
		TotalDays        int    `json:"total_days"`
		IsHalfDay        bool   `json:"is_half_day"`
		HalfDayPeriod    string `json:"half_day_period"`
		LeaveType        string `json:"leave_type"`
		LeaveCategory    string `json:"leave_category"`
		Reason           string `json:"reason"`
		ContactNumber    string `json:"contact_number"`
		ContactAddress   string `json:"contact_address"`
		EmergencyContact string `json:"emergency_contact"`
		Status           string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid request body",
		})
		return
	}

	// Apply updates
	if req.StartDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			existing.StartDate = parsed
		}
	}
	if req.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			existing.EndDate = parsed
		}
	}
	if req.TotalDays > 0 {
		existing.TotalDays = req.TotalDays
	}
	existing.IsHalfDay = req.IsHalfDay
	if req.HalfDayPeriod != "" {
		existing.HalfDayPeriod = req.HalfDayPeriod
	}
	if req.LeaveType != "" {
		existing.LeaveType = req.LeaveType
	}
	if req.LeaveCategory != "" {
		existing.LeaveCategory = req.LeaveCategory
	}
	if req.Reason != "" {
		existing.Reason = req.Reason
	}
	if req.ContactNumber != "" {
		existing.ContactNumber = req.ContactNumber
	}
	if req.ContactAddress != "" {
		existing.ContactAddress = req.ContactAddress
	}
	if req.EmergencyContact != "" {
		existing.EmergencyContact = req.EmergencyContact
	}
	if req.Status != "" {
		existing.Status = req.Status
	}
	existing.UpdatedBy = &userID

	if err := h.leaveRepo.UpdateLeave(existing); err != nil {
		log.Printf("[HR-Leave] UpdateLeave ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to update leave record",
		})
		return
	}

	log.Printf("[HR-Leave] UpdateLeave SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": existing,
	})
}

// DeleteLeave handles DELETE /hr/leaves/{id}
func (h *HRLeaveHandler) DeleteLeave(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid leave ID",
		})
		return
	}

	log.Printf("[HR-Leave] DeleteLeave: tenantID=%d, restaurantID=%d, id=%d\n", tenantID, restaurantID, id)

	// Cancel the leave (soft delete via status change)
	err = h.leaveRepo.CancelLeave(tenantID, restaurantID, id, userID, "Deleted by user")
	if err != nil {
		// If cancel fails (maybe already cancelled), try hard delete
		_, err = h.db.Exec(
			`DELETE FROM leaves WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`,
			id, tenantID, restaurantID,
		)
		if err != nil {
			log.Printf("[HR-Leave] DeleteLeave ERROR: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "error": "Failed to delete leave record",
			})
			return
		}
	}

	log.Printf("[HR-Leave] DeleteLeave SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// ApproveLeave handles POST /hr/leaves/{id}/approve
func (h *HRLeaveHandler) ApproveLeave(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid leave ID",
		})
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	log.Printf("[HR-Leave] ApproveLeave: tenantID=%d, restaurantID=%d, id=%d, approvedBy=%d\n",
		tenantID, restaurantID, id, userID)

	if err := h.leaveRepo.ApproveLeave(tenantID, restaurantID, id, userID, req.Notes); err != nil {
		log.Printf("[HR-Leave] ApproveLeave ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to approve leave",
		})
		return
	}

	log.Printf("[HR-Leave] ApproveLeave SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// RejectLeave handles POST /hr/leaves/{id}/reject
func (h *HRLeaveHandler) RejectLeave(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid leave ID",
		})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	log.Printf("[HR-Leave] RejectLeave: tenantID=%d, restaurantID=%d, id=%d, rejectedBy=%d\n",
		tenantID, restaurantID, id, userID)

	if err := h.leaveRepo.RejectLeave(tenantID, restaurantID, id, userID, req.Reason); err != nil {
		log.Printf("[HR-Leave] RejectLeave ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to reject leave",
		})
		return
	}

	log.Printf("[HR-Leave] RejectLeave SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
