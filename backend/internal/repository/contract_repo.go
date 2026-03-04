package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pos-saas/internal/domain"
)

type ContractRepository struct {
	db *sql.DB
}

func NewContractRepository(db *sql.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

func (r *ContractRepository) ListPaged(ctx context.Context, tenantID, restaurantID, page, pageSize int) ([]domain.EmploymentContract, int, error) {
	var total int
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(1) FROM employment_contracts WHERE tenant_id = $1 AND restaurant_id = $2`,
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
			id, tenant_id, restaurant_id, employee_id, contract_type, contract_number, start_date, end_date,
			salary_amount, currency, status, renewed_from_id, notes,
			created_at, updated_at, created_by, updated_by
		FROM employment_contracts
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY start_date DESC, id DESC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := make([]domain.EmploymentContract, 0, pageSize)
	for rows.Next() {
		var item domain.EmploymentContract
		var renewedFromID, createdBy, updatedBy sql.NullInt64
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.ContractType, &item.ContractNumber, &item.StartDate, &item.EndDate,
			&item.SalaryAmount, &item.Currency, &item.Status, &renewedFromID, &item.Notes,
			&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
		); err != nil {
			return nil, 0, err
		}
		if renewedFromID.Valid {
			v := int(renewedFromID.Int64)
			item.RenewedFromID = &v
		}
		if createdBy.Valid {
			v := int(createdBy.Int64)
			item.CreatedBy = &v
		}
		if updatedBy.Valid {
			v := int(updatedBy.Int64)
			item.UpdatedBy = &v
		}
		out = append(out, item)
	}
	return out, total, rows.Err()
}

func (r *ContractRepository) List(tenantID, restaurantID int) ([]domain.EmploymentContract, error) {
	out, _, err := r.ListPaged(context.Background(), tenantID, restaurantID, 1, 1000000)
	return out, err
}

func (r *ContractRepository) GetByID(tenantID, restaurantID, id int) (*domain.EmploymentContract, error) {
	var item domain.EmploymentContract
	var renewedFromID, createdBy, updatedBy sql.NullInt64
	err := r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, contract_type, contract_number, start_date, end_date,
			salary_amount, currency, status, renewed_from_id, notes,
			created_at, updated_at, created_by, updated_by
		FROM employment_contracts
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id).Scan(
		&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.ContractType, &item.ContractNumber, &item.StartDate, &item.EndDate,
		&item.SalaryAmount, &item.Currency, &item.Status, &renewedFromID, &item.Notes,
		&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if renewedFromID.Valid {
		v := int(renewedFromID.Int64)
		item.RenewedFromID = &v
	}
	if createdBy.Valid {
		v := int(createdBy.Int64)
		item.CreatedBy = &v
	}
	if updatedBy.Valid {
		v := int(updatedBy.Int64)
		item.UpdatedBy = &v
	}
	return &item, nil
}

func (r *ContractRepository) Create(item *domain.EmploymentContract) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO employment_contracts (
			tenant_id, restaurant_id, employee_id, contract_type, contract_number, start_date, end_date,
			salary_amount, currency, status, renewed_from_id, notes, created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id
	`, item.TenantID, item.RestaurantID, item.EmployeeID, item.ContractType, item.ContractNumber, item.StartDate, item.EndDate,
		item.SalaryAmount, item.Currency, item.Status, item.RenewedFromID, item.Notes, item.CreatedBy,
	).Scan(&id)
	return id, err
}

func (r *ContractRepository) Update(tenantID, restaurantID int, item *domain.EmploymentContract) error {
	res, err := r.db.Exec(`
		UPDATE employment_contracts
		SET employee_id = $1, contract_type = $2, contract_number = $3, start_date = $4, end_date = $5,
		    salary_amount = $6, currency = $7, status = $8, notes = $9, updated_by = $10, updated_at = CURRENT_TIMESTAMP
		WHERE id = $11 AND tenant_id = $12 AND restaurant_id = $13
	`, item.EmployeeID, item.ContractType, item.ContractNumber, item.StartDate, item.EndDate,
		item.SalaryAmount, item.Currency, item.Status, item.Notes, item.UpdatedBy, item.ID, tenantID, restaurantID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("contract not found")
	}
	return nil
}

func (r *ContractRepository) Delete(tenantID, restaurantID, id int) error {
	res, err := r.db.Exec(`
		DELETE FROM employment_contracts
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, id, tenantID, restaurantID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("contract not found")
	}
	return nil
}

func (r *ContractRepository) Renew(tenantID, restaurantID, id int, startDate time.Time, endDate *time.Time, salaryAmount float64, currency string, updatedBy int) (int, error) {
	existing, err := r.GetByID(tenantID, restaurantID, id)
	if err != nil {
		return 0, err
	}
	if existing == nil {
		return 0, errors.New("contract not found")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		UPDATE employment_contracts
		SET status = 'renewed', updated_by = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, updatedBy, id, tenantID, restaurantID); err != nil {
		return 0, err
	}

	var newID int
	err = tx.QueryRow(`
		INSERT INTO employment_contracts (
			tenant_id, restaurant_id, employee_id, contract_type, contract_number, start_date, end_date,
			salary_amount, currency, status, renewed_from_id, notes, created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,'active',$10,$11,$12)
		RETURNING id
	`,
		tenantID, restaurantID, existing.EmployeeID, existing.ContractType,
		existing.ContractNumber+"-R", startDate, endDate, salaryAmount, currency, id, existing.Notes, updatedBy,
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return newID, nil
}

func (r *ContractRepository) GetExpiring(tenantID, restaurantID int, withinDays int) ([]domain.EmploymentContract, error) {
	if withinDays <= 0 {
		withinDays = 30
	}
	cutoff := time.Now().UTC().AddDate(0, 0, withinDays)
	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, contract_type, contract_number, start_date, end_date,
			salary_amount, currency, status, renewed_from_id, notes,
			created_at, updated_at, created_by, updated_by
		FROM employment_contracts
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND status = 'active'
		  AND end_date IS NOT NULL
		  AND end_date <= $3
		ORDER BY end_date ASC, id ASC
	`, tenantID, restaurantID, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.EmploymentContract, 0)
	for rows.Next() {
		var item domain.EmploymentContract
		var renewedFromID, createdBy, updatedBy sql.NullInt64
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.ContractType, &item.ContractNumber, &item.StartDate, &item.EndDate,
			&item.SalaryAmount, &item.Currency, &item.Status, &renewedFromID, &item.Notes,
			&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
