'use client'

import { useEffect, useMemo, useState } from 'react'
import { AlertCircle, Loader2, Search } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { flattenAccounts, formatAccountLabel, formatCurrency, formatDate } from '@/components/accounting/accounting-utils'
import { getLedger } from '@/lib/api/accountingApi'
import { createTranslator } from '@/lib/translations'
import { debugError } from '@/lib/debug'
import type { Account, GeneralLedgerLine, LedgerPeriodSummary } from '@/types/accounting'

interface LedgerViewProps {
  accounts: Account[]
  locale: string
}

export function LedgerView({ accounts, locale }: LedgerViewProps) {
  const t = useMemo(() => createTranslator(locale), [locale])
  const selectableAccounts = useMemo(
    () =>
      flattenAccounts(accounts).filter((account) => account.is_active && !account.is_header),
    [accounts]
  )

  const [selectedAccountId, setSelectedAccountId] = useState<string>('')
  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')
  const [loading, setLoading] = useState(false)
  const [loadingMore, setLoadingMore] = useState(false)
  const [lines, setLines] = useState<GeneralLedgerLine[]>([])
  const [balance, setBalance] = useState(0)
  const [nextCursor, setNextCursor] = useState<number | null>(null)
  const [periodSummaries, setPeriodSummaries] = useState<LedgerPeriodSummary[]>([])
  const [visibleRows, setVisibleRows] = useState(120)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!selectedAccountId && selectableAccounts.length > 0) {
      setSelectedAccountId(String(selectableAccounts[0].id))
    }
  }, [selectedAccountId, selectableAccounts])

  useEffect(() => {
    if (selectedAccountId) {
      setVisibleRows(120)
      void loadLedger(selectedAccountId, { append: false })
    }
  }, [selectedAccountId])

  const selectedAccount = useMemo(
    () => selectableAccounts.find((account) => String(account.id) === selectedAccountId) ?? null,
    [selectableAccounts, selectedAccountId]
  )

  const displayedLines = useMemo(() => lines.slice(0, visibleRows), [lines, visibleRows])

  const loadLedger = async (
    accountId: string,
    options?: { append: boolean; cursor?: number | null }
  ) => {
    if (!accountId) {
      return
    }

    const append = options?.append ?? false
    if (append) {
      setLoadingMore(true)
    } else {
      setLoading(true)
      setError(null)
    }

    try {
      const response = await getLedger({
        account_id: Number(accountId),
        start: startDate || undefined,
        end: endDate || undefined,
        limit: 200,
        cursor: options?.cursor ?? undefined,
      })
      const incomingLines = response.lines ?? []
      setLines((prev) => (append ? [...prev, ...incomingLines] : incomingLines))
      setBalance(response.balance ?? 0)
      setNextCursor(response.next_cursor ?? null)
      if (!append) {
        setPeriodSummaries(response.period_summaries ?? [])
      }
    } catch (apiError: any) {
      const message =
        apiError?.response?.data?.error ||
        apiError?.message ||
        t('accounting.loadLedgerError', 'Failed to load general ledger.')
      setError(message)
      if (!append) {
        setLines([])
        setBalance(0)
        setNextCursor(null)
        setPeriodSummaries([])
      }
      debugError('[ACCOUNTING] load ledger failed', apiError)
      toast.error(message)
    } finally {
      if (append) {
        setLoadingMore(false)
      } else {
        setLoading(false)
      }
    }
  }

  return (
    <div className="space-y-4">
      <div className="grid gap-4 rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-gray-800 dark:bg-gray-900 lg:grid-cols-[2fr_1fr_1fr_auto]">
        <div className="space-y-2">
          <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('accounting.selectAccount', 'Select account')}
          </label>
          <Select value={selectedAccountId} onValueChange={setSelectedAccountId}>
            <SelectTrigger>
              <SelectValue placeholder={t('accounting.selectAccount', 'Select account')} />
            </SelectTrigger>
            <SelectContent>
              {selectableAccounts.map((account) => (
                <SelectItem key={account.id} value={String(account.id)}>
                  {formatAccountLabel(account)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('accounting.startDate', 'Start date')}
          </label>
          <Input type="date" value={startDate} onChange={(event) => setStartDate(event.target.value)} />
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('accounting.endDate', 'End date')}
          </label>
          <Input type="date" value={endDate} onChange={(event) => setEndDate(event.target.value)} />
        </div>

        <div className="flex items-end">
          <Button
            type="button"
            onClick={() => {
              setVisibleRows(120)
              void loadLedger(selectedAccountId, { append: false })
            }}
            disabled={!selectedAccountId || loading}
            className="w-full gap-2 lg:w-auto"
          >
            {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Search className="h-4 w-4" />}
            {t('accounting.filter', 'Filter')}
          </Button>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <div className="rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-gray-800 dark:bg-gray-900">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t('accounting.account', 'Account')}
          </p>
          <p className="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
            {selectedAccount ? formatAccountLabel(selectedAccount) : '-'}
          </p>
        </div>

        <div className="rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-gray-800 dark:bg-gray-900">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t('accounting.balance', 'Balance')}
          </p>
          <p className="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
            {formatCurrency(balance, locale)}
          </p>
        </div>
      </div>

      {periodSummaries.length > 0 ? (
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-gray-800 dark:bg-gray-900">
          <div className="border-b border-gray-200 px-4 py-3 text-sm font-semibold text-gray-700 dark:border-gray-800 dark:text-gray-200">
            {t('accounting.periodSummary', 'Fiscal period summary')}
          </div>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t('accounting.fiscalPeriod', 'Fiscal period')}</TableHead>
                  <TableHead className="text-right">{t('accounting.entries', 'Entries')}</TableHead>
                  <TableHead className="text-right">{t('accounting.totalDebit', 'Total debit')}</TableHead>
                  <TableHead className="text-right">{t('accounting.totalCredit', 'Total credit')}</TableHead>
                  <TableHead className="text-right">{t('accounting.netChange', 'Net change')}</TableHead>
                  <TableHead className="text-right">{t('accounting.closingBalance', 'Closing balance')}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {periodSummaries.map((summary) => (
                  <TableRow key={summary.fiscal_period}>
                    <TableCell className="font-medium">{summary.fiscal_period}</TableCell>
                    <TableCell className="text-right">{summary.entry_count}</TableCell>
                    <TableCell className="text-right">{formatCurrency(summary.debit_total, locale)}</TableCell>
                    <TableCell className="text-right">{formatCurrency(summary.credit_total, locale)}</TableCell>
                    <TableCell className="text-right">{formatCurrency(summary.net_change, locale)}</TableCell>
                    <TableCell className="text-right">{formatCurrency(summary.closing_balance, locale)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </div>
      ) : null}

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
              <TableHead>{t('accounting.entryDate', 'Entry date')}</TableHead>
              <TableHead>{t('accounting.entryNumber', 'Entry #')}</TableHead>
              <TableHead>{t('accounting.description', 'Description')}</TableHead>
              <TableHead>{t('accounting.reference', 'Reference')}</TableHead>
              <TableHead className="text-right">{t('accounting.totalDebit', 'Total debit')}</TableHead>
              <TableHead className="text-right">{t('accounting.totalCredit', 'Total credit')}</TableHead>
              <TableHead className="text-right">{t('accounting.runningBalance', 'Running balance')}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={7} className="py-12 text-center text-gray-500 dark:text-gray-400">
                  <div className="flex items-center justify-center gap-2">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    {t('accounting.loadingLedger', 'Loading ledger...')}
                  </div>
                </TableCell>
              </TableRow>
            ) : lines.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="py-12 text-center text-gray-500 dark:text-gray-400">
                  {t('accounting.noLedgerLines', 'No ledger activity matches the selected filters.')}
                </TableCell>
              </TableRow>
            ) : (
              displayedLines.map((line) => (
                <TableRow key={`${line.journal_entry_id}-${line.journal_line_id}`}>
                  <TableCell>{formatDate(line.entry_date, locale)}</TableCell>
                  <TableCell className="font-medium text-gray-900 dark:text-white">
                    {line.entry_number}
                  </TableCell>
                  <TableCell>
                    <div className="space-y-1">
                      <p className="text-gray-900 dark:text-white">{line.description}</p>
                      {line.line_description ? (
                        <p className="text-xs text-gray-500 dark:text-gray-400">
                          {line.line_description}
                        </p>
                      ) : null}
                    </div>
                  </TableCell>
                  <TableCell className="text-gray-600 dark:text-gray-400">
                    {line.reference_number || line.reference_type || '-'}
                  </TableCell>
                  <TableCell className="text-right font-medium text-gray-900 dark:text-white">
                    {formatCurrency(line.debit_amount, locale)}
                  </TableCell>
                  <TableCell className="text-right font-medium text-gray-900 dark:text-white">
                    {formatCurrency(line.credit_amount, locale)}
                  </TableCell>
                  <TableCell className="text-right font-semibold text-primary-700 dark:text-primary-300">
                    {formatCurrency(line.running_balance, locale)}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {lines.length > displayedLines.length ? (
        <div className="flex justify-center">
          <Button
            type="button"
            variant="outline"
            onClick={() => setVisibleRows((prev) => prev + 120)}
          >
            {t('accounting.showMoreRows', 'Show more rows')}
          </Button>
        </div>
      ) : null}

      {nextCursor ? (
        <div className="flex justify-center">
          <Button
            type="button"
            variant="secondary"
            onClick={() => void loadLedger(selectedAccountId, { append: true, cursor: nextCursor })}
            disabled={loadingMore}
            className="gap-2"
          >
            {loadingMore ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
            {t('accounting.loadMoreLedger', 'Load older entries')}
          </Button>
        </div>
      ) : null}
    </div>
  )
}

export default LedgerView
