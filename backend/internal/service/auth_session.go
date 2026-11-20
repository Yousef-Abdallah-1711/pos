package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const authSessionKeyPrefix = "auth-session:v1:"

const StaffSessionIdleTimeout = 30 * time.Minute

const StaffSessionAbsoluteLifetime = 8 * time.Hour

var ErrSessionNotFound = errors.New("session not found")

// AuthSessionStore stores only a hash of the browser session ID in shared KV storage.
type AuthSessionStore struct {
	cache KVCache
	ttl   time.Duration
}

type authSessionRecord struct {
	Version           int           `json:"v"`
	Credential        string        `json:"credential"`
	IdleTimeout       time.Duration `json:"idle_timeout"`
	AbsoluteExpiresAt time.Time     `json:"absolute_expires_at"`
}

func NewAuthSessionStore(cache KVCache, ttl time.Duration) *AuthSessionStore {
	return &AuthSessionStore{cache: cache, ttl: ttl}
}

func (s *AuthSessionStore) Create(ctx context.Context, credential string) (string, error) {
	if s == nil {
		return "", errors.New("session store is not configured")
	}
	return s.CreateWithLifetimes(ctx, credential, s.ttl, s.ttl)
}

// CreateWithLifetimes creates a browser session with a sliding idle timeout
// capped by an absolute lifetime from initial authentication.
func (s *AuthSessionStore) CreateWithLifetimes(ctx context.Context, credential string, idleTimeout, absoluteLifetime time.Duration) (string, error) {
	if s == nil || s.cache == nil || credential == "" || idleTimeout <= 0 || absoluteLifetime <= 0 {
		return "", errors.New("session store is not configured")
	}
	randomID := make([]byte, 32)
	if _, err := rand.Read(randomID); err != nil {
		return "", err
	}
	sessionID := base64.RawURLEncoding.EncodeToString(randomID)
	now := time.Now()
	record := authSessionRecord{
		Version:           1,
		Credential:        credential,
		IdleTimeout:       idleTimeout,
		AbsoluteExpiresAt: now.Add(absoluteLifetime),
	}
	value, err := encodeAuthSessionRecord(record)
	if err != nil {
		return "", err
	}
	if err := s.cache.Set(ctx, sessionKey(sessionID), value, sessionCacheTTL(record, now)); err != nil {
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
	oldKey := sessionKey(oldSessionID)
	value, err := s.cache.Get(ctx, oldKey)
	if err != nil {
		return "", err
	}
	if value == "" {
		return "", ErrSessionNotFound
	}
	record, err := decodeAuthSessionRecord(value)
	if err != nil || !time.Now().Before(record.AbsoluteExpiresAt) {
		_ = s.cache.Delete(ctx, oldKey)
		return "", ErrSessionNotFound
	}
	record.Credential = credential
	value, err = encodeAuthSessionRecord(record)
	if err != nil {
		return "", err
	}
	now := time.Now()
	ttl := sessionCacheTTL(record, now)
	if ttl.Milliseconds() <= 0 {
		_ = s.cache.Delete(ctx, oldKey)
		return "", ErrSessionNotFound
	}
	randomID := make([]byte, 32)
	if _, err := rand.Read(randomID); err != nil {
		return "", err
	}
	newSessionID := base64.RawURLEncoding.EncodeToString(randomID)
	if err := atomicCache.Rotate(ctx, oldKey, sessionKey(newSessionID), value, ttl); err != nil {
		return "", err
	}
	return newSessionID, nil
}

// RotateWithLifetimes atomically replaces an existing session with a fresh
// authenticated session and new lifetime bounds. It is used after a fresh
// password-verified login; refreshes use Rotate to retain the original cap.
func (s *AuthSessionStore) RotateWithLifetimes(ctx context.Context, oldSessionID, credential string, idleTimeout, absoluteLifetime time.Duration) (string, error) {
	if s == nil || s.cache == nil || oldSessionID == "" || credential == "" || idleTimeout <= 0 || absoluteLifetime <= 0 {
		return "", errors.New("session rotation is not configured")
	}
	atomicCache, ok := s.cache.(AtomicSessionCache)
	if !ok {
		return "", errors.New("session store does not support atomic rotation")
	}
	oldKey := sessionKey(oldSessionID)
	previous, err := s.cache.Get(ctx, oldKey)
	if err != nil {
		return "", err
	}
	if previous == "" {
		return "", ErrSessionNotFound
	}
	now := time.Now()
	record := authSessionRecord{
		Version:           1,
		Credential:        credential,
		IdleTimeout:       idleTimeout,
		AbsoluteExpiresAt: now.Add(absoluteLifetime),
	}
	value, err := encodeAuthSessionRecord(record)
	if err != nil {
		return "", err
	}
	ttl := sessionCacheTTL(record, now)
	if ttl.Milliseconds() <= 0 {
		return "", ErrSessionNotFound
	}
	randomID := make([]byte, 32)
	if _, err := rand.Read(randomID); err != nil {
		return "", err
	}
	newSessionID := base64.RawURLEncoding.EncodeToString(randomID)
	if err := atomicCache.Rotate(ctx, oldKey, sessionKey(newSessionID), value, ttl); err != nil {
		return "", err
	}
	return newSessionID, nil
}

func (s *AuthSessionStore) Resolve(ctx context.Context, sessionID string) (string, error) {
	if s == nil || s.cache == nil || sessionID == "" {
		return "", ErrSessionNotFound
	}
	key := sessionKey(sessionID)
	value, err := s.cache.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if value == "" {
		return "", ErrSessionNotFound
	}
	// Impersonation sessions use a separate, explicit short TTL and retain
	// their existing format. All staff and Super Admin sessions use records.
	if strings.HasPrefix(value, impersonationCredentialPrefix) {
		return value, nil
	}
	record, err := decodeAuthSessionRecord(value)
	if err != nil {
		return "", ErrSessionNotFound
	}
	now := time.Now()
	ttl := sessionCacheTTL(record, now)
	if ttl <= 0 {
		_ = s.cache.Delete(ctx, key)
		return "", ErrSessionNotFound
	}
	slidingCache, ok := s.cache.(SlidingSessionCache)
	if !ok {
		return "", errors.New("session store does not support sliding expiration")
	}
	if err := slidingCache.Touch(ctx, key, value, ttl); err != nil {
		return "", err
	}
	return record.Credential, nil
}

// AbsoluteLifetimeRemaining returns the hard time remaining for a managed
// browser session. It does not extend the idle timeout.
func (s *AuthSessionStore) AbsoluteLifetimeRemaining(ctx context.Context, sessionID string) (time.Duration, error) {
	if s == nil || s.cache == nil || sessionID == "" {
		return 0, ErrSessionNotFound
	}
	value, err := s.cache.Get(ctx, sessionKey(sessionID))
	if err != nil {
		return 0, err
	}
	if value == "" || strings.HasPrefix(value, impersonationCredentialPrefix) {
		return 0, ErrSessionNotFound
	}
	record, err := decodeAuthSessionRecord(value)
	if err != nil {
		return 0, ErrSessionNotFound
	}
	remaining := time.Until(record.AbsoluteExpiresAt)
	if remaining <= 0 {
		return 0, ErrSessionNotFound
	}
	return remaining, nil
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

func encodeAuthSessionRecord(record authSessionRecord) (string, error) {
	encoded, err := json.Marshal(record)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func decodeAuthSessionRecord(value string) (authSessionRecord, error) {
	var record authSessionRecord
	if err := json.Unmarshal([]byte(value), &record); err != nil {
		return authSessionRecord{}, err
	}
	if record.Version != 1 || record.Credential == "" || record.IdleTimeout <= 0 || record.AbsoluteExpiresAt.IsZero() {
		return authSessionRecord{}, errors.New("invalid session record")
	}
	return record, nil
}

func sessionCacheTTL(record authSessionRecord, now time.Time) time.Duration {
	remaining := record.AbsoluteExpiresAt.Sub(now)
	if remaining <= 0 {
		return 0
	}
	if remaining < record.IdleTimeout {
		return remaining
	}
	return record.IdleTimeout
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
