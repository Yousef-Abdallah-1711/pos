package domain

import "time"

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeLowStock   NotificationType = "low_stock"
	NotificationTypeOrder      NotificationType = "order"
	NotificationTypeEmployee   NotificationType = "employee"
	NotificationTypeLeave      NotificationType = "leave"
	NotificationTypeInventory  NotificationType = "inventory"
	NotificationTypeAttendance NotificationType = "attendance"
	NotificationTypeSystem     NotificationType = "system"
)

// NotificationModule represents the module that generated the notification
type NotificationModule string

const (
	ModuleProducts  NotificationModule = "products"
	ModuleOrders    NotificationModule = "orders"
	ModuleHR        NotificationModule = "hr"
	ModuleInventory NotificationModule = "inventory"
	ModuleSystem    NotificationModule = "system"
)

// NotificationPriority represents the priority level
type NotificationPriority string

const (
	PriorityLow      NotificationPriority = "low"
	PriorityNormal   NotificationPriority = "normal"
	PriorityHigh     NotificationPriority = "high"
	PriorityCritical NotificationPriority = "critical"
)

// Notification represents a system notification
type Notification struct {
	ID                int                  `json:"id"`
	TenantID          int                  `json:"tenant_id"`
	RestaurantID      int                  `json:"restaurant_id"`
	UserID            int                  `json:"user_id"`
	Type              NotificationType     `json:"type"`
	Module            NotificationModule   `json:"module"`
	Title             string               `json:"title"`
	Message           string               `json:"message"`
	Description       *string              `json:"description,omitempty"`
	IsRead            bool                 `json:"is_read"`
	ReadAt            *time.Time           `json:"read_at,omitempty"`
	RelatedEntityType *string              `json:"related_entity_type,omitempty"`
	RelatedEntityID   *int                 `json:"related_entity_id,omitempty"`
	Priority          NotificationPriority `json:"priority"`
	ActionURL         *string              `json:"action_url,omitempty"`
	ActionLabel       *string              `json:"action_label,omitempty"`
	IconName          *string              `json:"icon_name,omitempty"`
	Color             *string              `json:"color,omitempty"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	ExpiresAt         *time.Time           `json:"expires_at,omitempty"`
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	ID                             int       `json:"id"`
	UserID                         int       `json:"user_id"`
	LowStockEnabled                bool      `json:"low_stock_enabled"`
	OrderNotificationsEnabled      bool      `json:"order_notifications_enabled"`
	EmployeeNotificationsEnabled   bool      `json:"employee_notifications_enabled"`
	LeaveNotificationsEnabled      bool      `json:"leave_notifications_enabled"`
	InventoryNotificationsEnabled  bool      `json:"inventory_notifications_enabled"`
	AttendanceNotificationsEnabled bool      `json:"attendance_notifications_enabled"`
	SystemNotificationsEnabled     bool      `json:"system_notifications_enabled"`
	EmailNotifications             bool      `json:"email_notifications"`
	PushNotifications              bool      `json:"push_notifications"`
	InAppNotifications             bool      `json:"in_app_notifications"`
	SMSNotifications               bool      `json:"sms_notifications"`
	QuietHoursEnabled              bool      `json:"quiet_hours_enabled"`
	QuietHoursStart                *string   `json:"quiet_hours_start,omitempty"` // HH:MM format
	QuietHoursEnd                  *string   `json:"quiet_hours_end,omitempty"`   // HH:MM format
	AutoClearAfterDays             int       `json:"auto_clear_after_days"`
	CreatedAt                      time.Time `json:"created_at"`
	UpdatedAt                      time.Time `json:"updated_at"`
}

// NotificationListResponse represents paginated notification list
type NotificationListResponse struct {
	Total         int            `json:"total"`
	Unread        int            `json:"unread"`
	Page          int            `json:"page"`
	PageSize      int            `json:"page_size"`
	Notifications []Notification `json:"notifications"`
}

// NotificationFilters represents filters for notification queries
type NotificationFilters struct {
	Module    string
	Type      string
	IsRead    *bool
	Priority  string
	Page      int
	PageSize  int
	SortBy    string // "created_at", "updated_at", "priority"
	SortOrder string // "ASC", "DESC"
}

// CreateNotificationRequest represents a request to create a notification
type CreateNotificationRequest struct {
	Type              NotificationType     `json:"type" binding:"required"`
	Module            NotificationModule   `json:"module" binding:"required"`
	Title             string               `json:"title" binding:"required"`
	Message           string               `json:"message" binding:"required"`
	Description       *string              `json:"description,omitempty"`
	RelatedEntityType *string              `json:"related_entity_type,omitempty"`
	RelatedEntityID   *int                 `json:"related_entity_id,omitempty"`
	Priority          NotificationPriority `json:"priority,omitempty"`
	ActionURL         *string              `json:"action_url,omitempty"`
	ActionLabel       *string              `json:"action_label,omitempty"`
	IconName          *string              `json:"icon_name,omitempty"`
	Color             *string              `json:"color,omitempty"`
	UserID            int                  `json:"user_id" binding:"required"` // Target user
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	Total      int            `json:"total"`
	Unread     int            `json:"unread"`
	Read       int            `json:"read"`
	ByModule   map[string]int `json:"by_module"`
	ByPriority map[string]int `json:"by_priority"`
}
