'use client'

import { useCallback, useMemo, useState, type ReactNode } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Activity, ArrowRight, Loader, Lock, RefreshCw, TrendingUp } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { LastUpdatedIndicator } from '@/components/reports'
import { useReports } from '@/hooks/useReports'
import { getLocaleFromPath } from '@/lib/translations'
import { getReportsOverviewLabels } from '@/lib/reportsTranslations'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { useHasPermission } from '@/hooks/useRbac'
import { getReportsOverview } from '@/lib/api/reportsApi'
import { debugError } from '@/lib/debug'
import type { ReportsOverviewData } from '@/types/reports'

type OverviewFilters = { date_from: string; date_to: string }

function formatDateInput(date: Date): string {
  return date.toISOString().split('T')[0]
}

function defaultFilters(): OverviewFilters {
  const today = new Date()
  const from = new Date(today)
  from.setDate(today.getDate() - 29)
  return { date_from: formatDateInput(from), date_to: formatDateInput(today) }
}

export default function ReportsOverviewPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'

  const labels = useMemo(() => getReportsOverviewLabels(locale), [locale])

  const { hasPermission: canViewReports, loading: permissionLoading } = useHasPermission(MODULE_IDS.REPORTS, PERMISSION_LEVELS.READ)
  const [filters, setFilters] = useState<OverviewFilters>(() => defaultFilters())

  const currencyFormatter = useMemo(
    () => new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', { style: 'currency', currency: 'EGP', maximumFractionDigits: 2 }),
    [isRTL],
  )
  const numberFormatter = useMemo(
    () => new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', { maximumFractionDigits: 0 }),
    [isRTL],
  )

  const fetcher = useCallback(
    () => getReportsOverview(filters),
    [filters.date_from, filters.date_to]
  )

  const {
    data,
    loading,
    refreshing,
    error,
    lastUpdatedAt,
    refetch: fetchOverview,
  } = useReports<ReportsOverviewData>({
    enabled: !permissionLoading && canViewReports,
    deps: [filters.date_from, filters.date_to],
    fetcher,
    pollIntervalMs: 30000,
    getErrorMessage: (err: any) => err?.response?.data?.error || err?.message || labels.loadError,
    onError: (message, err, silent) => {
      if (!silent) toast.error(message)
      debugError('[REPORTS_OVERVIEW] fetch failed', err)
    },
  })

  if (permissionLoading) {
    return <div className="flex items-center justify-center h-96"><Loader className="h-8 w-8 animate-spin text-primary" /></div>
  }
  if (!canViewReports) {
    return (
      <div className="space-y-4">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{labels.title}</h1>
        <div className="rounded-xl border border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20 p-6">
          <div className="flex items-start gap-3">
            <Lock className="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5" />
            <div>
              <div className="font-semibold text-red-900 dark:text-red-300">{labels.permissionDeniedTitle}</div>
              <p className="text-sm text-red-800 dark:text-red-400 mt-1">{labels.permissionDeniedText}</p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <Activity className="h-7 w-7 text-primary" />
            {labels.title}
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">{labels.subtitle}</p>
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">{labels.autoRefresh}</p>
        </div>
        <div className="flex items-center gap-2">
          <Button type="button" variant="outline" onClick={() => void fetchOverview(true)} disabled={refreshing || loading}>
            <RefreshCw className={`h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
            {labels.refresh}
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
        <FilterCard label={labels.dateFrom}>
          <input type="date" value={filters.date_from} onChange={(e) => setFilters((prev) => ({ ...prev, date_from: e.target.value }))} className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 px-3 py-2 text-sm" />
        </FilterCard>
        <FilterCard label={labels.dateTo}>
          <input type="date" value={filters.date_to} onChange={(e) => setFilters((prev) => ({ ...prev, date_to: e.target.value }))} className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 px-3 py-2 text-sm" />
        </FilterCard>
        <div className="xl:col-span-2">
          <LastUpdatedIndicator
            locale={locale}
            generatedAtLabel={labels.generatedAt}
            lastUpdatedLabel={labels.lastUpdated}
            generatedAt={data?.generated_at}
            lastUpdatedAt={lastUpdatedAt}
          />
        </div>
      </div>

      {loading ? (
        <LoadingCard />
      ) : error ? (
        <ErrorCard message={error} />
      ) : !data ? null : (
        <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
          <OverviewCard
            title={labels.sales}
            icon={<TrendingUp className="h-5 w-5" />}
            href={`/${locale}/dashboard/reports/sales`}
            ctaLabel={labels.viewDetails}
            accent="blue"
          >
            <MetricLine label={labels.revenue} value={currencyFormatter.format(data.sales.summary.total_revenue)} />
            <MetricLine label={labels.orders} value={numberFormatter.format(data.sales.summary.total_orders)} />
            <MetricLine label={labels.avgTicket} value={currencyFormatter.format(data.sales.summary.avg_ticket)} />
            {data.sales.top_item && (
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-3">
                {labels.topSalesItem}: <span className="font-medium text-gray-900 dark:text-white">{data.sales.top_item.name}</span>
              </p>
            )}
          </OverviewCard>

          <OverviewCard
            title={labels.profitability}
            icon={<TrendingUp className="h-5 w-5" />}
            href={`/${locale}/dashboard/reports/profitability`}
            ctaLabel={labels.viewDetails}
            accent="emerald"
          >
            <MetricLine label={labels.grossMargin} value={`${data.profitability.summary.gross_margin_pct.toFixed(1)}%`} />
            <MetricLine label={labels.foodCost} value={`${data.profitability.summary.food_cost_pct.toFixed(1)}%`} />
            <MetricLine label={labels.grossProfit} value={currencyFormatter.format(data.profitability.summary.gross_profit)} />
            {data.profitability.most_profitable && (
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-3">
                {labels.topItem}: <span className="font-medium text-gray-900 dark:text-white">{data.profitability.most_profitable.name}</span>
              </p>
            )}
          </OverviewCard>

          <OverviewCard
            title={labels.inventory}
            icon={<TrendingUp className="h-5 w-5" />}
            href={`/${locale}/dashboard/reports/inventory`}
            ctaLabel={labels.viewDetails}
            accent="amber"
          >
            <MetricLine label={labels.avgTurnover} value={`${data.inventory.summary.average_monthly_turnover.toFixed(2)}x`} />
            <MetricLine label={labels.slowMoving} value={numberFormatter.format(data.inventory.summary.slow_moving_count)} />
            <MetricLine label={labels.inventoryValue} value={currencyFormatter.format(data.inventory.summary.total_current_inventory_value)} />
            {data.inventory.slowest_moving && (
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-3">
                {labels.slowItem}: <span className="font-medium text-gray-900 dark:text-white">{data.inventory.slowest_moving.name}</span>
              </p>
            )}
          </OverviewCard>
        </div>
      )}
    </div>
  )
}

function FilterCard({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4">
      <label className="text-sm text-gray-600 dark:text-gray-400 block mb-2">{label}</label>
      {children}
    </div>
  )
}

function OverviewCard({
  title,
  icon,
  children,
  href,
  ctaLabel,
  accent,
}: {
  title: string
  icon: ReactNode
  children: ReactNode
  href: string
  ctaLabel: string
  accent: 'blue' | 'emerald' | 'amber'
}) {
  const accents: Record<string, string> = {
    blue: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
    emerald: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400',
    amber: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400',
  }

  return (
    <div className="rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-5 flex flex-col">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{title}</h2>
        <div className={`p-2 rounded-lg ${accents[accent]}`}>{icon}</div>
      </div>
      <div className="space-y-2 flex-1">{children}</div>
      <Link href={href} className="mt-4 inline-flex items-center gap-2 text-sm font-medium text-primary hover:underline">
        {ctaLabel}
        <ArrowRight className="h-4 w-4" />
      </Link>
    </div>
  )
}

function MetricLine({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-3">
      <span className="text-sm text-gray-600 dark:text-gray-400">{label}</span>
      <span className="text-sm font-semibold text-gray-900 dark:text-white text-right">{value}</span>
    </div>
  )
}

function LoadingCard() {
  return <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-10 flex items-center justify-center"><Loader className="h-7 w-7 animate-spin text-primary" /></div>
}

function ErrorCard({ message }: { message: string }) {
  return <div className="rounded-xl border border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20 p-6 text-sm text-red-800 dark:text-red-300">{message}</div>
}
