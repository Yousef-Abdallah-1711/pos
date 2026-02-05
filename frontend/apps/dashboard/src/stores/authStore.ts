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

interface AuthState {
  user: User | null
  token: string | null
  setAuth: (user: User, token: string) => void
  logout: () => void
  isAuthenticated: () => boolean
  initializeFromCookie: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      setAuth: (user, token) => {
        set({ user, token })
      },
      logout: () => {
        set({ user: null, token: null })
      },
      isAuthenticated: () => {
        const state = get()
        return !!(state.user && state.token)
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

          const token = cookies['auth_token']
          const userStr = cookies['auth_user']

          if (token && userStr) {
            const user = JSON.parse(userStr)
            set({ user, token })
          }
        } catch (error) {
          // Silent fail - cookie initialization is optional
        }
      },
    }),
    {
      name: 'auth-storage',
    }
  )
)
