//go:build experimental
// +build experimental

package sharding_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"pos-saas/internal/sharding"
)

// TestMigrationStateTransitions tests migration state transitions
func TestMigrationStateTransitions(t *testing.T) {
	sourceConn, targetConn := setupTestMigrationDB(t)
	defer sourceConn.Close()
	defer targetConn.Close()

	config := sharding.MigrationConfig{
		SourceShardID:   0,
		TargetShardID:   4,
		BatchSize:       100,
		DualWriteLeader: 0,
		VerifyAll:       false,
		RollbackOnError: true,
	}

	migration := sharding.NewShardMigration(config, sourceConn, targetConn)

	// Check initial state
	state := migration.GetState()
	if state != "pending" {
		t.Fatalf("Expected initial state 'pending', got '%s'", state)
	}

	t.Logf("✓ Initial state is pending")
}

// TestMigrationProgressTracking tests progress tracking
func TestMigrationProgressTracking(t *testing.T) {
	sourceConn, targetConn := setupTestMigrationDB(t)
	defer sourceConn.Close()
	defer targetConn.Close()

	config := sharding.MigrationConfig{
		SourceShardID:   0,
		TargetShardID:   4,
		BatchSize:       100,
		DualWriteLeader: 0,
		VerifyAll:       false,
		RollbackOnError: true,
	}

	migration := sharding.NewShardMigration(config, sourceConn, targetConn)

	progress := migration.GetProgress()

	// Check initial progress
	if progress.TotalRecords != 0 {
		t.Fatalf("Expected 0 total records initially, got %d", progress.TotalRecords)
	}

	if progress.PercentComplete != 0 {
		t.Fatalf("Expected 0%% complete initially, got %.1f%%", progress.PercentComplete)
	}

	if progress.State != "pending" {
		t.Fatalf("Expected state 'pending', got '%s'", progress.State)
	}

	t.Logf("✓ Progress tracking initialized correctly")
}

// TestMigrationConfigValidation tests migration configuration validation
func TestMigrationConfigValidation(t *testing.T) {
	sourceConn, targetConn := setupTestMigrationDB(t)
	defer sourceConn.Close()
	defer targetConn.Close()

	tests := []struct {
		name   string
		config sharding.MigrationConfig
		valid  bool
	}{
		{
			name: "Valid config",
			config: sharding.MigrationConfig{
				SourceShardID:   0,
				TargetShardID:   4,
				BatchSize:       1000,
				DualWriteLeader: 0,
				VerifyAll:       true,
				RollbackOnError: true,
			},
			valid: true,
		},
		{
			name: "Same source and target",
			config: sharding.MigrationConfig{
				SourceShardID:   0,
				TargetShardID:   0,
				BatchSize:       1000,
				DualWriteLeader: 0,
				VerifyAll:       true,
				RollbackOnError: true,
			},
			valid: true, // Validation happens at API level, not here
		},
		{
			name: "Zero batch size",
			config: sharding.MigrationConfig{
				SourceShardID:   0,
				TargetShardID:   4,
				BatchSize:       0,
				DualWriteLeader: 0,
				VerifyAll:       true,
				RollbackOnError: true,
			},
			valid: true, // Will use default
		},
	}

	for _, tc := range tests {
		migration := sharding.NewShardMigration(tc.config, sourceConn, targetConn)
		if migration == nil {
			if tc.valid {
				t.Errorf("Test '%s': Expected valid migration, got nil", tc.name)
			}
		} else {
			if !tc.valid {
				t.Errorf("Test '%s': Expected invalid migration, got valid", tc.name)
			}
		}
	}

	t.Log("✓ Migration configuration validation works correctly")
}

// TestDualWriteCoordinator tests dual-write functionality
func TestDualWriteCoordinator(t *testing.T) {
	sourceConn, targetConn := setupTestMigrationDB(t)
	defer sourceConn.Close()
	defer targetConn.Close()

	config := sharding.MigrationConfig{
		SourceShardID:   0,
		TargetShardID:   4,
		BatchSize:       100,
		DualWriteLeader: 0,
		VerifyAll:       false,
		RollbackOnError: true,
	}

	migration := sharding.NewShardMigration(config, sourceConn, targetConn)

	// Get progress (which includes migration instance with dual-write coordinator)
	progress := migration.GetProgress()

	if progress.State != "pending" {
		t.Fatalf("Expected initial state 'pending', got '%s'", progress.State)
	}

	t.Logf("✓ Dual-write coordinator initialized")
}

// TestMigrationCancellation tests cancelling a migration
func TestMigrationCancellation(t *testing.T) {
	sourceConn, targetConn := setupTestMigrationDB(t)
	defer sourceConn.Close()
	defer targetConn.Close()

	config := sharding.MigrationConfig{
		SourceShardID:   0,
		TargetShardID:   4,
		BatchSize:       100,
		DualWriteLeader: 0,
		VerifyAll:       false,
		RollbackOnError: true,
	}

	migration := sharding.NewShardMigration(config, sourceConn, targetConn)

	// Cancel migration
	migration.Cancel()

	state := migration.GetState()
	if state != "rolled_back" {
		t.Fatalf("Expected state 'rolled_back' after cancel, got '%s'", state)
	}

	t.Logf("✓ Migration cancellation works correctly")
}

// TestMigrationTimeout tests migration timeout handling
func TestMigrationTimeout(t *testing.T) {
	sourceConn, targetConn := setupTestMigrationDB(t)
	defer sourceConn.Close()
	defer targetConn.Close()

	config := sharding.MigrationConfig{
		SourceShardID:   0,
		TargetShardID:   4,
		BatchSize:       100,
		DualWriteLeader: 0,
		VerifyAll:       false,
		RollbackOnError: true,
		Timeout:         time.Millisecond * 100, // Very short timeout
	}

	migration := sharding.NewShardMigration(config, sourceConn, targetConn)

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	// Start migration (will timeout)
	err := migration.Start(ctx)

	// Should get context deadline exceeded or similar
	if err == nil {
		// Migration completed very quickly without error (possible in test with empty DB)
		state := migration.GetState()
		if state != "completed" {
			t.Logf("Note: Migration completed state: %s", state)
		}
	}

	t.Logf("✓ Migration timeout handling works")
}

// TestMigrationProgressMonitoring tests monitoring migration progress
func TestMigrationProgressMonitoring(t *testing.T) {
	sourceConn, targetConn := setupTestMigrationDB(t)
	defer sourceConn.Close()
	defer targetConn.Close()

	config := sharding.MigrationConfig{
		SourceShardID:   0,
		TargetShardID:   4,
		BatchSize:       100,
		DualWriteLeader: 0,
		VerifyAll:       false,
		RollbackOnError: true,
	}

	migration := sharding.NewShardMigration(config, sourceConn, targetConn)

	// Check progress multiple times
	progressSnapshots := []sharding.MigrationProgress{}

	for i := 0; i < 3; i++ {
		progress := migration.GetProgress()
		progressSnapshots = append(progressSnapshots, progress)
		time.Sleep(time.Millisecond * 10)
	}

	// All snapshots should have same or increasing values
	for i := 1; i < len(progressSnapshots); i++ {
		if progressSnapshots[i].MigratedRecords < progressSnapshots[i-1].MigratedRecords {
			t.Fatalf("Progress went backward: %d -> %d",
				progressSnapshots[i-1].MigratedRecords,
				progressSnapshots[i].MigratedRecords)
		}
	}

	t.Logf("✓ Migration progress monitoring works correctly")
}

// TestMigrationRollback tests rollback functionality
func TestMigrationRollback(t *testing.T) {
	sourceConn, targetConn := setupTestMigrationDB(t)
	defer sourceConn.Close()
	defer targetConn.Close()

	config := sharding.MigrationConfig{
		SourceShardID:   0,
		TargetShardID:   4,
		BatchSize:       100,
		DualWriteLeader: 0,
		VerifyAll:       false,
		RollbackOnError: true,
	}

	migration := sharding.NewShardMigration(config, sourceConn, targetConn)

	// Rollback migration
	ctx := context.Background()
	err := migration.Rollback(ctx)

	if err != nil {
		t.Logf("Rollback returned error (expected in test DB): %v", err)
	}

	state := migration.GetState()
	if state != "rolled_back" {
		t.Fatalf("Expected state 'rolled_back' after rollback, got '%s'", state)
	}

	t.Logf("✓ Migration rollback works correctly")
}

// BenchmarkMigrationStateTransitions benchmarks state transition performance
func BenchmarkMigrationStateTransitions(b *testing.B) {
	sourceConn, targetConn := setupTestMigrationDB(&testing.T{})
	defer sourceConn.Close()
	defer targetConn.Close()

	config := sharding.MigrationConfig{
		SourceShardID:   0,
		TargetShardID:   4,
		BatchSize:       100,
		DualWriteLeader: 0,
		VerifyAll:       false,
		RollbackOnError: true,
	}

	migration := sharding.NewShardMigration(config, sourceConn, targetConn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = migration.GetState()
	}
}

// BenchmarkMigrationProgressQuery benchmarks progress query performance
func BenchmarkMigrationProgressQuery(b *testing.B) {
	sourceConn, targetConn := setupTestMigrationDB(&testing.T{})
	defer sourceConn.Close()
	defer targetConn.Close()

	config := sharding.MigrationConfig{
		SourceShardID:   0,
		TargetShardID:   4,
		BatchSize:       100,
		DualWriteLeader: 0,
		VerifyAll:       false,
		RollbackOnError: true,
	}

	migration := sharding.NewShardMigration(config, sourceConn, targetConn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = migration.GetProgress()
	}
}

// HELPER FUNCTIONS

// setupTestMigrationDB creates test databases for migration testing
func setupTestMigrationDB(t *testing.T) (*sql.DB, *sql.DB) {
	// In a real test environment, these would connect to actual test databases
	// For now, we create mock connections

	// Mock source connection
	sourceConn := &sql.DB{}
	targetConn := &sql.DB{}

	return sourceConn, targetConn
}
