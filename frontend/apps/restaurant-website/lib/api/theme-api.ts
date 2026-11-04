import { apiClient } from './api-client'
import { debugError } from '../debug'
import type { ThemeData } from '@/lib/store/theme-store'

/**
 * Custom error for theme API failures
 */
export class ThemeApiError extends Error {
  constructor(
    public message: string,
    public slug: string,
    public originalError?: unknown,
    public statusCode?: number
  ) {
    super(message)
    this.name = 'ThemeApiError'
  }
}

/**
 * Get a theme by its slug
 * @param slug - Theme identifier (e.g., "warm-comfort")
 * @returns Promise<ThemeData>
 * @throws {ThemeApiError} if fetch fails
 */
const inFlightThemeRequests = new Map<string, Promise<ThemeData>>()

async function fetchThemeBySlug(slug: string): Promise<ThemeData> {
  try {
    if (!slug) {
      throw new ThemeApiError('Theme slug is required', slug)
    }

    const data = await apiClient.get<ThemeData>(
      `/public/themes/${slug}`,
      {
        timeout: 5000, // 5-second timeout
      }
    )

    // Validate response has required fields
    if (!data.name || !data.slug || !data.colors) {
      throw new ThemeApiError(
        'Invalid theme data returned from API',
        slug
      )
    }

    return data
  } catch (error) {
    // Re-throw ThemeApiError as-is
    if (error instanceof ThemeApiError) {
      throw error
    }

    // Convert axios errors to ThemeApiError
    const axiosError = error as any
    const statusCode = axiosError?.response?.status
    const message =
      axiosError?.response?.data?.error ||
      axiosError?.response?.data?.message ||
      axiosError?.message ||
      'Failed to load theme'

    throw new ThemeApiError(
      message,
      slug,
      error,
      statusCode
    )
  }
}

/** Coalesce concurrent theme reads so multiple components share one request. */
export function getThemeBySlug(slug: string): Promise<ThemeData> {
  if (!slug) {
    return Promise.reject(new ThemeApiError('Theme slug is required', slug))
  }

  const existingRequest = inFlightThemeRequests.get(slug)
  if (existingRequest) {
    return existingRequest
  }

  const request = fetchThemeBySlug(slug)
  inFlightThemeRequests.set(slug, request)
  void request.then(
    () => {
      if (inFlightThemeRequests.get(slug) === request) inFlightThemeRequests.delete(slug)
    },
    () => {
      if (inFlightThemeRequests.get(slug) === request) inFlightThemeRequests.delete(slug)
    }
  )
  return request
}

/**
 * Get all available themes (optional)
 * Can be used to populate a dropdown/selector
 * @returns Promise<ThemeData[]>
 */
export async function getAllThemes(): Promise<ThemeData[]> {
  try {
    const themes = await apiClient.get<ThemeData[]>(
      '/public/themes',
      {
        timeout: 10000,
      }
    )
    return themes || []
  } catch (error) {
    debugError('Failed to fetch all themes:', error)
    return []
  }
}

/**
 * Type exports
 */
export type { ThemeData }
