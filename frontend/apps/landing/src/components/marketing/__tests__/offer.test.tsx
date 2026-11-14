import { render, screen, fireEvent } from '@testing-library/react'
import { vi } from 'vitest'
import { PricingOffer } from '../PricingOffer'
import { ProductExplorer } from '../ProductExplorer'
import { getFeatures } from '@/content/features'

beforeEach(() => {
  vi.stubGlobal('IntersectionObserver', class { observe() {} disconnect() {} })
})
afterEach(() => vi.unstubAllGlobals())

it('shows the annual charge and carries the preference into trial registration', () => {
  render(<PricingOffer locale="en" detailed />)
  fireEvent.click(screen.getByRole('radio', { name: 'Annual' }))
  expect(screen.getByText('5,000')).toBeInTheDocument()
  expect(screen.getByText('Save 4,000 EGP a year · approximately 44.4%')).toBeInTheDocument()
  expect(screen.getByRole('link', { name: /start your 14-day/i })).toHaveAttribute('href', '/en/auth/register?billing=annual')
  expect(screen.getByText(/Paid checkout is not available here/)).toBeInTheDocument()
})

it('supports keyboard navigation through the product walkthrough', () => {
  render(<ProductExplorer locale="en" features={getFeatures('en')} />)
  const first = screen.getByRole('tab', { name: 'Cashier POS' })
  fireEvent.keyDown(first, { key: 'ArrowRight' })
  expect(screen.getByRole('tab', { name: 'Kitchen display' })).toHaveAttribute('aria-selected', 'true')
  expect(screen.getByRole('tabpanel')).toHaveTextContent('A clear next action for every ticket.')
})
