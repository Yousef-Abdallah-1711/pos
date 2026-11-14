import { cleanup, render, screen, waitFor } from '@testing-library/react'
import '@testing-library/jest-dom/vitest'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { SuperAdminAuthProvider } from './SuperAdminAuthProvider'
import { useAuthStore } from '@/stores/authStore'

const { getSession, push } = vi.hoisted(() => ({ getSession: vi.fn(), push: vi.fn() }))
const router = { push }
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
})
