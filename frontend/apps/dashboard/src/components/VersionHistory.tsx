
import type { JSX } from 'react'
/**
 * VersionHistory Component
 * Displays theme version history with timeline, comparisons, and rollback
 */

'use client'

import React, { useState, useCallback } from 'react'
import type { ThemeVersion, ChangeType, VersionComparison } from '@/types/version-history'
import { CHANGE_TYPE_CONFIG, timeSinceVersion, formatVersionLabel } from '@/types/version-history'

interface VersionHistoryProps {
  /** Array of versions */
  versions: ThemeVersion[]
  /** Current active version */
  currentVersion?: ThemeVersion | null
  /** Handler when version is selected */
  onSelectVersion: (version: ThemeVersion) => void
  /** Handler to rollback to a version */
  onRollback: (version: ThemeVersion, reason?: string) => void
  /** Handler to compare versions */
  onCompare?: (fromVersion: ThemeVersion, toVersion: ThemeVersion) => void
  /** Loading state */
  isLoading?: boolean
  /** Error message */
  error?: string | null
  /** Custom CSS classes */
  className?: string
}

/**
 * VersionHistory Component
 */
export function VersionHistory({
  versions,
  currentVersion,
  onSelectVersion,
  onRollback,
  onCompare,
  isLoading = false,
  error,
  className = '',
}: VersionHistoryProps): JSX.Element {
  const [selectedVersion, setSelectedVersion] = useState<ThemeVersion | null>(null)
  const [comparisonMode, setComparisonMode] = useState(false)
  const [rollbackVersion, setRollbackVersion] = useState<ThemeVersion | null>(null)
  const [rollbackReason, setRollbackReason] = useState('')

  /**
   * Handle version selection
   */
  const handleSelectVersion = useCallback(
    (version: ThemeVersion) => {
      setSelectedVersion(version)
      onSelectVersion(version)
    },
    [onSelectVersion]
  )

  /**
   * Handle rollback
   */
  const handleRollback = useCallback(() => {
    if (rollbackVersion) {
      onRollback(rollbackVersion, rollbackReason || undefined)
      setRollbackVersion(null)
      setRollbackReason('')
    }
  }, [rollbackVersion, rollbackReason, onRollback])

  if (isLoading) {
    return (
      <div className={`flex items-center justify-center py-12 ${className}`}>
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mb-4" />
          <p className="text-gray-600 dark:text-gray-400">Loading version history...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className={`rounded-lg bg-red-50 dark:bg-red-900/20 p-6 ${className}`}>
        <div className="flex items-start gap-3">
          <span className="text-red-600 dark:text-red-400 text-2xl">⚠️</span>
          <div>
            <h3 className="font-semibold text-red-900 dark:text-red-100">Error Loading History</h3>
            <p className="text-sm text-red-700 dark:text-red-300">{error}</p>
          </div>
        </div>
      </div>
    )
  }

  if (versions.length === 0) {
    return (
      <div className={`rounded-lg bg-gray-50 dark:bg-gray-800 p-8 text-center ${className}`}>
        <div className="text-4xl mb-3">📋</div>
        <h3 className="font-semibold text-gray-900 dark:text-white mb-2">No Version History</h3>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          No versions recorded yet. Make changes to your theme to start version history.
        </p>
      </div>
    )
  }

  return (
    <div className={className}>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Timeline */}
        <div className="lg:col-span-2">
          <VersionTimeline
            versions={versions}
            selectedVersion={selectedVersion}
            currentVersion={currentVersion}
            onSelectVersion={handleSelectVersion}
          />
        </div>

        {/* Details Panel */}
        <div className="lg:col-span-1 space-y-4">
          {selectedVersion && (
            <VersionDetailsPanel
              version={selectedVersion}
              isCurrent={selectedVersion.id === currentVersion?.id}
              onRollback={() => setRollbackVersion(selectedVersion)}
            />
          )}

          {/* Rollback Confirmation */}
          {rollbackVersion && (
            <RollbackConfirmationPanel
              version={rollbackVersion}
              reason={rollbackReason}
              onReasonChange={setRollbackReason}
              onConfirm={handleRollback}
              onCancel={() => {
                setRollbackVersion(null)
                setRollbackReason('')
              }}
            />
          )}
        </div>
      </div>
    </div>
  )
}

/**
 * VersionTimeline Component
 */
function VersionTimeline({
  versions,
  selectedVersion,
  currentVersion,
  onSelectVersion,
}: {
  versions: ThemeVersion[]
  selectedVersion: ThemeVersion | null
  currentVersion?: ThemeVersion | null
  onSelectVersion: (version: ThemeVersion) => void
}): JSX.Element {
  return (
    <div className="space-y-2">
      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Version Timeline</h3>

      <div className="relative space-y-0">
        {versions.map((version, index) => {
          const config = CHANGE_TYPE_CONFIG[version.changeType]
          const isCurrent = version.id === currentVersion?.id
          const isSelected = version.id === selectedVersion?.id

          return (
            <div key={version.id} className="relative pl-8">
              {/* Timeline line */}
              {index < versions.length - 1 && (
                <div className="absolute left-3 top-10 bottom-0 w-0.5 bg-gray-300 dark:bg-gray-600" />
              )}

              {/* Timeline dot */}
              <div
                className={`
                  absolute left-0 top-2 w-6 h-6 rounded-full border-2 flex items-center justify-center
                  transition-all cursor-pointer
                  ${
                    isSelected
                      ? 'border-blue-600 bg-blue-100 dark:bg-blue-900'
                      : 'border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700'
                  }
                  ${isCurrent ? 'ring-2 ring-green-400' : ''}
                `}
                onClick={() => onSelectVersion(version)}
              >
                <span className="text-xs">{config?.icon || '📝'}</span>
              </div>

              {/* Version Card */}
              <button
                onClick={() => onSelectVersion(version)}
                className={`
                  w-full text-left p-3 rounded-lg border-2 transition-all
                  ${
                    isSelected
                      ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 shadow-md'
                      : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                  }
                `}
              >
                <div className="flex items-start justify-between mb-1">
                  <div>
                    <p className="font-semibold text-gray-900 dark:text-white">
                      {formatVersionLabel(version.versionNumber, isCurrent)}
                    </p>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      {config?.label || version.changeType}
                    </p>
                  </div>
                  {isCurrent && (
                    <span className="px-2 py-1 rounded-full text-xs font-medium bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-200">
                      Current
                    </span>
                  )}
                </div>

                <p className="text-sm text-gray-700 dark:text-gray-300 mb-2 line-clamp-2">
                  {version.changeSummary}
                </p>

                <div className="flex items-center justify-between text-xs text-gray-600 dark:text-gray-400">
                  <span>{version.authorName || 'System'}</span>
                  <span>{timeSinceVersion(version.createdAt)}</span>
                </div>
              </button>
            </div>
          )
        })}
      </div>
    </div>
  )
}

/**
 * VersionDetailsPanel Component
 */
function VersionDetailsPanel({
  version,
  isCurrent,
  onRollback,
}: {
  version: ThemeVersion
  isCurrent: boolean
  onRollback: () => void
}): JSX.Element {
  const config = CHANGE_TYPE_CONFIG[version.changeType]

  return (
    <div className="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4 space-y-4">
      <div>
        <h4 className="font-semibold text-gray-900 dark:text-white mb-2">Version Details</h4>

        <div className="space-y-2 text-sm">
          <DetailRow label="Version" value={formatVersionLabel(version.versionNumber, isCurrent)} />
          <DetailRow label="Type" value={config?.label || version.changeType} />
          <DetailRow label="Author" value={version.authorName || 'System'} />
          <DetailRow label="Date" value={new Date(version.createdAt).toLocaleString()} />
          {version.rollbackReason && (
            <DetailRow label="Reason" value={version.rollbackReason} />
          )}
        </div>
      </div>

      {/* Summary */}
      <div>
        <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 mb-2 uppercase">
          Summary
        </p>
        <p className="text-sm text-gray-700 dark:text-gray-300 p-2 rounded bg-gray-50 dark:bg-gray-700/50">
          {version.changeSummary}
        </p>
      </div>

      {/* Action Buttons */}
      {!isCurrent && (
        <button
          onClick={onRollback}
          className={`
            w-full px-4 py-2 rounded-lg font-medium text-sm
            bg-orange-600 hover:bg-orange-700 text-white
            transition-colors
          `}
        >
          ↩️ Rollback to This Version
        </button>
      )}

      {isCurrent && (
        <div className="px-4 py-2 rounded-lg bg-green-100 dark:bg-green-900/20 text-green-700 dark:text-green-200 text-sm text-center">
          ✓ This is your current version
        </div>
      )}
    </div>
  )
}

/**
 * RollbackConfirmationPanel Component
 */
function RollbackConfirmationPanel({
  version,
  reason,
  onReasonChange,
  onConfirm,
  onCancel,
}: {
  version: ThemeVersion
  reason: string
  onReasonChange: (reason: string) => void
  onConfirm: () => void
  onCancel: () => void
}): JSX.Element {
  return (
    <div className="rounded-lg border-2 border-orange-500 bg-orange-50 dark:bg-orange-900/20 p-4 space-y-3">
      <div className="flex items-start gap-3">
        <span className="text-2xl">⚠️</span>
        <div>
          <h4 className="font-semibold text-orange-900 dark:text-orange-100">Confirm Rollback</h4>
          <p className="text-sm text-orange-800 dark:text-orange-200 mt-1">
            Rolling back to {formatVersionLabel(version.versionNumber)} will replace your current theme.
            This action cannot be undone.
          </p>
        </div>
      </div>

      {/* Reason Input */}
      <textarea
        value={reason}
        onChange={(e) => onReasonChange(e.target.value)}
        placeholder="Why are you rolling back? (optional)"
        className={`
          w-full px-3 py-2 rounded-lg border border-orange-300 dark:border-orange-700
          bg-white dark:bg-gray-800 text-gray-900 dark:text-white
          placeholder-gray-500 dark:placeholder-gray-400
          focus:outline-none focus:ring-2 focus:ring-orange-500
          resize-none text-sm
        `}
        rows={2}
      />

      {/* Buttons */}
      <div className="flex gap-2">
        <button
          onClick={onConfirm}
          className={`
            flex-1 px-3 py-2 rounded-lg font-medium text-sm
            bg-orange-600 hover:bg-orange-700 text-white
            transition-colors
          `}
        >
          Yes, Rollback
        </button>
        <button
          onClick={onCancel}
          className={`
            flex-1 px-3 py-2 rounded-lg font-medium text-sm
            bg-gray-300 dark:bg-gray-700 text-gray-900 dark:text-white
            hover:bg-gray-400 dark:hover:bg-gray-600 transition-colors
          `}
        >
          Cancel
        </button>
      </div>
    </div>
  )
}

/**
 * DetailRow Component
 */
function DetailRow({ label, value }: { label: string; value: string }): JSX.Element {
  return (
    <div className="flex justify-between items-start gap-2">
      <span className="text-gray-600 dark:text-gray-400">{label}:</span>
      <span className="text-gray-900 dark:text-white font-medium text-right">{value}</span>
    </div>
  )
}

export default VersionHistory
