package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// AuditLog represents a record of permission changes for compliance and security
type PermissionAuditLog struct {
	ID        int64
	TenantID  int64
	AdminID   int64  // User ID of admin who made the change
	UserID    *int64 // User ID affected by the change (nullable)
	RoleID    *int64 // Role ID affected by the change (nullable)
	ModuleID  *int64 // Module ID affected by the change (nullable)
	Action    string // CREATE_ROLE, UPDATE_ROLE, DELETE_ROLE, ASSIGN_PERMISSION, REVOKE_PERMISSION, ASSIGN_USER_ROLE, REVOKE_USER_ROLE
	Details   AuditDetails
	IPAddress string // Source IP address
	UserAgent string // Browser/client user agent
	CreatedAt time.Time
}

// AuditDetails contains additional metadata about the change
type AuditDetails struct {
	Reason      string                 `json:"reason,omitempty"`
	OldValue    string                 `json:"old_value,omitempty"`
	NewValue    string                 `json:"new_value,omitempty"`
	Description string                 `json:"description,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
}

// Action type constants
const (
	ActionCreateRole       = "CREATE_ROLE"
	ActionUpdateRole       = "UPDATE_ROLE"
	ActionDeleteRole       = "DELETE_ROLE"
	ActionAssignPermission = "ASSIGN_PERMISSION"
	ActionRevokePermission = "REVOKE_PERMISSION"
	ActionAssignUserRole   = "ASSIGN_USER_ROLE"
	ActionRevokeUserRole   = "REVOKE_USER_ROLE"
	ActionCreateModule     = "CREATE_MODULE"
	ActionUpdateModule     = "UPDATE_MODULE"
	ActionDeleteModule     = "DELETE_MODULE"
)

// ValidActions defines all valid audit log actions
var ValidActions = []string{
	ActionCreateRole,
	ActionUpdateRole,
	ActionDeleteRole,
	ActionAssignPermission,
	ActionRevokePermission,
	ActionAssignUserRole,
	ActionRevokeUserRole,
	ActionCreateModule,
	ActionUpdateModule,
	ActionDeleteModule,
}

// Validate checks if audit log has required fields
func (pal *PermissionAuditLog) Validate() error {
	if pal.TenantID == 0 {
		return errors.New("tenant ID is required")
	}

	if pal.AdminID == 0 {
		return errors.New("admin ID is required")
	}

	if pal.Action == "" {
		return errors.New("action is required")
	}

	// Validate action is in allowed list
	isValidAction := false
	for _, validAction := range ValidActions {
		if pal.Action == validAction {
			isValidAction = true
			break
		}
	}

	if !isValidAction {
		return errors.New("invalid action type")
	}

	return nil
}

// String returns a human-readable description of the audit log
func (pal *PermissionAuditLog) String() string {
	description := ""

	switch pal.Action {
	case ActionCreateRole:
		description = "Created new role"
	case ActionUpdateRole:
		description = "Updated role"
	case ActionDeleteRole:
		description = "Deleted role"
	case ActionAssignPermission:
		description = "Assigned permission to role"
	case ActionRevokePermission:
		description = "Revoked permission from role"
	case ActionAssignUserRole:
		description = "Assigned role to user"
	case ActionRevokeUserRole:
		description = "Revoked role from user"
	case ActionCreateModule:
		description = "Created new module"
	case ActionUpdateModule:
		description = "Updated module"
	case ActionDeleteModule:
		description = "Deleted module"
	default:
		description = "Permission change"
	}

	return description
}

// Scan implements sql.Scanner interface for JSONB column
func (ad *AuditDetails) Scan(value interface{}) error {
	if value == nil {
		*ad = AuditDetails{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion []byte failed")
	}

	return json.Unmarshal(bytes, &ad)
}

// Value implements driver.Valuer interface for JSONB column
func (ad AuditDetails) Value() (driver.Value, error) {
	return json.Marshal(ad)
}

// NewPermissionAuditLog creates a new audit log entry
func NewPermissionAuditLog(
	tenantID int64,
	adminID int64,
	action string,
) *PermissionAuditLog {
	return &PermissionAuditLog{
		TenantID:  tenantID,
		AdminID:   adminID,
		Action:    action,
		Details:   AuditDetails{},
		CreatedAt: time.Now(),
	}
}
