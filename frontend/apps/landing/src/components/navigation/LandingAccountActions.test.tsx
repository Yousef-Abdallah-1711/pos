import { fireEvent, render, screen } from '@testing-library/react'
import type { AnchorHTMLAttributes, ReactNode } from 'react'
import { vi } from 'vitest'
import { LandingAccountActions } from './LandingAccountActions'

vi.mock('next/link', () => ({
  default: ({ children, href, prefetch: _prefetch, ...props }: { children: ReactNode; href: string; prefetch?: boolean } & AnchorHTMLAttributes<HTMLAnchorElement>) => <a href={href} {...props}>{children}</a>,
}))

const authenticatedSession = {
  status: 'authenticated' as const,
  profile: { id: 42, name: 'Sample Manager', email: 'manager@example.test', restaurantName: 'Sample Cafe' },
}

describe('landing account actions', () => {
  it('keeps guest actions hidden while the first session check is pending', () => {
    render(<LandingAccountActions locale="en" session={{ status: 'checking' }} isSigningOut={false} onRetry={vi.fn()} onSignOut={async () => true} />)

    expect(screen.getByRole('status')).toHaveTextContent('Checking sign-in')
    expect(screen.queryByRole('link', { name: 'Log in' })).not.toBeInTheDocument()
    expect(screen.queryByRole('link', { name: 'Sign up' })).not.toBeInTheDocument()
  })

  it('shows authorized account details and canonical dashboard links in the disclosure', () => {
    render(<LandingAccountActions locale="ar" session={authenticatedSession} isSigningOut={false} onRetry={vi.fn()} onSignOut={async () => true} />)
    const trigger = screen.getByRole('button', { name: 'قائمة الحساب لـ Sample Manager' })

    expect(trigger).toHaveAttribute('aria-expanded', 'false')
    fireEvent.click(trigger)

    expect(screen.getByText('manager@example.test')).toHaveAttribute('dir', 'auto')
    expect(screen.getByText('Sample Cafe')).toHaveAttribute('dir', 'auto')
    expect(screen.getByRole('link', { name: 'لوحة التحكم' })).toHaveAttribute('href', 'http://localhost:3002/ar/dashboard')
    expect(screen.getByRole('link', { name: 'إعدادات الحساب' })).toHaveAttribute('href', 'http://localhost:3002/ar/dashboard/settings/account')
    expect(trigger).toHaveAttribute('aria-expanded', 'true')

    fireEvent.keyDown(document, { key: 'Escape' })
    expect(trigger).toHaveAttribute('aria-expanded', 'false')
    expect(trigger).toHaveFocus()
  })

  it('offers a retry without presenting a temporary API failure as a guest session', () => {
    const onRetry = vi.fn()
    render(<LandingAccountActions locale="en" session={{ status: 'temporarily-unavailable' }} isSigningOut={false} onRetry={onRetry} onSignOut={async () => true} />)

    expect(screen.getByRole('status')).toHaveTextContent('Sign-in status is temporarily unavailable.')
    expect(screen.getByRole('link', { name: 'Log in' })).toHaveAttribute('href', '/en/auth/login')
    expect(screen.getByRole('link', { name: 'Sign up' })).toHaveAttribute('href', '/en/auth/register')
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }))
    expect(onRetry).toHaveBeenCalledOnce()
  })

  it('shows localized login and registration links only after the session is confirmed anonymous', () => {
    render(<LandingAccountActions locale="ar" session={{ status: 'guest' }} isSigningOut={false} onRetry={vi.fn()} onSignOut={async () => true} />)

    expect(screen.getByRole('link', { name: 'تسجيل الدخول' })).toHaveAttribute('href', '/ar/auth/login')
    expect(screen.getByRole('link', { name: 'إنشاء حساب' })).toHaveAttribute('href', '/ar/auth/register')
  })
})
