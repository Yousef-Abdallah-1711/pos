'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useHasPermission } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator } from '@/lib/translations'
import {
  User,
  Lock,
  Globe,
  Palette,
  Shield,
  Users,
  Key,
  Settings,
  Store,
  CreditCard
} from 'lucide-react'

interface SettingsSidebarProps {
  locale: string
}

export function SettingsSidebar({ locale }: SettingsSidebarProps) {
  const pathname = usePathname()
  const t = createTranslator(locale === 'ar' ? 'ar' : 'en')
  const isRTL = locale === 'ar'

  // Permission checks for different settings sections
  const { hasPermission: canManageSettings } = useHasPermission(MODULE_IDS.SETTINGS, PERMISSION_LEVELS.WRITE)
  const { hasPermission: canManageRoles } = useHasPermission(MODULE_IDS.ROLES, PERMISSION_LEVELS.WRITE)
  const { hasPermission: canManageUsers } = useHasPermission(MODULE_IDS.USERS, PERMISSION_LEVELS.WRITE)
  const { hasPermission: canManageThemes } = useHasPermission(MODULE_IDS.THEMES, PERMISSION_LEVELS.WRITE)

  // Admin has either SETTINGS WRITE or ROLES WRITE or USERS WRITE
  const isAdmin = canManageRoles || canManageUsers || canManageSettings

  const isActive = (path: string) => pathname.includes(path)

  // Settings links always visible (users can view their own settings)
  const settingsLinks = [
    {
      href: `/${locale}/dashboard/settings/account`,
      label: t('settings.sidebar.accountTitle', 'Account Settings'),
      icon: User,
      description: t('settings.sidebar.accountDesc', 'Profile, security, preferences'),
      requiresPermission: null, // Always visible
      hidden: false
    },
    {
      href: `/${locale}/dashboard/settings/restaurant`,
      label: t('settings.sidebar.restaurantTitle', 'Restaurant'),
      icon: Store,
      description: t('settings.sidebar.restaurantDesc', 'Contact info, maintenance page'),
      requiresPermission: canManageSettings, // Requires SETTINGS WRITE to edit
      hidden: false
    },
    {
      href: `/${locale}/dashboard/settings/appearance`,
      label: t('settings.sidebar.appearanceTitle', 'Appearance'),
      icon: Palette,
      description: t('settings.sidebar.appearanceDesc', 'Theme and colors'),
      requiresPermission: canManageThemes, // Requires THEMES WRITE to edit
      hidden: false
    },
    {
      href: `/${locale}/dashboard/settings/access-control`,
      label: t('settings.sidebar.accessControlTitle', 'Access Control'),
      icon: Shield,
      description: t('settings.sidebar.accessControlDesc', 'View your roles and permissions'),
      requiresPermission: null, // Always visible - user can view own permissions
      hidden: false
    },
    {
      href: `/${locale}/dashboard/settings/billing`,
      label: t('settings.sidebar.billingTitle', 'Billing & Subscription'),
      icon: CreditCard,
      description: t('settings.sidebar.billingDesc', 'Manage plans and payments'),
      requiresPermission: null, // Always visible - user can view own billing
      hidden: false
    },
  ]

  // Admin links - only visible to admins
  const adminLinks = isAdmin ? [
    {
      href: `/${locale}/dashboard/settings/roles`,
      label: t('settings.sidebar.roleManagementTitle', 'Role Management'),
      icon: Key,
      description: t('settings.sidebar.roleManagementDesc', 'Create and manage roles'),
      requiresPermission: canManageRoles, // Requires ROLES WRITE
      hidden: !canManageRoles,
      adminOnly: true
    },
    {
      href: `/${locale}/dashboard/settings/users`,
      label: t('settings.sidebar.userAssignmentsTitle', 'User Assignments'),
      icon: Users,
      description: t('settings.sidebar.userAssignmentsDesc', 'Assign roles to users'),
      requiresPermission: canManageUsers, // Requires USERS WRITE
      hidden: !canManageUsers,
      adminOnly: true
    },
  ] : []

  const allLinks = [...settingsLinks, ...adminLinks]

  return (
    <aside className={`w-full md:w-64 ${isRTL ? 'md:pl-6' : 'md:pr-6'}`}>
      <nav className="space-y-2">
        {/* Regular Settings */}
        <div className="mb-6">
          <h3 className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider px-3 py-2 mb-2">
            <Settings className={`inline h-3 w-3 ${isRTL ? 'ml-1' : 'mr-1'}`} />
            {t('settings.title', 'Settings')}
          </h3>
          <div className="space-y-1">
            {settingsLinks.map((link) => {
              const Icon = link.icon
              const active = isActive(link.href)
              const isDisabled = link.requiresPermission === false && link.requiresPermission !== null

              return (
                <Link key={link.href} href={isDisabled ? '#' : link.href}>
                  <div className={`flex items-start gap-3 px-3 py-2 rounded-lg transition-colors ${
                    isDisabled
                      ? 'opacity-50 cursor-not-allowed'
                      : active
                        ? 'bg-primary-100 dark:bg-primary-900/30'
                        : 'hover:bg-gray-100 dark:hover:bg-gray-800'
                  } ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
                    <Icon className={`h-4 w-4 mt-0.5 flex-shrink-0 ${
                      isDisabled
                        ? 'text-gray-400 dark:text-gray-500'
                        : active
                          ? 'text-primary-700 dark:text-primary-400'
                          : 'text-gray-500 dark:text-gray-400'
                    }`} />
                    <div className="flex-1 min-w-0">
                      <p className={`text-sm font-medium ${
                        isDisabled
                          ? 'text-gray-400 dark:text-gray-500'
                          : active
                            ? 'text-primary-700 dark:text-primary-300'
                            : 'text-gray-700 dark:text-gray-300'
                      }`}>
                        {link.label}
                      </p>
                      <p className="text-xs text-gray-500 dark:text-gray-400">
                        {link.description}
                      </p>
                    </div>
                  </div>
                </Link>
              )
            })}
          </div>
        </div>

        {/* Admin Section */}
        {isAdmin && adminLinks.filter(link => !link.hidden).length > 0 && (
          <div className="mb-6 pt-4 border-t border-gray-200 dark:border-gray-800">
            <h3 className={`text-xs font-semibold text-primary-600 dark:text-primary-400 uppercase tracking-wider px-3 py-2 mb-2 flex items-center gap-1 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
              <Lock className="h-3 w-3" />
              {t('settings.sidebar.administrationTitle', 'Administration')}
            </h3>
            <div className="space-y-1">
              {adminLinks
                .filter(link => !link.hidden)
                .map((link) => {
                  const Icon = link.icon
                  const active = isActive(link.href)
                  const isDisabled = link.requiresPermission === false && link.requiresPermission !== null

                  return (
                    <Link key={link.href} href={isDisabled ? '#' : link.href}>
                      <div className={`flex items-start gap-3 px-3 py-2 rounded-lg transition-colors ${
                        isDisabled
                          ? 'opacity-50 cursor-not-allowed'
                          : active
                            ? 'bg-primary-100 dark:bg-primary-900/30'
                            : 'hover:bg-gray-100 dark:hover:bg-gray-800'
                      } ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
                        <Icon className={`h-4 w-4 mt-0.5 flex-shrink-0 ${
                          isDisabled
                            ? 'text-gray-400 dark:text-gray-500'
                            : active
                              ? 'text-primary-700 dark:text-primary-400'
                              : 'text-gray-500 dark:text-gray-400'
                        }`} />
                        <div className="flex-1 min-w-0">
                          <p className={`text-sm font-medium ${
                            isDisabled
                              ? 'text-gray-400 dark:text-gray-500'
                              : active
                                ? 'text-primary-700 dark:text-primary-300'
                                : 'text-gray-700 dark:text-gray-300'
                          }`}>
                            {link.label}
                          </p>
                          <p className="text-xs text-gray-500 dark:text-gray-400">
                            {link.description}
                          </p>
                        </div>
                      </div>
                    </Link>
                  )
                })}
            </div>
          </div>
        )}

        {/* Help */}
        <div className="pt-4 border-t border-gray-200 dark:border-gray-800 mt-8">
          <button className={`w-full flex items-start gap-3 px-3 py-2 rounded-lg text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors ${isRTL ? 'text-right flex-row-reverse' : 'text-left'}`}>
            <Globe className="h-4 w-4 mt-0.5 flex-shrink-0" />
            <div>
              <p className="text-sm font-medium">{t('settings.sidebar.helpSupportTitle', 'Help & Support')}</p>
              <p className="text-xs text-gray-500 dark:text-gray-400">{t('settings.sidebar.helpSupportDesc', 'Contact administrator')}</p>
            </div>
          </button>
        </div>
      </nav>
    </aside>
  )
}
