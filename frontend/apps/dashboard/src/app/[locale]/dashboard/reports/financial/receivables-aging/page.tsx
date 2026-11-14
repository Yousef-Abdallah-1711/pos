'use client'

import { useMemo, useState } from 'react'
import { useTranslations } from 'next-intl'
import { usePathname } from 'next/navigation'
import { AgingReport } from '@/components/financial-reports/AgingReport'
import { ReportExportButton } from '@/components/financial-reports/ReportExportButton'
import { useFinancialReport } from '@/hooks/useFinancialReports'
import { getReceivablesAgingReport } from '@/lib/api/financialReportsApi'
import { getLocaleFromPath } from '@/lib/translations'

export default function ReceivablesAgingPage() {
  const t = useTranslations('reports.financial')
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const [asOf, setAsOf] = useState(() => new Date().toISOString().slice(0, 10))

  const { data, loading, error, refetch } = useFinancialReport({
    cacheKey: `ar-aging-${asOf}`,
    deps: [asOf],
    fetcher: () => getReceivablesAgingReport({ as_of: asOf }),
  })

  const exportRows = useMemo(() => {
    if (!data) return []
    return data.rows.map((row) => ({
      entity: row.entity_name,
      current: row.current,
      days_30: row.days_30,
      days_60: row.days_60,
      days_90: row.days_90,
      days_90_plus: row.days_90_plus,
      total: row.total,
    }))
  }, [data])

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 xl:flex-row xl:items-end xl:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{t('receivablesAging')}</h1>
          <p className="mt-1 text-gray-500 dark:text-gray-400">
            {t('receivablesAgingSubtitle')}
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <input
            type="date"
            value={asOf}
            onChange={(e) => setAsOf(e.target.value)}
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
              filename={`receivables-aging-${asOf}.pdf`}
              title={t('receivablesAgingExportTitle')}
              rows={exportRows}
            />
          ) : null}
        </div>
      </div>

      {loading ? <StateCard message={t('common.loadingReport')} /> : null}
      {error ? <StateCard message={error} error /> : null}
      {!loading && !error && data ? <AgingReport title={t('receivablesAging')} report={data} locale={locale} /> : null}
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
