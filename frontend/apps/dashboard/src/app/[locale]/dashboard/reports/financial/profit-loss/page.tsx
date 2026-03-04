'use client'

import { useMemo, useState } from 'react'
import { useTranslations } from 'next-intl'
import { usePathname } from 'next/navigation'
import { ProfitLossStatement } from '@/components/financial-reports/ProfitLossStatement'
import { ReportExportButton } from '@/components/financial-reports/ReportExportButton'
import { useFinancialReport } from '@/hooks/useFinancialReports'
import { getProfitLossReport } from '@/lib/api/financialReportsApi'
import { getLocaleFromPath } from '@/lib/translations'

function defaultDates() {
  const to = new Date()
  const from = new Date()
  from.setDate(to.getDate() - 29)
  const compareTo = new Date(from)
  compareTo.setDate(compareTo.getDate() - 1)
  const compareFrom = new Date(compareTo)
  compareFrom.setDate(compareFrom.getDate() - 29)
  return {
    date_from: from.toISOString().slice(0, 10),
    date_to: to.toISOString().slice(0, 10),
    compare_date_from: compareFrom.toISOString().slice(0, 10),
    compare_date_to: compareTo.toISOString().slice(0, 10),
  }
}

export default function ProfitLossPage() {
  const t = useTranslations('reports.financial')
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const [filters, setFilters] = useState(defaultDates())

  const { data, loading, error, refetch } = useFinancialReport({
    cacheKey: `pl-${filters.date_from}-${filters.date_to}-${filters.compare_date_from}-${filters.compare_date_to}`,
    deps: [filters.date_from, filters.date_to, filters.compare_date_from, filters.compare_date_to],
    fetcher: () => getProfitLossReport(filters),
  })

  const exportRows = useMemo(() => {
    if (!data) return []
    const revenueRows = data.revenue.map((r) => ({
      section: t('sections.revenue'),
      code: r.account_code,
      account: r.account_name,
      amount: r.amount,
    }))
    const expenseRows = data.expenses.map((r) => ({
      section: t('sections.expense'),
      code: r.account_code,
      account: r.account_name,
      amount: r.amount,
    }))
    return [...revenueRows, ...expenseRows]
  }, [data, t])

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{t('profitLoss')}</h1>
          <p className="mt-1 text-gray-500 dark:text-gray-400">
            {t('profitLossSubtitle')}
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <input
            type="date"
            value={filters.date_from}
            onChange={(e) => setFilters((prev) => ({ ...prev, date_from: e.target.value }))}
            className="rounded-md border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-900"
          />
          <input
            type="date"
            value={filters.date_to}
            onChange={(e) => setFilters((prev) => ({ ...prev, date_to: e.target.value }))}
            className="rounded-md border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-900"
          />
          <input
            type="date"
            value={filters.compare_date_from}
            onChange={(e) => setFilters((prev) => ({ ...prev, compare_date_from: e.target.value }))}
            className="rounded-md border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-900"
            title={t('compareFrom')}
          />
          <input
            type="date"
            value={filters.compare_date_to}
            onChange={(e) => setFilters((prev) => ({ ...prev, compare_date_to: e.target.value }))}
            className="rounded-md border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-900"
            title={t('compareTo')}
          />
          <button
            type="button"
            onClick={() => void refetch()}
            className="rounded-md border border-gray-300 px-3 py-2 text-sm dark:border-gray-700"
          >
            {t('common.refresh')}
          </button>
          {data ? (
            <ReportExportButton
              filename={`profit-loss-${filters.date_from}-to-${filters.date_to}.pdf`}
              title={t('profitLossExportTitle')}
              rows={exportRows}
              template="profit-loss"
            />
          ) : null}
        </div>
      </div>

      {loading ? <StateCard message={t('common.loadingReport')} /> : null}
      {error ? <StateCard message={error} error /> : null}
      {!loading && !error && data ? <ProfitLossStatement report={data} locale={locale} /> : null}
    </div>
  )
}

function StateCard({ message, error = false }: { message: string; error?: boolean }) {
  return (
    <div
      className={`rounded-2xl border p-5 ${
        error
          ? 'border-red-200 bg-red-50 text-red-700 dark:border-red-900 dark:bg-red-950/20 dark:text-red-300'
          : 'border-gray-200 bg-white text-gray-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300'
      }`}
    >
      {message}
    </div>
  )
}
