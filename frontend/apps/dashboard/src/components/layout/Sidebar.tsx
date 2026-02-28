'use client'

import { useState, useMemo } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { ChefHat, ShoppingCart, Users, BarChart3, Settings, Menu, ChevronDown, ChevronRight, ChevronLeft, Palette, Truck, Monitor, CreditCard } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { useTenantStatus } from '@/hooks/useTenantStatus'
import { useHasPermission } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { BlockedSidebar } from '@/components/blocked/BlockedSidebar'
import { SidebarSubscriptionBadge } from '@/components/layout/SidebarSubscriptionBadge'

interface SidebarChild {
  label: string
  href: string
  moduleId?: number
  requiredPermission?: string
}

interface SidebarItem {
  label: string
  href?: string
  icon: React.ReactNode
  children?: SidebarChild[]
  moduleId?: number  // For permission checking
  requiredPermission?: string  // Default: READ
}

export function Sidebar() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const { isBlocked } = useTenantStatus()
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [expandedItems, setExpandedItems] = useState<string[]>([
    t('navigation.reports'),
    t('navigation.orders'),
    t('navigation.products'),
    t('navigation.hr'),
  ])

  // Permission checks for each module
  const { hasPermission: canViewReports } = useHasPermission(MODULE_IDS.REPORTS, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewOrders } = useHasPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.READ)
  const { hasPermission: canWriteOrders } = useHasPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.WRITE)
  const { hasPermission: canViewProducts } = useHasPermission(MODULE_IDS.PRODUCTS, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewCategories } = useHasPermission(MODULE_IDS.CATEGORIES, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewHR } = useHasPermission(MODULE_IDS.HR, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewLeaves } = useHasPermission(MODULE_IDS.LEAVES, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewThemes } = useHasPermission(MODULE_IDS.THEMES, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewKDS } = useHasPermission(MODULE_IDS.KDS, PERMISSION_LEVELS.READ)
  const { hasPermission: canManageSalaries } = useHasPermission(MODULE_IDS.SALARIES, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewTables } = useHasPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.READ)

  // Show Products menu if user has either PRODUCTS or CATEGORIES permission
  const canViewProductsMenu = canViewProducts || canViewCategories
  // Show HR menu if user has HR, LEAVES, or SALARIES permission
  const canViewHRMenu = canViewHR || canViewLeaves || canManageSalaries
  // Show Orders menu if user has ORDERS permission
  const canViewOrdersMenu = canViewOrders

  // Show blocked sidebar if tenant is blocked
  if (isBlocked) {
    return <BlockedSidebar />
  }

  const allSidebarItems: SidebarItem[] = [
    {
      label: t('navigation.dashboard'),
      href: `/dashboard`,
      icon: <BarChart3 className="h-5 w-5" />,
      // Dashboard always visible
    },
    {
      label: t('navigation.kds', 'Kitchen Display'),
      href: `/dashboard/kds`,
      icon: <Monitor className="h-5 w-5" />,
      moduleId: MODULE_IDS.KDS,
      requiredPermission: PERMISSION_LEVELS.READ,
    },
    {
      label: t('navigation.pos', 'Cashier POS'),
      href: `/dashboard/pos`,
      icon: <CreditCard className="h-5 w-5" />,
      moduleId: MODULE_IDS.ORDERS,
      requiredPermission: PERMISSION_LEVELS.WRITE,
    },
    {
      label: t('navigation.reports'),
      icon: <BarChart3 className="h-5 w-5" />,
      moduleId: MODULE_IDS.REPORTS,
      requiredPermission: PERMISSION_LEVELS.READ,
      children: [
        { label: t('navigation.reportsOverview'), href: `/dashboard/reports/overview` },
        { label: t('navigation.reportsSales'), href: `/dashboard/reports/sales` },
        { label: t('navigation.reportsProfitability'), href: `/dashboard/reports/profitability` },
        { label: t('navigation.reportsInventory'), href: `/dashboard/reports/inventory` },
      ],
    },
    {
      label: t('navigation.orders'),
      icon: <Truck className="h-5 w-5" />,
      moduleId: MODULE_IDS.ORDERS,
      requiredPermission: PERMISSION_LEVELS.READ,
      children: [
        { label: t('navigation.orders'), href: `/dashboard/orders`, moduleId: MODULE_IDS.ORDERS },
        { label: t('navigation.tables'), href: `/dashboard/tables`, moduleId: MODULE_IDS.TABLES },
        { label: t('navigation.drivers'), href: `/dashboard/drivers`, moduleId: MODULE_IDS.ORDERS },
      ],
    },
    {
      label: t('navigation.products'),
      icon: <ShoppingCart className="h-5 w-5" />,
      moduleId: MODULE_IDS.PRODUCTS,
      requiredPermission: PERMISSION_LEVELS.READ,
      children: [
        { label: t('navigation.products'), href: `/dashboard/products`, moduleId: MODULE_IDS.PRODUCTS },
        { label: t('navigation.categories'), href: `/dashboard/categories`, moduleId: MODULE_IDS.CATEGORIES },
        { label: t('navigation.inventory'), href: `/dashboard/inventory`, moduleId: MODULE_IDS.PRODUCTS },
      ],
    },
    {
      label: t('navigation.hr'),
      icon: <Users className="h-5 w-5" />,
      moduleId: MODULE_IDS.HR,
      requiredPermission: PERMISSION_LEVELS.READ,
      children: [
        { label: t('navigation.employees'), href: `/dashboard/hr/employees` },
        { label: t('navigation.attendance'), href: `/dashboard/hr/attendance` },
        { label: t('navigation.leaves'), href: `/dashboard/hr/leaves`, moduleId: MODULE_IDS.LEAVES },
        { label: t('navigation.payroll'), href: `/dashboard/hr/payroll`, moduleId: MODULE_IDS.SALARIES },
      ],
    },
    {
      label: t('navigation.themeBuilder', 'Theme Builder'),
      href: `/dashboard/theme-builder`,
      icon: <Palette className="h-5 w-5" />,
      moduleId: MODULE_IDS.THEMES,
      requiredPermission: PERMISSION_LEVELS.READ,
    },
    {
      label: t('navigation.settings'),
      href: `/dashboard/settings`,
      icon: <Settings className="h-5 w-5" />,
      // Settings always visible
    },
  ]

  // Filter sidebar items based on user permissions
  const sidebarItems: SidebarItem[] = useMemo(() => {
    const permissionMap: Record<number, boolean> = {
      [MODULE_IDS.REPORTS]: canViewReports,
      [MODULE_IDS.ORDERS]: canViewOrders,
      [MODULE_IDS.PRODUCTS]: canViewProducts,
      [MODULE_IDS.CATEGORIES]: canViewCategories,
      [MODULE_IDS.HR]: canViewHR,
      [MODULE_IDS.LEAVES]: canViewLeaves,
      [MODULE_IDS.THEMES]: canViewThemes,
      [MODULE_IDS.KDS]: canViewKDS,
      [MODULE_IDS.SALARIES]: canManageSalaries,
      [MODULE_IDS.TABLES]: canViewTables,
    }

    return allSidebarItems
      .filter((item) => {
        // Always show items without moduleId requirement
        if (!item.moduleId) {
          return true
        }

        // Orders group should remain visible when user only has TABLES permission.
        if (item.moduleId === MODULE_IDS.ORDERS) {
          if (item.href === '/dashboard/pos') {
            return canWriteOrders
          }
          return (permissionMap[MODULE_IDS.ORDERS] ?? false) || (permissionMap[MODULE_IDS.TABLES] ?? false)
        }

        // Show item if user has permission
        return permissionMap[item.moduleId] ?? false
      })
      .map((item) => {
        // Filter children based on their individual permissions
        if (item.children) {
          return {
            ...item,
            children: item.children.filter((child) => {
              // If child has no moduleId, show it (user must have parent permission)
              if (!child.moduleId) {
                return true
              }
              // Show child if user has permission for that module
              return permissionMap[child.moduleId] ?? false
            }),
          }
        }
        return item
      })
  }, [canViewReports, canViewOrders, canWriteOrders, canViewProducts, canViewCategories, canViewHR, canViewLeaves, canViewThemes, canViewKDS, canManageSalaries, canViewTables])

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
          <Link
            href={`/${locale}/dashboard`}
            className={`flex items-center gap-2 hover:opacity-80 transition ${isRTL ? 'flex-row-reverse' : ''}`}
          >
            <ChefHat className="h-8 w-8 text-primary-600 dark:text-primary-400" />
            <span className="text-xl font-bold text-gray-900 dark:text-white">
              POS Dashboard
            </span>
          </Link>
        </div>

        <div className="flex flex-col h-[calc(100%-65px)]">
        <nav className="p-4 flex-1 overflow-y-auto">
          <ul className="space-y-2">
            {sidebarItems.map((item) => (
              <li key={item.label}>
                {item.children ? (
                  <>
                    <button
                      onClick={() => toggleExpanded(item.label)}
                      className={`w-full flex items-center justify-between gap-3 px-4 py-3 rounded-lg transition text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 ${isRTL ? 'text-right' : 'text-left'}`}
                    >
                      <div className={`flex items-center gap-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
                        {item.icon}
                        <span className="font-medium">{item.label}</span>
                      </div>
                      {expandedItems.includes(item.label) ? (
                        <ChevronDown className="h-4 w-4 sidebar-chevron-down" />
                      ) : (
                        (isRTL ? <ChevronLeft className="h-4 w-4 sidebar-chevron-right" /> : <ChevronRight className="h-4 w-4 sidebar-chevron-right" />)
                      )}
                    </button>
                    {expandedItems.includes(item.label) && (
                      <ul className={`mt-2 ${isRTL ? 'pr-4' : 'pl-4'} space-y-1`}>
                        {item.children.map((child) => (
                          <li key={child.href}>
                            <Link
                              href={`/${locale}${child.href}`}
                              className={`flex items-center gap-3 px-4 py-2 rounded-lg transition text-sm ${pathname === `/${locale}${child.href}`
                                ? 'bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 font-medium'
                                : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'
                                } ${isRTL ? 'justify-end text-right' : ''}`}
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
                      } ${isRTL ? 'flex-row-reverse text-right' : ''}`}
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
        <SidebarSubscriptionBadge />
        </div>
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
