import type { Supplier, CustomerAccount } from '@/types/accounting'

export type BankAccountType = 'checking' | 'savings' | 'credit_card' | 'petty_cash'
export type InvoiceStatus =
  | 'draft'
  | 'sent'
  | 'approved'
  | 'partially_paid'
  | 'paid'
  | 'overdue'
  | 'void'
export type AccountingPaymentStatus = 'completed' | 'void'

export interface BankAccount {
  id: number
  tenant_id: number
  restaurant_id: number
  account_name: string
  account_name_ar?: string
  bank_name: string
  bank_branch?: string
  account_number: string
  iban?: string
  swift_code?: string
  currency: string
  account_type: BankAccountType
  opening_balance: number
  current_balance: number
  opening_date: string
  gl_account_id?: number | null
  is_active: boolean
  is_primary: boolean
  last_reconciled_date?: string | null
  last_reconciled_balance?: number | null
  created_at?: string
  updated_at?: string
  created_by?: number | null
  updated_by?: number | null
}

export interface BankAccountFormInput {
  account_name: string
  account_name_ar?: string
  bank_name: string
  bank_branch?: string
  account_number: string
  iban?: string
  swift_code?: string
  currency: string
  account_type: BankAccountType
  opening_balance: number
  opening_date: string
  gl_account_id?: number | null
  is_active?: boolean
  is_primary: boolean
}

export interface BankTransaction {
  transaction_date: string
  entry_number: string
  reference_number?: string
  description: string
  debit_amount: number
  credit_amount: number
  running_balance: number
  journal_entry_id: number
}

export interface InvoiceLine {
  id?: number
  invoice_id?: number
  tenant_id?: number
  description: string
  description_ar?: string
  quantity: number
  unit_price: number
  tax_rate: number
  tax_amount: number
  discount_rate?: number
  discount_amount: number
  line_total: number
  account_id?: number | null
  line_order: number
  created_at?: string
}

export interface InvoiceLineFormInput {
  description: string
  description_ar?: string
  quantity: number
  unit_price: number
  tax_rate: number
  tax_amount: number
  discount_rate?: number
  discount_amount: number
  line_total: number
  account_id?: number | null
  line_order: number
}

export interface CustomerInvoice {
  id: number
  tenant_id: number
  restaurant_id: number
  invoice_number: string
  invoice_date: string
  due_date: string
  customer_id: number
  customer_name?: string
  order_id?: number | null
  subtotal: number
  tax_amount: number
  discount_amount: number
  total_amount: number
  amount_paid: number
  balance_due: number
  currency: string
  status: InvoiceStatus
  notes?: string
  terms?: string
  internal_notes?: string
  journal_entry_id?: number | null
  created_at?: string
  updated_at?: string
  created_by?: number | null
  updated_by?: number | null
  lines?: InvoiceLine[]
  customer?: CustomerAccount
}

export interface CustomerInvoiceFormInput {
  invoice_date: string
  due_date: string
  customer_id: number
  order_id?: number | null
  currency: string
  status?: InvoiceStatus
  notes?: string
  terms?: string
  internal_notes?: string
  lines: InvoiceLineFormInput[]
}

export interface SupplierInvoice {
  id: number
  tenant_id: number
  restaurant_id: number
  invoice_number: string
  supplier_ref?: string
  invoice_date: string
  due_date: string
  supplier_id: number
  supplier_name?: string
  subtotal: number
  tax_amount: number
  discount_amount: number
  total_amount: number
  amount_paid: number
  balance_due: number
  currency: string
  status: InvoiceStatus
  notes?: string
  internal_notes?: string
  journal_entry_id?: number | null
  created_at?: string
  updated_at?: string
  created_by?: number | null
  updated_by?: number | null
  lines?: InvoiceLine[]
  supplier?: Supplier
}

export interface SupplierInvoiceFormInput {
  supplier_ref?: string
  invoice_date: string
  due_date: string
  supplier_id: number
  currency: string
  status?: InvoiceStatus
  notes?: string
  internal_notes?: string
  lines: InvoiceLineFormInput[]
}

export interface PaymentReceived {
  id: number
  tenant_id: number
  restaurant_id: number
  payment_number: string
  payment_date: string
  customer_id: number
  customer_name?: string
  invoice_id?: number | null
  amount: number
  payment_method: string
  reference?: string
  bank_account_id?: number | null
  notes?: string
  journal_entry_id?: number | null
  status: AccountingPaymentStatus
  created_at?: string
  updated_at?: string
  created_by?: number | null
}

export interface PaymentReceivedFormInput {
  payment_date: string
  customer_id: number
  invoice_id?: number | null
  amount: number
  payment_method: string
  reference?: string
  bank_account_id?: number | null
  notes?: string
}

export interface PaymentMade {
  id: number
  tenant_id: number
  restaurant_id: number
  payment_number: string
  payment_date: string
  supplier_id: number
  supplier_name?: string
  invoice_id?: number | null
  amount: number
  payment_method: string
  reference?: string
  bank_account_id?: number | null
  notes?: string
  journal_entry_id?: number | null
  status: AccountingPaymentStatus
  created_at?: string
  updated_at?: string
  created_by?: number | null
}

export interface PaymentMadeFormInput {
  payment_date: string
  supplier_id: number
  invoice_id?: number | null
  amount: number
  payment_method: string
  reference?: string
  bank_account_id?: number | null
  notes?: string
}
