'use client'

import { useEffect, useRef } from 'react'
import { debugError } from '@/lib/debug'
import { useKDSStore } from '@/stores/kdsStore'

const ALERT_SOUND_PATH = '/sounds/new-order.mp3'
const ALERT_VOLUME = 0.85
const FALLBACK_TONE_FREQUENCY_HZ = 1046.5
const FALLBACK_TONE_DURATION_SECONDS = 0.24

type BrowserAudioContext = typeof AudioContext

const getAudioContextConstructor = (): BrowserAudioContext | null => {
  if (typeof window === 'undefined') {
    return null
  }

  return (window.AudioContext || (window as Window & { webkitAudioContext?: BrowserAudioContext }).webkitAudioContext || null)
}

export const useKDSAudio = () => {
  const soundEnabled = useKDSStore((state) => state.soundEnabled)
  const soundEnabledRef = useRef(soundEnabled)
  const audioRef = useRef<HTMLAudioElement | null>(null)
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

    // Try to unlock HTMLAudio playback on first user gesture.
    if (audioRef.current) {
      const audio = audioRef.current
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
        // Keep trying via future user interactions.
      } finally {
        audio.muted = previousMuted
        audio.volume = previousVolume
      }
    }

    isUnlockedRef.current = true
    return true
  }

  const playFallbackTone = async () => {
    const context = await ensureAudioContext()
    if (!context) {
      return false
    }

    const now = context.currentTime
    const gain = context.createGain()
    const oscillator = context.createOscillator()

    oscillator.type = 'triangle'
    oscillator.frequency.setValueAtTime(FALLBACK_TONE_FREQUENCY_HZ, now)

    gain.gain.setValueAtTime(0.0001, now)
    gain.gain.exponentialRampToValueAtTime(0.2, now + 0.02)
    gain.gain.exponentialRampToValueAtTime(0.0001, now + FALLBACK_TONE_DURATION_SECONDS)

    oscillator.connect(gain)
    gain.connect(context.destination)

    oscillator.start(now)
    oscillator.stop(now + FALLBACK_TONE_DURATION_SECONDS)
    return true
  }

  const playNewOrderAlert = async () => {
    if (!soundEnabledRef.current) {
      return false
    }

    if (!isUnlockedRef.current) {
      await unlockAudio()
    }

    if (!isUnlockedRef.current) {
      return false
    }

    if (audioRef.current) {
      try {
        const audio = audioRef.current
        audio.volume = ALERT_VOLUME
        audio.currentTime = 0
        await audio.play()
        return true
      } catch (error) {
        debugError('[KDS Audio] Audio file playback failed, using fallback tone', error)
      }
    }

    return playFallbackTone()
  }

  useEffect(() => {
    if (typeof window === 'undefined') {
      return
    }

    const audio = new Audio(ALERT_SOUND_PATH)
    audio.preload = 'auto'
    audio.volume = ALERT_VOLUME
    audioRef.current = audio

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

      if (audioRef.current) {
        audioRef.current.pause()
        audioRef.current.src = ''
      }
      audioRef.current = null

      if (audioContextRef.current) {
        void audioContextRef.current.close()
        audioContextRef.current = null
      }
    }
    // Keep this effect mount-only to avoid resetting audio objects.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return {
    playNewOrderAlert,
    unlockAudio,
  }
}
