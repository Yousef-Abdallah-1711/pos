'use client';

import { usePathname } from 'next/navigation';

export default function StorefrontLoading() {
  const pathname = usePathname() ?? '';
  const locale = pathname.split('/')[1] === 'ar' ? 'ar' : 'en';

  return (
    <main
      aria-busy="true"
      dir={locale === 'ar' ? 'rtl' : 'ltr'}
      className="mx-auto min-h-[60vh] w-full max-w-7xl space-y-8 px-4 py-8 sm:px-6 lg:px-8"
    >
      <p role="status" className="sr-only">
        {locale === 'ar' ? 'جاري التحميل...' : 'Loading...'}
      </p>
      <div aria-hidden="true" className="space-y-4">
        <div className="h-9 w-56 max-w-full animate-pulse rounded-lg bg-gray-200 motion-reduce:animate-none dark:bg-gray-700" />
        <div className="h-5 w-80 max-w-full animate-pulse rounded-lg bg-gray-100 motion-reduce:animate-none dark:bg-gray-800" />
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }, (_, index) => (
            <div key={index} className="h-64 animate-pulse rounded-2xl bg-gray-100 motion-reduce:animate-none dark:bg-gray-800" />
          ))}
        </div>
      </div>
    </main>
  );
}
