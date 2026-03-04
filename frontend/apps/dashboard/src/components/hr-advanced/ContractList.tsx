'use client'

import { useTranslations } from 'next-intl'
import type { EmploymentContract } from '@/types/hr-advanced'

interface Props {
  items: EmploymentContract[]
  onEdit: (item: EmploymentContract) => void
  onDelete: (id: number) => Promise<void>
  onRenew: (id: number) => Promise<void>
}

export function ContractList({ items, onEdit, onDelete, onRenew }: Props) {
  const t = useTranslations('hrAdvanced')
  return (
    <div className="overflow-x-auto rounded border bg-white dark:bg-gray-800">
      <table className="min-w-full text-sm">
        <thead className="bg-gray-100 dark:bg-gray-900/30">
          <tr>
            <th className="px-3 py-2 text-left">{t('table.employee')}</th>
            <th className="px-3 py-2 text-left">{t('table.contract')}</th>
            <th className="px-3 py-2 text-left">{t('table.dates')}</th>
            <th className="px-3 py-2 text-left">{t('table.status')}</th>
            <th className="px-3 py-2 text-right">{t('table.actions')}</th>
          </tr>
        </thead>
        <tbody>
          {items.map((item) => (
            <tr key={item.id} className="border-t">
              <td className="px-3 py-2">#{item.employee_id}</td>
              <td className="px-3 py-2">{item.contract_number}</td>
              <td className="px-3 py-2">{item.start_date} - {item.end_date || '-'}</td>
              <td className="px-3 py-2">{item.status}</td>
              <td className="px-3 py-2 text-right space-x-2">
                <button onClick={() => onEdit(item)} className="rounded border px-2 py-1">{t('actions.edit')}</button>
                <button onClick={() => void onRenew(item.id)} className="rounded border px-2 py-1">{t('actions.renew')}</button>
                <button onClick={() => void onDelete(item.id)} className="rounded border px-2 py-1 text-red-600">{t('actions.delete')}</button>
              </td>
            </tr>
          ))}
          {items.length === 0 ? <tr><td colSpan={5} className="px-3 py-6 text-center text-gray-500">{t('contracts.noData')}</td></tr> : null}
        </tbody>
      </table>
    </div>
  )
}
