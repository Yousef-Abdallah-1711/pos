package http

import (
	"net/http"
	"strconv"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

// RBACPermissionHandler manages RBAC permission endpoints.
type RBACPermissionHandler struct {
	rolePermissionRepo *repository.RolePermissionRepository
	userRoleRepo       *repository.UserRoleRepository
	roleRepo           *repository.RBACRoleRepository
}

// NewRBACPermissionHandler creates a new RBAC permission handler.
func NewRBACPermissionHandler(
	rolePermissionRepo *repository.RolePermissionRepository,
	userRoleRepo *repository.UserRoleRepository,
	roleRepo *repository.RBACRoleRepository,
) *RBACPermissionHandler {
	return &RBACPermissionHandler{
		rolePermissionRepo: rolePermissionRepo,
		userRoleRepo:       userRoleRepo,
		roleRepo:           roleRepo,
	}
}

// AssignPermission assigns or updates a role's permission level for a module.
// POST /api/v1/rbac/permissions/assign
func (h *RBACPermissionHandler) AssignPermission(w http.ResponseWriter, r *http.Request) {
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

	roleID, ok, err := getInt64Field(payload, "role_id", "roleId")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !ok {
		respondError(w, http.StatusBadRequest, "role_id is required")
		return
	}

	moduleID, ok, err := getInt64Field(payload, "module_id", "moduleId")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !ok {
		respondError(w, http.StatusBadRequest, "module_id is required")
		return
	}

	levelRaw := getStringField(payload, "permission_level", "level")
	level, err := normalizePermissionLevel(levelRaw)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	role, err := h.roleRepo.GetByID(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), roleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil {
		respondError(w, http.StatusNotFound, "Role not found")
		return
	}

	if err := h.rolePermissionRepo.Assign(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), roleID, moduleID, level); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to assign permission")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Permission assigned"})
}

// RevokePermission removes a role's permission for a module.
// DELETE /api/v1/rbac/permissions/{roleId}/{moduleId}
func (h *RBACPermissionHandler) RevokePermission(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	roleID, err := strconv.ParseInt(r.PathValue("roleId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	moduleID, err := strconv.ParseInt(r.PathValue("moduleId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid module ID")
		return
	}

	if err := h.rolePermissionRepo.Revoke(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), roleID, moduleID); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Permission revoked"})
}

// GetRolePermissions lists permissions for a role.
// GET /api/v1/rbac/roles/{roleId}/permissions
func (h *RBACPermissionHandler) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	roleID, err := strconv.ParseInt(r.PathValue("roleId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	perms, err := h.rolePermissionRepo.GetByRoleAndTenant(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), roleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": perms})
}

// CheckPermission checks if the current user has a permission level for a module.
// GET /api/v1/rbac/check-permission?moduleId=1&permissionLevel=READ
func (h *RBACPermissionHandler) CheckPermission(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	query := r.URL.Query()
	moduleStr := query.Get("moduleId")
	if moduleStr == "" {
		moduleStr = query.Get("module_id")
	}
	if moduleStr == "" {
		respondError(w, http.StatusBadRequest, "moduleId is required")
		return
	}
	moduleID, err := strconv.ParseInt(moduleStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid moduleId")
		return
	}

	level := query.Get("permissionLevel")
	if level == "" {
		level = query.Get("permission_level")
	}
	if level == "" {
		level = domain.PermissionRead
	}
	requiredLevel, err := normalizePermissionLevel(level)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	userRoles, err := h.userRoleRepo.GetUserRoles(r.Context(), int64(claims.UserID), int64(claims.TenantID), int64(claims.RestaurantID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to check permissions")
		return
	}

	highest := domain.PermissionNone
	for _, role := range userRoles {
		roleLevel, err := h.rolePermissionRepo.GetByRoleAndModule(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), int64(role.ID), moduleID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to check permissions")
			return
		}

		fallbackLevel := fallbackPermissionLevelForRole(role, moduleID)
		highest = maxPermissionLevel(highest, roleLevel)
		highest = maxPermissionLevel(highest, fallbackLevel)
	}

	hasPermission := domain.PermissionValues[highest] >= domain.PermissionValues[requiredLevel]
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"has_permission":   hasPermission,
		"module_id":        moduleID,
		"permission_level": highest,
	})
}

// GetUserPermissions returns a map of module_id -> permission_level for the current user.
// GET /api/v1/rbac/me/permissions
func (h *RBACPermissionHandler) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userRoles, err := h.userRoleRepo.GetUserRoles(r.Context(), int64(claims.UserID), int64(claims.TenantID), int64(claims.RestaurantID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load permissions")
		return
	}

	permissions := map[int64]string{}
	for _, role := range userRoles {
		rolePerms, err := h.rolePermissionRepo.GetByRoleAndTenant(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), int64(role.ID))
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to load permissions")
			return
		}
		for _, perm := range rolePerms {
			current := permissions[perm.ModuleID]
			if current == "" {
				current = domain.PermissionNone
			}
			permissions[perm.ModuleID] = maxPermissionLevel(current, perm.PermissionLevel)
		}

		// Compatibility fallback:
		// older tenants may have module 13 but missing explicit role_permissions rows.
		kdsFallback := fallbackPermissionLevelForRole(role, 13)
		if kdsFallback != domain.PermissionNone {
			current := permissions[13]
			if current == "" {
				current = domain.PermissionNone
			}
			permissions[13] = maxPermissionLevel(current, kdsFallback)
		}

		// Compatibility fallback:
		// older tenants may have module 14 (TABLES) without explicit role_permissions rows.
		tablesFallback := fallbackPermissionLevelForRole(role, 14)
		if tablesFallback != domain.PermissionNone {
			current := permissions[14]
			if current == "" {
				current = domain.PermissionNone
			}
			permissions[14] = maxPermissionLevel(current, tablesFallback)
		}
	}

	respondJSON(w, http.StatusOK, permissions)
}

func maxPermissionLevel(current string, candidate string) string {
	return domain.MaxPermissionLevel(current, candidate)
}

func fallbackPermissionLevelForRole(role domain.Role, moduleID int64) string {
	return domain.FallbackPermissionLevel(role, moduleID)
}
