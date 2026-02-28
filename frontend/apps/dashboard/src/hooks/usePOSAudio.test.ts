import { act, renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { usePOSAudio } from './usePOSAudio'
import { usePOSStore } from '@/stores/posStore'

type MockAudioInstance = {
  src: string
  preload: string
  volume: number
  muted: boolean
  currentTime: number
  play: ReturnType<typeof vi.fn>
  pause: ReturnType<typeof vi.fn>
}

const audioBySrc = new Map<string, MockAudioInstance>()

class MockAudioElement implements MockAudioInstance {
  src: string
  preload = ''
  volume = 1
  muted = false
  currentTime = 0
  play = vi.fn(async () => undefined)
  pause = vi.fn()

  constructor(src: string) {
    this.src = src
    audioBySrc.set(src, this)
  }
}

class MockAudioContext {
  state: AudioContextState = 'running'
  destination = {}
  resume = vi.fn(async () => undefined)
  close = vi.fn(async () => undefined)
  createGain = vi.fn(() => ({
    gain: {
      setValueAtTime: vi.fn(),
      exponentialRampToValueAtTime: vi.fn(),
    },
    connect: vi.fn(),
  }))
  createOscillator = vi.fn(() => ({
    type: 'sine' as OscillatorType,
    frequency: {
      setValueAtTime: vi.fn(),
      linearRampToValueAtTime: vi.fn(),
    },
    connect: vi.fn(),
    start: vi.fn(),
    stop: vi.fn(),
  }))
}

describe('usePOSAudio', () => {
  beforeEach(() => {
    audioBySrc.clear()
    usePOSStore.setState({ soundEnabled: true })

    Object.defineProperty(window, 'Audio', {
      value: MockAudioElement,
      configurable: true,
      writable: true,
    })

    Object.defineProperty(window, 'AudioContext', {
      value: MockAudioContext,
      configurable: true,
      writable: true,
    })
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('does not play sounds when sound toggle is off', async () => {
    usePOSStore.getState().setSoundEnabled(false)
    const { result } = renderHook(() => usePOSAudio())

    let played = false
    await act(async () => {
      played = await result.current.playSound('add')
    })

    expect(played).toBe(false)

    const addAudio = audioBySrc.get('/sounds/pos-add.mp3')
    expect(addAudio).toBeTruthy()
    expect(addAudio?.play).not.toHaveBeenCalled()
  })

  it('plays sound when enabled', async () => {
    const { result } = renderHook(() => usePOSAudio())

    let played = false
    await act(async () => {
      played = await result.current.playSound('add')
    })

    expect(played).toBe(true)

    const addAudio = audioBySrc.get('/sounds/pos-add.mp3')
    expect(addAudio).toBeTruthy()
    expect(addAudio?.play).toHaveBeenCalled()
  })

  it('respects runtime sound toggle changes', async () => {
    const { result } = renderHook(() => usePOSAudio())

    await act(async () => {
      await result.current.playSound('success')
    })

    await act(async () => {
      usePOSStore.getState().setSoundEnabled(false)
    })
    await act(async () => {
      await Promise.resolve()
    })

    let played = true
    await act(async () => {
      played = await result.current.playSound('remove')
    })

    expect(played).toBe(false)
  })
})
