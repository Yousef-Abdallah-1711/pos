import { afterEach, beforeEach, expect, it, vi } from 'vitest'
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { usePathname, useRouter } from 'next/navigation'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/authStore'
import { BlockedOverlay } from '../BlockedOverlay'
import { TrialExpiredScreen } from '../TrialExpiredScreen'

vi.mock('next/navigation', () => ({ usePathname: vi.fn(), useRouter: vi.fn() }))
vi.mock('@/lib/api', () => ({ default: { post: vi.fn() } }))
vi.mock('@/stores/authStore', () => ({ useAuthStore: vi.fn() }))
vi.mock('@/lib/authSession', () => ({ getDashboardLoginUrl: (pathname: string) => pathname.startsWith('/ar') ? '/ar/auth/login' : '/en/auth/login' }))

const logout = vi.fn()
const replace = vi.fn()
beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(usePathname).mockReturnValue('/en/dashboard')
  vi.mocked(useRouter).mockReturnValue({ replace } as ReturnType<typeof useRouter>)
  vi.mocked(useAuthStore).mockReturnValue(logout)
})
afterEach(cleanup)

it('routes an expired trial to billing without repricing historical entitlements', () => {
  render(<BlockedOverlay blockReason="trial_expired" />)
  expect(screen.getByRole('heading', { name: 'Your trial has expired' })).toBeInTheDocument()
  expect(screen.getByRole('link', { name: 'Review billing options' })).toHaveAttribute('href', '/en/dashboard/settings/billing')
  expect(document.body).not.toHaveTextContent(/500|15-day|support@pos-saas\.com/)
})

it('localizes expired access and billing in Arabic', () => {
  vi.mocked(usePathname).mockReturnValue('/ar/dashboard')
  render(<BlockedOverlay blockReason="expired" />)
  expect(screen.getByRole('heading', { name: 'انتهت صلاحية الاشتراك' })).toBeInTheDocument()
  expect(screen.getByRole('link', { name: 'مراجعة خيارات الفوترة' })).toHaveAttribute('href', '/ar/dashboard/settings/billing')
  expect(screen.getByRole('button', { name: 'تسجيل الخروج' })).toBeEnabled()
})

it('does not advertise purchasing or an invented mailbox for a disabled account', () => {
  render(<BlockedOverlay blockReason="tenant_disabled" />)
  expect(screen.queryByRole('link')).not.toBeInTheDocument()
  expect(screen.getByText(/Contact your account administrator/)).toBeInTheDocument()
})

it('keeps the legacy expiry entry point on the same accurate flow', () => {
  render(<TrialExpiredScreen trialEndDate="2026-01-01" />)
  expect(screen.getByRole('link', { name: 'Review billing options' })).toBeInTheDocument()
  expect(document.body).not.toHaveTextContent(/15-day|500 EGP/)
})

it('retains the session and offers retry when server revocation fails', async () => {
  vi.mocked(api.post).mockRejectedValueOnce(new Error('network failure'))
  render(<BlockedOverlay blockReason="trial_expired" />)
  fireEvent.click(screen.getByRole('button', { name: 'Sign out' }))
  expect(await screen.findByRole('alert')).toHaveTextContent('Sign out failed')
  expect(logout).not.toHaveBeenCalled()
  expect(replace).not.toHaveBeenCalled()
  expect(screen.getByRole('button', { name: 'Sign out' })).toBeEnabled()
})

it('waits for server revocation before clearing state and navigating', async () => {
  let revoke: (() => void) | undefined
  vi.mocked(api.post).mockImplementationOnce(() => new Promise(resolve => { revoke = () => resolve({ data: {} }) }))
  render(<BlockedOverlay blockReason="trial_expired" />)
  fireEvent.click(screen.getByRole('button', { name: 'Sign out' }))
  expect(api.post).toHaveBeenCalledWith('/auth/logout')
  expect(screen.getByRole('button', { name: 'Signing out…' })).toBeDisabled()
  expect(logout).not.toHaveBeenCalled()
  revoke?.()
  await waitFor(() => expect(replace).toHaveBeenCalledWith('/en/auth/login'))
  expect(logout).toHaveBeenCalledOnce()
})
