import React from 'react'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import POSPage from './page'
import { usePOSStore } from '@/stores/posStore'
import { useAuthStore } from '@/stores/authStore'

const mockCreateOrder = vi.fn()
const mockUpdatePaymentStatus = vi.fn()
const mockPrintReceipt = vi.fn(async (_receipt: {
  orderNumber: string
  paymentMethod: string
  items: Array<{ name: string }>
}) => true)
const mockPlaySound = vi.fn(async () => true)

vi.mock('next/navigation', () => ({
  usePathname: () => '/en/dashboard/pos',
  useSearchParams: () => new URLSearchParams(),
}))

vi.mock('sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}))

vi.mock('./components/POSTerminal', () => ({
  POSTerminal: ({
    onSubmitOrder,
    autoPrintEnabled,
  }: {
    onSubmitOrder: () => void
    autoPrintEnabled: boolean
  }) => (
    <div>
      <div data-testid="auto-print-state">{autoPrintEnabled ? 'on' : 'off'}</div>
      <button type="button" onClick={onSubmitOrder}>
        submit-order
      </button>
    </div>
  ),
}))

vi.mock('./components/POSReceipt', () => ({
  POSReceipt: () => null,
}))

vi.mock('@/hooks/usePOSData', () => ({
  usePOSData: () => ({
    products: [
      {
        id: 1,
        name: 'Burger',
        price: 12,
        barcode: '111111',
      },
    ],
    categories: [],
    tables: [],
    loading: false,
    error: null,
    refresh: vi.fn(),
  }),
}))

vi.mock('@/hooks/usePOSOffline', () => ({
  usePOSOffline: () => ({
    isOffline: false,
    queuedOrderCount: 0,
    isReplaying: false,
    enqueueOrder: vi.fn(async () => true),
    replayQueuedOrders: vi.fn(async () => ({ processed: 0, failed: 0 })),
  }),
}))

vi.mock('@/hooks/usePOSKeyboard', () => ({
  usePOSKeyboard: vi.fn(),
}))

vi.mock('@/hooks/usePOSBarcodeScanner', () => ({
  usePOSBarcodeScanner: vi.fn(),
}))

vi.mock('@/hooks/usePOSBackgroundSync', () => ({
  usePOSBackgroundSync: () => ({
    registerOrderSync: vi.fn(async () => false),
    syncSupported: false,
    syncTag: 'pos-order-sync',
  }),
}))

vi.mock('@/hooks/usePOSPerformance', () => ({
  usePOSPerformance: () => ({
    isDevMode: false,
    recentMetrics: [],
    markPageLoad: vi.fn(),
    measureSync: (_metric: string, callback: () => unknown) => callback(),
    measureAsync: async (_metric: string, callback: () => Promise<unknown>) => callback(),
    recordMetric: vi.fn(),
  }),
}))

vi.mock('@/hooks/usePOSAudio', () => ({
  usePOSAudio: () => ({
    playSound: mockPlaySound,
    unlockAudio: vi.fn(async () => true),
  }),
}))

vi.mock('@/hooks/usePOSPrint', () => ({
  usePOSPrint: () => ({
    printReceipt: mockPrintReceipt,
  }),
}))

vi.mock('@/hooks/usePOSProductIndex', () => ({
  usePOSProductIndex: (products: unknown[]) => ({
    search: () => products,
    getByBarcode: vi.fn(),
  }),
}))

vi.mock('@/hooks/useRbac', () => ({
  useHasPermission: () => ({
    hasPermission: true,
    loading: false,
  }),
}))

vi.mock('@/lib/translations', () => ({
  createTranslator:
    () =>
    (_path: string, fallback?: string) =>
      fallback ?? '',
  getLocaleFromPath: () => 'en',
}))

vi.mock('@/lib/api/posApi', () => ({
  posApi: {
    createOrder: (...args: unknown[]) => mockCreateOrder(...args),
    updateOrderPaymentStatus: (...args: unknown[]) => mockUpdatePaymentStatus(...args),
  },
}))

describe('POS page Phase 3 auto-print integration', () => {
  beforeEach(() => {
    vi.clearAllMocks()

    useAuthStore.setState({
      user: {
        id: 7,
        name: 'QA Cashier',
        email: 'qa@example.com',
        role: 'cashier',
        tenant_id: 1,
      },
      token: 'token',
    })

    usePOSStore.setState({
      items: [
        {
          productId: 1,
          name: 'Burger',
          price: 12,
          quantity: 1,
        },
      ],
      paymentMethod: 'cash',
      customerName: 'Walk-in',
      customerPhone: '',
      tableId: null,
      tableName: null,
      orderNotes: '',
      autoPrintEnabled: false,
    })

    mockCreateOrder.mockResolvedValue({
      id: 900,
      order_number: 'ORD-900',
      status: 'pending',
      total_amount: 13.68,
      created_at: '2026-02-26T10:00:00.000Z',
    })
    mockUpdatePaymentStatus.mockResolvedValue(undefined)
  })

  it('does not print when auto-print is disabled', async () => {
    usePOSStore.getState().setAutoPrintEnabled(false)
    render(<POSPage />)

    expect(screen.getByTestId('auto-print-state').textContent).toBe('off')
    fireEvent.click(screen.getByRole('button', { name: 'submit-order' }))

    await waitFor(() => {
      expect(mockCreateOrder).toHaveBeenCalledTimes(1)
    })
    expect(mockPrintReceipt).not.toHaveBeenCalled()
  })

  it('prints receipt when auto-print is enabled', async () => {
    usePOSStore.getState().setAutoPrintEnabled(true)
    render(<POSPage />)

    expect(screen.getByTestId('auto-print-state').textContent).toBe('on')
    fireEvent.click(screen.getByRole('button', { name: 'submit-order' }))

    await waitFor(() => {
      expect(mockCreateOrder).toHaveBeenCalledTimes(1)
    })

    await waitFor(() => {
      expect(mockPrintReceipt).toHaveBeenCalledTimes(1)
    })

    const receiptArg = mockPrintReceipt.mock.calls[0]?.[0]
    expect(receiptArg.orderNumber).toBe('ORD-900')
    expect(receiptArg.paymentMethod).toBe('cash')
    expect(Array.isArray(receiptArg.items)).toBe(true)
    expect(receiptArg.items[0].name).toBe('Burger')
  })
})
