package repository

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
	OpeningBalanceSourceUnavailable = "SOURCE_UNAVAILABLE"
	OpeningBalancePolicyPending     = "POLICY_PENDING"
)

// OpeningBalanceDiagnostic reports fields available for opening-balance review.
// bank_accounts_ledger has no source-reference or approval fields, so this
// diagnostic never describes a candidate amount as an accounting mismatch.
type OpeningBalanceDiagnostic struct {
	TenantID      int
	RestaurantID  int
	Status        string
	SourceModel   string
	Authoritative bool
	PolicyGaps    []string
	Accounts      []OpeningBalanceObservation
}

// OpeningBalanceObservation contains stored bank metadata and scoped GL context.
type OpeningBalanceObservation struct {
	BankAccountID                 int
	AccountName                   string
	Currency                      string
	OpeningBalanceCents           int64
	OpeningDate                   time.Time
	GLAccountID                   *int
	GLAccountScope                string
	GLAccountType                 string
	GLNormalBalance               string
	GLMovementOnOrBeforeDateCents int64
	GLRowsOnOpeningDate           int64
	HistoricalMovementBeforeDate  bool
	SourceReferenceAvailable      bool
	ApprovalMetadataAvailable     bool
}

type openingBalanceRow interface {
	Scan(dest ...any) error
}

// DiagnoseOpeningBalances reads bank account opening metadata and the scoped
// ledger snapshot using one read-only repeatable-read transaction. Since the
// schema has no approved source reference, approval state, or mapping contract,
// observations remain policy-pending and do not assert discrepancies.
func (r *BankAccountRepository) DiagnoseOpeningBalances(
	ctx context.Context,
	tenantID, restaurantID int,
) (OpeningBalanceDiagnostic, error) {
	if tenantID <= 0 || restaurantID <= 0 {
		return OpeningBalanceDiagnostic{}, fmt.Errorf("tenant and restaurant scope are required")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return OpeningBalanceDiagnostic{}, fmt.Errorf("begin opening-balance diagnostic snapshot: %w", err)
	}
	defer tx.Rollback()

	result := emptyOpeningBalanceDiagnostic(tenantID, restaurantID)
	rows, err := tx.QueryContext(ctx, openingBalanceObservationsQuery, tenantID, restaurantID)
	if err != nil {
		return OpeningBalanceDiagnostic{}, fmt.Errorf("read scoped opening-balance observations: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		observation, scanErr := scanOpeningBalanceObservation(rows)
		if scanErr != nil {
			return OpeningBalanceDiagnostic{}, fmt.Errorf("scan opening-balance observation: %w", scanErr)
		}
		result.Accounts = append(result.Accounts, observation)
	}
	if err := rows.Err(); err != nil {
		return OpeningBalanceDiagnostic{}, fmt.Errorf("iterate opening-balance observations: %w", err)
	}
	if len(result.Accounts) == 0 {
		result.Status = OpeningBalanceSourceUnavailable
	}
	if err := tx.Commit(); err != nil {
		return OpeningBalanceDiagnostic{}, fmt.Errorf("finish opening-balance diagnostic snapshot: %w", err)
	}
	return result, nil
}

const openingBalanceObservationsQuery = `
	SELECT ba.id, ba.account_name, ba.currency, ba.opening_balance::text,
	       ba.opening_date, ba.gl_account_id,
	       CASE
	         WHEN ba.gl_account_id IS NULL THEN 'missing'
	         WHEN coa.id IS NULL THEN 'missing'
	         WHEN coa.tenant_id <> $1 OR coa.restaurant_id <> $2 THEN 'out_of_scope'
	         ELSE 'scoped'
	       END AS gl_account_scope,
	       CASE WHEN coa.tenant_id = $1 AND coa.restaurant_id = $2
	            THEN coa.account_type ELSE '' END,
	       CASE WHEN coa.tenant_id = $1 AND coa.restaurant_id = $2
	            THEN coa.normal_balance ELSE '' END,
	       COALESCE(SUM(gl.debit_amount - gl.credit_amount), 0)::text,
	       COUNT(gl.id) FILTER (WHERE gl.entry_date = ba.opening_date),
	       COUNT(gl.id) FILTER (WHERE gl.entry_date < ba.opening_date) > 0
	FROM bank_accounts_ledger ba
	LEFT JOIN chart_of_accounts coa ON coa.id = ba.gl_account_id
	LEFT JOIN general_ledger gl
	  ON gl.account_id = ba.gl_account_id
	 AND gl.tenant_id = $1 AND gl.restaurant_id = $2
	 AND gl.entry_date <= ba.opening_date
	WHERE ba.tenant_id = $1 AND ba.restaurant_id = $2
	GROUP BY ba.id, ba.account_name, ba.currency, ba.opening_balance,
	         ba.opening_date, ba.gl_account_id, coa.id, coa.tenant_id,
	         coa.restaurant_id, coa.account_type, coa.normal_balance
	ORDER BY ba.id ASC`

func emptyOpeningBalanceDiagnostic(tenantID, restaurantID int) OpeningBalanceDiagnostic {
	return OpeningBalanceDiagnostic{
		TenantID: tenantID, RestaurantID: restaurantID,
		Status: OpeningBalancePolicyPending, SourceModel: "bank_accounts_ledger",
		Authoritative: false,
		PolicyGaps: []string{
			"opening_source_reference_unavailable",
			"opening_approval_metadata_unavailable",
			"opening_cutover_and_account_mapping_policy_pending",
		},
		Accounts: make([]OpeningBalanceObservation, 0),
	}
}

func scanOpeningBalanceObservation(row openingBalanceRow) (OpeningBalanceObservation, error) {
	var observation OpeningBalanceObservation
	var openingBalance, ledgerMovement string
	var glAccountID sql.NullInt64
	if err := row.Scan(
		&observation.BankAccountID, &observation.AccountName, &observation.Currency,
		&openingBalance, &observation.OpeningDate, &glAccountID,
		&observation.GLAccountScope, &observation.GLAccountType,
		&observation.GLNormalBalance, &ledgerMovement,
		&observation.GLRowsOnOpeningDate, &observation.HistoricalMovementBeforeDate,
	); err != nil {
		return OpeningBalanceObservation{}, err
	}
	var err error
	if observation.OpeningBalanceCents, err = numericTextToCents(openingBalance); err != nil {
		return OpeningBalanceObservation{}, fmt.Errorf("parse opening balance: %w", err)
	}
	if observation.GLMovementOnOrBeforeDateCents, err = numericTextToCents(ledgerMovement); err != nil {
		return OpeningBalanceObservation{}, fmt.Errorf("parse GL movement: %w", err)
	}
	if glAccountID.Valid {
		id := int(glAccountID.Int64)
		observation.GLAccountID = &id
	}
	return observation, nil
}

func numericTextToCents(value string) (int64, error) {
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
