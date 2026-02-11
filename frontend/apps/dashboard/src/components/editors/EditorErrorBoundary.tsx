'use client'

import React, { ReactNode } from 'react'
import { AlertTriangle } from 'lucide-react'
import { debugError } from '@/lib/debug'

interface EditorErrorBoundaryProps {
  children: ReactNode
  componentType?: string
  onReset?: () => void
}

interface EditorErrorBoundaryState {
  hasError: boolean
  error?: Error
  errorInfo?: React.ErrorInfo
}

/**
 * Error Boundary for Component Editors
 * Catches errors in editors and displays a user-friendly error message
 */
export class EditorErrorBoundary extends React.Component<
  EditorErrorBoundaryProps,
  EditorErrorBoundaryState
> {
  constructor(props: EditorErrorBoundaryProps) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: Error): EditorErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    this.setState({ error, errorInfo })
    debugError('Editor error:', error, errorInfo)
  }

  handleReset = () => {
    this.setState({ hasError: false, error: undefined, errorInfo: undefined })
    this.props.onReset?.()
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="w-full bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
          <div className="flex items-start gap-4">
            <AlertTriangle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <h3 className="font-semibold text-red-900 dark:text-red-200 mb-2">
                Editor Error
              </h3>
              <p className="text-sm text-red-700 dark:text-red-300 mb-3">
                {this.state.error?.message || 'An unexpected error occurred in the editor'}
              </p>
              {process.env.NODE_ENV === 'development' && this.state.errorInfo && (
                <details className="mb-4 text-xs text-red-600 dark:text-red-400">
                  <summary className="cursor-pointer font-mono mb-2">Error Details</summary>
                  <pre className="bg-red-100 dark:bg-red-800 p-2 rounded overflow-auto max-h-48 font-mono text-xs">
                    {this.state.errorInfo.componentStack}
                  </pre>
                </details>
              )}
              <button
                onClick={this.handleReset}
                className="inline-flex items-center gap-2 px-3 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg text-sm transition"
              >
                Try Again
              </button>
            </div>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}
