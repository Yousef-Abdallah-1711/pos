'use client'
import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { clearAuthScopedClientState, clearAuthScopedPersistentDatabases } from '@/lib/authScopeCleanup'
import { debugError } from '@/lib/debug'

export interface User {
  id: number
  email: string
  name: string
  role: string
  tenant_id: number
  restaurant_id?: number
  // Also support camelCase for compatibility
  tenantId?: number
  restaurantId?: number
}

export interface AuthState {
  user: User | null
  token: string | null
  setAuth: (user: User, token: string) => void
  logout: () => void
  isAuthenticated: () => boolean
}

export const COOKIE_SESSION_MARKER = 'cookie-session'

export function hasDifferentAuthScope(previous: User | null, next: User): boolean {
  if (!previous) return false

  const previousTenantID = previous.tenant_id ?? previous.tenantId ?? 0
  const nextTenantID = next.tenant_id ?? next.tenantId ?? 0
  const previousRestaurantID = previous.restaurant_id ?? previous.restaurantId ?? 0
  const nextRestaurantID = next.restaurant_id ?? next.restaurantId ?? 0

  return previous.id !== next.id
    || previousTenantID !== nextTenantID
    || previousRestaurantID !== nextRestaurantID
    || previous.role !== next.role
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      setAuth: (user) => {
        if (!get().user || hasDifferentAuthScope(get().user, user)) {
          clearAuthScopedClientState()
          void clearAuthScopedPersistentDatabases().catch((error: unknown) => {
            debugError('[Auth] Failed to clear scoped IndexedDB caches after identity change:', error)
          })
        }
        set({ user, token: COOKIE_SESSION_MARKER })
      },
      logout: () => {
        clearAuthScopedClientState()
        set({ user: null, token: null })
      },
      isAuthenticated: () => {
        const state = get()
        return !!(state.user && state.token)
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ user: state.user }),
      merge: (persisted, current) => {
        const legacy = persisted as Partial<AuthState>
        const user = legacy.user ?? null
        return { ...current, user, token: user ? COOKIE_SESSION_MARKER : null }
      },
    }
  )
)
