package usecase

import (
	"context"
	"fmt"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

// NotificationUseCase handles notification business logic
type NotificationUseCase struct {
	repo *repository.NotificationRepository
}

// NewNotificationUseCase creates new notification use case
func NewNotificationUseCase(repo *repository.NotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{repo: repo}
}

// CreateNotification creates a new notification
func (uc *NotificationUseCase) CreateNotification(ctx context.Context, notification *domain.Notification) (*domain.Notification, error) {
	if notification.Title == "" {
		return nil, fmt.Errorf("notification title is required")
	}

	if notification.Message == "" {
		return nil, fmt.Errorf("notification message is required")
	}

	return uc.repo.CreateNotification(ctx, notification)
}

// GetNotification retrieves a single notification
func (uc *NotificationUseCase) GetNotification(ctx context.Context, userID, tenantID, restaurantID, notificationID int) (*domain.Notification, error) {
	return uc.repo.GetNotification(ctx, userID, tenantID, restaurantID, notificationID)
}

// ListNotifications retrieves notifications with filters
func (uc *NotificationUseCase) ListNotifications(ctx context.Context, userID, tenantID, restaurantID int, filters *domain.NotificationFilters) (*domain.NotificationListResponse, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}
	if filters.PageSize > 100 {
		filters.PageSize = 100
	}

	if filters.SortBy == "" {
		filters.SortBy = "created_at"
	}
	if filters.SortOrder == "" {
		filters.SortOrder = "DESC"
	}

	return uc.repo.ListNotifications(ctx, userID, tenantID, restaurantID, filters)
}

// MarkAsRead marks a notification as read
func (uc *NotificationUseCase) MarkAsRead(ctx context.Context, userID, tenantID, restaurantID, notificationID int) error {
	return uc.repo.MarkAsRead(ctx, userID, tenantID, restaurantID, notificationID)
}

// MarkAsUnread marks a notification as unread
func (uc *NotificationUseCase) MarkAsUnread(ctx context.Context, userID, tenantID, restaurantID, notificationID int) error {
	return uc.repo.MarkAsUnread(ctx, userID, tenantID, restaurantID, notificationID)
}

// MarkAllAsRead marks all unread notifications as read
func (uc *NotificationUseCase) MarkAllAsRead(ctx context.Context, userID, tenantID, restaurantID int) error {
	return uc.repo.MarkAllAsRead(ctx, userID, tenantID, restaurantID)
}

// DeleteNotification deletes a notification
func (uc *NotificationUseCase) DeleteNotification(ctx context.Context, userID, tenantID, restaurantID, notificationID int) error {
	return uc.repo.DeleteNotification(ctx, userID, tenantID, restaurantID, notificationID)
}

// GetUnreadCount gets count of unread notifications
func (uc *NotificationUseCase) GetUnreadCount(ctx context.Context, userID, tenantID, restaurantID int) (int, error) {
	return uc.repo.GetUnreadCount(ctx, userID, tenantID, restaurantID)
}

// GetNotificationStats gets notification statistics
func (uc *NotificationUseCase) GetNotificationStats(ctx context.Context, userID, tenantID, restaurantID int) (*domain.NotificationStats, error) {
	return uc.repo.GetNotificationStats(ctx, userID, tenantID, restaurantID)
}
