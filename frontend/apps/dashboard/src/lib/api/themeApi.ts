/**
 * Theme API Service
 * Handles all API calls for theme management
 *
 * ⚠️ IMPORTANT: This service reads JWT token from Zustand auth-storage, not localStorage directly
 * The token is stored in Zustand's persist middleware under the key 'auth-storage'
 */

import axios, { AxiosError } from 'axios'
import { ThemeData, CreateThemeRequest, UpdateThemeRequest, ComponentConfig, ComponentType } from '@/types/theme'
import { debugLog, debugWarn, debugError } from '@/lib/debug'
import { getConfig } from '@/lib/config'

const API_BASE_URL_RAW = getConfig().apiUrl
const API_BASE_URL = API_BASE_URL_RAW.replace(/\/+$/, '')
const API_ROOT = API_BASE_URL.endsWith('/api/v1') ? API_BASE_URL : `${API_BASE_URL}/api/v1`

// Create axios instance with default config
const apiClient = axios.create({
  baseURL: API_ROOT,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
})

// ==============================================
// Authentication is carried by the API's HttpOnly session cookie.
// ==============================================
apiClient.interceptors.request.use((config) => {
  return config
})

// ==============================================
// ⚠️ ERROR INTERCEPTOR: Handle auth failures
// Clears auth and redirects to login on 401
// ==============================================
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    debugError('❌ [API Response] Error:', error.response?.status, error.response?.statusText)
    debugError('❌ [API Response] Full Data:', error.response?.data)

    // More detailed error logging
    if (error.response?.data) {
      const errorData = error.response.data as any
      debugError('❌ [API Response] Error Message:', errorData.error || errorData.message || 'Unknown error')
      debugError('❌ [API Response] Full Error Object:', JSON.stringify(errorData, null, 2))
    }

    // Handle 401 Unauthorized - session expired or invalid token
    if (error.response?.status === 401) {
      debugLog('🔑 [API Response] 401 Unauthorized - Clearing auth and redirecting to login')
      // Clear the auth storage
      localStorage.removeItem('auth-storage')
      // Redirect to login
      if (typeof window !== 'undefined') {
        window.location.href = '/en/login'
      }
    }

    return Promise.reject(error)
  }
)

// Types
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
  message?: string
}

export interface ApiError {
  message: string
  status?: number
  originalError?: AxiosError
}

const COMPONENT_TYPE_MAP: Record<string, ComponentType> = {
  hero: 'hero',
  products: 'products',
  'featured-items': 'products',
  featured_items: 'products',
  why_us: 'why_us',
  'why-us': 'why_us',
  'why-choose-us': 'why_us',
  why_choose_us: 'why_us',
  contact: 'contact',
  testimonials: 'testimonials',
  cta: 'cta',
  'call-to-action': 'cta',
  custom: 'custom',
}

const DEFAULT_COLORS = {
  primary: '#3b82f6',
  secondary: '#1e40af',
  accent: '#0ea5e9',
  background: '#ffffff',
  text: '#1f2937',
  border: '#e5e7eb',
  shadow: '#1f2937',
}

function extractText(value: any, fallback: string = ''): string {
  if (typeof value === 'string') return value
  if (value && typeof value === 'object') {
    return value.en || value.ar || Object.values(value)[0] || fallback
  }
  return fallback
}

function normalizeDate(value: any): string {
  if (!value) return new Date(0).toISOString()
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? new Date(0).toISOString() : date.toISOString()
}

function parseConfig(config: unknown): Record<string, any> {
  if (!config) return {}
  if (typeof config === 'string') {
    try {
      return JSON.parse(config)
    } catch (error) {
      debugWarn('⚠️ [themeApi] Failed to parse config JSON string:', error)
      return {}
    }
  }
  if (typeof config === 'object') {
    return config as Record<string, any>
  }
  return {}
}

function normalizeNavItems(items: any): Array<{ id?: string; label: string; href: string; order?: number }> {
  if (!Array.isArray(items)) return []
  return items.map((item: any, index: number) => ({
    id: item?.id ? String(item.id) : `${index + 1}`,
    label: extractText(item?.label, ''),
    href: typeof item?.href === 'string' ? item.href : '/',
    order: typeof item?.order === 'number' ? item.order : index + 1,
  }))
}

function normalizeSocialLinks(items: any): Array<{ id?: string; platform: string; url: string; order?: number }> {
  if (!Array.isArray(items)) return []
  return items.map((item: any, index: number) => ({
    id: item?.id ? String(item.id) : `${index + 1}`,
    platform: extractText(item?.platform, 'facebook').toLowerCase(),
    url: typeof item?.url === 'string' ? item.url : '',
    order: typeof item?.order === 'number' ? item.order : index + 1,
  }))
}

function normalizeFooterSections(items: any): Array<{ id?: string; title: string; links: Array<{ label: string; href: string }>; order?: number }> {
  if (!Array.isArray(items)) return []
  return items.map((section: any, index: number) => ({
    id: section?.id ? String(section.id) : `${index + 1}`,
    title: extractText(section?.title, `Section ${index + 1}`),
    links: Array.isArray(section?.links)
      ? section.links.map((link: any) => ({
          label: extractText(link?.label, 'Link'),
          href: typeof link?.href === 'string' ? link.href : '#',
        }))
      : [],
    order: typeof section?.order === 'number' ? section.order : index + 1,
  }))
}

function normalizeComponents(components: any): ComponentConfig[] {
  if (!Array.isArray(components)) return []

  return components.map((component: any, index: number) => {
    const mappedType = COMPONENT_TYPE_MAP[String(component?.type || '').toLowerCase()] || 'custom'
    const order = Number.isFinite(Number(component?.displayOrder))
      ? Number(component.displayOrder)
      : Number.isFinite(Number(component?.order))
        ? Number(component.order)
        : index

    return {
      id: component?.id ? String(component.id) : `component-${index + 1}`,
      type: mappedType,
      title: extractText(component?.title, mappedType),
      displayOrder: order,
      enabled: component?.enabled !== false,
      config: component?.config || component?.settings || {},
      settings: component?.settings || component?.config || {},
      styles: component?.styles || {},
      createdAt: normalizeDate(component?.createdAt || component?.created_at),
      updatedAt: normalizeDate(component?.updatedAt || component?.updated_at),
    } as ComponentConfig
  })
}

function normalizeTheme(rawTheme: any): ThemeData {
  const parsedConfig = parseConfig(rawTheme?.config)
  const merged = {
    ...(rawTheme || {}),
    ...parsedConfig,
  }

  const colors = {
    ...DEFAULT_COLORS,
    ...(merged.colors || {}),
  }

  const typographySource = merged.typography || {}
  const identitySource = merged.identity || {}
  const headerSource = merged.header || merged.header_config || {}
  const footerSource = merged.footer || merged.footer_config || {}

  const identity = {
    siteTitle: extractText(identitySource.siteTitle || identitySource.site_title, extractText(merged.name, 'Restaurant')),
    logoUrl: identitySource.logoUrl || identitySource.logo_url || '',
    faviconUrl: identitySource.faviconUrl || identitySource.favicon_url || '',
    domain: identitySource.domain || '',
  }

  const header = {
    logoUrl: headerSource.logoUrl || headerSource.logo_url || identity.logoUrl || '',
    logoText: extractText(headerSource.logoText || headerSource.logo_text, identity.siteTitle),
    logoHeight: Number(headerSource.logoHeight || headerSource.logo_height || 40),
    showLogo: headerSource.showLogo !== false,
    navigationItems: normalizeNavItems(headerSource.navigationItems || headerSource.navigation_items),
    navPosition: headerSource.navPosition || headerSource.nav_position || 'right',
    navAlignment: headerSource.navAlignment || headerSource.nav_alignment || 'horizontal',
    backgroundColor: headerSource.backgroundColor || headerSource.background_color || colors.primary,
    textColor: headerSource.textColor || headerSource.text_color || '#ffffff',
    height: Number(headerSource.height || headerSource.navHeight || headerSource.nav_height || 64),
    padding: Number(headerSource.padding || 16),
    stickyHeader: Boolean(headerSource.stickyHeader || headerSource.sticky),
    showShadow: headerSource.showShadow !== false && headerSource.shadow !== false,
    hideNavOnMobile: Boolean(headerSource.hideNavOnMobile || headerSource.hide_nav_on_mobile),
  }

  const footer = {
    companyName: extractText(
      footerSource.companyName || footerSource.company_name || footerSource.companyInfo?.name,
      identity.siteTitle
    ),
    companyDescription: extractText(
      footerSource.companyDescription || footerSource.company_description || footerSource.companyInfo?.description,
      ''
    ),
    address: extractText(footerSource.address || footerSource.companyInfo?.address, ''),
    phone: extractText(footerSource.phone || footerSource.companyInfo?.phone, ''),
    email: extractText(footerSource.email || footerSource.companyInfo?.email, ''),
    copyrightText: extractText(
      footerSource.copyrightText || footerSource.copyright_text,
      `© ${new Date().getFullYear()} ${identity.siteTitle}. All rights reserved.`
    ),
    socialLinks: normalizeSocialLinks(footerSource.socialLinks || footerSource.social_links),
    footerSections: normalizeFooterSections(
      footerSource.footerSections || footerSource.footer_sections || footerSource.sections
    ),
    legalLinks: Array.isArray(footerSource.legalLinks || footerSource.legal_links)
      ? (footerSource.legalLinks || footerSource.legal_links).map((link: any) => ({
          label: extractText(link?.label, 'Link'),
          href: typeof link?.href === 'string' ? link.href : '#',
        }))
      : [],
    backgroundColor: footerSource.backgroundColor || footerSource.background_color || '#111827',
    textColor: footerSource.textColor || footerSource.text_color || '#f3f4f6',
    linkColor: footerSource.linkColor || footerSource.link_color || colors.accent,
    padding: Number(footerSource.padding || 48),
    showLinks: footerSource.showLinks !== false,
    showLegal: footerSource.showLegal !== false,
    showBackToTop: footerSource.showBackToTop !== false,
    columns: Number(footerSource.columns || 3),
    layout: footerSource.layout || 'expanded',
  }

  return {
    id: String(merged.id || merged.themeId || ''),
    restaurantId: String(merged.restaurantId || merged.restaurant_id || ''),
    tenantId: String(merged.tenantId || merged.tenant_id || ''),
    name: extractText(merged.name, 'Untitled Theme'),
    slug: extractText(merged.slug, ''),
    description: extractText(merged.description, ''),
    version: Number(merged.version || 1),
    isActive: Boolean(merged.isActive || merged.is_active),
    isPublished: Boolean(merged.isPublished || merged.is_published),
    colors,
    typography: {
      fontFamily: typographySource.fontFamily || typographySource.font_family || typographySource.fontSans || 'Inter',
      baseFontSize: Number(typographySource.baseFontSize || typographySource.base_size || typographySource.base_font_size || 14),
      borderRadius: Number(typographySource.borderRadius || typographySource.border_radius || 8),
      lineHeight: Number(typographySource.lineHeight || typographySource.line_height || 1.5),
      fontWeights: Array.isArray(typographySource.fontWeights) ? typographySource.fontWeights : [400, 700],
      fontSize: Number(typographySource.fontSize || typographySource.baseFontSize || 14),
      fontStyle: typographySource.fontStyle || 'normal',
    },
    identity,
    websiteIdentity: identity,
    header,
    footer,
    components: normalizeComponents(merged.components || merged.homepage?.sections),
    customCss: merged.customCss || merged.custom_css || '',
    createdAt: normalizeDate(merged.createdAt || merged.created_at),
    updatedAt: normalizeDate(merged.updatedAt || merged.updated_at),
    publishedAt: merged.publishedAt || merged.published_at || undefined,
    themeId: merged.themeId ? String(merged.themeId) : undefined,
  }
}

// Theme API endpoints

/**
 * Get all theme presets (templates available for selection)
 */
export async function getThemePresets(limit: number = 100): Promise<any[]> {
  try {
    debugLog('📥 [getThemePresets] Fetching theme presets with limit:', limit)
    const response = await apiClient.get<any>(
      `/admin/theme-presets?limit=${limit}&page=1`
    )

    debugLog('📦 [getThemePresets] Full API Response:', response.data)

    let data = response.data?.data || response.data?.presets || []
    debugLog('✅ [getThemePresets] Got', Array.isArray(data) ? data.length : 0, 'presets')

    // Parse config JSON for each preset
    if (Array.isArray(data)) {
      data = data.map((preset: any) => {
        // If config is a string (JSON), parse it and merge it into the preset
        if (preset.config && typeof preset.config === 'string') {
          try {
            const parsedConfig = JSON.parse(preset.config)
            return {
              ...preset,
              ...parsedConfig,
              _rawConfig: preset.config,
            }
          } catch (err) {
            debugWarn('⚠️ [getThemePresets] Failed to parse preset config:', err)
            return preset
          }
        }
        return preset
      })
    }

    return Array.isArray(data) ? data : []
  } catch (error) {
    debugWarn('⚠️ [getThemePresets] Failed to fetch presets:', error)
    // Return empty array instead of throwing to prevent UI crashes
    return []
  }
}

/**
 * Get all themes for the restaurant (fetches all pages)
 */
export async function getThemes(limit: number = 100): Promise<ThemeData[]> {
  try {
    debugLog('📥 [getThemes] Fetching themes with limit:', limit)
    const response = await apiClient.get<any>(
      `/admin/themes?pageSize=${limit}&page=1`
    )
    
    // Log full response for debugging
    debugLog('📦 [getThemes] Full API Response:', response.data)
    
    // Handle all known response formats
    const data = response.data?.data || response.data?.themes || response.data?.items || []
    const pagination = response.data?.pagination || {
      total: response.data?.total || (Array.isArray(data) ? data.length : 0),
      totalPages: response.data?.total_pages || response.data?.totalPages || 1,
    }
    
    debugLog('✅ [getThemes] Response structure:', {
      totalThemes: pagination.total || 0,
      currentCount: Array.isArray(data) ? data.length : 0,
      pagination: pagination,
      allKeys: Object.keys(response.data || {})
    })
    
    // If we got all themes in one page, return them
    if (Array.isArray(data)) {
      // If there are more pages and we didn't fetch all, fetch remaining
      if (pagination.totalPages && pagination.totalPages > 1 && data.length < pagination.total) {
        debugLog(`📄 [getThemes] Multiple pages detected. Total: ${pagination.total}, Fetched: ${data.length}`)
        const allThemes = [...data]
        
        // Fetch remaining pages
        for (let page = 2; page <= pagination.totalPages; page++) {
          const pageResponse = await apiClient.get<any>(
            `/admin/themes?pageSize=${limit}&page=${page}`
          )
          const pageData = pageResponse.data?.data || pageResponse.data?.themes || pageResponse.data?.items || []
          if (Array.isArray(pageData)) {
            allThemes.push(...pageData)
          }
        }
        
        debugLog(`✅ [getThemes] Fetched all ${allThemes.length} themes across ${pagination.totalPages} pages`)
        return allThemes.map(normalizeTheme)
      }
      
      return data.map(normalizeTheme)
    }
    
    return []
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Get a single theme by ID
 */
export async function getTheme(themeId: string): Promise<ThemeData> {
  try {
    debugLog('📥 [getTheme] Fetching theme:', themeId)
    // Add timestamp to bypass cache
    const response = await apiClient.get<any>(
      `/admin/themes/${themeId}?_t=${Date.now()}`
    )

    // Log the full response for debugging
    debugLog('📦 [getTheme] Response status:', response.status)
    debugLog('📦 [getTheme] Response.data type:', typeof response.data)
    debugLog('📦 [getTheme] Response.data:', response.data)
    debugLog('📦 [getTheme] Response keys:', response.data ? Object.keys(response.data) : 'no data')

    // Handle different response formats - data, theme, or direct object
    let theme = response.data?.data || response.data?.theme || response.data

    // Validate theme exists
    if (!theme) {
      const errorMsg = `Theme not found in response. Response was: ${JSON.stringify(response.data)}`
      debugError('❌ [getTheme]', errorMsg)
      throw new Error(errorMsg)
    }

    // Validate theme is not an empty object
    if (typeof theme === 'object' && Object.keys(theme).length === 0) {
      const errorMsg = `Theme object is empty. Full response: ${JSON.stringify(response.data)}`
      debugError('❌ [getTheme]', errorMsg)
      throw new Error(errorMsg)
    }

    // Validate theme has required properties
    if (!theme.id && !theme.themeId) {
      const errorMsg = `Theme missing id. Theme: ${JSON.stringify(theme)}`
      debugError('❌ [getTheme]', errorMsg)
      throw new Error(errorMsg)
    }

    // Normalize all known API variants (snake_case/camelCase/config payloads)
    const normalizedTheme = normalizeTheme(theme)

    debugLog('✅ [getTheme] Received theme:', {
      id: normalizedTheme.id,
      name: normalizedTheme.name,
      components: Array.isArray(normalizedTheme.components) ? normalizedTheme.components.length : 0,
      hasColors: !!normalizedTheme.colors,
      hasTypography: !!normalizedTheme.typography,
    })
    return normalizedTheme
  } catch (error) {
    // Log comprehensive error info
    debugError('❌ [getTheme] Caught error')
    debugError('❌ [getTheme] Error instanceof Error:', error instanceof Error)
    debugError('❌ [getTheme] Error type:', typeof error)
    debugError('❌ [getTheme] Error constructor:', error && (error as any).constructor?.name)

    if (axios.isAxiosError(error)) {
      debugError('❌ [getTheme] Axios error status:', error.response?.status)
      debugError('❌ [getTheme] Axios error data:', error.response?.data)
      debugError('❌ [getTheme] Axios error message:', error.message)
    } else if (error instanceof Error) {
      debugError('❌ [getTheme] Error message:', error.message)
      debugError('❌ [getTheme] Error stack:', error.stack)
    } else {
      debugError('❌ [getTheme] Non-Error object:', error)
    }

    // Always throw a proper Error instance
    const apiError = handleApiError(error)
    debugError('❌ [getTheme] About to throw error with message:', apiError.message)
    throw apiError
  }
}

/**
 * Create a new theme
 */
export async function createTheme(
  request: CreateThemeRequest
): Promise<ThemeData> {
  try {
    const response = await apiClient.post<any>(
      '/admin/themes',
      request
    )
    // Handle different response formats - data, theme, or direct object
    const theme = response.data?.data || response.data?.theme || response.data
    if (!theme) {
      throw new Error('Failed to create theme')
    }
    return normalizeTheme(theme)
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Update an existing theme
 */
export async function updateTheme(
  themeId: string,
  request: UpdateThemeRequest
): Promise<ThemeData> {
  try {
    debugLog('🔄 [updateTheme] Sending to backend:', {
      themeId,
      componentCount: request.components?.length || 0,
      components: request.components?.map((c: any) => ({ id: c.id, type: c.type, order: c.order })) || []
    })
    const response = await apiClient.put<any>(
      `/admin/themes/${themeId}`,
      request
    )

    // Log response for debugging
    debugLog('📦 [updateTheme] Response status:', response.status)
    debugLog('📦 [updateTheme] Response.data:', response.data)

    // Handle different response formats - data, theme, or direct object
    const theme = response.data?.data || response.data?.theme || response.data

    if (!theme || (typeof theme === 'object' && Object.keys(theme).length === 0)) {
      throw new Error('Failed to update theme: empty response')
    }

    const normalizedTheme = normalizeTheme(theme)
    debugLog('✅ [updateTheme] Response from backend:', {
      componentCount: normalizedTheme.components?.length || 0,
      components: normalizedTheme.components?.map((c) => ({ id: c.id, type: c.type })) || []
    })
    return normalizedTheme
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Delete a theme
 */
export async function deleteTheme(themeId: string): Promise<void> {
  try {
    await apiClient.delete(`/admin/themes/${themeId}`)
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Activate a theme
 */
export async function activateTheme(themeId: string): Promise<void> {
  try {
    await apiClient.post(`/admin/themes/${themeId}/activate`)
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Duplicate a theme
 */
export async function duplicateTheme(
  themeId: string,
  newName?: string
): Promise<ThemeData> {
  try {
    const response = await apiClient.post<any>(
      `/admin/themes/${themeId}/duplicate`,
      { name: newName }
    )
    // Handle different response formats - data, theme, or direct object
    const theme = response.data?.data || response.data?.theme || response.data
    if (!theme) {
      throw new Error('Failed to duplicate theme')
    }
    return normalizeTheme(theme)
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Export theme as JSON
 */
export async function exportTheme(themeId: string): Promise<Blob> {
  try {
    const response = await apiClient.post(
      `/admin/themes/${themeId}/export`,
      {},
      { responseType: 'blob' }
    )
    return response.data
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Import theme from JSON
 */
export async function importTheme(themeJson: string): Promise<ThemeData> {
  try {
    const response = await apiClient.post<any>(
      '/admin/themes/import',
      { themeJson }
    )
    // Handle different response formats - data, theme, or direct object
    const theme = response.data?.data || response.data?.theme || response.data
    if (!theme) {
      throw new Error('Failed to import theme')
    }
    return normalizeTheme(theme)
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Get theme history
 */
export async function getThemeHistory(
  themeId: string
): Promise<Array<{ version: number; changeSummary: string; createdAt: string }>> {
  try {
    const response = await apiClient.get<ApiResponse<any[]>>(
      `/admin/themes/${themeId}/history`
    )
    return response.data.data || []
  } catch (error) {
    throw handleApiError(error)
  }
}

// Component API endpoints

/**
 * Get all components for a theme
 */
export async function getThemeComponents(themeId: string): Promise<any[]> {
  try {
    const response = await apiClient.get<ApiResponse<any[]>>(
      `/admin/themes/${themeId}/components`
    )
    return response.data.data || []
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Create a component in a theme
 */
export async function createComponent(
  themeId: string,
  component: any
): Promise<any> {
  try {
    const response = await apiClient.post<ApiResponse<any>>(
      `/admin/themes/${themeId}/components`,
      component
    )
    if (!response.data.data) {
      throw new Error('Failed to create component')
    }
    return response.data.data
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Update a component
 */
export async function updateComponent(
  themeId: string,
  componentId: string,
  updates: any
): Promise<any> {
  try {
    const response = await apiClient.put<ApiResponse<any>>(
      `/admin/themes/${themeId}/components/${componentId}`,
      updates
    )
    if (!response.data.data) {
      throw new Error('Failed to update component')
    }
    return response.data.data
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Delete a component
 */
export async function deleteComponent(
  themeId: string,
  componentId: string
): Promise<void> {
  try {
    await apiClient.delete(`/admin/themes/${themeId}/components/${componentId}`)
  } catch (error) {
    throw handleApiError(error)
  }
}

/**
 * Reorder components
 */
export async function reorderComponents(
  themeId: string,
  components: Array<{ id: string; displayOrder: number }>
): Promise<void> {
  try {
    await apiClient.post(`/admin/themes/${themeId}/components/reorder`, {
      components,
    })
  } catch (error) {
    throw handleApiError(error)
  }
}

// Library API endpoints

/**
 * Get component library
 */
export async function getComponentLibrary(): Promise<any[]> {
  try {
    const response = await apiClient.get<ApiResponse<any[]>>(
      '/public/component-library'
    )
    return response.data.data || []
  } catch (error) {
    throw handleApiError(error)
  }
}

// Error handling

function handleApiError(error: unknown): Error {
  // Handle Axios errors
  if (axios.isAxiosError(error)) {
    let message = 'An error occurred'

    // Try to extract message from response
    if (error.response?.data) {
      const data = error.response.data as any
      message = data.error || data.message || error.message || message
    } else if (error.message) {
      message = error.message
    }

    const apiError = new Error(message)
    ;(apiError as any).status = error.response?.status
    ;(apiError as any).originalError = error
    return apiError
  }

  // Handle Error instances
  if (error instanceof Error) {
    return error
  }

  // Handle string errors
  if (typeof error === 'string') {
    return new Error(error)
  }

  // Handle object errors
  if (error && typeof error === 'object') {
    const errorObj = error as any
    const message = errorObj.message || errorObj.error || JSON.stringify(error)
    return new Error(message || 'An unknown error occurred')
  }

  // Fallback for any other type
  return new Error('An unknown error occurred')
}

export default apiClient
