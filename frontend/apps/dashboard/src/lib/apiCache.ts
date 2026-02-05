import { debugLog, debugError } from '@/lib/debug'

/**
 * Simple API Response Caching Utility
 * Caches API responses for a specified duration (default 30 seconds)
 * Helps reduce unnecessary API calls when navigating between pages
 */

interface CacheEntry {
  data: any
  timestamp: number
  ttl: number // Time to live in milliseconds
}

class APICache {
  private cache = new Map<string, CacheEntry>()

  /**
   * Get cached data if it exists and hasn't expired
   */
  get(key: string): any | null {
    const entry = this.cache.get(key)

    if (!entry) {
      return null
    }

    // Check if cache has expired
    const now = Date.now()
    const age = now - entry.timestamp

    if (age > entry.ttl) {
      this.cache.delete(key)
      return null
    }

    debugLog(`[APICache] Cache HIT for key: ${key} (age: ${age}ms)`)
    return entry.data
  }

  /**
   * Set cache data with TTL (default 30 seconds)
   */
  set(key: string, data: any, ttlMs: number = 30000): void {
    debugLog(`[APICache] Cache SET for key: ${key} with TTL: ${ttlMs}ms`)
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl: ttlMs,
    })
  }

  /**
   * Clear specific cache entry
   */
  clear(key: string): void {
    this.cache.delete(key)
    debugLog(`[APICache] Cache CLEAR for key: ${key}`)
  }

  /**
   * Clear all cache
   */
  clearAll(): void {
    this.cache.clear()
    debugLog('[APICache] Cache CLEAR ALL')
  }

  /**
   * Get cache stats
   */
  getStats() {
    return {
      size: this.cache.size,
      keys: Array.from(this.cache.keys()),
    }
  }
}

// Export singleton instance
export const apiCache = new APICache()

/**
 * Hook-friendly wrapper for caching API calls
 * Usage:
 * const data = await cachedFetch('/api/products', { params: {...} }, 60000)
 */
export async function cachedFetch(
  url: string,
  options: any = {},
  ttlMs: number = 30000
): Promise<any> {
  // Create a cache key from URL and options
  const cacheKey = `${url}:${JSON.stringify(options)}`

  // Check cache first
  const cached = apiCache.get(cacheKey)
  if (cached) {
    return cached
  }

  // If not cached, make the request
  try {
    const response = await fetch(url, options)
    const data = await response.json()

    // Cache the successful response
    apiCache.set(cacheKey, data, ttlMs)

    return data
  } catch (error) {
    debugError('[APICache] Fetch error:', error)
    throw error
  }
}

/**
 * Generate cache key from URL and params
 */
export function generateCacheKey(url: string, params?: Record<string, any>): string {
  if (!params) return url
  const queryString = new URLSearchParams(params).toString()
  return `${url}?${queryString}`
}
