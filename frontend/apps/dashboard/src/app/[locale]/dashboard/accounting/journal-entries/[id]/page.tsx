'use client'

import Link from 'next/link'
import { use, useEffect, useMemo, useState } from 'react'
import { ArrowLeft, Loader2, Lock, Pencil, Slash, Upload } from 'lucide-react'
import { toast } from 'sonner'
import { JournalEntryForm } from '@/components/accounting/JournalEntryForm'
import {
  flattenAccounts,
  formatCurrency,
  formatDate,
  getJournalStatusClasses,
} from '@/components/accounting/accounting-utils'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useUserPermissions } from '@/hooks/useRbac'
import {
  getAccounts,
  getJournalEntry,
  postJournalEntry,
  updateJournalEntry,
  voidJournalEntry,
} from '@/lib/api/accountingApi'
import { cn } from '@/lib/utils'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'
import { usePathname } from 'next/navigation'
import type { Account, JournalEntry, JournalEntryFormInput } from '@/types/accounting'

interface JournalEntryDetailPageProps {
  params: Promise<{ id: string }>
}

export default function JournalEntryDetailPage({
  params: paramsPromise,
}: JournalEntryDetailPageProps) {
  const { id } = use(paramsPromise)
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canView = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.READ)
  const canWrite = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.WRITE)

  const [entry, setEntry] = useState<JournalEntry | null>(null)
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [showEditForm, setShowEditForm] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)

  useEffect(() => {
    if (permissionLoading || !canView) {
      return
    }

    void loadData()
  }, [canView, id, permissionLoading])

  const flatAccounts = useMemo(() => flattenAccounts(accounts), [accounts])

  const loadData = async () => {
    setLoading(true)
    setApiError(null)

    try {
      const numericId = Number(id)
      const [entryData, accountData] = await Promise.all([
        getJournalEntry(numericId),
        getAccounts({ include_inactive: false }),
      ])

      setEntry(entryData)
      setAccounts(accountData)
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to load journal entry'
      setApiError(message)
      debugError('[ACCOUNTING] load journal entry detail failed', error)
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  const handleUpdate = async (payload: JournalEntryFormInput) => {
    if (!entry) {
      return
    }

    setSubmitting(true)

    try {
      const updated = await updateJournalEntry(entry.id, payload)
      setEntry(updated)
      setShowEditForm(false)
      toast.success('Journal entry updated successfully')
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to update journal entry'
      debugError('[ACCOUNTING] update journal entry failed', error)
      toast.error(message)
      throw error
    } finally {
      setSubmitting(false)
    }
  }

  const handlePost = async () => {
    if (!entry) {
      return
    }

    setSubmitting(true)
    try {
      const posted = await postJournalEntry(entry.id)
      setEntry(posted)
      toast.success('Journal entry posted successfully')
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to post journal entry'
      debugError('[ACCOUNTING] post journal entry failed', error)
      toast.error(message)
    } finally {
      setSubmitting(false)
    }
  }

  const handleVoid = async () => {
    if (!entry) {
      return
    }

    const reason = window.prompt('Enter a reason for voiding this entry:')
    if (reason === null) {
      return
    }

    setSubmitting(true)
    try {
      const voided = await voidJournalEntry(entry.id, reason || undefined)
      setEntry(voided)
      toast.success('Journal entry voided successfully')
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to void journal entry'
      debugError('[ACCOUNTING] void journal entry failed', error)
      toast.error(message)
    } finally {
      setSubmitting(false)
    }
  }

  if (permissionLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!canView) {
    return (
      <div className="rounded-xl border border-red-200 bg-red-50 p-6 dark:border-red-900 dark:bg-red-950/30">
        <div className="flex items-start gap-3">
          <Lock className="mt-0.5 h-5 w-5 text-red-600 dark:text-red-400" />
          <div>
            <h2 className="font-semibold text-red-900 dark:text-red-300">
              {t('accounting.accessDeniedTitle', 'Access denied')}
            </h2>
            <p className="mt-1 text-sm text-red-700 dark:text-red-300">
              {t(
                'accounting.accessDeniedText',
                'You do not have permission to view accounting records.'
              )}
            </p>
          </div>
        </div>
      </div>
    )
  }

  if (loading && !entry) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (apiError && !entry) {
    return (
      <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-red-700 dark:border-red-900 dark:bg-red-950/30 dark:text-red-300">
        {apiError}
      </div>
    )
  }

  if (!entry) {
    return (
      <div className="space-y-4">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
          {t('accounting.journalEntries', 'Journal Entries')}
        </h1>
        <p className="text-gray-500 dark:text-gray-400">
          {t('accounting.entryNotFound', 'The requested journal entry was not found.')}
        </p>
      </div>
    )
  }

  const canEditDraft = canWrite && entry.status === 'draft'
  const canPostDraft = canWrite && entry.status === 'draft'
  const canVoidEntry = canWrite && entry.status !== 'void'

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div>
          <Link
            href={`/${locale}/dashboard/accounting/journal-entries`}
            className="inline-flex items-center gap-2 text-sm font-medium text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
          >
            <ArrowLeft className="h-4 w-4" />
            {t('accounting.backToJournalEntries', 'Back to journal entries')}
          </Link>
          <h1 className="mt-3 text-2xl font-bold text-gray-900 dark:text-white">
            {entry.entry_number}
          </h1>
          <p className="mt-1 text-gray-500 dark:text-gray-400">{entry.description}</p>
        </div>

        <div className="flex flex-wrap gap-2">
          {canEditDraft ? (
            <Button variant="outline" onClick={() => setShowEditForm(true)} disabled={submitting}>
              <Pencil className="mr-2 h-4 w-4" />
              {t('accounting.editEntry', 'Edit')}
            </Button>
          ) : null}
          {canPostDraft ? (
            <Button onClick={() => void handlePost()} disabled={submitting}>
              <Upload className="mr-2 h-4 w-4" />
              {t('accounting.postEntry', 'Post')}
            </Button>
          ) : null}
          {canVoidEntry ? (
            <Button
              variant="outline"
              onClick={() => void handleVoid()}
              disabled={submitting}
              className="text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
            >
              <Slash className="mr-2 h-4 w-4" />
              {t('accounting.voidEntry', 'Void')}
            </Button>
          ) : null}
        </div>
      </div>

      <div className="grid gap-4 lg:grid-cols-4">
        <SummaryCard label={t('accounting.entryDate', 'Entry date')} value={formatDate(entry.entry_date, locale)} />
        <SummaryCard label={t('accounting.entryType', 'Entry type')} value={entry.entry_type} />
        <SummaryCard label={t('accounting.totalDebit', 'Total debit')} value={formatCurrency(entry.total_debit, locale)} />
        <SummaryCard label={t('accounting.totalCredit', 'Total credit')} value={formatCurrency(entry.total_credit, locale)} />
      </div>

      <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
        <div className="flex flex-wrap items-center gap-3">
          <span
            className={cn(
              'inline-flex rounded-full border px-2.5 py-1 text-xs font-medium capitalize',
              getJournalStatusClasses(entry.status)
            )}
          >
            {entry.status}
          </span>
          {entry.reference_number ? (
            <span className="text-sm text-gray-500 dark:text-gray-400">
              Ref: {entry.reference_number}
            </span>
          ) : null}
          {entry.void_reason ? (
            <span className="text-sm text-red-600 dark:text-red-400">
              Void reason: {entry.void_reason}
            </span>
          ) : null}
        </div>
        {entry.notes ? (
          <p className="mt-4 text-sm text-gray-600 dark:text-gray-300">{entry.notes}</p>
        ) : null}
      </div>

      <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
        <div className="border-b border-gray-200 px-5 py-4 dark:border-gray-700">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            {t('accounting.journalLines', 'Journal lines')}
          </h2>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-900/40">
              <tr>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Account
                </th>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Description
                </th>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Debit
                </th>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Credit
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
              {entry.lines.map((line) => (
                <tr key={line.id || `${line.account_id}-${line.line_order}`}>
                  <td className="px-4 py-4">
                    <div className="font-medium text-gray-900 dark:text-white">
                      {line.account_code || `#${line.account_id}`}
                    </div>
                  </td>
                  <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                    {line.description || '-'}
                  </td>
                  <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                    {formatCurrency(line.debit_amount, locale)}
                  </td>
                  <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                    {formatCurrency(line.credit_amount, locale)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <Dialog open={showEditForm} onOpenChange={setShowEditForm}>
        <DialogContent className="max-w-6xl">
          <DialogHeader>
            <DialogTitle>{t('accounting.editEntry', 'Edit journal entry')}</DialogTitle>
          </DialogHeader>
          <JournalEntryForm
            locale={locale}
            accounts={flatAccounts}
            initialValues={{
              entry_date: entry.entry_date?.split('T')[0] || '',
              description: entry.description,
              description_ar: entry.description_ar,
              reference_type: entry.reference_type,
              reference_number: entry.reference_number,
              entry_type: entry.entry_type,
              notes: entry.notes,
              lines: entry.lines.map((line) => ({
                account_id: line.account_id,
                debit_amount: line.debit_amount,
                credit_amount: line.credit_amount,
                description: line.description,
                line_order: line.line_order,
              })),
            }}
            onSubmit={handleUpdate}
            onCancel={() => setShowEditForm(false)}
            submitting={submitting}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}

function SummaryCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
      <div className="text-sm text-gray-500 dark:text-gray-400">{label}</div>
      <div className="mt-2 text-lg font-semibold text-gray-900 dark:text-white">{value}</div>
    </div>
  )
}


