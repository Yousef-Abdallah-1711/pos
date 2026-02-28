package usecase

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

var tableTextSanitizer = regexp.MustCompile(`<[^>]*>`)

// TableUseCase handles table-management business logic.
type TableUseCase struct {
	tableRepo *repository.TableRepository
}

// NewTableUseCase creates a new table use case.
func NewTableUseCase(tableRepo *repository.TableRepository) *TableUseCase {
	return &TableUseCase{tableRepo: tableRepo}
}

// CreateTable validates and creates a new table.
func (uc *TableUseCase) CreateTable(tenantID, restaurantID int64, req *domain.CreateTableRequest) (*domain.RestaurantTable, error) {
	if uc.tableRepo == nil {
		return nil, errors.New("table repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}
	if req == nil {
		return nil, errors.New("table payload is required")
	}

	tableNumber := req.TableNumber
	if tableNumber <= 0 {
		return nil, errors.New("table_number must be greater than 0")
	}
	if req.Capacity <= 0 || req.Capacity > 50 {
		return nil, errors.New("capacity must be between 1 and 50")
	}

	location := strings.TrimSpace(req.Location)
	if location == "" {
		location = string(domain.TableLocationIndoor)
	}
	if !domain.ValidTableLocation(location) {
		return nil, errors.New("invalid table location")
	}

	shape := strings.TrimSpace(req.Shape)
	if shape == "" {
		shape = string(domain.TableShapeSquare)
	}
	if !domain.ValidTableShape(shape) {
		return nil, errors.New("invalid table shape")
	}

	table := &domain.RestaurantTable{
		TenantID:     tenantID,
		RestaurantID: restaurantID,
		TableNumber:  tableNumber,
		Name:         sanitizeTableText(req.Name, 120),
		Capacity:     req.Capacity,
		Location:     domain.TableLocation(location),
		Status:       domain.TableStatusAvailable,
		Description:  sanitizeTableText(req.Description, 500),
		Shape:        domain.TableShape(shape),
	}

	created, err := uc.tableRepo.CreateTable(table)
	if err != nil {
		return nil, err
	}
	return created, nil
}

// GetTable returns one table by ID.
func (uc *TableUseCase) GetTable(tenantID, restaurantID, tableID int64) (*domain.RestaurantTable, error) {
	if uc.tableRepo == nil {
		return nil, errors.New("table repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 || tableID <= 0 {
		return nil, errors.New("invalid table context")
	}
	return uc.tableRepo.GetTableByID(tenantID, restaurantID, tableID)
}

// ListTables returns tables for a restaurant, optionally filtered by status.
func (uc *TableUseCase) ListTables(tenantID, restaurantID int64, status string) ([]domain.RestaurantTable, error) {
	if uc.tableRepo == nil {
		return nil, errors.New("table repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}

	normalizedStatus := strings.TrimSpace(status)
	if normalizedStatus != "" && !domain.ValidTableStatus(normalizedStatus) {
		return nil, errors.New("invalid table status")
	}

	return uc.tableRepo.ListTables(tenantID, restaurantID, normalizedStatus)
}

// UpdateTable validates and updates table fields.
func (uc *TableUseCase) UpdateTable(tenantID, restaurantID, tableID int64, req *domain.UpdateTableRequest) error {
	if uc.tableRepo == nil {
		return errors.New("table repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 || tableID <= 0 {
		return errors.New("invalid table context")
	}
	if req == nil {
		return errors.New("update payload is required")
	}

	if req.TableNumber != nil && *req.TableNumber <= 0 {
		return errors.New("table_number must be greater than 0")
	}
	if req.Name != nil {
		sanitized := sanitizeTableText(*req.Name, 120)
		req.Name = &sanitized
	}
	if req.Description != nil {
		sanitized := sanitizeTableText(*req.Description, 500)
		req.Description = &sanitized
	}
	if req.Capacity != nil && (*req.Capacity <= 0 || *req.Capacity > 50) {
		return errors.New("capacity must be between 1 and 50")
	}
	if req.Location != nil && !domain.ValidTableLocation(strings.TrimSpace(*req.Location)) {
		return errors.New("invalid table location")
	}
	if req.Status != nil && !domain.ValidTableStatus(strings.TrimSpace(*req.Status)) {
		return errors.New("invalid table status")
	}
	if req.Shape != nil && !domain.ValidTableShape(strings.TrimSpace(*req.Shape)) {
		return errors.New("invalid table shape")
	}
	if req.LayoutID != nil && *req.LayoutID <= 0 {
		return errors.New("layout_id must be greater than 0")
	}

	return uc.tableRepo.UpdateTable(tenantID, restaurantID, tableID, req)
}

// DeleteTable deletes a table if there are no active orders linked to it.
func (uc *TableUseCase) DeleteTable(tenantID, restaurantID, tableID int64) error {
	if uc.tableRepo == nil {
		return errors.New("table repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 || tableID <= 0 {
		return errors.New("invalid table context")
	}

	hasActiveOrders, err := uc.tableRepo.HasActiveOrders(tenantID, restaurantID, tableID)
	if err != nil {
		return err
	}
	if hasActiveOrders {
		return errors.New("cannot delete table with active orders")
	}

	return uc.tableRepo.DeleteTable(tenantID, restaurantID, tableID)
}

// UpdateTableStatus changes a table status with transition validation.
func (uc *TableUseCase) UpdateTableStatus(tenantID, restaurantID, tableID int64, newStatus string) error {
	if uc.tableRepo == nil {
		return errors.New("table repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 || tableID <= 0 {
		return errors.New("invalid table context")
	}

	status := strings.TrimSpace(newStatus)
	if !domain.ValidTableStatus(status) {
		return errors.New("invalid table status")
	}

	current, err := uc.tableRepo.GetTableByID(tenantID, restaurantID, tableID)
	if err != nil {
		return err
	}
	if current.Status == domain.TableStatus(status) {
		return nil
	}

	if !canTransitionTableStatus(current.Status, domain.TableStatus(status)) {
		return fmt.Errorf("cannot transition table status from %s to %s", current.Status, status)
	}

	return uc.tableRepo.UpdateTableStatus(tenantID, restaurantID, tableID, status)
}

// FreeTable sets a table back to available.
func (uc *TableUseCase) FreeTable(tenantID, restaurantID, tableID int64) error {
	return uc.UpdateTableStatus(tenantID, restaurantID, tableID, string(domain.TableStatusAvailable))
}

// GetTableStats returns table statistics for a restaurant.
func (uc *TableUseCase) GetTableStats(tenantID, restaurantID int64) (*domain.TableStats, error) {
	if uc.tableRepo == nil {
		return nil, errors.New("table repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}
	return uc.tableRepo.GetTableStats(tenantID, restaurantID)
}

// UpdateTablePositions bulk updates floor-plan positions.
func (uc *TableUseCase) UpdateTablePositions(tenantID, restaurantID int64, positions []domain.TablePosition) error {
	if uc.tableRepo == nil {
		return errors.New("table repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return errors.New("invalid tenant context")
	}
	if len(positions) == 0 {
		return errors.New("positions payload is required")
	}
	for _, pos := range positions {
		if pos.TableID <= 0 {
			return errors.New("table_id must be greater than 0")
		}
	}
	return uc.tableRepo.UpdateTablePositions(tenantID, restaurantID, positions)
}

// CreateLayout creates a floor-plan layout.
func (uc *TableUseCase) CreateLayout(tenantID, restaurantID int64, req *domain.CreateLayoutRequest) (*domain.TableLayout, error) {
	if uc.tableRepo == nil {
		return nil, errors.New("table repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}
	if req == nil {
		return nil, errors.New("layout payload is required")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("layout name is required")
	}
	if req.FloorNumber <= 0 {
		return nil, errors.New("floor_number must be greater than 0")
	}
	if req.CanvasWidth <= 0 || req.CanvasHeight <= 0 {
		return nil, errors.New("canvas_width and canvas_height must be greater than 0")
	}

	layout := &domain.TableLayout{
		TenantID:        tenantID,
		RestaurantID:    restaurantID,
		Name:            name,
		FloorNumber:     req.FloorNumber,
		CanvasWidth:     req.CanvasWidth,
		CanvasHeight:    req.CanvasHeight,
		BackgroundImage: strings.TrimSpace(req.BackgroundImage),
		IsDefault:       req.IsDefault,
	}

	return uc.tableRepo.CreateLayout(layout)
}

// GetLayouts returns all layouts for a restaurant.
func (uc *TableUseCase) GetLayouts(tenantID, restaurantID int64) ([]domain.TableLayout, error) {
	if uc.tableRepo == nil {
		return nil, errors.New("table repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}
	return uc.tableRepo.GetLayouts(tenantID, restaurantID)
}

// GetTablePublicInfo returns public table metadata used by QR landing.
func (uc *TableUseCase) GetTablePublicInfo(tableID int64) (*domain.TablePublicInfo, error) {
	if uc.tableRepo == nil {
		return nil, errors.New("table repository is not configured")
	}
	if tableID <= 0 {
		return nil, errors.New("invalid table id")
	}
	return uc.tableRepo.GetTablePublicInfo(tableID)
}

func canTransitionTableStatus(oldStatus, newStatus domain.TableStatus) bool {
	if oldStatus == newStatus {
		return true
	}

	validTransitions := map[domain.TableStatus][]domain.TableStatus{
		domain.TableStatusAvailable: {
			domain.TableStatusOccupied,
			domain.TableStatusReserved,
			domain.TableStatusMaintenance,
		},
		domain.TableStatusOccupied: {
			domain.TableStatusAvailable,
			domain.TableStatusReserved,
			domain.TableStatusMaintenance,
		},
		domain.TableStatusReserved: {
			domain.TableStatusAvailable,
			domain.TableStatusOccupied,
			domain.TableStatusMaintenance,
		},
		domain.TableStatusMaintenance: {
			domain.TableStatusAvailable,
		},
	}

	candidates, exists := validTransitions[oldStatus]
	if !exists {
		return false
	}

	for _, status := range candidates {
		if status == newStatus {
			return true
		}
	}

	return false
}

func sanitizeTableText(input string, maxLen int) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return ""
	}

	cleaned := tableTextSanitizer.ReplaceAllString(trimmed, "")
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	if maxLen <= 0 {
		return cleaned
	}

	runes := []rune(cleaned)
	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}

	return cleaned
}
