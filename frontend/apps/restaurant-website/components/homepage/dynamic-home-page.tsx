'use client'

import type { Product } from '@/components/commerce/product-card'
import type { ThemeConfig, ThemeLocale, ThemeSectionConfig, SectionType } from '@/types/theme'
import { DEFAULT_TEMPLATE, TEMPLATE_REGISTRY } from '@/lib/utils/theme-templates'
import { HeroSection } from '@/components/sections/HeroSection'
import { FeaturedProductsSection } from '@/components/sections/FeaturedProductsSection'
import { CategoriesSection } from '@/components/sections/CategoriesSection'
import { WhyChooseUsSection } from '@/components/sections/WhyChooseUsSection'
import { GallerySection } from '@/components/sections/GallerySection'
import { TestimonialsSection } from '@/components/sections/TestimonialsSection'
import { CTASection } from '@/components/sections/CTASection'
import { ContactSection } from '@/components/sections/ContactSection'
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

interface DynamicHomePageProps {
  locale: ThemeLocale
  restaurantSlug: string
  themeConfig: ThemeConfig
  products: Product[]
  categories: CategoryLike[]
}

function mapSectionType(value: string): SectionType {
  const normalized = value.toLowerCase()
  if (normalized === 'featured-items' || normalized === 'featured_items') return 'featured_products'
  if (normalized === 'why-choose-us' || normalized === 'why_us' || normalized === 'why-us') return 'why_choose_us'
  if (normalized === 'call-to-action') return 'cta'
  if (normalized === 'food_gallery') return 'gallery'
  return normalized as SectionType
}

function getEnabledSections(themeConfig: ThemeConfig): ThemeSectionConfig[] {
  const template = TEMPLATE_REGISTRY[themeConfig.templateId ?? 'template-01'] ?? DEFAULT_TEMPLATE
  const sectionOrder = template.sectionOrder
  const items = (themeConfig.components || [])
    .map((section) => ({
      ...section,
      type: mapSectionType(section.type),
    }))
    .filter((section) => section.enabled)

  return items.sort((a, b) => {
    const indexA = sectionOrder.indexOf(a.type)
    const indexB = sectionOrder.indexOf(b.type)
    if (indexA !== -1 && indexB !== -1) return indexA - indexB
    if (indexA !== -1) return -1
    if (indexB !== -1) return 1
    return a.order - b.order
  })
}

export function DynamicHomePage({
  locale,
  restaurantSlug,
  themeConfig,
  products,
  categories,
}: DynamicHomePageProps) {
  const sections = getEnabledSections(themeConfig)
  const templateId = themeConfig.templateId ?? 'template-01'
  const identity = themeConfig.identity ?? {}
  const defaultHeadline = readLocalized(
    identity.siteTitle ?? identity.site_title ?? identity.name,
    locale,
    locale === 'ar' ? 'مرحباً بكم' : 'Welcome'
  )

  if (sections.length === 0) {
    return (
      <main data-template-id={templateId}>
        <div data-section-type="default-hero" data-section-order={1}>
          <HeroSection
            config={{
              id: 'default-home-hero',
              type: 'hero',
              enabled: true,
              order: 0,
              settings: { headline: defaultHeadline },
            }}
            locale={locale}
            themeConfig={themeConfig}
            restaurantSlug={restaurantSlug}
          />
        </div>
      </main>
    )
  }

  return (
    <main data-template-id={templateId}>
      {sections.map((section, index) => {
        if (section.type === 'hero') {
          return (
            <div key={section.id} data-section-type="hero" data-section-order={index + 1}>
              <HeroSection
                config={section}
                locale={locale}
                themeConfig={themeConfig}
                restaurantSlug={restaurantSlug}
              />
            </div>
          )
        }

        if (section.type === 'featured_products') {
          return (
            <div key={section.id} data-section-type="featured_products" data-section-order={index + 1}>
              <FeaturedProductsSection
                locale={locale}
                themeConfig={themeConfig}
                sectionConfig={section}
                products={products}
              />
            </div>
          )
        }

        if (section.type === 'categories') {
          return (
            <div key={section.id} data-section-type="categories" data-section-order={index + 1}>
              <CategoriesSection locale={locale} categories={categories} sectionConfig={section} />
            </div>
          )
        }

        if (section.type === 'why_choose_us') {
          return (
            <div key={section.id} data-section-type="why_choose_us" data-section-order={index + 1}>
              <WhyChooseUsSection locale={locale} sectionConfig={section} />
            </div>
          )
        }

        if (section.type === 'gallery') {
          return (
            <div key={section.id} data-section-type="gallery" data-section-order={index + 1}>
              <GallerySection locale={locale} sectionConfig={section} />
            </div>
          )
        }

        if (section.type === 'testimonials') {
          return (
            <div key={section.id} data-section-type="testimonials" data-section-order={index + 1}>
              <TestimonialsSection locale={locale} sectionConfig={section} />
            </div>
          )
        }

        if (section.type === 'cta') {
          return (
            <div key={section.id} data-section-type="cta" data-section-order={index + 1}>
              <CTASection locale={locale} sectionConfig={section} />
            </div>
          )
        }

        if (section.type === 'contact') {
          return (
            <div key={section.id} data-section-type="contact" data-section-order={index + 1}>
              <ContactSection locale={locale} sectionConfig={section} />
            </div>
          )
        }

        return null
      })}
    </main>
  )
}
