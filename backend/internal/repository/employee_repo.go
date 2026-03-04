package repository

import (
	"context"
	"database/sql"
	"log"
	"pos-saas/internal/domain"
)

// EmployeeRepository handles employee data operations
type EmployeeRepository struct {
	db *sql.DB
}

// NewEmployeeRepository creates a new employee repository
func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// ListEmployeesPaged retrieves active employees with database pagination.
func (r *EmployeeRepository) ListEmployeesPaged(ctx context.Context, tenantID, restaurantID, page, pageSize int) ([]domain.Employee, int, error) {
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
		FROM employees
		WHERE tenant_id = $1 AND restaurant_id = $2 AND is_active = true
	`

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, tenantID, restaurantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_code, first_name, last_name,
			first_name_ar, last_name_ar, email, phone, date_of_birth, gender,
			national_id, address, city, state, postal_code, country,
			hire_date, termination_date, employment_type, employment_status,
			position, department, manager_id, base_salary, salary_currency,
			payment_frequency, working_hours_per_week, shift_type,
			emergency_contact_name, emergency_contact_phone, emergency_contact_relationship,
			profile_image_url, documents, notes, is_active, created_at, updated_at,
			created_by, updated_by
		FROM employees
		WHERE tenant_id = $1 AND restaurant_id = $2 AND is_active = true
		ORDER BY first_name ASC, last_name ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var employees []domain.Employee
	for rows.Next() {
		var emp domain.Employee
		var firstNameAr, lastNameAr, phone, nationalID sql.NullString
		var address, city, state, postalCode, country sql.NullString
		var emergencyName, emergencyPhone, emergencyRelation sql.NullString
		var profileImageURL, notes sql.NullString
		var dateOfBirth, terminationDate sql.NullTime
		var gender, department sql.NullString
		var managerID, createdBy, updatedBy sql.NullInt64
		var documents sql.NullString

		err := rows.Scan(
			&emp.ID, &emp.TenantID, &emp.RestaurantID, &emp.EmployeeCode,
			&emp.FirstName, &emp.LastName, &firstNameAr, &lastNameAr,
			&emp.Email, &phone, &dateOfBirth, &gender, &nationalID,
			&address, &city, &state, &postalCode, &country,
			&emp.HireDate, &terminationDate, &emp.EmploymentType, &emp.EmploymentStatus,
			&emp.Position, &department, &managerID, &emp.BaseSalary, &emp.SalaryCurrency,
			&emp.PaymentFrequency, &emp.WorkingHoursPerWeek, &emp.ShiftType,
			&emergencyName, &emergencyPhone, &emergencyRelation,
			&profileImageURL, &documents, &notes, &emp.IsActive,
			&emp.CreatedAt, &emp.UpdatedAt, &createdBy, &updatedBy,
		)
		if err != nil {
			return nil, 0, err
		}

		// Handle nullable fields
		emp.FirstNameAr = firstNameAr.String
		emp.LastNameAr = lastNameAr.String
		emp.Phone = phone.String
		emp.NationalID = nationalID.String
		emp.Address = address.String
		emp.City = city.String
		emp.State = state.String
		emp.PostalCode = postalCode.String
		emp.Country = country.String
		emp.Department = department.String
		emp.Gender = gender.String
		emp.EmergencyContactName = emergencyName.String
		emp.EmergencyContactPhone = emergencyPhone.String
		emp.EmergencyContactRelationship = emergencyRelation.String
		emp.ProfileImageURL = profileImageURL.String
		emp.Notes = notes.String

		if dateOfBirth.Valid {
			emp.DateOfBirth = &dateOfBirth.Time
		}
		if terminationDate.Valid {
			emp.TerminationDate = &terminationDate.Time
		}
		if managerID.Valid {
			mid := int(managerID.Int64)
			emp.ManagerID = &mid
		}
		if createdBy.Valid {
			cb := int(createdBy.Int64)
			emp.CreatedBy = &cb
		}
		if updatedBy.Valid {
			ub := int(updatedBy.Int64)
			emp.UpdatedBy = &ub
		}
		if documents.Valid {
			emp.Documents = []byte(documents.String)
		}

		employees = append(employees, emp)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Return empty array instead of nil if no employees found
	if employees == nil {
		employees = []domain.Employee{}
	}

	return employees, total, nil
}

// ListEmployees retrieves all active employees for a restaurant.
func (r *EmployeeRepository) ListEmployees(tenantID, restaurantID int) ([]domain.Employee, error) {
	employees, _, err := r.ListEmployeesPaged(context.Background(), tenantID, restaurantID, 1, 1000000)
	return employees, err
}

// GetEmployeeByID retrieves a single employee by ID
func (r *EmployeeRepository) GetEmployeeByID(tenantID, restaurantID, id int) (*domain.Employee, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_code, first_name, last_name,
			first_name_ar, last_name_ar, email, phone, date_of_birth, gender,
			national_id, address, city, state, postal_code, country,
			hire_date, termination_date, employment_type, employment_status,
			position, department, manager_id, base_salary, salary_currency,
			payment_frequency, working_hours_per_week, shift_type,
			emergency_contact_name, emergency_contact_phone, emergency_contact_relationship,
			profile_image_url, documents, notes, is_active, created_at, updated_at,
			created_by, updated_by
		FROM employees
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	var emp domain.Employee
	var firstNameAr, lastNameAr, phone, nationalID sql.NullString
	var address, city, state, postalCode, country sql.NullString
	var emergencyName, emergencyPhone, emergencyRelation sql.NullString
	var profileImageURL, notes sql.NullString
	var dateOfBirth, terminationDate sql.NullTime
	var gender, department sql.NullString
	var managerID, createdBy, updatedBy sql.NullInt64
	var documents sql.NullString

	err := r.db.QueryRow(query, id, tenantID, restaurantID).Scan(
		&emp.ID, &emp.TenantID, &emp.RestaurantID, &emp.EmployeeCode,
		&emp.FirstName, &emp.LastName, &firstNameAr, &lastNameAr,
		&emp.Email, &phone, &dateOfBirth, &gender, &nationalID,
		&address, &city, &state, &postalCode, &country,
		&emp.HireDate, &terminationDate, &emp.EmploymentType, &emp.EmploymentStatus,
		&emp.Position, &department, &managerID, &emp.BaseSalary, &emp.SalaryCurrency,
		&emp.PaymentFrequency, &emp.WorkingHoursPerWeek, &emp.ShiftType,
		&emergencyName, &emergencyPhone, &emergencyRelation,
		&profileImageURL, &documents, &notes, &emp.IsActive,
		&emp.CreatedAt, &emp.UpdatedAt, &createdBy, &updatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	emp.FirstNameAr = firstNameAr.String
	emp.LastNameAr = lastNameAr.String
	emp.Phone = phone.String
	emp.NationalID = nationalID.String
	emp.Address = address.String
	emp.City = city.String
	emp.State = state.String
	emp.PostalCode = postalCode.String
	emp.Country = country.String
	emp.Department = department.String
	emp.Gender = gender.String
	emp.EmergencyContactName = emergencyName.String
	emp.EmergencyContactPhone = emergencyPhone.String
	emp.EmergencyContactRelationship = emergencyRelation.String
	emp.ProfileImageURL = profileImageURL.String
	emp.Notes = notes.String

	if dateOfBirth.Valid {
		emp.DateOfBirth = &dateOfBirth.Time
	}
	if terminationDate.Valid {
		emp.TerminationDate = &terminationDate.Time
	}
	if managerID.Valid {
		mid := int(managerID.Int64)
		emp.ManagerID = &mid
	}
	if createdBy.Valid {
		cb := int(createdBy.Int64)
		emp.CreatedBy = &cb
	}
	if updatedBy.Valid {
		ub := int(updatedBy.Int64)
		emp.UpdatedBy = &ub
	}
	if documents.Valid {
		emp.Documents = []byte(documents.String)
	}

	return &emp, nil
}

// CreateEmployee creates a new employee
func (r *EmployeeRepository) CreateEmployee(employee *domain.Employee) (int, error) {
	query := `
		INSERT INTO employees (
			tenant_id, restaurant_id, employee_code, first_name, last_name,
			first_name_ar, last_name_ar, email, phone, date_of_birth, gender,
			national_id, address, city, state, postal_code, country,
			hire_date, employment_type, employment_status, position, department,
			manager_id, base_salary, salary_currency, payment_frequency,
			working_hours_per_week, shift_type, emergency_contact_name,
			emergency_contact_phone, emergency_contact_relationship,
			profile_image_url, documents, notes, is_active, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30, $31, $32, $33, $34, $35, $36
		)
		RETURNING id
	`

	// Handle nullable gender - database check constraint only allows 'male', 'female', 'other'
	var gender interface{}
	if employee.Gender == "male" || employee.Gender == "female" || employee.Gender == "other" {
		gender = employee.Gender
	}

	// Handle nil/empty Documents - convert to valid JSON
	var documents interface{} = employee.Documents
	log.Printf("[REPO] CreateEmployee: Documents field length = %d\n", len(employee.Documents))
	log.Printf("[REPO] CreateEmployee: Documents value = %v\n", employee.Documents)

	if len(employee.Documents) == 0 {
		documents = "[]" // Default to empty JSON array
		log.Println("[REPO] CreateEmployee: Setting Documents to empty JSON array '[]'")
	}

	log.Println("[REPO] CreateEmployee: About to insert employee")
	log.Printf("[REPO] CreateEmployee: TenantID=%d, RestaurantID=%d, EmployeeCode=%s\n",
		employee.TenantID, employee.RestaurantID, employee.EmployeeCode)
	log.Printf("[REPO] CreateEmployee: FirstName=%s, Email=%s\n",
		employee.FirstName, employee.Email)
	log.Printf("[REPO] CreateEmployee: Final Documents value = %v (type: %T)\n", documents, documents)

	var id int
	err := r.db.QueryRow(
		query,
		employee.TenantID, employee.RestaurantID, employee.EmployeeCode,
		employee.FirstName, employee.LastName, employee.FirstNameAr, employee.LastNameAr,
		employee.Email, employee.Phone, employee.DateOfBirth, gender,
		employee.NationalID, employee.Address, employee.City, employee.State,
		employee.PostalCode, employee.Country, employee.HireDate, employee.EmploymentType,
		employee.EmploymentStatus, employee.Position, employee.Department,
		employee.ManagerID, employee.BaseSalary, employee.SalaryCurrency,
		employee.PaymentFrequency, employee.WorkingHoursPerWeek, employee.ShiftType,
		employee.EmergencyContactName, employee.EmergencyContactPhone,
		employee.EmergencyContactRelationship, employee.ProfileImageURL,
		documents, employee.Notes, employee.IsActive, employee.CreatedBy,
	).Scan(&id)

	if err != nil {
		log.Printf("[REPO] CreateEmployee: ERROR inserting employee = %v\n", err)
		return 0, err
	}

	log.Printf("[REPO] CreateEmployee: SUCCESS - Employee created with ID = %d\n", id)
	return id, nil
}

// UpdateEmployee updates an existing employee
func (r *EmployeeRepository) UpdateEmployee(employee *domain.Employee) error {
	query := `
		UPDATE employees
		SET
			first_name = $1, last_name = $2, first_name_ar = $3, last_name_ar = $4,
			email = $5, phone = $6, date_of_birth = $7, gender = $8, national_id = $9,
			address = $10, city = $11, state = $12, postal_code = $13, country = $14,
			employment_type = $15, employment_status = $16, position = $17,
			department = $18, manager_id = $19, base_salary = $20, salary_currency = $21,
			payment_frequency = $22, working_hours_per_week = $23, shift_type = $24,
			emergency_contact_name = $25, emergency_contact_phone = $26,
			emergency_contact_relationship = $27, profile_image_url = $28,
			documents = $29, notes = $30, is_active = $31, updated_by = $32,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $33 AND tenant_id = $34 AND restaurant_id = $35
	`

	_, err := r.db.Exec(
		query,
		employee.FirstName, employee.LastName, employee.FirstNameAr, employee.LastNameAr,
		employee.Email, employee.Phone, employee.DateOfBirth, employee.Gender, employee.NationalID,
		employee.Address, employee.City, employee.State, employee.PostalCode, employee.Country,
		employee.EmploymentType, employee.EmploymentStatus, employee.Position,
		employee.Department, employee.ManagerID, employee.BaseSalary, employee.SalaryCurrency,
		employee.PaymentFrequency, employee.WorkingHoursPerWeek, employee.ShiftType,
		employee.EmergencyContactName, employee.EmergencyContactPhone,
		employee.EmergencyContactRelationship, employee.ProfileImageURL,
		employee.Documents, employee.Notes, employee.IsActive, employee.UpdatedBy,
		employee.ID, employee.TenantID, employee.RestaurantID,
	)

	return err
}

// DeleteEmployee soft deletes an employee
func (r *EmployeeRepository) DeleteEmployee(tenantID, restaurantID, id int) error {
	query := `
		UPDATE employees
		SET is_active = false, employment_status = 'terminated',
		    termination_date = CURRENT_DATE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	_, err := r.db.Exec(query, id, tenantID, restaurantID)
	return err
}

// GetEmployeesByStatus retrieves employees by employment status
func (r *EmployeeRepository) GetEmployeesByStatus(tenantID, restaurantID int, status string) ([]domain.Employee, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_code, first_name, last_name,
			email, position, department, employment_type, employment_status,
			hire_date, base_salary, is_active, created_at
		FROM employees
		WHERE tenant_id = $1 AND restaurant_id = $2 AND employment_status = $3 AND is_active = true
		ORDER BY first_name ASC, last_name ASC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []domain.Employee
	for rows.Next() {
		var emp domain.Employee
		var department sql.NullString

		err := rows.Scan(
			&emp.ID, &emp.TenantID, &emp.RestaurantID, &emp.EmployeeCode,
			&emp.FirstName, &emp.LastName, &emp.Email, &emp.Position,
			&department, &emp.EmploymentType, &emp.EmploymentStatus,
			&emp.HireDate, &emp.BaseSalary, &emp.IsActive, &emp.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		emp.Department = department.String
		employees = append(employees, emp)
	}

	if employees == nil {
		employees = []domain.Employee{}
	}

	return employees, nil
}

// GetEmployeeByCode retrieves an employee by employee code
func (r *EmployeeRepository) GetEmployeeByCode(tenantID int, employeeCode string) (*domain.Employee, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_code, first_name, last_name,
			email, position, employment_status, is_active
		FROM employees
		WHERE tenant_id = $1 AND employee_code = $2
	`

	var emp domain.Employee
	err := r.db.QueryRow(query, tenantID, employeeCode).Scan(
		&emp.ID, &emp.TenantID, &emp.RestaurantID, &emp.EmployeeCode,
		&emp.FirstName, &emp.LastName, &emp.Email, &emp.Position,
		&emp.EmploymentStatus, &emp.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &emp, nil
}
