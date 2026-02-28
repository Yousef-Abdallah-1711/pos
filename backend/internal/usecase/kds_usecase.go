package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

// KDSBroadcaster publishes KDS websocket events.
type KDSBroadcaster interface {
	BroadcastKDSUpdate(ctx context.Context, tenantID, restaurantID int64, messageType string, payload interface{}) error
}

// KDSUseCase handles KDS business logic.
type KDSUseCase struct {
	repo        *repository.KDSRepository
	broadcaster KDSBroadcaster
}

// NewKDSUseCase creates a new KDS use case.
func NewKDSUseCase(repo *repository.KDSRepository) *KDSUseCase {
	return &KDSUseCase{repo: repo}
}

// SetBroadcaster connects websocket broadcasting to KDS events.
func (uc *KDSUseCase) SetBroadcaster(broadcaster KDSBroadcaster) {
	uc.broadcaster = broadcaster
}

// GetQueue fetches queue entries filtered by optional status filters.
func (uc *KDSUseCase) GetQueue(ctx context.Context, tenantID, restaurantID int64, filters domain.KDSFilters) ([]domain.KDSOrder, error) {
	if uc.repo == nil {
		return nil, errors.New("kds repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}

	statuses := make([]string, 0, len(filters.Statuses))
	if len(filters.Statuses) == 0 && !filters.IncludeCompleted {
		statuses = append(statuses, "new", "in_progress", "ready", "recalled")
	} else {
		for _, status := range filters.Statuses {
			value := strings.TrimSpace(string(status))
			if value == "" {
				continue
			}
			if !domain.ValidKDSStatus(value) {
				return nil, fmt.Errorf("invalid kds status: %s", value)
			}
			statuses = append(statuses, value)
		}
	}

	orders, err := uc.repo.GetQueue(ctx, tenantID, restaurantID, statuses)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for i := range orders {
		if orders[i].CompletedAt != nil {
			orders[i].ElapsedSeconds = int64(orders[i].CompletedAt.Sub(orders[i].ReceivedAt).Seconds())
		} else {
			orders[i].ElapsedSeconds = int64(now.Sub(orders[i].ReceivedAt).Seconds())
		}
		if orders[i].ElapsedSeconds < 0 {
			orders[i].ElapsedSeconds = 0
		}
	}

	return orders, nil
}

// BumpOrder advances a queue entry to the next status.
func (uc *KDSUseCase) BumpOrder(ctx context.Context, tenantID, restaurantID, queueID, userID int64) (*domain.KDSOrder, error) {
	if uc.repo == nil {
		return nil, errors.New("kds repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}
	if queueID <= 0 {
		return nil, errors.New("invalid queue id")
	}
	if userID <= 0 {
		return nil, errors.New("invalid user context")
	}

	order, err := uc.repo.BumpOrder(ctx, tenantID, restaurantID, queueID, userID)
	if err != nil {
		return nil, err
	}

	if order.CompletedAt != nil {
		order.ElapsedSeconds = int64(order.CompletedAt.Sub(order.ReceivedAt).Seconds())
	} else {
		order.ElapsedSeconds = int64(time.Since(order.ReceivedAt).Seconds())
	}
	if order.ElapsedSeconds < 0 {
		order.ElapsedSeconds = 0
	}

	uc.broadcast(ctx, tenantID, restaurantID, "order_bumped", order)
	uc.broadcast(ctx, tenantID, restaurantID, "order_updated", order)

	return order, nil
}

// RecallOrder reopens a completed queue entry.
func (uc *KDSUseCase) RecallOrder(ctx context.Context, tenantID, restaurantID, queueID int64) (*domain.KDSOrder, error) {
	if uc.repo == nil {
		return nil, errors.New("kds repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}
	if queueID <= 0 {
		return nil, errors.New("invalid queue id")
	}

	order, err := uc.repo.RecallOrder(ctx, tenantID, restaurantID, queueID)
	if err != nil {
		return nil, err
	}

	order.ElapsedSeconds = int64(time.Since(order.ReceivedAt).Seconds())
	if order.ElapsedSeconds < 0 {
		order.ElapsedSeconds = 0
	}

	uc.broadcast(ctx, tenantID, restaurantID, "order_recalled", order)
	uc.broadcast(ctx, tenantID, restaurantID, "order_updated", order)

	return order, nil
}

// GetStats returns queue statistics.
func (uc *KDSUseCase) GetStats(ctx context.Context, tenantID, restaurantID int64) (*domain.KDSStats, error) {
	if uc.repo == nil {
		return nil, errors.New("kds repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}
	return uc.repo.GetStats(ctx, tenantID, restaurantID)
}

// GetAnalytics returns KDS preparation analytics and workload patterns.
func (uc *KDSUseCase) GetAnalytics(ctx context.Context, tenantID, restaurantID int64) (*domain.KDSAnalytics, error) {
	if uc.repo == nil {
		return nil, errors.New("kds repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}
	return uc.repo.GetAnalytics(ctx, tenantID, restaurantID)
}

// GetSettings returns KDS settings for the current restaurant.
func (uc *KDSUseCase) GetSettings(ctx context.Context, tenantID, restaurantID int64) (*domain.KDSSettings, error) {
	if uc.repo == nil {
		return nil, errors.New("kds repository is not configured")
	}
	if tenantID <= 0 || restaurantID <= 0 {
		return nil, errors.New("invalid tenant context")
	}
	return uc.repo.GetSettings(ctx, tenantID, restaurantID)
}

// UpdateSettings validates and persists KDS settings.
func (uc *KDSUseCase) UpdateSettings(ctx context.Context, settings *domain.KDSSettings) error {
	if uc.repo == nil {
		return errors.New("kds repository is not configured")
	}
	if settings == nil {
		return errors.New("settings payload is required")
	}
	if settings.TenantID <= 0 || settings.RestaurantID <= 0 {
		return errors.New("invalid tenant context")
	}
	if settings.AutoBumpMinutes < 0 {
		return errors.New("auto_bump_minutes must be >= 0")
	}
	switch settings.DisplayDensity {
	case "compact", "comfortable", "spacious":
	default:
		return errors.New("display_density must be compact, comfortable, or spacious")
	}
	if settings.ColorWarningMinutes <= 0 {
		return errors.New("color_warning_minutes must be > 0")
	}
	if settings.ColorDangerMinutes <= settings.ColorWarningMinutes {
		return errors.New("color_danger_minutes must be greater than color_warning_minutes")
	}

	return uc.repo.UpdateSettings(ctx, settings)
}

func (uc *KDSUseCase) broadcast(ctx context.Context, tenantID, restaurantID int64, messageType string, payload interface{}) {
	if uc.broadcaster == nil {
		return
	}
	if err := uc.broadcaster.BroadcastKDSUpdate(ctx, tenantID, restaurantID, messageType, payload); err != nil {
		log.Printf("[KDS] broadcast failed (%s): tenant=%d restaurant=%d err=%v", messageType, tenantID, restaurantID, err)
	}
}
