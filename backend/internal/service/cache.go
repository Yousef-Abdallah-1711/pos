package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// KVCache is a minimal key-value cache abstraction for short-lived backend caching.
type KVCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// AtomicSessionCache rotates an existing session key into a new credential in
// one store operation. It is optional because general-purpose caches do not
// need session-specific atomicity.
type AtomicSessionCache interface {
	Rotate(ctx context.Context, oldKey, newKey, value string, ttl time.Duration) error
}

// SlidingSessionCache refreshes a session TTL only if its value has not been
// replaced or revoked since it was read.
type SlidingSessionCache interface {
	Touch(ctx context.Context, key, expectedValue string, ttl time.Duration) error
}

type memoryCacheItem struct {
	value     string
	expiresAt time.Time
}

// InMemoryRateLimitStore is a bounded, process-local fallback for non-production
// environments. A full store fails closed until expired entries are removed.
type InMemoryRateLimitStore struct {
	mu      sync.Mutex
	items   map[string]memoryRateLimitItem
	maxKeys int
}

type memoryRateLimitItem struct {
	count     int64
	expiresAt time.Time
}

// NewInMemoryRateLimitStore creates a bounded process-local rate limit store.
func NewInMemoryRateLimitStore(maxKeys int) *InMemoryRateLimitStore {
	if maxKeys <= 0 {
		maxKeys = 10000
	}
	return &InMemoryRateLimitStore{items: make(map[string]memoryRateLimitItem), maxKeys: maxKeys}
}

func (s *InMemoryRateLimitStore) IncrementWindow(_ context.Context, key string, window time.Duration) (int64, time.Duration, error) {
	if window <= 0 {
		return 0, 0, fmt.Errorf("rate limit window must be positive")
	}
	if window.Milliseconds() < 1 {
		return 0, 0, fmt.Errorf("rate limit window must be at least one millisecond")
	}
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[key]
	if ok && now.Before(item.expiresAt) {
		item.count++
		s.items[key] = item
		return item.count, time.Until(item.expiresAt), nil
	}
	for expiredKey, expiredItem := range s.items {
		if !now.Before(expiredItem.expiresAt) {
			delete(s.items, expiredKey)
		}
	}
	if len(s.items) >= s.maxKeys {
		return 0, 0, fmt.Errorf("in-memory rate limit store capacity reached")
	}
	item = memoryRateLimitItem{count: 1, expiresAt: now.Add(window)}
	s.items[key] = item
	return item.count, time.Until(item.expiresAt), nil
}

// InMemoryCache is a process-local TTL cache fallback.
type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]memoryCacheItem
}

// NewInMemoryCache returns an in-memory cache implementation.
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		items: make(map[string]memoryCacheItem),
	}
}

func (c *InMemoryCache) Get(_ context.Context, key string) (string, error) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return "", nil
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return "", nil
	}

	return item.value, nil
}

func (c *InMemoryCache) Set(_ context.Context, key string, value string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}

	c.mu.Lock()
	c.items[key] = memoryCacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()

	return nil
}

func (c *InMemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
	return nil
}

func (c *InMemoryCache) SetIfAbsent(_ context.Context, key, value string, ttl time.Duration) (bool, error) {
	if ttl <= 0 {
		return false, fmt.Errorf("cache TTL must be positive")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for expiredKey, item := range c.items {
		if !item.expiresAt.IsZero() && !now.Before(item.expiresAt) {
			delete(c.items, expiredKey)
		}
	}
	if item, ok := c.items[key]; ok && now.Before(item.expiresAt) {
		return false, nil
	}
	c.items[key] = memoryCacheItem{value: value, expiresAt: now.Add(ttl)}
	return true, nil
}

func (c *InMemoryCache) IncrementWindow(_ context.Context, key string, window time.Duration) (int64, time.Duration, error) {
	if window <= 0 {
		return 0, 0, fmt.Errorf("rate limit window must be positive")
	}
	if window.Milliseconds() < 1 {
		return 0, 0, fmt.Errorf("rate limit window must be at least one millisecond")
	}
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if !ok || !now.Before(item.expiresAt) {
		item = memoryCacheItem{value: "0", expiresAt: now.Add(window)}
	}
	var count int64
	if _, err := fmt.Sscan(item.value, &count); err != nil {
		return 0, 0, fmt.Errorf("invalid rate limit counter")
	}
	count++
	item.value = fmt.Sprint(count)
	c.items[key] = item
	return count, time.Until(item.expiresAt), nil
}

func (c *InMemoryCache) Rotate(_ context.Context, oldKey, newKey, value string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[oldKey]
	if !ok || (!item.expiresAt.IsZero() && time.Now().After(item.expiresAt)) {
		delete(c.items, oldKey)
		return ErrSessionNotFound
	}
	delete(c.items, oldKey)
	c.items[newKey] = memoryCacheItem{value: value, expiresAt: time.Now().Add(ttl)}
	return nil
}

func (c *InMemoryCache) Touch(_ context.Context, key, expectedValue string, ttl time.Duration) error {
	if ttl <= 0 {
		return ErrSessionNotFound
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if !ok {
		return ErrSessionNotFound
	}
	now := time.Now()
	if !item.expiresAt.IsZero() && !now.Before(item.expiresAt) {
		delete(c.items, key)
		return ErrSessionNotFound
	}
	if item.value != expectedValue {
		return ErrSessionNotFound
	}
	item.expiresAt = now.Add(ttl)
	c.items[key] = item
	return nil
}

// RedisCache stores values in Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a Redis cache using REDIS_URL.
func NewRedisCache(redisURL string) (*RedisCache, error) {
	options, err := redis.ParseURL(strings.TrimSpace(redisURL))
	if err != nil {
		return nil, fmt.Errorf("invalid redis url: %w", err)
	}

	client := redis.NewClient(options)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	if c == nil || c.client == nil {
		return "", nil
	}

	value, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if c == nil || c.client == nil {
		return nil
	}
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) SetIfAbsent(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	if c == nil || c.client == nil {
		return false, fmt.Errorf("redis one-time store is unavailable")
	}
	if ttl <= 0 {
		return false, fmt.Errorf("cache TTL must be positive")
	}
	return c.client.SetNX(ctx, key, value, ttl).Result()
}

var incrementWindowScript = redis.NewScript(`
local count = redis.call("INCR", KEYS[1])
if count == 1 then
	redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
local ttl = redis.call("PTTL", KEYS[1])
if ttl < 0 then
	redis.call("PEXPIRE", KEYS[1], ARGV[1])
	ttl = redis.call("PTTL", KEYS[1])
end
return {count, ttl}
`)

func (c *RedisCache) IncrementWindow(ctx context.Context, key string, window time.Duration) (int64, time.Duration, error) {
	if c == nil || c.client == nil {
		return 0, 0, fmt.Errorf("redis rate limit store is unavailable")
	}
	if window <= 0 {
		return 0, 0, fmt.Errorf("rate limit window must be positive")
	}
	if window.Milliseconds() < 1 {
		return 0, 0, fmt.Errorf("rate limit window must be at least one millisecond")
	}
	result, err := incrementWindowScript.Run(ctx, c.client, []string{key}, window.Milliseconds()).Slice()
	if err != nil {
		return 0, 0, err
	}
	if len(result) != 2 {
		return 0, 0, fmt.Errorf("unexpected Redis rate limit response")
	}
	count, countOK := result[0].(int64)
	ttlMillis, ttlOK := result[1].(int64)
	if !countOK || !ttlOK || count < 1 {
		return 0, 0, fmt.Errorf("invalid Redis rate limit response")
	}
	if ttlMillis < 0 {
		ttlMillis = 0
	}
	return count, time.Duration(ttlMillis) * time.Millisecond, nil
}

var rotateSessionScript = redis.NewScript(`
if redis.call("EXISTS", KEYS[1]) == 0 then
	return 0
end
redis.call("SET", KEYS[2], ARGV[1], "PX", ARGV[2])
redis.call("DEL", KEYS[1])
return 1
`)

func (c *RedisCache) Rotate(ctx context.Context, oldKey, newKey, value string, ttl time.Duration) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("redis session store is unavailable")
	}
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	ttlMilliseconds := ttl.Milliseconds()
	if ttlMilliseconds <= 0 {
		return ErrSessionNotFound
	}
	rotated, err := rotateSessionScript.Run(ctx, c.client, []string{oldKey, newKey}, value, ttlMilliseconds).Int()
	if err != nil {
		return err
	}
	if rotated != 1 {
		return ErrSessionNotFound
	}
	return nil
}

var touchSessionScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) ~= ARGV[1] then
	return 0
end
redis.call("PEXPIRE", KEYS[1], ARGV[2])
return 1
`)

func (c *RedisCache) Touch(ctx context.Context, key, expectedValue string, ttl time.Duration) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("redis session store is unavailable")
	}
	ttlMilliseconds := ttl.Milliseconds()
	if ttlMilliseconds <= 0 {
		return ErrSessionNotFound
	}
	touched, err := touchSessionScript.Run(ctx, c.client, []string{key}, expectedValue, ttlMilliseconds).Int()
	if err != nil {
		return err
	}
	if touched != 1 {
		return ErrSessionNotFound
	}
	return nil
}

func (c *RedisCache) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

// NewCacheFromEnv returns Redis cache when REDIS_URL is configured, otherwise in-memory cache.
func NewCacheFromEnv() KVCache {
	redisURL := strings.TrimSpace(os.Getenv("REDIS_URL"))
	if redisURL == "" {
		return NewInMemoryCache()
	}

	redisCache, err := NewRedisCache(redisURL)
	if err != nil {
		log.Printf("[Cache] Redis unavailable, falling back to in-memory cache: %v", err)
		return NewInMemoryCache()
	}
	return redisCache
}
