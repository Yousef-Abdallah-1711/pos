'use client'

import { useEffect, useRef } from 'react'
import type { POSProduct } from '@/lib/api/posApi'

export interface POSKeyboardOptions {
  visibleProducts: POSProduct[]
  cartItemsCount: number
  isSubmitting: boolean
  heldOrdersCount: number
  onQuickAddProduct: (product: POSProduct) => void
  onSubmitOrder: () => void
  onClearCart: () => void
  onHoldOrder: () => void
  onRecallLastHeldOrder: () => void
  onFocusSearch?: () => void
}

const isEditableTarget = (target: EventTarget | null): boolean => {
  if (!(target instanceof HTMLElement)) {
    return false
  }

  const tag = target.tagName.toUpperCase()
  if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') {
    return true
  }

  return target.isContentEditable
}

const isFunctionKey = (key: string) => /^F([1-9]|1[0-2])$/i.test(key)

const getFunctionIndex = (key: string): number | null => {
  if (!isFunctionKey(key)) {
    return null
  }
  const numeric = Number(key.slice(1))
  if (!Number.isFinite(numeric) || numeric < 1 || numeric > 12) {
    return null
  }
  return numeric - 1
}

export const usePOSKeyboard = (options: POSKeyboardOptions) => {
  const optionsRef = useRef(options)

  useEffect(() => {
    optionsRef.current = options
  }, [options])

  useEffect(() => {
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.defaultPrevented) {
        return
      }

      if (isEditableTarget(event.target)) {
        return
      }

      const current = optionsRef.current
      const key = event.key

      const functionIndex = getFunctionIndex(key)
      if (functionIndex !== null) {
        const product = current.visibleProducts[functionIndex]
        if (product && product.is_available !== false) {
          event.preventDefault()
          current.onQuickAddProduct(product)
        }
        return
      }

      if (key === ' ') {
        if (current.cartItemsCount > 0 && !current.isSubmitting) {
          event.preventDefault()
          current.onSubmitOrder()
        }
        return
      }

      if (key === 'Escape') {
        event.preventDefault()
        current.onClearCart()
        return
      }

      if (event.ctrlKey && !event.shiftKey && !event.altKey && key.toLowerCase() === 'h') {
        if (current.cartItemsCount > 0) {
          event.preventDefault()
          current.onHoldOrder()
        }
        return
      }

      if (event.ctrlKey && !event.shiftKey && !event.altKey && key.toLowerCase() === 'r') {
        if (current.heldOrdersCount > 0) {
          event.preventDefault()
          current.onRecallLastHeldOrder()
        }
        return
      }

      if ((event.ctrlKey && !event.shiftKey && !event.altKey && key.toLowerCase() === 'f') ||
          (event.ctrlKey && !event.shiftKey && !event.altKey && key.toLowerCase() === 'k') ||
          (!event.ctrlKey && !event.metaKey && !event.altKey && key === '/')) {
        event.preventDefault()
        current.onFocusSearch?.()
      }
    }

    window.addEventListener('keydown', onKeyDown)
    return () => {
      window.removeEventListener('keydown', onKeyDown)
    }
  }, [])
}

export default usePOSKeyboard
