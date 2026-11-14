package sharding

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// MigrationState tracks the state of shard migration
type MigrationState string

const (
	MigrationStatePending      MigrationState = "pending"      // Not started
	MigrationStateInitializing MigrationState = "initializing" // Setting up dual-write
	MigrationStateInProgress   MigrationState = "in_progress"  // Copying data
	MigrationStateVerifying    MigrationState = "verifying"    // Verifying data
	MigrationStateCutover      MigrationState = "cutover"      // Switching traffic
	MigrationStateCompleted    MigrationState = "completed"    // Migration done
	MigrationStateRolledBack   MigrationState = "rolled_back"  // Migration failed, rolled back
)

// MigrationProgress tracks migration progress metrics
type MigrationProgress struct {
	State            MigrationState
	TotalRecords     int64
	MigratedRecords  int64
	VerifiedRecords  int64
	FailedRecords    int64
	StartTime        time.Time
	EstimatedEndTime time.Time
	PercentComplete  float64
	ElapsedSeconds   int64
	EstimatedSeconds int64
}

// MigrationConfig holds migration configuration
type MigrationConfig struct {
	SourceShardID   int
	TargetShardID   int
	BatchSize       int
	DualWriteLeader int  // Which shard to lead for dual-write
	VerifyAll       bool // Verify all data after migration
	RollbackOnError bool // Rollback if errors occur
	Timeout         time.Duration
}

// ShardMigration manages live migration between shards
type ShardMigration struct {
	config     MigrationConfig
	sourceConn *sql.DB
	targetConn *sql.DB
	state      MigrationState
	progress   MigrationProgress
	dualWriter *DualWriteCoordinator
	mu         sync.RWMutex
	errorChan  chan error
	cancelFunc context.CancelFunc
	ctx        context.Context
}

// DualWriteCoordinator handles writing to both source and target shards
type DualWriteCoordinator struct {
	leaderShard int
	sourceConn  *sql.DB
	targetConn  *sql.DB
	mu          sync.RWMutex
	writeCount  int64
	errorCount  int64
}

// NewShardMigration creates a new migration instance
func NewShardMigration(config MigrationConfig, sourceConn, targetConn *sql.DB) *ShardMigration {
	ctx, cancel := context.WithCancel(context.Background())

	return &ShardMigration{
		config:     config,
		sourceConn: sourceConn,
		targetConn: targetConn,
		state:      MigrationStatePending,
		progress: MigrationProgress{
			State:     MigrationStatePending,
			StartTime: time.Now(),
		},
		dualWriter: &DualWriteCoordinator{
			leaderShard: config.DualWriteLeader,
			sourceConn:  sourceConn,
			targetConn:  targetConn,
		},
		errorChan:  make(chan error, 100),
		cancelFunc: cancel,
		ctx:        ctx,
	}
}

// Start begins the migration process
func (m *ShardMigration) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("[MIGRATION] Starting migration from shard %d to shard %d",
		m.config.SourceShardID, m.config.TargetShardID)

	// Phase 1: Initialize dual-write
	m.state = MigrationStateInitializing
	m.progress.State = MigrationStateInitializing
	log.Printf("[MIGRATION] Phase 1: Initializing dual-write coordinator")

	if err := m.initializeDualWrite(ctx); err != nil {
		log.Printf("[MIGRATION] ERROR: Failed to initialize dual-write: %v", err)
		m.state = MigrationStateRolledBack
		return err
	}

	log.Printf("[MIGRATION] Dual-write coordinator initialized successfully")

	// Phase 2: Count records to migrate
	log.Printf("[MIGRATION] Phase 2: Counting records in source shard")
	count, err := m.countRecords(ctx, m.sourceConn)
	if err != nil {
		log.Printf("[MIGRATION] ERROR: Failed to count records: %v", err)
		m.state = MigrationStateRolledBack
		return err
	}

	m.progress.TotalRecords = count
	m.progress.EstimatedSeconds = (count / int64(m.config.BatchSize)) + 60 // Add buffer
	log.Printf("[MIGRATION] Found %d records to migrate", count)

	// Phase 3: Migrate data
	m.state = MigrationStateInProgress
	m.progress.State = MigrationStateInProgress
	log.Printf("[MIGRATION] Phase 3: Starting data migration")

	if err := m.migrateData(ctx); err != nil {
		log.Printf("[MIGRATION] ERROR: Data migration failed: %v", err)
		if m.config.RollbackOnError {
			m.Rollback(ctx)
		}
		m.state = MigrationStateRolledBack
		return err
	}

	// Phase 4: Verify data
	if m.config.VerifyAll {
		m.state = MigrationStateVerifying
		m.progress.State = MigrationStateVerifying
		log.Printf("[MIGRATION] Phase 4: Verifying data in target shard")

		if err := m.verifyData(ctx); err != nil {
			log.Printf("[MIGRATION] ERROR: Data verification failed: %v", err)
			if m.config.RollbackOnError {
				m.Rollback(ctx)
			}
			m.state = MigrationStateRolledBack
			return err
		}
	}

	// Phase 5: Cutover
	m.state = MigrationStateCutover
	m.progress.State = MigrationStateCutover
	log.Printf("[MIGRATION] Phase 5: Starting traffic cutover")

	// Update router to use new shard (handled externally)
	// This is where traffic switches from source to target

	// Phase 6: Cleanup
	m.state = MigrationStateCompleted
	m.progress.State = MigrationStateCompleted
	log.Printf("[MIGRATION] Migration completed successfully")
	log.Printf("[MIGRATION] Migrated %d records in %.0f seconds",
		m.progress.MigratedRecords,
		time.Since(m.progress.StartTime).Seconds())

	return nil
}

// initializeDualWrite sets up dual-write to both shards
func (m *ShardMigration) initializeDualWrite(ctx context.Context) error {
	// Test both connections
	if err := m.sourceConn.PingContext(ctx); err != nil {
		return fmt.Errorf("source connection failed: %w", err)
	}

	if err := m.targetConn.PingContext(ctx); err != nil {
		return fmt.Errorf("target connection failed: %w", err)
	}

	log.Printf("[DUAL-WRITE] Coordinator initialized - leader: shard %d", m.dualWriter.leaderShard)
	return nil
}

// migrateData copies all records from source to target shard
func (m *ShardMigration) migrateData(ctx context.Context) error {
	log.Printf("[MIGRATE] Starting batch migration with batch size %d", m.config.BatchSize)

	var offset int64 = 0
	tables := []string{"orders", "order_items", "payments", "customers", "restaurants"}

	for _, table := range tables {
		log.Printf("[MIGRATE] Migrating table: %s", table)

		totalMigrated, err := m.migrateTable(ctx, table, offset)
		if err != nil {
			return fmt.Errorf("failed to migrate table %s: %w", table, err)
		}

		m.progress.MigratedRecords += totalMigrated
		m.progress.PercentComplete = float64(m.progress.MigratedRecords) / float64(m.progress.TotalRecords) * 100
		offset = m.progress.MigratedRecords

		log.Printf("[MIGRATE] Table %s: %d records migrated (%.1f%% complete)",
			table, totalMigrated, m.progress.PercentComplete)
	}

	return nil
}

// migrateTable migrates a specific table using batching
func (m *ShardMigration) migrateTable(ctx context.Context, table string, offset int64) (int64, error) {
	var totalMigrated int64

	for {
		select {
		case <-ctx.Done():
			return totalMigrated, ctx.Err()
		default:
		}

		// Query batch from source
		query := fmt.Sprintf(`
			SELECT * FROM %s
			ORDER BY id
			LIMIT %d OFFSET %d
		`, table, m.config.BatchSize, offset)

		rows, err := m.sourceConn.QueryContext(ctx, query)
		if err != nil {
			return totalMigrated, err
		}
		defer rows.Close()

		// Get column information
		columns, err := rows.Columns()
		if err != nil {
			return totalMigrated, err
		}

		// Read rows
		var rowsRead int
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for rows.Next() {
			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				m.progress.FailedRecords++
				log.Printf("[MIGRATE] Warning: Failed to scan row: %v", err)
				continue
			}

			// Insert into target
			if err := m.insertRow(ctx, table, columns, values); err != nil {
				m.progress.FailedRecords++
				log.Printf("[MIGRATE] Warning: Failed to insert row: %v", err)
				continue
			}

			rowsRead++
			totalMigrated++
		}

		if err := rows.Err(); err != nil {
			return totalMigrated, err
		}

		// If fewer rows than batch size, we're done with this table
		if rowsRead < m.config.BatchSize {
			break
		}

		offset += int64(rowsRead)
	}

	return totalMigrated, nil
}

// insertRow inserts a row into target shard
func (m *ShardMigration) insertRow(ctx context.Context, table string, columns []string, values []interface{}) error {
	// Build INSERT statement
	placeholders := ""
	for i := range columns {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (%s)
		ON CONFLICT DO NOTHING
	`, table, fmt.Sprintf("%s", columns), placeholders)

	_, err := m.targetConn.ExecContext(ctx, query, values...)
	return err
}

// verifyData verifies all data was migrated correctly
func (m *ShardMigration) verifyData(ctx context.Context) error {
	log.Printf("[VERIFY] Starting data verification")

	tables := []string{"orders", "order_items", "payments", "customers", "restaurants"}
	var totalVerified int64

	for _, table := range tables {
		sourceCount, err := m.countRecordsInTable(ctx, m.sourceConn, table)
		if err != nil {
			return fmt.Errorf("failed to count source records in %s: %w", table, err)
		}

		targetCount, err := m.countRecordsInTable(ctx, m.targetConn, table)
		if err != nil {
			return fmt.Errorf("failed to count target records in %s: %w", table, err)
		}

		if sourceCount != targetCount {
			return fmt.Errorf("table %s: source has %d records, target has %d records",
				table, sourceCount, targetCount)
		}

		totalVerified += sourceCount
		log.Printf("[VERIFY] Table %s: %d records verified ✓", table, sourceCount)
	}

	m.progress.VerifiedRecords = totalVerified
	log.Printf("[VERIFY] All %d records verified successfully", totalVerified)

	return nil
}

// countRecords returns total record count in a shard
func (m *ShardMigration) countRecords(ctx context.Context, conn *sql.DB) (int64, error) {
	var total int64
	tables := []string{"orders", "order_items", "payments", "customers", "restaurants"}

	for _, table := range tables {
		count, err := m.countRecordsInTable(ctx, conn, table)
		if err != nil {
			return 0, err
		}
		total += count
	}

	return total, nil
}

// countRecordsInTable counts records in a specific table
func (m *ShardMigration) countRecordsInTable(ctx context.Context, conn *sql.DB, table string) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	err := conn.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// Rollback reverts the migration and cleans up
func (m *ShardMigration) Rollback(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("[MIGRATION] Rolling back migration from shard %d to shard %d",
		m.config.SourceShardID, m.config.TargetShardID)

	// Delete all migrated data from target
	tables := []string{"order_items", "orders", "payments", "customers", "restaurants"}
	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s", table)
		if _, err := m.targetConn.ExecContext(ctx, query); err != nil {
			log.Printf("[MIGRATION] Warning: Failed to delete from %s: %v", table, err)
		}
	}

	m.state = MigrationStateRolledBack
	m.progress.State = MigrationStateRolledBack
	log.Printf("[MIGRATION] Rollback completed")

	return nil
}

// GetProgress returns current migration progress
func (m *ShardMigration) GetProgress() MigrationProgress {
	m.mu.RLock()
	defer m.mu.RUnlock()

	progress := m.progress
	progress.ElapsedSeconds = int64(time.Since(progress.StartTime).Seconds())

	return progress
}

// GetState returns current migration state
func (m *ShardMigration) GetState() MigrationState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.state
}

// Cancel cancels the migration process
func (m *ShardMigration) Cancel() {
	m.cancelFunc()
	m.state = MigrationStateRolledBack
	log.Printf("[MIGRATION] Migration cancelled by user")
}

// WriteToBothShards writes data to both source and target shards (during dual-write phase)
func (dw *DualWriteCoordinator) WriteData(ctx context.Context, table string, columns []string, values []interface{}) error {
	dw.mu.Lock()
	defer dw.mu.Unlock()

	// Build INSERT statement
	placeholders := ""
	for i := range columns {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (%s)
	`, table, fmt.Sprintf("%s", columns), placeholders)

	// Write to leader first
	var leaderConn, followerConn *sql.DB
	if dw.leaderShard == 0 {
		leaderConn = dw.sourceConn
		followerConn = dw.targetConn
	} else {
		leaderConn = dw.targetConn
		followerConn = dw.sourceConn
	}

	// Write to leader (blocking)
	if _, err := leaderConn.ExecContext(ctx, query, values...); err != nil {
		dw.errorCount++
		return err
	}

	// Write to follower (non-blocking, fire and forget)
	go func() {
		if _, err := followerConn.ExecContext(context.Background(), query, values...); err != nil {
			dw.errorCount++
			log.Printf("[DUAL-WRITE] Warning: Failed to write to follower: %v", err)
		}
	}()

	dw.writeCount++
	return nil
}

// GetDualWriteStats returns dual-write statistics
func (dw *DualWriteCoordinator) GetDualWriteStats() (writeCount, errorCount int64) {
	dw.mu.RLock()
	defer dw.mu.RUnlock()

	return dw.writeCount, dw.errorCount
}
