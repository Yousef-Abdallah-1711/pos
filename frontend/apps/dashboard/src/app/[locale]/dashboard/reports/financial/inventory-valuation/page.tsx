'use client'

import { useMemo, useState } from 'react'
import { useTranslations } from 'next-intl'
import { usePathname } from 'next/navigation'
import { InventoryValuationReport } from '@/components/financial-reports/InventoryValuationReport'
import { ReportExportButton } from '@/components/financial-reports/ReportExportButton'
import { useFinancialReport } from '@/hooks/useFinancialReports'
import { getInventoryValuationReport } from '@/lib/api/financialReportsApi'
import { getLocaleFromPath } from '@/lib/translations'

export default function InventoryValuationPage() {
  const t = useTranslations('reports.financial')
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const [asOf, setAsOf] = useState(() => new Date().toISOString().slice(0, 10))

  const { data, loading, error, refetch } = useFinancialReport({
    cacheKey: `inventory-valuation-${asOf}`,
    deps: [asOf],
    fetcher: () => getInventoryValuationReport({ as_of: asOf }),
  })

  const exportRows = useMemo(() => {
    if (!data) return []
    return data.rows.map((row) => ({
      product: row.product_name,
      quantity: row.quantity,
      cost_per_unit: row.cost_per_unit,
      total_value: row.total_value,
    }))
  }, [data])

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{t('inventoryValuation')}</h1>
          <p className="mt-1 text-gray-500 dark:text-gray-400">
            {t('inventoryValuationSubtitle')}
          </p>
        </div>
        <div className="flex gap-2">
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
              filename={`inventory-valuation-${asOf}.pdf`}
              title={t('inventoryValuationExportTitle')}
              rows={exportRows}
            />
          ) : null}
        </div>
      </div>

      {loading ? <StateCard message={t('common.loadingReport')} /> : null}
      {error ? <StateCard message={error} error /> : null}
      {!loading && !error && data ? <InventoryValuationReport report={data} locale={locale} /> : null}
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
