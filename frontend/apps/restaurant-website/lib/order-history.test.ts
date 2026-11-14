import { apiClient } from './api-client'
import { rememberOrderToken } from './order-access'
import { loadRememberedOrders } from './order-history'

jest.mock('./api-client', () => ({ apiClient: { get: jest.fn() } }))
const get = jest.mocked(apiClient.get)
beforeEach(() => { localStorage.clear(); jest.clearAllMocks() })
const response = { data: { id: 7, order_number: 'SYNTHETIC-7', status: 'pending', total_amount: 25, created_at: '2026-11-13T00:00:00Z' } }

describe('device-held public order history', () => {
  it('makes no request when the device has no remembered capability', async () => {
    expect(await loadRememberedOrders()).toEqual({ success: true, data: [] })
    expect(get).not.toHaveBeenCalled()
  })

  it('loads only a remembered order and maps the actual API response', async () => {
    rememberOrderToken({ id: 7, order_number: 'SYNTHETIC-7', tracking_token: 'synthetic-proof' })
    get.mockResolvedValue(response)
    const result = await loadRememberedOrders()
    expect(get).toHaveBeenCalledWith('/public/orders/7')
    expect(get).toHaveBeenCalledTimes(1)
    expect(result.data[0]).toMatchObject({ id: 7, orderNumber: 'SYNTHETIC-7', totalAmount: 25 })
  })

  it('omits revoked capabilities but surfaces server failures', async () => {
    rememberOrderToken({ id: 7, tracking_token: 'synthetic-proof' })
    get.mockRejectedValue({ isAxiosError: true, response: { status: 404 } })
    expect((await loadRememberedOrders()).data).toEqual([])
    get.mockRejectedValue(new Error('service unavailable'))
    await expect(loadRememberedOrders()).rejects.toThrow('service unavailable')
  })

  it('rejects malformed server data rather than displaying an invalid order', async () => {
    rememberOrderToken({ id: 7, tracking_token: 'synthetic-proof' })
    get.mockResolvedValue({ data: { id: 7 } })
    await expect(loadRememberedOrders()).rejects.toThrow()
  })
})
