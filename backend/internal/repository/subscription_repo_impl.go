package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"pos-saas/internal/domain"
)

// SubscriptionRepositoryImpl implements the SubscriptionRepository interface
type SubscriptionRepositoryImpl struct {
	db *sql.DB
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &SubscriptionRepositoryImpl{db: db}
}

// GetTenantSubscriptionInfo retrieves subscription information for a tenant
func (r *SubscriptionRepositoryImpl) GetTenantSubscriptionInfo(ctx context.Context, tenantID int64) (*domain.TenantSubscriptionInfo, error) {
	log.Printf("[SubscriptionRepo] GetTenantSubscriptionInfo - tenant_id: %d\n", tenantID)

	var info domain.TenantSubscriptionInfo
	var disabledReason sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			name,
			slug,
			subscription_status,
			trial_start_date,
			trial_end_date,
			trial_days,
			is_disabled,
			disabled_at,
			disabled_reason
		FROM tenants
		WHERE id = $1
	`, tenantID).Scan(
		&info.TenantID,
		&info.TenantName,
		&info.Slug,
		&info.SubscriptionStatus,
		&info.TrialStartDate,
		&info.TrialEndDate,
		&info.TrialDays,
		&info.IsDisabled,
		&info.DisabledAt,
		&disabledReason,
	)

	if err == sql.ErrNoRows {
		log.Printf("[SubscriptionRepo] ERROR: Tenant not found - tenant_id: %d\n", tenantID)
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to query tenant - tenant_id: %d, error: %v\n", tenantID, err)
		return nil, fmt.Errorf("failed to fetch tenant subscription: %w", err)
	}

	// Handle nullable disabled_reason
	if disabledReason.Valid {
		info.DisabledReason = &disabledReason.String
	}

	// Calculate days remaining
	if info.TrialEndDate != nil && time.Now().Before(*info.TrialEndDate) {
		days := info.TrialEndDate.Sub(time.Now()).Hours() / 24
		info.DaysRemaining = int(days) + 1 // Round up
	} else {
		info.DaysRemaining = 0
	}

	// Determine access status
	if info.IsDisabled {
		info.CanAccess = false
		blockReason := "tenant_disabled"
		info.BlockReason = &blockReason
	} else if info.SubscriptionStatus == "expired" {
		info.CanAccess = false
		blockReason := "expired"
		info.BlockReason = &blockReason
	} else if info.SubscriptionStatus == "trial" && info.TrialEndDate != nil && time.Now().After(*info.TrialEndDate) {
		info.CanAccess = false
		blockReason := "trial_expired"
		info.BlockReason = &blockReason
	} else {
		info.CanAccess = true
	}

	log.Printf("[SubscriptionRepo] Query result - status: %s, days_remaining: %d, can_access: %v\n",
		info.SubscriptionStatus, info.DaysRemaining, info.CanAccess)

	return &info, nil
}

// GetTenantSubscriptionBySlug retrieves subscription info by tenant slug
func (r *SubscriptionRepositoryImpl) GetTenantSubscriptionBySlug(ctx context.Context, slug string) (*domain.TenantSubscriptionInfo, error) {
	log.Printf("[SubscriptionRepo] GetTenantSubscriptionBySlug - slug: %s\n", slug)

	var tenantID int64
	err := r.db.QueryRowContext(ctx, `SELECT id FROM tenants WHERE slug = $1`, slug).Scan(&tenantID)

	if err == sql.ErrNoRows {
		log.Printf("[SubscriptionRepo] ERROR: Tenant not found - slug: %s\n", slug)
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to query tenant - slug: %s, error: %v\n", slug, err)
		return nil, fmt.Errorf("failed to fetch tenant: %w", err)
	}

	return r.GetTenantSubscriptionInfo(ctx, tenantID)
}

// UpdateSubscriptionStatus updates the subscription status
func (r *SubscriptionRepositoryImpl) UpdateSubscriptionStatus(ctx context.Context, tenantID int64, status string) error {
	log.Printf("[SubscriptionRepo] UpdateSubscriptionStatus - tenant_id: %d, status: %s\n", tenantID, status)

	if !domain.ValidateSubscriptionStatus(status) {
		log.Printf("[SubscriptionRepo] ERROR: Invalid subscription status - status: %s\n", status)
		return fmt.Errorf("invalid subscription status: %s", status)
	}

	result, err := r.db.ExecContext(ctx, `
		UPDATE tenants
		SET subscription_status = $1, updated_at = NOW()
		WHERE id = $2
	`, status, tenantID)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to update status - tenant_id: %d, error: %v\n", tenantID, err)
		return fmt.Errorf("failed to update subscription status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to get rows affected - error: %v\n", err)
		return fmt.Errorf("failed to verify update: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("[SubscriptionRepo] ERROR: Tenant not found - tenant_id: %d\n", tenantID)
		return fmt.Errorf("tenant not found")
	}

	log.Printf("[SubscriptionRepo] UpdateSubscriptionStatus success - tenant_id: %d\n", tenantID)
	return nil
}

// ExtendTrial extends the trial period and logs the action
func (r *SubscriptionRepositoryImpl) ExtendTrial(ctx context.Context, params *domain.ExtendTrialParams) error {
	log.Printf("[SubscriptionRepo] ExtendTrial - tenant_id: %d, days: %d\n", params.TenantID, params.Days)

	if params.Days <= 0 {
		log.Printf("[SubscriptionRepo] ERROR: Invalid days - days: %d\n", params.Days)
		return fmt.Errorf("days must be greater than 0")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to begin transaction - error: %v\n", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current trial end date
	var oldTrialEndDate sql.NullTime
	var currentStatus string
	err = tx.QueryRowContext(ctx, `
		SELECT trial_end_date, subscription_status
		FROM tenants
		WHERE id = $1
	`, params.TenantID).Scan(&oldTrialEndDate, &currentStatus)

	if err == sql.ErrNoRows {
		log.Printf("[SubscriptionRepo] ERROR: Tenant not found - tenant_id: %d\n", params.TenantID)
		return fmt.Errorf("tenant not found")
	}

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to query tenant - error: %v\n", err)
		return fmt.Errorf("failed to query tenant: %w", err)
	}

	// Calculate new trial end date
	newTrialEndDate := time.Now().AddDate(0, 0, params.Days)
	if oldTrialEndDate.Valid {
		newTrialEndDate = oldTrialEndDate.Time.AddDate(0, 0, params.Days)
	}

	// Update trial end date and reset status to trial if expired
	newStatus := currentStatus
	if currentStatus == "expired" {
		newStatus = "trial"
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE tenants
		SET
			trial_end_date = $1,
			subscription_status = $2,
			trial_days = trial_days + $3,
			updated_at = NOW()
		WHERE id = $4
	`, newTrialEndDate, newStatus, params.Days, params.TenantID)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to update trial - error: %v\n", err)
		return fmt.Errorf("failed to update trial: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Printf("[SubscriptionRepo] ERROR: No rows affected - error: %v\n", err)
		return fmt.Errorf("failed to update tenant trial")
	}

	// Create history record
	historyParams := &domain.CreateHistoryRecordParams{
		TenantID:        params.TenantID,
		Action:          domain.ActionTrialExtended,
		PerformedBy:     &params.PerformedBy,
		OldStatus:       &currentStatus,
		NewStatus:       &newStatus,
		OldTrialEndDate: &oldTrialEndDate.Time,
		NewTrialEndDate: &newTrialEndDate,
		DaysAdded:       &params.Days,
		Reason:          params.Reason,
		IPAddress:       &params.IPAddress,
		UserAgent:       &params.UserAgent,
	}

	if err := r.createHistoryRecordTx(ctx, tx, historyParams); err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to create history record - error: %v\n", err)
		return fmt.Errorf("failed to create history record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to commit transaction - error: %v\n", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("[SubscriptionRepo] ExtendTrial success - new_end_date: %s\n", newTrialEndDate.Format(time.RFC3339))
	return nil
}

// DisableTenant disables a tenant's access
func (r *SubscriptionRepositoryImpl) DisableTenant(ctx context.Context, params *domain.DisableTenantParams) error {
	log.Printf("[SubscriptionRepo] DisableTenant - tenant_id: %d\n", params.TenantID)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to begin transaction - error: %v\n", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current status
	var currentStatus string
	err = tx.QueryRowContext(ctx, `
		SELECT subscription_status FROM tenants WHERE id = $1
	`, params.TenantID).Scan(&currentStatus)

	if err == sql.ErrNoRows {
		log.Printf("[SubscriptionRepo] ERROR: Tenant not found - tenant_id: %d\n", params.TenantID)
		return fmt.Errorf("tenant not found")
	}

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to query tenant - error: %v\n", err)
		return fmt.Errorf("failed to query tenant: %w", err)
	}

	// Update tenant
	disabledStatus := "disabled"
	result, err := tx.ExecContext(ctx, `
		UPDATE tenants
		SET
			is_disabled = true,
			disabled_at = NOW(),
			disabled_by = $1,
			disabled_reason = $2,
			subscription_status = $3,
			updated_at = NOW()
		WHERE id = $4
	`, params.PerformedBy, params.Reason, disabledStatus, params.TenantID)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to disable tenant - error: %v\n", err)
		return fmt.Errorf("failed to disable tenant: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("tenant not found")
	}

	// Create history record
	historyParams := &domain.CreateHistoryRecordParams{
		TenantID:    params.TenantID,
		Action:      domain.ActionTenantDisabled,
		PerformedBy: &params.PerformedBy,
		OldStatus:   &currentStatus,
		NewStatus:   &disabledStatus,
		Reason:      &params.Reason,
		IPAddress:   &params.IPAddress,
		UserAgent:   &params.UserAgent,
	}

	if err := r.createHistoryRecordTx(ctx, tx, historyParams); err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to create history record - error: %v\n", err)
		return fmt.Errorf("failed to create history record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to commit transaction - error: %v\n", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("[SubscriptionRepo] DisableTenant success\n")
	return nil
}

// EnableTenant re-enables a disabled tenant
func (r *SubscriptionRepositoryImpl) EnableTenant(ctx context.Context, params *domain.EnableTenantParams) error {
	log.Printf("[SubscriptionRepo] EnableTenant - tenant_id: %d\n", params.TenantID)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to begin transaction - error: %v\n", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current status and check if disabled
	var currentStatus string
	var isDisabled bool
	err = tx.QueryRowContext(ctx, `
		SELECT subscription_status, is_disabled FROM tenants WHERE id = $1
	`, params.TenantID).Scan(&currentStatus, &isDisabled)

	if err == sql.ErrNoRows {
		log.Printf("[SubscriptionRepo] ERROR: Tenant not found - tenant_id: %d\n", params.TenantID)
		return fmt.Errorf("tenant not found")
	}

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to query tenant - error: %v\n", err)
		return fmt.Errorf("failed to query tenant: %w", err)
	}

	if !isDisabled {
		log.Printf("[SubscriptionRepo] ERROR: Tenant is not disabled - tenant_id: %d\n", params.TenantID)
		return fmt.Errorf("tenant is not disabled")
	}

	// Restore to trial if no paid subscription
	restoredStatus := "trial"
	if currentStatus == "active" || currentStatus == "suspended" {
		restoredStatus = currentStatus
	}

	// Update tenant
	result, err := tx.ExecContext(ctx, `
		UPDATE tenants
		SET
			is_disabled = false,
			enabled_at = NOW(),
			disabled_by = NULL,
			disabled_reason = NULL,
			subscription_status = $1,
			updated_at = NOW()
		WHERE id = $2
	`, restoredStatus, params.TenantID)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to enable tenant - error: %v\n", err)
		return fmt.Errorf("failed to enable tenant: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("tenant not found")
	}

	// Create history record
	historyParams := &domain.CreateHistoryRecordParams{
		TenantID:    params.TenantID,
		Action:      domain.ActionTenantEnabled,
		PerformedBy: &params.PerformedBy,
		OldStatus:   &currentStatus,
		NewStatus:   &restoredStatus,
		IPAddress:   &params.IPAddress,
		UserAgent:   &params.UserAgent,
	}

	if err := r.createHistoryRecordTx(ctx, tx, historyParams); err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to create history record - error: %v\n", err)
		return fmt.Errorf("failed to create history record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to commit transaction - error: %v\n", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("[SubscriptionRepo] EnableTenant success\n")
	return nil
}

// ScheduleDeletion schedules a tenant for deletion
func (r *SubscriptionRepositoryImpl) ScheduleDeletion(ctx context.Context, params *domain.DeleteTenantParams) error {
	log.Printf("[SubscriptionRepo] ScheduleDeletion - tenant_id: %d\n", params.TenantID)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to begin transaction - error: %v\n", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current status
	var currentStatus string
	err = tx.QueryRowContext(ctx, `
		SELECT subscription_status FROM tenants WHERE id = $1
	`, params.TenantID).Scan(&currentStatus)

	if err == sql.ErrNoRows {
		log.Printf("[SubscriptionRepo] ERROR: Tenant not found - tenant_id: %d\n", params.TenantID)
		return fmt.Errorf("tenant not found")
	}

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to query tenant - error: %v\n", err)
		return fmt.Errorf("failed to query tenant: %w", err)
	}

	// Schedule deletion
	_, err = tx.ExecContext(ctx, `
		UPDATE tenants
		SET deletion_scheduled_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, params.TenantID)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to schedule deletion - error: %v\n", err)
		return fmt.Errorf("failed to schedule deletion: %w", err)
	}

	// Create history record
	historyParams := &domain.CreateHistoryRecordParams{
		TenantID:    params.TenantID,
		Action:      domain.ActionDeletionScheduled,
		PerformedBy: &params.PerformedBy,
		OldStatus:   &currentStatus,
		Reason:      params.Reason,
		IPAddress:   &params.IPAddress,
		UserAgent:   &params.UserAgent,
	}

	if err := r.createHistoryRecordTx(ctx, tx, historyParams); err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to create history record - error: %v\n", err)
		return fmt.Errorf("failed to create history record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to commit transaction - error: %v\n", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("[SubscriptionRepo] ScheduleDeletion success\n")
	return nil
}

// CancelDeletion cancels a scheduled deletion
func (r *SubscriptionRepositoryImpl) CancelDeletion(ctx context.Context, tenantID int64, performedBy int64) error {
	log.Printf("[SubscriptionRepo] CancelDeletion - tenant_id: %d\n", tenantID)

	result, err := r.db.ExecContext(ctx, `
		UPDATE tenants
		SET deletion_scheduled_at = NULL, updated_at = NOW()
		WHERE id = $1 AND deletion_scheduled_at IS NOT NULL
	`, tenantID)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to cancel deletion - error: %v\n", err)
		return fmt.Errorf("failed to cancel deletion: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("[SubscriptionRepo] ERROR: No deletion scheduled for tenant - tenant_id: %d\n", tenantID)
		return fmt.Errorf("no deletion scheduled for this tenant")
	}

	log.Printf("[SubscriptionRepo] CancelDeletion success\n")
	return nil
}

// GetSubscriptionHistory retrieves paginated subscription change history
func (r *SubscriptionRepositoryImpl) GetSubscriptionHistory(ctx context.Context, tenantID int64, page int, pageSize int) ([]domain.SubscriptionHistory, int64, error) {
	log.Printf("[SubscriptionRepo] GetSubscriptionHistory - tenant_id: %d, page: %d, pageSize: %d\n", tenantID, page, pageSize)

	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int64
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM tenant_subscription_history WHERE tenant_id = $1
	`, tenantID).Scan(&totalCount)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to count history - error: %v\n", err)
		return nil, 0, fmt.Errorf("failed to count history: %w", err)
	}

	// Get paginated results
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id, tenant_id, action, performed_by, old_status, new_status,
			old_trial_end_date, new_trial_end_date, days_added, reason,
			ip_address, user_agent, metadata, created_at
		FROM tenant_subscription_history
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, tenantID, pageSize, offset)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to query history - error: %v\n", err)
		return nil, 0, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	var history []domain.SubscriptionHistory
	for rows.Next() {
		var h domain.SubscriptionHistory
		var metadataJSON sql.NullString
		var oldStatus, newStatus sql.NullString
		var oldTrialEndDate, newTrialEndDate sql.NullTime
		var daysAdded sql.NullInt64
		var reason, ipAddress, userAgent sql.NullString

		err := rows.Scan(
			&h.ID, &h.TenantID, &h.Action, &h.PerformedBy, &oldStatus, &newStatus,
			&oldTrialEndDate, &newTrialEndDate, &daysAdded, &reason,
			&ipAddress, &userAgent, &metadataJSON, &h.CreatedAt,
		)

		if err != nil {
			log.Printf("[SubscriptionRepo] ERROR: Failed to scan row - error: %v\n", err)
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}

		// Handle nullable fields
		if oldStatus.Valid {
			h.OldStatus = &oldStatus.String
		}
		if newStatus.Valid {
			h.NewStatus = &newStatus.String
		}
		if oldTrialEndDate.Valid {
			h.OldTrialEndDate = &oldTrialEndDate.Time
		}
		if newTrialEndDate.Valid {
			h.NewTrialEndDate = &newTrialEndDate.Time
		}
		if daysAdded.Valid {
			daysVal := int(daysAdded.Int64)
			h.DaysAdded = &daysVal
		}
		if reason.Valid {
			h.Reason = &reason.String
		}
		if ipAddress.Valid {
			h.IPAddress = &ipAddress.String
		}
		if userAgent.Valid {
			h.UserAgent = &userAgent.String
		}
		if metadataJSON.Valid {
			json.Unmarshal([]byte(metadataJSON.String), &h.Metadata)
		}

		history = append(history, h)
	}

	log.Printf("[SubscriptionRepo] GetSubscriptionHistory success - count: %d\n", len(history))
	return history, totalCount, nil
}

// CreateHistoryRecord creates a subscription history record
func (r *SubscriptionRepositoryImpl) CreateHistoryRecord(ctx context.Context, params *domain.CreateHistoryRecordParams) error {
	log.Printf("[SubscriptionRepo] CreateHistoryRecord - tenant_id: %d, action: %s\n", params.TenantID, params.Action)

	return r.createHistoryRecordTx(ctx, r.db, params)
}

// createHistoryRecordTx creates a history record (can be used within a transaction)
func (r *SubscriptionRepositoryImpl) createHistoryRecordTx(ctx context.Context, executor interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}, params *domain.CreateHistoryRecordParams) error {

	if !domain.ValidateSubscriptionAction(params.Action) {
		log.Printf("[SubscriptionRepo] ERROR: Invalid action - action: %s\n", params.Action)
		return fmt.Errorf("invalid subscription action: %s", params.Action)
	}

	var metadataJSON sql.NullString
	if params.Metadata != nil {
		data, err := json.Marshal(params.Metadata)
		if err == nil {
			metadataJSON = sql.NullString{String: string(data), Valid: true}
		}
	}

	_, err := executor.ExecContext(ctx, `
		INSERT INTO tenant_subscription_history
		(tenant_id, action, performed_by, old_status, new_status, old_trial_end_date, new_trial_end_date,
		 days_added, reason, ip_address, user_agent, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())
	`,
		params.TenantID, params.Action, params.PerformedBy, params.OldStatus, params.NewStatus,
		params.OldTrialEndDate, params.NewTrialEndDate, params.DaysAdded, params.Reason,
		params.IPAddress, params.UserAgent, metadataJSON,
	)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to create history record - error: %v\n", err)
		return fmt.Errorf("failed to create history record: %w", err)
	}

	log.Printf("[SubscriptionRepo] CreateHistoryRecord success\n")
	return nil
}

// GetExpiredTrials retrieves all tenants with expired trials
func (r *SubscriptionRepositoryImpl) GetExpiredTrials(ctx context.Context) ([]domain.TenantSubscriptionInfo, error) {
	log.Printf("[SubscriptionRepo] GetExpiredTrials\n")

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id, name, slug, subscription_status, trial_start_date, trial_end_date,
			trial_days, is_disabled, disabled_at, disabled_reason
		FROM tenants
		WHERE subscription_status = 'trial' AND trial_end_date < NOW()
	`)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to query expired trials - error: %v\n", err)
		return nil, fmt.Errorf("failed to query expired trials: %w", err)
	}
	defer rows.Close()

	var trials []domain.TenantSubscriptionInfo
	for rows.Next() {
		var info domain.TenantSubscriptionInfo
		var disabledReason sql.NullString

		err := rows.Scan(
			&info.TenantID, &info.TenantName, &info.Slug, &info.SubscriptionStatus,
			&info.TrialStartDate, &info.TrialEndDate, &info.TrialDays,
			&info.IsDisabled, &info.DisabledAt, &disabledReason,
		)

		if err != nil {
			log.Printf("[SubscriptionRepo] ERROR: Failed to scan row - error: %v\n", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if disabledReason.Valid {
			info.DisabledReason = &disabledReason.String
		}

		info.DaysRemaining = 0
		info.CanAccess = false
		blockReason := "trial_expired"
		info.BlockReason = &blockReason

		trials = append(trials, info)
	}

	log.Printf("[SubscriptionRepo] GetExpiredTrials success - count: %d\n", len(trials))
	return trials, nil
}

// GetTenantsScheduledForDeletion retrieves all tenants scheduled for deletion
func (r *SubscriptionRepositoryImpl) GetTenantsScheduledForDeletion(ctx context.Context) ([]domain.TenantSubscriptionInfo, error) {
	log.Printf("[SubscriptionRepo] GetTenantsScheduledForDeletion\n")

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id, name, slug, subscription_status, trial_start_date, trial_end_date,
			trial_days, is_disabled, disabled_at, disabled_reason
		FROM tenants
		WHERE deletion_scheduled_at IS NOT NULL
		ORDER BY deletion_scheduled_at ASC
	`)

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to query scheduled deletions - error: %v\n", err)
		return nil, fmt.Errorf("failed to query scheduled deletions: %w", err)
	}
	defer rows.Close()

	var tenants []domain.TenantSubscriptionInfo
	for rows.Next() {
		var info domain.TenantSubscriptionInfo
		var disabledReason sql.NullString

		err := rows.Scan(
			&info.TenantID, &info.TenantName, &info.Slug, &info.SubscriptionStatus,
			&info.TrialStartDate, &info.TrialEndDate, &info.TrialDays,
			&info.IsDisabled, &info.DisabledAt, &disabledReason,
		)

		if err != nil {
			log.Printf("[SubscriptionRepo] ERROR: Failed to scan row - error: %v\n", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if disabledReason.Valid {
			info.DisabledReason = &disabledReason.String
		}

		tenants = append(tenants, info)
	}

	log.Printf("[SubscriptionRepo] GetTenantsScheduledForDeletion success - count: %d\n", len(tenants))
	return tenants, nil
}

// CheckTenantAccess checks if a tenant can access the system
func (r *SubscriptionRepositoryImpl) CheckTenantAccess(ctx context.Context, tenantID int64) (bool, string, error) {
	log.Printf("[SubscriptionRepo] CheckTenantAccess - tenant_id: %d\n", tenantID)

	var isDisabled bool
	var status string
	var trialEndDate sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT is_disabled, subscription_status, trial_end_date
		FROM tenants
		WHERE id = $1
	`, tenantID).Scan(&isDisabled, &status, &trialEndDate)

	if err == sql.ErrNoRows {
		log.Printf("[SubscriptionRepo] ERROR: Tenant not found - tenant_id: %d\n", tenantID)
		return false, "tenant_not_found", fmt.Errorf("tenant not found")
	}

	if err != nil {
		log.Printf("[SubscriptionRepo] ERROR: Failed to query tenant - error: %v\n", err)
		return false, "", fmt.Errorf("failed to check access: %w", err)
	}

	if isDisabled {
		log.Printf("[SubscriptionRepo] Tenant is disabled - tenant_id: %d\n", tenantID)
		return false, "tenant_disabled", nil
	}

	if status == "trial" && trialEndDate.Valid && time.Now().After(trialEndDate.Time) {
		log.Printf("[SubscriptionRepo] Trial expired - tenant_id: %d\n", tenantID)
		return false, "trial_expired", nil
	}

	log.Printf("[SubscriptionRepo] Tenant has access - tenant_id: %d\n", tenantID)
	return true, "", nil
}
