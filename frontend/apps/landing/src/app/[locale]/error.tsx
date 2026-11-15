'use client'

import { usePathname } from 'next/navigation'
import Link from 'next/link'
import { AlertCircle } from 'lucide-react'
import english from '@/i18n/messages/en/errors.json'
import arabic from '@/i18n/messages/ar/errors.json'

interface LandingErrorProps {
  error: Error & { digest?: string }
  reset: () => void
}

export default function LandingError({ error, reset }: LandingErrorProps) {
  const pathname = usePathname() ?? ''
  const locale = pathname.split('/')[1] === 'ar' ? 'ar' : 'en'
  const messages = locale === 'ar' ? arabic : english

  return (
    <main dir={locale === 'ar' ? 'rtl' : 'ltr'} className="flex min-h-[60vh] items-center justify-center bg-stone-50 px-4 py-12">
      <section role="alert" className="w-full max-w-xl space-y-4 rounded-2xl border border-stone-200 bg-white p-6 text-center shadow-sm sm:p-8">
        <AlertCircle aria-hidden="true" className="mx-auto text-red-700" size={30} />
        <h1 className="text-2xl font-semibold text-stone-950">
          {messages.pageUnavailable}
        </h1>
        <p className="text-sm text-stone-700">
          {messages.retry}
        </p>
        <div className="flex flex-col justify-center gap-3 sm:flex-row">
          <button
            type="button"
            onClick={reset}
            className="min-h-12 rounded-lg bg-stone-900 px-5 py-3 font-medium text-white hover:bg-stone-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
          >
            {messages.retry}
          </button>
          <Link
            href={`/${locale}`}
            className="inline-flex min-h-12 items-center justify-center rounded-lg border border-stone-300 px-5 py-3 font-medium text-stone-900 hover:bg-stone-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
          >
            {messages.backHome}
          </Link>
        </div>
      </section>
    </main>
  )
}
