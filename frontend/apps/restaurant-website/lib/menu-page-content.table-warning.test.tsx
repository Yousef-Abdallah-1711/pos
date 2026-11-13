import { render, screen, waitFor } from '@testing-library/react'
import { MenuPageContent } from '@/app/[locale]/menu/menu-page-content'
import { useCartStore } from '@/lib/store/cart-store'
import { SubdomainProvider } from '@/lib/subdomain-context'

jest.mock('@/lib/hooks/use-api-queries', () => ({
  useProducts: jest.fn(() => ({
    data: [],
    isLoading: false,
  })),
}))

jest.mock('@/components/commerce/product-card-simple', () => ({
  ProductCard: () => null,
}))

describe('MenuPageContent table warning flow', () => {
  beforeEach(() => {
    localStorage.clear()
    useCartStore.setState({
      items: [],
      tableId: null,
      tableName: null,
    })
  })

  it('shows occupied warning and persists table context from QR landing data', async () => {
    render(
      <SubdomainProvider>
        <MenuPageContent
          locale="en"
          restaurantSlug="demo"
          themeData={{ colors: {} } as any}
          tableContext={{
            tableId: 8,
            tableNumber: 8,
            tableName: 'Garden',
            tableStatus: 'occupied',
          }}
          tableStatusNotice={{
            status: 'occupied',
            level: 'warning',
            translationKey: 'tableOrdering.statusWarnings.occupied',
          }}
        />
      </SubdomainProvider>
    )

    expect(screen.getByText('Dine-in Order - Table #8')).toBeInTheDocument()
    expect(
      screen.getByText('This table is currently marked as occupied. You can still place your order.')
    ).toBeInTheDocument()

    await waitFor(() => {
      expect(useCartStore.getState().tableId).toBe(8)
      expect(useCartStore.getState().tableName).toBe('Garden')
    })
  })
})
