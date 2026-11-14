//go:build integration

package usecase

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
)

func redriveMigrationFixture(t *testing.T) *sql.Tx {
	t.Helper()
	tx, err := openAccountingEventTestDB(t).Begin()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })
	schema := fmt.Sprintf("cfr004_redrive_rollback_%d", time.Now().UnixNano())
	// SQL identifiers cannot use value placeholders; quote this generated test name.
	if _, err := tx.Exec("CREATE SCHEMA " + pq.QuoteIdentifier(schema)); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(`SELECT set_config('search_path',$1,true)`, schema+",pg_catalog"); err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(`CREATE TABLE accounting_event_outbox (
  id BIGINT PRIMARY KEY,tenant_id BIGINT,restaurant_id BIGINT,updated_at TIMESTAMP,status TEXT)`)
	if err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../migrations/188_accounting_event_redrive_audit.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(string(migration)); err != nil {
		t.Fatal(err)
	}
	return tx
}

func redriveRollbackSQL(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("../../scripts/rollback_188_accounting_event_redrive_audit.sql")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	var statements []string
	for _, line := range lines {
		// Keep the guarded rollback inside this test's existing isolated transaction.
		if strings.TrimSpace(line) != "BEGIN;" && strings.TrimSpace(line) != "COMMIT;" {
			statements = append(statements, line)
		}
	}
	return strings.Join(statements, "\n")
}

func TestRedriveAuditMigrationRollbackPreservesHistory(t *testing.T) {
	for _, withHistory := range []bool{false, true} {
		t.Run(fmt.Sprint(withHistory), func(t *testing.T) {
			tx := redriveMigrationFixture(t)
			if withHistory {
				_, err := tx.Exec(`INSERT INTO accounting_event_outbox(id) VALUES(1);
     INSERT INTO accounting_event_redrive_audit(outbox_id,actor_user_id,previous_attempts,previous_error,reason)
     VALUES(1,91,8,'synthetic failure','synthetic audited correction')`)
				if err != nil {
					t.Fatal(err)
				}
			}
			if _, err := tx.Exec("SAVEPOINT before_rollback"); err != nil {
				t.Fatal(err)
			}
			_, err := tx.Exec(redriveRollbackSQL(t))
			if !withHistory {
				if err != nil {
					t.Fatal(err)
				}
				var table sql.NullString
				if err := tx.QueryRow(`SELECT to_regclass('accounting_event_redrive_audit')::text`).Scan(&table); err != nil || table.Valid {
					t.Fatalf("empty rollback left table: %v err=%v", table, err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), "audit history exists") {
				t.Fatalf("history rollback was not refused: %v", err)
			}
			if _, err := tx.Exec("ROLLBACK TO SAVEPOINT before_rollback"); err != nil {
				t.Fatal(err)
			}
			var count int
			if err := tx.QueryRow(`SELECT count(*) FROM accounting_event_redrive_audit`).Scan(&count); err != nil || count != 1 {
				t.Fatalf("history was not preserved: count=%d err=%v", count, err)
			}
		})
	}
}
