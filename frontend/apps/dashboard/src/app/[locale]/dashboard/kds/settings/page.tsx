'use client'

import { useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { AxiosError } from 'axios'
import { AlertCircle, ArrowLeft, Loader2, Lock, Save } from 'lucide-react'
import { useUserPermissions } from '@/hooks/useRbac'
import type { KDSDisplayDensity, KDSSettings } from '@/hooks/useKDS'
import { api } from '@/lib/api'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'

const DEFAULT_SETTINGS: KDSSettings = {
  id: 0,
  tenant_id: 0,
  restaurant_id: 0,
  auto_bump_minutes: 0,
  alert_sound_enabled: true,
  show_analytics_panel: false,
  display_density: 'comfortable',
  color_warning_minutes: 10,
  color_danger_minutes: 20,
}

const AUTO_BUMP_OPTIONS = [0, 5, 10, 15, 20] as const
const DENSITY_OPTIONS: KDSDisplayDensity[] = ['compact', 'comfortable', 'spacious']
const KDS_ANALYTICS_PANEL_PREF_KEY = 'kds-show-analytics-panel'

interface ApiResponse<T> {
  success?: boolean
  data?: T
  error?: string
  message?: string
}

const extractData = <T,>(payload: unknown, fallback: T): T => {
  if (payload && typeof payload === 'object' && 'data' in payload) {
    return ((payload as ApiResponse<T>).data ?? fallback) as T
  }
  return (payload as T) ?? fallback
}

const getErrorMessage = (error: unknown, fallback: string) => {
  if (!error || typeof error !== 'object') {
    return fallback
  }

  const axiosError = error as AxiosError<ApiResponse<unknown>>
  const responseError =
    axiosError.response?.data?.error ||
    axiosError.response?.data?.message

  return responseError || axiosError.message || fallback
}

const normalizeSettings = (input: KDSSettings): KDSSettings => ({
  ...DEFAULT_SETTINGS,
  ...input,
})

const readAnalyticsPanelPreference = () => {
  if (typeof window === 'undefined') {
    return null
  }
  try {
    const raw = window.localStorage.getItem(KDS_ANALYTICS_PANEL_PREF_KEY)
    if (raw === '1') {
      return true
    }
    if (raw === '0') {
      return false
    }
  } catch {
    // Ignore storage access issues.
  }
  return null
}

const writeAnalyticsPanelPreference = (value: boolean) => {
  if (typeof window === 'undefined') {
    return
  }
  try {
    window.localStorage.setItem(KDS_ANALYTICS_PANEL_PREF_KEY, value ? '1' : '0')
  } catch {
    // Ignore storage access issues.
  }
}

const resolveInitialSettings = (): KDSSettings => {
  const localPreference = readAnalyticsPanelPreference()
  if (localPreference === null) {
    return DEFAULT_SETTINGS
  }

  return {
    ...DEFAULT_SETTINGS,
    show_analytics_panel: localPreference,
  }
}

export default function KDSSettingsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canManageSettings = checkPermission(MODULE_IDS.KDS, PERMISSION_LEVELS.ADMIN)

  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  const [settings, setSettings] = useState<KDSSettings>(() => resolveInitialSettings())
  const [initialSettings, setInitialSettings] = useState<KDSSettings>(() => resolveInitialSettings())

  useEffect(() => {
    if (permissionLoading) {
      return
    }

    if (!canManageSettings) {
      setLoading(false)
      return
    }

    let active = true

    const loadSettings = async () => {
      try {
        setLoading(true)
        const response = await api.get('/kds/settings')
        const loadedPayload = extractData<KDSSettings>(response.data, DEFAULT_SETTINGS)
        const loaded = normalizeSettings(loadedPayload)

        const localPreference = readAnalyticsPanelPreference()
        if (localPreference !== null) {
          loaded.show_analytics_panel = localPreference
        }

        if (!active) {
          return
        }

        setSettings(loaded)
        setInitialSettings(loaded)
        writeAnalyticsPanelPreference(loaded.show_analytics_panel)
        setError(null)
      } catch (err) {
        if (!active) {
          return
        }
        setError(getErrorMessage(err, 'Failed to load KDS settings'))
      } finally {
        if (active) {
          setLoading(false)
        }
      }
    }

    void loadSettings()

    return () => {
      active = false
    }
  }, [canManageSettings, permissionLoading])

  const validationError = useMemo(() => {
    if (!AUTO_BUMP_OPTIONS.includes(settings.auto_bump_minutes as (typeof AUTO_BUMP_OPTIONS)[number])) {
      return t('kds.settingsPage.validation.autoBump', 'Auto bump must be one of: 0, 5, 10, 15, 20')
    }

    if (settings.color_warning_minutes <= 0) {
      return t('kds.settingsPage.validation.warning', 'Warning threshold must be greater than 0')
    }

    if (settings.color_danger_minutes <= settings.color_warning_minutes) {
      return t('kds.settingsPage.validation.danger', 'Danger threshold must be greater than warning threshold')
    }

    return null
  }, [settings, t])

  const isDirty = useMemo(
    () => JSON.stringify(settings) !== JSON.stringify(initialSettings),
    [settings, initialSettings]
  )

  const updateNumber = (field: 'color_warning_minutes' | 'color_danger_minutes', value: string) => {
    const parsed = Number.parseInt(value, 10)
    setSettings((current) => ({
      ...current,
      [field]: Number.isNaN(parsed) ? 0 : parsed,
    }))
  }

  const handleSave = async () => {
    if (validationError) {
      setError(validationError)
      setSuccess(null)
      return
    }

    setSaving(true)
    setError(null)
    setSuccess(null)

    try {
      const payload = {
        auto_bump_minutes: settings.auto_bump_minutes,
        alert_sound_enabled: settings.alert_sound_enabled,
        show_analytics_panel: settings.show_analytics_panel,
        display_density: settings.display_density,
        color_warning_minutes: settings.color_warning_minutes,
        color_danger_minutes: settings.color_danger_minutes,
      }

      const response = await api.put('/kds/settings', payload)
      const updated = normalizeSettings(extractData<KDSSettings>(response.data, settings))
      setSettings(updated)
      setInitialSettings(updated)
      writeAnalyticsPanelPreference(updated.show_analytics_panel)
      setSuccess(t('kds.settingsPage.saved', 'Settings saved'))
    } catch (err) {
      const axiosError = err as AxiosError<ApiResponse<unknown>>
      if (axiosError.response?.status === 400) {
        try {
          // Backward compatibility for older backend instances without show_analytics_panel support.
          const legacyPayload = {
            auto_bump_minutes: settings.auto_bump_minutes,
            alert_sound_enabled: settings.alert_sound_enabled,
            display_density: settings.display_density,
            color_warning_minutes: settings.color_warning_minutes,
            color_danger_minutes: settings.color_danger_minutes,
          }
          const response = await api.put('/kds/settings', legacyPayload)
          const updated = normalizeSettings(extractData<KDSSettings>(response.data, settings))
          updated.show_analytics_panel = settings.show_analytics_panel
          setSettings(updated)
          setInitialSettings(updated)
          writeAnalyticsPanelPreference(updated.show_analytics_panel)
          setSuccess(t('kds.settingsPage.saved', 'Settings saved'))
          return
        } catch {
          // Fall through to default error handling.
        }
      }
      setError(getErrorMessage(err, t('kds.settingsPage.saveError', 'Failed to save KDS settings')))
    } finally {
      setSaving(false)
    }
  }

  const handleReset = () => {
    setSettings(initialSettings)
    writeAnalyticsPanelPreference(initialSettings.show_analytics_panel)
    setError(null)
    setSuccess(null)
  }

  if (permissionLoading || loading) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!canManageSettings) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center px-4">
        <div className="w-full max-w-lg rounded-xl border border-red-400/40 bg-red-500/10 p-6">
          <div className={`flex items-start gap-3 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
            <Lock className="mt-0.5 h-5 w-5 flex-shrink-0 text-red-300" />
            <div>
              <h1 className="text-lg font-semibold text-red-100">
                {t('kds.accessDeniedTitle', 'Access denied')}
              </h1>
              <p className="mt-1 text-sm text-red-200">
                {t('kds.settingsPage.accessDenied', 'You do not have permission to manage KDS settings.')}
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="mx-auto w-full max-w-3xl space-y-6">
      <div className={`flex flex-wrap items-center justify-between gap-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
        <div className={isRTL ? 'text-right' : ''}>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {t('kds.settingsPage.title', 'KDS Settings')}
          </h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            {t('kds.settingsPage.subtitle', 'Configure timers, audio alerts, and display preferences')}
          </p>
        </div>

        <Link
          href={`/${locale}/dashboard/kds`}
          className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-800"
        >
          <ArrowLeft className="h-4 w-4" />
          {t('kds.settingsPage.backToKDS', 'Back to KDS')}
        </Link>
      </div>

      {(error || success) && (
        <div
          className={`rounded-xl border p-4 ${
            error
              ? 'border-red-300 bg-red-50 text-red-700 dark:border-red-500/40 dark:bg-red-500/10 dark:text-red-200'
              : 'border-emerald-300 bg-emerald-50 text-emerald-700 dark:border-emerald-500/40 dark:bg-emerald-500/10 dark:text-emerald-200'
          }`}
        >
          <div className={`flex items-start gap-2 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
            {error && <AlertCircle className="mt-0.5 h-4 w-4 flex-shrink-0" />}
            <span className="text-sm">{error || success}</span>
          </div>
        </div>
      )}

      <section className="space-y-5 rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-800 dark:bg-gray-900">
        <div className={isRTL ? 'text-right' : ''}>
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            {t('kds.settingsPage.displayAndAlerts', 'Display & Alerts')}
          </h2>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            {t('kds.settingsPage.displayAndAlertsText', 'Control how orders are shown and when alerts fire')}
          </p>
        </div>

        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <label className={`space-y-2 ${isRTL ? 'text-right' : ''}`}>
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('kds.settingsPage.autoBump', 'Auto bump timer')}
            </span>
            <select
              value={settings.auto_bump_minutes}
              onChange={(event) => {
                setSettings((current) => ({
                  ...current,
                  auto_bump_minutes: Number.parseInt(event.target.value, 10),
                }))
              }}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 dark:border-gray-700 dark:bg-gray-950 dark:text-gray-100"
            >
              {AUTO_BUMP_OPTIONS.map((minutes) => (
                <option key={minutes} value={minutes}>
                  {minutes === 0
                    ? t('kds.settingsPage.autoBumpDisabled', 'Disabled (0 minutes)')
                    : t('kds.settingsPage.autoBumpMinutes', '{minutes} minutes').replace('{minutes}', String(minutes))}
                </option>
              ))}
            </select>
          </label>

          <label className={`space-y-2 ${isRTL ? 'text-right' : ''}`}>
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('kds.settingsPage.displayDensity', 'Display density')}
            </span>
            <select
              value={settings.display_density}
              onChange={(event) => {
                setSettings((current) => ({
                  ...current,
                  display_density: event.target.value as KDSDisplayDensity,
                }))
              }}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 dark:border-gray-700 dark:bg-gray-950 dark:text-gray-100"
            >
              {DENSITY_OPTIONS.map((density) => (
                <option key={density} value={density}>
                  {t(`kds.settingsPage.density.${density}`, density)}
                </option>
              ))}
            </select>
          </label>

          <label className={`space-y-2 ${isRTL ? 'text-right' : ''}`}>
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('kds.settingsPage.warningThreshold', 'Warning threshold (minutes)')}
            </span>
            <input
              type="number"
              min={1}
              step={1}
              value={settings.color_warning_minutes}
              onChange={(event) => {
                updateNumber('color_warning_minutes', event.target.value)
              }}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 dark:border-gray-700 dark:bg-gray-950 dark:text-gray-100"
            />
          </label>

          <label className={`space-y-2 ${isRTL ? 'text-right' : ''}`}>
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('kds.settingsPage.dangerThreshold', 'Danger threshold (minutes)')}
            </span>
            <input
              type="number"
              min={2}
              step={1}
              value={settings.color_danger_minutes}
              onChange={(event) => {
                updateNumber('color_danger_minutes', event.target.value)
              }}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 dark:border-gray-700 dark:bg-gray-950 dark:text-gray-100"
            />
          </label>
        </div>

        <label className={`flex items-center gap-3 rounded-lg border border-gray-200 bg-gray-50 px-3 py-3 dark:border-gray-700 dark:bg-gray-800 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
          <input
            type="checkbox"
            checked={settings.alert_sound_enabled}
            onChange={(event) => {
              setSettings((current) => ({
                ...current,
                alert_sound_enabled: event.target.checked,
              }))
            }}
            className="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
          />
          <div>
            <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
              {t('kds.settingsPage.alertSound', 'Enable alert sound')}
            </p>
            <p className="text-xs text-gray-600 dark:text-gray-400">
              {t('kds.settingsPage.alertSoundText', 'Play a sound when a new order arrives')}
            </p>
          </div>
        </label>

        <label className={`flex items-center gap-3 rounded-lg border border-gray-200 bg-gray-50 px-3 py-3 dark:border-gray-700 dark:bg-gray-800 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
          <input
            type="checkbox"
            checked={settings.show_analytics_panel}
            onChange={(event) => {
              const checked = event.target.checked
              setSettings((current) => ({
                ...current,
                show_analytics_panel: checked,
              }))
              writeAnalyticsPanelPreference(checked)
            }}
            className="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
          />
          <div>
            <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
              {t('kds.settingsPage.showAnalytics', 'Show kitchen analytics panel')}
            </p>
            <p className="text-xs text-gray-600 dark:text-gray-400">
              {t('kds.settingsPage.showAnalyticsText', 'Disable this for a distraction-free kitchen queue view')}
            </p>
          </div>
        </label>

        <p className={`text-xs text-gray-500 dark:text-gray-400 ${isRTL ? 'text-right' : ''}`}>
          {t('kds.settingsPage.thresholdHint', 'Danger threshold must be greater than warning threshold.')}
        </p>
      </section>

      <div className={`flex flex-wrap gap-3 ${isRTL ? 'justify-start' : 'justify-end'}`}>
        <button
          type="button"
          onClick={handleReset}
          disabled={saving || !isDirty}
          className="inline-flex min-h-11 items-center rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-60 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-800"
        >
          {t('kds.settingsPage.reset', 'Reset')}
        </button>
        <button
          type="button"
          onClick={() => {
            void handleSave()
          }}
          disabled={saving || !!validationError || !isDirty}
          className="inline-flex min-h-11 items-center gap-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
          {saving
            ? t('kds.settingsPage.saving', 'Saving...')
            : t('kds.settingsPage.save', 'Save settings')}
        </button>
      </div>
    </div>
  )
}
