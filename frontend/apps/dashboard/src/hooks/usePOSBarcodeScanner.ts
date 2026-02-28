'use client'

import { useEffect, useRef } from 'react'
import type { POSProduct } from '@/lib/api/posApi'

export interface POSBarcodeScannerOptions {
  getByBarcode: (barcode: string) => POSProduct | undefined
  onProductFound: (product: POSProduct, barcode: string) => void
  onBarcodeNotFound: (barcode: string) => void
  enabled?: boolean
  keyGapMs?: number
  maxElapsedMs?: number
  minBarcodeLength?: number
  maxBufferLength?: number
}

const DEFAULT_KEY_GAP_MS = 50
const DEFAULT_MAX_ELAPSED_MS = 500
const DEFAULT_MIN_BARCODE_LENGTH = 6
const DEFAULT_MAX_BUFFER_LENGTH = 50

const isEditableTarget = (target: EventTarget | null) => {
  if (!(target instanceof HTMLElement)) {
    return false
  }

  const tag = target.tagName.toUpperCase()
  if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') {
    return true
  }

  return target.isContentEditable
}

const isScannableCharacter = (event: KeyboardEvent) =>
  !event.ctrlKey &&
  !event.altKey &&
  !event.metaKey &&
  event.key.length === 1 &&
  !event.repeat

export const usePOSBarcodeScanner = (options: POSBarcodeScannerOptions) => {
  const optionsRef = useRef(options)

  useEffect(() => {
    optionsRef.current = options
  }, [options])

  useEffect(() => {
    if (typeof window === 'undefined') {
      return
    }

    let buffer = ''
    let firstCharAt = 0
    let lastCharAt = 0

    const clearBuffer = () => {
      buffer = ''
      firstCharAt = 0
      lastCharAt = 0
    }

    const onKeyDown = (event: KeyboardEvent) => {
      const current = optionsRef.current
      if (current.enabled === false || event.defaultPrevented) {
        return
      }

      if (isEditableTarget(event.target)) {
        clearBuffer()
        return
      }

      const keyGapMs = current.keyGapMs ?? DEFAULT_KEY_GAP_MS
      const maxElapsedMs = current.maxElapsedMs ?? DEFAULT_MAX_ELAPSED_MS
      const minBarcodeLength =
        current.minBarcodeLength ?? DEFAULT_MIN_BARCODE_LENGTH
      const maxBufferLength =
        current.maxBufferLength ?? DEFAULT_MAX_BUFFER_LENGTH

      const now = performance.now()

      if (event.key === 'Escape') {
        clearBuffer()
        return
      }

      if (event.key === 'Enter') {
        const trimmed = buffer.trim()
        const elapsed = firstCharAt > 0 ? now - firstCharAt : Number.POSITIVE_INFINITY

        if (trimmed.length >= minBarcodeLength && elapsed <= maxElapsedMs) {
          event.preventDefault()
          const product = current.getByBarcode(trimmed)
          if (product) {
            current.onProductFound(product, trimmed)
          } else {
            current.onBarcodeNotFound(trimmed)
          }
        }

        clearBuffer()
        return
      }

      if (!isScannableCharacter(event)) {
        return
      }

      if (lastCharAt > 0 && now - lastCharAt > maxElapsedMs) {
        clearBuffer()
      }

      if (!buffer || firstCharAt <= 0) {
        buffer = event.key
        firstCharAt = now
        lastCharAt = now
        return
      }

      if (now - lastCharAt > keyGapMs) {
        buffer = event.key
        firstCharAt = now
        lastCharAt = now
        return
      }

      buffer = `${buffer}${event.key}`.slice(-maxBufferLength)
      lastCharAt = now
    }

    window.addEventListener('keydown', onKeyDown)
    return () => {
      window.removeEventListener('keydown', onKeyDown)
    }
  }, [])
}

export default usePOSBarcodeScanner
