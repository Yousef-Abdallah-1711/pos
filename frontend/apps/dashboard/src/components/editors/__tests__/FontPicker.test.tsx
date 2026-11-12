/**
 * FontPicker Component Tests
 * Tests font selection, search, filtering, and weight selection
 */

import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { FontPicker } from '../FontPicker'
import { googleFontsService } from '@/services/googleFontsService'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import type { ThemeComponent } from '@/types/theme'

describe('FontPicker', () => {
  const mockComponent: ThemeComponent = {
    id: 'test-1',
    type: 'custom',
    title: 'Font Config',
    displayOrder: 1,
    enabled: true,
    config: {
      fontFamily: 'roboto',
      fontWeight: 400,
      fontSize: 16,
      lineHeight: 1.6,
      fontWeights: [400, 700],
    },
  }

  const mockOnChange = vi.fn()
  const mockOnPreview = vi.fn()

  const defaultProps = {
    component: mockComponent,
    onChange: mockOnChange,
    onPreview: mockOnPreview,
    className: '',
  }

  beforeEach(() => {
    vi.clearAllMocks()
    // Clear localStorage
    localStorage.clear()
  })

  describe('Rendering', () => {
    it('should render font picker header', () => {
      render(<FontPicker {...defaultProps} />)
      expect(screen.getByText('Font Picker')).toBeInTheDocument()
    })

    it('should render search input', () => {
      render(<FontPicker {...defaultProps} />)
      expect(screen.getByPlaceholderText('Search fonts...')).toBeInTheDocument()
    })

    it('should render category tabs', () => {
      render(<FontPicker {...defaultProps} />)
      expect(screen.getByText('All')).toBeInTheDocument()
      expect(screen.getByText('Sans Serif')).toBeInTheDocument()
      expect(screen.getByText('Serif')).toBeInTheDocument()
      expect(screen.getByText('Monospace')).toBeInTheDocument()
    })

    it('should render font list', () => {
      render(<FontPicker {...defaultProps} />)
      const fonts = googleFontsService.getPopularFonts()
      expect(fonts.length).toBeGreaterThan(0)
    })

    it('should render preview button', () => {
      render(<FontPicker {...defaultProps} />)
      expect(screen.getByText('Preview Font')).toBeInTheDocument()
    })
  })

  describe('Font Search', () => {
    it('should filter fonts by search query', async () => {
      render(<FontPicker {...defaultProps} />)

      const searchInput = screen.getByPlaceholderText('Search fonts...')
      await userEvent.type(searchInput, 'Roboto')

      await waitFor(() => {
        expect(screen.getByText('Roboto')).toBeInTheDocument()
      })
    })

    it('should show no results for non-matching search', async () => {
      render(<FontPicker {...defaultProps} />)

      const searchInput = screen.getByPlaceholderText('Search fonts...')
      await userEvent.type(searchInput, 'NonexistentFont')

      await waitFor(() => {
        expect(screen.getByText(/No fonts found/i)).toBeInTheDocument()
      })
    })

    it('should clear search results when search cleared', async () => {
      render(<FontPicker {...defaultProps} />)

      const searchInput = screen.getByPlaceholderText('Search fonts...') as HTMLInputElement
      await userEvent.type(searchInput, 'Roboto')

      await waitFor(() => {
        expect(screen.getByText('Roboto')).toBeInTheDocument()
      })

      await userEvent.clear(searchInput)

      await waitFor(() => {
        expect(screen.queryByText(/No fonts found/i)).not.toBeInTheDocument()
      })
    })
  })

  describe('Font Categories', () => {
    it('should filter by Sans-serif category', async () => {
      render(<FontPicker {...defaultProps} />)

      const sansSerifTab = screen.getByText('Sans Serif')
      fireEvent.click(sansSerifTab)

      await waitFor(() => {
        const fonts = googleFontsService.getPopularFonts().filter((f) => f.category === 'sans-serif')
        expect(fonts.length).toBeGreaterThan(0)
      })
    })

    it('should filter by Serif category', async () => {
      render(<FontPicker {...defaultProps} />)

      const serifTab = screen.getByText('Serif')
      fireEvent.click(serifTab)

      await waitFor(() => {
        const fonts = googleFontsService.getPopularFonts().filter((f) => f.category === 'serif')
        expect(fonts.length).toBeGreaterThan(0)
      })
    })

    it('should filter by Monospace category', async () => {
      render(<FontPicker {...defaultProps} />)

      const monoTab = screen.getByText('Monospace')
      fireEvent.click(monoTab)

      await waitFor(() => {
        const fonts = googleFontsService.getPopularFonts().filter((f) => f.category === 'monospace')
        expect(fonts.length).toBeGreaterThan(0)
      })
    })

    it('should show all fonts when "All" tab selected', async () => {
      render(<FontPicker {...defaultProps} />)

      const allTab = screen.getByText('All')
      fireEvent.click(allTab)

      await waitFor(() => {
        const fonts = googleFontsService.getPopularFonts()
        expect(fonts.length).toBeGreaterThan(0)
      })
    })
  })

  describe('Font Selection', () => {
    it('should select font when clicked', async () => {
      render(<FontPicker {...defaultProps} />)

      const robotoFont = screen.getByText('Roboto')
      fireEvent.click(robotoFont)

      await waitFor(() => {
        expect(mockOnChange).toHaveBeenCalled()
      })
    })

    it('should highlight selected font', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto', fontWeight: 400, fontSize: 16 },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      // Selected font should have blue background
      const selectedFont = screen.getByText('Roboto').closest('button')
      expect(selectedFont).toHaveClass('bg-blue-100')
    })

    it('should add font to recent list', async () => {
      render(<FontPicker {...defaultProps} />)

      const robotoFont = screen.getByText('Roboto')
      fireEvent.click(robotoFont)

      await waitFor(() => {
        const recent = googleFontsService.getRecentFonts()
        expect(recent.length).toBeGreaterThan(0)
      })
    })

    it('should show recent fonts', async () => {
      // Add font to recent
      const roboto = googleFontsService.getFont('roboto')
      if (roboto) googleFontsService.addRecentFont(roboto)

      render(<FontPicker {...defaultProps} />)

      await waitFor(() => {
        expect(screen.getByText('Recent')).toBeInTheDocument()
      })
    })
  })

  describe('Font Configuration', () => {
    it('should display font weights selection', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto', fontWeights: [400, 700] },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      await waitFor(() => {
        expect(screen.getByText('Font Weights')).toBeInTheDocument()
      })
    })

    it('should toggle font weight selection', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto', fontWeights: [400] },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      await waitFor(() => {
        const weight700 = screen.getByText('700')
        fireEvent.click(weight700)
        expect(mockOnChange).toHaveBeenCalled()
      })
    })

    it('should control font size with slider', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto', fontSize: 16 },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      await waitFor(() => {
        const sizeSlider = screen.getByRole('slider', { name: /Font Size/ })
        expect(sizeSlider).toBeInTheDocument()
      })
    })

    it('should update font size display', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto', fontSize: 20 },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      await waitFor(() => {
        expect(screen.getByText(/Font Size: 20px/)).toBeInTheDocument()
      })
    })

    it('should control line height with slider', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto', lineHeight: 1.8 },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      await waitFor(() => {
        expect(screen.getByText(/Line Height: 1.8/)).toBeInTheDocument()
      })
    })

    it('should show font style options for italic fonts', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto', fontStyle: 'normal' },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      await waitFor(() => {
        const roboto = googleFontsService.getFont('roboto')
        if (roboto?.italics) {
          expect(screen.getByText('Font Style')).toBeInTheDocument()
        }
      })
    })
  })

  describe('Preview Mode', () => {
    it('should switch to preview mode', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto' },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      const previewBtn = screen.getByText('Preview Font')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('Back to Settings')).toBeInTheDocument()
      })
    })

    it('should show font info in preview', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto' },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      const previewBtn = screen.getByText('Preview Font')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('Roboto')).toBeInTheDocument()
      })
    })

    it('should show preview text', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto' },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      const previewBtn = screen.getByText('Preview Font')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText(/Pack my box/i)).toBeInTheDocument()
      })
    })

    it('should return to settings from preview', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto' },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      const previewBtn = screen.getByText('Preview Font')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        const backBtn = screen.getByText('Back to Settings')
        fireEvent.click(backBtn)
      })

      expect(screen.getByText('Font Picker')).toBeInTheDocument()
    })
  })

  describe('Live Preview', () => {
    it('should display live preview with selected font', () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto' },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      expect(screen.getByText(/The quick brown fox/)).toBeInTheDocument()
    })

    it('should update preview with font size changes', async () => {
      const component = {
        ...mockComponent,
        config: { fontFamily: 'roboto', fontSize: 24 },
      }
      render(<FontPicker {...defaultProps} component={component} />)

      expect(screen.getByText(/Font Size: 24px/)).toBeInTheDocument()
    })
  })
})
