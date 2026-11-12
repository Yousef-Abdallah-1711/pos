'use client'

import { memo } from 'react'
import type { ServerThemeData } from '@/lib/api/get-theme-server'
import { normalizeFooterConfig, type ThemeLocale } from './layout-normalizers'
import { ExpandedFooter, CompactFooter, CardsFooter, SocialHeavyFooter, MinimalFooter } from './footers'

export interface FooterProps {
  themeData: ServerThemeData
  locale?: ThemeLocale
}

function FooterComponent({
  themeData,
  locale = 'en',
}: FooterProps) {
  const safeLocale: ThemeLocale = locale === 'ar' ? 'ar' : 'en'
  const config = normalizeFooterConfig(themeData, safeLocale)

  const rawFooter = (themeData?.footer || {}) as Record<string, unknown>
  const layoutRaw =
    (themeData?.footerVariant as string) ||
    (themeData?.footer_variant as string) ||
    (rawFooter.layout as string) ||
    (rawFooter.footerVariant as string) ||
    'expanded'
  const layout = layoutRaw.toLowerCase()

  const props = { config, locale: safeLocale }

  if (layout === 'compact') return <div data-footer-variant="compact"><CompactFooter {...props} /></div>
  if (layout === 'cards') return <div data-footer-variant="cards"><CardsFooter {...props} /></div>
  if (layout === 'social-heavy' || layout === 'socialheavy') return <div data-footer-variant="social-heavy"><SocialHeavyFooter {...props} /></div>
  if (layout === 'minimal') return <div data-footer-variant="minimal"><MinimalFooter {...props} /></div>
  return <div data-footer-variant="expanded"><ExpandedFooter {...props} /></div>
}

export const Footer = memo(FooterComponent)
