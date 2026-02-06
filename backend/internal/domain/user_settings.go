package domain

import "time"

// UserSettings represents user preferences and settings
type UserSettings struct {
	ID                   int64     `json:"id"`
	UserID               int64     `json:"user_id"`
	Language             string    `json:"language"` // 'en' or 'ar'
	Theme                string    `json:"theme"`    // 'light', 'dark', 'system'
	PrimaryColor         string    `json:"primary_color"`
	SecondaryColor       string    `json:"secondary_color"`
	AccentColor          string    `json:"accent_color"`
	NotificationsEnabled bool      `json:"notifications_enabled"`
	EmailNotifications   bool      `json:"email_notifications"`
	SMSNotifications     bool      `json:"sms_notifications"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// SettingsAuditLog tracks changes to user settings
type SettingsAuditLog struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	SettingName string    `json:"setting_name"`
	OldValue    *string   `json:"old_value"`
	NewValue    *string   `json:"new_value"`
	ChangedBy   int64     `json:"changed_by"`
	ChangedAt   time.Time `json:"changed_at"`
	IPAddress   string    `json:"ip_address"`
}

// PasswordChangeHistory tracks password changes
type PasswordChangeHistory struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	PasswordHash string    `json:"-"` // Never expose hash
	ChangedAt    time.Time `json:"changed_at"`
	IPAddress    string    `json:"ip_address"`
}

// EmailChangeRequest tracks email change requests with verification
type EmailChangeRequest struct {
	ID                int64      `json:"id"`
	UserID            int64      `json:"user_id"`
	NewEmail          string     `json:"new_email"`
	VerificationToken string     `json:"-"` // Never expose token
	Verified          bool       `json:"verified"`
	VerifiedAt        *time.Time `json:"verified_at"`
	RequestedAt       time.Time  `json:"requested_at"`
	ExpiresAt         time.Time  `json:"expires_at"`
}

// UserProfile represents basic user profile information
type UserProfile struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	AvatarURL *string   `json:"avatar_url"`
	Phone     *string   `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
