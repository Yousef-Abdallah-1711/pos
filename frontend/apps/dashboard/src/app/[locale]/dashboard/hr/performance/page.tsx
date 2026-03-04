'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { PerformanceReviewForm } from '@/components/hr-advanced/PerformanceReviewForm'
import { PerformanceReviewList } from '@/components/hr-advanced/PerformanceReviewList'
import { usePerformanceReviews } from '@/hooks/usePerformanceReviews'
import type { PerformanceReview } from '@/types/hr-advanced'

export default function PerformancePage() {
  const t = useTranslations('hrAdvanced.performance')
  const { items, loading, error, createItem, updateItem, removeItem, acknowledge, complete } = usePerformanceReviews()
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<PerformanceReview | null>(null)
  const [saving, setSaving] = useState(false)

  const onSubmit = async (payload: Partial<PerformanceReview>) => {
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
          <h1 className="text-2xl font-bold">{t('title')}</h1>
          <p className="text-sm text-gray-500">{t('subtitle')}</p>
        </div>
        <button
          className="rounded bg-primary-600 px-4 py-2 text-white"
          onClick={() => {
            setEditing(null)
            setOpen(true)
          }}
        >
          {t('add')}
        </button>
      </div>
      {loading ? <div className="rounded border p-4 text-sm text-gray-500">{t('loading')}</div> : null}
      {error ? <div className="rounded border border-red-200 bg-red-50 p-4 text-sm text-red-700">{error}</div> : null}
      <PerformanceReviewList
        items={items}
        onEdit={(item) => {
          setEditing(item)
          setOpen(true)
        }}
        onDelete={removeItem}
        onAcknowledge={acknowledge}
        onComplete={complete}
      />

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-3xl">
          <DialogHeader>
            <DialogTitle>{editing ? t('edit') : t('new')}</DialogTitle>
          </DialogHeader>
          <PerformanceReviewForm initial={editing || undefined} onSubmit={onSubmit} submitting={saving} />
        </DialogContent>
      </Dialog>
    </div>
  )
}
