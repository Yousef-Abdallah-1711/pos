'use client'

import { useEffect, useMemo, useState } from 'react'
import { AlertCircle, Loader2, Lock, Plus, ReceiptText, Wallet } from 'lucide-react'
import { toast } from 'sonner'
import { InvoiceForm } from '@/components/accounting/InvoiceForm'
import { InvoiceList } from '@/components/accounting/InvoiceList'
import { PaymentForm } from '@/components/accounting/PaymentForm'
import { flattenAccounts, formatCurrency, formatDate } from '@/components/accounting/accounting-utils'
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
import {
  createCustomerInvoice,
  createPaymentReceived,
  getCustomerInvoices,
  getOverdueCustomerInvoices,
  getPaymentsReceived,
  sendCustomerInvoice,
  updateCustomerInvoice,
  voidCustomerInvoice,
} from '@/lib/api/accountingApi'
import { getAccounts } from '@/lib/api/accountingApi'
import { getBankAccounts } from '@/lib/api/bankApi'
import { getCustomerAccounts } from '@/lib/api/customersApi'
import { debugError } from '@/lib/debug'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { usePathname, useSearchParams } from 'next/navigation'
import type { Account, CustomerAccount } from '@/types/accounting'
import type {
  BankAccount,
  CustomerInvoice,
  CustomerInvoiceFormInput,
  PaymentReceived,
  PaymentReceivedFormInput,
} from '@/types/invoices'

export default function CustomerInvoicesPage() {
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])
  const initialTab = useMemo(() => {
    const tab = searchParams.get('tab')
    if (tab === 'payments') {
      return 'payments'
    }
    if (tab === 'overdue') {
      return 'overdue'
    }
    return 'invoices'
  }, [searchParams])

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canView = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.READ)
  const canWrite = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.WRITE)

  const [accounts, setAccounts] = useState<Account[]>([])
  const [customers, setCustomers] = useState<CustomerAccount[]>([])
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([])
  const [invoices, setInvoices] = useState<CustomerInvoice[]>([])
  const [overdueInvoices, setOverdueInvoices] = useState<CustomerInvoice[]>([])
  const [payments, setPayments] = useState<PaymentReceived[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [paymentSubmitting, setPaymentSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [showInvoiceForm, setShowInvoiceForm] = useState(false)
  const [showPaymentForm, setShowPaymentForm] = useState(false)
  const [selectedInvoice, setSelectedInvoice] = useState<CustomerInvoice | null>(null)
  const [activeTab, setActiveTab] = useState<'invoices' | 'overdue' | 'payments'>(initialTab)

  useEffect(() => {
    setActiveTab(initialTab)
  }, [initialTab])

  useEffect(() => {
    if (permissionLoading || !canView) {
      return
    }

    void loadData()
  }, [canView, permissionLoading])

  const flatAccounts = useMemo(() => flattenAccounts(accounts), [accounts])

  const openInvoiceCount = useMemo(
    () => invoices.filter((invoice) => invoice.balance_due > 0 && invoice.status !== 'void').length,
    [invoices]
  )
  const receivablesTotal = useMemo(
    () => invoices.reduce((sum, invoice) => sum + (invoice.status === 'void' ? 0 : invoice.balance_due), 0),
    [invoices]
  )

  const loadData = async () => {
    setLoading(true)
    setApiError(null)

    try {
      const [
        invoiceData,
        overdueData,
        paymentData,
        customerData,
        bankData,
        accountData,
      ] = await Promise.all([
        getCustomerInvoices(),
        getOverdueCustomerInvoices(),
        getPaymentsReceived(),
        getCustomerAccounts(),
        getBankAccounts(),
        getAccounts({ include_inactive: false }),
      ])

      setInvoices(invoiceData)
      setOverdueInvoices(overdueData)
      setPayments(paymentData)
      setCustomers(customerData)
      setBankAccounts(bankData)
      setAccounts(accountData)
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to load customer invoices'
      setApiError(message)
      debugError('[ACCOUNTING] load customer invoices failed', error)
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  const handleInvoiceSubmit = async (payload: CustomerInvoiceFormInput) => {
    setSubmitting(true)
    try {
      if (selectedInvoice) {
        await updateCustomerInvoice(selectedInvoice.id, payload)
        toast.success('Customer invoice updated successfully')
      } else {
        await createCustomerInvoice(payload)
        toast.success('Customer invoice created successfully')
      }

      setSelectedInvoice(null)
      setShowInvoiceForm(false)
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to save customer invoice'
      debugError('[ACCOUNTING] save customer invoice failed', error)
      toast.error(message)
      throw error
    } finally {
      setSubmitting(false)
    }
  }

  const handlePaymentSubmit = async (payload: PaymentReceivedFormInput) => {
    setPaymentSubmitting(true)
    try {
      await createPaymentReceived(payload)
      toast.success('Payment received recorded successfully')
      setShowPaymentForm(false)
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to record payment'
      debugError('[ACCOUNTING] record payment received failed', error)
      toast.error(message)
      throw error
    } finally {
      setPaymentSubmitting(false)
    }
  }

  const handleSendInvoice = async (invoice: CustomerInvoice) => {
    try {
      await sendCustomerInvoice(invoice.id)
      toast.success('Invoice marked as sent')
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to send invoice'
      debugError('[ACCOUNTING] send customer invoice failed', error)
      toast.error(message)
    }
  }

  const handleVoidInvoice = async (invoice: CustomerInvoice) => {
    const reason = window.prompt('Reason for voiding this invoice?', 'Voided by user')
    if (reason === null) {
      return
    }

    try {
      await voidCustomerInvoice(invoice.id, reason)
      toast.success('Invoice voided successfully')
      await loadData()
    } catch (error: any) {
      const message =
        error?.response?.data?.error || error?.message || 'Failed to void invoice'
      debugError('[ACCOUNTING] void customer invoice failed', error)
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
      {loading && invoices.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <>
          <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                {t('accounting.customerInvoices', 'Customer Invoices')}
              </h1>
              <p className="mt-1 text-gray-500 dark:text-gray-400">
                {t(
                  'accounting.customerInvoicesSubtitle',
                  'Manage receivable invoices, follow overdue balances, and record incoming settlements.'
                )}
              </p>
            </div>
            <div className="flex flex-wrap gap-3">
              <Button
                onClick={() => setShowPaymentForm(true)}
                disabled={!canWrite}
                variant="outline"
                className="gap-2"
              >
                <Wallet className="h-4 w-4" />
                {t('accounting.recordReceivedPayment', 'Record payment')}
              </Button>
              <Button
                onClick={() => {
                  setSelectedInvoice(null)
                  setShowInvoiceForm(true)
                }}
                disabled={!canWrite}
                className="gap-2"
              >
                <Plus className="h-4 w-4" />
                {t('accounting.newCustomerInvoice', 'New customer invoice')}
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

          <div className="grid gap-4 lg:grid-cols-3">
            <MetricCard
              icon={ReceiptText}
              label={t('accounting.openCustomerInvoices', 'Open invoices')}
              value={String(openInvoiceCount)}
            />
            <MetricCard
              icon={AlertCircle}
              label={t('accounting.overdueCustomerInvoices', 'Overdue invoices')}
              value={String(overdueInvoices.length)}
              tone="danger"
            />
            <MetricCard
              icon={Wallet}
              label={t('accounting.totalReceivables', 'Total receivables')}
              value={formatCurrency(receivablesTotal, locale)}
            />
          </div>

          <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'invoices' | 'overdue' | 'payments')} className="space-y-4">
            <TabsList>
              <TabsTrigger value="invoices">
                {t('accounting.invoiceRegister', 'Invoices')}
              </TabsTrigger>
              <TabsTrigger value="overdue">
                {t('accounting.overdueInvoices', 'Overdue')}
              </TabsTrigger>
              <TabsTrigger value="payments">
                {t('accounting.receivedPayments', 'Payments')}
              </TabsTrigger>
            </TabsList>

            <TabsContent value="invoices">
              <InvoiceList
                mode="customer"
                invoices={invoices}
                locale={locale}
                loading={loading}
                onEdit={(invoice) => {
                  setSelectedInvoice(invoice as CustomerInvoice)
                  setShowInvoiceForm(true)
                }}
                onPrimaryAction={(invoice) => void handleSendInvoice(invoice as CustomerInvoice)}
                onVoid={(invoice) => void handleVoidInvoice(invoice as CustomerInvoice)}
                canWrite={canWrite}
              />
            </TabsContent>

            <TabsContent value="overdue">
              <InvoiceList
                mode="customer"
                invoices={overdueInvoices}
                locale={locale}
                loading={loading}
                onEdit={(invoice) => {
                  setSelectedInvoice(invoice as CustomerInvoice)
                  setShowInvoiceForm(true)
                }}
                onPrimaryAction={(invoice) => void handleSendInvoice(invoice as CustomerInvoice)}
                onVoid={(invoice) => void handleVoidInvoice(invoice as CustomerInvoice)}
                canWrite={canWrite}
              />
            </TabsContent>

            <TabsContent value="payments">
              <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
                {payments.length === 0 ? (
                  <EmptyState
                    title={t('accounting.noReceivedPayments', 'No received payments yet')}
                    description={t(
                      'accounting.noReceivedPaymentsText',
                      'Recorded payments will appear here once receivables begin to settle.'
                    )}
                  />
                ) : (
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
                      <thead className="bg-gray-50 dark:bg-gray-900/40">
                        <tr>
                          <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                            Payment
                          </th>
                          <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                            Customer
                          </th>
                          <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                            Date
                          </th>
                          <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                            Method
                          </th>
                          <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                            Amount
                          </th>
                          <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                            Reference
                          </th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                        {payments.map((payment) => (
                          <tr key={payment.id} className="hover:bg-gray-50 dark:hover:bg-gray-900/20">
                            <td className="px-4 py-4 font-mono text-sm font-semibold text-gray-900 dark:text-white">
                              {payment.payment_number}
                            </td>
                            <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                              {payment.customer_name || '-'}
                            </td>
                            <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                              {formatDate(payment.payment_date, locale)}
                            </td>
                            <td className="px-4 py-4 capitalize text-gray-600 dark:text-gray-300">
                              {payment.payment_method.replace('_', ' ')}
                            </td>
                            <td className="px-4 py-4 font-medium text-gray-900 dark:text-white">
                              {formatCurrency(payment.amount, locale)}
                            </td>
                            <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                              {payment.reference || '-'}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </div>
            </TabsContent>
          </Tabs>
        </>
      )}

      <Dialog open={showInvoiceForm} onOpenChange={setShowInvoiceForm}>
        <DialogContent className="max-w-7xl">
          <DialogHeader>
            <DialogTitle>
              {selectedInvoice
                ? t('accounting.editCustomerInvoice', 'Edit customer invoice')
                : t('accounting.createCustomerInvoice', 'Create customer invoice')}
            </DialogTitle>
          </DialogHeader>
          <InvoiceForm
            mode="customer"
            accounts={flatAccounts}
            customers={customers}
            initialValues={selectedInvoice}
            onSubmit={handleInvoiceSubmit}
            onCancel={() => setShowInvoiceForm(false)}
            submitting={submitting}
            locale={locale}
          />
        </DialogContent>
      </Dialog>

      <Dialog open={showPaymentForm} onOpenChange={setShowPaymentForm}>
        <DialogContent className="max-w-4xl">
          <DialogHeader>
            <DialogTitle>
              {t('accounting.recordReceivedPayment', 'Record payment received')}
            </DialogTitle>
          </DialogHeader>
          <PaymentForm
            mode="received"
            bankAccounts={bankAccounts}
            customers={customers}
            customerInvoices={invoices}
            onSubmit={handlePaymentSubmit}
            onCancel={() => setShowPaymentForm(false)}
            submitting={paymentSubmitting}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}

function MetricCard({
  icon: Icon,
  label,
  value,
  tone = 'default',
}: {
  icon: React.ComponentType<{ className?: string }>
  label: string
  value: string
  tone?: 'default' | 'danger'
}) {
  return (
    <div className={`rounded-2xl border p-5 ${tone === 'danger'
      ? 'border-red-200 bg-red-50 dark:border-red-900 dark:bg-red-950/20'
      : 'border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800'
      }`}>
      <div className="flex items-start justify-between gap-3">
        <div>
          <div className="text-sm text-gray-500 dark:text-gray-400">{label}</div>
          <div className="mt-2 text-2xl font-bold text-gray-900 dark:text-white">{value}</div>
        </div>
        <div className={`rounded-full p-3 ${tone === 'danger'
          ? 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
          : 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300'
          }`}>
          <Icon className="h-5 w-5" />
        </div>
      </div>
    </div>
  )
}

function EmptyState({
  title,
  description,
}: {
  title: string
  description: string
}) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 px-6 py-16 text-center">
      <div className="rounded-full bg-gray-100 p-4 dark:bg-gray-900">
        <Wallet className="h-8 w-8 text-gray-400" />
      </div>
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{title}</h3>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{description}</p>
      </div>
    </div>
  )
}


