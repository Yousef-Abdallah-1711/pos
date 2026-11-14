'use client'

import React, { ReactNode, ErrorInfo } from 'react'
import { AlertCircle, RotateCcw } from 'lucide-react'
import { debugLog, debugError } from '@/lib/debug'

interface ErrorBoundaryProps {
  children: ReactNode
  fallback?: (error: Error, retry: () => void) => ReactNode
  onError?: (error: Error, errorInfo: ErrorInfo) => void
}

interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
  errorInfo: ErrorInfo | null
}

/**
 * Error Boundary Component
 * Catches React component errors and displays user-friendly error message
 *
 * Usage:
 * ```tsx
 * <ErrorBoundary onError={(err) => console.error(err)}>
 *   <ThemeEditor />
 * </ErrorBoundary>
 * ```
 */
export class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
    }
  }

  static getDerivedStateFromError(error: Error): Partial<ErrorBoundaryState> {
    return { hasError: true }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    debugError('🔴 [ErrorBoundary] Caught error:', error);
    debugError('🔴 [ErrorBoundary] Error info:', errorInfo);

    this.setState({
      error,
      errorInfo,
    })

    // Call optional error callback
    if (this.props.onError) {
      this.props.onError(error, errorInfo)
    }
  }

  handleRetry = () => {
    debugLog('🔄 [ErrorBoundary] Retrying...')
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
    })
  }

  render() {
    if (this.state.hasError && this.state.error) {
      // Use custom fallback if provided
      if (this.props.fallback) {
        return this.props.fallback(this.state.error, this.handleRetry)
      }

      // Default error UI
      return (
        <div className="w-full h-full flex items-center justify-center p-6">
          <div className="max-w-2xl w-full">
            {/* Error Card */}
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
              <div className="flex gap-4">
                <div className="flex-shrink-0">
                  <AlertCircle className="h-6 w-6 text-red-600 dark:text-red-400" />
                </div>

                <div className="flex-1 min-w-0">
                  <h2 className="text-lg font-semibold text-red-900 dark:text-red-200 mb-2">
                    Something Went Wrong
                  </h2>

                  {/* Error Message */}
                  <p className="text-sm text-red-800 dark:text-red-300 mb-4">
                    {this.state.error.message || 'An unexpected error occurred'}
                  </p>

                  {/* Error Details (Development Only) */}
                  {process.env.NODE_ENV === 'development' && this.state.errorInfo && (
                    <details className="mb-4">
                      <summary className="cursor-pointer text-sm font-medium text-red-700 dark:text-red-400 hover:text-red-800 dark:hover:text-red-300">
                        Error Details (Dev Only)
                      </summary>
                      <pre className="mt-2 p-3 bg-red-100 dark:bg-red-800 rounded text-xs text-red-900 dark:text-red-100 overflow-x-auto max-h-48 overflow-y-auto">
                        {this.state.errorInfo.componentStack}
                      </pre>
                    </details>
                  )}

                  {/* Action Buttons */}
                  <div className="flex gap-3">
                    <button
                      onClick={this.handleRetry}
                      className="flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 dark:bg-red-500 dark:hover:bg-red-600 text-white rounded-lg font-medium transition"
                    >
                      <RotateCcw className="h-4 w-4" />
                      Try Again
                    </button>

                    <button
                      onClick={() => window.location.href = '/en/dashboard'}
                      className="px-4 py-2 bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-900 dark:text-white rounded-lg font-medium transition"
                    >
                      Go to Dashboard
                    </button>
                  </div>

                  {/* Help Text */}
                  <p className="mt-4 text-xs text-red-700 dark:text-red-400">
                    If this problem persists, please refresh the page or contact support.
                  </p>
                </div>
              </div>
            </div>

            {/* Suggestions */}
            <div className="mt-6 space-y-2">
              <p className="text-sm font-medium text-gray-900 dark:text-white">Things you can try:</p>
              <ul className="text-sm text-gray-600 dark:text-gray-400 space-y-1 list-disc list-inside">
                <li>Click &quot;Try Again&quot; to retry the operation</li>
                <li>Refresh the page (Ctrl+R or Cmd+R)</li>
                <li>Check your internet connection</li>
                <li>Clear browser cache and refresh</li>
                <li>Log out and log back in</li>
              </ul>
            </div>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}

export default ErrorBoundary
