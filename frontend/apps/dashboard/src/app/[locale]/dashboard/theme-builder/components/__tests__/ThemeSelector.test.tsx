/**
 * ThemeSelector Component Tests
 * Tests preset themes, custom themes, and theme actions
 */

import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { ThemeSelector } from '../ThemeSelector'
import { describe, it, expect, beforeEach, vi } from 'vitest'

vi.mock('@/hooks/useThemeStore', () => ({
  useThemeStore: vi.fn(),
}))

import { useThemeStore } from '@/hooks/useThemeStore'

describe('ThemeSelector', () => {
  const mockCurrentTheme = {
    id: '1',
    name: 'My Theme',
    colors: {
      primary: '#3b82f6',
      secondary: '#10b981',
      accent: '#f59e0b',
      background: '#ffffff',
      text: '#000000',
      border: '#e5e7eb',
      shadow: '#00000029',
    },
  }

  const mockThemes = [mockCurrentTheme]
  const mockStoreActions = {
    setCurrentTheme: vi.fn(),
  }

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
      const store = {
        currentTheme: mockCurrentTheme,
        themes: mockThemes,
        ...mockStoreActions,
      }
      return selector(store)
    })
  })

  describe('Rendering', () => {
    it('should render theme selector header', () => {
      render(<ThemeSelector />)
      expect(screen.getByText('Themes')).toBeInTheDocument()
    })

    it('should render new theme button', () => {
      render(<ThemeSelector />)
      expect(screen.getByRole('button', { name: /New Theme/i })).toBeInTheDocument()
    })

    it('should display preset section', () => {
      render(<ThemeSelector />)
      expect(screen.getByText('Presets')).toBeInTheDocument()
    })

    it('should display custom themes section', () => {
      render(<ThemeSelector />)
      expect(screen.getByText(/Custom Themes/)).toBeInTheDocument()
    })

    it('should show import button', () => {
      render(<ThemeSelector />)
      expect(screen.getByRole('button', { name: /Import/i })).toBeInTheDocument()
    })
  })

  describe('Preset Themes', () => {
    it('should display all preset themes', () => {
      render(<ThemeSelector />)
      expect(screen.getByText('Modern')).toBeInTheDocument()
      expect(screen.getByText('Warm')).toBeInTheDocument()
      expect(screen.getByText('Fresh')).toBeInTheDocument()
    })

    it('should show color swatches for presets', () => {
      const { container } = render(<ThemeSelector />)
      const colorDivs = container.querySelectorAll('div[class*="-color-"]')
      expect(colorDivs.length).toBeGreaterThan(0)
    })

    it('should apply preset theme on click', () => {
      render(<ThemeSelector />)
      const modernButton = screen.getByText('Modern')

      fireEvent.click(modernButton)

      expect(mockStoreActions.setCurrentTheme).toHaveBeenCalled()
    })

    it('should apply correct preset colors', () => {
      render(<ThemeSelector />)
      const modernButton = screen.getByText('Modern')

      fireEvent.click(modernButton)

      expect(mockStoreActions.setCurrentTheme).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'Modern',
          colors: expect.objectContaining({
            primary: '#3b82f6',
          }),
        })
      )
    })
  })

  describe('Custom Themes', () => {
    it('should display custom theme count', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: mockCurrentTheme,
          themes: [mockCurrentTheme, { ...mockCurrentTheme, id: '2' }],
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<ThemeSelector />)
      expect(screen.getByText(/Custom Themes \(2\)/)).toBeInTheDocument()
    })

    it('should show empty message when no custom themes', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: null,
          themes: [],
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<ThemeSelector />)
      expect(screen.getByText('No custom themes yet')).toBeInTheDocument()
    })

    it('should display custom theme names', () => {
      render(<ThemeSelector />)
      expect(screen.getByText('My Theme')).toBeInTheDocument()
    })

    it('should show actions button for custom themes', () => {
      render(<ThemeSelector />)
      const actionsButtons = screen.getAllByText(/Actions/i)
      expect(actionsButtons.length).toBeGreaterThan(0)
    })

    it('should select custom theme on click', () => {
      render(<ThemeSelector />)
      const themeItem = screen.getByText('My Theme').closest('div')

      fireEvent.click(themeItem!)

      expect(mockStoreActions.setCurrentTheme).toHaveBeenCalledWith(mockCurrentTheme)
    })

    it('should highlight selected theme', () => {
      const { container } = render(<ThemeSelector />)
      const selected = container.querySelector('[class*="border-primary-600"]')

      expect(selected).toBeInTheDocument()
    })
  })

  describe('Theme Actions', () => {
    it('should toggle actions menu', () => {
      render(<ThemeSelector />)
      const actionsButton = screen.getByText(/Actions/i)

      fireEvent.click(actionsButton)

      expect(screen.getByText('Duplicate')).toBeInTheDocument()
      expect(screen.getByText('Export')).toBeInTheDocument()
      expect(screen.getByText('Delete')).toBeInTheDocument()
    })

    it('should close actions menu on second click', async () => {
      render(<ThemeSelector />)
      const actionsButton = screen.getByText(/Actions/i)

      fireEvent.click(actionsButton)
      expect(screen.getByText('Duplicate')).toBeInTheDocument()

      fireEvent.click(actionsButton)

      await waitFor(() => {
        expect(screen.queryByText('Duplicate')).not.toBeInTheDocument()
      })
    })

    it('should prevent event propagation on actions click', () => {
      render(<ThemeSelector />)
      const actionsButton = screen.getByText(/Actions/i)

      fireEvent.click(actionsButton)

      expect(mockStoreActions.setCurrentTheme).not.toHaveBeenCalled()
    })

    it('should show duplicate action', () => {
      render(<ThemeSelector />)
      const actionsButton = screen.getByText(/Actions/i)

      fireEvent.click(actionsButton)

      expect(screen.getByText('Duplicate')).toBeInTheDocument()
    })

    it('should show export action', () => {
      render(<ThemeSelector />)
      const actionsButton = screen.getByText(/Actions/i)

      fireEvent.click(actionsButton)

      expect(screen.getByText('Export')).toBeInTheDocument()
    })

    it('should show delete action', () => {
      render(<ThemeSelector />)
      const actionsButton = screen.getByText(/Actions/i)

      fireEvent.click(actionsButton)

      expect(screen.getByText('Delete')).toBeInTheDocument()
    })

    it('should style delete action in red', () => {
      render(<ThemeSelector />)
      const actionsButton = screen.getByText(/Actions/i)

      fireEvent.click(actionsButton)

      const deleteButton = screen.getByText('Delete')
      expect(deleteButton).toHaveClass('text-red-600')
    })
  })

  describe('Color Display', () => {
    it('should show color swatches for each theme', () => {
      const { container } = render(<ThemeSelector />)
      const colorBoxes = container.querySelectorAll('div[class*="-color-"]')

      expect(colorBoxes.length).toBeGreaterThan(0)
    })

    it('should render correct colors for preset', () => {
      const { container } = render(<ThemeSelector />)

      // Find Modern preset section
      const presetSection = container.textContent

      expect(presetSection).toContain('Modern')
    })
  })

  describe('Accessibility', () => {
    it('should have accessible theme list', () => {
      render(<ThemeSelector />)
      expect(screen.getByText('Themes')).toBeInTheDocument()
    })

    it('should have accessible buttons', () => {
      render(<ThemeSelector />)
      const newThemeButton = screen.getByRole('button', { name: /New Theme/i })
      expect(newThemeButton).toBeInTheDocument()
    })

    it('should have accessible preset buttons', () => {
      render(<ThemeSelector />)
      const modernButton = screen.getByText('Modern')
      expect(modernButton).toBeInTheDocument()
    })
  })

  describe('Dark Mode Support', () => {
    it('should have dark mode classes', () => {
      const { container } = render(<ThemeSelector />)
      const darkElements = container.querySelectorAll('[class*="dark:"]')
      expect(darkElements.length).toBeGreaterThan(0)
    })
  })

  describe('Layout', () => {
    it('should have scrollable custom themes area', () => {
      const { container } = render(<ThemeSelector />)
      const scrollArea = container.querySelector('div.overflow-y-auto')
      expect(scrollArea).toBeInTheDocument()
    })

    it('should organize content in sections', () => {
      render(<ThemeSelector />)
      expect(screen.getByText('Presets')).toBeInTheDocument()
      expect(screen.getByText(/Custom Themes/)).toBeInTheDocument()
    })
  })

  describe('Empty States', () => {
    it('should handle no themes gracefully', () => {
      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: null,
          themes: [],
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<ThemeSelector />)
      expect(screen.getByText('No custom themes yet')).toBeInTheDocument()
    })
  })

  describe('Theme Selection', () => {
    it('should support multiple custom themes', () => {
      const theme2 = { ...mockCurrentTheme, id: '2', name: 'Theme 2' }

      ;(useThemeStore as any).mockImplementation((selector: (state: unknown) => unknown) => {
        const store = {
          currentTheme: mockCurrentTheme,
          themes: [mockCurrentTheme, theme2],
          ...mockStoreActions,
        }
        return selector(store)
      })

      render(<ThemeSelector />)
      expect(screen.getByText('My Theme')).toBeInTheDocument()
      expect(screen.getByText('Theme 2')).toBeInTheDocument()
    })

    it('should call setCurrentTheme with correct theme', () => {
      render(<ThemeSelector />)
      const themeElement = screen.getByText('My Theme').closest('div')

      fireEvent.click(themeElement!)

      expect(mockStoreActions.setCurrentTheme).toHaveBeenCalledWith(mockCurrentTheme)
    })
  })

  describe('Preset Application', () => {
    it('should maintain current theme structure when applying preset', () => {
      render(<ThemeSelector />)
      const modernButton = screen.getByText('Modern')

      fireEvent.click(modernButton)

      expect(mockStoreActions.setCurrentTheme).toHaveBeenCalledWith(
        expect.objectContaining({
          id: mockCurrentTheme.id,
          colors: expect.any(Object),
        })
      )
    })

    it('should update theme name with preset name', () => {
      render(<ThemeSelector />)
      const modernButton = screen.getByText('Modern')

      fireEvent.click(modernButton)

      expect(mockStoreActions.setCurrentTheme).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'Modern',
        })
      )
    })
  })
})
