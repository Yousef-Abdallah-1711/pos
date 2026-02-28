import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { TableStatusBadge } from './TableStatusBadge'

describe('TableStatusBadge', () => {
  it('renders provided label', () => {
    render(<TableStatusBadge status="available" label="Available" />)

    expect(screen.getByText('Available')).toBeTruthy()
  })

  it('applies occupied status classes', () => {
    render(<TableStatusBadge status="occupied" label="Occupied" />)

    const badge = screen.getByText('Occupied')
    expect(badge.className).toContain('bg-red-100')
    expect(badge.className).toContain('text-red-800')
  })
})
