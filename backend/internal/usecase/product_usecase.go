package usecase

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
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
	ctx context.Context,
	tenantID, restaurantID, userID int,
	req *domain.CreateProductRequest,
	imageFile *multipart.FileHeader,
) (*domain.Product, error) {
	// Validate request
	if err := validateProductRequest(req); err != nil {
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

	if imageFile != nil {
		imageURL, uploadErr := uc.uploadProductImage(restaurantID, imageFile, tenantID)
		if uploadErr != nil {
			return nil, uploadErr
		}
		product.MainImageURL = imageURL
	}
	// Save to database
	preparedImageURL := product.MainImageURL
	product, err := uc.repo.CreateProduct(product)
	if err != nil {
		if imageFile != nil {
			uc.discardPreparedImage(preparedImageURL)
		}
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
		if _, err := uc.notificationRepo.CreateNotification(ctx, notification); err != nil {
			log.Printf("[PRODUCT] failed to create low-stock notification: tenant=%d restaurant=%d product=%d err=%v", tenantID, restaurantID, product.ID, err)
			// Don't fail the product creation if notification fails
		}
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (uc *ProductUseCase) UpdateProduct(
	ctx context.Context,
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
		product.StockAdjustmentRequested = *req.QuantityInStock != product.QuantityInStock
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

	if imageFile != nil {
		imageURL, uploadErr := uc.uploadProductImage(restaurantID, imageFile, tenantID)
		if uploadErr != nil {
			return nil, uploadErr
		}
		product.MainImageURL = imageURL
	}
	// Save changes
	preparedImageURL := product.MainImageURL
	product, err = uc.repo.UpdateProductContext(ctx, product)
	if err != nil {
		if imageFile != nil {
			uc.discardPreparedImage(preparedImageURL)
		}
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
		if _, err := uc.notificationRepo.CreateNotification(ctx, notification); err != nil {
			log.Printf("[PRODUCT] failed to create low-stock notification: tenant=%d restaurant=%d product=%d err=%v", tenantID, restaurantID, product.ID, err)
			// Don't fail the product update if notification fails
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
