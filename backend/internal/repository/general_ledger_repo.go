package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"pos-saas/internal/domain"
)

type GeneralLedgerRepository struct {
	db          *sql.DB
	journalRepo *JournalEntryRepository
}

type LedgerPeriodSummary struct {
	FiscalPeriod   string  `json:"fiscal_period"`
	EntryCount     int     `json:"entry_count"`
	DebitTotal     float64 `json:"debit_total"`
	CreditTotal    float64 `json:"credit_total"`
	NetChange      float64 `json:"net_change"`
	ClosingBalance float64 `json:"closing_balance"`
}

func NewGeneralLedgerRepository(db *sql.DB, journalRepo *JournalEntryRepository) *GeneralLedgerRepository {
	return &GeneralLedgerRepository{
		db:          db,
		journalRepo: journalRepo,
	}
}

// CreateAndPostAutoEntry creates an auto-posting journal and its ledger rows in
// one transaction. A transaction-scoped advisory lock serializes retries for
// the same source event so only one posting can be created for that identity.
func (r *GeneralLedgerRepository) CreateAndPostAutoEntry(entry *domain.JournalEntry, postedBy int) (int, error) {
	if r.db == nil || r.journalRepo == nil {
		return 0, errors.New("general ledger and journal repositories are required")
	}
	if entry == nil {
		return 0, errors.New("journal entry is required")
	}
	if entry.TenantID <= 0 || entry.RestaurantID <= 0 {
		return 0, errors.New("tenant and restaurant are required")
	}
	if entry.ReferenceID == nil || *entry.ReferenceID <= 0 || strings.TrimSpace(entry.ReferenceType) == "" || strings.TrimSpace(entry.EntryType) == "" {
		return 0, errors.New("auto-posting requires a stable source reference and entry type")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	lockKey := fmt.Sprintf("auto-posting:%d:%d:%s:%d:%s", entry.TenantID, entry.RestaurantID, entry.ReferenceType, *entry.ReferenceID, entry.EntryType)
	var lockAcquired int
	if err := tx.QueryRow(`SELECT 1 FROM (SELECT pg_advisory_xact_lock(hashtextextended($1, 0))) AS source_lock`, lockKey).Scan(&lockAcquired); err != nil {
		return 0, fmt.Errorf("failed to lock auto-posting source event: %w", err)
	}

	rows, err := tx.Query(`
		SELECT id, status
		FROM journal_entries
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND reference_type = $3
		  AND reference_id = $4
		  AND entry_type = $5
		ORDER BY id
		FOR UPDATE
	`, entry.TenantID, entry.RestaurantID, entry.ReferenceType, *entry.ReferenceID, entry.EntryType)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect existing auto-posting entries: %w", err)
	}
	type existingPosting struct {
		id     int
		status string
	}
	existing := make([]existingPosting, 0, 1)
	for rows.Next() {
		var posting existingPosting
		if err := rows.Scan(&posting.id, &posting.status); err != nil {
			rows.Close()
			return 0, err
		}
		existing = append(existing, posting)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return 0, err
	}
	if err := rows.Close(); err != nil {
		return 0, err
	}
	if len(existing) > 1 {
		return 0, fmt.Errorf("multiple auto-posting entries already exist for source %s:%d", entry.ReferenceType, *entry.ReferenceID)
	}
	if len(existing) == 1 {
		switch existing[0].status {
		case domain.JournalStatusPosted:
			if err := tx.Commit(); err != nil {
				return 0, err
			}
			return existing[0].id, nil
		case domain.JournalStatusDraft:
			return 0, fmt.Errorf("auto-posting draft %d already exists and requires recovery", existing[0].id)
		case domain.JournalStatusVoid:
			return 0, fmt.Errorf("auto-posting entry %d is void; a new source event is required", existing[0].id)
		default:
			return 0, fmt.Errorf("auto-posting entry %d has unsupported status %q", existing[0].id, existing[0].status)
		}
	}

	id, err := r.journalRepo.CreateJournalEntryTx(tx, entry)
	if err != nil {
		return 0, fmt.Errorf("failed to create auto-posting journal: %w", err)
	}
	if err := r.journalRepo.PostEntryTx(tx, entry.TenantID, entry.RestaurantID, id, postedBy); err != nil {
		return 0, fmt.Errorf("failed to post auto-posting journal %d: %w", id, err)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit auto-posting journal %d: %w", id, err)
	}
	return id, nil
}

func (r *GeneralLedgerRepository) AppendEntry(tenantID, restaurantID, journalEntryID, postedBy int) error {
	if r.journalRepo == nil {
		return errors.New("journal repository not configured")
	}
	return r.journalRepo.PostEntry(tenantID, restaurantID, journalEntryID, postedBy)
}

func (r *GeneralLedgerRepository) GetAccountBalance(tenantID, restaurantID, accountID int, asOf *time.Time) (float64, error) {
	query := `
		SELECT COALESCE(running_balance, 0)
		FROM general_ledger
		WHERE tenant_id = $1 AND restaurant_id = $2 AND account_id = $3
	`
	args := []interface{}{tenantID, restaurantID, accountID}
	if asOf != nil && !asOf.IsZero() {
		query += " AND entry_date <= $4"
		args = append(args, *asOf)
	}
	query += " ORDER BY entry_date DESC, id DESC LIMIT 1"

	var balance float64
	err := r.db.QueryRow(query, args...).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return balance, nil
}

func (r *GeneralLedgerRepository) GetLedgerByAccount(tenantID, restaurantID, accountID int, startDate, endDate *time.Time) ([]domain.GeneralLedgerLine, error) {
	lines, _, err := r.GetLedgerByAccountCursor(tenantID, restaurantID, accountID, startDate, endDate, 0, nil)
	return lines, err
}

func (r *GeneralLedgerRepository) GetLedgerByAccountCursor(
	tenantID, restaurantID, accountID int,
	startDate, endDate *time.Time,
	limit int,
	cursor *int,
) ([]domain.GeneralLedgerLine, *int, error) {
	var (
		builder strings.Builder
		args    = []interface{}{tenantID, restaurantID, accountID}
		argPos  = 4
	)
	if limit <= 0 {
		limit = 200
	}
	if limit > 500 {
		limit = 500
	}

	builder.WriteString(`
		SELECT
			gl.id,
			gl.account_id,
			coa.account_code,
			coa.account_name,
			gl.journal_entry_id,
			gl.journal_line_id,
			gl.entry_date,
			je.entry_number,
			COALESCE(je.description, ''),
			COALESCE(je.reference_type, ''),
			COALESCE(je.reference_number, ''),
			gl.debit_amount,
			gl.credit_amount,
			gl.running_balance,
			COALESCE(jel.description, '')
		FROM general_ledger gl
		JOIN chart_of_accounts coa ON coa.id = gl.account_id
		JOIN journal_entries je ON je.id = gl.journal_entry_id
		LEFT JOIN journal_entry_lines jel ON jel.id = gl.journal_line_id
		WHERE gl.tenant_id = $1 AND gl.restaurant_id = $2 AND gl.account_id = $3
	`)

	if startDate != nil && !startDate.IsZero() {
		builder.WriteString(fmt.Sprintf(" AND gl.entry_date >= $%d", argPos))
		args = append(args, *startDate)
		argPos++
	}
	if endDate != nil && !endDate.IsZero() {
		builder.WriteString(fmt.Sprintf(" AND gl.entry_date <= $%d", argPos))
		args = append(args, *endDate)
		argPos++
	}

	if cursor != nil && *cursor > 0 {
		builder.WriteString(fmt.Sprintf(" AND gl.id < $%d", argPos))
		args = append(args, *cursor)
		argPos++
	}
	builder.WriteString(fmt.Sprintf(" ORDER BY gl.entry_date DESC, gl.id DESC LIMIT $%d", argPos))
	args = append(args, limit+1)

	rows, err := r.db.Query(builder.String(), args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	lines := make([]domain.GeneralLedgerLine, 0)
	lineIDs := make([]int, 0)
	for rows.Next() {
		var ledgerID int
		var line domain.GeneralLedgerLine
		if err := rows.Scan(
			&ledgerID,
			&line.AccountID,
			&line.AccountCode,
			&line.AccountName,
			&line.JournalEntryID,
			&line.JournalLineID,
			&line.EntryDate,
			&line.EntryNumber,
			&line.Description,
			&line.ReferenceType,
			&line.ReferenceNumber,
			&line.DebitAmount,
			&line.CreditAmount,
			&line.RunningBalance,
			&line.LineDescription,
		); err != nil {
			return nil, nil, err
		}
		lines = append(lines, line)
		lineIDs = append(lineIDs, ledgerID)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	var nextCursor *int
	if len(lines) > limit {
		lines = lines[:limit]
		value := lineIDs[limit-1]
		nextCursor = &value
	}

	return lines, nextCursor, nil
}

func (r *GeneralLedgerRepository) GetLedgerPeriodSummaries(tenantID, restaurantID, accountID int, startDate, endDate *time.Time) ([]LedgerPeriodSummary, error) {
	var (
		builder strings.Builder
		args    = []interface{}{tenantID, restaurantID, accountID}
		argPos  = 4
	)

	builder.WriteString(`
		WITH filtered AS (
			SELECT id, fiscal_period, entry_date, debit_amount, credit_amount, running_balance
			FROM general_ledger
			WHERE tenant_id = $1 AND restaurant_id = $2 AND account_id = $3
	`)
	if startDate != nil && !startDate.IsZero() {
		builder.WriteString(fmt.Sprintf(" AND entry_date >= $%d", argPos))
		args = append(args, *startDate)
		argPos++
	}
	if endDate != nil && !endDate.IsZero() {
		builder.WriteString(fmt.Sprintf(" AND entry_date <= $%d", argPos))
		args = append(args, *endDate)
		argPos++
	}
	builder.WriteString(`
		),
		period_totals AS (
			SELECT
				fiscal_period,
				COUNT(*)::INT AS entry_count,
				COALESCE(SUM(debit_amount), 0) AS debit_total,
				COALESCE(SUM(credit_amount), 0) AS credit_total
			FROM filtered
			GROUP BY fiscal_period
		),
		period_closing AS (
			SELECT DISTINCT ON (fiscal_period)
				fiscal_period,
				running_balance
			FROM filtered
			ORDER BY fiscal_period, entry_date DESC, id DESC
		)
		SELECT
			pt.fiscal_period,
			pt.entry_count,
			pt.debit_total,
			pt.credit_total,
			pt.debit_total - pt.credit_total AS net_change,
			pc.running_balance
		FROM period_totals pt
		JOIN period_closing pc ON pc.fiscal_period = pt.fiscal_period
		ORDER BY pt.fiscal_period DESC
	`)

	rows, err := r.db.Query(builder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]LedgerPeriodSummary, 0)
	for rows.Next() {
		var row LedgerPeriodSummary
		if err := rows.Scan(
			&row.FiscalPeriod,
			&row.EntryCount,
			&row.DebitTotal,
			&row.CreditTotal,
			&row.NetChange,
			&row.ClosingBalance,
		); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GeneralLedgerRepository) GetTrialBalance(tenantID, restaurantID int, asOf time.Time) ([]domain.TrialBalanceRow, error) {
	if asOf.IsZero() {
		asOf = time.Now()
	}

	rows, err := r.db.Query(`
		WITH latest_balances AS (
			SELECT DISTINCT ON (gl.account_id)
				gl.account_id,
				gl.running_balance
			FROM general_ledger gl
			WHERE gl.tenant_id = $1
			  AND gl.restaurant_id = $2
			  AND gl.entry_date <= $3
			ORDER BY gl.account_id, gl.entry_date DESC, gl.id DESC
		)
		SELECT
			coa.id,
			coa.account_code,
			coa.account_name,
			coa.account_type,
			coa.normal_balance,
			COALESCE(lb.running_balance, 0) AS net_balance
		FROM chart_of_accounts coa
		LEFT JOIN latest_balances lb ON lb.account_id = coa.id
		WHERE coa.tenant_id = $1
		  AND coa.restaurant_id = $2
		  AND coa.is_active = TRUE
		  AND coa.is_header = FALSE
		ORDER BY coa.account_code ASC
	`, tenantID, restaurantID, asOf)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trialBalance := make([]domain.TrialBalanceRow, 0)
	for rows.Next() {
		var (
			row        domain.TrialBalanceRow
			netBalance float64
		)
		if err := rows.Scan(
			&row.AccountID,
			&row.AccountCode,
			&row.AccountName,
			&row.AccountType,
			&row.NormalBalance,
			&netBalance,
		); err != nil {
			return nil, err
		}

		row.AsOfDate = asOf
		row.NetBalance = netBalance
		if row.NormalBalance == domain.NormalBalanceCredit {
			if netBalance >= 0 {
				row.CreditBalance = netBalance
			} else {
				row.DebitBalance = -netBalance
			}
		} else {
			if netBalance >= 0 {
				row.DebitBalance = netBalance
			} else {
				row.CreditBalance = -netBalance
			}
		}

		trialBalance = append(trialBalance, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trialBalance, nil
}

func nullString(value string) interface{} {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
