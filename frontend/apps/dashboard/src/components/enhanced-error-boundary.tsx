/**
 * Enhanced Error Boundary Component
 * Catches React component errors and displays user-friendly error message
 * with recovery suggestions and error logging
 */

'use client'

import React, { ReactNode, ErrorInfo, useState } from 'react'
import {
  AlertCircle,
  RotateCcw,
  ChevronDown,
  AlertTriangle,
  Brain,
  LogOut,
} from 'lucide-react'
import { debugError } from '@/lib/debug'
import { getErrorLogger } from '@/lib/error-handling/error-logger'
import {
  getRecoverySuggestions,
  getErrorAdvice,
  type ErrorRecoveryContext,
} from '@/lib/error-handling/error-recovery-suggestions'

interface EnhancedErrorBoundaryProps {
  children: ReactNode
  fallback?: (error: Error, retry: () => void) => ReactNode
  onError?: (error: Error, errorInfo: ErrorInfo) => void
  context?: Record<string, any>
}

interface EnhancedErrorBoundaryState {
  hasError: boolean
  error: Error | null
  errorInfo: ErrorInfo | null
  showDetails: boolean
  retryCount: number
  errorReferenceId: string | null
}

/**
 * Enhanced Error Boundary Component
 * Features:
 * - Error logging service
 * - Context-aware recovery suggestions
 * - Error categorization and severity detection
 * - Development vs production error display
 * - Multiple recovery options
 */
export class EnhancedErrorBoundary extends React.Component<
  EnhancedErrorBoundaryProps,
  EnhancedErrorBoundaryState
> {
  private static errorReferenceSequence = 0
  private errorLogger = getErrorLogger()
  private readonly MAX_RETRIES = 3

  constructor(props: EnhancedErrorBoundaryProps) {
    super(props)
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
      showDetails: false,
      retryCount: 0,
      errorReferenceId: null,
    }
  }

  static getDerivedStateFromError(error: Error): Partial<EnhancedErrorBoundaryState> {
    return {
      hasError: true,
      error,
      errorReferenceId: `${Date.now()}-${++EnhancedErrorBoundary.errorReferenceSequence}`,
    }
  }

  componentDidUpdate(previousProps: EnhancedErrorBoundaryProps) {
    if (previousProps.children !== this.props.children && this.state.hasError) {
      this.setState({
        hasError: false,
        error: null,
        errorInfo: null,
        showDetails: false,
        retryCount: 0,
        errorReferenceId: null,
      })
    }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Log error
    this.errorLogger.logComponentError('ErrorBoundary', error, {
      componentStack: errorInfo.componentStack,
      ...this.props.context,
    })

    // Update state
    this.setState({
      error,
      errorInfo,
      showDetails: false,
    })

    // Call optional error callback
    if (this.props.onError) {
      this.props.onError(error, errorInfo)
    }

    // Log to console in development
    if (process.env.NODE_ENV === 'development') {
      debugError('🔴 [ErrorBoundary] Caught error:', error)
      debugError('🔴 [ErrorBoundary] Error info:', errorInfo)
    }
  }

  handleRetry = () => {
    const newRetryCount = this.state.retryCount + 1

    if (this.state.retryCount >= this.MAX_RETRIES) {
      return
    }

    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
      showDetails: false,
      retryCount: newRetryCount,
      errorReferenceId: null,
    })
  }

  handleShowDetails = () => {
    this.setState((state) => ({
      showDetails: !state.showDetails,
    }))
  }

  handleLogOut = () => {
    // Clear session/auth data
    localStorage.removeItem('authToken')
    sessionStorage.clear()
    window.location.href = '/login'
  }

  render() {
    if (this.state.hasError && this.state.error) {
      // Use custom fallback if provided
      if (this.props.fallback) {
        return this.props.fallback(this.state.error, this.handleRetry)
      }

      const error = this.state.error
      const context: ErrorRecoveryContext = {
        errorMessage: error.message,
        category: 'component', // Will be categorized by error logger
        severity: 'high',
        context: this.props.context,
      }

      const suggestions = getRecoverySuggestions(context)
      const advice = getErrorAdvice(context.category, context.severity)

      return (
        <div className="w-full min-h-screen flex items-center justify-center p-4 bg-gray-50 dark:bg-gray-900">
          <div className="max-w-3xl w-full">
            {/* Main Error Card */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
              {/* Header */}
              <div className="bg-gradient-to-r from-red-500 to-red-600 dark:from-red-600 dark:to-red-700 px-6 py-4">
                <div className="flex items-start gap-3">
                  <AlertTriangle className="h-6 w-6 text-white flex-shrink-0 mt-0.5" />
                  <div>
                    <h1 className="text-2xl font-bold text-white">
                      Something Went Wrong
                    </h1>
                    <p className="text-red-100 mt-1">
                      An unexpected error occurred. Please try one of the options below.
                    </p>
                  </div>
                </div>
              </div>

              {/* Content */}
              <div className="p-6 space-y-6">
                {/* Error Message */}
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                  <p className="text-sm text-red-900 dark:text-red-100 font-mono break-words">
                    {error.message || 'An unexpected error occurred'}
                  </p>
                </div>

                {/* Advice */}
                <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
                  <div className="flex gap-3">
                    <Brain className="h-5 w-5 text-blue-600 dark:text-blue-400 flex-shrink-0 mt-0.5" />
                    <div>
                      <p className="text-sm font-medium text-blue-900 dark:text-blue-100">
                        What happened?
                      </p>
                      <p className="text-sm text-blue-800 dark:text-blue-200 mt-1">
                        {advice}
                      </p>
                    </div>
                  </div>
                </div>

                {/* Recovery Suggestions */}
                {suggestions.length > 0 && (
                  <div className="space-y-3">
                    <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
                      Try one of these:
                    </h3>
                    <div className="grid gap-2">
                      {suggestions.map((suggestion, index) => (
                        <button
                          key={index}
                          onClick={suggestion.action.onClick}
                          className="text-left p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 transition"
                        >
                          <div className="flex items-start justify-between gap-2">
                            <div>
                              <p className="font-medium text-gray-900 dark:text-white text-sm">
                                {suggestion.title}
                              </p>
                              <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                                {suggestion.description}
                              </p>
                            </div>
                            <span className="px-2 py-1 bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 rounded text-xs font-medium flex-shrink-0">
                              {suggestion.action.label}
                            </span>
                          </div>
                        </button>
                      ))}
                    </div>
                  </div>
                )}

                {/* Error Details (Development Only) */}
                {process.env.NODE_ENV !== 'production' && this.state.errorInfo && (
                  <div className="border-t border-gray-200 dark:border-gray-700 pt-6">
                    <button
                      onClick={this.handleShowDetails}
                      className="flex items-center gap-2 text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition"
                    >
                      <ChevronDown
                        className={`h-4 w-4 transition-transform ${
                          this.state.showDetails ? 'rotate-180' : ''
                        }`}
                      />
                      {this.state.showDetails ? 'Hide' : 'Show'} Error Details
                    </button>

                    {this.state.showDetails && (
                      <div className="mt-4 space-y-4">
                        <div>
                          <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase mb-2">
                            Stack Trace
                          </p>
                          <pre className="p-3 bg-gray-100 dark:bg-gray-900 rounded text-xs text-gray-800 dark:text-gray-200 overflow-x-auto max-h-48 overflow-y-auto">
                            {error.stack}
                          </pre>
                        </div>

                        <div>
                          <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase mb-2">
                            Component Stack
                          </p>
                          <pre className="p-3 bg-gray-100 dark:bg-gray-900 rounded text-xs text-gray-800 dark:text-gray-200 overflow-x-auto max-h-48 overflow-y-auto">
                            {this.state.errorInfo.componentStack}
                          </pre>
                        </div>
                      </div>
                    )}
                  </div>
                )}

                {/* Actions */}
                <div className="flex gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                  <button
                    onClick={this.handleRetry}
                    disabled={this.state.retryCount >= this.MAX_RETRIES}
                    className="flex-1 flex items-center justify-center gap-2 px-4 py-2 bg-primary-600 hover:bg-primary-700 disabled:bg-gray-300 disabled:cursor-not-allowed dark:bg-primary-500 dark:hover:bg-primary-600 text-white rounded-lg font-medium transition"
                  >
                    <RotateCcw className="h-4 w-4" />
                    Try Again
                    {this.state.retryCount > 0 && (
                      <span className="text-xs">({this.state.retryCount}/{this.MAX_RETRIES})</span>
                    )}
                  </button>

                  <a
                    href="/en/dashboard"
                    className="flex-1 px-4 py-2 bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-900 dark:text-white rounded-lg font-medium transition"
                  >
                    Go to Dashboard
                  </a>

                  <button
                    onClick={this.handleLogOut}
                    className="px-4 py-2 bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-900 dark:text-white rounded-lg font-medium transition flex items-center gap-2"
                  >
                    <LogOut className="h-4 w-4" />
                    Log Out
                  </button>
                </div>

                {/* Footer Help Text */}
                <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                  <p className="text-xs text-gray-600 dark:text-gray-400">
                    If this problem persists, please{' '}
                    <a
                      href="mailto:support@example.com"
                      className="text-primary-600 dark:text-primary-400 hover:underline"
                    >
                      contact support
                    </a>
                    . Include the error details above to help us resolve the issue faster.
                  </p>
                </div>
              </div>
            </div>

            {/* Info Box */}
            <div className="mt-6 p-4 bg-gray-100 dark:bg-gray-800 rounded-lg">
              <p className="text-xs text-gray-600 dark:text-gray-400">
                <span className="font-medium">Error Reference ID:</span>{' '}
                <code className="text-gray-500 dark:text-gray-500">
                  {this.state.errorReferenceId}
                </code>
              </p>
            </div>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}

export default EnhancedErrorBoundary
