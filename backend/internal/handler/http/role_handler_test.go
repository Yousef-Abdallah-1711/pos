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

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
)

// MockRoleRepository is a mock implementation for testing
type MockRoleRepository struct {
	roles  map[int]*domain.Role
	nextID int
}

// NewMockRoleRepository creates a new mock repository
func NewMockRoleRepository() *MockRoleRepository {
	return &MockRoleRepository{
		roles:  make(map[int]*domain.Role),
		nextID: 1,
	}
}

// ListRoles returns all roles
func (m *MockRoleRepository) ListRoles(tenantID, restaurantID int) ([]domain.Role, error) {
	var result []domain.Role
	for _, role := range m.roles {
		if role.TenantID == tenantID && role.RestaurantID == restaurantID {
			result = append(result, *role)
		}
	}
	return result, nil
}

// GetRoleByID returns a single role
func (m *MockRoleRepository) GetRoleByID(tenantID, restaurantID int, id int) (*domain.Role, error) {
	if role, exists := m.roles[id]; exists {
		if role.TenantID == tenantID && role.RestaurantID == restaurantID {
			return role, nil
		}
	}
	return nil, nil
}

// CreateRole creates a new role
func (m *MockRoleRepository) CreateRole(role *domain.Role) (int, error) {
	id := m.nextID
	role.ID = id
	m.roles[id] = role
	m.nextID++
	return id, nil
}

// UpdateRole updates an existing role
func (m *MockRoleRepository) UpdateRole(role *domain.Role) error {
	if _, exists := m.roles[role.ID]; exists {
		m.roles[role.ID] = role
		return nil
	}
	return nil
}

// DeleteRole soft deletes a role
func (m *MockRoleRepository) DeleteRole(tenantID, restaurantID int, id int) error {
	if role, exists := m.roles[id]; exists {
		if role.TenantID == tenantID && role.RestaurantID == restaurantID {
			role.IsActive = false
		}
	}
	return nil
}

// Helper function to create a request with claims in context
func createRoleRequestWithClaims(method, url string, body interface{}, claims *jwt.Claims) *http.Request {
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

// Test CreateRole with minimal fields
func TestCreateRoleMinimalFields(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	inputBody := map[string]interface{}{
		"role_name": "Manager",
	}

	req := createRoleRequestWithClaims("POST", "/api/v1/hr/roles", inputBody, claims)
	w := httptest.NewRecorder()

	handler.CreateRole(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response domain.Role
	json.NewDecoder(w.Body).Decode(&response)

	if response.RoleName != "Manager" {
		t.Errorf("Expected RoleName 'Manager', got '%s'", response.RoleName)
	}
	if response.RoleCode == "" {
		t.Errorf("Expected auto-generated RoleCode, got empty string")
	}
	if response.AccessLevel != "1" {
		t.Errorf("Expected AccessLevel '1', got '%s'", response.AccessLevel)
	}
	if response.ID == 0 {
		t.Errorf("Expected valid ID, got 0")
	}
}

// Test CreateRole with no fields (defaults)
func TestCreateRoleNoFields(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	inputBody := map[string]interface{}{}

	req := createRoleRequestWithClaims("POST", "/api/v1/hr/roles", inputBody, claims)
	w := httptest.NewRecorder()

	handler.CreateRole(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response domain.Role
	json.NewDecoder(w.Body).Decode(&response)

	if response.RoleName != "Unnamed Role" {
		t.Errorf("Expected default RoleName 'Unnamed Role', got '%s'", response.RoleName)
	}
	if response.RoleCode == "" {
		t.Errorf("Expected auto-generated RoleCode, got empty string")
	}
}

// Test CreateRole without authentication
func TestCreateRoleNoAuth(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	inputBody := map[string]interface{}{
		"role_name": "Manager",
	}

	req := createRoleRequestWithClaims("POST", "/api/v1/hr/roles", inputBody, nil)
	w := httptest.NewRecorder()

	handler.CreateRole(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test CreateRole with full data
func TestCreateRoleFullData(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	minSalary := 30000.0
	maxSalary := 50000.0

	inputBody := map[string]interface{}{
		"role_name":            "Senior Manager",
		"role_code":            "ROLE-001",
		"description":          "Manages team operations",
		"access_level":         "5",
		"can_approve_leaves":   true,
		"can_approve_overtime": true,
		"can_manage_payroll":   true,
		"can_view_reports":     true,
		"min_salary":           minSalary,
		"max_salary":           maxSalary,
		"display_order":        1,
	}

	req := createRoleRequestWithClaims("POST", "/api/v1/hr/roles", inputBody, claims)
	w := httptest.NewRecorder()

	handler.CreateRole(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response domain.Role
	json.NewDecoder(w.Body).Decode(&response)

	if response.RoleName != "Senior Manager" {
		t.Errorf("Expected RoleName 'Senior Manager', got '%s'", response.RoleName)
	}
	if response.RoleCode != "ROLE-001" {
		t.Errorf("Expected RoleCode 'ROLE-001', got '%s'", response.RoleCode)
	}
	if !response.CanApproveLeaves {
		t.Errorf("Expected CanApproveLeaves true")
	}
	if !response.CanManagePayroll {
		t.Errorf("Expected CanManagePayroll true")
	}
}

// Test ListRoles with authentication
func TestListRolesWithAuth(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	// Add sample roles
	mockRepo.CreateRole(&domain.Role{
		TenantID:     "3",
		RestaurantID: "3",
		RoleName:     "Manager",
		RoleCode:     "ROLE-001",
		IsActive:     true,
	})

	mockRepo.CreateRole(&domain.Role{
		TenantID:     "3",
		RestaurantID: "3",
		RoleName:     "Staff",
		RoleCode:     "ROLE-002",
		IsActive:     true,
	})

	req := createRoleRequestWithClaims("GET", "/api/v1/hr/roles", nil, claims)
	w := httptest.NewRecorder()

	handler.ListRoles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []domain.Role
	json.NewDecoder(w.Body).Decode(&response)

	if len(response) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(response))
	}
}

// Test ListRoles without authentication
func TestListRolesNoAuth(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	req := createRoleRequestWithClaims("GET", "/api/v1/hr/roles", nil, nil)
	w := httptest.NewRecorder()

	handler.ListRoles(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test GetRole by ID
func TestGetRoleByID(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	// Create a role
	id, _ := mockRepo.CreateRole(&domain.Role{
		TenantID:     "3",
		RestaurantID: "3",
		RoleName:     "Manager",
		RoleCode:     "ROLE-001",
		IsActive:     true,
	})

	req := createRoleRequestWithClaims("GET", "/api/v1/hr/roles/"+string(rune(id)), nil, claims)
	req.SetPathValue("id", string(rune(id)))
	w := httptest.NewRecorder()

	handler.GetRole(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response domain.Role
	json.NewDecoder(w.Body).Decode(&response)

	if response.RoleName != "Manager" {
		t.Errorf("Expected RoleName 'Manager', got '%s'", response.RoleName)
	}
}

// Test DeleteRole
func TestDeleteRole(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	// Create a role
	id, _ := mockRepo.CreateRole(&domain.Role{
		TenantID:     "3",
		RestaurantID: "3",
		RoleName:     "Manager",
		RoleCode:     "ROLE-001",
		IsActive:     true,
	})

	req := createRoleRequestWithClaims("DELETE", "/api/v1/hr/roles/"+string(rune(id)), nil, claims)
	req.SetPathValue("id", string(rune(id)))
	w := httptest.NewRecorder()

	handler.DeleteRole(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Verify role is soft deleted
	role, _ := mockRepo.GetRoleByID(3, 3, id)
	if role != nil && role.IsActive {
		t.Errorf("Expected role to be soft deleted (IsActive=false)")
	}
}

// Test UpdateRole
func TestUpdateRole(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	// Create a role
	id, _ := mockRepo.CreateRole(&domain.Role{
		TenantID:     "3",
		RestaurantID: "3",
		RoleName:     "Manager",
		RoleCode:     "ROLE-001",
		IsActive:     true,
	})

	updateBody := map[string]interface{}{
		"role_name":   "Senior Manager",
		"description": "Manages operations",
		"is_active":   true,
	}

	req := createRoleRequestWithClaims("PUT", "/api/v1/hr/roles/"+string(rune(id)), updateBody, claims)
	req.SetPathValue("id", string(rune(id)))
	w := httptest.NewRecorder()

	handler.UpdateRole(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// Test invalid role ID
func TestGetRoleInvalidID(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	req := createRoleRequestWithClaims("GET", "/api/v1/hr/roles/invalid", nil, claims)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.GetRole(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Test role not found
func TestGetRoleNotFound(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	claims := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	req := createRoleRequestWithClaims("GET", "/api/v1/hr/roles/999", nil, claims)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.GetRole(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

// Test role multi-tenant isolation
func TestRoleMultiTenantIsolation(t *testing.T) {
	mockRepo := NewMockRoleRepository()
	handler := NewRoleHandler(mockRepo)

	claims1 := &jwt.Claims{
		UserID:       1,
		TenantID:     3,
		RestaurantID: 3,
	}

	claims2 := &jwt.Claims{
		UserID:       2,
		TenantID:     4,
		RestaurantID: 4,
	}

	// Create roles for different tenants
	mockRepo.CreateRole(&domain.Role{
		TenantID:     "3",
		RestaurantID: "3",
		RoleName:     "Manager",
		IsActive:     true,
	})

	mockRepo.CreateRole(&domain.Role{
		TenantID:     "4",
		RestaurantID: "4",
		RoleName:     "Staff",
		IsActive:     true,
	})

	// User from Tenant 3 should only see their role
	req := createRoleRequestWithClaims("GET", "/api/v1/hr/roles", nil, claims1)
	w := httptest.NewRecorder()
	handler.ListRoles(w, req)

	var response []domain.Role
	json.NewDecoder(w.Body).Decode(&response)

	if len(response) != 1 {
		t.Errorf("Expected 1 role for tenant 3, got %d", len(response))
	}

	if response[0].TenantID != "3" {
		t.Errorf("Expected TenantID '3', got '%s'", response[0].TenantID)
	}

	// User from Tenant 4 should only see their role
	req = createRoleRequestWithClaims("GET", "/api/v1/hr/roles", nil, claims2)
	w = httptest.NewRecorder()
	handler.ListRoles(w, req)

	response = []domain.Role{}
	json.NewDecoder(w.Body).Decode(&response)

	if len(response) != 1 {
		t.Errorf("Expected 1 role for tenant 4, got %d", len(response))
	}

	if response[0].TenantID != "4" {
		t.Errorf("Expected TenantID '4', got '%s'", response[0].TenantID)
	}
}
