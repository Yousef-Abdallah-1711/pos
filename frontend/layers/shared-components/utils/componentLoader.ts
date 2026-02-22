/**
 * Component Loader Utility
 * Handles dynamic loading of components from the registry
 */

import { ComponentDefinition, ComponentNotFoundError, componentRegistry } from '../registry/types'

/**
 * Component loader with error handling and fallback
 */
export class ComponentLoader {
  private static loadingPromises: Map<string, Promise<ComponentDefinition>> = new Map()

  /**
   * Load a component by ID and optional version
   * Async operation that resolves to the component definition
   */
  static async loadComponent(
    componentId: string,
    version?: string
  ): Promise<ComponentDefinition> {
    const key = version ? `${componentId}@${version}` : componentId

    // Check if already loading to avoid duplicate requests
    if (this.loadingPromises.has(key)) {
      return this.loadingPromises.get(key)!
    }

    // Create loading promise
    const loadingPromise = this.performLoad(componentId, version)
    this.loadingPromises.set(key, loadingPromise)

    try {
      const component = await loadingPromise
      return component
    } finally {
      // Clean up loading promise after completion
      setTimeout(() => {
        this.loadingPromises.delete(key)
      }, 100)
    }
  }

  /**
   * Internal method to perform the actual loading
   */
  private static async performLoad(
    componentId: string,
    version?: string
  ): Promise<ComponentDefinition> {
    // Simulate async loading (in real app, this would lazy-load modules)
    return new Promise((resolve, reject) => {
      // Small delay to simulate loading
      setTimeout(() => {
        const component = componentRegistry.getComponent(componentId, version)

        if (!component) {
          reject(
            new ComponentNotFoundError(componentId, version)
          )
          return
        }

        resolve(component)
      }, 50)
    })
  }

  /**
   * Load component asynchronously with timeout
   */
  static async loadComponentWithTimeout(
    componentId: string,
    version?: string,
    timeout: number = 5000
  ): Promise<ComponentDefinition> {
    return Promise.race([
      this.loadComponent(componentId, version),
      new Promise<ComponentDefinition>((_, reject) =>
        setTimeout(
          () => reject(new Error(`Component loading timeout for ${componentId}`)),
          timeout
        )
      ),
    ])
  }

  /**
   * Batch load multiple components
   */
  static async loadComponentsBatch(
    componentIds: Array<{ id: string; version?: string }>
  ): Promise<(ComponentDefinition | null)[]> {
    const promises = componentIds.map(({ id, version }) =>
      this.loadComponent(id, version)
        .then(comp => comp)
        .catch(() => null)
    )
    return Promise.all(promises)
  }

  /**
   * Preload critical components
   */
  static preloadComponents(componentIds: string[]): void {
    componentIds.forEach(id => {
      this.loadComponent(id).catch(err => {
        console.warn(`Failed to preload component ${id}:`, err)
      })
    })
  }

  /**
   * Load component with fallback
   */
  static async loadComponentWithFallback(
    componentId: string,
    fallbackId?: string,
    version?: string
  ): Promise<ComponentDefinition> {
    try {
      return await this.loadComponent(componentId, version)
    } catch (error) {
      if (fallbackId) {
        console.warn(
          `Failed to load ${componentId}, trying fallback ${fallbackId}`
        )
        return await this.loadComponent(fallbackId)
      }
      throw error
    }
  }

  /**
   * Clear loading cache
   */
  static clearCache(): void {
    this.loadingPromises.clear()
    componentRegistry.clearCache()
  }

  /**
   * Get loading status
   */
  static getLoadingStatus(): {
    loadingCount: number
    cachedComponents: number
  } {
    const cacheStats = componentRegistry.getCacheStats()
    return {
      loadingCount: this.loadingPromises.size,
      cachedComponents: cacheStats.cachedComponents,
    }
  }
}

/**
 * React Hook for loading components
 * Can be used in functional components
 */
export function useComponentLoader(
  componentId: string,
  version?: string
) {
  const [component, setComponent] = React.useState<ComponentDefinition | null>(null)
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState<Error | null>(null)

  React.useEffect(() => {
    let mounted = true

    const load = async () => {
      try {
        setLoading(true)
        const comp = await ComponentLoader.loadComponent(componentId, version)
        if (mounted) {
          setComponent(comp)
          setError(null)
        }
      } catch (err) {
        if (mounted) {
          setError(err instanceof Error ? err : new Error(String(err)))
          setComponent(null)
        }
      } finally {
        if (mounted) {
          setLoading(false)
        }
      }
    }

    load()

    return () => {
      mounted = false
    }
  }, [componentId, version])

  return { component, loading, error }
}

// Add React import for hook
import React from 'react'
