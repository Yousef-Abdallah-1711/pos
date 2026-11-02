package sharding

import (
	"crypto/md5"
	"database/sql"
	"encoding/binary"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

const (
	// MULTIPLIER used in jump hash algorithm
	MULTIPLIER = 1111111111111111111

	// INCREMENT used in jump hash algorithm
	INCREMENT = 11
)

// ShardInfo contains information about a database shard
type ShardInfo struct {
	ID       int
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Replica1 *ShardInfo // First replica
	Replica2 *ShardInfo // Second replica
}

// ShardRouter handles routing of restaurant data to correct shard
type ShardRouter struct {
	shards      []ShardInfo
	numShards   int
	connections map[int]*sql.DB
	mu          sync.RWMutex
}

// NewShardRouter creates a new shard router
func NewShardRouter(shards []ShardInfo) *ShardRouter {
	return &ShardRouter{
		shards:      shards,
		numShards:   len(shards),
		connections: make(map[int]*sql.DB),
	}
}

// GetShard returns the shard info for a given restaurant ID using jump hash
func (sr *ShardRouter) GetShard(restaurantID string) *ShardInfo {
	shardNumber := sr.getShardNumber(restaurantID)
	return &sr.shards[shardNumber]
}

// GetShardNumber returns the shard number (0-indexed) for a restaurant ID
func (sr *ShardRouter) getShardNumber(restaurantID string) int {
	return sr.jumpHash(restaurantID) % sr.numShards
}

// jumpHash maps a key to a shard in [0, numShards) with the Lamping-Veach jump consistent hash
// (https://arxiv.org/abs/1406.2294). The result depends on the key, is stable for a given
// (key, numShards), and when numShards grows only ~1/n of the keys move (always to the new shard).
func (sr *ShardRouter) jumpHash(key string) int {
	if sr.numShards <= 1 {
		return 0
	}

	// Spread the key over 64 bits (MD5 is used only as a uniform mixer, not for security).
	sum := md5.Sum([]byte(key))
	k := binary.BigEndian.Uint64(sum[:8])

	var b, j int64 = -1, 0
	for j < int64(sr.numShards) {
		b = j
		k = k*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((k>>33)+1)))
	}
	return int(b)
}

// GetDSN returns the database connection string for the shard
func (sr *ShardRouter) GetDSN(restaurantID string) string {
	shard := sr.GetShard(restaurantID)
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		shard.User,
		shard.Password,
		shard.Host,
		shard.Port,
		shard.Database,
	)
}

// GetConnection returns a database connection for the shard
func (sr *ShardRouter) GetConnection(restaurantID string) (*sql.DB, error) {
	shardNumber := sr.getShardNumber(restaurantID)

	sr.mu.RLock()
	if conn, exists := sr.connections[shardNumber]; exists {
		sr.mu.RUnlock()
		return conn, nil
	}
	sr.mu.RUnlock()

	// Create new connection
	dsn := sr.GetDSN(restaurantID)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to shard %d: %w", shardNumber, err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping shard %d: %w", shardNumber, err)
	}

	// Cache connection
	sr.mu.Lock()
	sr.connections[shardNumber] = db
	sr.mu.Unlock()

	return db, nil
}

// GetAllShards returns all shards (for operations like backup, maintenance)
func (sr *ShardRouter) GetAllShards() []ShardInfo {
	return sr.shards
}

// GetAllConnections returns connections to all shards
func (sr *ShardRouter) GetAllConnections() (map[int]*sql.DB, error) {
	connections := make(map[int]*sql.DB)

	for i, shard := range sr.shards {
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			shard.User,
			shard.Password,
			shard.Host,
			shard.Port,
			shard.Database,
		)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open connection to shard %d: %w", i, err)
		}

		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping shard %d: %w", i, err)
		}

		connections[i] = db
	}

	return connections, nil
}

// GetShardNumberForRestaurant returns the shard number for a restaurant
// (public method for routing logic)
func (sr *ShardRouter) GetShardNumberForRestaurant(restaurantID string) int {
	return sr.getShardNumber(restaurantID)
}

// Close closes all cached connections
func (sr *ShardRouter) Close() error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	for _, conn := range sr.connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}

	sr.connections = make(map[int]*sql.DB)
	return nil
}

// HealthCheck checks the health of a specific shard
func (sr *ShardRouter) HealthCheck(restaurantID string) error {
	conn, err := sr.GetConnection(restaurantID)
	if err != nil {
		return err
	}

	return conn.Ping()
}

// HealthCheckAllShards checks health of all shards
func (sr *ShardRouter) HealthCheckAllShards() map[int]error {
	results := make(map[int]error)

	for i, shard := range sr.shards {
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			shard.User,
			shard.Password,
			shard.Host,
			shard.Port,
			shard.Database,
		)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			results[i] = err
			continue
		}

		if err := db.Ping(); err != nil {
			results[i] = err
		} else {
			results[i] = nil
		}

		db.Close()
	}

	return results
}
