'use client'

import { Skeleton } from '@/components/ui/skeleton'

export function ProductCardSkeleton() {
  return (
    <div className="overflow-hidden rounded-xl border border-[hsl(var(--theme-border))] bg-[hsl(var(--theme-surface))]">
      <Skeleton className="h-44 w-full" />
      <div className="space-y-3 p-4">
        <Skeleton className="h-5 w-2/3" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-4/5" />
        <div className="flex items-center justify-between">
          <Skeleton className="h-8 w-24 rounded-full" />
          <Skeleton className="h-8 w-28" />
        </div>
        <Skeleton className="h-11 w-full" />
      </div>
    </div>
  )
}

