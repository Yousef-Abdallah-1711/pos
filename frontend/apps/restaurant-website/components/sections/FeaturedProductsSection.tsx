'use client'

import { Skeleton } from '@/components/ui/skeleton'
import type { Product } from '@/components/commerce/product-card'
import type { ThemeConfig, ThemeLocale, ThemeSectionConfig } from '@/types/theme'
import { ProductCarousel } from '@/components/sections/products/ProductCarousel'
import { ProductGrid } from '@/components/sections/products/ProductGrid'
import { ProductMasonry } from '@/components/sections/products/ProductMasonry'
import { DEFAULT_TEMPLATE, TEMPLATE_REGISTRY } from '@/lib/utils/theme-templates'
import { toRecord, toStringValue } from '@/components/sections/section-utils'

interface FeaturedProductsSectionProps {
  locale: ThemeLocale
  themeConfig: ThemeConfig
  sectionConfig?: ThemeSectionConfig
  products: Product[]
  isLoading?: boolean
}

export function FeaturedProductsSection({
  locale,
  themeConfig,
  sectionConfig,
  products,
  isLoading = false,
}: FeaturedProductsSectionProps) {
  const template = TEMPLATE_REGISTRY[themeConfig.templateId ?? 'template-01'] ?? DEFAULT_TEMPLATE
  const sectionOverride = toRecord(sectionConfig?.componentOverrides)
  const themeLevelOverride = toRecord((themeConfig as unknown as Record<string, unknown>).componentOverrides)
  const variant = toStringValue(
    sectionOverride.products ??
      sectionOverride.productsVariant ??
      themeLevelOverride.products ??
      themeLevelOverride.productsVariant ??
      (themeConfig as unknown as Record<string, unknown>).productsVariant ??
      (themeConfig as unknown as Record<string, unknown>).products_variant ??
      template.productsVariant,
    'carousel'
  ).toLowerCase()

  if (isLoading) {
    return (
      <section className="px-4 py-16 sm:px-6 lg:px-8">
        <div className="mx-auto grid max-w-6xl grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, index) => (
            <Skeleton key={index} className="h-72 w-full rounded-lg" />
          ))}
        </div>
      </section>
    )
  }

  const visibleProducts = products.slice(0, 8)
  if (visibleProducts.length === 0) {
    return null
  }

  if (variant === 'grid') {
    return <ProductGrid locale={locale} products={visibleProducts} />
  }
  if (variant === 'masonry') {
    return <ProductMasonry locale={locale} products={visibleProducts} />
  }
  return <ProductCarousel locale={locale} products={visibleProducts} />
}

