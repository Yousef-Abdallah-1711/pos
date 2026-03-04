'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import type { InsuranceTracking } from '@/types/hr-advanced'

interface Props {
  initial?: Partial<InsuranceTracking>
  onSubmit: (payload: Partial<InsuranceTracking>) => Promise<void>
  submitting?: boolean
}

export function InsuranceForm({ initial, onSubmit, submitting = false }: Props) {
  const t = useTranslations('hrAdvanced')
  const [form, setForm] = useState<Partial<InsuranceTracking>>({
    employee_id: initial?.employee_id ?? 0,
    policy_type: initial?.policy_type ?? '',
    provider_name: initial?.provider_name ?? '',
    policy_number: initial?.policy_number ?? '',
    coverage_start: initial?.coverage_start ?? '',
    coverage_end: initial?.coverage_end ?? '',
    monthly_premium: initial?.monthly_premium ?? 0,
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
        <input placeholder={t('forms.policyType')} value={form.policy_type || ''} onChange={(e) => setForm((p) => ({ ...p, policy_type: e.target.value }))} className="rounded border px-3 py-2" required />
        <input placeholder={t('forms.provider')} value={form.provider_name || ''} onChange={(e) => setForm((p) => ({ ...p, provider_name: e.target.value }))} className="rounded border px-3 py-2" required />
        <input placeholder={t('forms.policyNumber')} value={form.policy_number || ''} onChange={(e) => setForm((p) => ({ ...p, policy_number: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.coverage_start || ''} onChange={(e) => setForm((p) => ({ ...p, coverage_start: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.coverage_end || ''} onChange={(e) => setForm((p) => ({ ...p, coverage_end: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="number" step="0.01" placeholder={t('forms.monthlyPremium')} value={form.monthly_premium || 0} onChange={(e) => setForm((p) => ({ ...p, monthly_premium: Number(e.target.value) }))} className="rounded border px-3 py-2" />
        <select value={form.status || 'active'} onChange={(e) => setForm((p) => ({ ...p, status: e.target.value }))} className="rounded border px-3 py-2">
          <option value="active">{t('status.active')}</option>
          <option value="expired">{t('status.expired')}</option>
          <option value="cancelled">{t('status.cancelled')}</option>
        </select>
      </div>
      <textarea placeholder={t('forms.notes')} value={form.notes || ''} onChange={(e) => setForm((p) => ({ ...p, notes: e.target.value }))} className="w-full rounded border px-3 py-2" rows={2} />
      <button type="submit" disabled={submitting} className="rounded bg-primary-600 px-4 py-2 text-white disabled:opacity-50">{submitting ? t('saving') : t('save')}</button>
    </form>
  )
}
