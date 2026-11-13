'use client'

import type { JSX } from 'react'

import React from 'react'
import { ThemeComponent, ComponentType } from '@/types/theme'
import { HeroEditor } from './HeroEditor'
import { ProductGridEditor } from './ProductGridEditor'
import { WhyUsEditor } from './WhyUsEditor'
import { ContactEditor } from './ContactEditor'
import { TestimonialsEditor } from './TestimonialsEditor'
import { CTAEditor } from './CTAEditor'
import { GenericEditor } from './GenericEditor'

interface ComponentEditorRouterProps {
  component: ThemeComponent
  onChange: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
}

/**
 * Router component that selects the appropriate editor based on component type
 * Maps each component type to its specialized editor
 */
export function ComponentEditorRouter({
  component,
  onChange,
  onPreview,
  className = '',
}: ComponentEditorRouterProps): JSX.Element {
  const renderEditor = (type: ComponentType): JSX.Element | null => {
    switch (type) {
      case 'hero':
        return (
          <HeroEditor
            component={component}
            onChange={onChange}
            onPreview={onPreview}
            className={className}
          />
        )

      case 'products':
        return (
          <ProductGridEditor
            component={component}
            onChange={onChange}
            onPreview={onPreview}
            className={className}
          />
        )

      case 'why_us':
        return (
          <WhyUsEditor
            component={component}
            onChange={onChange}
            onPreview={onPreview}
            className={className}
          />
        )

      case 'contact':
        return (
          <ContactEditor
            component={component}
            onChange={onChange}
            onPreview={onPreview}
            className={className}
          />
        )

      case 'testimonials':
        return (
          <TestimonialsEditor
            component={component}
            onChange={onChange}
            onPreview={onPreview}
            className={className}
          />
        )

      case 'cta':
        return (
          <CTAEditor
            component={component}
            onChange={onChange}
            onPreview={onPreview}
            className={className}
          />
        )

      case 'custom':
      default:
        return (
          <GenericEditor
            component={component}
            onChange={onChange}
            onPreview={onPreview}
            className={className}
          />
        )
    }
  }

  return <>{renderEditor(component.type)}</>
}
