package usecase

import (
	"errors"
	"testing"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/jwt"
)

type fakeConfirmStore struct {
	users      []domain.User
	lastLogins []int
}

func (f *fakeConfirmStore) FindUserByEmailAllTenants(email string) ([]domain.User, error) {
	return f.users, nil
}
func (f *fakeConfirmStore) UpdateLastLogin(userID int) error {
	f.lastLogins = append(f.lastLogins, userID)
	return nil
}

func newConfirmUC(t *testing.T) (*AuthUseCase, *jwt.TokenService, *fakeConfirmStore) {
	t.Helper()
	ts, err := jwt.NewTokenService("unit-test-secret-unit-test-secret-1234567", "1h")
	if err != nil {
		t.Fatal(err)
	}
	store := &fakeConfirmStore{users: []domain.User{
		{ID: 1, TenantID: 3, Email: "a@example.test", Role: "owner"},
		{ID: 2, TenantID: 9, Email: "a@example.test", Role: "manager"},
	}}
	uc := &AuthUseCase{tokenService: ts}
	uc.confirmStore = store
	return uc, ts, store
}

// SEC-C: a bare email + tenant_id must not yield a session.
func TestLoginConfirm_WithoutSelectionTokenIsRejected(t *testing.T) {
	uc, _, store := newConfirmUC(t)
	resp, err := uc.LoginConfirm(&domain.TenantSelectionRequest{Email: "a@example.test", TenantID: 3})
	if err == nil || resp != nil {
		t.Fatalf("must be rejected, got resp=%v err=%v", resp, err)
	}
	if len(store.lastLogins) != 0 {
		t.Fatal("no login side effects allowed")
	}
}

func TestLoginConfirm_ForgedOrWrongSecretTokenIsRejected(t *testing.T) {
	uc, _, _ := newConfirmUC(t)
	other, _ := jwt.NewTokenService("another-secret-another-secret-another-1", "1h")
	forged, _ := other.GenerateSelectionToken("a@example.test", []int{3}, time.Minute)
	if _, err := uc.LoginConfirm(&domain.TenantSelectionRequest{Email: "a@example.test", TenantID: 3, SelectionToken: forged}); err == nil {
		t.Fatal("forged token must be rejected")
	}
	if _, err := uc.LoginConfirm(&domain.TenantSelectionRequest{Email: "a@example.test", TenantID: 3, SelectionToken: "garbage"}); err == nil {
		t.Fatal("garbage token must be rejected")
	}
}

func TestLoginConfirm_ExpiredTokenIsRejected(t *testing.T) {
	uc, ts, _ := newConfirmUC(t)
	tok, _ := ts.GenerateSelectionToken("a@example.test", []int{3}, -time.Second)
	if _, err := uc.LoginConfirm(&domain.TenantSelectionRequest{Email: "a@example.test", TenantID: 3, SelectionToken: tok}); err == nil {
		t.Fatal("expired token must be rejected")
	}
}

func TestLoginConfirm_TenantOutsideProofIsRejected(t *testing.T) {
	uc, ts, _ := newConfirmUC(t)
	tok, _ := ts.GenerateSelectionToken("a@example.test", []int{3}, time.Minute)
	if _, err := uc.LoginConfirm(&domain.TenantSelectionRequest{Email: "a@example.test", TenantID: 9, SelectionToken: tok}); err == nil {
		t.Fatal("tenant not covered by the proof must be rejected")
	}
}

func TestLoginConfirm_EmailMismatchIsRejected(t *testing.T) {
	uc, ts, _ := newConfirmUC(t)
	tok, _ := ts.GenerateSelectionToken("a@example.test", []int{3}, time.Minute)
	if _, err := uc.LoginConfirm(&domain.TenantSelectionRequest{Email: "victim@example.test", TenantID: 3, SelectionToken: tok}); err == nil {
		t.Fatal("email not matching the proof must be rejected")
	}
}

func TestLoginConfirm_ValidProofIssuesSessionForSelectedTenantOnce(t *testing.T) {
	uc, ts, store := newConfirmUC(t)
	tok, _ := ts.GenerateSelectionToken("A@example.test", []int{3, 9}, time.Minute)
	req := &domain.TenantSelectionRequest{Email: "a@example.test", TenantID: 9, SelectionToken: tok}

	resp, err := uc.LoginConfirm(req)
	if err != nil {
		t.Fatalf("valid confirm failed: %v", err)
	}
	claims, err := ts.ValidateToken(resp.Token)
	if err != nil || claims.TenantID != 9 || claims.UserID != 2 || claims.Role != "manager" {
		t.Fatalf("session must be for tenant 9 / user 2, got %+v err=%v", claims, err)
	}
	if resp.User.PasswordHash != "" {
		t.Fatal("password hash must be cleared")
	}
	if len(store.lastLogins) != 1 {
		t.Fatalf("expected one last-login update, got %d", len(store.lastLogins))
	}

	// replay of the same proof is rejected
	if _, err := uc.LoginConfirm(req); err == nil {
		t.Fatal("replayed selection token must be rejected")
	}
}

func TestLoginConfirm_StoreFailureDoesNotBurnProof(t *testing.T) {
	uc, ts, _ := newConfirmUC(t)
	uc.confirmStore = failingStore{}
	tok, _ := ts.GenerateSelectionToken("a@example.test", []int{3}, time.Minute)
	req := &domain.TenantSelectionRequest{Email: "a@example.test", TenantID: 3, SelectionToken: tok}
	if _, err := uc.LoginConfirm(req); err == nil {
		t.Fatal("expected store error")
	}
	uc.confirmStore = &fakeConfirmStore{users: []domain.User{{ID: 1, TenantID: 3, Email: "a@example.test", Role: "owner"}}}
	if _, err := uc.LoginConfirm(req); err != nil {
		t.Fatalf("proof should still be usable after a transient failure: %v", err)
	}
}

type failingStore struct{}

func (failingStore) FindUserByEmailAllTenants(string) ([]domain.User, error) {
	return nil, errors.New("db down")
}
func (failingStore) UpdateLastLogin(int) error { return nil }
