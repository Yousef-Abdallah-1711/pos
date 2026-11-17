import { useRef, useState } from 'react'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import type { AnchorHTMLAttributes, ReactNode } from 'react'
import { vi } from 'vitest'
import { Footer } from '@/components/Footer'
import type { LandingSessionState } from '@/lib/auth-session'
import { DesktopMegaMenu } from '../DesktopMegaMenu'
import { MobileNavigation } from '../MobileNavigation'

vi.mock('next/link', () => ({
  default: ({ children, href, prefetch: _prefetch, ...props }: { children: ReactNode; href: string; prefetch?: boolean } & AnchorHTMLAttributes<HTMLAnchorElement>) => <a href={href} {...props}>{children}</a>,
}))

function MobileNavigationHarness({ onClose, onToggleTheme, onSwitchLocale, session = { status: 'guest' } }: { onClose: () => void; onToggleTheme: () => void; onSwitchLocale: () => void; session?: LandingSessionState }) {
  const [open, setOpen] = useState(true)
  const triggerRef = useRef<HTMLButtonElement>(null)
  return <><button ref={triggerRef} onClick={() => setOpen(true)}>Open</button><MobileNavigation locale="ar" pathname="/ar/features/pos" open={open}
    onClose={() => { onClose(); setOpen(false) }} triggerRef={triggerRef} theme="dark" onToggleTheme={onToggleTheme} onSwitchLocale={onSwitchLocale}
    session={session} isSigningOut={false} onRetrySession={vi.fn()} onSignOut={async () => true} />
  </>
}

describe('marketing navigation interactions', () => {
  it('opens the desktop disclosure, exposes all feature routes and closes on Escape', () => {
    render(<DesktopMegaMenu locale="en" pathname="/en/features/pos" />)
    const trigger = screen.getByRole('button', { name: /open platform features/i })
    fireEvent.click(trigger)

    expect(trigger).toHaveAttribute('aria-expanded', 'true')
    const expectedFeatures = [
      ['Cashier POS', 'pos'],
      ['Kitchen display', 'kitchen-display'],
      ['Online ordering', 'online-ordering'],
      ['Storefront customization', 'website-builder'],
      ['Inventory', 'inventory'],
      ['Staff records', 'hr'],
      ['Operational overview', 'analytics'],
    ]
    for (const [label, slug] of expectedFeatures) {
      expect(screen.getByRole('link', { name: new RegExp(label, 'i') })).toHaveAttribute('href', `/en/features/${slug}`)
    }
    expect(screen.getByRole('link', { name: /explore all features/i })).toHaveAttribute('href', '/en/features')

    fireEvent.keyDown(document, { key: 'Escape' })
    expect(trigger).toHaveAttribute('aria-expanded', 'false')
    expect(trigger).toHaveFocus()
  })

  it('expands the mobile feature group and closes the drawer on Escape', async () => {
    const onClose = vi.fn()
    const onToggleTheme = vi.fn()
    const onSwitchLocale = vi.fn()
    render(<MobileNavigationHarness onClose={onClose} onToggleTheme={onToggleTheme} onSwitchLocale={onSwitchLocale} />)

    fireEvent.click(screen.getByRole('button', { name: /عرض روابط الميزات/i }))
    expect(screen.getByRole('link', { name: /نقاط البيع للكاشير/i })).toHaveAttribute('href', '/ar/features/pos')
    expect(screen.getByRole('button', { name: /إخفاء روابط الميزات/i })).toHaveAttribute('aria-current', 'page')
    expect(screen.getByRole('link', { name: /نقاط البيع للكاشير/i })).toHaveAttribute('aria-current', 'page')
    expect(screen.getByRole('link', { name: /النظرة التشغيلية/i })).toHaveAttribute('href', '/ar/features/analytics')
    fireEvent.click(screen.getByRole('button', { name: /التبديل إلى الوضع الفاتح/i }))
    expect(onToggleTheme).toHaveBeenCalledOnce()
    fireEvent.keyDown(document, { key: 'Escape' })
    expect(onClose).toHaveBeenCalledOnce()
    await waitFor(() => expect(screen.getByRole('button', { name: 'Open' })).toHaveFocus())
    fireEvent.click(screen.getByRole('button', { name: 'Open' }))
    fireEvent.click(screen.getByRole('button', { name: /التغيير إلى الإنجليزية/i }))
    expect(onSwitchLocale).toHaveBeenCalledOnce()
    expect(onClose).toHaveBeenCalledTimes(2)
  })

  it('places authenticated account actions inside the mobile drawer', async () => {
    const onClose = vi.fn()
    const profileSession: LandingSessionState = {
      status: 'authenticated',
      profile: { id: 7, name: 'Sample Manager', email: 'manager@example.test', restaurantName: 'Sample Cafe' },
    }
    render(<MobileNavigationHarness onClose={onClose} onToggleTheme={vi.fn()} onSwitchLocale={vi.fn()} session={profileSession} />)

    expect(screen.getByText('Sample Manager')).toHaveAttribute('dir', 'auto')
    expect(screen.getByRole('link', { name: 'لوحة التحكم' })).toHaveAttribute('href', 'http://localhost:3002/ar/dashboard')
    expect(screen.getByRole('link', { name: 'إعدادات الحساب' })).toHaveAttribute('href', 'http://localhost:3002/ar/dashboard/settings/account')
    expect(screen.queryByRole('link', { name: 'تسجيل الدخول' })).not.toBeInTheDocument()
    expect(screen.queryByRole('link', { name: 'ابدأ تجربتك المجانية لمدة 14 يومًا' })).not.toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: 'تسجيل الخروج' }))
    await waitFor(() => expect(onClose).toHaveBeenCalledOnce())
  })

  it('provides complete localized footer route groups without invented legal routes', () => {
    render(<Footer locale="ar" />)
    const links = screen.getAllByRole('link')
    expect(links.map(link => link.getAttribute('href'))).toContain('/ar/features/website-builder')
    expect(links.map(link => link.getAttribute('href'))).toContain('/ar/contact')
    expect(links.map(link => link.getAttribute('href'))).toContain('/ar/privacy')
    expect(links.map(link => link.getAttribute('href'))).not.toContain('/ar/terms')
    expect(screen.getByRole('heading', { name: 'المنتج' })).toBeInTheDocument()
  })
})
