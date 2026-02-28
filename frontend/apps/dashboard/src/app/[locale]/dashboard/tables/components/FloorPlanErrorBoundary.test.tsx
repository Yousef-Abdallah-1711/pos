import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { FloorPlanErrorBoundary } from './FloorPlanErrorBoundary'

const t = (path: string, fallback?: string) => fallback ?? path

function CrashComponent() {
  throw new Error('floor plan crash')
}

describe('FloorPlanErrorBoundary', () => {
  it('renders fallback when child crashes', () => {
    render(
      <FloorPlanErrorBoundary t={t}>
        <CrashComponent />
      </FloorPlanErrorBoundary>
    )

    expect(screen.getByText('Floor plan temporarily unavailable')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Retry' })).toBeTruthy()
  })
})
