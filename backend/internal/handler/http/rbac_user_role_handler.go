package http

import (
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

// RBACUserRoleHandler manages user-role assignment endpoints.
type RBACUserRoleHandler struct {
	userRoleRepo *repository.UserRoleRepository
	roleRepo     *repository.RBACRoleRepository
	userRepo     *repository.UserRepository
}

// NewRBACUserRoleHandler creates a new RBAC user-role handler.
func NewRBACUserRoleHandler(
	userRoleRepo *repository.UserRoleRepository,
	roleRepo *repository.RBACRoleRepository,
	userRepo *repository.UserRepository,
) *RBACUserRoleHandler {
	return &RBACUserRoleHandler{
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		userRepo:     userRepo,
	}
}

// AssignRoleToUser assigns a role to a user.
// POST /api/v1/rbac/users/{userId}/roles
func (h *RBACUserRoleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("userId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
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

	user, err := h.userRepo.GetByID(r.Context(), int(userID), claims.TenantID, claims.RestaurantID)
	if err != nil || user == nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}
	if int64(user.TenantID) != int64(claims.TenantID) {
		respondError(w, http.StatusForbidden, "User not in tenant")
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

	if err := h.userRoleRepo.Assign(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), userID, roleID, int64(claims.UserID)); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to assign role")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Role assigned"})
}

// RemoveRoleFromUser removes a role from a user.
// DELETE /api/v1/rbac/users/{userId}/roles/{roleId}
func (h *RBACUserRoleHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("userId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roleID, err := strconv.ParseInt(r.PathValue("roleId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	if err := h.userRoleRepo.Remove(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), userID, roleID); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Role removed"})
}

// GetUserRoles returns roles assigned to a user.
// GET /api/v1/rbac/users/{userId}/roles
func (h *RBACUserRoleHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("userId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roles, err := h.userRoleRepo.GetUserRoles(r.Context(), userID, int64(claims.TenantID), int64(claims.RestaurantID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": roles})
}

// GetMyRoles returns roles assigned to the current user.
// GET /api/v1/rbac/me/roles
func (h *RBACUserRoleHandler) GetMyRoles(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	roles, err := h.userRoleRepo.GetUserRoles(r.Context(), int64(claims.UserID), int64(claims.TenantID), int64(claims.RestaurantID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": roles})
}

// AssignRoleToUserByEmail assigns a role to a user by email.
// POST /api/v1/rbac/users-by-email/{email}/roles
func (h *RBACUserRoleHandler) AssignRoleToUserByEmail(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	email := strings.TrimSpace(r.PathValue("email"))
	if email == "" {
		respondError(w, http.StatusBadRequest, "Email is required")
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

	user, err := h.userRepo.GetByEmail(r.Context(), email)
	if err != nil || user == nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}
	if int64(user.TenantID) != int64(claims.TenantID) {
		respondError(w, http.StatusForbidden, "User not in tenant")
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

	if err := h.userRoleRepo.Assign(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), int64(user.ID), roleID, int64(claims.UserID)); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to assign role")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Role assigned"})
}

// RemoveRoleFromUserByEmail removes a role assignment by email.
// DELETE /api/v1/rbac/users-by-email/{email}/roles/{roleId}
func (h *RBACUserRoleHandler) RemoveRoleFromUserByEmail(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	email := strings.TrimSpace(r.PathValue("email"))
	if email == "" {
		respondError(w, http.StatusBadRequest, "Email is required")
		return
	}

	roleID, err := strconv.ParseInt(r.PathValue("roleId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	user, err := h.userRepo.GetByEmail(r.Context(), email)
	if err != nil || user == nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}
	if int64(user.TenantID) != int64(claims.TenantID) {
		respondError(w, http.StatusForbidden, "User not in tenant")
		return
	}

	if err := h.userRoleRepo.Remove(r.Context(), int64(claims.TenantID), int64(claims.RestaurantID), int64(user.ID), roleID); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Role removed"})
}

// GetUserRolesByEmail returns roles for a user by email.
// GET /api/v1/rbac/users-by-email/{email}/roles
func (h *RBACUserRoleHandler) GetUserRolesByEmail(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	email := strings.TrimSpace(r.PathValue("email"))
	if email == "" {
		respondError(w, http.StatusBadRequest, "Email is required")
		return
	}

	user, err := h.userRepo.GetByEmail(r.Context(), email)
	if err != nil || user == nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}
	if int64(user.TenantID) != int64(claims.TenantID) {
		respondError(w, http.StatusForbidden, "User not in tenant")
		return
	}

	roles, err := h.userRoleRepo.GetUserRoles(r.Context(), int64(user.ID), int64(claims.TenantID), int64(claims.RestaurantID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": roles})
}
