package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

// HRSalaryHandler handles HR salary/payroll management endpoints
type HRSalaryHandler struct {
	salaryRepo  *repository.SalaryRepository
	db          *sql.DB
	autoPosting usecase.AutoPostingService
}

// NewHRSalaryHandler creates a new HR salary handler
func NewHRSalaryHandler(salaryRepo *repository.SalaryRepository, db *sql.DB) *HRSalaryHandler {
	return &HRSalaryHandler{salaryRepo: salaryRepo, db: db}
}

// SetAutoPostingEngine links accounting auto-posting to salary payment updates.
func (h *HRSalaryHandler) SetAutoPostingEngine(engine usecase.AutoPostingService) {
	h.autoPosting = engine
}

// ListSalaries handles GET /hr/salaries?start_date=...&end_date=...
func (h *HRSalaryHandler) ListSalaries(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	log.Printf("[HR-Salary] ListSalaries: tenantID=%d, restaurantID=%d, start=%s, end=%s\n",
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

	salaries, err := h.salaryRepo.ListSalaries(tenantID, restaurantID, startDate, endDate)
	if err != nil {
		log.Printf("[HR-Salary] ListSalaries ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch salary records",
		})
		return
	}

	log.Printf("[HR-Salary] ListSalaries: found %d records\n", len(salaries))
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": salaries,
	})
}

// GetSalary handles GET /hr/salaries/{id}
func (h *HRSalaryHandler) GetSalary(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid salary ID",
		})
		return
	}

	salary, err := h.salaryRepo.GetSalaryByID(tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[HR-Salary] GetSalary ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch salary record",
		})
		return
	}

	if salary == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "error": "Salary record not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": salary,
	})
}

// CreateSalary handles POST /hr/salaries
func (h *HRSalaryHandler) CreateSalary(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req struct {
		EmployeeID         int     `json:"employee_id"`
		PayPeriodStart     string  `json:"pay_period_start"`
		PayPeriodEnd       string  `json:"pay_period_end"`
		BaseSalary         float64 `json:"base_salary"`
		Currency           string  `json:"currency"`
		OvertimeHours      float64 `json:"overtime_hours"`
		OvertimeRate       float64 `json:"overtime_rate"`
		Bonus              float64 `json:"bonus"`
		Commission         float64 `json:"commission"`
		Allowances         float64 `json:"allowances"`
		Tips               float64 `json:"tips"`
		OtherEarnings      float64 `json:"other_earnings"`
		Tax                float64 `json:"tax"`
		SocialInsurance    float64 `json:"social_insurance"`
		HealthInsurance    float64 `json:"health_insurance"`
		Pension            float64 `json:"pension"`
		LoanDeduction      float64 `json:"loan_deduction"`
		AdvanceDeduction   float64 `json:"advance_deduction"`
		OtherDeductions    float64 `json:"other_deductions"`
		DaysWorked         int     `json:"days_worked"`
		DaysAbsent         int     `json:"days_absent"`
		TotalHoursWorked   float64 `json:"total_hours_worked"`
		TotalOvertimeHours float64 `json:"total_overtime_hours"`
		PaymentMethod      string  `json:"payment_method"`
		Status             string  `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid request body",
		})
		return
	}

	log.Printf("[HR-Salary] CreateSalary: tenantID=%d, restaurantID=%d, employeeID=%d\n",
		tenantID, restaurantID, req.EmployeeID)

	if req.EmployeeID == 0 || req.PayPeriodStart == "" || req.PayPeriodEnd == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "employee_id, pay_period_start, and pay_period_end are required",
		})
		return
	}

	periodStart, err := time.Parse("2006-01-02", req.PayPeriodStart)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid pay_period_start format (expected YYYY-MM-DD)",
		})
		return
	}

	periodEnd, err := time.Parse("2006-01-02", req.PayPeriodEnd)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid pay_period_end format (expected YYYY-MM-DD)",
		})
		return
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}
	if req.PaymentMethod == "" {
		req.PaymentMethod = "bank_transfer"
	}
	if req.Status == "" {
		req.Status = "pending"
	}

	salary := &domain.Salary{
		TenantID:           tenantID,
		RestaurantID:       restaurantID,
		EmployeeID:         req.EmployeeID,
		PayPeriodStart:     periodStart,
		PayPeriodEnd:       periodEnd,
		BaseSalary:         req.BaseSalary,
		Currency:           req.Currency,
		OvertimeHours:      req.OvertimeHours,
		OvertimeRate:       req.OvertimeRate,
		Bonus:              req.Bonus,
		Commission:         req.Commission,
		Allowances:         req.Allowances,
		Tips:               req.Tips,
		OtherEarnings:      req.OtherEarnings,
		Tax:                req.Tax,
		SocialInsurance:    req.SocialInsurance,
		HealthInsurance:    req.HealthInsurance,
		Pension:            req.Pension,
		LoanDeduction:      req.LoanDeduction,
		AdvanceDeduction:   req.AdvanceDeduction,
		OtherDeductions:    req.OtherDeductions,
		DaysWorked:         req.DaysWorked,
		DaysAbsent:         req.DaysAbsent,
		TotalHoursWorked:   req.TotalHoursWorked,
		TotalOvertimeHours: req.TotalOvertimeHours,
		PaymentMethod:      req.PaymentMethod,
		Status:             req.Status,
		CreatedBy:          &userID,
	}

	id, err := h.salaryRepo.CreateSalary(salary)
	if err != nil {
		log.Printf("[HR-Salary] CreateSalary ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to create salary record: " + err.Error(),
		})
		return
	}

	salary.ID = id
	log.Printf("[HR-Salary] CreateSalary SUCCESS: id=%d\n", id)

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true, "data": salary,
	})
}

// UpdateSalary handles PUT /hr/salaries/{id}
func (h *HRSalaryHandler) UpdateSalary(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid salary ID",
		})
		return
	}

	// Get existing
	existing, err := h.salaryRepo.GetSalaryByID(tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[HR-Salary] UpdateSalary ERROR fetching: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch salary record",
		})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "error": "Salary record not found",
		})
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid request body",
		})
		return
	}

	// Apply updates
	if v, ok := req["base_salary"].(float64); ok {
		existing.BaseSalary = v
	}
	if v, ok := req["overtime_hours"].(float64); ok {
		existing.OvertimeHours = v
	}
	if v, ok := req["overtime_rate"].(float64); ok {
		existing.OvertimeRate = v
	}
	if v, ok := req["bonus"].(float64); ok {
		existing.Bonus = v
	}
	if v, ok := req["commission"].(float64); ok {
		existing.Commission = v
	}
	if v, ok := req["allowances"].(float64); ok {
		existing.Allowances = v
	}
	if v, ok := req["tips"].(float64); ok {
		existing.Tips = v
	}
	if v, ok := req["other_earnings"].(float64); ok {
		existing.OtherEarnings = v
	}
	if v, ok := req["tax"].(float64); ok {
		existing.Tax = v
	}
	if v, ok := req["social_insurance"].(float64); ok {
		existing.SocialInsurance = v
	}
	if v, ok := req["health_insurance"].(float64); ok {
		existing.HealthInsurance = v
	}
	if v, ok := req["pension"].(float64); ok {
		existing.Pension = v
	}
	if v, ok := req["loan_deduction"].(float64); ok {
		existing.LoanDeduction = v
	}
	if v, ok := req["advance_deduction"].(float64); ok {
		existing.AdvanceDeduction = v
	}
	if v, ok := req["other_deductions"].(float64); ok {
		existing.OtherDeductions = v
	}
	if v, ok := req["days_worked"].(float64); ok {
		existing.DaysWorked = int(v)
	}
	if v, ok := req["days_absent"].(float64); ok {
		existing.DaysAbsent = int(v)
	}
	if v, ok := req["total_hours_worked"].(float64); ok {
		existing.TotalHoursWorked = v
	}
	if v, ok := req["total_overtime_hours"].(float64); ok {
		existing.TotalOvertimeHours = v
	}
	if v, ok := req["payment_method"].(string); ok {
		existing.PaymentMethod = v
	}
	if v, ok := req["status"].(string); ok {
		existing.Status = v
	}
	existing.UpdatedBy = &userID

	if err := h.salaryRepo.UpdateSalary(existing); err != nil {
		log.Printf("[HR-Salary] UpdateSalary ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to update salary record",
		})
		return
	}

	log.Printf("[HR-Salary] UpdateSalary SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": existing,
	})
}

// DeleteSalary handles DELETE /hr/salaries/{id}
func (h *HRSalaryHandler) DeleteSalary(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid salary ID",
		})
		return
	}

	log.Printf("[HR-Salary] DeleteSalary: tenantID=%d, restaurantID=%d, id=%d\n", tenantID, restaurantID, id)

	_, err = h.db.Exec(
		`DELETE FROM salaries WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`,
		id, tenantID, restaurantID,
	)
	if err != nil {
		log.Printf("[HR-Salary] DeleteSalary ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to delete salary record",
		})
		return
	}

	log.Printf("[HR-Salary] DeleteSalary SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// MarkAsPaid handles POST /hr/salaries/{id}/mark-paid
func (h *HRSalaryHandler) MarkAsPaid(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid salary ID",
		})
		return
	}

	var req struct {
		PaymentReference string `json:"payment_reference"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid request body",
		})
		return
	}

	log.Printf("[HR-Salary] MarkAsPaid: tenantID=%d, restaurantID=%d, id=%d, paidBy=%d\n",
		tenantID, restaurantID, id, userID)

	if err := h.salaryRepo.MarkAsPaid(tenantID, restaurantID, id, userID, req.PaymentReference); err != nil {
		log.Printf("[HR-Salary] MarkAsPaid ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to mark salary as paid",
		})
		return
	}

	if h.autoPosting != nil {
		salary, err := h.salaryRepo.GetSalaryByID(tenantID, restaurantID, id)
		if err != nil {
			log.Printf("[AccountingAutoPost] salary reload failed tenant=%d restaurant=%d salary=%d err=%v", tenantID, restaurantID, id, err)
		} else if salary != nil {
			postedBy := &userID
			if userID <= 0 {
				postedBy = nil
			}
			if err := h.autoPosting.PostPayrollEntry(salary, postedBy); err != nil {
				log.Printf("[AccountingAutoPost] payroll posting failed tenant=%d restaurant=%d salary=%d err=%v", tenantID, restaurantID, id, err)
			}
		}
	}

	log.Printf("[HR-Salary] MarkAsPaid SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
