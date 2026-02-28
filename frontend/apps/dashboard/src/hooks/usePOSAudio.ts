'use client'

import { useEffect, useRef } from 'react'
import { debugError } from '@/lib/debug'
import { selectPOSSoundEnabled, usePOSStore } from '@/stores/posStore'

export type POSAudioType = 'add' | 'remove' | 'success' | 'error'

interface POSSoundConfig {
  path: string
  volume: number
  fallbackFrequencyHz: number
  fallbackDurationSeconds: number
  fallbackOscillator: OscillatorType
}

type BrowserAudioContext = typeof AudioContext

const SOUND_CONFIG: Record<POSAudioType, POSSoundConfig> = {
  add: {
    path: '/sounds/pos-add.mp3',
    volume: 0.75,
    fallbackFrequencyHz: 880,
    fallbackDurationSeconds: 0.11,
    fallbackOscillator: 'triangle',
  },
  remove: {
    path: '/sounds/pos-remove.mp3',
    volume: 0.72,
    fallbackFrequencyHz: 523.25,
    fallbackDurationSeconds: 0.12,
    fallbackOscillator: 'triangle',
  },
  success: {
    path: '/sounds/pos-success.mp3',
    volume: 0.8,
    fallbackFrequencyHz: 1174.66,
    fallbackDurationSeconds: 0.24,
    fallbackOscillator: 'sine',
  },
  error: {
    path: '/sounds/pos-error.mp3',
    volume: 0.8,
    fallbackFrequencyHz: 220,
    fallbackDurationSeconds: 0.22,
    fallbackOscillator: 'sawtooth',
  },
}

const getAudioContextConstructor = (): BrowserAudioContext | null => {
  if (typeof window === 'undefined') {
    return null
  }

  return (
    window.AudioContext ||
    (window as Window & { webkitAudioContext?: BrowserAudioContext }).webkitAudioContext ||
    null
  )
}

export const usePOSAudio = () => {
  const soundEnabled = usePOSStore(selectPOSSoundEnabled)
  const soundEnabledRef = useRef(soundEnabled)
  const audioMapRef = useRef<Map<POSAudioType, HTMLAudioElement>>(new Map())
  const audioContextRef = useRef<AudioContext | null>(null)
  const isUnlockedRef = useRef(false)

  useEffect(() => {
    soundEnabledRef.current = soundEnabled
  }, [soundEnabled])

  const ensureAudioContext = async () => {
    if (audioContextRef.current) {
      if (audioContextRef.current.state === 'suspended') {
        try {
          await audioContextRef.current.resume()
        } catch {
          return null
        }
      }
      return audioContextRef.current
    }

    const Ctx = getAudioContextConstructor()
    if (!Ctx) {
      return null
    }

    try {
      const context = new Ctx()
      audioContextRef.current = context
      if (context.state === 'suspended') {
        await context.resume()
      }
      return context
    } catch {
      return null
    }
  }

  const unlockAudio = async () => {
    if (isUnlockedRef.current) {
      return true
    }

    await ensureAudioContext()

    for (const audio of Array.from(audioMapRef.current.values())) {
      const previousMuted = audio.muted
      const previousVolume = audio.volume

      try {
        audio.currentTime = 0
        audio.muted = true
        audio.volume = 0
        await audio.play()
        audio.pause()
        audio.currentTime = 0
      } catch {
        // Keep trying on future interactions.
      } finally {
        audio.muted = previousMuted
        audio.volume = previousVolume
      }
    }

    isUnlockedRef.current = true
    return true
  }

  const playFallbackTone = async (type: POSAudioType) => {
    const context = await ensureAudioContext()
    if (!context) {
      return false
    }

    const config = SOUND_CONFIG[type]
    const now = context.currentTime
    const gain = context.createGain()
    const oscillator = context.createOscillator()

    oscillator.type = config.fallbackOscillator
    oscillator.frequency.setValueAtTime(config.fallbackFrequencyHz, now)
    if (type === 'success') {
      oscillator.frequency.linearRampToValueAtTime(
        config.fallbackFrequencyHz * 1.22,
        now + config.fallbackDurationSeconds
      )
    }

    gain.gain.setValueAtTime(0.0001, now)
    gain.gain.exponentialRampToValueAtTime(0.2, now + 0.015)
    gain.gain.exponentialRampToValueAtTime(
      0.0001,
      now + config.fallbackDurationSeconds
    )

    oscillator.connect(gain)
    gain.connect(context.destination)

    oscillator.start(now)
    oscillator.stop(now + config.fallbackDurationSeconds)
    return true
  }

  const playSound = async (type: POSAudioType) => {
    if (!soundEnabledRef.current) {
      return false
    }

    if (!isUnlockedRef.current) {
      await unlockAudio()
    }

    if (!isUnlockedRef.current) {
      return false
    }

    const audio = audioMapRef.current.get(type)
    if (audio) {
      try {
        audio.currentTime = 0
        audio.volume = SOUND_CONFIG[type].volume
        await audio.play()
        return true
      } catch (error) {
        debugError('[POS Audio] Audio file playback failed, using fallback tone', error)
      }
    }

    return playFallbackTone(type)
  }

  useEffect(() => {
    if (typeof window === 'undefined') {
      return
    }

    const map = new Map<POSAudioType, HTMLAudioElement>()
    ;(Object.keys(SOUND_CONFIG) as POSAudioType[]).forEach((type) => {
      const config = SOUND_CONFIG[type]
      const audio = new Audio(config.path)
      audio.preload = 'auto'
      audio.volume = config.volume
      map.set(type, audio)
    })
    audioMapRef.current = map

    const handleFirstInteraction = () => {
      void unlockAudio()
      window.removeEventListener('pointerdown', handleFirstInteraction)
      window.removeEventListener('keydown', handleFirstInteraction)
      window.removeEventListener('touchstart', handleFirstInteraction)
    }

    window.addEventListener('pointerdown', handleFirstInteraction, { passive: true })
    window.addEventListener('keydown', handleFirstInteraction)
    window.addEventListener('touchstart', handleFirstInteraction, { passive: true })

    return () => {
      window.removeEventListener('pointerdown', handleFirstInteraction)
      window.removeEventListener('keydown', handleFirstInteraction)
      window.removeEventListener('touchstart', handleFirstInteraction)

      audioMapRef.current.forEach((audio) => {
        audio.pause()
        audio.src = ''
      })
      audioMapRef.current.clear()

      if (audioContextRef.current) {
        void audioContextRef.current.close()
        audioContextRef.current = null
      }
    }
    // Keep this effect mount-only to avoid rebuilding audio objects.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return {
    playSound,
    unlockAudio,
  }
}

export default usePOSAudio
