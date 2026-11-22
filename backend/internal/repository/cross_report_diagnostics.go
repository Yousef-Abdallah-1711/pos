package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"
)

const (
	CrossReportReconciled    = "RECONCILED"
	CrossReportDiscrepancies = "DISCREPANCIES"
	CrossReportPolicyPending = "POLICY_PENDING"
	CrossReportGLBalanced    = "GL_BALANCED"
)

type CrossReportDiagnostics struct {
	TenantID                     int
	RestaurantID                 int
	DateFrom                     time.Time
	DateTo                       time.Time
	Status                       string
	ProfitLossStatus             string
	TrialBalanceStatus           string
	BalanceSheetStatus           string
	CashFlowStatus               string
	JournalRevenueCents          int64
	LedgerRevenueCents           int64
	JournalExpenseCents          int64
	LedgerExpenseCents           int64
	ProfitLossAccountMismatches  []CrossReportAccountMismatch
	TrialBalanceDebitsCents      int64
	TrialBalanceCreditsCents     int64
	BalanceSheetAssetsCents      int64
	BalanceSheetLiabilitiesCents int64
	BalanceSheetEquityCents      int64
	BalanceSheetResidualCents    int64
	CoveredFiscalPeriods         int64
	OverlappingFiscalPeriods     int64
	CrossPeriodVoidCount         int64
	PolicyNotes                  []string
}

type CrossReportAccountMismatch struct {
	AccountID       int
	AccountCode     string
	AccountType     string
	JournalCents    int64
	LedgerCents     int64
	DifferenceCents int64
}

type crossReportMovement struct {
	AccountID   int
	AccountCode string
	AccountType string
	AmountCents int64
}

type FinancialIntegrityDiagnosticRepository struct {
	db *sql.DB
}

func NewFinancialIntegrityDiagnosticRepository(db *sql.DB) *FinancialIntegrityDiagnosticRepository {
	return &FinancialIntegrityDiagnosticRepository{db: db}
}

// DiagnoseCrossReportIntegrity compares independent posted-journal and GL
// aggregates for P&L accounts, then reports exact GL control totals for Trial
// Balance and Balance Sheet checks. Cash Flow remains policy-pending because
// current production classification includes account-code heuristics.
func (r *FinancialIntegrityDiagnosticRepository) DiagnoseCrossReportIntegrity(
	ctx context.Context,
	tenantID, restaurantID int,
	from, to time.Time,
) (CrossReportDiagnostics, error) {
	if tenantID <= 0 || restaurantID <= 0 || from.IsZero() || to.IsZero() || to.Before(from) {
		return CrossReportDiagnostics{}, fmt.Errorf("valid tenant, restaurant, and inclusive date scope are required")
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return CrossReportDiagnostics{}, fmt.Errorf("begin cross-report snapshot: %w", err)
	}
	defer tx.Rollback()
	result := emptyCrossReportDiagnostics(tenantID, restaurantID, from, to)

	journalMovements, err := readCrossReportMovements(ctx, tx, crossReportJournalMovementsQuery, tenantID, restaurantID, from, to)
	if err != nil {
		return CrossReportDiagnostics{}, fmt.Errorf("read posted journal P&L movements: %w", err)
	}
	ledgerMovements, err := readCrossReportMovements(ctx, tx, crossReportLedgerMovementsQuery, tenantID, restaurantID, from, to)
	if err != nil {
		return CrossReportDiagnostics{}, fmt.Errorf("read GL P&L movements: %w", err)
	}
	compareCrossReportMovements(&result, journalMovements, ledgerMovements)
	if err := readCrossReportControlTotals(ctx, tx, tenantID, restaurantID, to, &result); err != nil {
		return CrossReportDiagnostics{}, err
	}
	if err := readCrossReportPeriodContext(ctx, tx, tenantID, restaurantID, from, to, &result); err != nil {
		return CrossReportDiagnostics{}, err
	}
	if err := tx.Commit(); err != nil {
		return CrossReportDiagnostics{}, fmt.Errorf("finish cross-report snapshot: %w", err)
	}
	setCrossReportStatuses(&result)
	return result, nil
}

const crossReportJournalMovementsQuery = `
	SELECT coa.id, coa.account_code, coa.account_type,
	       COALESCE(SUM(CASE WHEN coa.account_type = 'revenue'
	         THEN jel.credit_amount - jel.debit_amount
	         ELSE jel.debit_amount - jel.credit_amount END), 0)::text
	FROM journal_entry_lines jel
	JOIN journal_entries je ON je.id = jel.journal_entry_id
	JOIN chart_of_accounts coa ON coa.id = jel.account_id
	WHERE je.tenant_id = $1 AND je.restaurant_id = $2 AND je.status = 'posted'
	  AND je.entry_date >= $3::date AND je.entry_date <= $4::date
	  AND coa.tenant_id = $1 AND coa.restaurant_id = $2
	  AND coa.account_type IN ('revenue', 'expense')
	GROUP BY coa.id, coa.account_code, coa.account_type
	ORDER BY coa.account_code, coa.id`

const crossReportLedgerMovementsQuery = `
	SELECT coa.id, coa.account_code, coa.account_type,
	       COALESCE(SUM(CASE WHEN coa.account_type = 'revenue'
	         THEN gl.credit_amount - gl.debit_amount
	         ELSE gl.debit_amount - gl.credit_amount END), 0)::text
	FROM general_ledger gl
	JOIN chart_of_accounts coa ON coa.id = gl.account_id
	WHERE gl.tenant_id = $1 AND gl.restaurant_id = $2
	  AND gl.entry_date >= $3::date AND gl.entry_date <= $4::date
	  AND coa.tenant_id = $1 AND coa.restaurant_id = $2
	  AND coa.account_type IN ('revenue', 'expense')
	GROUP BY coa.id, coa.account_code, coa.account_type
	ORDER BY coa.account_code, coa.id`

func emptyCrossReportDiagnostics(tenantID, restaurantID int, from, to time.Time) CrossReportDiagnostics {
	return CrossReportDiagnostics{
		TenantID: tenantID, RestaurantID: restaurantID,
		DateFrom: from, DateTo: to, Status: CrossReportPolicyPending,
		ProfitLossStatus:            CrossReportReconciled,
		TrialBalanceStatus:          CrossReportGLBalanced,
		BalanceSheetStatus:          CrossReportPolicyPending,
		CashFlowStatus:              CrossReportPolicyPending,
		ProfitLossAccountMismatches: make([]CrossReportAccountMismatch, 0),
		PolicyNotes:                 []string{"cash_flow_mapping_and_semantics_not_approved", "opening_and_current_earnings_treatment_unapproved"},
	}
}

func readCrossReportMovements(ctx context.Context, tx *sql.Tx, query string, tenantID, restaurantID int, from, to time.Time) (map[int]crossReportMovement, error) {
	rows, err := tx.QueryContext(ctx, query, tenantID, restaurantID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	movements := make(map[int]crossReportMovement)
	for rows.Next() {
		var movement crossReportMovement
		var amount string
		if err := rows.Scan(&movement.AccountID, &movement.AccountCode, &movement.AccountType, &amount); err != nil {
			return nil, err
		}
		movement.AmountCents, err = numericTextToCents(amount)
		if err != nil {
			return nil, fmt.Errorf("parse account %s movement: %w", movement.AccountCode, err)
		}
		movements[movement.AccountID] = movement
	}
	return movements, rows.Err()
}

func compareCrossReportMovements(result *CrossReportDiagnostics, journals, ledger map[int]crossReportMovement) {
	accountIDs := make(map[int]struct{}, len(journals)+len(ledger))
	for accountID := range journals {
		accountIDs[accountID] = struct{}{}
	}
	for accountID := range ledger {
		accountIDs[accountID] = struct{}{}
	}
	for accountID := range accountIDs {
		journal := journals[accountID]
		gl := ledger[accountID]
		if journal.AccountType == "" {
			journal.AccountType, journal.AccountCode = gl.AccountType, gl.AccountCode
		}
		if journal.AccountType == "revenue" {
			result.JournalRevenueCents += journal.AmountCents
			result.LedgerRevenueCents += gl.AmountCents
		} else {
			result.JournalExpenseCents += journal.AmountCents
			result.LedgerExpenseCents += gl.AmountCents
		}
		if journal.AmountCents != gl.AmountCents {
			result.ProfitLossAccountMismatches = append(result.ProfitLossAccountMismatches, CrossReportAccountMismatch{
				AccountID: accountID, AccountCode: journal.AccountCode, AccountType: journal.AccountType,
				JournalCents: journal.AmountCents, LedgerCents: gl.AmountCents,
				DifferenceCents: journal.AmountCents - gl.AmountCents,
			})
		}
	}
	// Account IDs are not ordered in a Go map; sort for stable output.
	sortCrossReportMismatches(result.ProfitLossAccountMismatches)
}

func sortCrossReportMismatches(mismatches []CrossReportAccountMismatch) {
	sort.Slice(mismatches, func(left, right int) bool {
		if mismatches[left].AccountCode != mismatches[right].AccountCode {
			return mismatches[left].AccountCode < mismatches[right].AccountCode
		}
		return mismatches[left].AccountID < mismatches[right].AccountID
	})
}

func readCrossReportControlTotals(ctx context.Context, tx *sql.Tx, tenantID, restaurantID int, to time.Time, result *CrossReportDiagnostics) error {
	var debitText, creditText string
	err := tx.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(gl.debit_amount), 0)::text,
		       COALESCE(SUM(gl.credit_amount), 0)::text
		FROM general_ledger gl
		JOIN chart_of_accounts coa ON coa.id = gl.account_id
		WHERE gl.tenant_id = $1 AND gl.restaurant_id = $2 AND gl.entry_date <= $3::date
		  AND coa.tenant_id = $1 AND coa.restaurant_id = $2
	`, tenantID, restaurantID, to).Scan(&debitText, &creditText)
	if err != nil {
		return fmt.Errorf("read exact GL trial-balance totals: %w", err)
	}
	result.TrialBalanceDebitsCents, err = numericTextToCents(debitText)
	if err != nil {
		return fmt.Errorf("parse GL debit total: %w", err)
	}
	result.TrialBalanceCreditsCents, err = numericTextToCents(creditText)
	if err != nil {
		return fmt.Errorf("parse GL credit total: %w", err)
	}

	var assets, liabilities, equity string
	err = tx.QueryRowContext(ctx, `
		SELECT
		  COALESCE(SUM(CASE WHEN coa.account_type = 'asset' THEN
		    CASE WHEN coa.normal_balance = 'debit' THEN gl.debit_amount - gl.credit_amount ELSE gl.credit_amount - gl.debit_amount END
		    ELSE 0 END), 0)::text,
		  COALESCE(SUM(CASE WHEN coa.account_type = 'liability' THEN
		    CASE WHEN coa.normal_balance = 'credit' THEN gl.credit_amount - gl.debit_amount ELSE gl.debit_amount - gl.credit_amount END
		    ELSE 0 END), 0)::text,
		  COALESCE(SUM(CASE WHEN coa.account_type = 'equity' THEN
		    CASE WHEN coa.normal_balance = 'credit' THEN gl.credit_amount - gl.debit_amount ELSE gl.debit_amount - gl.credit_amount END
		    ELSE 0 END), 0)::text
		FROM chart_of_accounts coa
		LEFT JOIN general_ledger gl ON gl.account_id = coa.id
		  AND gl.tenant_id = $1 AND gl.restaurant_id = $2 AND gl.entry_date <= $3::date
		WHERE coa.tenant_id = $1 AND coa.restaurant_id = $2
		  AND coa.is_active = TRUE AND coa.is_header = FALSE
		  AND coa.account_type IN ('asset', 'liability', 'equity')
	`, tenantID, restaurantID, to).Scan(&assets, &liabilities, &equity)
	if err != nil {
		return fmt.Errorf("read exact GL balance-sheet totals: %w", err)
	}
	if result.BalanceSheetAssetsCents, err = numericTextToCents(assets); err != nil {
		return fmt.Errorf("parse GL asset total: %w", err)
	}
	if result.BalanceSheetLiabilitiesCents, err = numericTextToCents(liabilities); err != nil {
		return fmt.Errorf("parse GL liability total: %w", err)
	}
	if result.BalanceSheetEquityCents, err = numericTextToCents(equity); err != nil {
		return fmt.Errorf("parse GL equity total: %w", err)
	}
	result.BalanceSheetResidualCents = result.BalanceSheetAssetsCents - result.BalanceSheetLiabilitiesCents - result.BalanceSheetEquityCents
	return nil
}

func readCrossReportPeriodContext(ctx context.Context, tx *sql.Tx, tenantID, restaurantID int, from, to time.Time, result *CrossReportDiagnostics) error {
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FILTER (WHERE start_date <= $3::date AND end_date >= $4::date),
		       COUNT(*) FILTER (WHERE start_date <= $4::date AND end_date >= $3::date)
		FROM fiscal_periods
		WHERE tenant_id = $1 AND restaurant_id = $2
	`, tenantID, restaurantID, from, to).Scan(&result.CoveredFiscalPeriods, &result.OverlappingFiscalPeriods)
	if err != nil {
		return fmt.Errorf("read scoped fiscal-period context: %w", err)
	}
	err = tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM journal_entries
		WHERE tenant_id = $1 AND restaurant_id = $2 AND status = 'void'
		  AND voided_at IS NOT NULL AND (
		    (entry_date >= $3::date AND entry_date <= $4::date)
		    OR (voided_at::date >= $3::date AND voided_at::date <= $4::date)
		  ) AND (entry_date <> voided_at::date OR fiscal_period IS DISTINCT FROM to_char(voided_at, 'YYYY-MM'))
	`, tenantID, restaurantID, from, to).Scan(&result.CrossPeriodVoidCount)
	if err != nil {
		return fmt.Errorf("read cross-period reversal context: %w", err)
	}
	if result.CoveredFiscalPeriods != 1 || result.OverlappingFiscalPeriods != 1 {
		result.PolicyNotes = append(result.PolicyNotes, "fiscal_period_boundary_or_calendar_policy_pending")
	}
	if result.CrossPeriodVoidCount > 0 {
		result.PolicyNotes = append(result.PolicyNotes, "cross_period_reversal_presentation_policy_pending")
	}
	return nil
}

func setCrossReportStatuses(result *CrossReportDiagnostics) {
	if len(result.ProfitLossAccountMismatches) > 0 {
		result.ProfitLossStatus = CrossReportDiscrepancies
	}
	if result.TrialBalanceDebitsCents != result.TrialBalanceCreditsCents {
		result.TrialBalanceStatus = CrossReportDiscrepancies
	}
	if result.ProfitLossStatus == CrossReportDiscrepancies || result.TrialBalanceStatus == CrossReportDiscrepancies {
		result.Status = CrossReportDiscrepancies
		return
	}
	if result.CashFlowStatus == CrossReportPolicyPending || result.BalanceSheetStatus == CrossReportPolicyPending || len(result.PolicyNotes) > 0 {
		result.Status = CrossReportPolicyPending
		return
	}
	result.Status = CrossReportReconciled
}
