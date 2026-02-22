/**
 * Component Registration
 * Central place to register all shared components
 * This file should be called during app initialization
 */

import { ComponentRegistry } from '../registry/types'
import { definition as heroDefinition } from './hero/v1/Hero'
import { definition as productsDefinition } from './products/v1/Products'
import { definition as testimonialsDefinition } from './testimonials/v1/Testimonials'
import { definition as contactDefinition } from './contact/v1/Contact'
import { definition as ctaDefinition } from './cta/v1/CTA'
import { definition as whyUsDefinition } from './why_us/v1/WhyUs'

/**
 * Register all default components in the registry
 * Call this once during app startup
 */
export function registerDefaultComponents() {
  const registry = ComponentRegistry.getInstance()

  // Register all 6 components
  registry.registerComponent(heroDefinition)
  registry.registerComponent(productsDefinition)
  registry.registerComponent(testimonialsDefinition)
  registry.registerComponent(contactDefinition)
  registry.registerComponent(ctaDefinition)
  registry.registerComponent(whyUsDefinition)

  console.log('✅ Components registered successfully')
  const stats = registry.getComponentStats()
  console.log(`Total components: ${stats.totalComponents}`)
  console.log(`Components: Hero, Products, Testimonials, Contact, CTA, Why Us`)
}

/**
 * Export all component definitions for direct access
 */
export { heroDefinition, productsDefinition, testimonialsDefinition, contactDefinition, ctaDefinition, whyUsDefinition }

/**
 * Export components for React use
 */
export { Hero } from './hero/v1/Hero'
export { Products } from './products/v1/Products'
export { Testimonials } from './testimonials/v1/Testimonials'
export { Contact } from './contact/v1/Contact'
export { CTA } from './cta/v1/CTA'
export { WhyUs } from './why_us/v1/WhyUs'

// Export types
export type { HeroProps } from './hero/v1/types'
export type { ProductsProps } from './products/v1/types'
export type { TestimonialsProps } from './testimonials/v1/types'
export type { ContactProps } from './contact/v1/types'
export type { CTAProps } from './cta/v1/types'
export type { WhyUsProps } from './why_us/v1/types'
