package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/database"
)

const notificationColumns = `id, tenant_id, restaurant_id, user_id, type, module, title, message, description,
	is_read, read_at, related_entity_type, related_entity_id, priority, action_url, action_label,
	icon_name, color, created_at, updated_at, expires_at`

// NotificationRepository handles notification data operations.
type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, notification *domain.Notification) (*domain.Notification, error) {
	if notification == nil {
		return nil, fmt.Errorf("notification is required")
	}
	err := database.WithTenantUserTx(ctx, r.db, int64(notification.TenantID), int64(notification.RestaurantID), int64(notification.UserID), func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `INSERT INTO notifications (
			tenant_id, restaurant_id, user_id, type, module, title, message, description,
			related_entity_type, related_entity_id, priority, action_url, action_label, icon_name, color, expires_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		RETURNING id, created_at, updated_at`,
			notification.TenantID, notification.RestaurantID, notification.UserID,
			notification.Type, notification.Module, notification.Title, notification.Message, notification.Description,
			notification.RelatedEntityType, notification.RelatedEntityID, notification.Priority,
			notification.ActionURL, notification.ActionLabel, notification.IconName, notification.Color, notification.ExpiresAt,
		).Scan(&notification.ID, &notification.CreatedAt, &notification.UpdatedAt)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	return notification, nil
}

func (r *NotificationRepository) GetNotification(ctx context.Context, userID, tenantID, restaurantID, notificationID int) (*domain.Notification, error) {
	var notification domain.Notification
	err := database.WithTenantUserTx(ctx, r.db, int64(tenantID), int64(restaurantID), int64(userID), func(tx *sql.Tx) error {
		var err error
		notification, err = scanNotification(tx.QueryRowContext(ctx, `SELECT `+notificationColumns+`
			FROM notifications WHERE id = $1 AND user_id = $2 AND tenant_id = $3 AND restaurant_id = $4`,
			notificationID, userID, tenantID, restaurantID))
		return err
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return &notification, nil
}

func (r *NotificationRepository) ListNotifications(ctx context.Context, userID, tenantID, restaurantID int, filters *domain.NotificationFilters) (*domain.NotificationListResponse, error) {
	if filters == nil {
		filters = &domain.NotificationFilters{Page: 1, PageSize: 20}
	}
	var response domain.NotificationListResponse
	err := database.WithTenantUserTx(ctx, r.db, int64(tenantID), int64(restaurantID), int64(userID), func(tx *sql.Tx) error {
		where := " WHERE user_id = $1 AND tenant_id = $2 AND restaurant_id = $3"
		args := []any{userID, tenantID, restaurantID}
		appendFilter := func(column, value string) {
			args = append(args, value)
			where += fmt.Sprintf(" AND %s = $%d", column, len(args))
		}
		if filters.Module != "" {
			appendFilter("module", filters.Module)
		}
		if filters.Type != "" {
			appendFilter("type", filters.Type)
		}
		if filters.IsRead != nil {
			args = append(args, *filters.IsRead)
			where += fmt.Sprintf(" AND is_read = $%d", len(args))
		}
		if filters.Priority != "" {
			appendFilter("priority", filters.Priority)
		}
		if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM notifications"+where, args...).Scan(&response.Total); err != nil {
			return fmt.Errorf("count notifications: %w", err)
		}

		unreadWhere := " WHERE user_id = $1 AND tenant_id = $2 AND restaurant_id = $3 AND is_read = false"
		unreadArgs := []any{userID, tenantID, restaurantID}
		if filters.Module != "" {
			unreadArgs = append(unreadArgs, filters.Module)
			unreadWhere += fmt.Sprintf(" AND module = $%d", len(unreadArgs))
		}
		if filters.Type != "" {
			unreadArgs = append(unreadArgs, filters.Type)
			unreadWhere += fmt.Sprintf(" AND type = $%d", len(unreadArgs))
		}
		if filters.Priority != "" {
			unreadArgs = append(unreadArgs, filters.Priority)
			unreadWhere += fmt.Sprintf(" AND priority = $%d", len(unreadArgs))
		}
		if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM notifications"+unreadWhere, unreadArgs...).Scan(&response.Unread); err != nil {
			return fmt.Errorf("count unread notifications: %w", err)
		}

		sortFields := map[string]string{"created_at": "created_at", "updated_at": "updated_at", "priority": "priority", "type": "type", "title": "title"}
		sortField := sortFields[filters.SortBy]
		if sortField == "" {
			sortField = "created_at"
		}
		sortOrder := strings.ToUpper(filters.SortOrder)
		if sortOrder != "ASC" && sortOrder != "DESC" {
			sortOrder = "DESC"
		}
		page, pageSize := filters.Page, filters.PageSize
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}
		args = append(args, pageSize, (page-1)*pageSize)
		rows, err := tx.QueryContext(ctx, "SELECT "+notificationColumns+" FROM notifications"+where+
			fmt.Sprintf(" ORDER BY %s %s, id DESC LIMIT $%d OFFSET $%d", sortField, sortOrder, len(args)-1, len(args)), args...)
		if err != nil {
			return fmt.Errorf("list notifications: %w", err)
		}
		defer rows.Close()
		response.Notifications = []domain.Notification{}
		for rows.Next() {
			notification, err := scanNotification(rows)
			if err != nil {
				return fmt.Errorf("scan notification: %w", err)
			}
			response.Notifications = append(response.Notifications, notification)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("iterate notifications: %w", err)
		}
		response.Page, response.PageSize = page, pageSize
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, userID, tenantID, restaurantID, notificationID int) error {
	return r.updateReadState(ctx, userID, tenantID, restaurantID, notificationID, true)
}

func (r *NotificationRepository) MarkAsUnread(ctx context.Context, userID, tenantID, restaurantID, notificationID int) error {
	return r.updateReadState(ctx, userID, tenantID, restaurantID, notificationID, false)
}

func (r *NotificationRepository) updateReadState(ctx context.Context, userID, tenantID, restaurantID, notificationID int, read bool) error {
	err := database.WithTenantUserTx(ctx, r.db, int64(tenantID), int64(restaurantID), int64(userID), func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `UPDATE notifications SET is_read = $1,
			read_at = CASE WHEN $1 THEN CURRENT_TIMESTAMP ELSE NULL END
			WHERE id = $2 AND user_id = $3 AND tenant_id = $4 AND restaurant_id = $5`,
			read, notificationID, userID, tenantID, restaurantID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return fmt.Errorf("notification not found")
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update notification read state: %w", err)
	}
	return nil
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID, tenantID, restaurantID int) error {
	err := database.WithTenantUserTx(ctx, r.db, int64(tenantID), int64(restaurantID), int64(userID), func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE notifications SET is_read = true, read_at = CURRENT_TIMESTAMP
			WHERE user_id = $1 AND tenant_id = $2 AND restaurant_id = $3 AND is_read = false`, userID, tenantID, restaurantID)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) DeleteNotification(ctx context.Context, userID, tenantID, restaurantID, notificationID int) error {
	err := database.WithTenantUserTx(ctx, r.db, int64(tenantID), int64(restaurantID), int64(userID), func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `DELETE FROM notifications
			WHERE id = $1 AND user_id = $2 AND tenant_id = $3 AND restaurant_id = $4`, notificationID, userID, tenantID, restaurantID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return fmt.Errorf("notification not found")
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID, tenantID, restaurantID int) (int, error) {
	var count int
	err := database.WithTenantUserTx(ctx, r.db, int64(tenantID), int64(restaurantID), int64(userID), func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND tenant_id = $2
			AND restaurant_id = $3 AND is_read = false`, userID, tenantID, restaurantID).Scan(&count)
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

func (r *NotificationRepository) GetNotificationStats(ctx context.Context, userID, tenantID, restaurantID int) (*domain.NotificationStats, error) {
	stats := &domain.NotificationStats{ByModule: make(map[string]int), ByPriority: make(map[string]int)}
	err := database.WithTenantUserTx(ctx, r.db, int64(tenantID), int64(restaurantID), int64(userID), func(tx *sql.Tx) error {
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*), COUNT(*) FILTER (WHERE is_read = false),
			COUNT(*) FILTER (WHERE is_read = true) FROM notifications
			WHERE user_id = $1 AND tenant_id = $2 AND restaurant_id = $3`, userID, tenantID, restaurantID).
			Scan(&stats.Total, &stats.Unread, &stats.Read); err != nil {
			return fmt.Errorf("count notification stats: %w", err)
		}
		for _, grouping := range []struct {
			column string
			target map[string]int
		}{{"module", stats.ByModule}, {"priority", stats.ByPriority}} {
			rows, err := tx.QueryContext(ctx, "SELECT "+grouping.column+", COUNT(*) FROM notifications WHERE user_id = $1 AND tenant_id = $2 AND restaurant_id = $3 GROUP BY "+grouping.column,
				userID, tenantID, restaurantID)
			if err != nil {
				return fmt.Errorf("query notification %s stats: %w", grouping.column, err)
			}
			for rows.Next() {
				var key string
				var count int
				if err := rows.Scan(&key, &count); err != nil {
					rows.Close()
					return fmt.Errorf("scan notification %s stats: %w", grouping.column, err)
				}
				grouping.target[key] = count
			}
			if err := rows.Err(); err != nil {
				rows.Close()
				return fmt.Errorf("iterate notification %s stats: %w", grouping.column, err)
			}
			if err := rows.Close(); err != nil {
				return fmt.Errorf("close notification %s stats: %w", grouping.column, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get notification stats: %w", err)
	}
	return stats, nil
}

func scanNotification(row interface{ Scan(...any) error }) (domain.Notification, error) {
	var notification domain.Notification
	var nType, nModule, priority string
	err := row.Scan(&notification.ID, &notification.TenantID, &notification.RestaurantID, &notification.UserID,
		&nType, &nModule, &notification.Title, &notification.Message, &notification.Description,
		&notification.IsRead, &notification.ReadAt, &notification.RelatedEntityType, &notification.RelatedEntityID,
		&priority, &notification.ActionURL, &notification.ActionLabel, &notification.IconName, &notification.Color,
		&notification.CreatedAt, &notification.UpdatedAt, &notification.ExpiresAt)
	if err != nil {
		return domain.Notification{}, err
	}
	notification.Type = domain.NotificationType(nType)
	notification.Module = domain.NotificationModule(nModule)
	notification.Priority = domain.NotificationPriority(priority)
	return notification, nil
}
