package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"pos-saas/internal/domain"
)

var (
	ErrSalaryNotFound = errors.New("salary record not found")
	ErrSalaryPaid     = errors.New("paid salary is immutable")
	ErrSalaryState    = errors.New("salary state does not allow payment")
)

// SalaryRepository handles salary/payroll data operations
type SalaryRepository struct {
	db *sql.DB
}

// NewSalaryRepository creates a new salary repository
func NewSalaryRepository(db *sql.DB) *SalaryRepository {
	return &SalaryRepository{db: db}
}

// ListSalariesPaged retrieves salary records for a period using DB pagination.
func (r *SalaryRepository) ListSalariesPaged(ctx context.Context, tenantID, restaurantID int, startDate, endDate time.Time, page, pageSize int) ([]domain.Salary, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset := (page - 1) * pageSize

	const countQuery = `
		SELECT COUNT(*)
		FROM salaries
		WHERE tenant_id = $1 AND restaurant_id = $2
		  AND pay_period_start >= $3 AND pay_period_end <= $4
	`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, tenantID, restaurantID, startDate, endDate).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_id, pay_period_start,
			pay_period_end, payment_date, month, year, base_salary,
			currency, overtime_hours, overtime_rate, overtime_amount,
			bonus, commission, allowances, tips, other_earnings,
			tax, social_insurance, health_insurance, pension,
			loan_deduction, advance_deduction, other_deductions,
			gross_salary, total_deductions, net_salary, days_worked,
			days_absent, total_hours_worked, total_overtime_hours,
			payment_method, status, is_paid, created_at, updated_at
		FROM salaries
		WHERE tenant_id = $1 AND restaurant_id = $2
		  AND pay_period_start >= $3 AND pay_period_end <= $4
		ORDER BY pay_period_end DESC, employee_id ASC
		LIMIT $5 OFFSET $6
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID, startDate, endDate, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var salaries []domain.Salary
	for rows.Next() {
		var sal domain.Salary
		var paymentDate sql.NullTime

		err := rows.Scan(
			&sal.ID, &sal.TenantID, &sal.RestaurantID, &sal.EmployeeID,
			&sal.PayPeriodStart, &sal.PayPeriodEnd, &paymentDate,
			&sal.Month, &sal.Year, &sal.BaseSalary, &sal.Currency,
			&sal.OvertimeHours, &sal.OvertimeRate, &sal.OvertimeAmount,
			&sal.Bonus, &sal.Commission, &sal.Allowances, &sal.Tips,
			&sal.OtherEarnings, &sal.Tax, &sal.SocialInsurance,
			&sal.HealthInsurance, &sal.Pension, &sal.LoanDeduction,
			&sal.AdvanceDeduction, &sal.OtherDeductions, &sal.GrossSalary,
			&sal.TotalDeductions, &sal.NetSalary, &sal.DaysWorked,
			&sal.DaysAbsent, &sal.TotalHoursWorked, &sal.TotalOvertimeHours,
			&sal.PaymentMethod, &sal.Status, &sal.IsPaid,
			&sal.CreatedAt, &sal.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if paymentDate.Valid {
			sal.PaymentDate = &paymentDate.Time
		}

		salaries = append(salaries, sal)
	}

	if salaries == nil {
		salaries = []domain.Salary{}
	}

	return salaries, total, nil
}

// ListSalaries retrieves salary records for a period.
func (r *SalaryRepository) ListSalaries(tenantID, restaurantID int, startDate, endDate time.Time) ([]domain.Salary, error) {
	out, _, err := r.ListSalariesPaged(context.Background(), tenantID, restaurantID, startDate, endDate, 1, 1000000)
	return out, err
}

// GetSalaryByID retrieves a single salary record by ID
func (r *SalaryRepository) GetSalaryByID(tenantID, restaurantID, id int) (*domain.Salary, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_id, pay_period_start,
			pay_period_end, payment_date, month, year, base_salary,
			currency, overtime_hours, overtime_rate, overtime_amount,
			bonus, commission, allowances, tips, other_earnings,
			earnings_details, tax, social_insurance, health_insurance,
			pension, loan_deduction, advance_deduction, other_deductions,
			deductions_details, gross_salary, total_deductions, net_salary,
			days_worked, days_absent, total_hours_worked, total_overtime_hours,
			payment_method, payment_reference, bank_account_number,
			bank_name, status, is_paid, notes, created_at, updated_at
		FROM salaries
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	var sal domain.Salary
	var paymentDate sql.NullTime
	var paymentRef, bankAccount, bankName, notes sql.NullString
	var earningsDetails, deductionsDetails sql.NullString

	err := r.db.QueryRow(query, id, tenantID, restaurantID).Scan(
		&sal.ID, &sal.TenantID, &sal.RestaurantID, &sal.EmployeeID,
		&sal.PayPeriodStart, &sal.PayPeriodEnd, &paymentDate,
		&sal.Month, &sal.Year, &sal.BaseSalary, &sal.Currency,
		&sal.OvertimeHours, &sal.OvertimeRate, &sal.OvertimeAmount,
		&sal.Bonus, &sal.Commission, &sal.Allowances, &sal.Tips,
		&sal.OtherEarnings, &earningsDetails, &sal.Tax,
		&sal.SocialInsurance, &sal.HealthInsurance, &sal.Pension,
		&sal.LoanDeduction, &sal.AdvanceDeduction, &sal.OtherDeductions,
		&deductionsDetails, &sal.GrossSalary, &sal.TotalDeductions,
		&sal.NetSalary, &sal.DaysWorked, &sal.DaysAbsent,
		&sal.TotalHoursWorked, &sal.TotalOvertimeHours, &sal.PaymentMethod,
		&paymentRef, &bankAccount, &bankName, &sal.Status, &sal.IsPaid,
		&notes, &sal.CreatedAt, &sal.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if paymentDate.Valid {
		sal.PaymentDate = &paymentDate.Time
	}
	sal.PaymentReference = paymentRef.String
	sal.BankAccountNumber = bankAccount.String
	sal.BankName = bankName.String
	sal.Notes = notes.String
	if earningsDetails.Valid {
		sal.EarningsDetails = []byte(earningsDetails.String)
	}
	if deductionsDetails.Valid {
		sal.DeductionsDetails = []byte(deductionsDetails.String)
	}

	return &sal, nil
}

// CreateSalary creates a new salary record
func (r *SalaryRepository) CreateSalary(salary *domain.Salary) (int, error) {
	query := `
		INSERT INTO salaries (
			tenant_id, restaurant_id, employee_id, pay_period_start,
			pay_period_end, base_salary, currency, overtime_hours,
			overtime_rate, bonus, commission, allowances, tips,
			other_earnings, tax, social_insurance, health_insurance,
			pension, loan_deduction, advance_deduction, other_deductions,
			days_worked, days_absent, total_hours_worked, total_overtime_hours,
			payment_method, status, created_by
		)
		SELECT $1, $2, e.id, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28
		FROM employees e
		WHERE e.id = $3 AND e.tenant_id = $1 AND e.restaurant_id = $2
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(
		query,
		salary.TenantID, salary.RestaurantID, salary.EmployeeID,
		salary.PayPeriodStart, salary.PayPeriodEnd, salary.BaseSalary,
		salary.Currency, salary.OvertimeHours, salary.OvertimeRate,
		salary.Bonus, salary.Commission, salary.Allowances, salary.Tips,
		salary.OtherEarnings, salary.Tax, salary.SocialInsurance,
		salary.HealthInsurance, salary.Pension, salary.LoanDeduction,
		salary.AdvanceDeduction, salary.OtherDeductions, salary.DaysWorked,
		salary.DaysAbsent, salary.TotalHoursWorked, salary.TotalOvertimeHours,
		salary.PaymentMethod, salary.Status, salary.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateSalary updates an existing salary record
func (r *SalaryRepository) UpdateSalary(salary *domain.Salary) error {
	query := `
		UPDATE salaries
		SET
			base_salary = $1, overtime_hours = $2, overtime_rate = $3,
			bonus = $4, commission = $5, allowances = $6, tips = $7,
			other_earnings = $8, tax = $9, social_insurance = $10,
			health_insurance = $11, pension = $12, loan_deduction = $13,
			advance_deduction = $14, other_deductions = $15, days_worked = $16,
			days_absent = $17, total_hours_worked = $18, total_overtime_hours = $19,
			payment_method = $20, status = $21, updated_by = $22,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $23 AND tenant_id = $24 AND restaurant_id = $25
		  AND is_paid = FALSE AND LOWER(BTRIM(status)) <> 'paid'
	`

	result, err := r.db.Exec(
		query,
		salary.BaseSalary, salary.OvertimeHours, salary.OvertimeRate,
		salary.Bonus, salary.Commission, salary.Allowances, salary.Tips,
		salary.OtherEarnings, salary.Tax, salary.SocialInsurance,
		salary.HealthInsurance, salary.Pension, salary.LoanDeduction,
		salary.AdvanceDeduction, salary.OtherDeductions, salary.DaysWorked,
		salary.DaysAbsent, salary.TotalHoursWorked, salary.TotalOvertimeHours,
		salary.PaymentMethod, salary.Status, salary.UpdatedBy,
		salary.ID, salary.TenantID, salary.RestaurantID,
	)

	if err != nil {
		return fmt.Errorf("update salary: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check salary update result: %w", err)
	}
	if rows == 1 {
		return nil
	}
	return r.salaryMutationFailure(salary.TenantID, salary.RestaurantID, salary.ID)
}

// DeleteUnpaidSalary preserves payroll history once a salary has been paid.
func (r *SalaryRepository) DeleteUnpaidSalary(tenantID, restaurantID, id int) error {
	result, err := r.db.Exec(`
		DELETE FROM salaries
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
		  AND is_paid = FALSE AND LOWER(BTRIM(status)) <> 'paid'
	`, id, tenantID, restaurantID)
	if err != nil {
		return fmt.Errorf("delete salary: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check salary delete result: %w", err)
	}
	if rows == 1 {
		return nil
	}
	return r.salaryMutationFailure(tenantID, restaurantID, id)
}

func (r *SalaryRepository) salaryMutationFailure(tenantID, restaurantID, id int) error {
	var isPaid bool
	var status string
	err := r.db.QueryRow(`
		SELECT is_paid, status FROM salaries
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, id, tenantID, restaurantID).Scan(&isPaid, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrSalaryNotFound
	}
	if err != nil {
		return fmt.Errorf("inspect salary mutation result: %w", err)
	}
	if isPaid || strings.EqualFold(strings.TrimSpace(status), "paid") {
		return ErrSalaryPaid
	}
	return fmt.Errorf("salary mutation did not affect a row")
}

// MarkAsPaidTx locks and transitions one scoped salary. Repeated calls return
// the original paid record without changing its audit fields.
func (r *SalaryRepository) MarkAsPaidTx(ctx context.Context, tx *sql.Tx, tenantID, restaurantID, id, paidBy int, paymentRef string) (*domain.Salary, bool, error) {
	salary, err := r.getSalaryByID(ctx, tx, tenantID, restaurantID, id, true)
	if err != nil {
		return nil, false, err
	}
	if salary == nil {
		return nil, false, ErrSalaryNotFound
	}
	if salary.IsPaid || strings.EqualFold(strings.TrimSpace(salary.Status), "paid") {
		return salary, false, nil
	}
	if strings.EqualFold(strings.TrimSpace(salary.Status), "cancelled") || strings.EqualFold(strings.TrimSpace(salary.Status), "on_hold") {
		return nil, false, ErrSalaryState
	}
	previousStatus := salary.Status
	previousPaid := salary.IsPaid
	previousPaymentReference := salary.PaymentReference

	result, err := tx.ExecContext(ctx, `
		UPDATE salaries
		SET is_paid = TRUE, paid_at = CURRENT_TIMESTAMP, paid_by = NULLIF($1, 0),
		    payment_reference = $2, status = 'paid', updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5
		  AND is_paid = FALSE AND LOWER(BTRIM(status)) <> 'paid'
	`, paidBy, paymentRef, id, tenantID, restaurantID)
	if err != nil {
		return nil, false, fmt.Errorf("mark salary paid: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, false, fmt.Errorf("check salary payment transition: %w", err)
	}
	if rows != 1 {
		return nil, false, ErrSalaryPaid
	}
	salary, err = r.getSalaryByID(ctx, tx, tenantID, restaurantID, id, false)
	if err != nil {
		return nil, false, err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO hr_audit_logs (
		    tenant_id, restaurant_id, entity_type, entity_id, action,
		    action_description, old_values, new_values, changed_fields, user_id
		)
		VALUES (
		    $1, $2, 'salary', $3, 'update', 'Salary marked paid',
		    jsonb_build_object('status', $4::text, 'is_paid', $5::boolean, 'payment_reference', $6::text),
		    jsonb_build_object('status', $7::text, 'is_paid', $8::boolean, 'payment_reference', $9::text, 'paid_at', $10::timestamptz, 'paid_by', $11::integer),
		    ARRAY['status', 'is_paid', 'payment_reference', 'paid_at', 'paid_by'], NULLIF($11::integer, 0)
		)
	`, tenantID, restaurantID, id, previousStatus, previousPaid, previousPaymentReference,
		salary.Status, salary.IsPaid, salary.PaymentReference, salary.PaidAt, paidBy)
	if err != nil {
		return nil, false, fmt.Errorf("audit salary payment transition: %w", err)
	}
	return salary, true, nil
}

func (r *SalaryRepository) getSalaryByID(ctx context.Context, tx *sql.Tx, tenantID, restaurantID, id int, lock bool) (*domain.Salary, error) {
	query := `
		SELECT id, tenant_id, restaurant_id, employee_id, pay_period_start,
		       pay_period_end, payment_date, month, year, base_salary, currency,
		       overtime_hours, overtime_rate, overtime_amount, bonus, commission,
		       allowances, tips, other_earnings, earnings_details, tax,
		       social_insurance, health_insurance, pension, loan_deduction,
		       advance_deduction, other_deductions, deductions_details, gross_salary,
		       total_deductions, net_salary, days_worked, days_absent,
		       total_hours_worked, total_overtime_hours, payment_method,
	       payment_reference, bank_account_number, bank_name, status, is_paid, paid_at, paid_by,
		       notes, created_at, updated_at
		FROM salaries WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`
	if lock {
		query += ` FOR UPDATE`
	}
	row := tx.QueryRowContext(ctx, query, id, tenantID, restaurantID)
	var salary domain.Salary
	var paymentDate sql.NullTime
	var paidAt sql.NullTime
	var paidBy sql.NullInt64
	var paymentRef, bankAccount, bankName, notes sql.NullString
	var earningsDetails, deductionsDetails sql.NullString
	err := row.Scan(
		&salary.ID, &salary.TenantID, &salary.RestaurantID, &salary.EmployeeID,
		&salary.PayPeriodStart, &salary.PayPeriodEnd, &paymentDate,
		&salary.Month, &salary.Year, &salary.BaseSalary, &salary.Currency,
		&salary.OvertimeHours, &salary.OvertimeRate, &salary.OvertimeAmount,
		&salary.Bonus, &salary.Commission, &salary.Allowances, &salary.Tips,
		&salary.OtherEarnings, &earningsDetails, &salary.Tax,
		&salary.SocialInsurance, &salary.HealthInsurance, &salary.Pension,
		&salary.LoanDeduction, &salary.AdvanceDeduction, &salary.OtherDeductions,
		&deductionsDetails, &salary.GrossSalary, &salary.TotalDeductions,
		&salary.NetSalary, &salary.DaysWorked, &salary.DaysAbsent,
		&salary.TotalHoursWorked, &salary.TotalOvertimeHours, &salary.PaymentMethod,
		&paymentRef, &bankAccount, &bankName, &salary.Status, &salary.IsPaid, &paidAt, &paidBy,
		&notes, &salary.CreatedAt, &salary.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load salary in payment transaction: %w", err)
	}
	if paymentDate.Valid {
		salary.PaymentDate = &paymentDate.Time
	}
	if paidAt.Valid {
		salary.PaidAt = &paidAt.Time
	}
	if paidBy.Valid {
		value := int(paidBy.Int64)
		salary.PaidBy = &value
	}
	salary.PaymentReference = paymentRef.String
	salary.BankAccountNumber = bankAccount.String
	salary.BankName = bankName.String
	salary.Notes = notes.String
	if earningsDetails.Valid {
		salary.EarningsDetails = []byte(earningsDetails.String)
	}
	if deductionsDetails.Valid {
		salary.DeductionsDetails = []byte(deductionsDetails.String)
	}
	return &salary, nil
}

// GetEmployeeSalaries retrieves all salary records for an employee
func (r *SalaryRepository) GetEmployeeSalaries(tenantID, restaurantID, employeeID int) ([]domain.Salary, error) {
	query := `
		SELECT
			id, pay_period_start, pay_period_end, payment_date, month, year,
			gross_salary, total_deductions, net_salary, status, is_paid, created_at
		FROM salaries
		WHERE tenant_id = $1 AND restaurant_id = $2 AND employee_id = $3
		ORDER BY pay_period_end DESC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var salaries []domain.Salary
	for rows.Next() {
		var sal domain.Salary
		var paymentDate sql.NullTime

		err := rows.Scan(
			&sal.ID, &sal.PayPeriodStart, &sal.PayPeriodEnd, &paymentDate,
			&sal.Month, &sal.Year, &sal.GrossSalary, &sal.TotalDeductions,
			&sal.NetSalary, &sal.Status, &sal.IsPaid, &sal.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if paymentDate.Valid {
			sal.PaymentDate = &paymentDate.Time
		}

		salaries = append(salaries, sal)
	}

	if salaries == nil {
		salaries = []domain.Salary{}
	}

	return salaries, nil
}
