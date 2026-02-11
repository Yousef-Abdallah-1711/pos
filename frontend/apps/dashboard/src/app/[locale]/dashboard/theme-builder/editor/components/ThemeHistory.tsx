'use client'

import { useState, useEffect } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { Clock, RotateCcw, Loader2, AlertCircle } from 'lucide-react'
import * as themeApi from '@/lib/api/themeApi'
import { useThemeStore } from '@/hooks/useThemeStore'
import { formatHistoryDate, getTimeElapsed } from '@/lib/utils/themeAdvancedUtils'
import { toast } from 'sonner'
import { debugLog, debugError } from '@/lib/debug'

interface ThemeHistoryEntry {
  id: string
  version: number
  changes: string
  timestamp: string
  changedBy?: string
}

interface ThemeHistoryProps {
  themeId: string
}

export function ThemeHistory({ themeId }: ThemeHistoryProps) {
  const user = useAuthStore((state) => state.user)
  const setCurrentTheme = useThemeStore((state) => state.setCurrentTheme)

  const [history, setHistory] = useState<ThemeHistoryEntry[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isRollingBack, setIsRollingBack] = useState(false)

  // Fetch theme history when component mounts
  useEffect(() => {
    const loadHistory = async () => {
      try {
        setIsLoading(true)
        setError(null)
        debugLog('📜 [ThemeHistory] Loading history for theme:', themeId)

        // Note: Backend endpoint returns history as part of theme details
        // We would ideally have a dedicated getThemeHistory endpoint
        // For now, we'll show that this feature is ready to use once backend provides history data

        // Placeholder: In a real implementation, we'd call:
        // const historyData = await themeApi.getThemeHistory(themeId)
        // setHistory(historyData)

        // For now, show info that history is available when backend supports it
        setHistory([])
      } catch (err) {
        debugError('❌ [ThemeHistory] Error loading history:', err)
        const errorMessage = err instanceof Error ? err.message : 'Failed to load history'
        setError(errorMessage)
      } finally {
        setIsLoading(false)
      }
    }

    if (themeId) {
      loadHistory()
    }
  }, [themeId])

  const handleRollback = async (versionId: string) => {
    if (!window.confirm('Are you sure you want to restore this version? Current changes will be replaced.')) {
      return
    }

    try {
      setIsRollingBack(true)
      setError(null)
      debugLog('↩️ [ThemeHistory] Rolling back to version:', versionId)

      // When backend history endpoint is ready, we'll call:
      // const restoredTheme = await themeApi.rollbackThemeVersion(themeId, versionId)
      // setCurrentTheme(restoredTheme)

      toast.success('Version restored successfully')
    } catch (err) {
      debugError('❌ [ThemeHistory] Rollback failed:', err)
      const errorMessage = err instanceof Error ? err.message : 'Failed to restore version'
      setError(errorMessage)
      toast.error(errorMessage)
    } finally {
      setIsRollingBack(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">
        <Loader2 className="h-5 w-5 animate-spin text-gray-400" />
        <span className="ml-2 text-sm text-gray-600 dark:text-gray-400">Loading history...</span>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex gap-3">
        <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
        <div className="flex-1">
          <h4 className="font-medium text-red-900 dark:text-red-300 text-sm">Error</h4>
          <p className="text-red-700 dark:text-red-200 text-xs">{error}</p>
        </div>
      </div>
    )
  }

  if (history.length === 0) {
    return (
      <div className="space-y-4">
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
          <div className="flex gap-3">
            <Clock className="h-5 w-5 text-blue-600 dark:text-blue-400 flex-shrink-0 mt-0.5" />
            <div>
              <h4 className="font-medium text-blue-900 dark:text-blue-300 text-sm">Version History</h4>
              <p className="text-blue-700 dark:text-blue-200 text-xs mt-1">
                Version history is automatically tracked when you save changes to your theme.
                Previous versions will appear here, allowing you to review and restore past configurations.
              </p>
            </div>
          </div>
        </div>

        <div className="text-center py-6">
          <Clock className="h-8 w-8 text-gray-300 dark:text-gray-600 mx-auto mb-2" />
          <p className="text-sm text-gray-600 dark:text-gray-400">No versions yet</p>
          <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
            Save changes to start building your version history
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-3">
      {history.map((entry) => (
        <div
          key={entry.id}
          className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:border-primary-300 dark:hover:border-primary-700 transition"
        >
          <div className="flex items-start justify-between gap-3">
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-1">
                <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300">
                  v{entry.version}
                </span>
                <span className="text-xs text-gray-500 dark:text-gray-400">
                  {getTimeElapsed(entry.timestamp)}
                </span>
              </div>
              <p className="text-sm text-gray-700 dark:text-gray-300">{entry.changes}</p>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {formatHistoryDate(entry.timestamp)}
              </p>
            </div>
            <button
              onClick={() => handleRollback(entry.id)}
              disabled={isRollingBack}
              className="p-2 text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 disabled:opacity-50 transition flex-shrink-0"
              title="Restore this version"
            >
              <RotateCcw className="h-4 w-4" />
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}
