'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import type { EmploymentContract } from '@/types/hr-advanced'

interface Props {
  initial?: Partial<EmploymentContract>
  onSubmit: (payload: Partial<EmploymentContract>) => Promise<void>
  submitting?: boolean
}

export function ContractForm({ initial, onSubmit, submitting = false }: Props) {
  const t = useTranslations('hrAdvanced')
  const [form, setForm] = useState<Partial<EmploymentContract>>({
    employee_id: initial?.employee_id ?? 0,
    contract_type: initial?.contract_type ?? '',
    contract_number: initial?.contract_number ?? '',
    start_date: initial?.start_date ?? '',
    end_date: initial?.end_date ?? '',
    salary_amount: initial?.salary_amount ?? 0,
    currency: initial?.currency ?? 'EGP',
    status: initial?.status ?? 'active',
    notes: initial?.notes ?? '',
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await onSubmit(form)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
        <input type="number" placeholder={t('forms.employeeId')} value={form.employee_id || ''} onChange={(e) => setForm((p) => ({ ...p, employee_id: Number(e.target.value) }))} className="rounded border px-3 py-2" required />
        <input placeholder={t('forms.contractType')} value={form.contract_type || ''} onChange={(e) => setForm((p) => ({ ...p, contract_type: e.target.value }))} className="rounded border px-3 py-2" required />
        <input placeholder={t('forms.contractNumber')} value={form.contract_number || ''} onChange={(e) => setForm((p) => ({ ...p, contract_number: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.start_date || ''} onChange={(e) => setForm((p) => ({ ...p, start_date: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.end_date || ''} onChange={(e) => setForm((p) => ({ ...p, end_date: e.target.value }))} className="rounded border px-3 py-2" />
        <input type="number" step="0.01" placeholder={t('forms.salaryAmount')} value={form.salary_amount || 0} onChange={(e) => setForm((p) => ({ ...p, salary_amount: Number(e.target.value) }))} className="rounded border px-3 py-2" required />
        <input placeholder={t('forms.currency')} value={form.currency || ''} onChange={(e) => setForm((p) => ({ ...p, currency: e.target.value }))} className="rounded border px-3 py-2" />
        <select value={form.status || 'active'} onChange={(e) => setForm((p) => ({ ...p, status: e.target.value }))} className="rounded border px-3 py-2">
          <option value="active">{t('status.active')}</option>
          <option value="expired">{t('status.expired')}</option>
          <option value="terminated">{t('status.terminated')}</option>
          <option value="renewed">{t('status.renewed')}</option>
        </select>
      </div>
      <textarea placeholder={t('forms.notes')} value={form.notes || ''} onChange={(e) => setForm((p) => ({ ...p, notes: e.target.value }))} className="w-full rounded border px-3 py-2" rows={2} />
      <button type="submit" disabled={submitting} className="rounded bg-primary-600 px-4 py-2 text-white disabled:opacity-50">{submitting ? t('saving') : t('save')}</button>
    </form>
  )
}
