/**
 * Optimized Navigation Link Component
 *
 * This component preloads route code on hover/focus to improve perceived performance.
 * Usage:
 * <OptimizedNavLink href="/dashboard/products" label="Products" preload />
 */

import React, { ReactNode } from 'react'
import Link from 'next/link'
import { usePreloadRoute } from '@/hooks/usePreloadRoute'

interface OptimizedNavLinkProps {
  href: string
  label: ReactNode
  icon?: ReactNode
  badge?: number | string
  preload?: boolean
  preloadDelay?: number
  active?: boolean
  className?: string
  children?: ReactNode
}

export function OptimizedNavLink({
  href,
  label,
  icon,
  badge,
  preload = true,
  preloadDelay = 100,
  active = false,
  className = '',
  children,
}: OptimizedNavLinkProps) {
  // For demonstration purposes, we return the import function
  // In a real app, you'd map the href to the actual import
  const getImportFn = () => {
    // This would be replaced with actual route-specific imports
    return () => Promise.resolve()
  }

  const handlers = usePreloadRoute(getImportFn(), { delay: preloadDelay })
  const onMouseEnter = preload ? handlers.onMouseEnter : undefined
  const onFocus = preload ? handlers.onFocus : undefined

  return (
    <Link href={href}
        className={`
          flex items-center gap-3 px-4 py-2 rounded-lg transition-colors
          ${active
            ? 'bg-blue-50 text-blue-600 font-medium'
            : 'text-gray-700 hover:bg-gray-50'
          }
          ${className}
        `}
        onMouseEnter={onMouseEnter}
        onFocus={onFocus}
      >
        {icon && <span className="flex-shrink-0">{icon}</span>}
        <span className="flex-1">{label}</span>
        {badge && (
          <span className="flex items-center justify-center w-5 h-5 text-xs font-bold text-white bg-red-500 rounded-full">
            {badge}
          </span>
        )}
    </Link>
  )
}

/**
 * Sidebar Navigation Component with Route Preloading
 */
export function SidebarNavigation() {
  const navItems = [
    {
      label: 'Dashboard',
      href: '/dashboard',
      icon: '📊',
    },
    {
      label: 'Products',
      href: '/dashboard/products',
      icon: '📦',
    },
    {
      label: 'Categories',
      href: '/dashboard/categories',
      icon: '🏷️',
    },
    {
      label: 'Theme Builder',
      href: '/dashboard/theme-builder',
      icon: '🎨',
    },
    {
      label: 'Inventory',
      href: '/dashboard/inventory',
      icon: '📈',
    },
    {
      label: 'HR',
      href: '/dashboard/hr/employees',
      icon: '👥',
    },
    {
      label: 'Notifications',
      href: '/dashboard/notifications',
      icon: '🔔',
    },
    {
      label: 'Settings',
      href: '/dashboard/settings',
      icon: '⚙️',
    },
  ]

  return (
    <nav className="space-y-2">
      {navItems.map((item) => (
        <OptimizedNavLink
          key={item.href}
          href={item.href}
          label={item.label}
          icon={item.icon}
          preload
        />
      ))}
    </nav>
  )
}
