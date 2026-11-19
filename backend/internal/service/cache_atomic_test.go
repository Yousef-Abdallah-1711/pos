package service

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestInMemoryCacheSetIfAbsentIsAtomic(t *testing.T) {
	cache := NewInMemoryCache()
	var wait sync.WaitGroup
	results := make(chan bool, 32)
	for range 32 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			stored, err := cache.SetIfAbsent(context.Background(), "proof", "binding", time.Minute)
			if err != nil {
				t.Errorf("SetIfAbsent: %v", err)
				return
			}
			results <- stored
		}()
	}
	wait.Wait()
	close(results)
	created := 0
	for stored := range results {
		if stored {
			created++
		}
	}
	if created != 1 {
		t.Fatalf("SetIfAbsent successes = %d, want exactly one", created)
	}
}

func TestInMemoryRateLimitStoreIsBoundedAndExpires(t *testing.T) {
	store := NewInMemoryRateLimitStore(1)
	if _, _, err := store.IncrementWindow(context.Background(), "one", 20*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.IncrementWindow(context.Background(), "two", time.Minute); err == nil {
		t.Fatal("full store must reject a new key")
	}
	time.Sleep(25 * time.Millisecond)
	if _, _, err := store.IncrementWindow(context.Background(), "two", time.Minute); err != nil {
		t.Fatalf("expired keys should release bounded capacity: %v", err)
	}
}

func TestRedisAtomicAuthStoresAcrossClients(t *testing.T) {
	redisURL := os.Getenv("POOS_SESSION_REDIS_URL")
	if redisURL == "" {
		t.Skip("set POOS_SESSION_REDIS_URL to verify shared atomic auth stores")
	}
	first, err := NewRedisCache(redisURL)
	if err != nil {
		t.Fatalf("connect first Redis client: %v", err)
	}
	defer first.Close()
	second, err := NewRedisCache(redisURL)
	if err != nil {
		t.Fatalf("connect second Redis client: %v", err)
	}
	defer second.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	keyBase := fmt.Sprintf("test:auth-atomic:%d", time.Now().UnixNano())
	proofKey, counterKey := keyBase+":proof", keyBase+":counter"
	defer func() {
		_ = first.Delete(context.Background(), proofKey)
		_ = first.Delete(context.Background(), counterKey)
	}()

	var wait sync.WaitGroup
	proofResults := make(chan bool, 2)
	for _, cache := range []*RedisCache{first, second} {
		wait.Add(1)
		go func(cache *RedisCache) {
			defer wait.Done()
			stored, storeErr := cache.SetIfAbsent(ctx, proofKey, "tenant-bound-proof", time.Minute)
			if storeErr != nil {
				t.Errorf("SetIfAbsent: %v", storeErr)
				return
			}
			proofResults <- stored
		}(cache)
	}
	wait.Wait()
	close(proofResults)
	created := 0
	for stored := range proofResults {
		if stored {
			created++
		}
	}
	if created != 1 {
		t.Fatalf("cross-client SetIfAbsent successes = %d, want exactly one", created)
	}

	const increments = 24
	var countWait sync.WaitGroup
	errors := make(chan error, increments)
	for i := range increments {
		cache := first
		if i%2 == 1 {
			cache = second
		}
		countWait.Add(1)
		go func(cache *RedisCache) {
			defer countWait.Done()
			_, _, incrementErr := cache.IncrementWindow(ctx, counterKey, time.Minute)
			errors <- incrementErr
		}(cache)
	}
	countWait.Wait()
	close(errors)
	for incrementErr := range errors {
		if incrementErr != nil {
			t.Fatalf("IncrementWindow: %v", incrementErr)
		}
	}
	count, err := first.client.Get(ctx, counterKey).Int64()
	if err != nil {
		t.Fatalf("read shared rate-limit counter: %v", err)
	}
	if count != increments {
		t.Fatalf("shared counter = %d, want %d", count, increments)
	}
	ttl, err := first.client.PTTL(ctx, counterKey).Result()
	if err != nil {
		t.Fatalf("read shared rate-limit TTL: %v", err)
	}
	if ttl <= 0 || ttl > time.Minute {
		t.Fatalf("shared counter TTL = %s, want (0, 1m]", ttl)
	}
}
