import { act, renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockRecordPerformanceMetric = vi.fn()

vi.mock('@/lib/performance', () => ({
  recordPerformanceMetric: (...args: unknown[]) => mockRecordPerformanceMetric(...args),
}))

import { usePOSPerformance } from './usePOSPerformance'

describe('usePOSPerformance', () => {
  beforeEach(() => {
    mockRecordPerformanceMetric.mockReset()
  })

  it('measures sync callbacks and records metrics', () => {
    const { result } = renderHook(() => usePOSPerformance())

    let callbackResult = 0
    act(() => {
      callbackResult = result.current.measureSync('pos.add-to-cart', () => 42)
    })

    expect(callbackResult).toBe(42)
    expect(mockRecordPerformanceMetric).toHaveBeenCalled()
    expect(result.current.recentMetrics.length).toBeGreaterThan(0)
    expect(result.current.recentMetrics[0]?.metric).toBe('pos.add-to-cart')
  })

  it('marks page-load once', () => {
    const { result } = renderHook(() => usePOSPerformance())

    act(() => {
      result.current.markPageLoad()
      result.current.markPageLoad()
    })

    const pageLoadCalls = mockRecordPerformanceMetric.mock.calls.filter(
      (call) => call[0] === 'pos.page-load'
    )
    expect(pageLoadCalls).toHaveLength(1)
  })
})

