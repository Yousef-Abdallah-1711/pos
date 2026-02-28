'use client'

import { useEffect, useMemo, useState } from 'react'
import { useParams, usePathname, useRouter } from 'next/navigation'
import { ArrowLeft, Loader2, Lock } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { useHasPermission } from '@/hooks/useRbac'
import { AddTableForm } from '../../components/AddTableForm'
import {
  type RestaurantTable,
  tablesApi,
} from '@/lib/api/tablesApi'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'

export default function EditTablePage() {
  const params = useParams<{ id: string }>()
  const router = useRouter()
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])
  const isRTL = locale === 'ar'

  const tableID = Number.parseInt(String(params.id || ''), 10)

  const { hasPermission: canViewTables, loading: canViewLoading } =
    useHasPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.READ)
  const { hasPermission: canManageTables, loading: canManageLoading } =
    useHasPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.WRITE)

  const [table, setTable] = useState<RestaurantTable | null>(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    const fetchTable = async () => {
      if (!Number.isFinite(tableID) || tableID <= 0) {
        setLoading(false)
        return
      }

      setLoading(true)
      try {
        const tableData = await tablesApi.get(tableID)
        setTable(tableData)
      } catch (error) {
        toast.error(t('tables.loadFailed'))
        debugError('[TABLE_EDIT] Failed to load table:', error)
      } finally {
        setLoading(false)
      }
    }

    if (!canViewLoading && canViewTables) {
      void fetchTable()
    }
  }, [canViewLoading, canViewTables, tableID, t])

  const handleUpdateTable = async (
    payload: Parameters<typeof tablesApi.update>[1]
  ) => {
    if (!Number.isFinite(tableID) || tableID <= 0) {
      toast.error(t('tables.invalidTableNumber'))
      return
    }

    setSubmitting(true)
    try {
      await tablesApi.update(tableID, payload)
      toast.success(t('tables.updateSuccess'))
      router.push(`/${locale}/dashboard/tables`)
    } catch (error) {
      const message =
        error instanceof Error ? error.message : t('tables.saveFailed')
      toast.error(message)
      debugError('[TABLE_EDIT] Failed to update table:', error)
    } finally {
      setSubmitting(false)
    }
  }

  if (canViewLoading || canManageLoading) {
    return (
      <div className="flex items-center justify-center h-80">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!canViewTables || !canManageTables) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
          {t('tables.editTable')}
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

  if (!Number.isFinite(tableID) || tableID <= 0) {
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300">
        {t('tables.invalidTableNumber')}
      </div>
    )
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-80">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
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
      <div className="flex items-center gap-3">
        <Button
          variant="outline"
          onClick={() => router.push(`/${locale}/dashboard/tables`)}
        >
          <ArrowLeft className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
          {t('tables.backToTables')}
        </Button>
      </div>

      <div className="rounded-lg bg-white p-6 shadow dark:bg-gray-800">
        <h1 className="mb-4 text-2xl font-bold text-gray-900 dark:text-white">
          {t('tables.editTable')} #{table.table_number}
        </h1>
        <AddTableForm
          initialValues={table}
          submitting={submitting}
          showCancel={false}
          submitLabel={t('tables.update')}
          t={t}
          isRTL={isRTL}
          onSubmit={handleUpdateTable}
        />
      </div>
    </div>
  )
}
