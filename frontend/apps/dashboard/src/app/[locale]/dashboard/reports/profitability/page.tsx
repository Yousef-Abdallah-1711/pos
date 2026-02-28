'use client'

import { useCallback, useMemo, useState, type ReactNode } from 'react'
import { usePathname } from 'next/navigation'
import { BarChart, Bar, CartesianGrid, Legend, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import { DollarSign, Loader, Lock, RefreshCw, TrendingUp } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { ExportButton, KpiCard, LastUpdatedIndicator, ReportChart, ReportTable } from '@/components/reports'
import { useReports } from '@/hooks/useReports'
import { getLocaleFromPath } from '@/lib/translations'
import { getProfitabilityReportLabels } from '@/lib/reportsTranslations'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { useHasPermission } from '@/hooks/useRbac'
import { exportReportFile, getProfitabilityReport } from '@/lib/api/reportsApi'
import { debugError } from '@/lib/debug'
import type { ProfitabilityReportData, SalesGranularity } from '@/types/reports'

type FiltersState = {
  date_from: string
  date_to: string
  granularity: SalesGranularity
}

function formatDateInput(date: Date): string {
  return date.toISOString().split('T')[0]
}

function defaultFilters(): FiltersState {
  const today = new Date()
  const from = new Date(today)
  from.setDate(today.getDate() - 29)
  return {
    date_from: formatDateInput(from),
    date_to: formatDateInput(today),
    granularity: 'day',
  }
}

function downloadBlob(filename: string, blob: Blob) {
  const url = URL.createObjectURL(blob)
  const anchor = document.createElement('a')
  anchor.href = url
  anchor.download = filename
  anchor.click()
  URL.revokeObjectURL(url)
}

export default function ProfitabilityReportsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'

  const labels = useMemo(() => getProfitabilityReportLabels(locale), [locale])

  const { hasPermission: canViewReports, loading: permissionLoading } = useHasPermission(
    MODULE_IDS.REPORTS,
    PERMISSION_LEVELS.READ,
  )

  const [filters, setFilters] = useState<FiltersState>(() => defaultFilters())

  const currencyFormatter = useMemo(
    () => new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', { style: 'currency', currency: 'EGP', maximumFractionDigits: 2 }),
    [isRTL],
  )
  const numberFormatter = useMemo(
    () => new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', { maximumFractionDigits: 0 }),
    [isRTL],
  )

  const fetcher = useCallback(() => getProfitabilityReport(filters), [filters.date_from, filters.date_to])
  const getErrorMessage = useCallback(
    (err: any) => err?.response?.data?.error || err?.message || labels.loadError,
    [labels.loadError]
  )
  const onFetchError = useCallback((message: string, err: unknown, silent: boolean) => {
    if (!silent) toast.error(message)
    debugError('[REPORTS_PROFITABILITY] fetch failed', err)
  }, [])

  const { data, loading, refreshing, error, lastUpdatedAt, refetch: fetchReport } = useReports<ProfitabilityReportData>({
    enabled: !permissionLoading && canViewReports,
    deps: [filters.date_from, filters.date_to],
    fetcher,
    pollIntervalMs: 30000,
    getErrorMessage,
    onError: onFetchError,
  })

  const exportBackend = async (format: 'csv' | 'pdf') => {
    if (!data) return
    try {
      const file = await exportReportFile({
        report_type: 'profitability',
        format,
        date_from: filters.date_from,
        date_to: filters.date_to,
        locale: locale === 'ar' ? 'ar' : 'en',
      })
      downloadBlob(file.filename, file.blob)
      toast.success(labels.exportSuccess)
    } catch (err) {
      debugError(`[REPORTS_PROFITABILITY] ${format.toUpperCase()} export failed`, err)
      toast.error(labels.exportError)
    }
  }

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

  const topItems = data?.top_items ?? []
  const lowItems = data?.lowest_margin_items ?? []
  const chartData = topItems.slice(0, 10).map((item) => ({
    name: item.name.length > 18 ? `${item.name.slice(0, 15)}...` : item.name,
    gross_profit: Number(item.gross_profit.toFixed(2)),
    gross_margin_pct: Number(item.gross_margin_pct.toFixed(2)),
  }))

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{labels.title}</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">{labels.subtitle}</p>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <Button type="button" variant="outline" onClick={() => void fetchReport(true)} disabled={refreshing || loading}>
            <RefreshCw className={`h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
            {labels.refresh}
          </Button>
          <ExportButton format="csv" label={labels.csv} onClick={() => void exportBackend('csv')} disabled={!data} />
          <ExportButton format="pdf" label={labels.pdf} onClick={() => void exportBackend('pdf')} disabled={!data} />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
        <FilterCard label={labels.dateFrom}>
          <input
            type="date"
            value={filters.date_from}
            onChange={(e) => setFilters((prev) => ({ ...prev, date_from: e.target.value }))}
            className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 px-3 py-2 text-sm"
          />
        </FilterCard>
        <FilterCard label={labels.dateTo}>
          <input
            type="date"
            value={filters.date_to}
            onChange={(e) => setFilters((prev) => ({ ...prev, date_to: e.target.value }))}
            className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 px-3 py-2 text-sm"
          />
        </FilterCard>
        <KpiCard title={labels.grossMargin} value={data ? `${data.summary.gross_margin_pct.toFixed(1)}%` : '-'} icon={<TrendingUp className="h-5 w-5" />} accent="blue" />
        <LastUpdatedIndicator
          locale={locale}
          generatedAtLabel={labels.generatedAt}
          lastUpdatedLabel={labels.lastUpdated}
          generatedAt={data?.generated_at}
          lastUpdatedAt={lastUpdatedAt}
        />
      </div>

      {loading ? (
        <LoadingCard />
      ) : error ? (
        <ErrorCard message={error} />
      ) : !data ? (
        <EmptyCard text={labels.noData} />
      ) : (
        <>
          {data.summary.units_missing_cost > 0 && (
            <div className="rounded-xl border border-yellow-200 dark:border-yellow-800 bg-yellow-50 dark:bg-yellow-900/20 p-4">
              <p className="font-semibold text-yellow-900 dark:text-yellow-300">{labels.missingCost}</p>
              <p className="text-sm text-yellow-800 dark:text-yellow-400 mt-1">{labels.missingCostText}</p>
              <p className="text-xs text-yellow-700 dark:text-yellow-500 mt-2">
                {labels.missingCostUnits}: {numberFormatter.format(data.summary.units_missing_cost)}
              </p>
            </div>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
            <KpiCard title={labels.totalRevenue} value={currencyFormatter.format(data.summary.total_revenue)} icon={<DollarSign className="h-5 w-5" />} accent="green" />
            <KpiCard title={labels.totalCogs} value={currencyFormatter.format(data.summary.total_cogs)} icon={<DollarSign className="h-5 w-5" />} accent="amber" />
            <KpiCard title={labels.grossProfit} value={currencyFormatter.format(data.summary.gross_profit)} icon={<TrendingUp className="h-5 w-5" />} accent="blue" />
            <KpiCard title={labels.foodCost} value={`${data.summary.food_cost_pct.toFixed(1)}%`} icon={<TrendingUp className="h-5 w-5" />} accent="red" />
          </div>

          <ReportChart title={labels.topItems} bodyClassName="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#94a3b8" opacity={0.2} />
                  <XAxis dataKey="name" tick={{ fontSize: 12 }} interval={0} angle={-20} textAnchor="end" height={60} />
                  <YAxis yAxisId="left" tick={{ fontSize: 12 }} />
                  <YAxis yAxisId="right" orientation="right" tick={{ fontSize: 12 }} />
                  <Tooltip />
                  <Legend />
                  <Bar yAxisId="left" dataKey="gross_profit" fill="#2563eb" name={labels.profit} radius={[6, 6, 0, 0]} />
                  <Bar yAxisId="right" dataKey="gross_margin_pct" fill="#14b8a6" name={labels.margin} radius={[6, 6, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
          </ReportChart>

          <div className="grid grid-cols-1 2xl:grid-cols-2 gap-6">
            <ReportTable
              title={labels.topItems}
              columns={[labels.product, labels.qty, labels.revenue, labels.cogs, labels.profit, labels.margin]}
              emptyText={labels.noData}
              rows={topItems.map((item) => [
                item.name,
                numberFormatter.format(item.quantity),
                currencyFormatter.format(item.revenue),
                currencyFormatter.format(item.cogs),
                currencyFormatter.format(item.gross_profit),
                `${item.gross_margin_pct.toFixed(1)}%`,
              ])}
            />
            <ReportTable
              title={labels.lowMargin}
              columns={[labels.product, labels.qty, labels.revenue, labels.cogs, labels.profit, labels.margin]}
              emptyText={labels.noData}
              rows={lowItems.map((item) => [
                item.name,
                numberFormatter.format(item.quantity),
                currencyFormatter.format(item.revenue),
                currencyFormatter.format(item.cogs),
                currencyFormatter.format(item.gross_profit),
                `${item.gross_margin_pct.toFixed(1)}%`,
              ])}
            />
          </div>
        </>
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

function LoadingCard() {
  return (
    <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-10 flex items-center justify-center">
      <Loader className="h-7 w-7 animate-spin text-primary" />
    </div>
  )
}

function ErrorCard({ message }: { message: string }) {
  return <div className="rounded-xl border border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20 p-6 text-sm text-red-800 dark:text-red-300">{message}</div>
}

function EmptyCard({ text }: { text: string }) {
  return <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-6 text-sm text-gray-600 dark:text-gray-400">{text}</div>
}
