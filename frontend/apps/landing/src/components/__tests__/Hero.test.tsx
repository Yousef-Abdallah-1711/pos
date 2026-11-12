import { render, screen } from '@testing-library/react'
import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'
import { Hero } from '@/components/Hero'

vi.mock('@/contexts/LanguageContext', () => ({
  useLanguage: () => ({
    t: {
      hero: {
        title: 'All-in-one restaurant',
        titleHighlight: 'management.',
        description: 'One platform for the work at the front and back of the house.',
        cta: 'Create your restaurant workspace',
        demo: 'See how it works',
        captureLabel: 'POOS dashboard',
        captureSection: 'Role management',
        captureAlt: 'POOS dashboard role-management screen in a synthetic sample workspace',
        captureCaption: 'Role management in the POOS dashboard.',
        sampleLabel: 'Sample workspace',
        viewCapture: 'Open full screenshot',
        productAreas: 'Restaurant management areas',
        areas: [
          { label: 'Front of house', href: '#product' },
          { label: 'Back of house', href: '#product' },
        ],
      },
    },
    language: 'en',
    dir: 'ltr',
  }),
}))

describe('Hero', () => {
  let consoleError: ReturnType<typeof vi.spyOn>

  beforeEach(() => {
    consoleError = vi.spyOn(console, 'error').mockImplementation(() => undefined)
  })

  afterEach(() => {
    consoleError.mockRestore()
  })

  it('renders repeated section links without React duplicate-key warnings', () => {
    render(<Hero />)

    expect(screen.getByRole('navigation', { name: 'Restaurant management areas' })).toHaveTextContent('Front of house')
    expect(screen.getByRole('navigation', { name: 'Restaurant management areas' })).toHaveTextContent('Back of house')
    expect(consoleError.mock.calls.flat().join(' ')).not.toMatch(/same key|unique.*key/i)
  })
})
