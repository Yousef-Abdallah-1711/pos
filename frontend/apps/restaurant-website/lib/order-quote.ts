import { z } from 'zod'
import { apiClient } from './api-client'

export const orderQuoteSchema = z.object({
  subtotal: z.number().finite().nonnegative(),
  tax_amount: z.number().finite().nonnegative(),
  delivery_fee: z.number().finite().nonnegative(),
  total_amount: z.number().finite().nonnegative(),
  currency: z.literal('EGP'),
  items: z.array(z.object({ product_id: z.number().int().positive(), quantity: z.number().int().positive(), unit_price: z.number().finite().nonnegative(), total_price: z.number().finite().nonnegative() })),
}).refine(quote => Math.abs(quote.subtotal + quote.tax_amount + quote.delivery_fee - quote.total_amount) < 0.005, 'Inconsistent order quote')
export type OrderQuote = z.infer<typeof orderQuoteSchema>
export interface QuoteCartItem { productId: number; quantity: number }

export async function fetchOrderQuote(items: readonly QuoteCartItem[]): Promise<OrderQuote> {
  const raw = await apiClient.post<unknown>('/public/orders/quote', { items: items.map(item => ({ product_id: item.productId, quantity: item.quantity })) })
  return z.object({ success: z.literal(true), data: orderQuoteSchema }).parse(raw).data
}
