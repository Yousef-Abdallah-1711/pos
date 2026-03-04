'use client'

import { useEffect, useMemo, useState } from 'react'
import { MinusCircle, PlusCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import type { Account, CustomerAccount, Supplier } from '@/types/accounting'
import type {
  CustomerInvoice,
  CustomerInvoiceFormInput,
  SupplierInvoice,
  SupplierInvoiceFormInput,
} from '@/types/invoices'
import { formatAccountLabel, formatCurrency } from './accounting-utils'

type InvoiceFormMode = 'customer' | 'supplier'

interface InvoiceFormProps {
  mode: InvoiceFormMode
  accounts: Account[]
  customers?: CustomerAccount[]
  suppliers?: Supplier[]
  initialValues?: CustomerInvoice | SupplierInvoice | null
  onSubmit: (
    payload: CustomerInvoiceFormInput | SupplierInvoiceFormInput
  ) => Promise<void> | void
  onCancel: () => void
  submitting?: boolean
  locale?: string
}

type InvoiceLineState = {
  description: string
  description_ar: string
  quantity: string
  unit_price: string
  tax_rate: string
  discount_rate: string
  account_id: string
}

type InvoiceState = {
  invoice_date: string
  due_date: string
  customer_id: string
  supplier_id: string
  supplier_ref: string
  order_id: string
  currency: string
  status: string
  notes: string
  terms: string
  internal_notes: string
  lines: InvoiceLineState[]
}

export function InvoiceForm({
  mode,
  accounts,
  customers = [],
  suppliers = [],
  initialValues,
  onSubmit,
  onCancel,
  submitting = false,
  locale = 'en',
}: InvoiceFormProps) {
  const [form, setForm] = useState<InvoiceState>(() =>
    getInitialState(mode, initialValues)
  )
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setForm(getInitialState(mode, initialValues))
    setError(null)
  }, [initialValues, mode])

  const selectableAccounts = useMemo(
    () =>
      accounts.filter((account) => {
        if (account.is_header) {
          return false
        }
        if (mode === 'customer') {
          return account.account_type === 'revenue'
        }
        return account.account_type === 'expense' || account.account_type === 'asset'
      }),
    [accounts, mode]
  )

  const totals = useMemo(() => calculateInvoiceTotals(form.lines), [form.lines])

  const handleChange = <K extends keyof InvoiceState>(key: K, value: InvoiceState[K]) => {
    setForm((current) => ({ ...current, [key]: value }))
  }

  const handleLineChange = <K extends keyof InvoiceLineState>(
    index: number,
    key: K,
    value: InvoiceLineState[K]
  ) => {
    setForm((current) => ({
      ...current,
      lines: current.lines.map((line, lineIndex) =>
        lineIndex === index ? { ...line, [key]: value } : line
      ),
    }))
  }

  const addLine = () => {
    setForm((current) => ({
      ...current,
      lines: [...current.lines, createEmptyLine()],
    }))
  }

  const removeLine = (index: number) => {
    setForm((current) => ({
      ...current,
      lines:
        current.lines.length === 1
          ? current.lines
          : current.lines.filter((_, lineIndex) => lineIndex !== index),
    }))
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (mode === 'customer' && !form.customer_id) {
      setError('Customer is required.')
      return
    }
    if (mode === 'supplier' && !form.supplier_id) {
      setError('Supplier is required.')
      return
    }

    const validLines = form.lines.filter((line) => line.description.trim() !== '')
    if (validLines.length === 0) {
      setError('At least one invoice line is required.')
      return
    }

    setError(null)

    const lines = validLines.map((line, index) => mapLineToPayload(line, index + 1))

    if (mode === 'customer') {
      await onSubmit({
        invoice_date: form.invoice_date,
        due_date: form.due_date,
        customer_id: Number(form.customer_id),
        order_id: form.order_id ? Number(form.order_id) : null,
        currency: form.currency.trim() || 'EGP',
        status: form.status as CustomerInvoiceFormInput['status'],
        notes: form.notes.trim() || undefined,
        terms: form.terms.trim() || undefined,
        internal_notes: form.internal_notes.trim() || undefined,
        lines,
      })
      return
    }

    await onSubmit({
      supplier_ref: form.supplier_ref.trim() || undefined,
      invoice_date: form.invoice_date,
      due_date: form.due_date,
      supplier_id: Number(form.supplier_id),
      currency: form.currency.trim() || 'EGP',
      status: form.status as SupplierInvoiceFormInput['status'],
      notes: form.notes.trim() || undefined,
      internal_notes: form.internal_notes.trim() || undefined,
      lines,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2">
        <Field
          label="Invoice date"
          input={
            <Input
              type="date"
              value={form.invoice_date}
              onChange={(event) => handleChange('invoice_date', event.target.value)}
            />
          }
        />
        <Field
          label="Due date"
          input={
            <Input
              type="date"
              value={form.due_date}
              onChange={(event) => handleChange('due_date', event.target.value)}
            />
          }
        />

        {mode === 'customer' ? (
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

        {mode === 'customer' ? (
          <Field
            label="Order ID (optional)"
            input={
              <Input
                type="number"
                min="1"
                value={form.order_id}
                onChange={(event) => handleChange('order_id', event.target.value)}
                placeholder="Linked order id"
              />
            }
          />
        ) : (
          <Field
            label="Supplier reference"
            input={
              <Input
                value={form.supplier_ref}
                onChange={(event) => handleChange('supplier_ref', event.target.value)}
                placeholder="Vendor invoice reference"
              />
            }
          />
        )}

        <Field
          label="Currency"
          input={
            <Input
              value={form.currency}
              onChange={(event) => handleChange('currency', event.target.value)}
              placeholder="EGP"
            />
          }
        />
        <Field
          label="Status"
          input={
            <select
              value={form.status}
              onChange={(event) => handleChange('status', event.target.value)}
              disabled={mode === 'supplier'}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:bg-gray-100 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:disabled:bg-gray-800"
            >
              {mode === 'customer' ? (
                <>
                  <option value="draft">Draft</option>
                  <option value="sent">Sent</option>
                </>
              ) : (
                <option value="draft">Draft</option>
              )}
            </select>
          }
        />
      </div>

      <div className="rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
        <div className="flex items-center justify-between gap-4 border-b border-gray-200 px-6 py-4 dark:border-gray-700">
          <div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
              Invoice lines
            </h3>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Capture line descriptions, account mapping, and tax/discount percentages. Totals are calculated automatically.
            </p>
          </div>
          <Button type="button" variant="outline" onClick={addLine} className="gap-2">
            <PlusCircle className="h-4 w-4" />
            Add line
          </Button>
        </div>

        <div className="space-y-4 p-6">
          {form.lines.map((line, index) => {
            const lineTotals = calculateLineValues(line)
            return (
              <div
                key={`${index}-${line.description}`}
                className="rounded-2xl border border-gray-200 p-4 dark:border-gray-700"
              >
                <div className="mb-4 flex items-center justify-between gap-3">
                  <div className="text-sm font-medium text-gray-900 dark:text-white">
                    Line {index + 1}
                  </div>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => removeLine(index)}
                    disabled={form.lines.length === 1}
                    className="gap-2"
                  >
                    <MinusCircle className="h-4 w-4" />
                    Remove
                  </Button>
                </div>

                <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
                  <Field
                    label="Description"
                    input={
                      <Input
                        value={line.description}
                        onChange={(event) =>
                          handleLineChange(index, 'description', event.target.value)
                        }
                        placeholder="Consulting services / Raw materials"
                      />
                    }
                  />
                  <Field
                    label="Arabic description"
                    input={
                      <Input
                        value={line.description_ar}
                        onChange={(event) =>
                          handleLineChange(index, 'description_ar', event.target.value)
                        }
                        placeholder="وصف البند"
                      />
                    }
                  />
                  <Field
                    label={mode === 'customer' ? 'Revenue account' : 'Expense / inventory account'}
                    input={
                      <select
                        value={line.account_id}
                        onChange={(event) =>
                          handleLineChange(index, 'account_id', event.target.value)
                        }
                        className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
                      >
                        <option value="">Use system default</option>
                        {selectableAccounts.map((account) => (
                          <option key={account.id} value={account.id}>
                            {formatAccountLabel(account)}
                          </option>
                        ))}
                      </select>
                    }
                  />
                  <Field
                    label="Quantity"
                    input={
                      <Input
                        type="number"
                        min="0"
                        step="0.01"
                        value={line.quantity}
                        onChange={(event) =>
                          handleLineChange(index, 'quantity', event.target.value)
                        }
                      />
                    }
                  />
                  <Field
                    label="Unit price"
                    input={
                      <Input
                        type="number"
                        min="0"
                        step="0.01"
                        value={line.unit_price}
                        onChange={(event) =>
                          handleLineChange(index, 'unit_price', event.target.value)
                        }
                      />
                    }
                  />
                  <Field
                    label="Tax rate %"
                    input={
                      <Input
                        type="number"
                        min="0"
                        step="0.01"
                        value={line.tax_rate}
                        onChange={(event) =>
                          handleLineChange(index, 'tax_rate', event.target.value)
                        }
                      />
                    }
                  />
                  <Field
                    label="Discount rate %"
                    input={
                      <Input
                        type="number"
                        min="0"
                        step="0.01"
                        value={line.discount_rate}
                        onChange={(event) =>
                          handleLineChange(index, 'discount_rate', event.target.value)
                        }
                      />
                    }
                  />
                  <div className="rounded-xl border border-dashed border-gray-200 bg-gray-50 p-4 dark:border-gray-700 dark:bg-gray-900/40">
                    <div className="text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-gray-400">
                      Preview
                    </div>
                    <div className="mt-3 space-y-1 text-sm text-gray-600 dark:text-gray-300">
                      <div>Subtotal: {formatCurrency(lineTotals.subtotal, locale)}</div>
                      <div>Discount: {formatCurrency(lineTotals.discountAmount, locale)}</div>
                      <div>Tax: {formatCurrency(lineTotals.taxAmount, locale)}</div>
                      <div className="font-semibold text-gray-900 dark:text-white">
                        Total: {formatCurrency(lineTotals.lineTotal, locale)}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      </div>

      <div className="grid gap-4 lg:grid-cols-[minmax(0,1fr),320px]">
        <div className="space-y-4">
          {mode === 'customer' ? (
            <Field
              label="Terms"
              input={
                <Textarea
                  value={form.terms}
                  onChange={(event) => handleChange('terms', event.target.value)}
                  placeholder="Payment terms, invoice delivery terms, or commercial notes."
                />
              }
            />
          ) : null}
          <Field
            label="Notes"
            input={
              <Textarea
                value={form.notes}
                onChange={(event) => handleChange('notes', event.target.value)}
                placeholder="Visible notes for the invoice."
              />
            }
          />
          <Field
            label="Internal notes"
            input={
              <Textarea
                value={form.internal_notes}
                onChange={(event) => handleChange('internal_notes', event.target.value)}
                placeholder="Back-office comments, approval notes, or follow-up items."
              />
            }
          />
        </div>

        <div className="rounded-2xl border border-gray-200 bg-white p-6 dark:border-gray-700 dark:bg-gray-800">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            Totals
          </h3>
          <div className="mt-4 space-y-3 text-sm">
            <SummaryRow label="Subtotal" value={formatCurrency(totals.subtotal, locale)} />
            <SummaryRow label="Discount" value={formatCurrency(totals.discountAmount, locale)} />
            <SummaryRow label="Tax" value={formatCurrency(totals.taxAmount, locale)} />
            <SummaryRow
              label="Grand total"
              value={formatCurrency(totals.total, locale)}
              strong
            />
          </div>
        </div>
      </div>

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
            : initialValues
              ? `Update ${mode === 'customer' ? 'customer' : 'supplier'} invoice`
              : `Create ${mode === 'customer' ? 'customer' : 'supplier'} invoice`}
        </Button>
      </div>
    </form>
  )
}

function getInitialState(
  mode: InvoiceFormMode,
  initialValues?: CustomerInvoice | SupplierInvoice | null
): InvoiceState {
  const today = getTodayDate()
  const customerInvoice = mode === 'customer' && isCustomerInvoice(initialValues)
    ? initialValues
    : null
  const supplierInvoice = mode === 'supplier' && isSupplierInvoice(initialValues)
    ? initialValues
    : null
  return {
    invoice_date: initialValues?.invoice_date?.slice(0, 10) || today,
    due_date: initialValues?.due_date?.slice(0, 10) || today,
    customer_id: customerInvoice ? String(customerInvoice.customer_id || '') : '',
    supplier_id: supplierInvoice ? String(supplierInvoice.supplier_id || '') : '',
    supplier_ref: supplierInvoice?.supplier_ref || '',
    order_id: customerInvoice?.order_id ? String(customerInvoice.order_id) : '',
    currency: initialValues?.currency || 'EGP',
    status:
      initialValues?.status || (mode === 'customer' ? 'draft' : 'draft'),
    notes: initialValues?.notes || '',
    terms: customerInvoice?.terms || '',
    internal_notes: initialValues?.internal_notes || '',
    lines:
      initialValues?.lines?.length
        ? initialValues.lines.map((line) => ({
            description: line.description || '',
            description_ar: line.description_ar || '',
            quantity: String(line.quantity ?? 1),
            unit_price: String(line.unit_price ?? 0),
            tax_rate: String(line.tax_rate ?? 0),
            discount_rate: String(line.discount_rate ?? 0),
            account_id: line.account_id ? String(line.account_id) : '',
          }))
        : [createEmptyLine()],
  }
}

function isCustomerInvoice(
  invoice?: CustomerInvoice | SupplierInvoice | null
): invoice is CustomerInvoice {
  return Boolean(invoice && 'customer_id' in invoice)
}

function isSupplierInvoice(
  invoice?: CustomerInvoice | SupplierInvoice | null
): invoice is SupplierInvoice {
  return Boolean(invoice && 'supplier_id' in invoice)
}

function createEmptyLine(): InvoiceLineState {
  return {
    description: '',
    description_ar: '',
    quantity: '1',
    unit_price: '0',
    tax_rate: '0',
    discount_rate: '0',
    account_id: '',
  }
}

function calculateInvoiceTotals(lines: InvoiceLineState[]) {
  return lines.reduce(
    (totals, line) => {
      const values = calculateLineValues(line)
      return {
        subtotal: totals.subtotal + values.subtotal,
        discountAmount: totals.discountAmount + values.discountAmount,
        taxAmount: totals.taxAmount + values.taxAmount,
        total: totals.total + values.lineTotal,
      }
    },
    {
      subtotal: 0,
      discountAmount: 0,
      taxAmount: 0,
      total: 0,
    }
  )
}

function calculateLineValues(line: InvoiceLineState) {
  const quantity = toNumber(line.quantity, 1)
  const unitPrice = toNumber(line.unit_price)
  const taxRate = toNumber(line.tax_rate)
  const discountRate = toNumber(line.discount_rate)
  const subtotal = quantity * unitPrice
  const discountAmount = subtotal * (discountRate / 100)
  const taxableBase = subtotal - discountAmount
  const taxAmount = taxableBase * (taxRate / 100)
  return {
    subtotal,
    discountAmount,
    taxAmount,
    lineTotal: taxableBase + taxAmount,
  }
}

function mapLineToPayload(line: InvoiceLineState, lineOrder: number) {
  const values = calculateLineValues(line)
  return {
    description: line.description.trim(),
    description_ar: line.description_ar.trim() || undefined,
    quantity: toNumber(line.quantity, 1),
    unit_price: toNumber(line.unit_price),
    tax_rate: toNumber(line.tax_rate),
    tax_amount: values.taxAmount,
    discount_rate: toNumber(line.discount_rate),
    discount_amount: values.discountAmount,
    line_total: values.lineTotal,
    account_id: line.account_id ? Number(line.account_id) : null,
    line_order: lineOrder,
  }
}

function toNumber(value: string, fallback = 0) {
  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : fallback
}

function getTodayDate() {
  return new Date().toISOString().slice(0, 10)
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

function SummaryRow({
  label,
  value,
  strong = false,
}: {
  label: string
  value: string
  strong?: boolean
}) {
  return (
    <div className="flex items-center justify-between gap-3">
      <span className="text-gray-500 dark:text-gray-400">{label}</span>
      <span className={strong ? 'font-semibold text-gray-900 dark:text-white' : 'text-gray-700 dark:text-gray-200'}>
        {value}
      </span>
    </div>
  )
}
