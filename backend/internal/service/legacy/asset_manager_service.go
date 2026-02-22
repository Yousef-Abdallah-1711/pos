//go:build legacy
// +build legacy

package service

import (
	"crypto/sha256"
	"fmt"
	"io"
	"pos-saas/internal/repository"
	"time"
)

// AssetManagerService handles asset management operations
type AssetManagerService struct {
	assetRepo repository.AssetRepository
	cdnService *CDNService
}

// NewAssetManagerService creates a new asset manager service
func NewAssetManagerService(
	assetRepo repository.AssetRepository,
	cdnService *CDNService,
) *AssetManagerService {
	return &AssetManagerService{
		assetRepo: assetRepo,
		cdnService: cdnService,
	}
}

// UploadAsset uploads a new asset with deduplication and compression
func (s *AssetManagerService) UploadAsset(
	restaurantID int64,
	req *models.AssetUploadRequest,
) (*models.Asset, error) {
	fmt.Printf("DEBUG: Uploading asset %s for restaurant %d\n", req.OriginalFilename, restaurantID)

	// Calculate file hash for deduplication
	fileHash := s.calculateFileHash(req.FileData)

	// Check for existing asset with same hash
	existing, err := s.assetRepo.GetByFileHash(restaurantID, fileHash)
	if err == nil && existing != nil {
		fmt.Printf("DEBUG: Found duplicate asset, reusing %d\n", existing.ID)
		return existing, nil
	}

	// Determine file type
	fileType := s.determineFileType(req.MimeType)

	// Create storage key
	storageKey := s.generateStorageKey(restaurantID, req.OriginalFilename)

	// Compress asset based on type
	compressedData := req.FileData
	compressedSize := int64(len(req.FileData))
	metadata := models.AssetMetadata{
		Checksum: fileHash,
		CompressionType: "none",
	}

	if fileType == "image" {
		result, err := s.compressImage(req.FileData, req.MimeType)
		if err == nil && result != nil {
			compressedData = result.CompressedData
			compressedSize = result.CompressedSize
			metadata.Format = result.Format
			metadata.Quality = &result.Quality
			metadata.CompressionType = "auto"
		}
	}

	// Upload to CDN
	cdnUrl, err := s.cdnService.UploadAsset(storageKey, req.MimeType, compressedData)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to CDN: %w", err)
	}

	// Get image dimensions if applicable
	var width, height *int
	if fileType == "image" {
		w, h := s.getImageDimensions(compressedData, req.MimeType)
		if w > 0 && h > 0 {
			width = &w
			height = &h
		}
	}

	// Create asset record
	asset := &models.Asset{
		RestaurantID:      restaurantID,
		ComponentID:       req.ComponentID,
		ThemeID:           req.ThemeID,
		OriginalFilename:  req.OriginalFilename,
		StorageKey:        storageKey,
		FileHash:          fileHash,
		FileType:          fileType,
		MimeType:          req.MimeType,
		OriginalSize:      int64(len(req.FileData)),
		CompressedSize:    compressedSize,
		Width:             width,
		Height:            height,
		CDNUrl:            cdnUrl,
		IsPublic:          req.IsPublic,
		Tags:              req.Tags,
		Metadata:          metadata,
		IsDeuplicated:     false,
	}

	created, err := s.assetRepo.Create(asset)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset record: %w", err)
	}

	fmt.Printf("DEBUG: Created asset %d with hash %s\n", created.ID, fileHash)
	return created, nil
}

// GetAsset retrieves an asset by ID
func (s *AssetManagerService) GetAsset(assetID int64) (*models.Asset, error) {
	asset, err := s.assetRepo.GetByID(assetID)
	if err != nil {
		return nil, fmt.Errorf("asset not found: %w", err)
	}

	// Update last accessed time
	s.assetRepo.RecordUsage(assetID, "view", assetID, "asset_view")

	return asset, nil
}

// ListAssetsByRestaurant lists all assets for a restaurant
func (s *AssetManagerService) ListAssetsByRestaurant(
	restaurantID int64,
	limit int,
	offset int,
) ([]*models.Asset, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.assetRepo.ListByRestaurant(restaurantID, limit, offset)
}

// ListAssetsByComponent lists all assets for a component
func (s *AssetManagerService) ListAssetsByComponent(componentID int64) ([]*models.Asset, error) {
	return s.assetRepo.ListByComponent(componentID)
}

// ListAssetsByTheme lists all assets for a theme
func (s *AssetManagerService) ListAssetsByTheme(themeID int64) ([]*models.Asset, error) {
	return s.assetRepo.ListByTheme(themeID)
}

// DeleteAsset deletes an asset
func (s *AssetManagerService) DeleteAsset(assetID int64) error {
	fmt.Printf("DEBUG: Deleting asset %d\n", assetID)

	asset, err := s.assetRepo.GetByID(assetID)
	if err != nil {
		return fmt.Errorf("asset not found: %w", err)
	}

	// Check usage count
	usageCount, _ := s.assetRepo.GetUsageCount(assetID)
	if usageCount > 0 {
		return fmt.Errorf("cannot delete asset in use (%d usages)", usageCount)
	}

	// Delete from CDN
	_ = s.cdnService.DeleteAsset(asset.StorageKey)

	// Delete from database
	return s.assetRepo.Delete(assetID)
}

// GetAssetStatistics gets statistics for restaurant assets
func (s *AssetManagerService) GetAssetStatistics(restaurantID int64) (*models.AssetStatistics, error) {
	return s.assetRepo.GetStatistics(restaurantID)
}

// DeduplicateAssets finds and marks duplicate assets
func (s *AssetManagerService) DeduplicateAssets(restaurantID int64) (int, error) {
	fmt.Printf("DEBUG: Running deduplication for restaurant %d\n", restaurantID)

	assets, err := s.assetRepo.ListByRestaurant(restaurantID, 10000, 0)
	if err != nil {
		return 0, err
	}

	hashMap := make(map[string]*models.Asset)
	duplicateCount := 0

	for _, asset := range assets {
		if asset.IsDeuplicated {
			continue
		}

		if original, exists := hashMap[asset.FileHash]; exists {
			// Mark as duplicate
			_ = s.assetRepo.MarkAsDuplicate(asset.ID, original.ID)
			duplicateCount++
			fmt.Printf("DEBUG: Marked asset %d as duplicate of %d\n", asset.ID, original.ID)
		} else {
			hashMap[asset.FileHash] = asset
		}
	}

	return duplicateCount, nil
}

// CleanupUnusedAssets deletes assets not used for specified days
func (s *AssetManagerService) CleanupUnusedAssets(restaurantID int64, daysOld int) (int, error) {
	fmt.Printf("DEBUG: Cleaning up unused assets older than %d days\n", daysOld)

	deleted, err := s.assetRepo.DeleteUnused(restaurantID, daysOld)
	if err != nil {
		return 0, fmt.Errorf("failed to delete unused assets: %w", err)
	}

	fmt.Printf("DEBUG: Deleted %d unused assets\n", deleted)
	return deleted, nil
}

// RecordAssetUsage records that an asset is being used
func (s *AssetManagerService) RecordAssetUsage(
	assetID int64,
	usedInType string,
	usedInID int64,
	context string,
) error {
	return s.assetRepo.RecordUsage(assetID, usedInType, usedInID, context)
}

// Helper functions

// calculateFileHash calculates SHA-256 hash of file
func (s *AssetManagerService) calculateFileHash(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// determineFileType determines the type of file
func (s *AssetManagerService) determineFileType(mimeType string) string {
	switch {
	case mimeType[:6] == "image/":
		return "image"
	case mimeType[:5] == "video/":
		return "video"
	case mimeType[:5] == "audio/":
		return "audio"
	case mimeType == "application/pdf":
		return "document"
	default:
		return "file"
	}
}

// generateStorageKey generates a unique storage key for the asset
func (s *AssetManagerService) generateStorageKey(restaurantID int64, originalFilename string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("assets/%d/%d-%s", restaurantID, timestamp, originalFilename)
}

// compressImage compresses image data
func (s *AssetManagerService) compressImage(data []byte, mimeType string) (*ImageCompressionResult, error) {
	// This is a placeholder - actual implementation would use imaging library
	// For now, return original data
	return &ImageCompressionResult{
		CompressedData: data,
		CompressedSize: int64(len(data)),
		Format: mimeType,
		Quality: 85,
	}, nil
}

// getImageDimensions gets image width and height
func (s *AssetManagerService) getImageDimensions(data []byte, mimeType string) (int, int) {
	// Placeholder implementation
	// Real implementation would decode image and get dimensions
	return 1200, 800
}

// ImageCompressionResult holds compression results
type ImageCompressionResult struct {
	CompressedData []byte
	CompressedSize int64
	Format         string
	Quality        int
}
