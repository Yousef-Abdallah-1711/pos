'use client'

import { memo } from 'react'
import { Header } from '@pos-saas/ui-themes'
import type { ServerThemeData } from '@/lib/api/get-theme-server'
import { normalizeHeaderConfig, type ThemeLocale } from './theme-layout-normalizers'

export interface ThemeHeaderProps {
  themeData: ServerThemeData
  locale?: ThemeLocale
  onNavigationClick?: (href: string) => void
}

function ThemeHeaderComponent({
  themeData,
  locale = 'en',
  onNavigationClick,
}: ThemeHeaderProps) {
  const safeLocale: ThemeLocale = locale === 'ar' ? 'ar' : 'en'
  const headerConfig = normalizeHeaderConfig(themeData, safeLocale)

  return (
    <Header
      config={headerConfig}
      isArabic={safeLocale === 'ar'}
      onNavigationClick={onNavigationClick}
    />
  )
}

export const ThemeHeader = memo(ThemeHeaderComponent)
