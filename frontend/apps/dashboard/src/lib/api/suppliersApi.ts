import { api } from '@/lib/api'
import type { Supplier, SupplierFormInput } from '@/types/accounting'

const extractData = <T>(payload: any): T => (payload?.data ?? payload) as T

export async function getSuppliers(): Promise<Supplier[]> {
  const response = await api.get('/suppliers')
  return extractData<Supplier[]>(response.data)
}

export async function getSupplier(id: number): Promise<Supplier> {
  const response = await api.get(`/suppliers/${id}`)
  return extractData<Supplier>(response.data)
}

export async function createSupplier(
  payload: SupplierFormInput
): Promise<Supplier> {
  const response = await api.post('/suppliers', payload)
  return extractData<Supplier>(response.data)
}

export async function updateSupplier(
  id: number,
  payload: SupplierFormInput
): Promise<Supplier> {
  const response = await api.put(`/suppliers/${id}`, payload)
  return extractData<Supplier>(response.data)
}

export async function deleteSupplier(id: number): Promise<void> {
  await api.delete(`/suppliers/${id}`)
}
