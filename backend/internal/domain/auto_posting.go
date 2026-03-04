package domain

import (
	"encoding/json"
	"time"
)

const (
	AutoPostingEventOrderPaid       = "ORDER_PAID"
	AutoPostingEventOrderCOGS       = "ORDER_COGS"
	AutoPostingEventSalaryPaid      = "SALARY_PAID"
	AutoPostingEventSupplierInvoice = "SUPPLIER_INVOICE"
	AutoPostingEventSupplierPayment = "SUPPLIER_PAYMENT"
	AutoPostingEventDepreciation    = "DEPRECIATION"

	AutoPostingEntryTypeSales           = "auto_sales"
	AutoPostingEntryTypeCOGS            = "auto_cogs"
	AutoPostingEntryTypePayroll         = "auto_payroll"
	AutoPostingEntryTypeSupplierInvoice = "auto_supplier_invoice"
	AutoPostingEntryTypeSupplierPayment = "auto_supplier_payment"
	AutoPostingEntryTypeDepreciation    = "auto_depreciation"

	AutoPostingReferenceTypeOrder  = "order"
	AutoPostingReferenceTypeSalary = "salary"
	AutoPostingReferenceTypeAsset  = "asset"

	InventoryMovementTypeOpening    = "opening"
	InventoryMovementTypePurchase   = "purchase"
	InventoryMovementTypeAdjustment = "adjustment"
	InventoryMovementTypeSale       = "sale"
)

var DefaultPostingRuleAccounts = map[string]map[string]string{
	AutoPostingEventOrderPaid: {
		"cash":           "1110",
		"bank":           "1111",
		"card":           "1113",
		"online":         "1111",
		"wallet":         "1111",
		"receivable":     "1120",
		"sales":          "4100",
		"delivery_sales": "4120",
		"tax":            "2120",
		"discount":       "4150",
	},
	AutoPostingEventOrderCOGS: {
		"inventory":      "1130",
		"food_cogs":      "5110",
		"beverage_cogs":  "5120",
		"packaging_cogs": "5130",
	},
	AutoPostingEventSalaryPaid: {
		"cash":                      "1110",
		"bank":                      "1111",
		"salary_expense":            "5210",
		"overtime_expense":          "5220",
		"benefits_expense":          "5230",
		"salaries_payable":          "2130",
		"tips_payable":              "2140",
		"employee_advances":         "1121",
		"accrued_expenses":          "2110",
		"tax_payable":               "2110",
		"social_insurance_payable":  "2110",
		"health_insurance_payable":  "2110",
		"pension_payable":           "2110",
		"loan_clearing":             "2110",
		"other_deductions_payable":  "2110",
	},
	AutoPostingEventSupplierInvoice: {
		"accounts_payable": "2101",
		"inventory":        "1130",
		"input_vat":        "1150",
	},
	AutoPostingEventSupplierPayment: {
		"cash":             "1110",
		"bank":             "1111",
		"accounts_payable": "2101",
	},
	AutoPostingEventDepreciation: {
		"depreciation_expense":    "5550",
		"accumulated_depreciation": "1240",
	},
}

type AutoPostingRule struct {
	ID              int               `json:"id"`
	TenantID        int               `json:"tenant_id"`
	RestaurantID    int               `json:"restaurant_id"`
	RuleCode        string            `json:"rule_code"`
	RuleName        string            `json:"rule_name"`
	Description     string            `json:"description,omitempty"`
	IsActive        bool              `json:"is_active"`
	DefaultAccounts map[string]string `json:"default_accounts"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type PostingEvent struct {
	TenantID        int             `json:"tenant_id,omitempty"`
	RestaurantID    int             `json:"restaurant_id,omitempty"`
	Code            string          `json:"code"`
	EntryType       string          `json:"entry_type"`
	ReferenceType   string          `json:"reference_type"`
	ReferenceID     int             `json:"reference_id"`
	ReferenceNumber string          `json:"reference_number,omitempty"`
	Description     string          `json:"description"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	PostedBy        *int            `json:"posted_by,omitempty"`
}

type InventoryMovement struct {
	ID            int       `json:"id"`
	TenantID      int       `json:"tenant_id"`
	RestaurantID  int       `json:"restaurant_id"`
	ProductID     int64     `json:"product_id"`
	OrderID       *int64    `json:"order_id,omitempty"`
	ReferenceType string    `json:"reference_type,omitempty"`
	ReferenceID   *int      `json:"reference_id,omitempty"`
	MovementType  string    `json:"movement_type"`
	Quantity      float64   `json:"quantity"`
	UnitCost      float64   `json:"unit_cost"`
	TotalCost     float64   `json:"total_cost"`
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     *int      `json:"created_by,omitempty"`
}

type GeneralLedgerLine struct {
	AccountID        int       `json:"account_id"`
	AccountCode      string    `json:"account_code"`
	AccountName      string    `json:"account_name"`
	JournalEntryID   int       `json:"journal_entry_id"`
	JournalLineID    int       `json:"journal_line_id"`
	EntryDate        time.Time `json:"entry_date"`
	EntryNumber      string    `json:"entry_number"`
	Description      string    `json:"description"`
	ReferenceType    string    `json:"reference_type,omitempty"`
	ReferenceNumber  string    `json:"reference_number,omitempty"`
	DebitAmount      float64   `json:"debit_amount"`
	CreditAmount     float64   `json:"credit_amount"`
	RunningBalance   float64   `json:"running_balance"`
	LineDescription  string    `json:"line_description,omitempty"`
}

type TrialBalanceRow struct {
	AccountID       int       `json:"account_id"`
	AccountCode     string    `json:"account_code"`
	AccountName     string    `json:"account_name"`
	AccountType     string    `json:"account_type"`
	NormalBalance   string    `json:"normal_balance"`
	DebitBalance    float64   `json:"debit_balance"`
	CreditBalance   float64   `json:"credit_balance"`
	NetBalance      float64   `json:"net_balance"`
	AsOfDate        time.Time `json:"as_of_date"`
}
