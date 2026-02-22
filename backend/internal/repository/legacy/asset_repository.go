//go:build legacy
// +build legacy

package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// AssetRepository handles database operations for assets
type AssetRepository interface {
	// Create creates a new asset
	Create(asset *models.Asset) (*models.Asset, error)

	// GetByID gets an asset by ID
	GetByID(id int64) (*models.Asset, error)

	// GetByFileHash gets assets by file hash
	GetByFileHash(restaurantID int64, fileHash string) (*models.Asset, error)

	// ListByRestaurant lists all assets for a restaurant
	ListByRestaurant(restaurantID int64, limit int, offset int) ([]*models.Asset, error)

	// ListByComponent lists all assets for a component
	ListByComponent(componentID int64) ([]*models.Asset, error)

	// ListByTheme lists all assets for a theme
	ListByTheme(themeID int64) ([]*models.Asset, error)

	// Update updates an asset
	Update(asset *models.Asset) (*models.Asset, error)

	// Delete deletes an asset
	Delete(id int64) error

	// MarkAsDuplicate marks an asset as a duplicate
	MarkAsDuplicate(assetID int64, duplicateOfID int64) error

	// GetDuplicates gets all duplicates of an asset
	GetDuplicates(assetID int64) ([]*models.Asset, error)

	// GetStatistics gets asset statistics for a restaurant
	GetStatistics(restaurantID int64) (*models.AssetStatistics, error)

	// RecordUsage records that an asset is being used
	RecordUsage(assetID int64, usedInType string, usedInID int64, context string) error

	// GetUsageCount gets the count of usages for an asset
	GetUsageCount(assetID int64) (int, error)

	// DeleteUnused deletes assets that have no usages
	DeleteUnused(restaurantID int64, daysOld int) (int, error)
}

// MySQLAssetRepository implements AssetRepository for MySQL
type MySQLAssetRepository struct {
	db *sql.DB
}

// NewAssetRepository creates a new asset repository
func NewAssetRepository(db *sql.DB) AssetRepository {
	return &MySQLAssetRepository{db: db}
}

// Create implements AssetRepository.Create
func (r *MySQLAssetRepository) Create(asset *models.Asset) (*models.Asset, error) {
	query := `
		INSERT INTO assets (
			restaurant_id, component_id, theme_id,
			original_filename, storage_key, file_hash,
			file_type, mime_type, original_size, compressed_size,
			width, height, is_deduplicated, duplicate_of,
			cdn_url, is_public, tags, metadata,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	tagsJSON, _ := json.Marshal(asset.Tags)
	metadataJSON, _ := json.Marshal(asset.Metadata)

	err := r.db.QueryRow(
		query,
		asset.RestaurantID, asset.ComponentID, asset.ThemeID,
		asset.OriginalFilename, asset.StorageKey, asset.FileHash,
		asset.FileType, asset.MimeType, asset.OriginalSize, asset.CompressedSize,
		asset.Width, asset.Height, asset.IsDeuplicated, asset.DuplicateOf,
		asset.CDNUrl, asset.IsPublic, tagsJSON, metadataJSON,
		time.Now(), time.Now(),
	).Scan(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt)

	if err != nil {
		fmt.Printf("DEBUG: Failed to create asset: %v\n", err)
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	fmt.Printf("DEBUG: Created asset %d: %s\n", asset.ID, asset.OriginalFilename)
	return asset, nil
}

// GetByID implements AssetRepository.GetByID
func (r *MySQLAssetRepository) GetByID(id int64) (*models.Asset, error) {
	query := `
		SELECT id, restaurant_id, component_id, theme_id,
		       original_filename, storage_key, file_hash,
		       file_type, mime_type, original_size, compressed_size,
		       width, height, is_deduplicated, duplicate_of,
		       cdn_url, is_public, tags, metadata,
		       created_at, updated_at
		FROM assets
		WHERE id = ?
	`

	asset := &models.Asset{}
	var tagsJSON, metadataJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&asset.ID, &asset.RestaurantID, &asset.ComponentID, &asset.ThemeID,
		&asset.OriginalFilename, &asset.StorageKey, &asset.FileHash,
		&asset.FileType, &asset.MimeType, &asset.OriginalSize, &asset.CompressedSize,
		&asset.Width, &asset.Height, &asset.IsDeuplicated, &asset.DuplicateOf,
		&asset.CDNUrl, &asset.IsPublic, &tagsJSON, &metadataJSON,
		&asset.CreatedAt, &asset.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Parse JSON fields
	if tagsJSON != nil {
		json.Unmarshal(tagsJSON, &asset.Tags)
	}
	if metadataJSON != nil {
		json.Unmarshal(metadataJSON, &asset.Metadata)
	}

	return asset, nil
}

// GetByFileHash implements AssetRepository.GetByFileHash
func (r *MySQLAssetRepository) GetByFileHash(restaurantID int64, fileHash string) (*models.Asset, error) {
	query := `
		SELECT id, restaurant_id, component_id, theme_id,
		       original_filename, storage_key, file_hash,
		       file_type, mime_type, original_size, compressed_size,
		       width, height, is_deduplicated, duplicate_of,
		       cdn_url, is_public, tags, metadata,
		       created_at, updated_at
		FROM assets
		WHERE restaurant_id = ? AND file_hash = ? AND is_deduplicated = false
		LIMIT 1
	`

	asset := &models.Asset{}
	var tagsJSON, metadataJSON []byte

	err := r.db.QueryRow(query, restaurantID, fileHash).Scan(
		&asset.ID, &asset.RestaurantID, &asset.ComponentID, &asset.ThemeID,
		&asset.OriginalFilename, &asset.StorageKey, &asset.FileHash,
		&asset.FileType, &asset.MimeType, &asset.OriginalSize, &asset.CompressedSize,
		&asset.Width, &asset.Height, &asset.IsDeuplicated, &asset.DuplicateOf,
		&asset.CDNUrl, &asset.IsPublic, &tagsJSON, &metadataJSON,
		&asset.CreatedAt, &asset.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset not found: hash=%s", fileHash)
		}
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Parse JSON fields
	if tagsJSON != nil {
		json.Unmarshal(tagsJSON, &asset.Tags)
	}
	if metadataJSON != nil {
		json.Unmarshal(metadataJSON, &asset.Metadata)
	}

	return asset, nil
}

// ListByRestaurant implements AssetRepository.ListByRestaurant
func (r *MySQLAssetRepository) ListByRestaurant(restaurantID int64, limit int, offset int) ([]*models.Asset, error) {
	query := `
		SELECT id, restaurant_id, component_id, theme_id,
		       original_filename, storage_key, file_hash,
		       file_type, mime_type, original_size, compressed_size,
		       width, height, is_deduplicated, duplicate_of,
		       cdn_url, is_public, tags, metadata,
		       created_at, updated_at
		FROM assets
		WHERE restaurant_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, restaurantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	defer rows.Close()

	var assets []*models.Asset

	for rows.Next() {
		asset := &models.Asset{}
		var tagsJSON, metadataJSON []byte

		err := rows.Scan(
			&asset.ID, &asset.RestaurantID, &asset.ComponentID, &asset.ThemeID,
			&asset.OriginalFilename, &asset.StorageKey, &asset.FileHash,
			&asset.FileType, &asset.MimeType, &asset.OriginalSize, &asset.CompressedSize,
			&asset.Width, &asset.Height, &asset.IsDeuplicated, &asset.DuplicateOf,
			&asset.CDNUrl, &asset.IsPublic, &tagsJSON, &metadataJSON,
			&asset.CreatedAt, &asset.UpdatedAt,
		)

		if err != nil {
			fmt.Printf("DEBUG: Error scanning asset row: %v\n", err)
			continue
		}

		// Parse JSON fields
		if tagsJSON != nil {
			json.Unmarshal(tagsJSON, &asset.Tags)
		}
		if metadataJSON != nil {
			json.Unmarshal(metadataJSON, &asset.Metadata)
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

// ListByComponent implements AssetRepository.ListByComponent
func (r *MySQLAssetRepository) ListByComponent(componentID int64) ([]*models.Asset, error) {
	query := `
		SELECT id, restaurant_id, component_id, theme_id,
		       original_filename, storage_key, file_hash,
		       file_type, mime_type, original_size, compressed_size,
		       width, height, is_deduplicated, duplicate_of,
		       cdn_url, is_public, tags, metadata,
		       created_at, updated_at
		FROM assets
		WHERE component_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, componentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	defer rows.Close()

	var assets []*models.Asset

	for rows.Next() {
		asset := &models.Asset{}
		var tagsJSON, metadataJSON []byte

		err := rows.Scan(
			&asset.ID, &asset.RestaurantID, &asset.ComponentID, &asset.ThemeID,
			&asset.OriginalFilename, &asset.StorageKey, &asset.FileHash,
			&asset.FileType, &asset.MimeType, &asset.OriginalSize, &asset.CompressedSize,
			&asset.Width, &asset.Height, &asset.IsDeuplicated, &asset.DuplicateOf,
			&asset.CDNUrl, &asset.IsPublic, &tagsJSON, &metadataJSON,
			&asset.CreatedAt, &asset.UpdatedAt,
		)

		if err != nil {
			continue
		}

		// Parse JSON fields
		if tagsJSON != nil {
			json.Unmarshal(tagsJSON, &asset.Tags)
		}
		if metadataJSON != nil {
			json.Unmarshal(metadataJSON, &asset.Metadata)
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

// ListByTheme implements AssetRepository.ListByTheme
func (r *MySQLAssetRepository) ListByTheme(themeID int64) ([]*models.Asset, error) {
	query := `
		SELECT id, restaurant_id, component_id, theme_id,
		       original_filename, storage_key, file_hash,
		       file_type, mime_type, original_size, compressed_size,
		       width, height, is_deduplicated, duplicate_of,
		       cdn_url, is_public, tags, metadata,
		       created_at, updated_at
		FROM assets
		WHERE theme_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	defer rows.Close()

	var assets []*models.Asset

	for rows.Next() {
		asset := &models.Asset{}
		var tagsJSON, metadataJSON []byte

		err := rows.Scan(
			&asset.ID, &asset.RestaurantID, &asset.ComponentID, &asset.ThemeID,
			&asset.OriginalFilename, &asset.StorageKey, &asset.FileHash,
			&asset.FileType, &asset.MimeType, &asset.OriginalSize, &asset.CompressedSize,
			&asset.Width, &asset.Height, &asset.IsDeuplicated, &asset.DuplicateOf,
			&asset.CDNUrl, &asset.IsPublic, &tagsJSON, &metadataJSON,
			&asset.CreatedAt, &asset.UpdatedAt,
		)

		if err != nil {
			continue
		}

		// Parse JSON fields
		if tagsJSON != nil {
			json.Unmarshal(tagsJSON, &asset.Tags)
		}
		if metadataJSON != nil {
			json.Unmarshal(metadataJSON, &asset.Metadata)
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

// Update implements AssetRepository.Update
func (r *MySQLAssetRepository) Update(asset *models.Asset) (*models.Asset, error) {
	query := `
		UPDATE assets
		SET original_filename = ?, storage_key = ?, file_hash = ?,
		    file_type = ?, mime_type = ?, original_size = ?, compressed_size = ?,
		    width = ?, height = ?, is_deduplicated = ?, duplicate_of = ?,
		    cdn_url = ?, is_public = ?, tags = ?, metadata = ?,
		    component_id = ?, theme_id = ?, updated_at = ?
		WHERE id = ?
	`

	tagsJSON, _ := json.Marshal(asset.Tags)
	metadataJSON, _ := json.Marshal(asset.Metadata)

	_, err := r.db.Exec(
		query,
		asset.OriginalFilename, asset.StorageKey, asset.FileHash,
		asset.FileType, asset.MimeType, asset.OriginalSize, asset.CompressedSize,
		asset.Width, asset.Height, asset.IsDeuplicated, asset.DuplicateOf,
		asset.CDNUrl, asset.IsPublic, tagsJSON, metadataJSON,
		asset.ComponentID, asset.ThemeID, time.Now(),
		asset.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update asset: %w", err)
	}

	asset.UpdatedAt = time.Now()
	fmt.Printf("DEBUG: Updated asset %d\n", asset.ID)
	return asset, nil
}

// Delete implements AssetRepository.Delete
func (r *MySQLAssetRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM assets WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset not found: %d", id)
	}

	fmt.Printf("DEBUG: Deleted asset %d\n", id)
	return nil
}

// MarkAsDuplicate implements AssetRepository.MarkAsDuplicate
func (r *MySQLAssetRepository) MarkAsDuplicate(assetID int64, duplicateOfID int64) error {
	_, err := r.db.Exec(
		"UPDATE assets SET is_deduplicated = true, duplicate_of = ? WHERE id = ?",
		duplicateOfID, assetID,
	)

	if err != nil {
		return fmt.Errorf("failed to mark asset as duplicate: %w", err)
	}

	fmt.Printf("DEBUG: Marked asset %d as duplicate of %d\n", assetID, duplicateOfID)
	return nil
}

// GetDuplicates implements AssetRepository.GetDuplicates
func (r *MySQLAssetRepository) GetDuplicates(assetID int64) ([]*models.Asset, error) {
	query := `
		SELECT id, restaurant_id, component_id, theme_id,
		       original_filename, storage_key, file_hash,
		       file_type, mime_type, original_size, compressed_size,
		       width, height, is_deduplicated, duplicate_of,
		       cdn_url, is_public, tags, metadata,
		       created_at, updated_at
		FROM assets
		WHERE duplicate_of = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get duplicates: %w", err)
	}
	defer rows.Close()

	var assets []*models.Asset

	for rows.Next() {
		asset := &models.Asset{}
		var tagsJSON, metadataJSON []byte

		err := rows.Scan(
			&asset.ID, &asset.RestaurantID, &asset.ComponentID, &asset.ThemeID,
			&asset.OriginalFilename, &asset.StorageKey, &asset.FileHash,
			&asset.FileType, &asset.MimeType, &asset.OriginalSize, &asset.CompressedSize,
			&asset.Width, &asset.Height, &asset.IsDeuplicated, &asset.DuplicateOf,
			&asset.CDNUrl, &asset.IsPublic, &tagsJSON, &metadataJSON,
			&asset.CreatedAt, &asset.UpdatedAt,
		)

		if err != nil {
			continue
		}

		// Parse JSON fields
		if tagsJSON != nil {
			json.Unmarshal(tagsJSON, &asset.Tags)
		}
		if metadataJSON != nil {
			json.Unmarshal(metadataJSON, &asset.Metadata)
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

// GetStatistics implements AssetRepository.GetStatistics
func (r *MySQLAssetRepository) GetStatistics(restaurantID int64) (*models.AssetStatistics, error) {
	query := `
		SELECT
			COUNT(*) as total_assets,
			SUM(original_size) as total_size,
			SUM(compressed_size) as total_compressed,
			SUM(CASE WHEN is_deduplicated = true THEN 1 ELSE 0 END) as duplicate_assets
		FROM assets
		WHERE restaurant_id = ?
	`

	stats := &models.AssetStatistics{ByType: make(map[string]models.AssetTypeStats)}

	var totalAssets, duplicateAssets int
	var totalSize, totalCompressed sql.NullInt64

	err := r.db.QueryRow(query, restaurantID).Scan(
		&totalAssets, &totalSize, &totalCompressed, &duplicateAssets,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	stats.TotalAssets = totalAssets
	stats.TotalSize = totalSize.Int64
	stats.TotalCompressed = totalCompressed.Int64
	stats.DuplicateAssets = duplicateAssets

	if stats.TotalSize > 0 {
		stats.CompressionRatio = float64(stats.TotalCompressed) / float64(stats.TotalSize)
	}

	// Get stats by type
	typeQuery := `
		SELECT mime_type,
		       COUNT(*) as count,
		       SUM(original_size) as total_size,
		       SUM(compressed_size) as total_compressed
		FROM assets
		WHERE restaurant_id = ?
		GROUP BY mime_type
	`

	rows, err := r.db.Query(typeQuery, restaurantID)
	if err != nil {
		return stats, nil // Return partial stats
	}
	defer rows.Close()

	for rows.Next() {
		var mimeType string
		var count int
		var totalSize, totalCompressed sql.NullInt64

		if err := rows.Scan(&mimeType, &count, &totalSize, &totalCompressed); err != nil {
			continue
		}

		typeStats := models.AssetTypeStats{
			Count:          count,
			TotalSize:      totalSize.Int64,
			TotalCompressed: totalCompressed.Int64,
		}

		if typeStats.TotalSize > 0 {
			typeStats.CompressionRatio = float64(typeStats.TotalCompressed) / float64(typeStats.TotalSize)
		}

		stats.ByType[mimeType] = typeStats
	}

	return stats, nil
}

// RecordUsage implements AssetRepository.RecordUsage
func (r *MySQLAssetRepository) RecordUsage(assetID int64, usedInType string, usedInID int64, context string) error {
	query := `
		INSERT INTO asset_usage (asset_id, used_in_type, used_in_id, usage_context, last_used_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE last_used_at = ?
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		assetID, usedInType, usedInID, context, now, now, now,
	)

	if err != nil {
		return fmt.Errorf("failed to record asset usage: %w", err)
	}

	return nil
}

// GetUsageCount implements AssetRepository.GetUsageCount
func (r *MySQLAssetRepository) GetUsageCount(assetID int64) (int, error) {
	var count int

	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM asset_usage WHERE asset_id = ?",
		assetID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get usage count: %w", err)
	}

	return count, nil
}

// DeleteUnused implements AssetRepository.DeleteUnused
func (r *MySQLAssetRepository) DeleteUnused(restaurantID int64, daysOld int) (int, error) {
	query := `
		DELETE FROM assets
		WHERE restaurant_id = ?
		AND id NOT IN (SELECT asset_id FROM asset_usage WHERE last_used_at > DATE_SUB(NOW(), INTERVAL ? DAY))
		AND created_at < DATE_SUB(NOW(), INTERVAL ? DAY)
	`

	result, err := r.db.Exec(query, restaurantID, daysOld, daysOld)
	if err != nil {
		return 0, fmt.Errorf("failed to delete unused assets: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	fmt.Printf("DEBUG: Deleted %d unused assets\n", rowsAffected)
	return int(rowsAffected), nil
}
