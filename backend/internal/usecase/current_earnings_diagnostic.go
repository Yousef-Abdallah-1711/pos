package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	CurrentEarningsPolicyPending = "POLICY_PENDING"
	CurrentEarningsNoActivity    = "NO_ACTIVITY"
)

type CurrentEarningsDiagnostic struct {
	TenantID                      int
	RestaurantID                  int
	PeriodStart                   time.Time
	PeriodEnd                     time.Time
	Status                        string
	FiscalPeriodState             string
	RevenueCents                  int64
	ExpenseCents                  int64
	CandidateNetEarningsCents     int64
	RetainedEarningsMovementCents int64
	RetainedEarningsJournalCount  int64
	PotentialDoubleCountingRisk   bool
	CloseJournalIdentification    string
	CrossPeriodVoidCount          int64
	PolicyGaps                    []string
}

// DiagnosePotentialCurrentEarnings calculates a scoped candidate from posted
// GL movements without changing Balance Sheet presentation or posting a close.
// Retained-earnings movements are observations only because the schema has no
// approved journal marker that distinguishes a close from another adjustment.
func (uc *FinancialReportUseCase) DiagnosePotentialCurrentEarnings(
	ctx context.Context,
	tenantID, restaurantID int,
	periodStart, periodEnd time.Time,
) (CurrentEarningsDiagnostic, error) {
	if tenantID <= 0 || restaurantID <= 0 || periodStart.IsZero() || periodEnd.IsZero() || periodEnd.Before(periodStart) {
		return CurrentEarningsDiagnostic{}, fmt.Errorf("valid tenant, restaurant, and period boundaries are required")
	}
	periodStart, periodEnd = financialCalendarDate(periodStart), financialCalendarDate(periodEnd)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := uc.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return CurrentEarningsDiagnostic{}, fmt.Errorf("begin current-earnings snapshot: %w", err)
	}
	defer tx.Rollback()
	result := newCurrentEarningsDiagnostic(tenantID, restaurantID, periodStart, periodEnd)
	if err := readCurrentPeriodGL(ctx, tx, &result); err != nil {
		return CurrentEarningsDiagnostic{}, err
	}
	if err := readRetainedEarningsMovements(ctx, tx, &result); err != nil {
		return CurrentEarningsDiagnostic{}, err
	}
	if err := readCurrentFiscalPeriodState(ctx, tx, &result); err != nil {
		return CurrentEarningsDiagnostic{}, err
	}
	if err := readCurrentPeriodReversalContext(ctx, tx, &result); err != nil {
		return CurrentEarningsDiagnostic{}, err
	}
	if err := tx.Commit(); err != nil {
		return CurrentEarningsDiagnostic{}, fmt.Errorf("finish current-earnings snapshot: %w", err)
	}
	result.CandidateNetEarningsCents, err = subtractCurrentCents(result.RevenueCents, result.ExpenseCents)
	if err != nil {
		return CurrentEarningsDiagnostic{}, fmt.Errorf("calculate candidate net earnings: %w", err)
	}
	result.PotentialDoubleCountingRisk = result.CandidateNetEarningsCents != 0 && result.RetainedEarningsMovementCents != 0
	return result, nil
}

func newCurrentEarningsDiagnostic(tenantID, restaurantID int, start, end time.Time) CurrentEarningsDiagnostic {
	return CurrentEarningsDiagnostic{
		TenantID: tenantID, RestaurantID: restaurantID, PeriodStart: start, PeriodEnd: end,
		Status: CurrentEarningsPolicyPending, FiscalPeriodState: "POLICY_PENDING",
		CloseJournalIdentification: "UNAVAILABLE",
		PolicyGaps: []string{
			"EARN-01 earnings presentation approval pending",
			"fiscal year and period selection require accounting policy",
			"retained-earnings movements cannot be identified as close journals from current source fields",
		},
	}
}

func readCurrentPeriodGL(ctx context.Context, tx *sql.Tx, result *CurrentEarningsDiagnostic) error {
	var revenueText, expenseText string
	var movementCount int64
	err := tx.QueryRowContext(ctx, `
		SELECT
		  COUNT(*),
		  COALESCE(SUM(CASE WHEN coa.account_type = 'revenue' THEN gl.credit_amount - gl.debit_amount ELSE 0 END), 0)::text,
		  COALESCE(SUM(CASE WHEN coa.account_type = 'expense' THEN gl.debit_amount - gl.credit_amount ELSE 0 END), 0)::text
		FROM general_ledger gl
		JOIN chart_of_accounts coa ON coa.id = gl.account_id
		WHERE gl.tenant_id = $1 AND gl.restaurant_id = $2
		  AND gl.entry_date >= $3::date AND gl.entry_date <= $4::date
		  AND coa.tenant_id = $1 AND coa.restaurant_id = $2
		  AND coa.account_type IN ('revenue', 'expense')
	`, result.TenantID, result.RestaurantID, result.PeriodStart, result.PeriodEnd).Scan(&movementCount, &revenueText, &expenseText)
	if err != nil {
		return fmt.Errorf("read scoped posted GL earnings movements: %w", err)
	}
	if result.RevenueCents, err = currentNumericTextToCents(revenueText); err != nil {
		return fmt.Errorf("parse scoped GL revenue: %w", err)
	}
	if result.ExpenseCents, err = currentNumericTextToCents(expenseText); err != nil {
		return fmt.Errorf("parse scoped GL expenses: %w", err)
	}
	if movementCount == 0 {
		result.Status = CurrentEarningsNoActivity
	}
	return nil
}

func readRetainedEarningsMovements(ctx context.Context, tx *sql.Tx, result *CurrentEarningsDiagnostic) error {
	var amountText string
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT gl.journal_entry_id),
		       COALESCE(SUM(gl.credit_amount - gl.debit_amount), 0)::text
		FROM general_ledger gl
		JOIN chart_of_accounts coa ON coa.id = gl.account_id
		WHERE gl.tenant_id = $1 AND gl.restaurant_id = $2
		  AND gl.entry_date >= $3::date AND gl.entry_date <= $4::date
		  AND coa.tenant_id = $1 AND coa.restaurant_id = $2
		  AND coa.account_type = 'equity' AND coa.account_subtype = 'retained_earnings'
	`, result.TenantID, result.RestaurantID, result.PeriodStart, result.PeriodEnd).Scan(&result.RetainedEarningsJournalCount, &amountText)
	if err != nil {
		return fmt.Errorf("read scoped retained-earnings GL movements: %w", err)
	}
	result.RetainedEarningsMovementCents, err = currentNumericTextToCents(amountText)
	if err != nil {
		return fmt.Errorf("parse retained-earnings movements: %w", err)
	}
	if result.Status == CurrentEarningsNoActivity && result.RetainedEarningsJournalCount > 0 {
		result.Status = CurrentEarningsPolicyPending
	}
	return nil
}

func readCurrentFiscalPeriodState(ctx context.Context, tx *sql.Tx, result *CurrentEarningsDiagnostic) error {
	var covered, overlapping, closed int64
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FILTER (WHERE start_date <= $3::date AND end_date >= $4::date),
		       COUNT(*) FILTER (WHERE start_date <= $4::date AND end_date >= $3::date),
		       COUNT(*) FILTER (WHERE start_date <= $3::date AND end_date >= $4::date AND status = 'closed')
		FROM fiscal_periods
		WHERE tenant_id = $1 AND restaurant_id = $2
	`, result.TenantID, result.RestaurantID, result.PeriodStart, result.PeriodEnd).Scan(&covered, &overlapping, &closed)
	if err != nil {
		return fmt.Errorf("read scoped current fiscal-period status: %w", err)
	}
	switch {
	case covered == 1 && closed == 1:
		result.FiscalPeriodState = "CLOSED_PERIOD_OBSERVED"
	case covered == 1 && overlapping == 1:
		result.FiscalPeriodState = "OPEN_PERIOD_OBSERVED"
	case covered == 0 && overlapping == 0:
		result.FiscalPeriodState = "SOURCE_UNAVAILABLE"
	default:
		result.FiscalPeriodState = "AMBIGUOUS_PERIOD_SCOPE"
	}
	return nil
}

func readCurrentPeriodReversalContext(ctx context.Context, tx *sql.Tx, result *CurrentEarningsDiagnostic) error {
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM journal_entries
		WHERE tenant_id = $1 AND restaurant_id = $2 AND status = 'void' AND voided_at IS NOT NULL
		  AND ((entry_date >= $3::date AND entry_date <= $4::date)
		    OR (voided_at::date >= $3::date AND voided_at::date <= $4::date))
		  AND (entry_date <> voided_at::date OR fiscal_period IS DISTINCT FROM to_char(voided_at, 'YYYY-MM'))
	`, result.TenantID, result.RestaurantID, result.PeriodStart, result.PeriodEnd).Scan(&result.CrossPeriodVoidCount)
	if err != nil {
		return fmt.Errorf("read current-period reversal context: %w", err)
	}
	return nil
}

func financialCalendarDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

func subtractCurrentCents(left, right int64) (int64, error) {
	result := new(big.Int).Sub(big.NewInt(left), big.NewInt(right))
	if !result.IsInt64() {
		return 0, fmt.Errorf("candidate net earnings exceed the signed-cent range")
	}
	return result.Int64(), nil
}

func currentNumericTextToCents(value string) (int64, error) {
	value = strings.TrimSpace(value)
	negative := strings.HasPrefix(value, "-")
	original := value
	value = strings.TrimPrefix(strings.TrimPrefix(value, "-"), "+")
	parts := strings.Split(value, ".")
	if len(parts) > 2 || parts[0] == "" || (len(parts) == 2 && len(parts[1]) > 2) {
		return 0, fmt.Errorf("unsupported NUMERIC value %q", original)
	}
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse whole units: %w", err)
	}
	fraction := "00"
	if len(parts) == 2 {
		fraction = parts[1] + strings.Repeat("0", 2-len(parts[1]))
	}
	cents, err := strconv.ParseInt(fraction, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse fractional units: %w", err)
	}
	amount := new(big.Int).Mul(big.NewInt(whole), big.NewInt(100))
	amount.Add(amount, big.NewInt(cents))
	if negative {
		amount.Neg(amount)
	}
	if !amount.IsInt64() {
		return 0, fmt.Errorf("NUMERIC value %q exceeds the diagnostic cents range", original)
	}
	return amount.Int64(), nil
}
