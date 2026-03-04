'use client'

import { useEffect } from 'react'

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  useEffect(() => {
    console.error(error)
  }, [error])

  return (
    <div className="rounded-xl border border-red-200 bg-red-50 dark:border-red-800 dark:bg-red-900/20 p-6 space-y-3">
      <h2 className="text-lg font-semibold text-red-800 dark:text-red-300">Something went wrong</h2>
      <p className="text-sm text-red-700 dark:text-red-400">An unexpected error occurred while loading this section.</p>
      <button
        type="button"
        onClick={reset}
        className="inline-flex items-center rounded-lg bg-red-600 text-white px-4 py-2 text-sm font-medium hover:bg-red-700"
      >
        Try again
      </button>
    </div>
  )
}
