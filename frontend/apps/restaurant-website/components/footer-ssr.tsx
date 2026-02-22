'use client'

import { memo } from 'react'
import type { ServerThemeData } from '@/lib/api/get-theme-server'
import { ThemeFooter } from '@/components/layout'

interface FooterSSRProps {
  themeData: ServerThemeData
  locale?: 'en' | 'ar'
}

function FooterSSRComponent({ themeData, locale = 'en' }: FooterSSRProps) {
  return <ThemeFooter themeData={themeData} locale={locale} />
}

export const FooterSSR = memo(FooterSSRComponent)
