'use client'

import Link from 'next/link'
import type { ThemeLocale } from '@/types/theme'
import { readLocalized, toRecord, toStringValue } from '@/components/sections/section-utils'
import { GradientHero } from '@/components/sections/heroes/GradientHero'

interface VideoHeroProps {
  locale: ThemeLocale
  config?: Record<string, unknown>
  restaurantSlug?: string
}

export function VideoHero({ locale, config = {}, restaurantSlug }: VideoHeroProps) {
  const settings = toRecord(config)
  const videoUrl = toStringValue(settings.video_url ?? settings.videoUrl, '')
  const headline = readLocalized(settings.headline ?? settings.title, locale, 'Fresh food, fast delivery')
  const subtitle = readLocalized(settings.subtitle ?? settings.description, locale, 'Order your favorite dishes in minutes.')
  const menuHref = `/${locale}/menu`

  if (!videoUrl) {
    return <GradientHero locale={locale} config={settings} restaurantSlug={restaurantSlug} />
  }

  return (
    <section className="relative isolate flex min-h-[100svh] items-center overflow-hidden">
      <video
        className="absolute inset-0 -z-20 h-full w-full object-cover"
        autoPlay
        muted
        loop
        playsInline
      >
        <source src={videoUrl} />
      </video>
      <div className="absolute inset-0 -z-10 bg-gradient-to-b from-black/50 via-black/55 to-black/70" />
      <div className="mx-auto flex min-h-[100svh] w-full max-w-6xl items-center px-4 py-20 sm:px-6 lg:px-8">
        <div className={`max-w-3xl space-y-6 ${locale === 'ar' ? 'text-right' : 'text-left'}`}>
          <h1 className="text-4xl font-bold leading-tight text-white sm:text-5xl lg:text-6xl">
            {headline}
          </h1>
          <p className="max-w-2xl text-lg text-white/90 sm:text-xl">{subtitle}</p>
          <div className={`flex ${locale === 'ar' ? 'justify-end' : 'justify-start'}`}>
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
