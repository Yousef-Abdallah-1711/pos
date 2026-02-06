package domain

import (
	"time"
)

// Category represents a product category
type Category struct {
	ID            int       `json:"id"`
	TenantID      int       `json:"tenant_id"`
	RestaurantID  int       `json:"restaurant_id"`
	Name          string    `json:"name"`
	NameAr        string    `json:"name_ar"`
	Description   string    `json:"description,omitempty"`
	DescriptionAr string    `json:"description_ar,omitempty"`
	IconURL       string    `json:"icon_url,omitempty"`
	DisplayOrder  int       `json:"display_order"`
	IsActive      bool      `json:"is_active"`
	Language      string    `json:"language"` // 'en' or 'ar'
	CreatedBy     int       `json:"created_by"`
	UpdatedBy     *int      `json:"updated_by,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Product represents a menu item
type Product struct {
	ID           int    `json:"id"`
	TenantID     int    `json:"tenant_id"`
	RestaurantID int    `json:"restaurant_id"`
	CategoryID   int    `json:"category_id"`
	SKU          string `json:"sku,omitempty"`
	Barcode      string `json:"barcode,omitempty"`

	// Multilingual
	NameEn        string `json:"name_en"`
	NameAr        string `json:"name_ar"`
	DescriptionEn string `json:"description_en,omitempty"`
	DescriptionAr string `json:"description_ar,omitempty"`

	// Pricing
	Price              float64  `json:"price"`
	Cost               *float64 `json:"cost,omitempty"`
	DiscountPrice      *float64 `json:"discount_price,omitempty"`
	DiscountPercentage *float64 `json:"discount_percentage,omitempty"`

	// Nutritional Info
	Calories  *int     `json:"calories,omitempty"`
	ProteinG  *float64 `json:"protein_g,omitempty"`
	CarbsG    *float64 `json:"carbs_g,omitempty"`
	FatG      *float64 `json:"fat_g,omitempty"`
	FiberG    *float64 `json:"fiber_g,omitempty"`
	Allergens []string `json:"allergens,omitempty"`

	// Classification
	IsVegetarian bool `json:"is_vegetarian"`
	IsVegan      bool `json:"is_vegan"`
	IsSpicy      bool `json:"is_spicy"`
	IsGlutenFree bool `json:"is_gluten_free"`

	// Availability
	IsAvailable    bool     `json:"is_available"`
	AvailableFrom  *string  `json:"available_from,omitempty"`  // HH:MM format
	AvailableUntil *string  `json:"available_until,omitempty"` // HH:MM format
	AvailableDays  []string `json:"available_days,omitempty"`  // ["Monday", "Tuesday"]

	// Inventory
	TrackInventory    bool `json:"track_inventory"`
	QuantityInStock   int  `json:"quantity_in_stock"`
	LowStockThreshold int  `json:"low_stock_threshold"`
	ReorderQuantity   int  `json:"reorder_quantity"`

	// Display
	DisplayOrder int    `json:"display_order"`
	Featured     bool   `json:"featured"`
	MainImageURL string `json:"main_image_url,omitempty"`

	// Status
	Status    string    `json:"status"` // 'active', 'inactive', 'discontinued'
	CreatedBy int       `json:"created_by"`
	UpdatedBy *int      `json:"updated_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (populated on demand)
	Images   []ProductImage   `json:"images,omitempty"`
	Variants []ProductVariant `json:"variants,omitempty"`
	AddOns   []ProductAddOn   `json:"addons,omitempty"`
}

// ProductImage represents product images
type ProductImage struct {
	ID           int       `json:"id"`
	ProductID    int       `json:"product_id"`
	ImageURL     string    `json:"image_url"`
	AltTextEn    string    `json:"alt_text_en,omitempty"`
	AltTextAr    string    `json:"alt_text_ar,omitempty"`
	IsPrimary    bool      `json:"is_primary"`
	DisplayOrder int       `json:"display_order"`
	FileSize     int       `json:"file_size,omitempty"`
	Width        int       `json:"width,omitempty"`
	Height       int       `json:"height,omitempty"`
	MimeType     string    `json:"mime_type,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ProductVariant represents product variations (size, color, etc.)
type ProductVariant struct {
	ID              int       `json:"id"`
	ProductID       int       `json:"product_id"`
	NameEn          string    `json:"name_en"`
	NameAr          string    `json:"name_ar"`
	SKUSuffix       string    `json:"sku_suffix,omitempty"`
	PriceAdjustment float64   `json:"price_adjustment"`
	QuantityInStock int       `json:"quantity_in_stock"`
	IsAvailable     bool      `json:"is_available"`
	DisplayOrder    int       `json:"display_order"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ProductAddOn represents add-ons (toppings, sauces, etc.)
type ProductAddOn struct {
	ID                  int       `json:"id"`
	RestaurantID        int       `json:"restaurant_id"`
	NameEn              string    `json:"name_en"`
	NameAr              string    `json:"name_ar"`
	DescriptionEn       string    `json:"description_en,omitempty"`
	DescriptionAr       string    `json:"description_ar,omitempty"`
	Price               float64   `json:"price"`
	IsAvailable         bool      `json:"is_available"`
	MaxQuantityPerOrder int       `json:"max_quantity_per_order"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// CreateProductRequest is the request DTO for creating products
type CreateProductRequest struct {
	CategoryID         int      `json:"category_id"`
	SKU                string   `json:"sku"`
	Barcode            string   `json:"barcode"`
	NameEn             string   `json:"name_en"`
	NameAr             string   `json:"name_ar"`
	DescriptionEn      string   `json:"description_en"`
	DescriptionAr      string   `json:"description_ar"`
	Price              float64  `json:"price"`
	Cost               *float64 `json:"cost"`
	DiscountPrice      *float64 `json:"discount_price"`
	DiscountPercentage *float64 `json:"discount_percentage"`
	Calories           *int     `json:"calories"`
	ProteinG           *float64 `json:"protein_g"`
	CarbsG             *float64 `json:"carbs_g"`
	FatG               *float64 `json:"fat_g"`
	FiberG             *float64 `json:"fiber_g"`
	Allergens          []string `json:"allergens"`
	IsVegetarian       bool     `json:"is_vegetarian"`
	IsVegan            bool     `json:"is_vegan"`
	IsSpicy            bool     `json:"is_spicy"`
	IsGlutenFree       bool     `json:"is_gluten_free"`
	IsAvailable        bool     `json:"is_available"`
	AvailableFrom      *string  `json:"available_from"`
	AvailableUntil     *string  `json:"available_until"`
	AvailableDays      []string `json:"available_days"`
	TrackInventory     bool     `json:"track_inventory"`
	QuantityInStock    int      `json:"quantity_in_stock"`
	LowStockThreshold  int      `json:"low_stock_threshold"`
	ReorderQuantity    int      `json:"reorder_quantity"`
	Featured           bool     `json:"featured"`
	DisplayOrder       int      `json:"display_order"`
}

// UpdateProductRequest is the request DTO for updating products
type UpdateProductRequest struct {
	CategoryID        *int     `json:"category_id"`
	NameEn            *string  `json:"name_en"`
	NameAr            *string  `json:"name_ar"`
	DescriptionEn     *string  `json:"description_en"`
	DescriptionAr     *string  `json:"description_ar"`
	Price             *float64 `json:"price"`
	Cost              *float64 `json:"cost"`
	DiscountPrice     *float64 `json:"discount_price"`
	IsAvailable       *bool    `json:"is_available"`
	TrackInventory    *bool    `json:"track_inventory"`
	QuantityInStock   *int     `json:"quantity_in_stock"`
	LowStockThreshold *int     `json:"low_stock_threshold"`
	Featured          *bool    `json:"featured"`
	IsVegetarian      *bool    `json:"is_vegetarian"`
	IsVegan           *bool    `json:"is_vegan"`
	IsSpicy           *bool    `json:"is_spicy"`
	IsGlutenFree      *bool    `json:"is_gluten_free"`
}

// ProductListResponse is paginated list response
type ProductListResponse struct {
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Products []Product `json:"products"`
}
