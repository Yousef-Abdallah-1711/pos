package repository

import (
	"context"
	"database/sql"
	"time"

	"pos-saas/internal/domain"
)

// AttendanceRepository handles attendance data operations
type AttendanceRepository struct {
	db *sql.DB
}

// NewAttendanceRepository creates a new attendance repository
func NewAttendanceRepository(db *sql.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

// ListAttendancePaged retrieves attendance records for a date range using DB pagination.
func (r *AttendanceRepository) ListAttendancePaged(ctx context.Context, tenantID, restaurantID int, startDate, endDate time.Time, page, pageSize int) ([]domain.Attendance, int, error) {
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
		FROM attendance
		WHERE tenant_id = $1 AND restaurant_id = $2
		  AND attendance_date BETWEEN $3 AND $4
	`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, tenantID, restaurantID, startDate, endDate).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_id, attendance_date,
			clock_in, clock_out, scheduled_clock_in, scheduled_clock_out,
			break_start, break_end, total_break_minutes, total_hours,
			regular_hours, overtime_hours, status, is_late, late_by_minutes,
			is_early_departure, early_departure_minutes, is_overtime,
			overtime_approved, notes, created_at, updated_at
		FROM attendance
		WHERE tenant_id = $1 AND restaurant_id = $2
		  AND attendance_date BETWEEN $3 AND $4
		ORDER BY attendance_date DESC, employee_id ASC
		LIMIT $5 OFFSET $6
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID, startDate, endDate, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var attendances []domain.Attendance
	for rows.Next() {
		var att domain.Attendance
		var clockIn, clockOut, breakStart, breakEnd sql.NullTime
		var scheduledClockIn, scheduledClockOut sql.NullString
		var totalHours sql.NullFloat64
		var notes sql.NullString

		err := rows.Scan(
			&att.ID, &att.TenantID, &att.RestaurantID, &att.EmployeeID,
			&att.AttendanceDate, &clockIn, &clockOut, &scheduledClockIn,
			&scheduledClockOut, &breakStart, &breakEnd, &att.TotalBreakMinutes,
			&totalHours, &att.RegularHours, &att.OvertimeHours, &att.Status,
			&att.IsLate, &att.LateByMinutes, &att.IsEarlyDeparture,
			&att.EarlyDepartureMinutes, &att.IsOvertime, &att.OvertimeApproved,
			&notes, &att.CreatedAt, &att.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if clockIn.Valid {
			att.ClockIn = &clockIn.Time
		}
		if clockOut.Valid {
			att.ClockOut = &clockOut.Time
		}
		if breakStart.Valid {
			att.BreakStart = &breakStart.Time
		}
		if breakEnd.Valid {
			att.BreakEnd = &breakEnd.Time
		}
		if scheduledClockIn.Valid {
			att.ScheduledClockIn = &scheduledClockIn.String
		}
		if scheduledClockOut.Valid {
			att.ScheduledClockOut = &scheduledClockOut.String
		}
		if totalHours.Valid {
			att.TotalHours = &totalHours.Float64
		}
		att.Notes = notes.String

		attendances = append(attendances, att)
	}

	if attendances == nil {
		attendances = []domain.Attendance{}
	}

	return attendances, total, nil
}

// ListAttendance retrieves attendance records for a date range.
func (r *AttendanceRepository) ListAttendance(tenantID, restaurantID int, startDate, endDate time.Time) ([]domain.Attendance, error) {
	out, _, err := r.ListAttendancePaged(context.Background(), tenantID, restaurantID, startDate, endDate, 1, 1000000)
	return out, err
}

// GetAttendanceByID retrieves a single attendance record by ID
func (r *AttendanceRepository) GetAttendanceByID(tenantID, restaurantID, id int) (*domain.Attendance, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_id, attendance_date,
			clock_in, clock_out, scheduled_clock_in, scheduled_clock_out,
			break_start, break_end, total_break_minutes, total_hours,
			regular_hours, overtime_hours, status, is_late, late_by_minutes,
			is_early_departure, early_departure_minutes, is_overtime,
			overtime_approved, notes, admin_notes, created_at, updated_at
		FROM attendance
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	var att domain.Attendance
	var clockIn, clockOut, breakStart, breakEnd sql.NullTime
	var scheduledClockIn, scheduledClockOut sql.NullString
	var totalHours sql.NullFloat64
	var notes, adminNotes sql.NullString

	err := r.db.QueryRow(query, id, tenantID, restaurantID).Scan(
		&att.ID, &att.TenantID, &att.RestaurantID, &att.EmployeeID,
		&att.AttendanceDate, &clockIn, &clockOut, &scheduledClockIn,
		&scheduledClockOut, &breakStart, &breakEnd, &att.TotalBreakMinutes,
		&totalHours, &att.RegularHours, &att.OvertimeHours, &att.Status,
		&att.IsLate, &att.LateByMinutes, &att.IsEarlyDeparture,
		&att.EarlyDepartureMinutes, &att.IsOvertime, &att.OvertimeApproved,
		&notes, &adminNotes, &att.CreatedAt, &att.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if clockIn.Valid {
		att.ClockIn = &clockIn.Time
	}
	if clockOut.Valid {
		att.ClockOut = &clockOut.Time
	}
	if breakStart.Valid {
		att.BreakStart = &breakStart.Time
	}
	if breakEnd.Valid {
		att.BreakEnd = &breakEnd.Time
	}
	if scheduledClockIn.Valid {
		att.ScheduledClockIn = &scheduledClockIn.String
	}
	if scheduledClockOut.Valid {
		att.ScheduledClockOut = &scheduledClockOut.String
	}
	if totalHours.Valid {
		att.TotalHours = &totalHours.Float64
	}
	att.Notes = notes.String
	att.AdminNotes = adminNotes.String

	return &att, nil
}

// CreateAttendance creates a new attendance record
func (r *AttendanceRepository) CreateAttendance(attendance *domain.Attendance) (int, error) {
	query := `
		INSERT INTO attendance (
			tenant_id, restaurant_id, employee_id, attendance_date,
			clock_in, clock_out, scheduled_clock_in, scheduled_clock_out,
			status, notes, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(
		query,
		attendance.TenantID, attendance.RestaurantID, attendance.EmployeeID,
		attendance.AttendanceDate, attendance.ClockIn, attendance.ClockOut,
		attendance.ScheduledClockIn, attendance.ScheduledClockOut,
		attendance.Status, attendance.Notes, attendance.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateAttendance updates an existing attendance record
func (r *AttendanceRepository) UpdateAttendance(attendance *domain.Attendance) error {
	query := `
		UPDATE attendance
		SET
			clock_in = $1, clock_out = $2, scheduled_clock_in = $3,
			scheduled_clock_out = $4, break_start = $5, break_end = $6,
			total_break_minutes = $7, status = $8, notes = $9,
			admin_notes = $10, updated_by = $11, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12 AND tenant_id = $13 AND restaurant_id = $14
	`

	_, err := r.db.Exec(
		query,
		attendance.ClockIn, attendance.ClockOut, attendance.ScheduledClockIn,
		attendance.ScheduledClockOut, attendance.BreakStart, attendance.BreakEnd,
		attendance.TotalBreakMinutes, attendance.Status, attendance.Notes,
		attendance.AdminNotes, attendance.UpdatedBy,
		attendance.ID, attendance.TenantID, attendance.RestaurantID,
	)

	return err
}

// ClockIn records employee clock-in
func (r *AttendanceRepository) ClockIn(tenantID, restaurantID, employeeID int) (int, error) {
	query := `
		INSERT INTO attendance (
			tenant_id, restaurant_id, employee_id, attendance_date,
			clock_in, status
		) VALUES (
			$1, $2, $3, CURRENT_DATE, CURRENT_TIMESTAMP, 'present'
		)
		ON CONFLICT (employee_id, attendance_date)
		DO UPDATE SET clock_in = CURRENT_TIMESTAMP
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(query, tenantID, restaurantID, employeeID).Scan(&id)
	return id, err
}

// ClockOut records employee clock-out
func (r *AttendanceRepository) ClockOut(tenantID, restaurantID, employeeID int) error {
	query := `
		UPDATE attendance
		SET clock_out = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE employee_id = $1 AND tenant_id = $2 AND restaurant_id = $3
		  AND attendance_date = CURRENT_DATE AND clock_out IS NULL
	`

	_, err := r.db.Exec(query, employeeID, tenantID, restaurantID)
	return err
}

// GetTodayAttendance retrieves today's attendance for an employee
func (r *AttendanceRepository) GetTodayAttendance(tenantID, restaurantID, employeeID int) (*domain.Attendance, error) {
	query := `
		SELECT
			id, tenant_id, restaurant_id, employee_id, attendance_date,
			clock_in, clock_out, total_hours, status, created_at
		FROM attendance
		WHERE tenant_id = $1 AND restaurant_id = $2 AND employee_id = $3
		  AND attendance_date = CURRENT_DATE
	`

	var att domain.Attendance
	var clockIn, clockOut sql.NullTime
	var totalHours sql.NullFloat64

	err := r.db.QueryRow(query, tenantID, restaurantID, employeeID).Scan(
		&att.ID, &att.TenantID, &att.RestaurantID, &att.EmployeeID,
		&att.AttendanceDate, &clockIn, &clockOut, &totalHours,
		&att.Status, &att.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if clockIn.Valid {
		att.ClockIn = &clockIn.Time
	}
	if clockOut.Valid {
		att.ClockOut = &clockOut.Time
	}
	if totalHours.Valid {
		att.TotalHours = &totalHours.Float64
	}

	return &att, nil
}
