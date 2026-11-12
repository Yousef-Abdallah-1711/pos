'use client'

import { memo } from 'react'
import type { ServerThemeData } from '@/lib/api/get-theme-server'
import { normalizeHeaderConfig, type ThemeLocale } from './layout-normalizers'
import { GlassHeader, GradientHeader, SolidHeader, FloatingHeader, MinimalHeader } from './headers'

export interface HeaderProps {
  themeData: ServerThemeData
  locale?: ThemeLocale
  restaurantSlug?: string
}

function HeaderComponent({
  themeData,
  locale = 'en',
  restaurantSlug,
}: HeaderProps) {
  const safeLocale: ThemeLocale = locale === 'ar' ? 'ar' : 'en'
  const config = normalizeHeaderConfig(themeData, safeLocale)

  const rawHeader = (themeData?.header || {}) as Record<string, unknown>
  const variantRaw =
    (themeData?.headerVariant as string) ||
    (themeData?.header_variant as string) ||
    (rawHeader.headerVariant as string) ||
    (rawHeader.variant as string) ||
    'solid'
  const variant = variantRaw.toLowerCase()

  const whatsAppNumber =
    (rawHeader.whatsappNumber as string) ||
    (rawHeader.whatsapp_number as string) ||
    undefined
  const showWhatsAppButton =
    typeof rawHeader.showWhatsAppButton === 'boolean'
      ? rawHeader.showWhatsAppButton
      : typeof rawHeader.show_whatsapp_button === 'boolean'
        ? rawHeader.show_whatsapp_button
        : true

  const props = {
    config,
    locale: safeLocale,
    restaurantSlug,
    whatsAppNumber,
    showWhatsAppButton,
  }

  if (variant === 'glass') return <div data-header-variant="glass"><GlassHeader {...props} /></div>
  if (variant === 'gradient') return <div data-header-variant="gradient"><GradientHeader {...props} /></div>
  if (variant === 'floating') return <div data-header-variant="floating"><FloatingHeader {...props} /></div>
  if (variant === 'minimal') return <div data-header-variant="minimal"><MinimalHeader {...props} /></div>
  return <div data-header-variant="solid"><SolidHeader {...props} /></div>
}

export const Header = memo(HeaderComponent)
