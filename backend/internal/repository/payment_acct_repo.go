package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"pos-saas/internal/domain"
)

type PaymentAccountingRepository struct {
	db          *sql.DB
	journalRepo *JournalEntryRepository
	accountRepo *AccountRepository
	bankRepo    *BankAccountRepository
}

func NewPaymentAccountingRepository(db *sql.DB, journalRepo *JournalEntryRepository, accountRepo *AccountRepository, bankRepo *BankAccountRepository) *PaymentAccountingRepository {
	return &PaymentAccountingRepository{
		db:          db,
		journalRepo: journalRepo,
		accountRepo: accountRepo,
		bankRepo:    bankRepo,
	}
}

func (r *PaymentAccountingRepository) CreatePaymentReceived(payment *domain.PaymentReceived) (int, error) {
	if payment == nil {
		return 0, errors.New("payment is required")
	}
	if payment.Amount <= 0 {
		return 0, errors.New("payment amount must be greater than zero")
	}
	if payment.PaymentDate.IsZero() {
		payment.PaymentDate = time.Now()
	}
	if strings.TrimSpace(payment.PaymentMethod) == "" {
		return 0, errors.New("payment_method is required")
	}
	if payment.CustomerID <= 0 && payment.InvoiceID == nil {
		return 0, errors.New("customer_id or invoice_id is required")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	customerID, receivableAccountID, invoiceID, invoiceBalanceDue, err := r.resolveCustomerPaymentContextTx(tx, payment)
	if err != nil {
		return 0, err
	}
	if invoiceID != nil && payment.Amount > invoiceBalanceDue+0.009 {
		return 0, errors.New("payment amount exceeds invoice balance")
	}

	collectionAccountID, err := r.resolveCollectionAccountTx(tx, payment.TenantID, payment.RestaurantID, payment.BankAccountID, payment.PaymentMethod)
	if err != nil {
		return 0, err
	}

	number, err := r.generatePaymentNumberTx(tx, "received", payment.TenantID, payment.RestaurantID, payment.PaymentDate)
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRow(`
		INSERT INTO payments_received (
			tenant_id, restaurant_id, payment_number, payment_date, customer_id,
			invoice_id, amount, payment_method, reference, bank_account_id,
			notes, status, created_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13
		)
		RETURNING id
	`,
		payment.TenantID,
		payment.RestaurantID,
		number,
		payment.PaymentDate,
		customerID,
		invoiceID,
		payment.Amount,
		payment.PaymentMethod,
		nullString(payment.Reference),
		payment.BankAccountID,
		nullString(payment.Notes),
		domain.AccountingPaymentStatusCompleted,
		payment.CreatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	entry := domain.JournalEntry{
		TenantID:        payment.TenantID,
		RestaurantID:    payment.RestaurantID,
		EntryDate:       payment.PaymentDate,
		Description:     fmt.Sprintf("Customer payment %s", number),
		ReferenceType:   "payment_received",
		ReferenceID:     &id,
		ReferenceNumber: number,
		EntryType:       "payment_received",
		Status:          domain.JournalStatusDraft,
		CreatedBy:       payment.CreatedBy,
		UpdatedBy:       payment.CreatedBy,
		Lines: []domain.JournalEntryLine{
			{
				AccountID:   collectionAccountID,
				DebitAmount: roundAccountingAmount(payment.Amount),
				Description: fmt.Sprintf("Collection for %s", number),
				LineOrder:   1,
			},
			{
				AccountID:    receivableAccountID,
				CreditAmount: roundAccountingAmount(payment.Amount),
				Description:  fmt.Sprintf("Receivable settlement for %s", number),
				LineOrder:    2,
			},
		},
	}

	journalID, err := r.journalRepo.CreateJournalEntryTx(tx, &entry)
	if err != nil {
		return 0, err
	}
	postedBy := 0
	if payment.CreatedBy != nil {
		postedBy = *payment.CreatedBy
	}
	if err := r.journalRepo.PostEntryTx(tx, payment.TenantID, payment.RestaurantID, journalID, postedBy); err != nil {
		return 0, err
	}

	if _, err := tx.Exec(`
		UPDATE payments_received
		SET journal_entry_id = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, journalID, id, payment.TenantID, payment.RestaurantID); err != nil {
		return 0, err
	}

	if invoiceID != nil {
		if err := r.applyCustomerInvoicePaymentTx(tx, payment.TenantID, payment.RestaurantID, *invoiceID, payment.Amount); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`
		UPDATE customers_accounting
		SET current_balance = current_balance - $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, payment.Amount, customerID, payment.TenantID, payment.RestaurantID); err != nil {
		return 0, err
	}
	if payment.BankAccountID != nil && r.bankRepo != nil {
		if err := r.bankRepo.AdjustBalanceTx(tx, payment.TenantID, payment.RestaurantID, *payment.BankAccountID, payment.Amount); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PaymentAccountingRepository) CreatePaymentMade(payment *domain.PaymentMade) (int, error) {
	if payment == nil {
		return 0, errors.New("payment is required")
	}
	if payment.Amount <= 0 {
		return 0, errors.New("payment amount must be greater than zero")
	}
	if payment.PaymentDate.IsZero() {
		payment.PaymentDate = time.Now()
	}
	if strings.TrimSpace(payment.PaymentMethod) == "" {
		return 0, errors.New("payment_method is required")
	}
	if payment.SupplierID <= 0 && payment.InvoiceID == nil {
		return 0, errors.New("supplier_id or invoice_id is required")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	supplierID, payableAccountID, invoiceID, invoiceBalanceDue, err := r.resolveSupplierPaymentContextTx(tx, payment)
	if err != nil {
		return 0, err
	}
	if invoiceID != nil && payment.Amount > invoiceBalanceDue+0.009 {
		return 0, errors.New("payment amount exceeds invoice balance")
	}

	collectionAccountID, err := r.resolveCollectionAccountTx(tx, payment.TenantID, payment.RestaurantID, payment.BankAccountID, payment.PaymentMethod)
	if err != nil {
		return 0, err
	}

	number, err := r.generatePaymentNumberTx(tx, "made", payment.TenantID, payment.RestaurantID, payment.PaymentDate)
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRow(`
		INSERT INTO payments_made (
			tenant_id, restaurant_id, payment_number, payment_date, supplier_id,
			invoice_id, amount, payment_method, reference, bank_account_id,
			notes, status, created_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13
		)
		RETURNING id
	`,
		payment.TenantID,
		payment.RestaurantID,
		number,
		payment.PaymentDate,
		supplierID,
		invoiceID,
		payment.Amount,
		payment.PaymentMethod,
		nullString(payment.Reference),
		payment.BankAccountID,
		nullString(payment.Notes),
		domain.AccountingPaymentStatusCompleted,
		payment.CreatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	entry := domain.JournalEntry{
		TenantID:        payment.TenantID,
		RestaurantID:    payment.RestaurantID,
		EntryDate:       payment.PaymentDate,
		Description:     fmt.Sprintf("Supplier payment %s", number),
		ReferenceType:   "payment_made",
		ReferenceID:     &id,
		ReferenceNumber: number,
		EntryType:       "payment_made",
		Status:          domain.JournalStatusDraft,
		CreatedBy:       payment.CreatedBy,
		UpdatedBy:       payment.CreatedBy,
		Lines: []domain.JournalEntryLine{
			{
				AccountID:   payableAccountID,
				DebitAmount: roundAccountingAmount(payment.Amount),
				Description: fmt.Sprintf("Accounts payable settlement %s", number),
				LineOrder:   1,
			},
			{
				AccountID:    collectionAccountID,
				CreditAmount: roundAccountingAmount(payment.Amount),
				Description:  fmt.Sprintf("Cash/bank payment %s", number),
				LineOrder:    2,
			},
		},
	}

	journalID, err := r.journalRepo.CreateJournalEntryTx(tx, &entry)
	if err != nil {
		return 0, err
	}
	postedBy := 0
	if payment.CreatedBy != nil {
		postedBy = *payment.CreatedBy
	}
	if err := r.journalRepo.PostEntryTx(tx, payment.TenantID, payment.RestaurantID, journalID, postedBy); err != nil {
		return 0, err
	}

	if _, err := tx.Exec(`
		UPDATE payments_made
		SET journal_entry_id = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, journalID, id, payment.TenantID, payment.RestaurantID); err != nil {
		return 0, err
	}

	if invoiceID != nil {
		if err := r.applySupplierInvoicePaymentTx(tx, payment.TenantID, payment.RestaurantID, *invoiceID, payment.Amount); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`
		UPDATE suppliers
		SET current_balance = current_balance - $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, payment.Amount, supplierID, payment.TenantID, payment.RestaurantID); err != nil {
		return 0, err
	}
	if payment.BankAccountID != nil && r.bankRepo != nil {
		if err := r.bankRepo.AdjustBalanceTx(tx, payment.TenantID, payment.RestaurantID, *payment.BankAccountID, -payment.Amount); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PaymentAccountingRepository) ListPaymentsReceived(tenantID, restaurantID int) ([]domain.PaymentReceived, error) {
	items, _, err := r.ListPaymentsReceivedPaged(tenantID, restaurantID, 1, 100000)
	return items, err
}

func (r *PaymentAccountingRepository) ListPaymentsReceivedPaged(tenantID, restaurantID, page, pageSize int) ([]domain.PaymentReceived, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset := (page - 1) * pageSize

	var total int
	if err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM payments_received pr
		WHERE pr.tenant_id = $1 AND pr.restaurant_id = $2
	`, tenantID, restaurantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT
			pr.id, pr.tenant_id, pr.restaurant_id, pr.payment_number, pr.payment_date,
			pr.customer_id, COALESCE(c.customer_name, ''), pr.invoice_id, pr.amount,
			pr.payment_method, COALESCE(pr.reference, ''), pr.bank_account_id,
			COALESCE(pr.notes, ''), pr.journal_entry_id, pr.status, pr.created_at,
			pr.updated_at, pr.created_by
		FROM payments_received pr
		JOIN customers_accounting c ON c.id = pr.customer_id
		WHERE pr.tenant_id = $1 AND pr.restaurant_id = $2
		ORDER BY pr.payment_date DESC, pr.id DESC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.PaymentReceived, 0)
	for rows.Next() {
		item, err := scanPaymentReceived(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *PaymentAccountingRepository) ListPaymentsMade(tenantID, restaurantID int) ([]domain.PaymentMade, error) {
	items, _, err := r.ListPaymentsMadePaged(tenantID, restaurantID, 1, 100000)
	return items, err
}

func (r *PaymentAccountingRepository) ListPaymentsMadePaged(tenantID, restaurantID, page, pageSize int) ([]domain.PaymentMade, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset := (page - 1) * pageSize

	var total int
	if err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM payments_made pm
		WHERE pm.tenant_id = $1 AND pm.restaurant_id = $2
	`, tenantID, restaurantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT
			pm.id, pm.tenant_id, pm.restaurant_id, pm.payment_number, pm.payment_date,
			pm.supplier_id, COALESCE(s.company_name, ''), pm.invoice_id, pm.amount,
			pm.payment_method, COALESCE(pm.reference, ''), pm.bank_account_id,
			COALESCE(pm.notes, ''), pm.journal_entry_id, pm.status, pm.created_at,
			pm.updated_at, pm.created_by
		FROM payments_made pm
		JOIN suppliers s ON s.id = pm.supplier_id
		WHERE pm.tenant_id = $1 AND pm.restaurant_id = $2
		ORDER BY pm.payment_date DESC, pm.id DESC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.PaymentMade, 0)
	for rows.Next() {
		item, err := scanPaymentMade(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *PaymentAccountingRepository) resolveCustomerPaymentContextTx(tx *sql.Tx, payment *domain.PaymentReceived) (int, int, *int, float64, error) {
	if payment.InvoiceID != nil {
		var (
			customerID          int
			receivableAccountID sql.NullInt64
			balanceDue          float64
		)
		err := tx.QueryRow(`
			SELECT ic.customer_id, c.receivable_account_id, ic.balance_due
			FROM invoices_customer ic
			JOIN customers_accounting c ON c.id = ic.customer_id
			WHERE ic.id = $1 AND ic.tenant_id = $2 AND ic.restaurant_id = $3
		`, *payment.InvoiceID, payment.TenantID, payment.RestaurantID).Scan(&customerID, &receivableAccountID, &balanceDue)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, 0, nil, 0, errors.New("customer invoice not found")
			}
			return 0, 0, nil, 0, err
		}

		accountID := 0
		if receivableAccountID.Valid {
			accountID = int(receivableAccountID.Int64)
		} else {
			account, err := lookupAccountByCode(r.accountRepo, payment.TenantID, payment.RestaurantID, "1120")
			if err != nil {
				return 0, 0, nil, 0, err
			}
			accountID = account.ID
		}
		return customerID, accountID, payment.InvoiceID, balanceDue, nil
	}

	var receivableAccountID sql.NullInt64
	err := tx.QueryRow(`
		SELECT receivable_account_id
		FROM customers_accounting
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, payment.CustomerID, payment.TenantID, payment.RestaurantID).Scan(&receivableAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, nil, 0, errors.New("customer account not found")
		}
		return 0, 0, nil, 0, err
	}
	if receivableAccountID.Valid {
		return payment.CustomerID, int(receivableAccountID.Int64), nil, 0, nil
	}
	account, err := lookupAccountByCode(r.accountRepo, payment.TenantID, payment.RestaurantID, "1120")
	if err != nil {
		return 0, 0, nil, 0, err
	}
	return payment.CustomerID, account.ID, nil, 0, nil
}

func (r *PaymentAccountingRepository) resolveSupplierPaymentContextTx(tx *sql.Tx, payment *domain.PaymentMade) (int, int, *int, float64, error) {
	if payment.InvoiceID != nil {
		var (
			supplierID       int
			payableAccountID sql.NullInt64
			balanceDue       float64
		)
		err := tx.QueryRow(`
			SELECT inv.supplier_id, s.payable_account_id, inv.balance_due
			FROM invoices_supplier inv
			JOIN suppliers s ON s.id = inv.supplier_id
			WHERE inv.id = $1 AND inv.tenant_id = $2 AND inv.restaurant_id = $3
		`, *payment.InvoiceID, payment.TenantID, payment.RestaurantID).Scan(&supplierID, &payableAccountID, &balanceDue)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, 0, nil, 0, errors.New("supplier invoice not found")
			}
			return 0, 0, nil, 0, err
		}

		accountID := 0
		if payableAccountID.Valid {
			accountID = int(payableAccountID.Int64)
		} else {
			account, err := lookupAccountByCode(r.accountRepo, payment.TenantID, payment.RestaurantID, "2101")
			if err != nil {
				return 0, 0, nil, 0, err
			}
			accountID = account.ID
		}
		return supplierID, accountID, payment.InvoiceID, balanceDue, nil
	}

	var payableAccountID sql.NullInt64
	err := tx.QueryRow(`
		SELECT payable_account_id
		FROM suppliers
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, payment.SupplierID, payment.TenantID, payment.RestaurantID).Scan(&payableAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, nil, 0, errors.New("supplier not found")
		}
		return 0, 0, nil, 0, err
	}
	if payableAccountID.Valid {
		return payment.SupplierID, int(payableAccountID.Int64), nil, 0, nil
	}
	account, err := lookupAccountByCode(r.accountRepo, payment.TenantID, payment.RestaurantID, "2101")
	if err != nil {
		return 0, 0, nil, 0, err
	}
	return payment.SupplierID, account.ID, nil, 0, nil
}

func (r *PaymentAccountingRepository) resolveCollectionAccountTx(tx *sql.Tx, tenantID, restaurantID int, bankAccountID *int, paymentMethod string) (int, error) {
	if bankAccountID != nil {
		var glAccountID sql.NullInt64
		err := tx.QueryRow(`
			SELECT gl_account_id
			FROM bank_accounts_ledger
			WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
		`, *bankAccountID, tenantID, restaurantID).Scan(&glAccountID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, errors.New("bank account not found")
			}
			return 0, err
		}
		if glAccountID.Valid {
			return int(glAccountID.Int64), nil
		}
	}

	account, err := lookupAccountByCode(r.accountRepo, tenantID, restaurantID, fallbackCollectionAccountCode(paymentMethod))
	if err != nil {
		return 0, err
	}
	return account.ID, nil
}

func (r *PaymentAccountingRepository) applyCustomerInvoicePaymentTx(tx *sql.Tx, tenantID, restaurantID, invoiceID int, amount float64) error {
	_, err := tx.Exec(`
		UPDATE invoices_customer
		SET amount_paid = amount_paid + $1,
			balance_due = GREATEST(total_amount - (amount_paid + $1), 0),
			status = CASE
				WHEN total_amount - (amount_paid + $1) <= 0 THEN 'paid'
				WHEN amount_paid + $1 > 0 THEN 'partially_paid'
				WHEN due_date < CURRENT_DATE AND status <> 'draft' THEN 'overdue'
				ELSE status
			END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, amount, invoiceID, tenantID, restaurantID)
	return err
}

func (r *PaymentAccountingRepository) applySupplierInvoicePaymentTx(tx *sql.Tx, tenantID, restaurantID, invoiceID int, amount float64) error {
	_, err := tx.Exec(`
		UPDATE invoices_supplier
		SET amount_paid = amount_paid + $1,
			balance_due = GREATEST(total_amount - (amount_paid + $1), 0),
			status = CASE
				WHEN total_amount - (amount_paid + $1) <= 0 THEN 'paid'
				WHEN amount_paid + $1 > 0 THEN 'partially_paid'
				WHEN due_date < CURRENT_DATE AND status NOT IN ('draft', 'void') THEN 'overdue'
				ELSE status
			END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, amount, invoiceID, tenantID, restaurantID)
	return err
}

func (r *PaymentAccountingRepository) generatePaymentNumberTx(tx *sql.Tx, mode string, tenantID, restaurantID int, paymentDate time.Time) (string, error) {
	tableName := "payments_received"
	prefix := fmt.Sprintf("PMT-R-%d-", paymentDate.Year())
	if mode == "made" {
		tableName = "payments_made"
		prefix = fmt.Sprintf("PMT-M-%d-", paymentDate.Year())
	}

	query := fmt.Sprintf(`
		SELECT COALESCE(MAX(CAST(RIGHT(payment_number, 5) AS INT)), 0) + 1
		FROM %s
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND payment_number LIKE $3
	`, tableName)

	var nextNumber int
	if err := tx.QueryRow(query, tenantID, restaurantID, prefix+"%").Scan(&nextNumber); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%05d", prefix, nextNumber), nil
}

type paymentReceivedScanner interface {
	Scan(dest ...interface{}) error
}

func scanPaymentReceived(scanner paymentReceivedScanner) (domain.PaymentReceived, error) {
	var (
		item           domain.PaymentReceived
		invoiceID      sql.NullInt64
		bankAccountID  sql.NullInt64
		journalEntryID sql.NullInt64
		createdBy      sql.NullInt64
	)

	err := scanner.Scan(
		&item.ID,
		&item.TenantID,
		&item.RestaurantID,
		&item.PaymentNumber,
		&item.PaymentDate,
		&item.CustomerID,
		&item.CustomerName,
		&invoiceID,
		&item.Amount,
		&item.PaymentMethod,
		&item.Reference,
		&bankAccountID,
		&item.Notes,
		&journalEntryID,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
		&createdBy,
	)
	if err != nil {
		return domain.PaymentReceived{}, err
	}
	if invoiceID.Valid {
		value := int(invoiceID.Int64)
		item.InvoiceID = &value
	}
	if bankAccountID.Valid {
		value := int(bankAccountID.Int64)
		item.BankAccountID = &value
	}
	if journalEntryID.Valid {
		value := int(journalEntryID.Int64)
		item.JournalEntryID = &value
	}
	if createdBy.Valid {
		value := int(createdBy.Int64)
		item.CreatedBy = &value
	}
	return item, nil
}

type paymentMadeScanner interface {
	Scan(dest ...interface{}) error
}

func scanPaymentMade(scanner paymentMadeScanner) (domain.PaymentMade, error) {
	var (
		item           domain.PaymentMade
		invoiceID      sql.NullInt64
		bankAccountID  sql.NullInt64
		journalEntryID sql.NullInt64
		createdBy      sql.NullInt64
	)

	err := scanner.Scan(
		&item.ID,
		&item.TenantID,
		&item.RestaurantID,
		&item.PaymentNumber,
		&item.PaymentDate,
		&item.SupplierID,
		&item.SupplierName,
		&invoiceID,
		&item.Amount,
		&item.PaymentMethod,
		&item.Reference,
		&bankAccountID,
		&item.Notes,
		&journalEntryID,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
		&createdBy,
	)
	if err != nil {
		return domain.PaymentMade{}, err
	}
	if invoiceID.Valid {
		value := int(invoiceID.Int64)
		item.InvoiceID = &value
	}
	if bankAccountID.Valid {
		value := int(bankAccountID.Int64)
		item.BankAccountID = &value
	}
	if journalEntryID.Valid {
		value := int(journalEntryID.Int64)
		item.JournalEntryID = &value
	}
	if createdBy.Valid {
		value := int(createdBy.Int64)
		item.CreatedBy = &value
	}
	return item, nil
}
