'use client'

import { apiCache } from '@/lib/apiCache'
import { permissionCache } from '@/lib/rbac'
import { useKDSStore } from '@/stores/kdsStore'
import { usePOSStore } from '@/stores/posStore'
import { usePreferencesStore } from '@/stores/preferencesStore'
import { useThemeBuilderStore } from '@/stores/themeBuilderStore'

const AUTH_SCOPED_STORAGE_KEYS = [
  'auth-storage',
  'token',
  'user',
  'pos-storage',
  'pos-products-cache',
  'pos-categories-cache',
  'preferences-storage',
  'theme-builder-storage',
  'rbac-permission-cache-v1',
  'kds-bootstrap-snapshot-v1',
] as const

const AUTH_SCOPED_DATABASES = ['pos-product-db', 'kds-offline-db'] as const

export function clearAuthScopedClientState() {
  apiCache.clearAll()
  permissionCache.clear()

  if (typeof window !== 'undefined') {
    for (const key of AUTH_SCOPED_STORAGE_KEYS) {
      window.localStorage.removeItem(key)
    }
  }

  usePOSStore.getState().resetForAuthScopeChange()
  useKDSStore.getState().resetForAuthScopeChange()
  usePreferencesStore.getState().reset()
  useThemeBuilderStore.getState().resetTheme()
}

export async function clearAuthScopedPersistentDatabases() {
  if (typeof window === 'undefined' || typeof window.indexedDB === 'undefined') {
    return
  }

  await Promise.all(AUTH_SCOPED_DATABASES.map(deleteDatabase))
}

function deleteDatabase(name: string) {
  return new Promise<void>((resolve, reject) => {
    const request = window.indexedDB.deleteDatabase(name)
    request.onsuccess = () => resolve()
    request.onerror = () => reject(request.error ?? new Error(`Could not clear ${name}`))
    request.onblocked = () => reject(new Error(`Could not clear ${name}: database is in use`))
  })
}
