import { afterEach, describe, expect, it, vi } from 'vitest'
import { act, renderHook } from '@testing-library/react'
import { useKDSWebSocket } from '@/hooks/useKDSWebSocket'
import { useAuthStore } from '@/stores/authStore'

const socketURLs: string[] = []

class MockWebSocket {
  static CONNECTING = 0
  static OPEN = 1
  static CLOSING = 2

  readyState = MockWebSocket.CONNECTING
  onopen: (() => void) | null = null
  onmessage: ((event: MessageEvent) => void) | null = null
  onerror: (() => void) | null = null
  onclose: (() => void) | null = null

  constructor(url: string) {
    socketURLs.push(url)
  }

  addEventListener() {}
  close() {}
}

vi.mock('@/lib/config', () => ({
  getConfig: () => ({ apiUrl: 'https://api.example.test/api/v1' }),
}))

describe('useKDSWebSocket', () => {
  afterEach(() => {
    useAuthStore.getState().logout()
    socketURLs.length = 0
    vi.unstubAllGlobals()
  })

  it('opens the socket for a cookie-backed authenticated session', () => {
    act(() => {
      useAuthStore.getState().setAuth({
        id: 1,
        email: 'kds@example.test',
        name: 'KDS User',
        role: 'owner',
        tenant_id: 10,
        restaurant_id: 20,
      }, '')
    })
    vi.stubGlobal('WebSocket', MockWebSocket as unknown as typeof WebSocket)

    renderHook(() => useKDSWebSocket())

    expect(socketURLs).toEqual(['wss://api.example.test/api/v1/kds/ws'])
  })

  it('does not open the socket without an authenticated session', () => {
    vi.stubGlobal('WebSocket', MockWebSocket as unknown as typeof WebSocket)

    renderHook(() => useKDSWebSocket())

    expect(socketURLs).toEqual([])
  })
})
