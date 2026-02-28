'use client'

import { useEffect, useRef, useState, type PointerEvent as ReactPointerEvent } from 'react'
import { type RestaurantTable } from '@/lib/api/tablesApi'
import { TableStatusBadge } from './TableStatusBadge'

const statusShellStyles: Record<RestaurantTable['status'], string> = {
  available:
    'border-green-300 bg-green-50 text-green-900 dark:border-green-700 dark:bg-green-900/20 dark:text-green-100',
  occupied:
    'border-red-300 bg-red-50 text-red-900 dark:border-red-700 dark:bg-red-900/20 dark:text-red-100',
  reserved:
    'border-yellow-300 bg-yellow-50 text-yellow-900 dark:border-yellow-700 dark:bg-yellow-900/20 dark:text-yellow-100',
  maintenance:
    'border-gray-400 bg-gray-100 text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100',
}

const statusOptions: RestaurantTable['status'][] = [
  'available',
  'occupied',
  'reserved',
  'maintenance',
]

function getShapeClass(shape: RestaurantTable['shape']) {
  if (shape === 'circle') {
    return 'h-24 w-24 rounded-full'
  }
  if (shape === 'rectangle') {
    return 'h-20 w-36 rounded-lg'
  }
  return 'h-24 w-24 rounded-lg'
}

interface FloorPlanTableProps {
  table: RestaurantTable
  x: number
  y: number
  zoom: number
  canDrag: boolean
  canManageTables: boolean
  isSelected: boolean
  onMove: (tableID: number, nextX: number, nextY: number) => void
  onDoubleClick?: (tableID: number) => void
  onOpenDetails?: (tableID: number) => void
  onOpenOrders?: (table: RestaurantTable) => void
  onStatusChange?: (
    tableID: number,
    nextStatus: RestaurantTable['status'],
    currentStatus: RestaurantTable['status']
  ) => void
  onSelect: (tableID: number) => void
  t: (path: string, fallback?: string) => string
}

export function FloorPlanTable({
  table,
  x,
  y,
  zoom,
  canDrag,
  canManageTables,
  isSelected,
  onMove,
  onDoubleClick,
  onOpenDetails,
  onOpenOrders,
  onStatusChange,
  onSelect,
  t,
}: FloorPlanTableProps) {
  const [dragging, setDragging] = useState(false)
  const dragState = useRef<{
    startClientX: number
    startClientY: number
    originX: number
    originY: number
  } | null>(null)

  useEffect(() => {
    if (!dragging) {
      return
    }

    const handlePointerMove = (event: PointerEvent) => {
      if (!dragState.current) {
        return
      }

      const deltaX = (event.clientX - dragState.current.startClientX) / zoom
      const deltaY = (event.clientY - dragState.current.startClientY) / zoom
      onMove(table.id, dragState.current.originX + deltaX, dragState.current.originY + deltaY)
    }

    const stopDragging = () => {
      setDragging(false)
      dragState.current = null
    }

    window.addEventListener('pointermove', handlePointerMove)
    window.addEventListener('pointerup', stopDragging)

    return () => {
      window.removeEventListener('pointermove', handlePointerMove)
      window.removeEventListener('pointerup', stopDragging)
    }
  }, [dragging, onMove, table.id, zoom])

  const handlePointerDown = (event: ReactPointerEvent<HTMLDivElement>) => {
    if (!canDrag || event.button !== 0) {
      return
    }

    dragState.current = {
      startClientX: event.clientX,
      startClientY: event.clientY,
      originX: x,
      originY: y,
    }
    setDragging(true)
  }

  return (
    <div
      role="button"
      tabIndex={0}
      onPointerDown={handlePointerDown}
      onDoubleClick={() => onDoubleClick?.(table.id)}
      onClick={() => onSelect(table.id)}
      onKeyDown={(event) => {
        if (event.key === 'Enter' || event.key === ' ') {
          event.preventDefault()
          onSelect(table.id)
        }
      }}
      className={`absolute border p-2 text-left shadow transition ${
        getShapeClass(table.shape)
      } ${
        isSelected
          ? 'ring-2 ring-primary/70'
          : ''
      } ${
        statusShellStyles[table.status]
      } ${
        canDrag ? 'cursor-grab active:cursor-grabbing' : 'cursor-default'
      }`}
      style={{
        left: `${x}px`,
        top: `${y}px`,
      }}
    >
      <p className="text-xs font-semibold">
        #{table.table_number}
      </p>
      <p className="truncate text-[11px] opacity-90">
        {table.name || t('tables.tableName')}
      </p>
      <p className="mt-1 text-[10px] opacity-80">
        {table.capacity} {t('tables.persons')}
      </p>
      <div className="mt-1">
        <TableStatusBadge
          status={table.status}
          label={t(`tables.statuses.${table.status}`)}
        />
      </div>

      {isSelected && (
        <div
          className="absolute left-0 top-full z-30 mt-2 w-56 rounded-lg border border-gray-200 bg-white p-3 text-xs shadow-xl dark:border-gray-700 dark:bg-gray-900"
          onPointerDown={(event) => {
            event.stopPropagation()
          }}
          onClick={(event) => {
            event.stopPropagation()
          }}
        >
          <p className="text-sm font-semibold text-gray-900 dark:text-white">
            {table.name || `${t('tables.tableName')} #${table.table_number}`}
          </p>
          <p className="mt-1 text-gray-600 dark:text-gray-300">
            {t('tables.capacity')}: {table.capacity} {t('tables.persons')}
          </p>
          <p className="mt-1 text-gray-600 dark:text-gray-300">
            {t('tables.location')}: {t(`tables.locations.${table.location}`)}
          </p>
          <p className="mt-1 text-gray-600 dark:text-gray-300">
            {t('tables.status')}: {t(`tables.statuses.${table.status}`)}
          </p>

          {canManageTables && onStatusChange && (
            <div className="mt-3">
              <p className="mb-1 text-[11px] font-medium uppercase tracking-wide text-gray-500 dark:text-gray-400">
                {t('tables.status')}
              </p>
              <div className="flex flex-wrap gap-1">
                {statusOptions.map((status) => (
                  <button
                    key={status}
                    type="button"
                    className={`rounded border px-2 py-1 text-[11px] transition ${
                      table.status === status
                        ? 'border-primary bg-primary/10 text-primary'
                        : 'border-gray-300 text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-200 dark:hover:bg-gray-800'
                    }`}
                    onClick={(event) => {
                      event.stopPropagation()
                      onStatusChange(table.id, status, table.status)
                    }}
                    disabled={table.status === status}
                  >
                    {t(`tables.statuses.${status}`)}
                  </button>
                ))}
              </div>
            </div>
          )}

          <div className="mt-3 flex flex-wrap gap-2">
            {onOpenDetails && (
              <button
                type="button"
                className="rounded border border-gray-300 px-2 py-1 text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-200 dark:hover:bg-gray-800"
                onClick={(event) => {
                  event.stopPropagation()
                  onOpenDetails(table.id)
                }}
              >
                {t('tables.viewDetails')}
              </button>
            )}

            {canManageTables && onDoubleClick && (
              <button
                type="button"
                className="rounded border border-primary/40 px-2 py-1 text-primary hover:bg-primary/10"
                onClick={(event) => {
                  event.stopPropagation()
                  onDoubleClick(table.id)
                }}
              >
                {t('tables.editTable')}
              </button>
            )}

            {table.status === 'occupied' && onOpenOrders && (
              <button
                type="button"
                className="rounded border border-amber-400 px-2 py-1 text-amber-700 hover:bg-amber-50 dark:border-amber-600 dark:text-amber-300 dark:hover:bg-amber-900/30"
                onClick={(event) => {
                  event.stopPropagation()
                  onOpenOrders(table)
                }}
              >
                {t('tables.openCurrentOrder', 'Open Current Orders')}
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
