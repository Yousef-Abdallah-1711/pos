package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

// HREmployeeHandler handles HR employee management endpoints
type HREmployeeHandler struct {
	employeeRepo *repository.EmployeeRepository
}

// NewHREmployeeHandler creates a new HR employee handler
func NewHREmployeeHandler(employeeRepo *repository.EmployeeRepository) *HREmployeeHandler {
	return &HREmployeeHandler{employeeRepo: employeeRepo}
}

// ListEmployees handles GET /hr/employees
func (h *HREmployeeHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	log.Printf("[HR-Employee] ListEmployees: tenantID=%d, restaurantID=%d\n", tenantID, restaurantID)

	employees, err := h.employeeRepo.ListEmployees(tenantID, restaurantID)
	if err != nil {
		log.Printf("[HR-Employee] ListEmployees ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch employees",
		})
		return
	}

	log.Printf("[HR-Employee] ListEmployees: found %d employees\n", len(employees))
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": employees,
	})
}

// GetEmployee handles GET /hr/employees/{id}
func (h *HREmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid employee ID",
		})
		return
	}

	employee, err := h.employeeRepo.GetEmployeeByID(tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[HR-Employee] GetEmployee ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch employee",
		})
		return
	}

	if employee == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "error": "Employee not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": employee,
	})
}

// CreateEmployee handles POST /hr/employees
func (h *HREmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req struct {
		EmployeeCode    string  `json:"employee_code"`
		FirstName       string  `json:"first_name"`
		LastName        string  `json:"last_name"`
		FirstNameAr     string  `json:"first_name_ar"`
		LastNameAr      string  `json:"last_name_ar"`
		Email           string  `json:"email"`
		Phone           string  `json:"phone"`
		Gender          string  `json:"gender"`
		NationalID      string  `json:"national_id"`
		DateOfBirth     string  `json:"date_of_birth"`
		Address         string  `json:"address"`
		City            string  `json:"city"`
		State           string  `json:"state"`
		PostalCode      string  `json:"postal_code"`
		Country         string  `json:"country"`
		HireDate        string  `json:"hire_date"`
		EmploymentType  string  `json:"employment_type"`
		Position        string  `json:"position"`
		Department      string  `json:"department"`
		BaseSalary      float64 `json:"base_salary"`
		SalaryCurrency  string  `json:"salary_currency"`
		PaymentFrequency string `json:"payment_frequency"`
		WorkingHoursPerWeek float64 `json:"working_hours_per_week"`
		ShiftType       string  `json:"shift_type"`
		EmergencyContactName  string `json:"emergency_contact_name"`
		EmergencyContactPhone string `json:"emergency_contact_phone"`
		EmergencyContactRelationship string `json:"emergency_contact_relationship"`
		Notes           string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[HR-Employee] CreateEmployee: invalid request body: %v\n", err)
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid request body",
		})
		return
	}

	log.Printf("[HR-Employee] CreateEmployee: tenantID=%d, restaurantID=%d, code=%s, name=%s %s\n",
		tenantID, restaurantID, req.EmployeeCode, req.FirstName, req.LastName)

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "first_name, last_name, and email are required",
		})
		return
	}

	// Parse dates
	var hireDate time.Time
	if req.HireDate != "" {
		var err error
		hireDate, err = time.Parse("2006-01-02", req.HireDate)
		if err != nil {
			hireDate = time.Now()
		}
	} else {
		hireDate = time.Now()
	}

	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			dateOfBirth = &dob
		}
	}

	// Validate gender - database has check constraint for 'male', 'female', 'other'
	if req.Gender != "" && req.Gender != "male" && req.Gender != "female" && req.Gender != "other" {
		req.Gender = ""
	}

	// Set defaults
	if req.EmploymentType == "" {
		req.EmploymentType = "full_time"
	}
	if req.SalaryCurrency == "" {
		req.SalaryCurrency = "USD"
	}
	if req.PaymentFrequency == "" {
		req.PaymentFrequency = "monthly"
	}
	if req.ShiftType == "" {
		req.ShiftType = "day"
	}
	if req.Position == "" {
		req.Position = "Staff"
	}

	// Generate employee code if not provided
	if req.EmployeeCode == "" {
		req.EmployeeCode = "EMP-" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	}

	employee := &domain.Employee{
		TenantID:            tenantID,
		RestaurantID:        restaurantID,
		EmployeeCode:        req.EmployeeCode,
		FirstName:           req.FirstName,
		LastName:            req.LastName,
		FirstNameAr:         req.FirstNameAr,
		LastNameAr:          req.LastNameAr,
		Email:               req.Email,
		Phone:               req.Phone,
		DateOfBirth:         dateOfBirth,
		Gender:              req.Gender,
		NationalID:          req.NationalID,
		Address:             req.Address,
		City:                req.City,
		State:               req.State,
		PostalCode:          req.PostalCode,
		Country:             req.Country,
		HireDate:            hireDate,
		EmploymentType:      req.EmploymentType,
		EmploymentStatus:    "active",
		Position:            req.Position,
		Department:          req.Department,
		BaseSalary:          req.BaseSalary,
		SalaryCurrency:      req.SalaryCurrency,
		PaymentFrequency:    req.PaymentFrequency,
		WorkingHoursPerWeek: req.WorkingHoursPerWeek,
		ShiftType:           req.ShiftType,
		EmergencyContactName:         req.EmergencyContactName,
		EmergencyContactPhone:        req.EmergencyContactPhone,
		EmergencyContactRelationship: req.EmergencyContactRelationship,
		Notes:               req.Notes,
		IsActive:            true,
		CreatedBy:           &userID,
	}

	id, err := h.employeeRepo.CreateEmployee(employee)
	if err != nil {
		log.Printf("[HR-Employee] CreateEmployee ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to create employee: " + err.Error(),
		})
		return
	}

	employee.ID = id
	log.Printf("[HR-Employee] CreateEmployee SUCCESS: id=%d\n", id)

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true, "data": employee,
	})
}

// UpdateEmployee handles PUT /hr/employees/{id}
func (h *HREmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid employee ID",
		})
		return
	}

	// Get existing employee first
	existing, err := h.employeeRepo.GetEmployeeByID(tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[HR-Employee] UpdateEmployee ERROR fetching: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to fetch employee",
		})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "error": "Employee not found",
		})
		return
	}

	// Decode update request
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid request body",
		})
		return
	}

	// Apply updates to existing employee
	if v, ok := req["first_name"].(string); ok && v != "" {
		existing.FirstName = v
	}
	if v, ok := req["last_name"].(string); ok && v != "" {
		existing.LastName = v
	}
	if v, ok := req["first_name_ar"].(string); ok {
		existing.FirstNameAr = v
	}
	if v, ok := req["last_name_ar"].(string); ok {
		existing.LastNameAr = v
	}
	if v, ok := req["email"].(string); ok && v != "" {
		existing.Email = v
	}
	if v, ok := req["phone"].(string); ok {
		existing.Phone = v
	}
	if v, ok := req["position"].(string); ok {
		existing.Position = v
	}
	if v, ok := req["department"].(string); ok {
		existing.Department = v
	}
	if v, ok := req["employment_type"].(string); ok {
		existing.EmploymentType = v
	}
	if v, ok := req["employment_status"].(string); ok {
		existing.EmploymentStatus = v
	}
	if v, ok := req["base_salary"].(float64); ok {
		existing.BaseSalary = v
	}
	if v, ok := req["salary_currency"].(string); ok {
		existing.SalaryCurrency = v
	}
	if v, ok := req["payment_frequency"].(string); ok {
		existing.PaymentFrequency = v
	}
	if v, ok := req["shift_type"].(string); ok {
		existing.ShiftType = v
	}
	if v, ok := req["working_hours_per_week"].(float64); ok {
		existing.WorkingHoursPerWeek = v
	}
	if v, ok := req["gender"].(string); ok {
		existing.Gender = v
	}
	if v, ok := req["national_id"].(string); ok {
		existing.NationalID = v
	}
	if v, ok := req["address"].(string); ok {
		existing.Address = v
	}
	if v, ok := req["city"].(string); ok {
		existing.City = v
	}
	if v, ok := req["state"].(string); ok {
		existing.State = v
	}
	if v, ok := req["postal_code"].(string); ok {
		existing.PostalCode = v
	}
	if v, ok := req["country"].(string); ok {
		existing.Country = v
	}
	if v, ok := req["notes"].(string); ok {
		existing.Notes = v
	}
	if v, ok := req["emergency_contact_name"].(string); ok {
		existing.EmergencyContactName = v
	}
	if v, ok := req["emergency_contact_phone"].(string); ok {
		existing.EmergencyContactPhone = v
	}
	if v, ok := req["emergency_contact_relationship"].(string); ok {
		existing.EmergencyContactRelationship = v
	}

	existing.UpdatedBy = &userID

	if err := h.employeeRepo.UpdateEmployee(existing); err != nil {
		log.Printf("[HR-Employee] UpdateEmployee ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to update employee",
		})
		return
	}

	log.Printf("[HR-Employee] UpdateEmployee SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": existing,
	})
}

// DeleteEmployee handles DELETE /hr/employees/{id}
func (h *HREmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "error": "Invalid employee ID",
		})
		return
	}

	log.Printf("[HR-Employee] DeleteEmployee: tenantID=%d, restaurantID=%d, id=%d\n", tenantID, restaurantID, id)

	if err := h.employeeRepo.DeleteEmployee(tenantID, restaurantID, id); err != nil {
		log.Printf("[HR-Employee] DeleteEmployee ERROR: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "error": "Failed to delete employee",
		})
		return
	}

	log.Printf("[HR-Employee] DeleteEmployee SUCCESS: id=%d\n", id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
