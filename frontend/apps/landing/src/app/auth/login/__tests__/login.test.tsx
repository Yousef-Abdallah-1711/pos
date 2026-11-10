import React from 'react'
import { vi, type Mock } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import LoginPage from '@/app/[locale]/auth/login/page'
import * as api from '@/lib/api'
import { toast } from 'sonner'

describe('LoginPage - Multi-Tenant Authentication Flow', () => {
  const mockPush = vi.fn()
  const mockSetAuth = vi.fn()
  const mockApiPost = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useRouter as Mock).mockReturnValue({
      push: mockPush,
    })
    ;(useAuthStore as unknown as Mock).mockImplementation(() => mockSetAuth)
    ;(api.default.post as Mock) = mockApiPost
  })

  describe('Single Tenant Login', () => {
    it('should auto-login and redirect when user has only one tenant', async () => {
      const mockResponse = {
        status: 200,
        data: {
          user: {
            id: 1,
            email: 'user@test.com',
            name: 'Test User',
            role: 'owner',
            tenant_id: 1,
          },
          token: 'mock-jwt-token',
          multiple_tenants: false,
        },
      }

      mockApiPost.mockResolvedValueOnce(mockResponse)

      render(<LoginPage />)

      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(emailInput, 'user@test.com')
      await userEvent.type(passwordInput, 'password123')
      fireEvent.click(submitButton)

      await waitFor(() => {
        expect(mockApiPost).toHaveBeenCalledWith('/auth/login', {
          email: 'user@test.com',
          password: 'password123',
        })
      })

      await waitFor(() => {
        expect(mockSetAuth).toHaveBeenCalledWith(mockResponse.data.user, '')
      })
    })
  })

  describe('Multi-Tenant Login', () => {
    it('should show tenant selector modal when user has multiple tenants', async () => {
      const mockResponse = {
        status: 200,
        data: {
          success: true,
          multiple_tenants: true,
          tenants: [
            {
              tenant_id: 1,
              tenant_name: 'Company A',
              user_id: 1,
              roles: ['owner'],
            },
            {
              tenant_id: 2,
              tenant_name: 'Company B',
              user_id: 2,
              roles: ['manager'],
            },
          ],
          message: 'Multiple organizations found. Please select one.',
        },
      }

      mockApiPost.mockResolvedValueOnce(mockResponse)

      render(<LoginPage />)

      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(emailInput, 'multi@test.com')
      await userEvent.type(passwordInput, 'password123')
      fireEvent.click(submitButton)

      await waitFor(() => {
        expect(screen.getByText('Select Organization')).toBeInTheDocument()
      })

      // Check that both tenants are displayed
      await waitFor(() => {
        expect(screen.getByText('Company A')).toBeInTheDocument()
        expect(screen.getByText('Company B')).toBeInTheDocument()
      })
    })

    it('should display role badges for each tenant', async () => {
      const mockResponse = {
        status: 200,
        data: {
          success: true,
          multiple_tenants: true,
          tenants: [
            {
              tenant_id: 1,
              tenant_name: 'Company A',
              user_id: 1,
              roles: ['owner', 'admin'],
            },
            {
              tenant_id: 2,
              tenant_name: 'Company B',
              user_id: 2,
              roles: ['manager'],
            },
          ],
        },
      }

      mockApiPost.mockResolvedValueOnce(mockResponse)

      render(<LoginPage />)

      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(emailInput, 'multi@test.com')
      await userEvent.type(passwordInput, 'password123')
      fireEvent.click(submitButton)

      await waitFor(() => {
        expect(screen.getByText('owner')).toBeInTheDocument()
        expect(screen.getByText('admin')).toBeInTheDocument()
      })
    })
  })

  describe('Error Handling', () => {
    it('should display error message for invalid credentials', async () => {
      const mockError = {
        response: {
          status: 401,
          data: {
            error: 'Invalid email or password',
          },
        },
      }

      mockApiPost.mockRejectedValueOnce(mockError)

      render(<LoginPage />)

      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(emailInput, 'wrong@test.com')
      await userEvent.type(passwordInput, 'wrongpass')
      fireEvent.click(submitButton)

      await waitFor(() => {
        expect(toast.error).toHaveBeenCalledWith('Invalid email or password')
      })
    })

    it('should display error for empty email field', async () => {
      render(<LoginPage />)

      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(passwordInput, 'password123')
      fireEvent.click(submitButton)

      await waitFor(() => {
        expect(screen.getByText('Invalid email address')).toBeInTheDocument()
      })
    })

    it('should display error for empty password field', async () => {
      render(<LoginPage />)

      const emailInput = screen.getByPlaceholderText('your@email.com')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(emailInput, 'user@test.com')
      fireEvent.click(submitButton)

      await waitFor(() => {
        expect(screen.getByText('Password must be at least 6 characters')).toBeInTheDocument()
      })
    })
  })

  describe('Form Validation', () => {
    it('should show/hide password when toggle is clicked', async () => {
      render(<LoginPage />)

      const passwordInput = screen.getByPlaceholderText('Enter your password') as HTMLInputElement
      const toggleButton = screen.getByRole('button', { name: 'Show password' })

      expect(passwordInput.type).toBe('password')

      fireEvent.click(toggleButton)

      await waitFor(() => {
        expect(passwordInput.type).toBe('text')
      })

      fireEvent.click(screen.getByRole('button', { name: 'Hide password' }))

      await waitFor(() => {
        expect(passwordInput.type).toBe('password')
      })
    })

    it('should disable submit button while loading', async () => {
      mockApiPost.mockImplementation(
        () => new Promise(resolve => setTimeout(resolve, 1000))
      )

      render(<LoginPage />)

      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i }) as HTMLButtonElement

      await userEvent.type(emailInput, 'user@test.com')
      await userEvent.type(passwordInput, 'password123')

      fireEvent.click(submitButton)

      await waitFor(() => expect(submitButton).toBeDisabled())
    })
  })

  describe('Navigation', () => {
    it('should have link to register page', () => {
      render(<LoginPage />)

      const registerLink = screen.getByRole('link', { name: /create account/i })
      expect(registerLink).toHaveAttribute('href', '/en/auth/register')
    })

    it('does not advertise password reset until the end-to-end flow is supported', () => {
      render(<LoginPage />)

      expect(screen.queryByRole('link', { name: /forgot password/i })).not.toBeInTheDocument()
    })
  })
})
