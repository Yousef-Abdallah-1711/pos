'use client'

import { useEffect, useMemo, useState } from 'react'
import { AlertCircle, Loader2, RefreshCcw, Scale } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { formatCurrency } from '@/components/accounting/accounting-utils'
import { getTrialBalance } from '@/lib/api/accountingApi'
import { createTranslator } from '@/lib/translations'
import { debugError } from '@/lib/debug'
import type { TrialBalanceResponse, TrialBalanceRow } from '@/types/accounting'
import type { Locale } from '@/lib/translations'

interface TrialBalanceTableProps {
  locale: Locale
}

export function TrialBalanceTable({ locale }: TrialBalanceTableProps) {
  const t = useMemo(() => createTranslator(locale), [locale])
  const [asOfDate, setAsOfDate] = useState(() => new Date().toISOString().slice(0, 10))
  const [loading, setLoading] = useState(false)
  const [rows, setRows] = useState<TrialBalanceRow[]>([])
  const [summary, setSummary] = useState<TrialBalanceResponse | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    void loadTrialBalance(asOfDate)
  }, [])

  const loadTrialBalance = async (date: string) => {
    setLoading(true)
    setError(null)

    try {
      const response = await getTrialBalance({ as_of: date || undefined })
      setSummary(response)
      setRows(response.rows ?? [])
    } catch (apiError: any) {
      const message =
        apiError?.response?.data?.error ||
        apiError?.message ||
        t('accounting.loadTrialBalanceError', 'Failed to load trial balance.')
      setError(message)
      setRows([])
      setSummary(null)
      debugError('[ACCOUNTING] load trial balance failed', apiError)
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-4 rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-gray-800 dark:bg-gray-900 sm:flex-row sm:items-end">
        <div className="w-full max-w-sm space-y-2">
          <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('accounting.asOfDate', 'As of date')}
          </label>
          <Input type="date" value={asOfDate} onChange={(event) => setAsOfDate(event.target.value)} />
        </div>

        <Button type="button" onClick={() => void loadTrialBalance(asOfDate)} disabled={loading} className="gap-2">
          {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <RefreshCcw className="h-4 w-4" />}
          {t('accounting.refreshTrialBalance', 'Refresh')}
        </Button>

        <div className="sm:ml-auto">
          <div
            className={`inline-flex items-center gap-2 rounded-full px-3 py-2 text-sm font-medium ${
              summary?.is_balanced
                ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300'
                : 'bg-amber-100 text-amber-700 dark:bg-amber-950/40 dark:text-amber-300'
            }`}
          >
            <Scale className="h-4 w-4" />
            {summary?.is_balanced
              ? t('accounting.trialBalanceBalanced', 'Balanced')
              : t('accounting.trialBalanceUnbalanced', 'Out of balance')}
          </div>
        </div>
      </div>

      {error ? (
        <div className="rounded-xl border border-red-200 bg-red-50 p-4 dark:border-red-900 dark:bg-red-950/30">
          <div className="flex items-start gap-3">
            <AlertCircle className="mt-0.5 h-5 w-5 text-red-600 dark:text-red-400" />
            <div>
              <h3 className="font-semibold text-red-900 dark:text-red-300">
                {t('accounting.loadErrorTitle', 'Unable to load records')}
              </h3>
              <p className="mt-1 text-sm text-red-700 dark:text-red-300">{error}</p>
            </div>
          </div>
        </div>
      ) : null}

      <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-gray-800 dark:bg-gray-900">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t('accounting.accountCode', 'Account code')}</TableHead>
              <TableHead>{t('accounting.accountName', 'Account name')}</TableHead>
              <TableHead>{t('accounting.entryType', 'Entry type')}</TableHead>
              <TableHead className="text-right">{t('accounting.totalDebit', 'Total debit')}</TableHead>
              <TableHead className="text-right">{t('accounting.totalCredit', 'Total credit')}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={5} className="py-12 text-center text-gray-500 dark:text-gray-400">
                  <div className="flex items-center justify-center gap-2">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    {t('accounting.loadingTrialBalance', 'Loading trial balance...')}
                  </div>
                </TableCell>
              </TableRow>
            ) : rows.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="py-12 text-center text-gray-500 dark:text-gray-400">
                  {t('accounting.noTrialBalanceRows', 'No trial balance data is available yet.')}
                </TableCell>
              </TableRow>
            ) : (
              rows.map((row) => (
                <TableRow key={row.account_id}>
                  <TableCell className="font-medium text-gray-900 dark:text-white">
                    {row.account_code}
                  </TableCell>
                  <TableCell>{row.account_name}</TableCell>
                  <TableCell className="capitalize text-gray-600 dark:text-gray-400">
                    {row.account_type}
                  </TableCell>
                  <TableCell className="text-right font-medium text-gray-900 dark:text-white">
                    {formatCurrency(row.debit_balance, locale)}
                  </TableCell>
                  <TableCell className="text-right font-medium text-gray-900 dark:text-white">
                    {formatCurrency(row.credit_balance, locale)}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
          <TableFooter>
            <TableRow>
              <TableCell colSpan={3}>{t('accounting.total', 'Total')}</TableCell>
              <TableCell className="text-right">
                {formatCurrency(summary?.total_debit ?? 0, locale)}
              </TableCell>
              <TableCell className="text-right">
                {formatCurrency(summary?.total_credit ?? 0, locale)}
              </TableCell>
            </TableRow>
          </TableFooter>
        </Table>
      </div>
    </div>
  )
}

export default TrialBalanceTable
