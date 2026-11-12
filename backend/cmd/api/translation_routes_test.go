package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"pos-saas/internal/pkg/jwt"
)

type translationRouteStub struct{}

func (translationRouteStub) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (translationRouteStub) GetSupportedLanguages(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (translationRouteStub) GetTranslations(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (translationRouteStub) GetTranslationKey(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (translationRouteStub) UploadTranslations(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func TestTranslationRouteRegistrationPreservesPublicAndSuperAdminBoundaries(t *testing.T) {
	tokenService, err := jwt.NewSuperAdminTokenService("translation-route-test-secret-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create super-admin token service: %v", err)
	}
	mux := http.NewServeMux()
	registerTranslationRoutes(mux, translationRouteStub{}, tokenService)

	public := httptest.NewRecorder()
	mux.ServeHTTP(public, httptest.NewRequest(http.MethodGet, "/api/v1/translations/health", nil))
	if public.Code != http.StatusNoContent {
		t.Fatalf("public health route status = %d, want %d", public.Code, http.StatusNoContent)
	}

	unauthorized := httptest.NewRecorder()
	mux.ServeHTTP(unauthorized, httptest.NewRequest(http.MethodPost, "/api/v1/translations/en/upload", nil))
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("translation upload without credentials status = %d, want %d", unauthorized.Code, http.StatusUnauthorized)
	}

	accessToken, _, err := tokenService.GenerateSuperAdminToken(7, "translation-admin@example.test", 3)
	if err != nil {
		t.Fatalf("generate test super-admin token: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/translations/en/upload", nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	authorized := httptest.NewRecorder()
	mux.ServeHTTP(authorized, request)
	if authorized.Code != http.StatusNoContent {
		t.Fatalf("translation upload with super-admin token status = %d, want %d", authorized.Code, http.StatusNoContent)
	}
}
