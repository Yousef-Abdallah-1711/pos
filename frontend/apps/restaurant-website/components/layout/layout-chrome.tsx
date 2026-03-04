'use client'

import { useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import type { ServerThemeData } from '@/lib/api/get-theme-server'
import { Header } from '@/components/layout/Header'
import { Footer } from '@/components/layout/Footer'
import { MobileBottomNav } from '@/components/layout/MobileBottomNav'
import { CartBottomBar } from '@/components/layout/CartBottomBar'
import { WhatsAppFloating } from '@/components/layout/WhatsAppFloating'
import { TEMPLATE_REGISTRY } from '@/lib/utils/theme-templates'

interface LayoutChromeProps {
  children: React.ReactNode
  locale: 'en' | 'ar'
  themeData: ServerThemeData
}

export function LayoutChrome({ children, locale, themeData }: LayoutChromeProps) {
  const searchParams = useSearchParams()
  const isPreview = searchParams.get('preview') === 'true'
  const templateId = searchParams.get('template') || ''
  const templateOverride = TEMPLATE_REGISTRY[templateId]
  const themedData: ServerThemeData = templateOverride
    ? {
        ...themeData,
        templateId: templateOverride.id,
        template_id: templateOverride.id,
        headerVariant: templateOverride.headerVariant,
        header_variant: templateOverride.headerVariant,
        footerVariant: templateOverride.footerVariant,
        footer_variant: templateOverride.footerVariant,
      }
    : themeData

  const header = (themedData?.header || {}) as Record<string, unknown>
  const whatsAppNumber =
    (header.whatsappNumber as string) ||
    (header.whatsapp_number as string) ||
    undefined

  useEffect(() => {
    if (!isPreview || typeof window === 'undefined') {
      return
    }

    const applyCssVars = (payload: unknown) => {
      if (!payload || typeof payload !== 'object') {
        return
      }

      const record = payload as Record<string, unknown>
      const cssVars =
        (record.cssVars as Record<string, string>) ||
        (record.variables as Record<string, string>) ||
        (record.payload as Record<string, string>) ||
        null

      if (!cssVars || typeof cssVars !== 'object') {
        return
      }

      const root = document.documentElement
      requestAnimationFrame(() => {
        Object.entries(cssVars).forEach(([key, value]) => {
          if (typeof value !== 'string') return
          const cssKey = key.startsWith('--') ? key : `--${key}`
          root.style.setProperty(cssKey, value)
        })
      })
    }

    const handler = (event: MessageEvent) => {
      const data = event.data as Record<string, unknown> | undefined
      if (!data || typeof data !== 'object') return

      const type = String(data.type || '')
      if (type === 'theme-preview-update' || type === 'theme:update' || type === 'theme-preview-css-vars') {
        applyCssVars(data)
      } else {
        applyCssVars(data)
      }
    }

    window.addEventListener('message', handler)
    return () => {
      window.removeEventListener('message', handler)
    }
  }, [isPreview])

  if (isPreview) {
    return <>{children}</>
  }

  return (
    <>
      <Header themeData={themedData} locale={locale} />
      <main className="pb-[calc(4rem+env(safe-area-inset-bottom))] md:pb-0">{children}</main>
      <Footer themeData={themedData} locale={locale} />
      <MobileBottomNav />
      <CartBottomBar />
      <WhatsAppFloating phoneNumber={whatsAppNumber} show={Boolean(whatsAppNumber)} />
    </>
  )
}
