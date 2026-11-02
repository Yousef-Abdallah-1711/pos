/**
 * Enhanced Error Boundary Tests
 * Tests error catching, recovery suggestions, logging, and UI functionality
 */

import React, { ReactNode } from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { EnhancedErrorBoundary } from '../enhanced-error-boundary'
import { getErrorLogger } from '@/lib/error-handling/error-logger'
import { getRecoverySuggestions, getErrorAdvice } from '@/lib/error-handling/error-recovery-suggestions'

// Mock the error logger
vi.mock('@/lib/error-handling/error-logger', () => ({
  getErrorLogger: vi.fn(() => ({
    logComponentError: vi.fn(),
    log: vi.fn(),
    getAllLogs: vi.fn(() => []),
  })),
}))

// Mock the recovery suggestions
vi.mock('@/lib/error-handling/error-recovery-suggestions', () => ({
  getRecoverySuggestions: vi.fn(() => [
    {
      title: 'Try Again',
      description: 'Attempt to reload the component',
      action: { label: 'Retry', onClick: vi.fn() },
      priority: 1,
    },
    {
      title: 'Go to Dashboard',
      description: 'Return to the main dashboard',
      action: { label: 'Dashboard', onClick: vi.fn() },
      priority: 2,
    },
  ]),
  getErrorAdvice: vi.fn(() => 'An unexpected error occurred. Please try again or contact support.'),
}))

// Component that throws an error
interface ErrorThrowingComponentProps {
  shouldThrow?: boolean
  errorMessage?: string
}

const ErrorThrowingComponent: React.FC<ErrorThrowingComponentProps> = ({
  shouldThrow = false,
  errorMessage = 'Test error message',
}) => {
  if (shouldThrow) {
    throw new Error(errorMessage)
  }
  return <div>Component rendered successfully</div>
}

// Component that throws error on click
const ClickErrorComponent: React.FC = () => {
  const [shouldThrow, setShouldThrow] = React.useState(false)

  if (shouldThrow) {
    throw new Error('Click triggered error')
  }

  return (
    <button onClick={() => setShouldThrow(true)}>
      Click to trigger error
    </button>
  )
}

// These tests make components throw on purpose. React reports a boundary-caught error to window 'error' as
// well; Vitest's jsdom environment turns that into an "Uncaught" failure unless a listener is registered
// (Jest ignored it). Registering one keeps the expected throws from failing unrelated assertions.
const swallowExpectedWindowErrors = (event: ErrorEvent) => {
  event.preventDefault()
}
beforeEach(() => window.addEventListener('error', swallowExpectedWindowErrors))
afterEach(() => window.removeEventListener('error', swallowExpectedWindowErrors))

describe('EnhancedErrorBoundary', () => {
  // Suppress console.error for error boundary tests
  beforeAll(() => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
  })

  afterAll(() => {
    vi.restoreAllMocks()
  })

  afterEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    sessionStorage.clear()
  })

  describe('Normal Operation', () => {
    it('should render children when no error occurs', () => {
      render(
        <EnhancedErrorBoundary>
          <div>Test content</div>
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText('Test content')).toBeInTheDocument()
    })

    it('should render multiple children without error', () => {
      render(
        <EnhancedErrorBoundary>
          <h1>Title</h1>
          <p>Description</p>
          <button>Action</button>
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText('Title')).toBeInTheDocument()
      expect(screen.getByText('Description')).toBeInTheDocument()
      expect(screen.getByText('Action')).toBeInTheDocument()
    })
  })

  describe('Error Catching', () => {
    it('should catch and display component errors', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Something Went Wrong/i)).toBeInTheDocument()
      expect(screen.getByText(/Test error message/i)).toBeInTheDocument()
    })

    it('should display different error messages for different errors', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} errorMessage="Custom error" />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Custom error/i)).toBeInTheDocument()
    })

    it('should display error in red box with proper styling', () => {
      const { container } = render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const errorBox = container.querySelector('.bg-red-50')
      expect(errorBox).toBeInTheDocument()
    })
  })

  describe('Error Logging Integration', () => {
    it('should call error logger when error occurs', () => {
      const mockLogger = {
        logComponentError: vi.fn(),
        log: vi.fn(),
        getAllLogs: vi.fn(() => []),
      }
      ;(getErrorLogger as vi.Mock).mockReturnValue(mockLogger)

      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} errorMessage="Test error" />
        </EnhancedErrorBoundary>
      )

      expect(mockLogger.logComponentError).toHaveBeenCalledWith(
        'ErrorBoundary',
        expect.any(Error),
        expect.objectContaining({
          componentStack: expect.any(String),
        })
      )
    })

    it('should pass context to error logger', () => {
      const mockLogger = {
        logComponentError: vi.fn(),
        log: vi.fn(),
        getAllLogs: vi.fn(() => []),
      }
      ;(getErrorLogger as vi.Mock).mockReturnValue(mockLogger)

      const testContext = { userId: '123', theme: 'dark' }

      render(
        <EnhancedErrorBoundary context={testContext}>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(mockLogger.logComponentError).toHaveBeenCalledWith(
        'ErrorBoundary',
        expect.any(Error),
        expect.objectContaining(testContext)
      )
    })

    it('should call optional onError callback', () => {
      const onError = vi.fn()

      render(
        <EnhancedErrorBoundary onError={onError}>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(onError).toHaveBeenCalledWith(
        expect.any(Error),
        expect.objectContaining({
          componentStack: expect.any(String),
        })
      )
    })
  })

  describe('Recovery Suggestions', () => {
    it('should display recovery suggestions', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Try one of these/i)).toBeInTheDocument()
    })

    it('should call getRecoverySuggestions with error context', () => {
      render(
        <EnhancedErrorBoundary context={{ userId: '123' }}>
          <ErrorThrowingComponent shouldThrow={true} errorMessage="Test error" />
        </EnhancedErrorBoundary>
      )

      expect(getRecoverySuggestions).toHaveBeenCalledWith(
        expect.objectContaining({
          errorMessage: 'Test error',
          category: 'component',
          severity: 'high',
          context: expect.objectContaining({ userId: '123' }),
        })
      )
    })

    it('should render suggestion buttons', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText('Try Again')).toBeInTheDocument()
      expect(screen.getByText('Go to Dashboard')).toBeInTheDocument()
    })

    it('should display suggestion descriptions', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Attempt to reload the component/i)).toBeInTheDocument()
      expect(screen.getByText(/Return to the main dashboard/i)).toBeInTheDocument()
    })

    it('should not display suggestions when empty array returned', () => {
      ;(getRecoverySuggestions as vi.Mock).mockReturnValueOnce([])

      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.queryByText(/Try one of these/i)).not.toBeInTheDocument()
    })
  })

  describe('Error Advice', () => {
    it('should display error advice text', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/What happened/i)).toBeInTheDocument()
      expect(
        screen.getByText(/An unexpected error occurred. Please try again or contact support./i)
      ).toBeInTheDocument()
    })

    it('should call getErrorAdvice with category and severity', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(getErrorAdvice).toHaveBeenCalledWith('component', 'high')
    })

    it('should display advice in blue info box', () => {
      const { container } = render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const adviceBox = container.querySelector('.bg-blue-50')
      expect(adviceBox).toBeInTheDocument()
    })
  })

  describe('Retry Functionality', () => {
    it('should have Try Again button', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const retryButton = screen.getByRole('button', { name: /Try Again/i })
      expect(retryButton).toBeInTheDocument()
    })

    it('should reset error state on retry', () => {
      const { rerender } = render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Something Went Wrong/i)).toBeInTheDocument()

      const retryButton = screen.getByRole('button', { name: /Try Again/i })
      fireEvent.click(retryButton)

      // After retry, error should be cleared
      expect(screen.queryByText(/Something Went Wrong/i)).not.toBeInTheDocument()
    })

    it('should display retry count', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const retryButton = screen.getByRole('button', { name: /Try Again/i })
      fireEvent.click(retryButton)
      fireEvent.click(retryButton)

      // After 2 retries, should show (2/3)
      expect(screen.getByText(/\(2\/3\)/)).toBeInTheDocument()
    })

    it('should disable retry button after max retries', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const retryButton = screen.getByRole('button', { name: /Try Again/i })

      // Click 3 times (max retries)
      fireEvent.click(retryButton)
      fireEvent.click(retryButton)
      fireEvent.click(retryButton)

      expect(retryButton).toBeDisabled()
    })

    it('should show alert when max retries exceeded', () => {
      const alertSpy = vi.spyOn(window, 'alert').mockImplementation()

      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const retryButton = screen.getByRole('button', { name: /Try Again/i })

      // Click 4 times (exceeding max retries)
      fireEvent.click(retryButton)
      fireEvent.click(retryButton)
      fireEvent.click(retryButton)
      fireEvent.click(retryButton)

      expect(alertSpy).toHaveBeenCalledWith(
        expect.stringContaining('Maximum retry attempts reached')
      )

      alertSpy.mockRestore()
    })
  })

  describe('Action Buttons', () => {
    it('should have Go to Dashboard button', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const dashboardButton = screen.getByRole('button', { name: /Go to Dashboard/i })
      expect(dashboardButton).toBeInTheDocument()
    })

    it('should have Log Out button', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const logoutButton = screen.getByRole('button', { name: /Log Out/i })
      expect(logoutButton).toBeInTheDocument()
    })

    it('should navigate to dashboard on button click', () => {
      const originalHref = window.location.href
      delete (window.location as any).href
      window.location.href = 'http://localhost'

      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const dashboardButton = screen.getByRole('button', { name: /Go to Dashboard/i })
      fireEvent.click(dashboardButton)

      expect(window.location.href).toContain('/en/dashboard')
    })

    it('should clear auth data on logout', () => {
      const setItemSpy = vi.spyOn(Storage.prototype, 'setItem')
      const removeItemSpy = vi.spyOn(Storage.prototype, 'removeItem')

      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const logoutButton = screen.getByRole('button', { name: /Log Out/i })
      fireEvent.click(logoutButton)

      expect(removeItemSpy).toHaveBeenCalledWith('authToken')
      expect(setItemSpy).toHaveBeenCalled() // sessionStorage.clear() internally

      setItemSpy.mockRestore()
      removeItemSpy.mockRestore()
    })
  })

  describe('Error Details (Development Mode)', () => {
    const originalEnv = process.env.NODE_ENV

    beforeEach(() => {
      process.env.NODE_ENV = 'development'
    })

    afterEach(() => {
      process.env.NODE_ENV = originalEnv
    })

    it('should show error details button in development', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const detailsButton = screen.getByRole('button', { name: /Show Error Details/i })
      expect(detailsButton).toBeInTheDocument()
    })

    it('should toggle error details visibility', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const detailsButton = screen.getByRole('button', { name: /Show Error Details/i })

      // Initially hidden
      expect(screen.queryByText(/Stack Trace/i)).not.toBeInTheDocument()

      // Click to show
      fireEvent.click(detailsButton)
      expect(screen.getByText(/Stack Trace/i)).toBeInTheDocument()

      // Click to hide
      fireEvent.click(detailsButton)
      expect(screen.queryByText(/Stack Trace/i)).not.toBeInTheDocument()
    })

    it('should display stack trace in details', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const detailsButton = screen.getByRole('button', { name: /Show Error Details/i })
      fireEvent.click(detailsButton)

      expect(screen.getByText(/Stack Trace/i)).toBeInTheDocument()
    })

    it('should display component stack in details', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const detailsButton = screen.getByRole('button', { name: /Show Error Details/i })
      fireEvent.click(detailsButton)

      expect(screen.getByText(/Component Stack/i)).toBeInTheDocument()
    })
  })

  describe('Error Reference ID', () => {
    it('should display error reference ID', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Error Reference ID/i)).toBeInTheDocument()
    })

    it('should show unique reference ID', () => {
      const { rerender } = render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} errorMessage="Error 1" />
        </EnhancedErrorBoundary>
      )

      const refId1 = screen.getByText(/Error Reference ID/i).textContent

      vi.clearAllMocks()

      rerender(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} errorMessage="Error 2" />
        </EnhancedErrorBoundary>
      )

      const refId2 = screen.getByText(/Error Reference ID/i).textContent

      expect(refId1).not.toBe(refId2)
    })
  })

  describe('Support Contact', () => {
    it('should display support contact link', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const supportLink = screen.getByRole('link', { name: /contact support/i })
      expect(supportLink).toBeInTheDocument()
      expect(supportLink).toHaveAttribute('href', 'mailto:support@example.com')
    })

    it('should include error details in support message hint', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Include the error details above/i)).toBeInTheDocument()
    })
  })

  describe('Custom Fallback', () => {
    it('should use custom fallback when provided', () => {
      const customFallback = (error: Error, retry: () => void) => (
        <div>
          <h1>Custom Error UI</h1>
          <p>{error.message}</p>
          <button onClick={retry}>Custom Retry</button>
        </div>
      )

      render(
        <EnhancedErrorBoundary fallback={customFallback}>
          <ErrorThrowingComponent shouldThrow={true} errorMessage="Custom error" />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Custom Error UI/i)).toBeInTheDocument()
      expect(screen.getByText(/Custom error/i)).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Custom Retry/i })).toBeInTheDocument()
    })

    it('should use default UI when no fallback provided', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Something Went Wrong/i)).toBeInTheDocument()
    })
  })

  describe('Dark Mode Support', () => {
    it('should apply dark mode classes', () => {
      const { container } = render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      // Check for dark mode classes
      const errorCard = container.querySelector('.dark\\:bg-gray-800')
      expect(errorCard).toBeInTheDocument()
    })
  })

  describe('Edge Cases', () => {
    it('should handle errors with no message', () => {
      const ErrorNoMessage = () => {
        throw new Error()
      }

      render(
        <EnhancedErrorBoundary>
          <ErrorNoMessage />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Something Went Wrong/i)).toBeInTheDocument()
    })

    it('should handle null error gracefully', () => {
      render(
        <EnhancedErrorBoundary>
          <div>Safe content</div>
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Safe content/i)).toBeInTheDocument()
    })

    it('should recover from errors with async operations', async () => {
      const AsyncErrorComponent = () => {
        throw new Error('Async error')
      }

      render(
        <EnhancedErrorBoundary>
          <AsyncErrorComponent />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Async error/i)).toBeInTheDocument()

      const retryButton = screen.getByRole('button', { name: /Try Again/i })
      fireEvent.click(retryButton)

      // Error should be cleared after retry
      await waitFor(() => {
        expect(screen.queryByText(/Async error/i)).not.toBeInTheDocument()
      })
    })

    it('should handle rapid successive errors', () => {
      const QuickErrorComponent = () => {
        throw new Error('Quick error')
      }

      const { rerender } = render(
        <EnhancedErrorBoundary>
          <QuickErrorComponent />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByText(/Quick error/i)).toBeInTheDocument()

      // Rapid rerender with different error
      rerender(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} errorMessage="Another error" />
        </EnhancedErrorBoundary>
      )

      // Should show latest error
      expect(screen.getByText(/Another error/i)).toBeInTheDocument()
    })
  })

  describe('Accessibility', () => {
    it('should have proper heading hierarchy', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      const heading = screen.getByRole('heading', { name: /Something Went Wrong/i })
      expect(heading).toHaveClass('text-2xl')
    })

    it('should have accessible button labels', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      expect(screen.getByRole('button', { name: /Try Again/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Go to Dashboard/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Log Out/i })).toBeInTheDocument()
    })

    it('should provide alt text for icons', () => {
      render(
        <EnhancedErrorBoundary>
          <ErrorThrowingComponent shouldThrow={true} />
        </EnhancedErrorBoundary>
      )

      // Alert icons should have semantic meaning
      const alertTriangle = screen.getByRole('heading')
      expect(alertTriangle).toBeInTheDocument()
    })
  })
})
