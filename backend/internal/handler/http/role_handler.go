package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

// RoleHandler handles HTTP requests for roles
type RoleHandler struct {
	repo *repository.RoleRepository
}

// NewRoleHandler creates new role handler
func NewRoleHandler(repo *repository.RoleRepository) *RoleHandler {
	return &RoleHandler{repo: repo}
}

// ListRoles retrieves all roles for a restaurant
// GET /api/v1/hr/roles
func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		log.Println("[ListRoles] ERROR: Claims is nil - user not authenticated")
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	roles, err := h.repo.ListRoles(claims.TenantID, claims.RestaurantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, roles)
}

// GetRole retrieves a single role by ID
// GET /api/v1/hr/roles/{id}
func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		log.Println("[GetRole] ERROR: Claims is nil - user not authenticated")
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.repo.GetRoleByID(claims.TenantID, claims.RestaurantID, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if role == nil {
		respondError(w, http.StatusNotFound, "Role not found")
		return
	}

	respondJSON(w, http.StatusOK, role)
}

// CreateRole creates a new role
// POST /api/v1/hr/roles
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	log.Println("[CreateRole] Starting...")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		log.Println("[CreateRole] ERROR: Claims is nil - user not authenticated")
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	log.Printf("[CreateRole] Claims retrieved - TenantID: %d, RestaurantID: %d\n", claims.TenantID, claims.RestaurantID)

	var input struct {
		RoleName           string          `json:"role_name"`
		RoleNameAr         string          `json:"role_name_ar"`
		Description        string          `json:"description"`
		DescriptionAr      string          `json:"description_ar"`
		RoleCode           string          `json:"role_code"`
		Permissions        json.RawMessage `json:"permissions"`
		AccessLevel        string          `json:"access_level"`
		CanApproveLeaves   bool            `json:"can_approve_leaves"`
		CanApproveOvertime bool            `json:"can_approve_overtime"`
		CanManagePayroll   bool            `json:"can_manage_payroll"`
		CanViewReports     bool            `json:"can_view_reports"`
		MinSalary          *float64        `json:"min_salary"`
		MaxSalary          *float64        `json:"max_salary"`
		DisplayOrder       int             `json:"display_order"`
	}

	log.Println("[CreateRole] Decoding request body...")
	if err := decodeJSONStrict(r, &input); err != nil {
		log.Printf("[CreateRole] ERROR decoding body: %v\n", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	log.Printf("[CreateRole] Body decoded - RoleName: %s, RoleCode: %s\n", input.RoleName, input.RoleCode)

	// No required fields - backend accepts empty values
	// Set defaults for empty values
	if input.RoleName == "" {
		input.RoleName = "Unnamed Role"
		log.Println("[CreateRole] Set default RoleName: Unnamed Role")
	}
	if input.RoleCode == "" {
		input.RoleCode = fmt.Sprintf("ROLE-%d", time.Now().UnixNano()%10000)
		log.Printf("[CreateRole] Generated RoleCode: %s\n", input.RoleCode)
	}
	if input.AccessLevel == "" {
		input.AccessLevel = "basic"
		log.Println("[CreateRole] Set default AccessLevel: basic")
	}

	log.Println("[CreateRole] Building role object...")
	role := &domain.Role{
		TenantID:           claims.TenantID,
		RestaurantID:       claims.RestaurantID,
		RoleName:           input.RoleName,
		RoleNameAr:         input.RoleNameAr,
		Description:        input.Description,
		DescriptionAr:      input.DescriptionAr,
		RoleCode:           input.RoleCode,
		Permissions:        input.Permissions,
		AccessLevel:        input.AccessLevel,
		CanApproveLeaves:   input.CanApproveLeaves,
		CanApproveOvertime: input.CanApproveOvertime,
		CanManagePayroll:   input.CanManagePayroll,
		CanViewReports:     input.CanViewReports,
		MinSalary:          input.MinSalary,
		MaxSalary:          input.MaxSalary,
		DisplayOrder:       input.DisplayOrder,
		IsActive:           true,
		CreatedBy:          &claims.UserID,
	}
	log.Printf("[CreateRole] Role object built - Code: %s, Name: %s\n", role.RoleCode, role.RoleName)

	log.Println("[CreateRole] Calling repository CreateRole...")
	id, err := h.repo.CreateRole(role)
	if err != nil {
		log.Printf("[CreateRole] ERROR creating role: %v\n", err)
		// Check for specific database constraint violations
		errMsg := err.Error()
		if strings.Contains(errMsg, "unique_role_code_per_tenant") {
			respondError(w, http.StatusBadRequest, "A role with this code already exists in your organization")
			return
		}
		if strings.Contains(errMsg, "unique_role_name_per_tenant") {
			respondError(w, http.StatusBadRequest, "A role with this name already exists in your organization")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create role. Please try again.")
		return
	}
	log.Printf("[CreateRole] Role created successfully - ID: %d\n", id)

	role.ID = id
	respondJSON(w, http.StatusCreated, role)
	log.Println("[CreateRole] Response sent successfully")
}

// UpdateRole updates an existing role
// PUT /api/v1/hr/roles/{id}
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		log.Println("[UpdateRole] ERROR: Claims is nil - user not authenticated")
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var input struct {
		RoleName           string          `json:"role_name"`
		RoleNameAr         string          `json:"role_name_ar"`
		Description        string          `json:"description"`
		DescriptionAr      string          `json:"description_ar"`
		Permissions        json.RawMessage `json:"permissions"`
		AccessLevel        string          `json:"access_level"`
		CanApproveLeaves   bool            `json:"can_approve_leaves"`
		CanApproveOvertime bool            `json:"can_approve_overtime"`
		CanManagePayroll   bool            `json:"can_manage_payroll"`
		CanViewReports     bool            `json:"can_view_reports"`
		MinSalary          *float64        `json:"min_salary"`
		MaxSalary          *float64        `json:"max_salary"`
		DisplayOrder       int             `json:"display_order"`
		IsActive           bool            `json:"is_active"`
	}

	if err := decodeJSONStrict(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.RoleName == "" {
		respondError(w, http.StatusBadRequest, "Role name is required")
		return
	}

	role := &domain.Role{
		ID:                 id,
		TenantID:           claims.TenantID,
		RestaurantID:       claims.RestaurantID,
		RoleName:           input.RoleName,
		RoleNameAr:         input.RoleNameAr,
		Description:        input.Description,
		DescriptionAr:      input.DescriptionAr,
		Permissions:        input.Permissions,
		AccessLevel:        input.AccessLevel,
		CanApproveLeaves:   input.CanApproveLeaves,
		CanApproveOvertime: input.CanApproveOvertime,
		CanManagePayroll:   input.CanManagePayroll,
		CanViewReports:     input.CanViewReports,
		MinSalary:          input.MinSalary,
		MaxSalary:          input.MaxSalary,
		DisplayOrder:       input.DisplayOrder,
		IsActive:           input.IsActive,
		UpdatedBy:          &claims.UserID,
	}

	if err := h.repo.UpdateRole(role); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, role)
}

// DeleteRole soft deletes a role
// DELETE /api/v1/hr/roles/{id}
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		log.Println("[DeleteRole] ERROR: Claims is nil - user not authenticated")
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	if err := h.repo.DeleteRole(claims.TenantID, claims.RestaurantID, id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Role deleted successfully"})
}
