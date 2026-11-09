package repository

import (
	"database/sql"
	"errors"
	"time"

	"pos-saas/internal/domain"
)

type TrainingRepository struct {
	db *sql.DB
}

func NewTrainingRepository(db *sql.DB) *TrainingRepository {
	return &TrainingRepository{db: db}
}

func (r *TrainingRepository) List(tenantID, restaurantID int) ([]domain.TrainingRecord, error) {
	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, training_title, training_provider, training_type,
			start_date, end_date, status, certificate_number, certificate_expiry, notes,
			created_at, updated_at, created_by, updated_by
		FROM training_records
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY start_date DESC, id DESC
	`, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.TrainingRecord, 0)
	for rows.Next() {
		var item domain.TrainingRecord
		var createdBy, updatedBy sql.NullInt64
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.TrainingTitle, &item.TrainingProvider, &item.TrainingType,
			&item.StartDate, &item.EndDate, &item.Status, &item.CertificateNumber, &item.CertificateExpiry, &item.Notes,
			&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
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

func (r *TrainingRepository) GetByID(tenantID, restaurantID, id int) (*domain.TrainingRecord, error) {
	var item domain.TrainingRecord
	var createdBy, updatedBy sql.NullInt64
	err := r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, training_title, training_provider, training_type,
			start_date, end_date, status, certificate_number, certificate_expiry, notes,
			created_at, updated_at, created_by, updated_by
		FROM training_records
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id).Scan(
		&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.TrainingTitle, &item.TrainingProvider, &item.TrainingType,
		&item.StartDate, &item.EndDate, &item.Status, &item.CertificateNumber, &item.CertificateExpiry, &item.Notes,
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

func (r *TrainingRepository) Create(item *domain.TrainingRecord) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO training_records (
			tenant_id, restaurant_id, employee_id, training_title, training_provider, training_type,
			start_date, end_date, status, certificate_number, certificate_expiry, notes, created_by
		)
		SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
		WHERE EXISTS (
			SELECT 1 FROM employees e
			WHERE e.id = $3 AND e.tenant_id = $1 AND e.restaurant_id = $2
		)
		RETURNING id
	`, item.TenantID, item.RestaurantID, item.EmployeeID, item.TrainingTitle, item.TrainingProvider, item.TrainingType,
		item.StartDate, item.EndDate, item.Status, item.CertificateNumber, item.CertificateExpiry, item.Notes, item.CreatedBy,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("employee is outside tenant and restaurant scope")
	}
	return id, err
}

func (r *TrainingRepository) Update(tenantID, restaurantID int, item *domain.TrainingRecord) error {
	res, err := r.db.Exec(`
		UPDATE training_records
		SET employee_id = $1, training_title = $2, training_provider = $3, training_type = $4,
		    start_date = $5, end_date = $6, status = $7, certificate_number = $8,
		    certificate_expiry = $9, notes = $10, updated_by = $11, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12 AND tenant_id = $13 AND restaurant_id = $14
		  AND EXISTS (
			SELECT 1 FROM employees e
			WHERE e.id = $1 AND e.tenant_id = $13 AND e.restaurant_id = $14
		  )
	`, item.EmployeeID, item.TrainingTitle, item.TrainingProvider, item.TrainingType,
		item.StartDate, item.EndDate, item.Status, item.CertificateNumber,
		item.CertificateExpiry, item.Notes, item.UpdatedBy, item.ID, tenantID, restaurantID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("training record not found")
	}
	return nil
}

func (r *TrainingRepository) Delete(tenantID, restaurantID, id int) error {
	res, err := r.db.Exec(`
		DELETE FROM training_records
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
		return errors.New("training record not found")
	}
	return nil
}

func (r *TrainingRepository) GetExpiringCertificates(tenantID, restaurantID int, withinDays int) ([]domain.TrainingRecord, error) {
	if withinDays <= 0 {
		withinDays = 30
	}
	cutoff := time.Now().UTC().AddDate(0, 0, withinDays)
	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, training_title, training_provider, training_type,
			start_date, end_date, status, certificate_number, certificate_expiry, notes,
			created_at, updated_at, created_by, updated_by
		FROM training_records
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND certificate_expiry IS NOT NULL
		  AND certificate_expiry <= $3
		ORDER BY certificate_expiry ASC, id ASC
	`, tenantID, restaurantID, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.TrainingRecord, 0)
	for rows.Next() {
		var item domain.TrainingRecord
		var createdBy, updatedBy sql.NullInt64
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.TrainingTitle, &item.TrainingProvider, &item.TrainingType,
			&item.StartDate, &item.EndDate, &item.Status, &item.CertificateNumber, &item.CertificateExpiry, &item.Notes,
			&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
