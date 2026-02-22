'use client'

export function TableSkeleton({ rows = 5 }: { rows?: number }) {
  const delayRules = [0, 1, 2, 3, 4]
    .map((colIndex) => `.col-delay-${colIndex} { animation-delay: ${colIndex * 0.1}s; }`)
    .join('\n')

  return (
    <div className="table-skeleton space-y-3">
      {/* Table Header Skeleton */}
      <div className="grid grid-cols-5 gap-4 px-4 py-3 bg-gray-100 dark:bg-gray-700 rounded-lg">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="h-4 bg-gray-300 dark:bg-gray-600 rounded animate-pulse" />
        ))}
      </div>

      {/* Table Rows Skeleton */}
      {[...Array(rows)].map((_, rowIndex) => (
        <div key={rowIndex} className="grid grid-cols-5 gap-4 px-4 py-3 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
          {[...Array(5)].map((_, colIndex) => (
            <div
              key={colIndex}
              className={`h-4 bg-gray-200 dark:bg-gray-700 rounded animate-pulse col-delay-${colIndex}`}
            />
          ))}
        </div>
      ))}
      <style jsx>{`
        ${delayRules}
      `}</style>
    </div>
  )
}
