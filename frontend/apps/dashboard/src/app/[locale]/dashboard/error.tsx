'use client'

import { usePathname } from 'next/navigation'
import { AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'

interface DashboardErrorProps {
  error: Error & { digest?: string }
  reset: () => void
}

export default function DashboardError({ reset }: DashboardErrorProps) {
  const pathname = usePathname() ?? ''
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  return (
    <main dir={locale === 'ar' ? 'rtl' : 'ltr'} className="flex min-h-[40vh] items-center justify-center p-4 sm:p-6">
      <section role="alert" className="w-full max-w-xl space-y-4 rounded-xl border border-red-200 bg-white p-6 text-center shadow-sm dark:border-red-900 dark:bg-gray-900">
        <AlertCircle aria-hidden="true" className="mx-auto text-red-700 dark:text-red-300" size={28} />
        <h1 className="text-xl font-semibold text-gray-900 dark:text-gray-100">{t('common.sectionLoadErrorTitle')}</h1>
        <p className="text-sm text-gray-700 dark:text-gray-300">{t('common.sectionLoadError')}</p>
        <Button type="button" onClick={reset} className="min-h-12 min-w-28 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2">
          {t('common.tryAgain')}
        </Button>
      </section>
    </main>
  )
}
