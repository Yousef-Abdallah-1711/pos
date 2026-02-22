'use client'

import { useEffect, useState, useRef } from 'react'
import { useRouter, usePathname } from 'next/navigation'
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
import { debugError } from '@/lib/debug'
import { TrialWarningBanner } from '@/components/layout/TrialWarningBanner'
import { useTenantStatus } from '@/hooks/useTenantStatus'
import { BlockedOverlay } from '@/components/blocked/BlockedOverlay'

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const router = useRouter()
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const logout = useAuthStore((state) => state.logout)
  const { settings, setSettings, setLoading } = usePreferencesStore()
  const { theme, setTheme } = useTheme()
  const [isProfileOpen, setIsProfileOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)
  const { isBlocked, blockReason } = useTenantStatus()

  const isLoggedIn = !!(user && token)

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
    if (!isLoggedIn) {
      router.push(`/${locale}/auth/login`)
    }
  }, [isLoggedIn, router, locale])

  // Load user preferences on mount
  useEffect(() => {
    const loadPreferences = async () => {
      if (!user?.id || !token) return

      // Skip if settings already loaded
      if (settings && settings.userId === user.id) return

      setLoading(true)
      try {
        const response = await fetch(`http://localhost:8080/api/v1/users/${user.id}/settings`, {
          headers: { Authorization: `Bearer ${token}` },
        })

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
  }, [user?.id, token])

  // Sync theme and colors globally when settings change
  useEffect(() => {
    if (settings && settings.theme) {
      // Apply theme from preferences store
      setTheme(settings.theme)
    }
  }, [settings?.theme, setTheme])

  // Sync theme back from next-themes when changed externally
  useEffect(() => {
    if (theme && settings && theme !== settings.theme && (theme === 'light' || theme === 'dark')) {
      // Update store if next-themes changes don't match settings
      // This keeps the UI in sync with the store
    }
  }, [theme, settings])

  const handleLogout = () => {
    logout()
    router.push(`/${locale}/auth/login`)
  }

  if (!isLoggedIn) {
    return null
  }

  const isRTL = locale === 'ar'
  const blockedBypassPaths = ['/dashboard/settings/billing']
  const shouldBypassBlockedOverlay = blockedBypassPaths.some((path) => pathname.includes(path))

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex">
      {/* Sidebar */}
      <Sidebar />

      {/* Main Content Area */}
      <div className={`flex-1 ${isRTL ? 'md:mr-64' : 'md:ml-64'}`}>
        {/* Header */}
        <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-30">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center h-16">
              {/* Mobile menu button - now positioned on the left side of header */}
              <div className="md:hidden">
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => {
                    // This would toggle the sidebar on mobile
                    // For now, we'll just add a placeholder
                  }}
                  className="bg-white dark:bg-gray-800"
                >
                  <Menu className="h-5 w-5" />
                </Button>
              </div>

              {/* Right Side */}
              <div className={`flex items-center gap-4 ${isRTL ? 'mr-auto' : 'ml-auto'}`}>
                {/* Notification Center */}
                <NotificationCenter />

                {/* Language Switcher */}
                <LanguageSwitcher />

                {/* Theme Toggle */}
                <button
                  onClick={() => {
                    const newTheme = theme === 'dark' ? 'light' : 'dark'
                    // Update next-themes library
                    setTheme(newTheme)
                    // Also update preferences store for persistence
                    if (settings) {
                      setSettings({ ...settings, theme: newTheme as 'light' | 'dark' | 'system' })
                    }
                  }}
                  className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition"
                  title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
                >
                  {theme === 'dark' ? (
                    <Sun className="h-5 w-5 text-gray-700 dark:text-gray-300" />
                  ) : (
                    <Moon className="h-5 w-5 text-gray-700 dark:text-gray-300" />
                  )}
                </button>

                {/* Profile Dropdown */}
                <div className="relative" ref={dropdownRef}>
                  <button
                    onClick={() => setIsProfileOpen(!isProfileOpen)}
                    className="flex items-center gap-3 px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition"
                  >
                    <div className="w-8 h-8 rounded-full bg-primary-500 flex items-center justify-center text-white font-medium">
                      {user?.name?.charAt(0).toUpperCase() || 'U'}
                    </div>
                    <div className="hidden sm:block text-left">
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
                      <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
                        <p className="text-sm font-medium text-gray-900 dark:text-white">{user?.name}</p>
                        <p className="text-xs text-gray-500 dark:text-gray-400">{user?.email}</p>
                      </div>

                      {/* Settings Links */}
                      <div className="py-2">
                        <Link
                          href={`/${locale}/dashboard/settings?tab=profile`}
                          onClick={() => setIsProfileOpen(false)}
                          className="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition"
                        >
                          <User className="h-4 w-4" />
                          {t('common.profile')}
                        </Link>
                        <Link
                          href={`/${locale}/dashboard/settings?tab=security`}
                          onClick={() => setIsProfileOpen(false)}
                          className="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition"
                        >
                          <Shield className="h-4 w-4" />
                          {t('common.security')}
                        </Link>
                        <Link
                          href={`/${locale}/dashboard/settings?tab=preferences`}
                          onClick={() => setIsProfileOpen(false)}
                          className="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition"
                        >
                          <Globe className="h-4 w-4" />
                          {t('common.preferences')}
                        </Link>
                        <Link
                          href={`/${locale}/dashboard/settings?tab=appearance`}
                          onClick={() => setIsProfileOpen(false)}
                          className="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition"
                        >
                          <Palette className="h-4 w-4" />
                          {t('common.appearance')}
                        </Link>
                        <Link
                          href={`/${locale}/dashboard/settings?tab=rbac`}
                          onClick={() => setIsProfileOpen(false)}
                          className="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition"
                        >
                          <Lock className="h-4 w-4" />
                          Access Control
                        </Link>
                      </div>

                      {/* Logout */}
                      <div className="border-t border-gray-200 dark:border-gray-700 pt-2">
                        <button
                          onClick={() => {
                            setIsProfileOpen(false)
                            handleLogout()
                          }}
                          className="flex items-center gap-3 px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 w-full transition"
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

        {/* Trial Warning Banner */}
        {!isBlocked && <TrialWarningBanner />}

        {/* Main Content */}
        <main className="w-full">
          {isBlocked && blockReason && !shouldBypassBlockedOverlay ? (
            <BlockedOverlay blockReason={blockReason} />
          ) : (
            children
          )}
        </main>
      </div>
    </div>
  )
}
