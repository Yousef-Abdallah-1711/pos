'use client'

import { useMemo, useState } from 'react'
import { usePathname, useRouter } from 'next/navigation'
import { toast } from 'sonner'
import {
  FileText,
  LayoutGrid,
  Loader2,
  Lock,
  Plus,
  RefreshCcw,
  X,
} from 'lucide-react'
import { jsPDF } from 'jspdf'
import QRCode from 'qrcode'
import { Button } from '@/components/ui/button'
import { useHasPermission } from '@/hooks/useRbac'
import { useTables, type TableFilter } from '@/hooks/useTables'
import {
  type RestaurantTable,
  type TableStatus,
} from '@/lib/api/tablesApi'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'
import { buildTableMenuUrl } from '@/lib/table-qr'
import { AddTableForm } from './components/AddTableForm'
import { TableList } from './components/TableList'

const STATUSES: TableStatus[] = ['available', 'occupied', 'reserved', 'maintenance']

export default function TablesPage() {
  const router = useRouter()
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'

  const { hasPermission: canViewTables, loading: canViewLoading } =
    useHasPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.READ)
  const { hasPermission: canManageTables } = useHasPermission(
    MODULE_IDS.TABLES,
    PERMISSION_LEVELS.WRITE
  )
  const { hasPermission: canDeleteTables } = useHasPermission(
    MODULE_IDS.TABLES,
    PERMISSION_LEVELS.DELETE
  )

  const {
    tables,
    stats,
    statusFilter,
    setStatusFilter,
    loading,
    error,
    refresh,
    createTable,
    deleteTable,
    updateStatus,
    getNextTableNumber,
  } = useTables('all')

  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [printingAllQr, setPrintingAllQr] = useState(false)

  const statusFilters = useMemo<
    Array<{ value: TableFilter; label: string }>
  >(
    () => [
      { value: 'all', label: t('tables.all') },
      ...STATUSES.map((status) => ({
        value: status,
        label: t(`tables.statuses.${status}`),
      })),
    ],
    [t]
  )

  const handleCreateTable = async (
    payload: Parameters<typeof createTable>[0]
  ) => {
    setSubmitting(true)
    try {
      await createTable(payload)
      toast.success(t('tables.createSuccess'))
      setIsCreateModalOpen(false)
    } catch (err) {
      const message = err instanceof Error ? err.message : t('tables.saveFailed')
      toast.error(message)
      debugError('[TABLES_PAGE] Create failed:', err)
    } finally {
      setSubmitting(false)
    }
  }

  const handleEditTable = (tableID: number) => {
    router.push(`/${locale}/dashboard/tables/${tableID}/edit`)
  }

  const handleViewTable = (tableID: number) => {
    router.push(`/${locale}/dashboard/tables/${tableID}`)
  }

  const handleOpenQr = (table: RestaurantTable) => {
    router.push(`/${locale}/dashboard/tables/${table.id}/qr`)
  }

  const handleOpenPOS = (table: RestaurantTable) => {
    router.push(`/${locale}/dashboard/pos?table=${table.id}`)
  }

  const handleDeleteTable = async (tableID: number) => {
    if (!confirm(t('tables.confirmDelete'))) {
      return
    }

    try {
      await deleteTable(tableID)
      toast.success(t('tables.deleteSuccess'))
    } catch (err) {
      const message =
        err instanceof Error ? err.message : t('tables.deleteFailed')
      toast.error(message)
      debugError('[TABLES_PAGE] Delete failed:', err)
    }
  }

  const handleStatusChange = async (
    tableID: number,
    nextStatus: TableStatus,
    currentStatus: TableStatus
  ) => {
    if (nextStatus === currentStatus) {
      return
    }

    try {
      await updateStatus(tableID, nextStatus)
      toast.success(t('tables.statusUpdateSuccess'))
    } catch (err) {
      const message =
        err instanceof Error ? err.message : t('tables.statusUpdateFailed')
      toast.error(message)
      debugError('[TABLES_PAGE] Status update failed:', err)
    }
  }

  const handleFreeTable = async (tableID: number) => {
    try {
      await updateStatus(tableID, 'available')
      toast.success(t('tables.freeTableSuccess'))
    } catch (err) {
      const message =
        err instanceof Error ? err.message : t('tables.statusUpdateFailed')
      toast.error(message)
      debugError('[TABLES_PAGE] Free table failed:', err)
    }
  }

  const handlePrintAllQRCodes = async () => {
    if (tables.length === 0) {
      toast.error(t('tables.empty'))
      return
    }

    setPrintingAllQr(true)
    try {
      const doc = new jsPDF({
        orientation: 'portrait',
        unit: 'mm',
        format: 'a4',
      })

      const pageWidth = doc.internal.pageSize.getWidth()
      const pageHeight = doc.internal.pageSize.getHeight()
      const margin = 10
      const gap = 8
      const columns = 2
      const rows = 2
      const cardsPerPage = columns * rows
      const cardWidth = (pageWidth - margin * 2 - gap) / columns
      const cardHeight = (pageHeight - margin * 2 - gap) / rows
      const sortedTables = [...tables].sort((a, b) => a.table_number - b.table_number)

      for (let i = 0; i < sortedTables.length; i += 1) {
        if (i > 0 && i % cardsPerPage === 0) {
          doc.addPage()
        }

        const table = sortedTables[i]
        const indexInPage = i % cardsPerPage
        const row = Math.floor(indexInPage / columns)
        const col = indexInPage % columns
        const x = margin + col * (cardWidth + gap)
        const y = margin + row * (cardHeight + gap)

        doc.roundedRect(x, y, cardWidth, cardHeight, 3, 3, 'S')

        doc.setFontSize(12)
        doc.text(`${t('tables.tableNumber')} ${table.table_number}`, x + cardWidth / 2, y + 9, {
          align: 'center',
        })

        doc.setFontSize(9)
        doc.text(
          `${t('tables.tableName')}: ${table.name || t('common.notAvailable', 'N/A')}`,
          x + cardWidth / 2,
          y + 15,
          { align: 'center' }
        )
        doc.text(
          `${t('tables.capacity')}: ${table.capacity}`,
          x + cardWidth / 2,
          y + 20,
          { align: 'center' }
        )

        const tableUrl = buildTableMenuUrl(table.id, locale)
        const qrDataUrl = await QRCode.toDataURL(tableUrl, {
          width: 512,
          margin: 1,
          color: {
            dark: '#111827',
            light: '#FFFFFF',
          },
        })

        const qrSize = Math.min(cardWidth - 16, cardHeight - 44)
        const qrX = x + (cardWidth - qrSize) / 2
        const qrY = y + 24
        doc.addImage(qrDataUrl, 'PNG', qrX, qrY, qrSize, qrSize)
      }

      const dateLabel = new Date().toISOString().split('T')[0]
      doc.save(`tables-qr-codes-${dateLabel}.pdf`)
      toast.success(t('tables.printAllQRCodesSuccess', 'QR PDF generated successfully.'))
    } catch (error) {
      toast.error(t('tables.printAllQRCodesFailed', 'Failed to generate QR PDF.'))
      debugError('[TABLES_PAGE] Print all QR failed:', error)
    } finally {
      setPrintingAllQr(false)
    }
  }

  if (canViewLoading) {
    return (
      <div className="flex items-center justify-center h-80">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!canViewTables) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
          {t('tables.title')}
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
      <div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white sm:text-3xl">
            {t('tables.title')}
          </h1>
          <p className="mt-1 text-gray-600 dark:text-gray-400">
            {t('tables.manageTables')}
          </p>
        </div>

        <div
          className={`flex flex-wrap gap-2 ${isRTL ? 'justify-start lg:justify-end' : ''}`}
        >
          {canManageTables && (
            <Button onClick={() => setIsCreateModalOpen(true)}>
              <Plus className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
              {t('tables.addTable')}
            </Button>
          )}

          <Button
            variant="outline"
            onClick={() => {
              void refresh()
            }}
          >
            <RefreshCcw className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
            {t('tables.refresh')}
          </Button>

          <Button
            variant="outline"
            onClick={() => {
              void handlePrintAllQRCodes()
            }}
            disabled={printingAllQr}
            title={t('tables.printAllQRCodes', 'Print All QR Codes')}
          >
            {printingAllQr ? (
              <Loader2 className={isRTL ? 'ml-2 h-4 w-4 animate-spin' : 'mr-2 h-4 w-4 animate-spin'} />
            ) : (
              <FileText className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
            )}
            {t('tables.printAllQRCodes', 'Print All QR Codes')}
          </Button>

          <Button
            variant="outline"
            onClick={() => router.push(`/${locale}/dashboard/tables/layout`)}
            title={t('tables.layoutView')}
          >
            <LayoutGrid className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
            {t('tables.layoutView')}
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4 lg:grid-cols-5">
        <div className="rounded-lg bg-white p-4 shadow dark:bg-gray-800">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t('tables.stats.total')}
          </p>
          <p className="mt-1 text-2xl font-bold text-gray-900 dark:text-white">
            {stats.total_tables}
          </p>
        </div>
        <div className="rounded-lg bg-white p-4 shadow dark:bg-gray-800">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t('tables.stats.available')}
          </p>
          <p className="mt-1 text-2xl font-bold text-green-600">
            {stats.available}
          </p>
        </div>
        <div className="rounded-lg bg-white p-4 shadow dark:bg-gray-800">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t('tables.stats.occupied')}
          </p>
          <p className="mt-1 text-2xl font-bold text-red-600">
            {stats.occupied}
          </p>
        </div>
        <div className="rounded-lg bg-white p-4 shadow dark:bg-gray-800">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t('tables.stats.reserved')}
          </p>
          <p className="mt-1 text-2xl font-bold text-yellow-600">
            {stats.reserved}
          </p>
        </div>
        <div className="rounded-lg bg-white p-4 shadow dark:bg-gray-800">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t('tables.stats.maintenance')}
          </p>
          <p className="mt-1 text-2xl font-bold text-gray-700 dark:text-gray-200">
            {stats.maintenance}
          </p>
        </div>
      </div>

      <div className={`flex flex-wrap gap-2 ${isRTL ? 'justify-end' : ''}`}>
        {statusFilters.map((filter) => {
          const active = statusFilter === filter.value
          return (
            <button
              key={filter.value}
              type="button"
              onClick={() => setStatusFilter(filter.value)}
              className={`rounded-lg px-4 py-2 text-sm font-medium transition ${
                active
                  ? 'bg-primary text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600'
              }`}
            >
              {filter.label}
            </button>
          )
        })}
      </div>

      {error && (
        <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300">
          {t('tables.loadFailed')}: {error}
        </div>
      )}

      <TableList
        tables={tables}
        loading={loading}
        canViewTables={canViewTables}
        canManageTables={canManageTables}
        canDeleteTables={canDeleteTables}
        t={t}
        onEdit={(table) => handleEditTable(table.id)}
        onView={handleViewTable}
        onOpenQr={handleOpenQr}
        onOpenPOS={handleOpenPOS}
        onFreeTable={(tableID) => {
          void handleFreeTable(tableID)
        }}
        onDelete={(tableID) => {
          void handleDeleteTable(tableID)
        }}
        onStatusChange={(tableID, nextStatus, currentStatus) => {
          void handleStatusChange(tableID, nextStatus, currentStatus)
        }}
      />

      {isCreateModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 px-4">
          <div className="w-full max-w-2xl rounded-lg bg-white p-6 shadow-xl dark:bg-gray-800">
            <div className="mb-4 flex items-start justify-between">
              <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                {t('tables.addTable')}
              </h2>
              <button
                type="button"
                onClick={() => setIsCreateModalOpen(false)}
                className="rounded-md p-1 text-gray-500 transition hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-gray-700 dark:hover:text-gray-100"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            <AddTableForm
              suggestedTableNumber={getNextTableNumber()}
              submitting={submitting}
              t={t}
              isRTL={isRTL}
              onCancel={() => setIsCreateModalOpen(false)}
              onSubmit={handleCreateTable}
            />
          </div>
        </div>
      )}

    </div>
  )
}
