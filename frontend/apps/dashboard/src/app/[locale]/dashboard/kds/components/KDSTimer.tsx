'use client'

import { useEffect, useState } from 'react'

interface KDSTimerProps {
  elapsedSeconds: number
  warningMinutes?: number
  dangerMinutes?: number
  isStatic?: boolean
}

const formatElapsed = (seconds: number) => {
  const safeSeconds = Math.max(0, Math.floor(seconds))
  const minutes = Math.floor(safeSeconds / 60)
  const remainder = safeSeconds % 60
  return `${String(minutes).padStart(2, '0')}:${String(remainder).padStart(2, '0')}`
}

export function KDSTimer({
  elapsedSeconds,
  warningMinutes = 10,
  dangerMinutes = 20,
  isStatic = false,
}: KDSTimerProps) {
  const [seconds, setSeconds] = useState(Math.max(0, Math.floor(elapsedSeconds)))

  useEffect(() => {
    setSeconds(Math.max(0, Math.floor(elapsedSeconds)))
  }, [elapsedSeconds])

  useEffect(() => {
    if (isStatic) {
      return
    }

    const intervalId = window.setInterval(() => {
      setSeconds((current) => current + 1)
    }, 1000)

    return () => {
      window.clearInterval(intervalId)
    }
  }, [isStatic])

  const warningThreshold = warningMinutes * 60
  const dangerThreshold = dangerMinutes * 60

  const colorClass =
    seconds >= dangerThreshold
      ? 'text-red-700 dark:text-red-300'
      : seconds >= warningThreshold
        ? 'text-amber-700 dark:text-amber-300'
        : 'text-emerald-700 dark:text-emerald-300'

  return (
    <span className={`font-mono text-sm font-semibold tracking-wide ${colorClass}`}>
      {formatElapsed(seconds)}
    </span>
  )
}
