package domain

import "time"

const (
	BankAccountTypeChecking   = "checking"
	BankAccountTypeSavings    = "savings"
	BankAccountTypeCreditCard = "credit_card"
	BankAccountTypePettyCash  = "petty_cash"

	InvoiceStatusDraft         = "draft"
	InvoiceStatusSent          = "sent"
	InvoiceStatusApproved      = "approved"
	InvoiceStatusPartiallyPaid = "partially_paid"
	InvoiceStatusPaid          = "paid"
	InvoiceStatusOverdue       = "overdue"
	InvoiceStatusVoid          = "void"

	AccountingPaymentStatusCompleted = "completed"
	AccountingPaymentStatusVoid      = "void"
)

type BankAccount struct {
	ID                    int        `json:"id"`
	TenantID              int        `json:"tenant_id"`
	RestaurantID          int        `json:"restaurant_id"`
	AccountName           string     `json:"account_name"`
	AccountNameAr         string     `json:"account_name_ar,omitempty"`
	BankName              string     `json:"bank_name"`
	BankBranch            string     `json:"bank_branch,omitempty"`
	AccountNumber         string     `json:"account_number"`
	IBAN                  string     `json:"iban,omitempty"`
	SwiftCode             string     `json:"swift_code,omitempty"`
	Currency              string     `json:"currency"`
	AccountType           string     `json:"account_type"`
	OpeningBalance        float64    `json:"opening_balance"`
	CurrentBalance        float64    `json:"current_balance"`
	OpeningDate           time.Time  `json:"opening_date"`
	GLAccountID           *int       `json:"gl_account_id,omitempty"`
	IsActive              bool       `json:"is_active"`
	IsPrimary             bool       `json:"is_primary"`
	LastReconciledDate    *time.Time `json:"last_reconciled_date,omitempty"`
	LastReconciledBalance *float64   `json:"last_reconciled_balance,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	CreatedBy             *int       `json:"created_by,omitempty"`
	UpdatedBy             *int       `json:"updated_by,omitempty"`
}

type CreateBankAccountRequest struct {
	AccountName    string  `json:"account_name"`
	AccountNameAr  string  `json:"account_name_ar"`
	BankName       string  `json:"bank_name"`
	BankBranch     string  `json:"bank_branch"`
	AccountNumber  string  `json:"account_number"`
	IBAN           string  `json:"iban"`
	SwiftCode      string  `json:"swift_code"`
	Currency       string  `json:"currency"`
	AccountType    string  `json:"account_type"`
	OpeningBalance float64 `json:"opening_balance"`
	OpeningDate    string  `json:"opening_date"`
	GLAccountID    *int    `json:"gl_account_id"`
	IsActive       *bool   `json:"is_active"`
	IsPrimary      bool    `json:"is_primary"`
}

type UpdateBankAccountRequest = CreateBankAccountRequest

type BankTransaction struct {
	TransactionDate time.Time `json:"transaction_date"`
	EntryNumber     string    `json:"entry_number"`
	ReferenceNumber string    `json:"reference_number,omitempty"`
	Description     string    `json:"description"`
	DebitAmount     float64   `json:"debit_amount"`
	CreditAmount    float64   `json:"credit_amount"`
	RunningBalance  float64   `json:"running_balance"`
	JournalEntryID  int       `json:"journal_entry_id"`
}

type InvoiceLine struct {
	ID             int       `json:"id"`
	InvoiceID      int       `json:"invoice_id"`
	TenantID       int       `json:"tenant_id"`
	Description    string    `json:"description"`
	DescriptionAr  string    `json:"description_ar,omitempty"`
	Quantity       float64   `json:"quantity"`
	UnitPrice      float64   `json:"unit_price"`
	TaxRate        float64   `json:"tax_rate"`
	TaxAmount      float64   `json:"tax_amount"`
	DiscountRate   float64   `json:"discount_rate,omitempty"`
	DiscountAmount float64   `json:"discount_amount"`
	LineTotal      float64   `json:"line_total"`
	AccountID      *int      `json:"account_id,omitempty"`
	LineOrder      int       `json:"line_order"`
	CreatedAt      time.Time `json:"created_at"`
}

type InvoiceLineInput struct {
	Description    string  `json:"description"`
	DescriptionAr  string  `json:"description_ar"`
	Quantity       float64 `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	TaxRate        float64 `json:"tax_rate"`
	TaxAmount      float64 `json:"tax_amount"`
	DiscountRate   float64 `json:"discount_rate"`
	DiscountAmount float64 `json:"discount_amount"`
	LineTotal      float64 `json:"line_total"`
	AccountID      *int    `json:"account_id"`
	LineOrder      int     `json:"line_order"`
}

type CustomerInvoice struct {
	ID             int           `json:"id"`
	TenantID       int           `json:"tenant_id"`
	RestaurantID   int           `json:"restaurant_id"`
	InvoiceNumber  string        `json:"invoice_number"`
	InvoiceDate    time.Time     `json:"invoice_date"`
	DueDate        time.Time     `json:"due_date"`
	CustomerID     int           `json:"customer_id"`
	CustomerName   string        `json:"customer_name,omitempty"`
	OrderID        *int          `json:"order_id,omitempty"`
	Subtotal       float64       `json:"subtotal"`
	TaxAmount      float64       `json:"tax_amount"`
	DiscountAmount float64       `json:"discount_amount"`
	TotalAmount    float64       `json:"total_amount"`
	AmountPaid     float64       `json:"amount_paid"`
	BalanceDue     float64       `json:"balance_due"`
	Currency       string        `json:"currency"`
	Status         string        `json:"status"`
	Notes          string        `json:"notes,omitempty"`
	Terms          string        `json:"terms,omitempty"`
	InternalNotes  string        `json:"internal_notes,omitempty"`
	JournalEntryID *int          `json:"journal_entry_id,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	CreatedBy      *int          `json:"created_by,omitempty"`
	UpdatedBy      *int          `json:"updated_by,omitempty"`
	Lines          []InvoiceLine `json:"lines,omitempty"`
}

type CreateCustomerInvoiceRequest struct {
	InvoiceDate   string             `json:"invoice_date"`
	DueDate       string             `json:"due_date"`
	CustomerID    int                `json:"customer_id"`
	OrderID       *int               `json:"order_id"`
	Currency      string             `json:"currency"`
	Status        string             `json:"status"`
	Notes         string             `json:"notes"`
	Terms         string             `json:"terms"`
	InternalNotes string             `json:"internal_notes"`
	Lines         []InvoiceLineInput `json:"lines"`
}

type UpdateCustomerInvoiceRequest = CreateCustomerInvoiceRequest

type SupplierInvoice struct {
	ID             int           `json:"id"`
	TenantID       int           `json:"tenant_id"`
	RestaurantID   int           `json:"restaurant_id"`
	InvoiceNumber  string        `json:"invoice_number"`
	SupplierRef    string        `json:"supplier_ref,omitempty"`
	InvoiceDate    time.Time     `json:"invoice_date"`
	DueDate        time.Time     `json:"due_date"`
	SupplierID     int           `json:"supplier_id"`
	SupplierName   string        `json:"supplier_name,omitempty"`
	Subtotal       float64       `json:"subtotal"`
	TaxAmount      float64       `json:"tax_amount"`
	DiscountAmount float64       `json:"discount_amount"`
	TotalAmount    float64       `json:"total_amount"`
	AmountPaid     float64       `json:"amount_paid"`
	BalanceDue     float64       `json:"balance_due"`
	Currency       string        `json:"currency"`
	Status         string        `json:"status"`
	Notes          string        `json:"notes,omitempty"`
	InternalNotes  string        `json:"internal_notes,omitempty"`
	JournalEntryID *int          `json:"journal_entry_id,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	CreatedBy      *int          `json:"created_by,omitempty"`
	UpdatedBy      *int          `json:"updated_by,omitempty"`
	Lines          []InvoiceLine `json:"lines,omitempty"`
}

type CreateSupplierInvoiceRequest struct {
	SupplierRef   string             `json:"supplier_ref"`
	InvoiceDate   string             `json:"invoice_date"`
	DueDate       string             `json:"due_date"`
	SupplierID    int                `json:"supplier_id"`
	Currency      string             `json:"currency"`
	Status        string             `json:"status"`
	Notes         string             `json:"notes"`
	InternalNotes string             `json:"internal_notes"`
	Lines         []InvoiceLineInput `json:"lines"`
}

type UpdateSupplierInvoiceRequest = CreateSupplierInvoiceRequest

type PaymentReceived struct {
	ID             int       `json:"id"`
	TenantID       int       `json:"tenant_id"`
	RestaurantID   int       `json:"restaurant_id"`
	PaymentNumber  string    `json:"payment_number"`
	PaymentDate    time.Time `json:"payment_date"`
	CustomerID     int       `json:"customer_id"`
	CustomerName   string    `json:"customer_name,omitempty"`
	InvoiceID      *int      `json:"invoice_id,omitempty"`
	Amount         float64   `json:"amount"`
	PaymentMethod  string    `json:"payment_method"`
	Reference      string    `json:"reference,omitempty"`
	BankAccountID  *int      `json:"bank_account_id,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	JournalEntryID *int      `json:"journal_entry_id,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      *int      `json:"created_by,omitempty"`
}

type CreatePaymentReceivedRequest struct {
	PaymentDate   string  `json:"payment_date"`
	CustomerID    int     `json:"customer_id"`
	InvoiceID     *int    `json:"invoice_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	Reference     string  `json:"reference"`
	BankAccountID *int    `json:"bank_account_id"`
	Notes         string  `json:"notes"`
}

type PaymentMade struct {
	ID             int       `json:"id"`
	TenantID       int       `json:"tenant_id"`
	RestaurantID   int       `json:"restaurant_id"`
	PaymentNumber  string    `json:"payment_number"`
	PaymentDate    time.Time `json:"payment_date"`
	SupplierID     int       `json:"supplier_id"`
	SupplierName   string    `json:"supplier_name,omitempty"`
	InvoiceID      *int      `json:"invoice_id,omitempty"`
	Amount         float64   `json:"amount"`
	PaymentMethod  string    `json:"payment_method"`
	Reference      string    `json:"reference,omitempty"`
	BankAccountID  *int      `json:"bank_account_id,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	JournalEntryID *int      `json:"journal_entry_id,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      *int      `json:"created_by,omitempty"`
}

type CreatePaymentMadeRequest struct {
	PaymentDate   string  `json:"payment_date"`
	SupplierID    int     `json:"supplier_id"`
	InvoiceID     *int    `json:"invoice_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	Reference     string  `json:"reference"`
	BankAccountID *int    `json:"bank_account_id"`
	Notes         string  `json:"notes"`
}
