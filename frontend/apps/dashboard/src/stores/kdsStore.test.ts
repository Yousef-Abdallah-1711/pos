import { beforeEach, describe, expect, it, vi } from 'vitest'

const KDS_STORAGE_NAME = 'kds-storage'

const waitForHydration = async () => {
  await Promise.resolve()
  await Promise.resolve()
}

describe('kdsStore persistence', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.resetModules()
  })

  it('clears malformed persisted payloads without throwing', async () => {
    localStorage.setItem(KDS_STORAGE_NAME, '{"state":')

    const { useKDSStore } = await import('./kdsStore')
    await waitForHydration()

    const state = useKDSStore.getState()
    expect(state.filter).toBe('all')
    expect(state.soundEnabled).toBe(true)
    expect(state.isFullscreen).toBe(false)
    expect(localStorage.getItem(KDS_STORAGE_NAME)).toBeNull()
  })

  it('hydrates valid persisted UI preferences', async () => {
    localStorage.setItem(
      KDS_STORAGE_NAME,
      JSON.stringify({
        state: {
          filter: 'ready',
          soundEnabled: false,
          isFullscreen: true,
        },
        version: 0,
      })
    )

    const { useKDSStore } = await import('./kdsStore')
    await waitForHydration()

    const state = useKDSStore.getState()
    expect(state.filter).toBe('ready')
    expect(state.soundEnabled).toBe(false)
    expect(state.isFullscreen).toBe(true)
  })
})
