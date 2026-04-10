'use client'

import { Suspense, useEffect, useState, useRef } from 'react'
import { usePathname } from 'next/navigation'
import Link from 'next/link'
import { useAuthStore } from '@/stores/authStore'
import { usePreferencesStore } from '@/stores/preferencesStore'
import { LogOut, Moon, Sun, User, Menu, Shield, Globe, Palette, ChevronDown, Lock } from 'lucide-react'
import { useTheme } from 'next-themes'
import { Button } from '@/components/ui/button'
import { Sidebar } from '@/components/layout/Sidebar'
import { LanguageSwitcher } from '@/components/layout/LanguageSwitcher'
import { NotificationCenter } from '@/components/notifications/NotificationCenter'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { getDashboardLoginUrl } from '@/lib/authSession'
import { debugError } from '@/lib/debug'
import { getConfig } from '@/lib/config'
import { commonRoutes, preloadRoutes } from '@/lib/code-splitting'
import { TrialWarningBanner } from '@/components/layout/TrialWarningBanner'
import { useTenantStatus } from '@/hooks/useTenantStatus'
import { BlockedOverlay } from '@/components/blocked/BlockedOverlay'
import { authFetch } from '@/lib/authFetch'

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const user = useAuthStore((state) => state.user)
  const logout = useAuthStore((state) => state.logout)
  const { settings, setSettings, setLoading, updateTheme } = usePreferencesStore()
  const { resolvedTheme, setTheme } = useTheme()
  const [isProfileOpen, setIsProfileOpen] = useState(false)
  const [isSidebarMobileOpen, setIsSidebarMobileOpen] = useState(false)
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)
  const { isBlocked, blockReason } = useTenantStatus()

  const isLoggedIn = !!user

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsProfileOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  useEffect(() => {
    const saved = window.localStorage.getItem('dashboard.sidebar.collapsed')
    if (saved === '1') {
      setIsSidebarCollapsed(true)
    }
  }, [])

  useEffect(() => {
    window.localStorage.setItem('dashboard.sidebar.collapsed', isSidebarCollapsed ? '1' : '0')
  }, [isSidebarCollapsed])

  useEffect(() => {
    preloadRoutes(commonRoutes, { staggerDelay: 600, maxConcurrent: 2 })
  }, [])

  useEffect(() => {
    if (!isLoggedIn) {
      window.location.href = getDashboardLoginUrl(pathname)
    }
  }, [isLoggedIn, pathname])

  // Load user preferences on mount
  useEffect(() => {
    const loadPreferences = async () => {
      if (!user?.id) return

      // Skip if settings already loaded
      if (settings && settings.userId === user.id) return

      setLoading(true)
      try {
        const { apiUrl } = getConfig()
        const response = await authFetch(`${apiUrl}/users/${user.id}/settings`)

        if (response.ok) {
          const data = await response.json()
          setSettings({
            id: data.id,
            userId: data.user_id,
            language: data.language || 'en',
            theme: data.theme || 'system',
            primaryColor: data.primary_color || '#3B82F6',
            secondaryColor: data.secondary_color || '#6366F1',
            accentColor: data.accent_color || '#10B981',
            emailNotifications: data.email_notifications ?? true,
            pushNotifications: data.push_notifications ?? true,
            smsNotifications: data.sms_notifications ?? false,
          })
        }
      } catch (error) {
        debugError('[DashboardLayout] Failed to load preferences:', error)
      } finally {
        setLoading(false)
      }
    }

    loadPreferences()
  }, [user?.id, settings, setLoading, setSettings])

  // Sync theme and colors globally when settings change
  useEffect(() => {
    const themePreference = settings?.theme
    if (themePreference) {
      // Apply theme from preferences store
      setTheme(themePreference)
    }
  }, [settings, setTheme])

  const handleLogout = () => {
    logout()
    window.location.href = getDashboardLoginUrl(pathname)
  }

  if (!isLoggedIn) {
    return null
  }

  const isRTL = locale === 'ar'
  const effectiveTheme = resolvedTheme === 'dark' ? 'dark' : 'light'
  const blockedBypassPaths = ['/dashboard/settings/billing']
  const shouldBypassBlockedOverlay = blockedBypassPaths.some((path) => pathname.includes(path))
  const isThemeBuilderEditorRoute = pathname.includes('/dashboard/theme-builder/editor/')
  const isThemeBuilderListRoute = pathname.includes('/dashboard/theme-builder')
  const isKDSFullscreenRoute = pathname.includes('/dashboard/kds') && !pathname.includes('/dashboard/kds/settings')
  const isPOSFullscreenRoute = pathname.includes('/dashboard/pos') && !pathname.includes('/dashboard/pos/settings')
  const isFullscreenRoute = isKDSFullscreenRoute || isPOSFullscreenRoute

  const contentWrapperClass = isFullscreenRoute
    ? 'w-full h-full'
    : isThemeBuilderEditorRoute
      ? 'w-full'
      : isThemeBuilderListRoute
        ? 'mx-auto w-full max-w-[1800px] px-4 py-6 sm:px-6 lg:px-8'
        : 'mx-auto w-full max-w-7xl px-4 py-6 sm:px-6 lg:px-8'

  return (
    <div className={`min-h-screen flex ${isFullscreenRoute ? 'bg-gray-50 text-gray-900 dark:bg-slate-950 dark:text-slate-100' : 'bg-gray-50 dark:bg-gray-900'}`}>
      {/* Sidebar */}
      {!isFullscreenRoute && (
        <Sidebar
          isMobileOpen={isSidebarMobileOpen}
          onMobileToggle={() => setIsSidebarMobileOpen((prev) => !prev)}
          isCollapsed={isSidebarCollapsed}
          onCollapsedChange={setIsSidebarCollapsed}
        />
      )}

      {/* Main Content Area */}
      <div className={`flex-1 ${isFullscreenRoute ? '' : isRTL ? (isSidebarCollapsed ? 'md:mr-20' : 'md:mr-64') : (isSidebarCollapsed ? 'md:ml-20' : 'md:ml-64')} transition-all duration-300`}>
        {/* Header */}
        {!isFullscreenRoute && (
          <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-30">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
              <div className={`flex items-center justify-between h-16 gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
                {/* Mobile menu button */}
                <div className="md:hidden flex-shrink-0">
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => setIsSidebarMobileOpen((prev) => !prev)}
                    className="bg-white dark:bg-gray-800"
                  >
                    <Menu className="h-5 w-5" />
                  </Button>
                </div>

                <div className={`flex items-center gap-2 sm:gap-3 ${isRTL ? 'mr-auto flex-row-reverse' : 'ml-auto'}`}>
                  {/* Notification Center */}
                  <NotificationCenter />

                  {/* Language Switcher */}
                  <LanguageSwitcher />

                  {/* Theme Toggle */}
                  <button
                    onClick={() => {
                      const newTheme = effectiveTheme === 'dark' ? 'light' : 'dark'
                      // Update next-themes library
                      setTheme(newTheme)
                      // Also update preferences store for persistence
                      updateTheme(newTheme)
                    }}
                    className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition"
                    title={effectiveTheme === 'dark' ? t('common.lightMode', 'Switch to light mode') : t('common.darkMode', 'Switch to dark mode')}
                  >
                    {effectiveTheme === 'dark' ? (
                      <Sun className="h-5 w-5 text-gray-700 dark:text-gray-300" />
                    ) : (
                      <Moon className="h-5 w-5 text-gray-700 dark:text-gray-300" />
                    )}
                  </button>

                  {/* Profile Dropdown */}
                  <div className="relative" ref={dropdownRef}>
                    <button
                      onClick={() => setIsProfileOpen(!isProfileOpen)}
                      className={`flex items-center gap-2 sm:gap-3 px-2 sm:px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition ${isRTL ? 'flex-row-reverse' : ''}`}
                    >
                      <div className="w-8 h-8 rounded-full bg-primary-500 flex items-center justify-center text-white font-medium">
                        {user?.name?.charAt(0).toUpperCase() || 'U'}
                      </div>
                      <div className={`hidden lg:block ${isRTL ? 'text-right' : 'text-left'}`}>
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                          {user?.name}
                        </p>
                        <p className="text-xs text-gray-500 dark:text-gray-400 capitalize">
                          {user?.role}
                        </p>
                      </div>
                      <ChevronDown className={`h-4 w-4 text-gray-500 dark:text-gray-400 transition-transform ${isProfileOpen ? 'rotate-180' : ''}`} />
                    </button>

                    {/* Dropdown Menu */}
                    {isProfileOpen && (
                      <div className={`absolute ${isRTL ? 'left-0' : 'right-0'} mt-2 w-64 bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 py-2 z-50`}>
                        {/* User Info Header */}
                        <div className={`px-4 py-3 border-b border-gray-200 dark:border-gray-700 ${isRTL ? 'text-right' : 'text-left'}`}>
                          <p className="text-sm font-medium text-gray-900 dark:text-white">{user?.name}</p>
                          <p className="text-xs text-gray-500 dark:text-gray-400">{user?.email}</p>
                        </div>

                        {/* Settings Links */}
                        <div className="py-2">
                          <Link
                            href={`/${locale}/dashboard/settings?tab=profile`}
                            onClick={() => setIsProfileOpen(false)}
                            className={`flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition ${isRTL ? 'flex-row-reverse text-right' : ''}`}
                          >
                            <User className="h-4 w-4" />
                            {t('common.profile')}
                          </Link>
                          <Link
                            href={`/${locale}/dashboard/settings?tab=security`}
                            onClick={() => setIsProfileOpen(false)}
                            className={`flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition ${isRTL ? 'flex-row-reverse text-right' : ''}`}
                          >
                            <Shield className="h-4 w-4" />
                            {t('common.security')}
                          </Link>
                          <Link
                            href={`/${locale}/dashboard/settings?tab=preferences`}
                            onClick={() => setIsProfileOpen(false)}
                            className={`flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition ${isRTL ? 'flex-row-reverse text-right' : ''}`}
                          >
                            <Globe className="h-4 w-4" />
                            {t('common.preferences')}
                          </Link>
                          <Link
                            href={`/${locale}/dashboard/settings?tab=appearance`}
                            onClick={() => setIsProfileOpen(false)}
                            className={`flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition ${isRTL ? 'flex-row-reverse text-right' : ''}`}
                          >
                            <Palette className="h-4 w-4" />
                            {t('common.appearance')}
                          </Link>
                          <Link
                            href={`/${locale}/dashboard/settings?tab=rbac`}
                            onClick={() => setIsProfileOpen(false)}
                            className={`flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition ${isRTL ? 'flex-row-reverse text-right' : ''}`}
                          >
                            <Lock className="h-4 w-4" />
                            {t('settings.accessControl', 'Access Control')}
                          </Link>
                        </div>

                        {/* Logout */}
                        <div className="border-t border-gray-200 dark:border-gray-700 pt-2">
                          <button
                            onClick={() => {
                              setIsProfileOpen(false)
                              handleLogout()
                            }}
                            className={`flex items-center gap-3 px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 w-full transition ${isRTL ? 'flex-row-reverse text-right' : ''}`}
                          >
                            <LogOut className="h-4 w-4" />
                            {t('common.logout')}
                          </button>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </header>
        )}

        {/* Trial Warning Banner */}
        {!isBlocked && !isFullscreenRoute && <TrialWarningBanner />}

        {/* Main Content */}
        <main className={isFullscreenRoute ? 'w-full min-h-screen' : 'w-full min-h-[calc(100vh-4rem)]'}>
          {isBlocked && blockReason && !shouldBypassBlockedOverlay ? (
            <BlockedOverlay blockReason={blockReason} />
          ) : (
            <Suspense fallback={<div className="h-[60vh] w-full animate-pulse rounded-xl bg-gray-200 dark:bg-gray-800" />}>
              <div className={contentWrapperClass}>
                {children}
              </div>
            </Suspense>
          )}
        </main>
      </div>
    </div>
  )
}
