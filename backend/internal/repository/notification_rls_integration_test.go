//go:build integration

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/database"
)

func TestNotificationsRLSWithRestrictedRuntimeRole(t *testing.T) {
	admin := openAuthAuditTestDB(t, "POS_RLS_ADMIN_URL")
	appURL := os.Getenv("POS_RLS_APP_URL")
	if appURL == "" {
		t.Skip("POS_RLS_APP_URL is not set")
	}
	schema := fmt.Sprintf("notifications_rls_%d", time.Now().UnixNano())
	if !regexp.MustCompile(`^notifications_rls_[0-9]+$`).MatchString(schema) {
		t.Fatal("generated unsafe schema name")
	}
	ctx := context.Background()
	admin.SetMaxOpenConns(1)
	if _, err := admin.ExecContext(ctx, `CREATE SCHEMA `+schema); err != nil {
		t.Fatalf("create isolated schema: %v", err)
	}
	t.Cleanup(func() {
		if _, err := admin.ExecContext(context.Background(), `DROP SCHEMA `+schema+` CASCADE`); err != nil {
			t.Errorf("drop isolated schema: %v", err)
		}
	})
	if _, err := admin.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set admin search path: %v", err)
	}
	createNotificationRLSFixture(t, ctx, admin, schema)
	app := openAuthAuditURL(t, withSearchPath(t, appURL, schema))
	app.SetMaxOpenConns(1)
	if _, err := app.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set app search path: %v", err)
	}
	assertRestrictedRuntimeRole(t, ctx, app)
	seedNotificationRLSFixture(t, ctx, admin)
	repository := NewNotificationRepository(app)
	assertNotificationRepositoryIsolation(t, ctx, admin, app, repository)
	assertNotificationNoContextAndRollback(t, ctx, admin, app)
	rollbackNotificationsRLS(t, ctx, admin)
	assertNotificationMigrationRollback(t, ctx, app)
}

func createNotificationRLSFixture(t *testing.T, ctx context.Context, admin *sql.DB, schema string) {
	t.Helper()
	for _, statement := range []string{
		`CREATE TABLE tenants (id INT PRIMARY KEY)`,
		`CREATE TABLE restaurants (id INT PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id))`,
		`CREATE TABLE users (id INT PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id), restaurant_id INT NOT NULL REFERENCES restaurants(id))`,
		`CREATE TABLE notifications (
			id SERIAL PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id), restaurant_id INT NOT NULL REFERENCES restaurants(id),
			user_id INT NOT NULL REFERENCES users(id), type VARCHAR(50) NOT NULL, module VARCHAR(50) NOT NULL,
			title VARCHAR(255) NOT NULL, message TEXT NOT NULL, description TEXT, is_read BOOLEAN DEFAULT FALSE,
			read_at TIMESTAMP, related_entity_type VARCHAR(50), related_entity_id INT, priority VARCHAR(20) DEFAULT 'normal',
			action_url VARCHAR(512), action_label VARCHAR(100), icon_name VARCHAR(50), color VARCHAR(20),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, expires_at TIMESTAMP)`,
		`GRANT USAGE ON SCHEMA ` + schema + ` TO pos_app`,
		`GRANT SELECT ON tenants, restaurants, users TO pos_app`,
		`GRANT SELECT, INSERT, UPDATE, DELETE ON notifications TO pos_app`,
		`GRANT USAGE, SELECT ON SEQUENCE notifications_id_seq TO pos_app`,
	} {
		if _, err := admin.ExecContext(ctx, statement); err != nil {
			t.Fatalf("create notification fixture: %v", err)
		}
	}
	migration, err := os.ReadFile("../../migrations/172_rls_notifications.sql")
	if err != nil {
		t.Fatalf("read notification RLS migration: %v", err)
	}
	if _, err := admin.ExecContext(ctx, string(migration)); err != nil {
		t.Fatalf("apply notification RLS migration: %v", err)
	}
}

func seedNotificationRLSFixture(t *testing.T, ctx context.Context, admin *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`INSERT INTO tenants (id) VALUES (1), (2)`,
		`INSERT INTO restaurants (id, tenant_id) VALUES (11, 1), (12, 1), (21, 2)`,
		`INSERT INTO users (id, tenant_id, restaurant_id) VALUES (101, 1, 11), (102, 1, 11), (103, 1, 12), (104, 1, 21), (202, 2, 21)`,
		`INSERT INTO notifications (tenant_id, restaurant_id, user_id, type, module, title, message, priority) VALUES
			(1, 11, 101, 'system', 'products', 'A1 first', 'one', 'normal'),
			(1, 11, 101, 'system', 'orders', 'A1 second', 'two', 'high'),
			(1, 11, 102, 'system', 'products', 'A1 other user', 'three', 'normal'),
			(1, 12, 103, 'system', 'products', 'A2 sibling', 'four', 'normal'),
			(2, 21, 202, 'system', 'products', 'Tenant B', 'five', 'normal')`,
	} {
		if _, err := admin.ExecContext(ctx, statement); err != nil {
			t.Fatalf("seed notification fixture: %v", err)
		}
	}
}

func assertNotificationRepositoryIsolation(t *testing.T, ctx context.Context, admin, app *sql.DB, repository *NotificationRepository) {
	t.Helper()
	filters := &domain.NotificationFilters{Page: 1, PageSize: 10, SortBy: "created_at; DROP TABLE notifications", SortOrder: "DESC; DELETE FROM users"}
	list, err := repository.ListNotifications(ctx, 101, 1, 11, filters)
	if err != nil {
		t.Fatalf("list notifications under owner scope: %v", err)
	}
	if list.Total != 2 || list.Unread != 2 || len(list.Notifications) != 2 {
		t.Fatalf("owner scope returned unexpected notifications: %+v", list)
	}
	if list.Notifications[0].Title != "A1 second" || list.Notifications[1].Title != "A1 first" {
		t.Fatalf("invalid sort input was not safely replaced by created_at descending: %+v", list.Notifications)
	}
	if _, err := repository.GetNotification(ctx, 101, 1, 11, 3); err == nil {
		t.Fatal("same-restaurant notification for another user must not be readable")
	}
	if err := repository.MarkAsRead(ctx, 101, 1, 11, 3); err == nil {
		t.Fatal("same-restaurant notification for another user must not be mutable")
	}
	if err := repository.DeleteNotification(ctx, 101, 1, 11, 4); err == nil {
		t.Fatal("sibling-restaurant notification must not be deletable")
	}
	if _, err := repository.GetNotification(ctx, 101, 1, 11, 5); err == nil {
		t.Fatal("foreign-tenant notification must not be readable")
	}
	if err := repository.DeleteNotification(ctx, 101, 1, 11, 5); err == nil {
		t.Fatal("foreign-tenant notification must not be deletable")
	}
	if got := notificationReadState(t, ctx, admin, 3); got {
		t.Fatal("denied update changed another user's notification")
	}
	count := notificationCount(t, ctx, admin)
	wrongOwner := &domain.Notification{TenantID: 1, RestaurantID: 11, UserID: 103, Type: "system", Module: "products", Title: "forged", Message: "forged"}
	if _, err := repository.CreateNotification(ctx, wrongOwner); err == nil {
		t.Fatal("foreign/sibling recipient notification insert must fail RLS checks")
	}
	inconsistentParent := &domain.Notification{TenantID: 1, RestaurantID: 21, UserID: 104, Type: "system", Module: "products", Title: "forged parent", Message: "forged"}
	if _, err := repository.CreateNotification(ctx, inconsistentParent); err == nil {
		t.Fatal("notification must reject a restaurant owned by a different tenant")
	}
	if got := notificationCount(t, ctx, admin); got != count {
		t.Fatalf("rejected notification insert changed row count: before=%d after=%d", count, got)
	}
	created, err := repository.CreateNotification(ctx, &domain.Notification{TenantID: 1, RestaurantID: 11, UserID: 101, Type: "system", Module: "products", Title: "created", Message: "valid"})
	if err != nil || created.ID <= 0 {
		t.Fatalf("valid scoped notification create failed: value=%+v error=%v", created, err)
	}
	if err := repository.MarkAsRead(ctx, 101, 1, 11, created.ID); err != nil {
		t.Fatalf("mark owned notification read: %v", err)
	}
	if !notificationReadState(t, ctx, admin, created.ID) {
		t.Fatal("valid mark-as-read did not persist")
	}
	if err := repository.MarkAsUnread(ctx, 101, 1, 11, created.ID); err != nil {
		t.Fatalf("mark owned notification unread: %v", err)
	}
	stats, err := repository.GetNotificationStats(ctx, 101, 1, 11)
	if err != nil || stats.Total != 3 || stats.Unread != 3 || stats.ByModule["products"] != 2 {
		t.Fatalf("notification statistics escaped scope or were incorrect: stats=%+v err=%v", stats, err)
	}
}

func assertNotificationNoContextAndRollback(t *testing.T, ctx context.Context, admin, app *sql.DB) {
	t.Helper()
	var visible int
	if err := app.QueryRowContext(ctx, `SELECT COUNT(*) FROM notifications`).Scan(&visible); err != nil {
		t.Fatalf("unscoped restricted-role read: %v", err)
	}
	if visible != 0 {
		t.Fatalf("unscoped restricted-role read returned %d rows", visible)
	}
	if _, err := app.ExecContext(ctx, `INSERT INTO notifications (tenant_id, restaurant_id, user_id, type, module, title, message) VALUES (1,11,101,'system','system','unscoped','unscoped')`); err == nil {
		t.Fatal("notification writes without tenant, restaurant, and user context must fail")
	}
	before := notificationCount(t, ctx, admin)
	rollback := errors.New("force notification rollback")
	err := database.WithTenantUserTx(ctx, app, 1, 11, 101, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO notifications (tenant_id, restaurant_id, user_id, type, module, title, message) VALUES (1,11,101,'system','system','rollback','rollback')`); err != nil {
			return err
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatalf("transaction should return forced error and roll back, got %v", err)
	}
	if after := notificationCount(t, ctx, admin); after != before {
		t.Fatalf("rolled-back notification persisted: before=%d after=%d", before, after)
	}
	if err := app.QueryRowContext(ctx, `SELECT COUNT(*) FROM notifications`).Scan(&visible); err != nil || visible != 0 {
		t.Fatalf("transaction-local context leaked to pooled connection: visible=%d err=%v", visible, err)
	}
}

func rollbackNotificationsRLS(t *testing.T, ctx context.Context, admin *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`DROP POLICY IF EXISTS notifications_user_tenant_restaurant_scope ON notifications`,
		`ALTER TABLE notifications NO FORCE ROW LEVEL SECURITY`,
		`ALTER TABLE notifications DISABLE ROW LEVEL SECURITY`,
	} {
		if _, err := admin.ExecContext(ctx, statement); err != nil {
			t.Fatalf("rollback notification RLS migration: %v", err)
		}
	}
}

func assertNotificationMigrationRollback(t *testing.T, ctx context.Context, app *sql.DB) {
	t.Helper()
	var visible int
	if err := app.QueryRowContext(ctx, `SELECT COUNT(*) FROM notifications`).Scan(&visible); err != nil {
		t.Fatalf("query after isolated rollback: %v", err)
	}
	if visible != 6 {
		t.Fatalf("rollback did not restore unrestricted fixture visibility: got %d rows", visible)
	}
}

func notificationCount(t *testing.T, ctx context.Context, db *sql.DB) int {
	t.Helper()
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM notifications`).Scan(&count); err != nil {
		t.Fatalf("count notifications: %v", err)
	}
	return count
}

func notificationReadState(t *testing.T, ctx context.Context, db *sql.DB, id int) bool {
	t.Helper()
	var isRead bool
	if err := db.QueryRowContext(ctx, `SELECT is_read FROM notifications WHERE id = $1`, id).Scan(&isRead); err != nil {
		t.Fatalf("read notification state: %v", err)
	}
	return isRead
}
