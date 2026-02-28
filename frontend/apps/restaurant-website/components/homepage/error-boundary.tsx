'use client'

import React, { ReactNode, Component, ErrorInfo } from 'react'
import { AlertTriangle, RefreshCw } from 'lucide-react'
import { debugError } from '@/lib/debug'

interface Props {
  children: ReactNode
  locale?: 'en' | 'ar'
  fallback?: ReactNode
  onError?: (error: Error, info: ErrorInfo) => void
  section?: string
}

interface State {
  hasError: boolean
  error: Error | null
  errorCount: number
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = {
      hasError: false,
      error: null,
      errorCount: 0,
    }
  }

  static getDerivedStateFromError(error: Error): State {
    return {
      hasError: true,
      error,
      errorCount: 0,
    }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    if (process.env.NODE_ENV === 'development') {
      debugError(`[${this.props.section || 'ErrorBoundary'}] Caught error:`, error, info)
    }

    this.props.onError?.(error, info)

    this.setState((prev) => ({
      errorCount: prev.errorCount + 1,
    }))
  }

  handleRetry = () => {
    this.setState({
      hasError: false,
      error: null,
      errorCount: 0,
    })
  }

  render() {
    const { hasError, error } = this.state
    const { children, locale = 'en', fallback, section } = this.props

    if (!hasError) {
      return children
    }

    if (fallback) {
      return fallback
    }

    const isRTL = locale === 'ar'
    const labels = {
      en: {
        title: 'Something went wrong',
        message: `An error occurred while loading ${section || 'this section'}.`,
        retry: 'Try Again',
      },
      ar: {
        title: 'Ø­Ø¯Ø« Ø®Ø·Ø£ Ù…Ø§',
        message: `Ø­Ø¯Ø« Ø®Ø·Ø£ Ø£Ø«Ù†Ø§Ø¡ ØªØ­Ù…ÙŠÙ„ ${section || 'Ù‡Ø°Ø§ Ø§Ù„Ù‚Ø³Ù…'}.`,
        retry: 'Ø­Ø§ÙˆÙ„ Ù…Ø±Ø© Ø£Ø®Ø±Ù‰',
      },
    }

    const t = labels[locale as keyof typeof labels] || labels.en

    return (
      <div className={`error-boundary p-4 rounded-lg border ${isRTL ? 'text-right' : 'text-left'}`}>
        <div className="flex gap-3">
          <AlertTriangle className="error-boundary-icon h-5 w-5 flex-shrink-0 mt-0.5" />
          <div className="flex-1">
            <h3 className="error-boundary-title font-semibold text-sm">{t.title}</h3>
            <p className="error-boundary-message text-sm mt-1">{t.message}</p>
            {error && <p className="error-boundary-stack text-xs mt-2 font-mono overflow-auto max-h-12">{error.message}</p>}
            <button
              onClick={this.handleRetry}
              className="error-boundary-retry mt-3 inline-flex items-center gap-1 px-3 py-1 rounded text-xs font-medium transition-colors hover:opacity-90"
            >
              <RefreshCw className="h-3 w-3" />
              {t.retry}
            </button>
          </div>
        </div>

        <style jsx>{`
          .error-boundary {
            background-color: #fef2f2;
            border-color: #fecaca;
          }

          .error-boundary-icon {
            color: #dc2626;
          }

          .error-boundary-title {
            color: #7f1d1d;
          }

          .error-boundary-message {
            color: #991b1b;
          }

          .error-boundary-stack {
            color: #7f1d1d;
          }

          .error-boundary-retry {
            background-color: #dc2626;
            color: #ffffff;
          }
        `}</style>
      </div>
    )
  }
}

export function withErrorBoundary<P extends object>(
  Component: React.ComponentType<P>,
  errorBoundaryProps?: Omit<Props, 'children'>
) {
  return function WrappedComponent(props: P) {
    return (
      <ErrorBoundary {...errorBoundaryProps}>
        <Component {...props} />
      </ErrorBoundary>
    )
  }
}
