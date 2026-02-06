package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"time"
)

// SuperAdmin represents a SaaS platform super administrator
type SuperAdmin struct {
	ID                  int64      `json:"id"`
	UserID              int64      `json:"user_id"`
	SecretKey           *string    `json:"-"` // Never expose in JSON
	IPWhitelist         []string   `json:"ip_whitelist,omitempty"`
	IsActive            bool       `json:"is_active"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP         *string    `json:"last_login_ip,omitempty"`
	FailedLoginAttempts int        `json:"failed_login_attempts"`
	LockedUntil         *time.Time `json:"locked_until,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	// Joined fields (not stored in super_admins table)
	User *User `json:"user,omitempty"`
}

// IsLocked checks if the super admin account is currently locked
func (sa *SuperAdmin) IsLocked() bool {
	if sa.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*sa.LockedUntil)
}

// IsIPAllowed checks if the given IP is in the whitelist
// Returns true if whitelist is empty (all IPs allowed)
func (sa *SuperAdmin) IsIPAllowed(ip string) bool {
	if len(sa.IPWhitelist) == 0 {
		return true
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, allowed := range sa.IPWhitelist {
		// Check if it's a CIDR range
		if _, cidr, err := net.ParseCIDR(allowed); err == nil {
			if cidr.Contains(parsedIP) {
				return true
			}
		} else {
			// Direct IP comparison
			if net.ParseIP(allowed) != nil && allowed == ip {
				return true
			}
		}
	}
	return false
}

// ImpersonationSession represents a super admin impersonation session
type ImpersonationSession struct {
	ID                 int64      `json:"id"`
	SuperAdminID       int64      `json:"super_admin_id"`
	TargetUserID       int64      `json:"target_user_id"`
	TargetTenantID     int64      `json:"target_tenant_id"`
	ImpersonationToken string     `json:"-"` // Never expose raw token
	TokenHash          string     `json:"-"` // For validation
	IPAddress          string     `json:"ip_address"`
	UserAgent          string     `json:"user_agent,omitempty"`
	StartedAt          time.Time  `json:"started_at"`
	ExpiresAt          time.Time  `json:"expires_at"`
	EndedAt            *time.Time `json:"ended_at,omitempty"`
	IsActive           bool       `json:"is_active"`
	EndReason          *string    `json:"end_reason,omitempty"`

	// Joined fields
	SuperAdmin   *SuperAdmin `json:"super_admin,omitempty"`
	TargetUser   *User       `json:"target_user,omitempty"`
	TargetTenant *Tenant     `json:"target_tenant,omitempty"`
}

// ImpersonationEndReason constants
const (
	ImpersonationEndManual  = "MANUAL"
	ImpersonationEndExpired = "EXPIRED"
	ImpersonationEndForced  = "FORCED"
)

// IsExpired checks if the impersonation session has expired
func (is *ImpersonationSession) IsExpired() bool {
	return time.Now().After(is.ExpiresAt)
}

// IsValid checks if the session is still valid (active and not expired)
func (is *ImpersonationSession) IsValid() bool {
	return is.IsActive && !is.IsExpired()
}

// HashToken generates a SHA-256 hash of the token
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// SuperAdminAuditLog represents an audit log entry for super admin actions
type SuperAdminAuditLog struct {
	ID           int64                  `json:"id"`
	SuperAdminID int64                  `json:"super_admin_id"`
	Action       string                 `json:"action"`
	TargetType   *string                `json:"target_type,omitempty"`
	TargetID     *int64                 `json:"target_id,omitempty"`
	Details      map[string]interface{} `json:"details,omitempty"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`

	// Joined fields
	SuperAdmin *SuperAdmin `json:"super_admin,omitempty"`
}

// Audit action constants
const (
	AuditActionLoginSuccess      = "LOGIN_SUCCESS"
	AuditActionLoginFailed       = "LOGIN_FAILED"
	AuditActionLogout            = "LOGOUT"
	AuditActionImpersonateStart  = "IMPERSONATE_START"
	AuditActionImpersonateEnd    = "IMPERSONATE_END"
	AuditActionViewTenant        = "VIEW_TENANT"
	AuditActionViewUser          = "VIEW_USER"
	AuditActionViewTenantList    = "VIEW_TENANT_LIST"
	AuditActionViewUserList      = "VIEW_USER_LIST"
	AuditActionDisableTenant     = "DISABLE_TENANT"
	AuditActionEnableTenant      = "ENABLE_TENANT"
	AuditActionViewAuditLogs     = "VIEW_AUDIT_LOGS"
	AuditActionExportData        = "EXPORT_DATA"
	AuditActionUpdateIPWhitelist = "UPDATE_IP_WHITELIST"
	AuditActionUpdateSettings    = "UPDATE_SETTINGS"
	AuditActionAccessDashboard   = "ACCESS_DASHBOARD"
)

// Audit target type constants
const (
	AuditTargetTenant        = "TENANT"
	AuditTargetUser          = "USER"
	AuditTargetSystem        = "SYSTEM"
	AuditTargetImpersonation = "IMPERSONATION"
)

// TenantWithStats represents a tenant with additional statistics for super admin view
type TenantWithStats struct {
	Tenant
	TotalUsers       int        `json:"total_users"`
	TotalRestaurants int        `json:"total_restaurants"`
	TotalOrders      int        `json:"total_orders"`
	LastActivityAt   *time.Time `json:"last_activity_at,omitempty"`
}

// UserWithTenant represents a user with their tenant information for super admin view
type UserWithTenant struct {
	User
	TenantName string `json:"tenant_name"`
	TenantSlug string `json:"tenant_slug"`
}

// SuperAdminLoginRequest represents the login request for super admin
type SuperAdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SuperAdminLoginResponse represents the login response for super admin
type SuperAdminLoginResponse struct {
	SuperAdmin *SuperAdmin `json:"super_admin"`
	User       *User       `json:"user"`
	Token      string      `json:"token"`
	ExpiresAt  time.Time   `json:"expires_at"`
}

// ImpersonationRequest represents a request to impersonate a user
type ImpersonationRequest struct {
	TargetUserID int64 `json:"target_user_id"`
}

// ImpersonationResponse represents the response when starting impersonation
type ImpersonationResponse struct {
	Token              string    `json:"token"`
	TargetUser         *User     `json:"target_user"`
	TargetTenantID     int64     `json:"target_tenant_id"`
	TargetRestaurantID *int64    `json:"target_restaurant_id,omitempty"`
	ExpiresAt          time.Time `json:"expires_at"`
	DashboardURL       string    `json:"dashboard_url"`
}

// SuperAdminDashboardStats represents statistics for the super admin dashboard
type SuperAdminDashboardStats struct {
	TotalTenants         int       `json:"total_tenants"`
	ActiveTenants        int       `json:"active_tenants"`
	TotalUsers           int       `json:"total_users"`
	TotalOrders          int       `json:"total_orders"`
	NewTenantsToday      int       `json:"new_tenants_today"`
	NewUsersToday        int       `json:"new_users_today"`
	ActiveImpersonations int       `json:"active_impersonations"`
	LastUpdated          time.Time `json:"last_updated"`
}

// AuditLogFilter represents filters for querying audit logs
type AuditLogFilter struct {
	SuperAdminID *int64     `json:"super_admin_id,omitempty"`
	Action       *string    `json:"action,omitempty"`
	TargetType   *string    `json:"target_type,omitempty"`
	TargetID     *int64     `json:"target_id,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	IPAddress    *string    `json:"ip_address,omitempty"`
}

// TenantFilter represents filters for querying tenants
type TenantFilter struct {
	Search    *string    `json:"search,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// UserFilter represents filters for querying users across all tenants
type UserFilter struct {
	Search       *string    `json:"search,omitempty"`
	TenantID     *int64     `json:"tenant_id,omitempty"`
	Role         *string    `json:"role,omitempty"`
	IsSuperAdmin *bool      `json:"is_super_admin,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
}

// Pagination represents pagination parameters
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, total int64, page, pageSize int) *PaginatedResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
