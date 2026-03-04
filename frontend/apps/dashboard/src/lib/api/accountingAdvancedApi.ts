import { api } from '@/lib/api'
import type {
  BankReconciliationSession,
  BudgetVarianceRow,
  BreakEvenReport,
  EInvoiceSubmission,
  ExchangeRate,
  FinancialKPIReport,
  HRCostAnalysisReport,
  MenuProfitabilityReport,
  PredictiveCashFlowReport,
  RecurringJournalTemplate,
  SupplierPaymentSchedule,
} from '@/types/accounting-advanced'
import type { JournalEntryLine } from '@/types/accounting'

const extractData = <T>(payload: any): T => (payload?.data ?? payload) as T

export async function getBankReconciliationSessions(): Promise<BankReconciliationSession[]> {
  const response = await api.get('/accounting/bank-reconciliation/sessions')
  return extractData<BankReconciliationSession[]>(response.data)
}

export async function createBankReconciliationSession(payload: {
  bank_account_id: number
  statement_date: string
  statement_balance: number
  book_balance: number
  notes?: string
}): Promise<BankReconciliationSession> {
  const response = await api.post('/accounting/bank-reconciliation/sessions', payload)
  return extractData<BankReconciliationSession>(response.data)
}

export async function addBankReconciliationLine(
  sessionId: number,
  payload: {
    journal_entry_id?: number
    line_description?: string
    statement_amount: number
    book_amount: number
    match_status?: string
  }
): Promise<BankReconciliationSession> {
  const response = await api.post(`/accounting/bank-reconciliation/sessions/${sessionId}/lines`, payload)
  return extractData<BankReconciliationSession>(response.data)
}

export async function completeBankReconciliationSession(sessionId: number): Promise<BankReconciliationSession> {
  const response = await api.post(`/accounting/bank-reconciliation/sessions/${sessionId}/complete`)
  return extractData<BankReconciliationSession>(response.data)
}

export async function getRecurringJournals(): Promise<RecurringJournalTemplate[]> {
  const response = await api.get('/accounting/recurring-journals')
  return extractData<RecurringJournalTemplate[]>(response.data)
}

export async function createRecurringJournal(payload: {
  template_code: string
  template_name: string
  description?: string
  frequency: 'weekly' | 'monthly' | 'quarterly' | 'yearly'
  start_date: string
  next_run_date?: string
  auto_post?: boolean
  is_active?: boolean
  lines: JournalEntryLine[]
}): Promise<RecurringJournalTemplate> {
  const response = await api.post('/accounting/recurring-journals', payload)
  return extractData<RecurringJournalTemplate>(response.data)
}

export async function runRecurringJournal(id: number, runDate?: string): Promise<{ journal_entry_id: number }> {
  const response = await api.post(`/accounting/recurring-journals/${id}/run`, runDate ? { run_date: runDate } : {})
  return extractData<{ journal_entry_id: number }>(response.data)
}

export async function getBudgets(fiscalYear?: number): Promise<BudgetVarianceRow[]> {
  const response = await api.get('/accounting/budgets', { params: fiscalYear ? { fiscal_year: fiscalYear } : undefined })
  return extractData<BudgetVarianceRow[]>(response.data)
}

export async function upsertBudget(payload: {
  account_id: number
  fiscal_year: number
  period_month: number
  budget_amount: number
  currency?: string
  notes?: string
}): Promise<{ id: number }> {
  const response = await api.post('/accounting/budgets', payload)
  return extractData<{ id: number }>(response.data)
}

export async function getExchangeRates(baseCurrency?: string, quoteCurrency?: string): Promise<ExchangeRate[]> {
  const response = await api.get('/accounting/exchange-rates', {
    params: {
      base_currency: baseCurrency,
      quote_currency: quoteCurrency,
    },
  })
  return extractData<ExchangeRate[]>(response.data)
}

export async function upsertExchangeRate(payload: {
  base_currency: string
  quote_currency: string
  rate: number
  effective_date: string
  source?: string
  is_active?: boolean
}): Promise<{ id: number }> {
  const response = await api.post('/accounting/exchange-rates', payload)
  return extractData<{ id: number }>(response.data)
}

export async function convertCurrency(params: {
  base_currency: string
  quote_currency: string
  amount: number
  as_of?: string
}): Promise<{ converted: number; rate: number }> {
  const response = await api.get('/accounting/exchange-rates/convert', { params })
  return extractData<{ converted: number; rate: number }>(response.data)
}

export async function getSupplierPaymentSchedules(): Promise<SupplierPaymentSchedule[]> {
  const response = await api.get('/accounting/supplier-payment-schedules')
  return extractData<SupplierPaymentSchedule[]>(response.data)
}

export async function createSupplierPaymentSchedule(payload: {
  supplier_invoice_id?: number
  supplier_id: number
  scheduled_date: string
  amount: number
  currency?: string
  reminder_days?: number
  notes?: string
}): Promise<{ id: number }> {
  const response = await api.post('/accounting/supplier-payment-schedules', payload)
  return extractData<{ id: number }>(response.data)
}

export async function getSupplierPaymentReminders(dueDays: number = 7): Promise<SupplierPaymentSchedule[]> {
  const response = await api.get('/accounting/supplier-payment-reminders', { params: { due_days: dueDays } })
  return extractData<SupplierPaymentSchedule[]>(response.data)
}

export async function sendSupplierPaymentReminder(scheduleId: number): Promise<void> {
  await api.post(`/accounting/supplier-payment-schedules/${scheduleId}/send-reminder`)
}

export async function sendCustomerInvoiceEmail(id: number, payload: {
  to_email?: string
  subject?: string
  message?: string
  include_pdf?: boolean
}): Promise<{ message_id: string }> {
  const response = await api.post(`/accounting/customer-invoices/${id}/send-email`, payload)
  return extractData<{ message_id: string }>(response.data)
}

export async function submitEInvoice(invoiceId: number, invoiceType: 'customer' | 'supplier' = 'customer'): Promise<{ id: number; status: string }> {
  const response = await api.post(`/accounting/einvoice/${invoiceId}/submit`, { invoice_type: invoiceType })
  return extractData<{ id: number; status: string }>(response.data)
}

export async function getEInvoiceSubmissions(): Promise<EInvoiceSubmission[]> {
  const response = await api.get('/accounting/einvoice/submissions')
  return extractData<EInvoiceSubmission[]>(response.data)
}

export async function getFinancialKPIs(params: { date_from: string; date_to: string }): Promise<FinancialKPIReport> {
  const response = await api.get('/reports/financial/kpis', { params })
  return extractData<FinancialKPIReport>(response.data)
}

export async function getMenuProfitability(params: { date_from: string; date_to: string; limit?: number }): Promise<MenuProfitabilityReport> {
  const response = await api.get('/reports/financial/menu-profitability', { params })
  return extractData<MenuProfitabilityReport>(response.data)
}

export async function getBreakEvenReport(params: { date_from: string; date_to: string }): Promise<BreakEvenReport> {
  const response = await api.get('/reports/financial/break-even', { params })
  return extractData<BreakEvenReport>(response.data)
}

export async function getPredictiveCashFlow(params?: { history_months?: number; forecast_months?: number }): Promise<PredictiveCashFlowReport> {
  const response = await api.get('/reports/financial/predictive-cash-flow', { params })
  return extractData<PredictiveCashFlowReport>(response.data)
}

export async function getHRCostAnalysis(params: { date_from: string; date_to: string }): Promise<HRCostAnalysisReport> {
  const response = await api.get('/reports/hr/cost-analysis', { params })
  return extractData<HRCostAnalysisReport>(response.data)
}

export async function downloadCsv(path: string, params?: Record<string, string | number | boolean | undefined>): Promise<Blob> {
  const response = await api.get(path, { params, responseType: 'blob' })
  return response.data as Blob
}
