'use client'

import React, { useEffect, useState } from 'react'
import { componentRegistry, resolveDynamicComponent, ComponentErrorFallback, ComponentLoadingFallback } from '../../registry/componentRegistry'
import type { SectionRendererProps } from '../../types'

/**
 * SectionRenderer: Universal dynamic component renderer
 *
 * Renders different section types based on component configuration
 * using the Component Registry system for dynamic, scalable rendering.
 *
 * This eliminates hardcoded switch statements and enables:
 * - Lazy loading of components for better performance
 * - Easy addition of new components without code changes
 * - Component versioning and aliasing
 * - Auto-discovery of available components
 *
 * Used by both dashboard (preview) and restaurant website (production)
 */
export const SectionRenderer: React.FC<SectionRendererProps> = ({
  component,
  isArabic = false,
}) => {
  const [registryReady, setRegistryReady] = useState(false)
  const [initError, setInitError] = useState<Error | null>(null)

  // Initialize component registry on first render
  useEffect(() => {
    // Check if already initialized to avoid duplicate initialization
    if (componentRegistry.checkIfInitialized?.()) {
      setRegistryReady(true)
      return
    }

    const initializeRegistry = async () => {
      try {
        await componentRegistry.initialize()
        setRegistryReady(true)
      } catch (error) {
        const err = error instanceof Error ? error : new Error(String(error))
        console.error('❌ Failed to initialize component registry:', err)
        setInitError(err)
        // Still mark as ready to avoid infinite loading state
        setRegistryReady(true)
      }
    }

    initializeRegistry()
  }, [])

  // Skip rendering if component is disabled
  if (!component.enabled) {
    return null
  }

  // Show loading fallback while registry is initializing
  if (!registryReady) {
    return <ComponentLoadingFallback />
  }

  // Show error if registry initialization failed
  if (initError) {
    return (
      <ComponentErrorFallback
        componentType={component.type}
        error={new Error(`Registry initialization failed: ${initError.message}`)}
      />
    )
  }

  // Prepare settings/config
  const settings = component.settings || component.config || {}
  const commonProps = {
    title: component.title,
    isArabic,
    settings,
    config: settings,
  }

  // Resolve component from registry
  const { Component: ResolvedComponent, isMissing } = resolveDynamicComponent(component.type, isArabic)

  // If component is not in registry, show error
  if (isMissing) {
    return (
      <ComponentErrorFallback
        componentType={component.type}
        error={new Error(`Component type "${component.type}" is not registered`)}
      />
    )
  }

  // Render the dynamically resolved component with common props
  return <ResolvedComponent {...commonProps} component={component} />
}

export default SectionRenderer
