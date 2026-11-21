package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// JournalGLReconciliation summarizes read-only journal-to-ledger integrity for
// one authorized tenant and restaurant scope.
type JournalGLReconciliation struct {
	TenantID                int
	RestaurantID            int
	JournalCount            int64
	JournalLinesChecked     int64
	GLRowCount              int64
	PostedLinesMissingGL    int64
	DuplicateOrInvalidGL    int64
	GLScopeMismatch         int64
	GLAccountMismatch       int64
	GLAmountMismatch        int64
	GLDateMismatch          int64
	InvalidLedgerLinks      int64
	UnbalancedJournals      int64
	NonPostedJournalsWithGL int64
	Status                  string
}

// JournalGLReconciliationLinesQuery is the production snapshot query used to
// compare scoped journal lines with their linked immutable ledger effects.
const JournalGLReconciliationLinesQuery = `
		WITH scoped_lines AS (
			SELECT je.id AS journal_entry_id, je.status, je.entry_date,
			       je.fiscal_period, je.voided_at, je.void_reason, jel.id AS journal_line_id,
			       jel.tenant_id AS line_tenant_id, jel.account_id,
			       jel.debit_amount, jel.credit_amount,
			       coa.tenant_id AS account_tenant_id, coa.restaurant_id AS account_restaurant_id
			FROM journal_entries je
			JOIN journal_entry_lines jel ON jel.journal_entry_id = je.id
			JOIN chart_of_accounts coa ON coa.id = jel.account_id
			WHERE je.tenant_id = $1 AND je.restaurant_id = $2
		), line_effects AS (
			SELECT sl.*,
			       COUNT(gl.id) AS gl_count,
			       COUNT(gl.id) FILTER (
				       WHERE gl.tenant_id = $1 AND gl.restaurant_id = $2
				         AND gl.account_id = sl.account_id
			         AND gl.debit_amount = sl.debit_amount
			         AND gl.credit_amount = sl.credit_amount
			         AND gl.entry_date = sl.entry_date
			         AND gl.fiscal_period = COALESCE(sl.fiscal_period, to_char(sl.entry_date, 'YYYY-MM'))
			       ) AS original_matches,
			       COUNT(gl.id) FILTER (
				       WHERE sl.status = 'void'
			         AND sl.voided_at IS NOT NULL
			         AND btrim(COALESCE(sl.void_reason, '')) <> ''
			         AND gl.tenant_id = $1 AND gl.restaurant_id = $2
			         AND gl.account_id = sl.account_id
			         AND gl.debit_amount = sl.credit_amount
			         AND gl.credit_amount = sl.debit_amount
			         AND gl.entry_date = sl.voided_at::date
			         AND gl.fiscal_period = to_char(sl.voided_at, 'YYYY-MM')
			       ) AS reversal_matches,
			       COUNT(gl.id) FILTER (
				       WHERE gl.id IS NOT NULL AND (gl.tenant_id <> $1 OR gl.restaurant_id <> $2)
			       ) AS scope_mismatches,
			       COUNT(gl.id) FILTER (
				       WHERE gl.id IS NOT NULL AND gl.account_id <> sl.account_id
			       ) AS account_mismatches,
			       COUNT(gl.id) FILTER (
				       WHERE gl.debit_amount = sl.debit_amount AND gl.credit_amount = sl.credit_amount
			          OR (sl.status = 'void' AND gl.debit_amount = sl.credit_amount AND gl.credit_amount = sl.debit_amount)
			       ) AS amount_matches,
			       COUNT(gl.id) FILTER (
				       WHERE gl.id IS NOT NULL
			         AND gl.tenant_id = $1 AND gl.restaurant_id = $2
			         AND gl.account_id = sl.account_id
			         AND ((gl.debit_amount = sl.debit_amount AND gl.credit_amount = sl.credit_amount
			               AND (gl.entry_date <> sl.entry_date
			                 OR gl.fiscal_period <> COALESCE(sl.fiscal_period, to_char(sl.entry_date, 'YYYY-MM'))))
			           OR (sl.status = 'void' AND gl.debit_amount = sl.credit_amount
			               AND gl.credit_amount = sl.debit_amount
			               AND (gl.entry_date <> sl.voided_at::date
			                 OR gl.fiscal_period <> to_char(sl.voided_at, 'YYYY-MM'))))
			       ) AS date_mismatches
			FROM scoped_lines sl
			LEFT JOIN general_ledger gl
			  ON gl.journal_entry_id = sl.journal_entry_id
			 AND gl.journal_line_id = sl.journal_line_id
			GROUP BY sl.journal_entry_id, sl.status, sl.entry_date, sl.fiscal_period,
		         sl.voided_at, sl.void_reason, sl.journal_line_id, sl.line_tenant_id, sl.account_id,
		         sl.debit_amount, sl.credit_amount, sl.account_tenant_id,
		         sl.account_restaurant_id
		)
		SELECT COUNT(*),
		       COUNT(*) FILTER (WHERE status = 'posted' AND gl_count = 0),
		       COUNT(*) FILTER (
		         WHERE (status = 'posted' AND (gl_count <> 1 OR original_matches <> 1))
	            OR (status = 'void' AND (gl_count <> 2 OR original_matches <> 1 OR reversal_matches <> 1))
	       ),
	       COUNT(DISTINCT journal_entry_id) FILTER (WHERE status NOT IN ('posted', 'void') AND gl_count > 0),
	       COALESCE(SUM(scope_mismatches), 0),
	       COALESCE(SUM(account_mismatches), 0),
	       COALESCE(SUM(CASE WHEN status IN ('posted', 'void') THEN gl_count - amount_matches ELSE 0 END), 0),
	       COALESCE(SUM(date_mismatches), 0),
	       COUNT(*) FILTER (WHERE line_tenant_id <> $1),
	       COUNT(*) FILTER (WHERE account_tenant_id <> $1 OR account_restaurant_id <> $2)
	FROM line_effects`

// JournalGLReconciliationJournalCountQuery counts source journals in scope.
const JournalGLReconciliationJournalCountQuery = `SELECT COUNT(*) FROM journal_entries WHERE tenant_id = $1 AND restaurant_id = $2`

// JournalGLReconciliationLedgerCountQuery counts ledger effects in scope.
const JournalGLReconciliationLedgerCountQuery = `SELECT COUNT(*) FROM general_ledger WHERE tenant_id = $1 AND restaurant_id = $2`

// JournalGLReconciliationInvalidLinksQuery detects scoped orphan or mismatched ledger links.
const JournalGLReconciliationInvalidLinksQuery = `
		SELECT COUNT(*)
		FROM general_ledger gl
		LEFT JOIN journal_entries je ON je.id = gl.journal_entry_id
		LEFT JOIN journal_entry_lines jel
		  ON jel.id = gl.journal_line_id AND jel.journal_entry_id = gl.journal_entry_id
		WHERE ((gl.tenant_id = $1 AND gl.restaurant_id = $2)
		   OR (je.tenant_id = $1 AND je.restaurant_id = $2))
		AND (je.id IS NULL OR jel.id IS NULL
		  OR gl.tenant_id <> je.tenant_id OR gl.restaurant_id <> je.restaurant_id
		  OR jel.tenant_id <> je.tenant_id OR gl.account_id <> jel.account_id)`

// JournalGLReconciliationHeadersQuery counts posted journal header/line total disagreements.
const JournalGLReconciliationHeadersQuery = `
		SELECT COUNT(*)
		FROM (
			SELECT je.id
			FROM journal_entries je
			LEFT JOIN journal_entry_lines jel ON jel.journal_entry_id = je.id
			WHERE je.tenant_id = $1 AND je.restaurant_id = $2
			GROUP BY je.id
			HAVING je.status IN ('posted', 'void')
			   AND (COALESCE(SUM(jel.debit_amount), 0) <> COALESCE(SUM(jel.credit_amount), 0)
			     OR je.total_debit <> COALESCE(SUM(jel.debit_amount), 0)
			     OR je.total_credit <> COALESCE(SUM(jel.credit_amount), 0))
		) AS mismatched_journals`

// ReconcileJournalGL checks source journal lines against their immutable GL
// effects. It performs no writes and uses exact NUMERIC comparisons in SQL.
// Legacy journal types are classified by posting status and line linkage rather
// than rejected by entry_type, so supported historical/manual journals remain
// eligible for reconciliation.
func (r *JournalEntryRepository) ReconcileJournalGL(
	ctx context.Context,
	tenantID int,
	restaurantID int,
) (JournalGLReconciliation, error) {
	if tenantID <= 0 || restaurantID <= 0 {
		return JournalGLReconciliation{}, fmt.Errorf("tenant and restaurant scope are required")
	}
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return JournalGLReconciliation{}, fmt.Errorf("begin read-only reconciliation snapshot: %w", err)
	}
	defer tx.Rollback()
	result := JournalGLReconciliation{TenantID: tenantID, RestaurantID: restaurantID}

	linesQuery := JournalGLReconciliationLinesQuery
	var lineTenantMismatches, accountScopeMismatches int64
	if err := tx.QueryRowContext(ctx, linesQuery, tenantID, restaurantID).Scan(
		&result.JournalLinesChecked,
		&result.PostedLinesMissingGL,
		&result.DuplicateOrInvalidGL,
		&result.NonPostedJournalsWithGL,
		&result.GLScopeMismatch,
		&result.GLAccountMismatch,
		&result.GLAmountMismatch,
		&result.GLDateMismatch,
		&lineTenantMismatches,
		&accountScopeMismatches,
	); err != nil {
		return JournalGLReconciliation{}, fmt.Errorf("reconcile journal lines to general ledger: %w", err)
	}
	result.GLScopeMismatch += lineTenantMismatches
	result.GLAccountMismatch += accountScopeMismatches
	if err := tx.QueryRowContext(ctx, JournalGLReconciliationJournalCountQuery, tenantID, restaurantID).Scan(&result.JournalCount); err != nil {
		return JournalGLReconciliation{}, fmt.Errorf("count scoped journals: %w", err)
	}
	if err := tx.QueryRowContext(ctx, JournalGLReconciliationLedgerCountQuery, tenantID, restaurantID).Scan(&result.GLRowCount); err != nil {
		return JournalGLReconciliation{}, fmt.Errorf("count scoped ledger rows: %w", err)
	}
	if err := tx.QueryRowContext(ctx, JournalGLReconciliationInvalidLinksQuery, tenantID, restaurantID).Scan(&result.InvalidLedgerLinks); err != nil {
		return JournalGLReconciliation{}, fmt.Errorf("reconcile ledger source links: %w", err)
	}

	if err := tx.QueryRowContext(ctx, JournalGLReconciliationHeadersQuery, tenantID, restaurantID).Scan(&result.UnbalancedJournals); err != nil {
		return JournalGLReconciliation{}, fmt.Errorf("reconcile journal header totals: %w", err)
	}
	result.Status = "reconciled"
	if result.JournalCount == 0 && result.GLRowCount == 0 {
		result.Status = "empty"
	} else if result.PostedLinesMissingGL+result.DuplicateOrInvalidGL+result.GLScopeMismatch+
		result.GLAccountMismatch+result.GLAmountMismatch+result.GLDateMismatch+
		result.InvalidLedgerLinks+result.UnbalancedJournals+result.NonPostedJournalsWithGL > 0 {
		result.Status = "discrepancies"
	}
	if err := tx.Commit(); err != nil {
		return JournalGLReconciliation{}, fmt.Errorf("finish read-only reconciliation snapshot: %w", err)
	}
	return result, nil
}
