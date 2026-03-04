'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import type { TrainingRecord } from '@/types/hr-advanced'

interface Props {
  initial?: Partial<TrainingRecord>
  onSubmit: (payload: Partial<TrainingRecord>) => Promise<void>
  submitting?: boolean
}

export function TrainingForm({ initial, onSubmit, submitting = false }: Props) {
  const t = useTranslations('hrAdvanced')
  const [form, setForm] = useState<Partial<TrainingRecord>>({
    employee_id: initial?.employee_id ?? 0,
    training_title: initial?.training_title ?? '',
    training_provider: initial?.training_provider ?? '',
    training_type: initial?.training_type ?? '',
    start_date: initial?.start_date ?? '',
    end_date: initial?.end_date ?? '',
    status: initial?.status ?? 'scheduled',
    certificate_number: initial?.certificate_number ?? '',
    certificate_expiry: initial?.certificate_expiry ?? '',
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
        <input placeholder={t('forms.trainingTitle')} value={form.training_title || ''} onChange={(e) => setForm((p) => ({ ...p, training_title: e.target.value }))} className="rounded border px-3 py-2" required />
        <input placeholder={t('forms.provider')} value={form.training_provider || ''} onChange={(e) => setForm((p) => ({ ...p, training_provider: e.target.value }))} className="rounded border px-3 py-2" />
        <input placeholder={t('forms.trainingType')} value={form.training_type || ''} onChange={(e) => setForm((p) => ({ ...p, training_type: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.start_date || ''} onChange={(e) => setForm((p) => ({ ...p, start_date: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.end_date || ''} onChange={(e) => setForm((p) => ({ ...p, end_date: e.target.value }))} className="rounded border px-3 py-2" />
        <input placeholder={t('forms.certificateNumber')} value={form.certificate_number || ''} onChange={(e) => setForm((p) => ({ ...p, certificate_number: e.target.value }))} className="rounded border px-3 py-2" />
        <input type="date" value={form.certificate_expiry || ''} onChange={(e) => setForm((p) => ({ ...p, certificate_expiry: e.target.value }))} className="rounded border px-3 py-2" />
      </div>
      <textarea placeholder={t('forms.notes')} value={form.notes || ''} onChange={(e) => setForm((p) => ({ ...p, notes: e.target.value }))} className="w-full rounded border px-3 py-2" rows={2} />
      <button type="submit" disabled={submitting} className="rounded bg-primary-600 px-4 py-2 text-white disabled:opacity-50">{submitting ? t('saving') : t('save')}</button>
    </form>
  )
}
