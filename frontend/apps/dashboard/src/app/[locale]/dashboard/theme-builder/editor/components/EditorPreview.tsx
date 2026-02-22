'use client'

import { useState } from 'react'
import { useThemeStore } from '@/hooks/useThemeStore'
import { Footer, Header, SectionRenderer, type FooterConfig, type HeaderConfig } from '@pos-saas/ui-themes'
import { Smartphone, Monitor, Tablet, RefreshCw } from 'lucide-react'
import { mockRestaurantData } from '@/lib/mockData'
import type { ThemeData } from '@/types/theme'

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function extractLabel(label: unknown, fallback = ''): string {
  if (typeof label === 'string') {
    const trimmed = label.trim()
    return trimmed.length > 0 ? trimmed : fallback
  }

  if (isRecord(label)) {
    const english = label.en
    if (typeof english === 'string' && english.trim().length > 0) return english.trim()

    const arabic = label.ar
    if (typeof arabic === 'string' && arabic.trim().length > 0) return arabic.trim()

    for (const candidate of Object.values(label)) {
      if (typeof candidate === 'string' && candidate.trim().length > 0) {
        return candidate.trim()
      }
    }
  }

  return fallback
}

function normalizePreviewHref(href: string, locale: 'en' | 'ar'): string {
  if (!href || href === '#') return '#'
  if (!href.startsWith('/')) return href
  if (href === '/') return `/${locale}`
  if (href === `/${locale}` || href.startsWith(`/${locale}/`)) return href
  return `/${locale}${href}`
}

function normalizeHeaderConfig(theme: ThemeData, locale: 'en' | 'ar'): HeaderConfig {
  const defaultNavigation =
    locale === 'ar'
      ? [
          { id: 'home', label: 'الرئيسية', href: '/', order: 1 },
          { id: 'menu', label: 'القائمة', href: '/menu', order: 2 },
          { id: 'orders', label: 'الطلبات', href: '/orders', order: 3 },
          { id: 'cart', label: 'السلة', href: '/cart', order: 4 },
        ]
      : [
          { id: 'home', label: 'Home', href: '/', order: 1 },
          { id: 'menu', label: 'Menu', href: '/menu', order: 2 },
          { id: 'orders', label: 'Orders', href: '/orders', order: 3 },
          { id: 'cart', label: 'Cart', href: '/cart', order: 4 },
        ]

  const sourceNavigation =
    Array.isArray(theme.header?.navigationItems) && theme.header.navigationItems.length > 0
      ? theme.header.navigationItems
      : defaultNavigation

  const navigationItems = sourceNavigation
    .map((item, index) => ({
      id: item.id ? String(item.id) : `nav-${index + 1}`,
      label: extractLabel(item.label, locale === 'ar' ? 'رابط' : 'Link'),
      href: normalizePreviewHref(typeof item.href === 'string' ? item.href : '#', locale),
      order: typeof item.order === 'number' ? item.order : index + 1,
    }))
    .sort((a, b) => (a.order || 0) - (b.order || 0))

  return {
    homeHref: `/${locale}`,
    logoUrl: theme.header?.logoUrl || theme.identity?.logoUrl || '',
    logoText: extractLabel(theme.header?.logoText, '') || theme.identity?.siteTitle || 'Restaurant',
    logoHeight: theme.header?.logoHeight || 40,
    showLogo: theme.header?.showLogo !== false,
    navigationItems,
    navPosition: theme.header?.navPosition || 'right',
    navAlignment: theme.header?.navAlignment || 'horizontal',
    backgroundColor: theme.header?.backgroundColor || '#3b82f6',
    textColor: theme.header?.textColor || '#ffffff',
    height: theme.header?.height || 64,
    padding: theme.header?.padding || 16,
    stickyHeader: theme.header?.stickyHeader || false,
    showShadow: theme.header?.showShadow !== false,
    hideNavOnMobile: theme.header?.hideNavOnMobile || false,
  }
}

function normalizeFooterConfig(theme: ThemeData, locale: 'en' | 'ar'): FooterConfig {
  const companyName =
    extractLabel(theme.footer?.companyName, '') ||
    extractLabel(theme.identity?.siteTitle, '') ||
    'Restaurant'

  const footerSections = (theme.footer?.footerSections || []).map((section, sectionIndex) => ({
    id: section.id ? String(section.id) : `section-${sectionIndex + 1}`,
    title: extractLabel(section.title, locale === 'ar' ? 'قسم' : 'Section'),
    links: (section.links || []).map((link, linkIndex) => ({
      label: extractLabel(link.label, locale === 'ar' ? 'رابط' : 'Link'),
      href: normalizePreviewHref(link.href, locale) || `#link-${sectionIndex + 1}-${linkIndex + 1}`,
    })),
  }))

  const legalLinks = (theme.footer?.legalLinks || []).map((link) => ({
    label: extractLabel(link.label, locale === 'ar' ? 'قانوني' : 'Legal'),
    href: normalizePreviewHref(link.href, locale),
  }))

  const socialLinks = (theme.footer?.socialLinks || []).map((link, index) => ({
    id: link.id ? String(link.id) : `social-${index + 1}`,
    platform: extractLabel(link.platform, 'social').toLowerCase(),
    url: typeof link.url === 'string' ? link.url : '#',
  }))

  return {
    companyName,
    companyDescription: extractLabel(theme.footer?.companyDescription, ''),
    address: extractLabel(theme.footer?.address, ''),
    phone: extractLabel(theme.footer?.phone, ''),
    email: extractLabel(theme.footer?.email, ''),
    copyrightText:
      extractLabel(theme.footer?.copyrightText, '') ||
      `© ${new Date().getFullYear()} ${companyName}. All rights reserved.`,
    socialLinks,
    footerSections,
    legalLinks,
    backgroundColor: theme.footer?.backgroundColor || '#1f2937',
    textColor: theme.footer?.textColor || '#ffffff',
    linkColor: theme.footer?.linkColor || '#3b82f6',
    padding: theme.footer?.padding || 48,
    showLinks: theme.footer?.showLinks !== false,
    showLegal: theme.footer?.showLegal !== false,
    showBackToTop: theme.footer?.showBackToTop !== false,
    columns: theme.footer?.columns || 3,
    layout: theme.footer?.layout || 'expanded',
  }
}

function enrichComponentWithMockData(component: any) {
  const enriched = { ...component }

  switch (component.type?.toLowerCase()) {
    case 'hero':
      enriched.mockData = mockRestaurantData.hero
      break
    case 'featured-items':
    case 'products':
      enriched.mockData = {
        items: mockRestaurantData.menu?.items || [],
      }
      break
    case 'testimonials':
      enriched.mockData = {
        testimonials: mockRestaurantData.testimonials || [],
      }
      break
    case 'contact':
      enriched.mockData = {
        contactInfo: mockRestaurantData.contact || {},
      }
      break
    case 'cta':
      enriched.mockData = mockRestaurantData.cta
      break
    case 'why-choose-us':
      enriched.mockData = {
        reasons: mockRestaurantData.whyChooseUs || [],
      }
      break
    case 'info-cards':
      enriched.mockData = {
        cards: mockRestaurantData.infoCards || [],
      }
      break
  }

  return enriched
}

export function EditorPreview() {
  const theme = useThemeStore((state) => state.currentTheme)
  const [viewMode, setViewMode] = useState<'desktop' | 'tablet' | 'mobile'>('desktop')
  const [isArabic, setIsArabic] = useState(false)

  if (!theme) {
    return (
      <div className="flex items-center justify-center h-full bg-gray-100 dark:bg-gray-900 text-gray-500">
        No theme selected
      </div>
    )
  }

  const locale: 'en' | 'ar' = isArabic ? 'ar' : 'en'

  const colors = {
    primary: theme.colors?.primary || '#3b82f6',
    secondary: theme.colors?.secondary || '#1e40af',
    accent: theme.colors?.accent || '#0ea5e9',
    background: theme.colors?.background || '#ffffff',
    text: theme.colors?.text || '#1f2937',
    border: theme.colors?.border || '#e5e7eb',
    shadow: theme.colors?.shadow || '#1f2937',
  }

  const typography = {
    fontFamily: theme.typography?.fontFamily || 'Inter',
    baseFontSize: theme.typography?.baseFontSize || 14,
    borderRadius: theme.typography?.borderRadius || 8,
    lineHeight: theme.typography?.lineHeight || 1.5,
  }

  const headerConfig = normalizeHeaderConfig(theme, locale)
  const footerConfig = normalizeFooterConfig(theme, locale)

  const getCSSVariables = () => {
    return `
      :root {
        --color-primary: ${colors.primary};
        --color-secondary: ${colors.secondary};
        --color-accent: ${colors.accent};
        --color-background: ${colors.background};
        --color-text: ${colors.text};
        --color-border: ${colors.border};
        --color-shadow: ${colors.shadow};
        --font-family: ${typography.fontFamily};
        --font-size-base: ${typography.baseFontSize}px;
        --border-radius: ${typography.borderRadius}px;
        --line-height: ${typography.lineHeight};
      }
      ${theme.customCss || ''}
    `
  }

  return (
    <div className="h-full flex flex-col bg-gray-100 dark:bg-gray-900">
      <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between px-4 py-4">
          <h3 className="text-sm font-semibold text-gray-900 dark:text-white">Live Preview</h3>

          <div className="flex items-center gap-1 bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
            <button
              onClick={() => setIsArabic(false)}
              className={`px-2 py-1 rounded-md text-xs font-medium transition ${
                !isArabic
                  ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-white'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
              }`}
            >
              EN
            </button>
            <button
              onClick={() => setIsArabic(true)}
              className={`px-2 py-1 rounded-md text-xs font-medium transition ${
                isArabic
                  ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-white'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
              }`}
            >
              العربية
            </button>
          </div>

          <div className="flex items-center gap-1 bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
            <button
              onClick={() => setViewMode('desktop')}
              title="Desktop"
              className={`px-3 py-1 rounded-md text-sm transition flex items-center gap-1 ${
                viewMode === 'desktop'
                  ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-white'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
              }`}
            >
              <Monitor className="h-4 w-4" />
              Desktop
            </button>
            <button
              onClick={() => setViewMode('tablet')}
              title="Tablet"
              className={`px-3 py-1 rounded-md text-sm transition flex items-center gap-1 ${
                viewMode === 'tablet'
                  ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-white'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
              }`}
            >
              <Tablet className="h-4 w-4" />
              Tablet
            </button>
            <button
              onClick={() => setViewMode('mobile')}
              title="Mobile"
              className={`px-3 py-1 rounded-md text-sm transition flex items-center gap-1 ${
                viewMode === 'mobile'
                  ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-white'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
              }`}
            >
              <Smartphone className="h-4 w-4" />
              Mobile
            </button>
          </div>

          <button className="p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition" title="Refresh preview">
            <RefreshCw className="h-4 w-4" />
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6 flex items-start justify-center">
        <div
          className={`bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden transition-all ${
            viewMode === 'mobile' ? 'w-96 max-w-full' : viewMode === 'tablet' ? 'w-[768px] max-w-full' : 'w-full'
          }`}
        >
          <div
            dir={isArabic ? 'rtl' : 'ltr'}
            className="editor-preview-content min-h-screen"
          >
            <style>{getCSSVariables()}</style>
            <style jsx>{`
              .editor-preview-content {
                background-color: ${colors.background};
                color: ${colors.text};
                font-family: ${typography.fontFamily};
                font-size: ${typography.baseFontSize}px;
                line-height: ${typography.lineHeight};
              }
            `}</style>

            <Header
              config={headerConfig}
              isArabic={isArabic}
              onNavigationClick={() => {
                // Keep preview sandboxed inside the editor.
              }}
            />

            {theme.components && Array.isArray(theme.components) && theme.components.length > 0 && (
              <div>
                {theme.components
                  .filter((component) => component.enabled)
                  .sort((a: any, b: any) => {
                    const orderA =
                      typeof a.displayOrder === 'number'
                        ? a.displayOrder
                        : typeof a.order === 'number'
                          ? a.order
                          : 0
                    const orderB =
                      typeof b.displayOrder === 'number'
                        ? b.displayOrder
                        : typeof b.order === 'number'
                          ? b.order
                          : 0
                    return orderA - orderB
                  })
                  .map((component) => {
                    const enrichedComponent = enrichComponentWithMockData(component)
                    return (
                      <div key={component.id} className="w-full">
                        <SectionRenderer component={enrichedComponent} isArabic={isArabic} />
                      </div>
                    )
                  })}
              </div>
            )}

            <Footer
              config={footerConfig}
              isArabic={isArabic}
              onLinkClick={() => {
                // Keep preview sandboxed inside the editor.
              }}
            />
          </div>
        </div>
      </div>
    </div>
  )
}
