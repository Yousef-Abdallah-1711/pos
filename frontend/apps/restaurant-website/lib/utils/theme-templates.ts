import type { HeaderVariant, FooterVariant, SectionType } from '@/types/theme'

export interface ThemeTemplate {
  id: string
  name: string
  sectionOrder: SectionType[]
  headerVariant: HeaderVariant
  footerVariant: FooterVariant
  heroVariant: string
  productsVariant: string
}

export const TEMPLATE_REGISTRY: Record<string, ThemeTemplate> = {
  'template-01': {
    id: 'template-01',
    name: 'Modern Bistro',
    sectionOrder: ['hero', 'featured_products', 'why_choose_us', 'testimonials', 'cta', 'contact'],
    headerVariant: 'glass',
    heroVariant: 'video',
    productsVariant: 'carousel',
    footerVariant: 'expanded',
  },
  'template-02': {
    id: 'template-02',
    name: 'Luxury Gold',
    sectionOrder: ['hero', 'gallery', 'testimonials', 'featured_products', 'cta', 'contact'],
    headerVariant: 'gradient',
    heroVariant: 'gradient',
    productsVariant: 'grid',
    footerVariant: 'cards',
  },
  'template-03': {
    id: 'template-03',
    name: 'Urban Minimal',
    sectionOrder: ['hero', 'categories', 'featured_products', 'cta', 'contact'],
    headerVariant: 'minimal',
    heroVariant: 'gradient',
    productsVariant: 'grid',
    footerVariant: 'minimal',
  },
  'template-04': {
    id: 'template-04',
    name: 'Warm Rustic',
    sectionOrder: ['hero', 'why_choose_us', 'featured_products', 'gallery', 'cta', 'contact'],
    headerVariant: 'floating',
    heroVariant: 'parallax',
    productsVariant: 'masonry',
    footerVariant: 'social-heavy',
  },
  'template-05': {
    id: 'template-05',
    name: 'Family Classic',
    sectionOrder: ['hero', 'categories', 'featured_products', 'testimonials', 'contact', 'cta'],
    headerVariant: 'solid',
    heroVariant: 'slider',
    productsVariant: 'grid',
    footerVariant: 'expanded',
  },
  'template-06': {
    id: 'template-06',
    name: 'Vibrant Street',
    sectionOrder: ['hero', 'featured_products', 'categories', 'gallery', 'testimonials', 'cta'],
    headerVariant: 'gradient',
    heroVariant: 'video',
    productsVariant: 'carousel',
    footerVariant: 'cards',
  },
  'template-07': {
    id: 'template-07',
    name: 'Elegant Night',
    sectionOrder: ['hero', 'gallery', 'featured_products', 'why_choose_us', 'testimonials', 'contact'],
    headerVariant: 'glass',
    heroVariant: 'gradient',
    productsVariant: 'masonry',
    footerVariant: 'minimal',
  },
  'template-08': {
    id: 'template-08',
    name: 'Fast Casual',
    sectionOrder: ['hero', 'categories', 'featured_products', 'cta', 'contact', 'testimonials'],
    headerVariant: 'minimal',
    heroVariant: 'slider',
    productsVariant: 'grid',
    footerVariant: 'compact',
  },
  'template-09': {
    id: 'template-09',
    name: 'Ocean Fresh',
    sectionOrder: ['hero', 'featured_products', 'gallery', 'why_choose_us', 'cta', 'contact'],
    headerVariant: 'floating',
    heroVariant: 'video',
    productsVariant: 'carousel',
    footerVariant: 'expanded',
  },
  'template-10': {
    id: 'template-10',
    name: 'Desert Gold',
    sectionOrder: ['hero', 'testimonials', 'featured_products', 'categories', 'cta', 'contact'],
    headerVariant: 'solid',
    heroVariant: 'parallax',
    productsVariant: 'grid',
    footerVariant: 'social-heavy',
  },
}

export const DEFAULT_TEMPLATE = TEMPLATE_REGISTRY['template-01']
