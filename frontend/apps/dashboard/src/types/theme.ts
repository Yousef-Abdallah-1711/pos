// Theme and Component type definitions

export type Locale = 'en' | 'ar'

export type ComponentType = 'hero' | 'products' | 'why_us' | 'contact' | 'testimonials' | 'cta' | 'custom'

export interface GlobalColors {
  primary: string
  secondary: string
  accent: string
  background: string
  text: string
  border: string
  shadow: string
}

export interface TypographySettings {
  fontFamily: string
  baseFontSize: number
  borderRadius: number
  lineHeight: number
  fontWeights?: number[]
  fontSize?: number
  fontStyle?: string
}

export interface WebsiteIdentity {
  siteTitle: string
  logoUrl?: string
  faviconUrl?: string
  domain?: string
}

export interface HeaderConfig {
  logoUrl?: string
  logoText?: string
  logoHeight?: number
  showLogo: boolean
  navigationItems: Array<{ id?: string; label: string; href: string; order?: number }>
  navPosition?: 'left' | 'center' | 'right'
  navAlignment?: 'horizontal' | 'vertical'
  backgroundColor: string
  textColor: string
  height: number
  padding?: number
  stickyHeader?: boolean
  showShadow?: boolean
  hideNavOnMobile?: boolean
}

export interface FooterConfig {
  companyName: string
  companyDescription?: string
  address?: string
  phone?: string
  email?: string
  copyrightText: string
  socialLinks: Array<{ id?: string; platform: string; url: string; order?: number }>
  footerSections?: Array<{ id?: string; title: string; links: Array<{ label: string; href: string }>; order?: number }>
  legalLinks?: Array<{ label: string; href: string }>
  backgroundColor: string
  textColor: string
  linkColor?: string
  padding?: number
  showLinks: boolean
  showLegal?: boolean
  showBackToTop?: boolean
  columns?: number
  layout?: 'expanded' | 'compact'
}

export interface ComponentConfig {
  id: string
  type: ComponentType
  title: string
  displayOrder: number
  enabled: boolean
  config?: Record<string, any>
  settings?: Record<string, any> // Backward compatibility
  styles?: Record<string, any>
  createdAt?: string
  updatedAt?: string
}

export interface ThemeData {
  id: string
  restaurantId: string
  tenantId: string
  name: string
  slug: string
  description?: string
  version: number
  isActive: boolean
  isPublished: boolean
  colors: GlobalColors
  typography: TypographySettings
  identity: WebsiteIdentity
  websiteIdentity?: WebsiteIdentity // API uses this name
  header: HeaderConfig
  footer: FooterConfig
  components: ComponentConfig[]
  customCss?: string
  createdAt: string
  updatedAt: string
  publishedAt?: string
  themeId?: string // Some API responses include this
}

export interface ThemePreset {
  id: string
  name: string
  description: string
  thumbnail?: string
  colors: GlobalColors
  typography: TypographySettings
}

export interface ComponentLibraryItem {
  id: string
  type: ComponentType
  title: string
  description: string
  icon?: string
  defaultSettings: Record<string, any>
  settingsSchema: Record<string, any>
}

export interface ThemeHistoryEntry {
  id: string
  themeId: string
  version: number
  changes: string
  createdAt: string
  createdBy: string
}

export type ThemeComponent = ComponentConfig

export interface ThemeExport {
  version: '1.0'
  theme: ThemeData
  exportedAt: string
  exportedBy?: string
}

export interface CreateThemeRequest {
  name: string
  slug?: string
  description?: string
  colors: GlobalColors
  typography: TypographySettings
  identity?: WebsiteIdentity
  websiteIdentity?: WebsiteIdentity // API spec uses this name
  header?: HeaderConfig
  footer?: FooterConfig
  components?: ComponentConfig[]
  customCss?: string
  config?: Record<string, any>
}

export interface UpdateThemeRequest {
  name?: string
  slug?: string
  description?: string
  colors?: GlobalColors
  typography?: TypographySettings
  identity?: WebsiteIdentity
  websiteIdentity?: WebsiteIdentity // API spec uses this name
  header?: HeaderConfig
  footer?: FooterConfig
  header_config?: HeaderConfig
  footer_config?: FooterConfig
  components?: ComponentConfig[]
  customCss?: string
  config?: Record<string, any>
}

export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: string
  message?: string
  pagination?: {
    page: number
    limit: number
    total: number
    totalPages: number
  }
}

export interface ThemeListResponse {
  data: ThemeData[]
  message: string
  pagination: {
    page: number
    limit: number
    total: number
    totalPages: number
  }
}
