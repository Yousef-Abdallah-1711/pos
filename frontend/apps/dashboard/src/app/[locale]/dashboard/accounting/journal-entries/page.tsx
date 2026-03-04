'use client'

import { useEffect, useMemo, useState } from 'react'
import { AlertCircle, Loader2, Lock, Plus } from 'lucide-react'
import { toast } from 'sonner'
import { flattenAccounts } from '@/components/accounting/accounting-utils'
import { PageHeaderSkeleton, TableSkeleton } from '@/components/skeletons'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useUserPermissions } from '@/hooks/useRbac'
import {
  createJournalEntry,
  getAccounts,
  getJournalEntries,
} from '@/lib/api/accountingApi'
import { debugError } from '@/lib/debug'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { usePathname } from 'next/navigation'
import dynamic from 'next/dynamic'
import type { Account, JournalEntry, JournalEntryFilters, JournalEntryFormInput } from '@/types/accounting'

const JournalEntryForm = dynamic(
  () => import('@/components/accounting/JournalEntryForm').then((mod) => mod.JournalEntryForm),
  { ssr: false, loading: () => <TableSkeleton rows={4} /> }
)
const JournalEntryList = dynamic(
  () => import('@/components/accounting/JournalEntryList').then((mod) => mod.JournalEntryList),
  { ssr: false, loading: () => <TableSkeleton rows={8} /> }
)

export default function JournalEntriesPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canView = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.READ)
  const canWrite = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.WRITE)

  const [entries, setEntries] = useState<JournalEntry[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])
  const [filters, setFilters] = useState<JournalEntryFilters>({})
  const [loading, setLoading] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (permissionLoading || !canView) {
      return
    }

    void loadData()
  }, [canView, filters.end_date, filters.entry_type, filters.start_date, filters.status, permissionLoading])

  const flatAccounts = useMemo(() => flattenAccounts(accounts), [accounts])

  const loadData = async () => {
    setLoading(true)
    setApiError(null)

    try {
      const [entryData, accountData] = await Promise.all([
        getJournalEntries(filters),
        getAccounts({ include_inactive: false }),
      ])
      setEntries(entryData)
      setAccounts(accountData)
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to load journal entries'
      setApiError(message)
      debugError('[ACCOUNTING] load journal entries failed', error)
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (payload: JournalEntryFormInput) => {
    setSubmitting(true)

    try {
      await createJournalEntry(payload)
      toast.success('Journal entry created successfully')
      setShowForm(false)
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to create journal entry'
      debugError('[ACCOUNTING] create journal entry failed', error)
      toast.error(message)
      throw error
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

  return (
    <div className="space-y-6">
      {loading && entries.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <>
          <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                {t('accounting.journalEntries', 'Journal Entries')}
              </h1>
              <p className="mt-1 text-gray-500 dark:text-gray-400">
                {t(
                  'accounting.journalSubtitle',
                  'Create, review, and post manual accounting entries.'
                )}
              </p>
            </div>
            <Button onClick={() => setShowForm(true)} disabled={!canWrite} className="gap-2">
              <Plus className="h-4 w-4" />
              {t('accounting.newJournalEntry', 'New journal entry')}
            </Button>
          </div>

          {apiError ? (
            <div className="rounded-xl border border-red-200 bg-red-50 p-4 dark:border-red-900 dark:bg-red-950/30">
              <div className="flex items-start gap-3">
                <AlertCircle className="mt-0.5 h-5 w-5 text-red-600 dark:text-red-400" />
                <div>
                  <h3 className="font-semibold text-red-900 dark:text-red-300">
                    {t('accounting.loadErrorTitle', 'Unable to load records')}
                  </h3>
                  <p className="mt-1 text-sm text-red-700 dark:text-red-300">{apiError}</p>
                  <button
                    type="button"
                    onClick={() => void loadData()}
                    className="mt-2 text-sm font-medium text-red-700 underline dark:text-red-300"
                  >
                    {t('accounting.retry', 'Retry')}
                  </button>
                </div>
              </div>
            </div>
          ) : null}

          <JournalEntryList
            entries={entries}
            filters={filters}
            locale={locale}
            loading={loading}
            onFiltersChange={setFilters}
          />
        </>
      )}

      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-6xl">
          <DialogHeader>
            <DialogTitle>
              {t('accounting.createJournalEntry', 'Create journal entry')}
            </DialogTitle>
          </DialogHeader>
          <JournalEntryForm
            locale={locale}
            accounts={flatAccounts}
            onSubmit={handleCreate}
            onCancel={() => setShowForm(false)}
            submitting={submitting}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}


