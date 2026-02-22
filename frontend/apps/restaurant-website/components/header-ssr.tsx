'use client'

import { memo } from 'react'
import type { ServerThemeData } from '@/lib/api/get-theme-server'
import { ThemeHeader } from '@/components/layout'

interface HeaderSSRProps {
  themeData: ServerThemeData
  cartItemsCount?: number
  locale?: 'en' | 'ar'
}

function HeaderSSRComponent({ themeData, locale = 'en' }: HeaderSSRProps) {
  return <ThemeHeader themeData={themeData} locale={locale} />
}

export const HeaderSSR = memo(HeaderSSRComponent)
