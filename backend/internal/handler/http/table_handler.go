package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/usecase"
)

// TableHandler handles table-management API requests.
type TableHandler struct {
	tableUC *usecase.TableUseCase
}

// NewTableHandler creates a new table handler.
func NewTableHandler(tableUC *usecase.TableUseCase) *TableHandler {
	return &TableHandler{tableUC: tableUC}
}

// CreateTable creates a new table.
// POST /api/v1/tables
func (h *TableHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	tenantID, restaurantID, ok := tableTenantContext(r)
	if !ok {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	var req domain.CreateTableRequest
	if err := decodeTableBody(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	table, err := h.tableUC.CreateTable(tenantID, restaurantID, &req)
	if err != nil {
		respondTableError(w, err, "Failed to create table")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    table,
	})
}

// ListTables lists tables with optional status filter.
// GET /api/v1/tables
func (h *TableHandler) ListTables(w http.ResponseWriter, r *http.Request) {
	tenantID, restaurantID, ok := tableTenantContext(r)
	if !ok {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	status := strings.TrimSpace(r.URL.Query().Get("status"))
	tables, err := h.tableUC.ListTables(tenantID, restaurantID, status)
	if err != nil {
		respondTableError(w, err, "Failed to list tables")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    tables,
	})
}

// GetTable returns one table.
// GET /api/v1/tables/{id}
func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	tenantID, restaurantID, ok := tableTenantContext(r)
	if !ok {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	tableID, err := parseTableIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	table, err := h.tableUC.GetTable(tenantID, restaurantID, tableID)
	if err != nil {
		respondTableError(w, err, "Failed to get table")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    table,
	})
}

// UpdateTable updates one table.
// PUT /api/v1/tables/{id}
func (h *TableHandler) UpdateTable(w http.ResponseWriter, r *http.Request) {
	tenantID, restaurantID, ok := tableTenantContext(r)
	if !ok {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	tableID, err := parseTableIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	var req domain.UpdateTableRequest
	if err := decodeTableBody(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.tableUC.UpdateTable(tenantID, restaurantID, tableID, &req); err != nil {
		respondTableError(w, err, "Failed to update table")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// DeleteTable deletes one table.
// DELETE /api/v1/tables/{id}
func (h *TableHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	tenantID, restaurantID, ok := tableTenantContext(r)
	if !ok {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	tableID, err := parseTableIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	if err := h.tableUC.DeleteTable(tenantID, restaurantID, tableID); err != nil {
		respondTableError(w, err, "Failed to delete table")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// UpdateStatus updates one table status.
// PUT /api/v1/tables/{id}/status
func (h *TableHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, restaurantID, ok := tableTenantContext(r)
	if !ok {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	tableID, err := parseTableIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	var req domain.UpdateTableStatusRequest
	if err := decodeTableBody(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.tableUC.UpdateTableStatus(tenantID, restaurantID, tableID, req.Status); err != nil {
		respondTableError(w, err, "Failed to update table status")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// UpdatePositions bulk updates table positions.
// PUT /api/v1/tables/positions
func (h *TableHandler) UpdatePositions(w http.ResponseWriter, r *http.Request) {
	tenantID, restaurantID, ok := tableTenantContext(r)
	if !ok {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	var req domain.UpdateTablePositionsRequest
	if err := decodeTableBody(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.tableUC.UpdateTablePositions(tenantID, restaurantID, req.Positions); err != nil {
		respondTableError(w, err, "Failed to update table positions")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// GetTableStats returns aggregate table stats.
// GET /api/v1/tables/stats
func (h *TableHandler) GetTableStats(w http.ResponseWriter, r *http.Request) {
	tenantID, restaurantID, ok := tableTenantContext(r)
	if !ok {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	stats, err := h.tableUC.GetTableStats(tenantID, restaurantID)
	if err != nil {
		respondTableError(w, err, "Failed to fetch table stats")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// GetLayouts returns floor-plan layouts.
// GET /api/v1/tables/layouts
func (h *TableHandler) GetLayouts(w http.ResponseWriter, r *http.Request) {
	tenantID, restaurantID, ok := tableTenantContext(r)
	if !ok {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	layouts, err := h.tableUC.GetLayouts(tenantID, restaurantID)
	if err != nil {
		respondTableError(w, err, "Failed to fetch table layouts")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    layouts,
	})
}

// CreateLayout creates a table layout.
// POST /api/v1/tables/layouts
func (h *TableHandler) CreateLayout(w http.ResponseWriter, r *http.Request) {
	tenantID, restaurantID, ok := tableTenantContext(r)
	if !ok {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	var req domain.CreateLayoutRequest
	if err := decodeTableBody(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	layout, err := h.tableUC.CreateLayout(tenantID, restaurantID, &req)
	if err != nil {
		respondTableError(w, err, "Failed to create table layout")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    layout,
	})
}

// GetPublicTableInfo returns public table info by table ID.
// GET /api/v1/public/table/{id}/info
func (h *TableHandler) GetPublicTableInfo(w http.ResponseWriter, r *http.Request) {
	tableID, err := parseTableIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	info, err := h.tableUC.GetTablePublicInfo(tableID)
	if err != nil {
		respondTableError(w, err, "Failed to fetch table info")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    info,
	})
}

func tableTenantContext(r *http.Request) (int64, int64, bool) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		return 0, 0, false
	}
	return tenantID, restaurantID, true
}

func parseTableIDParam(r *http.Request) (int64, error) {
	tableID, err := strconv.ParseInt(strings.TrimSpace(r.PathValue("id")), 10, 64)
	if err != nil || tableID <= 0 {
		return 0, errors.New("invalid table id")
	}
	return tableID, nil
}

func decodeTableBody(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func respondTableError(w http.ResponseWriter, err error, fallbackMessage string) {
	if err == nil {
		respondError(w, http.StatusInternalServerError, fallbackMessage)
		return
	}

	msg := err.Error()
	lowerMsg := strings.ToLower(msg)

	switch {
	case strings.Contains(lowerMsg, "not found"):
		respondError(w, http.StatusNotFound, msg)
	case strings.Contains(lowerMsg, "already exists"), strings.Contains(lowerMsg, "duplicate"):
		respondError(w, http.StatusConflict, msg)
	case strings.Contains(lowerMsg, "active orders"):
		respondError(w, http.StatusConflict, msg)
	case strings.Contains(lowerMsg, "database is not initialized"),
		strings.Contains(lowerMsg, "not configured"):
		respondError(w, http.StatusServiceUnavailable, msg)
	case strings.Contains(lowerMsg, "invalid"),
		strings.Contains(lowerMsg, "required"),
		strings.Contains(lowerMsg, "cannot transition"),
		strings.Contains(lowerMsg, "must be"):
		respondError(w, http.StatusBadRequest, msg)
	default:
		respondError(w, http.StatusInternalServerError, fallbackMessage)
	}
}
