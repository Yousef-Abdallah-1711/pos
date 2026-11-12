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
        <Button type="button" variant="outline" className="h-12 min-w-20 gap-2 px-3">
          <Languages className="h-4 w-4" />
          <span className="text-xs font-semibold uppercase">{currentLocale}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => setLocale('en')}>EN</DropdownMenuItem>
        <DropdownMenuItem onClick={() => setLocale('ar')}>AR</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
