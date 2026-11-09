package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

// Session settings read by row-level-security policies. They are always set with
// set_config(..., is_local => true), i.e. scoped to ONE transaction, so a pooled connection can
// never carry a tenant context from one request into the next.
const (
	TenantSetting     = "app.tenant_id"
	RestaurantSetting = "app.restaurant_id"
	UserSetting       = "app.user_id"
)

// WithTenantTx runs fn inside a transaction whose RLS context is the given tenant (and restaurant).
// The context is transaction-local: it disappears on COMMIT or ROLLBACK, before the connection
// returns to the pool. tenantID must be positive; restaurantID may be 0 for tenant-wide work.
//
// This is the only supported way to touch RLS-protected tables. Policies compare against
// current_setting('app.tenant_id', true), which is NULL (no rows, writes rejected) when it is unset.
func WithTenantTx(ctx context.Context, db *sql.DB, tenantID, restaurantID int64, fn func(tx *sql.Tx) error) (err error) {
	return withTenantContextTx(ctx, db, tenantID, restaurantID, 0, fn)
}

// WithTenantUserTx sets tenant, restaurant, and user context for user-owned RLS rows.
func WithTenantUserTx(ctx context.Context, db *sql.DB, tenantID, restaurantID, userID int64, fn func(tx *sql.Tx) error) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user id %d", userID)
	}
	return withTenantContextTx(ctx, db, tenantID, restaurantID, userID, fn)
}

func withTenantContextTx(ctx context.Context, db *sql.DB, tenantID, restaurantID, userID int64, fn func(tx *sql.Tx) error) (err error) {
	if db == nil {
		return errors.New("database is not initialized")
	}
	if tenantID <= 0 {
		return fmt.Errorf("invalid tenant id %d", tenantID)
	}
	if restaurantID < 0 {
		return fmt.Errorf("invalid restaurant id %d", restaurantID)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tenant transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	restaurant := ""
	if restaurantID > 0 {
		restaurant = strconv.FormatInt(restaurantID, 10)
	}
	user := ""
	if userID > 0 {
		user = strconv.FormatInt(userID, 10)
	}
	if _, err = tx.ExecContext(ctx, `SELECT set_config($1, $2, true), set_config($3, $4, true), set_config($5, $6, true)`,
		TenantSetting, strconv.FormatInt(tenantID, 10), RestaurantSetting, restaurant, UserSetting, user); err != nil {
		return fmt.Errorf("set tenant context: %w", err)
	}

	if err = fn(tx); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tenant transaction: %w", err)
	}
	return nil
}
