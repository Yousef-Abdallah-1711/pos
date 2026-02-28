import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { AddTableForm } from './AddTableForm'

const toastError = vi.fn()

vi.mock('sonner', () => ({
  toast: {
    error: (...args: unknown[]) => toastError(...args),
  },
}))

const t = (path: string, fallback?: string) => fallback ?? path

describe('AddTableForm', () => {
  beforeEach(() => {
    toastError.mockReset()
  })

  it('shows validation error for invalid table number', async () => {
    const onSubmit = vi.fn()
    render(
      <AddTableForm
        t={t}
        isRTL={false}
        onSubmit={onSubmit}
      />
    )

    const [tableNumberInput] = screen.getAllByRole('spinbutton')
    fireEvent.change(tableNumberInput, {
      target: { value: '0' },
    })
    fireEvent.submit(screen.getByRole('button', { name: 'tables.save' }).closest('form')!)

    await waitFor(() => {
      expect(toastError).toHaveBeenCalledWith('tables.invalidTableNumber')
      expect(onSubmit).not.toHaveBeenCalled()
    })
  })

  it('shows validation error for invalid capacity', async () => {
    const onSubmit = vi.fn()
    render(
      <AddTableForm
        t={t}
        isRTL={false}
        onSubmit={onSubmit}
      />
    )

    const [, capacityInput] = screen.getAllByRole('spinbutton')
    fireEvent.change(capacityInput, {
      target: { value: '99' },
    })
    fireEvent.submit(screen.getByRole('button', { name: 'tables.save' }).closest('form')!)

    await waitFor(() => {
      expect(toastError).toHaveBeenCalledWith('tables.invalidCapacity')
      expect(onSubmit).not.toHaveBeenCalled()
    })
  })

  it('submits normalized payload when form data is valid', async () => {
    const onSubmit = vi.fn()
    render(
      <AddTableForm
        t={t}
        isRTL={false}
        onSubmit={onSubmit}
      />
    )

    const [tableNumberInput, capacityInput] = screen.getAllByRole('spinbutton')
    const textboxes = screen.getAllByRole('textbox')
    const nameInput = textboxes[0]
    const descriptionInput = textboxes[textboxes.length - 1]

    fireEvent.change(tableNumberInput, {
      target: { value: '7' },
    })
    fireEvent.change(capacityInput, {
      target: { value: '6' },
    })
    fireEvent.change(nameInput, {
      target: { value: '  Patio 7  ' },
    })
    fireEvent.change(descriptionInput, {
      target: { value: '  Near the window  ' },
    })

    fireEvent.submit(screen.getByRole('button', { name: 'tables.save' }).closest('form')!)

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledTimes(1)
      expect(onSubmit).toHaveBeenCalledWith({
        table_number: 7,
        capacity: 6,
        name: 'Patio 7',
        location: 'indoor',
        shape: 'square',
        description: 'Near the window',
      })
    })
  })
})
