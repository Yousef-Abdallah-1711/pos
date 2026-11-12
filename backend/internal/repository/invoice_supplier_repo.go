package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"pos-saas/internal/domain"
)

type SupplierInvoiceRepository struct {
	db          *sql.DB
	journalRepo *JournalEntryRepository
	accountRepo *AccountRepository
}

func NewSupplierInvoiceRepository(db *sql.DB, journalRepo *JournalEntryRepository, accountRepo *AccountRepository) *SupplierInvoiceRepository {
	return &SupplierInvoiceRepository{
		db:          db,
		journalRepo: journalRepo,
		accountRepo: accountRepo,
	}
}

func (r *SupplierInvoiceRepository) ListInvoices(tenantID, restaurantID int) ([]domain.SupplierInvoice, error) {
	invoices, _, err := r.ListInvoicesPaged(tenantID, restaurantID, 1, 100000)
	return invoices, err
}

func (r *SupplierInvoiceRepository) ListInvoicesPaged(tenantID, restaurantID, page, pageSize int) ([]domain.SupplierInvoice, int, error) {
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
		FROM invoices_supplier inv
		WHERE inv.tenant_id = $1 AND inv.restaurant_id = $2
	`, tenantID, restaurantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT
			inv.id, inv.tenant_id, inv.restaurant_id, inv.invoice_number,
			COALESCE(inv.supplier_ref, ''), inv.invoice_date, inv.due_date,
			inv.supplier_id, COALESCE(s.company_name, ''), inv.subtotal, inv.tax_amount,
			inv.discount_amount, inv.total_amount, inv.amount_paid, inv.balance_due,
			inv.currency,
			CASE
				WHEN inv.status NOT IN ('paid', 'void', 'draft') AND inv.balance_due > 0 AND inv.due_date < CURRENT_DATE
					THEN 'overdue'
				ELSE inv.status
			END AS status,
			COALESCE(inv.notes, ''), COALESCE(inv.internal_notes, ''),
			inv.journal_entry_id, inv.created_at, inv.updated_at, inv.created_by, inv.updated_by
		FROM invoices_supplier inv
		JOIN suppliers s ON s.id = inv.supplier_id
		 AND s.tenant_id = inv.tenant_id AND s.restaurant_id = inv.restaurant_id
		WHERE inv.tenant_id = $1 AND inv.restaurant_id = $2
		ORDER BY inv.invoice_date DESC, inv.id DESC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	invoices := make([]domain.SupplierInvoice, 0)
	for rows.Next() {
		invoice, err := scanSupplierInvoice(rows)
		if err != nil {
			return nil, 0, err
		}
		invoices = append(invoices, invoice)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *SupplierInvoiceRepository) GetInvoiceByID(tenantID, restaurantID, id int) (*domain.SupplierInvoice, error) {
	invoice, err := scanSupplierInvoice(r.db.QueryRow(`
		SELECT
			inv.id, inv.tenant_id, inv.restaurant_id, inv.invoice_number,
			COALESCE(inv.supplier_ref, ''), inv.invoice_date, inv.due_date,
			inv.supplier_id, COALESCE(s.company_name, ''), inv.subtotal, inv.tax_amount,
			inv.discount_amount, inv.total_amount, inv.amount_paid, inv.balance_due,
			inv.currency,
			CASE
				WHEN inv.status NOT IN ('paid', 'void', 'draft') AND inv.balance_due > 0 AND inv.due_date < CURRENT_DATE
					THEN 'overdue'
				ELSE inv.status
			END AS status,
			COALESCE(inv.notes, ''), COALESCE(inv.internal_notes, ''),
			inv.journal_entry_id, inv.created_at, inv.updated_at, inv.created_by, inv.updated_by
		FROM invoices_supplier inv
		JOIN suppliers s ON s.id = inv.supplier_id
		 AND s.tenant_id = inv.tenant_id AND s.restaurant_id = inv.restaurant_id
		WHERE inv.tenant_id = $1 AND inv.restaurant_id = $2 AND inv.id = $3
	`, tenantID, restaurantID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	lines, err := r.getInvoiceLines(id)
	if err != nil {
		return nil, err
	}
	invoice.Lines = lines
	return &invoice, nil
}

func (r *SupplierInvoiceRepository) CreateInvoice(invoice *domain.SupplierInvoice) (int, error) {
	if invoice == nil {
		return 0, errors.New("invoice is required")
	}
	if len(invoice.Lines) == 0 {
		return 0, errors.New("at least one invoice line is required")
	}
	if invoice.InvoiceDate.IsZero() {
		invoice.InvoiceDate = time.Now()
	}
	if invoice.DueDate.IsZero() {
		invoice.DueDate = invoice.InvoiceDate
	}
	if strings.TrimSpace(invoice.Currency) == "" {
		invoice.Currency = "EGP"
	}
	if strings.TrimSpace(invoice.Status) == "" {
		invoice.Status = domain.InvoiceStatusDraft
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	invoiceNumber, err := r.generateInvoiceNumberTx(tx, invoice.TenantID, invoice.RestaurantID, invoice.InvoiceDate)
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRow(`
		INSERT INTO invoices_supplier (
			tenant_id, restaurant_id, invoice_number, supplier_ref, invoice_date,
			due_date, supplier_id, currency, status, notes, internal_notes,
			created_by, updated_by
		) SELECT
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11,
			$12, $13
		WHERE EXISTS (
			SELECT 1 FROM suppliers s
			WHERE s.id = $7 AND s.tenant_id = $1 AND s.restaurant_id = $2
		)
		RETURNING id
	`,
		invoice.TenantID,
		invoice.RestaurantID,
		invoiceNumber,
		nullString(invoice.SupplierRef),
		invoice.InvoiceDate,
		invoice.DueDate,
		invoice.SupplierID,
		invoice.Currency,
		invoice.Status,
		nullString(invoice.Notes),
		nullString(invoice.InternalNotes),
		invoice.CreatedBy,
		invoice.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	if err := r.insertSupplierInvoiceLinesTx(tx, id, invoice.TenantID, invoice.RestaurantID, invoice.Lines); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SupplierInvoiceRepository) UpdateInvoice(invoice *domain.SupplierInvoice) error {
	if invoice == nil {
		return errors.New("invoice is required")
	}
	if len(invoice.Lines) == 0 {
		return errors.New("at least one invoice line is required")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	existing, err := r.getInvoiceByIDTx(tx, invoice.TenantID, invoice.RestaurantID, invoice.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("supplier invoice not found")
	}
	if existing.Status != domain.InvoiceStatusDraft {
		return errors.New("only draft supplier invoices can be edited")
	}

	if invoice.InvoiceDate.IsZero() {
		invoice.InvoiceDate = existing.InvoiceDate
	}
	if invoice.DueDate.IsZero() {
		invoice.DueDate = existing.DueDate
	}
	if strings.TrimSpace(invoice.Currency) == "" {
		invoice.Currency = existing.Currency
	}
	if strings.TrimSpace(invoice.Status) == "" {
		invoice.Status = existing.Status
	}

	result, err := tx.Exec(`
		UPDATE invoices_supplier
		SET supplier_ref = $1,
			invoice_date = $2,
			due_date = $3,
			supplier_id = $4,
			currency = $5,
			status = $6,
			notes = $7,
			internal_notes = $8,
			updated_by = $9,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $10 AND tenant_id = $11 AND restaurant_id = $12
		  AND EXISTS (
			SELECT 1 FROM suppliers s
			WHERE s.id = $4 AND s.tenant_id = $11 AND s.restaurant_id = $12
		  )
	`,
		nullString(invoice.SupplierRef),
		invoice.InvoiceDate,
		invoice.DueDate,
		invoice.SupplierID,
		invoice.Currency,
		invoice.Status,
		nullString(invoice.Notes),
		nullString(invoice.InternalNotes),
		invoice.UpdatedBy,
		invoice.ID,
		invoice.TenantID,
		invoice.RestaurantID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("supplier invoice not found or supplier is outside scope")
	}

	if _, err := tx.Exec(`DELETE FROM invoice_supplier_lines WHERE invoice_id = $1`, invoice.ID); err != nil {
		return err
	}
	if err := r.insertSupplierInvoiceLinesTx(tx, invoice.ID, invoice.TenantID, invoice.RestaurantID, invoice.Lines); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SupplierInvoiceRepository) ApproveInvoice(tenantID, restaurantID, id int, actor *int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	invoice, err := r.getInvoiceByIDTx(tx, tenantID, restaurantID, id)
	if err != nil {
		return err
	}
	if invoice == nil {
		return errors.New("supplier invoice not found")
	}
	if invoice.Status != domain.InvoiceStatusDraft {
		return errors.New("only draft supplier invoices can be approved")
	}

	entry, err := r.buildJournalEntry(tx, invoice, actor)
	if err != nil {
		return err
	}
	journalID, err := r.journalRepo.CreateJournalEntryTx(tx, &entry)
	if err != nil {
		return err
	}

	postedBy := 0
	if actor != nil {
		postedBy = *actor
	}
	if err := r.journalRepo.PostEntryTx(tx, tenantID, restaurantID, journalID, postedBy); err != nil {
		return err
	}

	nextStatus := domain.InvoiceStatusApproved
	if invoice.BalanceDue <= 0 {
		nextStatus = domain.InvoiceStatusPaid
	}

	_, err = tx.Exec(`
		UPDATE invoices_supplier
		SET journal_entry_id = $1,
			status = $2,
			updated_by = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND tenant_id = $5 AND restaurant_id = $6
	`, journalID, nextStatus, actor, id, tenantID, restaurantID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SupplierInvoiceRepository) VoidInvoice(tenantID, restaurantID, id int, voidedBy int, reason string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	invoice, err := r.getInvoiceByIDTx(tx, tenantID, restaurantID, id)
	if err != nil {
		return err
	}
	if invoice == nil {
		return errors.New("supplier invoice not found")
	}
	if invoice.AmountPaid > 0 {
		return errors.New("supplier invoice with payments cannot be voided")
	}
	if invoice.JournalEntryID != nil {
		if err := r.journalRepo.VoidEntryTx(tx, tenantID, restaurantID, *invoice.JournalEntryID, voidedBy, reason); err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
		UPDATE invoices_supplier
		SET status = $1,
			internal_notes = CONCAT(COALESCE(internal_notes, ''), CASE WHEN COALESCE(internal_notes, '') = '' THEN '' ELSE E'\n' END, $2),
			updated_by = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND tenant_id = $5 AND restaurant_id = $6
	`, domain.InvoiceStatusVoid, strings.TrimSpace(reason), voidedBy, id, tenantID, restaurantID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SupplierInvoiceRepository) getInvoiceLines(invoiceID int) ([]domain.InvoiceLine, error) {
	rows, err := r.db.Query(`
		SELECT
			id, invoice_id, tenant_id, description, COALESCE(description_ar, ''),
			quantity, unit_price, tax_rate, tax_amount, discount_rate,
			discount_amount, line_total, account_id, line_order, created_at
		FROM invoice_supplier_lines
		WHERE invoice_id = $1
		ORDER BY line_order ASC, id ASC
	`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanInvoiceLines(rows)
}

func (r *SupplierInvoiceRepository) getInvoiceByIDTx(tx *sql.Tx, tenantID, restaurantID, id int) (*domain.SupplierInvoice, error) {
	invoice, err := scanSupplierInvoice(tx.QueryRow(`
		SELECT
			inv.id, inv.tenant_id, inv.restaurant_id, inv.invoice_number,
			COALESCE(inv.supplier_ref, ''), inv.invoice_date, inv.due_date,
			inv.supplier_id, COALESCE(s.company_name, ''), inv.subtotal, inv.tax_amount,
			inv.discount_amount, inv.total_amount, inv.amount_paid, inv.balance_due,
			inv.currency, inv.status, COALESCE(inv.notes, ''), COALESCE(inv.internal_notes, ''),
			inv.journal_entry_id, inv.created_at, inv.updated_at, inv.created_by, inv.updated_by
		FROM invoices_supplier inv
		JOIN suppliers s ON s.id = inv.supplier_id
		 AND s.tenant_id = inv.tenant_id AND s.restaurant_id = inv.restaurant_id
		WHERE inv.tenant_id = $1 AND inv.restaurant_id = $2 AND inv.id = $3
	`, tenantID, restaurantID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	rows, err := tx.Query(`
		SELECT
			id, invoice_id, tenant_id, description, COALESCE(description_ar, ''),
			quantity, unit_price, tax_rate, tax_amount, discount_rate,
			discount_amount, line_total, account_id, line_order, created_at
		FROM invoice_supplier_lines
		WHERE invoice_id = $1
		ORDER BY line_order ASC, id ASC
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lines, err := scanInvoiceLines(rows)
	if err != nil {
		return nil, err
	}
	invoice.Lines = lines
	invoice.Status = effectiveInvoiceStatus(invoice.Status, invoice.DueDate, invoice.BalanceDue, domain.InvoiceStatusApproved)
	return &invoice, nil
}

func (r *SupplierInvoiceRepository) insertSupplierInvoiceLinesTx(tx *sql.Tx, invoiceID, tenantID, restaurantID int, lines []domain.InvoiceLine) error {
	for index, line := range lines {
		lineOrder := line.LineOrder
		if lineOrder <= 0 {
			lineOrder = index + 1
		}
		result, err := tx.Exec(`
			INSERT INTO invoice_supplier_lines (
				invoice_id, tenant_id, description, description_ar, quantity,
				unit_price, tax_rate, tax_amount, discount_rate, discount_amount,
				line_total, account_id, line_order
			) SELECT
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10,
				$11, $12, $13
			WHERE $12::integer IS NULL OR EXISTS (
				SELECT 1 FROM chart_of_accounts a
				WHERE a.id = $12 AND a.tenant_id = $2 AND a.restaurant_id = $14
			)
		`,
			invoiceID,
			tenantID,
			line.Description,
			nullString(line.DescriptionAr),
			line.Quantity,
			line.UnitPrice,
			line.TaxRate,
			line.TaxAmount,
			line.DiscountRate,
			line.DiscountAmount,
			line.LineTotal,
			line.AccountID,
			lineOrder,
			restaurantID,
		)
		if err != nil {
			return err
		}
		rowsAffected, affectedErr := result.RowsAffected()
		if affectedErr != nil {
			return affectedErr
		}
		if rowsAffected == 0 {
			return errors.New("supplier invoice line account is outside scope")
		}
	}
	return nil
}

func (r *SupplierInvoiceRepository) buildJournalEntry(tx *sql.Tx, invoice *domain.SupplierInvoice, actor *int) (domain.JournalEntry, error) {
	payableAccountID, err := r.lookupSupplierPayableAccountTx(tx, invoice.TenantID, invoice.RestaurantID, invoice.SupplierID)
	if err != nil {
		return domain.JournalEntry{}, err
	}
	defaultExpense, err := lookupAccountByCode(r.accountRepo, invoice.TenantID, invoice.RestaurantID, "5110")
	if err != nil {
		return domain.JournalEntry{}, err
	}
	vatAccount, err := lookupAccountByCode(r.accountRepo, invoice.TenantID, invoice.RestaurantID, "1150")
	if err != nil {
		return domain.JournalEntry{}, err
	}

	expenseByAccount := make(map[int]float64)
	for _, line := range invoice.Lines {
		accountID := defaultExpense.ID
		if line.AccountID != nil {
			accountID = *line.AccountID
		}
		netAmount := roundAccountingAmount(line.LineTotal - line.TaxAmount)
		if netAmount <= 0 {
			continue
		}
		expenseByAccount[accountID] = roundAccountingAmount(expenseByAccount[accountID] + netAmount)
	}

	accountIDs := make([]int, 0, len(expenseByAccount))
	for accountID := range expenseByAccount {
		accountIDs = append(accountIDs, accountID)
	}
	sort.Ints(accountIDs)

	lines := make([]domain.JournalEntryLine, 0, len(accountIDs)+2)
	lineOrder := 1
	for _, accountID := range accountIDs {
		lines = append(lines, domain.JournalEntryLine{
			AccountID:   accountID,
			DebitAmount: expenseByAccount[accountID],
			Description: fmt.Sprintf("Supplier invoice expense %s", invoice.InvoiceNumber),
			LineOrder:   lineOrder,
		})
		lineOrder++
	}
	if invoice.TaxAmount > 0 {
		lines = append(lines, domain.JournalEntryLine{
			AccountID:   vatAccount.ID,
			DebitAmount: roundAccountingAmount(invoice.TaxAmount),
			Description: fmt.Sprintf("Input VAT for supplier invoice %s", invoice.InvoiceNumber),
			LineOrder:   lineOrder,
		})
		lineOrder++
	}
	lines = append(lines, domain.JournalEntryLine{
		AccountID:    payableAccountID,
		CreditAmount: roundAccountingAmount(invoice.TotalAmount),
		Description:  fmt.Sprintf("Accounts payable for supplier invoice %s", invoice.InvoiceNumber),
		LineOrder:    lineOrder,
	})

	return domain.JournalEntry{
		TenantID:        invoice.TenantID,
		RestaurantID:    invoice.RestaurantID,
		EntryDate:       invoice.InvoiceDate,
		Description:     fmt.Sprintf("Supplier invoice %s", invoice.InvoiceNumber),
		ReferenceType:   "supplier_invoice",
		ReferenceID:     &invoice.ID,
		ReferenceNumber: invoice.InvoiceNumber,
		EntryType:       "supplier_invoice",
		Status:          domain.JournalStatusDraft,
		CreatedBy:       actor,
		UpdatedBy:       actor,
		Lines:           lines,
	}, nil
}

func (r *SupplierInvoiceRepository) lookupSupplierPayableAccountTx(tx *sql.Tx, tenantID, restaurantID, supplierID int) (int, error) {
	var payableAccountID sql.NullInt64
	err := tx.QueryRow(`
		SELECT payable_account_id
		FROM suppliers
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, supplierID, tenantID, restaurantID).Scan(&payableAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("supplier not found")
		}
		return 0, err
	}
	if payableAccountID.Valid {
		return int(payableAccountID.Int64), nil
	}

	account, err := lookupAccountByCode(r.accountRepo, tenantID, restaurantID, "2101")
	if err != nil {
		return 0, err
	}
	return account.ID, nil
}

func (r *SupplierInvoiceRepository) generateInvoiceNumberTx(tx *sql.Tx, tenantID, restaurantID int, invoiceDate time.Time) (string, error) {
	var nextNumber int
	prefix := fmt.Sprintf("PINV-%d-", invoiceDate.Year())
	err := tx.QueryRow(`
		SELECT COALESCE(MAX(CAST(RIGHT(invoice_number, 5) AS INT)), 0) + 1
		FROM invoices_supplier
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND invoice_number LIKE $3
	`, tenantID, restaurantID, prefix+"%").Scan(&nextNumber)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%05d", prefix, nextNumber), nil
}

type supplierInvoiceScanner interface {
	Scan(dest ...interface{}) error
}

func scanSupplierInvoice(scanner supplierInvoiceScanner) (domain.SupplierInvoice, error) {
	var (
		invoice        domain.SupplierInvoice
		journalEntryID sql.NullInt64
		createdBy      sql.NullInt64
		updatedBy      sql.NullInt64
	)

	err := scanner.Scan(
		&invoice.ID,
		&invoice.TenantID,
		&invoice.RestaurantID,
		&invoice.InvoiceNumber,
		&invoice.SupplierRef,
		&invoice.InvoiceDate,
		&invoice.DueDate,
		&invoice.SupplierID,
		&invoice.SupplierName,
		&invoice.Subtotal,
		&invoice.TaxAmount,
		&invoice.DiscountAmount,
		&invoice.TotalAmount,
		&invoice.AmountPaid,
		&invoice.BalanceDue,
		&invoice.Currency,
		&invoice.Status,
		&invoice.Notes,
		&invoice.InternalNotes,
		&journalEntryID,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		return domain.SupplierInvoice{}, err
	}

	if journalEntryID.Valid {
		value := int(journalEntryID.Int64)
		invoice.JournalEntryID = &value
	}
	if createdBy.Valid {
		value := int(createdBy.Int64)
		invoice.CreatedBy = &value
	}
	if updatedBy.Valid {
		value := int(updatedBy.Int64)
		invoice.UpdatedBy = &value
	}

	return invoice, nil
}
