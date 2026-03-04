'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { ExitChecklist } from '@/components/hr-advanced/ExitChecklist'
import { ResignationForm } from '@/components/hr-advanced/ResignationForm'
import { ResignationWorkflow } from '@/components/hr-advanced/ResignationWorkflow'
import { useResignations } from '@/hooks/useResignations'
import type { Resignation } from '@/types/hr-advanced'

export default function ResignationsPage() {
  const t = useTranslations('hrAdvanced')
  const { items, loading, error, createItem, updateItem, removeItem, acknowledge, complete, calculateSettlement } = useResignations()
  const defaultChecklist = [
    { key: 'return_assets', label: t('exitChecklist.returnAssets'), done: false },
    { key: 'final_attendance', label: t('exitChecklist.finalizeAttendance'), done: false },
    { key: 'payroll_clearance', label: t('exitChecklist.payrollClearance'), done: false },
  ]
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<Resignation | null>(null)
  const [saving, setSaving] = useState(false)
  const [checklist, setChecklist] = useState(defaultChecklist)

  const onSubmit = async (payload: Partial<Resignation>) => {
    setSaving(true)
    try {
      if (editing) {
        await updateItem(editing.id, payload)
      } else {
        await createItem(payload)
      }
      setOpen(false)
      setEditing(null)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t('resignations.title')}</h1>
          <p className="text-sm text-gray-500">{t('resignations.subtitle')}</p>
        </div>
        <button className="rounded bg-primary-600 px-4 py-2 text-white" onClick={() => setOpen(true)}>{t('resignations.add')}</button>
      </div>

      <ExitChecklist
        checklist={checklist}
        onToggle={(key) => setChecklist((prev) => prev.map((item) => (item.key === key ? { ...item, done: !item.done } : item)))}
      />

      {loading ? <div className="rounded border p-4 text-sm text-gray-500">{t('resignations.loading')}</div> : null}
      {error ? <div className="rounded border border-red-200 bg-red-50 p-4 text-sm text-red-700">{error}</div> : null}

      <div className="space-y-3">
        {items.map((item) => (
          <div key={item.id} className="rounded border p-4">
            <div className="mb-2 flex items-center justify-between">
              <div className="font-semibold">{t('resignations.employeeLabel', { id: item.employee_id })}</div>
              <div className="space-x-2">
                <button onClick={() => { setEditing(item); setOpen(true) }} className="rounded border px-2 py-1 text-sm">{t('actions.edit')}</button>
                <button onClick={() => void removeItem(item.id)} className="rounded border px-2 py-1 text-sm text-red-600">{t('actions.delete')}</button>
              </div>
            </div>
            <ResignationWorkflow
              item={item}
              onAcknowledge={acknowledge}
              onComplete={complete}
              onCalculate={calculateSettlement}
            />
          </div>
        ))}
        {items.length === 0 ? <div className="rounded border p-6 text-center text-sm text-gray-500">{t('resignations.noData')}</div> : null}
      </div>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-3xl">
          <DialogHeader>
            <DialogTitle>{editing ? t('resignations.edit') : t('resignations.new')}</DialogTitle>
          </DialogHeader>
          <ResignationForm initial={editing || undefined} onSubmit={onSubmit} submitting={saving} />
        </DialogContent>
      </Dialog>
    </div>
  )
}
