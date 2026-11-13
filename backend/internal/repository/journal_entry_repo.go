package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"pos-saas/internal/domain"
)

type JournalEntryRepository struct {
	db *sql.DB
}

func NewJournalEntryRepository(db *sql.DB) *JournalEntryRepository {
	return &JournalEntryRepository{db: db}
}

func (r *JournalEntryRepository) ListJournalEntries(tenantID, restaurantID int, filters domain.JournalEntryFilters) ([]domain.JournalEntry, error) {
	entries, _, err := r.ListJournalEntriesPaginated(tenantID, restaurantID, filters, 1000, 0)
	return entries, err
}

func (r *JournalEntryRepository) ListJournalEntriesPaginated(tenantID, restaurantID int, filters domain.JournalEntryFilters, limit, offset int) ([]domain.JournalEntry, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	var (
		builder strings.Builder
		args    = []interface{}{tenantID, restaurantID}
		idx     = 3
	)

	builder.WriteString(`
		SELECT
			id, tenant_id, restaurant_id, entry_number, entry_date, COALESCE(fiscal_period, ''),
			description, COALESCE(description_ar, ''), COALESCE(reference_type, ''),
			reference_id, COALESCE(reference_number, ''), entry_type, total_debit,
			total_credit, status, posted_at, posted_by, voided_at, voided_by,
			COALESCE(void_reason, ''), is_approved, approved_by, approved_at,
			COALESCE(notes, ''), COALESCE(attachments::text, ''), created_at, updated_at,
			created_by, updated_by
		FROM journal_entries
		WHERE tenant_id = $1 AND restaurant_id = $2
	`)

	if filters.StartDate != nil {
		builder.WriteString(fmt.Sprintf(" AND entry_date >= $%d", idx))
		args = append(args, *filters.StartDate)
		idx++
	}
	if filters.EndDate != nil {
		builder.WriteString(fmt.Sprintf(" AND entry_date <= $%d", idx))
		args = append(args, *filters.EndDate)
		idx++
	}
	if filters.Status != "" {
		builder.WriteString(fmt.Sprintf(" AND status = $%d", idx))
		args = append(args, filters.Status)
		idx++
	}
	if filters.EntryType != "" {
		builder.WriteString(fmt.Sprintf(" AND entry_type = $%d", idx))
		args = append(args, filters.EntryType)
		idx++
	}

	countQuery := "SELECT COUNT(*) FROM journal_entries WHERE tenant_id = $1 AND restaurant_id = $2"
	countArgs := []interface{}{tenantID, restaurantID}
	countIdx := 3

	if filters.StartDate != nil {
		countQuery += fmt.Sprintf(" AND entry_date >= $%d", countIdx)
		countArgs = append(countArgs, *filters.StartDate)
		countIdx++
	}
	if filters.EndDate != nil {
		countQuery += fmt.Sprintf(" AND entry_date <= $%d", countIdx)
		countArgs = append(countArgs, *filters.EndDate)
		countIdx++
	}
	if filters.Status != "" {
		countQuery += fmt.Sprintf(" AND status = $%d", countIdx)
		countArgs = append(countArgs, filters.Status)
		countIdx++
	}
	if filters.EntryType != "" {
		countQuery += fmt.Sprintf(" AND entry_type = $%d", countIdx)
		countArgs = append(countArgs, filters.EntryType)
		countIdx++
	}

	var total int
	if err := r.db.QueryRow(countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	builder.WriteString(fmt.Sprintf(" ORDER BY entry_date DESC, id DESC LIMIT $%d OFFSET $%d", idx, idx+1))
	args = append(args, limit, offset)

	rows, err := r.db.Query(builder.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	entries := make([]domain.JournalEntry, 0)
	for rows.Next() {
		entry, err := scanJournalEntry(rows)
		if err != nil {
			return nil, 0, err
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (r *JournalEntryRepository) GetJournalEntryByID(tenantID, restaurantID, id int) (*domain.JournalEntry, error) {
	entry, err := scanJournalEntry(r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, entry_number, entry_date, COALESCE(fiscal_period, ''),
			description, COALESCE(description_ar, ''), COALESCE(reference_type, ''),
			reference_id, COALESCE(reference_number, ''), entry_type, total_debit,
			total_credit, status, posted_at, posted_by, voided_at, voided_by,
			COALESCE(void_reason, ''), is_approved, approved_by, approved_at,
			COALESCE(notes, ''), COALESCE(attachments::text, ''), created_at, updated_at,
			created_by, updated_by
		FROM journal_entries
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	lines, err := r.getLines(id)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines

	return &entry, nil
}

func (r *JournalEntryRepository) CreateJournalEntry(entry *domain.JournalEntry) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	id, err := r.CreateJournalEntryTx(tx, entry)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *JournalEntryRepository) CreateJournalEntryTx(tx *sql.Tx, entry *domain.JournalEntry) (int, error) {
	if tx == nil {
		return 0, errors.New("transaction is required")
	}
	if entry.EntryDate.IsZero() {
		entry.EntryDate = time.Now()
	}
	if entry.EntryType == "" {
		entry.EntryType = "manual"
	}
	if entry.Status == "" {
		entry.Status = domain.JournalStatusDraft
	}
	entry.FiscalPeriod = formatFiscalPeriod(entry.EntryDate)

	totalDebit, totalCredit, err := validateJournalLines(entry.Lines)
	if err != nil {
		return 0, err
	}

	entryNumber, err := r.generateEntryNumberTx(tx, entry.TenantID, entry.RestaurantID, entry.EntryDate)
	if err != nil {
		return 0, err
	}
	var attachments interface{}
	if len(entry.Attachments) > 0 {
		attachments = entry.Attachments
	}

	var id int
	err = tx.QueryRow(`
		INSERT INTO journal_entries (
			tenant_id, restaurant_id, entry_number, entry_date, fiscal_period,
			description, description_ar, reference_type, reference_id, reference_number,
			entry_type, total_debit, total_credit, status, notes, attachments,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16,
			$17, $18
		)
		RETURNING id
	`,
		entry.TenantID,
		entry.RestaurantID,
		entryNumber,
		entry.EntryDate,
		entry.FiscalPeriod,
		entry.Description,
		entry.DescriptionAr,
		entry.ReferenceType,
		entry.ReferenceID,
		entry.ReferenceNumber,
		entry.EntryType,
		totalDebit,
		totalCredit,
		entry.Status,
		entry.Notes,
		attachments,
		entry.CreatedBy,
		entry.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	if err := r.insertLinesTx(tx, id, entry.TenantID, entry.RestaurantID, entry.Lines); err != nil {
		return 0, err
	}

	entry.ID = id
	entry.EntryNumber = entryNumber
	entry.TotalDebit = totalDebit
	entry.TotalCredit = totalCredit

	return id, nil
}

func (r *JournalEntryRepository) UpdateDraftEntry(entry *domain.JournalEntry) error {
	totalDebit, totalCredit, err := validateJournalLines(entry.Lines)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	err = tx.QueryRow(
		`SELECT status FROM journal_entries WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`,
		entry.ID,
		entry.TenantID,
		entry.RestaurantID,
	).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("journal entry not found")
		}
		return err
	}
	if status != domain.JournalStatusDraft {
		return errors.New("only draft journal entries can be updated")
	}

	if entry.EntryDate.IsZero() {
		entry.EntryDate = time.Now()
	}
	entry.FiscalPeriod = formatFiscalPeriod(entry.EntryDate)

	_, err = tx.Exec(`
		UPDATE journal_entries
		SET entry_date = $1,
			fiscal_period = $2,
			description = $3,
			description_ar = $4,
			reference_type = $5,
			reference_id = $6,
			reference_number = $7,
			entry_type = $8,
			total_debit = $9,
			total_credit = $10,
			notes = $11,
			attachments = $12,
			updated_by = $13,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $14 AND tenant_id = $15 AND restaurant_id = $16
	`,
		entry.EntryDate,
		entry.FiscalPeriod,
		entry.Description,
		entry.DescriptionAr,
		entry.ReferenceType,
		entry.ReferenceID,
		entry.ReferenceNumber,
		entry.EntryType,
		totalDebit,
		totalCredit,
		entry.Notes,
		entry.Attachments,
		entry.UpdatedBy,
		entry.ID,
		entry.TenantID,
		entry.RestaurantID,
	)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM journal_entry_lines WHERE journal_entry_id = $1`, entry.ID); err != nil {
		return err
	}
	if err := r.insertLinesTx(tx, entry.ID, entry.TenantID, entry.RestaurantID, entry.Lines); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *JournalEntryRepository) PostEntry(tenantID, restaurantID, id, postedBy int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := r.PostEntryTx(tx, tenantID, restaurantID, id, postedBy); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *JournalEntryRepository) PostEntryTx(tx *sql.Tx, tenantID, restaurantID, id, postedBy int) error {
	if tx == nil {
		return errors.New("transaction is required")
	}

	var (
		status       string
		entryDate    time.Time
		fiscalPeriod string
	)
	err := tx.QueryRow(
		`SELECT status, entry_date, COALESCE(fiscal_period, '') FROM journal_entries WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3 FOR UPDATE`,
		id,
		tenantID,
		restaurantID,
	).Scan(&status, &entryDate, &fiscalPeriod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("journal entry not found")
		}
		return err
	}
	if status != domain.JournalStatusDraft {
		return errors.New("only draft journal entries can be posted")
	}
	if fiscalPeriod == "" {
		fiscalPeriod = formatFiscalPeriod(entryDate)
	}
	if err := r.ensureFiscalPeriodOpenTx(tx, tenantID, restaurantID, entryDate, fiscalPeriod); err != nil {
		return err
	}

	lines, err := r.getPostingLinesTx(tx, id)
	if err != nil {
		return err
	}
	totalDebit, totalCredit, err := validatePostingLines(lines)
	if err != nil {
		return err
	}
	if math.Abs(totalDebit-totalCredit) > 0.009 {
		return errors.New("journal entry is not balanced")
	}

	_, err = tx.Exec(`
		UPDATE journal_entries
		SET status = $1,
			posted_at = CURRENT_TIMESTAMP,
			posted_by = $2,
			updated_by = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5
	`,
		domain.JournalStatusPosted,
		postedBy,
		id,
		tenantID,
		restaurantID,
	)
	if err != nil {
		return err
	}
	if err := r.appendGeneralLedgerTx(tx, tenantID, restaurantID, id, entryDate, fiscalPeriod, lines); err != nil {
		return err
	}

	return nil
}

func (r *JournalEntryRepository) VoidEntry(tenantID, restaurantID, id, voidedBy int, reason string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := r.VoidEntryTx(tx, tenantID, restaurantID, id, voidedBy, reason); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *JournalEntryRepository) VoidEntryTx(tx *sql.Tx, tenantID, restaurantID, id, voidedBy int, reason string) error {
	if tx == nil {
		return errors.New("transaction is required")
	}

	var status string
	err := tx.QueryRow(
		`SELECT status FROM journal_entries WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3 FOR UPDATE`,
		id,
		tenantID,
		restaurantID,
	).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("journal entry not found")
		}
		return err
	}
	if status == domain.JournalStatusVoid {
		return errors.New("journal entry is already void")
	}

	var reversed []postingLine
	var voidDate time.Time
	var fiscalPeriod string
	if status == domain.JournalStatusPosted {
		lines, err := r.getPostingLinesTx(tx, id)
		if err != nil {
			return err
		}

		voidDate = time.Now()
		fiscalPeriod = formatFiscalPeriod(voidDate)
		reversed = make([]postingLine, 0, len(lines))
		for _, line := range lines {
			reversed = append(reversed, postingLine{
				ID:            line.ID,
				AccountID:     line.AccountID,
				DebitAmount:   line.CreditAmount,
				CreditAmount:  line.DebitAmount,
				NormalBalance: line.NormalBalance,
			})
		}
	}

	_, err = tx.Exec(`
		UPDATE journal_entries
		SET status = $1,
			voided_at = CURRENT_TIMESTAMP,
			voided_by = $2,
			void_reason = $3,
			updated_by = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND tenant_id = $5 AND restaurant_id = $6
	`,
		domain.JournalStatusVoid,
		voidedBy,
		reason,
		id,
		tenantID,
		restaurantID,
	)
	if err != nil {
		return err
	}
	if status == domain.JournalStatusPosted {
		if err := r.appendGeneralLedgerTx(tx, tenantID, restaurantID, id, voidDate, fiscalPeriod, reversed); err != nil {
			return err
		}
	}

	return nil
}

type postingLine struct {
	ID            int
	AccountID     int
	DebitAmount   float64
	CreditAmount  float64
	NormalBalance string
}

func (r *JournalEntryRepository) getPostingLinesTx(tx *sql.Tx, journalEntryID int) ([]postingLine, error) {
	rows, err := tx.Query(`
		SELECT
			jel.id, jel.account_id, jel.debit_amount, jel.credit_amount,
			coa.normal_balance
		FROM journal_entry_lines jel
		JOIN chart_of_accounts coa ON coa.id = jel.account_id
		WHERE jel.journal_entry_id = $1
		ORDER BY jel.line_order ASC, jel.id ASC
	`, journalEntryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lines := make([]postingLine, 0)
	for rows.Next() {
		var line postingLine
		if err := rows.Scan(&line.ID, &line.AccountID, &line.DebitAmount, &line.CreditAmount, &line.NormalBalance); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func (r *JournalEntryRepository) appendGeneralLedgerTx(tx *sql.Tx, tenantID, restaurantID, journalEntryID int, entryDate time.Time, fiscalPeriod string, lines []postingLine) error {
	for _, line := range lines {
		var currentBalance float64
		err := tx.QueryRow(`
			SELECT COALESCE(running_balance, 0)
			FROM general_ledger
			WHERE tenant_id = $1 AND restaurant_id = $2 AND account_id = $3
			ORDER BY entry_date DESC, id DESC
			LIMIT 1
		`, tenantID, restaurantID, line.AccountID).Scan(&currentBalance)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if errors.Is(err, sql.ErrNoRows) {
			currentBalance = 0
		}

		newBalance := currentBalance + calculateBalanceDelta(line.NormalBalance, line.DebitAmount, line.CreditAmount)
		if _, err := tx.Exec(`
			INSERT INTO general_ledger (
				tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id,
				entry_date, fiscal_period, debit_amount, credit_amount, running_balance
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10
			)
		`,
			tenantID,
			restaurantID,
			line.AccountID,
			journalEntryID,
			line.ID,
			entryDate,
			fiscalPeriod,
			line.DebitAmount,
			line.CreditAmount,
			newBalance,
		); err != nil {
			return err
		}

		if _, err := tx.Exec(`
			UPDATE chart_of_accounts
			SET current_balance = $1,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
		`, newBalance, line.AccountID, tenantID, restaurantID); err != nil {
			return err
		}
	}
	return nil
}

func (r *JournalEntryRepository) getLines(journalEntryID int) ([]domain.JournalEntryLine, error) {
	rows, err := r.db.Query(`
		SELECT
			id, journal_entry_id, tenant_id, account_id, account_code,
			debit_amount, credit_amount, COALESCE(description, ''),
			COALESCE(description_ar, ''), line_order, COALESCE(tax_code, ''),
			COALESCE(tax_amount, 0), created_at
		FROM journal_entry_lines
		WHERE journal_entry_id = $1
		ORDER BY line_order ASC, id ASC
	`, journalEntryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lines := make([]domain.JournalEntryLine, 0)
	for rows.Next() {
		var line domain.JournalEntryLine
		if err := rows.Scan(
			&line.ID,
			&line.JournalEntryID,
			&line.TenantID,
			&line.AccountID,
			&line.AccountCode,
			&line.DebitAmount,
			&line.CreditAmount,
			&line.Description,
			&line.DescriptionAr,
			&line.LineOrder,
			&line.TaxCode,
			&line.TaxAmount,
			&line.CreatedAt,
		); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func (r *JournalEntryRepository) insertLinesTx(tx *sql.Tx, journalEntryID, tenantID, restaurantID int, lines []domain.JournalEntryLine) error {
	for index, line := range lines {
		accountCode, isHeader, err := lookupAccountForLine(tx, tenantID, restaurantID, line.AccountID)
		if err != nil {
			return err
		}
		if isHeader {
			return fmt.Errorf("account %d is a header account and cannot be posted to", line.AccountID)
		}

		lineOrder := line.LineOrder
		if lineOrder == 0 {
			lineOrder = index + 1
		}

		_, err = tx.Exec(`
			INSERT INTO journal_entry_lines (
				journal_entry_id, tenant_id, account_id, account_code,
				debit_amount, credit_amount, description, description_ar,
				line_order, tax_code, tax_amount
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8,
				$9, $10, $11
			)
		`,
			journalEntryID,
			tenantID,
			line.AccountID,
			accountCode,
			line.DebitAmount,
			line.CreditAmount,
			line.Description,
			line.DescriptionAr,
			lineOrder,
			line.TaxCode,
			line.TaxAmount,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func lookupAccountForLine(tx *sql.Tx, tenantID, restaurantID, accountID int) (string, bool, error) {
	var (
		accountCode string
		isHeader    bool
	)
	err := tx.QueryRow(`
		SELECT account_code, is_header
		FROM chart_of_accounts
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3 AND is_active = TRUE
	`, accountID, tenantID, restaurantID).Scan(&accountCode, &isHeader)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, fmt.Errorf("account %d not found", accountID)
		}
		return "", false, err
	}
	return accountCode, isHeader, nil
}

func validateJournalLines(lines []domain.JournalEntryLine) (float64, float64, error) {
	if len(lines) < 2 {
		return 0, 0, errors.New("journal entry must contain at least two lines")
	}

	var totalDebit float64
	var totalCredit float64
	for _, line := range lines {
		if line.AccountID <= 0 {
			return 0, 0, errors.New("each journal line requires an account_id")
		}
		if line.DebitAmount < 0 || line.CreditAmount < 0 {
			return 0, 0, errors.New("journal amounts cannot be negative")
		}
		if line.DebitAmount > 0 && line.CreditAmount > 0 {
			return 0, 0, errors.New("a journal line cannot be both debit and credit")
		}
		if line.DebitAmount == 0 && line.CreditAmount == 0 {
			return 0, 0, errors.New("a journal line must have a debit or credit amount")
		}
		totalDebit += line.DebitAmount
		totalCredit += line.CreditAmount
	}

	if math.Abs(totalDebit-totalCredit) > 0.009 {
		return 0, 0, errors.New("journal entry is not balanced")
	}

	return totalDebit, totalCredit, nil
}

func validatePostingLines(lines []postingLine) (float64, float64, error) {
	if len(lines) < 2 {
		return 0, 0, errors.New("journal entry must contain at least two lines")
	}
	var totalDebit float64
	var totalCredit float64
	for _, line := range lines {
		totalDebit += line.DebitAmount
		totalCredit += line.CreditAmount
	}
	return totalDebit, totalCredit, nil
}

func formatFiscalPeriod(entryDate time.Time) string {
	return entryDate.Format("2006-01")
}

func (r *JournalEntryRepository) ensureFiscalPeriodOpenTx(tx *sql.Tx, tenantID, restaurantID int, entryDate time.Time, fiscalPeriod string) error {
	var status string
	err := tx.QueryRow(`
		SELECT status
		FROM fiscal_periods
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND (period_name = $3 OR ($4 >= start_date AND $4 <= end_date))
		ORDER BY start_date DESC, id DESC
		LIMIT 1
	`, tenantID, restaurantID, fiscalPeriod, entryDate).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	if strings.EqualFold(status, "closed") {
		return errors.New("fiscal period is closed; posting is not allowed")
	}
	return nil
}

func calculateBalanceDelta(normalBalance string, debitAmount, creditAmount float64) float64 {
	if normalBalance == domain.NormalBalanceCredit {
		return creditAmount - debitAmount
	}
	return debitAmount - creditAmount
}

func (r *JournalEntryRepository) generateEntryNumberTx(tx *sql.Tx, tenantID, restaurantID int, entryDate time.Time) (string, error) {
	year := entryDate.Format("2006")
	prefix := fmt.Sprintf("JE-%s-", year)

	var nextNumber int
	err := tx.QueryRow(`
		SELECT COALESCE(MAX(CAST(RIGHT(entry_number, 5) AS INT)), 0) + 1
		FROM journal_entries
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND entry_number LIKE $3
	`, tenantID, restaurantID, prefix+"%").Scan(&nextNumber)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%05d", prefix, nextNumber), nil
}

type journalEntryScanner interface {
	Scan(dest ...interface{}) error
}

func scanJournalEntry(scanner journalEntryScanner) (domain.JournalEntry, error) {
	var (
		entry       domain.JournalEntry
		referenceID sql.NullInt64
		postedAt    sql.NullTime
		postedBy    sql.NullInt64
		voidedAt    sql.NullTime
		voidedBy    sql.NullInt64
		approvedAt  sql.NullTime
		approvedBy  sql.NullInt64
		createdBy   sql.NullInt64
		updatedBy   sql.NullInt64
		attachments string
	)

	err := scanner.Scan(
		&entry.ID,
		&entry.TenantID,
		&entry.RestaurantID,
		&entry.EntryNumber,
		&entry.EntryDate,
		&entry.FiscalPeriod,
		&entry.Description,
		&entry.DescriptionAr,
		&entry.ReferenceType,
		&referenceID,
		&entry.ReferenceNumber,
		&entry.EntryType,
		&entry.TotalDebit,
		&entry.TotalCredit,
		&entry.Status,
		&postedAt,
		&postedBy,
		&voidedAt,
		&voidedBy,
		&entry.VoidReason,
		&entry.IsApproved,
		&approvedBy,
		&approvedAt,
		&entry.Notes,
		&attachments,
		&entry.CreatedAt,
		&entry.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		return domain.JournalEntry{}, err
	}

	if referenceID.Valid {
		value := int(referenceID.Int64)
		entry.ReferenceID = &value
	}
	if postedAt.Valid {
		value := postedAt.Time
		entry.PostedAt = &value
	}
	if postedBy.Valid {
		value := int(postedBy.Int64)
		entry.PostedBy = &value
	}
	if voidedAt.Valid {
		value := voidedAt.Time
		entry.VoidedAt = &value
	}
	if voidedBy.Valid {
		value := int(voidedBy.Int64)
		entry.VoidedBy = &value
	}
	if approvedAt.Valid {
		value := approvedAt.Time
		entry.ApprovedAt = &value
	}
	if approvedBy.Valid {
		value := int(approvedBy.Int64)
		entry.ApprovedBy = &value
	}
	if createdBy.Valid {
		value := int(createdBy.Int64)
		entry.CreatedBy = &value
	}
	if updatedBy.Valid {
		value := int(updatedBy.Int64)
		entry.UpdatedBy = &value
	}
	if attachments != "" {
		entry.Attachments = json.RawMessage(attachments)
	}

	return entry, nil
}
