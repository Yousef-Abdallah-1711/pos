package service

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestAuthSessionStoreCreatesResolvesAndRevokesOpaqueSession(t *testing.T) {
	store := NewAuthSessionStore(NewInMemoryCache(), time.Minute)
	ctx := context.Background()
	const credential = "synthetic-jwt-for-session-test"

	sessionID, err := store.Create(ctx, credential)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if len(sessionID) < 40 || strings.Contains(sessionID, credential) {
		t.Fatalf("session id must be opaque and high entropy")
	}
	resolved, err := store.Resolve(ctx, sessionID)
	if err != nil || resolved != credential {
		t.Fatalf("resolve session = %q, %v", resolved, err)
	}
	if err := store.Revoke(ctx, sessionID); err != nil {
		t.Fatalf("revoke session: %v", err)
	}
	if _, err := store.Resolve(ctx, sessionID); err != ErrSessionNotFound {
		t.Fatalf("revoked session should be rejected, got %v", err)
	}
}

func TestAuthSessionStoreExpiresSession(t *testing.T) {
	store := NewAuthSessionStore(NewInMemoryCache(), 10*time.Millisecond)
	sessionID, err := store.Create(context.Background(), "synthetic-expiring-credential")
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	if _, err := store.Resolve(context.Background(), sessionID); err != ErrSessionNotFound {
		t.Fatalf("expired session should be rejected, got %v", err)
	}
}

func TestAuthSessionStoreRotationConsumesPriorSessionOnce(t *testing.T) {
	store := NewAuthSessionStore(NewInMemoryCache(), time.Minute)
	ctx := context.Background()
	oldID, err := store.Create(ctx, "pre-auth-session")
	if err != nil {
		t.Fatal(err)
	}
	const attempts = 20
	var successes int
	var winner string
	var mu sync.Mutex
	var wait sync.WaitGroup
	for range attempts {
		wait.Add(1)
		go func() {
			defer wait.Done()
			newID, rotateErr := store.Rotate(ctx, oldID, "authenticated-session")
			mu.Lock()
			defer mu.Unlock()
			if rotateErr == nil {
				successes++
				winner = newID
			} else if !errors.Is(rotateErr, ErrSessionNotFound) {
				t.Errorf("unexpected rotation failure: %v", rotateErr)
			}
		}()
	}
	wait.Wait()
	if successes != 1 {
		t.Fatalf("exactly one concurrent rotation should consume the old session; got %d", successes)
	}
	if _, err := store.Resolve(ctx, oldID); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("old session remained valid after rotation: %v", err)
	}
	if credential, err := store.Resolve(ctx, winner); err != nil || credential != "authenticated-session" {
		t.Fatalf("rotated session lookup = %q, %v", credential, err)
	}
}

func TestAuthSessionStoreRedisIntegration(t *testing.T) {
	redisURL := os.Getenv("POOS_SESSION_REDIS_URL")
	if redisURL == "" {
		t.Skip("set POOS_SESSION_REDIS_URL to run against the identified disposable Redis service")
	}
	firstCache, err := NewRedisCache(redisURL)
	if err != nil {
		t.Fatalf("connect to configured disposable Redis: %v", err)
	}
	defer firstCache.Close()
	secondCache, err := NewRedisCache(redisURL)
	if err != nil {
		t.Fatalf("connect second independent Redis client: %v", err)
	}
	defer secondCache.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	first := NewAuthSessionStore(firstCache, time.Minute)
	second := NewAuthSessionStore(secondCache, time.Minute)
	oldID, err := first.Create(ctx, "redis-integration-credential")
	if err != nil {
		t.Fatalf("create Redis session: %v", err)
	}
	t.Cleanup(func() { _ = first.Revoke(context.Background(), oldID) })
	if credential, err := second.Resolve(ctx, oldID); err != nil || credential != "redis-integration-credential" {
		t.Fatalf("cross-client lookup = %q, %v", credential, err)
	}

	const attempts = 20
	var successes int
	var winner string
	var mu sync.Mutex
	var wait sync.WaitGroup
	for range attempts {
		wait.Add(1)
		go func() {
			defer wait.Done()
			newID, rotateErr := second.Rotate(ctx, oldID, "redis-rotated-credential")
			mu.Lock()
			defer mu.Unlock()
			if rotateErr == nil {
				successes++
				winner = newID
			} else if !errors.Is(rotateErr, ErrSessionNotFound) {
				t.Errorf("unexpected Redis rotation failure: %v", rotateErr)
			}
		}()
	}
	wait.Wait()
	if successes != 1 {
		t.Fatalf("exactly one Redis rotation should consume the old session; got %d", successes)
	}
	t.Cleanup(func() { _ = first.Revoke(context.Background(), winner) })
	if _, err := first.Resolve(ctx, oldID); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("old Redis session remained valid: %v", err)
	}
	if credential, err := first.Resolve(ctx, winner); err != nil || credential != "redis-rotated-credential" {
		t.Fatalf("rotated Redis session lookup = %q, %v", credential, err)
	}
	if err := second.Revoke(ctx, winner); err != nil {
		t.Fatalf("Redis session revoke: %v", err)
	}
	if _, err := first.Resolve(ctx, winner); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("revocation was not visible to the other Redis client: %v", err)
	}
	if _, err := first.Resolve(ctx, "malformed session identifier"); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("malformed session identifier was not rejected: %v", err)
	}

	shortLived := NewAuthSessionStore(firstCache, 120*time.Millisecond)
	expiringID, err := shortLived.Create(ctx, "redis-expiring-credential")
	if err != nil {
		t.Fatalf("create expiring Redis session: %v", err)
	}
	t.Cleanup(func() { _ = first.Revoke(context.Background(), expiringID) })
	time.Sleep(180 * time.Millisecond)
	if _, err := second.Resolve(ctx, expiringID); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("Redis TTL was not enforced across clients: %v", err)
	}
	recoveryID, err := first.Create(ctx, "redis-recovery-credential")
	if err != nil {
		t.Fatalf("create recovery probe session: %v", err)
	}
	t.Cleanup(func() { _ = first.Revoke(context.Background(), recoveryID) })

	unavailableCache := &RedisCache{client: redis.NewClient(&redis.Options{Addr: "127.0.0.1:1", DialTimeout: 100 * time.Millisecond, ReadTimeout: 100 * time.Millisecond, WriteTimeout: 100 * time.Millisecond})}
	defer unavailableCache.Close()
	unavailableStore := NewAuthSessionStore(unavailableCache, time.Minute)
	if _, err := unavailableStore.Resolve(ctx, oldID); err == nil {
		t.Fatal("unavailable Redis must fail closed instead of authenticating")
	}
	if credential, err := first.Resolve(ctx, recoveryID); err != nil || credential != "redis-recovery-credential" {
		t.Fatalf("healthy Redis client did not recover after isolated store failure: %q, %v", credential, err)
	}
}
