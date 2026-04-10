package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// UserSettingsHandler handles HTTP requests for user settings
type UserSettingsHandler struct {
	settingsRepo *repository.UserSettingsRepository
	userRepo     *repository.UserRepository
}

// NewUserSettingsHandler creates new user settings handler
func NewUserSettingsHandler(settingsRepo *repository.UserSettingsRepository, userRepo *repository.UserRepository) *UserSettingsHandler {
	return &UserSettingsHandler{
		settingsRepo: settingsRepo,
		userRepo:     userRepo,
	}
}

// GetUserSettings retrieves complete user settings
// GET /api/v1/users/:userId/settings
func (h *UserSettingsHandler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	log.Println("[GetUserSettings] Starting...")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		log.Println("[GetUserSettings] ERROR: Claims is nil - user not authenticated")
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDStr := r.PathValue("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Printf("[GetUserSettings] ERROR: Invalid user ID: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Authorization: users can only view their own settings
	if userID != int64(claims.UserID) {
		log.Printf("[GetUserSettings] ERROR: User %d tried to access user %d settings", claims.UserID, userID)
		respondError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	settings, err := h.settingsRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		log.Printf("[GetUserSettings] ERROR: Failed to get settings: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch settings")
		return
	}

	log.Printf("[GetUserSettings] Settings retrieved for user %d", userID)
	respondJSON(w, http.StatusOK, settings)
}

// UpdateLanguage updates user's preferred language
// PUT /api/v1/users/:userId/settings/language
func (h *UserSettingsHandler) UpdateLanguage(w http.ResponseWriter, r *http.Request) {
	log.Println("[UpdateLanguage] Starting...")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		log.Println("[UpdateLanguage] ERROR: Claims is nil")
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDStr := r.PathValue("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if userID != int64(claims.UserID) {
		respondError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	var req struct {
		Language string `json:"language"`
	}

	if err := decodeJSONStrict(r, &req); err != nil {
		log.Printf("[UpdateLanguage] ERROR: Failed to decode request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate language
	if req.Language != "en" && req.Language != "ar" {
		log.Printf("[UpdateLanguage] ERROR: Invalid language: %s", req.Language)
		respondError(w, http.StatusBadRequest, "Language must be 'en' or 'ar'")
		return
	}

	settings, err := h.settingsRepo.UpdateLanguage(r.Context(), userID, req.Language)
	if err != nil {
		log.Printf("[UpdateLanguage] ERROR: Failed to update language: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to update language")
		return
	}

	// Log audit change
	h.settingsRepo.LogAuditChange(r.Context(), userID, int64(claims.UserID), "language", nil, &req.Language, getClientIP(r))

	log.Printf("[UpdateLanguage] Language updated for user %d to %s", userID, req.Language)
	respondJSON(w, http.StatusOK, settings)
}

// UpdateTheme updates user's theme preference
// PUT /api/v1/users/:userId/settings/theme
func (h *UserSettingsHandler) UpdateTheme(w http.ResponseWriter, r *http.Request) {
	log.Println("[UpdateTheme] Starting...")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDStr := r.PathValue("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if userID != int64(claims.UserID) {
		respondError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	var req struct {
		Theme string `json:"theme"`
	}

	if err := decodeJSONStrict(r, &req); err != nil {
		log.Printf("[UpdateTheme] ERROR: Failed to decode request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate theme
	if req.Theme != "light" && req.Theme != "dark" && req.Theme != "system" {
		respondError(w, http.StatusBadRequest, "Theme must be 'light', 'dark', or 'system'")
		return
	}

	settings, err := h.settingsRepo.UpdateTheme(r.Context(), userID, req.Theme)
	if err != nil {
		log.Printf("[UpdateTheme] ERROR: Failed to update theme: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to update theme")
		return
	}

	h.settingsRepo.LogAuditChange(r.Context(), userID, int64(claims.UserID), "theme", nil, &req.Theme, getClientIP(r))

	log.Printf("[UpdateTheme] Theme updated for user %d to %s", userID, req.Theme)
	respondJSON(w, http.StatusOK, settings)
}

// UpdateColors updates dashboard color scheme
// PUT /api/v1/users/:userId/settings/theme-colors
func (h *UserSettingsHandler) UpdateColors(w http.ResponseWriter, r *http.Request) {
	log.Println("[UpdateColors] Starting...")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDStr := r.PathValue("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if userID != int64(claims.UserID) {
		respondError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	var req struct {
		PrimaryColor   string `json:"primary_color"`
		SecondaryColor string `json:"secondary_color"`
		AccentColor    string `json:"accent_color"`
	}

	if err := decodeJSONStrict(r, &req); err != nil {
		log.Printf("[UpdateColors] ERROR: Failed to decode request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate colors are provided
	if req.PrimaryColor == "" || req.SecondaryColor == "" || req.AccentColor == "" {
		respondError(w, http.StatusBadRequest, "All colors are required")
		return
	}

	settings, err := h.settingsRepo.UpdateColors(r.Context(), userID, req.PrimaryColor, req.SecondaryColor, req.AccentColor)
	if err != nil {
		log.Printf("[UpdateColors] ERROR: Failed to update colors: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to update colors")
		return
	}

	h.settingsRepo.LogAuditChange(r.Context(), userID, int64(claims.UserID), "colors", nil,
		stringPtr(fmt.Sprintf("%s|%s|%s", req.PrimaryColor, req.SecondaryColor, req.AccentColor)), getClientIP(r))

	log.Printf("[UpdateColors] Colors updated for user %d", userID)
	respondJSON(w, http.StatusOK, settings)
}

// ChangePassword changes user's password
// POST /api/v1/users/:userId/settings/change-password
func (h *UserSettingsHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	log.Println("[ChangePassword] Starting...")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDStr := r.PathValue("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if userID != int64(claims.UserID) {
		respondError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := decodeJSONStrict(r, &req); err != nil {
		log.Printf("[ChangePassword] ERROR: Failed to decode request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate inputs
	if req.CurrentPassword == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		respondError(w, http.StatusBadRequest, "All password fields are required")
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		log.Printf("[ChangePassword] ERROR: Passwords do not match for user %d", userID)
		respondError(w, http.StatusBadRequest, "New passwords do not match")
		return
	}

	if len(req.NewPassword) < 8 {
		respondError(w, http.StatusBadRequest, "New password must be at least 8 characters")
		return
	}

	// Get user to verify current password
	user, err := h.userRepo.GetByID(r.Context(), int(userID))
	if err != nil {
		log.Printf("[ChangePassword] ERROR: User not found: %v", err)
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Verify current password
	if !user.VerifyPassword(req.CurrentPassword) {
		log.Printf("[ChangePassword] ERROR: Invalid current password for user %d", userID)
		respondError(w, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	// Hash new password
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		log.Printf("[ChangePassword] ERROR: Failed to hash password: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Update password in user repo
	if err := h.userRepo.UpdatePassword(r.Context(), int(userID), hashedPassword); err != nil {
		log.Printf("[ChangePassword] ERROR: Failed to update password: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	// Record password change in history
	if err := h.settingsRepo.RecordPasswordChange(r.Context(), userID, hashedPassword, getClientIP(r)); err != nil {
		log.Printf("[ChangePassword] ERROR: Failed to record password change: %v", err)
		// Don't fail the request, just log it
	}

	log.Printf("[ChangePassword] Password changed for user %d", userID)
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Password updated successfully",
	})
}

// GetUserProfile retrieves user's basic profile info
// GET /api/v1/users/:userId/profile
func (h *UserSettingsHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("[GetUserProfile] Starting...")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDStr := r.PathValue("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if userID != int64(claims.UserID) {
		respondError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), int(userID))
	if err != nil {
		log.Printf("[GetUserProfile] ERROR: Failed to get user: %v", err)
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Create profile with proper pointer handling
	var avatarPtr *string
	if user.AvatarURL != "" {
		avatarPtr = &user.AvatarURL
	}
	var phonePtr *string
	if user.Phone != "" {
		phonePtr = &user.Phone
	}

	profile := domain.UserProfile{
		ID:        int64(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		AvatarURL: avatarPtr,
		Phone:     phonePtr,
		CreatedAt: user.CreatedAt,
	}

	log.Printf("[GetUserProfile] Profile retrieved for user %d", userID)
	respondJSON(w, http.StatusOK, profile)
}

// UpdateProfile updates user's profile information
// PUT /api/v1/users/:userId/profile
func (h *UserSettingsHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("[UpdateProfile] Starting...")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDStr := r.PathValue("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if userID != int64(claims.UserID) {
		respondError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	var req struct {
		Name      string  `json:"name"`
		AvatarURL *string `json:"avatar_url"`
	}

	if err := decodeJSONStrict(r, &req); err != nil {
		log.Printf("[UpdateProfile] ERROR: Failed to decode request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update profile
	if err := h.userRepo.UpdateProfile(r.Context(), int(userID), req.Name, req.AvatarURL); err != nil {
		log.Printf("[UpdateProfile] ERROR: Failed to update profile: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), int(userID))
	if err != nil {
		log.Printf("[UpdateProfile] ERROR: Failed to get updated user: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve updated profile")
		return
	}

	// Create profile with proper pointer handling
	var avatarPtr *string
	if user.AvatarURL != "" {
		avatarPtr = &user.AvatarURL
	}
	var phonePtr *string
	if user.Phone != "" {
		phonePtr = &user.Phone
	}

	profile := domain.UserProfile{
		ID:        int64(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		AvatarURL: avatarPtr,
		Phone:     phonePtr,
		CreatedAt: user.CreatedAt,
	}

	log.Printf("[UpdateProfile] Profile updated for user %d", userID)
	respondJSON(w, http.StatusOK, profile)
}

// ListUsers retrieves all users for the tenant
// GET /api/v1/users
func (h *UserSettingsHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("[ListUsers] Starting...")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	tenantID := middleware.GetTenantID(r)
	log.Printf("[ListUsers] Listing users for tenant %d", tenantID)

	users, err := h.userRepo.ListByTenant(r.Context(), tenantID)
	if err != nil {
		log.Printf("[ListUsers] ERROR: Failed to list users: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	// Transform users to response format (without password hashes)
	type UserResponse struct {
		ID        int64   `json:"id"`
		Name      string  `json:"name"`
		Email     string  `json:"email"`
		Role      string  `json:"role"`
		Phone     *string `json:"phone,omitempty"`
		AvatarURL *string `json:"avatar_url,omitempty"`
		Status    string  `json:"status"`
		CreatedAt string  `json:"created_at"`
	}

	response := make([]UserResponse, 0, len(users))
	for _, user := range users {
		userResp := UserResponse{
			ID:        int64(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if user.Phone != "" {
			userResp.Phone = &user.Phone
		}
		if user.AvatarURL != "" {
			userResp.AvatarURL = &user.AvatarURL
		}

		response = append(response, userResp)
	}

	log.Printf("[ListUsers] Returning %d users for tenant %d", len(response), tenantID)
	respondJSON(w, http.StatusOK, response)
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}
