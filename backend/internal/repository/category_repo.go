package repository

import (
	"database/sql"
	"pos-saas/internal/domain"
)

// CategoryRepository handles category data operations
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// ListCategories retrieves all active categories for a restaurant
func (r *CategoryRepository) ListCategories(tenantID, restaurantID int) ([]domain.Category, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, name, name_ar,
			description, description_ar, icon_url, display_order,
			is_active, created_at, updated_at
		FROM categories
		WHERE tenant_id = $1 AND restaurant_id = $2 AND is_active = true
		ORDER BY display_order ASC, name ASC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var cat domain.Category
		var nameAr, description, descriptionAr, iconURL sql.NullString

		err := rows.Scan(
			&cat.ID, &cat.TenantID, &cat.RestaurantID,
			&cat.Name, &nameAr,
			&description, &descriptionAr,
			&iconURL, &cat.DisplayOrder, &cat.IsActive,
			&cat.CreatedAt, &cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		cat.NameAr = nameAr.String
		cat.Description = description.String
		cat.DescriptionAr = descriptionAr.String
		cat.IconURL = iconURL.String

		categories = append(categories, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return empty array instead of nil if no categories found
	if categories == nil {
		categories = []domain.Category{}
	}

	return categories, nil
}

// GetCategoryByID retrieves a single category by ID
func (r *CategoryRepository) GetCategoryByID(tenantID, restaurantID, id int) (*domain.Category, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, name, name_ar,
			description, description_ar, icon_url, display_order,
			is_active, created_at, updated_at
		FROM categories
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	var cat domain.Category
	var nameAr, description, descriptionAr, iconURL sql.NullString

	err := r.db.QueryRow(query, id, tenantID, restaurantID).Scan(
		&cat.ID, &cat.TenantID, &cat.RestaurantID,
		&cat.Name, &nameAr,
		&description, &descriptionAr,
		&iconURL, &cat.DisplayOrder, &cat.IsActive,
		&cat.CreatedAt, &cat.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	cat.NameAr = nameAr.String
	cat.Description = description.String
	cat.DescriptionAr = descriptionAr.String
	cat.IconURL = iconURL.String

	return &cat, nil
}

// CreateCategory creates a new category
func (r *CategoryRepository) CreateCategory(category *domain.Category) (int, error) {
	query := `
		INSERT INTO categories (
			tenant_id, restaurant_id, name, name_ar,
			description, description_ar, display_order, is_active
		)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8
		WHERE EXISTS (SELECT 1 FROM restaurants WHERE id = $2 AND tenant_id = $1)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(
		query,
		category.TenantID,
		category.RestaurantID,
		category.Name,
		category.NameAr,
		category.Description,
		category.DescriptionAr,
		category.DisplayOrder,
		category.IsActive,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, sql.ErrNoRows
	}
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateCategory updates an existing category
func (r *CategoryRepository) UpdateCategory(category *domain.Category) error {
	query := `
		UPDATE categories
		SET
			name = $1,
			name_ar = $2,
			description = $3,
			description_ar = $4,
			display_order = $5,
			is_active = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND tenant_id = $8 AND restaurant_id = $9
	`

	result, err := r.db.Exec(
		query,
		category.Name,
		category.NameAr,
		category.Description,
		category.DescriptionAr,
		category.DisplayOrder,
		category.IsActive,
		category.ID,
		category.TenantID,
		category.RestaurantID,
	)

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteCategory deletes a category (soft delete by setting is_active to false)
func (r *CategoryRepository) DeleteCategory(tenantID, restaurantID, id int) error {
	query := `
		UPDATE categories
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	result, err := r.db.Exec(query, id, tenantID, restaurantID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ListPublic returns active categories for public menu
func (r *CategoryRepository) ListPublic(tenantID, restaurantID int, lang string) ([]domain.Category, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, name, name_ar,
			description, description_ar, icon_url, display_order,
			is_active, created_at, updated_at
		FROM categories
		WHERE tenant_id = $1 AND restaurant_id = $2 AND is_active = true
		ORDER BY display_order ASC, name ASC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var cat domain.Category
		var nameAr, description, descriptionAr, iconURL sql.NullString

		err := rows.Scan(
			&cat.ID, &cat.TenantID, &cat.RestaurantID,
			&cat.Name, &nameAr,
			&description, &descriptionAr,
			&iconURL, &cat.DisplayOrder, &cat.IsActive,
			&cat.CreatedAt, &cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		cat.NameAr = nameAr.String
		cat.Description = description.String
		cat.DescriptionAr = descriptionAr.String
		cat.IconURL = iconURL.String

		categories = append(categories, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return empty array instead of nil if no categories found
	if categories == nil {
		categories = []domain.Category{}
	}

	return categories, nil
}
