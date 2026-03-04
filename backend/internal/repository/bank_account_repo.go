package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"pos-saas/internal/domain"
)

type BankAccountRepository struct {
	db *sql.DB
}

func NewBankAccountRepository(db *sql.DB) *BankAccountRepository {
	return &BankAccountRepository{db: db}
}

func (r *BankAccountRepository) ListBankAccountsPaged(ctx context.Context, tenantID, restaurantID, page, pageSize int) ([]domain.BankAccount, int, error) {
	var total int
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(1) FROM bank_accounts_ledger WHERE tenant_id = $1 AND restaurant_id = $2`,
		tenantID,
		restaurantID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id, tenant_id, restaurant_id, account_name, COALESCE(account_name_ar, ''),
			bank_name, COALESCE(bank_branch, ''), account_number, COALESCE(iban, ''),
			COALESCE(swift_code, ''), currency, account_type, opening_balance,
			current_balance, opening_date, gl_account_id, is_active, is_primary,
			last_reconciled_date, last_reconciled_balance, created_at, updated_at,
			created_by, updated_by
		FROM bank_accounts_ledger
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY is_primary DESC, account_name ASC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	accounts := make([]domain.BankAccount, 0, pageSize)
	for rows.Next() {
		account, err := scanBankAccount(rows)
		if err != nil {
			return nil, 0, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (r *BankAccountRepository) ListBankAccounts(tenantID, restaurantID int) ([]domain.BankAccount, error) {
	accounts, _, err := r.ListBankAccountsPaged(context.Background(), tenantID, restaurantID, 1, 1000000)
	return accounts, err
}

func (r *BankAccountRepository) GetBankAccountByID(tenantID, restaurantID, id int) (*domain.BankAccount, error) {
	account, err := scanBankAccount(r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, account_name, COALESCE(account_name_ar, ''),
			bank_name, COALESCE(bank_branch, ''), account_number, COALESCE(iban, ''),
			COALESCE(swift_code, ''), currency, account_type, opening_balance,
			current_balance, opening_date, gl_account_id, is_active, is_primary,
			last_reconciled_date, last_reconciled_balance, created_at, updated_at,
			created_by, updated_by
		FROM bank_accounts_ledger
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *BankAccountRepository) CreateBankAccount(account *domain.BankAccount) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO bank_accounts_ledger (
			tenant_id, restaurant_id, account_name, account_name_ar, bank_name,
			bank_branch, account_number, iban, swift_code, currency, account_type,
			opening_balance, current_balance, opening_date, gl_account_id,
			is_active, is_primary, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15,
			$16, $17, $18, $19
		)
		RETURNING id
	`,
		account.TenantID,
		account.RestaurantID,
		account.AccountName,
		nullString(account.AccountNameAr),
		account.BankName,
		nullString(account.BankBranch),
		account.AccountNumber,
		nullString(account.IBAN),
		nullString(account.SwiftCode),
		account.Currency,
		account.AccountType,
		account.OpeningBalance,
		account.CurrentBalance,
		account.OpeningDate,
		account.GLAccountID,
		account.IsActive,
		account.IsPrimary,
		account.CreatedBy,
		account.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *BankAccountRepository) UpdateBankAccount(account *domain.BankAccount) error {
	result, err := r.db.Exec(`
		UPDATE bank_accounts_ledger
		SET account_name = $1,
			account_name_ar = $2,
			bank_name = $3,
			bank_branch = $4,
			account_number = $5,
			iban = $6,
			swift_code = $7,
			currency = $8,
			account_type = $9,
			opening_balance = $10,
			opening_date = $11,
			gl_account_id = $12,
			is_active = $13,
			is_primary = $14,
			updated_by = $15,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $16 AND tenant_id = $17 AND restaurant_id = $18
	`,
		account.AccountName,
		nullString(account.AccountNameAr),
		account.BankName,
		nullString(account.BankBranch),
		account.AccountNumber,
		nullString(account.IBAN),
		nullString(account.SwiftCode),
		account.Currency,
		account.AccountType,
		account.OpeningBalance,
		account.OpeningDate,
		account.GLAccountID,
		account.IsActive,
		account.IsPrimary,
		account.UpdatedBy,
		account.ID,
		account.TenantID,
		account.RestaurantID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("bank account not found")
	}
	return nil
}

func (r *BankAccountRepository) DeleteBankAccount(tenantID, restaurantID, id int) error {
	result, err := r.db.Exec(`
		DELETE FROM bank_accounts_ledger
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
		return errors.New("bank account not found")
	}
	return nil
}

func (r *BankAccountRepository) GetBankTransactions(tenantID, restaurantID, id int) ([]domain.BankTransaction, error) {
	account, err := r.GetBankAccountByID(tenantID, restaurantID, id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("bank account not found")
	}
	if account.GLAccountID == nil {
		return []domain.BankTransaction{}, nil
	}

	rows, err := r.db.Query(`
		SELECT
			gl.entry_date,
			je.entry_number,
			COALESCE(je.reference_number, ''),
			COALESCE(jel.description, je.description, ''),
			gl.debit_amount,
			gl.credit_amount,
			gl.running_balance,
			gl.journal_entry_id
		FROM general_ledger gl
		JOIN journal_entries je ON je.id = gl.journal_entry_id
		LEFT JOIN journal_entry_lines jel ON jel.id = gl.journal_line_id
		WHERE gl.tenant_id = $1
		  AND gl.restaurant_id = $2
		  AND gl.account_id = $3
		ORDER BY gl.entry_date DESC, gl.id DESC
	`, tenantID, restaurantID, *account.GLAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]domain.BankTransaction, 0)
	for rows.Next() {
		var item domain.BankTransaction
		if err := rows.Scan(
			&item.TransactionDate,
			&item.EntryNumber,
			&item.ReferenceNumber,
			&item.Description,
			&item.DebitAmount,
			&item.CreditAmount,
			&item.RunningBalance,
			&item.JournalEntryID,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *BankAccountRepository) AdjustBalanceTx(tx *sql.Tx, tenantID, restaurantID, id int, delta float64) error {
	if tx == nil {
		return errors.New("transaction is required")
	}
	if id <= 0 || delta == 0 {
		return nil
	}

	result, err := tx.Exec(`
		UPDATE bank_accounts_ledger
		SET current_balance = current_balance + $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE tenant_id = $2 AND restaurant_id = $3 AND id = $4
	`, delta, tenantID, restaurantID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("bank account %d not found", id)
	}
	return nil
}

type bankAccountScanner interface {
	Scan(dest ...interface{}) error
}

func scanBankAccount(scanner bankAccountScanner) (domain.BankAccount, error) {
	var (
		account               domain.BankAccount
		glAccountID           sql.NullInt64
		lastReconciledDate    sql.NullTime
		lastReconciledBalance sql.NullFloat64
		createdBy             sql.NullInt64
		updatedBy             sql.NullInt64
	)

	err := scanner.Scan(
		&account.ID,
		&account.TenantID,
		&account.RestaurantID,
		&account.AccountName,
		&account.AccountNameAr,
		&account.BankName,
		&account.BankBranch,
		&account.AccountNumber,
		&account.IBAN,
		&account.SwiftCode,
		&account.Currency,
		&account.AccountType,
		&account.OpeningBalance,
		&account.CurrentBalance,
		&account.OpeningDate,
		&glAccountID,
		&account.IsActive,
		&account.IsPrimary,
		&lastReconciledDate,
		&lastReconciledBalance,
		&account.CreatedAt,
		&account.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		return domain.BankAccount{}, err
	}

	if glAccountID.Valid {
		value := int(glAccountID.Int64)
		account.GLAccountID = &value
	}
	if lastReconciledDate.Valid {
		account.LastReconciledDate = &lastReconciledDate.Time
	}
	if lastReconciledBalance.Valid {
		value := lastReconciledBalance.Float64
		account.LastReconciledBalance = &value
	}
	if createdBy.Valid {
		value := int(createdBy.Int64)
		account.CreatedBy = &value
	}
	if updatedBy.Valid {
		value := int(updatedBy.Int64)
		account.UpdatedBy = &value
	}

	return account, nil
}

func parseAccountingDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}
