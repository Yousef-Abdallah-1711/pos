'use client'

import { useEffect, useMemo, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import type { Account, AccountFormInput } from '@/types/accounting'
import {
  ACCOUNT_TYPE_OPTIONS,
  DEFAULT_NORMAL_BALANCE,
  formatAccountLabel,
} from './accounting-utils'

interface AccountFormProps {
  accounts: Account[]
  initialValues?: Account | null
  onSubmit: (payload: AccountFormInput) => Promise<void> | void
  onCancel: () => void
  submitting?: boolean
}

type AccountFormState = {
  account_code: string
  account_name: string
  account_name_ar: string
  description: string
  description_ar: string
  account_type: AccountFormInput['account_type']
  account_subtype: string
  normal_balance: AccountFormInput['normal_balance']
  parent_id: string
  is_header: boolean
  currency: string
  is_active: boolean
  is_reconcilable: boolean
  tax_rate: string
}

export function AccountForm({
  accounts,
  initialValues,
  onSubmit,
  onCancel,
  submitting = false,
}: AccountFormProps) {
  const [form, setForm] = useState<AccountFormState>(() => getInitialState(initialValues))
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    setForm(getInitialState(initialValues))
    setErrors({})
  }, [initialValues])

  const parentOptions = useMemo(
    () =>
      accounts.filter(
        (account) =>
          account.id !== initialValues?.id &&
          account.account_type === form.account_type &&
          account.is_header
      ),
    [accounts, form.account_type, initialValues?.id]
  )

  const handleChange = <K extends keyof AccountFormState>(
    key: K,
    value: AccountFormState[K]
  ) => {
    setForm((current) => {
      const next = { ...current, [key]: value }
      if (key === 'account_type') {
        next.normal_balance = DEFAULT_NORMAL_BALANCE[value as AccountFormInput['account_type']]
        next.parent_id = ''
      }
      return next
    })
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const nextErrors: Record<string, string> = {}
    if (!form.account_code.trim()) {
      nextErrors.account_code = 'Account code is required.'
    }
    if (!form.account_name.trim()) {
      nextErrors.account_name = 'Account name is required.'
    }

    const duplicate = accounts.find(
      (account) =>
        account.id !== initialValues?.id &&
        account.account_code.toLowerCase() === form.account_code.trim().toLowerCase()
    )
    if (duplicate) {
      nextErrors.account_code = 'Account code must be unique.'
    }

    if (form.parent_id) {
      const parent = accounts.find((account) => account.id === Number(form.parent_id))
      if (!parent) {
        nextErrors.parent_id = 'Selected parent account was not found.'
      } else if (parent.account_type !== form.account_type) {
        nextErrors.parent_id = 'Parent account must use the same account type.'
      }
    }

    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors)
      return
    }

    setErrors({})

    await onSubmit({
      account_code: form.account_code.trim(),
      account_name: form.account_name.trim(),
      account_name_ar: form.account_name_ar.trim() || undefined,
      description: form.description.trim() || undefined,
      description_ar: form.description_ar.trim() || undefined,
      account_type: form.account_type,
      account_subtype: form.account_subtype.trim() || undefined,
      normal_balance: form.normal_balance,
      parent_id: form.parent_id ? Number(form.parent_id) : null,
      is_header: form.is_header,
      currency: form.currency.trim() || 'EGP',
      is_active: form.is_active,
      is_reconcilable: form.is_reconcilable,
      tax_rate: form.tax_rate === '' ? null : Number(form.tax_rate),
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2">
        <Field
          label="Account code"
          error={errors.account_code}
          input={
            <Input
              value={form.account_code}
              onChange={(event) => handleChange('account_code', event.target.value)}
              placeholder="e.g. 1100"
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
          label="Account name"
          error={errors.account_name}
          input={
            <Input
              value={form.account_name}
              onChange={(event) => handleChange('account_name', event.target.value)}
              placeholder="Cash on Hand"
            />
          }
        />
        <Field
          label="Arabic name"
          input={
            <Input
              value={form.account_name_ar}
              onChange={(event) => handleChange('account_name_ar', event.target.value)}
              placeholder="النقدية بالصندوق"
            />
          }
        />
        <Field
          label="Account type"
          input={
            <select
              value={form.account_type}
              onChange={(event) =>
                handleChange(
                  'account_type',
                  event.target.value as AccountFormInput['account_type']
                )
              }
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              {ACCOUNT_TYPE_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          }
        />
        <Field
          label="Subtype"
          input={
            <Input
              value={form.account_subtype}
              onChange={(event) => handleChange('account_subtype', event.target.value)}
              placeholder="Current asset"
            />
          }
        />
        <Field
          label="Normal balance"
          input={
            <select
              value={form.normal_balance}
              onChange={(event) =>
                handleChange(
                  'normal_balance',
                  event.target.value as AccountFormInput['normal_balance']
                )
              }
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              <option value="debit">Debit</option>
              <option value="credit">Credit</option>
            </select>
          }
        />
        <Field
          label="Parent account"
          error={errors.parent_id}
          input={
            <select
              value={form.parent_id}
              onChange={(event) => handleChange('parent_id', event.target.value)}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              <option value="">No parent</option>
              {parentOptions.map((account) => (
                <option key={account.id} value={account.id}>
                  {formatAccountLabel(account)}
                </option>
              ))}
            </select>
          }
        />
        <Field
          label="Tax rate (%)"
          input={
            <Input
              type="number"
              step="0.01"
              value={form.tax_rate}
              onChange={(event) => handleChange('tax_rate', event.target.value)}
              placeholder="14"
            />
          }
        />
      </div>

      <Field
        label="Description"
        input={
          <Textarea
            value={form.description}
            onChange={(event) => handleChange('description', event.target.value)}
            placeholder="Used for till cash, petty cash, or bank float."
          />
        }
      />
      <Field
        label="Arabic description"
        input={
          <Textarea
            value={form.description_ar}
            onChange={(event) => handleChange('description_ar', event.target.value)}
            placeholder="يستخدم للنقدية أو العهدة أو رصيد الصندوق."
          />
        }
      />

      <div className="grid gap-4 md:grid-cols-3">
        <ToggleField
          label="Header account"
          description="Use this account as a grouping node for child accounts."
          checked={form.is_header}
          onCheckedChange={(checked) => handleChange('is_header', checked)}
        />
        <ToggleField
          label="Reconcilable"
          description="Mark this account for statement or balance reconciliation."
          checked={form.is_reconcilable}
          onCheckedChange={(checked) => handleChange('is_reconcilable', checked)}
        />
        <ToggleField
          label="Active"
          description="Inactive accounts stay visible in history but cannot be used."
          checked={form.is_active}
          onCheckedChange={(checked) => handleChange('is_active', checked)}
        />
      </div>

      <div className="flex justify-end gap-3 border-t border-gray-200 pt-4 dark:border-gray-700">
        <Button type="button" variant="outline" onClick={onCancel} disabled={submitting}>
          Cancel
        </Button>
        <Button type="submit" disabled={submitting}>
          {submitting ? 'Saving...' : initialValues ? 'Update account' : 'Create account'}
        </Button>
      </div>
    </form>
  )
}

function getInitialState(initialValues?: Account | null): AccountFormState {
  const accountType = initialValues?.account_type || 'asset'

  return {
    account_code: initialValues?.account_code || '',
    account_name: initialValues?.account_name || '',
    account_name_ar: initialValues?.account_name_ar || '',
    description: initialValues?.description || '',
    description_ar: initialValues?.description_ar || '',
    account_type: accountType,
    account_subtype: initialValues?.account_subtype || '',
    normal_balance:
      initialValues?.normal_balance ||
      DEFAULT_NORMAL_BALANCE[accountType as AccountFormInput['account_type']],
    parent_id: initialValues?.parent_id ? String(initialValues.parent_id) : '',
    is_header: initialValues?.is_header || false,
    currency: initialValues?.currency || 'EGP',
    is_active: initialValues?.is_active ?? true,
    is_reconcilable: initialValues?.is_reconcilable || false,
    tax_rate:
      typeof initialValues?.tax_rate === 'number' ? String(initialValues.tax_rate) : '',
  }
}

function Field({
  label,
  input,
  error,
}: {
  label: string
  input: React.ReactNode
  error?: string
}) {
  return (
    <div className="space-y-2">
      <Label>{label}</Label>
      {input}
      {error ? <p className="text-sm text-red-600 dark:text-red-400">{error}</p> : null}
    </div>
  )
}

function ToggleField({
  label,
  description,
  checked,
  onCheckedChange,
}: {
  label: string
  description: string
  checked: boolean
  onCheckedChange: (checked: boolean) => void
}) {
  return (
    <div className="rounded-xl border border-gray-200 p-4 dark:border-gray-700">
      <div className="flex items-start gap-3">
        <Checkbox
          checked={checked}
          onCheckedChange={(value) => onCheckedChange(value === true)}
          className="mt-1"
        />
        <div>
          <div className="font-medium text-gray-900 dark:text-white">{label}</div>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{description}</p>
        </div>
      </div>
    </div>
  )
}
