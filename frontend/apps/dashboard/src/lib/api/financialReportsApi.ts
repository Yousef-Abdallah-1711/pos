import { api } from '@/lib/api'
import type {
  AgingReport,
  BalanceSheetReport,
  CashFlowReport,
  FinancialReportFilters,
  InventoryValuationReport,
  ProfitLossReport,
  TaxPeriodInput,
  TaxPeriod,
  TrialBalanceReport,
  VATReport,
} from '@/types/financial-reports'

export async function getProfitLossReport(filters: FinancialReportFilters): Promise<ProfitLossReport> {
  const response = await api.get('/reports/financial/profit-loss', { params: filters })
  return response.data?.data as ProfitLossReport
}

export async function getBalanceSheetReport(filters: FinancialReportFilters): Promise<BalanceSheetReport> {
  const response = await api.get('/reports/financial/balance-sheet', { params: filters })
  return response.data?.data as BalanceSheetReport
}

export async function getCashFlowReport(filters: FinancialReportFilters): Promise<CashFlowReport> {
  const response = await api.get('/reports/financial/cash-flow', { params: filters })
  return response.data?.data as CashFlowReport
}

export async function getReceivablesAgingReport(filters: FinancialReportFilters): Promise<AgingReport> {
  const response = await api.get('/reports/financial/receivables-aging', { params: filters })
  return response.data?.data as AgingReport
}

export async function getPayablesAgingReport(filters: FinancialReportFilters): Promise<AgingReport> {
  const response = await api.get('/reports/financial/payables-aging', { params: filters })
  return response.data?.data as AgingReport
}

export async function getInventoryValuationReport(filters: FinancialReportFilters): Promise<InventoryValuationReport> {
  const response = await api.get('/reports/financial/inventory-valuation', { params: filters })
  return response.data?.data as InventoryValuationReport
}

export async function getTrialBalanceReport(filters: FinancialReportFilters): Promise<TrialBalanceReport> {
  const response = await api.get('/reports/financial/trial-balance', { params: filters })
  return response.data?.data as TrialBalanceReport
}

export async function getVATReport(filters: FinancialReportFilters): Promise<VATReport> {
  const response = await api.get('/accounting/tax/vat-report', { params: filters })
  return response.data?.data as VATReport
}

export async function getTaxPeriods(): Promise<TaxPeriod[]> {
  const response = await api.get('/accounting/tax/periods')
  return (response.data?.data || []) as TaxPeriod[]
}

export async function getTaxPeriod(id: number): Promise<TaxPeriod> {
  const response = await api.get(`/accounting/tax/periods/${id}`)
  return response.data?.data as TaxPeriod
}

export async function createTaxPeriod(data: TaxPeriodInput): Promise<TaxPeriod> {
  const response = await api.post('/accounting/tax/periods', data)
  return response.data?.data as TaxPeriod
}

export async function updateTaxPeriod(id: number, data: TaxPeriodInput): Promise<TaxPeriod> {
  const response = await api.put(`/accounting/tax/periods/${id}`, data)
  return response.data?.data as TaxPeriod
}

export async function deleteTaxPeriod(id: number): Promise<void> {
  await api.delete(`/accounting/tax/periods/${id}`)
}
