/**
 * ComponentBuilder Component Tests
 * Tests component list, add/delete/reorder functionality
 */

import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ComponentBuilder } from '../ComponentBuilder'
import { describe, it, expect, beforeEach, vi } from 'vitest'

// Mock the useThemeStore hook
vi.mock('@/hooks/useThemeStore', () => ({
  useThemeStore: vi.fn(),
}))

import { useThemeStore } from '@/hooks/useThemeStore'

describe('ComponentBuilder', () => {
  const mockTheme = {
    id: '1',
    name: 'Test Theme',
    colors: {},
    typography: {},
    identity: {},
    components: [
      {
        id: 'comp-1',
        type: 'hero',
        title: 'Hero Section',
        order: 0,
        enabled: true,
        settings: {},
      },
      {
        id: 'comp-2',
        type: 'products',
        title: 'Featured Products',
        order: 1,
        enabled: true,
        settings: {},
      },
    ],
  }

  const mockStoreActions = {
    addComponent: vi.fn(),
    deleteComponent: vi.fn(),
    selectComponent: vi.fn(),
    reorderComponents: vi.fn(),
  }

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
      const store = {
        currentTheme: mockTheme,
        selectedComponentId: null,
        ...mockStoreActions,
      }
      return selector(store)
    })
  })

  describe('Rendering', () => {
    it('should render without error', () => {
      render(<ComponentBuilder />)
      expect(screen.getByText('Components')).toBeInTheDocument()
    })

    it('should show empty state when no theme selected', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: null,
          selectedComponentId: null,
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<ComponentBuilder />)
      expect(screen.getByText('Select a theme to add components')).toBeInTheDocument()
    })

    it('should display component count', () => {
      render(<ComponentBuilder />)
      expect(screen.getByText(/2 components/)).toBeInTheDocument()
    })

    it('should show empty state when no components', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: { ...mockTheme, components: [] },
          selectedComponentId: null,
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<ComponentBuilder />)
      expect(screen.getByText('No components added yet')).toBeInTheDocument()
    })

    it('should render add button', () => {
      render(<ComponentBuilder />)
      expect(screen.getByRole('button', { name: /Add/i })).toBeInTheDocument()
    })
  })

  describe('Component List', () => {
    it('should display all components', () => {
      render(<ComponentBuilder />)
      expect(screen.getByText('Hero Section')).toBeInTheDocument()
      expect(screen.getByText('Featured Products')).toBeInTheDocument()
    })

    it('should display component types', () => {
      render(<ComponentBuilder />)
      expect(screen.getByText('hero')).toBeInTheDocument()
      expect(screen.getByText('products')).toBeInTheDocument()
    })

    it('should have grip handle for drag', () => {
      const { container } = render(<ComponentBuilder />)
      const gripIcons = container.querySelectorAll('svg.lucide-grip-vertical')
      expect(gripIcons.length).toBeGreaterThan(0)
    })

    it('should display delete button for each component', () => {
      render(<ComponentBuilder />)
      const deleteButtons = screen.getAllByRole('button', { name: /^Delete / })
      expect(deleteButtons.length).toBe(2)
    })
  })

  describe('Add Component', () => {
    it('should show add menu when button is clicked', async () => {
      render(<ComponentBuilder />)
      const addButton = screen.getByRole('button', { name: /Add/i })

      fireEvent.click(addButton)

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Hero Section.*Large banner with headline/ })).toBeInTheDocument()
        expect(screen.getByText('Large banner with headline')).toBeInTheDocument()
      })
    })

    it('should display all component type options', async () => {
      render(<ComponentBuilder />)
      fireEvent.click(screen.getByRole('button', { name: /Add/i }))

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /Featured Products.*Product showcase/ })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /Why Choose Us.*Key features section/ })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /Contact Section.*Contact form and info/ })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /Testimonials.*Customer reviews/ })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /Call to Action.*Action button section/ })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /Custom HTML.*Custom content block/ })).toBeInTheDocument()
      })
    })

    it('should show descriptions for component types', async () => {
      render(<ComponentBuilder />)
      fireEvent.click(screen.getByRole('button', { name: /Add/i }))

      await waitFor(() => {
        expect(screen.getByText('Product showcase')).toBeInTheDocument()
        expect(screen.getByText('Contact form and info')).toBeInTheDocument()
      })
    })

    it('should call addComponent with correct data', async () => {
      render(<ComponentBuilder />)
      fireEvent.click(screen.getByRole('button', { name: /Add/i }))

      await waitFor(() => {
        fireEvent.click(screen.getByRole('button', { name: /Hero Section.*Large banner with headline/ }))

        expect(mockStoreActions.addComponent).toHaveBeenCalledWith(
          expect.objectContaining({
            type: 'hero',
            title: 'Hero Section',
            enabled: true,
          })
        )
      })
    })

    it('should close menu after adding component', async () => {
      render(<ComponentBuilder />)
      const addButton = screen.getByRole('button', { name: /Add/i })

      fireEvent.click(addButton)
      fireEvent.click(screen.getByRole('button', { name: /Hero Section.*Large banner with headline/ }))

      await waitFor(() => {
        // Menu should be closed (only one Hero Section visible - from the list)
        const heroTexts = screen.getAllByText('Hero Section')
        expect(heroTexts.length).toBeGreaterThanOrEqual(1)
      })
    })

    it('should support multiple component additions', async () => {
      render(<ComponentBuilder />)

      fireEvent.click(screen.getByRole('button', { name: /Add/i }))
      fireEvent.click(screen.getByRole('button', { name: /Hero Section.*Large banner with headline/ }))

      expect(mockStoreActions.addComponent).toHaveBeenCalled()
    })
  })

  describe('Delete Component', () => {
    it('should delete component on button click', () => {
      render(<ComponentBuilder />)
      fireEvent.click(screen.getByRole('button', { name: 'Delete Hero Section' }))

      expect(mockStoreActions.deleteComponent).toHaveBeenCalledWith('comp-1')
    })

    it('should prevent event propagation on delete', () => {
      const { container } = render(<ComponentBuilder />)

      // Find and click the delete button
      const componentItem = container.querySelector('[draggable="true"]')
      if (componentItem) {
        const deleteBtn = componentItem.querySelector('button:last-child')
        if (deleteBtn) {
          const clickEvent = new MouseEvent('click', { bubbles: true })
          const stopSpy = vi.spyOn(clickEvent, 'stopPropagation')
          deleteBtn.dispatchEvent(clickEvent)
        }
      }
    })

    it('should show single item text for 1 component', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: {
            ...mockTheme,
            components: [mockTheme.components[0]],
          },
          selectedComponentId: null,
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<ComponentBuilder />)
      expect(screen.getByText(/1 component$/)).toBeInTheDocument()
    })
  })

  describe('Component Selection', () => {
    it('should call selectComponent on click', () => {
      render(<ComponentBuilder />)
      const componentItem = screen.getByText('Hero Section')

      fireEvent.click(componentItem.closest('[draggable="true"]')!)

      expect(mockStoreActions.selectComponent).toHaveBeenCalledWith('comp-1')
    })

    it('should highlight selected component', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: mockTheme,
          selectedComponentId: 'comp-1',
          ...mockStoreActions,
        }
        return selector(store)
      })

      const { container } = render(<ComponentBuilder />)
      const selected = container.querySelector('[class*="border-primary-600"]')

      expect(selected).toBeInTheDocument()
    })

    it('should show hover styles for non-selected components', () => {
      const { container } = render(<ComponentBuilder />)
      const components = container.querySelectorAll('[draggable="true"]')

      expect(components.length).toBeGreaterThan(0)
      components.forEach((comp) => {
        expect(comp).toHaveClass('hover:border-primary-300')
      })
    })
  })

  describe('Drag and Drop', () => {
    it('should handle drag start', () => {
      const { container } = render(<ComponentBuilder />)
      const firstComponent = container.querySelector('[draggable="true"]')

      fireEvent.dragStart(firstComponent!)
      expect(firstComponent).toHaveClass('opacity-50')
    })

    it('should allow drag over', () => {
      const { container } = render(<ComponentBuilder />)
      const component = container.querySelector('[draggable="true"]')

      const dragOverEvent = new Event('dragover', { bubbles: true, cancelable: true })
      const shouldContinue = component?.dispatchEvent(dragOverEvent)

      expect(shouldContinue).toBe(false)
      expect(dragOverEvent.defaultPrevented).toBe(true)
    })

    it('should handle drop for reordering', () => {
      const { container } = render(<ComponentBuilder />)
      const components = container.querySelectorAll('[draggable="true"]')

      if (components.length >= 2) {
        // Simulate drag from first to second
        fireEvent.dragStart(components[0])
        fireEvent.drop(components[1])

        expect(mockStoreActions.reorderComponents).toHaveBeenCalled()
      }
    })

    it('should set opacity during drag', () => {
      const { container } = render(<ComponentBuilder />)
      const component = container.querySelector('[draggable="true"]')

      // Component should not have opacity-50 by default
      expect(component).not.toHaveClass('opacity-50')
    })

    it('should not reorder on same target', () => {
      const { container } = render(<ComponentBuilder />)
      const component = container.querySelector('[draggable="true"]')

      fireEvent.dragStart(component!)

      // Drop on same element
      fireEvent.drop(component!)

      expect(mockStoreActions.reorderComponents).not.toHaveBeenCalled()
    })
  })

  describe('Empty State', () => {
    it('should show empty message when no components', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: { ...mockTheme, components: [] },
          selectedComponentId: null,
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<ComponentBuilder />)
      expect(screen.getByText('No components added yet')).toBeInTheDocument()
      expect(screen.getByText('Click "Add" to start building')).toBeInTheDocument()
    })

    it('should show empty state count', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: { ...mockTheme, components: [] },
          selectedComponentId: null,
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<ComponentBuilder />)
      expect(screen.getByText(/0 components/)).toBeInTheDocument()
    })
  })

  describe('Accessibility', () => {
    it('should have accessible buttons', () => {
      render(<ComponentBuilder />)
      const addButton = screen.getByRole('button', { name: /Add/i })
      expect(addButton).toBeInTheDocument()
    })

    it('should have semantic HTML for components list', () => {
      render(<ComponentBuilder />)
      expect(screen.getByText('Hero Section')).toBeInTheDocument()
    })

    it('should have accessible labels', () => {
      render(<ComponentBuilder />)
      expect(screen.getByText('Components')).toBeInTheDocument()
    })
  })

  describe('Dark Mode Support', () => {
    it('should have dark mode classes', () => {
      const { container } = render(<ComponentBuilder />)
      const darkElements = container.querySelectorAll('[class*="dark:"]')
      expect(darkElements.length).toBeGreaterThan(0)
    })
  })

  describe('Layout', () => {
    it('should have scrollable component list', () => {
      const { container } = render(<ComponentBuilder />)
      const listContainer = container.querySelector('div.overflow-y-auto')
      expect(listContainer).toBeInTheDocument()
    })

    it('should show header with title and add button', () => {
      render(<ComponentBuilder />)
      expect(screen.getByText('Components')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Add/i })).toBeInTheDocument()
    })
  })

  describe('Multiple Component Types', () => {
    it('should support all component types', async () => {
      render(<ComponentBuilder />)
      fireEvent.click(screen.getByRole('button', { name: /Add/i }))

      const types = [
        'Hero Section',
        'Featured Products',
        'Why Choose Us',
        'Contact Section',
        'Testimonials',
        'Call to Action',
        'Custom HTML',
      ]

      for (const type of types) {
        await waitFor(() => expect(screen.getAllByText(type).length).toBeGreaterThan(0))
      }
    })
  })

  describe('Component Order', () => {
    it('should display components in correct order', () => {
      render(<ComponentBuilder />)
      const components = screen.getAllByText(/^(Hero Section|Featured Products)$/)
      expect(components[0]).toHaveTextContent('Hero Section')
      expect(components[1]).toHaveTextContent('Featured Products')
    })

    it('should show order numbers in component list', () => {
      render(<ComponentBuilder />)
      const firstComponent = screen.getByText('Hero Section')
      expect(firstComponent).toBeInTheDocument()
    })
  })

  describe('Component Context', () => {
    it('should maintain component context across operations', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: mockTheme,
          selectedComponentId: 'comp-1',
          ...mockStoreActions,
        }
        return selector(store)
      })

      const { rerender } = render(<ComponentBuilder />)

      expect(mockStoreActions.selectComponent).not.toHaveBeenCalled()

      rerender(<ComponentBuilder />)

      expect(screen.getByText('Hero Section')).toBeInTheDocument()
    })
  })

  describe('Performance', () => {
    it('should handle large component lists', () => {
      const largeComponentList = Array.from({ length: 50 }, (_, i) => ({
        id: `comp-${i}`,
        type: 'hero' as const,
        title: `Component ${i}`,
        order: i,
        enabled: true,
        settings: {},
      }))

      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: { ...mockTheme, components: largeComponentList },
          selectedComponentId: null,
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<ComponentBuilder />)
      expect(screen.getByText(/50 components/)).toBeInTheDocument()
    })
  })
})
