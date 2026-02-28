package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/observability"
)

// KDSRepository handles kitchen display system data operations.
type KDSRepository struct {
	db *sql.DB
}

func defaultKDSSettings(tenantID, restaurantID int64) *domain.KDSSettings {
	return &domain.KDSSettings{
		TenantID:            tenantID,
		RestaurantID:        restaurantID,
		AutoBumpMinutes:     0,
		AlertSoundEnabled:   true,
		ShowAnalyticsPanel:  false,
		DisplayDensity:      "comfortable",
		ColorWarningMinutes: 10,
		ColorDangerMinutes:  20,
	}
}

func isMissingColumnError(err error, columnName string) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	col := strings.ToLower(columnName)
	return (strings.Contains(msg, "sqlstate 42703") || strings.Contains(msg, "does not exist")) &&
		strings.Contains(msg, col)
}

func isMissingRelationError(err error, relationName string) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	rel := strings.ToLower(relationName)
	return (strings.Contains(msg, "sqlstate 42p01") || strings.Contains(msg, "does not exist")) &&
		strings.Contains(msg, rel)
}

// NewKDSRepository creates a new KDS repository.
func NewKDSRepository(db *sql.DB) *KDSRepository {
	return &KDSRepository{db: db}
}

// EnqueueOrder inserts/refreshes a KDS queue entry for an order.
func (r *KDSRepository) EnqueueOrder(ctx context.Context, tenantID, restaurantID, orderID int64) (*domain.KDSOrder, error) {
	startedAt := time.Now()
	defer observability.RecordDB("kds.enqueue_order", time.Since(startedAt))

	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}

	orderNumber, customerName, notes, itemsJSON, err := r.buildQueuePayload(ctx, tenantID, restaurantID, orderID)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO kds_queue (
			tenant_id, restaurant_id, order_id, order_number, customer_name, status,
			priority, station_id, items, notes, received_at, started_at, completed_at, bumped_by
		) VALUES (
			$1, $2, $3, $4, $5, 'new',
			0, 'expo', $6, $7, CURRENT_TIMESTAMP, NULL, NULL, NULL
		)
		ON CONFLICT (tenant_id, restaurant_id, order_id)
		DO UPDATE SET
			order_number = EXCLUDED.order_number,
			customer_name = EXCLUDED.customer_name,
			status = 'new',
			priority = EXCLUDED.priority,
			station_id = EXCLUDED.station_id,
			items = EXCLUDED.items,
			notes = EXCLUDED.notes,
			received_at = CURRENT_TIMESTAMP,
			started_at = NULL,
			completed_at = NULL,
			bumped_by = NULL,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`

	var kdsID int64
	if err := r.db.QueryRowContext(
		ctx,
		query,
		tenantID,
		restaurantID,
		orderID,
		orderNumber,
		customerName,
		itemsJSON,
		notes,
	).Scan(&kdsID); err != nil {
		return nil, fmt.Errorf("failed to enqueue order in kds_queue: %w", err)
	}

	return r.getOrderByQueueID(ctx, tenantID, restaurantID, kdsID)
}

// GetQueue returns KDS queue entries filtered by status.
func (r *KDSRepository) GetQueue(ctx context.Context, tenantID, restaurantID int64, statuses []string) ([]domain.KDSOrder, error) {
	startedAt := time.Now()
	defer observability.RecordDB("kds.get_queue", time.Since(startedAt))

	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}

	baseQuery := `
		SELECT
			k.id, k.tenant_id, k.restaurant_id, k.order_id, k.order_number, COALESCE(o.order_source, ''), k.customer_name, k.status, k.priority,
			k.station_id, k.items, k.notes, k.received_at, k.started_at, k.completed_at, k.bumped_by, k.created_at, k.updated_at,
			o.table_id, t.table_number, COALESCE(t.name, '')
		FROM kds_queue k
		LEFT JOIN orders o
		  ON o.id = k.order_id
		 AND o.tenant_id = k.tenant_id
		 AND o.restaurant_id = k.restaurant_id
		LEFT JOIN restaurant_tables t
		  ON t.id = o.table_id
		 AND t.tenant_id = o.tenant_id
		 AND t.restaurant_id = o.restaurant_id
		WHERE k.tenant_id = $1 AND k.restaurant_id = $2
	`
	args := []interface{}{tenantID, restaurantID}

	if len(statuses) > 0 {
		valid := make([]string, 0, len(statuses))
		for _, status := range statuses {
			clean := strings.TrimSpace(status)
			if clean == "" {
				continue
			}
			if !domain.ValidKDSStatus(clean) {
				return nil, fmt.Errorf("invalid kds status filter: %s", clean)
			}
			valid = append(valid, clean)
		}

		if len(valid) > 0 {
			placeholders := make([]string, 0, len(valid))
			for _, status := range valid {
				args = append(args, status)
				placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
			}
			baseQuery += " AND k.status IN (" + strings.Join(placeholders, ",") + ")"
		}
	}

	baseQuery += " ORDER BY k.priority DESC, k.received_at ASC"

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch kds queue: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.KDSOrder, 0)
	for rows.Next() {
		order, err := scanKDSOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating kds queue rows: %w", err)
	}

	return orders, nil
}

// BumpOrder advances a KDS order through the kitchen statuses.
func (r *KDSRepository) BumpOrder(ctx context.Context, tenantID, restaurantID, queueID, userID int64) (*domain.KDSOrder, error) {
	startedAt := time.Now()
	defer observability.RecordDB("kds.bump_order", time.Since(startedAt))

	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	var currentStatus string
	var orderID int64
	selectForUpdateQuery := `
		SELECT status, order_id
		FROM kds_queue
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
		FOR UPDATE
	`
	err = tx.QueryRowContext(ctx, selectForUpdateQuery, queueID, tenantID, restaurantID).Scan(&currentStatus, &orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("kds order not found")
		}
		return nil, fmt.Errorf("failed to load kds order for bump: %w", err)
	}

	nextStatus, err := nextBumpStatus(currentStatus)
	if err != nil {
		return nil, err
	}

	updateQuery := `
		UPDATE kds_queue
		SET
			status = $1::text,
			started_at = CASE
				WHEN $1::text = 'in_progress'::text AND started_at IS NULL THEN CURRENT_TIMESTAMP
				ELSE started_at
			END,
			completed_at = CASE
				WHEN $1::text = 'completed'::text THEN CURRENT_TIMESTAMP
				ELSE completed_at
			END,
			bumped_by = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5 AND status = $6
	`
	result, err := tx.ExecContext(ctx, updateQuery, nextStatus, userID, queueID, tenantID, restaurantID, currentStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update kds bump status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check updated kds rows: %w", err)
	}
	if rowsAffected == 0 {
		return nil, errors.New("kds order status changed concurrently, retry")
	}

	if err := r.syncOrderStatusAfterBump(ctx, tx, tenantID, restaurantID, orderID, userID, currentStatus, nextStatus); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit bump transaction: %w", err)
	}

	return r.getOrderByQueueID(ctx, tenantID, restaurantID, queueID)
}

// RecallOrder moves a completed KDS order back to ready.
func (r *KDSRepository) RecallOrder(ctx context.Context, tenantID, restaurantID, queueID int64) (*domain.KDSOrder, error) {
	startedAt := time.Now()
	defer observability.RecordDB("kds.recall_order", time.Since(startedAt))

	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}

	query := `
		UPDATE kds_queue
		SET
			status = 'ready',
			completed_at = NULL,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
		  AND status IN ('completed', 'recalled')
	`
	result, err := r.db.ExecContext(ctx, query, queueID, tenantID, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to recall kds order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check recalled rows: %w", err)
	}
	if rowsAffected == 0 {
		return nil, errors.New("kds order not found or not recallable")
	}

	return r.getOrderByQueueID(ctx, tenantID, restaurantID, queueID)
}

// GetStats returns aggregate KDS queue statistics.
func (r *KDSRepository) GetStats(ctx context.Context, tenantID, restaurantID int64) (*domain.KDSStats, error) {
	startedAt := time.Now()
	defer observability.RecordDB("kds.get_stats", time.Since(startedAt))

	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}

	query := `
		SELECT
			COUNT(*) FILTER (WHERE status IN ('new', 'in_progress', 'ready', 'recalled')) AS active_orders,
			COALESCE(
				AVG(EXTRACT(EPOCH FROM (completed_at - received_at)) / 60.0)
				FILTER (WHERE status = 'completed' AND completed_at IS NOT NULL AND completed_at >= CURRENT_DATE),
				0
			) AS avg_prep_minutes,
			COUNT(*) FILTER (WHERE status = 'completed' AND completed_at >= CURRENT_DATE) AS completed_today,
			COUNT(*) FILTER (WHERE status = 'new') AS new_orders,
			COUNT(*) FILTER (WHERE status = 'in_progress') AS in_progress,
			COUNT(*) FILTER (WHERE status = 'ready') AS ready
		FROM kds_queue
		WHERE tenant_id = $1 AND restaurant_id = $2
	`

	stats := &domain.KDSStats{}
	if err := r.db.QueryRowContext(ctx, query, tenantID, restaurantID).Scan(
		&stats.ActiveOrders,
		&stats.AvgPrepMinutes,
		&stats.CompletedToday,
		&stats.NewOrders,
		&stats.InProgress,
		&stats.Ready,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch kds stats: %w", err)
	}

	return stats, nil
}

// GetAnalytics returns KDS preparation analytics for operational insights.
func (r *KDSRepository) GetAnalytics(ctx context.Context, tenantID, restaurantID int64) (*domain.KDSAnalytics, error) {
	startedAt := time.Now()
	defer observability.RecordDB("kds.get_analytics", time.Since(startedAt))

	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}

	analytics := &domain.KDSAnalytics{
		PrepTimeTrend: make([]domain.KDSPrepTrendPoint, 0, 7),
		BusiestHours:  make([]domain.KDSHourLoadPoint, 0, 24),
	}

	overallQuery := `
		SELECT
			COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - received_at)) / 60.0), 0) AS average_prep_minutes,
			COUNT(*) AS total_completed_orders
		FROM kds_queue
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND status = 'completed'
		  AND completed_at IS NOT NULL
	`

	if err := r.db.QueryRowContext(ctx, overallQuery, tenantID, restaurantID).Scan(
		&analytics.AveragePrepMinutes,
		&analytics.TotalCompletedOrders,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch kds analytics overview: %w", err)
	}

	trendQuery := `
		WITH days AS (
			SELECT generate_series(
				(CURRENT_DATE - INTERVAL '6 days')::date,
				CURRENT_DATE::date,
				INTERVAL '1 day'
			)::date AS day
		),
		completed AS (
			SELECT
				completed_at::date AS day,
				EXTRACT(EPOCH FROM (completed_at - received_at)) / 60.0 AS prep_minutes
			FROM kds_queue
			WHERE tenant_id = $1
			  AND restaurant_id = $2
			  AND status = 'completed'
			  AND completed_at IS NOT NULL
			  AND completed_at >= (CURRENT_DATE - INTERVAL '6 days')
		)
		SELECT
			TO_CHAR(days.day, 'YYYY-MM-DD') AS date,
			COALESCE(AVG(completed.prep_minutes), 0) AS avg_prep_minutes,
			COUNT(completed.prep_minutes) AS completed_orders
		FROM days
		LEFT JOIN completed ON completed.day = days.day
		GROUP BY days.day
		ORDER BY days.day ASC
	`

	trendRows, err := r.db.QueryContext(ctx, trendQuery, tenantID, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch kds preparation trend: %w", err)
	}
	defer trendRows.Close()

	for trendRows.Next() {
		point := domain.KDSPrepTrendPoint{}
		if err := trendRows.Scan(&point.Date, &point.AvgPrepMinutes, &point.CompletedOrders); err != nil {
			return nil, fmt.Errorf("failed to scan kds preparation trend row: %w", err)
		}
		analytics.PrepTimeTrend = append(analytics.PrepTimeTrend, point)
	}
	if err := trendRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating kds preparation trend rows: %w", err)
	}

	hourLoadQuery := `
		WITH hours AS (
			SELECT generate_series(0, 23) AS hour
		)
		SELECT
			hours.hour,
			COALESCE(COUNT(queue.id), 0) AS order_count
		FROM hours
		LEFT JOIN kds_queue AS queue
		  ON queue.tenant_id = $1
		 AND queue.restaurant_id = $2
		 AND queue.received_at >= (CURRENT_DATE - INTERVAL '6 days')
		 AND EXTRACT(HOUR FROM queue.received_at)::int = hours.hour
		GROUP BY hours.hour
		ORDER BY hours.hour ASC
	`

	hourRows, err := r.db.QueryContext(ctx, hourLoadQuery, tenantID, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch kds busiest hours: %w", err)
	}
	defer hourRows.Close()

	for hourRows.Next() {
		point := domain.KDSHourLoadPoint{}
		if err := hourRows.Scan(&point.Hour, &point.OrderCount); err != nil {
			return nil, fmt.Errorf("failed to scan kds busiest hour row: %w", err)
		}
		analytics.BusiestHours = append(analytics.BusiestHours, point)
	}
	if err := hourRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating kds busiest hour rows: %w", err)
	}

	return analytics, nil
}

// GetSettings returns KDS settings for the current tenant/restaurant.
func (r *KDSRepository) GetSettings(ctx context.Context, tenantID, restaurantID int64) (*domain.KDSSettings, error) {
	startedAt := time.Now()
	defer observability.RecordDB("kds.get_settings", time.Since(startedAt))

	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}

	query := `
		SELECT
			id, tenant_id, restaurant_id, auto_bump_minutes, alert_sound_enabled,
			show_analytics_panel, display_density, color_warning_minutes, color_danger_minutes, created_at, updated_at
		FROM kds_settings
		WHERE tenant_id = $1 AND restaurant_id = $2
	`

	settings := &domain.KDSSettings{}
	err := r.db.QueryRowContext(ctx, query, tenantID, restaurantID).Scan(
		&settings.ID,
		&settings.TenantID,
		&settings.RestaurantID,
		&settings.AutoBumpMinutes,
		&settings.AlertSoundEnabled,
		&settings.ShowAnalyticsPanel,
		&settings.DisplayDensity,
		&settings.ColorWarningMinutes,
		&settings.ColorDangerMinutes,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		if isMissingColumnError(err, "show_analytics_panel") {
			legacyQuery := `
				SELECT
					id, tenant_id, restaurant_id, auto_bump_minutes, alert_sound_enabled,
					display_density, color_warning_minutes, color_danger_minutes, created_at, updated_at
				FROM kds_settings
				WHERE tenant_id = $1 AND restaurant_id = $2
			`
			legacySettings := &domain.KDSSettings{}
			legacyErr := r.db.QueryRowContext(ctx, legacyQuery, tenantID, restaurantID).Scan(
				&legacySettings.ID,
				&legacySettings.TenantID,
				&legacySettings.RestaurantID,
				&legacySettings.AutoBumpMinutes,
				&legacySettings.AlertSoundEnabled,
				&legacySettings.DisplayDensity,
				&legacySettings.ColorWarningMinutes,
				&legacySettings.ColorDangerMinutes,
				&legacySettings.CreatedAt,
				&legacySettings.UpdatedAt,
			)
			if legacyErr != nil {
				if legacyErr == sql.ErrNoRows {
					return defaultKDSSettings(tenantID, restaurantID), nil
				}
				if isMissingRelationError(legacyErr, "kds_settings") {
					return defaultKDSSettings(tenantID, restaurantID), nil
				}
				return nil, fmt.Errorf("failed to fetch legacy kds settings: %w", legacyErr)
			}
			legacySettings.ShowAnalyticsPanel = false
			return legacySettings, nil
		}
		if isMissingRelationError(err, "kds_settings") {
			return defaultKDSSettings(tenantID, restaurantID), nil
		}
		if err == sql.ErrNoRows {
			return defaultKDSSettings(tenantID, restaurantID), nil
		}
		return nil, fmt.Errorf("failed to fetch kds settings: %w", err)
	}

	return settings, nil
}

// UpdateSettings upserts KDS settings.
func (r *KDSRepository) UpdateSettings(ctx context.Context, settings *domain.KDSSettings) error {
	startedAt := time.Now()
	defer observability.RecordDB("kds.update_settings", time.Since(startedAt))

	if r.db == nil {
		return errors.New("database is not initialized")
	}
	if settings == nil {
		return errors.New("settings payload is required")
	}

	query := `
		INSERT INTO kds_settings (
			tenant_id, restaurant_id, auto_bump_minutes, alert_sound_enabled, show_analytics_panel,
			display_density, color_warning_minutes, color_danger_minutes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (tenant_id, restaurant_id)
		DO UPDATE SET
			auto_bump_minutes = EXCLUDED.auto_bump_minutes,
			alert_sound_enabled = EXCLUDED.alert_sound_enabled,
			show_analytics_panel = EXCLUDED.show_analytics_panel,
			display_density = EXCLUDED.display_density,
			color_warning_minutes = EXCLUDED.color_warning_minutes,
			color_danger_minutes = EXCLUDED.color_danger_minutes,
			updated_at = CURRENT_TIMESTAMP
	`

	if _, err := r.db.ExecContext(
		ctx,
		query,
		settings.TenantID,
		settings.RestaurantID,
		settings.AutoBumpMinutes,
		settings.AlertSoundEnabled,
		settings.ShowAnalyticsPanel,
		settings.DisplayDensity,
		settings.ColorWarningMinutes,
		settings.ColorDangerMinutes,
	); err != nil {
		if isMissingColumnError(err, "show_analytics_panel") {
			legacyQuery := `
				INSERT INTO kds_settings (
					tenant_id, restaurant_id, auto_bump_minutes, alert_sound_enabled,
					display_density, color_warning_minutes, color_danger_minutes
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7
				)
				ON CONFLICT (tenant_id, restaurant_id)
				DO UPDATE SET
					auto_bump_minutes = EXCLUDED.auto_bump_minutes,
					alert_sound_enabled = EXCLUDED.alert_sound_enabled,
					display_density = EXCLUDED.display_density,
					color_warning_minutes = EXCLUDED.color_warning_minutes,
					color_danger_minutes = EXCLUDED.color_danger_minutes,
					updated_at = CURRENT_TIMESTAMP
			`
			if _, legacyErr := r.db.ExecContext(
				ctx,
				legacyQuery,
				settings.TenantID,
				settings.RestaurantID,
				settings.AutoBumpMinutes,
				settings.AlertSoundEnabled,
				settings.DisplayDensity,
				settings.ColorWarningMinutes,
				settings.ColorDangerMinutes,
			); legacyErr != nil {
				return fmt.Errorf("failed to update legacy kds settings: %w", legacyErr)
			}
			return nil
		}
		return fmt.Errorf("failed to update kds settings: %w", err)
	}

	return nil
}

func (r *KDSRepository) buildQueuePayload(ctx context.Context, tenantID, restaurantID, orderID int64) (string, string, string, []byte, error) {
	orderQuery := `
		SELECT order_number, customer_name, COALESCE(notes, '')
		FROM orders
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	var orderNumber string
	var customerName string
	var notes string
	err := r.db.QueryRowContext(ctx, orderQuery, orderID, tenantID, restaurantID).Scan(&orderNumber, &customerName, &notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", nil, errors.New("order not found")
		}
		return "", "", "", nil, fmt.Errorf("failed to fetch source order for kds enqueue: %w", err)
	}

	itemsQuery := `
		SELECT
			product_id,
			product_name,
			quantity,
			COALESCE(special_instructions, ''),
			COALESCE(addons::text, '')
		FROM order_items
		WHERE order_id = $1 AND tenant_id = $2 AND restaurant_id = $3
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, itemsQuery, orderID, tenantID, restaurantID)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("failed to fetch order items for kds enqueue: %w", err)
	}
	defer rows.Close()

	items := make([]domain.KDSItem, 0)
	for rows.Next() {
		var item domain.KDSItem
		var addOnsText string
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.Quantity, &item.SpecialInstructions, &addOnsText); err != nil {
			return "", "", "", nil, fmt.Errorf("failed to scan order item for kds enqueue: %w", err)
		}
		if addOnsText != "" && addOnsText != "null" {
			item.AddOns = addOnsText
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return "", "", "", nil, fmt.Errorf("error iterating order items for kds enqueue: %w", err)
	}

	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("failed to marshal kds items payload: %w", err)
	}

	return orderNumber, customerName, notes, itemsJSON, nil
}

func (r *KDSRepository) getOrderByQueueID(ctx context.Context, tenantID, restaurantID, queueID int64) (*domain.KDSOrder, error) {
	query := `
		SELECT
			k.id, k.tenant_id, k.restaurant_id, k.order_id, k.order_number, COALESCE(o.order_source, ''), k.customer_name, k.status, k.priority,
			k.station_id, k.items, k.notes, k.received_at, k.started_at, k.completed_at, k.bumped_by, k.created_at, k.updated_at,
			o.table_id, t.table_number, COALESCE(t.name, '')
		FROM kds_queue k
		LEFT JOIN orders o
		  ON o.id = k.order_id
		 AND o.tenant_id = k.tenant_id
		 AND o.restaurant_id = k.restaurant_id
		LEFT JOIN restaurant_tables t
		  ON t.id = o.table_id
		 AND t.tenant_id = o.tenant_id
		 AND t.restaurant_id = o.restaurant_id
		WHERE k.id = $1 AND k.tenant_id = $2 AND k.restaurant_id = $3
	`

	row := r.db.QueryRowContext(ctx, query, queueID, tenantID, restaurantID)
	order, err := scanKDSOrder(row)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *KDSRepository) syncOrderStatusAfterBump(
	ctx context.Context,
	tx *sql.Tx,
	tenantID, restaurantID, orderID, userID int64,
	oldStatus, newStatus string,
) error {
	var expectedOrderStatus string
	var targetOrderStatus string

	switch {
	case oldStatus == "new" && newStatus == "in_progress":
		expectedOrderStatus = "confirmed"
		targetOrderStatus = "preparing"
	case oldStatus == "in_progress" && newStatus == "ready":
		expectedOrderStatus = "preparing"
		targetOrderStatus = "ready"
	default:
		return nil
	}

	updateQuery := `
		UPDATE orders
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4 AND status = $5
	`
	result, err := tx.ExecContext(ctx, updateQuery, targetOrderStatus, orderID, tenantID, restaurantID, expectedOrderStatus)
	if err != nil {
		return fmt.Errorf("failed to sync order status from kds bump: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to inspect synced order rows: %w", err)
	}
	if rowsAffected == 0 {
		return nil
	}

	historyQuery := `
		INSERT INTO order_status_history (order_id, old_status, new_status, changed_by, change_reason)
		VALUES ($1, $2, $3, $4, $5)
	`
	changeReason := "Synced from KDS bump"
	if _, err := tx.ExecContext(ctx, historyQuery, orderID, expectedOrderStatus, targetOrderStatus, userID, changeReason); err != nil {
		return fmt.Errorf("failed to write order status history from kds bump: %w", err)
	}

	return nil
}

func nextBumpStatus(currentStatus string) (string, error) {
	switch currentStatus {
	case string(domain.KDSStatusNew):
		return string(domain.KDSStatusInProgress), nil
	case string(domain.KDSStatusInProgress):
		return string(domain.KDSStatusReady), nil
	case string(domain.KDSStatusReady):
		return string(domain.KDSStatusCompleted), nil
	case string(domain.KDSStatusCompleted):
		return "", errors.New("order is already completed; use recall endpoint to reopen")
	case string(domain.KDSStatusRecalled):
		return string(domain.KDSStatusReady), nil
	default:
		return "", fmt.Errorf("invalid current kds status: %s", currentStatus)
	}
}

type kdsScanner interface {
	Scan(dest ...interface{}) error
}

func scanKDSOrder(scanner kdsScanner) (*domain.KDSOrder, error) {
	order := &domain.KDSOrder{}

	var customerName sql.NullString
	var orderSource sql.NullString
	var status string
	var itemsRaw []byte
	var notes sql.NullString
	var startedAt sql.NullTime
	var completedAt sql.NullTime
	var bumpedBy sql.NullInt64
	var tableID sql.NullInt64
	var tableNumber sql.NullInt32
	var tableName sql.NullString

	err := scanner.Scan(
		&order.ID,
		&order.TenantID,
		&order.RestaurantID,
		&order.OrderID,
		&order.OrderNumber,
		&orderSource,
		&customerName,
		&status,
		&order.Priority,
		&order.StationID,
		&itemsRaw,
		&notes,
		&order.ReceivedAt,
		&startedAt,
		&completedAt,
		&bumpedBy,
		&order.CreatedAt,
		&order.UpdatedAt,
		&tableID,
		&tableNumber,
		&tableName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("kds order not found")
		}
		return nil, fmt.Errorf("failed to scan kds order: %w", err)
	}

	order.Status = domain.KDSStatus(status)
	if customerName.Valid {
		order.CustomerName = customerName.String
	}
	if orderSource.Valid {
		order.OrderSource = orderSource.String
	}
	if notes.Valid {
		order.Notes = notes.String
	}
	if startedAt.Valid {
		order.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		order.CompletedAt = &completedAt.Time
	}
	if bumpedBy.Valid {
		order.BumpedBy = &bumpedBy.Int64
	}
	if tableID.Valid {
		order.TableID = &tableID.Int64
	}
	if tableNumber.Valid {
		number := int(tableNumber.Int32)
		order.TableNumber = &number
	}
	if tableName.Valid {
		order.TableName = tableName.String
	}

	if len(itemsRaw) > 0 {
		if err := json.Unmarshal(itemsRaw, &order.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal kds items payload: %w", err)
		}
	}
	if order.Items == nil {
		order.Items = []domain.KDSItem{}
	}

	order.ElapsedSeconds = int64(time.Since(order.ReceivedAt).Seconds())
	if order.CompletedAt != nil {
		order.ElapsedSeconds = int64(order.CompletedAt.Sub(order.ReceivedAt).Seconds())
	}
	if order.ElapsedSeconds < 0 {
		order.ElapsedSeconds = 0
	}

	return order, nil
}
