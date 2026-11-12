import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import DashboardError from './error'
import DashboardLoading from './loading'

const { pathname } = vi.hoisted(() => ({ pathname: { value: '/ar/dashboard' } }))

vi.mock('next/navigation', () => ({
  usePathname: () => pathname.value,
}))

describe('dashboard route states', () => {
  beforeEach(() => {
    pathname.value = '/ar/dashboard'
  })

  it('shows an Arabic, direction-aware loading skeleton', () => {
    render(<DashboardLoading />)

    expect(screen.getByRole('main')).toHaveAttribute('dir', 'rtl')
    expect(screen.getByRole('main')).toHaveAttribute('aria-busy', 'true')
    expect(screen.getByRole('status')).toHaveTextContent('جاري التحميل...')
  })

  it('offers a localized retry action when a dashboard route fails', () => {
    const reset = vi.fn()
    render(<DashboardError error={new Error('route failed')} reset={reset} />)

    expect(screen.getByRole('main')).toHaveAttribute('dir', 'rtl')
    expect(screen.getByRole('alert')).toBeVisible()
    expect(screen.getByRole('heading', { name: 'تعذر تحميل هذا القسم.' })).toBeVisible()
    fireEvent.click(screen.getByRole('button', { name: 'حاول مرة أخرى' }))
    expect(reset).toHaveBeenCalledOnce()
  })
})
