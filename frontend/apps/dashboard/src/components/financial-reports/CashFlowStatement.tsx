'use client'

import { memo, useMemo } from 'react'
import { useTranslations } from 'next-intl'
import type { CashFlowReport } from '@/types/financial-reports'

interface CashFlowStatementProps {
  report: CashFlowReport
  locale: string
}

export function CashFlowStatement({ report, locale }: CashFlowStatementProps) {
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
    <div className="space-y-4">
      <CashFlowRow label={t('cashFlowMetrics.operating')} value={report.operating_cash_flow} formatter={formatter} />
      <CashFlowRow label={t('cashFlowMetrics.investing')} value={report.investing_cash_flow} formatter={formatter} />
      <CashFlowRow label={t('cashFlowMetrics.financing')} value={report.financing_cash_flow} formatter={formatter} />
      <CashFlowRow
        label={t('cashFlowMetrics.net')}
        value={report.net_cash_flow}
        formatter={formatter}
        strong
      />
    </div>
  )
}

export default memo(CashFlowStatement)

function CashFlowRow({
  label,
  value,
  formatter,
  strong = false,
}: {
  label: string
  value: number
  formatter: Intl.NumberFormat
  strong?: boolean
}) {
  const color = value >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'
  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
      <div className="text-sm text-gray-500 dark:text-gray-400">{label}</div>
      <div className={`mt-2 ${strong ? 'text-3xl' : 'text-2xl'} font-bold ${color}`}>
        {formatter.format(value)}
      </div>
    </div>
  )
}
