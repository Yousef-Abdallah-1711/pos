package repository

import (
	"database/sql"
	"errors"
	"sort"

	"pos-saas/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) ListAccounts(tenantID, restaurantID int, includeInactive bool) ([]domain.Account, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, account_code, account_name,
			COALESCE(account_name_ar, ''), COALESCE(description, ''),
			COALESCE(description_ar, ''), account_type,
			COALESCE(account_subtype, ''), normal_balance, parent_id,
			level, is_header, current_balance, currency, is_active,
			is_system, is_reconcilable, tax_rate, created_at, updated_at,
			created_by, updated_by
		FROM chart_of_accounts
		WHERE tenant_id = $1 AND restaurant_id = $2
	`
	if !includeInactive {
		query += ` AND is_active = TRUE`
	}
	query += ` ORDER BY account_code ASC`

	rows, err := r.db.Query(query, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]domain.Account, 0)
	for rows.Next() {
		account, err := scanAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepository) GetAccountByID(tenantID, restaurantID, id int) (*domain.Account, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, account_code, account_name,
			COALESCE(account_name_ar, ''), COALESCE(description, ''),
			COALESCE(description_ar, ''), account_type,
			COALESCE(account_subtype, ''), normal_balance, parent_id,
			level, is_header, current_balance, currency, is_active,
			is_system, is_reconcilable, tax_rate, created_at, updated_at,
			created_by, updated_by
		FROM chart_of_accounts
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`

	account, err := scanAccount(r.db.QueryRow(query, tenantID, restaurantID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) GetAccountByCode(tenantID, restaurantID int, code string) (*domain.Account, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, account_code, account_name,
			COALESCE(account_name_ar, ''), COALESCE(description, ''),
			COALESCE(description_ar, ''), account_type,
			COALESCE(account_subtype, ''), normal_balance, parent_id,
			level, is_header, current_balance, currency, is_active,
			is_system, is_reconcilable, tax_rate, created_at, updated_at,
			created_by, updated_by
		FROM chart_of_accounts
		WHERE tenant_id = $1 AND restaurant_id = $2 AND account_code = $3
	`

	account, err := scanAccount(r.db.QueryRow(query, tenantID, restaurantID, code))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) CreateAccount(account *domain.Account) (int, error) {
	query := `
		INSERT INTO chart_of_accounts (
			tenant_id, restaurant_id, account_code, account_name, account_name_ar,
			description, description_ar, account_type, account_subtype, normal_balance,
			parent_id, level, is_header, current_balance, currency, is_active,
			is_system, is_reconcilable, tax_rate, created_by, updated_by
		)
		SELECT
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21
		WHERE EXISTS (
			SELECT 1 FROM restaurants r WHERE r.id = $2 AND r.tenant_id = $1
		) AND ($11::integer IS NULL OR EXISTS (
			SELECT 1 FROM chart_of_accounts parent
			WHERE parent.id = $11 AND parent.tenant_id = $1 AND parent.restaurant_id = $2
		))
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(
		query,
		account.TenantID,
		account.RestaurantID,
		account.AccountCode,
		account.AccountName,
		account.AccountNameAr,
		account.Description,
		account.DescriptionAr,
		account.AccountType,
		account.AccountSubtype,
		account.NormalBalance,
		account.ParentID,
		account.Level,
		account.IsHeader,
		account.CurrentBalance,
		account.Currency,
		account.IsActive,
		account.IsSystem,
		account.IsReconcilable,
		account.TaxRate,
		account.CreatedBy,
		account.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AccountRepository) UpdateAccount(account *domain.Account) error {
	result, err := r.db.Exec(`
		UPDATE chart_of_accounts
		SET account_code = $1,
			account_name = $2,
			account_name_ar = $3,
			description = $4,
			description_ar = $5,
			account_type = $6,
			account_subtype = $7,
			normal_balance = $8,
			parent_id = $9,
			level = $10,
			is_header = $11,
			currency = $12,
			is_active = $13,
			is_reconcilable = $14,
			tax_rate = $15,
			updated_by = $16,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $17 AND tenant_id = $18 AND restaurant_id = $19
		  AND ($9 IS NULL OR EXISTS (
			SELECT 1 FROM chart_of_accounts parent
			WHERE parent.id = $9 AND parent.tenant_id = $18 AND parent.restaurant_id = $19
		  ))
	`,
		account.AccountCode,
		account.AccountName,
		account.AccountNameAr,
		account.Description,
		account.DescriptionAr,
		account.AccountType,
		account.AccountSubtype,
		account.NormalBalance,
		account.ParentID,
		account.Level,
		account.IsHeader,
		account.Currency,
		account.IsActive,
		account.IsReconcilable,
		account.TaxRate,
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
		return errors.New("account not found")
	}
	return nil
}

func (r *AccountRepository) DeleteAccount(tenantID, restaurantID, id int) error {
	result, err := r.db.Exec(
		`DELETE FROM chart_of_accounts WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3`,
		tenantID,
		restaurantID,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("account not found")
	}

	return nil
}

func (r *AccountRepository) GetAccountTree(tenantID, restaurantID int, includeInactive bool) ([]domain.Account, error) {
	accounts, err := r.ListAccounts(tenantID, restaurantID, includeInactive)
	if err != nil {
		return nil, err
	}

	childrenByParent := make(map[int][]domain.Account)
	roots := make([]domain.Account, 0)
	for _, account := range accounts {
		if account.ParentID == nil {
			roots = append(roots, account)
			continue
		}
		childrenByParent[*account.ParentID] = append(childrenByParent[*account.ParentID], account)
	}

	var attachChildren func(account domain.Account) domain.Account
	attachChildren = func(account domain.Account) domain.Account {
		children := childrenByParent[account.ID]
		sort.Slice(children, func(i, j int) bool {
			return children[i].AccountCode < children[j].AccountCode
		})
		for _, child := range children {
			account.Children = append(account.Children, attachChildren(child))
		}
		return account
	}

	sort.Slice(roots, func(i, j int) bool {
		return roots[i].AccountCode < roots[j].AccountCode
	})

	tree := make([]domain.Account, 0, len(roots))
	for _, root := range roots {
		tree = append(tree, attachChildren(root))
	}
	return tree, nil
}

func (r *AccountRepository) HasChildAccounts(tenantID, restaurantID, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chart_of_accounts
			WHERE tenant_id = $1 AND restaurant_id = $2 AND parent_id = $3
		)
	`, tenantID, restaurantID, id).Scan(&exists)
	return exists, err
}

func (r *AccountRepository) HasJournalEntries(tenantID, restaurantID, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM journal_entry_lines jel
			JOIN journal_entries je ON je.id = jel.journal_entry_id
			WHERE je.tenant_id = $1
			  AND je.restaurant_id = $2
			  AND jel.account_id = $3
		)
	`, tenantID, restaurantID, id).Scan(&exists)
	return exists, err
}

type accountScanner interface {
	Scan(dest ...interface{}) error
}

func scanAccount(scanner accountScanner) (domain.Account, error) {
	var (
		account   domain.Account
		parentID  sql.NullInt64
		createdBy sql.NullInt64
		updatedBy sql.NullInt64
		taxRate   sql.NullFloat64
	)

	err := scanner.Scan(
		&account.ID,
		&account.TenantID,
		&account.RestaurantID,
		&account.AccountCode,
		&account.AccountName,
		&account.AccountNameAr,
		&account.Description,
		&account.DescriptionAr,
		&account.AccountType,
		&account.AccountSubtype,
		&account.NormalBalance,
		&parentID,
		&account.Level,
		&account.IsHeader,
		&account.CurrentBalance,
		&account.Currency,
		&account.IsActive,
		&account.IsSystem,
		&account.IsReconcilable,
		&taxRate,
		&account.CreatedAt,
		&account.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		return domain.Account{}, err
	}

	if parentID.Valid {
		value := int(parentID.Int64)
		account.ParentID = &value
	}
	if createdBy.Valid {
		value := int(createdBy.Int64)
		account.CreatedBy = &value
	}
	if updatedBy.Valid {
		value := int(updatedBy.Int64)
		account.UpdatedBy = &value
	}
	if taxRate.Valid {
		value := taxRate.Float64
		account.TaxRate = &value
	}

	return account, nil
}
