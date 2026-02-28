import { act, cleanup, renderHook, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { useReports } from './useReports'

async function flushMicrotasks() {
  await Promise.resolve()
  await Promise.resolve()
}

describe('useReports', () => {
  beforeEach(() => {
    vi.useRealTimers()
  })

  afterEach(() => {
    cleanup()
    vi.clearAllTimers()
    vi.useRealTimers()
    vi.restoreAllMocks()
  })

  it('fetches data on mount and updates state', async () => {
    const fetcher = vi.fn().mockResolvedValue({ total: 42 })
    const getErrorMessage = vi.fn().mockReturnValue('failed')

    const { result } = renderHook(() =>
      useReports({
        enabled: true,
        deps: [],
        fetcher,
        getErrorMessage,
      })
    )

    expect(result.current.loading).toBe(true)

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
      expect(result.current.data).toEqual({ total: 42 })
    })

    expect(fetcher).toHaveBeenCalledTimes(1)
    expect(result.current.error).toBeNull()
    expect(result.current.lastUpdatedAt).toBeInstanceOf(Date)
  })

  it('handles fetch errors and calls onError with message', async () => {
    const err = new Error('network')
    const fetcher = vi.fn().mockRejectedValue(err)
    const getErrorMessage = vi.fn().mockReturnValue('Network Error')
    const onError = vi.fn()

    const { result } = renderHook(() =>
      useReports({
        enabled: true,
        deps: [],
        fetcher,
        getErrorMessage,
        onError,
      })
    )

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
      expect(result.current.error).toBe('Network Error')
    })

    expect(getErrorMessage).toHaveBeenCalledWith(err)
    expect(onError).toHaveBeenCalledWith('Network Error', err, false)
    expect(result.current.data).toBeNull()
  })

  it('does not refetch in a loop when callbacks change, and polling uses latest callback refs', async () => {
    vi.useFakeTimers()

    const fetcherA = vi.fn().mockResolvedValue({ version: 'A' })
    const fetcherB = vi.fn().mockResolvedValue({ version: 'B' })
    const onErrorA = vi.fn()
    const onErrorB = vi.fn()
    const getErrorMessageA = vi.fn().mockReturnValue('A')
    const getErrorMessageB = vi.fn().mockReturnValue('B')

    type Props = {
      fetcher: () => Promise<{ version: string }>
      getErrorMessage: (error: unknown) => string
      onError: (message: string, error: unknown, silent: boolean) => void
    }

    const { result, rerender, unmount } = renderHook(
      ({ fetcher, getErrorMessage, onError }: Props) =>
        useReports({
          enabled: true,
          deps: ['static-key'],
          fetcher,
          getErrorMessage,
          onError,
          pollIntervalMs: 1000,
        }),
      {
        initialProps: {
          fetcher: fetcherA,
          getErrorMessage: getErrorMessageA,
          onError: onErrorA,
        },
      }
    )

    await act(async () => {
      await flushMicrotasks()
    })
    expect(result.current.loading).toBe(false)
    expect(result.current.data).toEqual({ version: 'A' })
    expect(fetcherA).toHaveBeenCalledTimes(1)

    rerender({
      fetcher: fetcherB,
      getErrorMessage: getErrorMessageB,
      onError: onErrorB,
    })

    // Regression check: changing callback identities should not trigger immediate refetch.
    expect(fetcherA).toHaveBeenCalledTimes(1)
    expect(fetcherB).toHaveBeenCalledTimes(0)

    await act(async () => {
      await vi.advanceTimersByTimeAsync(1000)
      await flushMicrotasks()
    })

    expect(result.current.data).toEqual({ version: 'B' })
    expect(fetcherB).toHaveBeenCalledTimes(1)
    expect(result.current.refreshing).toBe(false)
    expect(onErrorA).not.toHaveBeenCalled()
    expect(onErrorB).not.toHaveBeenCalled()

    unmount()
  })

  it('does not fetch when disabled and marks loading false', async () => {
    const fetcher = vi.fn().mockResolvedValue({ total: 1 })
    const getErrorMessage = vi.fn().mockReturnValue('failed')

    const { result } = renderHook(() =>
      useReports({
        enabled: false,
        deps: [],
        fetcher,
        getErrorMessage,
      })
    )

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })

    expect(fetcher).not.toHaveBeenCalled()
    expect(result.current.data).toBeNull()
  })
})
