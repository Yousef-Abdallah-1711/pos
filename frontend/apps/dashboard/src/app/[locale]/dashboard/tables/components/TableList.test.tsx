import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { TableList } from './TableList'
import type { RestaurantTable } from '@/lib/api/tablesApi'

const t = (path: string, fallback?: string) => fallback ?? path

const sampleTable: RestaurantTable = {
  id: 15,
  tenant_id: 1,
  restaurant_id: 1,
  table_number: 15,
  name: 'Patio',
  capacity: 4,
  location: 'indoor',
  status: 'occupied',
  description: 'Near window',
  qr_code_config: {},
  position_x: 100,
  position_y: 120,
  shape: 'square',
  layout_id: null,
  created_at: '2026-02-25T10:00:00Z',
  updated_at: '2026-02-25T10:00:00Z',
}

describe('TableList', () => {
  it('renders table row and supports occupied quick action', () => {
    const onFreeTable = vi.fn()
    const onStatusChange = vi.fn()
    const onView = vi.fn()
    const onOpenQr = vi.fn()

    render(
      <TableList
        tables={[sampleTable]}
        loading={false}
        canManageTables
        canDeleteTables={false}
        canViewTables
        onEdit={vi.fn()}
        onDelete={vi.fn()}
        onView={onView}
        onOpenPOS={vi.fn()}
        onOpenQr={onOpenQr}
        onFreeTable={onFreeTable}
        onStatusChange={onStatusChange}
        onOpenPOS={vi.fn()}
        t={t}
      />
    )

    expect(screen.getByText('#15')).toBeTruthy()
    expect(screen.getByText('Patio')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'tables.freeTable' })).toBeTruthy()

    fireEvent.click(screen.getByRole('button', { name: 'tables.freeTable' }))
    expect(onFreeTable).toHaveBeenCalledWith(15)

    fireEvent.change(screen.getByRole('combobox'), {
      target: { value: 'available' },
    })
    expect(onStatusChange).toHaveBeenCalledWith(15, 'available', 'occupied')

    fireEvent.click(screen.getByTitle('tables.viewDetails'))
    expect(onView).toHaveBeenCalledWith(15)

    fireEvent.click(screen.getByTitle('tables.qrCode'))
    expect(onOpenQr).toHaveBeenCalledWith(sampleTable)
  })

  it('shows empty state when there are no tables', () => {
    render(
      <TableList
        tables={[]}
        loading={false}
        canManageTables={false}
        canDeleteTables={false}
        canViewTables={false}
        onEdit={vi.fn()}
        onDelete={vi.fn()}
        onView={vi.fn()}
        onOpenPOS={vi.fn()}
        onOpenQr={vi.fn()}
        onFreeTable={vi.fn()}
        onStatusChange={vi.fn()}
        onOpenPOS={vi.fn()}
        t={t}
      />
    )

    expect(screen.getByText('tables.empty')).toBeTruthy()
  })
})
