import { cleanup, render, screen, waitFor } from '@testing-library/react'
import '@testing-library/jest-dom/vitest'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { SuperAdminAuthProvider } from './SuperAdminAuthProvider'
import { useAuthStore } from '@/stores/authStore'

const { getSession, push } = vi.hoisted(() => ({ getSession: vi.fn(), push: vi.fn() }))
const router = { push }
const syntheticAdmin = {
  id: 1,
  user_id: 1,
  email: 'synthetic@example.invalid',
  ip_whitelist: [],
  is_active: true,
  failed_login_attempts: 0,
  created_at: '2026-01-01T00:00:00.000Z',
  updated_at: '2026-01-01T00:00:00.000Z',
}
vi.mock('next/navigation', () => ({ usePathname: () => '/dashboard', useRouter: () => router }))
vi.mock('@/lib/api', () => ({ superAdminAPI: { getSession } }))
afterEach(cleanup)
beforeEach(() => {
  vi.clearAllMocks()
  useAuthStore.setState({ token: null, superAdmin: null, isLoading: true })
})

describe('administrator session bootstrap', () => {
  it('does not request the session again when bootstrap sets its own authentication state', async () => {
    getSession.mockResolvedValue({ super_admin: { id: 1 }, email: 'synthetic@example.invalid', is_super_admin: true })
    render(<SuperAdminAuthProvider><p>Protected content</p></SuperAdminAuthProvider>)
    expect(await screen.findByText('Protected content')).toBeVisible()
    await waitFor(() => expect(useAuthStore.getState().token).toBe('cookie-session'))
    expect(getSession).toHaveBeenCalledOnce()
    expect(push).not.toHaveBeenCalled()
  })

  it('preserves the current identity and offers retry after a temporary session outage', async () => {
    useAuthStore.setState({ token: 'cookie-session', superAdmin: syntheticAdmin, isLoading: true })
    getSession.mockRejectedValueOnce({ response: { status: 503 } })
      .mockResolvedValueOnce({ super_admin: { id: 1 }, email: 'synthetic@example.invalid', is_super_admin: true })

    render(<SuperAdminAuthProvider><p>Protected content</p></SuperAdminAuthProvider>)

    const alert = await screen.findByRole('alert')
    expect(alert).toHaveTextContent(/temporarily unavailable/i)
    expect(useAuthStore.getState().token).toBe('cookie-session')
    expect(push).not.toHaveBeenCalled()
    expect(screen.queryByText('Protected content')).not.toBeInTheDocument()

    await screen.getByRole('button', { name: /retry/i }).click()
    expect(await screen.findByText('Protected content')).toBeVisible()
    expect(useAuthStore.getState().token).toBe('cookie-session')
    expect(getSession).toHaveBeenCalledTimes(2)
    expect(push).not.toHaveBeenCalled()
  })

  it('clears an invalid identity and redirects to login only on 401', async () => {
    useAuthStore.setState({ token: 'cookie-session', superAdmin: syntheticAdmin, isLoading: true })
    getSession.mockRejectedValue({ response: { status: 401 } })

    render(<SuperAdminAuthProvider><p>Protected content</p></SuperAdminAuthProvider>)

    await waitFor(() => expect(push).toHaveBeenCalledWith('/login'))
    expect(useAuthStore.getState().token).toBeNull()
  })
})
