'use client'

import { useCallback, useMemo, useState, type ReactNode } from 'react'
import { usePathname } from 'next/navigation'
import { Bar, BarChart, CartesianGrid, Legend, ResponsiveContainer, Tooltip, XAxis, YAxis } from '@/components/reports/lazy-recharts'
import { AlertTriangle, Boxes, Loader, Lock, RefreshCw, TrendingUp } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { ExportButton } from '@/components/reports/ExportButton'
import { KpiCard } from '@/components/reports/KpiCard'
import { LastUpdatedIndicator } from '@/components/reports/LastUpdatedIndicator'
import { ReportChart } from '@/components/reports/ReportChart'
import { ReportTable } from '@/components/reports/ReportTable'
import { useReports } from '@/hooks/useReports'
import { getLocaleFromPath } from '@/lib/translations'
import { getInventoryReportLabels } from '@/lib/reportsTranslations'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { useUserPermissions } from '@/hooks/useRbac'
import { exportReportFile, getInventoryReport } from '@/lib/api/reportsApi'
import { debugError } from '@/lib/debug'
import type { InventoryReportData, SalesGranularity } from '@/types/reports'

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
  return { date_from: formatDateInput(from), date_to: formatDateInput(today), granularity: 'day' }
}

function downloadBlob(filename: string, blob: Blob) {
  const url = URL.createObjectURL(blob)
  const anchor = document.createElement('a')
  anchor.href = url
  anchor.download = filename
  anchor.click()
  URL.revokeObjectURL(url)
}

export default function InventoryReportsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'

  const labels = useMemo(() => getInventoryReportLabels(locale), [locale])

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canViewReports = checkPermission(MODULE_IDS.REPORTS, PERMISSION_LEVELS.READ)
  const [filters, setFilters] = useState<FiltersState>(() => defaultFilters())

  const currencyFormatter = useMemo(
    () => new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', { style: 'currency', currency: 'EGP', maximumFractionDigits: 2 }),
    [isRTL],
  )
  const numberFormatter = useMemo(
    () => new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', { maximumFractionDigits: 0 }),
    [isRTL],
  )

  const fetcher = useCallback(() => getInventoryReport(filters), [filters.date_from, filters.date_to])
  const getErrorMessage = useCallback(
    (err: any) => err?.response?.data?.error || err?.message || labels.loadError,
    [labels.loadError]
  )
  const onFetchError = useCallback((message: string, err: unknown, silent: boolean) => {
    if (!silent) toast.error(message)
    debugError('[REPORTS_INVENTORY] fetch failed', err)
  }, [])

  const { data, loading, refreshing, error, lastUpdatedAt, refetch: fetchReport } = useReports<InventoryReportData>({
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
        report_type: 'inventory',
        format,
        date_from: filters.date_from,
        date_to: filters.date_to,
        locale: locale === 'ar' ? 'ar' : 'en',
      })
      downloadBlob(file.filename, file.blob)
      toast.success(labels.exportSuccess)
    } catch (err) {
      debugError(`[REPORTS_INVENTORY] ${format.toUpperCase()} export failed`, err)
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

  const topTurnover = data?.top_turnover ?? []
  const slowMovers = data?.slow_movers ?? []
  const chartData = topTurnover.slice(0, 10).map((item) => ({
    name: trim(item.name, 18),
    turnover: Number(item.monthly_turnover_ratio.toFixed(2)),
    sold_qty: item.sold_qty,
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
          <input type="date" value={filters.date_from} onChange={(e) => setFilters((prev) => ({ ...prev, date_from: e.target.value }))} className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 px-3 py-2 text-sm" />
        </FilterCard>
        <FilterCard label={labels.dateTo}>
          <input type="date" value={filters.date_to} onChange={(e) => setFilters((prev) => ({ ...prev, date_to: e.target.value }))} className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 px-3 py-2 text-sm" />
        </FilterCard>
        <KpiCard title={labels.avgTurnover} value={data ? `${data.summary.average_monthly_turnover.toFixed(2)}x` : '-'} icon={<TrendingUp className="h-5 w-5" />} accent="blue" />
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
          <div className="rounded-xl border border-blue-200 dark:border-blue-800 bg-blue-50 dark:bg-blue-900/20 p-4">
            <p className="text-sm text-blue-800 dark:text-blue-300">{labels.formulaNote}</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
            <KpiCard title={labels.inventoryValue} value={currencyFormatter.format(data.summary.total_current_inventory_value)} icon={<Boxes className="h-5 w-5" />} accent="green" />
            <KpiCard title={labels.estCogs} value={currencyFormatter.format(data.summary.total_estimated_cogs)} icon={<TrendingUp className="h-5 w-5" />} accent="amber" />
            <KpiCard title={labels.slowMoving} value={numberFormatter.format(data.summary.slow_moving_count)} icon={<AlertTriangle className="h-5 w-5" />} accent="red" />
            <KpiCard title={labels.trackedProducts} value={numberFormatter.format(data.summary.tracked_products_count)} icon={<Boxes className="h-5 w-5" />} accent="blue" />
          </div>

          <ReportChart title={labels.topTurnover} bodyClassName="h-80">
              <ResponsiveContainer width="100%" height="100%" minWidth={0} minHeight={280}>
                <BarChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#94a3b8" opacity={0.2} />
                  <XAxis dataKey="name" tick={{ fontSize: 12 }} interval={0} angle={-20} textAnchor="end" height={60} />
                  <YAxis yAxisId="left" tick={{ fontSize: 12 }} />
                  <YAxis yAxisId="right" orientation="right" tick={{ fontSize: 12 }} />
                  <Tooltip />
                  <Legend />
                  <Bar yAxisId="left" dataKey="turnover" fill="#2563eb" name={labels.turnover} radius={[6, 6, 0, 0]} />
                  <Bar yAxisId="right" dataKey="sold_qty" fill="#14b8a6" name={labels.soldQty} radius={[6, 6, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
          </ReportChart>

          <div className="grid grid-cols-1 2xl:grid-cols-2 gap-6">
            <ReportTable
              title={labels.slowItems}
              columns={[labels.product, labels.soldQty, labels.currentStock, labels.turnover, labels.waste]}
              emptyText={labels.noData}
              rows={slowMovers.map((item) => [
                item.name,
                numberFormatter.format(item.sold_qty),
                numberFormatter.format(item.current_stock_qty),
                `${item.monthly_turnover_ratio.toFixed(2)}x`,
                numberFormatter.format(item.waste_qty),
              ])}
            />
            <ReportTable
              title={labels.topTurnover}
              columns={[labels.product, labels.soldQty, labels.currentStock, labels.turnover, labels.waste]}
              emptyText={labels.noData}
              rows={topTurnover.map((item) => [
                item.name,
                numberFormatter.format(item.sold_qty),
                numberFormatter.format(item.current_stock_qty),
                `${item.monthly_turnover_ratio.toFixed(2)}x`,
                numberFormatter.format(item.waste_qty),
              ])}
            />
          </div>
        </>
      )}
    </div>
  )
}

function FilterCard({ label, children }: { label: string; children: ReactNode }) {
  return <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4"><label className="text-sm text-gray-600 dark:text-gray-400 block mb-2">{label}</label>{children}</div>
}

function LoadingCard() {
  return <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-10 flex items-center justify-center"><Loader className="h-7 w-7 animate-spin text-primary" /></div>
}
function ErrorCard({ message }: { message: string }) {
  return <div className="rounded-xl border border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20 p-6 text-sm text-red-800 dark:text-red-300">{message}</div>
}
function EmptyCard({ text }: { text: string }) {
  return <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-6 text-sm text-gray-600 dark:text-gray-400">{text}</div>
}

function trim(value: string, max: number): string {
  return value.length > max ? `${value.slice(0, max - 3)}...` : value
}
