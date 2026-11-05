package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/tests/testdb"
)

func TestAuthAuditRepositoryPersistsEventsToDisposablePostgres(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	defer db.Close()
	if err := db.PingContext(context.Background()); err != nil {
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}
	marker := fmt.Sprintf("auth-audit-integration-%d", time.Now().UnixNano())
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM user_auth_audit_logs WHERE user_agent = $1`, marker) })
	repository := NewAuthAuditRepository(db)
	for _, action := range []string{"login", "refresh", "logout"} {
		if err := repository.RecordAuthEvent(context.Background(), 0, 0, action, "success", "192.0.2.10", marker); err != nil {
			t.Fatalf("persist %s event: %v", action, err)
		}
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM user_auth_audit_logs WHERE user_agent = $1`, marker).Scan(&count); err != nil {
		t.Fatalf("count persisted auth events: %v", err)
	}
	if count != 3 {
		t.Fatalf("persisted auth event count = %d, want 3", count)
	}
}
