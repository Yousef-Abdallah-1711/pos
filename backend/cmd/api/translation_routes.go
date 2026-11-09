package main

import (
	"net/http"

	handler "pos-saas/internal/handler/http"
	"pos-saas/internal/middleware"
)

type translationRouteHandlers interface {
	HealthCheck(http.ResponseWriter, *http.Request)
	GetSupportedLanguages(http.ResponseWriter, *http.Request)
	GetTranslations(http.ResponseWriter, *http.Request)
	GetTranslationKey(http.ResponseWriter, *http.Request)
	UploadTranslations(http.ResponseWriter, *http.Request)
}

func registerTranslationRoutes(mux *http.ServeMux, translationHandler translationRouteHandlers, superAdminTokenService middleware.SuperAdminTokenValidator) {
	mux.HandleFunc("GET /api/v1/translations/health", translationHandler.HealthCheck)
	mux.HandleFunc("GET /api/v1/translations/languages", translationHandler.GetSupportedLanguages)
	mux.HandleFunc("GET /api/v1/translations/{language}", translationHandler.GetTranslations)
	mux.HandleFunc("GET /api/v1/translations/{language}/{namespace}/{key}", translationHandler.GetTranslationKey)
	mux.Handle("POST /api/v1/translations/{language}/upload", middleware.SuperAdminAuthMiddleware(superAdminTokenService)(http.HandlerFunc(translationHandler.UploadTranslations)))
}

var _ translationRouteHandlers = (*handler.TranslationHandler)(nil)
