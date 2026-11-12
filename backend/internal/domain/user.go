package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int        `json:"id"`
	TenantID     int        `json:"tenant_id"`
	RestaurantID *int       `json:"restaurant_id,omitempty"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Name         string     `json:"name"`
	Phone        string     `json:"phone"`
	AvatarURL    string     `json:"avatar_url"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	IsSuperAdmin bool       `json:"is_super_admin"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// VerifyPassword checks if the provided password matches the user's password hash
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type RegisterRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RestaurantName string `json:"restaurant_name"`
	RestaurantSlug string `json:"restaurant_slug"` // Subdomain for website (e.g., "my-restaurant")
	Phone          string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token,omitempty"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// TenantInfo holds information about a tenant for multi-tenant login
type TenantInfo struct {
	TenantID            int      `json:"tenant_id"`
	TenantName          string   `json:"tenant_name"`
	UserID              int      `json:"user_id"`
	Roles               []string `json:"roles"`
	DefaultRestaurantID *int     `json:"default_restaurant_id,omitempty"`
}

// MultiTenantLoginResponse is returned when a user has accounts in multiple tenants
type MultiTenantLoginResponse struct {
	Success         bool         `json:"success"`
	MultipleTenants bool         `json:"multiple_tenants"`
	Tenants         []TenantInfo `json:"tenants"`
	Message         string       `json:"message,omitempty"`
	// SelectionToken proves the password step succeeded; it must be echoed to /auth/login/confirm.
	SelectionToken string `json:"selection_token,omitempty"`
}

// TenantSelectionRequest is sent to /api/v1/auth/login/confirm
type TenantSelectionRequest struct {
	Email          string `json:"email"`
	TenantID       int    `json:"tenant_id"`
	SelectionToken string `json:"selection_token"`
}
