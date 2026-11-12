'use client'

import { useEffect, useMemo, useState } from 'react'
import useEmblaCarousel from 'embla-carousel-react'
import Autoplay from 'embla-carousel-autoplay'
import { Star } from 'lucide-react'
import type { ThemeLocale, ThemeSectionConfig } from '@/types/theme'
import { formatRating } from '@/lib/utils/format'
import { readLocalized, toArray, toRecord } from '@/components/sections/section-utils'

interface TestimonialsSectionProps {
  locale: ThemeLocale
  sectionConfig?: ThemeSectionConfig
}

export function TestimonialsSection({ locale, sectionConfig }: TestimonialsSectionProps) {
  const settings = toRecord(sectionConfig?.settings)
  const title = readLocalized(sectionConfig?.title ?? settings.title, locale, locale === 'ar' ? 'آراء العملاء' : 'Testimonials')
  const items = useMemo(() => {
    const configured = toArray<Record<string, unknown>>(settings.items)
    if (configured.length > 0) {
      return configured
    }
    return [
      {
        name: { en: 'Mona A.', ar: 'منى أ.' },
        review: { en: 'Great quality and very fast delivery.', ar: 'جودة ممتازة وتوصيل سريع جدا.' },
        rating: 4.8,
        count: 120,
      },
      {
        name: { en: 'Karim H.', ar: 'كريم ح.' },
        review: { en: 'The food is always fresh and delicious.', ar: 'الطعام دائما طازج ولذيذ.' },
        rating: 4.9,
        count: 86,
      },
      {
        name: { en: 'Nour S.', ar: 'نور س.' },
        review: { en: 'Best online ordering experience.', ar: 'أفضل تجربة طلب أونلاين.' },
        rating: 4.7,
        count: 59,
      },
    ]
  }, [settings])
  const [selectedIndex, setSelectedIndex] = useState(0)
  const [emblaRef, emblaApi] = useEmblaCarousel(
    { loop: true, align: 'start', direction: locale === 'ar' ? 'rtl' : 'ltr' },
    [Autoplay({ delay: 5000, stopOnMouseEnter: true })]
  )

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
    <section className="bg-[hsl(var(--theme-surface))] px-4 py-16 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-6xl">
        <h2 className={`mb-8 text-3xl font-bold text-[hsl(var(--theme-text))] ${locale === 'ar' ? 'text-right' : 'text-left'}`}>{title}</h2>

        <div className="overflow-hidden" ref={emblaRef}>
          <div className="flex">
            {items.map((item, index) => (
              <article key={index} className="min-w-0 flex-[0_0_100%] px-2 md:flex-[0_0_50%] lg:flex-[0_0_33.333%]">
                <div className="h-full rounded-xl border border-[hsl(var(--theme-border))] bg-[hsl(var(--theme-background))] p-6">
                  <div className="mb-4 flex items-center gap-1 text-amber-500">
                    {Array.from({ length: 5 }).map((_, starIndex) => (
                      <Star key={starIndex} className="h-4 w-4 fill-current" />
                    ))}
                  </div>
                  <p className="mb-4 text-sm text-[hsl(var(--theme-text))]/80">
                    {readLocalized(item.review, locale, '')}
                  </p>
                  <p className="font-semibold text-[hsl(var(--theme-text))]">{readLocalized(item.name, locale, '')}</p>
                  <p className="mt-1 text-xs text-[hsl(var(--theme-text))]/60">
                    {formatRating(Number(item.rating ?? 4.8), Number(item.count ?? 0), locale)}
                  </p>
                </div>
              </article>
            ))}
          </div>
        </div>

        {items.length > 1 ? (
          <div className="mt-6 flex items-center justify-center gap-2">
            {items.map((_, index) => (
              <button
                key={index}
                type="button"
                className={`h-2.5 rounded-full transition-all ${selectedIndex === index ? 'w-8 bg-[hsl(var(--theme-primary))]' : 'w-2.5 bg-[hsl(var(--theme-border))]'}`}
                onClick={() => emblaApi?.scrollTo(index)}
                aria-label={`Go to testimonial ${index + 1}`}
              />
            ))}
          </div>
        ) : null}
      </div>
    </section>
  )
}

