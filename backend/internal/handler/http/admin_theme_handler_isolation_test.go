package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/service"
)

// themeSvcStub embeds the interface so any method the handler must NOT reach
// panics (nil interface) instead of silently succeeding.
type themeSvcStub struct {
	service.ThemeService
	themes     map[int64]*domain.Theme
	presets    map[int64]*domain.Theme
	deleted    []int64
	activated  []int64
	updatedIDs []int64
}

func (s *themeSvcStub) GetTheme(_ context.Context, tenantID, restaurantID, id int64) (*domain.Theme, error) {
	if t, ok := s.themes[id]; ok {
		if t.TenantID == tenantID && t.RestaurantID == restaurantID {
			return t, nil
		}
	}
	return nil, errors.New("theme not found")
}
func (s *themeSvcStub) GetPreset(_ context.Context, id int64) (*domain.Theme, error) {
	if t, ok := s.presets[id]; ok {
		return t, nil
	}
	return nil, errors.New("preset not found")
}
func (s *themeSvcStub) DeleteTheme(_ context.Context, tenantID, restaurantID, id int64) error {
	if t, ok := s.themes[id]; ok && t.TenantID == tenantID && t.RestaurantID == restaurantID {
		s.deleted = append(s.deleted, id)
	}
	return nil
}
func (s *themeSvcStub) ActivateTheme(_ context.Context, tenantID, restaurantID, id int64) error {
	if t, ok := s.themes[id]; ok && t.TenantID == tenantID && t.RestaurantID == restaurantID {
		s.activated = append(s.activated, id)
	}
	return nil
}

// Caller: tenant 3, restaurant 5, user 7.
func themeRequest(method, path, id, body string, h func(http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.SetPathValue("id", id)
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, &jwt.Claims{UserID: 7, TenantID: 3, RestaurantID: 5})
	ctx = context.WithValue(ctx, middleware.TenantContextKey, int64(3))
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(5))
	rec := httptest.NewRecorder()
	h(rec, req.WithContext(ctx))
	return rec
}

func newThemeFixture() (*AdminThemeHandlerV2, *themeSvcStub) {
	stub := &themeSvcStub{
		themes: map[int64]*domain.Theme{
			1: {ID: 1, TenantID: 3, RestaurantID: 5, Name: "mine", Slug: "mine"},
			2: {ID: 2, TenantID: 9, RestaurantID: 11, Name: "other tenant", Slug: "other"},
			3: {ID: 3, TenantID: 3, RestaurantID: 6, Name: "same tenant, other restaurant", Slug: "sibling"},
		},
		presets: map[int64]*domain.Theme{},
	}
	return &AdminThemeHandlerV2{service: stub}, stub
}

// SEC-F: themes of another tenant or another restaurant must be invisible and immutable.
func TestThemeIsolation_ForeignThemesAre404(t *testing.T) {
	h, stub := newThemeFixture()
	for _, foreign := range []string{"2", "3"} {
		if rec := themeRequest("GET", "/api/v1/admin/themes/"+foreign, foreign, "", h.GetTheme); rec.Code != http.StatusNotFound {
			t.Errorf("GET theme %s: want 404, got %d", foreign, rec.Code)
		}
		if rec := themeRequest("PUT", "/api/v1/admin/themes/"+foreign, foreign, `{"name":"pwned"}`, h.UpdateTheme); rec.Code != http.StatusNotFound {
			t.Errorf("PUT theme %s: want 404, got %d", foreign, rec.Code)
		}
		if rec := themeRequest("DELETE", "/api/v1/admin/themes/"+foreign, foreign, "", h.DeleteTheme); rec.Code != http.StatusNotFound {
			t.Errorf("DELETE theme %s: want 404, got %d", foreign, rec.Code)
		}
		if rec := themeRequest("POST", "/api/v1/admin/themes/"+foreign+"/activate", foreign, "", h.ActivateTheme); rec.Code != http.StatusNotFound {
			t.Errorf("ACTIVATE theme %s: want 404, got %d", foreign, rec.Code)
		}
		if rec := themeRequest("POST", "/api/v1/admin/themes/"+foreign+"/duplicate", foreign, "", h.DuplicateTheme); rec.Code != http.StatusNotFound {
			t.Errorf("DUPLICATE theme %s: want 404, got %d", foreign, rec.Code)
		}
	}
	if len(stub.deleted) != 0 || len(stub.activated) != 0 {
		t.Fatalf("foreign themes must never reach delete/activate: deleted=%v activated=%v", stub.deleted, stub.activated)
	}
}

func TestThemeIsolation_MissingThemeLooksIdenticalToForeign(t *testing.T) {
	h, _ := newThemeFixture()
	missing := themeRequest("GET", "/api/v1/admin/themes/999", "999", "", h.GetTheme)
	foreign := themeRequest("GET", "/api/v1/admin/themes/2", "2", "", h.GetTheme)
	if missing.Code != foreign.Code || missing.Body.String() != foreign.Body.String() {
		t.Fatalf("responses must be indistinguishable: %d %q vs %d %q", missing.Code, missing.Body.String(), foreign.Code, foreign.Body.String())
	}
}

func TestThemeIsolation_OwnThemeStillWorks(t *testing.T) {
	h, stub := newThemeFixture()
	if rec := themeRequest("GET", "/api/v1/admin/themes/1", "1", "", h.GetTheme); rec.Code != http.StatusOK {
		t.Fatalf("own theme GET: want 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	if rec := themeRequest("POST", "/api/v1/admin/themes/1/activate", "1", "", h.ActivateTheme); rec.Code != http.StatusOK {
		t.Fatalf("own theme activate: want 200, got %d", rec.Code)
	}
	if rec := themeRequest("DELETE", "/api/v1/admin/themes/1", "1", "", h.DeleteTheme); rec.Code != http.StatusNoContent {
		t.Fatalf("own theme delete: want 204, got %d", rec.Code)
	}
	if len(stub.activated) != 1 || stub.activated[0] != 1 || len(stub.deleted) != 1 || stub.deleted[0] != 1 {
		t.Fatalf("service calls wrong: activated=%v deleted=%v", stub.activated, stub.deleted)
	}
}
