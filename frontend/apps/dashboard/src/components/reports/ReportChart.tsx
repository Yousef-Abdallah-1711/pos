'use client'

import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

interface ReportChartProps {
  title: string
  children: ReactNode
  className?: string
  bodyClassName?: string
  action?: ReactNode
}

export function ReportChart({ title, children, className, bodyClassName, action }: ReportChartProps) {
  return (
    <div
      className={cn(
        'rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-5',
        className
      )}
    >
      <div className="mb-4 flex items-center justify-between gap-3">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{title}</h2>
        {action}
      </div>
      <div className={bodyClassName}>{children}</div>
    </div>
  )
}

