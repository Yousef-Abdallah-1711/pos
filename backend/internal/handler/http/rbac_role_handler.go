package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

// RBACRoleHandler manages RBAC role endpoints.
type RBACRoleHandler struct {
	roleRepo           *repository.RBACRoleRepository
	rolePermissionRepo *repository.RolePermissionRepository
}

// NewRBACRoleHandler creates a new RBAC role handler.
func NewRBACRoleHandler(roleRepo *repository.RBACRoleRepository, rolePermissionRepo *repository.RolePermissionRepository) *RBACRoleHandler {
	return &RBACRoleHandler{
		roleRepo:           roleRepo,
		rolePermissionRepo: rolePermissionRepo,
	}
}

type roleListResponse struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	RoleName         string `json:"role_name"`
	RoleCode         string `json:"role_code"`
	Description      string `json:"description"`
	IsSystem         bool   `json:"is_system"`
	PermissionsCount int    `json:"permissions_count"`
	CreatedAt        string `json:"created_at"`
}

type rolePermissionResponse struct {
	ModuleID        int64  `json:"module_id"`
	PermissionLevel string `json:"permission_level"`
	Level           string `json:"level"`
}

type roleDetailResponse struct {
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	RoleName    string                   `json:"role_name"`
	RoleCode    string                   `json:"role_code"`
	Description string                   `json:"description"`
	IsSystem    bool                     `json:"is_system"`
	CreatedAt   string                   `json:"created_at"`
	Permissions []rolePermissionResponse `json:"permissions"`
}

// ListRoles returns all active roles for the tenant/restaurant.
// GET /api/v1/rbac/roles
func (h *RBACRoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	roles, err := h.roleRepo.ListWithPermissionCounts(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]roleListResponse, 0, len(roles))
	for _, role := range roles {
		response = append(response, roleListResponse{
			ID:               role.ID,
			Name:             role.RoleName,
			RoleName:         role.RoleName,
			RoleCode:         role.RoleCode,
			Description:      role.Description,
			IsSystem:         role.IsSystemRole,
			PermissionsCount: role.PermissionsCount,
			CreatedAt:        role.CreatedAt.Format(timeFormat),
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": response})
}

// GetRole returns a role with its permissions.
// GET /api/v1/rbac/roles/{id}
func (h *RBACRoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := r.PathValue("id")
	roleID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleRepo.GetByID(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), roleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil || !role.IsActive {
		respondError(w, http.StatusNotFound, "Role not found")
		return
	}

	rolePermissions, err := h.rolePermissionRepo.GetByRoleAndTenant(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), role.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	permissions := make([]rolePermissionResponse, 0, len(rolePermissions))
	for _, perm := range rolePermissions {
		permissions = append(permissions, rolePermissionResponse{
			ModuleID:        perm.ModuleID,
			PermissionLevel: perm.PermissionLevel,
			Level:           perm.PermissionLevel,
		})
	}

	response := roleDetailResponse{
		ID:          role.ID,
		Name:        role.RoleName,
		RoleName:    role.RoleName,
		RoleCode:    role.RoleCode,
		Description: role.Description,
		IsSystem:    role.IsSystemRole,
		CreatedAt:   role.CreatedAt.Format(timeFormat),
		Permissions: permissions,
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": response})
}

// CreateRole creates a new role and assigns permissions.
// POST /api/v1/rbac/roles
func (h *RBACRoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	payload, err := decodeJSONMap(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	roleName := getStringField(payload, "role_name", "roleName", "name")
	if roleName == "" {
		respondError(w, http.StatusBadRequest, "Role name is required")
		return
	}
	description := getStringField(payload, "description")
	roleCode := getStringField(payload, "role_code", "roleCode")
	if roleCode == "" {
		roleCode = buildRoleCode(roleName)
	}

	permissions, err := parsePermissionInputs(payload["permissions"])
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	role := &repository.RBACRole{
		TenantID:     int64(claims.TenantID),
		RestaurantID: int64(claims.RestaurantID),
		RoleName:     roleName,
		Description:  description,
		RoleCode:     roleCode,
		IsSystemRole: false,
		IsActive:     true,
	}

	roleID, err := h.roleRepo.CreateWithPermissions(r.Context(), role, int64(claims.UserID), toRBACPermissionMutations(permissions))
	if err != nil {
		if strings.Contains(err.Error(), "unique_role_name_per_tenant") {
			respondError(w, http.StatusBadRequest, "Role name already exists")
			return
		}
		if strings.Contains(err.Error(), "unique_role_code_per_tenant") {
			respondError(w, http.StatusBadRequest, "Role code already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create role")
		return
	}

	role.ID = roleID
	response := roleDetailResponse{
		ID:          role.ID,
		Name:        role.RoleName,
		RoleName:    role.RoleName,
		RoleCode:    role.RoleCode,
		Description: role.Description,
		IsSystem:    role.IsSystemRole,
		CreatedAt:   timeNowFormatted(),
	}
	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": response})
}

// UpdateRole updates role details and permissions.
// PUT /api/v1/rbac/roles/{id}
func (h *RBACRoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	roleID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	existingRole, err := h.roleRepo.GetByID(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), roleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existingRole == nil || !existingRole.IsActive {
		respondError(w, http.StatusNotFound, "Role not found")
		return
	}

	payload, err := decodeJSONMap(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	roleName := getStringField(payload, "role_name", "roleName", "name")
	if roleName == "" {
		roleName = existingRole.RoleName
	}
	description := getStringField(payload, "description")
	if description == "" {
		description = existingRole.Description
	}

	permissionsRaw, replacePermissions := payload["permissions"]
	var permissions []permissionInput
	if replacePermissions {
		permissions, err = parsePermissionInputs(permissionsRaw)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	updateRole := &repository.RBACRole{
		ID:           existingRole.ID,
		TenantID:     existingRole.TenantID,
		RestaurantID: existingRole.RestaurantID,
		RoleName:     roleName,
		Description:  description,
	}

	if err := h.roleRepo.UpdateWithPermissions(r.Context(), updateRole, int64(claims.UserID), toRBACPermissionMutations(permissions), replacePermissions); err != nil {
		if errors.Is(err, repository.ErrRBACRoleNotFound) {
			if existingRole.IsSystemRole {
				respondError(w, http.StatusForbidden, "Cannot update system roles")
				return
			}
			respondError(w, http.StatusNotFound, "Role not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update role")
		return
	}

	response := roleDetailResponse{
		ID:          roleID,
		Name:        roleName,
		RoleName:    roleName,
		RoleCode:    existingRole.RoleCode,
		Description: description,
		IsSystem:    existingRole.IsSystemRole,
		CreatedAt:   existingRole.CreatedAt.Format(timeFormat),
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": response})
}

// DeleteRole deactivates a role.
// DELETE /api/v1/rbac/roles/{id}
func (h *RBACRoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	roleID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleRepo.GetByID(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), roleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil || !role.IsActive {
		respondError(w, http.StatusNotFound, "Role not found")
		return
	}
	if role.IsSystemRole {
		respondError(w, http.StatusForbidden, "Cannot delete system roles")
		return
	}

	if err := h.roleRepo.DeactivateWithPermissions(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), roleID, int64(claims.UserID)); err != nil {
		if errors.Is(err, repository.ErrRBACRoleNotFound) {
			respondError(w, http.StatusNotFound, "Role not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete role")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Role deleted successfully"})
}

func toRBACPermissionMutations(permissions []permissionInput) []repository.RBACPermissionMutation {
	mutations := make([]repository.RBACPermissionMutation, 0, len(permissions))
	for _, permission := range permissions {
		mutations = append(mutations, repository.RBACPermissionMutation{ModuleID: permission.ModuleID, Level: permission.Level})
	}
	return mutations
}

func buildRoleCode(name string) string {
	upper := strings.ToUpper(strings.TrimSpace(name))
	builder := strings.Builder{}
	for _, r := range upper {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			builder.WriteRune('_')
		}
	}
	code := strings.Trim(builder.String(), "_")
	if code == "" {
		return "ROLE"
	}
	return code
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

func timeNowFormatted() string {
	return time.Now().Format(timeFormat)
}
