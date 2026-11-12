import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import LandingError from './error'
import LandingLoading from './loading'

const { pathname } = vi.hoisted(() => ({ pathname: { value: '/ar' } }))

vi.mock('next/navigation', () => ({
  usePathname: () => pathname.value,
}))

describe('landing route states', () => {
  beforeEach(() => {
    pathname.value = '/ar/auth/login'
  })

  it('shows a localized loading skeleton with the correct direction', () => {
    render(<LandingLoading />)

    expect(screen.getByRole('main')).toHaveAttribute('dir', 'rtl')
    expect(screen.getByRole('main')).toHaveAttribute('aria-busy', 'true')
    expect(screen.getByRole('status')).toHaveTextContent('جاري التحميل...')
  })

  it('offers localized recovery after a route error', () => {
    const reset = vi.fn()
    render(<LandingError error={new Error('route failure')} reset={reset} />)

    expect(screen.getByRole('main')).toHaveAttribute('dir', 'rtl')
    expect(screen.getByRole('heading', { name: 'تعذر تحميل هذه الصفحة.' })).toBeVisible()
    expect(screen.getByRole('link', { name: 'العودة للرئيسية' })).toHaveAttribute('href', '/ar')
    fireEvent.click(screen.getByRole('button', { name: 'حاول مرة أخرى' }))
    expect(reset).toHaveBeenCalledOnce()
  })
})
