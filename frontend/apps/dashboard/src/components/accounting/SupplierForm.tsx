'use client'

import { useEffect, useMemo, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import type { Account, Supplier, SupplierFormInput } from '@/types/accounting'
import {
  SUPPLIER_CATEGORY_OPTIONS,
  SUPPLIER_PAYMENT_TERMS,
  formatAccountLabel,
} from './accounting-utils'

interface SupplierFormProps {
  accounts: Account[]
  initialValues?: Supplier | null
  onSubmit: (payload: SupplierFormInput) => Promise<void> | void
  onCancel: () => void
  submitting?: boolean
}

type SupplierState = {
  company_name: string
  company_name_ar: string
  contact_person: string
  email: string
  phone: string
  mobile: string
  website: string
  address: string
  address_ar: string
  city: string
  state: string
  postal_code: string
  country: string
  payment_terms: string
  credit_limit: string
  currency: string
  payable_account_id: string
  category: string
  bank_name: string
  bank_account: string
  iban: string
  tax_number: string
  notes: string
  is_active: boolean
}

export function SupplierForm({
  accounts,
  initialValues,
  onSubmit,
  onCancel,
  submitting = false,
}: SupplierFormProps) {
  const [form, setForm] = useState<SupplierState>(() => getInitialState(initialValues))
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setForm(getInitialState(initialValues))
    setError(null)
  }, [initialValues])

  const payableAccounts = useMemo(
    () => accounts.filter((account) => account.account_type === 'liability'),
    [accounts]
  )

  const handleChange = <K extends keyof SupplierState>(key: K, value: SupplierState[K]) => {
    setForm((current) => ({ ...current, [key]: value }))
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (!form.company_name.trim()) {
      setError('Company name is required.')
      return
    }

    setError(null)

    await onSubmit({
      company_name: form.company_name.trim(),
      company_name_ar: form.company_name_ar.trim() || undefined,
      contact_person: form.contact_person.trim() || undefined,
      email: form.email.trim() || undefined,
      phone: form.phone.trim() || undefined,
      mobile: form.mobile.trim() || undefined,
      website: form.website.trim() || undefined,
      address: form.address.trim() || undefined,
      address_ar: form.address_ar.trim() || undefined,
      city: form.city.trim() || undefined,
      state: form.state.trim() || undefined,
      postal_code: form.postal_code.trim() || undefined,
      country: form.country.trim() || 'Egypt',
      tax_number: form.tax_number.trim() || undefined,
      payment_terms: form.payment_terms,
      credit_limit: Number(form.credit_limit) || 0,
      currency: form.currency.trim() || 'EGP',
      payable_account_id: form.payable_account_id ? Number(form.payable_account_id) : null,
      category: form.category || undefined,
      bank_name: form.bank_name.trim() || undefined,
      bank_account: form.bank_account.trim() || undefined,
      iban: form.iban.trim() || undefined,
      notes: form.notes.trim() || undefined,
      is_active: form.is_active,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2">
        <Field
          label="Company name"
          input={
            <Input
              value={form.company_name}
              onChange={(event) => handleChange('company_name', event.target.value)}
              placeholder="Fresh Produce Co."
            />
          }
        />
        <Field
          label="Arabic name"
          input={
            <Input
              value={form.company_name_ar}
              onChange={(event) => handleChange('company_name_ar', event.target.value)}
              placeholder="شركة الخضروات الطازجة"
            />
          }
        />
        <Field
          label="Contact person"
          input={
            <Input
              value={form.contact_person}
              onChange={(event) => handleChange('contact_person', event.target.value)}
              placeholder="Ahmed Hassan"
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
              placeholder="purchasing@supplier.com"
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
              placeholder="+20 100 000 0000"
            />
          }
        />
        <Field
          label="Website"
          input={
            <Input
              value={form.website}
              onChange={(event) => handleChange('website', event.target.value)}
              placeholder="https://supplier.example"
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
          label="Category"
          input={
            <select
              value={form.category}
              onChange={(event) => handleChange('category', event.target.value)}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              <option value="">Select category</option>
              {SUPPLIER_CATEGORY_OPTIONS.map((category) => (
                <option key={category} value={category}>
                  {category}
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
              {SUPPLIER_PAYMENT_TERMS.map((term) => (
                <option key={term} value={term}>
                  {term}
                </option>
              ))}
            </select>
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
          label="Payable account"
          input={
            <select
              value={form.payable_account_id}
              onChange={(event) => handleChange('payable_account_id', event.target.value)}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              <option value="">Select payable account</option>
              {payableAccounts.map((account) => (
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
      </div>

      <Field
        label="Address"
        input={
          <Textarea
            value={form.address}
            onChange={(event) => handleChange('address', event.target.value)}
            placeholder="Supplier office address"
          />
        }
      />
      <Field
        label="Arabic address"
        input={
          <Textarea
            value={form.address_ar}
            onChange={(event) => handleChange('address_ar', event.target.value)}
            placeholder="عنوان المورد"
          />
        }
      />

      <div className="grid gap-4 md:grid-cols-3">
        <Field
          label="City"
          input={
            <Input
              value={form.city}
              onChange={(event) => handleChange('city', event.target.value)}
              placeholder="Cairo"
            />
          }
        />
        <Field
          label="State"
          input={
            <Input
              value={form.state}
              onChange={(event) => handleChange('state', event.target.value)}
              placeholder="Cairo Governorate"
            />
          }
        />
        <Field
          label="Postal code"
          input={
            <Input
              value={form.postal_code}
              onChange={(event) => handleChange('postal_code', event.target.value)}
              placeholder="11511"
            />
          }
        />
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <Field
          label="Bank name"
          input={
            <Input
              value={form.bank_name}
              onChange={(event) => handleChange('bank_name', event.target.value)}
              placeholder="Banque Misr"
            />
          }
        />
        <Field
          label="Bank account"
          input={
            <Input
              value={form.bank_account}
              onChange={(event) => handleChange('bank_account', event.target.value)}
              placeholder="00123456789"
            />
          }
        />
        <Field
          label="IBAN"
          input={
            <Input
              value={form.iban}
              onChange={(event) => handleChange('iban', event.target.value)}
              placeholder="EG38..."
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
            placeholder="Preferred delivery windows or account comments."
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
            <div className="font-medium text-gray-900 dark:text-white">Active supplier</div>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Disable this supplier to prevent use in future purchasing workflows.
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
          {submitting ? 'Saving...' : initialValues ? 'Update supplier' : 'Create supplier'}
        </Button>
      </div>
    </form>
  )
}

function getInitialState(initialValues?: Supplier | null): SupplierState {
  return {
    company_name: initialValues?.company_name || '',
    company_name_ar: initialValues?.company_name_ar || '',
    contact_person: initialValues?.contact_person || '',
    email: initialValues?.email || '',
    phone: initialValues?.phone || '',
    mobile: initialValues?.mobile || '',
    website: initialValues?.website || '',
    address: initialValues?.address || '',
    address_ar: initialValues?.address_ar || '',
    city: initialValues?.city || '',
    state: initialValues?.state || '',
    postal_code: initialValues?.postal_code || '',
    country: initialValues?.country || 'Egypt',
    payment_terms: initialValues?.payment_terms || 'net_30',
    credit_limit: String(initialValues?.credit_limit ?? 0),
    currency: initialValues?.currency || 'EGP',
    payable_account_id: initialValues?.payable_account_id
      ? String(initialValues.payable_account_id)
      : '',
    category: initialValues?.category || '',
    bank_name: initialValues?.bank_name || '',
    bank_account: initialValues?.bank_account || '',
    iban: initialValues?.iban || '',
    tax_number: initialValues?.tax_number || '',
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
