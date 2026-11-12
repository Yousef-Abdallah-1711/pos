import { api } from '@/lib/api'
import { debugError } from '@/lib/debug'

export interface POSProduct {
  id: number
  name: string
  name_ar?: string
  description?: string
  price: number
  image?: string
  category_id?: number
  category_name?: string
  is_available?: boolean
  track_inventory?: boolean
  quantity_in_stock?: number
  barcode?: string
}

export interface POSCategory {
  id: number
  name: string
  name_ar?: string
  description?: string
  image?: string
  product_count?: number
}

export type POSPaymentMethod = 'cash' | 'card' | 'online' | 'wallet'
export type POSPaymentStatus = 'pending' | 'paid' | 'failed' | 'refunded'

export interface CreatePOSOrderItem {
  product_id: number
  quantity: number
  special_instructions?: string
}

export interface CreatePOSOrderPayload {
  order_source?: 'pos'
  customer_name: string
  customer_phone?: string
  payment_method: POSPaymentMethod
  table_id?: number
  notes?: string
  items: CreatePOSOrderItem[]
}

export interface POSOrderResponse {
  id: number
  order_number: string
  status: string
  total_amount: number
  customer_name?: string
  table_id?: number | null
  created_at?: string
}

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null

const toFiniteNumber = (value: unknown, fallback = 0): number => {
  const numeric = typeof value === 'number' ? value : Number(value)
  return Number.isFinite(numeric) ? numeric : fallback
}

const toText = (value: unknown, fallback = ''): string =>
  typeof value === 'string' ? value : fallback

const toOptionalText = (value: unknown): string | undefined => {
  if (typeof value !== 'string') {
    return undefined
  }
  const trimmed = value.trim()
  return trimmed.length > 0 ? trimmed : undefined
}

const toNullableInt = (value: unknown): number | null => {
  const numeric = toFiniteNumber(value, -1)
  if (numeric <= 0) {
    return null
  }
  return Math.floor(numeric)
}

const normalizeProduct = (value: unknown): POSProduct | null => {
  if (!isRecord(value)) {
    return null
  }

  const id = toNullableInt(value.id)
  const name = toOptionalText(value.name) || toOptionalText(value.name_en)
  const price = toFiniteNumber(value.price, -1)
  if (!id || !name || price < 0) {
    return null
  }

  return {
    id,
    name,
    name_ar: toOptionalText(value.name_ar),
    description: toOptionalText(value.description),
    price,
    image: toOptionalText(value.image) || toOptionalText(value.main_image_url),
    category_id: toNullableInt(value.category_id) ?? undefined,
    category_name: toOptionalText(value.category_name) || toOptionalText(value.category),
    is_available: typeof value.is_available === 'boolean' ? value.is_available : undefined,
    track_inventory:
      typeof value.track_inventory === 'boolean' ? value.track_inventory : undefined,
    quantity_in_stock:
      value.quantity_in_stock !== undefined
        ? Math.max(0, Math.floor(toFiniteNumber(value.quantity_in_stock, 0)))
        : undefined,
    barcode:
      typeof value.barcode === 'number'
        ? String(value.barcode)
        : toOptionalText(value.barcode),
  }
}

const normalizeCategory = (value: unknown): POSCategory | null => {
  if (!isRecord(value)) {
    return null
  }

  const id = toNullableInt(value.id)
  const name = toOptionalText(value.name)
  if (!id || !name) {
    return null
  }

  return {
    id,
    name,
    name_ar: toOptionalText(value.name_ar),
    description: toOptionalText(value.description),
    image: toOptionalText(value.image),
    product_count:
      value.product_count !== undefined
        ? Math.max(0, Math.floor(toFiniteNumber(value.product_count, 0)))
        : undefined,
  }
}

const extractListPayload = (payload: unknown, key: string): unknown[] => {
  // Required compatibility order from spec:
  // response.data.<key> || response.data.data || response.data
  if (!isRecord(payload)) {
    return Array.isArray(payload) ? payload : []
  }

  if (Array.isArray(payload[key])) {
    return payload[key] as unknown[]
  }

  if (Array.isArray(payload.data)) {
    return payload.data
  }

  if (isRecord(payload.data) && Array.isArray(payload.data[key])) {
    return payload.data[key] as unknown[]
  }

  return []
}

const extractOrderResponse = (payload: unknown): POSOrderResponse => {
  const data = isRecord(payload) && isRecord(payload.data) ? payload.data : payload
  if (!isRecord(data)) {
    throw new Error('Invalid POS create order response')
  }

  const id = toNullableInt(data.id)
  const orderNumber = toOptionalText(data.order_number)
  if (!id || !orderNumber) {
    throw new Error('POS create order response is missing required fields')
  }

  return {
    id,
    order_number: orderNumber,
    status: toText(data.status, 'pending'),
    total_amount: toFiniteNumber(data.total_amount, toFiniteNumber(data.total, 0)),
    customer_name: toOptionalText(data.customer_name),
    table_id: toNullableInt(data.table_id),
    created_at: toOptionalText(data.created_at),
  }
}

const sanitizeCreateOrderPayload = (
  payload: CreatePOSOrderPayload
): CreatePOSOrderPayload => {
  const items = Array.isArray(payload.items)
    ? payload.items
        .map((item) => ({
          product_id: Math.floor(toFiniteNumber(item.product_id, 0)),
          quantity: Math.max(1, Math.floor(toFiniteNumber(item.quantity, 1))),
          special_instructions: toOptionalText(item.special_instructions),
        }))
        .filter((item) => item.product_id > 0 && item.quantity > 0)
    : []

  if (items.length === 0) {
    throw new Error('Cannot create POS order without items')
  }

  const tableId = toNullableInt(payload.table_id)

  return {
    order_source: 'pos',
    customer_name: toOptionalText(payload.customer_name) || 'Walk-in',
    customer_phone: toOptionalText(payload.customer_phone),
    payment_method: payload.payment_method,
    table_id: tableId ?? undefined,
    notes: toOptionalText(payload.notes),
    items,
  }
}

export const posApi = {
  async getProducts(): Promise<POSProduct[]> {
    try {
      const response = await api.get('/products', {
        params: { page_size: 500 },
      })
      return extractListPayload(response.data, 'products')
        .map((item) => normalizeProduct(item))
        .filter((item): item is POSProduct => item !== null)
    } catch (error) {
      debugError('[posApi] Failed to fetch POS products:', error)
      throw error
    }
  },

  async getCategories(): Promise<POSCategory[]> {
    try {
      const response = await api.get('/categories')
      return extractListPayload(response.data, 'categories')
        .map((item) => normalizeCategory(item))
        .filter((item): item is POSCategory => item !== null)
    } catch (error) {
      debugError('[posApi] Failed to fetch POS categories:', error)
      throw error
    }
  },

  async createOrder(payload: CreatePOSOrderPayload): Promise<POSOrderResponse> {
    try {
      const sanitizedPayload = sanitizeCreateOrderPayload(payload)
      const response = await api.post('/pos/orders', sanitizedPayload)
      return extractOrderResponse(response.data)
    } catch (error) {
      debugError('[posApi] Failed to create POS order:', error)
      throw error
    }
  },

  async updateOrderPaymentStatus(
    orderId: number,
    status: POSPaymentStatus
  ): Promise<void> {
    try {
      const numericOrderId = Math.floor(toFiniteNumber(orderId, 0))
      if (numericOrderId <= 0) {
        throw new Error('Invalid order id for payment status update')
      }

      await api.put(`/admin/orders/${numericOrderId}/payment-status`, {
        payment_status: status,
      })
    } catch (error) {
      debugError('[posApi] Failed to update POS order payment status:', error)
      throw error
    }
  },
}

export default posApi
