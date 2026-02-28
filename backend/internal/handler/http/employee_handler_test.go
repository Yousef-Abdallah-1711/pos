//go:build legacy
// +build legacy

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
)

// MockEmployeeRepository is a mock implementation for testing
type MockEmployeeRepository struct {
	employees map[int]*domain.Employee
	nextID    int
}

// NewMockEmployeeRepository creates a new mock repository
func NewMockEmployeeRepository() *MockEmployeeRepository {
	return &MockEmployeeRepository{
		employees: make(map[int]*domain.Employee),
		nextID:    1,
	}
}

// ListEmployees returns all employees
func (m *MockEmployeeRepository) ListEmployees(tenantID, restaurantID int) ([]domain.Employee, error) {
	var result []domain.Employee
	for _, emp := range m.employees {
		if emp.TenantID == tenantID && emp.RestaurantID == restaurantID {
			result = append(result, *emp)
		}
	}
	return result, nil
}

// GetEmployeesByStatus returns employees by status
func (m *MockEmployeeRepository) GetEmployeesByStatus(tenantID, restaurantID int, status string) ([]domain.Employee, error) {
	var result []domain.Employee
	for _, emp := range m.employees {
		if emp.TenantID == tenantID && emp.RestaurantID == restaurantID && emp.EmploymentStatus == status {
			result = append(result, *emp)
		}
	}
	return result, nil
}

// GetEmployeeByID returns a single employee
func (m *MockEmployeeRepository) GetEmployeeByID(tenantID, restaurantID int, id int) (*domain.Employee, error) {
	if emp, exists := m.employees[id]; exists {
		if emp.TenantID == tenantID && emp.RestaurantID == restaurantID {
			return emp, nil
		}
	}
	return nil, nil
}

// CreateEmployee creates a new employee
func (m *MockEmployeeRepository) CreateEmployee(emp *domain.Employee) (int, error) {
	id := m.nextID
	emp.ID = id
	m.employees[id] = emp
	m.nextID++
	return id, nil
}

// UpdateEmployee updates an existing employee
func (m *MockEmployeeRepository) UpdateEmployee(emp *domain.Employee) error {
	if _, exists := m.employees[emp.ID]; exists {
		m.employees[emp.ID] = emp
		return nil
	}
	return nil
}

// DeleteEmployee soft deletes an employee
func (m *MockEmployeeRepository) DeleteEmployee(tenantID, restaurantID int, id int) error {
	if emp, exists := m.employees[id]; exists {
		if emp.TenantID == tenantID && emp.RestaurantID == restaurantID {
			emp.IsActive = false
		}
	}
	return nil
}

// Helper function to create a request with claims in context
func createRequestWithClaims(method, url string, body interface{}, claims *jwt.Claims) *http.Request {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, url, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Add claims to context
	if claims != nil {
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
		req = req.WithContext(ctx)
	}

	return req
}

// Test CreateEmployee with minimal fields (happy path)
func TestCreateEmployeeMinimalFields(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	inputBody := map[string]interface{}{
		"first_name": "John",
		"email":      "john@test.com",
	}

	req := createRequestWithClaims("POST", "/api/v1/hr/employees", inputBody, claims)
	w := httptest.NewRecorder()

	handler.CreateEmployee(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response domain.Employee
	json.NewDecoder(w.Body).Decode(&response)

	if response.FirstName != "John" {
		t.Errorf("Expected FirstName 'John', got '%s'", response.FirstName)
	}
	if response.Email != "john@test.com" {
		t.Errorf("Expected Email 'john@test.com', got '%s'", response.Email)
	}
	if response.EmployeeCode == "" {
		t.Errorf("Expected auto-generated EmployeeCode, got empty string")
	}
	if response.Position != "Staff" {
		t.Errorf("Expected Position 'Staff', got '%s'", response.Position)
	}
	if response.EmploymentStatus != "active" {
		t.Errorf("Expected EmploymentStatus 'active', got '%s'", response.EmploymentStatus)
	}
}

// Test CreateEmployee without authentication (nil claims)
func TestCreateEmployeeNoAuth(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	inputBody := map[string]interface{}{
		"first_name": "John",
		"email":      "john@test.com",
	}

	req := createRequestWithClaims("POST", "/api/v1/hr/employees", inputBody, nil)
	w := httptest.NewRecorder()

	handler.CreateEmployee(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test CreateEmployee without required FirstName
func TestCreateEmployeeNoFirstName(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	inputBody := map[string]interface{}{
		"email": "john@test.com",
	}

	req := createRequestWithClaims("POST", "/api/v1/hr/employees", inputBody, claims)
	w := httptest.NewRecorder()

	handler.CreateEmployee(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Test CreateEmployee without required Email
func TestCreateEmployeeNoEmail(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	inputBody := map[string]interface{}{
		"first_name": "John",
	}

	req := createRequestWithClaims("POST", "/api/v1/hr/employees", inputBody, claims)
	w := httptest.NewRecorder()

	handler.CreateEmployee(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Test CreateEmployee with all fields
func TestCreateEmployeeFullData(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	hireDate := time.Now()
	baseSalary := 50000.0

	inputBody := map[string]interface{}{
		"employee_code": "EMP-001",
		"first_name":    "John",
		"last_name":     "Doe",
		"email":         "john@test.com",
		"phone":         "1234567890",
		"gender":        "M",
		"position":      "Manager",
		"hire_date":     hireDate.Format(time.RFC3339),
		"base_salary":   baseSalary,
	}

	req := createRequestWithClaims("POST", "/api/v1/hr/employees", inputBody, claims)
	w := httptest.NewRecorder()

	handler.CreateEmployee(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response domain.Employee
	json.NewDecoder(w.Body).Decode(&response)

	if response.FirstName != "John" {
		t.Errorf("Expected FirstName 'John', got '%s'", response.FirstName)
	}
	if response.LastName != "Doe" {
		t.Errorf("Expected LastName 'Doe', got '%s'", response.LastName)
	}
	if response.Position != "Manager" {
		t.Errorf("Expected Position 'Manager', got '%s'", response.Position)
	}
	if response.BaseSalary != baseSalary {
		t.Errorf("Expected BaseSalary %.2f, got %.2f", baseSalary, response.BaseSalary)
	}
}

// Test ListEmployees with authentication
func TestListEmployeesWithAuth(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	// Add a sample employee
	mockRepo.CreateEmployee(&domain.Employee{
		TenantID:     3,
		RestaurantID: 3,
		FirstName:    "Test",
		Email:        "test@example.com",
		IsActive:     true,
	})

	req := createRequestWithClaims("GET", "/api/v1/hr/employees", nil, claims)
	w := httptest.NewRecorder()

	handler.ListEmployees(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []domain.Employee
	json.NewDecoder(w.Body).Decode(&response)

	if len(response) != 1 {
		t.Errorf("Expected 1 employee, got %d", len(response))
	}
}

// Test ListEmployees without authentication
func TestListEmployeesNoAuth(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	req := createRequestWithClaims("GET", "/api/v1/hr/employees", nil, nil)
	w := httptest.NewRecorder()

	handler.ListEmployees(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test GetEmployee by ID
func TestGetEmployeeByID(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	// Create an employee
	id, _ := mockRepo.CreateEmployee(&domain.Employee{
		TenantID:     "3",
		RestaurantID: "3",
		FirstName:    "Test",
		Email:        "test@example.com",
		IsActive:     true,
	})

	req := createRequestWithClaims("GET", "/api/v1/hr/employees/"+string(rune(id)), nil, claims)
	req.SetPathValue("id", string(rune(id)))
	w := httptest.NewRecorder()

	handler.GetEmployee(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response domain.Employee
	json.NewDecoder(w.Body).Decode(&response)

	if response.FirstName != "Test" {
		t.Errorf("Expected FirstName 'Test', got '%s'", response.FirstName)
	}
}

// Test DeleteEmployee
func TestDeleteEmployee(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	// Create an employee
	id, _ := mockRepo.CreateEmployee(&domain.Employee{
		TenantID:     "3",
		RestaurantID: "3",
		FirstName:    "Test",
		Email:        "test@example.com",
		IsActive:     true,
	})

	req := createRequestWithClaims("DELETE", "/api/v1/hr/employees/"+string(rune(id)), nil, claims)
	req.SetPathValue("id", string(rune(id)))
	w := httptest.NewRecorder()

	handler.DeleteEmployee(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Verify employee is soft deleted
	emp, _ := mockRepo.GetEmployeeByID(3, 3, id)
	if emp != nil && emp.IsActive {
		t.Errorf("Expected employee to be soft deleted (IsActive=false)")
	}
}

// Test invalid ID format
func TestGetEmployeeInvalidID(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	req := createRequestWithClaims("GET", "/api/v1/hr/employees/invalid", nil, claims)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.GetEmployee(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Test employee not found
func TestGetEmployeeNotFound(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	req := createRequestWithClaims("GET", "/api/v1/hr/employees/999", nil, claims)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.GetEmployee(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

// Test multi-tenant isolation
func TestMultiTenantIsolation(t *testing.T) {
	mockRepo := NewMockEmployeeRepository()
	handler := NewEmployeeHandler(mockRepo)

	claims1 := &middleware.UserClaims{
		UserID:       "1",
		TenantID:     "3",
		RestaurantID: "3",
	}

	claims2 := &middleware.UserClaims{
		UserID:       "2",
		TenantID:     "4",
		RestaurantID: "4",
	}

	// Create employees for different tenants
	mockRepo.CreateEmployee(&domain.Employee{
		TenantID:     "3",
		RestaurantID: "3",
		FirstName:    "Employee1",
		Email:        "emp1@test.com",
		IsActive:     true,
	})

	mockRepo.CreateEmployee(&domain.Employee{
		TenantID:     "4",
		RestaurantID: "4",
		FirstName:    "Employee2",
		Email:        "emp2@test.com",
		IsActive:     true,
	})

	// User from Tenant 3 should only see their employee
	req := createRequestWithClaims("GET", "/api/v1/hr/employees", nil, claims1)
	w := httptest.NewRecorder()
	handler.ListEmployees(w, req)

	var response []domain.Employee
	json.NewDecoder(w.Body).Decode(&response)

	if len(response) != 1 {
		t.Errorf("Expected 1 employee for tenant 3, got %d", len(response))
	}

	if response[0].TenantID != "3" {
		t.Errorf("Expected TenantID '3', got '%s'", response[0].TenantID)
	}

	// User from Tenant 4 should only see their employee
	req = createRequestWithClaims("GET", "/api/v1/hr/employees", nil, claims2)
	w = httptest.NewRecorder()
	handler.ListEmployees(w, req)

	response = []domain.Employee{}
	json.NewDecoder(w.Body).Decode(&response)

	if len(response) != 1 {
		t.Errorf("Expected 1 employee for tenant 4, got %d", len(response))
	}

	if response[0].TenantID != "4" {
		t.Errorf("Expected TenantID '4', got '%s'", response[0].TenantID)
	}
}
