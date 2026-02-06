package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/auth"
	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/usecase"
)

// NotificationHandler handles HTTP requests for notifications
type NotificationHandler struct {
	uc              *usecase.NotificationUseCase
	tenantValidator *auth.TenantValidator
}

// NewNotificationHandler creates new notification handler
func NewNotificationHandler(uc *usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		uc:              uc,
		tenantValidator: auth.NewTenantValidator(),
	}
}

// checkNotificationPermission verifies user has required permission for notification operations
func (h *NotificationHandler) checkNotificationPermission(
	r *http.Request,
	requiredLevel string,
) (int64, error) {
	// Step 1: Validate tenant context
	tenantID, err := h.tenantValidator.ValidateRequest(r)
	if err != nil {
		log.Printf("[NOTIFICATION HANDLER] Permission check failed: %v", err)
		return 0, fmt.Errorf("forbidden")
	}

	return tenantID, nil
}

// CreateNotification creates a new notification
// POST /api/v1/notifications
func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	var req domain.CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	notification := &domain.Notification{
		TenantID:          claims.TenantID,
		RestaurantID:      claims.RestaurantID,
		UserID:            req.UserID,
		Type:              req.Type,
		Module:            req.Module,
		Title:             req.Title,
		Message:           req.Message,
		Description:       req.Description,
		RelatedEntityType: req.RelatedEntityType,
		RelatedEntityID:   req.RelatedEntityID,
		Priority:          req.Priority,
		ActionURL:         req.ActionURL,
		ActionLabel:       req.ActionLabel,
		IconName:          req.IconName,
		Color:             req.Color,
	}

	created, err := h.uc.CreateNotification(notification)
	if err != nil {
		fmt.Printf("ERROR: CreateNotification failed: %v\n", err)
		respondError(w, http.StatusInternalServerError, "Failed to create notification")
		return
	}

	respondJSON(w, http.StatusCreated, created)
}

// ListNotifications retrieves notifications with filters
// GET /api/v1/notifications?module=products&type=low_stock&is_read=false&sort=created_at&page=1&limit=20
func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	// Parse query parameters
	module := r.URL.Query().Get("module")
	nType := r.URL.Query().Get("type")
	isReadStr := r.URL.Query().Get("is_read")
	priority := r.URL.Query().Get("priority")
	sortBy := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")
	page := 1
	pageSize := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := r.URL.Query().Get("limit"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	// Parse is_read filter
	var isRead *bool
	if isReadStr != "" {
		v := strings.ToLower(isReadStr) == "true"
		isRead = &v
	}

	filters := &domain.NotificationFilters{
		Module:    module,
		Type:      nType,
		IsRead:    isRead,
		Priority:  priority,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	response, err := h.uc.ListNotifications(claims.UserID, claims.TenantID, claims.RestaurantID, filters)
	if err != nil {
		fmt.Printf("ERROR: ListNotifications failed: %v\n", err)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve notifications")
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// GetNotification retrieves a single notification
// GET /api/v1/notifications/{id}
func (h *NotificationHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	notificationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	notification, err := h.uc.GetNotification(claims.UserID, notificationID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Notification not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to retrieve notification")
		}
		return
	}

	respondJSON(w, http.StatusOK, notification)
}

// MarkAsRead marks a notification as read
// POST /api/v1/notifications/{id}/read
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	notificationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := h.uc.MarkAsRead(claims.UserID, notificationID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Notification not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to mark notification as read")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

// MarkAsUnread marks a notification as unread
// POST /api/v1/notifications/{id}/unread
func (h *NotificationHandler) MarkAsUnread(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	notificationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := h.uc.MarkAsUnread(claims.UserID, notificationID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Notification not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to mark notification as unread")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Notification marked as unread"})
}

// MarkAllAsRead marks all unread notifications as read
// POST /api/v1/notifications/mark-all-read
func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	if err := h.uc.MarkAllAsRead(claims.UserID, claims.TenantID, claims.RestaurantID); err != nil {
		fmt.Printf("ERROR: MarkAllAsRead failed: %v\n", err)
		respondError(w, http.StatusInternalServerError, "Failed to mark notifications as read")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "All notifications marked as read"})
}

// DeleteNotification deletes a notification
// DELETE /api/v1/notifications/{id}
func (h *NotificationHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	notificationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := h.uc.DeleteNotification(claims.UserID, notificationID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Notification not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to delete notification")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Notification deleted"})
}

// GetStats gets notification statistics
// GET /api/v1/notifications/stats
func (h *NotificationHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	stats, err := h.uc.GetNotificationStats(claims.UserID, claims.TenantID, claims.RestaurantID)
	if err != nil {
		fmt.Printf("ERROR: GetStats failed: %v\n", err)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve notification stats")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

// GetUnreadCount gets count of unread notifications
// GET /api/v1/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)

	count, err := h.uc.GetUnreadCount(claims.UserID, claims.TenantID, claims.RestaurantID)
	if err != nil {
		fmt.Printf("ERROR: GetUnreadCount failed: %v\n", err)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve unread count")
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"unread": count})
}
