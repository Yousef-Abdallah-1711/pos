/**
 * CustomCSSEditor Component Tests
 * Tests CSS validation, editing, preview, and error detection
 */

import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { CustomCSSEditor } from '../CustomCSSEditor'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import type { ThemeComponent } from '@/types/theme'

describe('CustomCSSEditor', () => {
  const mockComponent: ThemeComponent = {
    id: 'test-1',
    type: 'custom',
    title: 'Custom CSS',
    displayOrder: 1,
    enabled: true,
    config: {
      customCss: '.button { color: blue; }',
      enableCustomCss: true,
      cssValidationEnabled: true,
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
    it('should render editor header', () => {
      render(<CustomCSSEditor {...defaultProps} />)
      expect(screen.getByText('Custom CSS Editor')).toBeInTheDocument()
    })

    it('should render enable toggle', () => {
      render(<CustomCSSEditor {...defaultProps} />)
      const checkbox = screen.getByLabelText(/Enable Custom CSS/i)
      expect(checkbox).toBeInTheDocument()
    })

    it('should render CSS textarea', () => {
      render(<CustomCSSEditor {...defaultProps} />)
      const textarea = screen.getByPlaceholderText(/Add your custom CSS here/)
      expect(textarea).toBeInTheDocument()
    })

    it('should render toolbar buttons', () => {
      render(<CustomCSSEditor {...defaultProps} />)
      expect(screen.getByText('Validate')).toBeInTheDocument()
      expect(screen.getByText(/Copy/i)).toBeInTheDocument()
      expect(screen.getByText('Clear')).toBeInTheDocument()
    })

    it('should render preview button', () => {
      render(<CustomCSSEditor {...defaultProps} />)
      expect(screen.getByText('Preview CSS')).toBeInTheDocument()
    })
  })

  describe('CSS Validation', () => {
    it('should validate CSS with !important warning', async () => {
      const component = {
        ...mockComponent,
        config: { customCss: '.button { color: red !important; }', enableCustomCss: true, cssValidationEnabled: true },
      }
      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const validateBtn = screen.getByText('Validate')
      fireEvent.click(validateBtn)

      await waitFor(() => {
        expect(screen.getByText(/Using !important is discouraged/i)).toBeInTheDocument()
      })
    })

    it('should detect inline styles error', async () => {
      const component = {
        ...mockComponent,
        config: {
          customCss: '<div style="color: red">Text</div>',
          enableCustomCss: true,
          cssValidationEnabled: true,
        },
      }
      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const validateBtn = screen.getByText('Validate')
      fireEvent.click(validateBtn)

      await waitFor(() => {
        expect(screen.getByText(/Inline styles should not appear/i)).toBeInTheDocument()
      })
    })

    it('should validate color values', async () => {
      const component = {
        ...mockComponent,
        config: {
          customCss: '.button { color: notacolor; }',
          enableCustomCss: true,
          cssValidationEnabled: true,
        },
      }
      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const validateBtn = screen.getByText('Validate')
      fireEvent.click(validateBtn)

      await waitFor(() => {
        expect(screen.getByText(/invalid color value/i)).toBeInTheDocument()
      })
    })

    it('should pass validation for valid CSS', async () => {
      const component = {
        ...mockComponent,
        config: {
          customCss: '.button { color: #ff0000; background: blue; }',
          enableCustomCss: true,
          cssValidationEnabled: true,
        },
      }
      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const validateBtn = screen.getByText('Validate')
      fireEvent.click(validateBtn)

      await waitFor(() => {
        expect(screen.getByText(/CSS looks good/i)).toBeInTheDocument()
      })
    })
  })

  describe('CSS Editing', () => {
    it('should update CSS on textarea change', async () => {
      render(<CustomCSSEditor {...defaultProps} />)

      const textarea = screen.getByPlaceholderText(/Add your custom CSS here/) as HTMLTextAreaElement
      fireEvent.change(textarea, { target: { value: '.new { color: green; }' } })

      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should disable textarea when CSS is disabled', () => {
      const component = {
        ...mockComponent,
        config: { customCss: '', enableCustomCss: false, cssValidationEnabled: false },
      }
      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const textarea = screen.getByPlaceholderText(/Add your custom CSS here/) as HTMLTextAreaElement
      expect(textarea).toBeDisabled()
    })

    it('should enable textarea when CSS is enabled', () => {
      const component = {
        ...mockComponent,
        config: { customCss: '', enableCustomCss: true, cssValidationEnabled: false },
      }
      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const textarea = screen.getByPlaceholderText(/Add your custom CSS here/) as HTMLTextAreaElement
      expect(textarea).not.toBeDisabled()
    })
  })

  describe('Toolbar Actions', () => {
    it('should toggle CSS enable/disable', async () => {
      render(<CustomCSSEditor {...defaultProps} />)

      const checkbox = screen.getByLabelText(/Enable Custom CSS/i)
      fireEvent.click(checkbox)

      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should toggle validation enable/disable', async () => {
      render(<CustomCSSEditor {...defaultProps} />)

      const checkbox = screen.getByLabelText(/Validate CSS/i)
      fireEvent.click(checkbox)

      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should copy CSS to clipboard', async () => {
      const component = {
        ...mockComponent,
        config: { customCss: '.test { color: red; }', enableCustomCss: true },
      }

      // Mock clipboard
      Object.assign(navigator, {
        clipboard: {
          writeText: vi.fn(),
        },
      })

      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const copyBtn = screen.getByText(/Copy/)
      fireEvent.click(copyBtn)

      await waitFor(() => {
        expect(navigator.clipboard.writeText).toHaveBeenCalledWith('.test { color: red; }')
      })
    })

    it('should show copy success message', async () => {
      const component = {
        ...mockComponent,
        config: { customCss: '.test { color: red; }', enableCustomCss: true },
      }

      Object.assign(navigator, {
        clipboard: {
          writeText: vi.fn(),
        },
      })

      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const copyBtn = screen.getByText(/Copy/)
      fireEvent.click(copyBtn)

      await waitFor(() => {
        expect(screen.getByText('Copied!')).toBeInTheDocument()
      })
    })

    it('should clear CSS with confirmation', async () => {
      const component = {
        ...mockComponent,
        config: { customCss: '.test { color: red; }', enableCustomCss: true },
      }

      // Mock window.confirm
      window.confirm = vi.fn(() => true)

      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const clearBtn = screen.getByText('Clear')
      fireEvent.click(clearBtn)

      expect(window.confirm).toHaveBeenCalled()
      expect(mockOnChange).toHaveBeenCalled()
    })

    it('should not clear CSS if confirmation declined', async () => {
      const component = {
        ...mockComponent,
        config: { customCss: '.test { color: red; }', enableCustomCss: true },
      }

      window.confirm = vi.fn(() => false)

      const { rerender } = render(<CustomCSSEditor {...defaultProps} component={component} />)

      const clearBtn = screen.getByText('Clear')
      fireEvent.click(clearBtn)

      expect(window.confirm).toHaveBeenCalled()
      // onChange should not be called for clear action if declined
      const clearCallCount = mockOnChange.mock.calls.length
      expect(clearCallCount).toBe(0)
    })
  })

  describe('Preview Mode', () => {
    it('should switch to preview mode when button clicked', async () => {
      render(<CustomCSSEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview CSS')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('Back to Editor')).toBeInTheDocument()
      })
    })

    it('should show CSS code in preview', async () => {
      const component = {
        ...mockComponent,
        config: { customCss: '.button { color: blue; }', enableCustomCss: true },
      }
      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const previewBtn = screen.getByText('Preview CSS')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText('.button { color: blue; }')).toBeInTheDocument()
      })
    })

    it('should return to settings from preview', async () => {
      render(<CustomCSSEditor {...defaultProps} />)

      const previewBtn = screen.getByText('Preview CSS')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        const backBtn = screen.getByText('Back to Editor')
        fireEvent.click(backBtn)
      })

      expect(screen.getByText('Custom CSS Editor')).toBeInTheDocument()
    })

    it('should show disabled message when CSS is disabled', async () => {
      const component = {
        ...mockComponent,
        config: { customCss: '.button { color: blue; }', enableCustomCss: false },
      }
      render(<CustomCSSEditor {...defaultProps} component={component} />)

      const previewBtn = screen.getByText('Preview CSS')
      fireEvent.click(previewBtn)

      await waitFor(() => {
        expect(screen.getByText(/Custom CSS is disabled/i)).toBeInTheDocument()
      })
    })
  })

  describe('Character Counter', () => {
    it('should display character count', () => {
      const component = {
        ...mockComponent,
        config: { customCss: '.test { }', enableCustomCss: true },
      }
      render(<CustomCSSEditor {...defaultProps} component={component} />)

      expect(screen.getByText('(9 characters)')).toBeInTheDocument()
    })

    it('should update character count on change', async () => {
      render(<CustomCSSEditor {...defaultProps} />)

      const textarea = screen.getByPlaceholderText(/Add your custom CSS here/) as HTMLTextAreaElement
      fireEvent.change(textarea, { target: { value: '.long { color: red; }' } })

      expect(screen.getByText(/21 characters/i)).toBeInTheDocument()
    })
  })
})
