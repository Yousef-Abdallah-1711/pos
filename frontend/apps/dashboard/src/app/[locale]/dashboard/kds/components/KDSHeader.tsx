'use client'

import Link from 'next/link'
import { Maximize2, Minimize2, Moon, Settings2, Sun, Volume2, VolumeX, Wifi, WifiOff } from 'lucide-react'
import { useTheme } from 'next-themes'
import { KDSFilter, KDSStats } from '@/stores/kdsStore'
import { usePreferencesStore } from '@/stores/preferencesStore'

interface KDSHeaderProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  stats: KDSStats | null
  filter: KDSFilter
  soundEnabled: boolean
  isFullscreen: boolean
  isFullscreenSupported?: boolean
  isConnected: boolean
  isOffline?: boolean
  queuedActionsCount?: number
  isReplayingOfflineQueue?: boolean
  canManageSettings?: boolean
  settingsHref?: string
  isRefreshing?: boolean
  onFilterChange: (filter: KDSFilter) => void
  onToggleSound: () => void
  onToggleFullscreen: () => void
}

const FILTERS: KDSFilter[] = ['all', 'new', 'in_progress', 'ready', 'completed']

const getFilterLabel = (filter: KDSFilter, t: (path: string, fallback?: string) => string) => {
  switch (filter) {
    case 'all':
      return t('kds.filters.all', 'All')
    case 'new':
      return t('kds.filters.new', 'New')
    case 'in_progress':
      return t('kds.filters.in_progress', 'In Progress')
    case 'ready':
      return t('kds.filters.ready', 'Ready')
    case 'completed':
      return t('kds.filters.completed', 'Completed')
    default:
      return filter
  }
}

export function KDSHeader({
  locale,
  t,
  stats,
  filter,
  soundEnabled,
  isFullscreen,
  isFullscreenSupported = true,
  isConnected,
  isOffline = false,
  queuedActionsCount = 0,
  isReplayingOfflineQueue = false,
  canManageSettings = false,
  settingsHref = '/dashboard/kds/settings',
  isRefreshing = false,
  onFilterChange,
  onToggleSound,
  onToggleFullscreen,
}: KDSHeaderProps) {
  const { resolvedTheme, setTheme } = useTheme()
  const updateTheme = usePreferencesStore((state) => state.updateTheme)
  const isRTL = locale === 'ar'
  const isDarkMode = resolvedTheme === 'dark'
  const safeStats = stats || {
    active_orders: 0,
    avg_prep_minutes: 0,
    completed_today: 0,
    new_orders: 0,
    in_progress: 0,
    ready: 0,
  }
  const toggleTheme = () => {
    const nextTheme = isDarkMode ? 'light' : 'dark'
    setTheme(nextTheme)
    updateTheme(nextTheme)
  }

  return (
    <header className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-slate-700 dark:bg-slate-900">
      <div className={`mb-4 flex flex-wrap items-center justify-between gap-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
        <div className={`${isRTL ? 'text-right' : 'text-left'}`}>
          <h1 className="text-xl font-bold text-gray-900 dark:text-white">
            {t('kds.title', 'Kitchen Display')}
          </h1>
          <p className="text-sm text-gray-600 dark:text-slate-400">
            {isRefreshing
              ? t('kds.refreshing', 'Refreshing...')
              : t('kds.subtitle', 'Live kitchen queue')}
          </p>
        </div>

        <div className={`flex items-center gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
          <span
            className={`inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs font-medium ${
              isOffline
                ? 'border-amber-300 bg-amber-100 text-amber-800 dark:border-amber-400/40 dark:bg-amber-500/20 dark:text-amber-200'
                : isConnected
                  ? 'border-emerald-300 bg-emerald-100 text-emerald-800 dark:border-emerald-400/40 dark:bg-emerald-500/20 dark:text-emerald-200'
                  : 'border-red-300 bg-red-100 text-red-800 dark:border-red-400/40 dark:bg-red-500/20 dark:text-red-200'
            }`}
          >
            {isOffline || !isConnected ? <WifiOff className="h-3.5 w-3.5" /> : <Wifi className="h-3.5 w-3.5" />}
            {isOffline
              ? t('kds.offline', 'Offline')
              : isConnected
                ? t('kds.connected', 'Connected')
                : t('kds.disconnected', 'Disconnected')}
          </span>

          <button
            type="button"
            onClick={onToggleSound}
            className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-gray-300 bg-white text-gray-800 transition hover:bg-gray-100 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700"
            title={soundEnabled ? t('kds.soundOn', 'Sound on') : t('kds.soundOff', 'Sound off')}
            aria-label={soundEnabled ? t('kds.soundOn', 'Sound on') : t('kds.soundOff', 'Sound off')}
          >
            {soundEnabled ? <Volume2 className="h-4 w-4" /> : <VolumeX className="h-4 w-4" />}
          </button>

          <button
            type="button"
            onClick={toggleTheme}
            className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-gray-300 bg-white text-gray-800 transition hover:bg-gray-100 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700"
            title={isDarkMode ? t('common.lightMode', 'Switch to light mode') : t('common.darkMode', 'Switch to dark mode')}
            aria-label={isDarkMode ? t('common.lightMode', 'Switch to light mode') : t('common.darkMode', 'Switch to dark mode')}
          >
            {isDarkMode ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
          </button>

          {isFullscreenSupported && (
            <button
              type="button"
              onClick={onToggleFullscreen}
              className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-gray-300 bg-white text-gray-800 transition hover:bg-gray-100 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700"
              title={isFullscreen ? t('kds.exitFullscreen', 'Exit fullscreen') : t('kds.fullscreen', 'Fullscreen')}
              aria-label={isFullscreen ? t('kds.exitFullscreen', 'Exit fullscreen') : t('kds.fullscreen', 'Fullscreen')}
            >
              {isFullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
            </button>
          )}

          {canManageSettings && (
            <Link
              href={`/${locale}${settingsHref}`}
              className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-gray-300 bg-white text-gray-800 transition hover:bg-gray-100 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700"
              title={t('kds.settings', 'Settings')}
              aria-label={t('kds.settings', 'Settings')}
            >
              <Settings2 className="h-4 w-4" />
            </Link>
          )}
        </div>
      </div>

      {isOffline && (
        <div className={`mb-4 rounded-lg border border-amber-300 bg-amber-50 px-3 py-2 text-xs text-amber-800 dark:border-amber-400/40 dark:bg-amber-500/10 dark:text-amber-100 ${isRTL ? 'text-right' : ''}`}>
          {t('kds.offlineBanner', 'Offline mode: actions are queued and will sync when connection returns.')}
          {queuedActionsCount > 0 ? ` ${queuedActionsCount} ${t('kds.queuedSuffix', 'queued')}.` : ''}
          {isReplayingOfflineQueue ? ` ${t('kds.syncingQueued', 'Syncing queued actions...')}` : ''}
        </div>
      )}

      <div className={`mb-4 grid grid-cols-2 gap-2 md:grid-cols-4 ${isRTL ? 'text-right' : 'text-left'}`}>
        <div className="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-slate-700 dark:bg-slate-800/80">
          <p className="text-xs uppercase tracking-wide text-gray-500 dark:text-slate-400">
            {t('kds.activeOrders', 'Active Orders')}
          </p>
          <p className="text-xl font-bold text-gray-900 dark:text-white">{safeStats.active_orders}</p>
        </div>
        <div className="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-slate-700 dark:bg-slate-800/80">
          <p className="text-xs uppercase tracking-wide text-gray-500 dark:text-slate-400">
            {t('kds.newOrders', 'New')}
          </p>
          <p className="text-xl font-bold text-sky-700 dark:text-sky-300">{safeStats.new_orders}</p>
        </div>
        <div className="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-slate-700 dark:bg-slate-800/80">
          <p className="text-xs uppercase tracking-wide text-gray-500 dark:text-slate-400">
            {t('kds.inProgress', 'In Progress')}
          </p>
          <p className="text-xl font-bold text-amber-700 dark:text-amber-300">{safeStats.in_progress}</p>
        </div>
        <div className="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-slate-700 dark:bg-slate-800/80">
          <p className="text-xs uppercase tracking-wide text-gray-500 dark:text-slate-400">
            {t('kds.avgPrepTime', 'Avg Prep')}
          </p>
          <p className="text-xl font-bold text-emerald-700 dark:text-emerald-300">
            {safeStats.avg_prep_minutes.toFixed(1)}
            {t('kds.minutesSuffix', 'm')}
          </p>
        </div>
      </div>

      <div className={`flex flex-wrap gap-2 ${isRTL ? 'justify-end' : ''}`}>
        {FILTERS.map((value) => (
          <button
            key={value}
            type="button"
            onClick={() => onFilterChange(value)}
            aria-pressed={filter === value}
            className={`touch-manipulation min-h-12 rounded-lg px-4 py-2 text-sm font-medium transition ${
              filter === value
                ? 'bg-primary-600 text-white'
                : 'border border-gray-300 bg-white text-gray-700 hover:bg-gray-100 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200 dark:hover:bg-slate-700'
            }`}
          >
            {getFilterLabel(value, t)}
          </button>
        ))}
      </div>
    </header>
  )
}
