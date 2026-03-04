'use client'

import { useMemo, useState, type ReactNode } from 'react'
import { BarChart3, Loader, Lock, RefreshCw, Users } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { usePathname } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { useHRReports } from '@/hooks/useHRReports'
import { useUserPermissions } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { getLocaleFromPath } from '@/lib/translations'

function toDateInput(date: Date): string {
  return date.toISOString().slice(0, 10)
}

function defaultFilters() {
  const today = new Date()
  const from = new Date(today)
  from.setDate(today.getDate() - 29)
  return {
    date_from: toDateInput(from),
    date_to: toDateInput(today),
    year: today.getFullYear(),
  }
}

export default function HRReportsPage() {
  const t = useTranslations('reports.hrAdvanced')
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'
  const [filters, setFilters] = useState(defaultFilters)
  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const hasPermission = checkPermission(MODULE_IDS.REPORTS, PERMISSION_LEVELS.READ)
  const { headcount, turnover, absenteeism, overtime, leaveBalances, trainingCompliance, loading, error, refresh } = useHRReports(filters)

  const numberFormat = useMemo(
    () => new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', { maximumFractionDigits: 2 }),
    [isRTL]
  )
  const currencyFormat = useMemo(
    () => new Intl.NumberFormat(isRTL ? 'ar-EG' : 'en-US', { style: 'currency', currency: 'EGP', maximumFractionDigits: 2 }),
    [isRTL]
  )

  if (permissionLoading) {
    return <div className="flex h-96 items-center justify-center"><Loader className="h-8 w-8 animate-spin text-primary" /></div>
  }
  if (!hasPermission) {
    return (
      <div className="space-y-4">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{t('title')}</h1>
        <div className="rounded-xl border border-red-200 bg-red-50 p-6 dark:border-red-800 dark:bg-red-900/20">
          <div className="flex items-start gap-3">
            <Lock className="mt-0.5 h-5 w-5 text-red-600 dark:text-red-400" />
            <div>
              <div className="font-semibold text-red-900 dark:text-red-300">{t('permissionDeniedTitle')}</div>
              <p className="mt-1 text-sm text-red-800 dark:text-red-400">{t('permissionDeniedText')}</p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
        <div>
          <h1 className="flex items-center gap-2 text-3xl font-bold text-gray-900 dark:text-white">
            <Users className="h-7 w-7 text-primary" />
            {t('title')}
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            {t('subtitle')}
          </p>
        </div>
        <Button type="button" variant="outline" onClick={() => void refresh()} disabled={loading}>
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          {t('refresh')}
        </Button>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-3 xl:grid-cols-4">
        <FilterCard label={t('filters.dateFrom')}>
          <input
            type="date"
            value={filters.date_from}
            onChange={(e) => setFilters((prev) => ({ ...prev, date_from: e.target.value }))}
            className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-900"
          />
        </FilterCard>
        <FilterCard label={t('filters.dateTo')}>
          <input
            type="date"
            value={filters.date_to}
            onChange={(e) => setFilters((prev) => ({ ...prev, date_to: e.target.value }))}
            className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-900"
          />
        </FilterCard>
        <FilterCard label={t('filters.year')}>
          <input
            type="number"
            value={filters.year}
            min={2000}
            max={3000}
            onChange={(e) => setFilters((prev) => ({ ...prev, year: Number(e.target.value || new Date().getFullYear()) }))}
            className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-900"
          />
        </FilterCard>
        <div className="rounded-xl border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
          <div className="text-sm text-gray-600 dark:text-gray-400">{t('generated')}</div>
          <div className="mt-2 text-sm font-medium text-gray-900 dark:text-white">
            {headcount?.generated_at ? new Date(headcount.generated_at).toLocaleString() : '-'}
          </div>
        </div>
      </div>

      {error ? (
        <div className="rounded-xl border border-red-200 bg-red-50 p-5 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300">
          {error}
        </div>
      ) : null}

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
        <MetricCard title={t('metrics.totalEmployees')} value={numberFormat.format(headcount?.summary.total_employees || 0)} />
        <MetricCard title={t('metrics.turnoverRate')} value={`${numberFormat.format(turnover?.summary.turnover_rate || 0)}%`} />
        <MetricCard title={t('metrics.absenteeismRate')} value={`${numberFormat.format(absenteeism?.summary.absenteeism_rate || 0)}%`} />
        <MetricCard title={t('metrics.overtimeCost')} value={currencyFormat.format(overtime?.summary.total_overtime_cost || 0)} />
      </div>

      <div className="grid grid-cols-1 gap-6 xl:grid-cols-2">
        <SimpleTableCard
          title={t('tables.headcountByDepartment.title')}
          headers={[t('tables.headcountByDepartment.department'), t('tables.headcountByDepartment.employees')]}
          rows={(headcount?.by_department || []).map((item) => [item.department, numberFormat.format(item.count)])}
          empty={t('common.noData')}
        />
        <SimpleTableCard
          title={t('tables.turnoverTrend.title')}
          headers={[t('tables.turnoverTrend.month'), t('tables.turnoverTrend.hires'), t('tables.turnoverTrend.departures')]}
          rows={(turnover?.trend || []).map((item) => [item.month, numberFormat.format(item.hires), numberFormat.format(item.departures)])}
          empty={t('common.noData')}
        />
        <SimpleTableCard
          title={t('tables.topAbsentees.title')}
          headers={[t('tables.topAbsentees.employee'), t('tables.topAbsentees.absentDays')]}
          rows={(absenteeism?.top_absentees || []).map((item) => [item.employee_name || t('employeeFallback', { id: item.employee_id }), numberFormat.format(item.absent_days)])}
          empty={t('common.noAbsences')}
        />
        <SimpleTableCard
          title={t('tables.overtimeByEmployee.title')}
          headers={[t('tables.overtimeByEmployee.employee'), t('tables.overtimeByEmployee.department'), t('tables.overtimeByEmployee.hours'), t('tables.overtimeByEmployee.cost')]}
          rows={(overtime?.items || []).slice(0, 10).map((item) => [
            item.employee_name || t('employeeFallback', { id: item.employee_id }),
            item.department,
            numberFormat.format(item.overtime_hours),
            currencyFormat.format(item.overtime_cost),
          ])}
          empty={t('common.noOvertime')}
        />
        <SimpleTableCard
          title={t('tables.leaveBalances.title')}
          headers={[t('tables.leaveBalances.employee'), t('tables.leaveBalances.annualUsed'), t('tables.leaveBalances.annualRemaining')]}
          rows={(leaveBalances?.items || []).map((item) => [
            item.employee_name || t('employeeFallback', { id: item.employee_id }),
            numberFormat.format(item.annual_leave_used),
            numberFormat.format(item.annual_leave_remaining),
          ])}
          empty={t('common.noData')}
        />
        <SimpleTableCard
          title={t('tables.trainingComplianceByDepartment.title')}
          headers={[t('tables.trainingComplianceByDepartment.department'), t('tables.trainingComplianceByDepartment.certified'), t('tables.trainingComplianceByDepartment.total'), t('tables.trainingComplianceByDepartment.rate')]}
          rows={(trainingCompliance?.by_department || []).map((item) => [
            item.department,
            numberFormat.format(item.certified_employees),
            numberFormat.format(item.total_employees),
            `${numberFormat.format(item.compliance_rate)}%`,
          ])}
          empty={t('common.noData')}
        />
      </div>

      <div className="rounded-xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
        <div className="mb-2 flex items-center gap-2 text-lg font-semibold text-gray-900 dark:text-white">
          <BarChart3 className="h-5 w-5 text-primary" />
          {t('trainingComplianceSummary.title')}
        </div>
        <div className="text-sm text-gray-600 dark:text-gray-400">
          {t('trainingComplianceSummary.certifiedEmployees')}: <span className="font-semibold text-gray-900 dark:text-white">{numberFormat.format(trainingCompliance?.summary.certified_employees || 0)}</span>
          {' '}{t('trainingComplianceSummary.of')}{' '}
          <span className="font-semibold text-gray-900 dark:text-white">{numberFormat.format(trainingCompliance?.summary.active_employees || 0)}</span>
          {' '}({numberFormat.format(trainingCompliance?.summary.compliance_rate || 0)}%)
        </div>
      </div>
    </div>
  )
}

function FilterCard({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
      <label className="mb-2 block text-sm text-gray-600 dark:text-gray-400">{label}</label>
      {children}
    </div>
  )
}

function MetricCard({ title, value }: { title: string; value: string }) {
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
      <div className="text-sm text-gray-500 dark:text-gray-400">{title}</div>
      <div className="mt-3 text-2xl font-bold text-gray-900 dark:text-white">{value}</div>
    </div>
  )
}

function SimpleTableCard({
  title,
  headers,
  rows,
  empty,
}: {
  title: string
  headers: string[]
  rows: string[][]
  empty: string
}) {
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
      <h2 className="mb-3 text-lg font-semibold text-gray-900 dark:text-white">{title}</h2>
      {rows.length === 0 ? (
        <div className="text-sm text-gray-500 dark:text-gray-400">{empty}</div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full min-w-[420px] text-sm">
            <thead>
              <tr className="border-b border-gray-200 text-left text-gray-500 dark:border-gray-700 dark:text-gray-400">
                {headers.map((header) => (
                  <th key={header} className="px-2 py-2 font-medium">{header}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((row, idx) => (
                <tr key={`${title}-${idx}`} className="border-b border-gray-100 last:border-0 dark:border-gray-800">
                  {row.map((cell, cellIdx) => (
                    <td key={`${title}-${idx}-${cellIdx}`} className="px-2 py-2 text-gray-900 dark:text-gray-100">{cell}</td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
