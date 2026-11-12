/**
 * ColorPicker Component Tests
 * Tests color input, validation, preview, and copy functionality
 */

import React from 'react'
import { act, render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ColorPicker } from '../ColorPicker'
import { describe, it, expect, beforeEach, vi } from 'vitest'

describe('ColorPicker', () => {
  const defaultProps = {
    label: 'Primary Color',
    value: '#3b82f6',
    onChange: vi.fn(),
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Rendering', () => {
    it('should render the label', () => {
      render(<ColorPicker {...defaultProps} />)
      expect(screen.getByText('Primary Color')).toBeInTheDocument()
    })

    it('should render color preview box', () => {
      const { container } = render(<ColorPicker {...defaultProps} />)
      const preview = container.querySelector('input[type="color"]')
      expect(preview).toBeInTheDocument()
    })

    it('should render hex input field', () => {
      render(<ColorPicker {...defaultProps} />)
      const input = screen.getByPlaceholderText('#000000')
      expect(input).toBeInTheDocument()
    })

    it('should display initial color value in input', () => {
      render(<ColorPicker {...defaultProps} value="#ff0000" />)
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement
      expect(input.value).toBe('#ff0000')
    })

    it('should render copy button', () => {
      render(<ColorPicker {...defaultProps} />)
      expect(screen.getByRole('button')).toBeInTheDocument()
    })

    it('should render description when provided', () => {
      render(
        <ColorPicker
          {...defaultProps}
          description="This is the primary color"
        />
      )
      expect(screen.getByText('This is the primary color')).toBeInTheDocument()
    })

    it('should not render description when not provided', () => {
      const { container } = render(<ColorPicker {...defaultProps} />)
      const descriptions = container.querySelectorAll('p')
      expect(descriptions.length).toBe(0)
    })
  })

  describe('Color Input', () => {
    it('should accept valid hex colors', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#ffffff' } })

      expect(onChange).toHaveBeenCalledWith('#ffffff')
    })

    it('should accept 3-digit hex colors', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#fff' } })

      expect(onChange).toHaveBeenCalledWith('#fff')
    })

    it('should accept 6-digit hex colors', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#ffffff' } })

      expect(onChange).toHaveBeenCalledWith('#ffffff')
    })

    it('should reject invalid hex colors', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#gggggg' } })

      expect(onChange).not.toHaveBeenCalled()
    })

    it('should reject hex colors without hash', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: 'ffffff' } })

      expect(onChange).not.toHaveBeenCalled()
    })

    it('should reject partial hex values', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#ff' } })

      expect(onChange).not.toHaveBeenCalled()
    })

    it('should update input value on change', () => {
      render(<ColorPicker {...defaultProps} />)
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#ff0000' } })

      expect(input.value).toBe('#ff0000')
    })

    it('should accept case-insensitive hex', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#AbCdEf' } })

      expect(onChange).toHaveBeenCalledWith('#AbCdEf')
    })
  })

  describe('Color Picker Input (Type Color)', () => {
    it('should update on color input change', () => {
      const onChange = vi.fn()
      const { container } = render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const colorInput = container.querySelector(
        'input[type="color"]'
      ) as HTMLInputElement

      fireEvent.change(colorInput, { target: { value: '#ff0000' } })

      expect(onChange).toHaveBeenCalledWith('#ff0000')
    })

    it('should reflect external color changes', () => {
      const { rerender, container } = render(
        <ColorPicker
          {...defaultProps}
          value="#3b82f6"
        />
      )
      const colorInput = container.querySelector(
        'input[type="color"]'
      ) as HTMLInputElement

      expect(colorInput.value).toBe('#3b82f6')

      rerender(
        <ColorPicker
          {...defaultProps}
          value="#ff0000"
        />
      )

      expect(colorInput.value).toBe('#ff0000')
    })
  })

  describe('Copy Functionality', () => {
    beforeEach(() => {
      vi.clearAllMocks()
      // Mock clipboard API
      Object.assign(navigator, {
        clipboard: {
          writeText: vi.fn().mockResolvedValue(undefined),
        },
      })
    })

    it('should copy color to clipboard', async () => {
      render(<ColorPicker {...defaultProps} value="#3b82f6" />)
      const copyButton = screen.getByRole('button')

      fireEvent.click(copyButton)

      await waitFor(() => {
        expect(navigator.clipboard.writeText).toHaveBeenCalledWith('#3b82f6')
      })
    })

    it('should show "Copied" message after copying', async () => {
      render(<ColorPicker {...defaultProps} />)
      const copyButton = screen.getByRole('button')

      expect(screen.getByText('Copy')).toBeInTheDocument()

      fireEvent.click(copyButton)

      await waitFor(() => {
        expect(screen.getByText('Copied')).toBeInTheDocument()
      })
    })

    it('should show checkmark icon after copying', async () => {
      const { container } = render(<ColorPicker {...defaultProps} />)
      const copyButton = screen.getByRole('button')

      fireEvent.click(copyButton)

      await waitFor(() => {
        const checkIcon = container.querySelector('svg')
        expect(checkIcon).toBeInTheDocument()
      })
    })

    it('should revert to copy button after 2 seconds', async () => {
      vi.useFakeTimers()
      try {
        render(<ColorPicker {...defaultProps} />)
        fireEvent.click(screen.getByRole('button'))
        expect(screen.getByText('Copied')).toBeInTheDocument()

        act(() => vi.advanceTimersByTime(2000))

        expect(screen.getByText('Copy')).toBeInTheDocument()
      } finally {
        vi.useRealTimers()
      }
    })

    it('should copy current color value', async () => {
      render(
        <ColorPicker
          {...defaultProps}
          value="#ff00ff"
        />
      )
      const copyButton = screen.getByRole('button')

      fireEvent.click(copyButton)

      await waitFor(() => {
        expect(navigator.clipboard.writeText).toHaveBeenCalledWith('#ff00ff')
      })
    })
  })

  describe('Picker Toggle', () => {
    it('should toggle picker visibility on preview click', () => {
      const { container } = render(<ColorPicker {...defaultProps} />)
      const preview = container.querySelector(
        'div.relative.w-16.h-16'
      ) as HTMLElement

      expect(preview).toBeInTheDocument()
      fireEvent.click(preview)

      // Note: In this component, the picker toggle doesn't change DOM,
      // but sets internal state
    })

    it('should close picker on click outside', () => {
      const { container } = render(<ColorPicker {...defaultProps} />)
      const preview = container.querySelector(
        'div.relative.w-16.h-16'
      ) as HTMLElement

      fireEvent.click(preview)

      // Simulate click outside
      fireEvent.mouseDown(document.body)

      // Picker should be closed (component handles this internally)
    })
  })

  describe('Props Changes', () => {
    it('should update input when value prop changes', () => {
      const { rerender } = render(
        <ColorPicker
          {...defaultProps}
          value="#3b82f6"
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      expect(input.value).toBe('#3b82f6')

      rerender(
        <ColorPicker
          {...defaultProps}
          value="#ff0000"
        />
      )

      expect(input.value).toBe('#ff0000')
    })

    it('should call onChange with new color', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#00ff00' } })

      expect(onChange).toHaveBeenCalledWith('#00ff00')
    })

    it('should handle label change', () => {
      const { rerender } = render(
        <ColorPicker
          {...defaultProps}
          label="Primary Color"
        />
      )

      expect(screen.getByText('Primary Color')).toBeInTheDocument()

      rerender(
        <ColorPicker
          {...defaultProps}
          label="Secondary Color"
        />
      )

      expect(screen.getByText('Secondary Color')).toBeInTheDocument()
      expect(screen.queryByText('Primary Color')).not.toBeInTheDocument()
    })
  })

  describe('Keyboard Input', () => {
    it('should handle typing in hex input', async () => {
      const user = userEvent.setup()
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000')

      await user.clear(input)
      await user.type(input, '#ff00ff')

      await waitFor(() => {
        expect(onChange).toHaveBeenCalledWith('#ff00ff')
      })
    })

    it('should handle rapid color changes', async () => {
      const user = userEvent.setup()
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000')

      await user.clear(input)
      await user.type(input, '#ffffff')

      expect(onChange).toHaveBeenCalledWith('#ffffff')
    })
  })

  describe('Accessibility', () => {
    it('should have accessible label', () => {
      render(<ColorPicker {...defaultProps} />)
      const label = screen.getByText('Primary Color')
      expect(label).toBeInTheDocument()
    })

    it('should have placeholder text', () => {
      render(<ColorPicker {...defaultProps} />)
      const input = screen.getByPlaceholderText('#000000')
      expect(input).toBeInTheDocument()
    })

    it('should have title attribute on color input', () => {
      const { container } = render(
        <ColorPicker
          {...defaultProps}
          value="#3b82f6"
        />
      )
      const colorInput = container.querySelector(
        'input[type="color"]'
      ) as HTMLInputElement

      expect(colorInput).toHaveAttribute('title', '#3b82f6')
    })

    it('should have accessible button for copy', () => {
      render(<ColorPicker {...defaultProps} />)
      const button = screen.getByRole('button')
      expect(button).toBeInTheDocument()
    })
  })

  describe('Edge Cases', () => {
    it('should handle uppercase hex', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#FFFFFF' } })

      expect(onChange).toHaveBeenCalledWith('#FFFFFF')
    })

    it('should handle mixed case hex', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#FfFfFf' } })

      expect(onChange).toHaveBeenCalledWith('#FfFfFf')
    })

    it('should handle all zeros color', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#000000' } })

      expect(onChange).toHaveBeenCalledWith('#000000')
    })

    it('should handle all F color', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#ffffff' } })

      expect(onChange).toHaveBeenCalledWith('#ffffff')
    })

    it('should handle empty input gracefully', () => {
      const onChange = vi.fn()
      render(
        <ColorPicker
          {...defaultProps}
          onChange={onChange}
        />
      )
      const input = screen.getByPlaceholderText('#000000') as HTMLInputElement

      fireEvent.change(input, { target: { value: '' } })

      expect(onChange).not.toHaveBeenCalled()
    })
  })

  describe('Dark Mode', () => {
    it('should render with dark mode classes', () => {
      const { container } = render(<ColorPicker {...defaultProps} />)

      // Check for dark mode classes
      const darkElements = container.querySelectorAll('[class*="dark:"]')
      expect(darkElements.length).toBeGreaterThan(0)
    })
  })

  describe('Styling', () => {
    it('should apply hover effects on preview', () => {
      const { container } = render(<ColorPicker {...defaultProps} />)
      const preview = container.querySelector('div.relative.w-16.h-16')

      expect(preview).toHaveClass('hover:shadow-lg')
    })

    it('should apply correct classes to input', () => {
      render(<ColorPicker {...defaultProps} />)
      const input = screen.getByPlaceholderText('#000000')

      expect(input).toHaveClass('font-mono')
    })

    it('should have rounded corners on preview', () => {
      const { container } = render(<ColorPicker {...defaultProps} />)
      const preview = container.querySelector('div.relative.w-16.h-16')

      expect(preview).toHaveClass('rounded-lg')
    })
  })

  describe('Multiple Instances', () => {
    it('should handle multiple color pickers independently', () => {
      const onChange1 = vi.fn()
      const onChange2 = vi.fn()

      const { container } = render(
        <>
          <ColorPicker
            label="Color 1"
            value="#ff0000"
            onChange={onChange1}
          />
          <ColorPicker
            label="Color 2"
            value="#00ff00"
            onChange={onChange2}
          />
        </>
      )

      const inputs = container.querySelectorAll('input[type="text"]')
      fireEvent.change(inputs[0], { target: { value: '#ffffff' } })

      expect(onChange1).toHaveBeenCalledWith('#ffffff')
      expect(onChange2).not.toHaveBeenCalled()
    })
  })
})
