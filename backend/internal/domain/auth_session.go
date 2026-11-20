package domain

// AuthSessionProfile contains only the tenant-scoped profile fields needed by
// browser portals to render an authenticated account area.
type AuthSessionProfile struct {
	UserID         int    `json:"id"`
	Email          string `json:"email"`
	DisplayName    string `json:"name"`
	TenantID       int    `json:"tenant_id"`
	RestaurantID   *int   `json:"restaurant_id,omitempty"`
	RestaurantName string `json:"restaurant_name,omitempty"`
	Role           string `json:"role"`
}
