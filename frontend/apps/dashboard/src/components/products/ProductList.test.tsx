import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import '@testing-library/jest-dom/vitest'
import { ProductList } from './ProductList'
import type { ProductListFilters } from './product-types'

vi.mock('next/navigation', () => ({ usePathname: () => '/en/dashboard/products' }))
afterEach(cleanup)
const filters: ProductListFilters = { page: 1, pageSize: 20, categoryId: null, searchText: '', sortBy: 'created_at', sortOrder: 'DESC' }

describe('product server pagination', () => {
  it('does not invent pages when the server reports one page', () => {
    render(<ProductList products={[{ id: 1, name_en: 'Synthetic product' }]} categories={[]} loading={false} filters={filters} onFilterChange={vi.fn()} totalPages={1} />)
    expect(screen.getByRole('link', { name: '1' })).toHaveAttribute('aria-current', 'page')
    expect(screen.queryByRole('link', { name: '2' })).not.toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Go to next page' })).toHaveAttribute('aria-disabled', 'true')
  })

  it('supports later pages and requests the selected valid page', async () => {
    const change = vi.fn()
    render(<ProductList products={[{ id: 1, name_en: 'Synthetic product' }]} categories={[]} loading={false} filters={{ ...filters, page: 8 }} onFilterChange={change} totalPages={10} />)
    await userEvent.click(screen.getByRole('link', { name: '9' }))
    expect(change).toHaveBeenCalledWith({ ...filters, page: 9 })
    expect(screen.getByRole('link', { name: '10' })).toBeInTheDocument()
    expect(screen.queryByRole('link', { name: '11' })).not.toBeInTheDocument()
  })
})
