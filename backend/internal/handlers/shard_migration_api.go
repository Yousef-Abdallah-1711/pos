//go:build legacy
// +build legacy

package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/sharding"
)

// ShardMigrationHandler manages HTTP API for shard migrations
type ShardMigrationHandler struct {
	shardRouter *sharding.ShardRouter
	migrations  map[string]*sharding.ShardMigration
}

// NewShardMigrationHandler creates a new migration handler
func NewShardMigrationHandler(router *sharding.ShardRouter) *ShardMigrationHandler {
	return &ShardMigrationHandler{
		shardRouter: router,
		migrations:  make(map[string]*sharding.ShardMigration),
	}
}

// StartMigrationRequest is the request body for starting migration
type StartMigrationRequest struct {
	SourceShardID int `json:"source_shard_id"`
	TargetShardID int `json:"target_shard_id"`
	BatchSize     int `json:"batch_size,omitempty"`
	VerifyAll     bool `json:"verify_all,omitempty"`
	RollbackOnError bool `json:"rollback_on_error,omitempty"`
}

// MigrationResponse is the response for migration status
type MigrationResponse struct {
	MigrationID  string  `json:"migration_id"`
	State        string  `json:"state"`
	Progress     float64 `json:"progress"`
	TotalRecords int64   `json:"total_records"`
	Migrated     int64   `json:"migrated"`
	Verified     int64   `json:"verified"`
	Failed       int64   `json:"failed"`
	Elapsed      int64   `json:"elapsed_seconds"`
	Message      string  `json:"message"`
}

// StartMigration starts a new shard migration
// POST /api/v1/admin/migrations
func (h *ShardMigrationHandler) StartMigration(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MIGRATION API] StartMigration called")

	var req StartMigrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("[MIGRATION API] Request: source=%d, target=%d", req.SourceShardID, req.TargetShardID)

	// Validate request
	if req.SourceShardID == req.TargetShardID {
		http.Error(w, "Source and target shards must be different", http.StatusBadRequest)
		return
	}

	// Get source and target connections
	shards := h.shardRouter.GetAllShards()

	var sourceConn, targetConn *sql.DB
	var sourceFound, targetFound bool

	for i, shard := range shards {
		if shard.ID == req.SourceShardID {
			conn, err := h.shardRouter.GetConnection(fmt.Sprintf("%d", shard.ID))
			if err != nil {
				http.Error(w, fmt.Sprintf("Cannot connect to source shard: %v", err), http.StatusInternalServerError)
				return
			}
			sourceConn = conn
			sourceFound = true
		}
		if shard.ID == req.TargetShardID {
			conn, err := h.shardRouter.GetConnection(fmt.Sprintf("%d", shard.ID))
			if err != nil {
				http.Error(w, fmt.Sprintf("Cannot connect to target shard: %v", err), http.StatusInternalServerError)
				return
			}
			targetConn = conn
			targetFound = true
		}
	}

	if !sourceFound || !targetFound {
		http.Error(w, "Source or target shard not found", http.StatusNotFound)
		return
	}

	// Set defaults
	batchSize := req.BatchSize
	if batchSize == 0 {
		batchSize = 1000
	}

	verifyAll := req.VerifyAll
	if !verifyAll {
		verifyAll = true // Default to verify
	}

	rollbackOnError := req.RollbackOnError
	if !rollbackOnError {
		rollbackOnError = true // Default to rollback
	}

	// Create migration
	config := sharding.MigrationConfig{
		SourceShardID:   req.SourceShardID,
		TargetShardID:   req.TargetShardID,
		BatchSize:       batchSize,
		DualWriteLeader: req.SourceShardID,
		VerifyAll:       verifyAll,
		RollbackOnError: rollbackOnError,
	}

	migration := sharding.NewShardMigration(config, sourceConn, targetConn)
	migrationID := fmt.Sprintf("migration_%d_to_%d_%d", req.SourceShardID, req.TargetShardID, time.Now().Unix())

	// Store migration
	h.migrations[migrationID] = migration

	// Start migration in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 24*time.Hour)
		defer cancel()

		log.Printf("[MIGRATION API] Starting background migration: %s", migrationID)
		if err := migration.Start(ctx); err != nil {
			log.Printf("[MIGRATION API] ERROR: Migration %s failed: %v", migrationID, err)
		} else {
			log.Printf("[MIGRATION API] Migration %s completed successfully", migrationID)
		}
	}()

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"migration_id": migrationID,
		"state":        "initializing",
		"message":      "Migration started in background",
	})

	log.Printf("[MIGRATION API] Migration %s created", migrationID)
}

// GetMigrationStatus returns the status of a migration
// GET /api/v1/admin/migrations/{migration_id}
func (h *ShardMigrationHandler) GetMigrationStatus(w http.ResponseWriter, r *http.Request) {
	migrationID := r.PathValue("migration_id")
	log.Printf("[MIGRATION API] GetMigrationStatus: %s", migrationID)

	migration, exists := h.migrations[migrationID]
	if !exists {
		http.Error(w, "Migration not found", http.StatusNotFound)
		return
	}

	progress := migration.GetProgress()

	response := MigrationResponse{
		MigrationID:  migrationID,
		State:        string(progress.State),
		Progress:     progress.PercentComplete,
		TotalRecords: progress.TotalRecords,
		Migrated:     progress.MigratedRecords,
		Verified:     progress.VerifiedRecords,
		Failed:       progress.FailedRecords,
		Elapsed:      progress.ElapsedSeconds,
		Message:      fmt.Sprintf("Migration in progress - %.1f%% complete", progress.PercentComplete),
	}

	// Set custom message based on state
	switch progress.State {
	case "completed":
		response.Message = fmt.Sprintf("Migration completed in %d seconds", progress.ElapsedSeconds)
	case "rolled_back":
		response.Message = "Migration rolled back due to error"
	case "verifying":
		response.Message = fmt.Sprintf("Verifying data - %d/%d verified", progress.VerifiedRecords, progress.TotalRecords)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CancelMigration cancels a running migration
// POST /api/v1/admin/migrations/{migration_id}/cancel
func (h *ShardMigrationHandler) CancelMigration(w http.ResponseWriter, r *http.Request) {
	migrationID := r.PathValue("migration_id")
	log.Printf("[MIGRATION API] CancelMigration: %s", migrationID)

	migration, exists := h.migrations[migrationID]
	if !exists {
		http.Error(w, "Migration not found", http.StatusNotFound)
		return
	}

	state := migration.GetState()
	if state == "completed" || state == "rolled_back" {
		http.Error(w, "Cannot cancel migration that is already completed", http.StatusBadRequest)
		return
	}

	migration.Cancel()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"migration_id": migrationID,
		"state":        "cancelled",
		"message":      "Migration cancellation initiated",
	})

	log.Printf("[MIGRATION API] Migration %s cancelled", migrationID)
}

// ListMigrations returns all active migrations
// GET /api/v1/admin/migrations
func (h *ShardMigrationHandler) ListMigrations(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MIGRATION API] ListMigrations called")

	var migrations []map[string]interface{}

	for id, migration := range h.migrations {
		progress := migration.GetProgress()
		migrations = append(migrations, map[string]interface{}{
			"migration_id":  id,
			"state":         progress.State,
			"progress":      progress.PercentComplete,
			"migrated":      progress.MigratedRecords,
			"total":         progress.TotalRecords,
			"elapsed":       progress.ElapsedSeconds,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":       len(migrations),
		"migrations":  migrations,
	})
}

// GetShardStats returns statistics for a shard
// GET /api/v1/admin/shards/{shard_id}/stats
func (h *ShardMigrationHandler) GetShardStats(w http.ResponseWriter, r *http.Request) {
	shardIDStr := r.PathValue("shard_id")
	shardID, err := strconv.Atoi(shardIDStr)
	if err != nil {
		http.Error(w, "Invalid shard ID", http.StatusBadRequest)
		return
	}

	log.Printf("[MIGRATION API] GetShardStats: shard=%d", shardID)

	// Get health check result
	healthResults := h.shardRouter.HealthCheckAllShards()
	shardHealth, exists := healthResults[shardID]

	stats := map[string]interface{}{
		"shard_id": shardID,
		"healthy":  exists && shardHealth == nil,
	}

	if !exists {
		http.Error(w, "Shard not found", http.StatusNotFound)
		return
	}

	if shardHealth != nil {
		stats["error"] = shardHealth.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ValidateShardData validates data integrity in a shard
// POST /api/v1/admin/shards/{shard_id}/validate
func (h *ShardMigrationHandler) ValidateShardData(w http.ResponseWriter, r *http.Request) {
	shardIDStr := r.PathValue("shard_id")
	shardID, err := strconv.Atoi(shardIDStr)
	if err != nil {
		http.Error(w, "Invalid shard ID", http.StatusBadRequest)
		return
	}

	log.Printf("[MIGRATION API] ValidateShardData: shard=%d", shardID)

	// TODO: Implement data validation logic
	// This would check for data consistency, missing records, etc.

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"shard_id":     shardID,
		"valid":        true,
		"record_count": 0,
		"message":      "Data validation pending implementation",
	})
}
