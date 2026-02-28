'use client'

import { useEffect, useRef, useState } from 'react'
import { usePathname } from 'next/navigation'
import { AlertCircle, Loader2, Lock } from 'lucide-react'
import { MODULE_IDS, PERMISSION_LEVELS, hasPermissionLevel } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { useKDSAnalytics } from '@/hooks/useKDSAnalytics'
import { useKDSAudio } from '@/hooks/useKDSAudio'
import { useUserPermissions } from '@/hooks/useRbac'
import { useKDS } from '@/hooks/useKDS'
import { useKDSWebSocket } from '@/hooks/useKDSWebSocket'
import { KDSAnalyticsPanel } from './components/KDSAnalyticsPanel'
import { KDSErrorBoundary } from './components/KDSErrorBoundary'
import { KDSGrid } from './components/KDSGrid'
import { KDSHeader } from './components/KDSHeader'

export default function KDSPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'
  const [isFullscreenSupported, setIsFullscreenSupported] = useState(true)

  const { permissions, loading: permissionsLoading, error: permissionsError } = useUserPermissions()
  const kdsPermissionLevel = permissions?.[MODULE_IDS.KDS]
  const canViewKDS = hasPermissionLevel(kdsPermissionLevel, PERMISSION_LEVELS.READ)
  const canBumpOrders = hasPermissionLevel(kdsPermissionLevel, PERMISSION_LEVELS.WRITE)
  const canRecallOrders = hasPermissionLevel(kdsPermissionLevel, PERMISSION_LEVELS.DELETE)
  const canManageKDSSettings = hasPermissionLevel(kdsPermissionLevel, PERMISSION_LEVELS.ADMIN)
  const canViewLoading = permissionsLoading
  const canBumpLoading = permissionsLoading
  const canRecallLoading = permissionsLoading
  const canManageSettingsLoading = permissionsLoading

  const {
    orders,
    filteredOrders,
    stats,
    settings,
    filter,
    soundEnabled,
    isFullscreen,
    isConnected,
    isOffline,
    queuedActionsCount,
    isReplayingOfflineQueue,
    loading,
    isRefreshing,
    error,
    pendingOrderIds,
    setFilter,
    toggleSound,
    setFullscreen,
    refresh,
    bumpOrder,
    recallOrder,
  } = useKDS()
  const {
    analytics,
    loading: analyticsLoading,
    error: analyticsError,
    analyticsAvailable,
    refreshAnalytics,
  } = useKDSAnalytics(canManageKDSSettings && !canManageSettingsLoading && settings.show_analytics_panel && !loading)
  const canViewAnalytics = canManageKDSSettings && !canManageSettingsLoading && settings.show_analytics_panel
  const { playNewOrderAlert } = useKDSAudio()
  useKDSWebSocket({ onOrderNew: playNewOrderAlert })
  const autoBumpRunningRef = useRef(false)
  const showShellLoading = loading && filteredOrders.length === 0

  useEffect(() => {
    setIsFullscreenSupported(Boolean(document.fullscreenEnabled))

    const onFullscreenChange = () => {
      setFullscreen(Boolean(document.fullscreenElement))
    }

    document.addEventListener('fullscreenchange', onFullscreenChange)
    return () => {
      document.removeEventListener('fullscreenchange', onFullscreenChange)
    }
  }, [setFullscreen])

  useEffect(() => {
    if (canBumpLoading || !canBumpOrders) {
      return
    }

    if (isOffline) {
      return
    }

    const autoBumpMinutes = settings.auto_bump_minutes
    if (autoBumpMinutes <= 0) {
      return
    }

    const thresholdMs = autoBumpMinutes * 60_000

    const tryAutoBump = async () => {
      if (autoBumpRunningRef.current) {
        return
      }

      autoBumpRunningRef.current = true
      try {
        const now = Date.now()
        const pendingSet = new Set(pendingOrderIds)
        const nextCandidate = orders
          .filter((order) => !pendingSet.has(order.id))
          .filter((order) => ['new', 'in_progress', 'ready', 'recalled'].includes(order.status))
          .sort((left, right) => new Date(left.received_at).getTime() - new Date(right.received_at).getTime())
          .find((order) => now - new Date(order.received_at).getTime() >= thresholdMs)

        if (!nextCandidate) {
          return
        }

        await bumpOrder(nextCandidate.id)
      } finally {
        autoBumpRunningRef.current = false
      }
    }

    const intervalId = window.setInterval(() => {
      void tryAutoBump()
    }, 15_000)

    void tryAutoBump()

    return () => {
      window.clearInterval(intervalId)
    }
  }, [
    orders,
    pendingOrderIds,
    settings.auto_bump_minutes,
    isOffline,
    canBumpLoading,
    canBumpOrders,
    bumpOrder,
  ])

  const handleFullscreenToggle = async () => {
    try {
      if (!document.fullscreenEnabled) {
        return
      }

      if (!document.fullscreenElement) {
        await document.documentElement.requestFullscreen()
      } else {
        await document.exitFullscreen()
      }
    } catch {
      // Keep local fullscreen preference in sync with actual browser state only.
      setFullscreen(Boolean(document.fullscreenElement))
    }
  }

  const refreshAll = async () => {
    await refresh()
    if (canViewAnalytics) {
      await refreshAnalytics()
    }
  }

  if (!canViewLoading && permissionsError && !canViewKDS) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4 dark:bg-slate-950">
        <div className="w-full max-w-lg rounded-xl border border-amber-300 bg-amber-50 p-6 dark:border-amber-400/40 dark:bg-amber-500/10">
          <div className={`flex items-start gap-3 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
            <AlertCircle className="mt-0.5 h-5 w-5 flex-shrink-0 text-amber-700 dark:text-amber-200" />
            <div className="space-y-3">
              <div>
                <h1 className="text-lg font-semibold text-amber-800 dark:text-amber-100">
                  {t('kds.permissionCheckFailedTitle', 'Permission check unavailable')}
                </h1>
                <p className="mt-1 text-sm text-amber-700 dark:text-amber-200">
                  {t('kds.permissionCheckFailedText', 'Unable to verify KDS access right now. Please retry in a few seconds.')}
                </p>
              </div>
              <button
                type="button"
                onClick={() => {
                  window.location.reload()
                }}
                className="touch-manipulation min-h-10 rounded-md border border-amber-300 px-3 py-1.5 text-xs font-medium text-amber-800 transition hover:bg-amber-100 dark:border-amber-300/50 dark:text-amber-100 dark:hover:bg-amber-500/20"
              >
                {t('kds.retry', 'Retry')}
              </button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!canViewLoading && !canViewKDS) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4 dark:bg-slate-950">
        <div className="w-full max-w-lg rounded-xl border border-red-300 bg-red-50 p-6 dark:border-red-400/40 dark:bg-red-500/10">
          <div className={`flex items-start gap-3 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
            <Lock className="mt-0.5 h-5 w-5 flex-shrink-0 text-red-600 dark:text-red-300" />
            <div>
              <h1 className="text-lg font-semibold text-red-700 dark:text-red-100">
                {t('kds.accessDeniedTitle', 'Access denied')}
              </h1>
              <p className="mt-1 text-sm text-red-700 dark:text-red-200">
                {t('kds.accessDeniedText', 'You do not have permission to view the kitchen display.')}
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div
      dir={isRTL ? 'rtl' : 'ltr'}
      className="min-h-screen space-y-4 bg-gray-50 p-4 text-gray-900 dark:bg-slate-950 dark:text-slate-100 sm:p-6"
    >
      <KDSHeader
        locale={locale}
        t={t}
        stats={stats}
        filter={filter}
        soundEnabled={soundEnabled}
        isFullscreen={isFullscreen}
        isFullscreenSupported={isFullscreenSupported}
        isConnected={isConnected}
        isOffline={isOffline}
        queuedActionsCount={queuedActionsCount}
        isReplayingOfflineQueue={isReplayingOfflineQueue}
        canManageSettings={canManageKDSSettings && !canManageSettingsLoading}
        settingsHref="/dashboard/kds/settings"
        isRefreshing={isRefreshing}
        onFilterChange={setFilter}
        onToggleSound={toggleSound}
        onToggleFullscreen={() => {
          void handleFullscreenToggle()
        }}
      />

      {showShellLoading && (
        <div className={`rounded-xl border border-blue-300 bg-blue-50 p-3 dark:border-blue-400/40 dark:bg-blue-500/10 ${isRTL ? 'text-right' : ''}`}>
          <div className={`flex items-center gap-2 text-sm text-blue-700 dark:text-blue-100 ${isRTL ? 'flex-row-reverse' : ''}`}>
            <Loader2 className="h-4 w-4 animate-spin" />
            <span>{t('kds.loadingData', 'Loading kitchen queue...')}</span>
          </div>
        </div>
      )}

      {error && (
        <div className="rounded-xl border border-red-300 bg-red-50 p-4 dark:border-red-400/40 dark:bg-red-500/10">
          <div className={`flex items-center justify-between gap-4 ${isRTL ? 'flex-row-reverse' : ''}`}>
            <div className={`flex items-center gap-2 text-red-700 dark:text-red-100 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
              <AlertCircle className="h-4 w-4" />
              <span className="text-sm">{error}</span>
            </div>
            <button
              type="button"
              onClick={() => {
                void refreshAll()
              }}
              className="touch-manipulation min-h-10 rounded-md border border-red-300 px-3 py-1.5 text-xs font-medium text-red-700 transition hover:bg-red-100 dark:border-red-300/50 dark:text-red-100 dark:hover:bg-red-500/20"
            >
              {t('kds.retry', 'Retry')}
            </button>
          </div>
        </div>
      )}

      {canViewAnalytics && analyticsAvailable && (
        <KDSAnalyticsPanel
          locale={locale}
          t={t}
          analytics={analytics}
          loading={analyticsLoading}
          error={analyticsError}
          onRefresh={() => {
            void refreshAnalytics()
          }}
        />
      )}

      <KDSErrorBoundary
        locale={locale}
        t={t}
        lastKnownOrderCount={filteredOrders.length}
        lastKnownStats={stats}
        onRetry={() => {
          void refreshAll()
        }}
      >
        <KDSGrid
          orders={filteredOrders}
          locale={locale}
          t={t}
          loading={showShellLoading}
          canBump={canBumpOrders && !canBumpLoading}
          canRecall={canRecallOrders && !canRecallLoading}
          pendingOrderIds={pendingOrderIds}
          displayDensity={settings.display_density}
          warningMinutes={settings.color_warning_minutes}
          dangerMinutes={settings.color_danger_minutes}
          onBump={bumpOrder}
          onRecall={recallOrder}
        />
      </KDSErrorBoundary>
    </div>
  )
}
