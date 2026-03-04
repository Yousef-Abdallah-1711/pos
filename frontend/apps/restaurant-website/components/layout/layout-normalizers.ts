import type { FooterConfig, HeaderConfig } from '@pos-saas/ui-themes'
import type { ServerThemeData } from '@/lib/api/get-theme-server'

export type ThemeLocale = 'en' | 'ar'

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function toRecord(value: unknown): Record<string, unknown> {
  return isRecord(value) ? value : {}
}

function toNumber(value: unknown, fallback: number): number {
  if (typeof value === 'number' && Number.isFinite(value)) return value
  if (typeof value === 'string') {
    const parsed = Number(value)
    if (Number.isFinite(parsed)) return parsed
  }
  return fallback
}

function toBoolean(value: unknown, fallback: boolean): boolean {
  if (typeof value === 'boolean') return value
  return fallback
}

export function readLocalizedText(
  value: unknown,
  locale: ThemeLocale,
  fallback = ''
): string {
  if (typeof value === 'string') {
    const trimmed = value.trim()
    return trimmed.length > 0 ? trimmed : fallback
  }

  if (!isRecord(value)) return fallback

  const localeValue = value[locale]
  if (typeof localeValue === 'string' && localeValue.trim().length > 0) {
    return localeValue.trim()
  }

  const englishValue = value.en
  if (typeof englishValue === 'string' && englishValue.trim().length > 0) {
    return englishValue.trim()
  }

  const arabicValue = value.ar
  if (typeof arabicValue === 'string' && arabicValue.trim().length > 0) {
    return arabicValue.trim()
  }

  for (const candidate of Object.values(value)) {
    if (typeof candidate === 'string' && candidate.trim().length > 0) {
      return candidate.trim()
    }
  }

  return fallback
}

function normalizeInternalHref(href: string, locale: ThemeLocale): string {
  if (!href || href === '#') return '#'
  if (!href.startsWith('/')) return href
  if (href === '/') return `/${locale}`
  if (href === `/${locale}` || href.startsWith(`/${locale}/`)) return href
  return `/${locale}${href}`
}

function toArray(value: unknown): Record<string, unknown>[] {
  return Array.isArray(value) ? value.map(toRecord) : []
}

function normalizeHeaderNavItems(
  header: Record<string, unknown>,
  locale: ThemeLocale
): NonNullable<HeaderConfig['navigationItems']> {
  const defaultNav = [
    { id: 'menu', label: 'Menu', href: '/menu', order: 1 },
    { id: 'orders', label: 'Orders', href: '/orders', order: 2 },
    { id: 'cart', label: 'Cart', href: '/cart', order: 3 },
  ]

  const rawNav = toArray(header.navigationItems || header.navigation_items)
  const source = rawNav.length > 0 ? rawNav : defaultNav

  return source
    .map((item, index) => {
      const href = typeof item.href === 'string' ? item.href : '#'
      return {
        id: item.id ? String(item.id) : `nav-${index + 1}`,
        label: readLocalizedText(item.label, locale, 'Link'),
        href: normalizeInternalHref(href, locale),
        order: toNumber(item.order, index + 1),
      }
    })
    .sort((a, b) => (a.order || 0) - (b.order || 0))
}

export function normalizeHeaderConfig(
  themeData: ServerThemeData,
  locale: ThemeLocale
): HeaderConfig {
  const identity = toRecord(themeData?.identity)
  const header = toRecord(themeData?.header)

  const navPositionRaw = typeof header.navPosition === 'string'
    ? header.navPosition
    : typeof header.nav_position === 'string'
      ? header.nav_position
      : 'right'

  const navAlignmentRaw = typeof header.navAlignment === 'string'
    ? header.navAlignment
    : typeof header.nav_alignment === 'string'
      ? header.nav_alignment
      : 'horizontal'

  const navPosition: HeaderConfig['navPosition'] =
    navPositionRaw === 'left' || navPositionRaw === 'center' || navPositionRaw === 'right'
      ? navPositionRaw
      : 'right'

  const navAlignment: HeaderConfig['navAlignment'] =
    navAlignmentRaw === 'vertical' || navAlignmentRaw === 'horizontal'
      ? navAlignmentRaw
      : 'horizontal'

  const showShadowCandidate =
    header.showShadow !== undefined ? header.showShadow : header.shadow

  return {
    homeHref: `/${locale}`,
    logoUrl:
      readLocalizedText(header.logoUrl || header.logo_url, locale, '') ||
      readLocalizedText(identity.logoUrl || identity.logo_url, locale, ''),
    logoText:
      readLocalizedText(header.logoText || header.logo_text, locale, '') ||
      readLocalizedText(identity.siteTitle || identity.site_title, locale, 'Restaurant'),
    logoHeight: toNumber(header.logoHeight || header.logo_height, 40),
    showLogo: toBoolean(header.showLogo, true),
    navigationItems: normalizeHeaderNavItems(header, locale),
    navPosition,
    navAlignment,
    backgroundColor: readLocalizedText(header.backgroundColor || header.background_color, locale, '#ffffff'),
    textColor: readLocalizedText(header.textColor || header.text_color, locale, '#111827'),
    height: toNumber(header.height, 64),
    padding: toNumber(header.padding, 16),
    stickyHeader: toBoolean(header.stickyHeader || header.sticky, false),
    showShadow: toBoolean(showShadowCandidate, true),
    hideNavOnMobile: toBoolean(header.hideNavOnMobile || header.hide_nav_on_mobile, false),
  }
}

function normalizeFooterLinks(
  value: unknown,
  locale: ThemeLocale
): Array<{ label: string; href: string }> {
  return toArray(value).map((item, index) => ({
    label: readLocalizedText(item.label, locale, 'Link'),
    href: normalizeInternalHref(
      typeof item.href === 'string' ? item.href : '#',
      locale
    ) || `#link-${index + 1}`,
  }))
}

function normalizeFooterSections(
  footer: Record<string, unknown>,
  locale: ThemeLocale
): NonNullable<FooterConfig['footerSections']> {
  const sections = toArray(
    footer.footerSections || footer.footer_sections || footer.sections
  )

  return sections.map((section, index) => ({
    id: section.id ? String(section.id) : `section-${index + 1}`,
    title: readLocalizedText(section.title, locale, 'Section'),
    links: normalizeFooterLinks(section.links, locale),
  }))
}

function normalizeSocialLinks(
  footer: Record<string, unknown>
): NonNullable<FooterConfig['socialLinks']> {
  const links = toArray(footer.socialLinks || footer.social_links)
  return links.map((link, index) => ({
    id: link.id ? String(link.id) : `social-${index + 1}`,
    platform:
      typeof link.platform === 'string' && link.platform.trim().length > 0
        ? link.platform.toLowerCase()
        : 'social',
    url: typeof link.url === 'string' && link.url.trim().length > 0 ? link.url : '#',
  }))
}

function resolveLayout(value: unknown): FooterConfig['layout'] {
  return value === 'compact' || value === 'minimal' ? value : 'expanded'
}

export function normalizeFooterConfig(
  themeData: ServerThemeData,
  locale: ThemeLocale
): FooterConfig {
  const identity = toRecord(themeData?.identity)
  const footer = toRecord(themeData?.footer)

  const companyName =
    readLocalizedText(footer.companyName || footer.company_name, locale, '') ||
    readLocalizedText(identity.siteTitle || identity.site_title, locale, 'Restaurant')

  return {
    companyName,
    companyDescription: readLocalizedText(
      footer.companyDescription || footer.company_description,
      locale,
      ''
    ),
    address: readLocalizedText(footer.address, locale, ''),
    phone: readLocalizedText(footer.phone, locale, ''),
    email: readLocalizedText(footer.email, locale, ''),
    copyrightText:
      readLocalizedText(footer.copyrightText || footer.copyright_text, locale, '') ||
      `(c) ${new Date().getFullYear()} ${companyName}. All rights reserved.`,
    socialLinks: normalizeSocialLinks(footer),
    footerSections: normalizeFooterSections(footer, locale),
    legalLinks: normalizeFooterLinks(footer.legalLinks || footer.legal_links, locale),
    backgroundColor: readLocalizedText(footer.backgroundColor || footer.background_color, locale, '#111827'),
    textColor: readLocalizedText(footer.textColor || footer.text_color, locale, '#f3f4f6'),
    linkColor:
      readLocalizedText(footer.linkColor || footer.link_color, locale, '') ||
      readLocalizedText(footer.textColor || footer.text_color, locale, '#f3f4f6'),
    padding: toNumber(footer.padding, 48),
    showLinks: toBoolean(footer.showLinks, true),
    showLegal: toBoolean(footer.showLegal, true),
    showBackToTop: toBoolean(footer.showBackToTop || footer.show_back_to_top, false),
    columns: toNumber(footer.columns, 3),
    layout: resolveLayout(footer.layout),
  }
}
