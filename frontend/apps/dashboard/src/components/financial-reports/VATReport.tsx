'use client'

import { memo, useMemo } from 'react'
import { useTranslations } from 'next-intl'
import type { VATReport as VATReportType } from '@/types/financial-reports'

interface VATReportProps {
  report: VATReportType
  locale: string
}

export function VATReport({ report, locale }: VATReportProps) {
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
    <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
      <VATCard label={t('outputVat')} value={formatter.format(report.output_vat)} />
      <VATCard label={t('inputVat')} value={formatter.format(report.input_vat)} />
      <VATCard
        label={t('netVat')}
        value={formatter.format(report.net_vat)}
        tone={report.net_vat >= 0 ? 'warn' : 'ok'}
      />
    </div>
  )
}

export default memo(VATReport)

function VATCard({
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
        ? 'text-amber-600 dark:text-amber-400'
        : 'text-gray-900 dark:text-white'

  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
      <div className="text-sm text-gray-500 dark:text-gray-400">{label}</div>
      <div className={`mt-2 text-2xl font-bold ${color}`}>{value}</div>
    </div>
  )
}
