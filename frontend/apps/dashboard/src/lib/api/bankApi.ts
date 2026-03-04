import { api } from '@/lib/api'
import type {
  BankAccount,
  BankAccountFormInput,
  BankTransaction,
} from '@/types/invoices'

const extractData = <T>(payload: any): T => (payload?.data ?? payload) as T

export async function getBankAccounts(): Promise<BankAccount[]> {
  const response = await api.get('/accounting/bank-accounts')
  return extractData<BankAccount[]>(response.data)
}

export async function getBankAccount(id: number): Promise<BankAccount> {
  const response = await api.get(`/accounting/bank-accounts/${id}`)
  return extractData<BankAccount>(response.data)
}

export async function createBankAccount(
  payload: BankAccountFormInput
): Promise<BankAccount> {
  const response = await api.post('/accounting/bank-accounts', payload)
  return extractData<BankAccount>(response.data)
}

export async function updateBankAccount(
  id: number,
  payload: BankAccountFormInput
): Promise<BankAccount> {
  const response = await api.put(`/accounting/bank-accounts/${id}`, payload)
  return extractData<BankAccount>(response.data)
}

export async function deleteBankAccount(id: number): Promise<void> {
  await api.delete(`/accounting/bank-accounts/${id}`)
}

export async function getBankTransactions(
  id: number
): Promise<BankTransaction[]> {
  const response = await api.get(`/accounting/bank-accounts/${id}/transactions`)
  return extractData<BankTransaction[]>(response.data)
}
