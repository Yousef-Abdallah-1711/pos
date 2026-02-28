'use client'
import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { debugLog } from '@/lib/debug'

export interface UserSettings {
  id: number
  userId: number
  language: 'en' | 'ar'
  theme: 'light' | 'dark' | 'system'
  primaryColor: string
  secondaryColor: string
  accentColor: string
  emailNotifications: boolean
  pushNotifications: boolean
  smsNotifications: boolean
}

interface PreferencesState {
  settings: UserSettings | null
  isLoading: boolean
  error: string | null

  // Actions
  setSettings: (settings: UserSettings) => void
  updateLanguage: (language: 'en' | 'ar') => void
  updateTheme: (theme: 'light' | 'dark' | 'system') => void
  updateColors: (primary: string, secondary: string, accent: string) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  reset: () => void
}

const defaultSettings: UserSettings = {
  id: 0,
  userId: 0,
  language: 'en',
  theme: 'system',
  primaryColor: '#3B82F6',
  secondaryColor: '#6366F1',
  accentColor: '#10B981',
  emailNotifications: true,
  pushNotifications: true,
  smsNotifications: false,
}

export const usePreferencesStore = create<PreferencesState>()(
  persist(
    (set, get) => ({
      settings: null,
      isLoading: false,
      error: null,

      setSettings: (settings) => {
        debugLog('[PreferencesStore] Setting preferences')
        set({ settings, error: null })
        // Apply colors immediately
        applyColors(settings.primaryColor, settings.secondaryColor, settings.accentColor)
        // Dispatch custom event for real-time sync across components
        window.dispatchEvent(new CustomEvent('preferencesUpdated', { detail: settings }))
      },

      updateLanguage: (language) => {
        const current = get().settings
        if (current) {
          set({ settings: { ...current, language } })
          // Update document direction for RTL support
          document.documentElement.dir = language === 'ar' ? 'rtl' : 'ltr'
          document.documentElement.lang = language
        }
      },

      updateTheme: (theme) => {
        const current = get().settings
        if (current) {
          set({ settings: { ...current, theme } })
        }
      },

      updateColors: (primaryColor, secondaryColor, accentColor) => {
        const current = get().settings
        if (current) {
          set({ settings: { ...current, primaryColor, secondaryColor, accentColor } })
          applyColors(primaryColor, secondaryColor, accentColor)
        }
      },

      setLoading: (isLoading) => set({ isLoading }),
      setError: (error) => set({ error }),
      reset: () => set({ settings: null, isLoading: false, error: null }),
    }),
    {
      name: 'preferences-storage',
      partialize: (state) => ({ settings: state.settings }),
    }
  )
)

// Helper function to generate color shades from a base color
function hexToHSL(hex: string): { h: number; s: number; l: number } {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  if (!result) return { h: 0, s: 0, l: 0 }

  let r = parseInt(result[1], 16) / 255
  let g = parseInt(result[2], 16) / 255
  let b = parseInt(result[3], 16) / 255

  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  let h = 0
  let s = 0
  const l = (max + min) / 2

  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)
    switch (max) {
      case r:
        h = ((g - b) / d + (g < b ? 6 : 0)) / 6
        break
      case g:
        h = ((b - r) / d + 2) / 6
        break
      case b:
        h = ((r - g) / d + 4) / 6
        break
    }
  }

  return { h: h * 360, s: s * 100, l: l * 100 }
}

function HSLToHex(h: number, s: number, l: number): string {
  s /= 100
  l /= 100
  const a = s * Math.min(l, 1 - l)
  const f = (n: number) => {
    const k = (n + h / 30) % 12
    const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1)
    return Math.round(255 * color).toString(16).padStart(2, '0')
  }
  return `#${f(0)}${f(8)}${f(4)}`
}

function generateColorShades(baseColor: string): Record<string, string> {
  const { h, s } = hexToHSL(baseColor)

  return {
    '50': HSLToHex(h, Math.max(s - 30, 10), 97),
    '100': HSLToHex(h, Math.max(s - 20, 15), 94),
    '200': HSLToHex(h, Math.max(s - 10, 20), 86),
    '300': HSLToHex(h, s, 76),
    '400': HSLToHex(h, s, 64),
    '500': HSLToHex(h, s, 50),
    '600': HSLToHex(h, s, 42),
    '700': HSLToHex(h, s, 34),
    '800': HSLToHex(h, s, 26),
    '900': HSLToHex(h, Math.min(s + 10, 100), 18),
  }
}

// Helper function to apply colors to CSS variables
function applyColors(primary: string, secondary: string, accent: string) {
  const root = document.documentElement

  const primaryShades = generateColorShades(primary)
  const secondaryShades = generateColorShades(secondary)
  const accentShades = generateColorShades(accent)

  // Apply primary color shades
  Object.entries(primaryShades).forEach(([shade, color]) => {
    root.style.setProperty(`--primary-${shade}`, color)
  })

  // Apply secondary color shades
  Object.entries(secondaryShades).forEach(([shade, color]) => {
    root.style.setProperty(`--secondary-${shade}`, color)
  })

  // Apply accent color shades
  Object.entries(accentShades).forEach(([shade, color]) => {
    root.style.setProperty(`--accent-${shade}`, color)
  })
}
