'use client'
import { create } from 'zustand'
import { persist } from 'zustand/middleware'

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

export interface AuthState {
  user: User | null
  token: string | null
  setAuth: (user: User, token: string) => void
  logout: () => void
  isAuthenticated: () => boolean
}

export const COOKIE_SESSION_MARKER = 'cookie-session'

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      setAuth: (user) => {
        set({ user, token: COOKIE_SESSION_MARKER })
      },
      logout: () => {
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
