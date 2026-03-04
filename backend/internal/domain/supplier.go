package domain

import "time"

type Supplier struct {
	ID               int       `json:"id"`
	TenantID         int       `json:"tenant_id"`
	RestaurantID     int       `json:"restaurant_id"`
	SupplierCode     string    `json:"supplier_code"`
	CompanyName      string    `json:"company_name"`
	CompanyNameAr    string    `json:"company_name_ar,omitempty"`
	ContactPerson    string    `json:"contact_person,omitempty"`
	Email            string    `json:"email,omitempty"`
	Phone            string    `json:"phone,omitempty"`
	Mobile           string    `json:"mobile,omitempty"`
	Website          string    `json:"website,omitempty"`
	Address          string    `json:"address,omitempty"`
	AddressAr        string    `json:"address_ar,omitempty"`
	City             string    `json:"city,omitempty"`
	State            string    `json:"state,omitempty"`
	PostalCode       string    `json:"postal_code,omitempty"`
	Country          string    `json:"country,omitempty"`
	TaxNumber        string    `json:"tax_number,omitempty"`
	PaymentTerms     string    `json:"payment_terms"`
	CreditLimit      float64   `json:"credit_limit"`
	CurrentBalance   float64   `json:"current_balance"`
	Currency         string    `json:"currency"`
	PayableAccountID *int      `json:"payable_account_id,omitempty"`
	Category         string    `json:"category,omitempty"`
	BankName         string    `json:"bank_name,omitempty"`
	BankAccount      string    `json:"bank_account,omitempty"`
	IBAN             string    `json:"iban,omitempty"`
	Notes            string    `json:"notes,omitempty"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedBy        *int      `json:"created_by,omitempty"`
	UpdatedBy        *int      `json:"updated_by,omitempty"`
}

type CreateSupplierRequest struct {
	CompanyName      string  `json:"company_name"`
	CompanyNameAr    string  `json:"company_name_ar"`
	ContactPerson    string  `json:"contact_person"`
	Email            string  `json:"email"`
	Phone            string  `json:"phone"`
	Mobile           string  `json:"mobile"`
	Website          string  `json:"website"`
	Address          string  `json:"address"`
	AddressAr        string  `json:"address_ar"`
	City             string  `json:"city"`
	State            string  `json:"state"`
	PostalCode       string  `json:"postal_code"`
	Country          string  `json:"country"`
	TaxNumber        string  `json:"tax_number"`
	PaymentTerms     string  `json:"payment_terms"`
	CreditLimit      float64 `json:"credit_limit"`
	Currency         string  `json:"currency"`
	PayableAccountID *int    `json:"payable_account_id"`
	Category         string  `json:"category"`
	BankName         string  `json:"bank_name"`
	BankAccount      string  `json:"bank_account"`
	IBAN             string  `json:"iban"`
	Notes            string  `json:"notes"`
	IsActive         *bool   `json:"is_active"`
}

type UpdateSupplierRequest = CreateSupplierRequest
