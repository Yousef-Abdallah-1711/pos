'use client'

import { useEffect, useMemo, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import type { CustomerAccount, Supplier } from '@/types/accounting'
import type {
  BankAccount,
  CustomerInvoice,
  PaymentMadeFormInput,
  PaymentReceivedFormInput,
  SupplierInvoice,
} from '@/types/invoices'
import { PAYMENT_METHOD_OPTIONS } from './accounting-utils'

type PaymentMode = 'received' | 'made'

interface PaymentFormProps {
  mode: PaymentMode
  bankAccounts: BankAccount[]
  customers?: CustomerAccount[]
  suppliers?: Supplier[]
  customerInvoices?: CustomerInvoice[]
  supplierInvoices?: SupplierInvoice[]
  onSubmit: (
    payload: PaymentReceivedFormInput | PaymentMadeFormInput
  ) => Promise<void> | void
  onCancel: () => void
  submitting?: boolean
}

type PaymentState = {
  payment_date: string
  customer_id: string
  supplier_id: string
  invoice_id: string
  amount: string
  payment_method: string
  reference: string
  bank_account_id: string
  notes: string
}

export function PaymentForm({
  mode,
  bankAccounts,
  customers = [],
  suppliers = [],
  customerInvoices = [],
  supplierInvoices = [],
  onSubmit,
  onCancel,
  submitting = false,
}: PaymentFormProps) {
  const [form, setForm] = useState<PaymentState>(() => getInitialState())
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setForm(getInitialState())
    setError(null)
  }, [mode])

  const openInvoices = useMemo(() => {
    return (mode === 'received' ? customerInvoices : supplierInvoices).filter(
      (invoice) => invoice.status !== 'void' && invoice.balance_due > 0
    )
  }, [customerInvoices, mode, supplierInvoices])

  const handleChange = <K extends keyof PaymentState>(
    key: K,
    value: PaymentState[K]
  ) => {
    setForm((current) => ({ ...current, [key]: value }))
  }

  const handleInvoiceChange = (invoiceID: string) => {
    const invoice = openInvoices.find((item) => item.id === Number(invoiceID))
    if (!invoice) {
      handleChange('invoice_id', invoiceID)
      return
    }

    setForm((current) => ({
      ...current,
      invoice_id: invoiceID,
      amount: current.amount || String(invoice.balance_due),
      customer_id:
        mode === 'received' && 'customer_id' in invoice
          ? String(invoice.customer_id)
          : current.customer_id,
      supplier_id:
        mode === 'made' && 'supplier_id' in invoice
          ? String(invoice.supplier_id)
          : current.supplier_id,
    }))
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (Number(form.amount) <= 0) {
      setError('Amount must be greater than zero.')
      return
    }
    if (!form.payment_method) {
      setError('Payment method is required.')
      return
    }
    if (mode === 'received' && !form.customer_id && !form.invoice_id) {
      setError('Select a customer or an invoice.')
      return
    }
    if (mode === 'made' && !form.supplier_id && !form.invoice_id) {
      setError('Select a supplier or an invoice.')
      return
    }

    setError(null)

    if (mode === 'received') {
      await onSubmit({
        payment_date: form.payment_date,
        customer_id: form.customer_id ? Number(form.customer_id) : 0,
        invoice_id: form.invoice_id ? Number(form.invoice_id) : null,
        amount: Number(form.amount),
        payment_method: form.payment_method,
        reference: form.reference.trim() || undefined,
        bank_account_id: form.bank_account_id ? Number(form.bank_account_id) : null,
        notes: form.notes.trim() || undefined,
      })
      return
    }

    await onSubmit({
      payment_date: form.payment_date,
      supplier_id: form.supplier_id ? Number(form.supplier_id) : 0,
      invoice_id: form.invoice_id ? Number(form.invoice_id) : null,
      amount: Number(form.amount),
      payment_method: form.payment_method,
      reference: form.reference.trim() || undefined,
      bank_account_id: form.bank_account_id ? Number(form.bank_account_id) : null,
      notes: form.notes.trim() || undefined,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2">
        <Field
          label="Payment date"
          input={
            <Input
              type="date"
              value={form.payment_date}
              onChange={(event) => handleChange('payment_date', event.target.value)}
            />
          }
        />
        <Field
          label={mode === 'received' ? 'Customer invoice' : 'Supplier invoice'}
          input={
            <select
              value={form.invoice_id}
              onChange={(event) => handleInvoiceChange(event.target.value)}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              <option value="">No linked invoice</option>
              {openInvoices.map((invoice) => (
                <option key={invoice.id} value={invoice.id}>
                  {invoice.invoice_number} • {mode === 'received' ? (invoice as CustomerInvoice).customer_name : (invoice as SupplierInvoice).supplier_name}
                </option>
              ))}
            </select>
          }
        />

        {mode === 'received' ? (
          <Field
            label="Customer"
            input={
              <select
                value={form.customer_id}
                onChange={(event) => handleChange('customer_id', event.target.value)}
                className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
              >
                <option value="">Select customer</option>
                {customers.map((customer) => (
                  <option key={customer.id} value={customer.id}>
                    {customer.customer_name}
                  </option>
                ))}
              </select>
            }
          />
        ) : (
          <Field
            label="Supplier"
            input={
              <select
                value={form.supplier_id}
                onChange={(event) => handleChange('supplier_id', event.target.value)}
                className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
              >
                <option value="">Select supplier</option>
                {suppliers.map((supplier) => (
                  <option key={supplier.id} value={supplier.id}>
                    {supplier.company_name}
                  </option>
                ))}
              </select>
            }
          />
        )}

        <Field
          label="Amount"
          input={
            <Input
              type="number"
              min="0"
              step="0.01"
              value={form.amount}
              onChange={(event) => handleChange('amount', event.target.value)}
              placeholder="0.00"
            />
          }
        />
        <Field
          label="Payment method"
          input={
            <select
              value={form.payment_method}
              onChange={(event) => handleChange('payment_method', event.target.value)}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              {PAYMENT_METHOD_OPTIONS.map((method) => (
                <option key={method} value={method}>
                  {method}
                </option>
              ))}
            </select>
          }
        />
        <Field
          label="Bank account"
          input={
            <select
              value={form.bank_account_id}
              onChange={(event) => handleChange('bank_account_id', event.target.value)}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              <option value="">System default account</option>
              {bankAccounts
                .filter((account) => account.is_active)
                .map((account) => (
                  <option key={account.id} value={account.id}>
                    {account.account_name} • {account.bank_name}
                  </option>
                ))}
            </select>
          }
        />
        <Field
          label="Reference"
          input={
            <Input
              value={form.reference}
              onChange={(event) => handleChange('reference', event.target.value)}
              placeholder="Transfer reference / cheque number"
            />
          }
        />
      </div>

      <Field
        label="Notes"
        input={
          <Textarea
            value={form.notes}
            onChange={(event) => handleChange('notes', event.target.value)}
            placeholder="Internal notes about this payment."
          />
        }
      />

      {error ? (
        <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/30 dark:text-red-300">
          {error}
        </div>
      ) : null}

      <div className="flex justify-end gap-3 border-t border-gray-200 pt-4 dark:border-gray-700">
        <Button type="button" variant="outline" onClick={onCancel} disabled={submitting}>
          Cancel
        </Button>
        <Button type="submit" disabled={submitting}>
          {submitting
            ? 'Saving...'
            : mode === 'received'
              ? 'Record payment received'
              : 'Record payment made'}
        </Button>
      </div>
    </form>
  )
}

function getInitialState(): PaymentState {
  return {
    payment_date: new Date().toISOString().slice(0, 10),
    customer_id: '',
    supplier_id: '',
    invoice_id: '',
    amount: '',
    payment_method: PAYMENT_METHOD_OPTIONS[0],
    reference: '',
    bank_account_id: '',
    notes: '',
  }
}

function Field({
  label,
  input,
}: {
  label: string
  input: React.ReactNode
}) {
  return (
    <div className="space-y-2">
      <Label>{label}</Label>
      {input}
    </div>
  )
}
