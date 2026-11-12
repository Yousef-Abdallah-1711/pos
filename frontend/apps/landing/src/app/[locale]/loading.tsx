'use client'

import { usePathname } from 'next/navigation'

export default function LandingLoading() {
  const locale = (usePathname() ?? '').split('/')[1] === 'ar' ? 'ar' : 'en'

  return (
    <main
      aria-busy="true"
      dir={locale === 'ar' ? 'rtl' : 'ltr'}
      className="mx-auto min-h-[60vh] w-full max-w-7xl space-y-8 px-4 py-10 sm:px-6 lg:px-8"
    >
      <p role="status" className="sr-only">{locale === 'ar' ? 'جاري التحميل...' : 'Loading...'}</p>
      <div aria-hidden="true" className="space-y-5">
        <div className="h-10 w-2/3 max-w-xl animate-pulse rounded-lg bg-stone-200 motion-reduce:animate-none" />
        <div className="h-5 w-full max-w-2xl animate-pulse rounded-lg bg-stone-100 motion-reduce:animate-none" />
        <div className="grid grid-cols-1 gap-5 pt-4 md:grid-cols-2">
          <div className="h-64 animate-pulse rounded-2xl bg-stone-100 motion-reduce:animate-none" />
          <div className="h-64 animate-pulse rounded-2xl bg-stone-100 motion-reduce:animate-none" />
        </div>
      </div>
    </main>
  )
}
