import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import '@testing-library/jest-dom/vitest'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import LoginPage from './page'
import { useAuthStore } from '@/stores/authStore'

const { login } = vi.hoisted(() => ({ login: vi.fn() }))
vi.mock('@/lib/api', () => ({ superAdminAPI: { login } }))
vi.mock('sonner', () => ({ toast: { error: vi.fn(), success: vi.fn() } }))
afterEach(cleanup)
beforeEach(() => {
  vi.clearAllMocks()
  useAuthStore.setState({ token: null, superAdmin: null, isLoading: false })
})

describe('super-admin login validation', () => {
  it('keeps the form mounted and displays the server rejection', async () => {
    let rejectLogin: (reason: unknown) => void = () => undefined
    login.mockReturnValue(new Promise((_resolve, reject) => { rejectLogin = reject }))
    render(<LoginPage />)
    fireEvent.change(screen.getByLabelText('Email Address'), { target: { value: 'invalid@example.invalid' } })
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'Synthetic-invalid-only-42!' } })
    fireEvent.click(screen.getByRole('button', { name: 'Login' }))
    expect(useAuthStore.getState().isLoading).toBe(false)
    expect(screen.getByRole('button', { name: 'Logging in...' })).toBeDisabled()
    await act(async () => { rejectLogin({ isAxiosError: true, response: { status: 401, data: { error: 'Invalid credentials' } } }) })
    expect(await screen.findByText('Invalid credentials')).toBeVisible()
    expect(screen.getByRole('button', { name: 'Login' })).toBeEnabled()
    expect(useAuthStore.getState().token).toBeNull()
  })

  it('rejects empty fields without requesting an administrator session', async () => {
    render(<LoginPage />)
    fireEvent.change(screen.getByLabelText('Email Address'), { target: { value: '' } })
    fireEvent.click(screen.getByRole('button', { name: 'Login' }))
    expect(await screen.findByText('Please enter both email and password')).toBeVisible()
    expect(login).not.toHaveBeenCalled()
    expect(useAuthStore.getState().isLoading).toBe(false)
  })
})
