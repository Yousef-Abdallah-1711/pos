/**
 * Component Resolver Utility
 * Handles component version resolution and compatibility checking
 */

import { ComponentDefinition, ComponentVersionInfo, componentRegistry } from '../registry/types'

/**
 * Component resolver for smart version resolution
 */
export class ComponentResolver {
  /**
   * Resolve component to specific version
   * Falls back to latest compatible version if exact version not found
   */
  static resolveComponent(
    componentId: string,
    targetVersion?: string
  ): ComponentDefinition | null {
    // If version specified, try to get exact match first
    if (targetVersion) {
      const exact = componentRegistry.getComponent(componentId, targetVersion)
      if (exact) {
        return exact
      }

      // Try to find compatible version
      const compatible = this.findCompatibleVersion(componentId, targetVersion)
      if (compatible) {
        console.warn(
          `Version ${targetVersion} of ${componentId} not found, using ${compatible.version}`
        )
        return compatible
      }
    }

    // Fall back to latest version
    return this.getLatestVersion(componentId)
  }

  /**
   * Find compatible version based on semantic versioning
   */
  static findCompatibleVersion(
    componentId: string,
    targetVersion: string
  ): ComponentDefinition | null {
    const versions = componentRegistry.getComponentVersions(componentId)

    if (versions.length === 0) {
      return null
    }

    // Parse target version
    const [majorTarget, minorTarget, patchTarget] = this.parseVersion(targetVersion)

    // Find best compatible match
    let bestMatch: ComponentDefinition | null = null
    let bestMatchScore = -1

    versions.forEach(version => {
      const [major, minor, patch] = this.parseVersion(version)
      const component = componentRegistry.getComponent(componentId, version)

      if (!component) return

      // Check compatibility: major version must match
      if (major === majorTarget) {
        // Prefer exact version, then same minor, then same major
        let score = 0
        if (major === majorTarget && minor === minorTarget && patch === patchTarget) {
          score = 1000 // Exact match
        } else if (major === majorTarget && minor === minorTarget) {
          score = 100 // Same major.minor
        } else if (major === majorTarget && minor >= minorTarget) {
          score = 50 // Same major, compatible minor
        } else {
          score = 10 // Same major only
        }

        if (score > bestMatchScore) {
          bestMatchScore = score
          bestMatch = component
        }
      }
    })

    return bestMatch
  }

  /**
   * Get latest version of a component
   */
  static getLatestVersion(componentId: string): ComponentDefinition | null {
    const versions = componentRegistry.getComponentVersions(componentId)

    if (versions.length === 0) {
      return null
    }

    // Versions are sorted in descending order
    return componentRegistry.getComponent(componentId, versions[0]) || null
  }

  /**
   * Get version information for a component
   */
  static getVersionInfo(componentId: string): ComponentVersionInfo | null {
    const versions = componentRegistry.getComponentVersions(componentId)

    if (versions.length === 0) {
      return null
    }

    const latest = this.getLatestVersion(componentId)
    if (!latest) {
      return null
    }

    return {
      componentId,
      latestVersion: versions[0],
      allVersions: [...versions],
      deprecated: latest.manifest.deprecated,
      latestDeprecatedAt: latest.manifest.deprecated ? new Date() : undefined,
    }
  }

  /**
   * Check if version is compatible with requirement
   * Returns true if targetVersion matches or is compatible with requiredVersion
   */
  static isCompatible(
    requiredVersion: string,
    targetVersion: string
  ): boolean {
    const [majorReq, minorReq] = this.parseVersion(requiredVersion)
    const [majorTgt, minorTgt] = this.parseVersion(targetVersion)

    // Major version must match
    if (majorReq !== majorTgt) {
      return false
    }

    // Target minor must be >= required minor (backward compatible)
    return minorTgt >= minorReq
  }

  /**
   * Parse semantic version string to numbers
   * Handles: "1.0.0", "1.0", "1"
   */
  private static parseVersion(versionString: string): [number, number, number] {
    const parts = versionString.split('.').map(p => parseInt(p || '0', 10))
    return [parts[0] || 0, parts[1] || 0, parts[2] || 0]
  }

  /**
   * Get migration path from source version to target version
   */
  static getMigrationPath(
    componentId: string,
    fromVersion: string,
    toVersion: string
  ): string[] {
    const versions = componentRegistry.getComponentVersions(componentId)

    if (!versions.includes(fromVersion) || !versions.includes(toVersion)) {
      return []
    }

    // Simple linear path (in production, could be more complex)
    const path: string[] = []
    let currentVersion = fromVersion

    while (currentVersion !== toVersion) {
      const currentIndex = versions.indexOf(currentVersion)
      if (currentIndex === -1) break

      // Get next version in sequence
      const nextVersion =
        versions[Math.min(currentIndex + 1, versions.length - 1)]
      if (nextVersion === currentVersion) break // Avoid infinite loop

      path.push(nextVersion)
      currentVersion = nextVersion
    }

    return path
  }

  /**
   * Check if component version is deprecated
   */
  static isDeprecated(componentId: string, version: string): boolean {
    const component = componentRegistry.getComponent(componentId, version)
    return component?.manifest.deprecated ?? false
  }

  /**
   * Get all versions of a component
   */
  static getAllVersions(componentId: string): string[] {
    return componentRegistry.getComponentVersions(componentId)
  }

  /**
   * Get components that are compatible with specific version requirement
   */
  static findCompatibleComponents(
    componentId: string,
    requiredVersion: string
  ): ComponentDefinition[] {
    const versions = componentRegistry.getComponentVersions(componentId)
    const compatible: ComponentDefinition[] = []

    versions.forEach(version => {
      if (this.isCompatible(requiredVersion, version)) {
        const component = componentRegistry.getComponent(componentId, version)
        if (component) {
          compatible.push(component)
        }
      }
    })

    return compatible
  }
}

/**
 * Version constraint validator
 * Validates if a version string matches a constraint pattern
 */
export class VersionConstraint {
  /**
   * Check if version matches constraint
   * Supports: "1.0.0", "^1.0.0", "~1.0.0", ">=1.0.0", "<=1.0.0", "1.x"
   */
  static matches(version: string, constraint: string): boolean {
    if (constraint === '*' || constraint === 'latest') {
      return true
    }

    // Exact match
    if (constraint === version) {
      return true
    }

    // Caret: ^1.0.0 - compatible with 1.x.x
    if (constraint.startsWith('^')) {
      const constraintVersion = constraint.substring(1)
      const [majReq, minReq] = this.parseVersion(constraintVersion)
      const [majVer, minVer] = this.parseVersion(version)

      return majVer === majReq && minVer >= minReq
    }

    // Tilde: ~1.0.0 - compatible with 1.0.x
    if (constraint.startsWith('~')) {
      const constraintVersion = constraint.substring(1)
      const [majReq, minReq, patchReq] = this.parseVersion(constraintVersion)
      const [majVer, minVer, patchVer] = this.parseVersion(version)

      return majVer === majReq && minVer === minReq && patchVer >= patchReq
    }

    // Greater than or equal
    if (constraint.startsWith('>=')) {
      const constraintVersion = constraint.substring(2)
      return this.compareVersions(version, constraintVersion) >= 0
    }

    // Less than or equal
    if (constraint.startsWith('<=')) {
      const constraintVersion = constraint.substring(2)
      return this.compareVersions(version, constraintVersion) <= 0
    }

    // Wildcard: 1.x
    if (constraint.includes('x')) {
      const [majReq] = this.parseVersion(constraint)
      const [majVer] = this.parseVersion(version)
      return majVer === majReq
    }

    return false
  }

  /**
   * Compare two versions
   * Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
   */
  private static compareVersions(v1: string, v2: string): number {
    const [maj1, min1, patch1] = this.parseVersion(v1)
    const [maj2, min2, patch2] = this.parseVersion(v2)

    if (maj1 !== maj2) return maj1 > maj2 ? 1 : -1
    if (min1 !== min2) return min1 > min2 ? 1 : -1
    if (patch1 !== patch2) return patch1 > patch2 ? 1 : -1
    return 0
  }

  /**
   * Parse semantic version string
   */
  private static parseVersion(versionString: string): [number, number, number] {
    const parts = versionString.split('.').map(p => {
      const num = parseInt(p.replace(/\D/g, ''), 10)
      return isNaN(num) ? 0 : num
    })
    return [parts[0] || 0, parts[1] || 0, parts[2] || 0]
  }
}
