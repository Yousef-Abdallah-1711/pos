'use client'

import { useEffect, useMemo, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import type { Account } from '@/types/accounting'
import type { BankAccount, BankAccountFormInput } from '@/types/invoices'
import {
  BANK_ACCOUNT_TYPE_OPTIONS,
  formatAccountLabel,
} from './accounting-utils'

interface BankAccountFormProps {
  accounts: Account[]
  initialValues?: BankAccount | null
  onSubmit: (payload: BankAccountFormInput) => Promise<void> | void
  onCancel: () => void
  submitting?: boolean
}

type BankAccountState = {
  account_name: string
  account_name_ar: string
  bank_name: string
  bank_branch: string
  account_number: string
  iban: string
  swift_code: string
  currency: string
  account_type: BankAccountFormInput['account_type']
  opening_balance: string
  opening_date: string
  gl_account_id: string
  is_active: boolean
  is_primary: boolean
}

export function BankAccountForm({
  accounts,
  initialValues,
  onSubmit,
  onCancel,
  submitting = false,
}: BankAccountFormProps) {
  const [form, setForm] = useState<BankAccountState>(() => getInitialState(initialValues))
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setForm(getInitialState(initialValues))
    setError(null)
  }, [initialValues])

  const assetAccounts = useMemo(
    () =>
      accounts.filter(
        (account) => account.account_type === 'asset' && !account.is_header
      ),
    [accounts]
  )

  const handleChange = <K extends keyof BankAccountState>(
    key: K,
    value: BankAccountState[K]
  ) => {
    setForm((current) => ({ ...current, [key]: value }))
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (!form.account_name.trim()) {
      setError('Account name is required.')
      return
    }
    if (!form.bank_name.trim()) {
      setError('Bank name is required.')
      return
    }
    if (!form.account_number.trim()) {
      setError('Account number is required.')
      return
    }

    setError(null)

    await onSubmit({
      account_name: form.account_name.trim(),
      account_name_ar: form.account_name_ar.trim() || undefined,
      bank_name: form.bank_name.trim(),
      bank_branch: form.bank_branch.trim() || undefined,
      account_number: form.account_number.trim(),
      iban: form.iban.trim() || undefined,
      swift_code: form.swift_code.trim() || undefined,
      currency: form.currency.trim() || 'EGP',
      account_type: form.account_type,
      opening_balance: Number(form.opening_balance) || 0,
      opening_date: form.opening_date,
      gl_account_id: form.gl_account_id ? Number(form.gl_account_id) : null,
      is_active: form.is_active,
      is_primary: form.is_primary,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2">
        <Field
          label="Account name"
          input={
            <Input
              value={form.account_name}
              onChange={(event) => handleChange('account_name', event.target.value)}
              placeholder="Main operating account"
            />
          }
        />
        <Field
          label="Arabic name"
          input={
            <Input
              value={form.account_name_ar}
              onChange={(event) => handleChange('account_name_ar', event.target.value)}
              placeholder="الحساب البنكي الرئيسي"
            />
          }
        />
        <Field
          label="Bank name"
          input={
            <Input
              value={form.bank_name}
              onChange={(event) => handleChange('bank_name', event.target.value)}
              placeholder="National Bank of Egypt"
            />
          }
        />
        <Field
          label="Branch"
          input={
            <Input
              value={form.bank_branch}
              onChange={(event) => handleChange('bank_branch', event.target.value)}
              placeholder="Maadi branch"
            />
          }
        />
        <Field
          label="Account number"
          input={
            <Input
              value={form.account_number}
              onChange={(event) => handleChange('account_number', event.target.value)}
              placeholder="001234567890"
            />
          }
        />
        <Field
          label="IBAN"
          input={
            <Input
              value={form.iban}
              onChange={(event) => handleChange('iban', event.target.value)}
              placeholder="EG380019000500000000263180002"
            />
          }
        />
        <Field
          label="Swift code"
          input={
            <Input
              value={form.swift_code}
              onChange={(event) => handleChange('swift_code', event.target.value)}
              placeholder="NBEGEGCX"
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
          label="Account type"
          input={
            <select
              value={form.account_type}
              onChange={(event) =>
                handleChange(
                  'account_type',
                  event.target.value as BankAccountState['account_type']
                )
              }
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              {BANK_ACCOUNT_TYPE_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          }
        />
        <Field
          label="Opening balance"
          input={
            <Input
              type="number"
              step="0.01"
              value={form.opening_balance}
              onChange={(event) => handleChange('opening_balance', event.target.value)}
              placeholder="0.00"
            />
          }
        />
        <Field
          label="Opening date"
          input={
            <Input
              type="date"
              value={form.opening_date}
              onChange={(event) => handleChange('opening_date', event.target.value)}
            />
          }
        />
        <Field
          label="Linked GL account"
          input={
            <select
              value={form.gl_account_id}
              onChange={(event) => handleChange('gl_account_id', event.target.value)}
              className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
            >
              <option value="">Select ledger account</option>
              {assetAccounts.map((account) => (
                <option key={account.id} value={account.id}>
                  {formatAccountLabel(account)}
                </option>
              ))}
            </select>
          }
        />
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <ToggleCard
          checked={form.is_active}
          onChange={(checked) => handleChange('is_active', checked)}
          title="Active account"
          description="Inactive bank accounts stay visible for history but disappear from new payment forms."
        />
        <ToggleCard
          checked={form.is_primary}
          onChange={(checked) => handleChange('is_primary', checked)}
          title="Primary settlement account"
          description="Mark the default bank account used most often for transfers and settlements."
        />
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
              ? 'Update bank account'
              : 'Create bank account'}
        </Button>
      </div>
    </form>
  )
}

function getInitialState(initialValues?: BankAccount | null): BankAccountState {
  return {
    account_name: initialValues?.account_name || '',
    account_name_ar: initialValues?.account_name_ar || '',
    bank_name: initialValues?.bank_name || '',
    bank_branch: initialValues?.bank_branch || '',
    account_number: initialValues?.account_number || '',
    iban: initialValues?.iban || '',
    swift_code: initialValues?.swift_code || '',
    currency: initialValues?.currency || 'EGP',
    account_type: initialValues?.account_type || 'checking',
    opening_balance: String(initialValues?.opening_balance ?? 0),
    opening_date: initialValues?.opening_date?.slice(0, 10) || getTodayDate(),
    gl_account_id: initialValues?.gl_account_id ? String(initialValues.gl_account_id) : '',
    is_active: initialValues?.is_active ?? true,
    is_primary: initialValues?.is_primary ?? false,
  }
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

function ToggleCard({
  checked,
  onChange,
  title,
  description,
}: {
  checked: boolean
  onChange: (checked: boolean) => void
  title: string
  description: string
}) {
  return (
    <div className="rounded-xl border border-gray-200 p-4 dark:border-gray-700">
      <div className="flex items-start gap-3">
        <Checkbox
          checked={checked}
          onCheckedChange={(value) => onChange(value === true)}
          className="mt-1"
        />
        <div>
          <div className="font-medium text-gray-900 dark:text-white">{title}</div>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{description}</p>
        </div>
      </div>
    </div>
  )
}
