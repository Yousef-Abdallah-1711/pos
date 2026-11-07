package repository

import (
	"database/sql"
	"errors"
	"time"

	"pos-saas/internal/domain"
)

type ResignationRepository struct {
	db *sql.DB
}

func NewResignationRepository(db *sql.DB) *ResignationRepository {
	return &ResignationRepository{db: db}
}

func (r *ResignationRepository) List(tenantID, restaurantID int) ([]domain.Resignation, error) {
	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, resignation_date, last_working_date,
			reason, status, notice_period_days, pending_salary, unused_leave_balance, settlement_amount,
			acknowledged_at, completed_at, created_at, updated_at, created_by, updated_by
		FROM resignations
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY resignation_date DESC, id DESC
	`, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.Resignation, 0)
	for rows.Next() {
		var item domain.Resignation
		var createdBy, updatedBy sql.NullInt64
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.ResignationDate, &item.LastWorkingDate,
			&item.Reason, &item.Status, &item.NoticePeriodDays, &item.PendingSalary, &item.UnusedLeaveBalance, &item.SettlementAmount,
			&item.AcknowledgedAt, &item.CompletedAt, &item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
		); err != nil {
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
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *ResignationRepository) GetByID(tenantID, restaurantID, id int) (*domain.Resignation, error) {
	var item domain.Resignation
	var createdBy, updatedBy sql.NullInt64
	err := r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, resignation_date, last_working_date,
			reason, status, notice_period_days, pending_salary, unused_leave_balance, settlement_amount,
			acknowledged_at, completed_at, created_at, updated_at, created_by, updated_by
		FROM resignations
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id).Scan(
		&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.ResignationDate, &item.LastWorkingDate,
		&item.Reason, &item.Status, &item.NoticePeriodDays, &item.PendingSalary, &item.UnusedLeaveBalance, &item.SettlementAmount,
		&item.AcknowledgedAt, &item.CompletedAt, &item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
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

func (r *ResignationRepository) Create(item *domain.Resignation) (int, error) {
	if item.Status == "" {
		item.Status = "submitted"
	}
	var id int
	err := r.db.QueryRow(`
		INSERT INTO resignations (
			tenant_id, restaurant_id, employee_id, resignation_date, last_working_date, reason, status,
			notice_period_days, pending_salary, unused_leave_balance, settlement_amount, created_by, updated_by
		)
		SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
		WHERE EXISTS (
			SELECT 1 FROM employees e
			WHERE e.id = $3 AND e.tenant_id = $1 AND e.restaurant_id = $2
		)
		RETURNING id
	`, item.TenantID, item.RestaurantID, item.EmployeeID, item.ResignationDate, item.LastWorkingDate, item.Reason, item.Status,
		item.NoticePeriodDays, item.PendingSalary, item.UnusedLeaveBalance, item.SettlementAmount, item.CreatedBy, item.UpdatedBy,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("employee is outside tenant and restaurant scope")
	}
	return id, err
}

func (r *ResignationRepository) Update(tenantID, restaurantID int, item *domain.Resignation) error {
	res, err := r.db.Exec(`
		UPDATE resignations
		SET resignation_date = $1, last_working_date = $2, reason = $3, status = $4,
		    notice_period_days = $5, pending_salary = $6, unused_leave_balance = $7,
		    settlement_amount = $8, updated_by = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10 AND tenant_id = $11 AND restaurant_id = $12
		  AND EXISTS (
			SELECT 1 FROM employees e
			WHERE e.id = $1 AND e.tenant_id = $11 AND e.restaurant_id = $12
		  )
	`, item.ResignationDate, item.LastWorkingDate, item.Reason, item.Status,
		item.NoticePeriodDays, item.PendingSalary, item.UnusedLeaveBalance, item.SettlementAmount,
		item.UpdatedBy, item.ID, tenantID, restaurantID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("resignation not found")
	}
	return nil
}

func (r *ResignationRepository) Delete(tenantID, restaurantID, id int) error {
	res, err := r.db.Exec(`
		DELETE FROM resignations
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
		return errors.New("resignation not found")
	}
	return nil
}

func (r *ResignationRepository) Acknowledge(tenantID, restaurantID, id int) error {
	res, err := r.db.Exec(`
		UPDATE resignations
		SET status = 'acknowledged', acknowledged_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
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
		return errors.New("resignation not found")
	}
	return nil
}

func (r *ResignationRepository) Complete(tenantID, restaurantID, id int, settlementAmount float64) error {
	res, err := r.db.Exec(`
		UPDATE resignations
		SET status = 'completed',
		    settlement_amount = $1,
		    completed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, settlementAmount, id, tenantID, restaurantID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("resignation not found")
	}
	return nil
}

func (r *ResignationRepository) CalculateSettlement(item *domain.Resignation) float64 {
	// Base formula: pending salary + unused leave payout.
	return item.PendingSalary + item.UnusedLeaveBalance
}

func calculateNoticePeriodDays(resignationDate, lastWorkingDate time.Time) int {
	return int(lastWorkingDate.Sub(resignationDate).Hours() / 24)
}
