package domain

import "time"

type PerformanceReview struct {
	ID                int        `json:"id"`
	TenantID          int        `json:"tenant_id"`
	RestaurantID      int        `json:"restaurant_id"`
	EmployeeID        int        `json:"employee_id"`
	ReviewerID        int        `json:"reviewer_id"`
	ReviewPeriodStart time.Time  `json:"review_period_start"`
	ReviewPeriodEnd   time.Time  `json:"review_period_end"`
	ReviewDate        time.Time  `json:"review_date"`
	OverallRating     float64    `json:"overall_rating"`
	GoalsAchievement  float64    `json:"goals_achievement"`
	TechnicalSkills   float64    `json:"technical_skills"`
	Communication     float64    `json:"communication"`
	Leadership        float64    `json:"leadership"`
	Strengths         string     `json:"strengths,omitempty"`
	AreasToImprove    string     `json:"areas_to_improve,omitempty"`
	ActionPlan        string     `json:"action_plan,omitempty"`
	Status            string     `json:"status"`
	AcknowledgedAt    *time.Time `json:"acknowledged_at,omitempty"`
	AcknowledgedBy    *int       `json:"acknowledged_by,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         *int       `json:"created_by,omitempty"`
	UpdatedBy         *int       `json:"updated_by,omitempty"`
}

type TrainingRecord struct {
	ID                int        `json:"id"`
	TenantID          int        `json:"tenant_id"`
	RestaurantID      int        `json:"restaurant_id"`
	EmployeeID        int        `json:"employee_id"`
	TrainingTitle     string     `json:"training_title"`
	TrainingProvider  string     `json:"training_provider,omitempty"`
	TrainingType      string     `json:"training_type"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           *time.Time `json:"end_date,omitempty"`
	Status            string     `json:"status"`
	CertificateNumber string     `json:"certificate_number,omitempty"`
	CertificateExpiry *time.Time `json:"certificate_expiry,omitempty"`
	Notes             string     `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         *int       `json:"created_by,omitempty"`
	UpdatedBy         *int       `json:"updated_by,omitempty"`
}

type InsuranceTracking struct {
	ID             int       `json:"id"`
	TenantID       int       `json:"tenant_id"`
	RestaurantID   int       `json:"restaurant_id"`
	EmployeeID     int       `json:"employee_id"`
	PolicyType     string    `json:"policy_type"`
	ProviderName   string    `json:"provider_name"`
	PolicyNumber   string    `json:"policy_number"`
	CoverageStart  time.Time `json:"coverage_start"`
	CoverageEnd    time.Time `json:"coverage_end"`
	MonthlyPremium float64   `json:"monthly_premium"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      *int      `json:"created_by,omitempty"`
	UpdatedBy      *int      `json:"updated_by,omitempty"`
}

type EmploymentContract struct {
	ID             int        `json:"id"`
	TenantID       int        `json:"tenant_id"`
	RestaurantID   int        `json:"restaurant_id"`
	EmployeeID     int        `json:"employee_id"`
	ContractType   string     `json:"contract_type"`
	ContractNumber string     `json:"contract_number"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	SalaryAmount   float64    `json:"salary_amount"`
	Currency       string     `json:"currency"`
	Status         string     `json:"status"`
	RenewedFromID  *int       `json:"renewed_from_id,omitempty"`
	Notes          string     `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      *int       `json:"created_by,omitempty"`
	UpdatedBy      *int       `json:"updated_by,omitempty"`
}

type Resignation struct {
	ID                 int        `json:"id"`
	TenantID           int        `json:"tenant_id"`
	RestaurantID       int        `json:"restaurant_id"`
	EmployeeID         int        `json:"employee_id"`
	ResignationDate    time.Time  `json:"resignation_date"`
	LastWorkingDate    time.Time  `json:"last_working_date"`
	Reason             string     `json:"reason"`
	Status             string     `json:"status"`
	NoticePeriodDays   int        `json:"notice_period_days"`
	PendingSalary      float64    `json:"pending_salary"`
	UnusedLeaveBalance float64    `json:"unused_leave_balance"`
	SettlementAmount   float64    `json:"settlement_amount"`
	AcknowledgedAt     *time.Time `json:"acknowledged_at,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	CreatedBy          *int       `json:"created_by,omitempty"`
	UpdatedBy          *int       `json:"updated_by,omitempty"`
}
