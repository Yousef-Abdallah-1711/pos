package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"pos-saas/internal/domain"
	"strings"
)

// ProductRepository handles product data operations
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates new product repository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// CreateProduct inserts a new product into database
func (r *ProductRepository) CreateProduct(product *domain.Product) (*domain.Product, error) {
	fmt.Println("DEBUG: Repo.CreateProduct started")

	query := `
		INSERT INTO products (
			tenant_id, restaurant_id, category_id, sku, barcode,
			name_en, name_ar, description_en, description_ar,
			price, cost, discount_price, discount_percentage,
			calories, protein_g, carbs_g, fat_g, fiber_g, allergens,
			is_vegetarian, is_vegan, is_spicy, is_gluten_free,
			is_available, available_from, available_until, available_days,
			track_inventory, quantity_in_stock, low_stock_threshold, reorder_quantity,
			display_order, featured, main_image_url, status, created_by,
			name
		)
		SELECT
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30, $31, $32, $33, $34, $35, $36, $37
		WHERE EXISTS (SELECT 1 FROM restaurants WHERE id = $2 AND tenant_id = $1)
		  AND ($3::integer IS NULL OR EXISTS (
			SELECT 1 FROM categories WHERE id = $3 AND tenant_id = $1 AND restaurant_id = $2
		  ))
		RETURNING id, created_at, updated_at
	`

	// Marshal allergens to JSON
	var allergensJSON interface{} = nil
	if len(product.Allergens) > 0 {
		bytes, err := json.Marshal(product.Allergens)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal allergens: %w", err)
		}
		allergensJSON = string(bytes)
	}

	// Marshal available days to JSON
	var availableDaysJSON interface{} = nil
	if len(product.AvailableDays) > 0 {
		bytes, err := json.Marshal(product.AvailableDays)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal available days: %w", err)
		}
		availableDaysJSON = string(bytes)
	}

	// Handle nullable/optional fields
	var categoryID sql.NullInt64
	if product.CategoryID > 0 {
		categoryID = sql.NullInt64{Int64: int64(product.CategoryID), Valid: true}
	}

	var sku, barcode, nameAr, descEn, descAr sql.NullString
	if product.SKU != "" {
		sku = sql.NullString{String: product.SKU, Valid: true}
	}
	if product.Barcode != "" {
		barcode = sql.NullString{String: product.Barcode, Valid: true}
	}
	if product.NameAr != "" {
		nameAr = sql.NullString{String: product.NameAr, Valid: true}
	}
	if product.DescriptionEn != "" {
		descEn = sql.NullString{String: product.DescriptionEn, Valid: true}
	}
	if product.DescriptionAr != "" {
		descAr = sql.NullString{String: product.DescriptionAr, Valid: true}
	}

	fmt.Printf("DEBUG: Executing INSERT query for product: %s\n, user id %d", product.NameEn, product.CreatedBy)
	err := r.db.QueryRow(query,
		product.TenantID,
		product.RestaurantID,
		categoryID,
		sku,
		barcode,
		product.NameEn,
		nameAr,
		descEn,
		descAr,
		product.Price,
		product.Cost,
		product.DiscountPrice,
		product.DiscountPercentage,
		product.Calories,
		product.ProteinG,
		product.CarbsG,
		product.FatG,
		product.FiberG,
		allergensJSON,
		product.IsVegetarian,
		product.IsVegan,
		product.IsSpicy,
		product.IsGlutenFree,
		product.IsAvailable,
		product.AvailableFrom,
		product.AvailableUntil,
		availableDaysJSON,
		product.TrackInventory,
		product.QuantityInStock,
		product.LowStockThreshold,
		product.ReorderQuantity,
		product.DisplayOrder,
		product.Featured,
		product.MainImageURL,
		product.Status,
		product.CreatedBy,
		product.NameEn,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		fmt.Printf("ERROR: INSERT failed: %v\n", err)
		if err == sql.ErrNoRows {
			return nil, errors.New("restaurant or category not found in the requested scope")
		}
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.New("product with this SKU already exists")
		}
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	fmt.Printf("DEBUG: Insert successful. ID: %d. Returning product.\n", product.ID)
	return product, nil
}

// GetProductByID retrieves a product by ID
func (r *ProductRepository) GetProductByID(tenantID, restaurantID, productID int) (*domain.Product, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, category_id, sku, barcode,
			name_en, name_ar, description_en, description_ar,
			price, cost, discount_price, discount_percentage,
			calories, protein_g, carbs_g, fat_g, fiber_g, allergens,
			is_vegetarian, is_vegan, is_spicy, is_gluten_free,
			is_available, available_from, available_until, available_days,
			track_inventory, quantity_in_stock, low_stock_threshold, reorder_quantity,
			display_order, featured, main_image_url, status,
			created_by, updated_by, created_at, updated_at
		FROM products
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	product := &domain.Product{}
	var allergensJSON, availableDaysJSON []byte
	var sku, barcode, nameAr, descEn, descAr, mainImageURL sql.NullString
	var categoryID, createdBy, updatedBy sql.NullInt64

	// Careful scan for nullable columns
	err := r.db.QueryRow(query, productID, tenantID, restaurantID).Scan(
		&product.ID, &product.TenantID, &product.RestaurantID, &categoryID,
		&sku, &barcode,
		&product.NameEn, &nameAr, &descEn, &descAr,
		&product.Price, &product.Cost, &product.DiscountPrice, &product.DiscountPercentage,
		&product.Calories, &product.ProteinG, &product.CarbsG, &product.FatG, &product.FiberG, &allergensJSON,
		&product.IsVegetarian, &product.IsVegan, &product.IsSpicy, &product.IsGlutenFree,
		&product.IsAvailable, &product.AvailableFrom, &product.AvailableUntil, &availableDaysJSON,
		&product.TrackInventory, &product.QuantityInStock, &product.LowStockThreshold, &product.ReorderQuantity,
		&product.DisplayOrder, &product.Featured, &mainImageURL, &product.Status,
		&createdBy, &updatedBy, &product.CreatedAt, &product.UpdatedAt,
	)

	if categoryID.Valid {
		product.CategoryID = int(categoryID.Int64)
	}

	if err == sql.ErrNoRows {
		// Do not follow a scoped miss with an ID-only lookup: it would reveal
		// whether another tenant owns the requested product.
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Assign nullable fields back to struct
	product.SKU = sku.String
	product.Barcode = barcode.String
	product.NameAr = nameAr.String
	product.DescriptionEn = descEn.String
	product.DescriptionAr = descAr.String
	product.MainImageURL = mainImageURL.String
	if createdBy.Valid {
		product.CreatedBy = int(createdBy.Int64)
	}
	if updatedBy.Valid {
		val := int(updatedBy.Int64)
		product.UpdatedBy = &val
	}

	// Unmarshal JSON fields
	if len(allergensJSON) > 0 {
		json.Unmarshal(allergensJSON, &product.Allergens)
	}
	if len(availableDaysJSON) > 0 {
		json.Unmarshal(availableDaysJSON, &product.AvailableDays)
	}

	return product, nil
}

// ProductFilters represents query filters for products
type ProductFilters struct {
	CategoryID  int
	Status      string // 'active', 'inactive', 'discontinued'
	IsAvailable *bool
	SearchText  string
	SortBy      string // 'created_at', 'name_en', 'price'
	SortOrder   string // 'ASC' or 'DESC'
	Page        int
	PageSize    int
}

// ListProducts retrieves products with filters, pagination, and sorting
func (r *ProductRepository) ListProducts(tenantID, restaurantID int, filters *ProductFilters) (*domain.ProductListResponse, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, category_id, sku,
			name_en, name_ar, price, is_available, track_inventory, quantity_in_stock, low_stock_threshold, featured,
			display_order, main_image_url, status, created_at
		FROM products
		WHERE tenant_id = $1 AND restaurant_id = $2 AND status != 'deleted'
	`
	args := []interface{}{tenantID, restaurantID}
	argNum := 3

	// Apply filters
	if filters.CategoryID > 0 {
		query += fmt.Sprintf(" AND category_id = $%d", argNum)
		args = append(args, filters.CategoryID)
		argNum++
	}

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, filters.Status)
		argNum++
	}

	if filters.IsAvailable != nil {
		query += fmt.Sprintf(" AND is_available = $%d", argNum)
		args = append(args, *filters.IsAvailable)
		argNum++
	}

	if filters.SearchText != "" {
		query += fmt.Sprintf(" AND (name_en ILIKE $%d OR name_ar ILIKE $%d)", argNum, argNum+1)
		searchText := "%" + filters.SearchText + "%"
		args = append(args, searchText, searchText)
		argNum += 2
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS counted"
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Apply sorting
	sortField := "created_at"
	sortOrder := "DESC"
	if filters.SortBy != "" {
		sortField = filters.SortBy
	}
	if filters.SortOrder != "" {
		sortOrder = filters.SortOrder
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)

	// Apply pagination
	offset := (filters.Page - 1) * filters.PageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, filters.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	products := []domain.Product{}
	for rows.Next() {
		var p domain.Product
		var categoryID sql.NullInt64
		var sku, nameAr, mainImageURL sql.NullString

		err := rows.Scan(
			&p.ID, &p.TenantID, &p.RestaurantID, &categoryID, &sku,
			&p.NameEn, &nameAr, &p.Price, &p.IsAvailable, &p.TrackInventory, &p.QuantityInStock, &p.LowStockThreshold, &p.Featured,
			&p.DisplayOrder, &mainImageURL, &p.Status, &p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		if categoryID.Valid {
			p.CategoryID = int(categoryID.Int64)
		}
		p.SKU = sku.String
		p.NameAr = nameAr.String
		p.MainImageURL = mainImageURL.String

		products = append(products, p)
	}

	return &domain.ProductListResponse{
		Total:    total,
		Page:     filters.Page,
		PageSize: filters.PageSize,
		Products: products,
	}, nil
}

// UpdateProduct updates an existing product
func (r *ProductRepository) UpdateProduct(product *domain.Product) (*domain.Product, error) {
	fmt.Printf("DEBUG: Repo.UpdateProduct started for ID: %d\n", product.ID)
	var categoryID sql.NullInt64
	if product.CategoryID > 0 {
		categoryID = sql.NullInt64{Int64: int64(product.CategoryID), Valid: true}
	}

	nameAr := sql.NullString{String: product.NameAr, Valid: product.NameAr != ""}
	descEn := sql.NullString{String: product.DescriptionEn, Valid: product.DescriptionEn != ""}
	descAr := sql.NullString{String: product.DescriptionAr, Valid: product.DescriptionAr != ""}

	query := `
		UPDATE products
		SET
			category_id = $1,
			name_en = $2,
			name_ar = $3,
			description_en = $4,
			description_ar = $5,
			price = $6,
			cost = $7,
			discount_price = $8,
			is_available = $9,
			track_inventory = $10,
			quantity_in_stock = $11,
			low_stock_threshold = $12,
			featured = $13,
			is_vegetarian = $14,
			is_vegan = $15,
			is_spicy = $16,
			is_gluten_free = $17,
			updated_by = $18,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $19 AND tenant_id = $20 AND restaurant_id = $21
		  AND ($22::integer IS NULL OR EXISTS (
			SELECT 1 FROM categories WHERE id = $22 AND tenant_id = $20 AND restaurant_id = $21
		  ))
		RETURNING updated_at
	`

	fmt.Printf("DEBUG: Executing UPDATE query: %s\n", query)
	err := r.db.QueryRow(
		query,
		categoryID,
		product.NameEn,
		nameAr,
		descEn,
		descAr,
		product.Price,
		product.Cost,
		product.DiscountPrice,
		product.IsAvailable,
		product.TrackInventory,
		product.QuantityInStock,
		product.LowStockThreshold,
		product.Featured,
		product.IsVegetarian,
		product.IsVegan,
		product.IsSpicy,
		product.IsGlutenFree,
		product.UpdatedBy,
		product.ID,
		product.TenantID,
		product.RestaurantID,
		categoryID,
	).Scan(&product.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		fmt.Printf("ERROR: Update failed: %v\n", err)
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// DeleteProduct soft-deletes a product
func (r *ProductRepository) DeleteProduct(tenantID, restaurantID, productID int, deletedBy int) error {
	fmt.Printf("DEBUG: Repo.DeleteProduct started for ID: %d\n", productID)
	query := `
		UPDATE products
		SET status = 'deleted', updated_by = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`

	result, err := r.db.Exec(query, deletedBy, productID, tenantID, restaurantID)
	if err != nil {
		fmt.Printf("ERROR: Delete failed: %v\n", err)
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

// ListPublic returns active products for public menu
func (r *ProductRepository) ListPublic(tenantID, restaurantID int, lang string) ([]domain.Product, error) {
	query := `
		SELECT 
			id, tenant_id, restaurant_id, category_id, sku,
			name_en, name_ar, description_en, description_ar,
			price, discount_price, discount_percentage,
			calories, protein_g, carbs_g, fat_g, fiber_g, allergens,
			is_vegetarian, is_vegan, is_spicy, is_gluten_free,
			is_available, available_from, available_until,
			track_inventory, quantity_in_stock, low_stock_threshold,
			main_image_url, featured, status,
			created_at
		FROM products
		WHERE tenant_id = $1 AND restaurant_id = $2
		  AND status = 'active'
		  AND is_available = true
		ORDER BY display_order ASC, featured DESC, created_at DESC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list public products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		var allergens, sku, nameAr, descEn, descAr, mainImageURL sql.NullString
		// Use NullString for potential nullable columns
		err := rows.Scan(
			&p.ID, &p.TenantID, &p.RestaurantID, &p.CategoryID, &sku,
			&p.NameEn, &nameAr, &descEn, &descAr,
			&p.Price, &p.DiscountPrice, &p.DiscountPercentage,
			&p.Calories, &p.ProteinG, &p.CarbsG, &p.FatG, &p.FiberG, &allergens,
			&p.IsVegetarian, &p.IsVegan, &p.IsSpicy, &p.IsGlutenFree,
			&p.IsAvailable, &p.AvailableFrom, &p.AvailableUntil,
			&p.TrackInventory, &p.QuantityInStock, &p.LowStockThreshold,
			&mainImageURL, &p.Featured, &p.Status,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		p.SKU = sku.String
		p.NameAr = nameAr.String
		p.DescriptionEn = descEn.String
		p.DescriptionAr = descAr.String
		p.MainImageURL = mainImageURL.String
		p.Allergens = []string{} // Reset/ignore for now to avoid nil pointer issues if any logic depends on it
		products = append(products, p)
	}
	return products, nil
}

// ListByCategoryPublic returns products for specific category
func (r *ProductRepository) ListByCategoryPublic(tenantID, restaurantID, categoryID int, lang string) ([]domain.Product, error) {
	query := `
		SELECT 
			id, tenant_id, restaurant_id, category_id, sku,
			name_en, name_ar, description_en, description_ar,
			price, discount_price, discount_percentage,
			calories, protein_g, carbs_g, fat_g, fiber_g, allergens,
			is_vegetarian, is_vegan, is_spicy, is_gluten_free,
			is_available, available_from, available_until,
			track_inventory, quantity_in_stock, low_stock_threshold,
			main_image_url, featured, status,
			created_at
		FROM products
		WHERE tenant_id = $1 AND restaurant_id = $2 AND category_id = $3
		  AND status = 'active'
		  AND is_available = true
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		var allergens, sku, nameAr, descEn, descAr, mainImageURL sql.NullString
		err := rows.Scan(
			&p.ID, &p.TenantID, &p.RestaurantID, &p.CategoryID, &sku,
			&p.NameEn, &nameAr, &descEn, &descAr,
			&p.Price, &p.DiscountPrice, &p.DiscountPercentage,
			&p.Calories, &p.ProteinG, &p.CarbsG, &p.FatG, &p.FiberG, &allergens,
			&p.IsVegetarian, &p.IsVegan, &p.IsSpicy, &p.IsGlutenFree,
			&p.IsAvailable, &p.AvailableFrom, &p.AvailableUntil,
			&p.TrackInventory, &p.QuantityInStock, &p.LowStockThreshold,
			&mainImageURL, &p.Featured, &p.Status,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		p.SKU = sku.String
		p.NameAr = nameAr.String
		p.DescriptionEn = descEn.String
		p.DescriptionAr = descAr.String
		p.MainImageURL = mainImageURL.String
		products = append(products, p)
	}
	return products, nil
}

// SearchPublic searches products
func (r *ProductRepository) SearchPublic(tenantID, restaurantID int, queryStr string, lang string) ([]domain.Product, error) {
	searchTerm := "%" + queryStr + "%"
	query := `
		SELECT 
			id, tenant_id, restaurant_id, category_id, sku,
			name_en, name_ar, description_en, description_ar,
			price, discount_price, discount_percentage,
			calories, protein_g, carbs_g, fat_g, fiber_g, allergens,
			is_vegetarian, is_vegan, is_spicy, is_gluten_free,
			is_available, available_from, available_until,
			track_inventory, quantity_in_stock, low_stock_threshold,
			main_image_url, featured, status,
			created_at
		FROM products
		WHERE tenant_id = $1 AND restaurant_id = $2
		  AND status = 'active'
		  AND is_available = true
		  AND (name_en ILIKE $3 OR name_ar ILIKE $3 OR description_en ILIKE $3)
		ORDER BY display_order ASC
	`
	rows, err := r.db.Query(query, tenantID, restaurantID, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		var allergens, sku, nameAr, descEn, descAr, mainImageURL sql.NullString
		err := rows.Scan(
			&p.ID, &p.TenantID, &p.RestaurantID, &p.CategoryID, &sku,
			&p.NameEn, &nameAr, &descEn, &descAr,
			&p.Price, &p.DiscountPrice, &p.DiscountPercentage,
			&p.Calories, &p.ProteinG, &p.CarbsG, &p.FatG, &p.FiberG, &allergens,
			&p.IsVegetarian, &p.IsVegan, &p.IsSpicy, &p.IsGlutenFree,
			&p.IsAvailable, &p.AvailableFrom, &p.AvailableUntil,
			&p.TrackInventory, &p.QuantityInStock, &p.LowStockThreshold,
			&mainImageURL, &p.Featured, &p.Status,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		p.SKU = sku.String
		p.NameAr = nameAr.String
		p.DescriptionEn = descEn.String
		p.DescriptionAr = descAr.String
		p.MainImageURL = mainImageURL.String
		products = append(products, p)
	}
	return products, nil
}

// GetProductImages returns images
func (r *ProductRepository) GetProductImages(tenantID, restaurantID, productID int) ([]domain.ProductImage, error) {
	query := `SELECT pi.id, pi.product_id, pi.image_url, pi.alt_text_en, pi.alt_text_ar, pi.is_primary, pi.display_order
		FROM product_images pi
		JOIN products p ON p.id = pi.product_id
		WHERE pi.product_id = $1 AND p.tenant_id = $2 AND p.restaurant_id = $3
		ORDER BY pi.display_order`
	rows, err := r.db.Query(query, productID, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []domain.ProductImage
	for rows.Next() {
		var img domain.ProductImage
		var altEn, altAr sql.NullString
		err := rows.Scan(&img.ID, &img.ProductID, &img.ImageURL, &altEn, &altAr, &img.IsPrimary, &img.DisplayOrder)
		if err != nil {
			return nil, err
		}
		img.AltTextEn = altEn.String
		img.AltTextAr = altAr.String
		images = append(images, img)
	}
	return images, nil
}

// GetLowStockAlerts returns products with stock below their threshold
func (r *ProductRepository) GetLowStockAlerts(tenantID, restaurantID int) ([]domain.Product, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, category_id, sku,
			name_en, name_ar, price, quantity_in_stock, low_stock_threshold,
			main_image_url, status, created_at
		FROM products
		WHERE tenant_id = $1
			AND restaurant_id = $2
			AND status != 'deleted'
			AND track_inventory = true
			AND quantity_in_stock < low_stock_threshold
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock alerts: %w", err)
	}
	defer rows.Close()

	products := []domain.Product{}
	for rows.Next() {
		var p domain.Product
		var categoryID sql.NullInt64
		var sku, nameAr, mainImageURL sql.NullString

		err := rows.Scan(
			&p.ID, &p.TenantID, &p.RestaurantID, &categoryID, &sku,
			&p.NameEn, &nameAr, &p.Price, &p.QuantityInStock, &p.LowStockThreshold,
			&mainImageURL, &p.Status, &p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan low stock product: %w", err)
		}

		if categoryID.Valid {
			p.CategoryID = int(categoryID.Int64)
		}
		p.SKU = sku.String
		p.NameAr = nameAr.String
		p.MainImageURL = mainImageURL.String

		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating low stock products: %w", err)
	}

	return products, nil
}
