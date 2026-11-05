package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"pos-saas/internal/domain"
)

// SuperAdminAuditRepository handles super admin audit log database operations
type SuperAdminAuditRepository struct {
	db *sql.DB
}

// NewSuperAdminAuditRepository creates a new super admin audit repository
func NewSuperAdminAuditRepository(db *sql.DB) *SuperAdminAuditRepository {
	return &SuperAdminAuditRepository{db: db}
}

// LogAction creates a new audit log entry
func (r *SuperAdminAuditRepository) LogAction(ctx context.Context, log_ *domain.SuperAdminAuditLog) error {
	log.Printf("[SuperAdminAuditRepository] LogAction: superAdminID=%d, action=%s\n", log_.SuperAdminID, log_.Action)

	// Convert details to JSON
	var detailsJSON []byte
	var err error
	if log_.Details != nil {
		detailsJSON, err = json.Marshal(log_.Details)
		if err != nil {
			log.Printf("[SuperAdminAuditRepository] ERROR: Failed to marshal details: %v\n", err)
			return fmt.Errorf("failed to marshal details: %w", err)
		}
	}
	var detailsValue interface{}
	if detailsJSON != nil {
		// lib/pq treats []byte as bytea; JSONB parameters must be sent as JSON text.
		detailsValue = string(detailsJSON)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO super_admin_audit_logs
			(super_admin_id, action, target_type, target_id, details, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`, log_.SuperAdminID, log_.Action, log_.TargetType, log_.TargetID, detailsValue, log_.IPAddress, log_.UserAgent)

	if err != nil {
		log.Printf("[SuperAdminAuditRepository] ERROR: Failed to create audit log: %v\n", err)
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	log.Printf("[SuperAdminAuditRepository] Audit log created for action: %s\n", log_.Action)
	return nil
}

// ListLogs retrieves audit logs with filtering and pagination
func (r *SuperAdminAuditRepository) ListLogs(ctx context.Context, filter domain.AuditLogFilter, pagination domain.Pagination) ([]domain.SuperAdminAuditLog, int64, error) {
	log.Printf("[SuperAdminAuditRepository] ListLogs: page=%d, pageSize=%d\n", pagination.Page, pagination.PageSize)

	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.SuperAdminID != nil {
		conditions = append(conditions, fmt.Sprintf("super_admin_id = $%d", argIndex))
		args = append(args, *filter.SuperAdminID)
		argIndex++
	}

	if filter.Action != nil {
		conditions = append(conditions, fmt.Sprintf("action = $%d", argIndex))
		args = append(args, *filter.Action)
		argIndex++
	}

	if filter.TargetType != nil {
		conditions = append(conditions, fmt.Sprintf("target_type = $%d", argIndex))
		args = append(args, *filter.TargetType)
		argIndex++
	}

	if filter.TargetID != nil {
		conditions = append(conditions, fmt.Sprintf("target_id = $%d", argIndex))
		args = append(args, *filter.TargetID)
		argIndex++
	}

	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.EndDate)
		argIndex++
	}

	if filter.IPAddress != nil {
		conditions = append(conditions, fmt.Sprintf("ip_address = $%d::inet", argIndex))
		args = append(args, *filter.IPAddress)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM super_admin_audit_logs %s`, whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("[SuperAdminAuditRepository] ERROR: Failed to count logs: %v\n", err)
		return nil, 0, fmt.Errorf("failed to count logs: %w", err)
	}

	// Main query
	offset := (pagination.Page - 1) * pagination.PageSize
	args = append(args, pagination.PageSize, offset)

	query := fmt.Sprintf(`
		SELECT id, super_admin_id, action, target_type, target_id, details, ip_address, user_agent, created_at
		FROM super_admin_audit_logs
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("[SuperAdminAuditRepository] ERROR: Failed to query logs: %v\n", err)
		return nil, 0, fmt.Errorf("failed to fetch logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.SuperAdminAuditLog
	for rows.Next() {
		var l domain.SuperAdminAuditLog
		var targetType, userAgent sql.NullString
		var targetID sql.NullInt64
		var detailsJSON []byte

		err := rows.Scan(
			&l.ID, &l.SuperAdminID, &l.Action, &targetType, &targetID,
			&detailsJSON, &l.IPAddress, &userAgent, &l.CreatedAt,
		)
		if err != nil {
			log.Printf("[SuperAdminAuditRepository] ERROR: Failed to scan log: %v\n", err)
			return nil, 0, fmt.Errorf("failed to scan log: %w", err)
		}

		if targetType.Valid {
			l.TargetType = &targetType.String
		}
		if targetID.Valid {
			l.TargetID = &targetID.Int64
		}
		if userAgent.Valid {
			l.UserAgent = userAgent.String
		}

		// Parse details JSON
		if len(detailsJSON) > 0 {
			var details map[string]interface{}
			if err := json.Unmarshal(detailsJSON, &details); err == nil {
				l.Details = details
			}
		}

		logs = append(logs, l)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[SuperAdminAuditRepository] ERROR: Row iteration error: %v\n", err)
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("[SuperAdminAuditRepository] Found %d logs (total: %d)\n", len(logs), total)
	return logs, total, nil
}

// GetLogsByTarget retrieves all logs for a specific target
func (r *SuperAdminAuditRepository) GetLogsByTarget(ctx context.Context, targetType string, targetID int64, pagination domain.Pagination) ([]domain.SuperAdminAuditLog, int64, error) {
	log.Printf("[SuperAdminAuditRepository] GetLogsByTarget: type=%s, id=%d\n", targetType, targetID)

	// Count
	var total int64
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM super_admin_audit_logs WHERE target_type = $1 AND target_id = $2
	`, targetType, targetID).Scan(&total)
	if err != nil {
		log.Printf("[SuperAdminAuditRepository] ERROR: Failed to count logs: %v\n", err)
		return nil, 0, fmt.Errorf("failed to count logs: %w", err)
	}

	// Query
	offset := (pagination.Page - 1) * pagination.PageSize
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, super_admin_id, action, target_type, target_id, details, ip_address, user_agent, created_at
		FROM super_admin_audit_logs
		WHERE target_type = $1 AND target_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, targetType, targetID, pagination.PageSize, offset)

	if err != nil {
		log.Printf("[SuperAdminAuditRepository] ERROR: Failed to query logs: %v\n", err)
		return nil, 0, fmt.Errorf("failed to fetch logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.SuperAdminAuditLog
	for rows.Next() {
		var l domain.SuperAdminAuditLog
		var targetTypeVal, userAgent sql.NullString
		var targetIDVal sql.NullInt64
		var detailsJSON []byte

		err := rows.Scan(
			&l.ID, &l.SuperAdminID, &l.Action, &targetTypeVal, &targetIDVal,
			&detailsJSON, &l.IPAddress, &userAgent, &l.CreatedAt,
		)
		if err != nil {
			log.Printf("[SuperAdminAuditRepository] ERROR: Failed to scan log: %v\n", err)
			return nil, 0, fmt.Errorf("failed to scan log: %w", err)
		}

		if targetTypeVal.Valid {
			l.TargetType = &targetTypeVal.String
		}
		if targetIDVal.Valid {
			l.TargetID = &targetIDVal.Int64
		}
		if userAgent.Valid {
			l.UserAgent = userAgent.String
		}

		if len(detailsJSON) > 0 {
			var details map[string]interface{}
			if err := json.Unmarshal(detailsJSON, &details); err == nil {
				l.Details = details
			}
		}

		logs = append(logs, l)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[SuperAdminAuditRepository] ERROR: Row iteration error: %v\n", err)
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return logs, total, nil
}

// GetRecentLogsForSuperAdmin retrieves recent logs for a specific super admin
func (r *SuperAdminAuditRepository) GetRecentLogsForSuperAdmin(ctx context.Context, superAdminID int64, limit int) ([]domain.SuperAdminAuditLog, error) {
	log.Printf("[SuperAdminAuditRepository] GetRecentLogsForSuperAdmin: superAdminID=%d, limit=%d\n", superAdminID, limit)

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, super_admin_id, action, target_type, target_id, details, ip_address, user_agent, created_at
		FROM super_admin_audit_logs
		WHERE super_admin_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, superAdminID, limit)

	if err != nil {
		log.Printf("[SuperAdminAuditRepository] ERROR: Failed to query logs: %v\n", err)
		return nil, fmt.Errorf("failed to fetch logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.SuperAdminAuditLog
	for rows.Next() {
		var l domain.SuperAdminAuditLog
		var targetType, userAgent sql.NullString
		var targetID sql.NullInt64
		var detailsJSON []byte

		err := rows.Scan(
			&l.ID, &l.SuperAdminID, &l.Action, &targetType, &targetID,
			&detailsJSON, &l.IPAddress, &userAgent, &l.CreatedAt,
		)
		if err != nil {
			log.Printf("[SuperAdminAuditRepository] ERROR: Failed to scan log: %v\n", err)
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}

		if targetType.Valid {
			l.TargetType = &targetType.String
		}
		if targetID.Valid {
			l.TargetID = &targetID.Int64
		}
		if userAgent.Valid {
			l.UserAgent = userAgent.String
		}

		if len(detailsJSON) > 0 {
			var details map[string]interface{}
			if err := json.Unmarshal(detailsJSON, &details); err == nil {
				l.Details = details
			}
		}

		logs = append(logs, l)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[SuperAdminAuditRepository] ERROR: Row iteration error: %v\n", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return logs, nil
}

// GetActionCounts retrieves counts of each action type for statistics
func (r *SuperAdminAuditRepository) GetActionCounts(ctx context.Context, superAdminID *int64) (map[string]int64, error) {
	log.Printf("[SuperAdminAuditRepository] GetActionCounts\n")

	var query string
	var args []interface{}

	if superAdminID != nil {
		query = `SELECT action, COUNT(*) as count FROM super_admin_audit_logs WHERE super_admin_id = $1 GROUP BY action`
		args = []interface{}{*superAdminID}
	} else {
		query = `SELECT action, COUNT(*) as count FROM super_admin_audit_logs GROUP BY action`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("[SuperAdminAuditRepository] ERROR: Failed to query action counts: %v\n", err)
		return nil, fmt.Errorf("failed to get action counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int64)
	for rows.Next() {
		var action string
		var count int64
		if err := rows.Scan(&action, &count); err != nil {
			log.Printf("[SuperAdminAuditRepository] ERROR: Failed to scan count: %v\n", err)
			return nil, fmt.Errorf("failed to scan count: %w", err)
		}
		counts[action] = count
	}

	return counts, nil
}
