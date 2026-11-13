import axios from 'axios'
import { normalizePublicOrder } from './order-response'
import { apiClient } from './api-client'
import { rememberedOrderIds } from './order-access'
import type { OrderListResponse } from './hooks/use-api-queries'

// There is deliberately no unauthenticated order-list endpoint. Each saved order
// is fetched with its own capability by the existing API client's token interceptor.
export async function loadRememberedOrders(): Promise<OrderListResponse> {
  const entries = await Promise.all(rememberedOrderIds().map(async id => {
    try {
      const response = await apiClient.get<unknown>(`/public/orders/${id}`)
      return normalizePublicOrder(response)
    } catch (error) {
      // Expired or revoked capabilities reveal no order and are omitted. Service
      // failures remain errors so the UI offers retry instead of a false empty state.
      if (axios.isAxiosError(error) && error.response?.status === 404) return null
      throw error
    }
  }))
  return { success: true, data: entries.filter((entry): entry is NonNullable<typeof entry> => entry !== null) }
}
