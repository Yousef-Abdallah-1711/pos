'use client'

import type { RestaurantTable } from '@/lib/api/tablesApi'
import { selectPOSTableId, usePOSStore } from '@/stores/posStore'

interface POSTableSelectorProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  tables: RestaurantTable[]
}

export function POSTableSelector({ locale, t, tables }: POSTableSelectorProps) {
  const tableId = usePOSStore(selectPOSTableId)
  const setTableContext = usePOSStore((state) => state.setTableContext)
  const isRTL = locale === 'ar'
  const selectableTables = tables.filter((table) => table.status !== 'maintenance')

  return (
    <div className="space-y-1">
      <label className={`block text-xs font-semibold text-gray-600 dark:text-gray-300 ${isRTL ? 'text-right' : ''}`}>
        {t('pos.selectTable', 'Table')}
      </label>
      <select
        value={tableId || ''}
        onChange={(event) => {
          const nextID = event.target.value ? Number(event.target.value) : null
          const selectedTable = tables.find((table) => table.id === nextID)
          setTableContext(nextID, selectedTable?.name || null)
        }}
        className={`h-9 w-full rounded border border-gray-300 bg-white px-2 text-sm text-gray-900 outline-none transition focus:border-primary focus:ring-2 focus:ring-primary/20 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-100 ${
          isRTL ? 'text-right' : 'text-left'
        }`}
      >
        <option value="">{t('pos.noTable', 'No table (Takeaway)')}</option>
        {selectableTables.map((table) => (
          <option key={table.id} value={table.id}>
            #{table.table_number} - {table.name || `Table ${table.table_number}`} ({table.capacity}{' '}
            {t('pos.seats', 'seats')}) - {t(`tables.statuses.${table.status}`, table.status)}
          </option>
        ))}
      </select>
    </div>
  )
}

export default POSTableSelector
