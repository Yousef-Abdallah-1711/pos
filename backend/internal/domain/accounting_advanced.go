package domain

import "time"

type BankReconciliationSession struct {
	ID               int        `json:"id"`
	TenantID         int        `json:"tenant_id"`
	RestaurantID     int        `json:"restaurant_id"`
	BankAccountID    int        `json:"bank_account_id"`
	StatementDate    time.Time  `json:"statement_date"`
	StatementBalance float64    `json:"statement_balance"`
	BookBalance      float64    `json:"book_balance"`
	Difference       float64    `json:"difference"`
	Status           string     `json:"status"`
	Notes            string     `json:"notes"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	CompletedBy      *int       `json:"completed_by,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedBy        *int       `json:"created_by,omitempty"`
	UpdatedBy        *int       `json:"updated_by,omitempty"`
}

type BankReconciliationLine struct {
	ID               int       `json:"id"`
	ReconciliationID int       `json:"reconciliation_id"`
	TenantID         int       `json:"tenant_id"`
	JournalEntryID   *int      `json:"journal_entry_id,omitempty"`
	LineDescription  string    `json:"line_description"`
	StatementAmount  float64   `json:"statement_amount"`
	BookAmount       float64   `json:"book_amount"`
	MatchStatus      string    `json:"match_status"`
	CreatedAt        time.Time `json:"created_at"`
}

type RecurringJournalTemplate struct {
	ID           int        `json:"id"`
	TenantID     int        `json:"tenant_id"`
	RestaurantID int        `json:"restaurant_id"`
	TemplateCode string     `json:"template_code"`
	TemplateName string     `json:"template_name"`
	Description  string     `json:"description"`
	Frequency    string     `json:"frequency"`
	StartDate    time.Time  `json:"start_date"`
	NextRunDate  time.Time  `json:"next_run_date"`
	LastRunDate  *time.Time `json:"last_run_date,omitempty"`
	IsActive     bool       `json:"is_active"`
	AutoPost     bool       `json:"auto_post"`
	LinesJSON    []byte     `json:"lines_json"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedBy    *int       `json:"created_by,omitempty"`
	UpdatedBy    *int       `json:"updated_by,omitempty"`
}

type BudgetEntry struct {
	ID             int       `json:"id"`
	TenantID       int       `json:"tenant_id"`
	RestaurantID   int       `json:"restaurant_id"`
	AccountID      int       `json:"account_id"`
	AccountCode    string    `json:"account_code"`
	AccountName    string    `json:"account_name"`
	FiscalYear     int       `json:"fiscal_year"`
	PeriodMonth    int       `json:"period_month"`
	BudgetAmount   float64   `json:"budget_amount"`
	ActualAmount   float64   `json:"actual_amount"`
	VarianceAmount float64   `json:"variance_amount"`
	VariancePct    float64   `json:"variance_pct"`
	Currency       string    `json:"currency"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ExchangeRate struct {
	ID            int       `json:"id"`
	TenantID      int       `json:"tenant_id"`
	RestaurantID  int       `json:"restaurant_id"`
	BaseCurrency  string    `json:"base_currency"`
	QuoteCurrency string    `json:"quote_currency"`
	Rate          float64   `json:"rate"`
	EffectiveDate time.Time `json:"effective_date"`
	Source        string    `json:"source"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SupplierPaymentSchedule struct {
	ID                int        `json:"id"`
	TenantID          int        `json:"tenant_id"`
	RestaurantID      int        `json:"restaurant_id"`
	SupplierInvoiceID *int       `json:"supplier_invoice_id,omitempty"`
	SupplierID        int        `json:"supplier_id"`
	SupplierName      string     `json:"supplier_name,omitempty"`
	ScheduledDate     time.Time  `json:"scheduled_date"`
	Amount            float64    `json:"amount"`
	Currency          string     `json:"currency"`
	ReminderDays      int        `json:"reminder_days"`
	Status            string     `json:"status"`
	ReminderSentAt    *time.Time `json:"reminder_sent_at,omitempty"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         *int       `json:"created_by,omitempty"`
	UpdatedBy         *int       `json:"updated_by,omitempty"`
}

type EInvoiceSubmission struct {
	ID              int       `json:"id"`
	TenantID        int       `json:"tenant_id"`
	RestaurantID    int       `json:"restaurant_id"`
	InvoiceType     string    `json:"invoice_type"`
	InvoiceID       int       `json:"invoice_id"`
	ETAUUID         string    `json:"eta_uuid"`
	Status          string    `json:"status"`
	SubmittedAt     time.Time `json:"submitted_at"`
	ResponsePayload []byte    `json:"response_payload"`
	ErrorMessage    string    `json:"error_message"`
	CreatedBy       *int      `json:"created_by,omitempty"`
}
