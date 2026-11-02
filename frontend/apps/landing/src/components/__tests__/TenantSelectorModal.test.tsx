import React from 'react'
import { vi, type Mock } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { TenantSelectorModal } from '../TenantSelectorModal'
import * as sonner from 'sonner'

// Mock fetch globally
global.fetch = vi.fn()

describe('TenantSelectorModal', () => {
  const mockPush = vi.fn()
  const mockSetAuth = vi.fn()
  const mockToast = {
    success: vi.fn(),
    error: vi.fn(),
  }

  const mockTenants = [
    {
      tenant_id: 1,
      tenant_name: 'Company A',
      user_id: 101,
      roles: ['owner'],
    },
    {
      tenant_id: 2,
      tenant_name: 'Company B',
      user_id: 202,
      roles: ['manager', 'admin'],
    },
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useRouter as Mock).mockReturnValue({
      push: mockPush,
    })
    ;(useAuthStore as unknown as Mock).mockImplementation(() => mockSetAuth)
    sonner.toast.success = mockToast.success
    sonner.toast.error = mockToast.error
    ;(global.fetch as Mock).mockReset()
  })

  describe('Rendering', () => {
    it('should not render when isOpen is false', () => {
      const { container } = render(
        <TenantSelectorModal
          isOpen={false}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
        />
      )

      expect(container.firstChild).toBeNull()
    })

    it('should not render when tenants list is empty', () => {
      const { container } = render(
        <TenantSelectorModal isOpen={true} email="user@test.com" selectionToken="proof-token" tenants={[]} />
      )

      expect(container.firstChild).toBeNull()
    })

    it('should render modal when isOpen and tenants exist', () => {
      render(
        <TenantSelectorModal
          isOpen={true}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
        />
      )

      expect(screen.getByText('Select Organization')).toBeInTheDocument()
      expect(screen.getByText('Company A')).toBeInTheDocument()
      expect(screen.getByText('Company B')).toBeInTheDocument()
    })

    it('should display user email in modal', () => {
      render(
        <TenantSelectorModal selectionToken="proof-token"
          isOpen={true}
          email="testuser@example.com"
          tenants={mockTenants}
        />
      )

      expect(screen.getByText(/testuser@example.com/)).toBeInTheDocument()
    })

    it('should display role badges for each tenant', () => {
      render(
        <TenantSelectorModal
          isOpen={true}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
        />
      )

      expect(screen.getByText('owner')).toBeInTheDocument()
      expect(screen.getByText('manager')).toBeInTheDocument()
      expect(screen.getByText('admin')).toBeInTheDocument()
    })
  })

  describe('Tenant Selection', () => {
    it('should handle tenant selection successfully', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          user: { id: 101, email: 'user@test.com' },
          token: 'mock-token',
        }),
      }

      ;(global.fetch as Mock).mockResolvedValueOnce(mockResponse)

      render(
        <TenantSelectorModal
          isOpen={true}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
        />
      )

      const companyAButton = screen.getByRole('button', { name: /company a/i })
      fireEvent.click(companyAButton)

      await waitFor(() => {
        expect(global.fetch).toHaveBeenCalledWith(
          expect.stringContaining('/auth/login/confirm'),
          expect.objectContaining({
            method: 'POST',
            body: JSON.stringify({
              email: 'user@test.com',
              tenant_id: 1,
              selection_token: 'proof-token',
            }),
          })
        )
      })

      await waitFor(() => {
        expect(mockSetAuth).toHaveBeenCalled()
        expect(mockToast.success).toHaveBeenCalledWith('Login successful!')
      })
    })

    it('should display error message on failed tenant selection', async () => {
      const mockError = {
        ok: false,
        json: vi.fn().mockResolvedValue({
          error: 'User not found in selected organization',
        }),
      }

      ;(global.fetch as Mock).mockResolvedValueOnce(mockError)

      render(
        <TenantSelectorModal
          isOpen={true}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
        />
      )

      const companyAButton = screen.getByRole('button', { name: /company a/i })
      fireEvent.click(companyAButton)

      await waitFor(() => {
        expect(screen.getByText('User not found in selected organization')).toBeInTheDocument()
      })

      expect(mockToast.error).toHaveBeenCalled()
    })

    it('should disable buttons while loading', async () => {
      ;(global.fetch as Mock).mockImplementation(
        () => new Promise(resolve => setTimeout(resolve, 1000))
      )

      render(
        <TenantSelectorModal
          isOpen={true}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
        />
      )

      const buttons = screen.getAllByRole('button')
      const selectButton = buttons[0]

      fireEvent.click(selectButton)

      // Check that other buttons are disabled during loading
      buttons.forEach(btn => {
        if (btn !== selectButton) {
          expect(btn).toBeDisabled()
        }
      })
    })
  })

  describe('Visual Feedback', () => {
    it('should highlight selected tenant', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          user: { id: 101 },
          token: 'mock-token',
        }),
      }

      ;(global.fetch as Mock).mockImplementation(
        () => new Promise(resolve => setTimeout(resolve, 100))
      )

      render(
        <TenantSelectorModal
          isOpen={true}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
        />
      )

      const companyAButton = screen.getByRole('button', { name: /company a/i })
      fireEvent.click(companyAButton)

      // Button should have visual indication of being selected
      await waitFor(() => {
        expect(companyAButton).toHaveClass('border-blue-500')
      })
    })

    it('should show loading spinner when selecting tenant', async () => {
      ;(global.fetch as Mock).mockImplementation(
        () => new Promise(resolve => setTimeout(resolve, 500))
      )

      render(
        <TenantSelectorModal
          isOpen={true}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
        />
      )

      const companyAButton = screen.getAllByRole('button')[0]
      fireEvent.click(companyAButton)

      // Check for spinner element
      await waitFor(() => {
        const spinner = document.querySelector('.animate-spin')
        expect(spinner).toBeInTheDocument()
      })
    })
  })

  describe('Accessibility', () => {
    it('should have proper heading for modal', () => {
      render(
        <TenantSelectorModal
          isOpen={true}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
        />
      )

      expect(screen.getByRole('heading', { name: /select organization/i })).toBeInTheDocument()
    })

    it('should have descriptive text for user guidance', () => {
      render(
        <TenantSelectorModal
          isOpen={true}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
        />
      )

      expect(
        screen.getByText(/you have accounts in multiple organizations/i)
      ).toBeInTheDocument()
    })
  })

  describe('Callback', () => {
    it('should call onTenantSelected callback after successful login', async () => {
      const mockCallback = vi.fn()
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          user: { id: 101 },
          token: 'mock-token',
        }),
      }

      ;(global.fetch as Mock).mockResolvedValueOnce(mockResponse)

      render(
        <TenantSelectorModal
          isOpen={true}
          email="user@test.com"
          selectionToken="proof-token"
          tenants={mockTenants}
          onTenantSelected={mockCallback}
        />
      )

      const companyAButton = screen.getByRole('button', { name: /company a/i })
      fireEvent.click(companyAButton)

      await waitFor(() => {
        expect(mockCallback).toHaveBeenCalledWith(1)
      })
    })
  })
})
