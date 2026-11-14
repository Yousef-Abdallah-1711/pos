package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/database"
)

type AccountingAdvancedRepository struct {
	db          *sql.DB
	journalRepo *JournalEntryRepository
}

func NewAccountingAdvancedRepository(db *sql.DB, journalRepo *JournalEntryRepository) *AccountingAdvancedRepository {
	return &AccountingAdvancedRepository{
		db:          db,
		journalRepo: journalRepo,
	}
}

func (r *AccountingAdvancedRepository) ListBankReconciliationSessions(tenantID, restaurantID int) ([]domain.BankReconciliationSession, error) {
	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, bank_account_id, statement_date, statement_balance,
			book_balance, difference, status, COALESCE(notes, ''), completed_at, completed_by,
			created_at, updated_at, created_by, updated_by
		FROM bank_reconciliation_sessions
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY statement_date DESC, id DESC
	`, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.BankReconciliationSession, 0)
	for rows.Next() {
		item, err := scanBankReconciliationSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AccountingAdvancedRepository) GetBankReconciliationSessionByID(tenantID, restaurantID, id int) (*domain.BankReconciliationSession, error) {
	item, err := scanBankReconciliationSession(r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, bank_account_id, statement_date, statement_balance,
			book_balance, difference, status, COALESCE(notes, ''), completed_at, completed_by,
			created_at, updated_at, created_by, updated_by
		FROM bank_reconciliation_sessions
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *AccountingAdvancedRepository) CreateBankReconciliationSession(item *domain.BankReconciliationSession) (int, error) {
	if item == nil {
		return 0, errors.New("reconciliation payload is required")
	}
	var id int
	err := r.db.QueryRow(`
		INSERT INTO bank_reconciliation_sessions (
			tenant_id, restaurant_id, bank_account_id, statement_date,
			statement_balance, book_balance, difference, status, notes,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, $9,
			$10, $11
		)
		RETURNING id
	`,
		item.TenantID,
		item.RestaurantID,
		item.BankAccountID,
		item.StatementDate,
		item.StatementBalance,
		item.BookBalance,
		item.Difference,
		defaultString(item.Status, "draft"),
		nullString(item.Notes),
		item.CreatedBy,
		item.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountingAdvancedRepository) AddBankReconciliationLine(
	tenantID, restaurantID, reconciliationID int,
	line *domain.BankReconciliationLine,
) error {
	if line == nil {
		return errors.New("reconciliation line is required")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	if err := tx.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM bank_reconciliation_sessions
			WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
		)
	`, tenantID, restaurantID, reconciliationID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return errors.New("reconciliation session not found")
	}

	_, err = tx.Exec(`
		INSERT INTO bank_reconciliation_lines (
			reconciliation_id, tenant_id, journal_entry_id, line_description,
			statement_amount, book_amount, match_status
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7
		)
	`,
		reconciliationID,
		tenantID,
		line.JournalEntryID,
		nullString(line.LineDescription),
		line.StatementAmount,
		line.BookAmount,
		defaultString(line.MatchStatus, "matched"),
	)
	if err != nil {
		return err
	}

	if err := refreshReconciliationSummaryTx(tx, tenantID, restaurantID, reconciliationID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AccountingAdvancedRepository) CompleteBankReconciliationSession(tenantID, restaurantID, reconciliationID, completedBy int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := refreshReconciliationSummaryTx(tx, tenantID, restaurantID, reconciliationID); err != nil {
		return err
	}

	var (
		bankAccountID int
		statementDate time.Time
		bookBalance   float64
		difference    float64
	)
	err = tx.QueryRow(`
		SELECT bank_account_id, statement_date, book_balance, difference
		FROM bank_reconciliation_sessions
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, reconciliationID).Scan(
		&bankAccountID,
		&statementDate,
		&bookBalance,
		&difference,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("reconciliation session not found")
		}
		return err
	}

	status := "in_progress"
	if math.Abs(difference) < 0.01 {
		status = "completed"
	}

	_, err = tx.Exec(`
		UPDATE bank_reconciliation_sessions
		SET status = $1,
			completed_at = CURRENT_TIMESTAMP,
			completed_by = $2,
			updated_by = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE tenant_id = $3 AND restaurant_id = $4 AND id = $5
	`, status, completedBy, tenantID, restaurantID, reconciliationID)
	if err != nil {
		return err
	}

	if status == "completed" {
		_, err = tx.Exec(`
			UPDATE bank_accounts_ledger
			SET last_reconciled_date = $1,
				last_reconciled_balance = $2,
				updated_at = CURRENT_TIMESTAMP
			WHERE tenant_id = $3 AND restaurant_id = $4 AND id = $5
		`, statementDate, bookBalance, tenantID, restaurantID, bankAccountID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *AccountingAdvancedRepository) ListRecurringJournalTemplates(tenantID, restaurantID int) ([]domain.RecurringJournalTemplate, error) {
	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, template_code, template_name, COALESCE(description, ''),
			frequency, start_date, next_run_date, last_run_date, is_active, auto_post, lines_json,
			created_at, updated_at, created_by, updated_by
		FROM recurring_journal_templates
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY next_run_date ASC, id DESC
	`, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.RecurringJournalTemplate, 0)
	for rows.Next() {
		item, err := scanRecurringTemplate(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AccountingAdvancedRepository) GetRecurringJournalTemplateByID(tenantID, restaurantID, id int) (*domain.RecurringJournalTemplate, error) {
	item, err := scanRecurringTemplate(r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, template_code, template_name, COALESCE(description, ''),
			frequency, start_date, next_run_date, last_run_date, is_active, auto_post, lines_json,
			created_at, updated_at, created_by, updated_by
		FROM recurring_journal_templates
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *AccountingAdvancedRepository) CreateRecurringJournalTemplate(item *domain.RecurringJournalTemplate) (int, error) {
	if item == nil {
		return 0, errors.New("template payload is required")
	}
	if len(item.LinesJSON) == 0 {
		item.LinesJSON = []byte("[]")
	}
	if item.NextRunDate.IsZero() {
		item.NextRunDate = item.StartDate
	}

	var id int
	err := r.db.QueryRow(`
		INSERT INTO recurring_journal_templates (
			tenant_id, restaurant_id, template_code, template_name, description, frequency,
			start_date, next_run_date, is_active, auto_post, lines_json, created_by, updated_by
		) SELECT
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13
		WHERE EXISTS (
			SELECT 1 FROM restaurants
			WHERE id = $2 AND tenant_id = $1
		)
		RETURNING id
	`,
		item.TenantID,
		item.RestaurantID,
		item.TemplateCode,
		item.TemplateName,
		nullString(item.Description),
		item.Frequency,
		item.StartDate,
		item.NextRunDate,
		item.IsActive,
		item.AutoPost,
		item.LinesJSON,
		item.CreatedBy,
		item.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountingAdvancedRepository) UpdateRecurringJournalTemplate(item *domain.RecurringJournalTemplate) error {
	if item == nil {
		return errors.New("template payload is required")
	}
	if len(item.LinesJSON) == 0 {
		item.LinesJSON = []byte("[]")
	}
	result, err := r.db.Exec(`
		UPDATE recurring_journal_templates
		SET template_code = $1,
			template_name = $2,
			description = $3,
			frequency = $4,
			start_date = $5,
			next_run_date = $6,
			is_active = $7,
			auto_post = $8,
			lines_json = $9,
			updated_by = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE tenant_id = $11 AND restaurant_id = $12 AND id = $13
	`,
		item.TemplateCode,
		item.TemplateName,
		nullString(item.Description),
		item.Frequency,
		item.StartDate,
		item.NextRunDate,
		item.IsActive,
		item.AutoPost,
		item.LinesJSON,
		item.UpdatedBy,
		item.TenantID,
		item.RestaurantID,
		item.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("recurring template not found")
	}
	return nil
}

func (r *AccountingAdvancedRepository) DeleteRecurringJournalTemplate(tenantID, restaurantID, id int) error {
	result, err := r.db.Exec(`
		DELETE FROM recurring_journal_templates
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("recurring template not found")
	}
	return nil
}

func (r *AccountingAdvancedRepository) RunRecurringJournalTemplate(
	tenantID, restaurantID, templateID, userID int,
	runDate time.Time,
) (int, error) {
	template, err := r.GetRecurringJournalTemplateByID(tenantID, restaurantID, templateID)
	if err != nil {
		return 0, err
	}
	if template == nil {
		return 0, errors.New("recurring template not found")
	}
	if !template.IsActive {
		return 0, errors.New("recurring template is inactive")
	}

	lines := make([]domain.JournalEntryLineInput, 0)
	if err := json.Unmarshal(template.LinesJSON, &lines); err != nil {
		return 0, fmt.Errorf("invalid recurring template lines_json: %w", err)
	}
	if len(lines) == 0 {
		return 0, errors.New("recurring template has no journal lines")
	}

	totalDebit := 0.0
	totalCredit := 0.0
	entryLines := make([]domain.JournalEntryLine, 0, len(lines))
	for idx, line := range lines {
		totalDebit += line.DebitAmount
		totalCredit += line.CreditAmount
		entryLines = append(entryLines, domain.JournalEntryLine{
			AccountID:     line.AccountID,
			DebitAmount:   line.DebitAmount,
			CreditAmount:  line.CreditAmount,
			Description:   line.Description,
			DescriptionAr: line.DescriptionAr,
			LineOrder:     idx + 1,
			TaxCode:       line.TaxCode,
			TaxAmount:     line.TaxAmount,
		})
	}
	if math.Abs(totalDebit-totalCredit) > 0.01 {
		return 0, errors.New("recurring template is not balanced")
	}

	if runDate.IsZero() {
		runDate = time.Now()
	}
	fiscalPeriod := fmt.Sprintf("%04d-%02d", runDate.Year(), runDate.Month())
	entry := &domain.JournalEntry{
		TenantID:      tenantID,
		RestaurantID:  restaurantID,
		EntryDate:     runDate,
		FiscalPeriod:  fiscalPeriod,
		Description:   fmt.Sprintf("%s (recurring)", template.TemplateName),
		ReferenceType: "recurring_template",
		ReferenceID:   &templateID,
		EntryType:     "recurring",
		Status:        domain.JournalStatusDraft,
		CreatedBy:     &userID,
		UpdatedBy:     &userID,
		Lines:         entryLines,
	}

	entryID, err := r.journalRepo.CreateJournalEntry(entry)
	if err != nil {
		return 0, err
	}
	if template.AutoPost {
		if err := r.journalRepo.PostEntry(tenantID, restaurantID, entryID, userID); err != nil {
			return 0, err
		}
	}

	nextRun := nextRecurringRunDate(template.Frequency, runDate)
	_, err = r.db.Exec(`
		UPDATE recurring_journal_templates
		SET last_run_date = $1,
			next_run_date = $2,
			updated_by = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE tenant_id = $4 AND restaurant_id = $5 AND id = $6
	`, runDate, nextRun, userID, tenantID, restaurantID, templateID)
	if err != nil {
		return 0, err
	}

	return entryID, nil
}

func (r *AccountingAdvancedRepository) UpsertBudget(
	tenantID, restaurantID, accountID, fiscalYear, periodMonth int,
	budgetAmount float64,
	currency, notes string,
) (int, error) {
	if fiscalYear < 2000 || periodMonth < 1 || periodMonth > 12 {
		return 0, errors.New("invalid fiscal period")
	}
	var id int
	err := r.db.QueryRow(`
		INSERT INTO account_budgets (
			tenant_id, restaurant_id, account_id, fiscal_year, period_month,
			budget_amount, currency, notes
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8
		)
		ON CONFLICT (tenant_id, restaurant_id, account_id, fiscal_year, period_month)
		DO UPDATE SET
			budget_amount = EXCLUDED.budget_amount,
			currency = EXCLUDED.currency,
			notes = EXCLUDED.notes,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`,
		tenantID,
		restaurantID,
		accountID,
		fiscalYear,
		periodMonth,
		budgetAmount,
		defaultString(strings.ToUpper(strings.TrimSpace(currency)), "EGP"),
		nullString(notes),
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountingAdvancedRepository) ListBudgetVariance(tenantID, restaurantID, fiscalYear int) ([]domain.BudgetEntry, error) {
	rows, err := r.db.Query(`
		SELECT
			b.id, b.tenant_id, b.restaurant_id, b.account_id,
			a.account_code, a.account_name,
			b.fiscal_year, b.period_month, b.budget_amount,
			COALESCE((
				SELECT SUM(gl.debit_amount - gl.credit_amount)
				FROM general_ledger gl
				WHERE gl.tenant_id = b.tenant_id
				  AND gl.restaurant_id = b.restaurant_id
				  AND gl.account_id = b.account_id
				  AND EXTRACT(YEAR FROM gl.entry_date) = b.fiscal_year
				  AND EXTRACT(MONTH FROM gl.entry_date) = b.period_month
			), 0) AS actual_amount,
			b.currency, COALESCE(b.notes, ''), b.created_at, b.updated_at
		FROM account_budgets b
		JOIN chart_of_accounts a
		  ON a.id = b.account_id
		 AND a.tenant_id = b.tenant_id
		 AND a.restaurant_id = b.restaurant_id
		WHERE b.tenant_id = $1
		  AND b.restaurant_id = $2
		  AND ($3 = 0 OR b.fiscal_year = $3)
		ORDER BY b.fiscal_year DESC, b.period_month ASC, a.account_code ASC
	`, tenantID, restaurantID, fiscalYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.BudgetEntry, 0)
	for rows.Next() {
		var item domain.BudgetEntry
		if err := rows.Scan(
			&item.ID,
			&item.TenantID,
			&item.RestaurantID,
			&item.AccountID,
			&item.AccountCode,
			&item.AccountName,
			&item.FiscalYear,
			&item.PeriodMonth,
			&item.BudgetAmount,
			&item.ActualAmount,
			&item.Currency,
			&item.Notes,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.VarianceAmount = item.ActualAmount - item.BudgetAmount
		if math.Abs(item.BudgetAmount) > 0.0001 {
			item.VariancePct = (item.VarianceAmount / item.BudgetAmount) * 100
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AccountingAdvancedRepository) UpsertExchangeRate(rate *domain.ExchangeRate) (int, error) {
	if rate == nil {
		return 0, errors.New("exchange rate payload is required")
	}
	var id int
	err := r.db.QueryRow(`
		INSERT INTO exchange_rates (
			tenant_id, restaurant_id, base_currency, quote_currency, rate,
			effective_date, source, is_active
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8
		)
		ON CONFLICT (tenant_id, restaurant_id, base_currency, quote_currency, effective_date)
		DO UPDATE SET
			rate = EXCLUDED.rate,
			source = EXCLUDED.source,
			is_active = EXCLUDED.is_active,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`,
		rate.TenantID,
		rate.RestaurantID,
		strings.ToUpper(strings.TrimSpace(rate.BaseCurrency)),
		strings.ToUpper(strings.TrimSpace(rate.QuoteCurrency)),
		rate.Rate,
		rate.EffectiveDate,
		nullString(rate.Source),
		rate.IsActive,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountingAdvancedRepository) ListExchangeRates(tenantID, restaurantID int, baseCurrency, quoteCurrency string) ([]domain.ExchangeRate, error) {
	baseCurrency = strings.ToUpper(strings.TrimSpace(baseCurrency))
	quoteCurrency = strings.ToUpper(strings.TrimSpace(quoteCurrency))
	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, base_currency, quote_currency, rate,
			effective_date, COALESCE(source, ''), is_active, created_at, updated_at
		FROM exchange_rates
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND ($3 = '' OR base_currency = $3)
		  AND ($4 = '' OR quote_currency = $4)
		ORDER BY effective_date DESC, id DESC
	`, tenantID, restaurantID, baseCurrency, quoteCurrency)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.ExchangeRate, 0)
	for rows.Next() {
		var item domain.ExchangeRate
		if err := rows.Scan(
			&item.ID,
			&item.TenantID,
			&item.RestaurantID,
			&item.BaseCurrency,
			&item.QuoteCurrency,
			&item.Rate,
			&item.EffectiveDate,
			&item.Source,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AccountingAdvancedRepository) ConvertAmount(
	tenantID, restaurantID int,
	baseCurrency, quoteCurrency string,
	amount float64,
	asOf time.Time,
) (float64, float64, error) {
	baseCurrency = strings.ToUpper(strings.TrimSpace(baseCurrency))
	quoteCurrency = strings.ToUpper(strings.TrimSpace(quoteCurrency))
	if baseCurrency == "" || quoteCurrency == "" {
		return 0, 0, errors.New("base_currency and quote_currency are required")
	}
	if baseCurrency == quoteCurrency {
		return amount, 1, nil
	}
	if asOf.IsZero() {
		asOf = time.Now()
	}
	var rate float64
	err := r.db.QueryRow(`
		SELECT rate
		FROM exchange_rates
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND base_currency = $3
		  AND quote_currency = $4
		  AND effective_date <= $5
		  AND is_active = TRUE
		ORDER BY effective_date DESC, id DESC
		LIMIT 1
	`, tenantID, restaurantID, baseCurrency, quoteCurrency, asOf).Scan(&rate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, fmt.Errorf("no exchange rate found for %s/%s", baseCurrency, quoteCurrency)
		}
		return 0, 0, err
	}
	return amount * rate, rate, nil
}

func (r *AccountingAdvancedRepository) CreateSupplierPaymentSchedule(item *domain.SupplierPaymentSchedule) (int, error) {
	if item == nil {
		return 0, errors.New("payment schedule payload is required")
	}
	var id int
	err := r.db.QueryRow(`
		INSERT INTO supplier_payment_schedules (
			tenant_id, restaurant_id, supplier_invoice_id, supplier_id, scheduled_date,
			amount, currency, reminder_days, status, notes, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11, $12
		)
		RETURNING id
	`,
		item.TenantID,
		item.RestaurantID,
		item.SupplierInvoiceID,
		item.SupplierID,
		item.ScheduledDate,
		item.Amount,
		defaultString(strings.ToUpper(strings.TrimSpace(item.Currency)), "EGP"),
		defaultInt(item.ReminderDays, 7),
		defaultString(item.Status, "scheduled"),
		nullString(item.Notes),
		item.CreatedBy,
		item.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountingAdvancedRepository) ListSupplierPaymentSchedules(tenantID, restaurantID int) ([]domain.SupplierPaymentSchedule, error) {
	rows, err := r.db.Query(`
		SELECT
			s.id, s.tenant_id, s.restaurant_id, s.supplier_invoice_id, s.supplier_id,
			COALESCE(sp.company_name, ''), s.scheduled_date, s.amount, s.currency, s.reminder_days,
			s.status, s.reminder_sent_at, COALESCE(s.notes, ''), s.created_at, s.updated_at, s.created_by, s.updated_by
		FROM supplier_payment_schedules s
		JOIN suppliers sp
		  ON sp.id = s.supplier_id
		 AND sp.tenant_id = s.tenant_id
		 AND sp.restaurant_id = s.restaurant_id
		WHERE s.tenant_id = $1 AND s.restaurant_id = $2
		ORDER BY s.scheduled_date ASC, s.id DESC
	`, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.SupplierPaymentSchedule, 0)
	for rows.Next() {
		item, err := scanSupplierPaymentSchedule(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AccountingAdvancedRepository) ListSupplierPaymentReminders(tenantID, restaurantID, dueWithinDays int) ([]domain.SupplierPaymentSchedule, error) {
	if dueWithinDays <= 0 {
		dueWithinDays = 7
	}
	rows, err := r.db.Query(`
		SELECT
			s.id, s.tenant_id, s.restaurant_id, s.supplier_invoice_id, s.supplier_id,
			COALESCE(sp.company_name, ''), s.scheduled_date, s.amount, s.currency, s.reminder_days,
			s.status, s.reminder_sent_at, COALESCE(s.notes, ''), s.created_at, s.updated_at, s.created_by, s.updated_by
		FROM supplier_payment_schedules s
		JOIN suppliers sp
		  ON sp.id = s.supplier_id
		 AND sp.tenant_id = s.tenant_id
		 AND sp.restaurant_id = s.restaurant_id
		WHERE s.tenant_id = $1
		  AND s.restaurant_id = $2
		  AND s.status IN ('scheduled', 'reminded')
		  AND s.scheduled_date <= (CURRENT_DATE + ($3 * INTERVAL '1 day'))
		  AND (s.reminder_sent_at IS NULL OR s.reminder_sent_at::date < CURRENT_DATE)
		ORDER BY s.scheduled_date ASC
	`, tenantID, restaurantID, dueWithinDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.SupplierPaymentSchedule, 0)
	for rows.Next() {
		item, err := scanSupplierPaymentSchedule(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AccountingAdvancedRepository) MarkSupplierPaymentReminderSent(tenantID, restaurantID, scheduleID int) error {
	result, err := r.db.Exec(`
		UPDATE supplier_payment_schedules
		SET status = CASE WHEN status = 'scheduled' THEN 'reminded' ELSE status END,
			reminder_sent_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, scheduleID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("supplier payment schedule not found")
	}
	return nil
}

func (r *AccountingAdvancedRepository) GetFinancialKPIs(tenantID, restaurantID int, dateFrom, dateTo time.Time) (map[string]interface{}, error) {
	revenue := 0.0
	expenses := 0.0
	cashPosition := 0.0

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM orders
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND payment_status = 'paid'
		  AND created_at >= $3
		  AND created_at < $4
	`, tenantID, restaurantID, dateFrom, dateTo).Scan(&revenue); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(jel.debit_amount - jel.credit_amount), 0)
		FROM journal_entry_lines jel
		JOIN journal_entries je ON je.id = jel.journal_entry_id
		JOIN chart_of_accounts coa
		  ON coa.id = jel.account_id
		 AND coa.tenant_id = jel.tenant_id
		WHERE jel.tenant_id = $1
		  AND je.restaurant_id = $2
		  AND je.status = 'posted'
		  AND coa.account_type = 'expense'
		  AND je.entry_date >= $3
		  AND je.entry_date < $4
	`, tenantID, restaurantID, dateFrom, dateTo).Scan(&expenses); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(current_balance), 0)
		FROM bank_accounts_ledger
		WHERE tenant_id = $1 AND restaurant_id = $2 AND is_active = TRUE
	`, tenantID, restaurantID).Scan(&cashPosition); err != nil {
		return nil, err
	}

	trendRows, err := r.db.Query(`
		SELECT
			DATE_TRUNC('month', created_at)::date AS month_bucket,
			COALESCE(SUM(total_amount), 0) AS revenue
		FROM orders
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND payment_status = 'paid'
		  AND created_at >= $3
		  AND created_at < $4
		GROUP BY 1
		ORDER BY 1 ASC
	`, tenantID, restaurantID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer trendRows.Close()

	trend := make([]map[string]interface{}, 0)
	for trendRows.Next() {
		var bucket time.Time
		var value float64
		if err := trendRows.Scan(&bucket, &value); err != nil {
			return nil, err
		}
		trend = append(trend, map[string]interface{}{
			"month":   bucket.Format("2006-01"),
			"revenue": value,
		})
	}
	if err := trendRows.Err(); err != nil {
		return nil, err
	}

	expenseRows, err := r.db.Query(`
		SELECT
			COALESCE(NULLIF(coa.account_subtype, ''), coa.account_name) AS label,
			COALESCE(SUM(jel.debit_amount - jel.credit_amount), 0) AS amount
		FROM journal_entry_lines jel
		JOIN journal_entries je ON je.id = jel.journal_entry_id
		JOIN chart_of_accounts coa
		  ON coa.id = jel.account_id
		 AND coa.tenant_id = jel.tenant_id
		WHERE jel.tenant_id = $1
		  AND je.restaurant_id = $2
		  AND je.status = 'posted'
		  AND coa.account_type = 'expense'
		  AND je.entry_date >= $3
		  AND je.entry_date < $4
		GROUP BY 1
		ORDER BY amount DESC
		LIMIT 10
	`, tenantID, restaurantID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer expenseRows.Close()

	breakdown := make([]map[string]interface{}, 0)
	for expenseRows.Next() {
		var label string
		var amount float64
		if err := expenseRows.Scan(&label, &amount); err != nil {
			return nil, err
		}
		breakdown = append(breakdown, map[string]interface{}{
			"label":  label,
			"amount": amount,
		})
	}
	if err := expenseRows.Err(); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"summary": map[string]interface{}{
			"revenue":       revenue,
			"expenses":      expenses,
			"net_profit":    revenue - expenses,
			"cash_position": cashPosition,
		},
		"revenue_trend":     trend,
		"expense_breakdown": breakdown,
		"generated_at":      time.Now().UTC(),
	}, nil
}

func (r *AccountingAdvancedRepository) GetMenuItemProfitability(tenantID, restaurantID int, dateFrom, dateTo time.Time, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.Query(`
		SELECT
			oi.product_id,
			COALESCE(NULLIF(oi.product_name, ''), p.name, CONCAT('Product #', oi.product_id::text)) AS product_name,
			COALESCE(SUM(oi.quantity - oi.cancelled_quantity), 0) AS quantity,
			COALESCE(SUM(oi.total_price), 0) AS revenue,
			COALESCE(SUM((oi.quantity - oi.cancelled_quantity) * COALESCE(p.cost_price, 0)), 0) AS cogs
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE o.tenant_id = $1
		  AND o.restaurant_id = $2
		  AND o.payment_status = 'paid'
		  AND o.created_at >= $3
		  AND o.created_at < $4
		GROUP BY oi.product_id, product_name
		ORDER BY revenue DESC
		LIMIT $5
	`, tenantID, restaurantID, dateFrom, dateTo, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		var (
			productID   int64
			productName string
			quantity    float64
			revenue     float64
			cogs        float64
		)
		if err := rows.Scan(&productID, &productName, &quantity, &revenue, &cogs); err != nil {
			return nil, err
		}
		grossProfit := revenue - cogs
		margin := 0.0
		if revenue > 0 {
			margin = (grossProfit / revenue) * 100
		}
		items = append(items, map[string]interface{}{
			"product_id":       productID,
			"name":             productName,
			"quantity":         quantity,
			"revenue":          revenue,
			"cogs":             cogs,
			"gross_profit":     grossProfit,
			"gross_margin_pct": margin,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AccountingAdvancedRepository) GetBreakEvenAnalysis(tenantID, restaurantID int, dateFrom, dateTo time.Time) (map[string]interface{}, error) {
	revenue := 0.0
	variableCosts := 0.0
	totalExpenses := 0.0

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM orders
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND payment_status = 'paid'
		  AND created_at >= $3
		  AND created_at < $4
	`, tenantID, restaurantID, dateFrom, dateTo).Scan(&revenue); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM((oi.quantity - oi.cancelled_quantity) * COALESCE(p.cost_price, 0)), 0)
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE o.tenant_id = $1
		  AND o.restaurant_id = $2
		  AND o.payment_status = 'paid'
		  AND o.created_at >= $3
		  AND o.created_at < $4
	`, tenantID, restaurantID, dateFrom, dateTo).Scan(&variableCosts); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(jel.debit_amount - jel.credit_amount), 0)
		FROM journal_entry_lines jel
		JOIN journal_entries je ON je.id = jel.journal_entry_id
		JOIN chart_of_accounts coa
		  ON coa.id = jel.account_id
		 AND coa.tenant_id = jel.tenant_id
		WHERE jel.tenant_id = $1
		  AND je.restaurant_id = $2
		  AND je.status = 'posted'
		  AND coa.account_type = 'expense'
		  AND je.entry_date >= $3
		  AND je.entry_date < $4
	`, tenantID, restaurantID, dateFrom, dateTo).Scan(&totalExpenses); err != nil {
		return nil, err
	}

	fixedCosts := totalExpenses - variableCosts
	if fixedCosts < 0 {
		fixedCosts = 0
	}
	contribution := revenue - variableCosts
	contributionMarginRatio := 0.0
	if revenue > 0 {
		contributionMarginRatio = contribution / revenue
	}
	breakEvenRevenue := 0.0
	if contributionMarginRatio > 0 {
		breakEvenRevenue = fixedCosts / contributionMarginRatio
	}
	breakEvenGap := revenue - breakEvenRevenue

	return map[string]interface{}{
		"summary": map[string]interface{}{
			"revenue":                   revenue,
			"variable_costs":            variableCosts,
			"fixed_costs":               fixedCosts,
			"contribution_margin":       contribution,
			"contribution_margin_ratio": contributionMarginRatio * 100,
			"break_even_revenue":        breakEvenRevenue,
			"break_even_gap":            breakEvenGap,
			"is_above_break_even":       breakEvenGap >= 0,
		},
		"generated_at": time.Now().UTC(),
	}, nil
}

func (r *AccountingAdvancedRepository) GetPredictiveCashFlow(tenantID, restaurantID, historyMonths, forecastMonths int) (map[string]interface{}, error) {
	if historyMonths <= 0 {
		historyMonths = 6
	}
	if forecastMonths <= 0 {
		forecastMonths = 3
	}

	startMonth := beginningOfMonth(time.Now().AddDate(0, -historyMonths+1, 0))
	rows, err := r.db.Query(`
		WITH received AS (
			SELECT DATE_TRUNC('month', payment_date)::date AS bucket, COALESCE(SUM(amount), 0) AS amount
			FROM payments_received
			WHERE tenant_id = $1
			  AND restaurant_id = $2
			  AND payment_date >= $3
			  AND status <> 'void'
			GROUP BY 1
		),
		made AS (
			SELECT DATE_TRUNC('month', payment_date)::date AS bucket, COALESCE(SUM(amount), 0) AS amount
			FROM payments_made
			WHERE tenant_id = $1
			  AND restaurant_id = $2
			  AND payment_date >= $3
			  AND status <> 'void'
			GROUP BY 1
		)
		SELECT
			COALESCE(r.bucket, m.bucket) AS bucket,
			COALESCE(r.amount, 0) AS cash_in,
			COALESCE(m.amount, 0) AS cash_out
		FROM received r
		FULL OUTER JOIN made m ON m.bucket = r.bucket
		ORDER BY bucket ASC
	`, tenantID, restaurantID, startMonth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	historyMap := make(map[string]map[string]float64)
	for rows.Next() {
		var (
			bucket  time.Time
			cashIn  float64
			cashOut float64
		)
		if err := rows.Scan(&bucket, &cashIn, &cashOut); err != nil {
			return nil, err
		}
		key := bucket.Format("2006-01")
		historyMap[key] = map[string]float64{
			"cash_in":  cashIn,
			"cash_out": cashOut,
			"net_cash": cashIn - cashOut,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	history := make([]map[string]interface{}, 0, historyMonths)
	netSeries := make([]float64, 0, historyMonths)
	for i := 0; i < historyMonths; i++ {
		month := startMonth.AddDate(0, i, 0)
		key := month.Format("2006-01")
		values := historyMap[key]
		if values == nil {
			values = map[string]float64{"cash_in": 0, "cash_out": 0, "net_cash": 0}
		}
		netSeries = append(netSeries, values["net_cash"])
		history = append(history, map[string]interface{}{
			"month":    key,
			"cash_in":  values["cash_in"],
			"cash_out": values["cash_out"],
			"net_cash": values["net_cash"],
		})
	}

	avgNet := averageFloat(netSeries)
	trendSlope := 0.0
	if len(netSeries) >= 2 {
		trendSlope = (netSeries[len(netSeries)-1] - netSeries[0]) / float64(len(netSeries)-1)
	}

	forecast := make([]map[string]interface{}, 0, forecastMonths)
	baseNet := 0.0
	if len(netSeries) > 0 {
		baseNet = netSeries[len(netSeries)-1]
	}
	for i := 1; i <= forecastMonths; i++ {
		month := beginningOfMonth(time.Now().AddDate(0, i, 0))
		projectedNet := (0.7 * baseNet) + (0.3 * avgNet) + (trendSlope * float64(i))
		forecast = append(forecast, map[string]interface{}{
			"month":              month.Format("2006-01"),
			"projected_net_cash": projectedNet,
		})
	}

	return map[string]interface{}{
		"history":  history,
		"forecast": forecast,
		"assumptions": map[string]interface{}{
			"history_months":   historyMonths,
			"forecast_months":  forecastMonths,
			"average_net_cash": avgNet,
			"trend_slope":      trendSlope,
		},
		"generated_at": time.Now().UTC(),
	}, nil
}

func (r *AccountingAdvancedRepository) GetHRCostAnalysis(tenantID, restaurantID int, dateFrom, dateTo time.Time) (map[string]interface{}, error) {
	revenue := 0.0
	payrollCost := 0.0
	overtimeCost := 0.0
	insuranceCost := 0.0

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM orders
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND payment_status = 'paid'
		  AND created_at >= $3
		  AND created_at < $4
	`, tenantID, restaurantID, dateFrom, dateTo).Scan(&revenue); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(net_salary), 0), COALESCE(SUM(overtime_amount), 0)
		FROM salaries
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND status = 'paid'
		  AND pay_period_end >= $3::date
		  AND pay_period_end < $4::date
	`, tenantID, restaurantID, dateFrom, dateTo).Scan(&payrollCost, &overtimeCost); err != nil {
		return nil, err
	}

	monthSpan := int(math.Max(1, math.Ceil(dateTo.Sub(dateFrom).Hours()/(24*30))))
	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(monthly_premium), 0)
		FROM insurance_tracking
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND status = 'active'
	`, tenantID, restaurantID).Scan(&insuranceCost); err != nil {
		return nil, err
	}
	insuranceCost *= float64(monthSpan)

	totalHRCost := payrollCost + insuranceCost
	laborPct := 0.0
	if revenue > 0 {
		laborPct = (totalHRCost / revenue) * 100
	}

	departmentRows, err := r.db.Query(`
		SELECT
			COALESCE(NULLIF(e.department, ''), 'Unassigned') AS department,
			COALESCE(SUM(s.net_salary), 0) AS payroll_cost
		FROM salaries s
		JOIN employees e ON e.id = s.employee_id
		WHERE s.tenant_id = $1
		  AND s.restaurant_id = $2
		  AND s.status = 'paid'
		  AND s.pay_period_end >= $3::date
		  AND s.pay_period_end < $4::date
		GROUP BY department
		ORDER BY payroll_cost DESC
	`, tenantID, restaurantID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer departmentRows.Close()

	byDepartment := make([]map[string]interface{}, 0)
	for departmentRows.Next() {
		var (
			department string
			cost       float64
		)
		if err := departmentRows.Scan(&department, &cost); err != nil {
			return nil, err
		}
		byDepartment = append(byDepartment, map[string]interface{}{
			"department":   department,
			"payroll_cost": cost,
		})
	}
	if err := departmentRows.Err(); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"summary": map[string]interface{}{
			"revenue":                   revenue,
			"payroll_cost":              payrollCost,
			"overtime_cost":             overtimeCost,
			"estimated_insurance_cost":  insuranceCost,
			"total_hr_cost":             totalHRCost,
			"labor_cost_pct_of_revenue": laborPct,
		},
		"by_department": byDepartment,
		"generated_at":  time.Now().UTC(),
	}, nil
}

func (r *AccountingAdvancedRepository) CreateEInvoiceSubmission(ctx context.Context, item *domain.EInvoiceSubmission) (int, error) {
	if item == nil {
		return 0, errors.New("e-invoice payload is required")
	}
	if item.ResponsePayload == nil {
		item.ResponsePayload = []byte("{}")
	}
	var id int
	err := database.WithTenantTx(ctx, r.db, int64(item.TenantID), int64(item.RestaurantID), func(tx *sql.Tx) error {
		invoiceTable := ""
		switch strings.ToLower(strings.TrimSpace(item.InvoiceType)) {
		case "customer":
			invoiceTable = "invoices_customer"
		case "supplier":
			invoiceTable = "invoices_supplier"
		default:
			return fmt.Errorf("unsupported e-invoice type %q", item.InvoiceType)
		}
		var invoiceExists bool
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS (
			SELECT 1 FROM `+invoiceTable+` WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
		)`, item.InvoiceID, item.TenantID, item.RestaurantID).Scan(&invoiceExists); err != nil {
			return err
		}
		if !invoiceExists {
			return sql.ErrNoRows
		}
		return tx.QueryRowContext(ctx, `
			INSERT INTO e_invoice_submissions (
				tenant_id, restaurant_id, invoice_type, invoice_id, eta_uuid,
				status, submitted_at, response_payload, error_message, created_by
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, COALESCE($7, CURRENT_TIMESTAMP), $8, $9, $10
			)
			RETURNING id
		`,
			item.TenantID,
			item.RestaurantID,
			item.InvoiceType,
			item.InvoiceID,
			nullString(item.ETAUUID),
			defaultString(item.Status, "pending"),
			nullTime(item.SubmittedAt),
			item.ResponsePayload,
			nullString(item.ErrorMessage),
			item.CreatedBy,
		).Scan(&id)
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountingAdvancedRepository) ListEInvoiceSubmissions(ctx context.Context, tenantID, restaurantID int) ([]domain.EInvoiceSubmission, error) {
	out := make([]domain.EInvoiceSubmission, 0)
	err := database.WithTenantTx(ctx, r.db, int64(tenantID), int64(restaurantID), func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, `
			SELECT
				id, tenant_id, restaurant_id, invoice_type, invoice_id, COALESCE(eta_uuid, ''),
				status, submitted_at, response_payload, COALESCE(error_message, ''), created_by
			FROM e_invoice_submissions
			WHERE tenant_id = $1 AND restaurant_id = $2
			ORDER BY submitted_at DESC, id DESC
		`, tenantID, restaurantID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var (
				item      domain.EInvoiceSubmission
				createdBy sql.NullInt64
			)
			if err := rows.Scan(
				&item.ID,
				&item.TenantID,
				&item.RestaurantID,
				&item.InvoiceType,
				&item.InvoiceID,
				&item.ETAUUID,
				&item.Status,
				&item.SubmittedAt,
				&item.ResponsePayload,
				&item.ErrorMessage,
				&createdBy,
			); err != nil {
				return err
			}
			if createdBy.Valid {
				val := int(createdBy.Int64)
				item.CreatedBy = &val
			}
			out = append(out, item)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AccountingAdvancedRepository) GetSocialInsuranceReportRows(tenantID, restaurantID, year, month int) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(`
		SELECT
			e.employee_code,
			TRIM(CONCAT(e.first_name, ' ', e.last_name)) AS employee_name,
			COALESCE(e.national_id, '') AS national_id,
			COALESCE(s.base_salary, 0) AS base_salary,
			COALESCE(s.social_insurance, 0) AS social_insurance_deduction,
			COALESCE(s.health_insurance, 0) AS health_insurance_deduction,
			COALESCE(it.policy_number, '') AS policy_number,
			COALESCE(it.provider_name, '') AS provider_name
		FROM salaries s
		JOIN employees e
		  ON e.id = s.employee_id
		 AND e.tenant_id = s.tenant_id
		 AND e.restaurant_id = s.restaurant_id
		LEFT JOIN insurance_tracking it
		  ON it.employee_id = e.id
		 AND it.tenant_id = e.tenant_id
		 AND it.restaurant_id = e.restaurant_id
		 AND it.status = 'active'
		WHERE s.tenant_id = $1
		  AND s.restaurant_id = $2
		  AND s.year = $3
		  AND s.month = $4
		ORDER BY e.employee_code ASC
	`, tenantID, restaurantID, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]map[string]interface{}, 0)
	for rows.Next() {
		var (
			employeeCode string
			employeeName string
			nationalID   string
			baseSalary   float64
			social       float64
			health       float64
			policyNumber string
			providerName string
		)
		if err := rows.Scan(
			&employeeCode,
			&employeeName,
			&nationalID,
			&baseSalary,
			&social,
			&health,
			&policyNumber,
			&providerName,
		); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"employee_code":              employeeCode,
			"employee_name":              employeeName,
			"national_id":                nationalID,
			"base_salary":                baseSalary,
			"social_insurance_deduction": social,
			"health_insurance_deduction": health,
			"policy_number":              policyNumber,
			"provider_name":              providerName,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AccountingAdvancedRepository) GetAuditExportRows(
	ctx context.Context,
	tenantID, restaurantID int,
	dateFrom, dateTo time.Time,
	exportType string,
) ([]map[string]interface{}, error) {
	normalized := strings.ToLower(strings.TrimSpace(exportType))
	if normalized == "" {
		normalized = "all"
	}

	rows := make([]map[string]interface{}, 0)
	addRow := func(entityType string, ts time.Time, reference, description string, amount float64, actor string) {
		rows = append(rows, map[string]interface{}{
			"entity_type": entityType,
			"timestamp":   ts.UTC().Format(time.RFC3339),
			"reference":   reference,
			"description": description,
			"amount":      amount,
			"actor":       actor,
		})
	}

	if normalized == "all" || normalized == "accounting" {
		entries, err := r.db.Query(`
			SELECT
				je.entry_date,
				je.entry_number,
				COALESCE(je.description, ''),
				COALESCE(je.total_debit, 0),
				COALESCE(u.name, u.email, CONCAT('user#', je.created_by::text), 'system')
			FROM journal_entries je
			LEFT JOIN users u ON u.id = je.created_by
			WHERE je.tenant_id = $1
			  AND je.restaurant_id = $2
			  AND je.entry_date >= $3::date
			  AND je.entry_date < $4::date
			ORDER BY je.entry_date ASC, je.id ASC
		`, tenantID, restaurantID, dateFrom, dateTo)
		if err != nil {
			return nil, err
		}
		for entries.Next() {
			var ts time.Time
			var reference, description, actor string
			var amount float64
			if err := entries.Scan(&ts, &reference, &description, &amount, &actor); err != nil {
				entries.Close()
				return nil, err
			}
			addRow("journal_entry", ts, reference, description, amount, actor)
		}
		if err := entries.Err(); err != nil {
			entries.Close()
			return nil, err
		}
		entries.Close()
	}

	if normalized == "all" || normalized == "rbac" {
		err := database.WithTenantTx(ctx, r.db, int64(tenantID), int64(restaurantID), func(tx *sql.Tx) error {
			auditRows, err := tx.QueryContext(ctx, `
				SELECT
					pal.created_at,
					COALESCE(pal.action, ''),
					COALESCE(pal.details::text, ''),
					0::numeric,
					COALESCE(u.name, u.email, CONCAT('user#', pal.admin_id::text), 'system')
				FROM permission_audit_logs pal
				LEFT JOIN users u ON u.id = pal.admin_id
				WHERE pal.tenant_id = $1
				  AND pal.created_at >= $2
				  AND pal.created_at < $3
				ORDER BY pal.created_at ASC, pal.id ASC
			`, tenantID, dateFrom, dateTo)
			if err != nil {
				return err
			}
			defer auditRows.Close()
			for auditRows.Next() {
				var ts time.Time
				var reference, description, actor string
				var amount float64
				if err := auditRows.Scan(&ts, &reference, &description, &amount, &actor); err != nil {
					return err
				}
				addRow("rbac_audit", ts, reference, description, amount, actor)
			}
			return auditRows.Err()
		})
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i]["timestamp"].(string) < rows[j]["timestamp"].(string)
	})
	return rows, nil
}

func (r *AccountingAdvancedRepository) LogAuditExport(
	ctx context.Context,
	tenantID, restaurantID, generatedBy int,
	exportType, exportFormat string,
	params map[string]interface{},
	rowCount int,
) error {
	payload, err := json.Marshal(params)
	if err != nil {
		payload = []byte("{}")
	}
	return database.WithTenantTx(ctx, r.db, int64(tenantID), int64(restaurantID), func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO audit_export_logs (
				tenant_id, restaurant_id, export_type, export_format,
				params_json, row_count, generated_by
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7
			)
		`, tenantID, restaurantID, exportType, exportFormat, payload, rowCount, generatedBy)
		return err
	})
}

func (r *AccountingAdvancedRepository) GetChartOfAccountsExcelRows(tenantID, restaurantID int) ([][]string, error) {
	rows, err := r.db.Query(`
		SELECT
			account_code,
			account_name,
			COALESCE(account_name_ar, ''),
			account_type,
			normal_balance,
			COALESCE(currency, 'EGP'),
			COALESCE(account_subtype, ''),
			COALESCE(description, ''),
			is_active
		FROM chart_of_accounts
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY account_code ASC
	`, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := [][]string{{
		"account_code",
		"account_name",
		"account_name_ar",
		"account_type",
		"normal_balance",
		"currency",
		"account_subtype",
		"description",
		"is_active",
	}}
	for rows.Next() {
		var (
			accountCode    string
			accountName    string
			accountNameAr  string
			accountType    string
			normalBalance  string
			currency       string
			accountSubtype string
			description    string
			isActive       bool
		)
		if err := rows.Scan(
			&accountCode,
			&accountName,
			&accountNameAr,
			&accountType,
			&normalBalance,
			&currency,
			&accountSubtype,
			&description,
			&isActive,
		); err != nil {
			return nil, err
		}
		out = append(out, []string{
			accountCode,
			accountName,
			accountNameAr,
			accountType,
			normalBalance,
			currency,
			accountSubtype,
			description,
			fmt.Sprintf("%t", isActive),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AccountingAdvancedRepository) GetJournalEntriesExcelRows(tenantID, restaurantID int, dateFrom, dateTo time.Time) ([][]string, error) {
	rows, err := r.db.Query(`
		SELECT
			je.entry_number,
			je.entry_date::date,
			COALESCE(je.description, ''),
			COALESCE(je.entry_type, ''),
			je.status,
			coa.account_code,
			COALESCE(jel.description, ''),
			jel.debit_amount,
			jel.credit_amount
		FROM journal_entries je
		JOIN journal_entry_lines jel ON jel.journal_entry_id = je.id
		JOIN chart_of_accounts coa
		  ON coa.id = jel.account_id
		 AND coa.tenant_id = je.tenant_id
		 AND coa.restaurant_id = je.restaurant_id
		WHERE je.tenant_id = $1
		  AND je.restaurant_id = $2
		  AND je.entry_date >= $3::date
		  AND je.entry_date < $4::date
		ORDER BY je.entry_date ASC, je.id ASC, jel.line_order ASC, jel.id ASC
	`, tenantID, restaurantID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := [][]string{{
		"entry_number",
		"entry_date",
		"description",
		"entry_type",
		"status",
		"account_code",
		"line_description",
		"debit_amount",
		"credit_amount",
	}}
	for rows.Next() {
		var (
			entryNumber     string
			entryDate       time.Time
			description     string
			entryType       string
			status          string
			accountCode     string
			lineDescription string
			debitAmount     float64
			creditAmount    float64
		)
		if err := rows.Scan(
			&entryNumber,
			&entryDate,
			&description,
			&entryType,
			&status,
			&accountCode,
			&lineDescription,
			&debitAmount,
			&creditAmount,
		); err != nil {
			return nil, err
		}
		out = append(out, []string{
			entryNumber,
			entryDate.Format("2006-01-02"),
			description,
			entryType,
			status,
			accountCode,
			lineDescription,
			fmt.Sprintf("%.2f", debitAmount),
			fmt.Sprintf("%.2f", creditAmount),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AccountingAdvancedRepository) GetAccountIDByCode(tenantID, restaurantID int, accountCode string) (int, error) {
	var accountID int
	err := r.db.QueryRow(`
		SELECT id
		FROM chart_of_accounts
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND account_code = $3
	`, tenantID, restaurantID, strings.TrimSpace(accountCode)).Scan(&accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("account code %s not found", accountCode)
		}
		return 0, err
	}
	return accountID, nil
}

func (r *AccountingAdvancedRepository) UpsertAccountFromImport(
	tenantID, restaurantID int,
	accountCode, accountName, accountNameAr, accountType, normalBalance, currency, accountSubtype, description string,
	isActive bool,
) error {
	accountCode = strings.TrimSpace(accountCode)
	if accountCode == "" {
		return errors.New("account_code is required")
	}
	accountName = strings.TrimSpace(accountName)
	if accountName == "" {
		return errors.New("account_name is required")
	}
	accountType = strings.TrimSpace(accountType)
	if accountType == "" {
		return errors.New("account_type is required")
	}
	normalBalance = strings.TrimSpace(normalBalance)
	if normalBalance == "" {
		normalBalance = domain.NormalBalanceDebit
	}
	currency = defaultString(strings.ToUpper(strings.TrimSpace(currency)), "EGP")

	var existingID int
	err := r.db.QueryRow(`
		SELECT id
		FROM chart_of_accounts
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND account_code = $3
	`, tenantID, restaurantID, accountCode).Scan(&existingID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		_, err = r.db.Exec(`
			INSERT INTO chart_of_accounts (
				tenant_id, restaurant_id, account_code, account_name, account_name_ar,
				description, account_type, account_subtype, normal_balance,
				level, is_header, current_balance, currency, is_active, is_system, is_reconcilable
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9,
				1, FALSE, 0, $10, $11, FALSE, FALSE
			)
		`,
			tenantID,
			restaurantID,
			accountCode,
			accountName,
			nullString(accountNameAr),
			nullString(description),
			accountType,
			nullString(accountSubtype),
			normalBalance,
			currency,
			isActive,
		)
		return err
	}

	_, err = r.db.Exec(`
		UPDATE chart_of_accounts
		SET account_name = $1,
			account_name_ar = $2,
			description = $3,
			account_type = $4,
			account_subtype = $5,
			normal_balance = $6,
			currency = $7,
			is_active = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE tenant_id = $9
		  AND restaurant_id = $10
		  AND id = $11
	`,
		accountName,
		nullString(accountNameAr),
		nullString(description),
		accountType,
		nullString(accountSubtype),
		normalBalance,
		currency,
		isActive,
		tenantID,
		restaurantID,
		existingID,
	)
	return err
}

func (r *AccountingAdvancedRepository) CreateImportedJournalEntry(
	tenantID, restaurantID, userID int,
	entryDate time.Time,
	description, entryType string,
	lines []domain.JournalEntryLineInput,
	autoPost bool,
) (int, error) {
	if len(lines) == 0 {
		return 0, errors.New("journal entry lines are required")
	}
	totalDebit := 0.0
	totalCredit := 0.0
	entryLines := make([]domain.JournalEntryLine, 0, len(lines))
	for idx, line := range lines {
		totalDebit += line.DebitAmount
		totalCredit += line.CreditAmount
		entryLines = append(entryLines, domain.JournalEntryLine{
			AccountID:     line.AccountID,
			DebitAmount:   line.DebitAmount,
			CreditAmount:  line.CreditAmount,
			Description:   line.Description,
			DescriptionAr: line.DescriptionAr,
			LineOrder:     idx + 1,
			TaxCode:       line.TaxCode,
			TaxAmount:     line.TaxAmount,
		})
	}
	if math.Abs(totalDebit-totalCredit) > 0.01 {
		return 0, errors.New("journal entry is not balanced")
	}
	fiscalPeriod := fmt.Sprintf("%04d-%02d", entryDate.Year(), entryDate.Month())
	entry := &domain.JournalEntry{
		TenantID:     tenantID,
		RestaurantID: restaurantID,
		EntryDate:    entryDate,
		FiscalPeriod: fiscalPeriod,
		Description:  strings.TrimSpace(description),
		EntryType:    defaultString(strings.TrimSpace(entryType), "manual"),
		Status:       domain.JournalStatusDraft,
		CreatedBy:    &userID,
		UpdatedBy:    &userID,
		Lines:        entryLines,
	}
	entryID, err := r.journalRepo.CreateJournalEntry(entry)
	if err != nil {
		return 0, err
	}
	if autoPost {
		if err := r.journalRepo.PostEntry(tenantID, restaurantID, entryID, userID); err != nil {
			return 0, err
		}
	}
	return entryID, nil
}

type bankReconciliationSessionScanner interface {
	Scan(dest ...interface{}) error
}

func scanBankReconciliationSession(scanner bankReconciliationSessionScanner) (domain.BankReconciliationSession, error) {
	var (
		item        domain.BankReconciliationSession
		completedAt sql.NullTime
		completedBy sql.NullInt64
		createdBy   sql.NullInt64
		updatedBy   sql.NullInt64
	)
	err := scanner.Scan(
		&item.ID,
		&item.TenantID,
		&item.RestaurantID,
		&item.BankAccountID,
		&item.StatementDate,
		&item.StatementBalance,
		&item.BookBalance,
		&item.Difference,
		&item.Status,
		&item.Notes,
		&completedAt,
		&completedBy,
		&item.CreatedAt,
		&item.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		return domain.BankReconciliationSession{}, err
	}
	if completedAt.Valid {
		item.CompletedAt = &completedAt.Time
	}
	if completedBy.Valid {
		value := int(completedBy.Int64)
		item.CompletedBy = &value
	}
	if createdBy.Valid {
		value := int(createdBy.Int64)
		item.CreatedBy = &value
	}
	if updatedBy.Valid {
		value := int(updatedBy.Int64)
		item.UpdatedBy = &value
	}
	return item, nil
}

type recurringTemplateScanner interface {
	Scan(dest ...interface{}) error
}

func scanRecurringTemplate(scanner recurringTemplateScanner) (domain.RecurringJournalTemplate, error) {
	var (
		item      domain.RecurringJournalTemplate
		lastRun   sql.NullTime
		createdBy sql.NullInt64
		updatedBy sql.NullInt64
	)
	err := scanner.Scan(
		&item.ID,
		&item.TenantID,
		&item.RestaurantID,
		&item.TemplateCode,
		&item.TemplateName,
		&item.Description,
		&item.Frequency,
		&item.StartDate,
		&item.NextRunDate,
		&lastRun,
		&item.IsActive,
		&item.AutoPost,
		&item.LinesJSON,
		&item.CreatedAt,
		&item.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		return domain.RecurringJournalTemplate{}, err
	}
	if lastRun.Valid {
		item.LastRunDate = &lastRun.Time
	}
	if createdBy.Valid {
		val := int(createdBy.Int64)
		item.CreatedBy = &val
	}
	if updatedBy.Valid {
		val := int(updatedBy.Int64)
		item.UpdatedBy = &val
	}
	return item, nil
}

type supplierPaymentScheduleScanner interface {
	Scan(dest ...interface{}) error
}

func scanSupplierPaymentSchedule(scanner supplierPaymentScheduleScanner) (domain.SupplierPaymentSchedule, error) {
	var (
		item              domain.SupplierPaymentSchedule
		supplierInvoiceID sql.NullInt64
		reminderSentAt    sql.NullTime
		createdBy         sql.NullInt64
		updatedBy         sql.NullInt64
	)
	err := scanner.Scan(
		&item.ID,
		&item.TenantID,
		&item.RestaurantID,
		&supplierInvoiceID,
		&item.SupplierID,
		&item.SupplierName,
		&item.ScheduledDate,
		&item.Amount,
		&item.Currency,
		&item.ReminderDays,
		&item.Status,
		&reminderSentAt,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		return domain.SupplierPaymentSchedule{}, err
	}
	if supplierInvoiceID.Valid {
		val := int(supplierInvoiceID.Int64)
		item.SupplierInvoiceID = &val
	}
	if reminderSentAt.Valid {
		item.ReminderSentAt = &reminderSentAt.Time
	}
	if createdBy.Valid {
		val := int(createdBy.Int64)
		item.CreatedBy = &val
	}
	if updatedBy.Valid {
		val := int(updatedBy.Int64)
		item.UpdatedBy = &val
	}
	return item, nil
}

func refreshReconciliationSummaryTx(tx *sql.Tx, tenantID, restaurantID, reconciliationID int) error {
	_, err := tx.Exec(`
		WITH totals AS (
			SELECT
				COALESCE(SUM(book_amount), 0) AS book_total
			FROM bank_reconciliation_lines
			WHERE reconciliation_id = $3
			  AND tenant_id = $1
		)
		UPDATE bank_reconciliation_sessions s
		SET book_balance = totals.book_total,
			difference = s.statement_balance - totals.book_total,
			status = CASE
				WHEN ABS(s.statement_balance - totals.book_total) < 0.01 THEN 'completed'
				ELSE 'in_progress'
			END,
			updated_at = CURRENT_TIMESTAMP
		FROM totals
		WHERE s.tenant_id = $1
		  AND s.restaurant_id = $2
		  AND s.id = $3
	`, tenantID, restaurantID, reconciliationID)
	return err
}

func nextRecurringRunDate(frequency string, from time.Time) time.Time {
	switch strings.ToLower(strings.TrimSpace(frequency)) {
	case "weekly":
		return from.AddDate(0, 0, 7)
	case "quarterly":
		return from.AddDate(0, 3, 0)
	case "yearly":
		return from.AddDate(1, 0, 0)
	default:
		return from.AddDate(0, 1, 0)
	}
}

func beginningOfMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func averageFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func defaultString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func defaultInt(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

func nullTime(value time.Time) interface{} {
	if value.IsZero() {
		return nil
	}
	return value
}
