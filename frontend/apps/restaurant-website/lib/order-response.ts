import { z } from 'zod'
import type { Order } from './hooks/use-api-queries'

const orderSchema = z.object({
  id: z.number().int().positive(), order_number: z.string(),
  status: z.enum(['pending', 'confirmed', 'preparing', 'ready', 'out_for_delivery', 'delivered', 'cancelled']),
  total_amount: z.number().finite().nonnegative(), created_at: z.string(),
  customer_name: z.string().optional(), customer_email: z.string().optional(), customer_phone: z.string().optional(), delivery_address: z.string().optional(),
  items: z.array(z.object({ id: z.number(), product_id: z.number(), product_name: z.string(), quantity: z.number().int().nonnegative(), unit_price: z.number().finite().nonnegative(), special_instructions: z.string().optional() })).optional(),
})

export function normalizePublicOrder(raw: unknown): Order {
  const data = z.object({ data: orderSchema }).parse(raw).data
  return { id: data.id, orderNumber: data.order_number, status: data.status, totalAmount: data.total_amount, createdAt: data.created_at,
    customerName: data.customer_name, customerEmail: data.customer_email, customerPhone: data.customer_phone, deliveryAddress: data.delivery_address,
    items: data.items?.map(item => ({ id: item.id, productId: item.product_id, name: item.product_name, quantity: item.quantity, price: item.unit_price, notes: item.special_instructions })) }
}
