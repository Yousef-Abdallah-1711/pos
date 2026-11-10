import { act, renderHook, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { useDebouncedAvailability } from '@pos-saas/form-foundation'

const valid = (value: string) => /^[a-z0-9-]{2,63}$/.test(value) && !value.includes('--')

describe('shared async availability validation', () => {
  afterEach(() => vi.useRealTimers())

  it('debounces input and marks malformed slugs without a request', async () => {
    vi.useFakeTimers()
    const check = vi.fn(async () => true)
    const { result, rerender } = renderHook(({ value }) => useDebouncedAvailability(value, valid, check, 500), { initialProps: { value: 'a' } })
    expect(result.current.state).toBe('invalid')
    rerender({ value: 'good' })
    act(() => { vi.advanceTimersByTime(499) })
    expect(check).not.toHaveBeenCalled()
    await act(async () => { vi.advanceTimersByTime(1); await Promise.resolve() })
    expect(check).toHaveBeenCalledTimes(1)
    expect(result.current.state).toBe('available')
  })

  it('aborts a stale request and ignores its late response', async () => {
    const requests: Array<{ value: string; signal: AbortSignal; resolve: (value: boolean) => void }> = []
    const check = vi.fn((value: string, signal: AbortSignal) => new Promise<boolean>((resolve) => requests.push({ value, signal, resolve })))
    const { result, rerender } = renderHook(({ value }) => useDebouncedAvailability(value, valid, check, 0), { initialProps: { value: 'first' } })
    await waitFor(() => expect(requests).toHaveLength(1))
    rerender({ value: 'second' })
    await waitFor(() => expect(requests).toHaveLength(2))
    expect(requests[0].signal.aborted).toBe(true)
    await act(async () => { requests[1].resolve(true); await Promise.resolve() })
    await act(async () => { requests[0].resolve(false); await Promise.resolve() })
    expect(result.current).toEqual({ state: 'available', checkedValue: 'second' })
  })

  it('distinguishes unavailable and network failure states', async () => {
    const firstCheck: (value: string, signal: AbortSignal) => Promise<boolean | 'reserved'> = async () => false
    const { result, rerender } = renderHook(({ value, check }: { value: string; check: (value: string, signal: AbortSignal) => Promise<boolean | 'reserved'> }) => useDebouncedAvailability(value, valid, check, 0), {
      initialProps: { value: 'taken', check: firstCheck },
    })
    await waitFor(() => expect(result.current.state).toBe('taken'))
    rerender({ value: 'reserved', check: async () => 'reserved' })
    await waitFor(() => expect(result.current.state).toBe('reserved'))
    rerender({ value: 'offline', check: async () => { throw new Error('offline') } })
    await waitFor(() => expect(result.current.state).toBe('error'))
  })
})
