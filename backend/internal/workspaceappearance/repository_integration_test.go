package workspaceappearance

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"testing"

	_ "github.com/lib/pq"
)

// This test requires the disposable fixture described in the migration test.
func TestRepositoryIsolationAndRevision(t *testing.T) {
	dsn := os.Getenv("WORKSPACE_APPEARANCE_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("disposable PostgreSQL fixture not configured")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := Repository{DB: db}
	ctx := context.Background()
	initial, err := r.Get(ctx, 1, 11)
	if err != nil || initial.Revision != 0 {
		t.Fatalf("initial: %+v %v", initial, err)
	}
	created, err := r.Save(ctx, 1, 11, 0, json.RawMessage(`{"colors":{"light":{"primary":"#722116"}}}`))
	if err != nil || created.Revision != 1 {
		t.Fatalf("create: %+v %v", created, err)
	}
	other, err := r.Get(ctx, 2, 22)
	if err != nil || other.Revision != 0 {
		t.Fatalf("tenant isolation: %+v %v", other, err)
	}
	if _, err := r.Save(ctx, 1, 11, 0, json.RawMessage(`{}`)); !errors.Is(err, ErrConflict) {
		t.Fatalf("stale create: %v", err)
	}
	if _, err := r.Save(ctx, 1, 22, 1, json.RawMessage(`{}`)); err == nil {
		t.Fatal("cross-tenant actor accepted")
	}
	reset, err := r.Save(ctx, 1, 11, 1, json.RawMessage(`{}`))
	if err != nil || reset.Revision != 2 {
		t.Fatalf("reset: %+v %v", reset, err)
	}
	if _, err := r.Save(ctx, 1, 11, 1, json.RawMessage(`{}`)); !errors.Is(err, ErrConflict) {
		t.Fatalf("stale update: %v", err)
	}
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() { defer wg.Done(); _, err := r.Save(ctx, 1, 11, 2, json.RawMessage(`{}`)); results <- err }()
	}
	wg.Wait()
	close(results)
	successes, conflicts := 0, 0
	for err := range results {
		if err == nil {
			successes++
		} else if errors.Is(err, ErrConflict) {
			conflicts++
		} else {
			t.Fatalf("concurrent save: %v", err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("concurrency: successes=%d conflicts=%d", successes, conflicts)
	}
	adminDSN := os.Getenv("WORKSPACE_APPEARANCE_TEST_ADMIN_URL")
	if adminDSN == "" {
		t.Log("audit rollback test skipped: disposable admin URL not configured")
		return
	}
	admin, err := sql.Open("postgres", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	if _, err := admin.ExecContext(ctx, `REVOKE INSERT ON tenant_workspace_appearance_audit FROM pos_app`); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.ExecContext(ctx, `GRANT INSERT ON tenant_workspace_appearance_audit TO pos_app`); err != nil {
			t.Errorf("restore fixture grant: %v", err)
		}
	}()
	if _, err := r.Save(ctx, 1, 11, 3, json.RawMessage(`{}`)); err == nil {
		t.Fatal("save succeeded with audit unavailable")
	}
	after, err := r.Get(ctx, 1, 11)
	if err != nil || after.Revision != 3 {
		t.Fatalf("audit failure did not roll back publication: %+v %v", after, err)
	}
}
