export interface FinancialReportFilters {
  date_from?: string
  date_to?: string
  as_of?: string
  compare_date_from?: string
  compare_date_to?: string
}

export interface FinancialLine {
  account_code: string
  account_name: string
  amount: number
}

export interface ProfitLossReport {
  date_from: string
  date_to: string
  revenue: FinancialLine[]
  expenses: FinancialLine[]
  total_revenue: number
  total_expenses: number
  net_profit: number
  comparison_date_from?: string
  comparison_date_to?: string
  comparison_total_revenue?: number
  comparison_total_expenses?: number
  comparison_net_profit?: number
  revenue_change_percent?: number
  expense_change_percent?: number
  net_profit_change_percent?: number
}

export interface TrialBalanceRow {
  account_id: number
  account_code: string
  account_name: string
  account_type: string
  normal_balance: string
  debit_balance: number
  credit_balance: number
  net_balance: number
}

export interface TrialBalanceReport {
  as_of_date: string
  rows: TrialBalanceRow[]
  total_debits: number
  total_credits: number
  is_balanced: boolean
}

export interface BalanceSheetReport {
  as_of_date: string
  assets: FinancialLine[]
  liabilities: FinancialLine[]
  equity: FinancialLine[]
  total_assets: number
  total_liabilities: number
  total_equity: number
  equation_balanced: boolean
  integrity_status: 'RECONCILED' | 'UNRECONCILED' | 'NOT_VERIFIED'
  equation_difference: number
  integrity_check_basis: string
}

export interface CashFlowReport {
  date_from: string
  date_to: string
  operating_cash_flow: number
  investing_cash_flow: number
  financing_cash_flow: number
  net_cash_flow: number
}

export interface AgingBucket {
  entity_id: number
  entity_name: string
  current: number
  days_30: number
  days_60: number
  days_90: number
  days_90_plus: number
  total: number
}

export interface AgingReport {
  as_of_date: string
  rows: AgingBucket[]
  total: number
}

export interface InventoryValuationRow {
  product_id: number
  product_name: string
  quantity: number
  cost_per_unit: number
  total_value: number
}

export interface InventoryValuationReport {
  as_of_date: string
  rows: InventoryValuationRow[]
  total_value: number
}

export interface VATReport {
  date_from: string
  date_to: string
  output_vat: number
  input_vat: number
  net_vat: number
}

export interface TaxPeriod {
  id: number
  period_name: string
  period_start: string
  period_end: string
  output_vat: number
  input_vat: number
  net_vat: number
  status: string
  closed_at?: string | null
  notes?: string
}

export interface TaxPeriodInput {
  period_name: string
  period_start: string
  period_end: string
  notes?: string
}
