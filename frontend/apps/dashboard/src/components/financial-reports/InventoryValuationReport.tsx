'use client'

import { memo, useMemo } from 'react'
import { useTranslations } from 'next-intl'
import type { InventoryValuationReport as InventoryValuationReportType } from '@/types/financial-reports'

interface InventoryValuationReportProps {
  report: InventoryValuationReportType
  locale: string
}

export function InventoryValuationReport({ report, locale }: InventoryValuationReportProps) {
  const t = useTranslations('reports.financial')
  const formatter = useMemo(
    () =>
      new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
        style: 'currency',
        currency: 'EGP',
        maximumFractionDigits: 2,
      }),
    [locale]
  )

  return (
    <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
      <div className="border-b border-gray-200 px-5 py-4 text-lg font-semibold text-gray-900 dark:border-gray-700 dark:text-white">
        {t('inventoryValuation')}
      </div>
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-900/40">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">{t('table.product')}</th>
              <th className="px-4 py-3 text-right font-medium text-gray-500 dark:text-gray-400">{t('table.qty')}</th>
              <th className="px-4 py-3 text-right font-medium text-gray-500 dark:text-gray-400">{t('table.costPerUnit')}</th>
              <th className="px-4 py-3 text-right font-medium text-gray-500 dark:text-gray-400">{t('table.totalValue')}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
            {report.rows.length === 0 ? (
              <tr>
                <td colSpan={4} className="px-4 py-6 text-center text-gray-500 dark:text-gray-400">
                  {t('common.noDataAvailable')}
                </td>
              </tr>
            ) : (
              report.rows.map((row) => (
                <tr key={row.product_id}>
                  <td className="px-4 py-3 text-gray-700 dark:text-gray-200">{row.product_name}</td>
                  <td className="px-4 py-3 text-right text-gray-700 dark:text-gray-200">{row.quantity.toFixed(2)}</td>
                  <td className="px-4 py-3 text-right text-gray-700 dark:text-gray-200">
                    {formatter.format(row.cost_per_unit)}
                  </td>
                  <td className="px-4 py-3 text-right font-semibold text-gray-900 dark:text-white">
                    {formatter.format(row.total_value)}
                  </td>
                </tr>
              ))
            )}
          </tbody>
          <tfoot className="bg-gray-50 dark:bg-gray-900/40">
            <tr>
              <td className="px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">{t('table.total')}</td>
              <td colSpan={2}></td>
              <td className="px-4 py-3 text-right font-bold text-gray-900 dark:text-white">
                {formatter.format(report.total_value)}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>
  )
}

export default memo(InventoryValuationReport)
