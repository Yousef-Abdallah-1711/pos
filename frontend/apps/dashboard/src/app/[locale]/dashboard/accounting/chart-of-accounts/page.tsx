'use client'

import { useEffect, useMemo, useState } from 'react'
import { AlertCircle, Loader2, Lock, Plus } from 'lucide-react'
import { toast } from 'sonner'
import { AccountForm } from '@/components/accounting/AccountForm'
import { AccountTree } from '@/components/accounting/AccountTree'
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
  createAccount,
  deleteAccount,
  getAccounts,
  updateAccount,
} from '@/lib/api/accountingApi'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'
import { usePathname } from 'next/navigation'
import type { Account, AccountFormInput } from '@/types/accounting'

export default function ChartOfAccountsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canView = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.READ)
  const canWrite = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.WRITE)
  const canDelete = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.DELETE)

  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [selectedAccount, setSelectedAccount] = useState<Account | null>(null)
  const [searchTerm, setSearchTerm] = useState('')
  const [typeFilter, setTypeFilter] = useState('')

  useEffect(() => {
    if (permissionLoading || !canView) {
      return
    }

    void loadAccounts()
  }, [canView, permissionLoading])

  const flatAccounts = useMemo(() => flattenAccounts(accounts), [accounts])

  const loadAccounts = async () => {
    setLoading(true)
    setApiError(null)

    try {
      const data = await getAccounts({ tree: true, include_inactive: true })
      setAccounts(data)
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to load accounts'
      setApiError(message)
      debugError('[ACCOUNTING] load accounts failed', error)
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  const openCreateDialog = () => {
    setSelectedAccount(null)
    setShowForm(true)
  }

  const openEditDialog = (account: Account) => {
    setSelectedAccount(account)
    setShowForm(true)
  }

  const handleSubmit = async (payload: AccountFormInput) => {
    setSubmitting(true)

    try {
      if (selectedAccount) {
        await updateAccount(selectedAccount.id, payload)
        toast.success('Account updated successfully')
      } else {
        await createAccount(payload)
        toast.success('Account created successfully')
      }

      setShowForm(false)
      setSelectedAccount(null)
      await loadAccounts()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to save account'
      debugError('[ACCOUNTING] save account failed', error)
      toast.error(message)
      throw error
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async (account: Account) => {
    const confirmed = window.confirm(
      `Delete account ${account.account_code} - ${account.account_name}?`
    )
    if (!confirmed) {
      return
    }

    try {
      await deleteAccount(account.id)
      toast.success('Account deleted successfully')
      await loadAccounts()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to delete account'
      debugError('[ACCOUNTING] delete account failed', error)
      toast.error(message)
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
      {loading && accounts.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <>
          <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                {t('accounting.chartOfAccounts', 'Chart of Accounts')}
              </h1>
              <p className="mt-1 text-gray-500 dark:text-gray-400">
                {t(
                  'accounting.chartSubtitle',
                  'Manage your account hierarchy, balances, and posting structure.'
                )}
              </p>
            </div>
            <Button onClick={openCreateDialog} disabled={!canWrite} className="gap-2">
              <Plus className="h-4 w-4" />
              {t('accounting.newAccount', 'New account')}
            </Button>
          </div>

          {apiError ? (
            <div className="rounded-xl border border-red-200 bg-red-50 p-4 dark:border-red-900 dark:bg-red-950/30">
              <div className="flex items-start gap-3">
                <AlertCircle className="mt-0.5 h-5 w-5 text-red-600 dark:text-red-400" />
                <div>
                  <h3 className="font-semibold text-red-900 dark:text-red-300">
                    {t('accounting.loadErrorTitle', 'Unable to load accounts')}
                  </h3>
                  <p className="mt-1 text-sm text-red-700 dark:text-red-300">{apiError}</p>
                  <button
                    type="button"
                    onClick={() => void loadAccounts()}
                    className="mt-2 text-sm font-medium text-red-700 underline dark:text-red-300"
                  >
                    {t('accounting.retry', 'Retry')}
                  </button>
                </div>
              </div>
            </div>
          ) : null}

          <AccountTree
            accounts={accounts}
            loading={loading}
            locale={locale}
            searchTerm={searchTerm}
            typeFilter={typeFilter}
            onSearchChange={setSearchTerm}
            onTypeFilterChange={setTypeFilter}
            onEdit={openEditDialog}
            onDelete={handleDelete}
            canWrite={canWrite}
            canDelete={canDelete}
          />
        </>
      )}

      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-4xl">
          <DialogHeader>
            <DialogTitle>
              {selectedAccount
                ? t('accounting.editAccount', 'Edit account')
                : t('accounting.createAccount', 'Create account')}
            </DialogTitle>
          </DialogHeader>
          <AccountForm
            accounts={flatAccounts}
            initialValues={selectedAccount}
            onSubmit={handleSubmit}
            onCancel={() => setShowForm(false)}
            submitting={submitting}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}


