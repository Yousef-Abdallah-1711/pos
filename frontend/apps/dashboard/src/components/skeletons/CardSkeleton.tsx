'use client'

export function CardSkeleton({ count = 3 }: { count?: number }) {
  return (
    <div className="card-skeleton grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {[...Array(count)].map((_, i) => (
        <div
          key={i}
          className="p-6 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 space-y-4"
        >
          {/* Header */}
          <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-3/4 animate-pulse" />

          {/* Content Lines */}
          <div className="space-y-2">
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" />
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-5/6 animate-pulse" />
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-4/6 animate-pulse" />
          </div>

          {/* Footer Button */}
          <div className="card-skeleton-footer h-10 bg-gray-200 dark:bg-gray-700 rounded w-full animate-pulse" />
        </div>
      ))}
      <style jsx>{`
        .card-skeleton-footer {
          margin-top: 1.5rem;
        }
      `}</style>
    </div>
  )
}
