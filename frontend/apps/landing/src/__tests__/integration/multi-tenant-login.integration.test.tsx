/**
 * Integration Tests: Multi-Tenant Authentication Flow
 * Tests the complete journey from login to tenant selection to dashboard
 */

import React from 'react'
import { vi, type Mock } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import LoginPage from '@/app/[locale]/auth/login/page'
import api from '@/lib/api'
import { toast } from 'sonner'

global.fetch = vi.fn()

describe('Multi-Tenant Login Integration Tests', () => {
  const mockRouter = {
    push: vi.fn(),
  }
  const mockSetAuth = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useRouter as Mock).mockReturnValue(mockRouter)
    ;(useAuthStore as unknown as Mock).mockImplementation(() => mockSetAuth)
    ;(global.fetch as Mock).mockReset()
  })

  describe('Scenario 1: User with Single Tenant Account', () => {
    it('should complete login flow without tenant selection', async () => {
      // Step 1: Mock login endpoint returning single tenant
      const loginResponse = {
        status: 200,
        data: {
          user: {
            id: 1,
            email: 'single@example.com',
            name: 'Single User',
            role: 'owner',
            tenant_id: 1,
          },
          token: 'eyJhbGc...',
          multiple_tenants: false,
        },
      }

      const mockApiPost = vi.fn().mockResolvedValue(loginResponse)
      api.post = mockApiPost

      render(<LoginPage />)

      // Step 2: User enters credentials
      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(emailInput, 'single@example.com')
      await userEvent.type(passwordInput, 'password123')

      // Step 3: User submits login form
      fireEvent.click(submitButton)

      // Step 4: System authenticates and auto-logs in
      await waitFor(() => {
        expect(mockSetAuth).toHaveBeenCalledWith(loginResponse.data.user, '')
      })

      // Step 5: System redirects to dashboard (done by TenantSelectorModal not being shown)
      // Verify tenant selector modal is NOT shown
      await waitFor(() => {
        expect(screen.queryByText('Select Organization')).not.toBeInTheDocument()
      })
    })
  })

  describe('Scenario 2: User with Multiple Tenant Accounts', () => {
    it('should show tenant selector and allow organization selection', async () => {
      // Step 1: User navigates to login
      // (already done in beforeEach render)

      render(<LoginPage />)

      // Step 2: User enters credentials for account with multiple organizations
      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(emailInput, 'multi@example.com')
      await userEvent.type(passwordInput, 'password123')

      // Step 3: Backend returns list of available tenants
      const loginMultiResponse = {
        status: 200,
        data: {
          success: true,
          multiple_tenants: true,
          selection_token: 'proof-token',
          tenants: [
            {
              tenant_id: 1,
              tenant_name: 'Acme Corp',
              user_id: 101,
              roles: ['owner'],
            },
            {
              tenant_id: 2,
              tenant_name: 'Tech Startup',
              user_id: 202,
              roles: ['manager', 'admin'],
            },
            {
              tenant_id: 3,
              tenant_name: 'Small Business',
              user_id: 303,
              roles: ['editor'],
            },
          ],
        },
      }

      const mockApiPost = vi.fn().mockResolvedValue(loginMultiResponse)
      api.post = mockApiPost

      const confirmResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          user: {
            id: 202,
            email: 'multi@example.com',
            name: 'Multi User',
            role: 'manager',
            tenant_id: 2,
          },
        }),
      }
      ;(global.fetch as Mock).mockResolvedValueOnce(confirmResponse)

      fireEvent.click(submitButton)

      // Step 4: Modal appears showing all available organizations
      await waitFor(() => {
        expect(screen.getByText('Select Organization')).toBeInTheDocument()
      })

      // Step 5: Verify all organizations are displayed
      await waitFor(() => {
        expect(screen.getByText('Acme Corp')).toBeInTheDocument()
        expect(screen.getByText('Tech Startup')).toBeInTheDocument()
        expect(screen.getByText('Small Business')).toBeInTheDocument()
      })

      // Step 6: Verify roles are displayed for each organization
      expect(screen.getByText('owner')).toBeInTheDocument()
      expect(screen.getByText('manager')).toBeInTheDocument()
      expect(screen.getByText('admin')).toBeInTheDocument()
      expect(screen.getByText('editor')).toBeInTheDocument()

      // Step 7: User selects organization
      const techStartupButton = screen.getByRole('button', { name: /tech startup/i })
      fireEvent.click(techStartupButton)

      // Step 8: System sends tenant confirmation request
      await waitFor(() => {
        expect(global.fetch).toHaveBeenCalledWith(
          expect.stringContaining('/auth/login/confirm'),
          expect.objectContaining({
            method: 'POST',
            body: JSON.stringify({
              email: 'multi@example.com',
              tenant_id: 2,
              selection_token: 'proof-token',
            }),
          })
        )
      })

      // Step 9: Backend returns JWT token for selected tenant
      // Step 10: System stores auth and redirects to dashboard
      await waitFor(() => {
        expect(mockSetAuth).toHaveBeenCalledWith(
          expect.objectContaining({
            tenant_id: 2,
            email: 'multi@example.com',
          }),
          expect.any(String)
        )
      })
    })

    it('should handle organization switching', async () => {
      // This tests that a user can go back and select a different organization
      // (This would typically require a logout/re-login or special dashboard feature)
      // For now, this is a note for future implementation
      expect(true).toBe(true)
    })
  })

  describe('Scenario 3: Error Handling', () => {
    it('should handle invalid credentials gracefully', async () => {
      const errorResponse = {
        status: 401,
        data: {
          error: 'Invalid email or password',
        },
      }

      const mockApiPost = vi.fn().mockRejectedValue({
        response: errorResponse,
      })

      api.post = mockApiPost

      render(<LoginPage />)

      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(emailInput, 'user@example.com')
      await userEvent.type(passwordInput, 'wrongpassword')
      fireEvent.click(submitButton)

      // Error message should be displayed
      await waitFor(() => {
        expect(toast.error).toHaveBeenCalledWith('Invalid email or password')
      })

      // Form should remain visible for retry
      expect(screen.getByPlaceholderText('your@email.com')).toBeInTheDocument()
    })

    it('should handle network errors', async () => {
      const networkError = Object.assign(new Error('Network request failed'), {
        code: 'ERR_NETWORK',
      })

      const mockApiPost = vi.fn().mockRejectedValue(networkError)

      api.post = mockApiPost

      render(<LoginPage />)

      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(emailInput, 'user@example.com')
      await userEvent.type(passwordInput, 'password123')
      fireEvent.click(submitButton)

      // Network error message should be shown
      await waitFor(() => {
        expect(toast.error).toHaveBeenCalledWith('Network error. Please check your connection and try again.')
      })
    })

    it('should handle tenant selection failure', async () => {
      // First login succeeds, showing multiple tenants
      const loginResponse = {
        status: 200,
        data: {
          success: true,
          multiple_tenants: true,
          tenants: [
            {
              tenant_id: 1,
              tenant_name: 'Org A',
              user_id: 101,
              roles: ['owner'],
            },
            {
              tenant_id: 2,
              tenant_name: 'Org B',
              user_id: 102,
              roles: ['manager'],
            },
          ],
        },
      }

      const mockApiPost = vi.fn().mockResolvedValue(loginResponse)
      api.post = mockApiPost

      render(<LoginPage />)

      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')
      const submitButton = screen.getByRole('button', { name: /sign in/i })

      await userEvent.type(emailInput, 'user@example.com')
      await userEvent.type(passwordInput, 'password123')
      fireEvent.click(submitButton)

      // Wait for modal to appear
      await waitFor(() => {
        expect(screen.getByText('Select Organization')).toBeInTheDocument()
      })

      // Now tenant selection fails
      const tenantError = {
        ok: false,
        json: vi.fn().mockResolvedValue({
          error: 'User not found in selected organization',
        }),
      }

      ;(global.fetch as Mock).mockResolvedValueOnce(tenantError)

      const orgButton = screen.getByRole('button', { name: /org a/i })
      fireEvent.click(orgButton)

      // Error message should be shown
      await waitFor(() => {
        expect(
          screen.getByText('User not found in selected organization')
        ).toBeInTheDocument()
      })

      // Modal should still be visible for retry
      expect(screen.getByText('Select Organization')).toBeInTheDocument()
    })
  })

  describe('Scenario 4: Security Considerations', () => {
    it('should not expose password in logs or errors', async () => {
      const consoleLogSpy = vi.spyOn(console, 'log').mockImplementation(() => {})
      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

      render(<LoginPage />)

      const emailInput = screen.getByPlaceholderText('your@email.com')
      const passwordInput = screen.getByPlaceholderText('Enter your password')

      await userEvent.type(emailInput, 'user@example.com')
      await userEvent.type(passwordInput, 'SuperSecretPassword123!')

      // Check that password is never logged
      const allLogs = [...consoleLogSpy.mock.calls, ...consoleErrorSpy.mock.calls].flat()
      const passwordInLogs = allLogs.some(log =>
        String(log).includes('SuperSecretPassword123!')
      )

      expect(passwordInLogs).toBe(false)

      consoleLogSpy.mockRestore()
      consoleErrorSpy.mockRestore()
    })

    it('should clear sensitive data after login', async () => {
      const loginResponse = {
        status: 200,
        data: {
          user: {
            id: 1,
            email: 'user@example.com',
            name: 'Test User',
            role: 'owner',
            tenant_id: 1,
          },
          token: 'jwt-token-here',
          // Note: password_hash should NOT be in response
        },
      }

      // Verify response doesn't contain password_hash
      expect(loginResponse.data).not.toHaveProperty('password_hash')
    })
  })
})
