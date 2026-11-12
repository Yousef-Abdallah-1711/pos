'use client'

import Link from 'next/link'
import { Truck } from 'lucide-react'
import type { ThemeLocale } from '@/types/theme'
import { readLocalized, toRecord, toStringValue } from '@/components/sections/section-utils'

interface GradientHeroProps {
  locale: ThemeLocale
  config?: Record<string, unknown>
  restaurantSlug?: string
}

export function GradientHero({ locale, config = {}, restaurantSlug }: GradientHeroProps) {
  const settings = toRecord(config)
  const headline = readLocalized(settings.headline ?? settings.title, locale, 'Fresh food, fast delivery')
  const subtitle = readLocalized(settings.subtitle ?? settings.description, locale, 'Order your favorite dishes in minutes.')
  const badge = readLocalized(settings.badge, locale, locale === 'ar' ? 'توصيل سريع' : 'Fast Delivery')
  const whatsappPhone = toStringValue(settings.whatsapp_number ?? settings.whatsappNumber, '')
  const menuHref = `/${locale}/menu`
  const waText = encodeURIComponent(locale === 'ar' ? 'مرحبا، أريد الطلب من المطعم' : 'Hello, I want to place an order')
  const waHref = whatsappPhone ? `https://wa.me/${whatsappPhone}?text=${waText}` : menuHref
  const backgroundImage = toStringValue(settings.background_image ?? settings.backgroundImage, '')

  return (
    <section className="relative isolate flex min-h-[100svh] items-center overflow-hidden bg-[linear-gradient(140deg,hsl(var(--theme-primary))_0%,hsl(var(--theme-accent))_100%)]">
      {backgroundImage ? (
        <div
          className="absolute inset-0 -z-20 bg-cover bg-center opacity-25"
          style={{ backgroundImage: `url(${backgroundImage})` }}
        />
      ) : null}
      <div className="absolute inset-0 -z-10 bg-black/45 dark:bg-black/60" />
      <div className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 lg:px-8">
        <div className={`max-w-3xl space-y-6 ${locale === 'ar' ? 'text-right' : 'text-left'}`}>
          <span className="inline-flex items-center gap-2 rounded-full bg-white/15 px-4 py-2 text-sm font-medium text-white backdrop-blur-sm">
            <Truck className="h-4 w-4" />
            {badge}
          </span>
          <h1 className="text-4xl font-bold leading-tight text-white sm:text-5xl lg:text-6xl">
            {headline}
          </h1>
          <p className="max-w-2xl text-lg text-white/90 sm:text-xl">{subtitle}</p>
          <div className={`flex flex-wrap items-center gap-3 ${locale === 'ar' ? 'justify-end' : 'justify-start'}`}>
            <Link
              href={waHref}
              target={whatsappPhone ? '_blank' : undefined}
              rel={whatsappPhone ? 'noreferrer noopener' : undefined}
              className="inline-flex min-h-12 items-center justify-center rounded-md bg-[#25D366] px-6 py-3 text-sm font-semibold text-white shadow-sm transition hover:bg-[#1faa52]"
            >
              {locale === 'ar' ? 'اطلب واتساب' : 'Order via WhatsApp'}
            </Link>
            <Link
              href={menuHref}
              className="inline-flex min-h-12 items-center justify-center rounded-md bg-white px-6 py-3 text-sm font-semibold text-black transition hover:bg-white/90"
            >
              {locale === 'ar' ? 'اطلب الآن' : 'Order Now'}
            </Link>
          </div>
          {restaurantSlug ? (
            <p className="text-sm text-white/70">{locale === 'ar' ? `مطعم ${restaurantSlug}` : restaurantSlug}</p>
          ) : null}
        </div>
      </div>
    </section>
  )
}

