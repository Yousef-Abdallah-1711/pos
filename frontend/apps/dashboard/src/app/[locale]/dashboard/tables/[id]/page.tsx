'use client'

import { useEffect, useState } from 'react'
import { useParams, usePathname, useRouter } from 'next/navigation'
import { ArrowLeft, Loader2, Lock, PencilLine, QrCode } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { useHasPermission } from '@/hooks/useRbac'
import {
  type RestaurantTable,
  type TableStatus,
  tablesApi,
} from '@/lib/api/tablesApi'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'
import { TableStatusBadge } from '../components/TableStatusBadge'

const STATUS_ACTIONS: TableStatus[] = [
  'available',
  'occupied',
  'reserved',
  'maintenance',
]

export default function TableDetailsPage() {
  const params = useParams<{ id: string }>()
  const router = useRouter()
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'

  const tableID = Number.parseInt(String(params.id || ''), 10)

  const { hasPermission: canViewTables, loading: canViewLoading } =
    useHasPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.READ)
  const { hasPermission: canManageTables, loading: canManageLoading } =
    useHasPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.WRITE)

  const [table, setTable] = useState<RestaurantTable | null>(null)
  const [loading, setLoading] = useState(true)
  const [statusLoading, setStatusLoading] = useState(false)

  const loadTable = async () => {
    if (!Number.isFinite(tableID) || tableID <= 0) {
      setLoading(false)
      return
    }

    setLoading(true)
    try {
      const data = await tablesApi.get(tableID)
      setTable(data)
    } catch (error) {
      toast.error(t('tables.loadFailed'))
      debugError('[TABLE_DETAILS] Failed to load table:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (!canViewLoading && canViewTables) {
      void loadTable()
    }
  }, [canViewLoading, canViewTables, tableID])

  const handleStatusAction = async (nextStatus: TableStatus) => {
    if (!table || nextStatus === table.status) {
      return
    }

    setStatusLoading(true)
    try {
      await tablesApi.updateStatus(table.id, nextStatus)
      toast.success(t('tables.statusUpdateSuccess'))
      await loadTable()
    } catch (error) {
      const message =
        error instanceof Error ? error.message : t('tables.statusUpdateFailed')
      toast.error(message)
      debugError('[TABLE_DETAILS] Failed to update status:', error)
    } finally {
      setStatusLoading(false)
    }
  }

  const handleFreeTable = async () => {
    if (!table || table.status === 'available') {
      return
    }

    setStatusLoading(true)
    try {
      await tablesApi.updateStatus(table.id, 'available')
      toast.success(t('tables.freeTableSuccess'))
      await loadTable()
    } catch (error) {
      const message =
        error instanceof Error ? error.message : t('tables.statusUpdateFailed')
      toast.error(message)
      debugError('[TABLE_DETAILS] Failed to free table:', error)
    } finally {
      setStatusLoading(false)
    }
  }

  if (canViewLoading || canManageLoading || loading) {
    return (
      <div className="flex h-80 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!canViewTables) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
          {t('tables.tableDetails')}
        </h1>

        <div className="rounded-lg border border-red-200 bg-red-50 p-6 dark:border-red-800 dark:bg-red-900/20">
          <div className="flex items-start gap-3">
            <Lock className="mt-0.5 h-5 w-5 text-red-600 dark:text-red-400" />
            <div>
              <p className="font-semibold text-red-900 dark:text-red-200">
                {t('tables.accessDeniedTitle')}
              </p>
              <p className="mt-1 text-sm text-red-800 dark:text-red-300">
                {t('tables.accessDeniedText')}
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!table) {
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300">
        {t('tables.loadFailed')}
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            onClick={() => router.push(`/${locale}/dashboard/tables`)}
          >
            <ArrowLeft className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
            {t('tables.backToTables')}
          </Button>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {t('tables.tableDetails')} #{table.table_number}
          </h1>
        </div>

        <div className="flex flex-wrap gap-2">
          <Button
            variant="outline"
            onClick={() => router.push(`/${locale}/dashboard/tables/${table.id}/qr`)}
            title={t('tables.qrCode')}
          >
            <QrCode className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
            {t('tables.qrCode')}
          </Button>
          {canManageTables && (
            <Button
              variant="outline"
              onClick={() => router.push(`/${locale}/dashboard/tables/${table.id}/edit`)}
            >
              <PencilLine className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
              {t('tables.editTable')}
            </Button>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        <div className="rounded-lg bg-white p-4 shadow dark:bg-gray-800">
          <p className="text-sm text-gray-500 dark:text-gray-400">{t('tables.tableName')}</p>
          <p className="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
            {table.name || '-'}
          </p>
        </div>
        <div className="rounded-lg bg-white p-4 shadow dark:bg-gray-800">
          <p className="text-sm text-gray-500 dark:text-gray-400">{t('tables.capacity')}</p>
          <p className="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
            {table.capacity} {t('tables.persons')}
          </p>
        </div>
        <div className="rounded-lg bg-white p-4 shadow dark:bg-gray-800">
          <p className="text-sm text-gray-500 dark:text-gray-400">{t('tables.location')}</p>
          <p className="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
            {t(`tables.locations.${table.location}`)}
          </p>
        </div>
        <div className="rounded-lg bg-white p-4 shadow dark:bg-gray-800">
          <p className="text-sm text-gray-500 dark:text-gray-400">{t('tables.status')}</p>
          <div className="mt-1">
            <TableStatusBadge
              status={table.status}
              label={t(`tables.statuses.${table.status}`)}
            />
          </div>
        </div>
      </div>

      <div className="rounded-lg bg-white p-6 shadow dark:bg-gray-800">
        <h2 className="mb-3 text-lg font-semibold text-gray-900 dark:text-white">
          {t('tables.quickActions')}
        </h2>
        {canManageTables ? (
          <div className="flex flex-wrap gap-2">
            {STATUS_ACTIONS.map((status) => (
              <Button
                key={status}
                variant={table.status === status ? 'default' : 'outline'}
                disabled={statusLoading || table.status === status}
                onClick={() => {
                  void handleStatusAction(status)
                }}
              >
                {t(`tables.statuses.${status}`)}
              </Button>
            ))}
            {table.status === 'occupied' && (
              <Button
                variant="outline"
                disabled={statusLoading}
                onClick={() => {
                  void handleFreeTable()
                }}
              >
                {t('tables.freeTable')}
              </Button>
            )}
          </div>
        ) : (
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t('tables.readOnlyStatusHint')}
          </p>
        )}
      </div>

      <div className="rounded-lg bg-white p-6 shadow dark:bg-gray-800">
        <h2 className="mb-3 text-lg font-semibold text-gray-900 dark:text-white">
          {t('tables.description')}
        </h2>
        <p className="text-sm text-gray-700 dark:text-gray-300">
          {table.description || t('tables.noDescription')}
        </p>
      </div>
    </div>
  )
}
