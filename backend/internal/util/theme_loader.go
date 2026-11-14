package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// ThemeConfig represents the complete theme configuration from JSON file
type ThemeConfig struct {
	Meta struct {
		Name        string   `json:"name"`
		Slug        string   `json:"slug"`
		Version     string   `json:"version"`
		Author      string   `json:"author"`
		Description string   `json:"description"`
		Category    string   `json:"category"`
		Tags        []string `json:"tags"`
		Preview     string   `json:"preview"`
		CreatedAt   string   `json:"created_at"`
		Bilingual   bool     `json:"bilingual"`
	} `json:"meta"`

	Identity struct {
		SiteTitle  interface{} `json:"siteTitle"`
		LogoUrl    string      `json:"logoUrl"`
		FaviconUrl string      `json:"faviconUrl"`
		Domain     string      `json:"domain"`
	} `json:"identity"`

	Colors struct {
		Primary    string `json:"primary"`
		Secondary  string `json:"secondary"`
		Accent     string `json:"accent"`
		Background string `json:"background"`
		Text       string `json:"text"`
		Border     string `json:"border"`
		Shadow     string `json:"shadow"`
	} `json:"colors"`

	Typography struct {
		FontFamily   string                 `json:"fontFamily"`
		BaseFontSize int                    `json:"baseFontSize"`
		LineHeight   float64                `json:"lineHeight"`
		BorderRadius int                    `json:"borderRadius"`
		Headings     map[string]interface{} `json:"headings"`
	} `json:"typography"`

	Header struct {
		Style                string                   `json:"style"`
		Layout               string                   `json:"layout"`
		Position             string                   `json:"position"`
		Height               int                      `json:"height"`
		Padding              int                      `json:"padding"`
		BackgroundColor      string                   `json:"backgroundColor"`
		TextColor            string                   `json:"textColor"`
		LogoPosition         string                   `json:"logoPosition"`
		ShowLogo             bool                     `json:"showLogo"`
		ShowSearch           bool                     `json:"showSearch"`
		ShowLanguageSwitcher bool                     `json:"showLanguageSwitcher"`
		NavigationItems      []map[string]interface{} `json:"navigationItems"`
		CtaButton            map[string]interface{}   `json:"ctaButton"`
	} `json:"header"`

	Footer struct {
		Style           string                   `json:"style"`
		Layout          string                   `json:"layout"`
		Columns         int                      `json:"columns"`
		BackgroundColor string                   `json:"backgroundColor"`
		TextColor       string                   `json:"textColor"`
		LinkColor       string                   `json:"linkColor"`
		ShowBackToTop   bool                     `json:"showBackToTop"`
		CompanyInfo     map[string]interface{}   `json:"companyInfo"`
		Sections        []map[string]interface{} `json:"sections"`
		SocialLinks     []map[string]interface{} `json:"socialLinks"`
		LegalLinks      []map[string]interface{} `json:"legalLinks"`
		CopyrightText   interface{}              `json:"copyrightText"`
	} `json:"footer"`

	Components []map[string]interface{} `json:"components"`

	Customization struct {
		AllowColorChange      bool `json:"allowColorChange"`
		AllowFontChange       bool `json:"allowFontChange"`
		AllowLayoutChange     bool `json:"allowLayoutChange"`
		AllowComponentReorder bool `json:"allowComponentReorder"`
		AllowComponentDisable bool `json:"allowComponentDisable"`
	} `json:"customization"`
}

// LoadThemeFromFile loads a single theme from JSON file
func LoadThemeFromFile(themesDir, themeSlug string) (*ThemeConfig, error) {
	themeFile := filepath.Join(themesDir, themeSlug, "theme.json")

	fmt.Printf("📂 Loading theme from: %s\n", themeFile)

	data, err := os.ReadFile(themeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file %s: %w", themeFile, err)
	}

	var theme ThemeConfig
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme JSON %s: %w", themeFile, err)
	}

	fmt.Printf("✅ Loaded theme: %s (slug: %s)\n", theme.Meta.Name, theme.Meta.Slug)
	return &theme, nil
}

// LoadAllThemes loads all themes from the themes directory
func LoadAllThemes(themesDir string) (map[string]*ThemeConfig, error) {
	themes := make(map[string]*ThemeConfig)

	fmt.Printf("🔍 Scanning themes directory: %s\n", themesDir)

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read themes directory: %w", err)
	}

	for _, entry := range entries {
		// Skip non-directories and special directories
		if !entry.IsDir() || entry.Name() == "_template" {
			continue
		}

		themeSlug := entry.Name()
		theme, err := LoadThemeFromFile(themesDir, themeSlug)
		if err != nil {
			fmt.Printf("⚠️  Skipping theme %s: %v\n", themeSlug, err)
			continue
		}

		themes[themeSlug] = theme
	}

	fmt.Printf("\n✅ Successfully loaded %d themes\n", len(themes))
	return themes, nil
}

// GetThemeMetadata extracts just the metadata from a theme
func GetThemeMetadata(theme *ThemeConfig) map[string]interface{} {
	return map[string]interface{}{
		"name":        theme.Meta.Name,
		"slug":        theme.Meta.Slug,
		"description": theme.Meta.Description,
		"category":    theme.Meta.Category,
		"version":     theme.Meta.Version,
		"author":      theme.Meta.Author,
		"tags":        theme.Meta.Tags,
		"preview":     theme.Meta.Preview,
		"bilingual":   theme.Meta.Bilingual,
	}
}

// ThemeToJSON converts a theme to a JSON byte array (for database storage)
func ThemeToJSON(theme *ThemeConfig) ([]byte, error) {
	return json.Marshal(theme)
}

// ValidateTheme checks if a theme has all required fields
func ValidateTheme(theme *ThemeConfig) []string {
	var errors []string

	if theme.Meta.Name == "" {
		errors = append(errors, "meta.name is required")
	}
	if theme.Meta.Slug == "" {
		errors = append(errors, "meta.slug is required")
	}
	if theme.Colors.Primary == "" {
		errors = append(errors, "colors.primary is required")
	}
	if theme.Colors.Background == "" {
		errors = append(errors, "colors.background is required")
	}
	if len(theme.Components) == 0 {
		errors = append(errors, "at least one component is required")
	}

	return errors
}

// ListThemes returns a sorted list of all theme slugs
func ListThemes(themes map[string]*ThemeConfig) []string {
	slugs := make([]string, 0, len(themes))
	for slug := range themes {
		slugs = append(slugs, slug)
	}
	sort.Strings(slugs)
	return slugs
}

// GetThemeBySlug retrieves a theme by its slug
func GetThemeBySlug(themes map[string]*ThemeConfig, slug string) *ThemeConfig {
	return themes[slug]
}
