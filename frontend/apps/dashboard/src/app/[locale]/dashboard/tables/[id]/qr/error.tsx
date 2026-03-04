'use client'

type SegmentErrorProps = {
  error: Error & { digest?: string }
  reset: () => void
}

export default function Error({ reset }: SegmentErrorProps) {
  return (
    <div className="mx-auto max-w-2xl p-6">
      <div className="rounded-xl border border-red-200 bg-red-50 p-5 dark:border-red-900/60 dark:bg-red-950/30">
        <h2 className="text-lg font-semibold text-red-700 dark:text-red-300">Something went wrong</h2>
        <p className="mt-2 text-sm text-red-600 dark:text-red-200">Please retry. If the issue continues, check logs.</p>
        <button
          type="button"
          onClick={reset}
          className="mt-4 inline-flex rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
        >
          Try again
        </button>
      </div>
    </div>
  )
}

