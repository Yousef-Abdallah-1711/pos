'use client'

import { useState } from 'react'
import { useTranslations } from 'next-intl'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { ContractForm } from '@/components/hr-advanced/ContractForm'
import { ContractList } from '@/components/hr-advanced/ContractList'
import { useContracts } from '@/hooks/useContracts'
import type { EmploymentContract } from '@/types/hr-advanced'

export default function ContractsPage() {
  const t = useTranslations('hrAdvanced.contracts')
  const { items, expiring, loading, error, createItem, updateItem, removeItem, renewItem } = useContracts()
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<EmploymentContract | null>(null)
  const [saving, setSaving] = useState(false)

  const onSubmit = async (payload: Partial<EmploymentContract>) => {
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
        <button className="rounded bg-primary-600 px-4 py-2 text-white" onClick={() => setOpen(true)}>{t('add')}</button>
      </div>
      {expiring.length > 0 ? (
        <div className="rounded border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
          {t('expiringBanner', { count: expiring.length })}
        </div>
      ) : null}
      {loading ? <div className="rounded border p-4 text-sm text-gray-500">{t('loading')}</div> : null}
      {error ? <div className="rounded border border-red-200 bg-red-50 p-4 text-sm text-red-700">{error}</div> : null}
      <ContractList
        items={items}
        onEdit={(item) => {
          setEditing(item)
          setOpen(true)
        }}
        onDelete={removeItem}
        onRenew={(id) =>
          renewItem(id, {
            start_date: new Date().toISOString().slice(0, 10),
            salary_amount: items.find((item) => item.id === id)?.salary_amount || 0,
            currency: items.find((item) => item.id === id)?.currency || 'EGP',
          })
        }
      />

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-3xl">
          <DialogHeader>
            <DialogTitle>{editing ? t('edit') : t('new')}</DialogTitle>
          </DialogHeader>
          <ContractForm initial={editing || undefined} onSubmit={onSubmit} submitting={saving} />
        </DialogContent>
      </Dialog>
    </div>
  )
}
