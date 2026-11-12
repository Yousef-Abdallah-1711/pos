package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"pos-saas/internal/middleware"
)

// InventoryAdjustmentHandler applies a stock change and its audit movement atomically.
type InventoryAdjustmentHandler struct {
	db *sql.DB
}

func NewInventoryAdjustmentHandler(db *sql.DB) *InventoryAdjustmentHandler {
	return &InventoryAdjustmentHandler{db: db}
}

type inventoryAdjustmentRequest struct {
	ProductID      int64  `json:"product_id"`
	QuantityChange int64  `json:"quantity_change"`
	Reason         string `json:"reason"`
	Notes          string `json:"notes"`
}

func (h *InventoryAdjustmentHandler) Adjust(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	if h == nil || h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Inventory service is unavailable")
		return
	}

	var request inventoryAdjustmentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	request.Reason = strings.TrimSpace(request.Reason)
	if request.ProductID <= 0 || request.QuantityChange == 0 || !validInventoryAdjustmentReason(request.Reason) {
		respondError(w, http.StatusBadRequest, "A product, non-zero quantity change, and valid reason are required")
		return
	}

	ctx := r.Context()
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to start inventory adjustment")
		return
	}
	defer tx.Rollback()

	var currentStock int64
	err = tx.QueryRowContext(ctx, `
		SELECT quantity_in_stock
		FROM products
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3 AND track_inventory = TRUE
		FOR UPDATE
	`, request.ProductID, claims.TenantID, claims.RestaurantID).Scan(&currentStock)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load product stock")
		return
	}
	quantityAfter := currentStock + request.QuantityChange
	if quantityAfter < 0 {
		respondError(w, http.StatusConflict, "Inventory adjustment cannot make stock negative")
		return
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE products SET quantity_in_stock = $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, quantityAfter, request.ProductID, claims.TenantID, claims.RestaurantID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update product stock")
		return
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO inventory_movements (
			tenant_id, restaurant_id, product_id, movement_type, quantity, notes, created_by
		) VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7)
	`, claims.TenantID, claims.RestaurantID, request.ProductID, request.Reason, request.QuantityChange, strings.TrimSpace(request.Notes), claims.UserID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to record inventory movement")
		return
	}
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit inventory adjustment")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":        "Inventory adjusted successfully",
		"quantity_after": quantityAfter,
	})
}

func validInventoryAdjustmentReason(reason string) bool {
	switch reason {
	case "purchase", "sale", "adjustment", "waste", "return":
		return true
	default:
		return false
	}
}
