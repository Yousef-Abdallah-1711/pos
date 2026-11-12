/**
 * Service Worker Management Hook
 *
 * Handles:
 * - Service worker registration
 * - Update detection
 * - Cache management
 * - Offline detection
 */

import { useEffect, useState, useCallback } from 'react'
import { debugLog, debugWarn, debugError } from '@/lib/debug'

interface ServiceWorkerState {
  isSupported: boolean
  isReady: boolean
  isUpdateAvailable: boolean
  isOffline: boolean
  registration?: ServiceWorkerRegistration
}

interface UseServiceWorkerOptions {
  onUpdateAvailable?: () => void
  onInstalled?: () => void
  onUpdating?: () => void
  enableAutoUpdate?: boolean
  checkInterval?: number
}

export function useServiceWorker(options: UseServiceWorkerOptions = {}) {
  const {
    onUpdateAvailable,
    onInstalled,
    onUpdating,
    enableAutoUpdate = true,
    checkInterval = 60000, // Check every minute
  } = options

  const [state, setState] = useState<ServiceWorkerState>({
    isSupported: false,
    isReady: false,
    isUpdateAvailable: false,
    isOffline: false,
  })

  // Register service worker
  useEffect(() => {
    if (typeof window === 'undefined') return

    const isSupported = 'serviceWorker' in navigator

    if (!isSupported) {
      debugWarn('[SW] Service Workers not supported')
      setState((prev) => ({ ...prev, isSupported: false }))
      return
    }

    setState((prev) => ({ ...prev, isSupported: true }))

    const registerSW = async () => {
      try {
        const registration = await navigator.serviceWorker.register('/sw.js', {
          scope: '/',
        })

        debugLog('[SW] Registered successfully:', registration)
        setState((prev) => ({
          ...prev,
          isReady: true,
          registration,
        }))

        onInstalled?.()

        // Check for updates
        registration.addEventListener('updatefound', () => {
          const newWorker = registration.installing

          if (newWorker) {
            newWorker.addEventListener('statechange', () => {
              if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
                debugLog('[SW] Update available')
                setState((prev) => ({ ...prev, isUpdateAvailable: true }))
                onUpdateAvailable?.()
              }
            })
          }
        })

        // Periodically check for updates
        if (enableAutoUpdate && checkInterval > 0) {
          const checkUpdates = () => {
            registration.update().catch(() => {})
          }

          const interval = setInterval(checkUpdates, checkInterval)
          return () => clearInterval(interval)
        }
      } catch (error) {
        debugError('[SW] Registration failed:', error)
      }
    }

    registerSW()
  }, [enableAutoUpdate, checkInterval, onInstalled, onUpdateAvailable])

  // Track online/offline status
  useEffect(() => {
    const handleOnline = () => {
      setState((prev) => ({ ...prev, isOffline: false }))
    }

    const handleOffline = () => {
      setState((prev) => ({ ...prev, isOffline: true }))
    }

    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)

    // Set initial offline state
    setState((prev) => ({ ...prev, isOffline: !navigator.onLine }))

    return () => {
      window.removeEventListener('online', handleOnline)
      window.removeEventListener('offline', handleOffline)
    }
  }, [])

  // Apply update
  const applyUpdate = useCallback(async () => {
    if (!state.registration?.waiting) return

    debugLog('[SW] Applying update...')
    onUpdating?.()

    // Tell waiting service worker to skip waiting and activate
    state.registration.waiting.postMessage({ type: 'SKIP_WAITING' })

    // Reload when the service worker is activated
    let reloaded = false
    navigator.serviceWorker.addEventListener('controllerchange', () => {
      if (!reloaded) {
        reloaded = true
        window.location.reload()
      }
    })
  }, [state.registration, onUpdating])

  // Clear all caches
  const clearCache = useCallback(async () => {
    if (!state.registration) return

    debugLog('[SW] Clearing cache...')
    state.registration.active?.postMessage({ type: 'CLEAR_CACHE' })

    // Also clear through Cache API directly
    const cacheNames = await caches.keys()
    await Promise.all(cacheNames.map((name) => caches.delete(name)))
  }, [state.registration])

  // Get cache size (approximate)
  const getCacheSize = useCallback(async () => {
    const cacheNames = await caches.keys()
    let totalSize = 0

    for (const cacheName of cacheNames) {
      const cache = await caches.open(cacheName)
      const keys = await cache.keys()

      for (const request of keys) {
        const response = await cache.match(request)
        if (response) {
          const blob = await response.blob()
          totalSize += blob.size
        }
      }
    }

    return totalSize
  }, [])

  return {
    ...state,
    applyUpdate,
    clearCache,
    getCacheSize,
  }
}

/**
 * Hook to detect offline state
 */
export function useOffline() {
  const [isOffline, setIsOffline] = useState(false)

  useEffect(() => {
    const handleOnline = () => setIsOffline(false)
    const handleOffline = () => setIsOffline(true)

    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)

    setIsOffline(!navigator.onLine)

    return () => {
      window.removeEventListener('online', handleOnline)
      window.removeEventListener('offline', handleOffline)
    }
  }, [])

  return isOffline
}

/**
 * Hook for periodic sync (background sync)
 */
export function useBackgroundSync(tag: string) {
  const [isSupported, setIsSupported] = useState(false)
  const [isRegistered, setIsRegistered] = useState(false)

  useEffect(() => {
    if (typeof window === 'undefined') return

    const supported = 'serviceWorker' in navigator && 'SyncManager' in window

    if (!supported) {
      debugWarn('[BG-SYNC] Not supported')
      return
    }

    setIsSupported(true)
  }, [])

  const register = useCallback(async () => {
    if (!isSupported) return

    try {
      const registration = await navigator.serviceWorker.ready
      await (registration as ServiceWorkerRegistration & { sync?: { register(tag: string): Promise<void> } }).sync?.register(tag)
      debugLog('[BG-SYNC] Registered:', tag)
      setIsRegistered(true)
    } catch (error) {
      debugError('[BG-SYNC] Registration failed:', error)
    }
  }, [isSupported, tag])

  return {
    isSupported,
    isRegistered,
    register,
  }
}

/**
 * Hook to request persistent storage
 */
export function usePersistentStorage() {
  const [isPersisted, setIsPersisted] = useState(false)
  const [isSupported, setIsSupported] = useState(false)

  useEffect(() => {
    if (typeof navigator === 'undefined' || !navigator.storage?.persist) {
      return
    }

    setIsSupported(true)

    // Check current persistence status
    navigator.storage.persisted().then(setIsPersisted)
  }, [])

  const requestPersistence = useCallback(async () => {
    if (!isSupported) return false

    try {
      const isPersisted = await navigator.storage.persist()
      setIsPersisted(isPersisted)
      debugLog('[STORAGE] Persistence:', isPersisted)
      return isPersisted
    } catch (error) {
      debugError('[STORAGE] Failed:', error)
      return false
    }
  }, [isSupported])

  return {
    isSupported,
    isPersisted,
    requestPersistence,
  }
}
