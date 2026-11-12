/**
 * Theme Advanced Features Utilities
 * Handles cloning, versioning, and component toggling
 */

import { ThemeData, ThemeComponent } from '@/types/theme'

/**
 * Clone/Duplicate Theme
 * Creates a copy with a new name and incremented version
 */
export function createThemeClone(
  originalTheme: ThemeData,
  newName: string,
  newDescription?: string
): Partial<ThemeData> {
  return {
    name: newName,
    description: newDescription || `Copy of ${originalTheme.name}`,
    colors: { ...originalTheme.colors },
    typography: { ...originalTheme.typography },
    identity: { ...originalTheme.identity },
    header: originalTheme.header ? { ...originalTheme.header } : undefined,
    footer: originalTheme.footer ? { ...originalTheme.footer } : undefined,
    components: originalTheme.components
      ? originalTheme.components.map((comp) => ({
          ...comp,
          id: `component-${Date.now()}-${Math.floor(Math.random() * 1_000_000)}`, // new string id (ComponentConfig.id is a string)
        }))
      : [],
    isActive: false,
    isPublished: false,
  }
}

/**
 * Generate theme version summary for history
 */
export function generateVersionSummary(
  previousTheme: ThemeData | null,
  currentTheme: ThemeData
): string {
  if (!previousTheme) {
    return `Created theme "${currentTheme.name}"`
  }

  const changes: string[] = []

  // Check what changed
  if (previousTheme.name !== currentTheme.name) {
    changes.push(`Renamed to "${currentTheme.name}"`)
  }

  if (JSON.stringify(previousTheme.colors) !== JSON.stringify(currentTheme.colors)) {
    changes.push('Updated colors')
  }

  if (JSON.stringify(previousTheme.typography) !== JSON.stringify(currentTheme.typography)) {
    changes.push('Updated typography')
  }

  if (previousTheme.components?.length !== currentTheme.components?.length) {
    changes.push(
      `Components changed (${previousTheme.components?.length || 0} → ${currentTheme.components?.length || 0})`
    )
  }

  return changes.length > 0
    ? changes.join(', ')
    : `Updated theme "${currentTheme.name}"`
}

/**
 * Toggle component visibility (enabled/disabled)
 */
export function toggleComponentEnabled(
  components: ThemeComponent[],
  componentId: string | number
): ThemeComponent[] {
  return components.map((comp) =>
    comp.id === componentId ? { ...comp, enabled: !comp.enabled } : comp
  )
}

/**
 * Enable all components
 */
export function enableAllComponents(components: ThemeComponent[]): ThemeComponent[] {
  return components.map((comp) => ({ ...comp, enabled: true }))
}

/**
 * Disable all components
 */
export function disableAllComponents(components: ThemeComponent[]): ThemeComponent[] {
  return components.map((comp) => ({ ...comp, enabled: false }))
}

/**
 * Get enabled components only
 */
export function getEnabledComponents(components: ThemeComponent[]): ThemeComponent[] {
  return components.filter((comp) => comp.enabled)
}

/**
 * Get disabled components
 */
export function getDisabledComponents(components: ThemeComponent[]): ThemeComponent[] {
  return components.filter((comp) => !comp.enabled)
}

/**
 * Count enabled vs disabled components
 */
export function getComponentStats(components: ThemeComponent[]): {
  total: number
  enabled: number
  disabled: number
} {
  const enabled = components.filter((c) => c.enabled).length
  return {
    total: components.length,
    enabled,
    disabled: components.length - enabled,
  }
}

/**
 * Format theme history date
 */
export function formatHistoryDate(dateString: string): string {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

/**
 * Get time elapsed since a date
 */
export function getTimeElapsed(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  if (seconds < 60) return 'just now'
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
  if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`
  return `${Math.floor(seconds / 604800)}w ago`
}
