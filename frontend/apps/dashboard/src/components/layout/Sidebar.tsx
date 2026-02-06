'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { ChefHat, Store, ShoppingCart, Users, BarChart3, Settings, Menu, ChevronDown, ChevronRight, Package, FolderOpen, Archive, Palette, Truck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { useTenantStatus } from '@/hooks/useTenantStatus'
import { BlockedSidebar } from '@/components/blocked/BlockedSidebar'

interface SidebarItem {
  label: string
  href?: string
  icon: React.ReactNode
  children?: { label: string; href: string }[]
}

export function Sidebar() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const { isBlocked } = useTenantStatus()
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [expandedItems, setExpandedItems] = useState<string[]>([
    t('navigation.orders'),
    t('navigation.products'),
    t('navigation.hr'),
  ])

  // Show blocked sidebar if tenant is blocked
  if (isBlocked) {
    return <BlockedSidebar />
  }

  const sidebarItems: SidebarItem[] = [
    {
      label: t('navigation.dashboard'),
      href: `/dashboard`,
      icon: <BarChart3 className="h-5 w-5" />,
    },
    {
      label: t('navigation.orders'),
      icon: <Truck className="h-5 w-5" />,
      children: [
        { label: t('navigation.orders'), href: `/dashboard/orders` },
        { label: t('navigation.drivers'), href: `/dashboard/drivers` },
      ],
    },
    {
      label: t('navigation.products'),
      icon: <ShoppingCart className="h-5 w-5" />,
      children: [
        { label: t('navigation.products'), href: `/dashboard/products` },
        { label: t('navigation.categories'), href: `/dashboard/categories` },
        { label: t('navigation.inventory'), href: `/dashboard/inventory` },
      ],
    },
    {
      label: t('navigation.hr'),
      icon: <Users className="h-5 w-5" />,
      children: [
        { label: t('navigation.employees'), href: `/dashboard/hr/employees` },
        { label: t('navigation.roles'), href: `/dashboard/hr/roles` },
        { label: t('navigation.attendance'), href: `/dashboard/hr/attendance` },
        { label: t('navigation.leaves'), href: `/dashboard/hr/leaves` },
        { label: t('navigation.payroll'), href: `/dashboard/hr/payroll` },
      ],
    },
    {
      label: 'Theme Builder',
      href: `/dashboard/theme-builder`,
      icon: <Palette className="h-5 w-5" />,
    },
    {
      label: t('navigation.settings'),
      href: `/dashboard/settings`,
      icon: <Settings className="h-5 w-5" />,
    },
  ]

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen)
  }

  const toggleExpanded = (label: string) => {
    setExpandedItems((prev) =>
      prev.includes(label) ? prev.filter((item) => item !== label) : [...prev, label]
    )
  }

  const isRTL = locale === 'ar'

  return (
    <>
      {/* Mobile menu button */}
      <div className={`md:hidden fixed top-4 ${isRTL ? 'right-4' : 'left-4'} z-50`}>
        <Button
          variant="outline"
          size="icon"
          onClick={toggleMobileMenu}
          className="bg-white dark:bg-gray-800"
        >
          <Menu className="h-5 w-5" />
        </Button>
      </div>

      {/* Sidebar */}
      <aside
        className={`fixed inset-y-0 ${isRTL ? 'right-0' : 'left-0'} z-40 w-64 bg-white dark:bg-gray-800 ${isRTL ? 'border-l' : 'border-r'} border-gray-200 dark:border-gray-700 transform ${isMobileMenuOpen ? 'translate-x-0' : isRTL ? 'translate-x-full' : '-translate-x-full'
          } md:translate-x-0 transition-transform duration-300 ease-in-out`}
      >
        <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
          <Link href={`/${locale}/dashboard`} className="flex items-center gap-2 hover:opacity-80 transition">
            <ChefHat className="h-8 w-8 text-primary-600 dark:text-primary-400" />
            <span className="text-xl font-bold text-gray-900 dark:text-white">
              POS Dashboard
            </span>
          </Link>
        </div>

        <nav className="p-4">
          <ul className="space-y-2">
            {sidebarItems.map((item) => (
              <li key={item.label}>
                {item.children ? (
                  <>
                    <button
                      onClick={() => toggleExpanded(item.label)}
                      className={`w-full flex items-center justify-between gap-3 px-4 py-3 rounded-lg transition text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700`}
                    >
                      <div className="flex items-center gap-3">
                        {item.icon}
                        <span className="font-medium">{item.label}</span>
                      </div>
                      {expandedItems.includes(item.label) ? (
                        <ChevronDown className="h-4 w-4 sidebar-chevron-down" />
                      ) : (
                        <ChevronRight className="h-4 w-4 sidebar-chevron-right" />
                      )}
                    </button>
                    {expandedItems.includes(item.label) && (
                      <ul className={`mt-2 ${isRTL ? 'mr-4' : 'ml-4'} space-y-1`}>
                        {item.children.map((child) => (
                          <li key={child.href}>
                            <Link
                              href={`/${locale}${child.href}`}
                              className={`flex items-center gap-3 px-4 py-2 rounded-lg transition text-sm ${pathname === `/${locale}${child.href}`
                                ? 'bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 font-medium'
                                : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'
                                }`}
                              onClick={() => setIsMobileMenuOpen(false)}
                            >
                              {child.label}
                            </Link>
                          </li>
                        ))}
                      </ul>
                    )}
                  </>
                ) : (
                  <Link
                    href={`/${locale}${item.href}`}
                    className={`flex items-center gap-3 px-4 py-3 rounded-lg transition ${pathname === `/${locale}${item.href}`
                      ? 'bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400'
                      : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                      }`}
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    {item.icon}
                    <span className="font-medium">{item.label}</span>
                  </Link>
                )}
              </li>
            ))}
          </ul>
        </nav>
      </aside>

      {/* Overlay for mobile */}
      {isMobileMenuOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-30 md:hidden"
          onClick={toggleMobileMenu}
        ></div>
      )}
    </>
  )
}
