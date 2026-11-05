package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAuthAuditRepositoryPersistsLoginAndLogoutEvents(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repository := NewAuthAuditRepository(db)
	statement := regexp.QuoteMeta(`
		INSERT INTO user_auth_audit_logs (tenant_id, user_id, action, status, ip_address, user_agent)
		VALUES (NULLIF($1, 0), NULLIF($2, 0), $3, $4, NULLIF($5, ''), NULLIF($6, ''))`)
	for _, action := range []string{"login", "logout", "refresh"} {
		mock.ExpectExec(statement).
			WithArgs(int64(4), int64(9), action, "success", "192.0.2.8", "test-agent").
			WillReturnResult(sqlmock.NewResult(1, 1))
		if err := repository.RecordAuthEvent(context.Background(), 4, 9, action, "success", "192.0.2.8", "test-agent"); err != nil {
			t.Fatalf("record %s event: %v", action, err)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAuthAuditRepositoryRejectsUnknownAction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repository := NewAuthAuditRepository(db)
	if err := repository.RecordAuthEvent(context.Background(), 0, 0, "password-dump", "success", "", ""); err == nil {
		t.Fatal("unknown auth audit action should be rejected before database access")
	}
	if err := repository.RecordAuthEvent(context.Background(), 0, 0, "login", "failed", "", ""); err == nil {
		t.Fatal("unsupported auth audit status should be rejected before database access")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
