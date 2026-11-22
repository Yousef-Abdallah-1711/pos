package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"pos-saas/internal/config"
	"pos-saas/internal/pkg/database"
)

func authLookupDatabaseConfig(base config.DatabaseConfig) (config.DatabaseConfig, error) {
	user := strings.TrimSpace(os.Getenv("AUTH_LOOKUP_DB_USER"))
	password := os.Getenv("AUTH_LOOKUP_DB_PASSWORD")
	if user != "pos_auth_lookup" || password == "" {
		return base, fmt.Errorf("dedicated AUTH_LOOKUP_DB_USER=pos_auth_lookup and AUTH_LOOKUP_DB_PASSWORD are required")
	}
	base.User = user
	base.Password = password
	base.MaxOpenConns = 5
	base.MaxIdleConns = 2
	return base, nil
}

// connectAuthLookupDatabase is deliberately opt-in. The legacy auth path remains
// the only wired login path until the identity runtime change is separately reviewed.
func connectAuthLookupDatabase(primary *sql.DB, base config.DatabaseConfig) (*sql.DB, error) {
	user := strings.TrimSpace(os.Getenv("AUTH_LOOKUP_DB_USER"))
	password := os.Getenv("AUTH_LOOKUP_DB_PASSWORD")
	if user == "" && password == "" {
		return nil, nil
	}
	if primary == nil {
		return nil, fmt.Errorf("primary database is required for the auth lookup pool")
	}
	cfg, err := authLookupDatabaseConfig(base)
	if err != nil {
		return nil, err
	}
	db, err := database.Connect(&cfg)
	if err != nil {
		return nil, fmt.Errorf("auth lookup database connection failed: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := validateAuthLookupDatabase(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func validateAuthLookupDatabase(ctx context.Context, db *sql.DB) error {
	var role string
	var privileged, databaseCreate, authSchemaUsage, authSchemaCreate bool
	var defaultReadOnly, transactionReadOnly string
	var canLookupIdentity, canLookupMemberships bool
	var inheritedRole bool
	err := db.QueryRowContext(ctx, `
		SELECT r.rolname,
		       r.rolsuper OR r.rolcreatedb OR r.rolcreaterole OR r.rolreplication OR r.rolbypassrls,
	       has_database_privilege(current_user, current_database(), 'CREATE'),
	       has_schema_privilege(current_user, 'auth_private', 'USAGE'),
	       has_schema_privilege(current_user, 'auth_private', 'CREATE'),
	       current_setting('default_transaction_read_only'),
	       current_setting('transaction_read_only'),
	       has_function_privilege(current_user, 'auth_private.lookup_verified_identity_by_email(text)', 'EXECUTE'),
	       has_function_privilege(current_user, 'auth_private.lookup_eligible_memberships(uuid)', 'EXECUTE'),
	       EXISTS (SELECT 1 FROM pg_auth_members am WHERE am.member = r.oid OR am.roleid = r.oid)
		FROM pg_roles r
		WHERE r.rolname = current_user
	`).Scan(
		&role, &privileged, &databaseCreate, &authSchemaUsage, &authSchemaCreate,
		&defaultReadOnly, &transactionReadOnly, &canLookupIdentity, &canLookupMemberships, &inheritedRole,
	)
	if err != nil {
		return fmt.Errorf("inspect auth lookup role: %w", err)
	}
	if role != "pos_auth_lookup" || privileged || databaseCreate || !authSchemaUsage || authSchemaCreate ||
		defaultReadOnly != "on" || transactionReadOnly != "on" || !canLookupIdentity || !canLookupMemberships || inheritedRole {
		return fmt.Errorf("auth lookup database must use the isolated read-only pos_auth_lookup role")
	}

	var tablePrivileges, sequencePrivileges, extraPrivateFunctions, callableDefiners int
	if err := db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = 'public'
		  AND c.relkind IN ('r', 'p', 'v', 'm', 'f')
		  AND (has_table_privilege(current_user, c.oid, 'SELECT')
		    OR has_table_privilege(current_user, c.oid, 'INSERT')
		    OR has_table_privilege(current_user, c.oid, 'UPDATE')
		    OR has_table_privilege(current_user, c.oid, 'DELETE')
		    OR has_table_privilege(current_user, c.oid, 'TRUNCATE')
		    OR has_table_privilege(current_user, c.oid, 'REFERENCES')
		    OR has_table_privilege(current_user, c.oid, 'TRIGGER'))
	`).Scan(&tablePrivileges); err != nil {
		return fmt.Errorf("inspect auth lookup table privileges: %w", err)
	}
	if err := db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = 'public'
		  AND c.relkind = 'S'
		  AND (has_sequence_privilege(current_user, c.oid, 'USAGE')
		    OR has_sequence_privilege(current_user, c.oid, 'SELECT')
		    OR has_sequence_privilege(current_user, c.oid, 'UPDATE'))
	`).Scan(&sequencePrivileges); err != nil {
		return fmt.Errorf("inspect auth lookup sequence privileges: %w", err)
	}
	if err := db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM pg_proc p
		JOIN pg_namespace n ON n.oid = p.pronamespace
		WHERE n.nspname = 'auth_private'
		  AND p.oid NOT IN (
		      'auth_private.lookup_verified_identity_by_email(text)'::regprocedure,
		      'auth_private.lookup_eligible_memberships(uuid)'::regprocedure)
		  AND has_function_privilege(current_user, p.oid, 'EXECUTE')
	`).Scan(&extraPrivateFunctions); err != nil {
		return fmt.Errorf("inspect auth lookup function privileges: %w", err)
	}
	if err := db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM pg_proc p
		JOIN pg_namespace n ON n.oid = p.pronamespace
	WHERE p.prosecdef
	  AND p.prorettype <> 'trigger'::regtype
	  AND n.nspname NOT IN ('pg_catalog', 'information_schema', 'auth_private')
	  AND has_function_privilege(current_user, p.oid, 'EXECUTE')
	`).Scan(&callableDefiners); err != nil {
		return nilSafeAuthLookupError("inspect callable security-definer functions", err)
	}
	if tablePrivileges != 0 || sequencePrivileges != 0 || extraPrivateFunctions != 0 || callableDefiners != 0 {
		return fmt.Errorf("auth lookup role has privileges beyond its two read-only lookup functions")
	}
	return nil
}

func nilSafeAuthLookupError(label string, err error) error {
	return fmt.Errorf("%s: %w", label, err)
}
