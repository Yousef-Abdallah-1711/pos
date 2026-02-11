import { debugLog, debugWarn, debugError } from './debug'

/**
 * Utility function to check tenant status from server-side
 * Used by middleware to determine if tenant can access their website
 */

const statusCache = new Map<
  string,
  {
    data: TenantStatusResult
    timestamp: number
  }
>();
const CACHE_TTL = 30 * 1000; // 30 seconds

export interface TenantStatusResult {
  canAccess: boolean
  blockReason: string | null
  isMaintenanceMode?: boolean
  tenantInfo?: {
    name: string
    contactEmail?: string
    contactPhone?: string
  }
}

/**
 * Clear status cache (for manual resets)
 */
export function clearStatusCache(slug?: string) {
  if (slug) {
    statusCache.delete(slug)
  } else {
    statusCache.clear()
  }
}

/**
 * Check tenant status with caching
 * Fetches from public API endpoint
 *
 * @param slug - Restaurant/tenant slug
 * @returns Status result with access flag and block reason
 */
export async function checkTenantStatus(slug: string): Promise<TenantStatusResult> {
  if (!slug) {
    debugWarn('[checkTenantStatus] No slug provided')
    return {
      canAccess: false,
      blockReason: 'invalid_tenant',
    }
  }

  // Check cache first
  const cached = statusCache.get(slug)
  if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
    debugLog(`[checkTenantStatus] Cache HIT for slug: ${slug}`)
    return cached.data
  }

  try {
    debugLog(`[checkTenantStatus] Checking status for slug: ${slug}`)

    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'
    const url = `${apiUrl}/public/tenant-status/${encodeURIComponent(slug)}`

    debugLog(`[checkTenantStatus] Fetching from: ${url}`)

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      // Cache for 30 seconds at the fetch level as well
      next: {
        revalidate: 30,
      },
    })

    if (!response.ok) {
      debugWarn(`[checkTenantStatus] API returned status: ${response.status}`)

      // If 404 or 400, tenant doesn't exist
      if (response.status === 404 || response.status === 400) {
        return {
          canAccess: false,
          blockReason: 'tenant_not_found',
        }
      }

      // For other errors, assume blocked (fail secure)
      return {
        canAccess: false,
        blockReason: 'api_error',
      }
    }

    const responseData = await response.json()
    debugLog(`[checkTenantStatus] API Response:`, responseData)

    // Extract status from API response
    const result: TenantStatusResult = {
      canAccess: responseData.data?.can_access ?? false,
      blockReason: responseData.data?.block_reason ?? null,
      isMaintenanceMode: responseData.data?.is_maintenance ?? false,
    }

    // Cache the result
    statusCache.set(slug, {
      data: result,
      timestamp: Date.now(),
    })

    debugLog(
      `[checkTenantStatus] Status for ${slug}: canAccess=${result.canAccess}, reason=${result.blockReason}`
    )

    return result
  } catch (error) {
    debugError(`[checkTenantStatus] Error checking status for ${slug}:`, error)

    // On error, deny access (fail secure)
    return {
      canAccess: false,
      blockReason: 'error',
    }
  }
}

/**
 * Check if a path should skip tenant status check
 * Allows maintenance page and other excluded paths to load even when blocked
 */
export function shouldSkipStatusCheck(pathname: string): boolean {
  // Exact paths to skip
  const skipPaths = ['/maintenance', '/favicon.ico', '/robots.txt', '/sitemap.xml']

  for (const path of skipPaths) {
    if (pathname === path || pathname.startsWith(path)) {
      return true
    }
  }

  // Skip path prefixes
  const skipPrefixes = ['/_next/', '/api/', '/public/', '/health']

  for (const prefix of skipPrefixes) {
    if (pathname.startsWith(prefix)) {
      return true
    }
  }

  return false
}
