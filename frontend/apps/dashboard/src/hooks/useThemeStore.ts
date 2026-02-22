'use client'

import { create } from 'zustand'
import {
  ThemeData,
  GlobalColors,
  TypographySettings,
  WebsiteIdentity,
  ComponentConfig,
  HeaderConfig,
  FooterConfig,
} from '@/types/theme'

interface ThemeStore {
  // State
  currentTheme: ThemeData | null
  themes: ThemeData[]
  selectedComponentId: string | null
  isDirty: boolean
  isSaving: boolean
  error: string | null

  // Actions
  setCurrentTheme: (theme: ThemeData) => void
  setThemes: (themes: ThemeData[]) => void
  updateThemeColors: (colors: Partial<GlobalColors>) => void
  updateTypography: (typography: Partial<TypographySettings>) => void
  updateIdentity: (identity: Partial<WebsiteIdentity>) => void
  updateHeader: (header: Partial<HeaderConfig>) => void
  updateFooter: (footer: Partial<FooterConfig>) => void
  updateCustomCss: (css: string) => void
  addComponent: (component: ComponentConfig) => void
  updateComponent: (id: string, updates: Partial<ComponentConfig>) => void
  deleteComponent: (id: string) => void
  reorderComponents: (components: ComponentConfig[]) => void
  selectComponent: (id: string | null) => void
  setDirty: (dirty: boolean) => void
  setSaving: (saving: boolean) => void
  setError: (error: string | null) => void
  reset: () => void
}

const DEFAULT_THEME: ThemeData = {
  id: '',
  restaurantId: '',
  tenantId: '',
  name: '',
  slug: '',
  version: 1,
  isActive: false,
  isPublished: false,
  colors: {
    primary: '#3b82f6',
    secondary: '#1e40af',
    accent: '#0ea5e9',
    background: '#ffffff',
    text: '#1f2937',
    border: '#e5e7eb',
    shadow: '#1f2937',
  },
  typography: {
    fontFamily: 'Inter',
    baseFontSize: 14,
    borderRadius: 8,
    lineHeight: 1.5,
    fontWeights: [400, 700],
    fontSize: 14,
    fontStyle: 'normal',
  },
  identity: {
    siteTitle: '',
  },
  header: {
    logoUrl: '',
    logoText: '',
    logoHeight: 40,
    showLogo: true,
    navigationItems: [],
    navPosition: 'right',
    navAlignment: 'horizontal',
    backgroundColor: '#3b82f6',
    textColor: '#ffffff',
    height: 64,
    padding: 16,
    stickyHeader: false,
    showShadow: true,
    hideNavOnMobile: false,
  },
  footer: {
    companyName: '',
    companyDescription: '',
    address: '',
    phone: '',
    email: '',
    copyrightText: '',
    socialLinks: [],
    footerSections: [],
    legalLinks: [],
    backgroundColor: '#1f2937',
    textColor: '#ffffff',
    linkColor: '#3b82f6',
    padding: 48,
    showLinks: true,
    showLegal: true,
    showBackToTop: true,
    columns: 3,
    layout: 'expanded',
  },
  components: [],
  customCss: '',
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
}

export const useThemeStore = create<ThemeStore>((set) => ({
  currentTheme: null,
  themes: [],
  selectedComponentId: null,
  isDirty: false,
  isSaving: false,
  error: null,

  setCurrentTheme: (theme) =>
    set({ currentTheme: theme, isDirty: false, selectedComponentId: null }),

  setThemes: (themes) => set({ themes }),

  updateThemeColors: (colorUpdates) =>
    set((state) => {
      if (!state.currentTheme) return state
      return {
        currentTheme: {
          ...state.currentTheme,
          colors: {
            ...state.currentTheme.colors,
            ...colorUpdates,
          },
        },
        isDirty: true,
      }
    }),

  updateTypography: (typographyUpdates) =>
    set((state) => {
      if (!state.currentTheme) return state
      return {
        currentTheme: {
          ...state.currentTheme,
          typography: {
            ...state.currentTheme.typography,
            ...typographyUpdates,
          },
        },
        isDirty: true,
      }
    }),

  updateIdentity: (identityUpdates) =>
    set((state) => {
      if (!state.currentTheme) return state
      return {
        currentTheme: {
          ...state.currentTheme,
          identity: {
            ...state.currentTheme.identity,
            ...identityUpdates,
          },
        },
        isDirty: true,
      }
    }),

  updateHeader: (headerUpdates) =>
    set((state) => {
      if (!state.currentTheme) return state
      return {
        currentTheme: {
          ...state.currentTheme,
          header: {
            ...state.currentTheme.header,
            ...headerUpdates,
          },
        },
        isDirty: true,
      }
    }),

  updateFooter: (footerUpdates) =>
    set((state) => {
      if (!state.currentTheme) return state
      return {
        currentTheme: {
          ...state.currentTheme,
          footer: {
            ...state.currentTheme.footer,
            ...footerUpdates,
          },
        },
        isDirty: true,
      }
    }),

  updateCustomCss: (css) =>
    set((state) => {
      if (!state.currentTheme) return state
      return {
        currentTheme: {
          ...state.currentTheme,
          customCss: css,
        },
        isDirty: true,
      }
    }),

  addComponent: (component) =>
    set((state) => {
      if (!state.currentTheme) return state
      const newComponent = {
        ...component,
        id: component.id ? String(component.id) : `component-${state.currentTheme.components.length + 1}`,
        displayOrder: state.currentTheme.components.length,
      }
      return {
        currentTheme: {
          ...state.currentTheme,
          components: [...state.currentTheme.components, newComponent],
        },
        isDirty: true,
      }
    }),

  updateComponent: (id, updates) =>
    set((state) => {
      if (!state.currentTheme) return state
      return {
        currentTheme: {
          ...state.currentTheme,
          components: state.currentTheme.components.map((comp) =>
            String(comp.id) === String(id) ? { ...comp, ...updates } : comp
          ),
        },
        isDirty: true,
      }
    }),

  deleteComponent: (id) =>
    set((state) => {
      if (!state.currentTheme) return state
      return {
        currentTheme: {
          ...state.currentTheme,
          components: state.currentTheme.components.filter((comp) => String(comp.id) !== String(id)),
        },
        isDirty: true,
        selectedComponentId: state.selectedComponentId === String(id) ? null : state.selectedComponentId,
      }
    }),

  reorderComponents: (components) =>
    set((state) => {
      if (!state.currentTheme) return state
      return {
        currentTheme: {
          ...state.currentTheme,
          components: components.map((comp, idx) => ({
            ...comp,
            id: comp.id ? String(comp.id) : `component-${idx + 1}`,
            displayOrder: idx,
            ...(typeof (comp as any).order === 'number' ? { order: idx } : {}),
          })),
        },
        isDirty: true,
      }
    }),

  selectComponent: (id) => set({ selectedComponentId: id }),

  setDirty: (dirty) => set({ isDirty: dirty }),

  setSaving: (saving) => set({ isSaving: saving }),

  setError: (error) => set({ error }),

  reset: () =>
    set({
      currentTheme: null,
      themes: [],
      selectedComponentId: null,
      isDirty: false,
      isSaving: false,
      error: null,
    }),
}))
