/**
 * FooterEditor Component Tests
 * Tests footer configuration, sections, and social links
 */

import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { FooterEditor } from '../FooterEditor'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import type { ThemeComponent } from '@/types/theme'

describe('FooterEditor', () => {
  const mockComponent: ThemeComponent = {
    id: 'test-1',
    type: 'custom',
    title: 'Footer Config',
    displayOrder: 1,
    enabled: true,
    config: {
      companyName: 'My Restaurant',
      companyDescription: 'Great food',
      phone: '+1-555-1234',
      email: 'info@restaurant.com',
      address: '123 Main St, City',
      footerSections: [],
      socialLinks: [],
      backgroundColor: '#1f2937',
      textColor: '#ffffff',
      linkColor: '#3b82f6',
      padding: 48,
      copyrightText: '© 2024 My Restaurant',
      showLegal: true,
      legalLinks: [
        { label: 'Privacy', href: '/privacy' },
        { label: 'Terms', href: '/terms' },
      ],
      columns: 3,
      layout: 'expanded',
      showBackToTop: true,
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
    it('should render footer editor title', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText('Footer Configuration')).toBeInTheDocument()
    })

    it('should render company information section', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText('Company Information')).toBeInTheDocument()
    })

    it('should render footer sections manager', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText('Footer Sections')).toBeInTheDocument()
    })

    it('should render social links manager', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText('Social Media Links')).toBeInTheDocument()
    })

    it('should render legal section', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText('Legal & Copyright')).toBeInTheDocument()
    })

    it('should render styling section', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText('Styling')).toBeInTheDocument()
    })
  })

  describe('Company Information', () => {
    it('should display company name input', () => {
      render(<FooterEditor {...defaultProps} />)
      const input = screen.getByPlaceholderText('Your Company') as HTMLInputElement
      expect(input.value).toBe('My Restaurant')
    })

    it('should update company name', async () => {
      render(<FooterEditor {...defaultProps} />)
      const input = screen.getByPlaceholderText('Your Company')
      await userEvent.clear(input)
      await userEvent.type(input, 'New Restaurant')
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should display company description', () => {
      render(<FooterEditor {...defaultProps} />)
      const textarea = screen.getByPlaceholderText('Brief company description') as HTMLTextAreaElement
      expect(textarea.value).toBe('Great food')
    })

    it('should display phone input', () => {
      render(<FooterEditor {...defaultProps} />)
      const input = screen.getByPlaceholderText('+1 (555) 123-4567') as HTMLInputElement
      expect(input.value).toBe('+1-555-1234')
    })

    it('should display email input', () => {
      render(<FooterEditor {...defaultProps} />)
      const input = screen.getByPlaceholderText('contact@example.com') as HTMLInputElement
      expect(input.value).toBe('info@restaurant.com')
    })

    it('should display address input', () => {
      render(<FooterEditor {...defaultProps} />)
      const textarea = screen.getByPlaceholderText('123 Main St, City, State 12345') as HTMLTextAreaElement
      expect(textarea.value).toBe('123 Main St, City')
    })
  })

  describe('Footer Sections', () => {
    it('should display add section button', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText(/Add Section/)).toBeInTheDocument()
    })

    it('should add new section', async () => {
      render(<FooterEditor {...defaultProps} />)
      const addBtn = screen.getByText(/Add Section/)
      fireEvent.click(addBtn)
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should display footer sections with link count', () => {
      const component = {
        ...mockComponent,
        config: {
          ...mockComponent.config,
          footerSections: [
            { id: '1', title: 'Products', links: [{ label: 'Burgers', href: '/burgers' }], order: 1 },
          ],
        },
      }
      render(<FooterEditor {...defaultProps} component={component} />)
      expect(screen.getByText('Products')).toBeInTheDocument()
    })
  })

  describe('Social Links', () => {
    it('should display add social link button', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText(/Add Link/)).toBeInTheDocument()
    })

    it('should add new social link', async () => {
      render(<FooterEditor {...defaultProps} />)
      fireEvent.click(screen.getByRole('button', { name: /Add Link/ }))
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should display social link platform select', () => {
      const component = {
        ...mockComponent,
        config: {
          ...mockComponent.config,
          socialLinks: [{ id: '1', platform: 'facebook', url: 'https://facebook.com', order: 1 }],
        },
      }
      render(<FooterEditor {...defaultProps} component={component} />)
      expect(screen.getByRole('combobox')).toHaveValue('facebook')
    })

    it('should change social platform', async () => {
      const component = {
        ...mockComponent,
        config: {
          ...mockComponent.config,
          socialLinks: [{ id: '1', platform: 'facebook', url: 'https://facebook.com', order: 1 }],
        },
      }
      render(<FooterEditor {...defaultProps} component={component} />)
      const select = screen.getByRole('combobox') as HTMLSelectElement
      fireEvent.change(select, { target: { value: 'twitter' } })
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should update social link URL', async () => {
      const component = {
        ...mockComponent,
        config: {
          ...mockComponent.config,
          socialLinks: [{ id: '1', platform: 'facebook', url: 'https://facebook.com', order: 1 }],
        },
      }
      render(<FooterEditor {...defaultProps} component={component} />)
      const input = screen.getByDisplayValue('https://facebook.com') as HTMLInputElement
      await userEvent.clear(input)
      await userEvent.type(input, 'https://newfacebook.com')
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should delete social link', async () => {
      const component = {
        ...mockComponent,
        config: {
          ...mockComponent.config,
          socialLinks: [{ id: '1', platform: 'facebook', url: 'https://facebook.com', order: 1 }],
        },
      }
      render(<FooterEditor {...defaultProps} component={component} />)
      fireEvent.click(screen.getByRole('button', { name: 'Delete facebook social link' }))
      expect(mockOnChange).toHaveBeenCalled()
    })
  })

  describe('Legal Section', () => {
    it('should display copyright text input', () => {
      render(<FooterEditor {...defaultProps} />)
      const input = screen.getByDisplayValue('© 2024 My Restaurant') as HTMLInputElement
      expect(input).toBeInTheDocument()
    })

    it('should toggle show legal links', () => {
      render(<FooterEditor {...defaultProps} />)
      const checkbox = screen.getByLabelText(/Show Legal Links/i)
      fireEvent.click(checkbox)
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should update copyright text', async () => {
      render(<FooterEditor {...defaultProps} />)
      const input = screen.getByDisplayValue('© 2024 My Restaurant')
      await userEvent.clear(input)
      await userEvent.type(input, '© 2025 New Name')
      expect(mockOnChange).toHaveBeenCalled()
    })
  })

  describe('Styling', () => {
    it('should have background color picker', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText('Background')).toBeInTheDocument()
    })

    it('should have text color picker', () => {
      render(<FooterEditor {...defaultProps} />)
      const labels = screen.getAllByText('Text Color')
      expect(labels.length).toBeGreaterThan(0)
    })

    it('should have link color picker', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText('Link Color')).toBeInTheDocument()
    })

    it('should have padding slider', () => {
      render(<FooterEditor {...defaultProps} />)
      expect(screen.getByText(/Padding:/)).toBeInTheDocument()
    })

    it('should toggle back to top button', () => {
      render(<FooterEditor {...defaultProps} />)
      const checkbox = screen.getByLabelText(/Show Back to Top/i)
      fireEvent.click(checkbox)
      expect(mockOnChange).toHaveBeenCalled()
    })
  })

  describe('Preview', () => {
    it('should switch to preview mode', async () => {
      render(<FooterEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview Footer')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('Back to Settings')).toBeInTheDocument()
      })
    })

    it('should display company info in preview', async () => {
      render(<FooterEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview Footer')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('My Restaurant')).toBeInTheDocument()
        expect(screen.getByText('Great food')).toBeInTheDocument()
      })
    })

    it('should display copyright in preview', async () => {
      render(<FooterEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview Footer')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('© 2024 My Restaurant')).toBeInTheDocument()
      })
    })

    it('should return to settings from preview', async () => {
      render(<FooterEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview Footer')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        const backBtn = screen.getByText('Back to Settings')
        fireEvent.click(backBtn)
      })

      expect(screen.getByText('Footer Configuration')).toBeInTheDocument()
    })
  })
})
