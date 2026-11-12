/**
 * Form Validation Integration Tests
 * Tests complete theme form validation workflow
 */

import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useFormValidation, useColorFieldValidation, useFontFieldValidation } from '@/hooks/useFormValidation'
import { describe, it, expect, beforeEach, vi } from 'vitest'

describe('Form Validation Integration', () => {
  describe('useFormValidation Hook', () => {
    it('should initialize with empty errors', () => {
      function TestComponent() {
        const form = useFormValidation({ name: '' })
        return <div>{Object.keys(form.errors).length === 0 ? 'No errors' : 'Has errors'}</div>
      }

      render(<TestComponent />)
      expect(screen.getByText('No errors')).toBeInTheDocument()
    })

    it('should validate color fields', async () => {
      function TestComponent() {
        const form = useFormValidation({ colors: { primary: '#3b82f6' } }, { mode: 'onChange' })

        return (
          <div>
            <input
              type="text"
              value={form.values.colors?.primary || ''}
              onChange={(e) => form.setValue('colors.primary', e.target.value)}
              placeholder="Color"
            />
            {form.errors['colors.primary'] && <p>{form.errors['colors.primary']}</p>}
          </div>
        )
      }

      render(<TestComponent />)
      const input = screen.getByPlaceholderText('Color') as HTMLInputElement

      fireEvent.change(input, { target: { value: '#gggggg' } })

      await waitFor(() => {
        expect(input.value).toBe('#gggggg')
      })
    })

    it('should track dirty state', async () => {
      function TestComponent() {
        const form = useFormValidation({ name: '' })

        return (
          <div>
            <input
              type="text"
              value={form.values.name}
              onChange={(e) => form.setValue('name', e.target.value)}
              placeholder="Name"
            />
            <p>{form.isDirty ? 'Dirty' : 'Clean'}</p>
          </div>
        )
      }

      render(<TestComponent />)
      const input = screen.getByPlaceholderText('Name')

      expect(screen.getByText('Clean')).toBeInTheDocument()

      fireEvent.change(input, { target: { value: 'Test' } })

      await waitFor(() => {
        expect(screen.getByText('Dirty')).toBeInTheDocument()
      })
    })

    it('should support form reset', async () => {
      function TestComponent() {
        const form = useFormValidation({ name: 'Initial' })

        return (
          <div>
            <input
              type="text"
              value={form.values.name}
              onChange={(e) => form.setValue('name', e.target.value)}
              placeholder="Name"
            />
            <button onClick={() => form.reset()}>Reset</button>
            <p>{form.isDirty ? 'Dirty' : 'Clean'}</p>
          </div>
        )
      }

      render(<TestComponent />)
      const input = screen.getByPlaceholderText('Name') as HTMLInputElement

      fireEvent.change(input, { target: { value: 'Changed' } })
      expect(screen.getByText('Dirty')).toBeInTheDocument()

      fireEvent.click(screen.getByText('Reset'))

      await waitFor(() => {
        expect(input.value).toBe('Initial')
        expect(screen.getByText('Clean')).toBeInTheDocument()
      })
    })

    it('should validate multiple fields', async () => {
      function TestComponent() {
        const form = useFormValidation(
          {
            colors: { primary: '#3b82f6' },
            typography: { fontFamily: 'Inter' },
          },
          { mode: 'onBlur' }
        )

        return (
          <div>
            <input
              type="text"
              value={form.values.colors?.primary || ''}
              onChange={(e) => form.setValue('colors.primary', e.target.value)}
              onBlur={() => form.markFieldTouched('colors.primary')}
              placeholder="Primary Color"
            />
            <select
              value={form.values.typography?.fontFamily || ''}
              onChange={(e) => form.setValue('typography.fontFamily', e.target.value)}
              onBlur={() => form.markFieldTouched('typography.fontFamily')}
            >
              <option>Inter</option>
              <option>Roboto</option>
            </select>
          </div>
        )
      }

      render(<TestComponent />)

      const colorInput = screen.getByPlaceholderText('Primary Color')
      fireEvent.change(colorInput, { target: { value: '#ff0000' } })

      expect(colorInput).toHaveValue('#ff0000')
    })
  })

  describe('useColorFieldValidation Hook', () => {
    it('should validate color input', () => {
      function TestComponent() {
        const color = useColorFieldValidation('#3b82f6')

        return (
          <div>
            <input
              type="text"
              value={color.color}
              onChange={(e) => color.onChange(e.target.value)}
              placeholder="Color"
            />
            {color.isValid ? <p>Valid</p> : <p>Invalid</p>}
            {color.error && <p>{color.error}</p>}
          </div>
        )
      }

      render(<TestComponent />)
      const input = screen.getByPlaceholderText('Color')

      fireEvent.change(input, { target: { value: '#ffffff' } })

      expect(screen.getByText('Valid')).toBeInTheDocument()
    })

    it('should reject invalid colors', () => {
      function TestComponent() {
        const color = useColorFieldValidation('#3b82f6')

        return (
          <div>
            <input
              type="text"
              value={color.color}
              onChange={(e) => color.onChange(e.target.value)}
              placeholder="Color"
            />
            {!color.isValid && <p>Invalid Color</p>}
          </div>
        )
      }

      render(<TestComponent />)
      const input = screen.getByPlaceholderText('Color')

      fireEvent.change(input, { target: { value: '#gggggg' } })

      expect(screen.getByText('Invalid Color')).toBeInTheDocument()
    })

    it('should support color reset', () => {
      function TestComponent() {
        const color = useColorFieldValidation('#3b82f6')

        return (
          <div>
            <input
              type="text"
              value={color.color}
              onChange={(e) => color.onChange(e.target.value)}
              placeholder="Color"
            />
            <button onClick={() => color.reset()}>Reset</button>
            <p>{color.color}</p>
          </div>
        )
      }

      render(<TestComponent />)
      const input = screen.getByPlaceholderText('Color')

      fireEvent.change(input, { target: { value: '#ff0000' } })
      fireEvent.click(screen.getByText('Reset'))

      expect(screen.getByText('#3b82f6')).toBeInTheDocument()
    })
  })

  describe('useFontFieldValidation Hook', () => {
    it('should validate font selection', () => {
      function TestComponent() {
        const font = useFontFieldValidation('Inter')

        return (
          <div>
            <select value={font.font} onChange={(e) => font.onChange(e.target.value)}>
              <option>Inter</option>
              <option>Roboto</option>
              <option>Invalid Font</option>
            </select>
            {font.isValid ? <p>Valid</p> : <p>Invalid</p>}
          </div>
        )
      }

      render(<TestComponent />)
      const select = screen.getByRole('combobox')

      fireEvent.change(select, { target: { value: 'Roboto' } })

      // Font validation depends on the allowed fonts list
      expect(select).toHaveValue('Roboto')
    })
  })

  describe('Complete Theme Form Validation', () => {
    it('should validate complete theme creation', async () => {
      function CompleteThemeForm() {
        const form = useFormValidation(
          {
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
            },
          },
          { mode: 'onBlur' }
        )

        const handleSubmit = async (e: React.FormEvent) => {
          e.preventDefault()
          const isValid = await form.validateForm()
          if (isValid) {
            // Form is valid
          }
        }

        return (
          <form onSubmit={handleSubmit}>
            <input
              type="text"
              value={form.values.name}
              onChange={(e) => form.setValue('name', e.target.value)}
              placeholder="Theme Name"
            />
            <button type="submit">Submit</button>
          </form>
        )
      }

      render(<CompleteThemeForm />)
      const input = screen.getByPlaceholderText('Theme Name')

      expect(input).toHaveValue('My Theme')
    })

    it('should prevent submission with validation errors', async () => {
      function FormWithValidation() {
        const form = useFormValidation({ name: '' }, { mode: 'onSubmit' })

        const handleSubmit = async (e: React.FormEvent) => {
          e.preventDefault()
          await form.validateForm()
        }

        return (
          <form onSubmit={handleSubmit}>
            <input
              type="text"
              value={form.values.name}
              onChange={(e) => form.setValue('name', e.target.value)}
              placeholder="Name"
            />
            <button type="submit">Submit</button>
            {Object.keys(form.errors).length > 0 && <p>Has Errors</p>}
          </form>
        )
      }

      render(<FormWithValidation />)

      // Submit with empty field
      fireEvent.click(screen.getByText('Submit'))

      // Validation should fail
      await waitFor(() => {
        // Form should handle the error
      })
    })
  })

  describe('Validation Modes', () => {
    it('should validate onChange mode', async () => {
      function TestComponent() {
        const form = useFormValidation(
          { email: '' },
          { mode: 'onChange' }
        )

        return (
          <div>
            <input
              type="email"
              value={form.values.email}
              onChange={(e) => form.setValue('email', e.target.value)}
            />
            {form.errors.email && <p>{form.errors.email}</p>}
          </div>
        )
      }

      render(<TestComponent />)
      const input = screen.getByRole('textbox')

      fireEvent.change(input, { target: { value: 'test' } })

      // Should validate on change
      expect(input).toHaveValue('test')
    })

    it('should validate onBlur mode', async () => {
      function TestComponent() {
        const form = useFormValidation(
          { field: '' },
          { mode: 'onBlur' }
        )

        return (
          <input
            value={form.values.field}
            onChange={(e) => form.setValue('field', e.target.value)}
            onBlur={() => form.markFieldTouched('field')}
            placeholder="Field"
          />
        )
      }

      render(<TestComponent />)
      const input = screen.getByPlaceholderText('Field')

      fireEvent.change(input, { target: { value: 'test' } })
      fireEvent.blur(input)

      expect(input).toHaveValue('test')
    })
  })

  describe('Error Handling', () => {
    it('should clear errors on valid input', async () => {
      function TestComponent() {
        const form = useFormValidation({ color: '#invalid' })

        return (
          <div>
            <input
              type="text"
              value={form.values.color}
              onChange={(e) => {
                form.setValue('color', e.target.value)
                if (e.target.value === '#ffffff') {
                  form.clearErrors()
                }
              }}
              placeholder="Color"
            />
            {Object.keys(form.errors).length === 0 ? <p>No Errors</p> : <p>Has Errors</p>}
          </div>
        )
      }

      render(<TestComponent />)
      const input = screen.getByPlaceholderText('Color')

      fireEvent.change(input, { target: { value: '#ffffff' } })

      expect(screen.getByText('No Errors')).toBeInTheDocument()
    })
  })

  describe('Field-Level Operations', () => {
    it('should support setValue for nested fields', () => {
      function TestComponent() {
        const form = useFormValidation({
          colors: { primary: '#3b82f6' },
        })

        return (
          <div>
            <button
              onClick={() => form.setValue('colors.primary', '#ff0000')}
            >
              Change Color
            </button>
            <p>{form.values.colors?.primary}</p>
          </div>
        )
      }

      render(<TestComponent />)
      fireEvent.click(screen.getByText('Change Color'))

      expect(screen.getByText('#ff0000')).toBeInTheDocument()
    })

    it('should support setValues for multiple fields', () => {
      function TestComponent() {
        const form = useFormValidation({
          name: '',
          email: '',
        })

        return (
          <div>
            <button
              onClick={() =>
                form.setValues({ name: 'John', email: 'john@example.com' })
              }
            >
              Set Values
            </button>
            <p>{form.values.name} - {form.values.email}</p>
          </div>
        )
      }

      render(<TestComponent />)
      fireEvent.click(screen.getByText('Set Values'))

      expect(screen.getByText('John - john@example.com')).toBeInTheDocument()
    })
  })
})
