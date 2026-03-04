'use client'

import { useEffect, useMemo, useState } from 'react'
import { AlertCircle, Loader2, Lock, Plus } from 'lucide-react'
import Link from 'next/link'
import { toast } from 'sonner'
import { SupplierForm } from '@/components/accounting/SupplierForm'
import { SupplierList } from '@/components/accounting/SupplierList'
import { flattenAccounts } from '@/components/accounting/accounting-utils'
import { PageHeaderSkeleton, TableSkeleton } from '@/components/skeletons'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useUserPermissions } from '@/hooks/useRbac'
import { getAccounts } from '@/lib/api/accountingApi'
import {
  createSupplier,
  deleteSupplier,
  getSuppliers,
  updateSupplier,
} from '@/lib/api/suppliersApi'
import { debugError } from '@/lib/debug'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { usePathname } from 'next/navigation'
import type { Account, Supplier, SupplierFormInput } from '@/types/accounting'

export default function SuppliersPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canView = checkPermission(MODULE_IDS.SUPPLIERS, PERMISSION_LEVELS.READ)
  const canWrite = checkPermission(MODULE_IDS.SUPPLIERS, PERMISSION_LEVELS.WRITE)
  const canDelete = checkPermission(MODULE_IDS.SUPPLIERS, PERMISSION_LEVELS.DELETE)

  const [suppliers, setSuppliers] = useState<Supplier[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [selectedSupplier, setSelectedSupplier] = useState<Supplier | null>(null)

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
      const [supplierData, accountData] = await Promise.all([
        getSuppliers(),
        getAccounts({ include_inactive: false }),
      ])
      setSuppliers(supplierData)
      setAccounts(accountData)
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to load suppliers'
      setApiError(message)
      debugError('[ACCOUNTING] load suppliers failed', error)
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (payload: SupplierFormInput) => {
    setSubmitting(true)

    try {
      if (selectedSupplier) {
        await updateSupplier(selectedSupplier.id, payload)
        toast.success('Supplier updated successfully')
      } else {
        await createSupplier(payload)
        toast.success('Supplier created successfully')
      }

      setSelectedSupplier(null)
      setShowForm(false)
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to save supplier'
      debugError('[ACCOUNTING] save supplier failed', error)
      toast.error(message)
      throw error
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async (supplier: Supplier) => {
    const confirmed = window.confirm(`Delete supplier ${supplier.company_name}?`)
    if (!confirmed) {
      return
    }

    try {
      await deleteSupplier(supplier.id)
      toast.success('Supplier deleted successfully')
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to delete supplier'
      debugError('[ACCOUNTING] delete supplier failed', error)
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
                'accounting.suppliersAccessDenied',
                'You do not have permission to view suppliers.'
              )}
            </p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {loading && suppliers.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <>
          <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                {t('accounting.suppliers', 'Suppliers')}
              </h1>
              <p className="mt-1 text-gray-500 dark:text-gray-400">
                {t(
                  'accounting.suppliersSubtitle',
                  'Maintain supplier master data and payable account links.'
                )}
              </p>
            </div>
            <div className="flex flex-wrap gap-3">
              <Button asChild variant="outline" className="gap-2">
                <Link href={`/${locale}/dashboard/accounting/supplier-invoices`}>
                  {t('accounting.viewSupplierInvoices', 'Supplier invoices')}
                </Link>
              </Button>
              <Button
                onClick={() => {
                  setSelectedSupplier(null)
                  setShowForm(true)
                }}
                disabled={!canWrite}
                className="gap-2"
              >
                <Plus className="h-4 w-4" />
                {t('accounting.newSupplier', 'New supplier')}
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

          <Tabs defaultValue="suppliers" className="space-y-4">
            <TabsList>
              <TabsTrigger value="suppliers">
                {t('accounting.suppliers', 'Suppliers')}
              </TabsTrigger>
              <TabsTrigger value="purchase-invoices">
                {t('navigation.purchaseInvoices', 'Purchase Invoices')}
              </TabsTrigger>
            </TabsList>

            <TabsContent value="suppliers">
              <SupplierList
                suppliers={suppliers}
                locale={locale}
                loading={loading}
                onEdit={(supplier) => {
                  setSelectedSupplier(supplier)
                  setShowForm(true)
                }}
                onDelete={handleDelete}
                canDelete={canDelete}
              />
            </TabsContent>

            <TabsContent value="purchase-invoices">
              <div className="rounded-2xl border border-gray-200 bg-white p-6 dark:border-gray-700 dark:bg-gray-800">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                  {t('accounting.supplierInvoices', 'Supplier Invoices')}
                </h3>
                <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                  {t(
                    'accounting.supplierInvoicesSubtitle',
                    'Track vendor bills, approve payable entries, and record outgoing settlements.'
                  )}
                </p>
                <div className="mt-4">
                  <Button asChild>
                    <Link href={`/${locale}/dashboard/accounting/supplier-invoices`}>
                      {t('accounting.viewSupplierInvoices', 'Supplier invoices')}
                    </Link>
                  </Button>
                </div>
              </div>
            </TabsContent>
          </Tabs>
        </>
      )}

      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-5xl">
          <DialogHeader>
            <DialogTitle>
              {selectedSupplier
                ? t('accounting.editSupplier', 'Edit supplier')
                : t('accounting.createSupplier', 'Create supplier')}
            </DialogTitle>
          </DialogHeader>
          <SupplierForm
            accounts={flatAccounts}
            initialValues={selectedSupplier}
            onSubmit={handleSubmit}
            onCancel={() => setShowForm(false)}
            submitting={submitting}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}


