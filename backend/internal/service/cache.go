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

type memoryCacheItem struct {
	value     string
	expiresAt time.Time
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
	rotated, err := rotateSessionScript.Run(ctx, c.client, []string{oldKey, newKey}, value, ttl.Milliseconds()).Int()
	if err != nil {
		return err
	}
	if rotated != 1 {
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
