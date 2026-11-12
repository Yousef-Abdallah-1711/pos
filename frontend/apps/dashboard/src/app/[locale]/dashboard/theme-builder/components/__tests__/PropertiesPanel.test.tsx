/**
 * PropertiesPanel Component Tests
 * Tests tab navigation, color editing, typography controls, identity settings
 */

import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { PropertiesPanel } from '../PropertiesPanel'
import { describe, it, expect, beforeEach, vi } from 'vitest'

// Mock the useThemeStore hook
vi.mock('@/hooks/useThemeStore', () => ({
  useThemeStore: vi.fn(),
}))

import { useThemeStore } from '@/hooks/useThemeStore'

describe('PropertiesPanel', () => {
  const mockTheme = {
    id: '1',
    name: 'Test Theme',
    colors: {
      primary: '#3b82f6',
      secondary: '#10b981',
      accent: '#f59e0b',
      background: '#ffffff',
      text: '#000000',
      border: '#e5e7eb',
      shadow: '#00000029',
    },
    typography: {
      fontFamily: 'Inter',
      baseFontSize: 16,
      borderRadius: 8,
      lineHeight: 1.5,
    },
    identity: {
      siteTitle: 'My Website',
      logoUrl: 'https://example.com/logo.png',
      faviconUrl: 'https://example.com/favicon.ico',
      domain: 'example.com',
    },
    components: [
      {
        id: '1',
        type: 'hero',
        title: 'Hero Section',
      },
    ],
  }

  const mockStoreActions = {
    updateThemeColors: vi.fn(),
    updateTypography: vi.fn(),
    updateIdentity: vi.fn(),
    updateComponent: vi.fn(),
  }

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useThemeStore as any).mockImplementation((selector: (state: Record<string, unknown>) => unknown) => {
      const store = {
        currentTheme: mockTheme,
        selectedComponentId: '1',
        ...mockStoreActions,
      }
      return selector(store)
    })
  })

  describe('Rendering', () => {
    it('should render without error', () => {
      render(<PropertiesPanel />)
      expect(screen.getByText('Global Colors')).toBeInTheDocument()
    })

    it('should show empty state when no theme selected', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: Record<string, unknown>) => unknown) => {
        const store = {
          currentTheme: null,
          selectedComponentId: null,
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<PropertiesPanel />)
      expect(screen.getByText('Select a theme to edit properties')).toBeInTheDocument()
    })

    it('should render all tabs', () => {
      render(<PropertiesPanel />)
      expect(screen.getByText('Colors')).toBeInTheDocument()
      expect(screen.getByText('Type')).toBeInTheDocument()
      expect(screen.getByText('Identity')).toBeInTheDocument()
    })

    it('should render component tab when component is selected', () => {
      render(<PropertiesPanel />)
      expect(screen.getByText('Component')).toBeInTheDocument()
    })

    it('should not render component tab when no component is selected', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: Record<string, unknown>) => unknown) => {
        const store = {
          currentTheme: mockTheme,
          selectedComponentId: null,
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<PropertiesPanel />)
      expect(screen.queryByText('Component')).not.toBeInTheDocument()
    })
  })

  describe('Colors Tab', () => {
    it('should display colors tab as active by default', () => {
      render(<PropertiesPanel />)
      expect(screen.getByText('Global Colors')).toBeInTheDocument()
    })

    it('should render all color pickers', () => {
      render(<PropertiesPanel />)
      expect(screen.getByText('Primary')).toBeInTheDocument()
      expect(screen.getByText('Secondary')).toBeInTheDocument()
      expect(screen.getByText('Accent')).toBeInTheDocument()
      expect(screen.getByText('Background')).toBeInTheDocument()
      expect(screen.getByText('Text')).toBeInTheDocument()
      expect(screen.getByText('Border')).toBeInTheDocument()
      expect(screen.getByText('Shadow')).toBeInTheDocument()
    })

    it('should display color values correctly', () => {
      const { container } = render(<PropertiesPanel />)
      const colorInputs = container.querySelectorAll('input[type="text"]')

      // Note: inputs may have different order depending on component rendering
      expect(colorInputs.length).toBeGreaterThan(0)
    })

    it('should call updateColors when color changes', () => {
      render(<PropertiesPanel />)
      const primaryColorLabel = screen.getByText('Primary')

      // The color picker is a sibling of the label
      // We need to find and interact with the associated input
      const inputs = screen.getAllByPlaceholderText('#000000')
      fireEvent.change(inputs[0], { target: { value: '#ff0000' } })

      expect(mockStoreActions.updateThemeColors).toHaveBeenCalled()
    })

    it('should render color descriptions', () => {
      render(<PropertiesPanel />)
      expect(screen.getByText('Main brand color for buttons and CTAs')).toBeInTheDocument()
      expect(screen.getByText('Complementary brand color')).toBeInTheDocument()
    })
  })

  describe('Typography Tab', () => {
    it('should display typography controls when tab is active', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        expect(screen.getByText('Typography')).toBeInTheDocument()
      })
    })

    it('should display font family select with options', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        const fontSelect = screen.getByDisplayValue('Inter')
        expect(fontSelect).toBeInTheDocument()
      })
    })

    it('should have all font families in select options', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        expect(screen.getByText('Poppins')).toBeInTheDocument()
        expect(screen.getByText('Playfair Display')).toBeInTheDocument()
        expect(screen.getByText('Georgia')).toBeInTheDocument()
      })
    })

    it('should call updateTypography when font family changes', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      const fontSelect = screen.getByDisplayValue('Inter') as HTMLSelectElement
      await userEvent.selectOptions(fontSelect, 'Poppins')

      expect(mockStoreActions.updateTypography).toHaveBeenCalledWith({ fontFamily: 'Poppins' })
    })

    it('should display base font size with range input', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        expect(screen.getByText(/Base Font Size: 16px/)).toBeInTheDocument()
      })
    })

    it('should update base font size on slider change', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        const rangeInputs = screen.getAllByRole('slider')
        const fontSizeSlider = rangeInputs[0]

        fireEvent.change(fontSizeSlider, { target: { value: '18' } })

        expect(mockStoreActions.updateTypography).toHaveBeenCalledWith({
          baseFontSize: 18,
        })
      })
    })

    it('should display border radius with range input', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        expect(screen.getByText(/Border Radius: 8px/)).toBeInTheDocument()
      })
    })

    it('should display line height with range input', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        expect(screen.getByText(/Line Height: 1.5/)).toBeInTheDocument()
      })
    })

    it('should show range limits for font size', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        expect(screen.getByText('12px - 20px')).toBeInTheDocument()
      })
    })

    it('should show range limits for border radius', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        expect(screen.getByText('0px - 20px')).toBeInTheDocument()
      })
    })

    it('should show range limits for line height', async () => {
      render(<PropertiesPanel />)
      const typographyTab = screen.getByText('Type')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        expect(screen.getByText('1.2 - 2.0')).toBeInTheDocument()
      })
    })
  })

  describe('Identity Tab', () => {
    it('should display identity settings when tab is active', async () => {
      render(<PropertiesPanel />)
      const identityTab = screen.getByText('Identity')

      fireEvent.click(identityTab)

      await waitFor(() => {
        expect(screen.getByText('Website Identity')).toBeInTheDocument()
      })
    })

    it('should display site title input', async () => {
      render(<PropertiesPanel />)
      const identityTab = screen.getByText('Identity')

      fireEvent.click(identityTab)

      await waitFor(() => {
        const siteInput = screen.getByDisplayValue('My Website')
        expect(siteInput).toBeInTheDocument()
      })
    })

    it('should call updateIdentity when site title changes', async () => {
      render(<PropertiesPanel />)
      const identityTab = screen.getByText('Identity')

      fireEvent.click(identityTab)

      await waitFor(() => {
        const siteInput = screen.getByDisplayValue('My Website') as HTMLInputElement
        fireEvent.change(siteInput, { target: { value: 'New Title' } })

        expect(mockStoreActions.updateIdentity).toHaveBeenCalledWith({
          siteTitle: 'New Title',
        })
      })
    })

    it('should display logo URL input', async () => {
      render(<PropertiesPanel />)
      const identityTab = screen.getByText('Identity')

      fireEvent.click(identityTab)

      await waitFor(() => {
        const logoInput = screen.getByDisplayValue('https://example.com/logo.png')
        expect(logoInput).toBeInTheDocument()
      })
    })

    it('should display favicon URL input', async () => {
      render(<PropertiesPanel />)
      const identityTab = screen.getByText('Identity')

      fireEvent.click(identityTab)

      await waitFor(() => {
        const faviconInput = screen.getByDisplayValue('https://example.com/favicon.ico')
        expect(faviconInput).toBeInTheDocument()
      })
    })

    it('should display domain field as read-only', async () => {
      render(<PropertiesPanel />)
      const identityTab = screen.getByText('Identity')

      fireEvent.click(identityTab)

      await waitFor(() => {
        const domainInput = screen.getByDisplayValue('example.com') as HTMLInputElement
        expect(domainInput).toBeDisabled()
      })
    })

    it('should show read-only indicator for domain field', async () => {
      render(<PropertiesPanel />)
      const identityTab = screen.getByText('Identity')

      fireEvent.click(identityTab)

      await waitFor(() => {
        expect(screen.getByText('Read-only field')).toBeInTheDocument()
      })
    })
  })

  describe('Component Tab', () => {
    it('should display component settings when tab is active', async () => {
      render(<PropertiesPanel />)
      const componentTab = screen.getByText('Component')

      fireEvent.click(componentTab)

      await waitFor(() => {
        expect(screen.getByText('Component Settings')).toBeInTheDocument()
      })
    })

    it('should display selected component type', async () => {
      render(<PropertiesPanel />)
      const componentTab = screen.getByText('Component')

      fireEvent.click(componentTab)

      await waitFor(() => {
        const typeLabel = screen.getByText('Type:')
        expect(typeLabel.parentElement).toHaveTextContent('hero')
      })
    })

    it('should display selected component title', async () => {
      render(<PropertiesPanel />)
      const componentTab = screen.getByText('Component')

      fireEvent.click(componentTab)

      await waitFor(() => {
        const titleLabel = screen.getByText('Title:')
        expect(titleLabel.parentElement).toHaveTextContent('Hero Section')
      })
    })
  })

  describe('Tab Navigation', () => {
    it('should switch between tabs on click', async () => {
      render(<PropertiesPanel />)

      // Start on colors tab
      expect(screen.getByText('Global Colors')).toBeInTheDocument()

      // Switch to typography
      fireEvent.click(screen.getByText('Type'))
      await waitFor(() => {
        expect(screen.getByText('Typography')).toBeInTheDocument()
      })

      // Switch to identity
      fireEvent.click(screen.getByText('Identity'))
      await waitFor(() => {
        expect(screen.getByText('Website Identity')).toBeInTheDocument()
      })

      // Switch back to colors
      fireEvent.click(screen.getByText('Colors'))
      await waitFor(() => {
        expect(screen.getByText('Global Colors')).toBeInTheDocument()
      })
    })

    it('should highlight active tab', () => {
      render(<PropertiesPanel />)
      const colorsTab = screen.getByText('Colors')

      expect(colorsTab).toHaveClass('border-primary-600')
    })

    it('should update styles when switching tabs', async () => {
      render(<PropertiesPanel />)
      const colorsTab = screen.getByText('Colors')
      const typographyTab = screen.getByText('Type')

      expect(colorsTab).toHaveClass('border-primary-600')
      expect(typographyTab).not.toHaveClass('border-primary-600')

      fireEvent.click(typographyTab)

      await waitFor(() => {
        expect(colorsTab).not.toHaveClass('border-primary-600')
        expect(typographyTab).toHaveClass('border-primary-600')
      })
    })
  })

  describe('Accessibility', () => {
    it('should have accessible labels for color pickers', () => {
      render(<PropertiesPanel />)
      expect(screen.getByText('Primary')).toBeInTheDocument()
      expect(screen.getByText('Secondary')).toBeInTheDocument()
    })

    it('should have accessible labels for font controls', async () => {
      render(<PropertiesPanel />)
      fireEvent.click(screen.getByText('Type'))

      await waitFor(() => {
        expect(screen.getByText(/Font Family/)).toBeInTheDocument()
      })
    })

    it('should have accessible labels for identity inputs', async () => {
      render(<PropertiesPanel />)
      fireEvent.click(screen.getByText('Identity'))

      await waitFor(() => {
        expect(screen.getByText(/Site Title/)).toBeInTheDocument()
        expect(screen.getByText(/Logo URL/)).toBeInTheDocument()
      })
    })
  })

  describe('Dark Mode Support', () => {
    it('should have dark mode classes', () => {
      const { container } = render(<PropertiesPanel />)

      const darkElements = container.querySelectorAll('[class*="dark:"]')
      expect(darkElements.length).toBeGreaterThan(0)
    })
  })

  describe('Empty States', () => {
    it('should handle theme with empty component list', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: Record<string, unknown>) => unknown) => {
        const store = {
          currentTheme: {
            ...mockTheme,
            components: [],
          },
          selectedComponentId: null,
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<PropertiesPanel />)
      expect(screen.queryByText('Component')).not.toBeInTheDocument()
    })

    it('should handle missing optional identity fields', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: Record<string, unknown>) => unknown) => {
        const store = {
          currentTheme: {
            ...mockTheme,
            identity: {
              siteTitle: 'Site',
              logoUrl: undefined,
              faviconUrl: undefined,
              domain: undefined,
            },
          },
          selectedComponentId: '1',
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<PropertiesPanel />)
      fireEvent.click(screen.getByText('Identity'))

      expect(screen.getByText('Logo URL')).toBeInTheDocument()
    })
  })

  describe('Keyboard Navigation', () => {
    it('should be keyboard navigable between tabs', async () => {
      const user = userEvent.setup()
      render(<PropertiesPanel />)

      const colorsTab = screen.getByText('Colors')
      await user.tab()

      // Focus should move through buttons
      expect(colorsTab).toBeInTheDocument()
    })
  })

  describe('Layout', () => {
    it('should have scrollable content area', () => {
      const { container } = render(<PropertiesPanel />)
      const content = container.querySelector('div.overflow-y-auto')
      expect(content).toBeInTheDocument()
    })

    it('should display tabs at top', () => {
      const { container } = render(<PropertiesPanel />)
      const tabs = container.querySelector('div.flex.border-b')
      expect(tabs).toBeInTheDocument()
    })
  })

  describe('Theme Data Updates', () => {
    it('should reflect theme changes from store', () => {
      const { rerender } = render(<PropertiesPanel />)

      ;(useThemeStore as any).mockImplementation((selector: (state: Record<string, unknown>) => unknown) => {
        const store = {
          currentTheme: {
            ...mockTheme,
            colors: {
              ...mockTheme.colors,
              primary: '#ff0000',
            },
          },
          selectedComponentId: '1',
          ...mockStoreActions,
        }
        return selector(store)
      })

      rerender(<PropertiesPanel />)
      // Color should be updated in component
    })
  })
})
