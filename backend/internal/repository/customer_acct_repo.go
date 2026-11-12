package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"pos-saas/internal/domain"
)

type CustomerAccountRepository struct {
	db *sql.DB
}

func NewCustomerAccountRepository(db *sql.DB) *CustomerAccountRepository {
	return &CustomerAccountRepository{db: db}
}

func (r *CustomerAccountRepository) ListCustomers(tenantID, restaurantID int) ([]domain.CustomerAccount, error) {
	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, customer_code, customer_name,
			COALESCE(customer_name_ar, ''), COALESCE(customer_type, 'individual'),
			COALESCE(contact_person, ''), COALESCE(email, ''), COALESCE(phone, ''),
			COALESCE(mobile, ''), COALESCE(address, ''), COALESCE(address_ar, ''),
			COALESCE(city, ''), COALESCE(country, 'Egypt'), COALESCE(tax_number, ''),
			COALESCE(payment_terms, 'immediate'), COALESCE(credit_limit, 0),
			COALESCE(current_balance, 0), COALESCE(currency, 'EGP'),
			receivable_account_id, COALESCE(notes, ''), is_active, created_at, updated_at,
			created_by, updated_by
		FROM customers_accounting
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY customer_name ASC
	`, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make([]domain.CustomerAccount, 0)
	for rows.Next() {
		customer, err := scanCustomerAccount(rows)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomerAccountRepository) GetCustomerByID(tenantID, restaurantID, id int) (*domain.CustomerAccount, error) {
	customer, err := scanCustomerAccount(r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, customer_code, customer_name,
			COALESCE(customer_name_ar, ''), COALESCE(customer_type, 'individual'),
			COALESCE(contact_person, ''), COALESCE(email, ''), COALESCE(phone, ''),
			COALESCE(mobile, ''), COALESCE(address, ''), COALESCE(address_ar, ''),
			COALESCE(city, ''), COALESCE(country, 'Egypt'), COALESCE(tax_number, ''),
			COALESCE(payment_terms, 'immediate'), COALESCE(credit_limit, 0),
			COALESCE(current_balance, 0), COALESCE(currency, 'EGP'),
			receivable_account_id, COALESCE(notes, ''), is_active, created_at, updated_at,
			created_by, updated_by
		FROM customers_accounting
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerAccountRepository) CreateCustomer(customer *domain.CustomerAccount) (int, error) {
	code, err := r.generateCustomerCode(customer.TenantID, customer.RestaurantID)
	if err != nil {
		return 0, err
	}

	var id int
	err = r.db.QueryRow(`
		INSERT INTO customers_accounting (
			tenant_id, restaurant_id, customer_code, customer_name, customer_name_ar,
			customer_type, contact_person, email, phone, mobile, address, address_ar,
			city, country, tax_number, payment_terms, credit_limit, current_balance,
			currency, receivable_account_id, notes, is_active, created_by, updated_by
		) SELECT
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24
		WHERE EXISTS (SELECT 1 FROM restaurants r WHERE r.id = $2 AND r.tenant_id = $1)
		  AND ($20::integer IS NULL OR EXISTS (
			SELECT 1 FROM chart_of_accounts a
			WHERE a.id = $20 AND a.tenant_id = $1 AND a.restaurant_id = $2
		  ))
		RETURNING id
	`,
		customer.TenantID,
		customer.RestaurantID,
		code,
		customer.CustomerName,
		customer.CustomerNameAr,
		customer.CustomerType,
		customer.ContactPerson,
		customer.Email,
		customer.Phone,
		customer.Mobile,
		customer.Address,
		customer.AddressAr,
		customer.City,
		customer.Country,
		customer.TaxNumber,
		customer.PaymentTerms,
		customer.CreditLimit,
		customer.CurrentBalance,
		customer.Currency,
		customer.ReceivableAccountID,
		customer.Notes,
		customer.IsActive,
		customer.CreatedBy,
		customer.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	customer.CustomerCode = code
	return id, nil
}

func (r *CustomerAccountRepository) UpdateCustomer(customer *domain.CustomerAccount) error {
	result, err := r.db.Exec(`
		UPDATE customers_accounting
		SET customer_name = $1,
			customer_name_ar = $2,
			customer_type = $3,
			contact_person = $4,
			email = $5,
			phone = $6,
			mobile = $7,
			address = $8,
			address_ar = $9,
			city = $10,
			country = $11,
			tax_number = $12,
			payment_terms = $13,
			credit_limit = $14,
			current_balance = $15,
			currency = $16,
			receivable_account_id = $17,
			notes = $18,
			is_active = $19,
			updated_by = $20,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $21 AND tenant_id = $22 AND restaurant_id = $23
		  AND ($17::integer IS NULL OR EXISTS (
			SELECT 1 FROM chart_of_accounts a
			WHERE a.id = $17 AND a.tenant_id = $22 AND a.restaurant_id = $23
		  ))
	`,
		customer.CustomerName,
		customer.CustomerNameAr,
		customer.CustomerType,
		customer.ContactPerson,
		customer.Email,
		customer.Phone,
		customer.Mobile,
		customer.Address,
		customer.AddressAr,
		customer.City,
		customer.Country,
		customer.TaxNumber,
		customer.PaymentTerms,
		customer.CreditLimit,
		customer.CurrentBalance,
		customer.Currency,
		customer.ReceivableAccountID,
		customer.Notes,
		customer.IsActive,
		customer.UpdatedBy,
		customer.ID,
		customer.TenantID,
		customer.RestaurantID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("customer account not found")
	}
	return nil
}

func (r *CustomerAccountRepository) DeleteCustomer(tenantID, restaurantID, id int) error {
	result, err := r.db.Exec(`DELETE FROM customers_accounting WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3`, tenantID, restaurantID, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("customer account not found")
	}
	return nil
}

func (r *CustomerAccountRepository) GetCustomerStatement(tenantID, restaurantID, customerID int) ([]domain.CustomerStatementLine, error) {
	customer, err := r.GetCustomerByID(tenantID, restaurantID, customerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("customer account not found")
	}
	if customer.ReceivableAccountID == nil {
		return []domain.CustomerStatementLine{}, nil
	}

	rows, err := r.db.Query(`
		SELECT
			gl.entry_date,
			je.entry_number,
			COALESCE(je.reference_number, ''),
			COALESCE(je.description, ''),
			gl.debit_amount,
			gl.credit_amount,
			gl.running_balance
		FROM general_ledger gl
		JOIN journal_entries je ON je.id = gl.journal_entry_id
		WHERE gl.tenant_id = $1
		  AND gl.restaurant_id = $2
		  AND gl.account_id = $3
		ORDER BY gl.entry_date DESC, gl.id DESC
	`, tenantID, restaurantID, *customer.ReceivableAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statement := make([]domain.CustomerStatementLine, 0)
	for rows.Next() {
		var line domain.CustomerStatementLine
		if err := rows.Scan(
			&line.EntryDate,
			&line.EntryNumber,
			&line.ReferenceNumber,
			&line.Description,
			&line.DebitAmount,
			&line.CreditAmount,
			&line.RunningBalance,
		); err != nil {
			return nil, err
		}
		statement = append(statement, line)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statement, nil
}

func (r *CustomerAccountRepository) generateCustomerCode(tenantID, restaurantID int) (string, error) {
	var nextNumber int
	err := r.db.QueryRow(`
		SELECT COALESCE(MAX(CAST(RIGHT(customer_code, 4) AS INT)), 0) + 1
		FROM customers_accounting
		WHERE tenant_id = $1 AND restaurant_id = $2 AND customer_code LIKE 'CUST-%'
	`, tenantID, restaurantID).Scan(&nextNumber)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("CUST-%04d", nextNumber), nil
}

type customerAccountScanner interface {
	Scan(dest ...interface{}) error
}

func scanCustomerAccount(scanner customerAccountScanner) (domain.CustomerAccount, error) {
	var (
		customer            domain.CustomerAccount
		receivableAccountID sql.NullInt64
		createdBy           sql.NullInt64
		updatedBy           sql.NullInt64
	)

	err := scanner.Scan(
		&customer.ID,
		&customer.TenantID,
		&customer.RestaurantID,
		&customer.CustomerCode,
		&customer.CustomerName,
		&customer.CustomerNameAr,
		&customer.CustomerType,
		&customer.ContactPerson,
		&customer.Email,
		&customer.Phone,
		&customer.Mobile,
		&customer.Address,
		&customer.AddressAr,
		&customer.City,
		&customer.Country,
		&customer.TaxNumber,
		&customer.PaymentTerms,
		&customer.CreditLimit,
		&customer.CurrentBalance,
		&customer.Currency,
		&receivableAccountID,
		&customer.Notes,
		&customer.IsActive,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		return domain.CustomerAccount{}, err
	}

	if receivableAccountID.Valid {
		value := int(receivableAccountID.Int64)
		customer.ReceivableAccountID = &value
	}
	if createdBy.Valid {
		value := int(createdBy.Int64)
		customer.CreatedBy = &value
	}
	if updatedBy.Valid {
		value := int(updatedBy.Int64)
		customer.UpdatedBy = &value
	}

	return customer, nil
}
