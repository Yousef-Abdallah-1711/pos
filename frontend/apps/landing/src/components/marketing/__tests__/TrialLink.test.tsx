import { render, screen } from '@testing-library/react'
import type { ReactNode } from 'react'
import { vi } from 'vitest'
import { TrialLink } from '../TrialLink'

vi.mock('next/link', () => ({
  default: ({ children, href, prefetch }: { children: ReactNode; href: string; prefetch?: boolean }) => (
    <a href={href} data-prefetch={String(prefetch)}>{children}</a>
  ),
}))

it('does not prefetch the registration form from marketing calls to action', () => {
  render(<TrialLink locale="en" source="home-hero" />)
  expect(screen.getByRole('link', { name: /create your restaurant account/i })).toHaveAttribute('href', '/en/auth/register?billing=monthly')
  expect(screen.getByRole('link', { name: /create your restaurant account/i })).toHaveAttribute('data-prefetch', 'false')
})
