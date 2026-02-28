import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { KDSOrderCard } from './KDSOrderCard'
import type { KDSOrder } from '@/stores/kdsStore'

const t = (_path: string, fallback?: string) => fallback ?? ''

const buildOrder = (overrides: Partial<KDSOrder> = {}): KDSOrder => ({
  id: 1,
  tenant_id: 1,
  restaurant_id: 1,
  order_id: 2001,
  order_number: '2001',
  customer_name: 'Ahmed',
  status: 'new',
  priority: 0,
  station_id: 'expo',
  items: [
    {
      product_name: 'Chicken Shawarma',
      quantity: 2,
      special_instructions: 'No pickles',
    },
  ],
  notes: 'Call customer on arrival',
  received_at: '2026-02-24T10:00:00Z',
  started_at: null,
  completed_at: null,
  bumped_by: null,
  elapsed_seconds: 120,
  created_at: '2026-02-24T10:00:00Z',
  updated_at: '2026-02-24T10:00:00Z',
  ...overrides,
})

describe('KDSOrderCard', () => {
  it('renders order details, instructions, and rush badge', () => {
    render(
      <KDSOrderCard
        order={buildOrder({ priority: 1 })}
        locale="en"
        t={t}
        canBump
        canRecall
        onBump={vi.fn()}
        onRecall={vi.fn()}
      />
    )

    expect(screen.getByText('#2001')).toBeTruthy()
    expect(screen.getByText('Rush')).toBeTruthy()
    expect(screen.getByText(/Customer: Ahmed/)).toBeTruthy()
    expect(screen.getByText('2x')).toBeTruthy()
    expect(screen.getByText('Chicken Shawarma')).toBeTruthy()
    expect(screen.getByText('No pickles')).toBeTruthy()
    expect(screen.getByText('Call customer on arrival')).toBeTruthy()
  })

  it('calls bump handler for bumpable statuses', () => {
    const onBump = vi.fn()

    render(
      <KDSOrderCard
        order={buildOrder({ id: 17, status: 'new' })}
        locale="en"
        t={t}
        canBump
        canRecall={false}
        onBump={onBump}
        onRecall={vi.fn()}
      />
    )

    fireEvent.click(screen.getByRole('button', { name: 'Bump' }))
    expect(onBump).toHaveBeenCalledTimes(1)
    expect(onBump).toHaveBeenCalledWith(17)
  })

  it('shows recall for completed orders and calls recall handler', () => {
    const onRecall = vi.fn()

    render(
      <KDSOrderCard
        order={buildOrder({ id: 88, status: 'completed', completed_at: '2026-02-24T10:04:00Z' })}
        locale="en"
        t={t}
        canBump
        canRecall
        onBump={vi.fn()}
        onRecall={onRecall}
      />
    )

    fireEvent.click(screen.getByRole('button', { name: 'Recall' }))
    expect(onRecall).toHaveBeenCalledTimes(1)
    expect(onRecall).toHaveBeenCalledWith(88)
  })

  it('shows busy labels while action is pending', () => {
    render(
      <KDSOrderCard
        order={buildOrder({ status: 'new' })}
        locale="en"
        t={t}
        canBump
        canRecall
        isBusy
        onBump={vi.fn()}
        onRecall={vi.fn()}
      />
    )

    expect(screen.getByText('Bumping...')).toBeTruthy()
  })
})
