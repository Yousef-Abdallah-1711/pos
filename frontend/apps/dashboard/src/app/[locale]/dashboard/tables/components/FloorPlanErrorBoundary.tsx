'use client'

import { Component, type ErrorInfo, type ReactNode } from 'react'
import { AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface FloorPlanErrorBoundaryProps {
  children: ReactNode
  t: (path: string, fallback?: string) => string
}

interface FloorPlanErrorBoundaryState {
  hasError: boolean
}

export class FloorPlanErrorBoundary extends Component<
  FloorPlanErrorBoundaryProps,
  FloorPlanErrorBoundaryState
> {
  state: FloorPlanErrorBoundaryState = {
    hasError: false,
  }

  static getDerivedStateFromError(): FloorPlanErrorBoundaryState {
    return {
      hasError: true,
    }
  }

  componentDidCatch(_error: Error, _errorInfo: ErrorInfo) {
    // Intentionally silent: fallback UI is rendered for the user and error logs are captured elsewhere.
  }

  private handleRetry = () => {
    this.setState({ hasError: false })
  }

  render() {
    const { children, t } = this.props
    const { hasError } = this.state

    if (hasError) {
      return (
        <div className="rounded-lg border border-red-200 bg-red-50 p-6 dark:border-red-800 dark:bg-red-900/20">
          <div className="flex items-start gap-3">
            <AlertTriangle className="mt-0.5 h-5 w-5 text-red-600 dark:text-red-400" />
            <div>
              <p className="font-semibold text-red-900 dark:text-red-200">
                {t('tables.floorPlanErrorTitle', 'Floor plan temporarily unavailable')}
              </p>
              <p className="mt-1 text-sm text-red-800 dark:text-red-300">
                {t(
                  'tables.floorPlanErrorDescription',
                  'An unexpected issue occurred while rendering the floor plan.'
                )}
              </p>
              <Button
                type="button"
                variant="outline"
                className="mt-3"
                onClick={this.handleRetry}
              >
                {t('common.retry', 'Retry')}
              </Button>
            </div>
          </div>
        </div>
      )
    }

    return children
  }
}
