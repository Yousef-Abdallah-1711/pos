import { renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { usePOSBarcodeScanner } from './usePOSBarcodeScanner'
import type { POSProduct } from '@/lib/api/posApi'

const dispatchKeyDown = (key: string, target: EventTarget = window) => {
  const event = new KeyboardEvent('keydown', { key, bubbles: true })
  target.dispatchEvent(event)
}

const makeProduct = (id: number): POSProduct => ({
  id,
  name: `Product ${id}`,
  price: 10,
  barcode: `${id}`,
})

describe('usePOSBarcodeScanner', () => {
  let now = 0
  let nowSpy: ReturnType<typeof vi.spyOn>

  beforeEach(() => {
    now = 0
    nowSpy = vi.spyOn(performance, 'now').mockImplementation(() => now)
  })

  afterEach(() => {
    nowSpy.mockRestore()
  })

  it('treats rapid input + Enter as barcode and resolves product', () => {
    const product = makeProduct(123456)
    const getByBarcode = vi.fn().mockReturnValue(product)
    const onProductFound = vi.fn()
    const onBarcodeNotFound = vi.fn()

    renderHook(() =>
      usePOSBarcodeScanner({
        getByBarcode,
        onProductFound,
        onBarcodeNotFound,
      })
    )

    for (const char of '123456') {
      now += 10
      dispatchKeyDown(char)
    }
    now += 10
    dispatchKeyDown('Enter')

    expect(getByBarcode).toHaveBeenCalledWith('123456')
    expect(onProductFound).toHaveBeenCalledTimes(1)
    expect(onProductFound).toHaveBeenCalledWith(product, '123456')
    expect(onBarcodeNotFound).not.toHaveBeenCalled()
  })

  it('ignores slow typing (not scanner-like)', () => {
    const getByBarcode = vi.fn()
    const onProductFound = vi.fn()
    const onBarcodeNotFound = vi.fn()

    renderHook(() =>
      usePOSBarcodeScanner({
        getByBarcode,
        onProductFound,
        onBarcodeNotFound,
      })
    )

    for (const char of '123456') {
      now += 120
      dispatchKeyDown(char)
    }
    now += 120
    dispatchKeyDown('Enter')

    expect(getByBarcode).not.toHaveBeenCalled()
    expect(onProductFound).not.toHaveBeenCalled()
    expect(onBarcodeNotFound).not.toHaveBeenCalled()
  })

  it('calls not-found callback for unknown rapid barcode', () => {
    const getByBarcode = vi.fn().mockReturnValue(undefined)
    const onProductFound = vi.fn()
    const onBarcodeNotFound = vi.fn()

    renderHook(() =>
      usePOSBarcodeScanner({
        getByBarcode,
        onProductFound,
        onBarcodeNotFound,
      })
    )

    for (const char of '987654') {
      now += 8
      dispatchKeyDown(char)
    }
    now += 8
    dispatchKeyDown('Enter')

    expect(getByBarcode).toHaveBeenCalledWith('987654')
    expect(onProductFound).not.toHaveBeenCalled()
    expect(onBarcodeNotFound).toHaveBeenCalledTimes(1)
    expect(onBarcodeNotFound).toHaveBeenCalledWith('987654')
  })

  it('ignores scanner capture when typing in input fields', () => {
    const getByBarcode = vi.fn()
    const onProductFound = vi.fn()
    const onBarcodeNotFound = vi.fn()

    renderHook(() =>
      usePOSBarcodeScanner({
        getByBarcode,
        onProductFound,
        onBarcodeNotFound,
      })
    )

    const input = document.createElement('input')
    document.body.appendChild(input)
    input.focus()

    for (const char of '123456') {
      now += 5
      dispatchKeyDown(char, input)
    }
    now += 5
    dispatchKeyDown('Enter', input)

    expect(getByBarcode).not.toHaveBeenCalled()
    expect(onProductFound).not.toHaveBeenCalled()
    expect(onBarcodeNotFound).not.toHaveBeenCalled()

    document.body.removeChild(input)
  })
})

