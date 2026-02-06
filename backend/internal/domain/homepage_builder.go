package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// LegacyTheme represents a restaurant's website theme configuration (legacy, use Theme from theme.go instead)
type LegacyTheme struct {
	ID              int64     `json:"id"`
	RestaurantID    int64     `json:"restaurant_id"`
	TenantID        int64     `json:"tenant_id"`
	TemplateName    string    `json:"template_name"` // modern, classic, minimalist
	PrimaryColor    string    `json:"primary_color"`
	SecondaryColor  string    `json:"secondary_color"`
	AccentColor     string    `json:"accent_color"`
	BackgroundColor string    `json:"background_color"`
	TextColor       string    `json:"text_color"`
	FontFamily      string    `json:"font_family"`
	FontSizeBase    int       `json:"font_size_base"`
	BorderRadius    int       `json:"border_radius"`
	LogoURL         string    `json:"logo_url"`
	FaviconURL      string    `json:"favicon_url"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// HomepageSection represents a customizable homepage section
type HomepageSection struct {
	ID                 int64           `json:"id"`
	RestaurantID       int64           `json:"restaurant_id"`
	TenantID           int64           `json:"tenant_id"`
	SectionType        string          `json:"section_type"` // hero, featured_products, why_choose_us, contact, testimonials, cta, custom
	SectionName        string          `json:"section_name"`
	DisplayOrder       int             `json:"display_order"`
	IsVisible          bool            `json:"is_visible"`
	TitleEN            string          `json:"title_en"`
	TitleAR            string          `json:"title_ar"`
	SubtitleEN         string          `json:"subtitle_en"`
	SubtitleAR         string          `json:"subtitle_ar"`
	DescriptionEN      string          `json:"description_en"`
	DescriptionAR      string          `json:"description_ar"`
	BackgroundImageURL string          `json:"background_image_url"`
	IconURL            string          `json:"icon_url"`
	Config             json.RawMessage `json:"config"` // Flexible JSON config per section type
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// WebsiteSettings represents global website configuration
type WebsiteSettings struct {
	ID                 int64     `json:"id"`
	RestaurantID       int64     `json:"restaurant_id"`
	TenantID           int64     `json:"tenant_id"`
	SiteTitle          string    `json:"site_title"`
	SiteDescription    string    `json:"site_description"`
	SiteKeywords       string    `json:"site_keywords"`
	Phone              string    `json:"phone"`
	Email              string    `json:"email"`
	Address            string    `json:"address"`
	City               string    `json:"city"`
	Country            string    `json:"country"`
	OpeningTime        string    `json:"opening_time"` // HH:MM format
	ClosingTime        string    `json:"closing_time"` // HH:MM format
	Timezone           string    `json:"timezone"`
	EnableOrders       bool      `json:"enable_orders"`
	EnableDelivery     bool      `json:"enable_delivery"`
	EnableReservations bool      `json:"enable_reservations"`
	FacebookURL        string    `json:"facebook_url"`
	InstagramURL       string    `json:"instagram_url"`
	TwitterURL         string    `json:"twitter_url"`
	MetaTitle          string    `json:"meta_title"`
	MetaDescription    string    `json:"meta_description"`
	OGImageURL         string    `json:"og_image_url"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ThemeTemplate defines pre-built theme templates
type ThemeTemplate struct {
	Name            string `json:"name"`
	PrimaryColor    string `json:"primary_color"`
	SecondaryColor  string `json:"secondary_color"`
	AccentColor     string `json:"accent_color"`
	BackgroundColor string `json:"background_color"`
	TextColor       string `json:"text_color"`
	FontFamily      string `json:"font_family"`
}

// ThemeTemplates contains all available preset themes
var ThemeTemplates = map[string]ThemeTemplate{
	"modern": {
		Name:            "Modern",
		PrimaryColor:    "#3b82f6",
		SecondaryColor:  "#1e40af",
		AccentColor:     "#0ea5e9",
		BackgroundColor: "#ffffff",
		TextColor:       "#1f2937",
		FontFamily:      "Inter",
	},
	"classic": {
		Name:            "Classic",
		PrimaryColor:    "#d97706",
		SecondaryColor:  "#b45309",
		AccentColor:     "#f59e0b",
		BackgroundColor: "#fefce8",
		TextColor:       "#78350f",
		FontFamily:      "Georgia",
	},
	"minimalist": {
		Name:            "Minimalist",
		PrimaryColor:    "#059669",
		SecondaryColor:  "#047857",
		AccentColor:     "#10b981",
		BackgroundColor: "#f0fdf4",
		TextColor:       "#1f2937",
		FontFamily:      "Helvetica",
	},
}

// HomepageData represents complete homepage data for rendering
type HomepageData struct {
	Theme    *LegacyTheme      `json:"theme"`
	Sections []HomepageSection `json:"sections"`
	Settings *WebsiteSettings  `json:"settings"`
	Products []Product         `json:"products"`
}

// SectionConfig helpers for different section types

// HeroSectionConfig configuration for hero section
type HeroSectionConfig struct {
	CTAButtonText  string  `json:"cta_button_text"`
	CTAButtonURL   string  `json:"cta_button_url"`
	Height         string  `json:"height"`          // e.g., "500px"
	OverlayOpacity float64 `json:"overlay_opacity"` // 0-1
	TextAlignment  string  `json:"text_alignment"`  // left, center, right
}

// FeaturedProductsConfig configuration for featured products section
type FeaturedProductsConfig struct {
	MaxProducts      int     `json:"max_products"` // 3-12
	DisplayType      string  `json:"display_type"` // grid, carousel, list
	ShowPrice        bool    `json:"show_price"`
	ShowCategory     bool    `json:"show_category"`
	CategoriesToShow []int64 `json:"categories_to_show"` // null means all
	BackgroundStyle  string  `json:"background_style"`   // light, dark
}

// TestimonialsConfig configuration for testimonials section
type TestimonialsConfig struct {
	MaxTestimonials int    `json:"max_testimonials"`
	DisplayType     string `json:"display_type"` // carousel
	ShowRating      bool   `json:"show_rating"`
	AutoRotate      bool   `json:"auto_rotate"`
	RotationSpeed   int    `json:"rotation_speed"` // milliseconds
}

// CustomSectionConfig configuration for custom HTML section
type CustomSectionConfig struct {
	HTMLContent      string   `json:"html_content"`
	CSSClasses       string   `json:"css_classes"`
	AllowedVariables []string `json:"allowed_variables"` // restaurant_name, phone, email, etc.
}

// Scan implements sql.Scanner interface for json.RawMessage
func (c *HomepageSection) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, c)
}

// Value implements driver.Valuer interface for json.RawMessage
func (c *HomepageSection) Value() (driver.Value, error) {
	return json.Marshal(c)
}
