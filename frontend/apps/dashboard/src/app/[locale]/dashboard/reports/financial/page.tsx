'use client'

import Link from 'next/link'
import { useMemo } from 'react'
import { BarChart3, FileText, Landmark, Package, Receipt, Wallet } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { usePathname } from 'next/navigation'
import { useFinancialReport } from '@/hooks/useFinancialReports'
import { getVATReport } from '@/lib/api/financialReportsApi'
import { getLocaleFromPath } from '@/lib/translations'

function getDateRange() {
  const to = new Date()
  const from = new Date()
  from.setDate(to.getDate() - 29)
  return {
    date_from: toDateInputValue(from),
    date_to: toDateInputValue(to),
  }
}

function toDateInputValue(date: Date) {
  return date.toISOString().slice(0, 10)
}

export default function FinancialReportsHubPage() {
  const t = useTranslations('reports.financial')
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const dates = useMemo(() => getDateRange(), [])
  const formatter = useMemo(
    () =>
      new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
        style: 'currency',
        currency: 'EGP',
        maximumFractionDigits: 2,
      }),
    [locale]
  )

  const { data: vatSummary } = useFinancialReport({
    cacheKey: `financial-vat-${dates.date_from}-${dates.date_to}`,
    deps: [dates.date_from, dates.date_to],
    fetcher: () => getVATReport(dates),
  })

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{t('title')}</h1>
        <p className="mt-1 text-gray-500 dark:text-gray-400">
          {t('subtitle')}
        </p>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
        <KpiCard
          icon={Receipt}
          title={t('netVat')}
          value={vatSummary ? formatter.format(vatSummary.net_vat) : formatter.format(0)}
        />
        <KpiCard
          icon={BarChart3}
          title={t('outputVat')}
          value={vatSummary ? formatter.format(vatSummary.output_vat) : formatter.format(0)}
        />
        <KpiCard
          icon={FileText}
          title={t('inputVat')}
          value={vatSummary ? formatter.format(vatSummary.input_vat) : formatter.format(0)}
        />
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
        <ReportCard
          href={`/${locale}/dashboard/reports/financial/profit-loss`}
          icon={FileText}
          title={t('profitLoss')}
          description={t('descriptions.profitLoss')}
        />
        <ReportCard
          href={`/${locale}/dashboard/reports/financial/balance-sheet`}
          icon={Landmark}
          title={t('balanceSheet')}
          description={t('descriptions.balanceSheet')}
        />
        <ReportCard
          href={`/${locale}/dashboard/reports/financial/cash-flow`}
          icon={BarChart3}
          title={t('cashFlow')}
          description={t('descriptions.cashFlow')}
        />
        <ReportCard
          href={`/${locale}/dashboard/reports/financial/vat`}
          icon={Receipt}
          title={t('vat')}
          description={t('descriptions.vat')}
        />
        <ReportCard
          href={`/${locale}/dashboard/reports/financial/receivables-aging`}
          icon={Wallet}
          title={t('receivablesAging')}
          description={t('descriptions.receivablesAging')}
        />
        <ReportCard
          href={`/${locale}/dashboard/reports/financial/payables-aging`}
          icon={Wallet}
          title={t('payablesAging')}
          description={t('descriptions.payablesAging')}
        />
        <ReportCard
          href={`/${locale}/dashboard/reports/financial/inventory-valuation`}
          icon={Package}
          title={t('inventoryValuation')}
          description={t('descriptions.inventoryValuation')}
        />
        <ReportCard
          href={`/${locale}/dashboard/reports/financial/advanced`}
          icon={BarChart3}
          title={locale === 'ar' ? 'تحليلات متقدمة' : 'Advanced Insights'}
          description={locale === 'ar' ? 'مؤشرات مالية، نقطة التعادل، تدفق نقدي متوقع وتكلفة العمالة.' : 'KPI dashboard, break-even, predictive cash flow, and labor cost analytics.'}
        />
      </div>
    </div>
  )
}

function KpiCard({
  icon: Icon,
  title,
  value,
}: {
  icon: React.ComponentType<{ className?: string }>
  title: string
  value: string
}) {
  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
      <div className="flex items-center justify-between">
        <div className="text-sm text-gray-500 dark:text-gray-400">{title}</div>
        <div className="rounded-xl bg-primary-100 p-2 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300">
          <Icon className="h-4 w-4" />
        </div>
      </div>
      <div className="mt-3 text-2xl font-bold text-gray-900 dark:text-white">{value}</div>
    </div>
  )
}

function ReportCard({
  href,
  icon: Icon,
  title,
  description,
}: {
  href: string
  icon: React.ComponentType<{ className?: string }>
  title: string
  description: string
}) {
  return (
    <Link
      href={href}
      className="block rounded-2xl border border-gray-200 bg-white p-5 transition hover:border-primary-300 hover:shadow-sm dark:border-gray-700 dark:bg-gray-800 dark:hover:border-primary-800"
    >
      <div className="rounded-xl bg-gray-100 p-2 text-gray-700 dark:bg-gray-900 dark:text-gray-200 inline-flex">
        <Icon className="h-4 w-4" />
      </div>
      <h3 className="mt-4 text-lg font-semibold text-gray-900 dark:text-white">{title}</h3>
      <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{description}</p>
    </Link>
  )
}
