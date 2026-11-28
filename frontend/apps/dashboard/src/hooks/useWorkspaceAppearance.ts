'use client'

import { useEffect, useState } from 'react'
import { authFetch } from '@/lib/authFetch'
import { getConfig } from '@/lib/config'
import { debugError } from '@/lib/debug'
import { resolveWorkspaceAppearance, type WorkspaceAppearance } from '@/lib/workspaceAppearance'

const FALLBACK = resolveWorkspaceAppearance({})
type ScopedAppearance = { scope: string; value: WorkspaceAppearance }

export function useWorkspaceAppearance(tenantId: number | undefined, userId: number | undefined): WorkspaceAppearance {
  const scope = tenantId && userId ? `${tenantId}:${userId}` : ''
  const [loaded, setLoaded] = useState<ScopedAppearance | null>(null)

  useEffect(() => {
    if (!scope) return
    const controller = new AbortController()
    let requestVersion = 0
    async function load() {
      const currentVersion = ++requestVersion
      try {
        const response = await authFetch(`${getConfig().apiUrl}/workspace-appearance`, { signal: controller.signal, cache: 'no-store' })
        if (!response.ok) throw new Error(`Workspace appearance request failed (${response.status})`)
        const payload: unknown = await response.json()
        if (!payload || typeof payload !== 'object' || !('schemaVersion' in payload) || payload.schemaVersion !== 1 || !('overrides' in payload)) throw new Error('Invalid workspace appearance response')
        const value = resolveWorkspaceAppearance(payload.overrides)
        if (!controller.signal.aborted && currentVersion === requestVersion) setLoaded({ scope, value })
      } catch (error) {
        if (!controller.signal.aborted) debugError('[WorkspaceAppearance] Using defaults:', error)
      }
    }
    const refresh = () => { void load() }
    void load()
    window.addEventListener('workspace-appearance-updated', refresh)
    return () => { controller.abort(); window.removeEventListener('workspace-appearance-updated', refresh) }
  }, [scope])

  return loaded?.scope === scope ? loaded.value : FALLBACK
}
