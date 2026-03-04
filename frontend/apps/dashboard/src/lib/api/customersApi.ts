import { api } from '@/lib/api'
import type {
  CustomerAccount,
  CustomerAccountFormInput,
  CustomerStatementLine,
} from '@/types/accounting'

const extractData = <T>(payload: any): T => (payload?.data ?? payload) as T

export async function getCustomerAccounts(): Promise<CustomerAccount[]> {
  const response = await api.get('/accounting/customers')
  return extractData<CustomerAccount[]>(response.data)
}

export async function getCustomerAccount(id: number): Promise<CustomerAccount> {
  const response = await api.get(`/accounting/customers/${id}`)
  return extractData<CustomerAccount>(response.data)
}

export async function createCustomerAccount(
  payload: CustomerAccountFormInput
): Promise<CustomerAccount> {
  const response = await api.post('/accounting/customers', payload)
  return extractData<CustomerAccount>(response.data)
}

export async function updateCustomerAccount(
  id: number,
  payload: CustomerAccountFormInput
): Promise<CustomerAccount> {
  const response = await api.put(`/accounting/customers/${id}`, payload)
  return extractData<CustomerAccount>(response.data)
}

export async function deleteCustomerAccount(id: number): Promise<void> {
  await api.delete(`/accounting/customers/${id}`)
}

export async function getCustomerStatement(
  id: number
): Promise<CustomerStatementLine[]> {
  const response = await api.get(`/accounting/customers/${id}/statement`)
  return extractData<CustomerStatementLine[]>(response.data)
}
