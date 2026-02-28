'use client'

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { Maximize2, MoveDiagonal2, RotateCcw, Save, ZoomIn, ZoomOut } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  type RestaurantTable,
  type TablePositionPayload,
} from '@/lib/api/tablesApi'
import { FloorPlanTable } from './FloorPlanTable'

interface FloorPlanEditorProps {
  tables: RestaurantTable[]
  canManageTables: boolean
  saving: boolean
  canvasWidth?: number
  canvasHeight?: number
  onEditTable?: (tableID: number) => void
  onOpenTableDetails?: (tableID: number) => void
  onOpenOrders?: (table: RestaurantTable) => void
  onStatusChange?: (
    tableID: number,
    nextStatus: RestaurantTable['status'],
    currentStatus: RestaurantTable['status']
  ) => void
  onSave: (positions: TablePositionPayload[]) => Promise<void> | void
  t: (path: string, fallback?: string) => string
  isRTL: boolean
}

const DEFAULT_CANVAS_WIDTH = 980
const DEFAULT_CANVAS_HEIGHT = 560
const GRID_SIZE = 24
const MIN_COORD = 0
const TABLE_WIDTH_BUFFER = 150
const TABLE_HEIGHT_BUFFER = 110
const MIN_ZOOM = 0.5
const MAX_ZOOM = 2
const ZOOM_STEP = 0.1

type PositionMap = Record<number, { x: number; y: number }>

function clamp(value: number, min: number, max: number) {
  if (value < min) {
    return min
  }
  if (value > max) {
    return max
  }
  return value
}

function buildPositions(tables: RestaurantTable[]): PositionMap {
  return tables.reduce<PositionMap>((acc, table) => {
    acc[table.id] = {
      x: Number.isFinite(table.position_x) ? table.position_x : 0,
      y: Number.isFinite(table.position_y) ? table.position_y : 0,
    }
    return acc
  }, {})
}

export function FloorPlanEditor({
  tables,
  canManageTables,
  saving,
  canvasWidth,
  canvasHeight,
  onEditTable,
  onOpenTableDetails,
  onOpenOrders,
  onStatusChange,
  onSave,
  t,
  isRTL,
}: FloorPlanEditorProps) {
  const [positions, setPositions] = useState<PositionMap>(() => buildPositions(tables))
  const [selectedTableID, setSelectedTableID] = useState<number | null>(null)
  const [snapToGrid, setSnapToGrid] = useState(true)
  const [zoom, setZoom] = useState(1)
  const viewportRef = useRef<HTMLDivElement | null>(null)
  const effectiveCanvasWidth = Math.max(
    DEFAULT_CANVAS_WIDTH,
    Math.round(canvasWidth || DEFAULT_CANVAS_WIDTH)
  )
  const effectiveCanvasHeight = Math.max(
    DEFAULT_CANVAS_HEIGHT,
    Math.round(canvasHeight || DEFAULT_CANVAS_HEIGHT)
  )
  const maxX = Math.max(MIN_COORD, effectiveCanvasWidth - TABLE_WIDTH_BUFFER)
  const maxY = Math.max(MIN_COORD, effectiveCanvasHeight - TABLE_HEIGHT_BUFFER)

  useEffect(() => {
    setPositions(buildPositions(tables))
    setSelectedTableID((prev) => {
      if (prev === null) {
        return prev
      }
      return tables.some((table) => table.id === prev) ? prev : null
    })
  }, [tables])

  const selectedTable = useMemo(() => {
    if (!selectedTableID) {
      return null
    }
    return tables.find((table) => table.id === selectedTableID) || null
  }, [selectedTableID, tables])

  const fitToScreen = useCallback(() => {
    const viewport = viewportRef.current
    if (!viewport) {
      return
    }

    const horizontalRatio = (viewport.clientWidth - 32) / effectiveCanvasWidth
    const verticalRatio = (viewport.clientHeight - 32) / effectiveCanvasHeight
    const fittedZoom = clamp(Math.min(horizontalRatio, verticalRatio), MIN_ZOOM, MAX_ZOOM)
    setZoom(fittedZoom)
  }, [effectiveCanvasHeight, effectiveCanvasWidth])

  const moveSelected = (deltaX: number, deltaY: number) => {
    if (!selectedTableID || !canManageTables) {
      return
    }

    setPositions((prev) => {
      const current = prev[selectedTableID]
      if (!current) {
        return prev
      }

      return {
        ...prev,
        [selectedTableID]: {
          x: clamp(current.x + deltaX, MIN_COORD, maxX),
          y: clamp(current.y + deltaY, MIN_COORD, maxY),
        },
      }
    })
  }

  const handleMoveTable = (tableID: number, nextX: number, nextY: number) => {
    if (!canManageTables) {
      return
    }

    const normalizeCoordinate = (value: number) => {
      if (!snapToGrid) {
        return value
      }
      return Math.round(value / GRID_SIZE) * GRID_SIZE
    }

    setPositions((prev) => ({
      ...prev,
      [tableID]: {
        x: clamp(normalizeCoordinate(nextX), MIN_COORD, maxX),
        y: clamp(normalizeCoordinate(nextY), MIN_COORD, maxY),
      },
    }))
  }

  useEffect(() => {
    fitToScreen()
  }, [fitToScreen])

  useEffect(() => {
    const handleResize = () => {
      fitToScreen()
    }

    window.addEventListener('resize', handleResize)
    return () => {
      window.removeEventListener('resize', handleResize)
    }
  }, [fitToScreen])

  const handleReset = () => {
    setPositions(buildPositions(tables))
  }

  const handleSave = async () => {
    const payload: TablePositionPayload[] = tables.map((table) => {
      const position = positions[table.id] || { x: 0, y: 0 }
      return {
        table_id: table.id,
        position_x: position.x,
        position_y: position.y,
      }
    })

    await onSave(payload)
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            {t('tables.floorPlan')}
          </h2>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t('tables.floorPlanHint')}
          </p>
        </div>

        <div className={`flex flex-wrap gap-2 ${isRTL ? 'justify-start' : ''}`}>
          <Button
            variant="outline"
            onClick={() => setZoom((prev) => clamp(prev - ZOOM_STEP, MIN_ZOOM, MAX_ZOOM))}
          >
            <ZoomOut className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
            {t('tables.zoomOut', 'Zoom Out')}
          </Button>
          <Button
            variant="outline"
            onClick={() => setZoom((prev) => clamp(prev + ZOOM_STEP, MIN_ZOOM, MAX_ZOOM))}
          >
            <ZoomIn className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
            {t('tables.zoomIn', 'Zoom In')}
          </Button>
          <Button variant="outline" onClick={fitToScreen}>
            <Maximize2 className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
            {t('tables.fitToScreen', 'Fit')}
          </Button>
          <span className="inline-flex items-center rounded-md border border-gray-200 px-3 text-sm text-gray-600 dark:border-gray-700 dark:text-gray-300">
            {Math.round(zoom * 100)}%
          </span>
          <Button
            variant={snapToGrid ? 'default' : 'outline'}
            onClick={() => setSnapToGrid((prev) => !prev)}
            disabled={!canManageTables}
          >
            {t('tables.gridSnap', 'Grid Snap')}
          </Button>
          <Button variant="outline" onClick={handleReset} disabled={saving}>
            <RotateCcw className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
            {t('tables.resetPositions')}
          </Button>
          <Button onClick={() => void handleSave()} disabled={!canManageTables || saving}>
            <Save className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
            {t('tables.savePositions')}
          </Button>
        </div>
      </div>

      <div className="rounded-lg border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
        <div className="mb-4 flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300">
          <MoveDiagonal2 className="h-4 w-4" />
          {canManageTables ? t('tables.layoutInstruction') : t('tables.readOnlyLayoutHint')}
        </div>
        {canManageTables && (
          <p className="mb-4 text-xs text-gray-500 dark:text-gray-400">
            {t('tables.doubleClickToEdit', 'Double-click a table to open edit view.')}
          </p>
        )}
        <div
          ref={viewportRef}
          className="relative overflow-auto rounded-lg border border-dashed border-gray-300 bg-gray-50 dark:border-gray-700 dark:bg-gray-900"
          style={{
            width: '100%',
            height: '60vh',
            minHeight: '420px',
            maxHeight: '700px',
          }}
        >
          <div
            className="relative"
            style={{
              width: `${effectiveCanvasWidth * zoom}px`,
              height: `${effectiveCanvasHeight * zoom}px`,
            }}
          >
            <div
              className="relative"
              style={{
                width: `${effectiveCanvasWidth}px`,
                height: `${effectiveCanvasHeight}px`,
                transform: `scale(${zoom})`,
                transformOrigin: 'top left',
                backgroundImage:
                  'linear-gradient(to right, rgba(148,163,184,0.15) 1px, transparent 1px), linear-gradient(to bottom, rgba(148,163,184,0.15) 1px, transparent 1px)',
                backgroundSize: '24px 24px',
              }}
            >
              {tables.map((table) => {
                const position = positions[table.id] || { x: 0, y: 0 }
                return (
                  <FloorPlanTable
                    key={table.id}
                    table={table}
                    x={position.x}
                    y={position.y}
                    zoom={zoom}
                    canDrag={canManageTables}
                    canManageTables={canManageTables}
                    isSelected={table.id === selectedTableID}
                    onMove={handleMoveTable}
                    onDoubleClick={onEditTable}
                    onOpenDetails={onOpenTableDetails}
                    onOpenOrders={onOpenOrders}
                    onStatusChange={onStatusChange}
                    onSelect={setSelectedTableID}
                    t={t}
                  />
                )
              })}
            </div>
          </div>
        </div>
      </div>

      <div className="rounded-lg bg-white p-4 shadow dark:bg-gray-800">
        <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
          {t('tables.selectedTable')}
        </h3>
        {selectedTable ? (
          <div className="mt-3 space-y-3">
            <p className="text-sm text-gray-700 dark:text-gray-300">
              #{selectedTable.table_number} - {selectedTable.name || '-'}
            </p>
            <div className="flex flex-wrap gap-2">
              <Button
                variant="outline"
                size="sm"
                disabled={!canManageTables}
                onClick={() => moveSelected(-20, 0)}
              >
                {t('tables.moveLeft')}
              </Button>
              <Button
                variant="outline"
                size="sm"
                disabled={!canManageTables}
                onClick={() => moveSelected(20, 0)}
              >
                {t('tables.moveRight')}
              </Button>
              <Button
                variant="outline"
                size="sm"
                disabled={!canManageTables}
                onClick={() => moveSelected(0, -20)}
              >
                {t('tables.moveUp')}
              </Button>
              <Button
                variant="outline"
                size="sm"
                disabled={!canManageTables}
                onClick={() => moveSelected(0, 20)}
              >
                {t('tables.moveDown')}
              </Button>
            </div>
          </div>
        ) : (
          <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
            {t('tables.selectTableHint')}
          </p>
        )}
      </div>
    </div>
  )
}
