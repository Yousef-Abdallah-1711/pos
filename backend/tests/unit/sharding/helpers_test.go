//go:build experimental
// +build experimental

package sharding_test

import (
	"testing"

	"pos-saas/internal/sharding"
)

// createTestShardRouter builds a router over four fake shard descriptors. No connections are
// opened until a shard is used, so tests that only exercise routing and bookkeeping need no database.
// (health_test.go referenced this helper but it was never defined in this package.)
func createTestShardRouter(t *testing.T) *sharding.ShardRouter {
	t.Helper()
	shards := make([]sharding.ShardInfo, 0, 4)
	for i := 0; i < 4; i++ {
		shards = append(shards, sharding.ShardInfo{
			ID:       i,
			Host:     "shard" + string(rune('0'+i)) + ".example.invalid",
			Port:     5432,
			User:     "postgres",
			Password: "test-only",
			Database: "shard" + string(rune('0'+i)) + "_db",
		})
	}
	return sharding.NewShardRouter(shards)
}
