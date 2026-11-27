package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/workspaceappearance"
)

type fakeWorkspaceStore struct {
	saveErr error
	saves   int
}

func (f *fakeWorkspaceStore) Get(context.Context, int64, int64) (workspaceappearance.Published, error) {
	return workspaceappearance.Published{SchemaVersion: 1, Overrides: json.RawMessage(`{}`)}, nil
}
func (f *fakeWorkspaceStore) Save(_ context.Context, _, _, revision int64, _ json.RawMessage) (workspaceappearance.Published, error) {
	f.saves++
	return workspaceappearance.Published{SchemaVersion: 1, Revision: revision + 1, Overrides: json.RawMessage(`{}`)}, f.saveErr
}
func workspaceRequest(body string) *http.Request {
	r := httptest.NewRequest(http.MethodPut, "/api/v1/workspace-appearance", strings.NewReader(body))
	ctx := context.WithValue(r.Context(), middleware.UserContextKey, &jwt.Claims{UserID: 11, TenantID: 1})
	return r.WithContext(context.WithValue(ctx, middleware.TenantContextKey, int64(1)))
}
func TestWorkspaceAppearanceRejectsInvalidInputAndStaleRevision(t *testing.T) {
	store := &fakeWorkspaceStore{}
	handler := NewWorkspaceAppearanceHandler(store)
	for _, body := range []string{`{}`, `{"expectedRevision":0,"overrides":{"colors":{"light":{"primary":"red"}}}}`} {
		w := httptest.NewRecorder()
		handler.Put(w, workspaceRequest(body))
		if w.Code != http.StatusBadRequest || store.saves != 0 {
			t.Fatalf("invalid body accepted: %d", w.Code)
		}
	}
	store.saveErr = workspaceappearance.ErrConflict
	w := httptest.NewRecorder()
	handler.Put(w, workspaceRequest(`{"expectedRevision":0,"overrides":{}}`))
	if w.Code != http.StatusConflict {
		t.Fatalf("expected conflict, got %d", w.Code)
	}
}
func TestWorkspaceAppearanceRequiresMatchingTenant(t *testing.T) {
	store := &fakeWorkspaceStore{}
	r := workspaceRequest(`{"expectedRevision":0,"overrides":{}}`)
	r = r.WithContext(context.WithValue(r.Context(), middleware.TenantContextKey, int64(2)))
	w := httptest.NewRecorder()
	NewWorkspaceAppearanceHandler(store).Put(w, r)
	if w.Code != http.StatusForbidden || store.saves != 0 {
		t.Fatalf("cross-tenant request: %d", w.Code)
	}
}
