package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pos-saas/internal/domain"
)

type SupplierRepository struct {
	db *sql.DB
}

func NewSupplierRepository(db *sql.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) ListSuppliersPaged(ctx context.Context, tenantID, restaurantID, page, pageSize int) ([]domain.Supplier, int, error) {
	var total int
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(1) FROM suppliers WHERE tenant_id = $1 AND restaurant_id = $2`,
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
			id, tenant_id, restaurant_id, supplier_code, company_name,
			COALESCE(company_name_ar, ''), COALESCE(contact_person, ''),
			COALESCE(email, ''), COALESCE(phone, ''), COALESCE(mobile, ''),
			COALESCE(website, ''), COALESCE(address, ''), COALESCE(address_ar, ''),
			COALESCE(city, ''), COALESCE(state, ''), COALESCE(postal_code, ''),
			COALESCE(country, 'Egypt'), COALESCE(tax_number, ''),
			COALESCE(payment_terms, 'net_30'), COALESCE(credit_limit, 0),
			COALESCE(current_balance, 0), COALESCE(currency, 'EGP'),
			payable_account_id, COALESCE(category, ''), COALESCE(bank_name, ''),
			COALESCE(bank_account, ''), COALESCE(iban, ''), COALESCE(notes, ''),
			is_active, created_at, updated_at, created_by, updated_by
		FROM suppliers
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY company_name ASC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	suppliers := make([]domain.Supplier, 0, pageSize)
	for rows.Next() {
		supplier, err := scanSupplier(rows)
		if err != nil {
			return nil, 0, err
		}
		suppliers = append(suppliers, supplier)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

func (r *SupplierRepository) ListSuppliers(tenantID, restaurantID int) ([]domain.Supplier, error) {
	suppliers, _, err := r.ListSuppliersPaged(context.Background(), tenantID, restaurantID, 1, 1000000)
	return suppliers, err
}

func (r *SupplierRepository) GetSupplierByID(tenantID, restaurantID, id int) (*domain.Supplier, error) {
	supplier, err := scanSupplier(r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, supplier_code, company_name,
			COALESCE(company_name_ar, ''), COALESCE(contact_person, ''),
			COALESCE(email, ''), COALESCE(phone, ''), COALESCE(mobile, ''),
			COALESCE(website, ''), COALESCE(address, ''), COALESCE(address_ar, ''),
			COALESCE(city, ''), COALESCE(state, ''), COALESCE(postal_code, ''),
			COALESCE(country, 'Egypt'), COALESCE(tax_number, ''),
			COALESCE(payment_terms, 'net_30'), COALESCE(credit_limit, 0),
			COALESCE(current_balance, 0), COALESCE(currency, 'EGP'),
			payable_account_id, COALESCE(category, ''), COALESCE(bank_name, ''),
			COALESCE(bank_account, ''), COALESCE(iban, ''), COALESCE(notes, ''),
			is_active, created_at, updated_at, created_by, updated_by
		FROM suppliers
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

func (r *SupplierRepository) CreateSupplier(supplier *domain.Supplier) (int, error) {
	code, err := r.generateSupplierCode(supplier.TenantID, supplier.RestaurantID)
	if err != nil {
		return 0, err
	}

	var id int
	err = r.db.QueryRow(`
		INSERT INTO suppliers (
			tenant_id, restaurant_id, supplier_code, company_name, company_name_ar,
			contact_person, email, phone, mobile, website, address, address_ar,
			city, state, postal_code, country, tax_number, payment_terms,
			credit_limit, current_balance, currency, payable_account_id, category,
			bank_name, bank_account, iban, notes, is_active, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23,
			$24, $25, $26, $27, $28, $29, $30
		)
		RETURNING id
	`,
		supplier.TenantID,
		supplier.RestaurantID,
		code,
		supplier.CompanyName,
		supplier.CompanyNameAr,
		supplier.ContactPerson,
		supplier.Email,
		supplier.Phone,
		supplier.Mobile,
		supplier.Website,
		supplier.Address,
		supplier.AddressAr,
		supplier.City,
		supplier.State,
		supplier.PostalCode,
		supplier.Country,
		supplier.TaxNumber,
		supplier.PaymentTerms,
		supplier.CreditLimit,
		supplier.CurrentBalance,
		supplier.Currency,
		supplier.PayableAccountID,
		supplier.Category,
		supplier.BankName,
		supplier.BankAccount,
		supplier.IBAN,
		supplier.Notes,
		supplier.IsActive,
		supplier.CreatedBy,
		supplier.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	supplier.SupplierCode = code
	return id, nil
}

func (r *SupplierRepository) UpdateSupplier(supplier *domain.Supplier) error {
	result, err := r.db.Exec(`
		UPDATE suppliers
		SET company_name = $1,
			company_name_ar = $2,
			contact_person = $3,
			email = $4,
			phone = $5,
			mobile = $6,
			website = $7,
			address = $8,
			address_ar = $9,
			city = $10,
			state = $11,
			postal_code = $12,
			country = $13,
			tax_number = $14,
			payment_terms = $15,
			credit_limit = $16,
			current_balance = $17,
			currency = $18,
			payable_account_id = $19,
			category = $20,
			bank_name = $21,
			bank_account = $22,
			iban = $23,
			notes = $24,
			is_active = $25,
			updated_by = $26,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $27 AND tenant_id = $28 AND restaurant_id = $29
	`,
		supplier.CompanyName,
		supplier.CompanyNameAr,
		supplier.ContactPerson,
		supplier.Email,
		supplier.Phone,
		supplier.Mobile,
		supplier.Website,
		supplier.Address,
		supplier.AddressAr,
		supplier.City,
		supplier.State,
		supplier.PostalCode,
		supplier.Country,
		supplier.TaxNumber,
		supplier.PaymentTerms,
		supplier.CreditLimit,
		supplier.CurrentBalance,
		supplier.Currency,
		supplier.PayableAccountID,
		supplier.Category,
		supplier.BankName,
		supplier.BankAccount,
		supplier.IBAN,
		supplier.Notes,
		supplier.IsActive,
		supplier.UpdatedBy,
		supplier.ID,
		supplier.TenantID,
		supplier.RestaurantID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("supplier not found")
	}
	return nil
}

func (r *SupplierRepository) DeleteSupplier(tenantID, restaurantID, id int) error {
	result, err := r.db.Exec(`DELETE FROM suppliers WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3`, tenantID, restaurantID, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("supplier not found")
	}
	return nil
}

func (r *SupplierRepository) generateSupplierCode(tenantID, restaurantID int) (string, error) {
	var nextNumber int
	err := r.db.QueryRow(`
		SELECT COALESCE(MAX(CAST(RIGHT(supplier_code, 4) AS INT)), 0) + 1
		FROM suppliers
		WHERE tenant_id = $1 AND restaurant_id = $2 AND supplier_code LIKE 'SUP-%'
	`, tenantID, restaurantID).Scan(&nextNumber)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("SUP-%04d", nextNumber), nil
}

type supplierScanner interface {
	Scan(dest ...interface{}) error
}

func scanSupplier(scanner supplierScanner) (domain.Supplier, error) {
	var (
		supplier         domain.Supplier
		payableAccountID sql.NullInt64
		createdBy        sql.NullInt64
		updatedBy        sql.NullInt64
	)

	err := scanner.Scan(
		&supplier.ID,
		&supplier.TenantID,
		&supplier.RestaurantID,
		&supplier.SupplierCode,
		&supplier.CompanyName,
		&supplier.CompanyNameAr,
		&supplier.ContactPerson,
		&supplier.Email,
		&supplier.Phone,
		&supplier.Mobile,
		&supplier.Website,
		&supplier.Address,
		&supplier.AddressAr,
		&supplier.City,
		&supplier.State,
		&supplier.PostalCode,
		&supplier.Country,
		&supplier.TaxNumber,
		&supplier.PaymentTerms,
		&supplier.CreditLimit,
		&supplier.CurrentBalance,
		&supplier.Currency,
		&payableAccountID,
		&supplier.Category,
		&supplier.BankName,
		&supplier.BankAccount,
		&supplier.IBAN,
		&supplier.Notes,
		&supplier.IsActive,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		return domain.Supplier{}, err
	}

	if payableAccountID.Valid {
		value := int(payableAccountID.Int64)
		supplier.PayableAccountID = &value
	}
	if createdBy.Valid {
		value := int(createdBy.Int64)
		supplier.CreatedBy = &value
	}
	if updatedBy.Valid {
		value := int(updatedBy.Int64)
		supplier.UpdatedBy = &value
	}

	return supplier, nil
}
