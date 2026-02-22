//go:build legacy
// +build legacy

package models

import (
	"encoding/json"
	"time"
)

// BilingualText represents text content in both English and Arabic
type BilingualText struct {
	EN string `json:"en" db:"text_en"`   // English text
	AR string `json:"ar" db:"text_ar"`   // Arabic text (RTL)
}

// IsEmpty returns true if both English and Arabic text are empty
func (bt *BilingualText) IsEmpty() bool {
	if bt == nil {
		return true
	}
	return bt.EN == "" && bt.AR == ""
}

// GetText returns the text for the specified language, defaults to English if not available
func (bt *BilingualText) GetText(lang string) string {
	if bt == nil {
		return ""
	}
	if lang == "ar" || lang == "AR" {
		if bt.AR != "" {
			return bt.AR
		}
	}
	return bt.EN
}

// ThemeColors represents the 7 global colors used in a theme
type ThemeColors struct {
	Primary     string `json:"primary" db:"primary_color"`
	Secondary   string `json:"secondary" db:"secondary_color"`
	Accent      string `json:"accent" db:"accent_color"`
	Background  string `json:"background" db:"background_color"`
	Text        string `json:"text" db:"text_color"`
	Border      string `json:"border" db:"border_color"`
	Shadow      string `json:"shadow" db:"shadow_color"`
}

// TypographySettings represents typography configuration
type TypographySettings struct {
	FontFamily  string  `json:"fontFamily" db:"font_family"`
	BaseFontSize int   `json:"baseFontSize" db:"base_font_size"`
	BorderRadius int   `json:"borderRadius" db:"border_radius"`
	LineHeight  float64 `json:"lineHeight" db:"line_height"`
}

// WebsiteIdentity represents brand and identity information
type WebsiteIdentity struct {
	SiteTitle  string `json:"siteTitle" db:"site_title"`
	LogoURL    string `json:"logoUrl" db:"logo_url"`
	FaviconURL string `json:"faviconUrl" db:"favicon_url"`
	Domain     string `json:"domain" db:"domain"`
}

// Note: Database fields may be NULL, so when scanning we need to handle sql.NullString
// The JSON marshaling will handle empty strings appropriately

// HeaderConfig represents header configuration
type HeaderConfig struct {
	LogoURL         string `json:"logoUrl,omitempty"`
	LogoText        string `json:"logoText,omitempty"`
	ShowLogo        bool   `json:"showLogo"`
	LogoHeight      int    `json:"logoHeight"`
	NavigationItems []NavItem `json:"navigationItems"`
	BackgroundColor string `json:"backgroundColor"`
	TextColor       string `json:"textColor"`
	Height          int    `json:"height"`
	Padding         int    `json:"padding"`
	ShowShadow      bool   `json:"showShadow"`
	StickyHeader    bool   `json:"stickyHeader"`
	HideNavOnMobile bool   `json:"hideNavOnMobile"`
	NavPosition     string `json:"navPosition"` // left, center, right
}

// NavItem represents a navigation menu item
type NavItem struct {
	ID    string `json:"id,omitempty"`
	Label string `json:"label"`
	Href  string `json:"href"`
	Order int    `json:"order"`
}

// FooterConfig represents footer configuration
type FooterConfig struct {
	CompanyName    string `json:"companyName,omitempty"`
	CompanyDescription string `json:"companyDescription,omitempty"`
	Address        string `json:"address,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Email          string `json:"email,omitempty"`
	CopyrightText  string `json:"copyrightText"`
	SocialLinks    []SocialLink `json:"socialLinks"`
	FooterSections []FooterSection `json:"footerSections,omitempty"`
	LegalLinks     []FooterLink `json:"legalLinks"`
	BackgroundColor string `json:"backgroundColor"`
	TextColor      string `json:"textColor"`
	ShowLinks      bool   `json:"showLinks"`
}

// SocialLink represents a social media link
type SocialLink struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

// FooterSection represents a section in the footer
type FooterSection struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Links []FooterLink `json:"links"`
}

// FooterLink represents a link in the footer
type FooterLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// Theme represents a complete theme configuration
type Theme struct {
	ID             int64                 `json:"id" db:"id"`
	RestaurantID   int64                 `json:"restaurantId" db:"restaurant_id"`
	TenantID       int64                 `json:"tenantId" db:"tenant_id"`
	Name           string                `json:"name" db:"name"`
	Slug           string                `json:"slug" db:"slug"`
	Description    string                `json:"description" db:"description"`
	Version        int                   `json:"version" db:"version"`
	IsActive       bool                  `json:"isActive" db:"is_active"`
	IsPublished    bool                  `json:"isPublished" db:"is_published"`
	Colors         ThemeColors           `json:"colors"`
	Typography     TypographySettings    `json:"typography"`
	Identity       WebsiteIdentity       `json:"identity"`
	Header         *HeaderConfig         `json:"header,omitempty"`
	Footer         *FooterConfig         `json:"footer,omitempty"`
	Components     []ThemeComponent      `json:"components,omitempty"`
	CreatedAt      time.Time             `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time             `json:"updatedAt" db:"updated_at"`
	PublishedAt    *time.Time            `json:"publishedAt" db:"published_at"`
}

// ThemeComponent represents a component within a theme
type ThemeComponent struct {
	ID            int64             `json:"id" db:"id"`
	ThemeID       int64             `json:"themeId" db:"theme_id"`
	ComponentType string            `json:"type" db:"component_type"`
	Title         BilingualText     `json:"title"`                      // EN/AR bilingual title
	Subtitle      *BilingualText    `json:"subtitle,omitempty"`         // Optional EN/AR subtitle
	Description   *BilingualText    `json:"description,omitempty"`      // Optional EN/AR description
	ButtonText    *BilingualText    `json:"buttonText,omitempty"`       // Optional EN/AR button text
	DisplayOrder  int               `json:"order" db:"display_order"`
	Enabled       bool              `json:"enabled" db:"enabled"`
	Settings      json.RawMessage   `json:"settings" db:"settings"`
	CreatedAt     time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time         `json:"updatedAt" db:"updated_at"`
}

// ComponentLibrary represents a reusable component type definition
type ComponentLibrary struct {
	ID               int64           `json:"id" db:"id"`
	ComponentType    string          `json:"componentType" db:"component_type"`
	Title            string          `json:"title" db:"title"`
	Description      string          `json:"description" db:"description"`
	IconURL          string          `json:"iconUrl" db:"icon_url"`
	DefaultSettings  json.RawMessage `json:"defaultSettings" db:"default_settings"`
	SettingsSchema   json.RawMessage `json:"settingsSchema" db:"settings_schema"`
	IsActive         bool            `json:"isActive" db:"is_active"`
	CreatedAt        time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time       `json:"updatedAt" db:"updated_at"`
}

// CreateThemeRequest represents the request to create a new theme
type CreateThemeRequest struct {
	Name          string                `json:"name" validate:"required"`
	Description   string                `json:"description"`
	Colors        ThemeColors           `json:"colors" validate:"required"`
	Typography    TypographySettings    `json:"typography" validate:"required"`
	Identity      WebsiteIdentity       `json:"identity" validate:"required"`
	Components    []ThemeComponent      `json:"components"`
}

// UpdateThemeRequest represents the request to update an existing theme
type UpdateThemeRequest struct {
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	Colors        *ThemeColors          `json:"colors"`
	Typography    *TypographySettings   `json:"typography"`
	Identity      *WebsiteIdentity      `json:"identity"`
	Header        *HeaderConfig         `json:"header"`
	Footer        *FooterConfig         `json:"footer"`
	Components    []ThemeComponent      `json:"components"`
}

// ThemeResponse represents the response with a theme
type ThemeResponse struct {
	Data  *Theme `json:"data"`
	Error string `json:"error,omitempty"`
}

// ThemesResponse represents the response with multiple themes
type ThemesResponse struct {
	Data  []*Theme `json:"data"`
	Error string   `json:"error,omitempty"`
}
