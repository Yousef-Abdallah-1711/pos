package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"pos-saas/internal/domain"
)

func TestSuperAdminAuditRepositoryPassesJSONBAsTextAndNull(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}
	defer db.Close()

	repository := NewSuperAdminAuditRepository(db)
	query := regexp.QuoteMeta(`
		INSERT INTO super_admin_audit_logs
			(super_admin_id, action, target_type, target_id, details, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`)

	mock.ExpectExec(query).
		WithArgs(int64(4), domain.AuditActionLoginSuccess, nil, nil, nil, "127.0.0.1", "test-agent").
		WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repository.LogAction(context.Background(), &domain.SuperAdminAuditLog{
		SuperAdminID: 4,
		Action:       domain.AuditActionLoginSuccess,
		IPAddress:    "127.0.0.1",
		UserAgent:    "test-agent",
	}); err != nil {
		t.Fatalf("log successful login without details: %v", err)
	}

	mock.ExpectExec(query).
		WithArgs(int64(4), domain.AuditActionLoginFailed, nil, nil, `{"reason":"invalid credentials"}`, "127.0.0.1", "test-agent").
		WillReturnResult(sqlmock.NewResult(2, 1))
	if err := repository.LogAction(context.Background(), &domain.SuperAdminAuditLog{
		SuperAdminID: 4,
		Action:       domain.AuditActionLoginFailed,
		Details:      map[string]interface{}{"reason": "invalid credentials"},
		IPAddress:    "127.0.0.1",
		UserAgent:    "test-agent",
	}); err != nil {
		t.Fatalf("log failed login with details: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL expectations: %v", err)
	}
}
