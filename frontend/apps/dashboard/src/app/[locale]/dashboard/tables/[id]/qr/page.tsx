'use client'

import { useEffect, useMemo, useState } from 'react'
import { useParams, usePathname, useRouter } from 'next/navigation'
import { Loader2, Lock } from 'lucide-react'
import { toast } from 'sonner'
import { useUserPermissions } from '@/hooks/useRbac'
import { type RestaurantTable, tablesApi } from '@/lib/api/tablesApi'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'
import { QRStudio } from '../../components/QRStudio'

export default function TableQrStudioPage() {
  const params = useParams<{ id: string }>()
  const router = useRouter()
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])
  const isRTL = locale === 'ar'
  const tableID = Number.parseInt(String(params.id || ''), 10)

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canViewTables = checkPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.READ)

  const [table, setTable] = useState<RestaurantTable | null>(null)
  const [loading, setLoading] = useState(true)

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
        debugError('[TABLE_QR_STUDIO] Failed to load table:', error)
      } finally {
        setLoading(false)
      }
    }

    if (!permissionLoading && canViewTables) {
      void fetchTable()
    }
  }, [permissionLoading, canViewTables, tableID, t])

  if (permissionLoading || loading) {
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
          {t('tables.qrBuilder')}
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

  if (!table) {
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300">
        {t('tables.loadFailed')}
      </div>
    )
  }

  return (
    <QRStudio
      table={table}
      locale={locale}
      t={t}
      isRTL={isRTL}
      onBack={() => router.push(`/${locale}/dashboard/tables`)}
    />
  )
}
