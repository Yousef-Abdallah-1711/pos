'use client'

import { useMemo, useState } from 'react'
import { CheckCircle2, Edit2, FileText, Mail, Search, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import type { CustomerInvoice, SupplierInvoice } from '@/types/invoices'
import {
  formatCurrency,
  formatDate,
  getInvoiceStatusClasses,
} from './accounting-utils'

type InvoiceMode = 'customer' | 'supplier'
type InvoiceRecord = CustomerInvoice | SupplierInvoice

interface InvoiceListProps {
  mode: InvoiceMode
  invoices: InvoiceRecord[]
  locale: string
  loading?: boolean
  onEdit?: (invoice: InvoiceRecord) => void
  onPrimaryAction?: (invoice: InvoiceRecord) => void
  onVoid?: (invoice: InvoiceRecord) => void
  canWrite?: boolean
}

export function InvoiceList({
  mode,
  invoices,
  locale,
  loading = false,
  onEdit,
  onPrimaryAction,
  onVoid,
  canWrite = false,
}: InvoiceListProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState('')

  const filteredInvoices = useMemo(() => {
    return invoices.filter((invoice) => {
      const partyName =
        mode === 'customer'
          ? (invoice as CustomerInvoice).customer_name
          : (invoice as SupplierInvoice).supplier_name

      const searchHaystack = [
        invoice.invoice_number,
        partyName,
        invoice.notes,
        'supplier_ref' in invoice ? invoice.supplier_ref : undefined,
      ]
        .filter(Boolean)
        .join(' ')
        .toLowerCase()

      const matchesSearch =
        searchTerm.trim() === '' ||
        searchHaystack.includes(searchTerm.trim().toLowerCase())

      const matchesStatus = statusFilter === '' || invoice.status === statusFilter

      return matchesSearch && matchesStatus
    })
  }, [invoices, mode, searchTerm, statusFilter])

  const primaryLabel = mode === 'customer' ? 'Send' : 'Approve'
  const PrimaryIcon = mode === 'customer' ? Mail : CheckCircle2

  if (loading) {
    return (
      <div className="rounded-2xl border border-gray-200 bg-white p-6 dark:border-gray-700 dark:bg-gray-800">
        <div className="space-y-3">
          <div className="h-10 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-14 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-14 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-14 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="rounded-2xl border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
        <div className="grid gap-3 lg:grid-cols-[minmax(0,1fr),220px]">
          <div className="relative">
            <Search className="pointer-events-none absolute left-3 top-3.5 h-4 w-4 text-gray-400" />
            <Input
              value={searchTerm}
              onChange={(event) => setSearchTerm(event.target.value)}
              placeholder={`Search ${mode} invoices by number, party, or note`}
              className="pl-9"
            />
          </div>
          <select
            value={statusFilter}
            onChange={(event) => setStatusFilter(event.target.value)}
            className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
          >
            <option value="">All statuses</option>
            <option value="draft">Draft</option>
            <option value="sent">Sent</option>
            <option value="approved">Approved</option>
            <option value="partially_paid">Partially paid</option>
            <option value="paid">Paid</option>
            <option value="overdue">Overdue</option>
            <option value="void">Void</option>
          </select>
        </div>
      </div>

      <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
        {filteredInvoices.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-3 px-6 py-16 text-center">
            <div className="rounded-full bg-gray-100 p-4 dark:bg-gray-900">
              <FileText className="h-8 w-8 text-gray-400" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                No invoices found
              </h3>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                Create an invoice or widen the filters to show more records.
              </p>
            </div>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
              <thead className="bg-gray-50 dark:bg-gray-900/40">
                <tr>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Invoice
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Party
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Invoice date
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Due date
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Total
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Balance due
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Status
                  </th>
                  <th className="px-4 py-3 text-right font-medium text-gray-500 dark:text-gray-400">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                {filteredInvoices.map((invoice) => {
                  const partyName =
                    mode === 'customer'
                      ? (invoice as CustomerInvoice).customer_name
                      : (invoice as SupplierInvoice).supplier_name
                  const invoiceReference =
                    mode === 'supplier' && 'supplier_ref' in invoice
                      ? invoice.supplier_ref
                      : undefined
                  const canTriggerPrimary =
                    mode === 'customer'
                      ? invoice.status === 'draft'
                      : invoice.status === 'draft'
                  const canEditInvoice =
                    mode === 'customer'
                      ? invoice.status !== 'paid' && invoice.status !== 'void'
                      : invoice.status === 'draft'

                  return (
                    <tr
                      key={invoice.id}
                      className={cn(
                        'hover:bg-gray-50 dark:hover:bg-gray-900/20',
                        invoice.status === 'overdue' &&
                          'bg-red-50/70 dark:bg-red-950/10'
                      )}
                    >
                      <td className="px-4 py-4">
                        <div className="font-mono text-sm font-semibold text-gray-900 dark:text-white">
                          {invoice.invoice_number}
                        </div>
                        {invoiceReference ? (
                          <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                            Ref: {invoiceReference}
                          </div>
                        ) : null}
                      </td>
                      <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                        {partyName || '-'}
                      </td>
                      <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                        {formatDate(invoice.invoice_date, locale)}
                      </td>
                      <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                        {formatDate(invoice.due_date, locale)}
                      </td>
                      <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                        {formatCurrency(invoice.total_amount, locale, invoice.currency)}
                      </td>
                      <td className="px-4 py-4 font-medium text-gray-900 dark:text-white">
                        {formatCurrency(invoice.balance_due, locale, invoice.currency)}
                      </td>
                      <td className="px-4 py-4">
                        <span
                          className={cn(
                            'inline-flex rounded-full border px-2.5 py-1 text-xs font-medium capitalize',
                            getInvoiceStatusClasses(invoice.status)
                          )}
                        >
                          {invoice.status.replace('_', ' ')}
                        </span>
                      </td>
                      <td className="px-4 py-4">
                        <div className="flex justify-end gap-2">
                          {canWrite && onEdit && canEditInvoice ? (
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              onClick={() => onEdit(invoice)}
                            >
                              <Edit2 className="mr-2 h-4 w-4" />
                              Edit
                            </Button>
                          ) : null}
                          {canWrite && onPrimaryAction && canTriggerPrimary ? (
                            <Button
                              type="button"
                              size="sm"
                              variant="outline"
                              onClick={() => onPrimaryAction(invoice)}
                              className="gap-2"
                            >
                              <PrimaryIcon className="h-4 w-4" />
                              {primaryLabel}
                            </Button>
                          ) : null}
                          {canWrite && onVoid && invoice.status !== 'void' && invoice.status !== 'paid' ? (
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              onClick={() => onVoid(invoice)}
                              className="text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              Void
                            </Button>
                          ) : null}
                        </div>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
