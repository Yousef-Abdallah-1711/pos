'use client'

import { memo, useMemo } from 'react'
import { useTranslations } from 'next-intl'
import type { BalanceSheetReport } from '@/types/financial-reports'

interface BalanceSheetProps {
  report: BalanceSheetReport
  locale: string
}

export function BalanceSheet({ report, locale }: BalanceSheetProps) {
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
    <div className="space-y-6">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
        <Metric label={t('metrics.totalAssets')} value={formatter.format(report.total_assets)} />
        <Metric label={t('metrics.totalLiabilities')} value={formatter.format(report.total_liabilities)} />
        <Metric label={t('metrics.totalEquity')} value={formatter.format(report.total_equity)} />
        <Metric
          label={t('metrics.equation')}
          value={report.equation_balanced ? t('metrics.balanced') : t('metrics.unbalanced')}
          tone={report.equation_balanced ? 'ok' : 'warn'}
        />
      </div>

      <div className="grid grid-cols-1 gap-6 xl:grid-cols-3">
        <Block title={t('sections.assets')} rows={report.assets} formatter={formatter} />
        <Block title={t('sections.liabilities')} rows={report.liabilities} formatter={formatter} />
        <Block title={t('sections.equity')} rows={report.equity} formatter={formatter} />
      </div>
    </div>
  )
}

export default memo(BalanceSheet)

function Metric({
  label,
  value,
  tone = 'default',
}: {
  label: string
  value: string
  tone?: 'default' | 'ok' | 'warn'
}) {
  const color =
    tone === 'ok'
      ? 'text-emerald-600 dark:text-emerald-400'
      : tone === 'warn'
        ? 'text-red-600 dark:text-red-400'
        : 'text-gray-900 dark:text-white'
  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
      <div className="text-sm text-gray-500 dark:text-gray-400">{label}</div>
      <div className={`mt-2 text-2xl font-bold ${color}`}>{value}</div>
    </div>
  )
}

function Block({
  title,
  rows,
  formatter,
}: {
  title: string
  rows: { account_code: string; account_name: string; amount: number }[]
  formatter: Intl.NumberFormat
}) {
  const t = useTranslations('reports.financial')
  return (
    <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
      <div className="border-b border-gray-200 px-5 py-4 text-lg font-semibold text-gray-900 dark:border-gray-700 dark:text-white">
        {title}
      </div>
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-900/40">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">{t('table.code')}</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">{t('table.account')}</th>
              <th className="px-4 py-3 text-right font-medium text-gray-500 dark:text-gray-400">{t('table.amount')}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
            {rows.length === 0 ? (
              <tr>
                <td colSpan={3} className="px-4 py-6 text-center text-gray-500 dark:text-gray-400">
                  {t('common.noDataAvailable')}
                </td>
              </tr>
            ) : (
              rows.map((row) => (
                <tr key={`${title}-${row.account_code}-${row.account_name}`}>
                  <td className="px-4 py-3 font-mono text-gray-700 dark:text-gray-200">{row.account_code}</td>
                  <td className="px-4 py-3 text-gray-700 dark:text-gray-200">{row.account_name}</td>
                  <td className="px-4 py-3 text-right font-medium text-gray-900 dark:text-white">
                    {formatter.format(row.amount)}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
