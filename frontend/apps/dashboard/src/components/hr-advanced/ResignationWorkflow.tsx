'use client'

import { useTranslations } from 'next-intl'
import type { Resignation } from '@/types/hr-advanced'

interface Props {
  item: Resignation
  onAcknowledge: (id: number) => Promise<void>
  onComplete: (id: number) => Promise<void>
  onCalculate: (id: number) => Promise<number>
}

export function ResignationWorkflow({ item, onAcknowledge, onComplete, onCalculate }: Props) {
  const t = useTranslations('hrAdvanced')
  return (
    <div className="rounded border bg-white p-4 dark:bg-gray-800 space-y-3">
      <div className="text-sm text-gray-600 dark:text-gray-300">{t('workflow.status')}: <span className="font-semibold">{item.status}</span></div>
      <div className="text-sm">{t('workflow.notice')}: {item.notice_period_days} {t('workflow.days')}</div>
      <div className="flex flex-wrap gap-2">
        <button onClick={() => void onAcknowledge(item.id)} className="rounded border px-3 py-1">{t('actions.acknowledge')}</button>
        <button onClick={() => void onComplete(item.id)} className="rounded border px-3 py-1">{t('actions.complete')}</button>
        <button
          onClick={async () => {
            const settlement = await onCalculate(item.id)
            alert(t('workflow.settlementAlert', { amount: settlement.toFixed(2) }))
          }}
          className="rounded border px-3 py-1"
        >
          {t('actions.calculateSettlement')}
        </button>
      </div>
    </div>
  )
}
