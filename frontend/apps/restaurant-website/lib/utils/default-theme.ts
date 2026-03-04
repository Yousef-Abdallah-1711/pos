import type { ThemeConfig } from '@/types/theme'
import type { ThemeData } from '@/lib/store/theme-store'

export const DEFAULT_THEME: ThemeConfig = {
  id: 'default-theme-v2',
  name: 'Default Theme',
  slug: 'default-theme',
  templateId: 'template-01',
  headerVariant: 'solid',
  footerVariant: 'expanded',
  colors: {
    primary: '#3b82f6',
    secondary: '#1e40af',
    accent: '#f59e0b',
    background: '#ffffff',
    text: '#111827',
    border: '#e5e7eb',
    shadow: '#111827',
    primary_dark: '#60a5fa',
    secondary_dark: '#2563eb',
    accent_dark: '#fbbf24',
    background_dark: '#0f172a',
    text_dark: '#e2e8f0',
    border_dark: '#334155',
    shadow_dark: '#020617',
    surface_color: '#f8fafc',
    surface_dark: '#1e293b',
    header_bg: '#3b82f6',
    header_text: '#ffffff',
    header_bg_dark: '#111827',
    header_text_dark: '#f8fafc',
    footer_bg: '#111827',
    footer_text: '#f8fafc',
    footer_bg_dark: '#020617',
    footer_text_dark: '#e2e8f0',
    dark_mode_enabled: 'true',
  },
  typography: {
    font_family: 'Inter, sans-serif',
    base_font_size: 16,
    line_height: 1.5,
    border_radius: 8,
  },
  identity: {
    site_title: 'Restaurant',
    logo_url: '',
    favicon_url: '',
  },
  header: {
    headerVariant: 'solid',
    showLogo: true,
    stickyHeader: true,
    showShadow: true,
    logoText: { en: 'Restaurant', ar: 'مطعم' },
    navigationItems: [
      { id: 'menu', label: { en: 'Menu', ar: 'القائمة' }, href: '/menu', order: 1 },
      { id: 'orders', label: { en: 'Orders', ar: 'الطلبات' }, href: '/orders', order: 2 },
      { id: 'cart', label: { en: 'Cart', ar: 'السلة' }, href: '/cart', order: 3 },
    ],
    backgroundColor: '#3b82f6',
    textColor: '#ffffff',
    height: 64,
  },
  footer: {
    footerVariant: 'expanded',
    companyName: { en: 'Restaurant', ar: 'مطعم' },
    companyDescription: {
      en: 'Fresh food delivered daily.',
      ar: 'طعام طازج يوميا.',
    },
    socialLinks: [],
    footerSections: [],
    legalLinks: [],
    backgroundColor: '#111827',
    textColor: '#f8fafc',
  },
  components: [
    { id: 'hero-1', type: 'hero', enabled: true, order: 1 },
    { id: 'featured-products-1', type: 'featured_products', enabled: true, order: 2 },
    { id: 'cta-1', type: 'cta', enabled: true, order: 3 },
  ],
}

export function getDefaultTheme(): ThemeData {
  return {
    name: DEFAULT_THEME.name,
    slug: DEFAULT_THEME.slug,
    description: 'Default theme fallback',
    colors: {
      primary: DEFAULT_THEME.colors.primary,
      secondary: DEFAULT_THEME.colors.secondary,
      accent: DEFAULT_THEME.colors.accent,
      background: DEFAULT_THEME.colors.background,
      text: DEFAULT_THEME.colors.text,
      border: DEFAULT_THEME.colors.border,
      shadow: DEFAULT_THEME.colors.shadow,
    },
    typography: {
      font_family: String(DEFAULT_THEME.typography.font_family ?? 'Inter, sans-serif'),
      base_font_size: Number(DEFAULT_THEME.typography.base_font_size ?? 16),
      line_height: Number(DEFAULT_THEME.typography.line_height ?? 1.5),
      border_radius: Number(DEFAULT_THEME.typography.border_radius ?? 8),
    },
    identity: {
      site_title: String((DEFAULT_THEME.identity as Record<string, unknown>)?.site_title ?? 'Restaurant'),
      logo_url: String((DEFAULT_THEME.identity as Record<string, unknown>)?.logo_url ?? ''),
      favicon_url: String((DEFAULT_THEME.identity as Record<string, unknown>)?.favicon_url ?? ''),
    },
    header: {
      style: DEFAULT_THEME.header.headerVariant ?? 'solid',
      sticky_nav: DEFAULT_THEME.header.stickyHeader ?? true,
      show_search: true,
      show_language: true,
      logo_url: DEFAULT_THEME.header.logoUrl ?? '',
      show_logo: DEFAULT_THEME.header.showLogo ?? true,
      navigation_items: (DEFAULT_THEME.header.navigationItems ?? []).map((item) => ({
        id: item.id ?? item.href,
        label: typeof item.label === 'string' ? item.label : item.label.en ?? item.label.ar ?? 'Link',
        href: item.href,
        order: item.order ?? 0,
      })),
      background_color: DEFAULT_THEME.header.backgroundColor ?? '#3b82f6',
      text_color: DEFAULT_THEME.header.textColor ?? '#ffffff',
      height: DEFAULT_THEME.header.height ?? 64,
      padding: 16,
      show_shadow: DEFAULT_THEME.header.showShadow ?? true,
      sticky_header: DEFAULT_THEME.header.stickyHeader ?? true,
      hide_nav_on_mobile: false,
      nav_position: 'right',
    },
    footer: {
      style: DEFAULT_THEME.footer.footerVariant ?? 'expanded',
      columns: 3,
      show_social: true,
      company_name: typeof DEFAULT_THEME.footer.companyName === 'string' ? DEFAULT_THEME.footer.companyName : DEFAULT_THEME.name,
      company_description: typeof DEFAULT_THEME.footer.companyDescription === 'string' ? DEFAULT_THEME.footer.companyDescription : '',
      address: typeof DEFAULT_THEME.footer.address === 'string' ? DEFAULT_THEME.footer.address : '',
      phone: DEFAULT_THEME.footer.phone ?? '',
      email: DEFAULT_THEME.footer.email ?? '',
      copyright_text: `© ${new Date().getFullYear()} ${DEFAULT_THEME.name}. All rights reserved.`,
      social_links: {},
      footer_sections: [],
      legal_links: [],
      background_color: DEFAULT_THEME.footer.backgroundColor ?? '#111827',
      text_color: DEFAULT_THEME.footer.textColor ?? '#f8fafc',
      show_links: true,
    },
    components: DEFAULT_THEME.components.map((component) => ({
      id: component.id,
      type: component.type,
      enabled: component.enabled,
      order: component.order,
    })),
  }
}

export function getDefaultDarkTheme(): ThemeData {
  const base = getDefaultTheme()
  return {
    ...base,
    name: 'Default Theme Dark',
    slug: 'default-theme-dark',
    colors: {
      ...base.colors,
      primary: DEFAULT_THEME.colors.primary_dark,
      secondary: DEFAULT_THEME.colors.secondary_dark,
      accent: DEFAULT_THEME.colors.accent_dark,
      background: DEFAULT_THEME.colors.background_dark,
      text: DEFAULT_THEME.colors.text_dark,
      border: DEFAULT_THEME.colors.border_dark,
      shadow: DEFAULT_THEME.colors.shadow_dark,
    },
    header: {
      ...base.header,
      background_color: DEFAULT_THEME.colors.header_bg_dark,
      text_color: DEFAULT_THEME.colors.header_text_dark,
    },
    footer: {
      ...base.footer,
      background_color: DEFAULT_THEME.colors.footer_bg_dark,
      text_color: DEFAULT_THEME.colors.footer_text_dark,
    },
  }
}

export function getDefaultThemeForSystem(): ThemeData {
  if (typeof window === 'undefined') {
    return getDefaultTheme()
  }

  return window.matchMedia('(prefers-color-scheme: dark)').matches
    ? getDefaultDarkTheme()
    : getDefaultTheme()
}
