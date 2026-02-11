/**
 * Theme API Service
 * Handles all API calls for theme management
 *
 * ⚠️ IMPORTANT: This service reads JWT token from Zustand auth-storage, not localStorage directly
 * The token is stored in Zustand's persist middleware under the key 'auth-storage'
 */

import axios, { AxiosError } from 'axios'
import { ThemeData, CreateThemeRequest, UpdateThemeRequest } from '@/types/theme'
import { debugLog, debugWarn, debugError } from '@/lib/debug'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
const API_VERSION = '/api/v1'

// Create axios instance with default config
const apiClient = axios.create({
  baseURL: `${API_BASE_URL}${API_VERSION}`,
  headers: {
    'Content-Type': 'application/json',
  },
})

// ==============================================
// 🔐 AUTH INTERCEPTOR: Get token from Zustand auth-storage
// This reads the token from localStorage where Zustand persists it
// Key: 'auth-storage' (contains entire auth state as JSON)
// ==============================================
apiClient.interceptors.request.use((config) => {
  try {
    // Get the full auth state from localStorage
    const authStorageRaw = localStorage.getItem('auth-storage')
    if (!authStorageRaw) {
      debugWarn('⚠️ [API Request] No auth-storage found in localStorage')
      return config
    }

    // Parse the JSON to get the token
    const authState = JSON.parse(authStorageRaw)
    const token = authState?.state?.token

    if (token) {
      debugLog('✅ [API Request] Token found, adding to Authorization header')
      config.headers.Authorization = `Bearer ${token}`
    } else {
      debugWarn('⚠️ [API Request] Token not found in auth-storage.state.token')
      debugLog('📦 [API Request] Auth storage contents:', authState)
    }
  } catch (err) {
    debugError('❌ [API Request] Error parsing auth-storage:', err)
  }

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
      `/admin/themes?limit=${limit}&page=1`
    )
    
    // Log full response for debugging
    debugLog('📦 [getThemes] Full API Response:', response.data)
    
    // Handle both response formats
    const data = response.data?.data || response.data?.themes || []
    const pagination = response.data?.pagination || {}
    
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
            `/admin/themes?limit=${limit}&page=${page}`
          )
          const pageData = pageResponse.data?.data || pageResponse.data?.themes || []
          if (Array.isArray(pageData)) {
            allThemes.push(...pageData)
          }
        }
        
        debugLog(`✅ [getThemes] Fetched all ${allThemes.length} themes across ${pagination.totalPages} pages`)
        return allThemes
      }
      
      return data
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
    if (!theme.id) {
      const errorMsg = `Theme missing id. Theme: ${JSON.stringify(theme)}`
      debugError('❌ [getTheme]', errorMsg)
      throw new Error(errorMsg)
    }

    // IMPORTANT: Parse config JSON if it's a string
    // When created from a preset, config contains the full theme JSON
    if (theme.config && typeof theme.config === 'string') {
      try {
        debugLog('🔄 [getTheme] Parsing config JSON from preset...')
        const parsedConfig = JSON.parse(theme.config)
        debugLog('✅ [getTheme] Parsed config:', {
          hasComponents: !!parsedConfig.components,
          componentCount: Array.isArray(parsedConfig.components) ? parsedConfig.components.length : 0
        })
        // Merge the parsed config into the theme so components, colors, etc. are available at top level
        theme = {
          ...theme,
          ...parsedConfig,
          // Keep original config as the raw JSON for saving
          _rawConfig: theme.config,
        }
      } catch (parseErr) {
        debugWarn('⚠️ [getTheme] Failed to parse config JSON:', parseErr)
        // Continue with theme as-is if JSON parsing fails
      }
    }

    debugLog('✅ [getTheme] Received theme:', {
      id: theme.id,
      name: theme.name,
      components: Array.isArray(theme.components) ? theme.components.length : 0,
      hasColors: !!theme.colors,
      hasTypography: !!theme.typography,
    })
    return theme
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
    return theme
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
      components: request.components?.map((c) => ({ id: c.id, type: c.type, order: c.order })) || []
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

    debugLog('✅ [updateTheme] Response from backend:', {
      componentCount: theme.components?.length || 0,
      components: theme.components?.map((c) => ({ id: c.id, type: c.type })) || []
    })
    return theme
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
    let theme = response.data?.data || response.data?.theme || response.data
    if (!theme) {
      throw new Error('Failed to duplicate theme')
    }

    // IMPORTANT: Parse config JSON if it's a string
    // When duplicating from a preset, config contains the full theme JSON
    if (theme.config && typeof theme.config === 'string') {
      try {
        debugLog('🔄 [duplicateTheme] Parsing config JSON from preset...')
        const parsedConfig = JSON.parse(theme.config)
        debugLog('✅ [duplicateTheme] Parsed config:', {
          hasComponents: !!parsedConfig.components,
          componentCount: Array.isArray(parsedConfig.components) ? parsedConfig.components.length : 0
        })
        // Merge the parsed config into the theme so components, colors, etc. are available at top level
        theme = {
          ...theme,
          ...parsedConfig,
          // Keep original config as the raw JSON for saving
          _rawConfig: theme.config,
        }
      } catch (parseErr) {
        debugWarn('⚠️ [duplicateTheme] Failed to parse config JSON:', parseErr)
        // Continue with theme as-is if JSON parsing fails
      }
    }

    return theme
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
    return theme
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
