'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { cn } from '@/lib/utils'

export function ThemeBuilderSubTabs() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  const tabs = [
    {
      key: 'themes',
      href: `/${locale}/dashboard/theme-builder`,
      label: t('themeBuilder.tabs.themes', 'Themes'),
      active: pathname === `/${locale}/dashboard/theme-builder`,
    },
    {
      key: 'qr',
      href: `/${locale}/dashboard/theme-builder/qr`,
      label: t('themeBuilder.tabs.qrBuilder', 'QR Builder'),
      active: pathname.startsWith(`/${locale}/dashboard/theme-builder/qr`),
    },
  ]

  return (
    <div className="mb-6 border-b border-gray-200 dark:border-gray-700">
      <nav className="-mb-px flex gap-2 overflow-x-auto">
        {tabs.map((tab) => (
          <Link
            key={tab.key}
            href={tab.href}
            className={cn(
              'inline-flex items-center rounded-t-lg border-b-2 px-4 py-2 text-sm font-medium transition-colors',
              tab.active
                ? 'border-primary text-primary'
                : 'border-transparent text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white'
            )}
          >
            {tab.label}
          </Link>
        ))}
      </nav>
    </div>
  )
}
