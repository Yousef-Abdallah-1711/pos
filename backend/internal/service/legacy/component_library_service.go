//go:build legacy
// +build legacy

package service

import (
	"errors"
	"fmt"

	"pos-saas/internal/repository"
)

// ComponentLibraryService handles business logic for component library
type ComponentLibraryService struct {
	libraryRepo repository.ComponentLibraryRepositoryV2
}

// NewComponentLibraryService creates a new component library service
func NewComponentLibraryService(libraryRepo repository.ComponentLibraryRepositoryV2) *ComponentLibraryService {
	return &ComponentLibraryService{
		libraryRepo: libraryRepo,
	}
}

// GetAllLibrary retrieves all component library entries
func (s *ComponentLibraryService) GetAllLibrary() ([]*models.ComponentLibrary, error) {
	return s.libraryRepo.List()
}

// GetActiveLibrary retrieves only active component library entries
func (s *ComponentLibraryService) GetActiveLibrary() ([]*models.ComponentLibrary, error) {
	return s.libraryRepo.GetActive()
}

// GetLibraryByType retrieves a library entry by component type
func (s *ComponentLibraryService) GetLibraryByType(componentType string) (*models.ComponentLibrary, error) {
	if componentType == "" {
		return nil, errors.New("component type is required")
	}

	return s.libraryRepo.GetByType(componentType)
}

// GetLibraryByID retrieves a library entry by ID
func (s *ComponentLibraryService) GetLibraryByID(id int64) (*models.ComponentLibrary, error) {
	if id == 0 {
		return nil, errors.New("library id is required")
	}

	return s.libraryRepo.GetByID(id)
}

// CreateLibraryEntry creates a new library entry
func (s *ComponentLibraryService) CreateLibraryEntry(library *models.ComponentLibrary) (*models.ComponentLibrary, error) {
	if library == nil {
		return nil, errors.New("library cannot be nil")
	}

	if library.ComponentType == "" {
		return nil, errors.New("component type is required")
	}

	if library.Title == "" {
		return nil, errors.New("title is required")
	}

	if !library.IsActive {
		library.IsActive = true
	}

	return s.libraryRepo.Create(library)
}

// UpdateLibraryEntry updates an existing library entry
func (s *ComponentLibraryService) UpdateLibraryEntry(library *models.ComponentLibrary) (*models.ComponentLibrary, error) {
	if library == nil {
		return nil, errors.New("library cannot be nil")
	}

	if library.ID == 0 {
		return nil, errors.New("library id is required")
	}

	return s.libraryRepo.Update(library)
}

// DeleteLibraryEntry deletes a library entry
func (s *ComponentLibraryService) DeleteLibraryEntry(id int64) error {
	if id == 0 {
		return errors.New("library id is required")
	}

	return s.libraryRepo.Delete(id)
}

// ValidateLibraryEntry validates a library entry
func (s *ComponentLibraryService) ValidateLibraryEntry(library *models.ComponentLibrary) error {
	if library == nil {
		return errors.New("library cannot be nil")
	}

	if library.ComponentType == "" {
		return errors.New("component type is required")
	}

	if library.Title == "" {
		return errors.New("title is required")
	}

	return nil
}

// GetComponentTypeCount returns the number of component types in the library
func (s *ComponentLibraryService) GetComponentTypeCount() (int, error) {
	libraries, err := s.libraryRepo.List()
	if err != nil {
		return 0, fmt.Errorf("failed to count component types: %w", err)
	}

	return len(libraries), nil
}

// GetAvailableComponentTypes returns a list of available component types
func (s *ComponentLibraryService) GetAvailableComponentTypes() ([]string, error) {
	libraries, err := s.libraryRepo.GetActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get available types: %w", err)
	}

	types := make([]string, len(libraries))
	for i, lib := range libraries {
		types[i] = lib.ComponentType
	}

	return types, nil
}
