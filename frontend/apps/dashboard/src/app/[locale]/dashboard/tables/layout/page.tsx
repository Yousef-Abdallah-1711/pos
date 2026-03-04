'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import { usePathname, useRouter } from 'next/navigation'
import { ArrowLeft, Loader2, Lock, Plus } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useUserPermissions } from '@/hooks/useRbac'
import { useTables } from '@/hooks/useTables'
import {
  type RestaurantTable,
  type TableLayout,
  type TablePositionPayload,
  tablesApi,
} from '@/lib/api/tablesApi'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'
import { FloorPlanEditor } from '../components/FloorPlanEditor'
import { FloorPlanErrorBoundary } from '../components/FloorPlanErrorBoundary'

type LayoutTab = 'all' | 'unassigned' | number

const defaultLayoutForm = {
  name: '',
  floorNumber: '1',
  canvasWidth: '980',
  canvasHeight: '560',
}

export default function TablesLayoutPage() {
  const pathname = usePathname()
  const router = useRouter()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canViewTables = checkPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.READ)
  const canManageTables = checkPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.WRITE)

  const { tables, loading, refresh } = useTables('all')

  const [saving, setSaving] = useState(false)
  const [layouts, setLayouts] = useState<TableLayout[]>([])
  const [layoutsLoading, setLayoutsLoading] = useState(true)
  const [layoutsError, setLayoutsError] = useState<string | null>(null)
  const [selectedLayout, setSelectedLayout] = useState<LayoutTab>('all')
  const [isCreateLayoutOpen, setIsCreateLayoutOpen] = useState(false)
  const [creatingLayout, setCreatingLayout] = useState(false)
  const [layoutForm, setLayoutForm] = useState(defaultLayoutForm)

  const loadLayouts = useCallback(async () => {
    setLayoutsLoading(true)
    try {
      const data = await tablesApi.getLayouts()
      setLayouts(data)
      setLayoutsError(null)

      setSelectedLayout((prev) => {
        if (prev === 'all' || prev === 'unassigned') {
          return prev
        }

        if (typeof prev === 'number' && data.some((layout) => layout.id === prev)) {
          return prev
        }

        if (data.length === 0) {
          return 'all'
        }

        const defaultLayout = data.find((layout) => layout.is_default)
        return defaultLayout ? defaultLayout.id : data[0].id
      })
    } catch (error) {
      const message =
        error instanceof Error ? error.message : 'Failed to load layouts'
      setLayoutsError(message)
      debugError('[TABLE_LAYOUT] Failed to load layouts:', error)
    } finally {
      setLayoutsLoading(false)
    }
  }, [])

  useEffect(() => {
    if (permissionLoading) {
      return
    }

    if (!canViewTables) {
      setLayoutsLoading(false)
      return
    }

    void loadLayouts()
  }, [permissionLoading, canViewTables, loadLayouts])

  const selectedLayoutData = useMemo(() => {
    if (typeof selectedLayout !== 'number') {
      return null
    }
    return layouts.find((layout) => layout.id === selectedLayout) || null
  }, [layouts, selectedLayout])

  const unassignedTablesCount = useMemo(() => {
    return tables.filter((table) => !table.layout_id).length
  }, [tables])

  const tableCountByLayout = useMemo(() => {
    return tables.reduce<Record<number, number>>((acc, table) => {
      if (table.layout_id) {
        acc[table.layout_id] = (acc[table.layout_id] || 0) + 1
      }
      return acc
    }, {})
  }, [tables])

  const layoutTabs = useMemo(() => {
    const baseTabs: Array<{ value: LayoutTab; label: string; count: number }> = [
      { value: 'all', label: t('tables.all'), count: tables.length },
    ]

    for (const layout of layouts) {
      baseTabs.push({
        value: layout.id,
        label: layout.name,
        count: tableCountByLayout[layout.id] || 0,
      })
    }

    baseTabs.push({
      value: 'unassigned',
      label: t('tables.unassignedLayout', 'Unassigned'),
      count: unassignedTablesCount,
    })

    return baseTabs
  }, [layouts, tableCountByLayout, tables.length, t, unassignedTablesCount])

  const filteredTables = useMemo(() => {
    if (selectedLayout === 'all') {
      return tables
    }

    if (selectedLayout === 'unassigned') {
      return tables.filter((table) => !table.layout_id)
    }

    return tables.filter((table) => table.layout_id === selectedLayout)
  }, [selectedLayout, tables])

  const handleSavePositions = async (positions: TablePositionPayload[]) => {
    setSaving(true)
    try {
      await tablesApi.updatePositions(positions)
      toast.success(t('tables.positionsSaved'))
      await refresh()
    } catch (error) {
      const message =
        error instanceof Error ? error.message : t('tables.positionsSaveFailed')
      toast.error(message)
      debugError('[FLOOR_PLAN] Failed to save positions:', error)
    } finally {
      setSaving(false)
    }
  }

  const handleStatusChange = async (
    tableID: number,
    nextStatus: RestaurantTable['status'],
    currentStatus: RestaurantTable['status']
  ) => {
    if (nextStatus === currentStatus) {
      return
    }

    try {
      await tablesApi.updateStatus(tableID, nextStatus)
      toast.success(t('tables.statusUpdateSuccess'))
      await refresh()
    } catch (error) {
      const message =
        error instanceof Error ? error.message : t('tables.statusUpdateFailed')
      toast.error(message)
      debugError('[TABLE_LAYOUT] Failed to update table status:', error)
    }
  }

  const handleCreateLayout = async () => {
    const name = layoutForm.name.trim()
    if (!name) {
      toast.error(t('tables.layoutNameRequired', 'Layout name is required.'))
      return
    }

    const floorNumber = Number.parseInt(layoutForm.floorNumber, 10)
    const canvasWidth = Number.parseInt(layoutForm.canvasWidth, 10)
    const canvasHeight = Number.parseInt(layoutForm.canvasHeight, 10)

    if (
      !Number.isFinite(floorNumber) ||
      !Number.isFinite(canvasWidth) ||
      !Number.isFinite(canvasHeight) ||
      floorNumber < 1 ||
      canvasWidth < 500 ||
      canvasHeight < 300
    ) {
      toast.error(
        t('tables.invalidLayoutFields', 'Please provide valid floor and canvas values.')
      )
      return
    }

    setCreatingLayout(true)
    try {
      const created = await tablesApi.createLayout({
        name,
        floor_number: floorNumber,
        canvas_width: canvasWidth,
        canvas_height: canvasHeight,
      })
      toast.success(t('tables.layoutCreated', 'Layout created successfully.'))
      setLayoutForm(defaultLayoutForm)
      setIsCreateLayoutOpen(false)
      await loadLayouts()
      setSelectedLayout(created.id)
    } catch (error) {
      const message =
        error instanceof Error ? error.message : t('tables.layoutCreateFailed', 'Failed to create layout.')
      toast.error(message)
      debugError('[TABLE_LAYOUT] Failed to create layout:', error)
    } finally {
      setCreatingLayout(false)
    }
  }

  if (permissionLoading || loading || layoutsLoading) {
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
          {t('tables.floorPlan')}
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
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              {t('tables.floorPlan')}
            </h1>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {t('tables.floorPlanSubTitle')}
            </p>
          </div>
        </div>
      </div>

      <div className="rounded-lg border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
        <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
          <div className={`flex flex-wrap gap-2 ${isRTL ? 'justify-end' : ''}`}>
            {layoutTabs.map((tab) => {
              const active = selectedLayout === tab.value
              return (
                <button
                  key={String(tab.value)}
                  type="button"
                  onClick={() => setSelectedLayout(tab.value)}
                  className={`rounded-lg border px-3 py-2 text-sm transition ${
                    active
                      ? 'border-primary bg-primary text-white'
                      : 'border-gray-200 bg-gray-50 text-gray-700 hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-700'
                  }`}
                >
                  {tab.label} ({tab.count})
                </button>
              )
            })}
          </div>

          {canManageTables && (
            <Button variant="outline" onClick={() => setIsCreateLayoutOpen(true)}>
              <Plus className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
              {t('tables.newLayout', 'New Layout')}
            </Button>
          )}
        </div>

        <p className="mt-3 text-xs text-gray-500 dark:text-gray-400">
          {selectedLayoutData
            ? `${selectedLayoutData.name} | ${t('tables.floorNumber', 'Floor')} ${selectedLayoutData.floor_number} | ${Math.round(selectedLayoutData.canvas_width)}x${Math.round(selectedLayoutData.canvas_height)}`
            : selectedLayout === 'unassigned'
            ? t('tables.unassignedLayoutHint', 'Tables without a layout assignment.')
            : t('tables.layoutAllHint', 'Showing tables across all layouts.')}
        </p>

        {layoutsError && (
          <p className="mt-2 text-xs text-red-600 dark:text-red-400">
            {layoutsError}
          </p>
        )}
      </div>

      <FloorPlanErrorBoundary t={t}>
        <FloorPlanEditor
          tables={filteredTables}
          canManageTables={canManageTables}
          saving={saving}
          canvasWidth={selectedLayoutData?.canvas_width}
          canvasHeight={selectedLayoutData?.canvas_height}
          onEditTable={(tableID) => router.push(`/${locale}/dashboard/tables/${tableID}/edit`)}
          onOpenTableDetails={(tableID) => router.push(`/${locale}/dashboard/tables/${tableID}`)}
          onOpenOrders={(table) => router.push(`/${locale}/dashboard/orders?table_id=${table.id}`)}
          onStatusChange={(tableID, nextStatus, currentStatus) => {
            void handleStatusChange(tableID, nextStatus, currentStatus)
          }}
          onSave={handleSavePositions}
          t={t}
          isRTL={isRTL}
        />
      </FloorPlanErrorBoundary>

      <Dialog
        open={isCreateLayoutOpen}
        onOpenChange={(open) => {
          if (creatingLayout) {
            return
          }
          setIsCreateLayoutOpen(open)
        }}
      >
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{t('tables.createLayout', 'Create Layout')}</DialogTitle>
          </DialogHeader>

          <form
            className="space-y-4"
            onSubmit={(event) => {
              event.preventDefault()
              void handleCreateLayout()
            }}
          >
            <div className="space-y-2">
              <Label htmlFor="layout-name">{t('tables.layoutName', 'Layout Name')}</Label>
              <Input
                id="layout-name"
                value={layoutForm.name}
                onChange={(event) =>
                  setLayoutForm((prev) => ({ ...prev, name: event.target.value }))
                }
                placeholder={t('tables.layoutNamePlaceholder', 'Main Floor')}
                disabled={creatingLayout}
              />
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
              <div className="space-y-2">
                <Label htmlFor="layout-floor">{t('tables.floorNumber', 'Floor')}</Label>
                <Input
                  id="layout-floor"
                  type="number"
                  min={1}
                  value={layoutForm.floorNumber}
                  onChange={(event) =>
                    setLayoutForm((prev) => ({ ...prev, floorNumber: event.target.value }))
                  }
                  disabled={creatingLayout}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="layout-width">{t('tables.canvasWidth', 'Canvas Width')}</Label>
                <Input
                  id="layout-width"
                  type="number"
                  min={500}
                  value={layoutForm.canvasWidth}
                  onChange={(event) =>
                    setLayoutForm((prev) => ({ ...prev, canvasWidth: event.target.value }))
                  }
                  disabled={creatingLayout}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="layout-height">{t('tables.canvasHeight', 'Canvas Height')}</Label>
                <Input
                  id="layout-height"
                  type="number"
                  min={300}
                  value={layoutForm.canvasHeight}
                  onChange={(event) =>
                    setLayoutForm((prev) => ({ ...prev, canvasHeight: event.target.value }))
                  }
                  disabled={creatingLayout}
                />
              </div>
            </div>

            <div className="flex justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => setIsCreateLayoutOpen(false)}
                disabled={creatingLayout}
              >
                {t('tables.cancel')}
              </Button>
              <Button type="submit" disabled={creatingLayout}>
                {creatingLayout ? t('common.saving') : t('tables.createLayout', 'Create Layout')}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
