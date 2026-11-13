import { cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import '@testing-library/jest-dom/vitest'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { DashboardNav } from './DashboardNav'
import { useAuthStore } from '@/stores/authStore'

const { logout, replace } = vi.hoisted(() => ({ logout: vi.fn(), replace: vi.fn() }))
vi.mock('next/navigation', () => ({ usePathname: () => '/dashboard', useRouter: () => ({ replace }) }))
vi.mock('@/lib/api', () => ({ superAdminAPI: { logout } }))

afterEach(cleanup)
beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
  useAuthStore.setState({ token: 'cookie-session' })
})

describe('administrator logout', () => {
  it('revokes the server session before clearing local authentication', async () => {
    logout.mockResolvedValue(undefined)
    render(<DashboardNav />)
    await userEvent.click(screen.getByRole('button', { name: 'Logout' }))
    await waitFor(() => expect(replace).toHaveBeenCalledWith('/login'))
    expect(logout).toHaveBeenCalledOnce()
    expect(useAuthStore.getState().token).toBeNull()
    expect(localStorage.getItem('super-admin-auth-store')).toBeNull()
  })

  it('keeps authentication and offers retry when revocation fails', async () => {
    logout.mockRejectedValue(new Error('service unavailable'))
    render(<DashboardNav />)
    await userEvent.click(screen.getByRole('button', { name: 'Logout' }))
    expect(await screen.findByRole('alert')).toHaveTextContent('Unable to sign out')
    expect(useAuthStore.getState().token).toBe('cookie-session')
    expect(replace).not.toHaveBeenCalled()
    expect(screen.getByRole('button', { name: 'Logout' })).toBeEnabled()
  })
})
