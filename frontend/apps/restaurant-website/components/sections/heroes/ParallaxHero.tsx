'use client'

import Link from 'next/link'
import { motion, useScroll, useTransform } from 'framer-motion'
import type { ThemeLocale } from '@/types/theme'
import { readLocalized, toRecord, toStringValue } from '@/components/sections/section-utils'

interface ParallaxHeroProps {
  locale: ThemeLocale
  config?: Record<string, unknown>
}

export function ParallaxHero({ locale, config = {} }: ParallaxHeroProps) {
  const settings = toRecord(config)
  const headline = readLocalized(settings.headline ?? settings.title, locale, 'Crafted flavors for every table')
  const subtitle = readLocalized(settings.subtitle ?? settings.description, locale, 'Explore signature dishes made fresh every day.')
  const backgroundImage = toStringValue(settings.background_image ?? settings.backgroundImage, '')
  const { scrollYProgress } = useScroll()
  const translateY = useTransform(scrollYProgress, [0, 1], ['0%', '20%'])

  return (
    <section className="relative isolate min-h-[100svh] overflow-hidden bg-black">
      <motion.div
        className="absolute inset-0 -z-20"
        style={{ y: translateY }}
      >
        <div
          className="h-[120%] w-full bg-cover bg-center"
          style={{
            backgroundImage: backgroundImage
              ? `url(${backgroundImage})`
              : 'linear-gradient(135deg, hsl(var(--theme-primary)) 0%, hsl(var(--theme-accent)) 100%)',
          }}
        />
      </motion.div>
      <div className="absolute inset-0 -z-10 bg-black/55 dark:bg-black/65" />

      <div className="mx-auto flex min-h-[100svh] w-full max-w-6xl items-center px-4 py-20 sm:px-6 lg:px-8">
        <div className={`max-w-3xl space-y-6 ${locale === 'ar' ? 'text-right' : 'text-left'}`}>
          <h1 className="text-4xl font-bold text-white sm:text-5xl lg:text-6xl">{headline}</h1>
          <p className="text-lg text-white/85 sm:text-xl">{subtitle}</p>
          <div className={`flex ${locale === 'ar' ? 'justify-end' : 'justify-start'}`}>
            <Link href={`/${locale}/menu`} className="inline-flex min-h-12 items-center rounded-md bg-white px-6 py-3 text-sm font-semibold text-black">
              {locale === 'ar' ? 'اطلب الآن' : 'Order Now'}
            </Link>
          </div>
        </div>
      </div>
    </section>
  )
}

