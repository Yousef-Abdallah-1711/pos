import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Sidebar } from './Sidebar'

const state = vi.hoisted(() => ({
  isBlocked: false,
  permissionsLoading: true,
}))

vi.mock('@/hooks/useTenantStatus', () => ({
  useTenantStatus: () => ({ isBlocked: state.isBlocked }),
}))

vi.mock('@/hooks/useRbac', () => ({
  useUserPermissions: () => ({
    checkPermission: () => true,
    loading: state.permissionsLoading,
  }),
}))

vi.mock('@/components/blocked/BlockedSidebar', () => ({
  BlockedSidebar: () => <aside>Tenant blocked</aside>,
}))

vi.mock('@/components/layout/SidebarSubscriptionBadge', () => ({
  SidebarSubscriptionBadge: () => null,
}))

const props = {
  isMobileOpen: false,
  onMobileToggle: vi.fn(),
  isCollapsed: false,
  onCollapsedChange: vi.fn(),
}

describe('Sidebar hook ordering', () => {
  beforeEach(() => {
    state.isBlocked = false
    state.permissionsLoading = true
  })

  it('keeps hook order stable when permissions finish loading', () => {
    const view = render(<Sidebar {...props} />)
    expect(view.container.querySelector('aside')).not.toBeNull()

    state.permissionsLoading = false
    expect(() => view.rerender(<Sidebar {...props} />)).not.toThrow()
    expect(screen.getByText('POS Dashboard')).toBeInTheDocument()
  })

  it('keeps hook order stable when tenant access becomes blocked', () => {
    state.permissionsLoading = false
    const view = render(<Sidebar {...props} />)
    expect(screen.getByText('POS Dashboard')).toBeInTheDocument()

    state.isBlocked = true
    expect(() => view.rerender(<Sidebar {...props} />)).not.toThrow()
    expect(screen.getByText('Tenant blocked')).toBeInTheDocument()
  })
})
