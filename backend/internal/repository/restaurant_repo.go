package repository

import (
	"database/sql"
	"fmt"
	"pos-saas/internal/domain"
)

// RestaurantRepository handles restaurant data operations
type RestaurantRepository struct {
	db *sql.DB
}

// NewRestaurantRepository creates new restaurant repository
func NewRestaurantRepository(db *sql.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

// GetBySlug retrieves a restaurant by slug
func (r *RestaurantRepository) GetBySlug(slug string) (*domain.Restaurant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("restaurant repository unavailable: database not initialized")
	}

	query := `
		SELECT 
			id, tenant_id, name, slug, description, logo_url, hero_image_url,
			email, phone, address, city, theme,
			website_enabled, pos_enabled, delivery_enabled, reservation_enabled,
			status, created_at, updated_at
		FROM restaurants
		WHERE slug = $1
	`

	restaurant := &domain.Restaurant{}
	err := r.db.QueryRow(query, slug).Scan(
		&restaurant.ID,
		&restaurant.TenantID,
		&restaurant.Name,
		&restaurant.Slug,
		&restaurant.Description,
		&restaurant.LogoURL,
		&restaurant.HeroImageURL,
		&restaurant.Email,
		&restaurant.Phone,
		&restaurant.Address,
		&restaurant.City,
		&restaurant.Theme,
		&restaurant.WebsiteEnabled,
		&restaurant.POSEnabled,
		&restaurant.DeliveryEnabled,
		&restaurant.ReservationEnabled,
		&restaurant.Status,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("restaurant not found")
		}
		return nil, fmt.Errorf("failed to get restaurant: %w", err)
	}

	return restaurant, nil
}
