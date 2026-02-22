'use client'

import { memo } from 'react'
import { Footer } from '@pos-saas/ui-themes'
import type { ServerThemeData } from '@/lib/api/get-theme-server'
import { normalizeFooterConfig, type ThemeLocale } from './theme-layout-normalizers'

export interface ThemeFooterProps {
  themeData: ServerThemeData
  locale?: ThemeLocale
  onLinkClick?: (href: string) => void
}

function ThemeFooterComponent({
  themeData,
  locale = 'en',
  onLinkClick,
}: ThemeFooterProps) {
  const safeLocale: ThemeLocale = locale === 'ar' ? 'ar' : 'en'
  const footerConfig = normalizeFooterConfig(themeData, safeLocale)

  return (
    <Footer
      config={footerConfig}
      isArabic={safeLocale === 'ar'}
      onLinkClick={onLinkClick}
    />
  )
}

export const ThemeFooter = memo(ThemeFooterComponent)
