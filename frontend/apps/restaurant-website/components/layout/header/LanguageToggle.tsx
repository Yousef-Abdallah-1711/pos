'use client'

import { Languages } from 'lucide-react'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { usePreferencesStore } from '@/lib/store/preferences-store'

function switchLocaleInPath(pathname: string, locale: 'en' | 'ar') {
  if (!pathname) return `/${locale}`
  if (/^\/(en|ar)(\/|$)/.test(pathname)) {
    return pathname.replace(/^\/(en|ar)/, `/${locale}`)
  }
  return `/${locale}${pathname.startsWith('/') ? '' : '/'}${pathname}`
}

export function LanguageToggle() {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const updateLanguage = usePreferencesStore((state) => state.updateLanguage)

  const currentLocale = pathname?.startsWith('/ar') ? 'ar' : 'en'

  const setLocale = (nextLocale: 'en' | 'ar') => {
    const targetPath = switchLocaleInPath(pathname || '/', nextLocale)
    const query = searchParams?.toString()
    updateLanguage(nextLocale)
    router.push(query ? `${targetPath}?${query}` : targetPath)
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button type="button" variant="outline" aria-label={currentLocale === 'ar' ? 'تغيير اللغة' : 'Change language'} className="h-12 w-12 gap-2 px-2 sm:w-auto sm:min-w-20 sm:px-3">
          <Languages className="h-4 w-4" />
          <span className="hidden text-xs font-semibold uppercase sm:inline">{currentLocale}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => setLocale('en')}>EN</DropdownMenuItem>
        <DropdownMenuItem onClick={() => setLocale('ar')}>AR</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
