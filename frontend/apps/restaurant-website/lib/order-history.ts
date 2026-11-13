import axios from 'axios'
import { z } from 'zod'
import { apiClient } from './api-client'
import { rememberedOrderIds } from './order-access'
import type { OrderListResponse } from './hooks/use-api-queries'

const responseSchema = z.object({
  data: z.object({
    id: z.number().int().positive(),
    order_number: z.string(),
    status: z.enum(['pending', 'confirmed', 'preparing', 'ready', 'out_for_delivery', 'delivered', 'cancelled']),
    total_amount: z.number().finite().nonnegative(),
    created_at: z.string(),
  }),
})

// There is deliberately no unauthenticated order-list endpoint. Each saved order
// is fetched with its own capability by the existing API client's token interceptor.
export async function loadRememberedOrders(): Promise<OrderListResponse> {
  const entries = await Promise.all(rememberedOrderIds().map(async id => {
    try {
      const response = await apiClient.get<unknown>(`/public/orders/${id}`)
      const order = responseSchema.parse(response).data
      return { id: order.id, orderNumber: order.order_number, status: order.status, totalAmount: order.total_amount, createdAt: order.created_at }
    } catch (error) {
      // Expired or revoked capabilities reveal no order and are omitted. Service
      // failures remain errors so the UI offers retry instead of a false empty state.
      if (axios.isAxiosError(error) && error.response?.status === 404) return null
      throw error
    }
  }))
  return { success: true, data: entries.filter((entry): entry is NonNullable<typeof entry> => entry !== null) }
}
