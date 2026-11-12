package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"unicode"
)

// Public order access.
//
// Public (unauthenticated) order endpoints are keyed by sequential database ids and
// guessable order numbers, so the id/number alone must never be enough to read or
// cancel an order. Every order gets an unguessable tracking token when it is created:
//   - the raw token (256 bits, URL-safe) is returned exactly once, in the create response;
//   - only its SHA-256 hash is stored, so a database leak does not expose live tokens;
//   - public endpoints require the token (header X-Order-Token), or, for the manual
//     "track my order" lookup by order number, the customer phone number used at checkout.
// Orders without a stored hash (created before this mechanism existed) are not
// publicly accessible.

// PublicTokenBytes is the number of random bytes in a tracking token.
const PublicTokenBytes = 32

// GeneratePublicTrackingToken returns a new raw token and the hash to store.
func GeneratePublicTrackingToken() (raw string, hash string, err error) {
	b := make([]byte, PublicTokenBytes)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	return raw, HashPublicTrackingToken(raw), nil
}

// HashPublicTrackingToken returns the hex SHA-256 of a raw token.
func HashPublicTrackingToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// VerifyPublicTrackingToken compares a presented raw token with the stored hash in constant time.
// An empty presented token or empty stored hash never verifies.
func VerifyPublicTrackingToken(raw, storedHash string) bool {
	if raw == "" || storedHash == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(HashPublicTrackingToken(raw)), []byte(storedHash)) == 1
}

// NormalizePhone keeps only digits, so "+20 100-000 0000" and "201000000000" compare equal.
func NormalizePhone(phone string) string {
	var sb strings.Builder
	for _, r := range phone {
		if unicode.IsDigit(r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// PhonesMatch reports whether two phone numbers are the same after normalisation.
// Numbers shorter than 7 digits never match (too weak to act as a proof).
func PhonesMatch(a, b string) bool {
	na, nb := NormalizePhone(a), NormalizePhone(b)
	if len(na) < 7 || len(nb) < 7 {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(na), []byte(nb)) == 1
}
