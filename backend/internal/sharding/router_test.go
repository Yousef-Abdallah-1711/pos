package sharding

import (
	"strconv"
	"testing"
)

func routerWithShards(n int) *ShardRouter {
	shards := make([]ShardInfo, n)
	for i := range shards {
		shards[i] = ShardInfo{ID: i, Host: "shard" + strconv.Itoa(i) + ".invalid", Port: 5432, Database: "db"}
	}
	return NewShardRouter(shards)
}

// DB-006: the previous implementation ignored the key and routed every restaurant to shard 0.
func TestJumpHash_DependsOnKeyAndStaysInRange(t *testing.T) {
	r := routerWithShards(4)
	seen := map[int]int{}
	for i := 0; i < 4000; i++ {
		n := r.GetShardNumberForRestaurant(strconv.Itoa(i))
		if n < 0 || n >= 4 {
			t.Fatalf("shard %d out of range for key %d", n, i)
		}
		seen[n]++
	}
	if len(seen) != 4 {
		t.Fatalf("keys must spread over all 4 shards, got %v", seen)
	}
	for shard, count := range seen {
		// uniform expectation is 1000 per shard; allow generous slack
		if count < 800 || count > 1200 {
			t.Errorf("shard %d got %d of 4000 keys (expected ~1000)", shard, count)
		}
	}
}

func TestJumpHash_IsDeterministic(t *testing.T) {
	a, b := routerWithShards(8), routerWithShards(8)
	for i := 0; i < 500; i++ {
		key := "restaurant-" + strconv.Itoa(i)
		if a.GetShardNumberForRestaurant(key) != b.GetShardNumberForRestaurant(key) {
			t.Fatalf("routing for %q must be stable", key)
		}
	}
}

// Growing the cluster must only move keys to the NEW shard, and only about 1/(n+1) of them.
func TestJumpHash_MinimalMovementWhenAddingAShard(t *testing.T) {
	before, after := routerWithShards(4), routerWithShards(5)
	moved := 0
	const total = 5000
	for i := 0; i < total; i++ {
		key := strconv.Itoa(i)
		from, to := before.GetShardNumberForRestaurant(key), after.GetShardNumberForRestaurant(key)
		if from != to {
			moved++
			if to != 4 {
				t.Fatalf("key %s moved from shard %d to existing shard %d; only the new shard may receive keys", key, from, to)
			}
		}
	}
	// expected 1/5 = 1000 moved
	if moved < 800 || moved > 1200 {
		t.Fatalf("expected ~%d keys to move, %d moved", total/5, moved)
	}
}

func TestJumpHash_SingleShardAlwaysZero(t *testing.T) {
	r := routerWithShards(1)
	for i := 0; i < 50; i++ {
		if n := r.GetShardNumberForRestaurant(strconv.Itoa(i)); n != 0 {
			t.Fatalf("single shard must always be 0, got %d", n)
		}
	}
}
