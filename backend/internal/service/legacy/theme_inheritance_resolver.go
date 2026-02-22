//go:build legacy
// +build legacy

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"pos-saas/internal/repository"
)

// ThemeInheritanceResolver handles theme property resolution with inheritance chain
type ThemeInheritanceResolver struct {
	themeRepo      *repository.ThemeRepository
	inheritanceRepo *repository.ThemeInheritanceRepository
	cache          map[string]*models.ResolvedTheme
	cacheMutex     sync.RWMutex
}

// NewThemeInheritanceResolver creates a new resolver instance
func NewThemeInheritanceResolver(
	themeRepo *repository.ThemeRepository,
	inheritanceRepo *repository.ThemeInheritanceRepository,
) *ThemeInheritanceResolver {
	return &ThemeInheritanceResolver{
		themeRepo:       themeRepo,
		inheritanceRepo: inheritanceRepo,
		cache:           make(map[string]*models.ResolvedTheme),
	}
}

// ResolveTheme returns a fully resolved theme with all inherited properties
func (r *ThemeInheritanceResolver) ResolveTheme(ctx context.Context, themeID string) (*models.ResolvedTheme, error) {
	// Check cache first
	r.cacheMutex.RLock()
	if cached, exists := r.cache[themeID]; exists {
		r.cacheMutex.RUnlock()
		return cached, nil
	}
	r.cacheMutex.RUnlock()

	// Build inheritance chain and resolve
	inheritanceChain, err := r.buildInheritanceChain(ctx, themeID)
	if err != nil {
		return nil, err
	}

	// Resolve properties from bottom to top (child to parent)
	resolved, err := r.resolveProperties(ctx, inheritanceChain)
	if err != nil {
		return nil, err
	}

	// Cache the result
	r.cacheMutex.Lock()
	r.cache[themeID] = resolved
	r.cacheMutex.Unlock()

	return resolved, nil
}

// buildInheritanceChain builds the inheritance chain from child to root
func (r *ThemeInheritanceResolver) buildInheritanceChain(ctx context.Context, themeID string) ([]*models.ThemeInheritanceNode, error) {
	var chain []*models.ThemeInheritanceNode
	visitedIDs := make(map[string]bool)
	currentID := themeID
	level := 0

	// Prevent infinite loops from circular references
	maxLevels := 10

	for level < maxLevels {
		// Check for circular reference
		if visitedIDs[currentID] {
			return nil, fmt.Errorf("circular theme inheritance detected at level %d", level)
		}
		visitedIDs[currentID] = true

		// Get current theme
		theme, err := r.themeRepo.GetThemeByID(ctx, currentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get theme %s: %w", currentID, err)
		}

		// Add to chain
		node := &models.ThemeInheritanceNode{
			ID:        theme.ID,
			Name:      theme.Name,
			Level:     level,
			Overrides: make(map[string]interface{}),
		}

		// Get property overrides for this theme
		overrides, err := r.inheritanceRepo.GetPropertyOverrides(ctx, currentID)
		if err == nil && overrides != nil {
			for _, override := range overrides {
				if override.OverrideValue.Valid {
					var value interface{}
					json.Unmarshal([]byte(override.OverrideValue.String), &value)
					node.Overrides[override.PropertyPath] = value
				}
			}
		}

		chain = append(chain, node)

		// Move to parent if exists
		if theme.ParentThemeID.Valid && theme.ParentThemeID.String != "" {
			currentID = theme.ParentThemeID.String
			level++
		} else {
			break
		}
	}

	if level >= maxLevels {
		return nil, fmt.Errorf("theme inheritance too deep (max %d levels)", maxLevels)
	}

	return chain, nil
}

// resolveProperties resolves all properties from the inheritance chain
func (r *ThemeInheritanceResolver) resolveProperties(ctx context.Context, chain []*models.ThemeInheritanceNode) (*models.ResolvedTheme, error) {
	if len(chain) == 0 {
		return nil, fmt.Errorf("empty inheritance chain")
	}

	// Start with the root theme (end of chain) and apply overrides going back to child
	rootNode := chain[len(chain)-1]
	rootTheme, err := r.themeRepo.GetThemeByID(ctx, rootNode.ID)
	if err != nil {
		return nil, err
	}

	// Build resolved theme starting from root
	resolved := &models.ResolvedTheme{
		ID:              rootTheme.ID,
		Name:            rootTheme.Name,
		Description:     rootTheme.Description,
		Slug:            rootTheme.Slug,
		SiteTitle:       rootTheme.WebsiteIdentity.SiteTitle,
		PrimaryColor:    rootTheme.Colors.Primary,
		SecondaryColor:  rootTheme.Colors.Secondary,
		AccentColor:     rootTheme.Colors.Accent,
		BackgroundColor: rootTheme.Colors.Background,
		TextColor:       rootTheme.Colors.Text,
		BorderColor:     rootTheme.Colors.Border,
		ShadowColor:     rootTheme.Colors.Shadow,
		FontFamily:      rootTheme.Typography.FontFamily,
		BaseFontSize:    rootTheme.Typography.BaseFontSize,
		BorderRadius:    rootTheme.Typography.BorderRadius,
		LineHeight:      rootTheme.Typography.LineHeight,
		LogoURL:         rootTheme.WebsiteIdentity.LogoURL,
		FaviconURL:      rootTheme.WebsiteIdentity.FaviconURL,
		IsPublished:     rootTheme.IsPublished,
		IsActive:        rootTheme.IsActive,
		CreatedAt:       rootTheme.CreatedAt,
		UpdatedAt:       rootTheme.UpdatedAt,
		InheritanceChain: chain,
	}

	// Apply overrides from parent to child (reverse order)
	for i := len(chain) - 1; i >= 0; i-- {
		node := chain[i]
		r.applyOverrides(resolved, node.Overrides)
	}

	return resolved, nil
}

// applyOverrides applies property overrides to a resolved theme
func (r *ThemeInheritanceResolver) applyOverrides(resolved *models.ResolvedTheme, overrides map[string]interface{}) {
	for path, value := range overrides {
		// Simple property mapping (can be expanded for nested properties)
		switch path {
		case "primaryColor":
			if str, ok := value.(string); ok {
				resolved.PrimaryColor = str
			}
		case "secondaryColor":
			if str, ok := value.(string); ok {
				resolved.SecondaryColor = str
			}
		case "accentColor":
			if str, ok := value.(string); ok {
				resolved.AccentColor = str
			}
		case "backgroundColor":
			if str, ok := value.(string); ok {
				resolved.BackgroundColor = str
			}
		case "textColor":
			if str, ok := value.(string); ok {
				resolved.TextColor = str
			}
		case "borderColor":
			if str, ok := value.(string); ok {
				resolved.BorderColor = str
			}
		case "shadowColor":
			if str, ok := value.(string); ok {
				resolved.ShadowColor = str
			}
		case "fontFamily":
			if str, ok := value.(string); ok {
				resolved.FontFamily = str
			}
		case "baseFontSize":
			if num, ok := value.(float64); ok {
				resolved.BaseFontSize = int(num)
			}
		case "borderRadius":
			if num, ok := value.(float64); ok {
				resolved.BorderRadius = int(num)
			}
		case "lineHeight":
			if num, ok := value.(float64); ok {
				resolved.LineHeight = num
			}
		case "siteTitle":
			if str, ok := value.(string); ok {
				resolved.SiteTitle = str
			}
		case "logoUrl":
			if str, ok := value.(string); ok {
				resolved.LogoURL = str
			}
		case "faviconUrl":
			if str, ok := value.(string); ok {
				resolved.FaviconURL = str
			}
		}
	}
}

// ClearCache clears the theme resolution cache
func (r *ThemeInheritanceResolver) ClearCache() {
	r.cacheMutex.Lock()
	r.cache = make(map[string]*models.ResolvedTheme)
	r.cacheMutex.Unlock()
}

// InvalidateThemeCache invalidates cache for a specific theme and its children
func (r *ThemeInheritanceResolver) InvalidateThemeCache(themeID string) {
	r.cacheMutex.Lock()
	delete(r.cache, themeID)
	r.cacheMutex.Unlock()
	// TODO: Also invalidate all child themes
}
