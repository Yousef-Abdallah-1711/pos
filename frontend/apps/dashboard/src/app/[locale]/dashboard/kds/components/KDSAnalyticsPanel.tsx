'use client'

import { Activity, Clock3, Flame } from 'lucide-react'
import { KDSAnalytics } from '@/hooks/useKDSAnalytics'

interface KDSAnalyticsPanelProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  analytics: KDSAnalytics
  loading: boolean
  error: string | null
  onRefresh: () => void
}

const formatTrendLabel = (dateValue: string) => {
  const parsed = new Date(dateValue)
  if (Number.isNaN(parsed.getTime())) {
    return dateValue
  }

  return parsed.toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
  })
}

export function KDSAnalyticsPanel({
  locale,
  t,
  analytics,
  loading,
  error,
  onRefresh,
}: KDSAnalyticsPanelProps) {
  const isRTL = locale === 'ar'
  const trend = analytics.prep_time_trend || []
  const busiest = analytics.busiest_hours || []

  const maxTrendMinutes = trend.reduce((max, point) => Math.max(max, point.avg_prep_minutes || 0), 0)
  const maxHourOrders = busiest.reduce((max, point) => Math.max(max, point.order_count || 0), 0)

  return (
    <section
      dir={isRTL ? 'rtl' : 'ltr'}
      className="space-y-4 rounded-xl border border-gray-200 bg-white p-4 dark:border-slate-700 dark:bg-slate-900/80"
    >
      <div className={`flex flex-wrap items-center justify-between gap-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
        <div className={isRTL ? 'text-right' : ''}>
          <h2 className="text-lg font-semibold text-gray-900 dark:text-slate-100">
            {t('kds.analytics.title', 'Kitchen Analytics')}
          </h2>
          <p className="text-xs text-gray-600 dark:text-slate-400">
            {t('kds.analytics.subtitle', 'Preparation trend and workload pattern over the last 7 days')}
          </p>
        </div>

        <button
          type="button"
          onClick={onRefresh}
          disabled={loading}
          className="touch-manipulation rounded-md border border-gray-300 px-3 py-1.5 text-xs font-medium text-gray-800 transition hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-600 dark:text-slate-100 dark:hover:bg-slate-800"
        >
          {loading
            ? t('kds.analytics.refreshing', 'Refreshing...')
            : t('kds.analytics.refresh', 'Refresh analytics')}
        </button>
      </div>

      {error && (
        <div className={`rounded-lg border border-red-300 bg-red-50 p-3 text-xs text-red-700 dark:border-red-400/40 dark:bg-red-500/10 dark:text-red-200 ${isRTL ? 'text-right' : ''}`}>
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
        <div className={`rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-slate-700 dark:bg-slate-800/70 ${isRTL ? 'text-right' : ''}`}>
          <p className={`mb-1 flex items-center gap-1 text-xs uppercase tracking-wide text-gray-500 dark:text-slate-400 ${isRTL ? 'flex-row-reverse justify-end' : ''}`}>
            <Clock3 className="h-3.5 w-3.5" />
            {t('kds.analytics.avgPrep', 'Average prep time')}
          </p>
          <p className="text-xl font-bold text-emerald-700 dark:text-emerald-300">
            {(analytics.average_prep_minutes || 0).toFixed(1)}
            {t('kds.minutesSuffix', 'm')}
          </p>
        </div>

        <div className={`rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-slate-700 dark:bg-slate-800/70 ${isRTL ? 'text-right' : ''}`}>
          <p className={`mb-1 flex items-center gap-1 text-xs uppercase tracking-wide text-gray-500 dark:text-slate-400 ${isRTL ? 'flex-row-reverse justify-end' : ''}`}>
            <Activity className="h-3.5 w-3.5" />
            {t('kds.analytics.completedOrders', 'Completed orders')}
          </p>
          <p className="text-xl font-bold text-sky-700 dark:text-sky-300">
            {analytics.total_completed_orders || 0}
          </p>
        </div>

        <div className={`rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-slate-700 dark:bg-slate-800/70 ${isRTL ? 'text-right' : ''}`}>
          <p className={`mb-1 flex items-center gap-1 text-xs uppercase tracking-wide text-gray-500 dark:text-slate-400 ${isRTL ? 'flex-row-reverse justify-end' : ''}`}>
            <Flame className="h-3.5 w-3.5" />
            {t('kds.analytics.peakHour', 'Peak hour')}
          </p>
          <p className="text-xl font-bold text-amber-700 dark:text-amber-300">
            {(() => {
              if (!busiest.length) {
                return '--'
              }
              const peak = [...busiest].sort((a, b) => b.order_count - a.order_count)[0]
              return `${String(peak.hour).padStart(2, '0')}:00`
            })()}
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <div className="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-slate-700 dark:bg-slate-800/60">
          <p className={`mb-3 text-xs font-semibold uppercase tracking-wide text-gray-700 dark:text-slate-300 ${isRTL ? 'text-right' : ''}`}>
            {t('kds.analytics.prepTrend', '7-day prep trend')}
          </p>
          <div className="space-y-2">
            {trend.map((point) => {
              const width = maxTrendMinutes > 0
                ? Math.max(6, Math.round((point.avg_prep_minutes / maxTrendMinutes) * 100))
                : 6

              return (
                <div key={point.date} className={`space-y-1 ${isRTL ? 'text-right' : ''}`}>
                  <div className={`flex items-center justify-between text-xs text-gray-700 dark:text-slate-300 ${isRTL ? 'flex-row-reverse' : ''}`}>
                    <span>{formatTrendLabel(point.date)}</span>
                    <span>
                      {(point.avg_prep_minutes || 0).toFixed(1)}
                      {t('kds.minutesSuffix', 'm')}
                    </span>
                  </div>
                  <div className={`h-2 rounded bg-gray-200 dark:bg-slate-700 ${isRTL ? 'flex justify-end' : ''}`}>
                    <div
                      className="h-2 rounded bg-emerald-500/80"
                      style={{ width: `${width}%` }}
                    />
                  </div>
                </div>
              )
            })}
            {trend.length === 0 && (
              <p className={`text-xs text-gray-500 dark:text-slate-400 ${isRTL ? 'text-right' : ''}`}>
                {t('kds.analytics.noTrend', 'No prep trend data available yet')}
              </p>
            )}
          </div>
        </div>

        <div className="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-slate-700 dark:bg-slate-800/60">
          <p className={`mb-3 text-xs font-semibold uppercase tracking-wide text-gray-700 dark:text-slate-300 ${isRTL ? 'text-right' : ''}`}>
            {t('kds.analytics.busiestHours', 'Busiest hours heatmap')}
          </p>
          <div className={`grid grid-cols-6 gap-2 sm:grid-cols-8 md:grid-cols-12 ${isRTL ? 'text-right' : ''}`}>
            {busiest.map((point) => {
              const intensity = maxHourOrders > 0 ? point.order_count / maxHourOrders : 0
              const alpha = Math.min(0.9, Math.max(0.15, intensity))
              return (
                <div
                  key={point.hour}
                  className="rounded border border-gray-300 p-1 text-center dark:border-slate-700"
                  style={{ backgroundColor: `rgba(56, 189, 248, ${alpha})` }}
                  title={`${String(point.hour).padStart(2, '0')}:00 - ${point.order_count}`}
                >
                  <p className="text-[10px] font-semibold text-gray-900 dark:text-slate-100">{String(point.hour).padStart(2, '0')}</p>
                  <p className="text-[10px] text-gray-800 dark:text-slate-100/90">{point.order_count}</p>
                </div>
              )
            })}
            {busiest.length === 0 && (
              <p className={`col-span-full text-xs text-gray-500 dark:text-slate-400 ${isRTL ? 'text-right' : ''}`}>
                {t('kds.analytics.noHours', 'No hourly workload data available yet')}
              </p>
            )}
          </div>
        </div>
      </div>
    </section>
  )
}
