'use client'

import { useEffect, useMemo, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import type { Account, CustomerAccount, CustomerAccountFormInput } from '@/types/accounting'
import {
  CUSTOMER_PAYMENT_TERMS,
  CUSTOMER_TYPE_OPTIONS,
  formatAccountLabel,
} from './accounting-utils'

interface CustomerFormProps {
  accounts: Account[]
  initialValues?: CustomerAccount | null
  onSubmit: (payload: CustomerAccountFormInput) => Promise<void> | void
  onCancel: () => void
  submitting?: boolean
}

type CustomerState = {
  customer_name: string
  customer_name_ar: string
  customer_type: string
  contact_person: string
  email: string
  phone: string
  mobile: string
  address: string
  address_ar: string
  city: string
  country: string
  tax_number: string
  payment_terms: string
  credit_limit: string
  currency: string
  receivable_account_id: string
  notes: string
  is_active: boolean
}

export function CustomerForm({
  accounts,
  initialValues,
  onSubmit,
  onCancel,
  submitting = false,
}: CustomerFormProps) {
  const [form, setForm] = useState<CustomerState>(() => getInitialState(initialValues))
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setForm(getInitialState(initialValues))
    setError(null)
  }, [initialValues])

  const receivableAccounts = useMemo(
    () => accounts.filter((account) => account.account_type === 'asset'),
    [accounts]
  )

  const handleChange = <K extends keyof CustomerState>(key: K, value: CustomerState[K]) => {
    setForm((current) => ({ ...current, [key]: value }))
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (!form.customer_name.trim()) {
      setError('Customer name is required.')
      return
    }

    setError(null)

    await onSubmit({
      customer_name: form.customer_name.trim(),
      customer_name_ar: form.customer_name_ar.trim() || undefined,
      customer_type: form.customer_type,
      contact_person: form.contact_person.trim() || undefined,
      email: form.email.trim() || undefined,
      phone: form.phone.trim() || undefined,
      mobile: form.mobile.trim() || undefined,
      address: form.address.trim() || undefined,
      address_ar: form.address_ar.trim() || undefined,
      city: form.city.trim() || undefined,
      country: form.country.trim() || 'Egypt',
      tax_number: form.tax_number.trim() || undefined,
      payment_terms: form.payment_terms,
      credit_limit: Number(form.credit_limit) || 0,
      currency: form.currency.trim() || 'EGP',
      receivable_account_id: form.receivable_account_id
        ? Number(form.receivable_account_id)
        : null,
      notes: form.notes.trim() || undefined,
      is_active: form.is_active,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2">
        <Field
          label="Customer name"
          input={
            <Input
              value={form.customer_name}
              onChange={(event) => handleChange('customer_name', event.target.value)}
              placeholder="Corporate Catering Client"
            />
          }
        />
        <Field
          label="Arabic name"
          input={
            <Input
              value={form.customer_name_ar}
              onChange={(event) => handleChange('customer_name_ar', event.target.value)}
              placeholder="عميل تموين الشركات"
            />
          }
        />
        <Field
          label="Customer type"
          input={
            <select
              value={form.customer_type}
              onChange={(event) => handleChange('customer_type', event.target.value)}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              {CUSTOMER_TYPE_OPTIONS.map((type) => (
                <option key={type} value={type}>
                  {type}
                </option>
              ))}
            </select>
          }
        />
        <Field
          label="Payment terms"
          input={
            <select
              value={form.payment_terms}
              onChange={(event) => handleChange('payment_terms', event.target.value)}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              {CUSTOMER_PAYMENT_TERMS.map((term) => (
                <option key={term} value={term}>
                  {term}
                </option>
              ))}
            </select>
          }
        />
        <Field
          label="Contact person"
          input={
            <Input
              value={form.contact_person}
              onChange={(event) => handleChange('contact_person', event.target.value)}
              placeholder="Mona Ali"
            />
          }
        />
        <Field
          label="Email"
          input={
            <Input
              type="email"
              value={form.email}
              onChange={(event) => handleChange('email', event.target.value)}
              placeholder="finance@client.com"
            />
          }
        />
        <Field
          label="Phone"
          input={
            <Input
              value={form.phone}
              onChange={(event) => handleChange('phone', event.target.value)}
              placeholder="+20 2 1234 5678"
            />
          }
        />
        <Field
          label="Mobile"
          input={
            <Input
              value={form.mobile}
              onChange={(event) => handleChange('mobile', event.target.value)}
              placeholder="+20 101 000 0000"
            />
          }
        />
        <Field
          label="Credit limit"
          input={
            <Input
              type="number"
              min="0"
              step="0.01"
              value={form.credit_limit}
              onChange={(event) => handleChange('credit_limit', event.target.value)}
              placeholder="0.00"
            />
          }
        />
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
          label="Receivable account"
          input={
            <select
              value={form.receivable_account_id}
              onChange={(event) => handleChange('receivable_account_id', event.target.value)}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              <option value="">Select receivable account</option>
              {receivableAccounts.map((account) => (
                <option key={account.id} value={account.id}>
                  {formatAccountLabel(account)}
                </option>
              ))}
            </select>
          }
        />
        <Field
          label="Country"
          input={
            <Input
              value={form.country}
              onChange={(event) => handleChange('country', event.target.value)}
              placeholder="Egypt"
            />
          }
        />
        <Field
          label="Tax number"
          input={
            <Input
              value={form.tax_number}
              onChange={(event) => handleChange('tax_number', event.target.value)}
              placeholder="123-456-789"
            />
          }
        />
        <Field
          label="City"
          input={
            <Input
              value={form.city}
              onChange={(event) => handleChange('city', event.target.value)}
              placeholder="Alexandria"
            />
          }
        />
      </div>

      <Field
        label="Address"
        input={
          <Textarea
            value={form.address}
            onChange={(event) => handleChange('address', event.target.value)}
            placeholder="Customer billing address"
          />
        }
      />
      <Field
        label="Arabic address"
        input={
          <Textarea
            value={form.address_ar}
            onChange={(event) => handleChange('address_ar', event.target.value)}
            placeholder="عنوان العميل"
          />
        }
      />
      <Field
        label="Notes"
        input={
          <Textarea
            value={form.notes}
            onChange={(event) => handleChange('notes', event.target.value)}
            placeholder="Credit handling notes or billing instructions."
          />
        }
      />

      <div className="rounded-xl border border-gray-200 p-4 dark:border-gray-700">
        <div className="flex items-start gap-3">
          <Checkbox
            checked={form.is_active}
            onCheckedChange={(value) => handleChange('is_active', value === true)}
            className="mt-1"
          />
          <div>
            <div className="font-medium text-gray-900 dark:text-white">Active customer</div>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Disabled customer accounts remain visible for reporting but cannot be reused.
            </p>
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
          {submitting ? 'Saving...' : initialValues ? 'Update customer' : 'Create customer'}
        </Button>
      </div>
    </form>
  )
}

function getInitialState(initialValues?: CustomerAccount | null): CustomerState {
  return {
    customer_name: initialValues?.customer_name || '',
    customer_name_ar: initialValues?.customer_name_ar || '',
    customer_type: initialValues?.customer_type || 'individual',
    contact_person: initialValues?.contact_person || '',
    email: initialValues?.email || '',
    phone: initialValues?.phone || '',
    mobile: initialValues?.mobile || '',
    address: initialValues?.address || '',
    address_ar: initialValues?.address_ar || '',
    city: initialValues?.city || '',
    country: initialValues?.country || 'Egypt',
    tax_number: initialValues?.tax_number || '',
    payment_terms: initialValues?.payment_terms || 'immediate',
    credit_limit: String(initialValues?.credit_limit ?? 0),
    currency: initialValues?.currency || 'EGP',
    receivable_account_id: initialValues?.receivable_account_id
      ? String(initialValues.receivable_account_id)
      : '',
    notes: initialValues?.notes || '',
    is_active: initialValues?.is_active ?? true,
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
