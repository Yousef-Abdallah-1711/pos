'use client'

import Link from 'next/link'
import type { ThemeLocale, ThemeSectionConfig } from '@/types/theme'
import { readLocalized } from '@/components/sections/section-utils'

interface CategoryLike {
  id?: number | string
  name?: string
  name_en?: string
  name_ar?: string
  item_count?: number
  count?: number
  emoji?: string
}

interface CategoriesSectionProps {
  locale: ThemeLocale
  categories: CategoryLike[]
  sectionConfig?: ThemeSectionConfig
}

export function CategoriesSection({ locale, categories, sectionConfig }: CategoriesSectionProps) {
  const title = readLocalized(sectionConfig?.title, locale, locale === 'ar' ? 'تصفح الفئات' : 'Browse Categories')

  if (categories.length === 0) {
    return null
  }

  return (
    <section className="px-4 py-14 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-6xl">
        <h2 className={`mb-6 text-2xl font-bold text-[hsl(var(--theme-text))] ${locale === 'ar' ? 'text-right' : 'text-left'}`}>{title}</h2>
        <div className="flex snap-x items-stretch gap-3 overflow-x-auto pb-2">
          {categories.map((category, index) => {
            const id = category.id ?? index
            const label = locale === 'ar'
              ? category.name_ar || category.name || category.name_en || 'قسم'
              : category.name_en || category.name || category.name_ar || 'Category'
            const count = Number(category.item_count ?? category.count ?? 0)
            const emoji = category.emoji || '🍽️'
            return (
              <Link
                key={id}
                href={`/${locale}/menu?category=${encodeURIComponent(String(id))}`}
                className="snap-start rounded-full border border-[hsl(var(--theme-border))] bg-[hsl(var(--theme-background))] px-4 py-3 text-sm transition hover:border-[hsl(var(--theme-primary))]"
              >
                <span className="inline-flex items-center gap-2 whitespace-nowrap text-[hsl(var(--theme-text))]">
                  <span>{emoji}</span>
                  <span>{label}</span>
                  <span className="rounded-full bg-[hsl(var(--theme-surface))] px-2 py-0.5 text-xs text-[hsl(var(--theme-text))]/75">
                    {count}
                  </span>
                </span>
              </Link>
            )
          })}
        </div>
      </div>
    </section>
  )
}

