import React, { lazy, Suspense } from 'react'
import type { SectionRendererProps } from '../types'
import { ComponentCache, type CacheStats } from './componentCache'
import { componentDiscovery, type DiscoveredComponent } from './componentDiscovery'

/**
 * Component Registry System
 *
 * Provides dynamic, lazy-loaded component registration and resolution.
 * Eliminates the need for hardcoded switch statements in SectionRenderer.
 *
 * Features:
 * - Lazy-loading with React.lazy() for code splitting
 * - LRU in-memory cache with configurable size limits
 * - Auto-discovery of components from sections folder
 * - Version support for component variants
 * - Fallback components for missing types
 * - Advanced memory management and cache statistics
 */

// Component metadata interface
export interface ComponentDefinition {
  id: string
  name: string
  aliases: string[] // Alternative names for the component
  version: string
  component: React.FC<any> | React.LazyExoticComponent<React.FC<any>>
  mockData?: Record<string, any>
  description?: string
}

// Global registry cache
class ComponentRegistryManager {
  private registry: Map<string, ComponentDefinition> = new Map()
  private componentCache: ComponentCache<React.FC<any> | React.LazyExoticComponent<React.FC<any>>>
  private isInitialized = false
  private discoveredComponents: DiscoveredComponent[] = []
  private useAutoDiscovery = true // Default: use auto-discovery

  constructor(useAutoDiscovery: boolean = true) {
    // Initialize cache with max 100 components in memory
    this.componentCache = new ComponentCache(100)
    this.useAutoDiscovery = useAutoDiscovery
  }

  /**
   * Initialize the registry with all available components
   * This should be called once at app startup
   * Supports auto-discovery mode
   */
  async initialize(useDiscovery: boolean = true): Promise<void> {
    if (this.isInitialized) return

    console.log('📦 Initializing Component Registry...')
    console.log('💾 Component cache system active (Max: 100 components, LRU eviction enabled)')

    if (useDiscovery && this.useAutoDiscovery) {
      console.log('🔍 Auto-discovery enabled - scanning for components...')
      await this.initializeWithDiscovery()
    } else {
      await this.initializeWithManualRegistration()
    }

    this.isInitialized = true
    console.log(`✅ Component Registry initialized with ${this.getAllComponents().length} unique components (${this.registry.size} total with aliases)`)
  }

  /**
   * Initialize using auto-discovery system
   */
  private async initializeWithDiscovery(): Promise<void> {
    try {
      const discovered = await componentDiscovery.discoverComponents()
      this.discoveredComponents = discovered

      // Register components dynamically based on discovery
      const componentMap: Record<string, () => Promise<any>> = {
        hero: () => import('../sections/marketing/HeroSection').then(m => ({ default: m.HeroSection })),
        products: () => import('../sections/commerce/ProductsSection').then(m => ({ default: m.ProductsSection })),
        categories: () => import('../sections/commerce/ProductsSection').then(m => ({ default: m.ProductsSection })),
        gallery: () => import('../sections/commerce/ProductsSection').then(m => ({ default: m.ProductsSection })),
        why_us: () => import('../sections/content/WhyChooseUsSection').then(m => ({ default: m.WhyChooseUsSection })),
        contact: () => import('../sections/forms/ContactSection').then(m => ({ default: m.ContactSection })),
        testimonials: () => import('../sections/marketing/TestimonialsSection').then(m => ({ default: m.TestimonialsSection })),
        cta: () => import('../sections/marketing/CTASection').then(m => ({ default: m.CTASection })),
        custom: () => import('../sections/custom/CustomSection').then(m => ({ default: m.CustomSection })),
      }

      for (const discovered of this.discoveredComponents) {
        const importer = componentMap[discovered.id]
        if (!importer) {
          console.warn(`⚠️ No importer found for discovered component: ${discovered.id}`)
          continue
        }

        this.registerComponent({
          id: discovered.id,
          name: discovered.name,
          aliases: discovered.aliases,
          version: discovered.version,
          description: discovered.description,
          component: lazy(importer),
        })
      }

      console.log(`✅ Auto-discovery: Registered ${this.discoveredComponents.length} discovered components`)
    } catch (error) {
      console.error('❌ Auto-discovery failed:', error)
      console.log('Falling back to manual registration...')
      await this.initializeWithManualRegistration()
    }
  }

  /**
   * Initialize with manual component registration (fallback)
   */
  private async initializeWithManualRegistration(): Promise<void> {
    // Register all components with their aliases
    this.registerComponent({
      id: 'hero',
      name: 'Hero Section',
      aliases: ['hero-section', 'hero_section'],
      version: '1.0.0',
      component: lazy(() => import('../sections/marketing/HeroSection').then(m => ({ default: m.HeroSection }))),
      description: 'Large hero banner with title, subtitle, description, and CTA button',
    })

    this.registerComponent({
      id: 'products',
      name: 'Products Section',
      aliases: [
        'products-section',
        'products_section',
        'featured_products',
        'product_grid',
        'categories',
        'categories-section',
        'browse_menu',
        'gallery',
        'gallery-section',
        'food_gallery',
      ],
      version: '1.0.0',
      component: lazy(() => import('../sections/commerce/ProductsSection').then(m => ({ default: m.ProductsSection }))),
      description: 'Display products in grid or list layout',
    })

    this.registerComponent({
      id: 'categories',
      name: 'Categories Section',
      aliases: ['categories-section', 'categories_section', 'browse_menu'],
      version: '1.0.0',
      component: lazy(() => import('../sections/commerce/ProductsSection').then(m => ({ default: m.ProductsSection }))),
      description: 'Compatibility mapping for categories sections',
    })

    this.registerComponent({
      id: 'gallery',
      name: 'Gallery Section',
      aliases: ['gallery-section', 'gallery_section', 'food_gallery'],
      version: '1.0.0',
      component: lazy(() => import('../sections/commerce/ProductsSection').then(m => ({ default: m.ProductsSection }))),
      description: 'Compatibility mapping for gallery sections',
    })

    this.registerComponent({
      id: 'why_us',
      name: 'Why Choose Us Section',
      aliases: ['why-us', 'why_choose_us', 'why-choose-us', 'features'],
      version: '1.0.0',
      component: lazy(() => import('../sections/content/WhyChooseUsSection').then(m => ({ default: m.WhyChooseUsSection }))),
      description: 'Highlight features or benefits with multiple columns',
    })

    this.registerComponent({
      id: 'contact',
      name: 'Contact Section',
      aliases: ['contact-section', 'contact_section', 'contact_us'],
      version: '1.0.0',
      component: lazy(() => import('../sections/forms/ContactSection').then(m => ({ default: m.ContactSection }))),
      description: 'Contact information with optional form and map',
    })

    this.registerComponent({
      id: 'testimonials',
      name: 'Testimonials Section',
      aliases: ['testimonials-section', 'testimonials_section', 'reviews'],
      version: '1.0.0',
      component: lazy(() => import('../sections/marketing/TestimonialsSection').then(m => ({ default: m.TestimonialsSection }))),
      description: 'Display customer testimonials in grid or carousel',
    })

    this.registerComponent({
      id: 'cta',
      name: 'Call To Action Section',
      aliases: ['cta-section', 'cta_section', 'call_to_action'],
      version: '1.0.0',
      component: lazy(() => import('../sections/marketing/CTASection').then(m => ({ default: m.CTASection }))),
      description: 'Prominent call-to-action section with button',
    })

    this.registerComponent({
      id: 'custom',
      name: 'Custom Section',
      aliases: ['custom-section', 'custom_section', 'html'],
      version: '1.0.0',
      component: lazy(() => import('../sections/custom/CustomSection').then(m => ({ default: m.CustomSection }))),
      description: 'Custom HTML or markdown content',
    })

    console.log('✅ Manual registration: Registered 7 components')
  }

  /**
   * Register a new component in the registry
   */
  private registerComponent(definition: ComponentDefinition): void {
    this.registry.set(definition.id, definition)

    // Also register all aliases pointing to the same component
    definition.aliases.forEach(alias => {
      this.registry.set(alias.toLowerCase(), definition)
    })
  }

  /**
   * Resolve a component by type name
   * Supports exact match and alias matching
   */
  getComponent(componentType: string): ComponentDefinition | undefined {
    return this.registry.get(componentType.toLowerCase())
  }

  /**
   * Get all registered components (without aliases)
   */
  getAllComponents(): ComponentDefinition[] {
    const seen = new Set<string>()
    const components: ComponentDefinition[] = []

    for (const def of this.registry.values()) {
      if (!seen.has(def.id)) {
        seen.add(def.id)
        components.push(def)
      }
    }

    return components
  }

  /**
   * Check if a component type is registered
   */
  hasComponent(componentType: string): boolean {
    return this.registry.has(componentType.toLowerCase())
  }

  /**
   * Check if registry is initialized
   */
  checkIfInitialized(): boolean {
    return this.isInitialized
  }

  /**
   * Get registry statistics for debugging
   */
  getStats() {
    const cacheStats = this.componentCache.getStats()
    return {
      totalRegistered: this.registry.size,
      uniqueComponents: this.getAllComponents().length,
      isInitialized: this.isInitialized,
      cache: {
        loaded: cacheStats.totalLoaded,
        maxSize: cacheStats.maxSize,
        hitRate: `${cacheStats.hitRate}%`,
        missRate: `${cacheStats.missRate}%`,
        evictions: cacheStats.evictions,
        memoryUsage: cacheStats.memoryUsage,
      },
      components: this.getAllComponents().map(c => ({
        id: c.id,
        name: c.name,
        version: c.version,
        aliases: c.aliases,
      })),
    }
  }

  /**
   * Get detailed cache statistics
   */
  getCacheStats(): CacheStats {
    return this.componentCache.getStats()
  }

  /**
   * Get cache debug information
   */
  getCacheDebug() {
    return this.componentCache.debug()
  }

  /**
   * Clear the component cache (useful for development/testing)
   */
  clearCache(): void {
    this.componentCache.clear()
    console.log('🗑️ Component registry cache cleared')
  }

  /**
   * Get discovered components
   */
  getDiscoveredComponents(): DiscoveredComponent[] {
    return this.discoveredComponents
  }

  /**
   * Get discovery statistics
   */
  getDiscoveryStats() {
    return {
      isDiscovered: this.discoveredComponents.length > 0,
      totalDiscovered: this.discoveredComponents.length,
      components: this.discoveredComponents.map(c => ({
        id: c.id,
        filename: c.filename,
        name: c.name,
        description: c.description,
        aliases: c.aliases,
        version: c.version,
      })),
    }
  }

  /**
   * Enable or disable auto-discovery mode
   */
  setAutoDiscovery(enabled: boolean): void {
    this.useAutoDiscovery = enabled
    console.log(`Auto-discovery ${enabled ? 'enabled' : 'disabled'}`)
  }
}

// Singleton instance
export const componentRegistry = new ComponentRegistryManager()

/**
 * Loading fallback component shown while sections are lazy-loading
 */
export const ComponentLoadingFallback: React.FC = () =>
  React.createElement(
    'div',
    { className: 'min-h-96 flex items-center justify-center bg-gray-50 animate-pulse' },
    React.createElement(
      'div',
      { className: 'text-center' },
      React.createElement('div', { className: 'h-8 bg-gray-200 rounded w-32 mx-auto mb-2' }),
      React.createElement('div', { className: 'h-4 bg-gray-200 rounded w-48 mx-auto' })
    )
  )

/**
 * Error fallback component shown when a section fails to load
 */
export const ComponentErrorFallback: React.FC<{ componentType: string; error?: Error }> = ({
  componentType,
  error,
}) =>
  React.createElement(
    'div',
    { className: 'min-h-96 flex items-center justify-center bg-red-50 border-2 border-red-200 rounded-lg' },
    React.createElement(
      'div',
      { className: 'text-center' },
      React.createElement('p', { className: 'text-red-600 font-semibold' }, 'Failed to load component'),
      React.createElement('p', { className: 'text-red-500 text-sm mt-2' }, componentType),
      error && React.createElement('p', { className: 'text-red-400 text-xs mt-2 font-mono' }, error.message)
    )
  )

/**
 * Dynamic component resolver
 * Returns the appropriate component with lazy loading and error boundary support
 */
export function resolveDynamicComponent(
  componentType: string,
  isArabic: boolean = false
): { Component: React.FC<SectionRendererProps>; isMissing: boolean } {
  const definition = componentRegistry.getComponent(componentType)

  if (!definition) {
    console.warn(`⚠️ Component type '${componentType}' not found in registry`)
    return {
      Component: ({ component }) =>
        React.createElement(ComponentErrorFallback, {
          componentType: component.type,
        }),
      isMissing: true,
    }
  }

  // Return a wrapper component that handles lazy loading
  const RenderedComponent: React.FC<SectionRendererProps> = (props) => {
    const Component = definition.component

    return React.createElement(
      Suspense,
      { fallback: React.createElement(ComponentLoadingFallback) },
      React.createElement(Component, { ...props })
    )
  }

  return {
    Component: RenderedComponent,
    isMissing: false,
  }
}

export default componentRegistry
