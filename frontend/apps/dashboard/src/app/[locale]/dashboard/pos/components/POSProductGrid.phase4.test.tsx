import React from 'react'
import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { POSProductGrid } from './POSProductGrid'
import type { POSProduct } from '@/lib/api/posApi'

vi.mock('./POSProductCard', () => ({
  POSProductCard: ({ product }: { product: POSProduct }) => (
    <div data-testid="pos-product-card">{product.name}</div>
  ),
}))

vi.mock('./POSVirtualGrid', () => ({
  POSVirtualGrid: ({ products }: { products: POSProduct[] }) => (
    <div data-testid="pos-virtual-grid">{products.length}</div>
  ),
}))

const t = (_path: string, fallback?: string) => fallback ?? ''

const makeProducts = (count: number): POSProduct[] =>
  Array.from({ length: count }).map((_, index) => ({
    id: index + 1,
    name: `Product ${index + 1}`,
    price: 10 + index,
  }))

describe('POSProductGrid Phase 4 virtualization', () => {
  it('keeps regular grid when product count is 50 or less', () => {
    render(
      <POSProductGrid
        locale="en"
        t={t}
        products={makeProducts(50)}
        onAddProduct={vi.fn()}
      />
    )

    expect(screen.queryByTestId('pos-virtual-grid')).toBeNull()
    expect(screen.getAllByTestId('pos-product-card')).toHaveLength(50)
  })

  it('switches to virtual grid when product count exceeds 50', () => {
    render(
      <POSProductGrid
        locale="en"
        t={t}
        products={makeProducts(51)}
        onAddProduct={vi.fn()}
      />
    )

    const virtualGrid = screen.getByTestId('pos-virtual-grid')
    expect(virtualGrid.textContent).toBe('51')
    expect(screen.queryByTestId('pos-product-card')).toBeNull()
  })
})

