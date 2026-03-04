'use client'

import type { ThemeConfig, ThemeLocale, ThemeSectionConfig } from '@/types/theme'
import { GradientHero } from '@/components/sections/heroes/GradientHero'
import { VideoHero } from '@/components/sections/heroes/VideoHero'
import { SliderHero } from '@/components/sections/heroes/SliderHero'
import { ParallaxHero } from '@/components/sections/heroes/ParallaxHero'
import { toRecord, toStringValue } from '@/components/sections/section-utils'
import { DEFAULT_TEMPLATE, TEMPLATE_REGISTRY } from '@/lib/utils/theme-templates'

interface HeroSectionProps {
  config?: ThemeSectionConfig
  locale: ThemeLocale
  themeConfig: ThemeConfig
  restaurantSlug: string
}

export function HeroSection({ config, locale, themeConfig, restaurantSlug }: HeroSectionProps) {
  const sectionSettings = toRecord(config?.settings)
  const sectionOverride = toRecord(config?.componentOverrides)
  const template = TEMPLATE_REGISTRY[themeConfig.templateId ?? 'template-01'] ?? DEFAULT_TEMPLATE

  const themeLevelOverride = toRecord((themeConfig as unknown as Record<string, unknown>).componentOverrides)
  const variant = toStringValue(
    sectionOverride.hero ??
      sectionOverride.heroVariant ??
      themeLevelOverride.hero ??
      themeLevelOverride.heroVariant ??
      (themeConfig as unknown as Record<string, unknown>).heroVariant ??
      (themeConfig as unknown as Record<string, unknown>).hero_variant ??
      template.heroVariant,
    'gradient'
  ).toLowerCase()

  if (variant === 'video') {
    return <VideoHero locale={locale} config={sectionSettings} restaurantSlug={restaurantSlug} />
  }
  if (variant === 'slider') {
    return <SliderHero locale={locale} config={sectionSettings} />
  }
  if (variant === 'parallax') {
    return <ParallaxHero locale={locale} config={sectionSettings} />
  }
  return <GradientHero locale={locale} config={sectionSettings} restaurantSlug={restaurantSlug} />
}

