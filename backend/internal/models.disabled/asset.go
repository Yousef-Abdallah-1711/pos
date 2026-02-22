//go:build legacy
// +build legacy

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Asset represents a media asset (image, CSS, JS, etc.)
type Asset struct {
	ID                int64           `db:"id" json:"id"`
	RestaurantID      int64           `db:"restaurant_id" json:"restaurant_id"`
	ComponentID       *int64          `db:"component_id" json:"component_id,omitempty"`
	ThemeID           *int64          `db:"theme_id" json:"theme_id,omitempty"`
	OriginalFilename  string          `db:"original_filename" json:"original_filename"`
	StorageKey        string          `db:"storage_key" json:"storage_key"`
	FileHash          string          `db:"file_hash" json:"file_hash"`
	FileType          string          `db:"file_type" json:"file_type"`
	MimeType          string          `db:"mime_type" json:"mime_type"`
	OriginalSize      int64           `db:"original_size" json:"original_size"`
	CompressedSize    int64           `db:"compressed_size" json:"compressed_size"`
	Width             *int            `db:"width" json:"width,omitempty"`
	Height            *int            `db:"height" json:"height,omitempty"`
	IsDeuplicated     bool            `db:"is_deduplicated" json:"is_deduplicated"`
	DuplicateOf       *int64          `db:"duplicate_of" json:"duplicate_of,omitempty"`
	Versions          []AssetVersion  `db:"versions" json:"versions,omitempty"`
	CDNUrl            string          `db:"cdn_url" json:"cdn_url"`
	IsPublic          bool            `db:"is_public" json:"is_public"`
	Tags              []string        `db:"tags" json:"tags"`
	Metadata          AssetMetadata   `db:"metadata" json:"metadata"`
	CreatedAt         time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time       `db:"updated_at" json:"updated_at"`
}

// AssetVersion represents a version of an asset
type AssetVersion struct {
	ID            int64     `json:"id"`
	Version       int       `json:"version"`
	StorageKey    string    `json:"storage_key"`
	CompressedSize int64    `json:"compressed_size"`
	CDNUrl        string    `json:"cdn_url"`
	CreatedAt     time.Time `json:"created_at"`
}

// AssetMetadata stores additional asset information
type AssetMetadata struct {
	Format          string            `json:"format"`
	ColorSpace      string            `json:"color_space,omitempty"`
	Orientation     int               `json:"orientation,omitempty"`
	DPI             *int              `json:"dpi,omitempty"`
	BitDepth        *int              `json:"bit_depth,omitempty"`
	HasAlpha        bool              `json:"has_alpha"`
	Quality         *int              `json:"quality,omitempty"`
	ExifData        map[string]string `json:"exif_data,omitempty"`
	CompressionType string            `json:"compression_type"`
	Checksum        string            `json:"checksum"`
}

// Scan implements sql.Scanner for AssetMetadata
func (am *AssetMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), am)
}

// Value implements driver.Valuer for AssetMetadata
func (am AssetMetadata) Value() (driver.Value, error) {
	return json.Marshal(am)
}

// Scan implements sql.Scanner for asset versions
func (versions *[]AssetVersion) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), versions)
}

// Value implements driver.Valuer for asset versions
func (versions []AssetVersion) Value() (driver.Value, error) {
	return json.Marshal(versions)
}

// Scan implements sql.Scanner for tags
func (tags *[]string) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), tags)
}

// Value implements driver.Valuer for tags
func (tags []string) Value() (driver.Value, error) {
	return json.Marshal(tags)
}

// AssetUploadRequest is the request for uploading an asset
type AssetUploadRequest struct {
	OriginalFilename string `json:"original_filename"`
	FileData         []byte `json:"file_data"`
	MimeType         string `json:"mime_type"`
	ComponentID      *int64 `json:"component_id,omitempty"`
	ThemeID          *int64 `json:"theme_id,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	IsPublic         bool   `json:"is_public"`
}

// AssetCompressionOptions defines compression settings
type AssetCompressionOptions struct {
	Quality   int    `json:"quality"`
	Format    string `json:"format"`
	MaxWidth  *int   `json:"max_width,omitempty"`
	MaxHeight *int   `json:"max_height,omitempty"`
	Method    string `json:"method"`
}

// AssetCompressionResult contains compression results
type AssetCompressionResult struct {
	OriginalSize   int64
	CompressedSize int64
	CompressionRatio float64
	Format         string
	Width          int
	Height         int
	Duration       int64
}

// CDNUploadRequest is sent to CDN provider
type CDNUploadRequest struct {
	Key         string
	ContentType string
	Data        []byte
	Metadata    map[string]string
}

// CDNUploadResponse is received from CDN provider
type CDNUploadResponse struct {
	URL       string
	Key       string
	ETag      string
	Size      int64
	Timestamp time.Time
}

// AssetStatistics contains asset usage statistics
type AssetStatistics struct {
	TotalAssets      int
	TotalSize        int64
	TotalCompressed  int64
	DuplicateAssets  int
	CompressionRatio float64
	ByType           map[string]AssetTypeStats
}

// AssetTypeStats contains statistics for a specific asset type
type AssetTypeStats struct {
	Count            int
	TotalSize        int64
	TotalCompressed  int64
	CompressionRatio float64
}
