package domain

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Employee represents an employee in the system
type Employee struct {
	ID           int `json:"id"`
	TenantID     int `json:"tenant_id"`
	RestaurantID int `json:"restaurant_id"`

	// Personal Information
	EmployeeCode string     `json:"employee_code"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	FirstNameAr  string     `json:"first_name_ar,omitempty"`
	LastNameAr   string     `json:"last_name_ar,omitempty"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone,omitempty"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
	Gender       string     `json:"gender,omitempty"` // 'male', 'female', 'other'
	NationalID   string     `json:"national_id,omitempty"`

	// Address Information
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`

	// Employment Information
	HireDate         time.Time  `json:"hire_date"`
	TerminationDate  *time.Time `json:"termination_date,omitempty"`
	EmploymentType   string     `json:"employment_type"`   // 'full_time', 'part_time', 'contract', 'intern'
	EmploymentStatus string     `json:"employment_status"` // 'active', 'on_leave', 'suspended', 'terminated'
	Position         string     `json:"position"`
	Department       string     `json:"department,omitempty"`
	ManagerID        *int       `json:"manager_id,omitempty"`

	// Compensation
	BaseSalary       float64 `json:"base_salary"`
	SalaryCurrency   string  `json:"salary_currency"`
	PaymentFrequency string  `json:"payment_frequency"` // 'hourly', 'daily', 'weekly', 'bi_weekly', 'monthly'

	// Work Schedule
	WorkingHoursPerWeek float64 `json:"working_hours_per_week"`
	ShiftType           string  `json:"shift_type"` // 'day', 'night', 'rotating'

	// Emergency Contact
	EmergencyContactName         string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone        string `json:"emergency_contact_phone,omitempty"`
	EmergencyContactRelationship string `json:"emergency_contact_relationship,omitempty"`

	// Documents and Notes
	ProfileImageURL string          `json:"profile_image_url,omitempty"`
	Documents       json.RawMessage `json:"documents,omitempty"` // JSON array
	Notes           string          `json:"notes,omitempty"`

	// System Fields
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy *int      `json:"created_by,omitempty"`
	UpdatedBy *int      `json:"updated_by,omitempty"`

	// Relations (populated on demand)
	Roles []Role `json:"roles,omitempty"`
}

// Role represents an employee role with permissions
type Role struct {
	ID           int `json:"id"`
	TenantID     int `json:"tenant_id"`
	RestaurantID int `json:"restaurant_id"`

	// Role Information
	RoleName      string `json:"role_name"`
	RoleNameAr    string `json:"role_name_ar,omitempty"`
	Description   string `json:"description,omitempty"`
	DescriptionAr string `json:"description_ar,omitempty"`
	RoleCode      string `json:"role_code"`

	// Permissions
	Permissions json.RawMessage `json:"permissions"` // JSON object

	// Access Levels
	AccessLevel        string `json:"access_level"` // 'basic', 'supervisor', 'manager', 'admin', 'owner'
	CanApproveLeaves   bool   `json:"can_approve_leaves"`
	CanApproveOvertime bool   `json:"can_approve_overtime"`
	CanManagePayroll   bool   `json:"can_manage_payroll"`
	CanViewReports     bool   `json:"can_view_reports"`

	// Salary Range
	MinSalary *float64 `json:"min_salary,omitempty"`
	MaxSalary *float64 `json:"max_salary,omitempty"`

	// System Fields
	IsActive     bool      `json:"is_active"`
	IsSystemRole bool      `json:"is_system_role"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    *int      `json:"created_by,omitempty"`
	UpdatedBy    *int      `json:"updated_by,omitempty"`
}

// EmployeeRole represents the junction between employees and roles
type EmployeeRole struct {
	ID              int        `json:"id"`
	TenantID        int        `json:"tenant_id"`
	EmployeeID      int        `json:"employee_id"`
	RoleID          int        `json:"role_id"`
	AssignedDate    time.Time  `json:"assigned_date"`
	EffectiveFrom   time.Time  `json:"effective_from"`
	EffectiveTo     *time.Time `json:"effective_to,omitempty"`
	IsPrimary       bool       `json:"is_primary"`
	AssignedBy      *int       `json:"assigned_by,omitempty"`
	AssignmentNotes string     `json:"assignment_notes,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Attendance represents employee daily attendance
type Attendance struct {
	ID           int `json:"id"`
	TenantID     int `json:"tenant_id"`
	RestaurantID int `json:"restaurant_id"`
	EmployeeID   int `json:"employee_id"`

	// Attendance Date
	AttendanceDate time.Time `json:"attendance_date"`

	// Clock Times
	ClockIn           *time.Time `json:"clock_in,omitempty"`
	ClockOut          *time.Time `json:"clock_out,omitempty"`
	ScheduledClockIn  *string    `json:"scheduled_clock_in,omitempty"`  // TIME format
	ScheduledClockOut *string    `json:"scheduled_clock_out,omitempty"` // TIME format

	// Break Times
	BreakStart        *time.Time `json:"break_start,omitempty"`
	BreakEnd          *time.Time `json:"break_end,omitempty"`
	TotalBreakMinutes int        `json:"total_break_minutes"`

	// Calculated Fields
	TotalHours    *float64 `json:"total_hours,omitempty"`
	RegularHours  float64  `json:"regular_hours"`
	OvertimeHours float64  `json:"overtime_hours"`

	// Attendance Status
	Status                string `json:"status"` // 'present', 'absent', 'late', 'half_day', 'on_leave', 'holiday', 'weekend'
	IsLate                bool   `json:"is_late"`
	LateByMinutes         int    `json:"late_by_minutes"`
	IsEarlyDeparture      bool   `json:"is_early_departure"`
	EarlyDepartureMinutes int    `json:"early_departure_minutes"`

	// Overtime
	IsOvertime         bool       `json:"is_overtime"`
	OvertimeApproved   bool       `json:"overtime_approved"`
	OvertimeApprovedBy *int       `json:"overtime_approved_by,omitempty"`
	OvertimeApprovedAt *time.Time `json:"overtime_approved_at,omitempty"`

	// Location Tracking
	ClockInLocation  json.RawMessage `json:"clock_in_location,omitempty"`  // JSON object
	ClockOutLocation json.RawMessage `json:"clock_out_location,omitempty"` // JSON object
	ClockInIP        string          `json:"clock_in_ip,omitempty"`
	ClockOutIP       string          `json:"clock_out_ip,omitempty"`

	// Device Information
	ClockInDevice  string `json:"clock_in_device,omitempty"`
	ClockOutDevice string `json:"clock_out_device,omitempty"`

	// Notes
	Notes         string `json:"notes,omitempty"`
	AdminNotes    string `json:"admin_notes,omitempty"`
	AbsenceReason string `json:"absence_reason,omitempty"`

	// Approval
	RequiresApproval bool       `json:"requires_approval"`
	IsApproved       bool       `json:"is_approved"`
	ApprovedBy       *int       `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`

	// System Fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy *int      `json:"created_by,omitempty"`
	UpdatedBy *int      `json:"updated_by,omitempty"`
}

// Salary represents employee payroll records
type Salary struct {
	ID           int `json:"id"`
	TenantID     int `json:"tenant_id"`
	RestaurantID int `json:"restaurant_id"`
	EmployeeID   int `json:"employee_id"`

	// Payroll Period
	PayPeriodStart time.Time  `json:"pay_period_start"`
	PayPeriodEnd   time.Time  `json:"pay_period_end"`
	PaymentDate    *time.Time `json:"payment_date,omitempty"`
	Month          int        `json:"month"`
	Year           int        `json:"year"`

	// Base Salary
	BaseSalary float64 `json:"base_salary"`
	Currency   string  `json:"currency"`

	// Additional Earnings
	OvertimeHours  float64 `json:"overtime_hours"`
	OvertimeRate   float64 `json:"overtime_rate"`
	OvertimeAmount float64 `json:"overtime_amount"`
	Bonus          float64 `json:"bonus"`
	Commission     float64 `json:"commission"`
	Allowances     float64 `json:"allowances"`
	Tips           float64 `json:"tips"`
	OtherEarnings  float64 `json:"other_earnings"`

	// Earnings Details
	EarningsDetails json.RawMessage `json:"earnings_details,omitempty"` // JSON object

	// Deductions
	Tax              float64 `json:"tax"`
	SocialInsurance  float64 `json:"social_insurance"`
	HealthInsurance  float64 `json:"health_insurance"`
	Pension          float64 `json:"pension"`
	LoanDeduction    float64 `json:"loan_deduction"`
	AdvanceDeduction float64 `json:"advance_deduction"`
	OtherDeductions  float64 `json:"other_deductions"`

	// Deductions Details
	DeductionsDetails json.RawMessage `json:"deductions_details,omitempty"` // JSON object

	// Calculated Totals
	GrossSalary     float64 `json:"gross_salary"`
	TotalDeductions float64 `json:"total_deductions"`
	NetSalary       float64 `json:"net_salary"`

	// Attendance Summary
	DaysWorked         int     `json:"days_worked"`
	DaysAbsent         int     `json:"days_absent"`
	TotalHoursWorked   float64 `json:"total_hours_worked"`
	TotalOvertimeHours float64 `json:"total_overtime_hours"`

	// Payment Information
	PaymentMethod     string `json:"payment_method"` // 'cash', 'bank_transfer', 'check', 'mobile_money'
	PaymentReference  string `json:"payment_reference,omitempty"`
	BankAccountNumber string `json:"bank_account_number,omitempty"`
	BankName          string `json:"bank_name,omitempty"`

	// Status
	Status string     `json:"status"` // 'pending', 'processing', 'paid', 'cancelled', 'on_hold'
	IsPaid bool       `json:"is_paid"`
	PaidAt *time.Time `json:"paid_at,omitempty"`
	PaidBy *int       `json:"paid_by,omitempty"`

	// Approval
	IsApproved    bool       `json:"is_approved"`
	ApprovedBy    *int       `json:"approved_by,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	ApprovalNotes string     `json:"approval_notes,omitempty"`

	// Notes
	Notes            string `json:"notes,omitempty"`
	CalculationNotes string `json:"calculation_notes,omitempty"`

	// System Fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy *int      `json:"created_by,omitempty"`
	UpdatedBy *int      `json:"updated_by,omitempty"`
}

// Leave represents employee leave requests
type Leave struct {
	ID           int `json:"id"`
	TenantID     int `json:"tenant_id"`
	RestaurantID int `json:"restaurant_id"`
	EmployeeID   int `json:"employee_id"`

	// Leave Period
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	TotalDays     int       `json:"total_days"`
	IsHalfDay     bool      `json:"is_half_day"`
	HalfDayPeriod string    `json:"half_day_period,omitempty"` // 'morning', 'afternoon'

	// Leave Type
	LeaveType     string `json:"leave_type"`     // 'annual', 'sick', 'casual', 'maternity', 'paternity', etc.
	LeaveCategory string `json:"leave_category"` // 'paid', 'unpaid'

	// Leave Details
	Reason           string `json:"reason"`
	ContactNumber    string `json:"contact_number,omitempty"`
	ContactAddress   string `json:"contact_address,omitempty"`
	EmergencyContact string `json:"emergency_contact,omitempty"`

	// Supporting Documents
	Attachments json.RawMessage `json:"attachments,omitempty"` // JSON array

	// Status and Approval
	Status     string `json:"status"` // 'pending', 'approved', 'rejected', 'cancelled', 'on_hold'
	IsApproved bool   `json:"is_approved"`

	// Approval Chain
	ApprovedBy    *int       `json:"approved_by,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	ApprovalNotes string     `json:"approval_notes,omitempty"`

	RejectedBy      *int       `json:"rejected_by,omitempty"`
	RejectedAt      *time.Time `json:"rejected_at,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`

	CancelledBy        *int       `json:"cancelled_by,omitempty"`
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty"`

	// Leave Balance Impact
	DeductedFromBalance bool     `json:"deducted_from_balance"`
	BalanceBefore       *float64 `json:"balance_before,omitempty"`
	BalanceAfter        *float64 `json:"balance_after,omitempty"`

	// HR Notes
	HRNotes       string `json:"hr_notes,omitempty"`
	HandoverNotes string `json:"handover_notes,omitempty"`

	// Replacement/Coverage
	ReplacementEmployeeID *int `json:"replacement_employee_id,omitempty"`
	HandoverCompleted     bool `json:"handover_completed"`

	// Notification Settings
	NotifyManager      bool       `json:"notify_manager"`
	NotifyHR           bool       `json:"notify_hr"`
	NotificationSentAt *time.Time `json:"notification_sent_at,omitempty"`

	// System Fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy *int      `json:"created_by,omitempty"`
	UpdatedBy *int      `json:"updated_by,omitempty"`
}

// NullTime is a helper type for handling nullable time fields
type NullTime sql.NullTime

// Scan implements the sql.Scanner interface
func (nt *NullTime) Scan(value interface{}) error {
	var t sql.NullTime
	if err := t.Scan(value); err != nil {
		return err
	}
	*nt = NullTime(t)
	return nil
}

// MarshalJSON implements json.Marshaler interface
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}
