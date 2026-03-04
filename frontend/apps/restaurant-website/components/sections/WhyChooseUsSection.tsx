'use client'

import { Clock3, ShieldCheck, Truck, UtensilsCrossed } from 'lucide-react'
import type { ThemeLocale, ThemeSectionConfig } from '@/types/theme'
import { readLocalized, toArray, toRecord } from '@/components/sections/section-utils'

interface WhyChooseUsSectionProps {
  locale: ThemeLocale
  sectionConfig?: ThemeSectionConfig
}

const icons = [Truck, Clock3, ShieldCheck, UtensilsCrossed]

export function WhyChooseUsSection({ locale, sectionConfig }: WhyChooseUsSectionProps) {
  const settings = toRecord(sectionConfig?.settings)
  const items = toArray<Record<string, unknown>>(settings.items)
  const title = readLocalized(sectionConfig?.title ?? settings.title, locale, locale === 'ar' ? 'لماذا تختارنا' : 'Why Choose Us')

  const resolvedItems = (items.length > 0 ? items : [
    {
      title: { en: 'Fast Delivery', ar: 'توصيل سريع' },
      description: { en: 'Your order reaches you quickly and safely.', ar: 'طلبك يصلك بسرعة وأمان.' },
    },
    {
      title: { en: 'Fresh Ingredients', ar: 'مكونات طازجة' },
      description: { en: 'Prepared daily with quality ingredients.', ar: 'يتم التحضير يوميا بمكونات عالية الجودة.' },
    },
    {
      title: { en: 'Reliable Service', ar: 'خدمة موثوقة' },
      description: { en: 'Consistent quality every time you order.', ar: 'جودة ثابتة في كل طلب.' },
    },
    {
      title: { en: 'Great Taste', ar: 'مذاق رائع' },
      description: { en: 'Chef-crafted meals your family will love.', ar: 'وجبات محضرة بعناية يحبها الجميع.' },
    },
  ]).slice(0, 4)

  return (
    <section className="bg-[hsl(var(--theme-surface))] px-4 py-16 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-6xl">
        <h2 className={`mb-8 text-3xl font-bold text-[hsl(var(--theme-text))] ${locale === 'ar' ? 'text-right' : 'text-left'}`}>{title}</h2>
        <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
          {resolvedItems.map((item, index) => {
            const Icon = icons[index % icons.length]
            return (
              <article
                key={index}
                className="rounded-xl border border-[hsl(var(--theme-border))] bg-[hsl(var(--theme-background))] p-5 shadow-sm transition-transform hover:-translate-y-1"
              >
                <Icon className="mb-3 h-6 w-6 text-[hsl(var(--theme-primary))]" />
                <h3 className="mb-2 text-base font-semibold text-[hsl(var(--theme-text))]">
                  {readLocalized(item.title, locale, locale === 'ar' ? 'ميزة' : 'Feature')}
                </h3>
                <p className="text-sm text-[hsl(var(--theme-text))]/75">
                  {readLocalized(item.description, locale, '')}
                </p>
              </article>
            )
          })}
        </div>
      </div>
    </section>
  )
}

