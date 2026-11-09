/**
 * HeaderEditor Component Tests
 * Tests header configuration, navigation management, and preview
 */

import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { HeaderEditor } from '../HeaderEditor'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import type { ThemeComponent } from '@/types/theme'

describe('HeaderEditor', () => {
  const mockComponent: ThemeComponent = {
    id: 'test-1',
    type: 'custom',
    title: 'Header Config',
    displayOrder: 1,
    enabled: true,
    config: {
      navigationItems: [
        { id: '1', label: 'Home', href: '/', order: 1 },
        { id: '2', label: 'About', href: '/about', order: 2 },
      ],
      navPosition: 'right',
      navAlignment: 'horizontal',
      backgroundColor: '#ffffff',
      textColor: '#000000',
      height: 80,
      padding: 16,
      showLogo: true,
      logoHeight: 40,
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
  })

  describe('Rendering', () => {
    it('should render header editor title', () => {
      render(<HeaderEditor {...defaultProps} />)
      expect(screen.getByText('Header Configuration')).toBeInTheDocument()
    })

    it('should render logo settings section', () => {
      render(<HeaderEditor {...defaultProps} />)
      expect(screen.getByText('Logo Settings')).toBeInTheDocument()
    })

    it('should render navigation section', () => {
      render(<HeaderEditor {...defaultProps} />)
      expect(screen.getByText('Navigation Items')).toBeInTheDocument()
    })

    it('should render styling section', () => {
      render(<HeaderEditor {...defaultProps} />)
      expect(screen.getByText('Styling')).toBeInTheDocument()
    })

    it('should render preview button', () => {
      render(<HeaderEditor {...defaultProps} />)
      expect(screen.getByText('Preview Header')).toBeInTheDocument()
    })
  })

  describe('Logo Configuration', () => {
    it('should toggle show logo', () => {
      render(<HeaderEditor {...defaultProps} />)
      const checkbox = screen.getByLabelText(/Show Logo/i)
      fireEvent.click(checkbox)
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should display logo inputs when enabled', () => {
      const component = { ...mockComponent, config: { ...mockComponent.config, showLogo: true } }
      render(<HeaderEditor {...defaultProps} component={component} />)
      expect(screen.getByPlaceholderText('https://example.com/logo.png')).toBeInTheDocument()
    })

    it('should hide logo inputs when disabled', () => {
      const component = { ...mockComponent, config: { ...mockComponent.config, showLogo: false } }
      render(<HeaderEditor {...defaultProps} component={component} />)
      expect(screen.queryByPlaceholderText('https://example.com/logo.png')).not.toBeInTheDocument()
    })

    it('should update logo URL', async () => {
      const component = { ...mockComponent, config: { ...mockComponent.config, showLogo: true } }
      render(<HeaderEditor {...defaultProps} component={component} />)

      const input = screen.getByPlaceholderText('https://example.com/logo.png') as HTMLInputElement
      await userEvent.clear(input)
      await userEvent.type(input, 'https://example.com/newlogo.png')

      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should update logo height with slider', async () => {
      const component = { ...mockComponent, config: { ...mockComponent.config, showLogo: true } }
      render(<HeaderEditor {...defaultProps} component={component} />)

      expect(screen.getByText(/Logo Height:/)).toBeInTheDocument()
    })
  })

  describe('Navigation Management', () => {
    it('should display navigation items', () => {
      render(<HeaderEditor {...defaultProps} />)
      expect(screen.getByText('Home')).toBeInTheDocument()
      expect(screen.getByText('About')).toBeInTheDocument()
    })

    it('should add new navigation item', async () => {
      render(<HeaderEditor {...defaultProps} />)

      const addBtn = screen.getByText(/Add Item/)
      fireEvent.click(addBtn)

      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should delete navigation item', async () => {
      render(<HeaderEditor {...defaultProps} />)

      fireEvent.click(screen.getByRole('button', { name: 'Delete Home' }))
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should move navigation item up', async () => {
      render(<HeaderEditor {...defaultProps} />)

      fireEvent.click(screen.getByRole('button', { name: 'Move About up' }))
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should move navigation item down', async () => {
      render(<HeaderEditor {...defaultProps} />)

      fireEvent.click(screen.getByRole('button', { name: 'Move Home down' }))
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should change navigation position', async () => {
      render(<HeaderEditor {...defaultProps} />)

      const positionSelect = screen.getByDisplayValue('Right') as HTMLSelectElement
      fireEvent.change(positionSelect, { target: { value: 'left' } })

      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should change navigation alignment', async () => {
      render(<HeaderEditor {...defaultProps} />)

      const alignSelect = screen.getByDisplayValue('Horizontal') as HTMLSelectElement
      fireEvent.change(alignSelect, { target: { value: 'vertical' } })

      expect(mockOnChange).toHaveBeenCalled()
    })
  })

  describe('Styling', () => {
    it('should have background color picker', () => {
      render(<HeaderEditor {...defaultProps} />)
      expect(screen.getByText('Background Color')).toBeInTheDocument()
    })

    it('should have text color picker', () => {
      render(<HeaderEditor {...defaultProps} />)
      expect(screen.getByText('Text Color')).toBeInTheDocument()
    })

    it('should have height slider', () => {
      render(<HeaderEditor {...defaultProps} />)
      expect(screen.getByText(/^Height:/)).toBeInTheDocument()
    })

    it('should have padding slider', () => {
      render(<HeaderEditor {...defaultProps} />)
      expect(screen.getByText(/Padding:/)).toBeInTheDocument()
    })

    it('should toggle sticky header', () => {
      render(<HeaderEditor {...defaultProps} />)
      const stickyCheckbox = screen.getByLabelText(/Sticky Header/i)
      fireEvent.click(stickyCheckbox)
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should toggle show shadow', () => {
      render(<HeaderEditor {...defaultProps} />)
      const shadowCheckbox = screen.getByLabelText(/Show Shadow/i)
      fireEvent.click(shadowCheckbox)
      expect(mockOnChange).toHaveBeenCalled()
    })
  })

  describe('Preview', () => {
    it('should switch to preview mode', async () => {
      render(<HeaderEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview Header')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('Back to Settings')).toBeInTheDocument()
      })
    })

    it('should show desktop preview', async () => {
      render(<HeaderEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview Header')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('Desktop Preview')).toBeInTheDocument()
      })
    })

    it('should show mobile preview', async () => {
      render(<HeaderEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview Header')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('Mobile Preview')).toBeInTheDocument()
      })
    })

    it('should display navigation items in preview', async () => {
      render(<HeaderEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview Header')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('Home')).toBeInTheDocument()
        expect(screen.getByText('About')).toBeInTheDocument()
      })
    })

    it('should return to settings from preview', async () => {
      render(<HeaderEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview Header')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        const backBtn = screen.getByText('Back to Settings')
        fireEvent.click(backBtn)
      })

      expect(screen.getByText('Header Configuration')).toBeInTheDocument()
    })
  })
})
