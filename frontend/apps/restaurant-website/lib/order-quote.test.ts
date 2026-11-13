import { apiClient } from './api-client'
import { fetchOrderQuote } from './order-quote'
import { normalizePublicOrder } from './order-response'

jest.mock('./api-client', () => ({ apiClient: { post: jest.fn() } }))
const post = jest.mocked(apiClient.post)
beforeEach(() => jest.clearAllMocks())

describe('authoritative checkout pricing', () => {
  it('requests prices by product and quantity and uses server tax and delivery', async () => {
    post.mockResolvedValue({ success: true, data: { subtotal: 25, tax_amount: 2.5, delivery_fee: 0, total_amount: 27.5, currency: 'EGP', items: [{ product_id: 7, quantity: 1, unit_price: 25, total_price: 25 }] } })
    const quote = await fetchOrderQuote([{ productId: 7, quantity: 1 }])
    expect(post).toHaveBeenCalledWith('/public/orders/quote', { items: [{ product_id: 7, quantity: 1 }] })
    expect(quote.total_amount).toBe(27.5)
    expect(quote.delivery_fee).toBe(0)
  })

  it('rejects inconsistent or malformed totals instead of showing a fabricated estimate', async () => {
    post.mockResolvedValue({ success: true, data: { subtotal: 25, tax_amount: 2.5, delivery_fee: 0, total_amount: 30.99, currency: 'EGP', items: [] } })
    await expect(fetchOrderQuote([{ productId: 7, quantity: 1 }])).rejects.toThrow('Inconsistent order quote')
  })

  it('normalizes the actual protected detail response and line snapshots', () => {
    const order = normalizePublicOrder({ success: true, data: { id: 1, order_number: 'QA-1', status: 'pending', total_amount: 27.5, created_at: '2026-11-13T00:00:00Z', items: [{ id: 2, product_id: 7, product_name: 'Synthetic', quantity: 1, unit_price: 25 }] } })
    expect(order.orderNumber).toBe('QA-1')
    expect(order.totalAmount).toBe(27.5)
    expect(order.items?.[0]).toMatchObject({ name: 'Synthetic', quantity: 1, price: 25 })
  })
})
