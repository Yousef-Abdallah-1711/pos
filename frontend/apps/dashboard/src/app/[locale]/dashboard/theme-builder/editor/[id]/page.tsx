'use client'

import { useState, useEffect, use } from 'react'
import { usePathname, useRouter } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { Menu, X, Save, ArrowLeft, Loader2, AlertCircle } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { useThemeStore } from '@/hooks/useThemeStore'
import { ThemeData } from '@/types/theme'
import * as themeApi from '@/lib/api/themeApi'
import { toast } from 'sonner'
import { debugLog, debugWarn, debugError } from '@/lib/debug'

import { EditorSidebar } from '../components/EditorSidebar'
import { EditorPreview } from '../components/EditorPreview'
import { ErrorBoundary } from '@/components/error-boundary'

interface EditorPageProps {
  params: Promise<{
    id: string
  }>
}

// Helper function to create default components for preview
function createDefaultComponents() {
  const components = [
    {
      id: 1,
      type: 'hero',
      title: 'Hero Section',
      enabled: true,
      displayOrder: 0,
      config: {
        title_en: 'Welcome to Our Restaurant',
        title_ar: 'مرحبا بك في مطعمنا',
        subtitle_en: 'Experience authentic flavors and exceptional service',
        subtitle_ar: 'اختبر النكهات الأصلية والخدمة الاستثنائية',
        description_en: 'Join us for an unforgettable dining experience',
        description_ar: 'انضم إلينا لتجربة طعام لا تُنسى',
        cta_button_text: 'Learn More',
        cta_button_url: '#',
        height: 'medium',
        overlay_opacity: 0.3,
      },
    },
    {
      id: 2,
      type: 'products',
      title: 'Our Products',
      enabled: true,
      displayOrder: 1,
      config: {
        title_en: 'Our Products',
        title_ar: 'منتجاتنا',
        layout: 'grid',
        columns: 3,
        show_prices: true,
        show_images: true,
      },
    },
    {
      id: 3,
      type: 'why_us',
      title: 'Why Choose Us',
      enabled: true,
      displayOrder: 2,
      config: {
        title_en: 'Why Choose Us',
        title_ar: 'لماذا اختيارنا',
        layout: 'grid',
        columns: 4,
      },
    },
    {
      id: 4,
      type: 'testimonials',
      title: 'Customer Testimonials',
      enabled: true,
      displayOrder: 3,
      config: {
        title_en: 'What Our Customers Say',
        title_ar: 'ما يقوله عملاؤنا',
        layout: 'grid',
      },
    },
    {
      id: 5,
      type: 'contact',
      title: 'Contact Us',
      enabled: true,
      displayOrder: 4,
      config: {
        title_en: 'Contact Us',
        title_ar: 'اتصل بنا',
        show_form: true,
        show_map: true,
      },
    },
    {
      id: 6,
      type: 'cta',
      title: 'Call to Action',
      enabled: true,
      displayOrder: 5,
      config: {
        title_en: 'Ready to Get Started?',
        title_ar: 'هل أنت مستعد للبدء؟',
        description_en: 'Order now and get 10% off your first meal',
        description_ar: 'اطلب الآن واحصل على خصم 10% على وجبتك الأولى',
        button_text_en: 'Order Now',
        button_text_ar: 'اطلب الآن',
        button_url: '#',
        background_color: '#3b82f6',
      },
    },
  ]
  return components
}

// Helper function to extract label text (handles both string and translation objects)
// Can be translation objects like {en: "...", ar: "..."} or plain strings
function extractLabel(label: any): string {
  if (typeof label === 'string') return label
  if (label && typeof label === 'object') {
    // Handle translation objects like {en: "...", ar: "..."}
    return label.en || label.ar || Object.values(label)[0] || 'Link'
  }
  return ''
}

// Helper function to ensure theme has all required fields
function normalizeTheme(theme: any): ThemeData {
  // Normalize navigation items from backend
  const normalizeNavItems = (items: any) => {
    if (!Array.isArray(items)) return []
    return items.map((item: any) => ({
      ...item,
      label: extractLabel(item.label),
      id: item?.id ? String(item.id) : undefined,
      href: typeof item?.href === 'string' ? item.href : '/',
    }))
  }

  const themeHeader = theme.header || theme.header_config || {}
  const themeFooter = theme.footer || theme.footer_config || {}
  const themeIdentity = theme.identity || {}
  const themeTypography = theme.typography || {}
  const themeColors = theme.colors || {}

  const normalizedTheme = {
    ...theme,
    // Normalize top-level string fields that might be translation objects
    name: extractLabel(theme.name),
    slug: extractLabel(theme.slug),
    description: extractLabel(theme.description),
    header: {
      logoUrl: '',
      logoText: 'My Restaurant',
      logoHeight: 40,
      showLogo: true,
      navigationItems: [
        { id: '1', label: 'Home', href: '/', order: 1 },
        { id: '2', label: 'Menu', href: '/menu', order: 2 },
        { id: '3', label: 'About', href: '/about', order: 3 },
        { id: '4', label: 'Contact', href: '/contact', order: 4 },
      ],
      navPosition: 'right',
      navAlignment: 'horizontal',
      backgroundColor: '#3b82f6',
      textColor: '#ffffff',
      height: 64,
      padding: 16,
      stickyHeader: false,
      showShadow: true,
      hideNavOnMobile: false,
      ...themeHeader,
      logoUrl: themeHeader.logoUrl || themeHeader.logo_url || '',
      logoText: extractLabel(themeHeader.logoText || themeHeader.logo_text),
      navigationItems: normalizeNavItems(themeHeader.navigationItems || themeHeader.navigation_items),
      backgroundColor: themeHeader.backgroundColor || themeHeader.background_color || '#3b82f6',
      textColor: themeHeader.textColor || themeHeader.text_color || '#ffffff',
      height: Number(themeHeader.height || themeHeader.navHeight || themeHeader.nav_height || 64),
      stickyHeader: Boolean(themeHeader.stickyHeader || themeHeader.sticky),
      showShadow: themeHeader.showShadow !== false && themeHeader.shadow !== false,
    },
    footer: {
      companyName: 'My Restaurant',
      companyDescription: 'Delicious food for everyone',
      address: '123 Main Street',
      phone: '(555) 123-4567',
      email: 'contact@restaurant.com',
      copyrightText: '© 2025 My Restaurant. All rights reserved.',
      socialLinks: [],
      footerSections: [],
      legalLinks: [
        { label: 'Privacy', href: '/privacy' },
        { label: 'Terms', href: '/terms' },
      ],
      backgroundColor: '#1f2937',
      textColor: '#ffffff',
      linkColor: '#3b82f6',
      padding: 48,
      showLinks: true,
      showLegal: true,
      showBackToTop: true,
      columns: 3,
      layout: 'expanded',
      ...themeFooter,
      // Normalize all footer string fields that might be translation objects
      companyName: extractLabel(themeFooter.companyName || themeFooter.company_name || themeIdentity.siteTitle || themeIdentity.site_title),
      companyDescription: extractLabel(themeFooter.companyDescription || themeFooter.company_description),
      address: extractLabel(themeFooter.address),
      phone: extractLabel(themeFooter.phone),
      email: extractLabel(themeFooter.email),
      copyrightText: extractLabel(themeFooter.copyrightText || themeFooter.copyright_text),
      legalLinks: normalizeNavItems(themeFooter.legalLinks || themeFooter.legal_links),
      footerSections: (themeFooter.footerSections || themeFooter.footer_sections || themeFooter.sections || []).map((section: any) => ({
        ...section,
        title: extractLabel(section.title),
        links: normalizeNavItems(section.links),
      })),
      socialLinks: (themeFooter.socialLinks || themeFooter.social_links || []).map((link: any) => ({
        ...link,
        id: link?.id ? String(link.id) : undefined,
        platform: extractLabel(link.platform),
      })),
      backgroundColor: themeFooter.backgroundColor || themeFooter.background_color || '#1f2937',
      textColor: themeFooter.textColor || themeFooter.text_color || '#ffffff',
    },
    identity: {
      siteTitle: extractLabel(themeIdentity.siteTitle || themeIdentity.site_title),
      logoUrl: themeIdentity.logoUrl || themeIdentity.logo_url || '',
      faviconUrl: themeIdentity.faviconUrl || themeIdentity.favicon_url || '',
      domain: themeIdentity.domain || '',
    },
    customCss: theme.customCss || theme.custom_css || '',
    typography: {
      fontFamily: 'Inter',
      baseFontSize: 14,
      borderRadius: 8,
      lineHeight: 1.5,
      fontWeights: [400, 700],
      fontSize: 14,
      fontStyle: 'normal',
      ...themeTypography,
      fontFamily: themeTypography.fontFamily || themeTypography.font_family || 'Inter',
      baseFontSize: Number(themeTypography.baseFontSize || themeTypography.base_size || themeTypography.base_font_size || 14),
      borderRadius: Number(themeTypography.borderRadius || themeTypography.border_radius || 8),
      lineHeight: Number(themeTypography.lineHeight || themeTypography.line_height || 1.5),
    },
    colors: {
      primary: '#3b82f6',
      secondary: '#1e40af',
      accent: '#0ea5e9',
      background: '#ffffff',
      text: '#1f2937',
      border: '#e5e7eb',
      shadow: '#1f2937',
      ...themeColors,
    },
    // Use actual components from backend - DO NOT auto-fill with defaults
    // If backend says 0 components, show 0 components
    components: Array.isArray(theme.components)
      ? theme.components.map((component: any, index: number) => ({
          ...component,
          id: component?.id ? String(component.id) : `component-${index + 1}`,
          title: extractLabel(component?.title || component?.type || 'Component'),
          displayOrder: Number(
            component?.displayOrder !== undefined ? component.displayOrder : component?.order !== undefined ? component.order : index
          ),
          config: component?.config || component?.settings || {},
          settings: component?.settings || component?.config || {},
          enabled: component?.enabled !== false,
        }))
      : [],
  }

  return normalizedTheme
}

export default function ThemeEditorPage({ params: paramsPromise }: EditorPageProps) {
  // ⚠️ FIX: In Next.js 15.5.9, params is now a Promise, use React.use() to unwrap it
  const params = use(paramsPromise)

  const pathname = usePathname()
  const router = useRouter()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const user = useAuthStore((state) => state.user)

  const currentTheme = useThemeStore((state) => state.currentTheme)
  const setCurrentTheme = useThemeStore((state) => state.setCurrentTheme)
  const isDirty = useThemeStore((state) => state.isDirty)
  const isSaving = useThemeStore((state) => state.isSaving)
  const setSaving = useThemeStore((state) => state.setSaving)
  const setDirty = useThemeStore((state) => state.setDirty)

  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // ✅ FIX: Load theme once params and user are available
  useEffect(() => {
    const loadTheme = async () => {
      if (!user?.restaurant_id) {
        debugWarn('⚠️ [ThemeEditor] Missing restaurant_id')
        setIsLoading(false)
        return
      }

      try {
        setIsLoading(true)
        setError(null)
        debugLog('📡 [ThemeEditor] Loading theme:', params.id)
        const theme = await themeApi.getTheme(params.id)
        debugLog('✅ [ThemeEditor] Theme loaded:', theme)
        debugLog('📊 [ThemeEditor] Components from backend:', theme.components?.length || 0, 'components')
        if (theme.components) {
          theme.components.forEach((comp, idx) => {
            debugLog(`  Component ${idx}: id=${comp.id}, type=${comp.type}, enabled=${comp.enabled}`)
          })
        }
        const normalizedTheme = normalizeTheme(theme)
        debugLog('📊 [ThemeEditor] After normalize:', normalizedTheme.components?.length || 0, 'components')
        setCurrentTheme(normalizedTheme)
      } catch (err) {
        debugError('❌ [ThemeEditor] Caught error during load')
        debugError('❌ [ThemeEditor] Error instanceof Error:', err instanceof Error)
        debugError('❌ [ThemeEditor] Error type:', typeof err)
        debugError('❌ [ThemeEditor] Error constructor:', err && (err as any).constructor?.name)
        debugError('❌ [ThemeEditor] Full error object:', err)

        let errorMessage = 'Failed to load theme'
        if (err instanceof Error) {
          errorMessage = err.message
          debugError('❌ [ThemeEditor] Error message:', err.message)
          debugError('❌ [ThemeEditor] Error stack:', err.stack)
        } else if (typeof err === 'string') {
          errorMessage = err
          debugError('❌ [ThemeEditor] String error:', err)
        } else if (err && typeof err === 'object') {
          const errObj = err as any
          errorMessage = errObj.message || errObj.error || JSON.stringify(err)
          debugError('❌ [ThemeEditor] Object error:', errorMessage)
        }

        debugError('❌ [ThemeEditor] Final error message to display:', errorMessage)
        setError(errorMessage)
      } finally {
        setIsLoading(false)
      }
    }

    loadTheme()
  }, [user?.restaurant_id, params.id, setCurrentTheme])

  const handleSave = async () => {
    if (!currentTheme) return

    try {
      setSaving(true)
      setError(null)

      // Transform components to match backend schema
      // Backend expects: id, type, title (BilingualText), order, enabled, settings (not config)
      const components = currentTheme.components.map((comp, index) => {
        const config = comp.config || {}
        const displayOrder = Number(
          comp.displayOrder !== undefined ? comp.displayOrder : (comp as any).order !== undefined ? (comp as any).order : index
        )

        // Extract bilingual fields from config
        const title = {
          en: config.title_en || config.title || '',
          ar: config.title_ar || '',
        }

        const subtitle = config.subtitle_en || config.subtitle_ar ? {
          en: config.subtitle_en || '',
          ar: config.subtitle_ar || '',
        } : undefined

        const description = config.description_en || config.description_ar ? {
          en: config.description_en || '',
          ar: config.description_ar || '',
        } : undefined

        const buttonText = config.button_text_en || config.button_text_ar ? {
          en: config.button_text_en || '',
          ar: config.button_text_ar || '',
        } : undefined

        const component: any = {
          id: comp.id ? String(comp.id) : `component-${index + 1}`,
          type: comp.type,
          title,
          order: displayOrder,
          displayOrder,
          enabled: comp.enabled,
          config,
          settings: config,
        }

        // Only include optional fields if they have values
        if (subtitle && (subtitle.en || subtitle.ar)) component.subtitle = subtitle
        if (description && (description.en || description.ar)) component.description = description
        if (buttonText && (buttonText.en || buttonText.ar)) component.buttonText = buttonText

        return component
      }) as any

      const themeConfig = {
        colors: currentTheme.colors,
        typography: {
          fontFamily: currentTheme.typography.fontFamily,
          baseFontSize: currentTheme.typography.baseFontSize,
          borderRadius: currentTheme.typography.borderRadius,
          lineHeight: currentTheme.typography.lineHeight,
          fontWeights: currentTheme.typography.fontWeights || [400, 700],
          fontStyle: currentTheme.typography.fontStyle || 'normal',
        },
        identity: currentTheme.identity,
        header: currentTheme.header,
        footer: currentTheme.footer,
        components,
        customCss: currentTheme.customCss || '',
      }

      const updateData = {
        name: currentTheme.name,
        slug: currentTheme.slug,
        description: currentTheme.description,
        colors: themeConfig.colors,
        typography: themeConfig.typography,
        identity: currentTheme.identity,
        header_config: currentTheme.header,
        footer_config: currentTheme.footer,
        components,
        customCss: currentTheme.customCss || '',
        config: themeConfig,
      }

      debugLog('📤 [ThemeEditor] Saving theme with data:', updateData)
      debugLog('📤 [ThemeEditor] Components being saved:', updateData.components)
      debugLog('📤 [ThemeEditor] Component count:', updateData.components?.length || 0)
      if (updateData.components && updateData.components.length > 0) {
        debugLog('📤 [ThemeEditor] First component details:', updateData.components[0])
      }
      await themeApi.updateTheme(params.id, updateData)
      setDirty(false)
      toast.success(`Theme "${currentTheme.name}" saved successfully`)
    } catch (err) {
      debugError('Failed to save theme:', err)
      const errorMessage = err instanceof Error ? err.message : 'Failed to save theme'
      setError(errorMessage)
      toast.error(errorMessage)
    } finally {
      setSaving(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin text-primary-600 mx-auto mb-4" />
          <p className="text-gray-600 dark:text-gray-400">Loading theme editor...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="p-6 max-w-4xl mx-auto">
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-semibold text-red-900 dark:text-red-300">Error</h3>
            <p className="text-red-700 dark:text-red-200 text-sm">{error}</p>
            <button
              onClick={() => window.location.reload()}
              className="mt-2 text-sm text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300 font-medium"
            >
              Retry
            </button>
          </div>
        </div>
      </div>
    )
  }

  if (!currentTheme) {
    return (
      <div className="p-6 max-w-4xl mx-auto">
        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
          <p className="text-yellow-700 dark:text-yellow-200">Theme not found</p>
        </div>
      </div>
    )
  }

  return (
    <ErrorBoundary
      onError={(error) => {
        debugError('Theme editor error:', error)
        setError(error.message)
      }}
    >
      <div className="w-full h-screen flex flex-col bg-gray-50 dark:bg-gray-900">
        {/* Error Bar */}
        {error && (
          <div className="bg-red-50 dark:bg-red-900/20 border-b border-red-200 dark:border-red-800 px-6 py-3 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <AlertCircle className="h-4 w-4 text-red-600 dark:text-red-400" />
              <p className="text-sm text-red-700 dark:text-red-200">{error}</p>
            </div>
            <button
              onClick={() => setError(null)}
              className="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300"
            >
              ×
            </button>
          </div>
        )}

        {/* Header - Sticky */}
        <div className="sticky top-0 z-40 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4 flex items-center justify-between shadow-sm">
          <div className="flex items-center gap-4 min-w-0">
            <button
              onClick={() => router.back()}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition flex-shrink-0"
              title="Go back"
            >
              <ArrowLeft className="h-5 w-5" />
            </button>
            <div className="min-w-0">
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white truncate">
                {extractLabel(currentTheme?.name)}
              </h1>
              <p className="text-sm text-gray-600 dark:text-gray-400">Editing theme</p>
            </div>
          </div>

          <div className="flex items-center gap-3 flex-shrink-0">
            {isDirty && (
              <button
                onClick={handleSave}
                disabled={isSaving}
                className="flex items-center gap-2 px-4 py-2 bg-primary-600 hover:bg-primary-700 disabled:bg-primary-400 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-colors"
                title="Save changes"
              >
                {isSaving ? (
                  <>
                    <Loader2 className="h-4 w-4 animate-spin" />
                    <span className="hidden sm:inline">Saving...</span>
                  </>
                ) : (
                  <>
                    <Save className="h-4 w-4" />
                    <span className="hidden sm:inline">Save</span>
                  </>
                )}
              </button>
            )}

            <button
              onClick={() => setSidebarOpen(!sidebarOpen)}
              className="md:hidden p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition"
              title={sidebarOpen ? 'Hide sidebar' : 'Show sidebar'}
            >
              {sidebarOpen ? (
                <X className="h-5 w-5" />
              ) : (
                <Menu className="h-5 w-5" />
              )}
            </button>
          </div>
        </div>

        {/* Main Content */}
        <div className="flex flex-1 overflow-hidden gap-0 w-full">
          {/* Sidebar - Collapsible with Smooth Transitions */}
          {sidebarOpen && (
            <div className="w-full md:w-96 overflow-hidden border-r border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 flex-shrink-0">
              <EditorSidebar />
            </div>
          )}

          {/* Preview Area - Full Width */}
          <div className="flex-1 overflow-hidden bg-gray-50 dark:bg-gray-900">
            <EditorPreview />
          </div>
        </div>
      </div>
    </ErrorBoundary>
  )
}
