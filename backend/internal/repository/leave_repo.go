package repository

import (
	"database/sql"
	"time"

	"pos-saas/internal/domain"
)

// LeaveRepository handles leave data operations
type LeaveRepository struct {
	db *sql.DB
}

// NewLeaveRepository creates a new leave repository
func NewLeaveRepository(db *sql.DB) *LeaveRepository {
	return &LeaveRepository{db: db}
}

// ListLeaves retrieves leave records for a date range
func (r *LeaveRepository) ListLeaves(tenantID, restaurantID int, startDate, endDate time.Time) ([]domain.Leave, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_id, start_date, end_date,
			total_days, is_half_day, half_day_period, leave_type, leave_category,
			reason, status, is_approved, created_at, updated_at
		FROM leaves
		WHERE tenant_id = $1 AND restaurant_id = $2
		  AND start_date >= $3 AND end_date <= $4
		ORDER BY start_date DESC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaves []domain.Leave
	for rows.Next() {
		var leave domain.Leave
		var halfDayPeriod sql.NullString

		err := rows.Scan(
			&leave.ID, &leave.TenantID, &leave.RestaurantID, &leave.EmployeeID,
			&leave.StartDate, &leave.EndDate, &leave.TotalDays, &leave.IsHalfDay,
			&halfDayPeriod, &leave.LeaveType, &leave.LeaveCategory, &leave.Reason,
			&leave.Status, &leave.IsApproved, &leave.CreatedAt, &leave.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		leave.HalfDayPeriod = halfDayPeriod.String
		leaves = append(leaves, leave)
	}

	if leaves == nil {
		leaves = []domain.Leave{}
	}

	return leaves, nil
}

// GetLeaveByID retrieves a single leave record by ID
func (r *LeaveRepository) GetLeaveByID(tenantID, restaurantID, id int) (*domain.Leave, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_id, start_date, end_date,
			total_days, is_half_day, half_day_period, leave_type, leave_category,
			reason, contact_number, contact_address, emergency_contact, attachments,
			status, is_approved, approved_by, approved_at, approval_notes,
			rejected_by, rejected_at, rejection_reason, hr_notes, handover_notes,
			replacement_employee_id, handover_completed, created_at, updated_at
		FROM leaves
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	var leave domain.Leave
	var halfDayPeriod, contactNumber, contactAddress, emergencyContact sql.NullString
	var approvalNotes, rejectionReason, hrNotes, handoverNotes sql.NullString
	var attachments sql.NullString
	var approvedBy, rejectedBy, replacementEmployeeID sql.NullInt64
	var approvedAt, rejectedAt sql.NullTime

	err := r.db.QueryRow(query, id, tenantID, restaurantID).Scan(
		&leave.ID, &leave.TenantID, &leave.RestaurantID, &leave.EmployeeID,
		&leave.StartDate, &leave.EndDate, &leave.TotalDays, &leave.IsHalfDay,
		&halfDayPeriod, &leave.LeaveType, &leave.LeaveCategory, &leave.Reason,
		&contactNumber, &contactAddress, &emergencyContact, &attachments,
		&leave.Status, &leave.IsApproved, &approvedBy, &approvedAt, &approvalNotes,
		&rejectedBy, &rejectedAt, &rejectionReason, &hrNotes, &handoverNotes,
		&replacementEmployeeID, &leave.HandoverCompleted, &leave.CreatedAt, &leave.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	leave.HalfDayPeriod = halfDayPeriod.String
	leave.ContactNumber = contactNumber.String
	leave.ContactAddress = contactAddress.String
	leave.EmergencyContact = emergencyContact.String
	leave.ApprovalNotes = approvalNotes.String
	leave.RejectionReason = rejectionReason.String
	leave.HRNotes = hrNotes.String
	leave.HandoverNotes = handoverNotes.String

	if attachments.Valid {
		leave.Attachments = []byte(attachments.String)
	}
	if approvedBy.Valid {
		ab := int(approvedBy.Int64)
		leave.ApprovedBy = &ab
	}
	if rejectedBy.Valid {
		rb := int(rejectedBy.Int64)
		leave.RejectedBy = &rb
	}
	if replacementEmployeeID.Valid {
		re := int(replacementEmployeeID.Int64)
		leave.ReplacementEmployeeID = &re
	}
	if approvedAt.Valid {
		leave.ApprovedAt = &approvedAt.Time
	}
	if rejectedAt.Valid {
		leave.RejectedAt = &rejectedAt.Time
	}

	return &leave, nil
}

// CreateLeave creates a new leave request
func (r *LeaveRepository) CreateLeave(leave *domain.Leave) (int, error) {
	query := `
		INSERT INTO leaves (
			tenant_id, restaurant_id, employee_id, start_date, end_date,
			total_days, is_half_day, half_day_period, leave_type, leave_category,
			reason, contact_number, contact_address, emergency_contact,
			status, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
		RETURNING id
	`

	// Handle half_day_period: must be NULL when is_half_day is false (DB constraint)
	var halfDayPeriod interface{}
	if leave.IsHalfDay && leave.HalfDayPeriod != "" {
		halfDayPeriod = leave.HalfDayPeriod
	}

	var id int
	err := r.db.QueryRow(
		query,
		leave.TenantID, leave.RestaurantID, leave.EmployeeID, leave.StartDate,
		leave.EndDate, leave.TotalDays, leave.IsHalfDay, halfDayPeriod,
		leave.LeaveType, leave.LeaveCategory, leave.Reason, leave.ContactNumber,
		leave.ContactAddress, leave.EmergencyContact, leave.Status, leave.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateLeave updates an existing leave request
func (r *LeaveRepository) UpdateLeave(leave *domain.Leave) error {
	query := `
		UPDATE leaves
		SET
			start_date = $1, end_date = $2, total_days = $3, is_half_day = $4,
			half_day_period = $5, leave_type = $6, leave_category = $7, reason = $8,
			contact_number = $9, contact_address = $10, emergency_contact = $11,
			status = $12, updated_by = $13, updated_at = CURRENT_TIMESTAMP
		WHERE id = $14 AND tenant_id = $15 AND restaurant_id = $16
	`

	_, err := r.db.Exec(
		query,
		leave.StartDate, leave.EndDate, leave.TotalDays, leave.IsHalfDay,
		leave.HalfDayPeriod, leave.LeaveType, leave.LeaveCategory, leave.Reason,
		leave.ContactNumber, leave.ContactAddress, leave.EmergencyContact,
		leave.Status, leave.UpdatedBy, leave.ID, leave.TenantID, leave.RestaurantID,
	)

	return err
}

// ApproveLeave approves a leave request
func (r *LeaveRepository) ApproveLeave(tenantID, restaurantID, id int, approvedBy int, notes string) error {
	query := `
		UPDATE leaves
		SET
			is_approved = true, status = 'approved', approved_by = $1,
			approved_at = CURRENT_TIMESTAMP, approval_notes = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5 AND status = 'pending'
	`

	_, err := r.db.Exec(query, approvedBy, notes, id, tenantID, restaurantID)
	return err
}

// RejectLeave rejects a leave request
func (r *LeaveRepository) RejectLeave(tenantID, restaurantID, id int, rejectedBy int, reason string) error {
	query := `
		UPDATE leaves
		SET
			is_approved = false, status = 'rejected', rejected_by = $1,
			rejected_at = CURRENT_TIMESTAMP, rejection_reason = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5 AND status = 'pending'
	`

	_, err := r.db.Exec(query, rejectedBy, reason, id, tenantID, restaurantID)
	return err
}

// CancelLeave cancels a leave request
func (r *LeaveRepository) CancelLeave(tenantID, restaurantID, id int, cancelledBy int, reason string) error {
	query := `
		UPDATE leaves
		SET
			status = 'cancelled', cancelled_by = $1, cancelled_at = CURRENT_TIMESTAMP,
			cancellation_reason = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5
		  AND status IN ('pending', 'approved')
	`

	_, err := r.db.Exec(query, cancelledBy, reason, id, tenantID, restaurantID)
	return err
}

// GetPendingLeaves retrieves all pending leave requests
func (r *LeaveRepository) GetPendingLeaves(tenantID, restaurantID int) ([]domain.Leave, error) {
	query := `
		SELECT
			id, employee_id, start_date, end_date, total_days,
			leave_type, leave_category, reason, created_at
		FROM leaves
		WHERE tenant_id = $1 AND restaurant_id = $2 AND status = 'pending'
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaves []domain.Leave
	for rows.Next() {
		var leave domain.Leave

		err := rows.Scan(
			&leave.ID, &leave.EmployeeID, &leave.StartDate, &leave.EndDate,
			&leave.TotalDays, &leave.LeaveType, &leave.LeaveCategory,
			&leave.Reason, &leave.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		leaves = append(leaves, leave)
	}

	if leaves == nil {
		leaves = []domain.Leave{}
	}

	return leaves, nil
}

// GetEmployeeLeaves retrieves all leave records for an employee
func (r *LeaveRepository) GetEmployeeLeaves(tenantID, restaurantID, employeeID int) ([]domain.Leave, error) {
	query := `
		SELECT
			id, start_date, end_date, total_days, leave_type, leave_category,
			reason, status, is_approved, created_at
		FROM leaves
		WHERE tenant_id = $1 AND restaurant_id = $2 AND employee_id = $3
		ORDER BY start_date DESC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaves []domain.Leave
	for rows.Next() {
		var leave domain.Leave

		err := rows.Scan(
			&leave.ID, &leave.StartDate, &leave.EndDate, &leave.TotalDays,
			&leave.LeaveType, &leave.LeaveCategory, &leave.Reason,
			&leave.Status, &leave.IsApproved, &leave.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		leaves = append(leaves, leave)
	}

	if leaves == nil {
		leaves = []domain.Leave{}
	}

	return leaves, nil
}
