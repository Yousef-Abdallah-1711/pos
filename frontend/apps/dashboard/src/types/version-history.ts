/**
 * Version History Types
 * Types for theme versioning and rollback functionality
 */

import type { ThemeData as Theme } from './theme'

/**
 * Type of change made to a theme
 */
export type ChangeType = 'create' | 'update' | 'rollback' | 'activate' | 'publish' | 'duplicate'

/**
 * Represents a single version/snapshot of a theme
 */
export interface ThemeVersion {
  id: string | number
  themeId: string | number
  versionNumber: number
  changes: string // JSON string of changes
  changeSummary: string
  changeType: ChangeType
  authorName?: string
  authorEmail?: string
  createdBy?: string | number
  createdAt: string
  isCurrent: boolean
  rollbackReason?: string
}

/**
 * Represents a complete theme snapshot at a point in time
 */
export interface ThemeSnapshot {
  id: string | number
  themeId: string | number
  versionNumber: number
  themeSnapshot: Partial<Theme>
  componentsSnapshot: any[]
  snapshotSizeBytes?: number
  createdAt: string
}

/**
 * Version history entry with snapshot
 */
export interface VersionHistoryEntry {
  version: ThemeVersion
  snapshot?: ThemeSnapshot
}

/**
 * Change details for a specific property
 */
export interface PropertyChange {
  path: string
  previousValue: any
  newValue: any
  type: 'modified' | 'added' | 'removed'
}

/**
 * Detailed changes breakdown
 */
export interface ChangesBreakdown {
  colors?: PropertyChange[]
  typography?: PropertyChange[]
  identity?: PropertyChange[]
  components?: PropertyChange[]
  other?: PropertyChange[]
}

/**
 * Version comparison result
 */
export interface VersionComparison {
  fromVersion: ThemeVersion
  toVersion: ThemeVersion
  changes: ChangesBreakdown
  changeCount: number
  percentageDifference: number
}

/**
 * Rollback options
 */
export interface RollbackOptions {
  targetVersionNumber: number
  reason?: string
  confirmationRequired?: boolean
}

/**
 * Version history state
 */
export interface VersionHistoryState {
  versions: ThemeVersion[]
  snapshots: Map<number, ThemeSnapshot>
  currentVersion: ThemeVersion | null
  isLoading: boolean
  error: string | null
  selectedVersion: ThemeVersion | null
  comparisonVersion?: ThemeVersion
}

/**
 * Version history action types
 */
export type VersionHistoryAction =
  | { type: 'SET_VERSIONS'; payload: ThemeVersion[] }
  | { type: 'SET_SNAPSHOTS'; payload: Map<number, ThemeSnapshot> }
  | { type: 'SET_CURRENT_VERSION'; payload: ThemeVersion | null }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'SELECT_VERSION'; payload: ThemeVersion | null }
  | { type: 'SET_COMPARISON'; payload: ThemeVersion | undefined }
  | { type: 'RESET' }

/**
 * Hook for version history data
 */
export interface UseVersionHistoryOptions {
  themeId?: string | number
  includeSnapshots?: boolean
  limit?: number
}

/**
 * Rollback result
 */
export interface RollbackResult {
  success: boolean
  message: string
  newVersion?: ThemeVersion
  previousVersion?: ThemeVersion
  rollbackTime?: number
}

/**
 * Change type icons and labels
 */
export const CHANGE_TYPE_CONFIG: Record<ChangeType, { icon: string; label: string; color: string }> = {
  create: { icon: '✨', label: 'Created', color: 'green' },
  update: { icon: '✏️', label: 'Updated', color: 'blue' },
  rollback: { icon: '↩️', label: 'Rolled Back', color: 'orange' },
  activate: { icon: '🚀', label: 'Activated', color: 'purple' },
  publish: { icon: '📢', label: 'Published', color: 'cyan' },
  duplicate: { icon: '📋', label: 'Duplicated', color: 'gray' },
}

/**
 * Format version change type
 */
export function getChangeTypeConfig(type: ChangeType) {
  return CHANGE_TYPE_CONFIG[type] || { icon: '📝', label: type, color: 'gray' }
}

/**
 * Parse changes JSON string to object
 */
export function parseChanges(changesJson: string): ChangesBreakdown {
  try {
    return JSON.parse(changesJson)
  } catch {
    return {}
  }
}

/**
 * Calculate version difference percentage
 */
export function calculateDifference(from: any, to: any): number {
  if (!from || !to) return 100
  const diff = JSON.stringify(from) !== JSON.stringify(to)
  return diff ? 100 : 0
}

/**
 * Format version label
 */
export function formatVersionLabel(versionNumber: number, isCurrent?: boolean): string {
  const label = `v${versionNumber}`
  return isCurrent ? `${label} (Current)` : label
}

/**
 * Time since version
 */
export function timeSinceVersion(versionDate: string): string {
  const date = new Date(versionDate)
  const now = new Date()
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  if (seconds < 60) return `${seconds}s ago`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
  if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`
  return date.toLocaleDateString()
}
