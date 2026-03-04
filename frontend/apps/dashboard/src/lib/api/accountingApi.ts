import { api } from '@/lib/api'
import type {
  Account,
  AccountFormInput,
  JournalEntry,
  JournalEntryFilters,
  JournalEntryFormInput,
  LedgerResponse,
  TrialBalanceResponse,
} from '@/types/accounting'
import type {
  CustomerInvoice,
  CustomerInvoiceFormInput,
  PaymentMade,
  PaymentMadeFormInput,
  PaymentReceived,
  PaymentReceivedFormInput,
  SupplierInvoice,
  SupplierInvoiceFormInput,
} from '@/types/invoices'

const extractData = <T>(payload: any): T => (payload?.data ?? payload) as T

export async function getAccounts(params?: {
  tree?: boolean
  include_inactive?: boolean
}): Promise<Account[]> {
  const response = await api.get('/accounting/accounts', { params })
  return extractData<Account[]>(response.data)
}

export async function getAccount(id: number): Promise<Account> {
  const response = await api.get(`/accounting/accounts/${id}`)
  return extractData<Account>(response.data)
}

export async function createAccount(payload: AccountFormInput): Promise<Account> {
  const response = await api.post('/accounting/accounts', payload)
  return extractData<Account>(response.data)
}

export async function updateAccount(
  id: number,
  payload: AccountFormInput
): Promise<Account> {
  const response = await api.put(`/accounting/accounts/${id}`, payload)
  return extractData<Account>(response.data)
}

export async function deleteAccount(id: number): Promise<void> {
  await api.delete(`/accounting/accounts/${id}`)
}

export async function getJournalEntries(
  filters?: JournalEntryFilters
): Promise<JournalEntry[]> {
  const response = await api.get('/accounting/journal-entries', {
    params: filters,
  })
  return extractData<JournalEntry[]>(response.data)
}

export async function getJournalEntry(id: number): Promise<JournalEntry> {
  const response = await api.get(`/accounting/journal-entries/${id}`)
  return extractData<JournalEntry>(response.data)
}

export async function createJournalEntry(
  payload: JournalEntryFormInput
): Promise<JournalEntry> {
  const response = await api.post('/accounting/journal-entries', payload)
  return extractData<JournalEntry>(response.data)
}

export async function updateJournalEntry(
  id: number,
  payload: JournalEntryFormInput
): Promise<JournalEntry> {
  const response = await api.put(`/accounting/journal-entries/${id}`, payload)
  return extractData<JournalEntry>(response.data)
}

export async function postJournalEntry(id: number): Promise<JournalEntry> {
  const response = await api.post(`/accounting/journal-entries/${id}/post`)
  return extractData<JournalEntry>(response.data)
}

export async function voidJournalEntry(
  id: number,
  reason?: string
): Promise<JournalEntry> {
  const response = await api.post(`/accounting/journal-entries/${id}/void`, {
    reason,
  })
  return extractData<JournalEntry>(response.data)
}

export async function getLedger(params: {
  account_id: number
  start?: string
  end?: string
  limit?: number
  cursor?: number
}): Promise<LedgerResponse> {
  const response = await api.get('/accounting/ledger', { params })
  return extractData<LedgerResponse>(response.data)
}

export async function getTrialBalance(params?: {
  as_of?: string
}): Promise<TrialBalanceResponse> {
  const response = await api.get('/accounting/trial-balance', { params })
  return extractData<TrialBalanceResponse>(response.data)
}

export async function getCustomerInvoices(): Promise<CustomerInvoice[]> {
  const response = await api.get('/accounting/customer-invoices')
  return extractData<CustomerInvoice[]>(response.data)
}

export async function getCustomerInvoice(id: number): Promise<CustomerInvoice> {
  const response = await api.get(`/accounting/customer-invoices/${id}`)
  return extractData<CustomerInvoice>(response.data)
}

export async function getOverdueCustomerInvoices(): Promise<CustomerInvoice[]> {
  const response = await api.get('/accounting/customer-invoices/overdue')
  return extractData<CustomerInvoice[]>(response.data)
}

export async function createCustomerInvoice(
  payload: CustomerInvoiceFormInput
): Promise<CustomerInvoice> {
  const response = await api.post('/accounting/customer-invoices', payload)
  return extractData<CustomerInvoice>(response.data)
}

export async function updateCustomerInvoice(
  id: number,
  payload: CustomerInvoiceFormInput
): Promise<CustomerInvoice> {
  const response = await api.put(`/accounting/customer-invoices/${id}`, payload)
  return extractData<CustomerInvoice>(response.data)
}

export async function sendCustomerInvoice(id: number): Promise<CustomerInvoice> {
  const response = await api.post(`/accounting/customer-invoices/${id}/send`)
  return extractData<CustomerInvoice>(response.data)
}

export async function voidCustomerInvoice(
  id: number,
  reason?: string
): Promise<CustomerInvoice> {
  const response = await api.post(`/accounting/customer-invoices/${id}/void`, {
    reason,
  })
  return extractData<CustomerInvoice>(response.data)
}

export async function getSupplierInvoices(): Promise<SupplierInvoice[]> {
  const response = await api.get('/accounting/supplier-invoices')
  return extractData<SupplierInvoice[]>(response.data)
}

export async function getSupplierInvoice(id: number): Promise<SupplierInvoice> {
  const response = await api.get(`/accounting/supplier-invoices/${id}`)
  return extractData<SupplierInvoice>(response.data)
}

export async function createSupplierInvoice(
  payload: SupplierInvoiceFormInput
): Promise<SupplierInvoice> {
  const response = await api.post('/accounting/supplier-invoices', payload)
  return extractData<SupplierInvoice>(response.data)
}

export async function updateSupplierInvoice(
  id: number,
  payload: SupplierInvoiceFormInput
): Promise<SupplierInvoice> {
  const response = await api.put(`/accounting/supplier-invoices/${id}`, payload)
  return extractData<SupplierInvoice>(response.data)
}

export async function approveSupplierInvoice(id: number): Promise<SupplierInvoice> {
  const response = await api.post(`/accounting/supplier-invoices/${id}/approve`)
  return extractData<SupplierInvoice>(response.data)
}

export async function voidSupplierInvoice(
  id: number,
  reason?: string
): Promise<SupplierInvoice> {
  const response = await api.post(`/accounting/supplier-invoices/${id}/void`, {
    reason,
  })
  return extractData<SupplierInvoice>(response.data)
}

export async function getPaymentsReceived(): Promise<PaymentReceived[]> {
  const response = await api.get('/accounting/payments/received')
  return extractData<PaymentReceived[]>(response.data)
}

export async function createPaymentReceived(
  payload: PaymentReceivedFormInput
): Promise<{ id: number }> {
  const response = await api.post('/accounting/payments/received', payload)
  return extractData<{ id: number }>(response.data)
}

export async function getPaymentsMade(): Promise<PaymentMade[]> {
  const response = await api.get('/accounting/payments/made')
  return extractData<PaymentMade[]>(response.data)
}

export async function createPaymentMade(
  payload: PaymentMadeFormInput
): Promise<{ id: number }> {
  const response = await api.post('/accounting/payments/made', payload)
  return extractData<{ id: number }>(response.data)
}
