/**
 * Global color palette configuration
 */
export interface GlobalColors {
  primary: string
  secondary: string
  accent: string
  background: string
  text: string
  border: string
  shadow: string
}

/**
 * Typography settings for theming
 */
export interface TypographySettings {
  fontFamily: string
  baseFontSize: number
  borderRadius: number
  lineHeight: number
  fontWeights?: number[]
  fontSize?: number
  fontStyle?: string
}

/**
 * Website identity (branding)
 */
export interface WebsiteIdentity {
  siteTitle: string
  logoUrl?: string
  faviconUrl?: string
}

/**
 * Component configuration and settings
 */
export interface ThemeComponent {
  id: number
  themeId?: number
  type: string
  title: string
  enabled: boolean
  displayOrder: number
  settings?: Record<string, any>
  config?: Record<string, any>
}

/**
 * Complete theme configuration
 */
export interface Theme {
  id: number
  restaurantId: number
  name: string
  slug: string
  colors: GlobalColors
  typography: TypographySettings
  identity: WebsiteIdentity
  isActive: boolean
  isPublished: boolean
  components: ThemeComponent[]
  createdAt: string
  updatedAt: string
}

/**
 * Theme provider props
 */
export interface ThemeProviderProps {
  theme: Theme | { colors: GlobalColors; typography: TypographySettings }
  children: React.ReactNode
}

/**
 * Section renderer props
 */
export interface SectionRendererProps {
  component: ThemeComponent
  isArabic?: boolean
}

/**
 * Common section component props
 */
export interface BaseSectionProps {
  title?: string
  isArabic?: boolean
  settings?: Record<string, any>
  config?: Record<string, any>
}

/**
 * Hero section props
 */
export interface HeroSectionProps extends BaseSectionProps {
  title_en?: string
  title_ar?: string
  subtitle_en?: string
  subtitle_ar?: string
  description_en?: string
  description_ar?: string
  background_image_url?: string
  config?: {
    cta_button_text?: string
    cta_button_url?: string
    height?: string
    overlay_opacity?: number
    text_alignment?: string
  }
}

/**
 * Products section props
 */
export interface ProductsSectionProps extends BaseSectionProps {
  title_en?: string
  title_ar?: string
  description_en?: string
  description_ar?: string
  category_id?: number
  max_items?: number
  config?: {
    layout?: 'grid' | 'list'
    columns?: number
    show_prices?: boolean
    show_images?: boolean
    items?: Array<{
      id?: string
      name?: string
      description?: string
      price?: number
      image?: string
      category?: string
    }>
  }
}

/**
 * Why choose us section props
 */
export interface WhyChooseUsSectionProps extends BaseSectionProps {
  title_en?: string
  title_ar?: string
  description_en?: string
  description_ar?: string
  config?: {
    layout?: 'grid' | 'flex'
    columns?: number
    items?: Array<{
      title_en?: string
      title_ar?: string
      description_en?: string
      description_ar?: string
      icon?: string
    }>
  }
}

/**
 * Contact section props
 */
export interface ContactSectionProps extends BaseSectionProps {
  title_en?: string
  title_ar?: string
  phone?: string
  email?: string
  address_en?: string
  address_ar?: string
  config?: {
    show_form?: boolean
    show_map?: boolean
  }
}

/**
 * Testimonials section props
 */
export interface TestimonialsSectionProps extends BaseSectionProps {
  title_en?: string
  title_ar?: string
  config?: {
    layout?: 'carousel' | 'grid'
    items?: Array<{
      author?: string
      content_en?: string
      content_ar?: string
      rating?: number
    }>
  }
}

/**
 * CTA section props
 */
export interface CTASectionProps extends BaseSectionProps {
  title_en?: string
  title_ar?: string
  description_en?: string
  description_ar?: string
  config?: {
    button_text_en?: string
    button_text_ar?: string
    button_url?: string
    background_color?: string
  }
}

/**
 * Custom section props
 */
export interface CustomSectionProps extends BaseSectionProps {
  content?: string
  html_content?: string
  config?: Record<string, any>
}

/**
 * Header component configuration
 */
export interface HeaderConfig {
  homeHref?: string
  logoUrl?: string
  logoText?: string
  logoHeight?: number
  showLogo?: boolean
  navigationItems?: Array<{
    id?: string | number
    label: string
    href: string
    order?: number
  }>
  navPosition?: 'left' | 'center' | 'right'
  navAlignment?: 'horizontal' | 'vertical'
  backgroundColor?: string
  textColor?: string
  height?: number
  padding?: number
  stickyHeader?: boolean
  showShadow?: boolean
  hideNavOnMobile?: boolean
}

/**
 * Footer component configuration
 */
export interface FooterConfig {
  companyName?: string
  companyDescription?: string
  address?: string
  phone?: string
  email?: string
  copyrightText?: string
  socialLinks?: Array<{
    id?: string | number
    platform: string
    url: string
  }>
  footerSections?: Array<{
    id?: string | number
    title: string
    links?: Array<{
      label: string
      href: string
    }>
  }>
  legalLinks?: Array<{
    label: string
    href: string
  }>
  backgroundColor?: string
  textColor?: string
  linkColor?: string
  padding?: number
  showLinks?: boolean
  showLegal?: boolean
  showBackToTop?: boolean
  columns?: number
  layout?: 'expanded' | 'compact' | 'minimal'
}

/**
 * Complete theme data configuration
 */
export interface ThemeData {
  id?: number
  restaurantId?: number
  tenantId?: number
  name: string
  slug: string
  version?: number
  colors: GlobalColors
  typography: TypographySettings
  identity: WebsiteIdentity
  header?: HeaderConfig
  footer?: FooterConfig
  isActive: boolean
  isPublished: boolean
  components?: ThemeComponent[]
  customCss?: string
  createdAt?: string
  updatedAt?: string
}

/**
 * Homepage data response from API
 */
export interface HomepageData {
  theme: Theme
  components: ThemeComponent[]
  restaurant?: {
    id?: number
    slug: string
    name?: string
  }
}
