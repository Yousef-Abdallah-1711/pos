package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

const authSessionKeyPrefix = "auth-session:v1:"

var ErrSessionNotFound = errors.New("session not found")

// AuthSessionStore stores only a hash of the browser session ID in shared KV storage.
type AuthSessionStore struct {
	cache KVCache
	ttl   time.Duration
}

func NewAuthSessionStore(cache KVCache, ttl time.Duration) *AuthSessionStore {
	return &AuthSessionStore{cache: cache, ttl: ttl}
}

func (s *AuthSessionStore) Create(ctx context.Context, credential string) (string, error) {
	if s == nil || s.cache == nil || s.ttl <= 0 || credential == "" {
		return "", errors.New("session store is not configured")
	}
	randomID := make([]byte, 32)
	if _, err := rand.Read(randomID); err != nil {
		return "", err
	}
	sessionID := base64.RawURLEncoding.EncodeToString(randomID)
	if err := s.cache.Set(ctx, sessionKey(sessionID), credential, s.ttl); err != nil {
		return "", err
	}
	return sessionID, nil
}

// Rotate atomically consumes an existing browser session ID and returns a new
// random ID bound to the replacement credential. Stores without atomic support
// fail closed rather than exposing a partially rotated session.
func (s *AuthSessionStore) Rotate(ctx context.Context, oldSessionID, credential string) (string, error) {
	if s == nil || s.cache == nil || s.ttl <= 0 || oldSessionID == "" || credential == "" {
		return "", errors.New("session rotation is not configured")
	}
	atomicCache, ok := s.cache.(AtomicSessionCache)
	if !ok {
		return "", errors.New("session store does not support atomic rotation")
	}
	randomID := make([]byte, 32)
	if _, err := rand.Read(randomID); err != nil {
		return "", err
	}
	newSessionID := base64.RawURLEncoding.EncodeToString(randomID)
	if err := atomicCache.Rotate(ctx, sessionKey(oldSessionID), sessionKey(newSessionID), credential, s.ttl); err != nil {
		return "", err
	}
	return newSessionID, nil
}

func (s *AuthSessionStore) Resolve(ctx context.Context, sessionID string) (string, error) {
	if s == nil || s.cache == nil || sessionID == "" {
		return "", ErrSessionNotFound
	}
	credential, err := s.cache.Get(ctx, sessionKey(sessionID))
	if err != nil {
		return "", err
	}
	if credential == "" {
		return "", ErrSessionNotFound
	}
	return credential, nil
}

func (s *AuthSessionStore) Revoke(ctx context.Context, sessionID string) error {
	if s == nil || s.cache == nil || sessionID == "" {
		return nil
	}
	return s.cache.Delete(ctx, sessionKey(sessionID))
}

func sessionKey(sessionID string) string {
	digest := sha256.Sum256([]byte(sessionID))
	return authSessionKeyPrefix + hex.EncodeToString(digest[:])
}

// AuthSessionCookieName is intentionally host-only: no Domain attribute is set.
const AuthSessionCookieName = "poos_session"
const SuperAdminSessionCookieName = "poos_super_admin_session"
const ImpersonationSessionCookieName = "poos_impersonation"
const impersonationCredentialPrefix = "impersonation-session:"

func (s *AuthSessionStore) CreateImpersonation(ctx context.Context, token string, ttl time.Duration) (string, error) {
	if ttl <= 0 || s == nil || s.cache == nil || token == "" {
		return "", errors.New("impersonation session is not configured")
	}
	randomID := make([]byte, 32)
	if _, err := rand.Read(randomID); err != nil {
		return "", err
	}
	sessionID := base64.RawURLEncoding.EncodeToString(randomID)
	if err := s.cache.Set(ctx, sessionKey(sessionID), impersonationCredentialPrefix+token, ttl); err != nil {
		return "", err
	}
	return sessionID, nil
}

func IsImpersonationCredential(credential string) bool {
	return strings.HasPrefix(credential, impersonationCredentialPrefix)
}

func UnwrapImpersonationCredential(credential string) string {
	return strings.TrimPrefix(credential, impersonationCredentialPrefix)
}

func (s *AuthSessionStore) TTL() time.Duration {
	if s == nil {
		return 0
	}
	return s.ttl
}
