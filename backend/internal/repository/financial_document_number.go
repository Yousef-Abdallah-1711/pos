package repository

import (
	"database/sql"
	"fmt"
)

// lockFinancialDocumentNumberTx serializes one document-number series inside
// the caller's transaction. The lock key preserves the existing tenant,
// restaurant, family and year scope without changing historical numbers.
func lockFinancialDocumentNumberTx(tx *sql.Tx, family string, tenantID, restaurantID int, prefix string) error {
	if tx == nil {
		return fmt.Errorf("transaction is required to allocate %s number", family)
	}
	lockKey := fmt.Sprintf("financial-document-number:%s:%d:%d:%s", family, tenantID, restaurantID, prefix)
	var lockAcquired int
	if err := tx.QueryRow(`
		SELECT 1 FROM (SELECT pg_advisory_xact_lock(hashtextextended($1, 0))) AS number_lock
	`, lockKey).Scan(&lockAcquired); err != nil {
		return fmt.Errorf("lock %s document-number series: %w", family, err)
	}
	return nil
}
