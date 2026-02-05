'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useHasPermission } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import {
  User,
  Lock,
  Globe,
  Palette,
  Shield,
  Users,
  Key,
  Settings
} from 'lucide-react'

interface SettingsSidebarProps {
  locale: string
}

export function SettingsSidebar({ locale }: SettingsSidebarProps) {
  const pathname = usePathname()

  // Check if user is admin (can manage roles)
  const { hasPermission: isAdmin } = useHasPermission(MODULE_IDS.SETTINGS, PERMISSION_LEVELS.WRITE)

  const isActive = (path: string) => pathname.includes(path)

  const settingsLinks = [
    {
      href: `/${locale}/dashboard/settings/account`,
      label: 'Account Settings',
      icon: User,
      description: 'Profile, security, preferences'
    },
    {
      href: `/${locale}/dashboard/settings/appearance`,
      label: 'Appearance',
      icon: Palette,
      description: 'Theme and colors'
    },
    {
      href: `/${locale}/dashboard/settings/access-control`,
      label: 'Access Control',
      icon: Shield,
      description: 'View your roles and permissions'
    },
  ]

  const adminLinks = isAdmin ? [
    {
      href: `/${locale}/dashboard/settings/roles`,
      label: 'Role Management',
      icon: Key,
      description: 'Create and manage roles',
      adminOnly: true
    },
    {
      href: `/${locale}/dashboard/settings/users`,
      label: 'User Assignments',
      icon: Users,
      description: 'Assign roles to users',
      adminOnly: true
    },
  ] : []

  const allLinks = [...settingsLinks, ...adminLinks]

  return (
    <aside className="w-full md:w-64 md:pr-6">
      <nav className="space-y-2">
        {/* Regular Settings */}
        <div className="mb-6">
          <h3 className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider px-3 py-2 mb-2">
            <Settings className="inline h-3 w-3 mr-1" />
            Settings
          </h3>
          <div className="space-y-1">
            {settingsLinks.map((link) => {
              const Icon = link.icon
              const active = isActive(link.href)
              return (
                <Link key={link.href} href={link.href}>
                  <div className={`flex items-start gap-3 px-3 py-2 rounded-lg transition-colors ${
                    active
                      ? 'bg-primary-100 dark:bg-primary-900/30'
                      : 'hover:bg-gray-100 dark:hover:bg-gray-800'
                  }`}>
                    <Icon className={`h-4 w-4 mt-0.5 flex-shrink-0 ${
                      active
                        ? 'text-primary-700 dark:text-primary-400'
                        : 'text-gray-500 dark:text-gray-400'
                    }`} />
                    <div className="flex-1 min-w-0">
                      <p className={`text-sm font-medium ${
                        active
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
        {isAdmin && adminLinks.length > 0 && (
          <div className="mb-6 pt-4 border-t border-gray-200 dark:border-gray-800">
            <h3 className="text-xs font-semibold text-primary-600 dark:text-primary-400 uppercase tracking-wider px-3 py-2 mb-2 flex items-center gap-1">
              <Lock className="h-3 w-3" />
              Administration
            </h3>
            <div className="space-y-1">
              {adminLinks.map((link) => {
                const Icon = link.icon
                const active = isActive(link.href)
                return (
                  <Link key={link.href} href={link.href}>
                    <div className={`flex items-start gap-3 px-3 py-2 rounded-lg transition-colors ${
                      active
                        ? 'bg-primary-100 dark:bg-primary-900/30'
                        : 'hover:bg-gray-100 dark:hover:bg-gray-800'
                    }`}>
                      <Icon className={`h-4 w-4 mt-0.5 flex-shrink-0 ${
                        active
                          ? 'text-primary-700 dark:text-primary-400'
                          : 'text-gray-500 dark:text-gray-400'
                      }`} />
                      <div className="flex-1 min-w-0">
                        <p className={`text-sm font-medium ${
                          active
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
          <button className="w-full flex items-start gap-3 px-3 py-2 rounded-lg text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-left">
            <Globe className="h-4 w-4 mt-0.5 flex-shrink-0" />
            <div>
              <p className="text-sm font-medium">Help & Support</p>
              <p className="text-xs text-gray-500 dark:text-gray-400">Contact administrator</p>
            </div>
          </button>
        </div>
      </nav>
    </aside>
  )
}
