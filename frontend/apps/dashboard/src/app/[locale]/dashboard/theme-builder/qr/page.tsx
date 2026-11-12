'use client'

import { useEffect, useMemo, useState } from 'react'
import { usePathname } from 'next/navigation'
import { Loader2 } from 'lucide-react'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'
import { debugError } from '@/lib/debug'
import * as themeApi from '@/lib/api/themeApi'
import { type ThemeData } from '@/types/theme'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { ThemeBuilderSubTabs } from '../components/ThemeBuilderSubTabs'
import { WebsiteQRStudio } from '../components/WebsiteQRStudio'
import { hasAuthenticatedSession } from '@/lib/authSession'

function normalizeDomain(value: string) {
  return value.trim().replace(/^https?:\/\//, '').replace(/\/+$/, '')
}

export default function ThemeBuilderQrPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const hasSession = hasAuthenticatedSession(user, token)
  const [themes, setThemes] = useState<ThemeData[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false

    const loadThemes = async () => {
      if (!hasSession) {
        setLoading(false)
        return
      }

      setLoading(true)
      try {
        const response = await themeApi.getThemes()
        if (!cancelled) {
          setThemes(response)
        }
      } catch (error) {
        debugError('[THEME_QR_BUILDER] Failed to load themes:', error)
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    void loadThemes()

    return () => {
      cancelled = true
    }
  }, [hasSession])

  const websiteHost = useMemo(() => {
    const activeTheme = themes.find((theme) => theme.isActive)
    const activeThemeDomain =
      activeTheme?.identity?.domain?.trim() ||
      activeTheme?.websiteIdentity?.domain?.trim()

    if (activeThemeDomain) {
      return activeThemeDomain
    }

    const restaurantDomain = normalizeDomain(getConfig().restaurantDomain)
    const restaurantID = user?.restaurant_id || user?.restaurantId

    if (restaurantID) {
      return `restaurant-${restaurantID}.${restaurantDomain}`
    }

    return `demo.${restaurantDomain}`
  }, [themes, user?.restaurant_id, user?.restaurantId])

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <ThemeBuilderSubTabs />

      {loading ? (
        <div className="flex h-64 items-center justify-center">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : (
        <WebsiteQRStudio
          locale={locale}
          websiteHost={websiteHost}
          t={t}
          isRTL={isRTL}
        />
      )}
    </div>
  )
}
