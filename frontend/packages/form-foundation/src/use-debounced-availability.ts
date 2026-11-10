'use client'

import { useEffect, useRef, useState } from 'react'

export type AvailabilityState = 'idle' | 'invalid' | 'checking' | 'available' | 'taken' | 'reserved' | 'error'
export interface AvailabilityResult { state: AvailabilityState; checkedValue: string }
export type AvailabilityCheckResult = boolean | 'reserved'

export function useDebouncedAvailability(
  value: string,
  isValidFormat: (value: string) => boolean,
  check: (value: string, signal: AbortSignal) => Promise<AvailabilityCheckResult>,
  delay = 500,
): AvailabilityResult {
  const [result, setResult] = useState<AvailabilityResult>({ state: 'idle', checkedValue: '' })
  const sequence = useRef(0)

  useEffect(() => {
    const candidate = value.trim().toLowerCase()
    const currentSequence = ++sequence.current
    if (!candidate) {
      setResult({ state: 'idle', checkedValue: '' })
      return
    }
    if (!isValidFormat(candidate)) {
      setResult({ state: 'invalid', checkedValue: candidate })
      return
    }

    const controller = new AbortController()
    setResult({ state: 'checking', checkedValue: candidate })
    const timeout = setTimeout(() => {
      void check(candidate, controller.signal).then((available) => {
        if (sequence.current === currentSequence && !controller.signal.aborted) {
          const state = available === 'reserved' ? 'reserved' : available ? 'available' : 'taken'
          setResult({ state, checkedValue: candidate })
        }
      }).catch(() => {
        if (sequence.current === currentSequence && !controller.signal.aborted) {
          setResult({ state: 'error', checkedValue: candidate })
        }
      })
    }, delay)

    return () => {
      clearTimeout(timeout)
      controller.abort()
    }
  }, [value, isValidFormat, check, delay])

  return result
}
