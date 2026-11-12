/**
 * Component Registry Manager
 *
 * Runtime manager for loading, searching, and accessing components
 * from the auto-generated registry.json file.
 */

import React from 'react'
import registryData from './registry.json'
import {
  ComponentRegistry,
  ComponentRegistryEntry,
  ComponentCategory,
  IndustryType,
} from './registry-types'

type ComponentLoader = () => Promise<Record<string, unknown>>

const componentLoaders: Record<string, ComponentLoader> = {
  'hero-section': () => import('../organisms/HeroSection'),
}

const editorLoaders: Record<string, ComponentLoader> = {}

export const registerEditorLoader = (name: string, loader: ComponentLoader) => {
  editorLoaders[name] = loader
}

const resolveModuleExport = (module: Record<string, unknown>, exportName?: string) => {
  if (exportName && typeof module[exportName] === 'function') {
    return module[exportName] as React.ComponentType<any>
  }

  if (typeof module.default === 'function') {
    return module.default as React.ComponentType<any>
  }

  const firstExport = Object.values(module).find((value) => typeof value === 'function')
  if (typeof firstExport === 'function') {
    return firstExport as React.ComponentType<any>
  }

  return null
}

/**
 * Registry Manager Class
 *
 * Provides singleton access to component registry with methods for:
 * - Loading all components
 * - Filtering by category or industry
 * - Searching by name/tags
 * - Lazy loading components and editors
 */
class ComponentRegistryManager {
  private registry: ComponentRegistry
  private loadedComponents: Map<string, React.ComponentType<any>> = new Map()
  private loadedEditors: Map<string, React.ComponentType<any>> = new Map()

  constructor() {
    this.registry = registryData as unknown as ComponentRegistry
    console.log(
      `📦 Component Registry initialized with ${Object.keys(this.registry.components).length} components`
    )
  }

  /**
   * Get all registered components
   */
  getAllComponents(): ComponentRegistryEntry[] {
    return Object.values(this.registry.components) as ComponentRegistryEntry[]
  }

  /**
   * Get components by category
   */
  getComponentsByCategory(category: ComponentCategory): ComponentRegistryEntry[] {
    const componentIds = this.registry.categories[category] || []
    return componentIds
      .map((id) => this.registry.components[id])
      .filter(Boolean) as ComponentRegistryEntry[]
  }

  /**
   * Get components by industry
   */
  getComponentsByIndustry(industry: IndustryType): ComponentRegistryEntry[] {
    const componentIds = this.registry.industries[industry] || []
    return componentIds
      .map((id) => this.registry.components[id])
      .filter(Boolean) as ComponentRegistryEntry[]
  }

  /**
   * Get only draggable components (for ComponentLibrary)
   */
  getDraggableComponents(): ComponentRegistryEntry[] {
    return this.getAllComponents().filter((c) => c.isDraggable) as ComponentRegistryEntry[]
  }

  /**
   * Search components by name, tags, or description
   */
  searchComponents(query: string): ComponentRegistryEntry[] {
    const lowerQuery = query.toLowerCase()
    return this.getAllComponents().filter((component) => {
      const metadata = component as unknown as any
      return (
        metadata.name.toLowerCase().includes(lowerQuery) ||
        metadata.tags?.some((tag: string) => tag.toLowerCase().includes(lowerQuery)) ||
        metadata.description?.toLowerCase().includes(lowerQuery)
      )
    }) as ComponentRegistryEntry[]
  }

  /**
   * Lazy load component by ID
   */
  async loadComponent(componentId: string): Promise<React.ComponentType<any>> {
    if (this.loadedComponents.has(componentId)) {
      return this.loadedComponents.get(componentId)!
    }

    const entry = this.registry.components[componentId]
    if (!entry) {
      throw new Error(`Component not found: ${componentId}`)
    }

    const loader = componentLoaders[componentId]
    if (!loader) {
      throw new Error(`No static component loader registered for ${componentId}`)
    }

    try {
      const module = await loader()
      const Component = resolveModuleExport(module)

      if (!Component) {
        throw new Error(`No component export found for ${componentId}`)
      }

      this.loadedComponents.set(componentId, Component)
      return Component
    } catch (err) {
      console.error(`Failed to load component ${componentId}:`, err)
      throw err
    }
  }

  /**
   * Lazy load editor by component ID
   */
  async loadEditor(componentId: string): Promise<React.ComponentType<any> | null> {
    if (this.loadedEditors.has(componentId)) {
      return this.loadedEditors.get(componentId) || null
    }

    const entry = this.registry.components[componentId]
    if (!entry) {
      console.warn(`Component not found: ${componentId}`)
      return null
    }

    const metadata = entry as unknown as any
    if (!metadata.editorComponent) {
      return null
    }

    const loader = editorLoaders[metadata.editorComponent]
    if (!loader) {
      console.warn(`No static editor loader registered for ${metadata.editorComponent}`)
      return null
    }

    try {
      const module = await loader()
      const Editor = resolveModuleExport(module, metadata.editorComponent)

      if (!Editor) {
        console.warn(`No editor found for ${componentId}`)
        return null
      }

      this.loadedEditors.set(componentId, Editor)
      return Editor
    } catch (err) {
      console.warn(`Failed to load editor for ${componentId}:`, err)
      return null
    }
  }

  /**
   * Get component metadata
   */
  getMetadata(componentId: string): ComponentRegistryEntry | null {
    const entry = this.registry.components[componentId]
    return entry as ComponentRegistryEntry | null
  }

  /**
   * Check if component exists
   */
  hasComponent(componentId: string): boolean {
    return componentId in this.registry.components
  }

  /**
   * Get registry version
   */
  getVersion(): string {
    return this.registry.version
  }

  /**
   * Get registry generation timestamp
   */
  getGeneratedAt(): string {
    return this.registry.generatedAt
  }

  /**
   * Get statistics
   */
  getStats() {
    const components = this.getAllComponents()
    const categories = Object.keys(this.registry.categories).length
    const industries = Object.keys(this.registry.industries).length
    const draggable = components.filter((c) => (c as unknown as any).isDraggable).length

    return {
      totalComponents: components.length,
      draggableComponents: draggable,
      categories,
      industries,
      generatedAt: this.registry.generatedAt,
    }
  }
}

/**
 * Singleton instance of registry manager
 */
let instance: ComponentRegistryManager | null = null

/**
 * Get the component registry singleton
 */
export function getComponentRegistry(): ComponentRegistryManager {
  if (!instance) {
    instance = new ComponentRegistryManager()
  }
  return instance
}

/**
 * Alias for easier access
 */
export const componentRegistry = getComponentRegistry()

/**
 * React hook for accessing registry in components
 */
export function useComponentRegistry(): ComponentRegistryManager {
  return componentRegistry
}
