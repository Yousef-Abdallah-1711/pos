'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import type { Resignation } from '@/types/hr-advanced'

interface Props {
  initial?: Partial<Resignation>
  onSubmit: (payload: Partial<Resignation>) => Promise<void>
  submitting?: boolean
}

export function ResignationForm({ initial, onSubmit, submitting = false }: Props) {
  const t = useTranslations('hrAdvanced')
  const [form, setForm] = useState<Partial<Resignation>>({
    employee_id: initial?.employee_id ?? 0,
    resignation_date: initial?.resignation_date ?? '',
    last_working_date: initial?.last_working_date ?? '',
    reason: initial?.reason ?? '',
    status: initial?.status ?? 'submitted',
    notice_period_days: initial?.notice_period_days ?? 0,
    pending_salary: initial?.pending_salary ?? 0,
    unused_leave_balance: initial?.unused_leave_balance ?? 0,
    settlement_amount: initial?.settlement_amount ?? 0,
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await onSubmit(form)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
        <input type="number" placeholder={t('forms.employeeId')} value={form.employee_id || ''} onChange={(e) => setForm((p) => ({ ...p, employee_id: Number(e.target.value) }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.resignation_date || ''} onChange={(e) => setForm((p) => ({ ...p, resignation_date: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.last_working_date || ''} onChange={(e) => setForm((p) => ({ ...p, last_working_date: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="number" placeholder={t('forms.noticePeriodDays')} value={form.notice_period_days || 0} onChange={(e) => setForm((p) => ({ ...p, notice_period_days: Number(e.target.value) }))} className="rounded border px-3 py-2" />
        <input type="number" step="0.01" placeholder={t('forms.pendingSalary')} value={form.pending_salary || 0} onChange={(e) => setForm((p) => ({ ...p, pending_salary: Number(e.target.value) }))} className="rounded border px-3 py-2" />
        <input type="number" step="0.01" placeholder={t('forms.unusedLeaveBalance')} value={form.unused_leave_balance || 0} onChange={(e) => setForm((p) => ({ ...p, unused_leave_balance: Number(e.target.value) }))} className="rounded border px-3 py-2" />
      </div>
      <textarea placeholder={t('forms.reason')} value={form.reason || ''} onChange={(e) => setForm((p) => ({ ...p, reason: e.target.value }))} className="w-full rounded border px-3 py-2" rows={3} required />
      <button type="submit" disabled={submitting} className="rounded bg-primary-600 px-4 py-2 text-white disabled:opacity-50">{submitting ? t('saving') : t('save')}</button>
    </form>
  )
}
