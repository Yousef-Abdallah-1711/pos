package repository

import (
	"context"
	"database/sql"
	"fmt"

	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/database"
)

const userSettingsColumns = `id, user_id, language, theme, primary_color, secondary_color, accent_color,
	notifications_enabled, email_notifications, sms_notifications, created_at, updated_at`

// UserSettingsRepository handles user settings database operations.
type UserSettingsRepository struct {
	db *sql.DB
}

// NewUserSettingsRepository creates new user settings repository.
func NewUserSettingsRepository(db *sql.DB) *UserSettingsRepository {
	return &UserSettingsRepository{db: db}
}

// GetByUserID retrieves settings for a tenant user and creates defaults if absent.
func (r *UserSettingsRepository) GetByUserID(ctx context.Context, tenantID, userID int64) (*domain.UserSettings, error) {
	var settings domain.UserSettings
	err := r.withUserTransaction(ctx, tenantID, userID, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO user_settings (user_id)
			VALUES ($1)
			ON CONFLICT (user_id) DO NOTHING`, userID)
		if err != nil {
			return fmt.Errorf("create default settings: %w", err)
		}
		return scanUserSettings(tx.QueryRowContext(ctx, `SELECT `+userSettingsColumns+`
			FROM user_settings WHERE user_id = $1`, userID), &settings)
	})
	if err != nil {
		return nil, fmt.Errorf("fetch user settings: %w", err)
	}
	return &settings, nil
}

// UpdateLanguage updates a user's language preference.
func (r *UserSettingsRepository) UpdateLanguage(ctx context.Context, tenantID, userID int64, language string) (*domain.UserSettings, error) {
	if language != "en" && language != "ar" {
		return nil, fmt.Errorf("invalid language: %s", language)
	}
	return r.updateSettings(ctx, tenantID, userID, "update language",
		`UPDATE user_settings SET language = $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2 RETURNING `+userSettingsColumns,
		language, userID)
}

// UpdateTheme updates a user's theme preference.
func (r *UserSettingsRepository) UpdateTheme(ctx context.Context, tenantID, userID int64, theme string) (*domain.UserSettings, error) {
	if theme != "light" && theme != "dark" && theme != "system" {
		return nil, fmt.Errorf("invalid theme: %s", theme)
	}
	return r.updateSettings(ctx, tenantID, userID, "update theme",
		`UPDATE user_settings SET theme = $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2 RETURNING `+userSettingsColumns,
		theme, userID)
}

// UpdateColors updates a user's dashboard color scheme.
func (r *UserSettingsRepository) UpdateColors(ctx context.Context, tenantID, userID int64, primaryColor, secondaryColor, accentColor string) (*domain.UserSettings, error) {
	if !isValidHexColor(primaryColor) || !isValidHexColor(secondaryColor) || !isValidHexColor(accentColor) {
		return nil, fmt.Errorf("invalid color format provided")
	}
	return r.updateSettings(ctx, tenantID, userID, "update colors",
		`UPDATE user_settings SET primary_color = $1, secondary_color = $2, accent_color = $3,
		updated_at = CURRENT_TIMESTAMP WHERE user_id = $4 RETURNING `+userSettingsColumns,
		primaryColor, secondaryColor, accentColor, userID)
}

// UpdateNotifications updates a user's notification preferences.
func (r *UserSettingsRepository) UpdateNotifications(ctx context.Context, tenantID, userID int64, enabled, emailNotif, smsNotif bool) (*domain.UserSettings, error) {
	return r.updateSettings(ctx, tenantID, userID, "update notifications",
		`UPDATE user_settings SET notifications_enabled = $1, email_notifications = $2,
		sms_notifications = $3, updated_at = CURRENT_TIMESTAMP WHERE user_id = $4 RETURNING `+userSettingsColumns,
		enabled, emailNotif, smsNotif, userID)
}

// LogAuditChange logs a settings change inside the user's tenant transaction.
func (r *UserSettingsRepository) LogAuditChange(ctx context.Context, tenantID, userID, changedBy int64, settingName string, oldValue, newValue *string, ipAddress string) error {
	return r.withUserTransaction(ctx, tenantID, userID, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO user_settings_audit_log
			(user_id, setting_name, old_value, new_value, changed_by, ip_address)
			VALUES ($1, $2, $3, $4, $5, $6)`, userID, settingName, oldValue, newValue, changedBy, ipAddress)
		return err
	})
}

// RecordPasswordChange records a password change in the user's tenant transaction.
func (r *UserSettingsRepository) RecordPasswordChange(ctx context.Context, tenantID, userID int64, passwordHash, ipAddress string) error {
	return r.withUserTransaction(ctx, tenantID, userID, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO password_change_history (user_id, password_hash, ip_address)
			VALUES ($1, $2, $3)`, userID, passwordHash, ipAddress)
		return err
	})
}

// GetPasswordChangeHistory retrieves password change history for a tenant user.
func (r *UserSettingsRepository) GetPasswordChangeHistory(ctx context.Context, tenantID, userID int64, limit int) ([]domain.PasswordChangeHistory, error) {
	var history []domain.PasswordChangeHistory
	err := r.withUserTransaction(ctx, tenantID, userID, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, `SELECT id, user_id, changed_at, ip_address
			FROM password_change_history WHERE user_id = $1 ORDER BY changed_at DESC LIMIT $2`, userID, limit)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var record domain.PasswordChangeHistory
			if err := rows.Scan(&record.ID, &record.UserID, &record.ChangedAt, &record.IPAddress); err != nil {
				return fmt.Errorf("scan password history: %w", err)
			}
			history = append(history, record)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, fmt.Errorf("fetch password history: %w", err)
	}
	return history, nil
}

func (r *UserSettingsRepository) updateSettings(ctx context.Context, tenantID, userID int64, operation, query string, args ...any) (*domain.UserSettings, error) {
	var settings domain.UserSettings
	err := r.withUserTransaction(ctx, tenantID, userID, func(tx *sql.Tx) error {
		return scanUserSettings(tx.QueryRowContext(ctx, query, args...), &settings)
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return &settings, nil
}

func (r *UserSettingsRepository) withUserTransaction(ctx context.Context, tenantID, userID int64, operation func(*sql.Tx) error) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user id %d", userID)
	}
	return database.WithTenantTx(ctx, r.db, tenantID, 0, operation)
}

func scanUserSettings(row interface{ Scan(...any) error }, settings *domain.UserSettings) error {
	return row.Scan(&settings.ID, &settings.UserID, &settings.Language, &settings.Theme,
		&settings.PrimaryColor, &settings.SecondaryColor, &settings.AccentColor,
		&settings.NotificationsEnabled, &settings.EmailNotifications, &settings.SMSNotifications,
		&settings.CreatedAt, &settings.UpdatedAt)
}

// isValidHexColor validates hex color format.
func isValidHexColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	for i := 1; i < len(color); i++ {
		if (color[i] < '0' || color[i] > '9') && (color[i] < 'a' || color[i] > 'f') && (color[i] < 'A' || color[i] > 'F') {
			return false
		}
	}
	return true
}
