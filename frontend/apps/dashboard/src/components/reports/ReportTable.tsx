'use client'

import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

type CellValue = ReactNode

interface ReportTableProps {
  title?: string
  columns: string[]
  rows: CellValue[][]
  emptyText: string
  className?: string
  compact?: boolean
}

export function ReportTable({
  title,
  columns,
  rows,
  emptyText,
  className,
  compact = false,
}: ReportTableProps) {
  const cellY = compact ? 'py-2' : 'py-3'
  return (
    <div
      className={cn(
        'rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-5',
        className
      )}
    >
      {title ? <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">{title}</h2> : null}
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-gray-200 dark:border-gray-700 text-left">
              {columns.map((col, colIndex) => (
                <th
                  key={`col-${colIndex}-${col}`}
                  className="py-2 pr-3 text-gray-600 dark:text-gray-400 font-medium whitespace-nowrap"
                >
                  {col}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="py-4 text-gray-500 dark:text-gray-400">
                  {emptyText}
                </td>
              </tr>
            ) : (
              rows.map((row, rowIndex) => (
                <tr
                  key={`row-${rowIndex}`}
                  className="border-b border-gray-100 dark:border-gray-800 last:border-b-0"
                >
                  {row.map((cell, cellIndex) => (
                    <td
                      key={`row-${rowIndex}-cell-${cellIndex}`}
                      className={cn(
                        cellY,
                        'pr-3',
                        cellIndex === 0
                          ? 'text-gray-900 dark:text-gray-100'
                          : 'text-gray-700 dark:text-gray-300',
                        cellIndex === row.length - 1 ? 'font-medium text-gray-900 dark:text-gray-100' : undefined
                      )}
                    >
                      {cell}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
