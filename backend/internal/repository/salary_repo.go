package repository

import (
	"database/sql"
	"time"

	"pos-saas/internal/domain"
)

// SalaryRepository handles salary/payroll data operations
type SalaryRepository struct {
	db *sql.DB
}

// NewSalaryRepository creates a new salary repository
func NewSalaryRepository(db *sql.DB) *SalaryRepository {
	return &SalaryRepository{db: db}
}

// ListSalaries retrieves salary records for a period
func (r *SalaryRepository) ListSalaries(tenantID, restaurantID int, startDate, endDate time.Time) ([]domain.Salary, error) {
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
	`

	rows, err := r.db.Query(query, tenantID, restaurantID, startDate, endDate)
	if err != nil {
		return nil, err
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
	`

	_, err := r.db.Exec(
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

	return err
}

// MarkAsPaid marks a salary record as paid
func (r *SalaryRepository) MarkAsPaid(tenantID, restaurantID, id int, paidBy int, paymentRef string) error {
	query := `
		UPDATE salaries
		SET
			is_paid = true, paid_at = CURRENT_TIMESTAMP, paid_by = $1,
			payment_reference = $2, status = 'paid', updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5
	`

	_, err := r.db.Exec(query, paidBy, paymentRef, id, tenantID, restaurantID)
	return err
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
