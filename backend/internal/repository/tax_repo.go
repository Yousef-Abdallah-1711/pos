package repository

import (
	"database/sql"
	"errors"
	"time"
)

type TaxPeriod struct {
	ID           int        `json:"id"`
	TenantID     int        `json:"tenant_id"`
	RestaurantID int        `json:"restaurant_id"`
	PeriodName   string     `json:"period_name"`
	PeriodStart  time.Time  `json:"period_start"`
	PeriodEnd    time.Time  `json:"period_end"`
	OutputVAT    float64    `json:"output_vat"`
	InputVAT     float64    `json:"input_vat"`
	NetVAT       float64    `json:"net_vat"`
	Status       string     `json:"status"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	ClosedBy     *int       `json:"closed_by,omitempty"`
	Notes        string     `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type TaxSummary struct {
	DateFrom  string  `json:"date_from"`
	DateTo    string  `json:"date_to"`
	OutputVAT float64 `json:"output_vat"`
	InputVAT  float64 `json:"input_vat"`
	NetVAT    float64 `json:"net_vat"`
}

type TaxRepository struct {
	db *sql.DB
}

type TaxPeriodInput struct {
	PeriodName  string
	PeriodStart time.Time
	PeriodEnd   time.Time
	Notes       string
}

func NewTaxRepository(db *sql.DB) *TaxRepository {
	return &TaxRepository{db: db}
}

func (r *TaxRepository) ListTaxPeriods(tenantID, restaurantID int) ([]TaxPeriod, error) {
	items, _, err := r.ListTaxPeriodsPaged(tenantID, restaurantID, 1, 100000)
	return items, err
}

func (r *TaxRepository) ListTaxPeriodsPaged(tenantID, restaurantID, page, pageSize int) ([]TaxPeriod, int, error) {
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
		FROM tax_periods
		WHERE tenant_id = $1 AND restaurant_id = $2
	`, tenantID, restaurantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, period_name, period_start, period_end,
			output_vat, input_vat, net_vat, status, closed_at, closed_by, notes,
			created_at, updated_at
		FROM tax_periods
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY period_start DESC, id DESC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]TaxPeriod, 0)
	for rows.Next() {
		var (
			item     TaxPeriod
			closedBy sql.NullInt64
		)
		if err := rows.Scan(
			&item.ID,
			&item.TenantID,
			&item.RestaurantID,
			&item.PeriodName,
			&item.PeriodStart,
			&item.PeriodEnd,
			&item.OutputVAT,
			&item.InputVAT,
			&item.NetVAT,
			&item.Status,
			&item.ClosedAt,
			&closedBy,
			&item.Notes,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		if closedBy.Valid {
			v := int(closedBy.Int64)
			item.ClosedBy = &v
		}
		out = append(out, item)
	}
	return out, total, rows.Err()
}

func (r *TaxRepository) GetTaxPeriodByID(tenantID, restaurantID, id int) (TaxPeriod, error) {
	var (
		item     TaxPeriod
		closedBy sql.NullInt64
	)
	err := r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, period_name, period_start, period_end,
			output_vat, input_vat, net_vat, status, closed_at, closed_by, notes,
			created_at, updated_at
		FROM tax_periods
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id).Scan(
		&item.ID,
		&item.TenantID,
		&item.RestaurantID,
		&item.PeriodName,
		&item.PeriodStart,
		&item.PeriodEnd,
		&item.OutputVAT,
		&item.InputVAT,
		&item.NetVAT,
		&item.Status,
		&item.ClosedAt,
		&closedBy,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return item, err
	}
	if closedBy.Valid {
		v := int(closedBy.Int64)
		item.ClosedBy = &v
	}
	return item, nil
}

func (r *TaxRepository) CreateTaxPeriod(tenantID, restaurantID int, input TaxPeriodInput) (TaxPeriod, error) {
	summary, err := r.GetVATSummary(tenantID, restaurantID, input.PeriodStart, input.PeriodEnd)
	if err != nil {
		return TaxPeriod{}, err
	}

	var created TaxPeriod
	err = r.db.QueryRow(`
		INSERT INTO tax_periods (
			tenant_id, restaurant_id, period_name, period_start, period_end,
			output_vat, input_vat, net_vat, status, notes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'open', $9)
		RETURNING id, tenant_id, restaurant_id, period_name, period_start, period_end,
		          output_vat, input_vat, net_vat, status, closed_at, closed_by, notes,
				  created_at, updated_at
	`, tenantID, restaurantID, input.PeriodName, input.PeriodStart, input.PeriodEnd,
		summary.OutputVAT, summary.InputVAT, summary.NetVAT, input.Notes,
	).Scan(
		&created.ID,
		&created.TenantID,
		&created.RestaurantID,
		&created.PeriodName,
		&created.PeriodStart,
		&created.PeriodEnd,
		&created.OutputVAT,
		&created.InputVAT,
		&created.NetVAT,
		&created.Status,
		&created.ClosedAt,
		new(sql.NullInt64),
		&created.Notes,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	return created, err
}

func (r *TaxRepository) UpdateTaxPeriod(tenantID, restaurantID, id int, input TaxPeriodInput) (TaxPeriod, error) {
	summary, err := r.GetVATSummary(tenantID, restaurantID, input.PeriodStart, input.PeriodEnd)
	if err != nil {
		return TaxPeriod{}, err
	}
	res, err := r.db.Exec(`
		UPDATE tax_periods
		SET period_name = $1,
			period_start = $2,
			period_end = $3,
			output_vat = $4,
			input_vat = $5,
			net_vat = $6,
			notes = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		  AND tenant_id = $9
		  AND restaurant_id = $10
		  AND status = 'open'
	`, input.PeriodName, input.PeriodStart, input.PeriodEnd, summary.OutputVAT, summary.InputVAT, summary.NetVAT, input.Notes, id, tenantID, restaurantID)
	if err != nil {
		return TaxPeriod{}, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return TaxPeriod{}, err
	}
	if affected == 0 {
		return TaxPeriod{}, errors.New("tax period not found or not editable")
	}
	return r.GetTaxPeriodByID(tenantID, restaurantID, id)
}

func (r *TaxRepository) DeleteTaxPeriod(tenantID, restaurantID, id int) error {
	res, err := r.db.Exec(`
		DELETE FROM tax_periods
		WHERE id = $1
		  AND tenant_id = $2
		  AND restaurant_id = $3
		  AND status = 'open'
	`, id, tenantID, restaurantID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("tax period not found or not deletable")
	}
	return nil
}

func (r *TaxRepository) GetVATSummary(tenantID, restaurantID int, from, to time.Time) (TaxSummary, error) {
	var out TaxSummary
	out.DateFrom = from.Format("2006-01-02")
	out.DateTo = to.Format("2006-01-02")

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(tax_amount), 0)
		FROM invoices_customer
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND status <> 'void'
		  AND invoice_date >= $3
		  AND invoice_date <= $4
	`, tenantID, restaurantID, from, to).Scan(&out.OutputVAT); err != nil {
		return out, err
	}

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(tax_amount), 0)
		FROM invoices_supplier
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND status <> 'void'
		  AND invoice_date >= $3
		  AND invoice_date <= $4
	`, tenantID, restaurantID, from, to).Scan(&out.InputVAT); err != nil {
		return out, err
	}

	out.NetVAT = out.OutputVAT - out.InputVAT
	return out, nil
}

func (r *TaxRepository) CloseTaxPeriod(tenantID, restaurantID, periodID, closedBy int) error {
	res, err := r.db.Exec(`
		UPDATE tax_periods
		SET status = 'closed',
			closed_at = CURRENT_TIMESTAMP,
			closed_by = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		  AND tenant_id = $3
		  AND restaurant_id = $4
		  AND status = 'open'
	`, closedBy, periodID, tenantID, restaurantID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("tax period not found or already closed")
	}

	return nil
}
