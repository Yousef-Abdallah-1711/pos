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
	"pos-saas/internal/usecase"
)

// TestRBACRoleHandler_CreateRole tests role creation via HTTP handler
func TestRBACRoleHandler_CreateRole(t *testing.T) {
	// Create mock repositories
	roleRepo := NewMockRoleRepository()
	auditRepo := NewMockAuditLogRepository()

	// Create usecase and handler
	roleUC := usecase.NewRoleUseCase(roleRepo, auditRepo)
	roleHandler := NewRBACRoleHandler(roleUC)

	// Create request
	reqBody := CreateRoleRequest{
		RoleName:    "Test Manager",
		Description: "Test manager role",
	}
	body, _ := json.Marshal(reqBody)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/api/v1/rbac/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Add user claims to context
	claims := &jwt.Claims{
		UserID: 1,
		Email:  "admin@example.com",
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	// Set tenant ID header
	req.Header.Set("X-Tenant-ID", "1")

	// Create response writer
	recorder := httptest.NewRecorder()

	// Call handler
	roleHandler.CreateRole(recorder, req)

	// Check response status
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", recorder.Code)
	}

	// Parse response
	var response SuccessResponse
	json.NewDecoder(recorder.Body).Decode(&response)

	if !response.Success {
		t.Error("Expected success response")
	}
}

// TestRBACRoleHandler_DeleteSystemRole tests that system roles cannot be deleted
func TestRBACRoleHandler_DeleteSystemRole(t *testing.T) {
	// Create mock repositories
	roleRepo := NewMockRoleRepository()
	auditRepo := NewMockAuditLogRepository()

	// Create a system role
	systemRole := &domain.Role{
		ID:       1,
		TenantID: 1,
		RoleName: "Admin",
		IsSystem: true,
	}
	roleRepo.roles[1] = systemRole

	// Create usecase and handler
	roleUC := usecase.NewRoleUseCase(roleRepo, auditRepo)
	roleHandler := NewRBACRoleHandler(roleUC)

	// Create HTTP request
	req := httptest.NewRequest("DELETE", "/api/v1/rbac/roles/1", nil)
	req.RequestURI = "/api/v1/rbac/roles/1"

	// Add user claims to context
	claims := &jwt.Claims{
		UserID: 1,
		Email:  "admin@example.com",
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	// Set tenant ID header
	req.Header.Set("X-Tenant-ID", "1")

	// Create response writer
	recorder := httptest.NewRecorder()

	// Call handler
	roleHandler.DeleteRole(recorder, req)

	// Check response status (should be forbidden)
	if recorder.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", recorder.Code)
	}
}

// TestRBACPermissionHandler_AssignPermission tests permission assignment
func TestRBACPermissionHandler_AssignPermission(t *testing.T) {
	// Create mock repositories
	rolePermRepo := NewMockRolePermissionRepository()
	userRoleRepo := NewMockUserRoleRepository()
	auditRepo := NewMockAuditLogRepository()

	// Create usecase and handler
	permUC := usecase.NewPermissionUseCase(rolePermRepo, userRoleRepo, nil, auditRepo)
	permHandler := NewRBACPermissionHandler(permUC)

	// Create request
	reqBody := AssignPermissionRequest{
		RoleID:          1,
		ModuleID:        1,
		PermissionLevel: domain.PermissionRead,
	}
	body, _ := json.Marshal(reqBody)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/api/v1/rbac/permissions/assign", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Add user claims to context
	claims := &jwt.Claims{
		UserID: 1,
		Email:  "admin@example.com",
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	// Set tenant ID header
	req.Header.Set("X-Tenant-ID", "1")

	// Create response writer
	recorder := httptest.NewRecorder()

	// Call handler
	permHandler.AssignPermission(recorder, req)

	// Check response status
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
}

// TestRBACUserRoleHandler_AssignRoleToUser tests role assignment to user
func TestRBACUserRoleHandler_AssignRoleToUser(t *testing.T) {
	// Create mock repositories
	userRoleRepo := NewMockUserRoleRepository()
	auditRepo := NewMockAuditLogRepository()

	// Create usecase and handler
	userRoleUC := usecase.NewUserRoleUseCase(userRoleRepo, auditRepo)
	userRoleHandler := NewRBACUserRoleHandler(userRoleUC)

	// Create request
	reqBody := AssignUserRoleRequest{
		RoleID: 1,
	}
	body, _ := json.Marshal(reqBody)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/api/v1/rbac/users/1/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = "/api/v1/rbac/users/1/roles"

	// Add user claims to context
	claims := &jwt.Claims{
		UserID: 1,
		Email:  "admin@example.com",
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	// Set tenant ID header
	req.Header.Set("X-Tenant-ID", "1")

	// Create response writer
	recorder := httptest.NewRecorder()

	// Call handler
	userRoleHandler.AssignRoleToUser(recorder, req)

	// Check response status
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
}

// Mock repositories for handler testing
type MockRoleRepository struct {
	roles  map[int64]*domain.Role
	nextID int64
}

func NewMockRoleRepository() *MockRoleRepository {
	return &MockRoleRepository{
		roles:  make(map[int64]*domain.Role),
		nextID: 1,
	}
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	return m.roles[id], nil
}

func (m *MockRoleRepository) GetByName(ctx context.Context, tenantID int64, name string) (*domain.Role, error) {
	for _, role := range m.roles {
		if role.TenantID == tenantID && role.RoleName == name {
			return role, nil
		}
	}
	return nil, nil
}

func (m *MockRoleRepository) GetByTenant(ctx context.Context, tenantID int64) ([]domain.Role, error) {
	var roles []domain.Role
	for _, role := range m.roles {
		if role.TenantID == tenantID {
			roles = append(roles, *role)
		}
	}
	return roles, nil
}

func (m *MockRoleRepository) Create(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	role.ID = m.nextID
	m.nextID++
	m.roles[role.ID] = role
	return role, nil
}

func (m *MockRoleRepository) Update(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	m.roles[role.ID] = role
	return role, nil
}

func (m *MockRoleRepository) Delete(ctx context.Context, id int64) error {
	delete(m.roles, id)
	return nil
}

type MockRolePermissionRepository struct {
	permissions map[string]*domain.RolePermission
}

func NewMockRolePermissionRepository() *MockRolePermissionRepository {
	return &MockRolePermissionRepository{
		permissions: make(map[string]*domain.RolePermission),
	}
}

func (m *MockRolePermissionRepository) GetByRoleAndModule(ctx context.Context, tenantID, restaurantID, roleID, moduleID int64) (string, error) {
	key := generateKey(roleID, moduleID)
	if perm, exists := m.permissions[key]; exists {
		return perm.PermissionLevel, nil
	}
	return domain.PermissionNone, nil
}

func (m *MockRolePermissionRepository) GetByRole(ctx context.Context, tenantID, restaurantID, roleID int64) ([]domain.RolePermission, error) {
	var perms []domain.RolePermission
	for _, perm := range m.permissions {
		if perm.RoleID == roleID {
			perms = append(perms, *perm)
		}
	}
	return perms, nil
}

func (m *MockRolePermissionRepository) GetByModule(ctx context.Context, tenantID, restaurantID, moduleID int64) ([]domain.RolePermission, error) {
	var perms []domain.RolePermission
	for _, perm := range m.permissions {
		if perm.ModuleID == moduleID {
			perms = append(perms, *perm)
		}
	}
	return perms, nil
}

func (m *MockRolePermissionRepository) Assign(ctx context.Context, tenantID, restaurantID, roleID, moduleID int64, level string) error {
	key := generateKey(roleID, moduleID)
	m.permissions[key] = &domain.RolePermission{
		TenantID:        tenantID,
		RoleID:          roleID,
		ModuleID:        moduleID,
		PermissionLevel: level,
	}
	return nil
}

func (m *MockRolePermissionRepository) Revoke(ctx context.Context, tenantID, restaurantID, roleID, moduleID int64) error {
	key := generateKey(roleID, moduleID)
	delete(m.permissions, key)
	return nil
}

func (m *MockRolePermissionRepository) RevokeAll(ctx context.Context, tenantID, restaurantID, roleID int64) error {
	for key := range m.permissions {
		if m.permissions[key].RoleID == roleID {
			delete(m.permissions, key)
		}
	}
	return nil
}

func (m *MockRolePermissionRepository) GetByRoleAndTenant(ctx context.Context, tenantID, restaurantID, roleID int64) ([]domain.RolePermission, error) {
	var perms []domain.RolePermission
	for _, perm := range m.permissions {
		if perm.TenantID == tenantID && perm.RoleID == roleID {
			perms = append(perms, *perm)
		}
	}
	return perms, nil
}

type MockUserRoleRepository struct {
	userRoles map[string]*domain.UserRole
}

func NewMockUserRoleRepository() *MockUserRoleRepository {
	return &MockUserRoleRepository{
		userRoles: make(map[string]*domain.UserRole),
	}
}

func (m *MockUserRoleRepository) GetUserRoles(ctx context.Context, userID, tenantID, restaurantID int64) ([]domain.Role, error) {
	var roles []domain.Role
	for _, ur := range m.userRoles {
		if ur.UserID == userID && ur.TenantID == tenantID {
			// Would normally fetch from role repository
		}
	}
	return roles, nil
}

func (m *MockUserRoleRepository) GetRoleUsers(ctx context.Context, tenantID, restaurantID, roleID int64) ([]domain.User, error) {
	var users []domain.User
	for _, ur := range m.userRoles {
		if ur.RoleID == roleID {
			// Would normally fetch from user repository
		}
	}
	return users, nil
}

func (m *MockUserRoleRepository) Assign(ctx context.Context, tenantID, restaurantID, userID, roleID, assignedBy int64) error {
	key := generateUserRoleKey(tenantID, userID, roleID)
	m.userRoles[key] = &domain.UserRole{
		TenantID: tenantID,
		UserID:   userID,
		RoleID:   roleID,
	}
	return nil
}

func (m *MockUserRoleRepository) Remove(ctx context.Context, tenantID, restaurantID, userID, roleID int64) error {
	key := generateUserRoleKey(tenantID, userID, roleID)
	delete(m.userRoles, key)
	return nil
}

func (m *MockUserRoleRepository) RemoveAll(ctx context.Context, tenantID, restaurantID, userID int64) error {
	for key := range m.userRoles {
		if m.userRoles[key].TenantID == tenantID && m.userRoles[key].UserID == userID {
			delete(m.userRoles, key)
		}
	}
	return nil
}

func (m *MockUserRoleRepository) HasRole(ctx context.Context, tenantID int64, userID int64, roleID int64) (bool, error) {
	key := generateUserRoleKey(tenantID, userID, roleID)
	_, exists := m.userRoles[key]
	return exists, nil
}

func (m *MockUserRoleRepository) CountUserRoles(ctx context.Context, tenantID int64, userID int64) (int, error) {
	count := 0
	for _, ur := range m.userRoles {
		if ur.TenantID == tenantID && ur.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *MockUserRoleRepository) GetUsersWithRole(ctx context.Context, tenantID, restaurantID, roleID int64) ([]domain.User, error) {
	var users []domain.User
	for _, ur := range m.userRoles {
		if ur.TenantID == tenantID && ur.RoleID == roleID {
			// Would normally fetch from user repository
		}
	}
	return users, nil
}

type MockAuditLogRepository struct {
	logs []*domain.PermissionAuditLog
}

func NewMockAuditLogRepository() *MockAuditLogRepository {
	return &MockAuditLogRepository{
		logs: make([]*domain.PermissionAuditLog, 0),
	}
}

func (m *MockAuditLogRepository) Log(ctx context.Context, log *domain.PermissionAuditLog) error {
	m.logs = append(m.logs, log)
	return nil
}

func (m *MockAuditLogRepository) GetByTenant(ctx context.Context, tenantID int64) ([]domain.PermissionAuditLog, error) {
	var logs []domain.PermissionAuditLog
	for _, log := range m.logs {
		if log.TenantID == tenantID {
			logs = append(logs, *log)
		}
	}
	return logs, nil
}

func (m *MockAuditLogRepository) GetByUser(ctx context.Context, userID int64) ([]domain.PermissionAuditLog, error) {
	var logs []domain.PermissionAuditLog
	for _, log := range m.logs {
		if log.UserID == userID {
			logs = append(logs, *log)
		}
	}
	return logs, nil
}

func (m *MockAuditLogRepository) GetByAction(ctx context.Context, action string) ([]domain.PermissionAuditLog, error) {
	var logs []domain.PermissionAuditLog
	for _, log := range m.logs {
		if log.Action == action {
			logs = append(logs, *log)
		}
	}
	return logs, nil
}

func (m *MockAuditLogRepository) GetRecent(ctx context.Context, limit int) ([]domain.PermissionAuditLog, error) {
	if limit > len(m.logs) {
		limit = len(m.logs)
	}
	var logs []domain.PermissionAuditLog
	for i := 0; i < limit; i++ {
		logs = append(logs, *m.logs[len(m.logs)-i-1])
	}
	return logs, nil
}

func (m *MockAuditLogRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.logs)), nil
}

// Helper functions
func generateKey(roleID, moduleID int64) string {
	return string(rune(roleID)) + "_" + string(rune(moduleID))
}

func generateUserRoleKey(tenantID, userID, roleID int64) string {
	return string(rune(tenantID)) + "_" + string(rune(userID)) + "_" + string(rune(roleID))
}
