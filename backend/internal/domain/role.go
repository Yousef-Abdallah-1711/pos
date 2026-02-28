package domain

import (
	"errors"
	"time"
)

// SimpleRole represents a basic role that groups permissions together (DEPRECATED - use HR Role instead)
// Roles can be system-defined (immutable) or custom (tenant-specific)
type SimpleRole struct {
	ID          int64
	TenantID    int64
	Name        string // ADMIN, MANAGER, STAFF, VIEWER, or custom name
	Description string
	IsSystem    bool // System roles cannot be deleted
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RoleConstants defines standard role names
const (
	RoleAdmin   = "ADMIN"
	RoleManager = "MANAGER"
	RoleStaff   = "STAFF"
	RoleViewer  = "VIEWER"
)

// Validate checks if role has required fields
func (r *SimpleRole) Validate() error {
	if r.Name == "" {
		return errors.New("role name is required")
	}

	if r.TenantID == 0 {
		return errors.New("tenant ID is required")
	}

	if len(r.Name) > 100 {
		return errors.New("role name must be less than 100 characters")
	}

	return nil
}

// IsSystemRole checks if this is a built-in system role
func (r *SimpleRole) IsSystemRole() bool {
	return r.IsSystem
}

// CanDelete checks if role can be deleted
func (r *SimpleRole) CanDelete() bool {
	return !r.IsSystem
}

// ModuleDefinition represents a system module that can have permissions assigned
// Examples: Products, HR, Notifications, Settings, etc.
type ModuleDefinition struct {
	ID          int64
	Name        string // PRODUCTS, HR, NOTIFICATIONS, SETTINGS, etc.
	DisplayName string
	Description string
	Icon        string // Icon name for UI display
	Path        string // URL path like "/dashboard/products"
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ModuleNames defines standard module names
// Note: These are duplicate definitions - see notification.go for the canonical versions
// const (
// 	ModuleProducts      = "PRODUCTS"
// 	ModuleCategories    = "CATEGORIES"
// 	ModuleHR            = "HR"
// 	ModuleNotifications = "NOTIFICATIONS"
// 	ModuleSettings      = "SETTINGS"
// 	ModuleUsers         = "USERS"
// 	ModuleRoles         = "ROLES"
// 	ModuleReports       = "REPORTS"
// 	ModuleThemeBuilder  = "THEMES"
// )

// Validate checks if module has required fields
func (m *ModuleDefinition) Validate() error {
	if m.Name == "" {
		return errors.New("module name is required")
	}

	if m.DisplayName == "" {
		return errors.New("module display name is required")
	}

	if len(m.Name) > 100 {
		return errors.New("module name must be less than 100 characters")
	}

	return nil
}

// PermissionLevel represents the level of access a role has to a module
type PermissionLevel struct {
	ID          int64
	Name        string // READ, EXPORT, WRITE, DELETE, ADMIN, NONE
	Value       int    // Numeric value for comparison hierarchy
	Description string
}

// PermissionLevelConstants defines permission levels
const (
	PermissionNone   = "NONE"   // No access (value 0)
	PermissionRead   = "READ"   // Read-only (value 1)
	PermissionExport = "EXPORT" // Read + export reports/files (value 2)
	PermissionWrite  = "WRITE"  // Read and write (value 3)
	PermissionDelete = "DELETE" // Read, write, and delete (value 4)
	PermissionAdmin  = "ADMIN"  // Full admin access (value 5)
)

// PermissionValues maps permission names to numeric values
var PermissionValues = map[string]int{
	PermissionNone:   0,
	PermissionRead:   1,
	PermissionExport: 2,
	PermissionWrite:  3,
	PermissionDelete: 4,
	PermissionAdmin:  5,
}

// Validate checks if permission level has required fields
func (p *PermissionLevel) Validate() error {
	if p.Name == "" {
		return errors.New("permission level name is required")
	}

	if p.Value < 0 {
		return errors.New("permission level value cannot be negative")
	}

	return nil
}

// CanRead checks if this permission level allows reading
func (p *PermissionLevel) CanRead() bool {
	return p.Value >= PermissionValues[PermissionRead]
}

// CanWrite checks if this permission level allows writing
func (p *PermissionLevel) CanWrite() bool {
	return p.Value >= PermissionValues[PermissionWrite]
}

// CanDelete checks if this permission level allows deleting
func (p *PermissionLevel) CanDelete() bool {
	return p.Value >= PermissionValues[PermissionDelete]
}

// IsAdmin checks if this is admin level access
func (p *PermissionLevel) IsAdmin() bool {
	return p.Value >= PermissionValues[PermissionAdmin]
}

// RolePermission represents the permission a role has for a specific module
type RolePermission struct {
	ID              int64
	TenantID        int64
	RoleID          int64
	ModuleID        int64
	PermissionLevel string // READ, WRITE, DELETE, ADMIN, NONE
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Validate checks if role permission has required fields
func (rp *RolePermission) Validate() error {
	if rp.RoleID == 0 {
		return errors.New("role ID is required")
	}

	if rp.ModuleID == 0 {
		return errors.New("module ID is required")
	}

	if rp.TenantID == 0 {
		return errors.New("tenant ID is required")
	}

	validLevels := []string{
		PermissionNone,
		PermissionRead,
		PermissionExport,
		PermissionWrite,
		PermissionDelete,
		PermissionAdmin,
	}

	isValid := false
	for _, level := range validLevels {
		if rp.PermissionLevel == level {
			isValid = true
			break
		}
	}

	if !isValid {
		return errors.New("invalid permission level")
	}

	return nil
}

// HasPermissionLevel checks if user permission meets or exceeds required level
func (rp *RolePermission) HasPermissionLevel(requiredLevel string) bool {
	requiredValue, exists := PermissionValues[requiredLevel]
	if !exists {
		return false
	}

	currentValue, exists := PermissionValues[rp.PermissionLevel]
	if !exists {
		return false
	}

	return currentValue >= requiredValue
}

// UserRole represents a user assigned to a role
type UserRole struct {
	ID         int64
	TenantID   int64
	UserID     int64
	RoleID     int64
	AssignedBy int64 // User ID of admin who assigned this role
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Validate checks if user role has required fields
func (ur *UserRole) Validate() error {
	if ur.UserID == 0 {
		return errors.New("user ID is required")
	}

	if ur.RoleID == 0 {
		return errors.New("role ID is required")
	}

	if ur.TenantID == 0 {
		return errors.New("tenant ID is required")
	}

	return nil
}
