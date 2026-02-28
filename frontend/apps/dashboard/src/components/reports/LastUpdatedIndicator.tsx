'use client'

interface LastUpdatedIndicatorProps {
  locale: string
  generatedAtLabel: string
  lastUpdatedLabel: string
  generatedAt?: string | null
  lastUpdatedAt?: Date | null
}

export function LastUpdatedIndicator({
  locale,
  generatedAtLabel,
  lastUpdatedLabel,
  generatedAt,
  lastUpdatedAt,
}: LastUpdatedIndicatorProps) {
  const localeCode = locale === 'ar' ? 'ar-EG' : 'en-US'
  return (
    <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4 flex flex-col justify-center">
      <div className="text-xs text-gray-500 dark:text-gray-400">{generatedAtLabel}</div>
      <div className="mt-1 text-sm font-medium text-gray-900 dark:text-gray-100">
        {generatedAt ? new Date(generatedAt).toLocaleString(localeCode) : '—'}
      </div>
      <div className="mt-3 text-xs text-gray-500 dark:text-gray-400">{lastUpdatedLabel}</div>
      <div className="mt-1 text-sm font-medium text-gray-900 dark:text-gray-100">
        {lastUpdatedAt ? lastUpdatedAt.toLocaleString(localeCode) : '—'}
      </div>
    </div>
  )
}

