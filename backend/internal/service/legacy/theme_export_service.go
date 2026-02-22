//go:build legacy
// +build legacy

package service

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ThemeExportService handles theme export and import operations
type ThemeExportService struct {
	themeService *ThemeService
}

// NewThemeExportService creates a new theme export service
func NewThemeExportService(themeService *ThemeService) *ThemeExportService {
	return &ThemeExportService{
		themeService: themeService,
	}
}

// ExportThemeMetadata exports theme metadata as JSON
func (s *ThemeExportService) ExportThemeMetadata(themeID int64) (map[string]interface{}, error) {
	theme, err := s.themeService.GetThemeByID(themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch theme: %w", err)
	}

	return map[string]interface{}{
		"id":          theme.ID,
		"name":        theme.Name,
		"slug":        theme.Slug,
		"description": theme.Description,
		"version":     theme.Version,
		"isActive":    theme.IsActive,
		"isPublished": theme.IsPublished,
		"createdAt":   theme.CreatedAt,
		"updatedAt":   theme.UpdatedAt,
		"publishedAt": theme.PublishedAt,
	}, nil
}

// ExportThemeColors exports only the colors as JSON
func (s *ThemeExportService) ExportThemeColors(themeID int64) (map[string]string, error) {
	theme, err := s.themeService.GetThemeByID(themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch theme: %w", err)
	}

	return map[string]string{
		"primary":     theme.Colors.Primary,
		"secondary":   theme.Colors.Secondary,
		"accent":      theme.Colors.Accent,
		"background":  theme.Colors.Background,
		"text":        theme.Colors.Text,
		"border":      theme.Colors.Border,
		"shadow":      theme.Colors.Shadow,
	}, nil
}

// ExportCompleteTheme exports a complete theme with all related data
func (s *ThemeExportService) ExportCompleteTheme(themeID int64) (map[string]interface{}, error) {
	theme, err := s.themeService.GetThemeByID(themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch theme: %w", err)
	}

	// Get history
	history, err := s.themeService.GetThemeHistory(themeID)
	if err != nil {
		history = []*models.ThemeHistory{}
	}

	return map[string]interface{}{
		"theme": map[string]interface{}{
			"id":          theme.ID,
			"name":        theme.Name,
			"slug":        theme.Slug,
			"description": theme.Description,
			"version":     theme.Version,
			"colors":      theme.Colors,
			"typography":  theme.Typography,
			"identity":    theme.Identity,
			"isActive":    theme.IsActive,
			"isPublished": theme.IsPublished,
			"createdAt":   theme.CreatedAt,
			"updatedAt":   theme.UpdatedAt,
		},
		"components": theme.Components,
		"history":    history,
		"exportDate": time.Now(),
	}, nil
}

// ImportThemeFromJSONString imports a theme from a JSON string
func (s *ThemeExportService) ImportThemeFromJSONString(restaurantID, tenantID int64, jsonString string) (*models.Theme, error) {
	return s.themeService.ImportThemeFromJSON(restaurantID, tenantID, []byte(jsonString))
}

// ExportThemeAsJSONString exports a theme as a JSON string
func (s *ThemeExportService) ExportThemeAsJSONString(themeID int64) (string, error) {
	theme, err := s.themeService.GetThemeByID(themeID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch theme: %w", err)
	}

	data, err := json.MarshalIndent(theme, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal theme: %w", err)
	}

	return string(data), nil
}

// ValidateThemeJSON validates that JSON data contains valid theme information
func (s *ThemeExportService) ValidateThemeJSON(jsonData []byte) error {
	var data struct {
		Name       string `json:"name"`
		Colors     map[string]string `json:"colors"`
		Typography map[string]interface{} `json:"typography"`
		Identity   map[string]string `json:"identity"`
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	if strings.TrimSpace(data.Name) == "" {
		return errors.New("theme name is required")
	}

	if len(data.Colors) == 0 {
		return errors.New("colors are required")
	}

	if data.Colors["primary"] == "" {
		return errors.New("primary color is required")
	}

	if data.Identity["siteTitle"] == "" {
		return errors.New("site title is required")
	}

	return nil
}

// ExportThemeComponentsAsCSV exports theme components as CSV
func (s *ThemeExportService) ExportThemeComponentsAsCSV(themeID int64) (string, error) {
	theme, err := s.themeService.GetThemeByID(themeID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch theme: %w", err)
	}

	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	// Write header
	header := []string{"ID", "Type", "Title", "Order", "Enabled", "Settings"}
	writer.Write(header)

	// Write component rows
	for _, comp := range theme.Components {
		row := []string{
			fmt.Sprintf("%d", comp.ID),
			comp.ComponentType,
			fmt.Sprintf("%v", comp.Title),
			fmt.Sprintf("%d", comp.DisplayOrder),
			fmt.Sprintf("%v", comp.Enabled),
			string(comp.Settings),
		}
		writer.Write(row)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("failed to write CSV: %w", err)
	}

	return sb.String(), nil
}

// ThemeComparison compares two themes
type ThemeComparison struct {
	Theme1ID     int64                          `json:"theme1Id"`
	Theme2ID     int64                          `json:"theme2Id"`
	Theme1Name   string                         `json:"theme1Name"`
	Theme2Name   string                         `json:"theme2Name"`
	Differences  map[string]interface{}         `json:"differences"`
	IsSimilar    bool                           `json:"isSimilar"`
}

// CompareThemes compares two themes and returns their differences
func (s *ThemeExportService) CompareThemes(themeID1, themeID2 int64) (*ThemeComparison, error) {
	theme1, err := s.themeService.GetThemeByID(themeID1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch first theme: %w", err)
	}

	theme2, err := s.themeService.GetThemeByID(themeID2)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch second theme: %w", err)
	}

	comparison := &ThemeComparison{
		Theme1ID:    themeID1,
		Theme2ID:    themeID2,
		Theme1Name:  theme1.Name,
		Theme2Name:  theme2.Name,
		Differences: make(map[string]interface{}),
		IsSimilar:   true,
	}

	// Compare colors
	if theme1.Colors != theme2.Colors {
		comparison.Differences["colors"] = map[string]interface{}{
			"theme1": theme1.Colors,
			"theme2": theme2.Colors,
		}
		comparison.IsSimilar = false
	}

	// Compare typography
	if theme1.Typography != theme2.Typography {
		comparison.Differences["typography"] = map[string]interface{}{
			"theme1": theme1.Typography,
			"theme2": theme2.Typography,
		}
		comparison.IsSimilar = false
	}

	// Compare identity
	if theme1.Identity != theme2.Identity {
		comparison.Differences["identity"] = map[string]interface{}{
			"theme1": theme1.Identity,
			"theme2": theme2.Identity,
		}
		comparison.IsSimilar = false
	}

	// Compare number of components
	if len(theme1.Components) != len(theme2.Components) {
		comparison.Differences["componentCount"] = map[string]int{
			"theme1": len(theme1.Components),
			"theme2": len(theme2.Components),
		}
		comparison.IsSimilar = false
	}

	return comparison, nil
}

// ThemeTemplate represents a theme template for cloning
type ThemeTemplate struct {
	Name        string
	Description string
	Colors      models.ThemeColors
	Typography  models.TypographySettings
	Identity    models.WebsiteIdentity
}

// GetThemeTemplates returns predefined theme templates
func (s *ThemeExportService) GetThemeTemplates() []ThemeTemplate {
	return []ThemeTemplate{
		{
			Name:        "Modern Blue",
			Description: "Clean and professional blue palette",
			Colors: models.ThemeColors{
				Primary:     "#3b82f6",
				Secondary:   "#1e40af",
				Accent:      "#0ea5e9",
				Background:  "#f8fafc",
				Text:        "#1e293b",
				Border:      "#e2e8f0",
				Shadow:      "#0000001a",
			},
			Typography: models.TypographySettings{
				FontFamily:   "Inter, system-ui, sans-serif",
				BaseFontSize: 16,
				BorderRadius: 8,
				LineHeight:   1.5,
			},
			Identity: models.WebsiteIdentity{
				SiteTitle: "Modern Restaurant",
				Domain:    "",
			},
		},
		{
			Name:        "Warm Orange",
			Description: "Warm and inviting orange palette",
			Colors: models.ThemeColors{
				Primary:     "#f97316",
				Secondary:   "#ea580c",
				Accent:      "#fb923c",
				Background:  "#fef7f0",
				Text:        "#1f2937",
				Border:      "#fed7aa",
				Shadow:      "#0000001a",
			},
			Typography: models.TypographySettings{
				FontFamily:   "Poppins, sans-serif",
				BaseFontSize: 16,
				BorderRadius: 12,
				LineHeight:   1.6,
			},
			Identity: models.WebsiteIdentity{
				SiteTitle: "Warm & Cozy",
				Domain:    "",
			},
		},
		{
			Name:        "Elegant Green",
			Description: "Elegant and natural green palette",
			Colors: models.ThemeColors{
				Primary:     "#16a34a",
				Secondary:   "#15803d",
				Accent:      "#22c55e",
				Background:  "#f0fdf4",
				Text:        "#1e293b",
				Border:      "#86efac",
				Shadow:      "#0000001a",
			},
			Typography: models.TypographySettings{
				FontFamily:   "Sora, sans-serif",
				BaseFontSize: 16,
				BorderRadius: 10,
				LineHeight:   1.55,
			},
			Identity: models.WebsiteIdentity{
				SiteTitle: "Natural & Fresh",
				Domain:    "",
			},
		},
		{
			Name:        "Luxury Purple",
			Description: "Luxury and premium purple palette",
			Colors: models.ThemeColors{
				Primary:     "#7c3aed",
				Secondary:   "#6d28d9",
				Accent:      "#a855f7",
				Background:  "#faf5ff",
				Text:        "#1e1b4b",
				Border:      "#e9d5ff",
				Shadow:      "#0000001a",
			},
			Typography: models.TypographySettings{
				FontFamily:   "Cormorant Garamond, serif",
				BaseFontSize: 18,
				BorderRadius: 6,
				LineHeight:   1.7,
			},
			Identity: models.WebsiteIdentity{
				SiteTitle: "Premium Dining",
				Domain:    "",
			},
		},
	}
}
