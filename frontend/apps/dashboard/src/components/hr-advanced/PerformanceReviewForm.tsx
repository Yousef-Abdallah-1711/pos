'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import type { PerformanceReview } from '@/types/hr-advanced'

interface Props {
  initial?: Partial<PerformanceReview>
  onSubmit: (payload: Partial<PerformanceReview>) => Promise<void>
  submitting?: boolean
}

export function PerformanceReviewForm({ initial, onSubmit, submitting = false }: Props) {
  const t = useTranslations('hrAdvanced')
  const [form, setForm] = useState<Partial<PerformanceReview>>({
    employee_id: initial?.employee_id ?? 0,
    reviewer_id: initial?.reviewer_id ?? 0,
    review_period_start: initial?.review_period_start ?? '',
    review_period_end: initial?.review_period_end ?? '',
    review_date: initial?.review_date ?? '',
    goals_achievement: initial?.goals_achievement ?? 0,
    technical_skills: initial?.technical_skills ?? 0,
    communication: initial?.communication ?? 0,
    leadership: initial?.leadership ?? 0,
    strengths: initial?.strengths ?? '',
    areas_to_improve: initial?.areas_to_improve ?? '',
    action_plan: initial?.action_plan ?? '',
    status: initial?.status ?? 'draft',
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await onSubmit(form)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
        <input type="number" placeholder={t('forms.employeeId')} value={form.employee_id || ''} onChange={(e) => setForm((p) => ({ ...p, employee_id: Number(e.target.value) }))} className="rounded border px-3 py-2" required />
        <input type="number" placeholder={t('forms.reviewerId')} value={form.reviewer_id || ''} onChange={(e) => setForm((p) => ({ ...p, reviewer_id: Number(e.target.value) }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.review_period_start || ''} onChange={(e) => setForm((p) => ({ ...p, review_period_start: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.review_period_end || ''} onChange={(e) => setForm((p) => ({ ...p, review_period_end: e.target.value }))} className="rounded border px-3 py-2" required />
        <input type="date" value={form.review_date || ''} onChange={(e) => setForm((p) => ({ ...p, review_date: e.target.value }))} className="rounded border px-3 py-2" required />
        <select value={form.status || 'draft'} onChange={(e) => setForm((p) => ({ ...p, status: e.target.value }))} className="rounded border px-3 py-2">
          <option value="draft">{t('status.draft')}</option>
          <option value="acknowledged">{t('status.acknowledged')}</option>
          <option value="completed">{t('status.completed')}</option>
        </select>
      </div>
      <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
        <input type="number" step="0.01" placeholder={t('forms.goals')} value={form.goals_achievement || 0} onChange={(e) => setForm((p) => ({ ...p, goals_achievement: Number(e.target.value) }))} className="rounded border px-3 py-2" />
        <input type="number" step="0.01" placeholder={t('forms.technical')} value={form.technical_skills || 0} onChange={(e) => setForm((p) => ({ ...p, technical_skills: Number(e.target.value) }))} className="rounded border px-3 py-2" />
        <input type="number" step="0.01" placeholder={t('forms.communication')} value={form.communication || 0} onChange={(e) => setForm((p) => ({ ...p, communication: Number(e.target.value) }))} className="rounded border px-3 py-2" />
        <input type="number" step="0.01" placeholder={t('forms.leadership')} value={form.leadership || 0} onChange={(e) => setForm((p) => ({ ...p, leadership: Number(e.target.value) }))} className="rounded border px-3 py-2" />
      </div>
      <textarea placeholder={t('forms.strengths')} value={form.strengths || ''} onChange={(e) => setForm((p) => ({ ...p, strengths: e.target.value }))} className="w-full rounded border px-3 py-2" rows={2} />
      <textarea placeholder={t('forms.areasToImprove')} value={form.areas_to_improve || ''} onChange={(e) => setForm((p) => ({ ...p, areas_to_improve: e.target.value }))} className="w-full rounded border px-3 py-2" rows={2} />
      <textarea placeholder={t('forms.actionPlan')} value={form.action_plan || ''} onChange={(e) => setForm((p) => ({ ...p, action_plan: e.target.value }))} className="w-full rounded border px-3 py-2" rows={2} />
      <button type="submit" disabled={submitting} className="rounded bg-primary-600 px-4 py-2 text-white disabled:opacity-50">{submitting ? t('saving') : t('save')}</button>
    </form>
  )
}
