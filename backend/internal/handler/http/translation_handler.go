package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type TranslationHandler struct {
	translations map[string]map[string]interface{}
}

func NewTranslationHandler() *TranslationHandler {
	return &TranslationHandler{
		translations: make(map[string]map[string]interface{}),
	}
}

// LoadTranslations loads translation files from the frontend
// In a real application, these would be loaded from files or a database
func (h *TranslationHandler) LoadTranslations() error {
	// For now, we'll return an empty map
	// In production, translations should be loaded from the frontend's i18n/messages directory
	h.translations = map[string]map[string]interface{}{
		"en": {},
		"ar": {},
	}
	return nil
}

// GetTranslations returns all translations for a specific language
func (h *TranslationHandler) GetTranslations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract language from URL path: /api/v1/translations/{language}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/translations"), "/")
	if len(parts) < 2 {
		respondError(w, http.StatusBadRequest, "Language parameter required")
		return
	}

	language := parts[1]
	if language == "" {
		respondError(w, http.StatusBadRequest, "Language parameter required")
		return
	}

	// Validate language
	if language != "en" && language != "ar" {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Unsupported language: %s. Supported languages: en, ar", language))
		return
	}

	log.Printf("[GetTranslations] Fetching translations for language: %s", language)

	// Load translations if not already loaded
	if len(h.translations) == 0 {
		if err := h.LoadTranslations(); err != nil {
			log.Printf("[GetTranslations] ERROR: Failed to load translations: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to load translations")
			return
		}
	}

	// Return translations for the specified language
	translations := h.translations[language]
	if translations == nil {
		translations = make(map[string]interface{})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"language":     language,
		"translations": translations,
	})
}

// GetSupportedLanguages returns list of supported languages
func (h *TranslationHandler) GetSupportedLanguages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	log.Printf("[GetSupportedLanguages] Fetching supported languages")

	languages := []map[string]string{
		{
			"code": "en",
			"name": "English",
		},
		{
			"code": "ar",
			"name": "العربية",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"languages": languages,
	})
}

// GetTranslationKey returns a specific translation key for a language
func (h *TranslationHandler) GetTranslationKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract parameters from URL: /api/v1/translations/{language}/{namespace}/{key}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/translations"), "/")
	if len(parts) < 4 {
		respondError(w, http.StatusBadRequest, "Missing required parameters: language, namespace, key")
		return
	}

	language := parts[1]
	namespace := parts[2]
	key := parts[3]

	if language == "" || namespace == "" || key == "" {
		respondError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	// Validate language
	if language != "en" && language != "ar" {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Unsupported language: %s", language))
		return
	}

	log.Printf("[GetTranslationKey] Fetching translation: language=%s, namespace=%s, key=%s", language, namespace, key)

	// Load translations if not already loaded
	if len(h.translations) == 0 {
		if err := h.LoadTranslations(); err != nil {
			log.Printf("[GetTranslationKey] ERROR: Failed to load translations: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to load translations")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"language":  language,
		"namespace": namespace,
		"key":       key,
		"value":     "", // In production, fetch from translations map
	})
}

// UploadTranslations allows uploading new translations (Admin only)
func (h *TranslationHandler) UploadTranslations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Validate file upload
	file, _, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	language := r.FormValue("language")
	if language == "" {
		respondError(w, http.StatusBadRequest, "Language parameter required")
		return
	}

	if language != "en" && language != "ar" {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Unsupported language: %s", language))
		return
	}

	// Read and parse JSON from file
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[UploadTranslations] ERROR: Failed to read file: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to read file")
		return
	}

	var translations map[string]interface{}
	if err := json.Unmarshal(fileContent, &translations); err != nil {
		log.Printf("[UploadTranslations] ERROR: Failed to parse JSON: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Store translations
	h.translations[language] = translations

	log.Printf("[UploadTranslations] Successfully uploaded translations for language: %s", language)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Translations uploaded successfully",
		"language": language,
	})
}

// Health check endpoint for translation service
func (h *TranslationHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":             "ok",
		"service":            "translation-service",
		"supportedLanguages": []string{"en", "ar"},
	})
}
