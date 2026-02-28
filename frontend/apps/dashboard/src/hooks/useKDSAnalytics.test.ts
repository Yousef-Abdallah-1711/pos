import { renderHook, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
  },
}))

import { api } from '@/lib/api'
import { useKDSAnalytics } from './useKDSAnalytics'

const mockedApi = api as unknown as { get: ReturnType<typeof vi.fn> }

describe('useKDSAnalytics', () => {
  beforeEach(() => {
    mockedApi.get.mockReset()
  })

  it('fetches analytics when enabled', async () => {
    mockedApi.get.mockResolvedValue({
      data: {
        success: true,
        data: {
          average_prep_minutes: 9.5,
          total_completed_orders: 25,
          prep_time_trend: [{ date: '2026-02-24', avg_prep_minutes: 9.5, completed_orders: 4 }],
          busiest_hours: [{ hour: 12, order_count: 8 }],
        },
      },
    })

    const { result } = renderHook(() => useKDSAnalytics(true))

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    expect(mockedApi.get).toHaveBeenCalledWith('/kds/analytics')
    expect(result.current.analytics.average_prep_minutes).toBe(9.5)
    expect(result.current.analytics.total_completed_orders).toBe(25)
    expect(result.current.analyticsAvailable).toBe(true)
    expect(result.current.error).toBeNull()
  })

  it('does not fetch when disabled', async () => {
    const { result } = renderHook(() => useKDSAnalytics(false))

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    expect(mockedApi.get).not.toHaveBeenCalled()
    expect(result.current.analytics.total_completed_orders).toBe(0)
    expect(result.current.analyticsAvailable).toBe(true)
  })

  it('surfaces API error message when request fails', async () => {
    mockedApi.get.mockRejectedValue({
      response: { status: 403, data: { error: 'Access denied' } },
      message: 'request failed',
    })

    const { result } = renderHook(() => useKDSAnalytics(true))

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
      expect(result.current.error).toBe('Access denied')
    })
  })

  it('marks analytics unavailable when endpoint is missing', async () => {
    mockedApi.get.mockRejectedValue({
      response: { status: 404, data: {} },
      message: 'not found',
    })

    const { result } = renderHook(() => useKDSAnalytics(true))

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    expect(result.current.analyticsAvailable).toBe(false)
    expect(result.current.error).toBeNull()
    expect(result.current.analytics.total_completed_orders).toBe(0)
  })
})
