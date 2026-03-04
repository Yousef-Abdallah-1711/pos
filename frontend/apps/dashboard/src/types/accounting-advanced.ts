export interface BankReconciliationSession {
  id: number
  bank_account_id: number
  statement_date: string
  statement_balance: number
  book_balance: number
  difference: number
  status: 'draft' | 'in_progress' | 'completed' | 'cancelled'
  notes?: string
  completed_at?: string | null
}

export interface RecurringJournalTemplate {
  id: number
  template_code: string
  template_name: string
  description?: string
  frequency: 'weekly' | 'monthly' | 'quarterly' | 'yearly'
  start_date: string
  next_run_date: string
  last_run_date?: string | null
  auto_post: boolean
  is_active: boolean
}

export interface BudgetVarianceRow {
  id: number
  account_id: number
  account_code: string
  account_name: string
  fiscal_year: number
  period_month: number
  budget_amount: number
  actual_amount: number
  variance_amount: number
  variance_pct: number
  currency: string
}

export interface ExchangeRate {
  id: number
  base_currency: string
  quote_currency: string
  rate: number
  effective_date: string
  source?: string
  is_active: boolean
}

export interface SupplierPaymentSchedule {
  id: number
  supplier_invoice_id?: number | null
  supplier_id: number
  supplier_name?: string
  scheduled_date: string
  amount: number
  currency: string
  reminder_days: number
  status: 'scheduled' | 'reminded' | 'paid' | 'cancelled'
  reminder_sent_at?: string | null
  notes?: string
}

export interface EInvoiceSubmission {
  id: number
  invoice_type: 'customer' | 'supplier'
  invoice_id: number
  eta_uuid?: string
  status: 'pending' | 'accepted' | 'rejected' | 'error'
  submitted_at: string
  error_message?: string
}

export interface FinancialKPIReport {
  summary: {
    revenue: number
    expenses: number
    net_profit: number
    cash_position: number
  }
  revenue_trend: Array<{ month: string; revenue: number }>
  expense_breakdown: Array<{ label: string; amount: number }>
  generated_at: string
}

export interface BreakEvenReport {
  summary: {
    revenue: number
    variable_costs: number
    fixed_costs: number
    contribution_margin: number
    contribution_margin_ratio: number
    break_even_revenue: number
    break_even_gap: number
    is_above_break_even: boolean
  }
  generated_at: string
}

export interface PredictiveCashFlowReport {
  history: Array<{ month: string; cash_in: number; cash_out: number; net_cash: number }>
  forecast: Array<{ month: string; projected_net_cash: number }>
  assumptions: {
    history_months: number
    forecast_months: number
    average_net_cash: number
    trend_slope: number
  }
  generated_at: string
}

export interface HRCostAnalysisReport {
  summary: {
    revenue: number
    payroll_cost: number
    overtime_cost: number
    estimated_insurance_cost: number
    total_hr_cost: number
    labor_cost_pct_of_revenue: number
  }
  by_department: Array<{ department: string; payroll_cost: number }>
  generated_at: string
}

export interface MenuProfitabilityReport {
  items: Array<{
    product_id: number
    name: string
    quantity: number
    revenue: number
    cogs: number
    gross_profit: number
    gross_margin_pct: number
  }>
  generated_at: string
}
