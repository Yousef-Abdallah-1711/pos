'use client'

import { usePathname } from 'next/navigation'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'

export default function DashboardLoading() {
  const pathname = usePathname() ?? ''
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  return (
    <main aria-busy="true" dir={locale === 'ar' ? 'rtl' : 'ltr'} className="space-y-6 p-4 sm:p-6">
      <p role="status" className="sr-only">{t('common.loading')}</p>
      <div aria-hidden="true" className="space-y-6">
        <div className="space-y-2">
          <div className="h-8 w-52 max-w-full animate-pulse rounded-md bg-gray-200 motion-reduce:animate-none dark:bg-gray-700" />
          <div className="h-4 w-72 max-w-full animate-pulse rounded-md bg-gray-200 motion-reduce:animate-none dark:bg-gray-700" />
        </div>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {Array.from({ length: 4 }, (_, index) => (
            <div key={index} className="h-28 animate-pulse rounded-xl bg-gray-100 motion-reduce:animate-none dark:bg-gray-800" />
          ))}
        </div>
        <div className="grid grid-cols-1 gap-4 xl:grid-cols-2">
          <div className="h-72 animate-pulse rounded-xl bg-gray-100 motion-reduce:animate-none dark:bg-gray-800" />
          <div className="h-72 animate-pulse rounded-xl bg-gray-100 motion-reduce:animate-none dark:bg-gray-800" />
        </div>
      </div>
    </main>
  )
}
