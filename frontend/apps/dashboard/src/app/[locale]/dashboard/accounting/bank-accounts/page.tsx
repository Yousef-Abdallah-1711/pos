'use client'

import { useEffect, useMemo, useState } from 'react'
import { AlertCircle, Landmark, Loader2, Lock, Pencil, Plus, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import { BankAccountForm } from '@/components/accounting/BankAccountForm'
import { BankTransactionList } from '@/components/accounting/BankTransactionList'
import { flattenAccounts, formatCurrency, formatDate } from '@/components/accounting/accounting-utils'
import { PageHeaderSkeleton, TableSkeleton } from '@/components/skeletons'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useUserPermissions } from '@/hooks/useRbac'
import { getAccounts } from '@/lib/api/accountingApi'
import {
  createBankAccount,
  deleteBankAccount,
  getBankAccounts,
  getBankTransactions,
  updateBankAccount,
} from '@/lib/api/bankApi'
import { debugError } from '@/lib/debug'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { usePathname } from 'next/navigation'
import type { Account } from '@/types/accounting'
import type {
  BankAccount,
  BankAccountFormInput,
  BankTransaction,
} from '@/types/invoices'

export default function BankAccountsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canView = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.READ)
  const canWrite = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.WRITE)
  const canDelete = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.DELETE)

  const [accounts, setAccounts] = useState<Account[]>([])
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([])
  const [transactions, setTransactions] = useState<BankTransaction[]>([])
  const [loading, setLoading] = useState(false)
  const [transactionsLoading, setTransactionsLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [selectedBankAccountID, setSelectedBankAccountID] = useState<number | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [selectedBankAccount, setSelectedBankAccount] = useState<BankAccount | null>(null)

  useEffect(() => {
    if (permissionLoading || !canView) {
      return
    }

    void loadData()
  }, [canView, permissionLoading])

  useEffect(() => {
    if (!selectedBankAccountID || !canView) {
      setTransactions([])
      return
    }

    void loadTransactions(selectedBankAccountID)
  }, [canView, selectedBankAccountID])

  const flatAccounts = useMemo(() => flattenAccounts(accounts), [accounts])
  const activeSelection = useMemo(
    () =>
      bankAccounts.find((account) => account.id === selectedBankAccountID) || null,
    [bankAccounts, selectedBankAccountID]
  )

  const loadData = async () => {
    setLoading(true)
    setApiError(null)

    try {
      const [bankData, accountData] = await Promise.all([
        getBankAccounts(),
        getAccounts({ include_inactive: false }),
      ])
      setBankAccounts(bankData)
      setAccounts(accountData)
      setSelectedBankAccountID((current) => current ?? bankData[0]?.id ?? null)
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || t('bankAccountMessages.loadFailed', 'Failed to load bank accounts')
      setApiError(message)
      debugError('[ACCOUNTING] load bank accounts failed', error)
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  const loadTransactions = async (bankAccountID: number) => {
    setTransactionsLoading(true)
    try {
      const data = await getBankTransactions(bankAccountID)
      setTransactions(data)
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || t('bankAccountMessages.loadTransactionsFailed', 'Failed to load bank transactions')
      debugError('[ACCOUNTING] load bank transactions failed', error)
      toast.error(message)
      setTransactions([])
    } finally {
      setTransactionsLoading(false)
    }
  }

  const handleSubmit = async (payload: BankAccountFormInput) => {
    setSubmitting(true)
    try {
      if (selectedBankAccount) {
        await updateBankAccount(selectedBankAccount.id, payload)
        toast.success(t('bankAccountMessages.updated', 'Bank account updated successfully'))
      } else {
        const created = await createBankAccount(payload)
        toast.success(t('bankAccountMessages.created', 'Bank account created successfully'))
        setSelectedBankAccountID(created.id)
      }

      setSelectedBankAccount(null)
      setShowForm(false)
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || t('bankAccountMessages.saveFailed', 'Failed to save bank account')
      debugError('[ACCOUNTING] save bank account failed', error)
      toast.error(message)
      throw error
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async (bankAccount: BankAccount) => {
    if (!window.confirm(t('bankAccountMessages.confirmDelete', 'Delete bank account {name}?').replace('{name}', bankAccount.account_name))) {
      return
    }

    try {
      await deleteBankAccount(bankAccount.id)
      toast.success(t('bankAccountMessages.deleted', 'Bank account deleted successfully'))
      if (selectedBankAccountID === bankAccount.id) {
        setSelectedBankAccountID(null)
      }
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || t('bankAccountMessages.deleteFailed', 'Failed to delete bank account')
      debugError('[ACCOUNTING] delete bank account failed', error)
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
      {loading && bankAccounts.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={6} />
        </>
      ) : (
        <>
          <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                {t('accounting.bankAccounts', 'Bank Accounts')}
              </h1>
              <p className="mt-1 text-gray-500 dark:text-gray-400">
                {t(
                  'accounting.bankAccountsSubtitle',
                  'Manage settlement accounts, link them to the ledger, and review posted bank activity.'
                )}
              </p>
            </div>
            <Button
              onClick={() => {
                setSelectedBankAccount(null)
                setShowForm(true)
              }}
              disabled={!canWrite}
              className="gap-2"
            >
              <Plus className="h-4 w-4" />
              {t('accounting.newBankAccount', 'New bank account')}
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

          <div className="grid gap-6 xl:grid-cols-[380px,minmax(0,1fr)]">
            <div className="space-y-4">
              {bankAccounts.length === 0 ? (
                <div className="rounded-2xl border border-gray-200 bg-white p-8 text-center dark:border-gray-700 dark:bg-gray-800">
                  <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-900">
                    <Landmark className="h-6 w-6 text-gray-400" />
                  </div>
                  <h3 className="mt-4 text-lg font-semibold text-gray-900 dark:text-white">
                    {t('accounting.noBankAccounts', 'No bank accounts yet')}
                  </h3>
                  <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                    {t(
                      'accounting.noBankAccountsText',
                      'Create your first bank account to start routing payments and transfers.'
                    )}
                  </p>
                </div>
              ) : (
                bankAccounts.map((bankAccount) => {
                  const isSelected = bankAccount.id === selectedBankAccountID
                  return (
                    <button
                      key={bankAccount.id}
                      type="button"
                      onClick={() => setSelectedBankAccountID(bankAccount.id)}
                      className={`w-full rounded-2xl border p-5 text-left transition ${isSelected
                        ? 'border-primary-500 bg-primary-50 shadow-sm dark:border-primary-400 dark:bg-primary-950/20'
                        : 'border-gray-200 bg-white hover:border-primary-200 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:hover:border-primary-900 dark:hover:bg-gray-900/40'
                        }`}
                    >
                      <div className="flex items-start justify-between gap-3">
                        <div>
                          <div className="text-lg font-semibold text-gray-900 dark:text-white">
                            {bankAccount.account_name}
                          </div>
                          <div className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                            {bankAccount.bank_name}
                          </div>
                        </div>
                        <span className="rounded-full border border-gray-200 bg-white px-2.5 py-1 text-xs font-medium capitalize text-gray-600 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300">
                          {bankAccount.account_type.replace('_', ' ')}
                        </span>
                      </div>

                      <div className="mt-4 grid gap-3 text-sm text-gray-600 dark:text-gray-300">
                        <div className="flex items-center justify-between gap-3">
                          <span>{t('balance', 'Balance')}</span>
                          <span className="font-semibold text-gray-900 dark:text-white">
                            {formatCurrency(
                              bankAccount.current_balance,
                              locale,
                              bankAccount.currency
                            )}
                          </span>
                        </div>
                        <div className="flex items-center justify-between gap-3">
                          <span>{t('opened', 'Opened')}</span>
                          <span>{formatDate(bankAccount.opening_date, locale)}</span>
                        </div>
                        <div className="flex items-center justify-between gap-3">
                          <span>{t('primary', 'Primary')}</span>
                          <span>{bankAccount.is_primary ? t('yes', 'Yes') : t('no', 'No')}</span>
                        </div>
                      </div>

                      <div className="mt-4 flex items-center justify-end gap-2">
                        {canWrite ? (
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={(event) => {
                              event.stopPropagation()
                              setSelectedBankAccount(bankAccount)
                              setShowForm(true)
                            }}
                          >
                            <Pencil className="mr-2 h-4 w-4" />
                            {t('edit', 'Edit')}
                          </Button>
                        ) : null}
                        {canDelete ? (
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={(event) => {
                              event.stopPropagation()
                              void handleDelete(bankAccount)
                            }}
                            className="text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
                          >
                            <Trash2 className="mr-2 h-4 w-4" />
                            {t('delete', 'Delete')}
                          </Button>
                        ) : null}
                      </div>
                    </button>
                  )
                })
              )}
            </div>

            <BankTransactionList
              account={activeSelection}
              transactions={transactions}
              locale={locale}
              loading={transactionsLoading}
            />
          </div>
        </>
      )}

      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-5xl">
          <DialogHeader>
            <DialogTitle>
              {selectedBankAccount
                ? t('accounting.editBankAccount', 'Edit bank account')
                : t('accounting.createBankAccount', 'Create bank account')}
            </DialogTitle>
          </DialogHeader>
          <BankAccountForm
            accounts={flatAccounts}
            initialValues={selectedBankAccount}
            onSubmit={handleSubmit}
            onCancel={() => setShowForm(false)}
            submitting={submitting}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}


