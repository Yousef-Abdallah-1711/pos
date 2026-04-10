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

type CustomerInvoiceRepository struct {
	db          *sql.DB
	journalRepo *JournalEntryRepository
	accountRepo *AccountRepository
}

func NewCustomerInvoiceRepository(db *sql.DB, journalRepo *JournalEntryRepository, accountRepo *AccountRepository) *CustomerInvoiceRepository {
	return &CustomerInvoiceRepository{
		db:          db,
		journalRepo: journalRepo,
		accountRepo: accountRepo,
	}
}

func (r *CustomerInvoiceRepository) ListInvoices(tenantID, restaurantID int) ([]domain.CustomerInvoice, error) {
	invoices, _, err := r.ListInvoicesPaged(tenantID, restaurantID, 1, 100000)
	return invoices, err
}

func (r *CustomerInvoiceRepository) ListInvoicesPaged(tenantID, restaurantID, page, pageSize int) ([]domain.CustomerInvoice, int, error) {
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
		FROM invoices_customer ic
		WHERE ic.tenant_id = $1 AND ic.restaurant_id = $2
	`, tenantID, restaurantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT
			ic.id, ic.tenant_id, ic.restaurant_id, ic.invoice_number,
			ic.invoice_date, ic.due_date, ic.customer_id, COALESCE(c.customer_name, ''),
			ic.order_id, ic.subtotal, ic.tax_amount, ic.discount_amount,
			ic.total_amount, ic.amount_paid, ic.balance_due, ic.currency,
			CASE
				WHEN ic.status NOT IN ('paid', 'void') AND ic.balance_due > 0 AND ic.due_date < CURRENT_DATE AND ic.status <> 'draft'
					THEN 'overdue'
				ELSE ic.status
			END AS status,
			COALESCE(ic.notes, ''), COALESCE(ic.terms, ''), COALESCE(ic.internal_notes, ''),
			ic.journal_entry_id, ic.created_at, ic.updated_at, ic.created_by, ic.updated_by
		FROM invoices_customer ic
		JOIN customers_accounting c ON c.id = ic.customer_id
		WHERE ic.tenant_id = $1 AND ic.restaurant_id = $2
		ORDER BY ic.invoice_date DESC, ic.id DESC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	invoices := make([]domain.CustomerInvoice, 0)
	for rows.Next() {
		invoice, err := scanCustomerInvoice(rows)
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

func (r *CustomerInvoiceRepository) GetInvoiceByID(tenantID, restaurantID, id int) (*domain.CustomerInvoice, error) {
	invoice, err := scanCustomerInvoice(r.db.QueryRow(`
		SELECT
			ic.id, ic.tenant_id, ic.restaurant_id, ic.invoice_number,
			ic.invoice_date, ic.due_date, ic.customer_id, COALESCE(c.customer_name, ''),
			ic.order_id, ic.subtotal, ic.tax_amount, ic.discount_amount,
			ic.total_amount, ic.amount_paid, ic.balance_due, ic.currency,
			CASE
				WHEN ic.status NOT IN ('paid', 'void') AND ic.balance_due > 0 AND ic.due_date < CURRENT_DATE AND ic.status <> 'draft'
					THEN 'overdue'
				ELSE ic.status
			END AS status,
			COALESCE(ic.notes, ''), COALESCE(ic.terms, ''), COALESCE(ic.internal_notes, ''),
			ic.journal_entry_id, ic.created_at, ic.updated_at, ic.created_by, ic.updated_by
		FROM invoices_customer ic
		JOIN customers_accounting c ON c.id = ic.customer_id
		WHERE ic.tenant_id = $1 AND ic.restaurant_id = $2 AND ic.id = $3
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

func (r *CustomerInvoiceRepository) CreateInvoice(invoice *domain.CustomerInvoice) (int, error) {
	if invoice == nil {
		return 0, errors.New("invoice is required")
	}
	if len(invoice.Lines) == 0 {
		return 0, errors.New("at least one invoice line is required")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	id, err := r.createInvoiceTx(tx, invoice)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CustomerInvoiceRepository) UpdateInvoice(invoice *domain.CustomerInvoice) error {
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
		return errors.New("customer invoice not found")
	}
	if existing.Status == domain.InvoiceStatusVoid || existing.AmountPaid > 0 {
		return errors.New("customer invoice cannot be edited after payment or void")
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

	_, err = tx.Exec(`
		UPDATE invoices_customer
		SET invoice_date = $1,
			due_date = $2,
			customer_id = $3,
			order_id = $4,
			currency = $5,
			status = $6,
			notes = $7,
			terms = $8,
			internal_notes = $9,
			updated_by = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $11 AND tenant_id = $12 AND restaurant_id = $13
	`,
		invoice.InvoiceDate,
		invoice.DueDate,
		invoice.CustomerID,
		invoice.OrderID,
		invoice.Currency,
		invoice.Status,
		nullString(invoice.Notes),
		nullString(invoice.Terms),
		nullString(invoice.InternalNotes),
		invoice.UpdatedBy,
		invoice.ID,
		invoice.TenantID,
		invoice.RestaurantID,
	)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM invoice_customer_lines WHERE invoice_id = $1`, invoice.ID); err != nil {
		return err
	}
	if err := r.insertCustomerInvoiceLinesTx(tx, invoice.ID, invoice.TenantID, invoice.Lines); err != nil {
		return err
	}

	if existing.JournalEntryID != nil {
		voidedBy := 0
		if invoice.UpdatedBy != nil {
			voidedBy = *invoice.UpdatedBy
		}
		if err := r.journalRepo.VoidEntryTx(tx, invoice.TenantID, invoice.RestaurantID, *existing.JournalEntryID, voidedBy, "Customer invoice updated"); err != nil {
			return err
		}
	}

	updated, err := r.getInvoiceByIDTx(tx, invoice.TenantID, invoice.RestaurantID, invoice.ID)
	if err != nil {
		return err
	}
	if updated == nil {
		return errors.New("customer invoice not found")
	}

	journalID, err := r.syncJournalEntryTx(tx, updated, invoice.UpdatedBy)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`
		UPDATE invoices_customer
		SET journal_entry_id = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, journalID, invoice.ID, invoice.TenantID, invoice.RestaurantID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CustomerInvoiceRepository) MarkAsSent(tenantID, restaurantID, id int, updatedBy *int) error {
	status, err := r.GetInvoiceByID(tenantID, restaurantID, id)
	if err != nil {
		return err
	}
	if status == nil {
		return errors.New("customer invoice not found")
	}
	nextStatus := domain.InvoiceStatusSent
	if status.BalanceDue > 0 && status.DueDate.Before(time.Now()) {
		nextStatus = domain.InvoiceStatusOverdue
	}

	_, err = r.db.Exec(`
		UPDATE invoices_customer
		SET status = $1,
			updated_by = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5
	`, nextStatus, updatedBy, id, tenantID, restaurantID)
	return err
}

func (r *CustomerInvoiceRepository) VoidInvoice(tenantID, restaurantID, id int, voidedBy int, reason string) error {
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
		return errors.New("customer invoice not found")
	}
	if invoice.AmountPaid > 0 {
		return errors.New("customer invoice with payments cannot be voided")
	}
	if invoice.JournalEntryID != nil {
		if err := r.journalRepo.VoidEntryTx(tx, tenantID, restaurantID, *invoice.JournalEntryID, voidedBy, reason); err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
		UPDATE invoices_customer
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

func (r *CustomerInvoiceRepository) GetOverdueInvoices(tenantID, restaurantID int) ([]domain.CustomerInvoice, error) {
	invoices, _, err := r.GetOverdueInvoicesPaged(tenantID, restaurantID, 1, 100000)
	return invoices, err
}

func (r *CustomerInvoiceRepository) GetOverdueInvoicesPaged(tenantID, restaurantID, page, pageSize int) ([]domain.CustomerInvoice, int, error) {
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
		FROM invoices_customer ic
		WHERE ic.tenant_id = $1
		  AND ic.restaurant_id = $2
		  AND ic.status NOT IN ('paid', 'void', 'draft')
		  AND ic.balance_due > 0
		  AND ic.due_date < CURRENT_DATE
	`, tenantID, restaurantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT
			ic.id, ic.tenant_id, ic.restaurant_id, ic.invoice_number,
			ic.invoice_date, ic.due_date, ic.customer_id, COALESCE(c.customer_name, ''),
			ic.order_id, ic.subtotal, ic.tax_amount, ic.discount_amount,
			ic.total_amount, ic.amount_paid, ic.balance_due, ic.currency,
			'overdue' AS status,
			COALESCE(ic.notes, ''), COALESCE(ic.terms, ''), COALESCE(ic.internal_notes, ''),
			ic.journal_entry_id, ic.created_at, ic.updated_at, ic.created_by, ic.updated_by
		FROM invoices_customer ic
		JOIN customers_accounting c ON c.id = ic.customer_id
		WHERE ic.tenant_id = $1
		  AND ic.restaurant_id = $2
		  AND ic.status NOT IN ('paid', 'void', 'draft')
		  AND ic.balance_due > 0
		  AND ic.due_date < CURRENT_DATE
		ORDER BY ic.due_date ASC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	invoices := make([]domain.CustomerInvoice, 0)
	for rows.Next() {
		invoice, err := scanCustomerInvoice(rows)
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

func (r *CustomerInvoiceRepository) createInvoiceTx(tx *sql.Tx, invoice *domain.CustomerInvoice) (int, error) {
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

	invoiceNumber, err := r.generateInvoiceNumberTx(tx, invoice.TenantID, invoice.RestaurantID, invoice.InvoiceDate)
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRow(`
		INSERT INTO invoices_customer (
			tenant_id, restaurant_id, invoice_number, invoice_date, due_date,
			customer_id, order_id, currency, status, notes, terms, internal_notes,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11, $12,
			$13, $14
		)
		RETURNING id
	`,
		invoice.TenantID,
		invoice.RestaurantID,
		invoiceNumber,
		invoice.InvoiceDate,
		invoice.DueDate,
		invoice.CustomerID,
		invoice.OrderID,
		invoice.Currency,
		invoice.Status,
		nullString(invoice.Notes),
		nullString(invoice.Terms),
		nullString(invoice.InternalNotes),
		invoice.CreatedBy,
		invoice.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	if err := r.insertCustomerInvoiceLinesTx(tx, id, invoice.TenantID, invoice.Lines); err != nil {
		return 0, err
	}

	created, err := r.getInvoiceByIDTx(tx, invoice.TenantID, invoice.RestaurantID, id)
	if err != nil {
		return 0, err
	}
	if created == nil {
		return 0, errors.New("customer invoice not found")
	}

	journalID, err := r.syncJournalEntryTx(tx, created, invoice.CreatedBy)
	if err != nil {
		return 0, err
	}

	if _, err := tx.Exec(`
		UPDATE invoices_customer
		SET journal_entry_id = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, journalID, id, invoice.TenantID, invoice.RestaurantID); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *CustomerInvoiceRepository) syncJournalEntryTx(tx *sql.Tx, invoice *domain.CustomerInvoice, actor *int) (int, error) {
	entry, err := r.buildJournalEntry(tx, invoice, actor)
	if err != nil {
		return 0, err
	}

	journalID, err := r.journalRepo.CreateJournalEntryTx(tx, &entry)
	if err != nil {
		return 0, err
	}

	postedBy := 0
	if actor != nil {
		postedBy = *actor
	}
	if err := r.journalRepo.PostEntryTx(tx, invoice.TenantID, invoice.RestaurantID, journalID, postedBy); err != nil {
		return 0, err
	}

	return journalID, nil
}

func (r *CustomerInvoiceRepository) buildJournalEntry(tx *sql.Tx, invoice *domain.CustomerInvoice, actor *int) (domain.JournalEntry, error) {
	receivableAccountID, err := r.lookupCustomerReceivableAccountTx(tx, invoice.TenantID, invoice.RestaurantID, invoice.CustomerID)
	if err != nil {
		return domain.JournalEntry{}, err
	}
	defaultRevenue, err := lookupAccountByCode(r.accountRepo, invoice.TenantID, invoice.RestaurantID, "4100")
	if err != nil {
		return domain.JournalEntry{}, err
	}
	vatAccount, err := lookupAccountByCode(r.accountRepo, invoice.TenantID, invoice.RestaurantID, "2120")
	if err != nil {
		return domain.JournalEntry{}, err
	}

	lines := []domain.JournalEntryLine{
		{
			AccountID:   receivableAccountID,
			DebitAmount: roundAccountingAmount(invoice.TotalAmount),
			Description: fmt.Sprintf("Accounts receivable for invoice %s", invoice.InvoiceNumber),
			LineOrder:   1,
		},
	}

	revenueByAccount := make(map[int]float64)
	for _, line := range invoice.Lines {
		accountID := defaultRevenue.ID
		if line.AccountID != nil {
			accountID = *line.AccountID
		}
		netRevenue := roundAccountingAmount(line.LineTotal - line.TaxAmount)
		if netRevenue <= 0 {
			continue
		}
		revenueByAccount[accountID] = roundAccountingAmount(revenueByAccount[accountID] + netRevenue)
	}

	accountIDs := make([]int, 0, len(revenueByAccount))
	for accountID := range revenueByAccount {
		accountIDs = append(accountIDs, accountID)
	}
	sort.Ints(accountIDs)

	lineOrder := 2
	for _, accountID := range accountIDs {
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    accountID,
			CreditAmount: revenueByAccount[accountID],
			Description:  fmt.Sprintf("Revenue recognition for invoice %s", invoice.InvoiceNumber),
			LineOrder:    lineOrder,
		})
		lineOrder++
	}
	if invoice.TaxAmount > 0 {
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    vatAccount.ID,
			CreditAmount: roundAccountingAmount(invoice.TaxAmount),
			Description:  fmt.Sprintf("Output VAT for invoice %s", invoice.InvoiceNumber),
			LineOrder:    lineOrder,
		})
	}

	return domain.JournalEntry{
		TenantID:        invoice.TenantID,
		RestaurantID:    invoice.RestaurantID,
		EntryDate:       invoice.InvoiceDate,
		Description:     fmt.Sprintf("Customer invoice %s", invoice.InvoiceNumber),
		ReferenceType:   "customer_invoice",
		ReferenceID:     &invoice.ID,
		ReferenceNumber: invoice.InvoiceNumber,
		EntryType:       "customer_invoice",
		Status:          domain.JournalStatusDraft,
		CreatedBy:       actor,
		UpdatedBy:       actor,
		Lines:           lines,
	}, nil
}

func (r *CustomerInvoiceRepository) lookupCustomerReceivableAccountTx(tx *sql.Tx, tenantID, restaurantID, customerID int) (int, error) {
	var receivableAccountID sql.NullInt64
	err := tx.QueryRow(`
		SELECT receivable_account_id
		FROM customers_accounting
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, customerID, tenantID, restaurantID).Scan(&receivableAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("customer account not found")
		}
		return 0, err
	}
	if receivableAccountID.Valid {
		return int(receivableAccountID.Int64), nil
	}

	account, err := lookupAccountByCode(r.accountRepo, tenantID, restaurantID, "1120")
	if err != nil {
		return 0, err
	}
	return account.ID, nil
}

func (r *CustomerInvoiceRepository) getInvoiceLines(invoiceID int) ([]domain.InvoiceLine, error) {
	rows, err := r.db.Query(`
		SELECT
			id, invoice_id, tenant_id, description, COALESCE(description_ar, ''),
			quantity, unit_price, tax_rate, tax_amount, discount_rate,
			discount_amount, line_total, account_id, line_order, created_at
		FROM invoice_customer_lines
		WHERE invoice_id = $1
		ORDER BY line_order ASC, id ASC
	`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanInvoiceLines(rows)
}

func (r *CustomerInvoiceRepository) getInvoiceByIDTx(tx *sql.Tx, tenantID, restaurantID, id int) (*domain.CustomerInvoice, error) {
	invoice, err := scanCustomerInvoice(tx.QueryRow(`
		SELECT
			ic.id, ic.tenant_id, ic.restaurant_id, ic.invoice_number,
			ic.invoice_date, ic.due_date, ic.customer_id, COALESCE(c.customer_name, ''),
			ic.order_id, ic.subtotal, ic.tax_amount, ic.discount_amount,
			ic.total_amount, ic.amount_paid, ic.balance_due, ic.currency,
			ic.status, COALESCE(ic.notes, ''), COALESCE(ic.terms, ''),
			COALESCE(ic.internal_notes, ''), ic.journal_entry_id, ic.created_at, ic.updated_at,
			ic.created_by, ic.updated_by
		FROM invoices_customer ic
		JOIN customers_accounting c ON c.id = ic.customer_id
		WHERE ic.tenant_id = $1 AND ic.restaurant_id = $2 AND ic.id = $3
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
		FROM invoice_customer_lines
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
	invoice.Status = effectiveInvoiceStatus(invoice.Status, invoice.DueDate, invoice.BalanceDue, domain.InvoiceStatusSent)
	return &invoice, nil
}

func (r *CustomerInvoiceRepository) insertCustomerInvoiceLinesTx(tx *sql.Tx, invoiceID, tenantID int, lines []domain.InvoiceLine) error {
	for index, line := range lines {
		lineOrder := line.LineOrder
		if lineOrder <= 0 {
			lineOrder = index + 1
		}
		if _, err := tx.Exec(`
			INSERT INTO invoice_customer_lines (
				invoice_id, tenant_id, description, description_ar, quantity,
				unit_price, tax_rate, tax_amount, discount_rate, discount_amount,
				line_total, account_id, line_order
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10,
				$11, $12, $13
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
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *CustomerInvoiceRepository) generateInvoiceNumberTx(tx *sql.Tx, tenantID, restaurantID int, invoiceDate time.Time) (string, error) {
	var nextNumber int
	prefix := fmt.Sprintf("INV-%d-", invoiceDate.Year())
	err := tx.QueryRow(`
		SELECT COALESCE(MAX(CAST(RIGHT(invoice_number, 5) AS INT)), 0) + 1
		FROM invoices_customer
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND invoice_number LIKE $3
	`, tenantID, restaurantID, prefix+"%").Scan(&nextNumber)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%05d", prefix, nextNumber), nil
}

type customerInvoiceScanner interface {
	Scan(dest ...interface{}) error
}

func scanCustomerInvoice(scanner customerInvoiceScanner) (domain.CustomerInvoice, error) {
	var (
		invoice        domain.CustomerInvoice
		orderID        sql.NullInt64
		journalEntryID sql.NullInt64
		createdBy      sql.NullInt64
		updatedBy      sql.NullInt64
	)

	err := scanner.Scan(
		&invoice.ID,
		&invoice.TenantID,
		&invoice.RestaurantID,
		&invoice.InvoiceNumber,
		&invoice.InvoiceDate,
		&invoice.DueDate,
		&invoice.CustomerID,
		&invoice.CustomerName,
		&orderID,
		&invoice.Subtotal,
		&invoice.TaxAmount,
		&invoice.DiscountAmount,
		&invoice.TotalAmount,
		&invoice.AmountPaid,
		&invoice.BalanceDue,
		&invoice.Currency,
		&invoice.Status,
		&invoice.Notes,
		&invoice.Terms,
		&invoice.InternalNotes,
		&journalEntryID,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		return domain.CustomerInvoice{}, err
	}

	if orderID.Valid {
		value := int(orderID.Int64)
		invoice.OrderID = &value
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

func scanInvoiceLines(rows *sql.Rows) ([]domain.InvoiceLine, error) {
	lines := make([]domain.InvoiceLine, 0)
	for rows.Next() {
		var (
			line      domain.InvoiceLine
			accountID sql.NullInt64
		)
		if err := rows.Scan(
			&line.ID,
			&line.InvoiceID,
			&line.TenantID,
			&line.Description,
			&line.DescriptionAr,
			&line.Quantity,
			&line.UnitPrice,
			&line.TaxRate,
			&line.TaxAmount,
			&line.DiscountRate,
			&line.DiscountAmount,
			&line.LineTotal,
			&accountID,
			&line.LineOrder,
			&line.CreatedAt,
		); err != nil {
			return nil, err
		}
		if accountID.Valid {
			value := int(accountID.Int64)
			line.AccountID = &value
		}
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
