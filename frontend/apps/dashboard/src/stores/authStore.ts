'use client'
import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export const COOKIE_SESSION_TOKEN = '__cookie_session__'

interface User {
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

interface AuthState {
  user: User | null
  token: string | null
  setAuth: (user: User, token?: string | null) => void
  logout: () => void
  isAuthenticated: () => boolean
  initializeFromCookie: () => void
}

const shouldPersistToken = process.env.NODE_ENV !== 'production'

const normalizeSessionToken = (user: User | null, token?: string | null) => {
  if (token) {
    return token
  }

  return user ? COOKIE_SESSION_TOKEN : null
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      setAuth: (user, token) => {
        // Persist the token in local development; production relies on cookie-backed auth.
        set({ user, token: normalizeSessionToken(user, token) })
      },
      logout: () => {
        set({ user: null, token: null })
      },
      isAuthenticated: () => {
        const state = get()
        return !!state.user
      },
      initializeFromCookie: () => {
        try {
          const cookieString = document.cookie
          const cookies = cookieString.split('; ').reduce((acc, cookie) => {
            const [key, value] = cookie.split('=')
            if (key && value) {
              acc[key] = decodeURIComponent(value)
            }
            return acc
          }, {} as Record<string, string>)

          const userStr = cookies['auth_user']
          if (userStr) {
            const user = JSON.parse(userStr)
            set({ user, token: COOKIE_SESSION_TOKEN })
          }
        } catch (error) {
          // Silent fail - cookie initialization is optional
        }
      },
    }),
    {
      name: 'auth-storage',
      merge: (persistedState, currentState) => {
        const persisted = (persistedState as Partial<AuthState>) ?? {}
        const merged = {
          ...currentState,
          ...persisted,
        }

        return {
          ...merged,
          token: normalizeSessionToken(merged.user ?? null, merged.token ?? null),
        }
      },
      partialize: (state) => ({
        user: state.user,
        ...(shouldPersistToken ? { token: state.token } : {}),
      }),
    }
  )
)
