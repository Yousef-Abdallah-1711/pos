'use client'

import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

type Accent = 'green' | 'blue' | 'amber' | 'red' | 'emerald' | 'violet'

const ACCENT_CLASSES: Record<Accent, string> = {
  green: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
  blue: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
  amber: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400',
  red: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
  emerald: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400',
  violet: 'bg-violet-100 text-violet-700 dark:bg-violet-900/30 dark:text-violet-400',
}

interface KpiCardProps {
  title: string
  value: string
  icon?: ReactNode
  accent?: Accent
  subtext?: string
  className?: string
}

export function KpiCard({
  title,
  value,
  icon,
  accent = 'blue',
  subtext,
  className,
}: KpiCardProps) {
  return (
    <div
      className={cn(
        'rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-5',
        className
      )}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <div className="text-sm text-gray-600 dark:text-gray-400">{title}</div>
          <div className="mt-2 text-2xl font-bold text-gray-900 dark:text-white break-words">
            {value}
          </div>
          {subtext ? (
            <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">{subtext}</div>
          ) : null}
        </div>
        {icon ? (
          <div className={cn('p-2.5 rounded-lg shrink-0', ACCENT_CLASSES[accent])}>{icon}</div>
        ) : null}
      </div>
    </div>
  )
}

