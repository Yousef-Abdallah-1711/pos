
import type { JSX } from 'react'
/**
 * VersionComparison Component
 * Displays side-by-side comparison of two theme versions
 */

'use client'

import React from 'react'
import type { ThemeVersion, VersionComparison as VersionComparisonType } from '@/types/version-history'
import { formatVersionLabel } from '@/types/version-history'

interface VersionComparisonProps {
  /** From version */
  fromVersion: ThemeVersion
  /** To version */
  toVersion: ThemeVersion
  /** Comparison data */
  comparison?: VersionComparisonType
  /** Custom CSS classes */
  className?: string
}

/**
 * VersionComparison Component
 */
export function VersionComparison({
  fromVersion,
  toVersion,
  comparison,
  className = '',
}: VersionComparisonProps): JSX.Element {
  const progressWidth = comparison ? Math.min(comparison.percentageDifference, 100) : 0

  return (
    <div className={`version-comparison rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden ${className}`}>
      {/* Header */}
      <div className="bg-gray-50 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Version Comparison</h3>
          {comparison && (
            <div className="flex items-center gap-4">
              <span className="text-sm text-gray-600 dark:text-gray-400">
                {comparison.changeCount} change{comparison.changeCount !== 1 ? 's' : ''}
              </span>
              <div className="w-24 h-2 bg-gray-300 dark:bg-gray-600 rounded-full overflow-hidden">
                <div className="comparison-progress-fill h-full bg-blue-500 transition-all" />
              </div>
            </div>
          )}
        </div>

        {/* Version Labels */}
        <div className="grid grid-cols-2 gap-4">
          <div className="flex items-center gap-2">
            <span className="text-2xl">←</span>
            <div>
              <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase">From</p>
              <p className="font-mono font-semibold text-gray-900 dark:text-white">
                {formatVersionLabel(fromVersion.versionNumber)}
              </p>
            </div>
          </div>
          <div className="flex items-center justify-end gap-2">
            <div className="text-right">
              <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase">To</p>
              <p className="font-mono font-semibold text-gray-900 dark:text-white">
                {formatVersionLabel(toVersion.versionNumber)}
              </p>
            </div>
            <span className="text-2xl">→</span>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="divide-y divide-gray-200 dark:divide-gray-700">
        {/* Colors Section */}
        {comparison?.changes.colors && comparison.changes.colors.length > 0 && (
          <ComparisonSection
            title="🎨 Colors"
            changes={comparison.changes.colors}
          />
        )}

        {/* Typography Section */}
        {comparison?.changes.typography && comparison.changes.typography.length > 0 && (
          <ComparisonSection
            title="✏️ Typography"
            changes={comparison.changes.typography}
          />
        )}

        {/* Identity Section */}
        {comparison?.changes.identity && comparison.changes.identity.length > 0 && (
          <ComparisonSection
            title="🏷️ Identity"
            changes={comparison.changes.identity}
          />
        )}

        {/* Components Section */}
        {comparison?.changes.components && comparison.changes.components.length > 0 && (
          <ComparisonSection
            title="📦 Components"
            changes={comparison.changes.components}
          />
        )}

        {/* Other Changes */}
        {comparison?.changes.other && comparison.changes.other.length > 0 && (
          <ComparisonSection
            title="📋 Other"
            changes={comparison.changes.other}
          />
        )}

        {/* No Changes */}
        {(!comparison || Object.keys(comparison.changes).length === 0) && (
          <div className="px-6 py-8 text-center">
            <p className="text-gray-600 dark:text-gray-400">No differences found between versions</p>
          </div>
        )}
      </div>

      <style jsx>{`
        .comparison-progress-fill {
          width: ${progressWidth}%;
        }
      `}</style>
    </div>
  )
}

/**
 * ComparisonSection Component
 */
function ComparisonSection({
  title,
  changes,
}: {
  title: string
  changes: any[]
}): JSX.Element {
  return (
    <div>
      <div className="bg-gray-50 dark:bg-gray-900 px-6 py-3 border-b border-gray-200 dark:border-gray-700">
        <h4 className="font-semibold text-gray-900 dark:text-white text-sm">{title}</h4>
      </div>

      <div className="divide-y divide-gray-100 dark:divide-gray-700">
        {changes.map((change, index) => (
          <ChangeRow key={index} change={change} />
        ))}
      </div>
    </div>
  )
}

/**
 * ChangeRow Component
 */
function ChangeRow({ change }: { change: any }): JSX.Element {
  const typeIcons = {
    modified: '🔄',
    added: '✨',
    removed: '🗑️',
  }

  const typeColors = {
    modified: 'blue',
    added: 'green',
    removed: 'red',
  }

  const icon = typeIcons[change.type as keyof typeof typeIcons] || '📝'
  const color = typeColors[change.type as keyof typeof typeColors] || 'gray'

  return (
    <div className="px-6 py-4 hover:bg-gray-50 dark:hover:bg-gray-700/30 transition-colors">
      <div className="flex items-start gap-3 mb-2">
        <span className="text-lg">{icon}</span>
        <div className="flex-1 min-w-0">
          <p className="font-mono text-sm font-semibold text-gray-900 dark:text-white break-all">
            {change.path}
          </p>
          <p className="text-xs text-gray-600 dark:text-gray-400 mt-1 capitalize">
            {change.type}
          </p>
        </div>
      </div>

      {/* Values */}
      <div className="grid grid-cols-2 gap-3 ml-7 mt-2">
        {/* Previous Value */}
        {change.previousValue !== undefined && change.type !== 'added' && (
          <div>
            <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 mb-1">
              Before
            </p>
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-2">
              <code className="text-xs text-gray-900 dark:text-gray-200 break-all">
                {JSON.stringify(change.previousValue)}
              </code>
            </div>
          </div>
        )}

        {/* New Value */}
        {change.newValue !== undefined && change.type !== 'removed' && (
          <div>
            <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 mb-1">
              After
            </p>
            <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded p-2">
              <code className="text-xs text-gray-900 dark:text-gray-200 break-all">
                {JSON.stringify(change.newValue)}
              </code>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default VersionComparison