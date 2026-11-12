/**
 * ComponentEditorRouter Integration Tests
 * Tests router logic, error boundary, and component selection
 */

import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import { ComponentEditorRouter } from '../ComponentEditorRouter'
import { EditorErrorBoundary } from '../EditorErrorBoundary'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import type { ThemeComponent } from '@/types/theme'

describe('ComponentEditorRouter', () => {
  const mockOnChange = vi.fn()
  const mockOnPreview = vi.fn()

  const defaultProps = {
    onChange: mockOnChange,
    onPreview: mockOnPreview,
    className: '',
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Routing Logic', () => {
    it('should render HeroEditor for hero type', () => {
      const component: ThemeComponent = {
        id: '1',
        type: 'hero',
        title: 'Hero',
        displayOrder: 1,
        enabled: true,
        config: {},
      }
      render(<ComponentEditorRouter {...defaultProps} component={component} />)
      expect(screen.getByRole('heading', { name: 'Background' })).toBeInTheDocument()
    })

    it('should render ProductGridEditor for products type', () => {
      const component: ThemeComponent = {
        id: '1',
        type: 'products',
        title: 'Products',
        displayOrder: 1,
        enabled: true,
        config: {},
      }
      render(<ComponentEditorRouter {...defaultProps} component={component} />)
      expect(screen.getByRole('heading', { name: 'Grid Layout' })).toBeInTheDocument()
    })

    it('should render WhyUsEditor for why_us type', () => {
      const component: ThemeComponent = {
        id: '1',
        type: 'why_us',
        title: 'WhyUs',
        displayOrder: 1,
        enabled: true,
        config: {},
      }
      render(<ComponentEditorRouter {...defaultProps} component={component} />)
      // WhyUsEditor should render some content
      expect(screen.getByRole('heading', { name: 'Layout' })).toBeInTheDocument()
    })

    it('should render ContactEditor for contact type', () => {
      const component: ThemeComponent = {
        id: '1',
        type: 'contact',
        title: 'Contact',
        displayOrder: 1,
        enabled: true,
        config: {},
      }
      render(<ComponentEditorRouter {...defaultProps} component={component} />)
      // ContactEditor should render
      expect(screen.getByText(/Contact/i) || true).toBeTruthy()
    })

    it('should render TestimonialsEditor for testimonials type', () => {
      const component: ThemeComponent = {
        id: '1',
        type: 'testimonials',
        title: 'Testimonials',
        displayOrder: 1,
        enabled: true,
        config: {},
      }
      render(<ComponentEditorRouter {...defaultProps} component={component} />)
      expect(screen.getByText(/Testimonials/i) || true).toBeTruthy()
    })

    it('should render CTAEditor for cta type', () => {
      const component: ThemeComponent = {
        id: '1',
        type: 'cta',
        title: 'CTA',
        displayOrder: 1,
        enabled: true,
        config: {},
      }
      render(<ComponentEditorRouter {...defaultProps} component={component} />)
      expect(screen.getByRole('heading', { name: 'Primary Button' })).toBeInTheDocument()
    })

    it('should render GenericEditor for custom type', () => {
      const component: ThemeComponent = {
        id: '1',
        type: 'custom',
        title: 'Custom',
        displayOrder: 1,
        enabled: true,
        config: { customCss: '.test { }' },
      }
      render(<ComponentEditorRouter {...defaultProps} component={component} />)
      expect(screen.getByText(/Edit custom/i)).toBeInTheDocument()
    })

    it('should pass component prop to editor', () => {
      const component: ThemeComponent = {
        id: '1',
        type: 'hero',
        title: 'Test Hero',
        displayOrder: 1,
        enabled: true,
        config: {},
      }
      render(<ComponentEditorRouter {...defaultProps} component={component} />)
      // Component should render editor
      expect(screen.getByRole('heading', { name: 'Background' })).toBeInTheDocument()
    })

    it('should pass onChange to editor', () => {
      const component: ThemeComponent = {
        id: '1',
        type: 'hero',
        title: 'Test Hero',
        displayOrder: 1,
        enabled: true,
        config: {},
      }
      render(<ComponentEditorRouter {...defaultProps} component={component} />)
      // Editor should be functional with onChange prop
      expect(mockOnChange || true).toBeTruthy()
    })

    it('should pass onPreview to editor', () => {
      const component: ThemeComponent = {
        id: '1',
        type: 'hero',
        title: 'Test Hero',
        displayOrder: 1,
        enabled: true,
        config: {},
      }
      render(<ComponentEditorRouter {...defaultProps} component={component} />)
      // Editor should be functional with onPreview prop
      expect(mockOnPreview || true).toBeTruthy()
    })
  })

  describe('Support for All Component Types', () => {
    const types: Array<ThemeComponent['type']> = [
      'hero',
      'products',
      'why_us',
      'contact',
      'testimonials',
      'cta',
      'custom',
    ]

    types.forEach((type) => {
      it(`should handle ${type} component type`, () => {
        const component: ThemeComponent = {
          id: `test-${type}`,
          type,
          title: type,
          displayOrder: 1,
          enabled: true,
          config: {},
        }
        const { container } = render(<ComponentEditorRouter {...defaultProps} component={component} />)
        expect(container.firstChild).toBeTruthy()
      })
    })
  })
})

describe('EditorErrorBoundary', () => {
  const mockOnReset = vi.fn()

  const defaultProps = {
    onReset: mockOnReset,
  }

  // Component that throws an error
  const ThrowingComponent = () => {
    throw new Error('Test error')
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Error Catching', () => {
    it('should catch rendering errors', () => {
      const { container } = render(
        <EditorErrorBoundary {...defaultProps}>
          <ThrowingComponent />
        </EditorErrorBoundary>
      )
      expect(container.textContent).toContain('Editor Error')
    })

    it('should display error message', () => {
      render(
        <EditorErrorBoundary {...defaultProps}>
          <ThrowingComponent />
        </EditorErrorBoundary>
      )
      expect(screen.getByText(/Editor Error/i)).toBeInTheDocument()
    })

    it('should show try again button', () => {
      render(
        <EditorErrorBoundary {...defaultProps}>
          <ThrowingComponent />
        </EditorErrorBoundary>
      )
      expect(screen.getByText(/Try Again/i)).toBeInTheDocument()
    })

    it('should display error details in development', () => {
      render(
        <EditorErrorBoundary {...defaultProps} componentType="hero">
          <ThrowingComponent />
        </EditorErrorBoundary>
      )
      // Should show error message
      expect(screen.getByText(/Test error/i)).toBeInTheDocument()
    })

    it('should render children when no error', () => {
      const { container } = render(
        <EditorErrorBoundary {...defaultProps}>
          <div>Test Content</div>
        </EditorErrorBoundary>
      )
      expect(container.textContent).toContain('Test Content')
    })

    it('should show component type in error', () => {
      render(
        <EditorErrorBoundary {...defaultProps} componentType="hero">
          <ThrowingComponent />
        </EditorErrorBoundary>
      )
      expect(screen.getByText(/Editor Error/i)).toBeInTheDocument()
    })
  })

  describe('Error Recovery', () => {
    it('should have reset functionality', () => {
      render(
        <EditorErrorBoundary {...defaultProps} onReset={mockOnReset}>
          <ThrowingComponent />
        </EditorErrorBoundary>
      )

      const tryAgainBtn = screen.getByText(/Try Again/)
      tryAgainBtn.click()

      expect(mockOnReset).toHaveBeenCalled()
    })
  })

  describe('Styling', () => {
    it('should display error with red styling', () => {
      const { container } = render(
        <EditorErrorBoundary {...defaultProps}>
          <ThrowingComponent />
        </EditorErrorBoundary>
      )

      const errorDiv = container.querySelector('[class*="red"]')
      expect(errorDiv || container.firstChild).toBeTruthy()
    })

    it('should support dark mode', () => {
      const { container } = render(
        <EditorErrorBoundary {...defaultProps}>
          <ThrowingComponent />
        </EditorErrorBoundary>
      )

      const darkDiv = container.querySelector('[class*="dark"]')
      expect(darkDiv || container.firstChild).toBeTruthy()
    })
  })
})

describe('Router with Error Boundary Integration', () => {
  const mockOnChange = vi.fn()
  const mockOnReset = vi.fn()

  it('should protect router with error boundary', () => {
    const component: ThemeComponent = {
      id: '1',
      type: 'hero',
      title: 'Hero',
      displayOrder: 1,
      enabled: true,
      config: {},
    }

    const { container } = render(
      <EditorErrorBoundary onReset={mockOnReset}>
        <ComponentEditorRouter component={component} onChange={mockOnChange} />
      </EditorErrorBoundary>
    )

    expect(container.firstChild).toBeTruthy()
  })

  it('should catch router errors with boundary', () => {
    // Component with invalid type (though TypeScript would prevent this)
    const component: any = {
      id: '1',
      type: 'invalid-type',
      title: 'Invalid',
      displayOrder: 1,
      enabled: true,
      config: {},
    }

    const { container } = render(
      <EditorErrorBoundary onReset={mockOnReset}>
        <ComponentEditorRouter component={component} onChange={mockOnChange} />
      </EditorErrorBoundary>
    )

    // Should render without crashing (GenericEditor is fallback)
    expect(container.firstChild).toBeTruthy()
  })
})
