'use client'

import { useEffect, useState } from 'react'
import { authFetch } from '@/lib/authFetch'
import { getConfig } from '@/lib/config'
import { resolveWorkspaceAppearance, workspaceAppearanceVariables, type WorkspaceAppearance } from '@/lib/workspaceAppearance'
import { useUserPermissions } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'

type Mode = 'light' | 'dark'
type Token = keyof WorkspaceAppearance['colors']['light']
const TOKENS: Token[] = ['background', 'surface', 'foreground', 'muted', 'border', 'primary']
const defaults = () => resolveWorkspaceAppearance({})

export function WorkspaceAppearanceEditor({ locale }: { locale: 'en' | 'ar' }) {
  const ar = locale === 'ar'
  const { checkPermission, loading } = useUserPermissions()
  const canEdit = !loading && checkPermission(MODULE_IDS.SETTINGS, PERMISSION_LEVELS.WRITE)
  const [draft, setDraft] = useState<WorkspaceAppearance>(defaults)
  const [published, setPublished] = useState<WorkspaceAppearance>(defaults)
  const [revision, setRevision] = useState<number | null>(null)
  const [mode, setMode] = useState<Mode>('light')
  const [status, setStatus] = useState<'loading' | 'idle' | 'saving' | 'error'>('loading')
  const [message, setMessage] = useState('')
  const dirty = JSON.stringify(draft) !== JSON.stringify(published)

  useEffect(() => {
    const controller = new AbortController()
    async function load() {
      try {
        const response = await authFetch(`${getConfig().apiUrl}/workspace-appearance`, { signal: controller.signal, cache: 'no-store' })
        if (!response.ok) throw new Error('load')
        const value: unknown = await response.json()
        if (!value || typeof value !== 'object' || !('schemaVersion' in value) || value.schemaVersion !== 1 || !('revision' in value) || typeof value.revision !== 'number' || !Number.isSafeInteger(value.revision) || value.revision < 0 || !('overrides' in value)) throw new Error('response')
        const appearance = resolveWorkspaceAppearance(value.overrides)
        if (!controller.signal.aborted) { setDraft(appearance); setPublished(appearance); setRevision(value.revision); setStatus('idle') }
      } catch { if (!controller.signal.aborted) { setStatus('error'); setMessage(ar ? 'تعذر تحميل مظهر مساحة العمل.' : 'Could not load workspace appearance.') } }
    }
    void load()
    return () => controller.abort()
  }, [ar])

  useEffect(() => {
    const warn = (event: BeforeUnloadEvent) => { if (dirty) event.preventDefault() }
    window.addEventListener('beforeunload', warn)
    return () => window.removeEventListener('beforeunload', warn)
  }, [dirty])

  const save = async (reset = false) => {
    if (!canEdit || revision === null) return
    try {
      const candidate = reset ? defaults() : resolveWorkspaceAppearance({ colors: draft.colors, typography: draft.typography })
      setStatus('saving'); setMessage('')
      const response = await authFetch(`${getConfig().apiUrl}/workspace-appearance`, {
        method: 'PUT', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ expectedRevision: revision, overrides: reset ? {} : { colors: candidate.colors, typography: candidate.typography } }),
      })
      if (response.status === 409) throw new Error('conflict')
      if (!response.ok) throw new Error('save')
      const result: unknown = await response.json()
      if (!result || typeof result !== 'object' || !('revision' in result) || typeof result.revision !== 'number' || result.revision !== revision + 1) throw new Error('response')
      setRevision(result.revision); setDraft(candidate); setPublished(candidate); setStatus('idle')
      setMessage(ar ? 'تم حفظ مظهر مساحة العمل.' : 'Workspace appearance saved.')
      window.dispatchEvent(new Event('workspace-appearance-updated'))
    } catch (error) {
      setStatus('error')
      setMessage(error instanceof Error && error.message === 'conflict' ? (ar ? 'تغير المظهر. أعد تحميل الصفحة قبل الحفظ.' : 'Appearance changed. Reload before saving.') : (ar ? 'تعذر الحفظ. تحقق من الألوان والتباين.' : 'Could not save. Check the colors and contrast.'))
    }
  }

  if (status === 'loading') return <div aria-live="polite">{ar ? 'جارٍ تحميل مظهر مساحة العمل…' : 'Loading workspace appearance…'}</div>
  if (revision === null) return <p role="alert">{message}</p>
  const preview = workspaceAppearanceVariables(draft, mode)
  return <section className="space-y-4 rounded-lg border border-gray-300 p-5 dark:border-gray-700" aria-labelledby="workspace-appearance-title">
    <h2 id="workspace-appearance-title" className="text-xl font-semibold">{ar ? 'مظهر مساحة العمل' : 'Workspace appearance'}</h2>
    <p className="text-sm">{ar ? 'تنطبق هذه الإعدادات على مساحة العمل، ولا تغيّر مظهر المتجر.' : 'These settings apply to the workspace and do not change the storefront.'}</p>
    {!canEdit && <p>{ar ? 'يمكن لمسؤولي الإعدادات فقط تعديل المظهر.' : 'Only users with settings write permission can edit appearance.'}</p>}
    <label className="block">{ar ? 'معاينة الوضع' : 'Preview mode'} <select value={mode} onChange={event => setMode(event.target.value as Mode)} className="min-h-12 rounded border p-2 text-gray-900"><option value="light">{ar ? 'فاتح' : 'Light'}</option><option value="dark">{ar ? 'داكن' : 'Dark'}</option></select></label>
    <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">{TOKENS.map(token => <label key={token} className="flex min-h-12 items-center justify-between gap-3"><span>{token}</span><input type="color" disabled={!canEdit} value={draft.colors[mode][token]} onChange={event => setDraft(current => ({ ...current, colors: { ...current.colors, [mode]: { ...current.colors[mode], [token]: event.target.value } } }))} className="h-12 w-12" /></label>)}</div>
    <label className="block">{ar ? 'حجم الخط الأساسي' : 'Base font size'} <select disabled={!canEdit} value={draft.typography.baseSize} onChange={event => setDraft(current => ({ ...current, typography: { ...current.typography, baseSize: Number(event.target.value) } }))} className="min-h-12 rounded border p-2 text-gray-900">{[14, 15, 16, 17, 18].map(size => <option key={size} value={size}>{size}px</option>)}</select></label>
    <div className="grid gap-3 sm:grid-cols-3">
      <label>{ar ? 'الخط الإنجليزي' : 'English font'} <select disabled={!canEdit} value={draft.typography.english} onChange={event => setDraft(current => ({ ...current, typography: { ...current.typography, english: event.target.value as WorkspaceAppearance['typography']['english'] } }))} className="min-h-12 rounded border p-2 text-gray-900"><option value="manrope">Manrope</option><option value="system">System</option></select></label>
      <label>{ar ? 'الخط العربي' : 'Arabic font'} <select disabled={!canEdit} value={draft.typography.arabic} onChange={event => setDraft(current => ({ ...current, typography: { ...current.typography, arabic: event.target.value as WorkspaceAppearance['typography']['arabic'] } }))} className="min-h-12 rounded border p-2 text-gray-900"><option value="ibm-plex-arabic">IBM Plex Sans Arabic</option><option value="system">System</option></select></label>
      <label>{ar ? 'خط العناوين الإنجليزية' : 'English heading font'} <select disabled={!canEdit} value={draft.typography.heading} onChange={event => setDraft(current => ({ ...current, typography: { ...current.typography, heading: event.target.value as WorkspaceAppearance['typography']['heading'] } }))} className="min-h-12 rounded border p-2 text-gray-900"><option value="fraunces">Fraunces</option><option value="manrope">Manrope</option></select></label>
    </div>
    <div style={{ ...preview, background: 'var(--workspace-background)', color: 'var(--workspace-foreground)', fontSize: 'var(--workspace-base-size)', fontFamily: ar ? 'var(--workspace-font-ar)' : 'var(--workspace-font-en)' } as React.CSSProperties} className="rounded-lg border p-4"><strong>{ar ? 'معاينة مساحة العمل' : 'Workspace preview'}</strong><p>{ar ? 'الطلبات والمخزون والتقارير' : 'Orders, inventory and reports'}</p><span style={{ color: 'var(--workspace-primary)' }}>{ar ? 'إجراء أساسي' : 'Primary action'}</span></div>
    <p aria-live="polite">{message}</p>
    <div className="flex flex-wrap gap-2"><button type="button" disabled={!canEdit || !dirty || status === 'saving'} onClick={() => void save()} className="min-h-12 rounded bg-red-800 px-4 text-white disabled:opacity-50">{ar ? 'حفظ' : 'Save'}</button><button type="button" disabled={!canEdit || status === 'saving'} onClick={() => void save(true)} className="min-h-12 rounded border px-4 disabled:opacity-50">{ar ? 'استعادة الافتراضي' : 'Reset to default'}</button><button type="button" disabled={!dirty || status === 'saving'} onClick={() => setDraft(published)} className="min-h-12 rounded border px-4 disabled:opacity-50">{ar ? 'تجاهل' : 'Discard'}</button></div>
  </section>
}
