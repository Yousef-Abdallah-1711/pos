import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { QRBuilderModal } from './QRBuilderModal'
import type { RestaurantTable } from '@/lib/api/tablesApi'

const toastSuccess = vi.fn()
const toastError = vi.fn()
const clipboardWriteText = vi.fn()

vi.mock('@/lib/config', () => ({
  getConfig: () => ({
    restaurantDomain: 'demo.example.com',
  }),
}))

vi.mock('sonner', () => ({
  toast: {
    success: (...args: unknown[]) => toastSuccess(...args),
    error: (...args: unknown[]) => toastError(...args),
  },
}))

vi.mock('./QRPreview', () => ({
  QRPreview: () => <div data-testid="qr-preview">qr-preview</div>,
}))

const t = (path: string, fallback?: string) => fallback ?? path

const sampleTable: RestaurantTable = {
  id: 5,
  tenant_id: 1,
  restaurant_id: 1,
  table_number: 5,
  name: 'Garden',
  capacity: 4,
  location: 'indoor',
  status: 'available',
  description: '',
  qr_code_config: {},
  position_x: 0,
  position_y: 0,
  shape: 'square',
  layout_id: null,
  created_at: '2026-02-25T10:00:00Z',
  updated_at: '2026-02-25T10:00:00Z',
}

describe('QRBuilderModal', () => {
  beforeEach(() => {
    clipboardWriteText.mockReset()
    toastSuccess.mockReset()
    toastError.mockReset()
  })

  it('returns null when closed', () => {
    const { container } = render(
      <QRBuilderModal
        isOpen={false}
        table={sampleTable}
        locale="en"
        onClose={vi.fn()}
        t={t}
        isRTL={false}
      />
    )

    expect(container.firstChild).toBeNull()
  })

  it('renders QR URL and copies link to clipboard', async () => {
    Object.defineProperty(navigator, 'clipboard', {
      value: {
        writeText: clipboardWriteText,
      },
      configurable: true,
    })
    clipboardWriteText.mockResolvedValue(undefined)

    render(
      <QRBuilderModal
        isOpen
        table={sampleTable}
        locale="en"
        onClose={vi.fn()}
        t={t}
        isRTL={false}
      />
    )

    const expectedUrl = 'https://demo.example.com/en/table/5/menu'
    expect(screen.getByText(expectedUrl)).toBeTruthy()
    expect(screen.getByTestId('qr-preview')).toBeTruthy()

    fireEvent.click(screen.getByRole('button', { name: 'tables.copyLink' }))

    await waitFor(() => {
      expect(clipboardWriteText).toHaveBeenCalledWith(expectedUrl)
      expect(toastSuccess).toHaveBeenCalledWith('tables.linkCopied')
    })
  })
})
