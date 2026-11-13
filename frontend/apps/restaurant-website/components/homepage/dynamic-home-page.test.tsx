import { render, screen } from '@testing-library/react'
import { DynamicHomePage } from './dynamic-home-page'
import type { ThemeConfig } from '@/types/theme'

const emptyTheme: ThemeConfig = {
  name: 'Default Theme',
  slug: 'default-theme',
  templateId: 'template-01',
  colors: { primary: '#2563eb', accent: '#0ea5e9' },
  typography: {},
  identity: { siteTitle: 'WS05 Tenant A Restaurant' },
  header: {},
  footer: {},
  components: [],
}

describe('DynamicHomePage', () => {
  it('renders a branded ordering hero when the saved theme has no sections', () => {
    render(
      <DynamicHomePage
        locale="en"
        restaurantSlug="phase2a-restaurant-a"
        themeConfig={emptyTheme}
        products={[]}
        categories={[]}
      />
    )

    expect(
      screen.getByRole('heading', { level: 1, name: 'WS05 Tenant A Restaurant' })
    ).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Order Now' })).toHaveAttribute('href', '/en/menu')
  })

  it('uses the Arabic site title for an empty Arabic storefront', () => {
    render(
      <DynamicHomePage
        locale="ar"
        restaurantSlug="phase2a-restaurant-a"
        themeConfig={{ ...emptyTheme, identity: { siteTitle: { ar: 'مطعم الاختبار' } } }}
        products={[]}
        categories={[]}
      />
    )

    expect(screen.getByRole('heading', { level: 1, name: 'مطعم الاختبار' })).toBeInTheDocument()
  })
})
