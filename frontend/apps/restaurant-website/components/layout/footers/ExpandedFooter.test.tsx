import { render, screen } from '@testing-library/react'
import { ExpandedFooter } from './ExpandedFooter'

it('renders Arabic footer controls and navigable localized links without invented opening hours', () => {
  render(<ExpandedFooter locale="ar" config={{ companyName: 'مطعم تجريبي', legalLinks: [{ label: 'القائمة', href: '/menu' }] }} />)
  expect(screen.getByRole('heading', { name: 'روابط سريعة' })).toBeVisible()
  expect(screen.getByRole('heading', { name: 'ساعات العمل' })).toBeVisible()
  expect(screen.getByRole('heading', { name: 'تواصل معنا' })).toBeVisible()
  for (const link of screen.getAllByRole('link', { name: 'القائمة' })) {
    expect(link).toHaveAttribute('href', '/ar/menu')
  }
  expect(screen.getByRole('button', { name: 'الأعلى' })).toBeVisible()
  expect(screen.queryByText(/10:00 AM/)).not.toBeInTheDocument()
})
