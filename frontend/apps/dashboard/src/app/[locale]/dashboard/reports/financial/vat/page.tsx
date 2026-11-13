'use client'

import { useMemo, useState } from 'react'
import { useTranslations } from 'next-intl'
import { usePathname } from 'next/navigation'
import { VATReport } from '@/components/financial-reports/VATReport'
import { ReportExportButton } from '@/components/financial-reports/ReportExportButton'
import { useFinancialReport } from '@/hooks/useFinancialReports'
import { getVATReport } from '@/lib/api/financialReportsApi'
import { getLocaleFromPath } from '@/lib/translations'

function defaultDates() {
  const to = new Date()
  const from = new Date()
  from.setDate(to.getDate() - 29)
  return { date_from: from.toISOString().slice(0, 10), date_to: to.toISOString().slice(0, 10) }
}

export default function VATReportPage() {
  const t = useTranslations('reports.financial')
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const [filters, setFilters] = useState(defaultDates())

  const { data, loading, error, refetch } = useFinancialReport({
    cacheKey: `vat-${filters.date_from}-${filters.date_to}`,
    deps: [filters.date_from, filters.date_to],
    fetcher: () => getVATReport(filters),
  })

  const exportRows = useMemo(() => {
    if (!data) return []
    return [
      { metric: t('outputVat'), amount: data.output_vat },
      { metric: t('inputVat'), amount: data.input_vat },
      { metric: t('netVat'), amount: data.net_vat },
    ]
  }, [data, t])

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 xl:flex-row xl:items-end xl:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{t('vat')}</h1>
          <p className="mt-1 text-gray-500 dark:text-gray-400">
            {t('vatSubtitle')}
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
          <button
            type="button"
            onClick={() => void refetch()}
            className="rounded-md border border-gray-300 px-3 py-2 text-sm dark:border-gray-700"
          >
            {t('common.refresh')}
          </button>
          {data ? (
            <ReportExportButton
              filename={`vat-${filters.date_from}-to-${filters.date_to}.pdf`}
              title={t('vatExportTitle')}
              rows={exportRows}
              template="vat"
            />
          ) : null}
        </div>
      </div>

      {loading ? <StateCard message={t('common.loadingReport')} /> : null}
      {error ? <StateCard message={error} error /> : null}
      {!loading && !error && data ? <VATReport report={data} locale={locale} /> : null}
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
