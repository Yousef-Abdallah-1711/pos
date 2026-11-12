'use client'

import { useCallback, useMemo, useState } from 'react'
import { usePathname } from 'next/navigation'
import {
  BarChart as BarChartIcon,
  Download,
  Loader,
  Lock,
  RefreshCw,
  TrendingUp,
  ShoppingBag,
} from 'lucide-react'
import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
} from '@/components/reports/lazy-recharts'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { ExportButton } from '@/components/reports/ExportButton'
import { KpiCard } from '@/components/reports/KpiCard'
import { LastUpdatedIndicator } from '@/components/reports/LastUpdatedIndicator'
import { ReportChart } from '@/components/reports/ReportChart'
import { ReportTable } from '@/components/reports/ReportTable'
import { useReports } from '@/hooks/useReports'
import { getLocaleFromPath } from '@/lib/translations'
import { getSalesReportLabels } from '@/lib/reportsTranslations'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { exportReportFile, getSalesReport } from '@/lib/api/reportsApi'
import { useUserPermissions } from '@/hooks/useRbac'
import { debugError } from '@/lib/debug'
import type { SalesReportData, SalesGranularity } from '@/types/reports'

type SalesFiltersState = {
  date_from: string
  date_to: string
  granularity: SalesGranularity
}

const CATEGORY_COLORS = ['#2563eb', '#059669', '#f59e0b', '#dc2626', '#7c3aed', '#0891b2']

function formatDateInput(date: Date): string {
  return date.toISOString().split('T')[0]
}

function defaultFilters(): SalesFiltersState {
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

export default function SalesReportsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'
  const labels = useMemo(() => getSalesReportLabels(locale), [locale])

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canViewReports = checkPermission(MODULE_IDS.REPORTS, PERMISSION_LEVELS.READ)

  const [filters, setFilters] = useState<SalesFiltersState>(() => defaultFilters())

  const currencyFormatter = useMemo(
    () =>
      new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', {
        style: 'currency',
        currency: 'EGP',
        maximumFractionDigits: 2,
      }),
    [isRTL]
  )

  const numberFormatter = useMemo(
    () =>
      new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', {
        maximumFractionDigits: 0,
      }),
    [isRTL]
  )

  const fetcher = useCallback(
    () => getSalesReport(filters),
    [filters.date_from, filters.date_to, filters.granularity]
  )
  const getErrorMessage = useCallback(
    (err: any) => err?.response?.data?.error || err?.message || labels.loadError,
    [labels.loadError]
  )
  const onFetchError = useCallback((message: string, err: unknown, silent: boolean) => {
    if (!silent) toast.error(message)
    debugError('[REPORTS_SALES] fetch failed', err)
  }, [])

  const { data, loading, refreshing, error, lastUpdatedAt, refetch: fetchReport } = useReports<SalesReportData>({
    enabled: !permissionLoading && canViewReports,
    deps: [filters.date_from, filters.date_to, filters.granularity],
    fetcher,
    pollIntervalMs: 30000,
    getErrorMessage,
    onError: onFetchError,
  })

  const exportBackend = async (format: 'csv' | 'pdf') => {
    try {
      const file = await exportReportFile({
        report_type: 'sales',
        format,
        date_from: filters.date_from,
        date_to: filters.date_to,
        granularity: filters.granularity,
        locale: locale === 'ar' ? 'ar' : 'en',
      })
      downloadBlob(file.filename, file.blob)
      toast.success(labels.exportSuccess)
    } catch (err) {
      debugError(`[REPORTS_SALES] ${format.toUpperCase()} export failed`, err)
      toast.error(labels.exportError)
    }
  }

  const exportCsv = () => {
    if (!data) return
    void exportBackend('csv')
  }

  const exportPdf = async () => {
    if (!data) return
    await exportBackend('pdf')
  }

  if (permissionLoading) {
    return (
      <div className="flex items-center justify-center h-96">
        <Loader className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!canViewReports) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{labels.title}</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">{labels.subtitle}</p>
        </div>
        <div className="rounded-xl border border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20 p-6">
          <div className="flex items-start gap-3">
            <Lock className="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5" />
            <div>
              <h2 className="font-semibold text-red-900 dark:text-red-300">{labels.permissionDeniedTitle}</h2>
              <p className="text-sm text-red-800 dark:text-red-400 mt-1">{labels.permissionDeniedText}</p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  const trendData = data?.trend ?? []
  const hourlyData = data?.hourly ?? []
  const topItems = data?.top_items ?? []
  const topCategories = data?.top_categories ?? []

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <BarChartIcon className="h-7 w-7 text-primary" />
            {labels.title}
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">{labels.subtitle}</p>
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">{labels.autoRefresh}</p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Button type="button" variant="outline" onClick={() => void fetchReport(true)} disabled={refreshing || loading}>
            <RefreshCw className={`h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
            {labels.refresh}
          </Button>
          <ExportButton format="csv" label={labels.csv} onClick={exportCsv} disabled={!data} />
          <ExportButton format="pdf" label={labels.pdf} onClick={() => void exportPdf()} disabled={!data} />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
        <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4">
          <label className="text-sm text-gray-600 dark:text-gray-400 block mb-2">{labels.dateFrom}</label>
          <input
            type="date"
            value={filters.date_from}
            onChange={(e) => setFilters((prev) => ({ ...prev, date_from: e.target.value }))}
            className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 px-3 py-2 text-sm"
          />
        </div>
        <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4">
          <label className="text-sm text-gray-600 dark:text-gray-400 block mb-2">{labels.dateTo}</label>
          <input
            type="date"
            value={filters.date_to}
            onChange={(e) => setFilters((prev) => ({ ...prev, date_to: e.target.value }))}
            className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 px-3 py-2 text-sm"
          />
        </div>
        <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4">
          <label className="text-sm text-gray-600 dark:text-gray-400 block mb-2">{labels.granularity}</label>
          <select
            value={filters.granularity}
            onChange={(e) =>
              setFilters((prev) => ({ ...prev, granularity: e.target.value as SalesGranularity }))
            }
            className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 px-3 py-2 text-sm"
          >
            <option value="day">{labels.day}</option>
            <option value="week">{labels.week}</option>
            <option value="month">{labels.month}</option>
          </select>
        </div>
        <LastUpdatedIndicator
          locale={locale}
          generatedAtLabel={labels.generatedAt}
          lastUpdatedLabel={labels.lastUpdated}
          generatedAt={data?.generated_at}
          lastUpdatedAt={lastUpdatedAt}
        />
      </div>

      {loading ? (
        <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-10 flex items-center justify-center">
          <Loader className="h-7 w-7 animate-spin text-primary" />
        </div>
      ) : error ? (
        <div className="rounded-xl border border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20 p-6 text-sm text-red-800 dark:text-red-300">
          {error}
        </div>
      ) : !data ? (
        <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-6 text-sm text-gray-600 dark:text-gray-400">
          {labels.noData}
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <KpiCard
              title={labels.totalRevenue}
              value={currencyFormatter.format(data.summary.total_revenue)}
              icon={<TrendingUp className="h-5 w-5" />}
              accent="green"
            />
            <KpiCard
              title={labels.totalOrders}
              value={numberFormatter.format(data.summary.total_orders)}
              icon={<ShoppingBag className="h-5 w-5" />}
              accent="blue"
            />
            <KpiCard
              title={labels.avgTicket}
              value={currencyFormatter.format(data.summary.avg_ticket)}
              icon={<Download className="h-5 w-5" />}
              accent="amber"
            />
          </div>

          <div className="grid grid-cols-1 2xl:grid-cols-3 gap-6">
            <ReportChart title={labels.salesTrend} className="2xl:col-span-2" bodyClassName="h-80">
              <ResponsiveContainer width="100%" height="100%" minWidth={0} minHeight={280}>
                <LineChart data={trendData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#94a3b8" opacity={0.2} />
                  <XAxis dataKey="bucket" tick={{ fontSize: 12 }} />
                  <YAxis tickFormatter={(value) => compactCurrency(value)} tick={{ fontSize: 12 }} />
                  <Tooltip
                    formatter={(value, name) => [
                      name === 'revenue' ? currencyFormatter.format(Number(value)) : numberFormatter.format(Number(value)),
                      name === 'revenue' ? labels.revenue : labels.orders,
                    ]}
                  />
                  <Legend />
                  <Line type="monotone" dataKey="revenue" name={labels.revenue} stroke="#2563eb" strokeWidth={2.5} dot={false} />
                  <Line type="monotone" dataKey="orders" name={labels.orders} stroke="#059669" strokeWidth={2} dot={false} />
                </LineChart>
              </ResponsiveContainer>
            </ReportChart>

            <ReportChart title={labels.topCategories}>
              <div className="h-80">
                <ResponsiveContainer width="100%" height="100%" minWidth={0} minHeight={280}>
                  <PieChart>
                    <Pie data={topCategories} dataKey="revenue" nameKey="name" outerRadius={100} innerRadius={56}>
                      {topCategories.map((_, index) => (
                        <Cell key={`cell-${index}`} fill={CATEGORY_COLORS[index % CATEGORY_COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip formatter={(value) => currencyFormatter.format(Number(value))} />
                  </PieChart>
                </ResponsiveContainer>
              </div>
              <div className="mt-2 space-y-2">
                {topCategories.slice(0, 4).map((category, index) => (
                  <div key={`${category.name}-${index}`} className="flex items-center justify-between gap-3 text-sm">
                    <div className="flex items-center gap-2 min-w-0">
                      <span
                        className={`inline-block h-2.5 w-2.5 rounded-full flex-shrink-0 ${categoryDotClass(index)}`}
                      />
                      <span className="truncate text-gray-700 dark:text-gray-300">{category.name}</span>
                    </div>
                    <span className="font-medium text-gray-900 dark:text-white">
                      {currencyFormatter.format(category.revenue)}
                    </span>
                  </div>
                ))}
              </div>
            </ReportChart>
          </div>

          <ReportChart title={labels.hourlySales} bodyClassName="h-80">
            <ResponsiveContainer width="100%" height="100%" minWidth={0} minHeight={280}>
              <BarChart data={hourlyData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#94a3b8" opacity={0.2} />
                <XAxis dataKey="hour" tickFormatter={(h) => `${h}:00`} tick={{ fontSize: 12 }} />
                <YAxis tickFormatter={(value) => compactCurrency(value)} tick={{ fontSize: 12 }} />
                <Tooltip
                  labelFormatter={(label) => `${label}:00`}
                  formatter={(value, name) => [
                    name === 'revenue' ? currencyFormatter.format(Number(value)) : numberFormatter.format(Number(value)),
                    name === 'revenue' ? labels.revenue : labels.orders,
                  ]}
                />
                <Legend />
                <Bar dataKey="revenue" name={labels.revenue} fill="#2563eb" radius={[6, 6, 0, 0]} />
                <Bar dataKey="orders" name={labels.orders} fill="#14b8a6" radius={[6, 6, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </ReportChart>

          <div className="grid grid-cols-1 2xl:grid-cols-2 gap-6">
            <ReportTable
              title={labels.topItems}
              columns={[labels.product, labels.qty, labels.revenue]}
              emptyText={labels.noData}
              rows={topItems.map((item) => [
                item.name,
                numberFormatter.format(item.quantity),
                currencyFormatter.format(item.revenue),
              ])}
            />
            <ReportTable
              title={labels.topCategories}
              columns={[labels.category, labels.qty, labels.revenue]}
              emptyText={labels.noData}
              rows={topCategories.map((category) => [
                category.name,
                numberFormatter.format(category.quantity),
                currencyFormatter.format(category.revenue),
              ])}
            />
          </div>
        </>
      )}
    </div>
  )
}

function compactCurrency(value: number): string {
  if (Math.abs(value) >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}M`
  if (Math.abs(value) >= 1_000) return `${(value / 1_000).toFixed(1)}k`
  return `${Math.round(value)}`
}

function categoryDotClass(index: number): string {
  const classes = [
    'bg-blue-600',
    'bg-emerald-600',
    'bg-amber-500',
    'bg-red-600',
    'bg-violet-600',
    'bg-cyan-600',
  ]
  return classes[index % classes.length]
}
