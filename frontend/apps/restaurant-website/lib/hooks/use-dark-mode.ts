'use client'

import { useTheme as useNextTheme } from 'next-themes'
import { usePreferencesStore, type DarkModePreference } from '@/lib/store/preferences-store'

export function useDarkMode() {
  const { theme, resolvedTheme, setTheme } = useNextTheme()
  const storedMode = usePreferencesStore((state) => state.settings.darkMode)
  const setDarkModePreference = usePreferencesStore((state) => state.setDarkMode)

  const mode = (theme ?? storedMode ?? 'system') as DarkModePreference
  const isDark = resolvedTheme === 'dark'

  const setMode = (nextMode: DarkModePreference) => {
    setTheme(nextMode)
    setDarkModePreference(nextMode)
  }

  const toggle = () => {
    if (mode === 'light') {
      setMode('dark')
      return
    }
    if (mode === 'dark') {
      setMode('system')
      return
    }
    setMode('light')
  }

  return { mode, setMode, isDark, toggle }
}
