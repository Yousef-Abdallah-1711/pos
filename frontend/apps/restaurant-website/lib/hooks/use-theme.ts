'use client'

import { useEffect } from 'react'
import { useThemeStore } from '@/lib/store/theme-store'
import type { ThemeData } from '@/lib/store/theme-store'

// ============================================================================
// TypeScript Interfaces (support both camelCase and snake_case from API)
// ============================================================================

export interface ThemeIdentityData {
  siteTitle: string;
  logoUrl: string;
  faviconUrl: string;
  isLoading: boolean;
}

export interface ThemeHeaderData {
  backgroundColor: string;
  textColor: string;
  tagline: string;
  height: number;
  navigationItems: Array<{
    id?: string;
    label: string;
    href: string;
    order?: number;
  }>;
  isLoading: boolean;
}

export interface ThemeFooterData {
  backgroundColor: string;
  textColor: string;
  companyDescription: string;
  address: string;
  phone: string;
  email: string;
  copyrightText: string;
  socialLinks: Array<{
    platform: string;
    url: string;
  }>;
  footerSections: Array<{
    id: string;
    title: string;
    links: Array<{ label: string; href: string }>;
  }>;
  legalLinks: Array<{ label: string; href: string }>;
  isLoading: boolean;
}

export interface ThemeColorsData {
  primary: string;
  secondary: string;
  accent: string;
  background: string;
  text: string;
  border: string;
  shadow: string;
  isLoading: boolean;
}

// ============================================================================
// Core Theme Hook
// ============================================================================

/**
 * Custom hook to access theme from Zustand store
 * Automatically loads theme on first use
 * Caching is handled internally by the store (memory cache -> localStorage -> API)
 */
export function useTheme(defaultSlug: string = 'warm-comfort') {
  const currentTheme = useThemeStore((state) => state.currentTheme)
  const isLoading = useThemeStore((state) => state.isLoading)
  const error = useThemeStore((state) => state.error)
  const loadTheme = useThemeStore((state) => state.loadTheme)
  const setTheme = useThemeStore((state) => state.setTheme)
  const clearError = useThemeStore((state) => state.clearError)

  useEffect(() => {
    if (currentTheme) {
      return
    }
    void loadTheme(defaultSlug)
  }, [currentTheme, defaultSlug, loadTheme])

  return {
    currentTheme,
    isLoading,
    error,
    loadTheme,
    setTheme,
    clearError,
  }
}

// ============================================================================
// Typed Theme Hooks with Loading States
// ============================================================================

/**
 * Hook to access current theme identity (site title, logo, favicon)
 * Returns properly typed data with loading state
 * Note: isLoading is true when no theme is loaded yet (initial state)
 */
export function useThemeIdentity(): ThemeIdentityData {
  const { currentTheme, isLoading } = useTheme()
  const identity = currentTheme?.identity as any || {}

  // Consider "no theme yet" as loading
  const stillLoading = isLoading || !currentTheme

  return {
    siteTitle: identity.site_title || identity.siteTitle || 'Restaurant',
    logoUrl: identity.logo_url || identity.logoUrl || '',
    faviconUrl: identity.favicon_url || identity.faviconUrl || '',
    isLoading: stillLoading,
  }
}

/**
 * Hook to access current theme header config
 * Returns properly typed data with loading state
 * Note: isLoading is true when no theme is loaded yet (initial state)
 */
export function useThemeHeader(): ThemeHeaderData {
  const { currentTheme, isLoading } = useTheme()
  const header = currentTheme?.header as any || {}

  // Consider "no theme yet" as loading
  const stillLoading = isLoading || !currentTheme

  return {
    backgroundColor: header.background_color || header.backgroundColor || '#ffffff',
    textColor: header.text_color || header.textColor || '#1f2937',
    tagline: header.tagline || '',
    height: header.height || header.nav_height || 64,
    navigationItems: header.navigation_items || header.navigationItems || [],
    isLoading: stillLoading,
  }
}

/**
 * Hook to access current theme footer config
 * Returns properly typed data with loading state
 * Note: isLoading is true when no theme is loaded yet (initial state)
 */
export function useThemeFooter(): ThemeFooterData {
  const { currentTheme, isLoading } = useTheme()
  const footer = currentTheme?.footer as any || {}

  // Consider "no theme yet" as loading
  const stillLoading = isLoading || !currentTheme

  return {
    backgroundColor: footer.background_color || footer.backgroundColor || '#111827',
    textColor: footer.text_color || footer.textColor || '#f3f4f6',
    companyDescription: footer.company_description || footer.companyDescription || '',
    address: footer.address || '',
    phone: footer.phone || '',
    email: footer.email || '',
    copyrightText: footer.copyright_text || footer.copyrightText || '',
    socialLinks: footer.social_links || footer.socialLinks || [],
    footerSections: footer.footer_sections || footer.footerSections || footer.sections || [],
    legalLinks: footer.legal_links || footer.legalLinks || [],
    isLoading: stillLoading,
  }
}

/**
 * Hook to access current theme colors
 * Returns properly typed data with loading state
 * Note: isLoading is true when no theme is loaded yet (initial state)
 */
export function useThemeColors(): ThemeColorsData {
  const { currentTheme, isLoading } = useTheme()
  const colors = currentTheme?.colors as any || {}

  // Consider "no theme yet" as loading
  const stillLoading = isLoading || !currentTheme

  return {
    primary: colors.primary || '#f97316',
    secondary: colors.secondary || '#0ea5e9',
    accent: colors.accent || '#fbbf24',
    background: colors.background || '#ffffff',
    text: colors.text || '#1f2937',
    border: colors.border || '#e5e7eb',
    shadow: colors.shadow || '#1f2937',
    isLoading: stillLoading,
  }
}

/**
 * Hook to access current theme typography
 */
export function useThemeTypography() {
  const { currentTheme, isLoading } = useTheme()
  const typography = currentTheme?.typography as any || {}

  return {
    fontFamily: typography.font_family || typography.fontFamily || 'Inter, sans-serif',
    baseFontSize: typography.base_font_size || typography.baseFontSize || 14,
    borderRadius: typography.border_radius || typography.borderRadius || 8,
    lineHeight: typography.line_height || typography.lineHeight || 1.5,
    isLoading,
  }
}

// ============================================================================
// Utility Hooks
// ============================================================================

/**
 * Hook to load a specific theme by slug
 */
export function useLoadTheme(slug: string) {
  const store = useThemeStore()

  const load = async () => {
    await store.loadTheme(slug)
  }

  return {
    isLoading: store.isLoading,
    error: store.error,
    load,
  }
}

/**
 * Hook to switch theme
 */
export function useSwitchTheme() {
  const store = useThemeStore()

  const switchTheme = async (slug: string) => {
    await store.loadTheme(slug)
  }

  return switchTheme
}

/**
 * Hook to check if theme is loading
 */
export function useThemeLoading() {
  const { isLoading } = useTheme()
  return isLoading
}

/**
 * Hook to check if theme has error
 */
export function useThemeError() {
  const { error } = useTheme()
  return error
}

/**
 * Hook to get theme loading progress
 */
export function useThemeProgress() {
  const { isLoading, error } = useTheme()
  return { isLoading, error }
}

/**
 * Hook to check if theme is ready (loaded and no error)
 */
export function useThemeReady() {
  const { currentTheme, isLoading, error } = useTheme()
  return {
    isReady: !!currentTheme && !isLoading && !error,
    isLoading,
    hasError: !!error,
    error,
  }
}

/**
 * Type exports
 */
export type { ThemeData }
