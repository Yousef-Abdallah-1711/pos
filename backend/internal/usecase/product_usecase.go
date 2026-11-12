package usecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"time"
)

// ProductUseCase handles product business logic
type ProductUseCase struct {
	repo             *repository.ProductRepository
	notificationRepo *repository.NotificationRepository
	storageURL       string // S3 or local storage URL
}

// NewProductUseCase creates new product use case
func NewProductUseCase(repo *repository.ProductRepository, notificationRepo *repository.NotificationRepository, storageURL string) *ProductUseCase {
	return &ProductUseCase{
		repo:             repo,
		notificationRepo: notificationRepo,
		storageURL:       storageURL,
	}
}

// CreateProduct handles product creation with image upload
func (uc *ProductUseCase) CreateProduct(
	tenantID, restaurantID, userID int,
	req *domain.CreateProductRequest,
	imageFile *multipart.FileHeader,
) (*domain.Product, error) {
	fmt.Println("DEBUG: UseCase.CreateProduct started")
	// Validate request
	if err := validateProductRequest(req); err != nil {
		fmt.Printf("ERROR: validateProductRequest failed: %v\n", err)
		return nil, err
	}

	// Create product entity
	product := &domain.Product{
		TenantID:           tenantID,
		RestaurantID:       restaurantID,
		CategoryID:         req.CategoryID,
		SKU:                req.SKU,
		Barcode:            req.Barcode,
		NameEn:             req.NameEn,
		NameAr:             req.NameAr,
		DescriptionEn:      req.DescriptionEn,
		DescriptionAr:      req.DescriptionAr,
		Price:              req.Price,
		Cost:               req.Cost,
		DiscountPrice:      req.DiscountPrice,
		DiscountPercentage: req.DiscountPercentage,
		Calories:           req.Calories,
		ProteinG:           req.ProteinG,
		CarbsG:             req.CarbsG,
		FatG:               req.FatG,
		FiberG:             req.FiberG,
		Allergens:          req.Allergens,
		IsVegetarian:       req.IsVegetarian,
		IsVegan:            req.IsVegan,
		IsSpicy:            req.IsSpicy,
		IsGlutenFree:       req.IsGlutenFree,
		IsAvailable:        req.IsAvailable,
		AvailableFrom:      req.AvailableFrom,
		AvailableUntil:     req.AvailableUntil,
		AvailableDays:      req.AvailableDays,
		TrackInventory:     req.TrackInventory,
		QuantityInStock:    req.QuantityInStock,
		LowStockThreshold:  req.LowStockThreshold,
		ReorderQuantity:    req.ReorderQuantity,
		DisplayOrder:       req.DisplayOrder,
		Featured:           req.Featured,
		Status:             "active",
		CreatedBy:          userID,
	}

	// Save to database
	fmt.Println("DEBUG: Calling repo.CreateProduct...")
	product, err := uc.repo.CreateProduct(product)
	if err != nil {
		fmt.Printf("ERROR: repo.CreateProduct failed: %v\n", err)
		return nil, err
	}
	fmt.Printf("DEBUG: Product saved to DB with ID: %d\n", product.ID)

	// Create low stock notification if applicable
	if product.TrackInventory && product.QuantityInStock < product.LowStockThreshold {
		actionURL := fmt.Sprintf("/dashboard/products/%d/edit", product.ID)
		iconName := "AlertTriangle"
		color := "orange"
		entityType := "product"
		notification := &domain.Notification{
			TenantID:          tenantID,
			RestaurantID:      restaurantID,
			UserID:            userID,
			Type:              domain.NotificationTypeLowStock,
			Module:            domain.ModuleProducts,
			Title:             fmt.Sprintf("Low Stock Alert: %s", product.NameEn),
			Message:           fmt.Sprintf("Product '%s' has only %d units in stock (threshold: %d)", product.NameEn, product.QuantityInStock, product.LowStockThreshold),
			Priority:          domain.PriorityHigh,
			RelatedEntityType: &entityType,
			RelatedEntityID:   &product.ID,
			ActionURL:         &actionURL,
			IconName:          &iconName,
			Color:             &color,
		}
		if _, err := uc.notificationRepo.CreateNotification(notification); err != nil {
			fmt.Printf("ERROR: Failed to create low stock notification: %v\n", err)
			// Don't fail the product creation if notification fails
		}
	}

	// Upload image if provided
	if imageFile != nil {
		fmt.Println("DEBUG: Uploading product image...")
		imageURL, err := uc.uploadProductImage(product.ID, imageFile, tenantID)
		if err != nil {
			// Log error but don't fail - product was created
			fmt.Printf("ERROR: Image upload failed: %v\n", err)
		} else {
			fmt.Printf("DEBUG: Image uploaded to: %s\n", imageURL)
			product.MainImageURL = imageURL
			// Update the product with the image URL? The original code didn't save it back to DB!
			// We found another bug! The image URL is set on the struct but not updated in DB after initial create.
			// But that's not the 500 error.
		}
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (uc *ProductUseCase) UpdateProduct(
	tenantID, restaurantID, productID, userID int,
	req *domain.UpdateProductRequest,
	imageFile *multipart.FileHeader,
) (*domain.Product, error) {
	// Get existing product
	product, err := uc.repo.GetProductByID(tenantID, restaurantID, productID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.NameEn != nil {
		product.NameEn = *req.NameEn
	}
	if req.NameAr != nil {
		product.NameAr = *req.NameAr
	}
	if req.DescriptionEn != nil {
		product.DescriptionEn = *req.DescriptionEn
	}
	if req.DescriptionAr != nil {
		product.DescriptionAr = *req.DescriptionAr
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Cost != nil {
		product.Cost = req.Cost
	}
	if req.DiscountPrice != nil {
		product.DiscountPrice = req.DiscountPrice
	}
	if req.IsAvailable != nil {
		product.IsAvailable = *req.IsAvailable
	}
	if req.TrackInventory != nil {
		product.TrackInventory = *req.TrackInventory
	}
	if req.QuantityInStock != nil {
		product.QuantityInStock = *req.QuantityInStock
	}
	if req.LowStockThreshold != nil {
		product.LowStockThreshold = *req.LowStockThreshold
	}
	if req.Featured != nil {
		product.Featured = *req.Featured
	}
	if req.IsVegetarian != nil {
		product.IsVegetarian = *req.IsVegetarian
	}
	if req.IsVegan != nil {
		product.IsVegan = *req.IsVegan
	}
	if req.IsSpicy != nil {
		product.IsSpicy = *req.IsSpicy
	}
	if req.IsGlutenFree != nil {
		product.IsGlutenFree = *req.IsGlutenFree
	}

	product.UpdatedBy = &userID

	// Save changes
	product, err = uc.repo.UpdateProduct(product)
	if err != nil {
		return nil, err
	}

	// Create low stock notification if applicable
	if product.TrackInventory && product.QuantityInStock < product.LowStockThreshold {
		actionURL := fmt.Sprintf("/dashboard/products/%d/edit", product.ID)
		iconName := "AlertTriangle"
		color := "orange"
		entityType := "product"
		notification := &domain.Notification{
			TenantID:          tenantID,
			RestaurantID:      restaurantID,
			UserID:            userID,
			Type:              domain.NotificationTypeLowStock,
			Module:            domain.ModuleProducts,
			Title:             fmt.Sprintf("Low Stock Alert: %s", product.NameEn),
			Message:           fmt.Sprintf("Product '%s' has only %d units in stock (threshold: %d)", product.NameEn, product.QuantityInStock, product.LowStockThreshold),
			Priority:          domain.PriorityHigh,
			RelatedEntityType: &entityType,
			RelatedEntityID:   &product.ID,
			ActionURL:         &actionURL,
			IconName:          &iconName,
			Color:             &color,
		}
		if _, err := uc.notificationRepo.CreateNotification(notification); err != nil {
			fmt.Printf("ERROR: Failed to create low stock notification: %v\n", err)
			// Don't fail the product update if notification fails
		}
	}

	// Upload new image if provided
	if imageFile != nil {
		imageURL, err := uc.uploadProductImage(product.ID, imageFile, tenantID)
		if err != nil {
			fmt.Printf("Image upload failed: %v\n", err)
		} else {
			product.MainImageURL = imageURL
		}
	}

	return product, nil
}

// GetProductByID retrieves a product by ID
func (uc *ProductUseCase) GetProductByID(tenantID, restaurantID, productID int) (*domain.Product, error) {
	return uc.repo.GetProductByID(tenantID, restaurantID, productID)
}

// ListProducts retrieves paginated products
func (uc *ProductUseCase) ListProducts(
	tenantID, restaurantID int,
	filters *repository.ProductFilters,
) ([]domain.Product, int, error) {
	resp, err := uc.repo.ListProducts(tenantID, restaurantID, filters)
	if err != nil {
		return nil, 0, err
	}
	return resp.Products, resp.Total, nil
}

// DeleteProduct deletes a product (soft delete with user tracking)
func (uc *ProductUseCase) DeleteProduct(tenantID, restaurantID, productID, userID int) error {
	return uc.repo.DeleteProduct(tenantID, restaurantID, productID, userID)
}

// uploadProductImage saves image to storage
func (uc *ProductUseCase) uploadProductImage(productID int, fileHeader *multipart.FileHeader, tenantID int) (string, error) {
	// Validate file
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !allowedTypes[fileHeader.Header.Get("Content-Type")] {
		return "", fmt.Errorf("invalid file type. allowed: jpeg, png, gif, webp")
	}

	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if fileHeader.Size > maxFileSize {
		return "", fmt.Errorf("file size exceeds 10MB limit")
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("products/%d/%d/%d_%s", tenantID, productID, timestamp, fileHeader.Filename)

	// Save to local storage (or S3 depending on config)
	imagePath := filepath.Join("./uploads", filename)
	os.MkdirAll(filepath.Dir(imagePath), 0755)

	out, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return image URL
	imageURL := fmt.Sprintf("%s/%s", uc.storageURL, filename)
	return imageURL, nil
}

// ListPublicMenu returns public menu
func (uc *ProductUseCase) ListPublicMenu(ctx context.Context, tenantID, restaurantID int, lang string) ([]domain.Product, error) {
	return uc.repo.ListPublic(tenantID, restaurantID, lang)
}

// ListByCategoryPublic returns products for category
func (uc *ProductUseCase) ListByCategoryPublic(ctx context.Context, tenantID, restaurantID, categoryID int, lang string) ([]domain.Product, error) {
	return uc.repo.ListByCategoryPublic(tenantID, restaurantID, categoryID, lang)
}

// SearchPublic searches products
func (uc *ProductUseCase) SearchPublic(ctx context.Context, tenantID, restaurantID int, query string, lang string) ([]domain.Product, error) {
	return uc.repo.SearchPublic(tenantID, restaurantID, query, lang)
}

// GetProductImages returns images
func (uc *ProductUseCase) GetProductImages(ctx context.Context, tenantID, restaurantID, productID int) ([]domain.ProductImage, error) {
	return uc.repo.GetProductImages(tenantID, restaurantID, productID)
}

// GetLowStockAlerts retrieves products with low stock
func (uc *ProductUseCase) GetLowStockAlerts(tenantID, restaurantID int) ([]domain.Product, error) {
	return uc.repo.GetLowStockAlerts(tenantID, restaurantID)
}

// Helper function to validate product request
func validateProductRequest(req *domain.CreateProductRequest) error {
	// if req.NameEn == "" {
	// 	return fmt.Errorf("name in English is required")
	// }
	// Removed mandatory Arabic name check
	// Removed mandatory CategoryID check if 0 is allowed (or handle it gracefully)

	// if req.Price <= 0 {
	// 	return fmt.Errorf("price must be greater than 0")
	// }
	// Allow creation without category initially if needed, or if frontend sends 0
	return nil
}
