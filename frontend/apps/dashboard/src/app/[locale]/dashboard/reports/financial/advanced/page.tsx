'use client'

import { useMemo, useState } from 'react'
import { usePathname } from 'next/navigation'
import { Loader, RefreshCw } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { getLocaleFromPath } from '@/lib/translations'
import {
  downloadCsv,
  getBreakEvenReport,
  getFinancialKPIs,
  getHRCostAnalysis,
  getMenuProfitability,
  getPredictiveCashFlow,
} from '@/lib/api/accountingAdvancedApi'

function toDateInput(value: Date): string {
  return value.toISOString().slice(0, 10)
}

function defaultRange() {
  const today = new Date()
  const from = new Date(today)
  from.setDate(today.getDate() - 29)
  return { date_from: toDateInput(from), date_to: toDateInput(today) }
}

function downloadBlob(filename: string, blob: Blob) {
  const url = URL.createObjectURL(blob)
  const anchor = document.createElement('a')
  anchor.href = url
  anchor.download = filename
  anchor.click()
  URL.revokeObjectURL(url)
}

export default function AdvancedFinancialInsightsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'
  const [range, setRange] = useState(defaultRange)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [data, setData] = useState<any>(null)

  const currencyFormat = useMemo(
    () => new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', { style: 'currency', currency: 'EGP', maximumFractionDigits: 2 }),
    [isRTL]
  )
  const numberFormat = useMemo(
    () => new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', { maximumFractionDigits: 2 }),
    [isRTL]
  )

  const load = async () => {
    setLoading(true)
    setError(null)
    try {
      const [kpis, menuProfitability, breakEven, predictiveCashFlow, hrCost] = await Promise.all([
        getFinancialKPIs(range),
        getMenuProfitability({ ...range, limit: 10 }),
        getBreakEvenReport(range),
        getPredictiveCashFlow({ history_months: 6, forecast_months: 3 }),
        getHRCostAnalysis(range),
      ])
      setData({ kpis, menuProfitability, breakEven, predictiveCashFlow, hrCost })
    } catch (err: any) {
      setError(err?.response?.data?.error || err?.message || 'Failed to load advanced financial insights')
    } finally {
      setLoading(false)
    }
  }

  const exportAudit = async () => {
    const blob = await downloadCsv('/audit/export', {
      type: 'all',
      date_from: range.date_from,
      date_to: range.date_to,
    })
    downloadBlob(`audit-trail-${range.date_from}-to-${range.date_to}.csv`, blob)
  }

  const exportSocialInsurance = async () => {
    const from = new Date(range.date_from)
    const blob = await downloadCsv('/reports/hr/social-insurance', {
      year: from.getFullYear(),
      month: from.getMonth() + 1,
    })
    downloadBlob(`social-insurance-${from.getFullYear()}-${from.getMonth() + 1}.csv`, blob)
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 xl:flex-row xl:items-end xl:justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Advanced Financial Insights</h1>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            Phase 8 analytics + Phase 9 compliance export center.
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Button variant="outline" onClick={exportAudit}>Export Audit CSV</Button>
          <Button variant="outline" onClick={exportSocialInsurance}>Export Social Insurance CSV</Button>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
        <input
          type="date"
          value={range.date_from}
          onChange={(e) => setRange((prev) => ({ ...prev, date_from: e.target.value }))}
          className="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-900"
        />
        <input
          type="date"
          value={range.date_to}
          onChange={(e) => setRange((prev) => ({ ...prev, date_to: e.target.value }))}
          className="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-900"
        />
        <Button onClick={load} disabled={loading}>
          {loading ? <Loader className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
          Refresh
        </Button>
      </div>

      {error ? (
        <div className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-800 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300">{error}</div>
      ) : null}

      {!data ? (
        <div className="rounded-xl border border-gray-200 bg-white p-6 text-sm text-gray-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400">
          Click refresh to load the advanced reports.
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-6 xl:grid-cols-2">
          <Card title="Financial KPI Summary">
            <Metric label="Revenue" value={currencyFormat.format(data.kpis.summary.revenue || 0)} />
            <Metric label="Expenses" value={currencyFormat.format(data.kpis.summary.expenses || 0)} />
            <Metric label="Net Profit" value={currencyFormat.format(data.kpis.summary.net_profit || 0)} />
            <Metric label="Cash Position" value={currencyFormat.format(data.kpis.summary.cash_position || 0)} />
          </Card>

          <Card title="Break-even Analysis">
            <Metric label="Break-even Revenue" value={currencyFormat.format(data.breakEven.summary.break_even_revenue || 0)} />
            <Metric label="Current Revenue" value={currencyFormat.format(data.breakEven.summary.revenue || 0)} />
            <Metric label="Break-even Gap" value={currencyFormat.format(data.breakEven.summary.break_even_gap || 0)} />
            <Metric label="Contribution Margin" value={`${numberFormat.format(data.breakEven.summary.contribution_margin_ratio || 0)}%`} />
          </Card>

          <Card title="Predictive Cash Flow">
            {(data.predictiveCashFlow.forecast || []).map((row: any) => (
              <Metric key={row.month} label={row.month} value={currencyFormat.format(row.projected_net_cash || 0)} />
            ))}
          </Card>

          <Card title="HR Cost Analysis">
            <Metric label="Total HR Cost" value={currencyFormat.format(data.hrCost.summary.total_hr_cost || 0)} />
            <Metric label="Payroll Cost" value={currencyFormat.format(data.hrCost.summary.payroll_cost || 0)} />
            <Metric label="Overtime Cost" value={currencyFormat.format(data.hrCost.summary.overtime_cost || 0)} />
            <Metric label="Labor % of Revenue" value={`${numberFormat.format(data.hrCost.summary.labor_cost_pct_of_revenue || 0)}%`} />
          </Card>

          <Card title="Top Menu Profitability">
            {(data.menuProfitability.items || []).slice(0, 10).map((item: any) => (
              <Metric
                key={`${item.product_id}-${item.name}`}
                label={`${item.name} (${numberFormat.format(item.gross_margin_pct || 0)}%)`}
                value={currencyFormat.format(item.gross_profit || 0)}
              />
            ))}
          </Card>
        </div>
      )}
    </div>
  )
}

function Card({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
      <h2 className="mb-3 text-lg font-semibold text-gray-900 dark:text-white">{title}</h2>
      <div className="space-y-2">{children}</div>
    </div>
  )
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-3">
      <span className="text-sm text-gray-600 dark:text-gray-400">{label}</span>
      <span className="text-sm font-semibold text-gray-900 dark:text-white">{value}</span>
    </div>
  )
}
