'use client'

import { useEffect, useMemo, useState } from 'react'
import { Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import type { Account, JournalEntryFormInput, JournalEntryLine } from '@/types/accounting'
import {
  formatAccountLabel,
  formatCurrency,
  parseNumberInput,
} from './accounting-utils'

interface JournalEntryFormProps {
  locale: string
  accounts: Account[]
  initialValues?: Partial<JournalEntryFormInput>
  onSubmit: (payload: JournalEntryFormInput) => Promise<void> | void
  onCancel: () => void
  submitting?: boolean
}

type FormState = {
  entry_date: string
  description: string
  description_ar: string
  reference_type: string
  reference_number: string
  entry_type: string
  notes: string
  lines: JournalEntryLine[]
}

export function JournalEntryForm({
  locale,
  accounts,
  initialValues,
  onSubmit,
  onCancel,
  submitting = false,
}: JournalEntryFormProps) {
  const [form, setForm] = useState<FormState>(() => getInitialState(initialValues))
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setForm(getInitialState(initialValues))
    setError(null)
  }, [initialValues])

  const totals = useMemo(() => {
    const totalDebit = form.lines.reduce((sum, line) => sum + line.debit_amount, 0)
    const totalCredit = form.lines.reduce((sum, line) => sum + line.credit_amount, 0)
    return {
      totalDebit,
      totalCredit,
      isBalanced:
        form.lines.length >= 2 &&
        totalDebit > 0 &&
        totalCredit > 0 &&
        Math.abs(totalDebit - totalCredit) < 0.001,
    }
  }, [form.lines])

  const handleHeaderChange = <K extends Exclude<keyof FormState, 'lines'>>(
    key: K,
    value: FormState[K]
  ) => {
    setForm((current) => ({ ...current, [key]: value }))
  }

  const handleLineChange = (
    index: number,
    key: keyof JournalEntryLine,
    value: string | number
  ) => {
    setForm((current) => {
      const lines = current.lines.map((line, lineIndex) => {
        if (lineIndex !== index) {
          return line
        }

        const nextLine = { ...line, [key]: value }
        if (key === 'debit_amount' && Number(value) > 0) {
          nextLine.credit_amount = 0
        }
        if (key === 'credit_amount' && Number(value) > 0) {
          nextLine.debit_amount = 0
        }
        return nextLine
      })

      return { ...current, lines }
    })
  }

  const addLine = () => {
    setForm((current) => ({
      ...current,
      lines: [...current.lines, createEmptyLine(current.lines.length + 1)],
    }))
  }

  const removeLine = (index: number) => {
    setForm((current) => ({
      ...current,
      lines: current.lines
        .filter((_, lineIndex) => lineIndex !== index)
        .map((line, lineIndex) => ({ ...line, line_order: lineIndex + 1 })),
    }))
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (!form.description.trim()) {
      setError('Description is required.')
      return
    }

    if (!totals.isBalanced) {
      setError('Journal entry must be balanced before saving.')
      return
    }

    const lines = form.lines
      .map((line, index) => ({
        ...line,
        line_order: index + 1,
        description: line.description?.trim() || undefined,
      }))
      .filter(
        (line) =>
          line.account_id > 0 && (line.debit_amount > 0 || line.credit_amount > 0)
      )

    if (lines.length < 2) {
      setError('At least two valid lines are required.')
      return
    }

    setError(null)

    await onSubmit({
      entry_date: form.entry_date,
      description: form.description.trim(),
      description_ar: form.description_ar.trim() || undefined,
      reference_type: form.reference_type.trim() || undefined,
      reference_number: form.reference_number.trim() || undefined,
      entry_type: form.entry_type.trim() || 'manual',
      notes: form.notes.trim() || undefined,
      lines,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2">
        <Field
          label="Entry date"
          input={
            <Input
              type="date"
              value={form.entry_date}
              onChange={(event) => handleHeaderChange('entry_date', event.target.value)}
            />
          }
        />
        <Field
          label="Entry type"
          input={
            <Input
              value={form.entry_type}
              onChange={(event) => handleHeaderChange('entry_type', event.target.value)}
              placeholder="manual"
            />
          }
        />
        <Field
          label="Description"
          input={
            <Input
              value={form.description}
              onChange={(event) => handleHeaderChange('description', event.target.value)}
              placeholder="Record petty cash replenishment"
            />
          }
        />
        <Field
          label="Arabic description"
          input={
            <Input
              value={form.description_ar}
              onChange={(event) => handleHeaderChange('description_ar', event.target.value)}
              placeholder="قيد تغذية عهدة المصروفات"
            />
          }
        />
        <Field
          label="Reference type"
          input={
            <Input
              value={form.reference_type}
              onChange={(event) => handleHeaderChange('reference_type', event.target.value)}
              placeholder="expense_claim"
            />
          }
        />
        <Field
          label="Reference number"
          input={
            <Input
              value={form.reference_number}
              onChange={(event) => handleHeaderChange('reference_number', event.target.value)}
              placeholder="EXP-2026-0012"
            />
          }
        />
      </div>

      <Field
        label="Notes"
        input={
          <Textarea
            value={form.notes}
            onChange={(event) => handleHeaderChange('notes', event.target.value)}
            placeholder="Optional notes for approvers or auditors."
          />
        }
      />

      <div className="overflow-hidden rounded-2xl border border-gray-200 dark:border-gray-700">
        <div className="flex flex-col gap-3 border-b border-gray-200 bg-gray-50 p-4 dark:border-gray-700 dark:bg-gray-900/40 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h3 className="font-semibold text-gray-900 dark:text-white">Journal lines</h3>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Add debit and credit lines until the entry balances.
            </p>
          </div>
          <div className="flex flex-wrap gap-3">
            <div className="rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-800">
              Debit: {formatCurrency(totals.totalDebit, locale)}
            </div>
            <div className="rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-800">
              Credit: {formatCurrency(totals.totalCredit, locale)}
            </div>
            <div
              className={`rounded-lg px-3 py-2 text-sm font-medium ${
                totals.isBalanced
                  ? 'border border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900 dark:bg-emerald-950/30 dark:text-emerald-300'
                  : 'border border-red-200 bg-red-50 text-red-700 dark:border-red-900 dark:bg-red-950/30 dark:text-red-300'
              }`}
            >
              {totals.isBalanced ? 'Balanced' : 'Not balanced'}
            </div>
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
            <thead className="bg-white dark:bg-gray-800">
              <tr>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Account
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
                <th className="px-4 py-3 text-right font-medium text-gray-500 dark:text-gray-400">
                  Action
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 bg-white dark:divide-gray-700 dark:bg-gray-800">
              {form.lines.map((line, index) => (
                <tr key={`line-${index}`}>
                  <td className="px-4 py-3 align-top">
                    <select
                      value={line.account_id || ''}
                      onChange={(event) =>
                        handleLineChange(index, 'account_id', Number(event.target.value))
                      }
                      className="flex h-10 w-72 rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
                    >
                      <option value="">Select account</option>
                      {accounts.map((account) => (
                        <option key={account.id} value={account.id}>
                          {formatAccountLabel(account)}
                        </option>
                      ))}
                    </select>
                  </td>
                  <td className="px-4 py-3 align-top">
                    <Input
                      value={line.description || ''}
                      onChange={(event) =>
                        handleLineChange(index, 'description', event.target.value)
                      }
                      placeholder="Optional line memo"
                    />
                  </td>
                  <td className="px-4 py-3 align-top">
                    <Input
                      type="number"
                      min="0"
                      step="0.01"
                      value={line.debit_amount === 0 ? '' : line.debit_amount}
                      onChange={(event) =>
                        handleLineChange(
                          index,
                          'debit_amount',
                          parseNumberInput(event.target.value)
                        )
                      }
                      placeholder="0.00"
                    />
                  </td>
                  <td className="px-4 py-3 align-top">
                    <Input
                      type="number"
                      min="0"
                      step="0.01"
                      value={line.credit_amount === 0 ? '' : line.credit_amount}
                      onChange={(event) =>
                        handleLineChange(
                          index,
                          'credit_amount',
                          parseNumberInput(event.target.value)
                        )
                      }
                      placeholder="0.00"
                    />
                  </td>
                  <td className="px-4 py-3 text-right align-top">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => removeLine(index)}
                      disabled={form.lines.length <= 2}
                      className="text-red-600 hover:text-red-700 dark:text-red-400"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="border-t border-gray-200 p-4 dark:border-gray-700">
          <Button type="button" variant="outline" onClick={addLine}>
            <Plus className="mr-2 h-4 w-4" />
            Add line
          </Button>
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
          {submitting ? 'Saving...' : 'Save journal entry'}
        </Button>
      </div>
    </form>
  )
}

function getInitialState(initialValues?: Partial<JournalEntryFormInput>): FormState {
  return {
    entry_date: initialValues?.entry_date || new Date().toISOString().split('T')[0],
    description: initialValues?.description || '',
    description_ar: initialValues?.description_ar || '',
    reference_type: initialValues?.reference_type || '',
    reference_number: initialValues?.reference_number || '',
    entry_type: initialValues?.entry_type || 'manual',
    notes: initialValues?.notes || '',
    lines:
      initialValues?.lines?.length && initialValues.lines.length > 0
        ? initialValues.lines.map((line, index) => ({
            ...line,
            line_order: index + 1,
          }))
        : [createEmptyLine(1), createEmptyLine(2)],
  }
}

function createEmptyLine(lineOrder: number): JournalEntryLine {
  return {
    account_id: 0,
    debit_amount: 0,
    credit_amount: 0,
    description: '',
    line_order: lineOrder,
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
