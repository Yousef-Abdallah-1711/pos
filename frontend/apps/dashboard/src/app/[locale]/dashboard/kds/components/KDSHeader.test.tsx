import React from 'react'
import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { KDSHeader } from './KDSHeader'
import type { KDSFilter, KDSStats } from '@/stores/kdsStore'

const mockSetTheme = vi.fn()
const mockUpdateTheme = vi.fn()

vi.mock('next/link', () => ({
  default: ({
    href,
    children,
    ...props
  }: {
    href: string
    children: React.ReactNode
    [key: string]: unknown
  }) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
}))

vi.mock('next-themes', () => ({
  useTheme: () => ({
    resolvedTheme: 'dark',
    setTheme: mockSetTheme,
  }),
}))

vi.mock('@/stores/preferencesStore', () => ({
  usePreferencesStore: (selector: (state: { updateTheme: (theme: 'light' | 'dark' | 'system') => void }) => unknown) =>
    selector({
      updateTheme: mockUpdateTheme,
    }),
}))

const t = (_path: string, fallback?: string) => fallback ?? ''

const stats: KDSStats = {
  active_orders: 12,
  avg_prep_minutes: 8.6,
  completed_today: 45,
  new_orders: 3,
  in_progress: 5,
  ready: 4,
}

const renderHeader = (overrides: Partial<React.ComponentProps<typeof KDSHeader>> = {}) => {
  const onFilterChange = vi.fn()
  const onToggleSound = vi.fn()
  const onToggleFullscreen = vi.fn()

  render(
    <KDSHeader
      locale="en"
      t={t}
      stats={stats}
      filter={'all' as KDSFilter}
      soundEnabled
      isFullscreen={false}
      isConnected
      onFilterChange={onFilterChange}
      onToggleSound={onToggleSound}
      onToggleFullscreen={onToggleFullscreen}
      {...overrides}
    />
  )

  return { onFilterChange, onToggleSound, onToggleFullscreen }
}

describe('KDSHeader', () => {
  beforeEach(() => {
    mockSetTheme.mockClear()
    mockUpdateTheme.mockClear()
  })

  it('renders title, counters, and connection state', () => {
    renderHeader()

    expect(screen.getByText('Kitchen Display')).toBeTruthy()
    expect(screen.getByText('Live kitchen queue')).toBeTruthy()
    expect(screen.getByText('Connected')).toBeTruthy()
    expect(screen.getByText('12')).toBeTruthy()
    expect(screen.getByText('3')).toBeTruthy()
    expect(screen.getByText('5')).toBeTruthy()
    expect(screen.getByText('8.6m')).toBeTruthy()
  })

  it('triggers filter callback when a filter button is clicked', () => {
    const { onFilterChange } = renderHeader()

    fireEvent.click(screen.getByRole('button', { name: 'Ready' }))
    expect(onFilterChange).toHaveBeenCalledTimes(1)
    expect(onFilterChange).toHaveBeenCalledWith('ready')
  })

  it('triggers sound and fullscreen toggles', () => {
    const { onToggleSound, onToggleFullscreen } = renderHeader()

    fireEvent.click(screen.getByRole('button', { name: /Sound on/i }))
    fireEvent.click(screen.getByRole('button', { name: /Fullscreen/i }))

    expect(onToggleSound).toHaveBeenCalledTimes(1)
    expect(onToggleFullscreen).toHaveBeenCalledTimes(1)
  })

  it('toggles theme and syncs preferences store', () => {
    renderHeader()

    fireEvent.click(screen.getByRole('button', { name: /Switch to light mode/i }))

    expect(mockSetTheme).toHaveBeenCalledWith('light')
    expect(mockUpdateTheme).toHaveBeenCalledWith('light')
  })

  it('shows settings link with locale prefix when enabled', () => {
    renderHeader({
      canManageSettings: true,
      settingsHref: '/dashboard/kds/settings',
    })

    const settingsLink = screen.getByRole('link', { name: /Settings/i }) as HTMLAnchorElement
    expect(settingsLink.getAttribute('href')).toBe('/en/dashboard/kds/settings')
  })

  it('shows offline queue banner details while offline', () => {
    renderHeader({
      isConnected: false,
      isOffline: true,
      queuedActionsCount: 2,
      isReplayingOfflineQueue: true,
    })

    const banner = screen.getByText(/Offline mode: actions are queued and will sync when connection returns\./i)
    expect(banner).toBeTruthy()
    expect(banner.textContent).toContain('2 queued.')
    expect(banner.textContent).toContain('Syncing queued actions...')
  })
})
