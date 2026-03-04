'use client'

import { memo, useMemo, useState } from 'react'
import { useTranslations } from 'next-intl'
import type { ProfitLossReport } from '@/types/financial-reports'

interface ProfitLossStatementProps {
  report: ProfitLossReport
  locale: string
}

export function ProfitLossStatement({ report, locale }: ProfitLossStatementProps) {
  const t = useTranslations('reports.financial')
  const [showRevenueDetails, setShowRevenueDetails] = useState(true)
  const [showExpenseDetails, setShowExpenseDetails] = useState(true)
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
      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
        <SummaryCard
          label={t('metrics.totalRevenue')}
          value={formatter.format(report.total_revenue)}
          comparisonValue={report.comparison_total_revenue}
          changePercent={report.revenue_change_percent}
          formatter={formatter}
        />
        <SummaryCard
          label={t('metrics.totalExpenses')}
          value={formatter.format(report.total_expenses)}
          comparisonValue={report.comparison_total_expenses}
          changePercent={report.expense_change_percent}
          formatter={formatter}
        />
        <SummaryCard
          label={t('metrics.netProfit')}
          value={formatter.format(report.net_profit)}
          positive={report.net_profit >= 0}
          comparisonValue={report.comparison_net_profit}
          changePercent={report.net_profit_change_percent}
          formatter={formatter}
        />
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <SectionTable
          title={t('sections.revenueAccounts')}
          rows={report.revenue}
          formatter={formatter}
          t={t}
          expanded={showRevenueDetails}
          onToggle={() => setShowRevenueDetails((prev) => !prev)}
        />
        <SectionTable
          title={t('sections.expenseAccounts')}
          rows={report.expenses}
          formatter={formatter}
          t={t}
          expanded={showExpenseDetails}
          onToggle={() => setShowExpenseDetails((prev) => !prev)}
        />
      </div>
    </div>
  )
}

export default memo(ProfitLossStatement)

function SummaryCard({
  label,
  value,
  positive,
  comparisonValue,
  changePercent,
  formatter,
}: {
  label: string
  value: string
  positive?: boolean
  comparisonValue?: number
  changePercent?: number
  formatter?: Intl.NumberFormat
}) {
  const t = useTranslations('reports.financial')
  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
      <div className="text-sm text-gray-500 dark:text-gray-400">{label}</div>
      <div
        className={`mt-2 text-2xl font-bold ${
          positive === undefined
            ? 'text-gray-900 dark:text-white'
            : positive
              ? 'text-emerald-600 dark:text-emerald-400'
              : 'text-red-600 dark:text-red-400'
        }`}
      >
        {value}
      </div>
      {comparisonValue !== undefined && formatter ? (
        <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">
          {t('common.prev')}: {formatter.format(comparisonValue)}{' '}
          {changePercent !== undefined ? `(${changePercent.toFixed(2)}%)` : ''}
        </div>
      ) : null}
    </div>
  )
}

function SectionTable({
  title,
  rows,
  formatter,
  t,
  expanded,
  onToggle,
}: {
  title: string
  rows: { account_code: string; account_name: string; amount: number }[]
  formatter: Intl.NumberFormat
  t: (key: string, values?: Record<string, string | number>) => string
  expanded: boolean
  onToggle: () => void
}) {
  return (
    <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
      <div className="border-b border-gray-200 px-5 py-4 dark:border-gray-700">
        <button
          type="button"
          onClick={onToggle}
          className="flex w-full items-center justify-between text-left text-lg font-semibold text-gray-900 dark:text-white"
        >
          <span>{title}</span>
          <span className="text-sm text-gray-500 dark:text-gray-400">{expanded ? t('common.hide') : t('common.show')}</span>
        </button>
      </div>
      {expanded ? <div className="overflow-x-auto">
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
                <tr key={`${row.account_code}-${row.account_name}`}>
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
      </div> : null}
    </div>
  )
}
