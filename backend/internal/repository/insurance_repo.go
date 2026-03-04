package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pos-saas/internal/domain"
)

type InsuranceRepository struct {
	db *sql.DB
}

func NewInsuranceRepository(db *sql.DB) *InsuranceRepository {
	return &InsuranceRepository{db: db}
}

func (r *InsuranceRepository) ListPaged(ctx context.Context, tenantID, restaurantID, page, pageSize int) ([]domain.InsuranceTracking, int, error) {
	var total int
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(1) FROM insurance_tracking WHERE tenant_id = $1 AND restaurant_id = $2`,
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
			id, tenant_id, restaurant_id, employee_id, policy_type, provider_name, policy_number,
			coverage_start, coverage_end, monthly_premium, status, notes,
			created_at, updated_at, created_by, updated_by
		FROM insurance_tracking
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY coverage_end ASC, id DESC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := make([]domain.InsuranceTracking, 0, pageSize)
	for rows.Next() {
		var item domain.InsuranceTracking
		var createdBy, updatedBy sql.NullInt64
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.PolicyType, &item.ProviderName, &item.PolicyNumber,
			&item.CoverageStart, &item.CoverageEnd, &item.MonthlyPremium, &item.Status, &item.Notes,
			&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
		); err != nil {
			return nil, 0, err
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

func (r *InsuranceRepository) List(tenantID, restaurantID int) ([]domain.InsuranceTracking, error) {
	out, _, err := r.ListPaged(context.Background(), tenantID, restaurantID, 1, 1000000)
	return out, err
}

func (r *InsuranceRepository) GetByID(tenantID, restaurantID, id int) (*domain.InsuranceTracking, error) {
	var item domain.InsuranceTracking
	var createdBy, updatedBy sql.NullInt64
	err := r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, policy_type, provider_name, policy_number,
			coverage_start, coverage_end, monthly_premium, status, notes,
			created_at, updated_at, created_by, updated_by
		FROM insurance_tracking
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id).Scan(
		&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.PolicyType, &item.ProviderName, &item.PolicyNumber,
		&item.CoverageStart, &item.CoverageEnd, &item.MonthlyPremium, &item.Status, &item.Notes,
		&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
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

func (r *InsuranceRepository) Create(item *domain.InsuranceTracking) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO insurance_tracking (
			tenant_id, restaurant_id, employee_id, policy_type, provider_name, policy_number,
			coverage_start, coverage_end, monthly_premium, status, notes, created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id
	`, item.TenantID, item.RestaurantID, item.EmployeeID, item.PolicyType, item.ProviderName, item.PolicyNumber,
		item.CoverageStart, item.CoverageEnd, item.MonthlyPremium, item.Status, item.Notes, item.CreatedBy,
	).Scan(&id)
	return id, err
}

func (r *InsuranceRepository) Update(tenantID, restaurantID int, item *domain.InsuranceTracking) error {
	res, err := r.db.Exec(`
		UPDATE insurance_tracking
		SET employee_id = $1, policy_type = $2, provider_name = $3, policy_number = $4,
		    coverage_start = $5, coverage_end = $6, monthly_premium = $7, status = $8,
		    notes = $9, updated_by = $10, updated_at = CURRENT_TIMESTAMP
		WHERE id = $11 AND tenant_id = $12 AND restaurant_id = $13
	`, item.EmployeeID, item.PolicyType, item.ProviderName, item.PolicyNumber,
		item.CoverageStart, item.CoverageEnd, item.MonthlyPremium, item.Status,
		item.Notes, item.UpdatedBy, item.ID, tenantID, restaurantID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("insurance record not found")
	}
	return nil
}

func (r *InsuranceRepository) Delete(tenantID, restaurantID, id int) error {
	res, err := r.db.Exec(`
		DELETE FROM insurance_tracking
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
		return errors.New("insurance record not found")
	}
	return nil
}

func (r *InsuranceRepository) GetExpiring(tenantID, restaurantID int, withinDays int) ([]domain.InsuranceTracking, error) {
	if withinDays <= 0 {
		withinDays = 30
	}
	cutoff := time.Now().UTC().AddDate(0, 0, withinDays)
	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, policy_type, provider_name, policy_number,
			coverage_start, coverage_end, monthly_premium, status, notes,
			created_at, updated_at, created_by, updated_by
		FROM insurance_tracking
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND status = 'active'
		  AND coverage_end <= $3
		ORDER BY coverage_end ASC, id ASC
	`, tenantID, restaurantID, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.InsuranceTracking, 0)
	for rows.Next() {
		var item domain.InsuranceTracking
		var createdBy, updatedBy sql.NullInt64
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.PolicyType, &item.ProviderName, &item.PolicyNumber,
			&item.CoverageStart, &item.CoverageEnd, &item.MonthlyPremium, &item.Status, &item.Notes,
			&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
