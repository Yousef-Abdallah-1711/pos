import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { QRBuilderModal } from './QRBuilderModal'
import type { RestaurantTable } from '@/lib/api/tablesApi'

vi.mock('@/lib/config', () => ({
  getConfig: () => ({
    restaurantDomain: 'demo.example.com',
  }),
}))

vi.mock('sonner', () => ({
  toast: {
    success: () => undefined,
    error: () => undefined,
  },
}))

vi.mock('./QRPreview', () => ({
  QRPreview: () => {
    throw new Error('QR render failed')
  },
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

describe('QRBuilderModal fallback', () => {
  it('shows fallback UI when QR preview rendering fails', () => {
    render(
      <QRBuilderModal
        isOpen
        table={sampleTable}
        locale="en"
        onClose={() => undefined}
        t={t}
        isRTL={false}
      />
    )

    expect(screen.getAllByText('tables.qrUnavailable').length).toBeGreaterThan(0)
  })
})
