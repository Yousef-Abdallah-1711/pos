package domain

import (
	"encoding/json"
	"time"
)

const (
	AccountTypeAsset     = "asset"
	AccountTypeLiability = "liability"
	AccountTypeEquity    = "equity"
	AccountTypeRevenue   = "revenue"
	AccountTypeExpense   = "expense"

	NormalBalanceDebit  = "debit"
	NormalBalanceCredit = "credit"

	JournalStatusDraft  = "draft"
	JournalStatusPosted = "posted"
	JournalStatusVoid   = "void"
)

type Account struct {
	ID             int       `json:"id"`
	TenantID       int       `json:"tenant_id"`
	RestaurantID   int       `json:"restaurant_id"`
	AccountCode    string    `json:"account_code"`
	AccountName    string    `json:"account_name"`
	AccountNameAr  string    `json:"account_name_ar,omitempty"`
	Description    string    `json:"description,omitempty"`
	DescriptionAr  string    `json:"description_ar,omitempty"`
	AccountType    string    `json:"account_type"`
	AccountSubtype string    `json:"account_subtype,omitempty"`
	NormalBalance  string    `json:"normal_balance"`
	ParentID       *int      `json:"parent_id,omitempty"`
	Level          int       `json:"level"`
	IsHeader       bool      `json:"is_header"`
	CurrentBalance float64   `json:"current_balance"`
	Currency       string    `json:"currency"`
	IsActive       bool      `json:"is_active"`
	IsSystem       bool      `json:"is_system"`
	IsReconcilable bool      `json:"is_reconcilable"`
	TaxRate        *float64  `json:"tax_rate,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      *int      `json:"created_by,omitempty"`
	UpdatedBy      *int      `json:"updated_by,omitempty"`
	Children       []Account `json:"children,omitempty"`
}

type CreateAccountRequest struct {
	AccountCode    string   `json:"account_code"`
	AccountName    string   `json:"account_name"`
	AccountNameAr  string   `json:"account_name_ar"`
	Description    string   `json:"description"`
	DescriptionAr  string   `json:"description_ar"`
	AccountType    string   `json:"account_type"`
	AccountSubtype string   `json:"account_subtype"`
	NormalBalance  string   `json:"normal_balance"`
	ParentID       *int     `json:"parent_id"`
	IsHeader       bool     `json:"is_header"`
	Currency       string   `json:"currency"`
	IsActive       *bool    `json:"is_active"`
	IsReconcilable bool     `json:"is_reconcilable"`
	TaxRate        *float64 `json:"tax_rate"`
}

type UpdateAccountRequest = CreateAccountRequest

type JournalEntry struct {
	ID              int                `json:"id"`
	TenantID        int                `json:"tenant_id"`
	RestaurantID    int                `json:"restaurant_id"`
	EntryNumber     string             `json:"entry_number"`
	EntryDate       time.Time          `json:"entry_date"`
	FiscalPeriod    string             `json:"fiscal_period"`
	Description     string             `json:"description"`
	DescriptionAr   string             `json:"description_ar,omitempty"`
	ReferenceType   string             `json:"reference_type,omitempty"`
	ReferenceID     *int               `json:"reference_id,omitempty"`
	ReferenceNumber string             `json:"reference_number,omitempty"`
	EntryType       string             `json:"entry_type"`
	TotalDebit      float64            `json:"total_debit"`
	TotalCredit     float64            `json:"total_credit"`
	Status          string             `json:"status"`
	PostedAt        *time.Time         `json:"posted_at,omitempty"`
	PostedBy        *int               `json:"posted_by,omitempty"`
	VoidedAt        *time.Time         `json:"voided_at,omitempty"`
	VoidedBy        *int               `json:"voided_by,omitempty"`
	VoidReason      string             `json:"void_reason,omitempty"`
	IsApproved      bool               `json:"is_approved"`
	ApprovedBy      *int               `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time         `json:"approved_at,omitempty"`
	Notes           string             `json:"notes,omitempty"`
	Attachments     json.RawMessage    `json:"attachments,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	CreatedBy       *int               `json:"created_by,omitempty"`
	UpdatedBy       *int               `json:"updated_by,omitempty"`
	Lines           []JournalEntryLine `json:"lines,omitempty"`
}

type JournalEntryLine struct {
	ID             int       `json:"id"`
	JournalEntryID int       `json:"journal_entry_id"`
	TenantID       int       `json:"tenant_id"`
	AccountID      int       `json:"account_id"`
	AccountCode    string    `json:"account_code"`
	DebitAmount    float64   `json:"debit_amount"`
	CreditAmount   float64   `json:"credit_amount"`
	Description    string    `json:"description,omitempty"`
	DescriptionAr  string    `json:"description_ar,omitempty"`
	LineOrder      int       `json:"line_order"`
	TaxCode        string    `json:"tax_code,omitempty"`
	TaxAmount      float64   `json:"tax_amount"`
	CreatedAt      time.Time `json:"created_at"`
}

type JournalEntryLineInput struct {
	AccountID     int     `json:"account_id"`
	DebitAmount   float64 `json:"debit_amount"`
	CreditAmount  float64 `json:"credit_amount"`
	Description   string  `json:"description"`
	DescriptionAr string  `json:"description_ar"`
	LineOrder     int     `json:"line_order"`
	TaxCode       string  `json:"tax_code"`
	TaxAmount     float64 `json:"tax_amount"`
}

type CreateJournalEntryRequest struct {
	EntryDate       string                  `json:"entry_date"`
	Description     string                  `json:"description"`
	DescriptionAr   string                  `json:"description_ar"`
	ReferenceType   string                  `json:"reference_type"`
	ReferenceID     *int                    `json:"reference_id"`
	ReferenceNumber string                  `json:"reference_number"`
	EntryType       string                  `json:"entry_type"`
	Notes           string                  `json:"notes"`
	Attachments     json.RawMessage         `json:"attachments"`
	Lines           []JournalEntryLineInput `json:"lines"`
}

type UpdateJournalEntryRequest = CreateJournalEntryRequest

type JournalEntryFilters struct {
	StartDate *time.Time
	EndDate   *time.Time
	Status    string
	EntryType string
}

type GeneralLedgerEntry struct {
	ID             int       `json:"id"`
	TenantID       int       `json:"tenant_id"`
	RestaurantID   int       `json:"restaurant_id"`
	AccountID      int       `json:"account_id"`
	JournalEntryID int       `json:"journal_entry_id"`
	JournalLineID  int       `json:"journal_line_id"`
	EntryDate      time.Time `json:"entry_date"`
	FiscalPeriod   string    `json:"fiscal_period"`
	DebitAmount    float64   `json:"debit_amount"`
	CreditAmount   float64   `json:"credit_amount"`
	RunningBalance float64   `json:"running_balance"`
	CreatedAt      time.Time `json:"created_at"`
}

type CustomerAccount struct {
	ID                  int       `json:"id"`
	TenantID            int       `json:"tenant_id"`
	RestaurantID        int       `json:"restaurant_id"`
	CustomerCode        string    `json:"customer_code"`
	CustomerName        string    `json:"customer_name"`
	CustomerNameAr      string    `json:"customer_name_ar,omitempty"`
	CustomerType        string    `json:"customer_type"`
	ContactPerson       string    `json:"contact_person,omitempty"`
	Email               string    `json:"email,omitempty"`
	Phone               string    `json:"phone,omitempty"`
	Mobile              string    `json:"mobile,omitempty"`
	Address             string    `json:"address,omitempty"`
	AddressAr           string    `json:"address_ar,omitempty"`
	City                string    `json:"city,omitempty"`
	Country             string    `json:"country,omitempty"`
	TaxNumber           string    `json:"tax_number,omitempty"`
	PaymentTerms        string    `json:"payment_terms"`
	CreditLimit         float64   `json:"credit_limit"`
	CurrentBalance      float64   `json:"current_balance"`
	Currency            string    `json:"currency"`
	ReceivableAccountID *int      `json:"receivable_account_id,omitempty"`
	Notes               string    `json:"notes,omitempty"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	CreatedBy           *int      `json:"created_by,omitempty"`
	UpdatedBy           *int      `json:"updated_by,omitempty"`
}

type CreateCustomerAccountRequest struct {
	CustomerName        string  `json:"customer_name"`
	CustomerNameAr      string  `json:"customer_name_ar"`
	CustomerType        string  `json:"customer_type"`
	ContactPerson       string  `json:"contact_person"`
	Email               string  `json:"email"`
	Phone               string  `json:"phone"`
	Mobile              string  `json:"mobile"`
	Address             string  `json:"address"`
	AddressAr           string  `json:"address_ar"`
	City                string  `json:"city"`
	Country             string  `json:"country"`
	TaxNumber           string  `json:"tax_number"`
	PaymentTerms        string  `json:"payment_terms"`
	CreditLimit         float64 `json:"credit_limit"`
	Currency            string  `json:"currency"`
	ReceivableAccountID *int    `json:"receivable_account_id"`
	Notes               string  `json:"notes"`
	IsActive            *bool   `json:"is_active"`
}

type UpdateCustomerAccountRequest = CreateCustomerAccountRequest

type CustomerStatementLine struct {
	EntryDate       time.Time `json:"entry_date"`
	EntryNumber     string    `json:"entry_number"`
	ReferenceNumber string    `json:"reference_number,omitempty"`
	Description     string    `json:"description"`
	DebitAmount     float64   `json:"debit_amount"`
	CreditAmount    float64   `json:"credit_amount"`
	RunningBalance  float64   `json:"running_balance"`
}
