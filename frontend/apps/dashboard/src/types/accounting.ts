export type AccountType = 'asset' | 'liability' | 'equity' | 'revenue' | 'expense'
export type NormalBalance = 'debit' | 'credit'
export type JournalEntryStatus = 'draft' | 'posted' | 'void'

export interface Account {
  id: number
  tenant_id: number
  restaurant_id: number
  account_code: string
  account_name: string
  account_name_ar?: string
  description?: string
  description_ar?: string
  account_type: AccountType
  account_subtype?: string
  normal_balance: NormalBalance
  parent_id?: number | null
  level: number
  is_header: boolean
  current_balance: number
  currency: string
  is_active: boolean
  is_system: boolean
  is_reconcilable: boolean
  tax_rate?: number | null
  created_at?: string
  updated_at?: string
  children?: Account[]
}

export interface AccountFormInput {
  account_code: string
  account_name: string
  account_name_ar?: string
  description?: string
  description_ar?: string
  account_type: AccountType
  account_subtype?: string
  normal_balance: NormalBalance
  parent_id?: number | null
  is_header: boolean
  currency: string
  is_active?: boolean
  is_reconcilable: boolean
  tax_rate?: number | null
}

export interface JournalEntryLine {
  id?: number
  journal_entry_id?: number
  tenant_id?: number
  account_id: number
  account_code?: string
  debit_amount: number
  credit_amount: number
  description?: string
  description_ar?: string
  line_order: number
  tax_code?: string
  tax_amount?: number
  created_at?: string
}

export interface JournalEntry {
  id: number
  tenant_id: number
  restaurant_id: number
  entry_number: string
  entry_date: string
  fiscal_period?: string
  description: string
  description_ar?: string
  reference_type?: string
  reference_id?: number | null
  reference_number?: string
  entry_type: string
  total_debit: number
  total_credit: number
  status: JournalEntryStatus
  posted_at?: string | null
  posted_by?: number | null
  voided_at?: string | null
  voided_by?: number | null
  void_reason?: string
  is_approved?: boolean
  approved_by?: number | null
  approved_at?: string | null
  notes?: string
  attachments?: unknown
  created_at?: string
  updated_at?: string
  created_by?: number | null
  updated_by?: number | null
  lines: JournalEntryLine[]
}

export interface JournalEntryFormInput {
  entry_date: string
  description: string
  description_ar?: string
  reference_type?: string
  reference_id?: number | null
  reference_number?: string
  entry_type: string
  notes?: string
  attachments?: unknown
  lines: JournalEntryLine[]
}

export interface JournalEntryFilters {
  start_date?: string
  end_date?: string
  status?: string
  entry_type?: string
}

export interface Supplier {
  id: number
  tenant_id: number
  restaurant_id: number
  supplier_code: string
  company_name: string
  company_name_ar?: string
  contact_person?: string
  email?: string
  phone?: string
  mobile?: string
  website?: string
  address?: string
  address_ar?: string
  city?: string
  state?: string
  postal_code?: string
  country?: string
  tax_number?: string
  payment_terms?: string
  credit_limit: number
  current_balance: number
  currency: string
  payable_account_id?: number | null
  category?: string
  bank_name?: string
  bank_account?: string
  iban?: string
  notes?: string
  is_active: boolean
  created_at?: string
  updated_at?: string
}

export interface SupplierFormInput {
  company_name: string
  company_name_ar?: string
  contact_person?: string
  email?: string
  phone?: string
  mobile?: string
  website?: string
  address?: string
  address_ar?: string
  city?: string
  state?: string
  postal_code?: string
  country?: string
  tax_number?: string
  payment_terms?: string
  credit_limit: number
  currency: string
  payable_account_id?: number | null
  category?: string
  bank_name?: string
  bank_account?: string
  iban?: string
  notes?: string
  is_active?: boolean
}

export interface CustomerAccount {
  id: number
  tenant_id: number
  restaurant_id: number
  customer_code: string
  customer_name: string
  customer_name_ar?: string
  customer_type: string
  contact_person?: string
  email?: string
  phone?: string
  mobile?: string
  address?: string
  address_ar?: string
  city?: string
  country?: string
  tax_number?: string
  payment_terms?: string
  credit_limit: number
  current_balance: number
  currency: string
  receivable_account_id?: number | null
  notes?: string
  is_active: boolean
  created_at?: string
  updated_at?: string
}

export interface CustomerAccountFormInput {
  customer_name: string
  customer_name_ar?: string
  customer_type?: string
  contact_person?: string
  email?: string
  phone?: string
  mobile?: string
  address?: string
  address_ar?: string
  city?: string
  country?: string
  tax_number?: string
  payment_terms?: string
  credit_limit: number
  currency: string
  receivable_account_id?: number | null
  notes?: string
  is_active?: boolean
}

export interface CustomerStatementLine {
  entry_date: string
  entry_number: string
  reference_number?: string
  description: string
  debit_amount: number
  credit_amount: number
  running_balance: number
}

export interface GeneralLedgerLine {
  account_id: number
  account_code: string
  account_name: string
  journal_entry_id: number
  journal_line_id: number
  entry_date: string
  entry_number: string
  description: string
  reference_type?: string
  reference_number?: string
  debit_amount: number
  credit_amount: number
  running_balance: number
  line_description?: string
}

export interface LedgerResponse {
  account_id: number
  lines: GeneralLedgerLine[]
  balance: number
  next_cursor?: number | null
  limit?: number
  period_summaries?: LedgerPeriodSummary[]
}

export interface LedgerPeriodSummary {
  fiscal_period: string
  entry_count: number
  debit_total: number
  credit_total: number
  net_change: number
  closing_balance: number
}

export interface TrialBalanceRow {
  account_id: number
  account_code: string
  account_name: string
  account_type: AccountType
  normal_balance?: NormalBalance
  debit_balance: number
  credit_balance: number
  net_balance?: number
  as_of_date?: string
}

export interface TrialBalanceResponse {
  as_of: string
  rows: TrialBalanceRow[]
  total_debit: number
  total_credit: number
  is_balanced: boolean
}
