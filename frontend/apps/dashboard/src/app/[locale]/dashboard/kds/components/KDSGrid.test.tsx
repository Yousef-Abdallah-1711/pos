import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { KDSGrid } from './KDSGrid'
import type { KDSOrder } from '@/stores/kdsStore'

const t = (_path: string, fallback?: string) => fallback ?? ''

const order = (id: number, status: KDSOrder['status'], priority: number, receivedAt: string): KDSOrder => ({
  id,
  tenant_id: 1,
  restaurant_id: 1,
  order_id: 1000 + id,
  order_number: String(1000 + id),
  customer_name: `Customer ${id}`,
  status,
  priority,
  station_id: 'expo',
  items: [{ product_name: `Item ${id}`, quantity: 1 }],
  notes: '',
  received_at: receivedAt,
  started_at: null,
  completed_at: status === 'completed' ? '2026-02-24T10:10:00Z' : null,
  bumped_by: null,
  elapsed_seconds: 30,
  created_at: receivedAt,
  updated_at: receivedAt,
})

describe('KDSGrid', () => {
  it('shows empty state when there are no cards', () => {
    render(
      <KDSGrid
        orders={[]}
        locale="en"
        t={t}
        canBump
        canRecall
        pendingOrderIds={[]}
        onBump={vi.fn()}
        onRecall={vi.fn()}
      />
    )

    expect(screen.getByText('No active orders')).toBeTruthy()
    expect(screen.getByText('New kitchen orders will appear here as soon as they are confirmed.')).toBeTruthy()
  })

  it('sorts cards by priority then received time', async () => {
    render(
      <KDSGrid
        orders={[
          order(1, 'new', 0, '2026-02-24T10:00:00Z'),
          order(2, 'new', 1, '2026-02-24T10:05:00Z'),
          order(3, 'new', 0, '2026-02-24T10:08:00Z'),
        ]}
        locale="en"
        t={t}
        canBump
        canRecall
        pendingOrderIds={[]}
        onBump={vi.fn()}
        onRecall={vi.fn()}
      />
    )

    await waitFor(() => {
      expect(screen.getByText('#1002')).toBeTruthy()
      expect(screen.getByText('#1001')).toBeTruthy()
      expect(screen.getByText('#1003')).toBeTruthy()
    })

    const headings = screen.getAllByRole('heading', { level: 3 }).map((el) => el.textContent)
    expect(headings[0]).toBe('#1002')
    expect(headings[1]).toBe('#1001')
    expect(headings[2]).toBe('#1003')
  })

  it('renders busy state for pending bump action', async () => {
    render(
      <KDSGrid
        orders={[order(10, 'new', 0, '2026-02-24T10:00:00Z')]}
        locale="en"
        t={t}
        canBump
        canRecall
        pendingOrderIds={[10]}
        onBump={vi.fn()}
        onRecall={vi.fn()}
      />
    )

    await waitFor(() => {
      expect(screen.getByText('Bumping...')).toBeTruthy()
    })
  })

  it('wires bump and recall actions to callbacks', async () => {
    const onBump = vi.fn()
    const onRecall = vi.fn()

    render(
      <KDSGrid
        orders={[
          order(21, 'new', 0, '2026-02-24T10:00:00Z'),
          order(22, 'completed', 0, '2026-02-24T09:55:00Z'),
        ]}
        locale="en"
        t={t}
        canBump
        canRecall
        pendingOrderIds={[]}
        onBump={onBump}
        onRecall={onRecall}
      />
    )

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Bump' })).toBeTruthy()
      expect(screen.getByRole('button', { name: 'Recall' })).toBeTruthy()
    })

    fireEvent.click(screen.getByRole('button', { name: 'Bump' }))
    fireEvent.click(screen.getByRole('button', { name: 'Recall' }))

    expect(onBump).toHaveBeenCalledWith(21)
    expect(onRecall).toHaveBeenCalledWith(22)
  })
})
