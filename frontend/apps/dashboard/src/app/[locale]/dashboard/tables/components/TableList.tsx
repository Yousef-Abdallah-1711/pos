'use client'

import { Edit, Eye, Loader2, QrCode, ShoppingCart, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  type RestaurantTable,
  type TableStatus,
} from '@/lib/api/tablesApi'
import { TableStatusBadge } from './TableStatusBadge'

const STATUS_OPTIONS: TableStatus[] = [
  'available',
  'occupied',
  'reserved',
  'maintenance',
]

interface TableListProps {
  tables: RestaurantTable[]
  loading: boolean
  canManageTables: boolean
  canDeleteTables: boolean
  canViewTables: boolean
  onEdit: (table: RestaurantTable) => void
  onDelete: (tableID: number) => void
  onView: (tableID: number) => void
  onOpenQr: (table: RestaurantTable) => void
  onOpenPOS: (table: RestaurantTable) => void
  onFreeTable: (tableID: number) => void
  onStatusChange: (
    tableID: number,
    nextStatus: TableStatus,
    currentStatus: TableStatus
  ) => void
  t: (path: string, fallback?: string) => string
}

export function TableList({
  tables,
  loading,
  canManageTables,
  canDeleteTables,
  canViewTables,
  onEdit,
  onDelete,
  onView,
  onOpenQr,
  onOpenPOS,
  onFreeTable,
  onStatusChange,
  t,
}: TableListProps) {
  return (
    <div className="overflow-hidden rounded-lg bg-white shadow dark:bg-gray-800">
      <div className="overflow-x-auto">
        <table className="w-full min-w-[860px]">
          <thead className="bg-gray-50 dark:bg-gray-700">
            <tr>
              <th className="px-5 py-3 text-left text-xs font-semibold uppercase text-gray-600 dark:text-gray-300">
                {t('tables.tableNumber')}
              </th>
              <th className="px-5 py-3 text-left text-xs font-semibold uppercase text-gray-600 dark:text-gray-300">
                {t('tables.tableName')}
              </th>
              <th className="px-5 py-3 text-left text-xs font-semibold uppercase text-gray-600 dark:text-gray-300">
                {t('tables.capacity')}
              </th>
              <th className="px-5 py-3 text-left text-xs font-semibold uppercase text-gray-600 dark:text-gray-300">
                {t('tables.location')}
              </th>
              <th className="px-5 py-3 text-left text-xs font-semibold uppercase text-gray-600 dark:text-gray-300">
                {t('tables.shape')}
              </th>
              <th className="px-5 py-3 text-left text-xs font-semibold uppercase text-gray-600 dark:text-gray-300">
                {t('tables.status')}
              </th>
              <th className="px-5 py-3 text-left text-xs font-semibold uppercase text-gray-600 dark:text-gray-300">
                {t('tables.actions')}
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
            {loading ? (
              <tr>
                <td colSpan={7} className="px-6 py-12 text-center">
                  <Loader2 className="mx-auto h-6 w-6 animate-spin text-primary" />
                </td>
              </tr>
            ) : tables.length === 0 ? (
              <tr>
                <td
                  colSpan={7}
                  className="px-6 py-10 text-center text-gray-500 dark:text-gray-400"
                >
                  {t('tables.empty')}
                </td>
              </tr>
            ) : (
              tables.map((table) => (
                <tr key={table.id} className="hover:bg-gray-50 dark:hover:bg-gray-750">
                  <td className="px-5 py-4 font-semibold text-gray-900 dark:text-white">
                    #{table.table_number}
                  </td>
                  <td className="px-5 py-4 text-gray-900 dark:text-white">
                    {table.name || '-'}
                  </td>
                  <td className="px-5 py-4 text-gray-700 dark:text-gray-300">
                    {table.capacity} {t('tables.persons')}
                  </td>
                  <td className="px-5 py-4 text-gray-700 dark:text-gray-300">
                    {t(`tables.locations.${table.location}`)}
                  </td>
                  <td className="px-5 py-4 text-gray-700 dark:text-gray-300">
                    {t(`tables.shapes.${table.shape}`)}
                  </td>
                  <td className="px-5 py-4">
                    {canManageTables ? (
                      <select
                        value={table.status}
                        onChange={(event) => {
                          onStatusChange(
                            table.id,
                            event.target.value as TableStatus,
                            table.status
                          )
                        }}
                        className="rounded-md border border-gray-300 bg-white px-2 py-1 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary dark:border-gray-600 dark:bg-gray-700 dark:text-white"
                      >
                        {STATUS_OPTIONS.map((status) => (
                          <option key={status} value={status}>
                            {t(`tables.statuses.${status}`)}
                          </option>
                        ))}
                      </select>
                    ) : (
                      <TableStatusBadge
                        status={table.status}
                        label={t(`tables.statuses.${table.status}`)}
                      />
                    )}
                  </td>
                  <td className="px-5 py-4">
                    <div className="flex flex-wrap gap-2">
                      {canViewTables && (
                        <Button
                          size="sm"
                          variant="outline"
                          title={t('tables.viewDetails')}
                          onClick={() => onView(table.id)}
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                      )}

                      {canViewTables && (
                        <Button
                          size="sm"
                          variant="outline"
                          title={t('tables.qrCode')}
                          onClick={() => onOpenQr(table)}
                        >
                          <QrCode className="h-4 w-4" />
                        </Button>
                      )}

                      {canManageTables && (
                        <Button
                          size="sm"
                          variant="outline"
                          title={t('tables.assignInPOS', 'Assign in POS')}
                          onClick={() => onOpenPOS(table)}
                        >
                          <ShoppingCart className="h-4 w-4" />
                        </Button>
                      )}

                      {canManageTables && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => onEdit(table)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                      )}
                      {canDeleteTables && (
                        <Button
                          size="sm"
                          variant="destructive"
                          onClick={() => onDelete(table.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      )}

                      {canManageTables && table.status === 'occupied' && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => onFreeTable(table.id)}
                        >
                          {t('tables.freeTable')}
                        </Button>
                      )}
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
