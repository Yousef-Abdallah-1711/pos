import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import {
  ThemeData,
  GlobalColors,
  TypographySettings,
  HeaderConfig,
  FooterConfig,
  ComponentConfig,
  WebsiteIdentity,
} from '@/types/theme'

export interface ThemeBuilderState {
  // Current theme being edited
  currentTheme: Partial<ThemeData> | null
  isLoading: boolean
  isSaving: boolean
  error: string | null
  successMessage: string | null
  isDirty: boolean

  // Form data
  colors: GlobalColors
  typography: TypographySettings
  identity: WebsiteIdentity
  header: HeaderConfig
  footer: FooterConfig
  components: ComponentConfig[]
  customCss: string

  // Actions - Update fields
  updateColors: (colors: Partial<GlobalColors>) => void
  updateTypography: (typography: Partial<TypographySettings>) => void
  updateIdentity: (identity: Partial<WebsiteIdentity>) => void
  updateHeader: (header: Partial<HeaderConfig>) => void
  updateFooter: (footer: Partial<FooterConfig>) => void
  updateComponents: (components: ComponentConfig[]) => void
  updateCustomCss: (css: string) => void

  // Actions - Theme management
  loadTheme: (theme: ThemeData) => void
  resetTheme: () => void
  markDirty: () => void
  markClean: () => void

  // Actions - Status
  setLoading: (loading: boolean) => void
  setSaving: (saving: boolean) => void
  setError: (error: string | null) => void
  setSuccess: (message: string | null) => void

  // Validation
  validateColors: () => boolean
  validateTypography: () => boolean
  validateIdentity: () => boolean
  validateHeader: () => boolean
  validateFooter: () => boolean
  validateTheme: () => boolean

  // Selectors
  getColorValue: (colorKey: keyof GlobalColors) => string
  getComponentById: (id: string) => ComponentConfig | undefined
  getFormData: () => Partial<ThemeData>
  hasErrors: () => boolean
}

const defaultColors: GlobalColors = {
  primary: '#007bff',
  secondary: '#6c757d',
  accent: '#ff6b6b',
  background: '#ffffff',
  text: '#212529',
  border: '#dee2e6',
  shadow: '#00000020',
}

const defaultTypography: TypographySettings = {
  fontFamily: 'Segoe UI, Tahoma, Geneva, Verdana, sans-serif',
  baseFontSize: 16,
  borderRadius: 4,
  lineHeight: 1.5,
  fontWeights: [400, 500, 600, 700],
}

const defaultIdentity: WebsiteIdentity = {
  siteTitle: 'My Restaurant',
  logoUrl: undefined,
  faviconUrl: undefined,
  domain: undefined,
}

const defaultHeader: HeaderConfig = {
  showLogo: true,
  navigationItems: [
    { id: '1', label: 'Menu', href: '/menu', order: 1 },
    { id: '2', label: 'About', href: '/about', order: 2 },
    { id: '3', label: 'Contact', href: '/contact', order: 3 },
  ],
  navPosition: 'right',
  navAlignment: 'horizontal',
  backgroundColor: '#ffffff',
  textColor: '#212529',
  height: 80,
  padding: 16,
  stickyHeader: true,
  showShadow: true,
  hideNavOnMobile: false,
}

const defaultFooter: FooterConfig = {
  companyName: 'My Restaurant',
  companyDescription: 'Welcome to our restaurant',
  address: '123 Main Street',
  phone: '+1 (555) 000-0000',
  email: 'info@restaurant.com',
  copyrightText: `© ${new Date().getFullYear()} My Restaurant. All rights reserved.`,
  socialLinks: [],
  footerSections: [],
  legalLinks: [
    { label: 'Privacy Policy', href: '/privacy' },
    { label: 'Terms of Service', href: '/terms' },
  ],
  backgroundColor: '#f8f9fa',
  textColor: '#212529',
  linkColor: '#007bff',
  padding: 24,
  showLinks: true,
  showLegal: true,
  showBackToTop: true,
  columns: 3,
  layout: 'expanded',
}

export const useThemeBuilderStore = create<ThemeBuilderState>()(
  persist(
    (set, get) => ({
      currentTheme: null,
      isLoading: false,
      isSaving: false,
      error: null,
      successMessage: null,
      isDirty: false,

      colors: defaultColors,
      typography: defaultTypography,
      identity: defaultIdentity,
      header: defaultHeader,
      footer: defaultFooter,
      components: [],
      customCss: '',

      // Update actions
      updateColors: (colors) => {
        set((state) => ({
          colors: { ...state.colors, ...colors },
          isDirty: true,
        }))
      },

      updateTypography: (typography) => {
        set((state) => ({
          typography: { ...state.typography, ...typography },
          isDirty: true,
        }))
      },

      updateIdentity: (identity) => {
        set((state) => ({
          identity: { ...state.identity, ...identity },
          isDirty: true,
        }))
      },

      updateHeader: (header) => {
        set((state) => ({
          header: { ...state.header, ...header },
          isDirty: true,
        }))
      },

      updateFooter: (footer) => {
        set((state) => ({
          footer: { ...state.footer, ...footer },
          isDirty: true,
        }))
      },

      updateComponents: (components) => {
        set({
          components,
          isDirty: true,
        })
      },

      updateCustomCss: (css) => {
        set({
          customCss: css,
          isDirty: true,
        })
      },

      // Theme management
      loadTheme: (theme) => {
        set({
          currentTheme: theme,
          colors: theme.colors || defaultColors,
          typography: theme.typography || defaultTypography,
          identity: theme.identity || theme.websiteIdentity || defaultIdentity,
          header: theme.header || defaultHeader,
          footer: theme.footer || defaultFooter,
          components: theme.components || [],
          customCss: theme.customCss || '',
          isDirty: false,
          error: null,
        })
      },

      resetTheme: () => {
        set({
          currentTheme: null,
          colors: defaultColors,
          typography: defaultTypography,
          identity: defaultIdentity,
          header: defaultHeader,
          footer: defaultFooter,
          components: [],
          customCss: '',
          isDirty: false,
          error: null,
          successMessage: null,
        })
      },

      markDirty: () => set({ isDirty: true }),
      markClean: () => set({ isDirty: false }),

      // Status actions
      setLoading: (loading) => set({ isLoading: loading }),
      setSaving: (saving) => set({ isSaving: saving }),
      setError: (error) => set({ error }),
      setSuccess: (successMessage) => set({ successMessage }),

      // Validation methods
      validateColors: () => {
        const state = get()
        const { colors } = state

        if (!colors) return false

        const hexRegex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/
        const requiredColors: (keyof GlobalColors)[] = [
          'primary',
          'secondary',
          'background',
          'text',
        ]

        return requiredColors.every((key) => {
          const value = colors[key]
          return value && hexRegex.test(value)
        })
      },

      validateTypography: () => {
        const state = get()
        const { typography } = state

        if (!typography) return false

        return (
          Boolean(typography.fontFamily) &&
          typography.fontFamily.trim().length > 0 &&
          typography.baseFontSize > 0 &&
          typography.baseFontSize < 100 &&
          typography.borderRadius >= 0 &&
          typography.borderRadius <= 50 &&
          typography.lineHeight > 0 &&
          typography.lineHeight < 5
        )
      },

      validateIdentity: () => {
        const state = get()
        const { identity } = state

        if (!identity) return false

        return Boolean(identity.siteTitle) && identity.siteTitle.trim().length > 0
      },

      validateHeader: () => {
        const state = get()
        const { header } = state

        if (!header) return false

        const hexRegex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/

        return (
          hexRegex.test(header.backgroundColor) &&
          hexRegex.test(header.textColor) &&
          header.height > 0 &&
          header.height < 200 &&
          header.navigationItems &&
          header.navigationItems.length > 0
        )
      },

      validateFooter: () => {
        const state = get()
        const { footer } = state

        if (!footer) return false

        const hexRegex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/

        return (
          Boolean(footer.companyName) &&
          footer.companyName.trim().length > 0 &&
          hexRegex.test(footer.backgroundColor) &&
          hexRegex.test(footer.textColor) &&
          Boolean(footer.copyrightText) &&
          footer.copyrightText.trim().length > 0
        )
      },

      validateTheme: () => {
        const state = get()
        return (
          state.validateColors() &&
          state.validateTypography() &&
          state.validateIdentity() &&
          state.validateHeader() &&
          state.validateFooter()
        )
      },

      // Selector methods
      getColorValue: (colorKey) => {
        const state = get()
        return state.colors[colorKey] || defaultColors[colorKey]
      },

      getComponentById: (id) => {
        const state = get()
        return state.components.find((comp) => comp.id === id)
      },

      getFormData: () => {
        const state = get()
        return {
          name: state.currentTheme?.name || '',
          slug: state.currentTheme?.slug || '',
          description: state.currentTheme?.description || '',
          colors: state.colors,
          typography: state.typography,
          identity: state.identity,
          header: state.header,
          footer: state.footer,
          components: state.components,
          customCss: state.customCss,
        }
      },

      hasErrors: () => {
        const state = get()
        return state.error !== null && state.error.length > 0
      },
    }),
    {
      name: 'theme-builder-storage',
      // Only persist essential state, not transient flags
      partialize: (state) => ({
        currentTheme: state.currentTheme,
        colors: state.colors,
        typography: state.typography,
        identity: state.identity,
        header: state.header,
        footer: state.footer,
        components: state.components,
        customCss: state.customCss,
      }),
    }
  )
)
