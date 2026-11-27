import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { WorkspaceAppearanceEditor } from './WorkspaceAppearanceEditor'

const { request, permission } = vi.hoisted(() => ({ request: vi.fn(), permission: vi.fn() }))
vi.mock('@/lib/authFetch', () => ({ authFetch: request }))
vi.mock('@/lib/config', () => ({ getConfig: () => ({ apiUrl: '/api/v1' }) }))
vi.mock('@/hooks/useRbac', () => ({ useUserPermissions: () => ({ checkPermission: permission, loading: false }) }))

describe('workspace appearance editor', () => {
  beforeEach(() => {
    request.mockReset()
    permission.mockReset().mockReturnValue(true)
    request.mockResolvedValueOnce({ ok: true, json: async () => ({ schemaVersion: 1, revision: 0, overrides: {} }) })
  })

  it('saves a validated draft and announces success', async () => {
    request.mockResolvedValueOnce({ ok: true, status: 200, json: async () => ({ revision: 1 }) })
    render(<WorkspaceAppearanceEditor locale="en" />)
    fireEvent.change(await screen.findByLabelText('Base font size'), { target: { value: '17' } })
    fireEvent.click(screen.getByRole('button', { name: 'Save' }))
    await screen.findByText('Workspace appearance saved.')
    const payload = JSON.parse(request.mock.calls[1][1].body)
    expect(payload.expectedRevision).toBe(0)
    expect(payload.overrides.typography.baseSize).toBe(17)
  })

  it('does not offer writes to users without settings permission', async () => {
    permission.mockReturnValue(false)
    render(<WorkspaceAppearanceEditor locale="en" />)
    const button = await screen.findByRole('button', { name: 'Reset to default' })
    expect(button).toHaveProperty('disabled', true)
    fireEvent.click(button)
    await waitFor(() => expect(request).toHaveBeenCalledTimes(1))
  })

  it('reports a revision conflict without replacing the draft', async () => {
    request.mockResolvedValueOnce({ ok: false, status: 409 })
    render(<WorkspaceAppearanceEditor locale="en" />)
    const select = await screen.findByLabelText('Base font size')
    fireEvent.change(select, { target: { value: '17' } })
    fireEvent.click(screen.getByRole('button', { name: 'Save' }))
    await screen.findByText('Appearance changed. Reload before saving.')
    expect((select as HTMLSelectElement).value).toBe('17')
  })
})
