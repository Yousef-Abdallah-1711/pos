'use client'

import { Component, type ErrorInfo, type ReactNode } from 'react'
import { AlertTriangle, RefreshCw } from 'lucide-react'
import type { KDSStats } from '@/stores/kdsStore'

interface KDSErrorBoundaryProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  lastKnownOrderCount: number
  lastKnownStats: KDSStats | null
  onRetry: () => void
  children: ReactNode
}

interface KDSErrorBoundaryState {
  hasError: boolean
}

export class KDSErrorBoundary extends Component<KDSErrorBoundaryProps, KDSErrorBoundaryState> {
  state: KDSErrorBoundaryState = {
    hasError: false,
  }

  static getDerivedStateFromError(): KDSErrorBoundaryState {
    return { hasError: true }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Keep logging visible in production while presenting a safe fallback.
    console.error('[KDS] Grid render error captured by boundary', error, errorInfo)
  }

  private handleRetry = () => {
    this.setState({ hasError: false })
    this.props.onRetry()
  }

  render() {
    if (!this.state.hasError) {
      return this.props.children
    }

    const { locale, t, lastKnownOrderCount, lastKnownStats } = this.props
    const isRTL = locale === 'ar'

    return (
      <div className="rounded-xl border border-red-300 bg-red-50 p-5 dark:border-red-400/40 dark:bg-red-500/10">
        <div className={`mb-4 flex items-start gap-3 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
          <AlertTriangle className="mt-0.5 h-5 w-5 flex-shrink-0 text-red-600 dark:text-red-200" />
          <div>
            <h3 className="text-base font-semibold text-red-700 dark:text-red-100">
              {t('kds.errorBoundaryTitle', 'Kitchen view failed to render')}
            </h3>
            <p className="mt-1 text-sm text-red-700 dark:text-red-200">
              {t('kds.errorBoundaryText', 'Showing last known state. Retry to recover the grid.')}
            </p>
          </div>
        </div>

        <div className={`mb-4 text-xs text-red-700 dark:text-red-100 ${isRTL ? 'text-right' : ''}`}>
          <p>
            {t('kds.lastKnownOrders', 'Last known orders')}: {lastKnownOrderCount}
          </p>
          <p>
            {t('kds.lastKnownActive', 'Last known active orders')}: {lastKnownStats?.active_orders ?? 0}
          </p>
        </div>

        <button
          type="button"
          onClick={this.handleRetry}
          className="touch-manipulation inline-flex min-h-11 items-center gap-2 rounded-lg border border-red-300 bg-red-100 px-4 py-2 text-sm font-medium text-red-700 transition hover:bg-red-200 dark:border-red-300/40 dark:bg-red-500/10 dark:text-red-100 dark:hover:bg-red-500/20"
        >
          <RefreshCw className="h-4 w-4" />
          {t('kds.retry', 'Retry')}
        </button>
      </div>
    )
  }
}
