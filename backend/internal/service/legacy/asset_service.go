//go:build legacy
// +build legacy

package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AssetService provides image and asset management functionality
type AssetService interface {
	// ProcessImage resizes an image to multiple sizes and returns URLs
	ProcessImage(ctx context.Context, file multipart.File, filename string, assetType string, tenantID, themeID int64) (map[string]string, error)

	// UploadToStorage stores a file and returns its URL
	UploadToStorage(ctx context.Context, filepath string, data []byte) (string, error)

	// DeleteAsset removes an asset file
	DeleteAsset(ctx context.Context, filepath string) error

	// ValidateImage checks if image is valid (MIME type, size)
	ValidateImage(ctx context.Context, file multipart.File, filename string) error
}

// assetServiceImpl implements AssetService
type assetServiceImpl struct {
	storagePath string // Base directory for storing assets
}

// NewAssetService creates a new asset service
func NewAssetService(storagePath string) AssetService {
	return &assetServiceImpl{
		storagePath: storagePath,
	}
}

// ImageSize represents an image size variant
type ImageSize struct {
	Name   string // "full", "large", "medium", "thumb"
	Width  int    // Width in pixels
	Height int    // Height in pixels (0 = auto)
	Quality int   // Compression quality 0-100
}

// Image size configurations
var ImageSizes = []ImageSize{
	{
		Name:    "full",
		Width:   2000,
		Height:  0, // Auto height
		Quality: 85,
	},
	{
		Name:    "large",
		Width:   500,
		Height:  0,
		Quality: 85,
	},
	{
		Name:    "medium",
		Width:   200,
		Height:  0,
		Quality: 85,
	},
	{
		Name:    "thumb",
		Width:   100,
		Height:  0,
		Quality: 85,
	},
}

// ValidateImage checks if image is valid
func (s *assetServiceImpl) ValidateImage(ctx context.Context, file multipart.File, filename string) error {
	log.Printf("[AssetService] Validating image: %s", filename)

	// Check file size (max 10MB) by seeking to end
	fileSize, err := file.Seek(0, 2) // Seek to end
	if err != nil {
		return fmt.Errorf("failed to determine file size: %w", err)
	}

	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if fileSize > maxFileSize {
		return fmt.Errorf("file size exceeds maximum limit of 10MB")
	}

	// Check MIME type by file extension
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
		".gif":  "image/gif",
	}

	if _, valid := validExtensions[ext]; !valid {
		return fmt.Errorf("unsupported file format: %s (allowed: jpg, jpeg, png, webp, gif)", ext)
	}

	// Reset file position to beginning for future reads
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file position: %w", err)
	}

	log.Printf("[AssetService] Image validation passed: %s (size: %d bytes)", filename, fileSize)
	return nil
}

// ProcessImage resizes an image to multiple sizes
func (s *assetServiceImpl) ProcessImage(
	ctx context.Context,
	file multipart.File,
	filename string,
	assetType string,
	tenantID, themeID int64,
) (map[string]string, error) {
	log.Printf("[AssetService] Processing image: %s for tenant: %d, theme: %d", filename, tenantID, themeID)

	// Validate image first
	if err := s.ValidateImage(ctx, file, filename); err != nil {
		log.Printf("[AssetService] Image validation failed: %v", err)
		return nil, err
	}

	// Create storage directory
	storageDir := filepath.Join(
		s.storagePath,
		"themes",
		fmt.Sprintf("%d", tenantID),
		fmt.Sprintf("%d", themeID),
		assetType,
	)

	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Read file data
	fileData := make([]byte, 0)
	buf := make([]byte, 32*1024) // 32KB chunks
	for {
		n, err := file.Read(buf)
		if n > 0 {
			fileData = append(fileData, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	// Process image for each size
	urls := make(map[string]string)
	timestamp := time.Now().Unix()

	for _, size := range ImageSizes {
		// Generate filename for this size
		ext := filepath.Ext(filename)
		baseName := strings.TrimSuffix(filename, ext)
		sizedFilename := fmt.Sprintf("%s_%s_%d%s", baseName, size.Name, timestamp, ext)
		sizedPath := filepath.Join(storageDir, sizedFilename)

		log.Printf("[AssetService] Processing %s image: %s (width: %d)", size.Name, sizedFilename, size.Width)

		// In production, this would use an image library like github.com/disintegration/imaging
		// For now, we'll create a placeholder that copies the file
		// The actual image processing would happen here:
		// 1. Decode image
		// 2. Resize to specified width
		// 3. Compress with specified quality
		// 4. Encode and save

		processedData, err := s.resizeAndCompress(fileData, size.Width, size.Height, size.Quality)
		if err != nil {
			log.Printf("[AssetService] Failed to process %s image: %v", size.Name, err)
			// Continue with other sizes even if one fails
			continue
		}

		// Upload to storage
		url, err := s.UploadToStorage(ctx, sizedPath, processedData)
		if err != nil {
			log.Printf("[AssetService] Failed to upload %s image: %v", size.Name, err)
			continue
		}

		urls[size.Name] = url
		log.Printf("[AssetService] Successfully processed %s image: %s", size.Name, url)
	}

	if len(urls) == 0 {
		return nil, errors.New("failed to process any image sizes")
	}

	log.Printf("[AssetService] Image processing complete: %d sizes processed", len(urls))
	return urls, nil
}

// resizeAndCompress simulates image resizing and compression
// In production, this would use github.com/disintegration/imaging or similar
func (s *assetServiceImpl) resizeAndCompress(
	imageData []byte,
	width int,
	height int,
	quality int,
) ([]byte, error) {
	// Placeholder implementation
	// In production:
	// 1. Decode image (imageData) into an image.Image
	// 2. Resize using imaging.Resize(img, width, height, imaging.Lanczos)
	// 3. Encode with quality setting

	// For now, just return a copy of the data
	// This allows the code to compile and work for testing
	processedData := make([]byte, len(imageData))
	copy(processedData, imageData)

	log.Printf("[AssetService] Image processed (mock): width=%d, height=%d, quality=%d", width, height, quality)
	return processedData, nil
}

// UploadToStorage stores a file in the storage system
func (s *assetServiceImpl) UploadToStorage(
	ctx context.Context,
	filePath string,
	data []byte,
) (string, error) {
	log.Printf("[AssetService] Uploading file to storage: %s", filePath)

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Generate URL (in production, this might be a CDN URL or S3 URL)
	// For local storage, we construct a path-based URL
	url := fmt.Sprintf("/assets/%s", filePath[strings.Index(filePath, "themes"):])

	log.Printf("[AssetService] File uploaded successfully: %s", url)
	return url, nil
}

// DeleteAsset removes an asset file
func (s *assetServiceImpl) DeleteAsset(ctx context.Context, filePath string) error {
	log.Printf("[AssetService] Deleting asset: %s", filePath)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("[AssetService] Asset not found: %s", filePath)
			return nil // Not an error if file doesn't exist
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	log.Printf("[AssetService] Asset deleted successfully: %s", filePath)
	return nil
}

// GetAssetPath constructs the path for an asset
func (s *assetServiceImpl) GetAssetPath(tenantID, themeID int64, assetType, filename string) string {
	return filepath.Join(
		s.storagePath,
		"themes",
		fmt.Sprintf("%d", tenantID),
		fmt.Sprintf("%d", themeID),
		assetType,
		filename,
	)
}
