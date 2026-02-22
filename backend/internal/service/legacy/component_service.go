//go:build legacy
// +build legacy

package service

import (
	"context"
	"encoding/json"
	"fmt"

	"pos-saas/internal/repository"
)

// ComponentService handles business logic for components
type ComponentService struct {
	componentRepo *repository.ComponentRepository
}

// NewComponentService creates a new component service
func NewComponentService(componentRepo *repository.ComponentRepository) *ComponentService {
	return &ComponentService{
		componentRepo: componentRepo,
	}
}

// RegisterComponent registers a new component
func (s *ComponentService) RegisterComponent(ctx context.Context, component *models.ComponentRegistry) error {
	// Validate component
	if component.Name == "" {
		return fmt.Errorf("component name is required")
	}
	if component.ComponentID == "" {
		return fmt.Errorf("component ID is required")
	}

	return s.componentRepo.CreateComponentRegistry(ctx, component)
}

// GetComponent retrieves a component by ID
func (s *ComponentService) GetComponent(ctx context.Context, id string) (*models.ComponentRegistry, error) {
	return s.componentRepo.GetComponentByID(ctx, id)
}

// GetComponentByComponentID retrieves a component by its component ID (e.g., "hero")
func (s *ComponentService) GetComponentByComponentID(ctx context.Context, componentID string) (*models.ComponentRegistry, error) {
	return s.componentRepo.GetComponentByComponentID(ctx, componentID)
}

// ListComponents retrieves all components with optional filtering
func (s *ComponentService) ListComponents(ctx context.Context, category string, limit int, offset int) ([]*models.ComponentRegistry, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	return s.componentRepo.ListComponents(ctx, category, limit, offset)
}

// UpdateComponent updates an existing component
func (s *ComponentService) UpdateComponent(ctx context.Context, component *models.ComponentRegistry) error {
	return s.componentRepo.UpdateComponentRegistry(ctx, component)
}

// DeleteComponent soft-deletes a component
func (s *ComponentService) DeleteComponent(ctx context.Context, id string) error {
	return s.componentRepo.DeleteComponentRegistry(ctx, id)
}

// AddComponentToTheme adds a component to a theme
func (s *ComponentService) AddComponentToTheme(ctx context.Context, themeID string, componentID string, config map[string]interface{}) error {
	// Convert config to JSON
	var configJSON *string
	if config != nil {
		jsonBytes, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}
		jsonStr := string(jsonBytes)
		configJSON = &jsonStr
	}

	mapping := &models.ComponentThemeMapping{
		ThemeID:     themeID,
		ComponentID: componentID,
		IsEnabled:   true,
		Config:      *configJSON,
	}

	return s.componentRepo.CreateComponentThemeMapping(ctx, mapping)
}

// GetThemeComponents retrieves all components used in a theme
func (s *ComponentService) GetThemeComponents(ctx context.Context, themeID string) ([]*models.ComponentThemeMapping, error) {
	return s.componentRepo.GetComponentsByTheme(ctx, themeID)
}

// SyncDiscoveredComponents registers all discovered components from the frontend
func (s *ComponentService) SyncDiscoveredComponents(ctx context.Context, discoveredComponents []map[string]interface{}) error {
	for _, comp := range discoveredComponents {
		component := &models.ComponentRegistry{
			ComponentID: comp["id"].(string),
			Name:        comp["name"].(string),
			Description: comp["description"].(string),
			Version:     comp["version"].(string),
			Category:    "sections",
			IsActive:    true,
		}

		// Handle aliases
		if aliases, ok := comp["aliases"].([]string); ok {
			aliasesJSON, _ := json.Marshal(aliases)
			component.Aliases.String = string(aliasesJSON)
			component.Aliases.Valid = true
		}

		if err := s.componentRepo.CreateComponentRegistry(ctx, component); err != nil {
			// Component already exists, that's fine
			continue
		}
	}

	return nil
}
