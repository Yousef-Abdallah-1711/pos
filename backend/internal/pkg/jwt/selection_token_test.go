package jwt

import (
	"testing"
	"time"
)

func newSvc(t *testing.T) *TokenService {
	t.Helper()
	s, err := NewTokenService("unit-test-secret-unit-test-secret-1234567", "1h")
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSelectionToken_RoundTrip(t *testing.T) {
	s := newSvc(t)
	tok, err := s.GenerateSelectionToken("A@Example.test", []int{3, 9}, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	c, err := s.ValidateSelectionToken(tok)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if c.Email != "a@example.test" || len(c.TenantIDs) != 2 || c.ID == "" {
		t.Fatalf("unexpected claims %+v", c)
	}
	if !c.AllowsTenant(9) || c.AllowsTenant(4) {
		t.Fatalf("tenant allow-list wrong: %+v", c.TenantIDs)
	}
}

func TestSelectionToken_ExpiredIsRejected(t *testing.T) {
	s := newSvc(t)
	tok, _ := s.GenerateSelectionToken("a@example.test", []int{3}, -time.Minute)
	if _, err := s.ValidateSelectionToken(tok); err == nil {
		t.Fatal("expired selection token must be rejected")
	}
}

func TestSelectionToken_TamperedOrForeignSecretRejected(t *testing.T) {
	s := newSvc(t)
	other, _ := NewTokenService("a-completely-different-secret-value-xyz", "1h")
	tok, _ := other.GenerateSelectionToken("a@example.test", []int{3}, time.Minute)
	if _, err := s.ValidateSelectionToken(tok); err == nil {
		t.Fatal("token signed with another secret must be rejected")
	}
	if _, err := s.ValidateSelectionToken("not.a.jwt"); err == nil {
		t.Fatal("garbage must be rejected")
	}
}

// A normal session token must not be usable as a selection proof, and vice versa.
func TestSelectionToken_NotInterchangeableWithSessionToken(t *testing.T) {
	s := newSvc(t)
	rid := 5
	session, _ := s.GenerateToken(7, 3, &rid, "a@example.test", "owner")
	if _, err := s.ValidateSelectionToken(session); err == nil {
		t.Fatal("session token must not validate as a selection token")
	}
	sel, _ := s.GenerateSelectionToken("a@example.test", []int{3}, time.Minute)
	if _, err := s.ValidateToken(sel); err == nil {
		t.Fatal("selection token must not validate as a session token")
	}
}
