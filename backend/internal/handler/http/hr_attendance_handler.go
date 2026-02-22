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

// HRAttendanceHandler handles HR attendance management endpoints
type HRAttendanceHandler struct {
	attendanceRepo *repository.AttendanceRepository
	db             *sql.DB
}

// NewHRAttendanceHandler creates a new HR attendance handler
func NewHRAttendanceHandler(attendanceRepo *repository.AttendanceRepository, db *sql.DB) *HRAttendanceHandler {
	return &HRAttendanceHandler{attendanceRepo: attendanceRepo, db: db}
}

// ListAttendance handles GET /hr/attendance?start_date=...&end_date=...
func (h *HRAttendanceHandler) ListAttendance(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	log.Printf("[HR-Attendance] ListAttendance: tenantID=%d, restaurantID=%d, start=%s, end=%s\n",
		tenantID, restaurantID, startDateStr, endDateStr)

	// Default to current month if no dates provided
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

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

	records, err := h.attendanceRepo.ListAttendance(tenantID, restaurantID, startDate, endDate)
	if err != nil {
		log.Printf("[HR-Attendance] ListAttendance ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch attendance records",
		})
		return
	}

	log.Printf("[HR-Attendance] ListAttendance: found %d records\n", len(records))
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": records,
	})
}

// GetAttendance handles GET /hr/attendance/{id}
func (h *HRAttendanceHandler) GetAttendance(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid attendance ID",
		})
		return
	}

	record, err := h.attendanceRepo.GetAttendanceByID(tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[HR-Attendance] GetAttendance ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch attendance record",
		})
		return
	}

	if record == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "error": "Attendance record not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": record,
	})
}

// ClockIn handles POST /hr/attendance/clock-in
func (h *HRAttendanceHandler) ClockIn(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	var req struct {
		EmployeeID int    `json:"employee_id"`
		Notes      string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid request body",
		})
		return
	}

	if req.EmployeeID == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "employee_id is required",
		})
		return
	}

	log.Printf("[HR-Attendance] ClockIn: tenantID=%d, restaurantID=%d, employeeID=%d\n",
		tenantID, restaurantID, req.EmployeeID)

	id, err := h.attendanceRepo.ClockIn(tenantID, restaurantID, req.EmployeeID)
	if err != nil {
		log.Printf("[HR-Attendance] ClockIn ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to clock in: " + err.Error(),
		})
		return
	}

	log.Printf("[HR-Attendance] ClockIn SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":          id,
			"employee_id": req.EmployeeID,
			"status":      "present",
			"clock_in":    time.Now(),
		},
	})
}

// CreateAttendance handles POST /hr/attendance
func (h *HRAttendanceHandler) CreateAttendance(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req struct {
		EmployeeID     int    `json:"employee_id"`
		AttendanceDate string `json:"attendance_date"`
		ClockIn        string `json:"clock_in"`
		ClockOut       string `json:"clock_out"`
		Status         string `json:"status"`
		Notes          string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid request body",
		})
		return
	}

	if req.EmployeeID == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "employee_id is required",
		})
		return
	}

	// Parse attendance date
	attDate := time.Now()
	if req.AttendanceDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.AttendanceDate); err == nil {
			attDate = parsed
		}
	}

	// Parse clock times
	var clockIn, clockOut *time.Time
	if req.ClockIn != "" {
		if parsed, err := time.Parse(time.RFC3339, req.ClockIn); err == nil {
			clockIn = &parsed
		}
	}
	if req.ClockOut != "" {
		if parsed, err := time.Parse(time.RFC3339, req.ClockOut); err == nil {
			clockOut = &parsed
		}
	}

	if req.Status == "" {
		req.Status = "present"
	}

	attendance := &domain.Attendance{
		TenantID:       tenantID,
		RestaurantID:   restaurantID,
		EmployeeID:     req.EmployeeID,
		AttendanceDate: attDate,
		ClockIn:        clockIn,
		ClockOut:       clockOut,
		Status:         req.Status,
		Notes:          req.Notes,
		CreatedBy:      &userID,
	}

	id, err := h.attendanceRepo.CreateAttendance(attendance)
	if err != nil {
		log.Printf("[HR-Attendance] CreateAttendance ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to create attendance record: " + err.Error(),
		})
		return
	}

	attendance.ID = id
	log.Printf("[HR-Attendance] CreateAttendance SUCCESS: id=%d\n", id)

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true, "data": attendance,
	})
}

// UpdateAttendance handles PUT /hr/attendance/{id}
func (h *HRAttendanceHandler) UpdateAttendance(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid attendance ID",
		})
		return
	}

	// Get existing record
	existing, err := h.attendanceRepo.GetAttendanceByID(tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[HR-Attendance] UpdateAttendance ERROR fetching: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch attendance record",
		})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "error": "Attendance record not found",
		})
		return
	}

	var req struct {
		ClockIn    string `json:"clock_in"`
		ClockOut   string `json:"clock_out"`
		Status     string `json:"status"`
		Notes      string `json:"notes"`
		AdminNotes string `json:"admin_notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid request body",
		})
		return
	}

	// Apply updates
	if req.ClockIn != "" {
		if parsed, err := time.Parse(time.RFC3339, req.ClockIn); err == nil {
			existing.ClockIn = &parsed
		}
	}
	if req.ClockOut != "" {
		if parsed, err := time.Parse(time.RFC3339, req.ClockOut); err == nil {
			existing.ClockOut = &parsed
		}
	}
	if req.Status != "" {
		existing.Status = req.Status
	}
	if req.Notes != "" {
		existing.Notes = req.Notes
	}
	if req.AdminNotes != "" {
		existing.AdminNotes = req.AdminNotes
	}
	existing.UpdatedBy = &userID

	if err := h.attendanceRepo.UpdateAttendance(existing); err != nil {
		log.Printf("[HR-Attendance] UpdateAttendance ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to update attendance record",
		})
		return
	}

	log.Printf("[HR-Attendance] UpdateAttendance SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": existing,
	})
}

// DeleteAttendance handles DELETE /hr/attendance/{id}
func (h *HRAttendanceHandler) DeleteAttendance(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid attendance ID",
		})
		return
	}

	log.Printf("[HR-Attendance] DeleteAttendance: tenantID=%d, restaurantID=%d, id=%d\n", tenantID, restaurantID, id)

	_, err = h.db.Exec(
		`DELETE FROM attendance WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`,
		id, tenantID, restaurantID,
	)
	if err != nil {
		log.Printf("[HR-Attendance] DeleteAttendance ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to delete attendance record",
		})
		return
	}

	log.Printf("[HR-Attendance] DeleteAttendance SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
