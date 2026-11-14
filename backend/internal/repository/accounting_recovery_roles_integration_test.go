//go:build integration

package repository

import "testing"

func TestAccountingRecoveryCorePrivilegesAreNarrow(t *testing.T) {
	db := openRLSRolePermissionDB(t, safeRLSDatabaseURL(t, "POS_RLS_APP_URL"))
	var role string
	if err := db.QueryRow(`SELECT current_user`).Scan(&role); err != nil || role != "pos_app" {
		t.Fatalf("restricted core identity required: %v", err)
	}
	permissions := map[string]map[string]bool{
		"accounting_event_outbox":              {"SELECT": true, "INSERT": true, "UPDATE": true, "DELETE": false, "TRUNCATE": false, "TRIGGER": false},
		"accounting_inventory_movement_events": {"SELECT": true, "INSERT": true, "UPDATE": false, "DELETE": false, "TRUNCATE": false, "TRIGGER": false},
		"general_ledger":                       {"INSERT": false, "UPDATE": false, "DELETE": false},
	}
	for table, operations := range permissions {
		for operation, want := range operations {
			var got bool
			if err := db.QueryRow(`SELECT has_table_privilege(current_user,$1,$2)`, table, operation).Scan(&got); err != nil || got != want {
				t.Fatalf("%s %s granted=%v want=%v err=%v", table, operation, got, want, err)
			}
		}
	}
	var sequence bool
	if err := db.QueryRow(`SELECT has_sequence_privilege(current_user,'accounting_event_outbox_id_seq','USAGE')`).Scan(&sequence); err != nil || !sequence {
		t.Fatalf("outbox identity sequence unavailable: %v", err)
	}
}
