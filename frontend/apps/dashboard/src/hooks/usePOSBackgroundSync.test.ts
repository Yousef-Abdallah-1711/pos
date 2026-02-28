import { act, renderHook, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { POS_ORDER_SYNC_TAG, usePOSBackgroundSync } from './usePOSBackgroundSync'

type MockServiceWorkerContainer = {
  ready: Promise<ServiceWorkerRegistration>
  addEventListener: ReturnType<typeof vi.fn>
  removeEventListener: ReturnType<typeof vi.fn>
}

const setServiceWorkerContainer = (container: MockServiceWorkerContainer | undefined) => {
  Object.defineProperty(navigator, 'serviceWorker', {
    value: container,
    configurable: true,
    writable: true,
  })
}

describe('usePOSBackgroundSync', () => {
  afterEach(() => {
    setServiceWorkerContainer(undefined)
  })

  it('registers background sync when SyncManager is supported', async () => {
    const registerSync = vi.fn(async () => undefined)
    const registration = {
      sync: {
        register: registerSync,
      },
    } as unknown as ServiceWorkerRegistration

    setServiceWorkerContainer({
      ready: Promise.resolve(registration),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    })

    const { result } = renderHook(() => usePOSBackgroundSync())

    await waitFor(() => {
      expect(result.current.syncSupported).toBe(true)
    })

    let registered = false
    await act(async () => {
      registered = await result.current.registerOrderSync()
    })

    expect(registered).toBe(true)
    expect(registerSync).toHaveBeenCalledWith(POS_ORDER_SYNC_TAG)
  })

  it('forwards sync message events to callback', async () => {
    const registerSync = vi.fn(async () => undefined)
    const addEventListener = vi.fn()
    const removeEventListener = vi.fn()

    setServiceWorkerContainer({
      ready: Promise.resolve({
        sync: {
          register: registerSync,
        },
      } as unknown as ServiceWorkerRegistration),
      addEventListener,
      removeEventListener,
    })

    const onSyncMessage = vi.fn()
    const { result } = renderHook(() => usePOSBackgroundSync({ onSyncMessage }))

    await waitFor(() => {
      expect(result.current.syncSupported).toBe(true)
    })

    const handler = addEventListener.mock.calls.find(
      (call) => call[0] === 'message'
    )?.[1] as ((event: MessageEvent) => void) | undefined

    expect(typeof handler).toBe('function')
    if (!handler) {
      return
    }

    act(() => {
      handler(
        new MessageEvent('message', {
          data: {
            type: 'pos-sync-complete',
            processed: 2,
            failed: 0,
          },
        })
      )
    })

    expect(onSyncMessage).toHaveBeenCalledTimes(1)
    expect(onSyncMessage).toHaveBeenCalledWith({
      type: 'pos-sync-complete',
      processed: 2,
      failed: 0,
    })
  })
})
