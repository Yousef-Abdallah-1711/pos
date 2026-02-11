'use client'

import { useState } from 'react'
import { AlertTriangle, X } from 'lucide-react'
import { useTenantStatus } from '@/hooks/useTenantStatus'

export function TrialWarningBanner() {
  const { status, daysRemaining, isLoading } = useTenantStatus()
  const [dismissed, setDismissed] = useState(false)

  if (dismissed || isLoading) return null
  if (status !== 'trial') return null
  if (daysRemaining === null || daysRemaining > 7 || daysRemaining <= 0) return null

  return (
    <div className="bg-amber-50 dark:bg-amber-900/20 border-b border-amber-200 dark:border-amber-800 px-4 py-3">
      <div className="max-w-7xl mx-auto flex items-center justify-between gap-3">
        <div className="flex items-center gap-2 text-amber-800 dark:text-amber-300">
          <AlertTriangle className="h-4 w-4 flex-shrink-0" />
          <p className="text-sm">
            Your trial expires in <strong>{daysRemaining} day{daysRemaining === 1 ? '' : 's'}</strong>.
            Contact support to upgrade your plan.
          </p>
        </div>
        <button
          onClick={() => setDismissed(true)}
          className="text-amber-600 dark:text-amber-400 hover:text-amber-800 dark:hover:text-amber-200 transition flex-shrink-0"
          aria-label="Dismiss"
        >
          <X className="h-4 w-4" />
        </button>
      </div>
    </div>
  )
}
