'use client'

import { useEffect, useMemo, useState } from 'react'
import useEmblaCarousel from 'embla-carousel-react'
import Autoplay from 'embla-carousel-autoplay'
import Link from 'next/link'
import type { ThemeLocale } from '@/types/theme'
import { readLocalized, toArray, toRecord } from '@/components/sections/section-utils'

interface SliderHeroProps {
  locale: ThemeLocale
  config?: Record<string, unknown>
}

export function SliderHero({ locale, config = {} }: SliderHeroProps) {
  const settings = toRecord(config)
  const slides = useMemo(() => {
    const images = toArray<string>(settings.images)
    if (images.length > 0) {
      return images
    }
    const fallback = settings.background_image || settings.backgroundImage
    return typeof fallback === 'string' && fallback.length > 0 ? [fallback] : []
  }, [settings])
  const headline = readLocalized(settings.headline ?? settings.title, locale, 'Taste the best dishes in town')
  const subtitle = readLocalized(settings.subtitle ?? settings.description, locale, 'Fresh ingredients, authentic flavor, fast service.')
  const [selectedIndex, setSelectedIndex] = useState(0)
  const [emblaRef, emblaApi] = useEmblaCarousel({ loop: true, align: 'start', direction: locale === 'ar' ? 'rtl' : 'ltr' }, [Autoplay({ delay: 5000 })])

  useEffect(() => {
    if (!emblaApi) return
    const onSelect = () => setSelectedIndex(emblaApi.selectedScrollSnap())
    onSelect()
    emblaApi.on('select', onSelect)
    return () => {
      emblaApi.off('select', onSelect)
    }
  }, [emblaApi])

  return (
    <section className="relative isolate min-h-[100svh] overflow-hidden bg-black">
      {slides.length > 0 ? (
        <div className="absolute inset-0 -z-20 overflow-hidden" ref={emblaRef}>
          <div className="flex h-full">
            {slides.map((slide, index) => (
              <div key={`${slide}-${index}`} className="relative min-w-full">
                <div
                  className="h-full min-h-[100svh] w-full bg-cover bg-center"
                  style={{ backgroundImage: `url(${slide})` }}
                />
              </div>
            ))}
          </div>
        </div>
      ) : (
        <div className="absolute inset-0 -z-20 bg-[linear-gradient(130deg,hsl(var(--theme-primary))_0%,hsl(var(--theme-accent))_100%)]" />
      )}

      <div className="absolute inset-0 -z-10 bg-black/50" />
      <div className="mx-auto flex min-h-[100svh] w-full max-w-6xl items-center px-4 py-20 sm:px-6 lg:px-8">
        <div className={`max-w-3xl space-y-6 ${locale === 'ar' ? 'text-right' : 'text-left'}`}>
          <h1 className="text-4xl font-bold text-white sm:text-5xl lg:text-6xl">{headline}</h1>
          <p className="text-lg text-white/85 sm:text-xl">{subtitle}</p>
          <div className={`flex gap-3 ${locale === 'ar' ? 'justify-end' : 'justify-start'}`}>
            <Link href={`/${locale}/menu`} className="inline-flex min-h-12 items-center rounded-md bg-white px-6 py-3 text-sm font-semibold text-black">
              {locale === 'ar' ? 'اطلب الآن' : 'Order Now'}
            </Link>
          </div>
        </div>
      </div>

      {slides.length > 1 ? (
        <div className="absolute bottom-6 left-1/2 z-10 flex -translate-x-1/2 items-center gap-2">
          {slides.map((_, index) => (
            <button
              key={index}
              type="button"
              className={`h-2.5 rounded-full transition-all ${selectedIndex === index ? 'w-8 bg-white' : 'w-2.5 bg-white/50'}`}
              onClick={() => emblaApi?.scrollTo(index)}
              aria-label={`Go to slide ${index + 1}`}
            />
          ))}
        </div>
      ) : null}
    </section>
  )
}

