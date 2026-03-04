'use client'

import { useEffect, useMemo, useState } from 'react'
import { AlertCircle, Loader2, Lock, Plus } from 'lucide-react'
import Link from 'next/link'
import { toast } from 'sonner'
import { CustomerForm } from '@/components/accounting/CustomerForm'
import { CustomerList } from '@/components/accounting/CustomerList'
import {
  flattenAccounts,
  formatCurrency,
  formatDate,
} from '@/components/accounting/accounting-utils'
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
  createCustomerAccount,
  deleteCustomerAccount,
  getCustomerAccounts,
  getCustomerStatement,
  updateCustomerAccount,
} from '@/lib/api/customersApi'
import { debugError } from '@/lib/debug'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { usePathname } from 'next/navigation'
import type {
  Account,
  CustomerAccount,
  CustomerAccountFormInput,
  CustomerStatementLine,
} from '@/types/accounting'

export default function CustomersPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canView = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.READ)
  const canWrite = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.WRITE)
  const canDelete = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.DELETE)

  const [customers, setCustomers] = useState<CustomerAccount[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [selectedCustomer, setSelectedCustomer] = useState<CustomerAccount | null>(null)
  const [statementCustomer, setStatementCustomer] = useState<CustomerAccount | null>(null)
  const [statementLines, setStatementLines] = useState<CustomerStatementLine[]>([])
  const [statementLoading, setStatementLoading] = useState(false)

  useEffect(() => {
    if (permissionLoading || !canView) {
      return
    }

    void loadData()
  }, [canView, permissionLoading])

  const flatAccounts = useMemo(() => flattenAccounts(accounts), [accounts])

  const loadData = async () => {
    setLoading(true)
    setApiError(null)

    try {
      const [customerData, accountData] = await Promise.all([
        getCustomerAccounts(),
        getAccounts({ include_inactive: false }),
      ])
      setCustomers(customerData)
      setAccounts(accountData)
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to load customers'
      setApiError(message)
      debugError('[ACCOUNTING] load customers failed', error)
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (payload: CustomerAccountFormInput) => {
    setSubmitting(true)

    try {
      if (selectedCustomer) {
        await updateCustomerAccount(selectedCustomer.id, payload)
        toast.success('Customer updated successfully')
      } else {
        await createCustomerAccount(payload)
        toast.success('Customer created successfully')
      }

      setSelectedCustomer(null)
      setShowForm(false)
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to save customer'
      debugError('[ACCOUNTING] save customer failed', error)
      toast.error(message)
      throw error
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async (customer: CustomerAccount) => {
    const confirmed = window.confirm(`Delete customer ${customer.customer_name}?`)
    if (!confirmed) {
      return
    }

    try {
      await deleteCustomerAccount(customer.id)
      toast.success('Customer deleted successfully')
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to delete customer'
      debugError('[ACCOUNTING] delete customer failed', error)
      toast.error(message)
    }
  }

  const handleViewStatement = async (customer: CustomerAccount) => {
    setStatementCustomer(customer)
    setStatementLoading(true)
    setStatementLines([])

    try {
      const data = await getCustomerStatement(customer.id)
      setStatementLines(data)
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to load statement'
      debugError('[ACCOUNTING] load customer statement failed', error)
      toast.error(message)
    } finally {
      setStatementLoading(false)
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
                'accounting.customersAccessDenied',
                'You do not have permission to view customer accounts.'
              )}
            </p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {loading && customers.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <>
          <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                {t('accounting.customers', 'Customer Accounts')}
              </h1>
              <p className="mt-1 text-gray-500 dark:text-gray-400">
                {t(
                  'accounting.customersSubtitle',
                  'Track receivable customers, credit terms, and account statements.'
                )}
              </p>
            </div>
            <div className="flex flex-wrap gap-3">
              <Button
                asChild
                variant="outline"
                className="gap-2"
              >
                <Link href={`/${locale}/dashboard/accounting/customer-invoices`}>
                  {t('accounting.viewCustomerInvoices', 'Customer invoices')}
                </Link>
              </Button>
              <Button
                onClick={() => {
                  setSelectedCustomer(null)
                  setShowForm(true)
                }}
                disabled={!canWrite}
                className="gap-2"
              >
                <Plus className="h-4 w-4" />
                {t('accounting.newCustomer', 'New customer')}
              </Button>
            </div>
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

          <CustomerList
            customers={customers}
            locale={locale}
            loading={loading}
            onEdit={(customer) => {
              setSelectedCustomer(customer)
              setShowForm(true)
            }}
            onDelete={handleDelete}
            onViewStatement={(customer) => void handleViewStatement(customer)}
            canDelete={canDelete}
          />
        </>
      )}

      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-5xl">
          <DialogHeader>
            <DialogTitle>
              {selectedCustomer
                ? t('accounting.editCustomer', 'Edit customer')
                : t('accounting.createCustomer', 'Create customer')}
            </DialogTitle>
          </DialogHeader>
          <CustomerForm
            accounts={flatAccounts}
            initialValues={selectedCustomer}
            onSubmit={handleSubmit}
            onCancel={() => setShowForm(false)}
            submitting={submitting}
          />
        </DialogContent>
      </Dialog>

      <Dialog
        open={statementCustomer !== null}
        onOpenChange={(open) => {
          if (!open) {
            setStatementCustomer(null)
            setStatementLines([])
          }
        }}
      >
        <DialogContent className="max-w-5xl">
          <DialogHeader>
            <DialogTitle>
              {statementCustomer
                ? `${t('accounting.customerStatement', 'Customer statement')}: ${statementCustomer.customer_name}`
                : t('accounting.customerStatement', 'Customer statement')}
            </DialogTitle>
          </DialogHeader>

          {statementLoading ? (
            <div className="flex h-40 items-center justify-center">
              <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
            </div>
          ) : statementLines.length === 0 ? (
            <div className="rounded-xl border border-gray-200 bg-gray-50 p-6 text-sm text-gray-500 dark:border-gray-700 dark:bg-gray-900/40 dark:text-gray-400">
              {t(
                'accounting.noStatementLines',
                'No statement lines are available for this customer yet.'
              )}
            </div>
          ) : (
            <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
                  <thead className="bg-gray-50 dark:bg-gray-900/40">
                    <tr>
                      <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                        Date
                      </th>
                      <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                        Entry
                      </th>
                      <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                        Reference
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
                      <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                        Running balance
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                    {statementLines.map((line, index) => (
                      <tr key={`${line.entry_number}-${index}`}>
                        <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                          {formatDate(line.entry_date, locale)}
                        </td>
                        <td className="px-4 py-4 font-mono text-gray-900 dark:text-white">
                          {line.entry_number}
                        </td>
                        <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                          {line.reference_number || '-'}
                        </td>
                        <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                          {line.description}
                        </td>
                        <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                          {formatCurrency(line.debit_amount, locale)}
                        </td>
                        <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                          {formatCurrency(line.credit_amount, locale)}
                        </td>
                        <td className="px-4 py-4 font-semibold text-gray-900 dark:text-white">
                          {formatCurrency(line.running_balance, locale)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}


