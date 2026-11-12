'use client'

import { useMemo, useState } from 'react'
import Link from 'next/link'
import { usePathname, useSearchParams } from 'next/navigation'
import { ChefHat, ShoppingCart, Users, BarChart3, Settings, ChevronDown, ChevronRight, ChevronLeft, Palette, Truck, Monitor, CreditCard, BookOpen, Building2, PanelLeftClose, PanelLeftOpen } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { useTenantStatus } from '@/hooks/useTenantStatus'
import { useUserPermissions } from '@/hooks/useRbac'
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

interface SidebarProps {
  isMobileOpen: boolean
  onMobileToggle: () => void
  isCollapsed: boolean
  onCollapsedChange: (collapsed: boolean) => void
}

export function Sidebar({ isMobileOpen, onMobileToggle, isCollapsed, onCollapsedChange }: SidebarProps) {
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const { isBlocked } = useTenantStatus()
  const [expandedItems, setExpandedItems] = useState<string[]>([])
  const { checkPermission, loading: permissionsLoading } = useUserPermissions()

  const canViewReports = checkPermission(MODULE_IDS.REPORTS, PERMISSION_LEVELS.READ)
  const canViewOrders = checkPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.READ)
  const canWriteOrders = checkPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.WRITE)
  const canViewProducts = checkPermission(MODULE_IDS.PRODUCTS, PERMISSION_LEVELS.READ)
  const canViewCategories = checkPermission(MODULE_IDS.CATEGORIES, PERMISSION_LEVELS.READ)
  const canViewHR = checkPermission(MODULE_IDS.HR, PERMISSION_LEVELS.READ)
  const canViewLeaves = checkPermission(MODULE_IDS.LEAVES, PERMISSION_LEVELS.READ)
  const canViewThemes = checkPermission(MODULE_IDS.THEMES, PERMISSION_LEVELS.READ)
  const canViewKDS = checkPermission(MODULE_IDS.KDS, PERMISSION_LEVELS.READ)
  const canManageSalaries = checkPermission(MODULE_IDS.SALARIES, PERMISSION_LEVELS.READ)
  const canViewTables = checkPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.READ)
  const canViewAccounting = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.READ)
  const canViewSuppliers = checkPermission(MODULE_IDS.SUPPLIERS, PERMISSION_LEVELS.READ)
  const canViewFinancialReports = checkPermission(MODULE_IDS.FINANCIAL_REPORTS, PERMISSION_LEVELS.READ)
  const canViewPerformance = checkPermission(MODULE_IDS.PERFORMANCE, PERMISSION_LEVELS.READ)
  const canViewTraining = checkPermission(MODULE_IDS.TRAINING, PERMISSION_LEVELS.READ)
  const canViewContracts = checkPermission(MODULE_IDS.CONTRACTS, PERMISSION_LEVELS.READ)

  const allSidebarItems: SidebarItem[] = [
    {
      label: t('navigation.dashboard'),
      href: `/dashboard`,
      icon: <BarChart3 className="h-5 w-5" />,
      // Dashboard always visible
    },
    {
      label: t('navigation.pos', 'Cashier POS'),
      href: `/dashboard/pos`,
      icon: <CreditCard className="h-5 w-5" />,
      moduleId: MODULE_IDS.ORDERS,
      requiredPermission: PERMISSION_LEVELS.WRITE,
    },
    {
      label: t('navigation.kds', 'Kitchen Display'),
      href: `/dashboard/kds`,
      icon: <Monitor className="h-5 w-5" />,
      moduleId: MODULE_IDS.KDS,
      requiredPermission: PERMISSION_LEVELS.READ,
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
        { label: t('navigation.performance', 'Performance'), href: `/dashboard/hr/performance`, moduleId: MODULE_IDS.PERFORMANCE },
        { label: t('navigation.training', 'Training'), href: `/dashboard/hr/training`, moduleId: MODULE_IDS.TRAINING },
        { label: t('navigation.insurance', 'Insurance'), href: `/dashboard/hr/insurance`, moduleId: MODULE_IDS.CONTRACTS },
        { label: t('navigation.contracts', 'Contracts'), href: `/dashboard/hr/contracts`, moduleId: MODULE_IDS.CONTRACTS },
        { label: t('navigation.resignations', 'Resignations'), href: `/dashboard/hr/resignations`, moduleId: MODULE_IDS.CONTRACTS },
      ],
    },
    {
      label: t('navigation.accounting', 'Accounting'),
      icon: <BookOpen className="h-5 w-5" />,
      moduleId: MODULE_IDS.ACCOUNTING,
      requiredPermission: PERMISSION_LEVELS.READ,
      children: [
        {
          label: t('navigation.chartOfAccounts', 'Chart of Accounts'),
          href: `/dashboard/accounting/chart-of-accounts`,
          moduleId: MODULE_IDS.ACCOUNTING,
        },
        {
          label: t('navigation.journalEntries', 'Journal Entries'),
          href: `/dashboard/accounting/journal-entries`,
          moduleId: MODULE_IDS.ACCOUNTING,
        },
        {
          label: t('navigation.generalLedger', 'General Ledger'),
          href: `/dashboard/accounting/general-ledger`,
          moduleId: MODULE_IDS.ACCOUNTING,
        },
        {
          label: t('navigation.bankAccounts', 'Bank Accounts'),
          href: `/dashboard/accounting/bank-accounts`,
          moduleId: MODULE_IDS.ACCOUNTING,
        },
        {
          label: t('navigation.advancedAccounting', 'Advanced Accounting'),
          href: `/dashboard/accounting/advanced`,
          moduleId: MODULE_IDS.ACCOUNTING,
        },
        {
          label: t('navigation.customersAccounting', 'Customers'),
          href: `/dashboard/accounting/customers`,
          moduleId: MODULE_IDS.ACCOUNTING,
        },
        {
          label: t('navigation.customerInvoices', 'Customer Invoices'),
          href: `/dashboard/accounting/customer-invoices`,
          moduleId: MODULE_IDS.ACCOUNTING,
        },
        {
          label: t('navigation.receivedPayments', 'Received Payments'),
          href: `/dashboard/accounting/customer-invoices?tab=payments`,
          moduleId: MODULE_IDS.ACCOUNTING,
        },
      ],
    },
    {
      label: t('navigation.suppliersAccounting', 'Suppliers'),
      icon: <Building2 className="h-5 w-5" />,
      moduleId: MODULE_IDS.SUPPLIERS,
      requiredPermission: PERMISSION_LEVELS.READ,
      children: [
        {
          label: t('navigation.supplierList', 'Supplier List'),
          href: `/dashboard/accounting/suppliers`,
          moduleId: MODULE_IDS.SUPPLIERS,
        },
        {
          label: t('navigation.supplierInvoices', 'Supplier Invoices'),
          href: `/dashboard/accounting/supplier-invoices`,
          moduleId: MODULE_IDS.ACCOUNTING,
        },
        {
          label: t('navigation.madePayments', 'Payments Made'),
          href: `/dashboard/accounting/supplier-invoices?tab=payments`,
          moduleId: MODULE_IDS.ACCOUNTING,
        },
      ],
    },
    {
      label: t('navigation.financialReports', 'Financial Reports'),
      icon: <BarChart3 className="h-5 w-5" />,
      moduleId: MODULE_IDS.FINANCIAL_REPORTS,
      requiredPermission: PERMISSION_LEVELS.READ,
      children: [
        { label: t('navigation.financialHub', 'Financial Hub'), href: `/dashboard/reports/financial` },
        { label: t('navigation.profitLoss', 'Profit & Loss'), href: `/dashboard/reports/financial/profit-loss` },
        { label: t('navigation.balanceSheet', 'Balance Sheet'), href: `/dashboard/reports/financial/balance-sheet` },
        { label: t('navigation.cashFlow', 'Cash Flow'), href: `/dashboard/reports/financial/cash-flow` },
        { label: t('navigation.vatReport', 'VAT Report'), href: `/dashboard/reports/financial/vat` },
        { label: t('reports.financial.receivablesAging', 'Receivables Aging'), href: `/dashboard/reports/financial/receivables-aging` },
        { label: t('reports.financial.payablesAging', 'Payables Aging'), href: `/dashboard/reports/financial/payables-aging` },
        { label: t('reports.financial.inventoryValuation', 'Inventory Valuation'), href: `/dashboard/reports/financial/inventory-valuation` },
        { label: t('navigation.financialAdvanced', 'Advanced Insights'), href: `/dashboard/reports/financial/advanced` },
      ],
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
        { label: t('navigation.hrReports', 'HR Reports'), href: `/dashboard/reports/hr` },
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
      [MODULE_IDS.ACCOUNTING]: canViewAccounting,
      [MODULE_IDS.SUPPLIERS]: canViewSuppliers,
      [MODULE_IDS.FINANCIAL_REPORTS]: canViewFinancialReports,
      [MODULE_IDS.PERFORMANCE]: canViewPerformance,
      [MODULE_IDS.TRAINING]: canViewTraining,
      [MODULE_IDS.CONTRACTS]: canViewContracts,
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

        if (item.moduleId === MODULE_IDS.SUPPLIERS) {
          return (permissionMap[MODULE_IDS.SUPPLIERS] ?? false) || (permissionMap[MODULE_IDS.ACCOUNTING] ?? false)
        }

        // Keep Financial Reports visible for existing report-enabled roles
        // even if the dedicated FINANCIAL_REPORTS module wasn't seeded yet.
        if (item.moduleId === MODULE_IDS.FINANCIAL_REPORTS) {
          return (permissionMap[MODULE_IDS.FINANCIAL_REPORTS] ?? false) || (permissionMap[MODULE_IDS.REPORTS] ?? false)
        }

        // Show item if user has permission
        return permissionMap[item.moduleId] ?? false
      })
      .map((item) => {
        // Filter children based on their individual permissions
        if (item.children) {
          const filteredChildren = item.children.filter((child) => {
            // If child has no moduleId, show it (user must have parent permission)
            if (!child.moduleId) {
              return true
            }
            const isAdvancedHRChild =
              child.moduleId === MODULE_IDS.PERFORMANCE ||
              child.moduleId === MODULE_IDS.TRAINING ||
              child.moduleId === MODULE_IDS.CONTRACTS

            // Show advanced HR pages when base HR access is granted.
            // This avoids hiding newly added pages when module 18-20
            // permissions are missing for older seeded roles.
            if (isAdvancedHRChild && (permissionMap[MODULE_IDS.HR] ?? false)) {
              return true
            }
            // Show child if user has permission for that module
            return permissionMap[child.moduleId] ?? false
          })

          return {
            ...item,
            children: filteredChildren,
          }
        }
        return item
      })
      .filter((item) => !item.children || item.children.length > 0)
  }, [canViewReports, canViewOrders, canWriteOrders, canViewProducts, canViewCategories, canViewHR, canViewLeaves, canViewThemes, canViewKDS, canManageSalaries, canViewTables, canViewAccounting, canViewSuppliers, canViewFinancialReports, canViewPerformance, canViewTraining, canViewContracts])

  // Keep hooks above conditional returns so loading/status changes preserve hook order.
  if (isBlocked) {
    return <BlockedSidebar />
  }

  if (permissionsLoading) {
    return (
      <aside
        className={`fixed inset-y-0 ${locale === 'ar' ? 'right-0' : 'left-0'} z-40 ${isCollapsed ? 'w-20' : 'w-64'} bg-white dark:bg-gray-800 ${locale === 'ar' ? 'border-l' : 'border-r'} border-gray-200 dark:border-gray-700 transform ${isMobileOpen ? 'translate-x-0' : locale === 'ar' ? 'translate-x-full' : '-translate-x-full'} md:translate-x-0 transition-all duration-300 ease-in-out`}
      >
        <div className={`flex h-[89px] items-center border-b border-gray-200 p-4 dark:border-gray-700 ${isCollapsed ? 'justify-center' : 'justify-between'}`}>
          <div className={`h-8 rounded bg-gray-200 dark:bg-gray-700 animate-pulse ${isCollapsed ? 'w-10' : 'w-40'}`} />
          {!isCollapsed && <div className="h-8 w-8 rounded bg-gray-200 dark:bg-gray-700 animate-pulse" />}
        </div>
        <div className={`flex h-[calc(100%-89px)] flex-col p-3 ${isCollapsed ? 'px-2' : 'px-4'}`}>
          <div className="space-y-2">
          {Array.from({ length: 10 }).map((_, i) => (
            <div
              key={i}
              className={`rounded-lg bg-gray-200 dark:bg-gray-700 animate-pulse ${isCollapsed ? 'h-11 w-11 mx-auto' : 'h-11 w-full'}`}
            />
          ))}
          </div>
        </div>
      </aside>
    )
  }

  const toggleExpanded = (label: string) => {
    if (isCollapsed) return
    setExpandedItems((prev) =>
      prev.includes(label) ? prev.filter((item) => item !== label) : [...prev, label]
    )
  }

  const isRTL = locale === 'ar'
  const currentTab = searchParams.get('tab') || ''

  const isChildActive = (childHref: string) => {
    const [basePath, queryString] = childHref.split('?')
    if (pathname !== `/${locale}${basePath}`) {
      return false
    }
    if (!queryString) {
      return true
    }
    const query = new URLSearchParams(queryString)
    const childTab = query.get('tab') || ''
    return childTab === currentTab
  }

  return (
    <>
      {/* Sidebar */}
      <aside
        className={`fixed inset-y-0 ${isRTL ? 'right-0' : 'left-0'} z-40 ${isCollapsed ? 'w-20' : 'w-64'} bg-white dark:bg-gray-800 ${isRTL ? 'border-l' : 'border-r'} border-gray-200 dark:border-gray-700 transform ${isMobileOpen ? 'translate-x-0' : isRTL ? 'translate-x-full' : '-translate-x-full'} md:translate-x-0 transition-all duration-300 ease-in-out`}
      >
        <div className={`flex h-[89px] items-center ${isCollapsed ? 'justify-center' : 'justify-between'} p-4 border-b border-gray-200 dark:border-gray-700`}>
          <Link
            href={`/${locale}/dashboard`}
            className={`flex items-center ${isCollapsed ? 'justify-center' : 'gap-2'} hover:opacity-80 transition ${isRTL && !isCollapsed ? 'flex-row-reverse' : ''}`}
            title={t('navigation.dashboard')}
          >
            <ChefHat className="h-8 w-8 text-primary-600 dark:text-primary-400" />
            {!isCollapsed && <span className="text-xl font-bold text-gray-900 dark:text-white">
              POS Dashboard
            </span>}
          </Link>
          <button
            type="button"
            onClick={() => onCollapsedChange(!isCollapsed)}
            className={`hidden md:inline-flex items-center justify-center rounded-lg p-2 text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 transition ${isRTL ? 'mr-1' : 'ml-1'}`}
            title={isCollapsed ? t('common.open', 'Expand sidebar') : t('common.close', 'Collapse sidebar')}
          >
            {isCollapsed ? <PanelLeftOpen className="h-4 w-4" /> : <PanelLeftClose className="h-4 w-4" />}
          </button>
        </div>

        <div className="flex flex-col h-[calc(100%-89px)]">
        <nav className={`p-3 flex-1 overflow-y-auto ${isCollapsed ? 'px-2' : 'px-4'}`}>
          <ul className="space-y-2">
            {sidebarItems.map((item) => (
              <li key={item.label}>
                {item.children ? (
                  <>
                    <button
                      onClick={() => toggleExpanded(item.label)}
                      title={item.label}
                      className={`w-full flex items-center ${isCollapsed ? 'justify-center' : 'justify-between'} gap-3 ${isCollapsed ? 'px-2 py-3' : 'px-4 py-3'} rounded-lg transition text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 ${isRTL && !isCollapsed ? 'text-right' : 'text-left'}`}
                    >
                      <div className={`flex items-center gap-3 ${isRTL && !isCollapsed ? 'flex-row-reverse' : ''}`}>
                        {item.icon}
                        {!isCollapsed && <span className="font-medium">{item.label}</span>}
                      </div>
                      {!isCollapsed && (
                        expandedItems.includes(item.label) ? (
                          <ChevronDown className="h-4 w-4 sidebar-chevron-down" />
                        ) : (
                          (isRTL ? <ChevronLeft className="h-4 w-4 sidebar-chevron-right" /> : <ChevronRight className="h-4 w-4 sidebar-chevron-right" />)
                        )
                      )}
                    </button>
                    {!isCollapsed && expandedItems.includes(item.label) && (
                      <ul className={`mt-2 ${isRTL ? 'pr-4' : 'pl-4'} space-y-1`}>
                        {item.children.map((child) => (
                          <li key={child.href}>
                            <Link
                              href={`/${locale}${child.href}`}
                              className={`flex items-center gap-3 px-4 py-2 rounded-lg transition text-sm ${isChildActive(child.href)
                                ? 'bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 font-medium'
                                : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'
                                } ${isRTL ? 'justify-end text-right' : ''}`}
                              onClick={() => {
                                if (isMobileOpen) onMobileToggle()
                              }}
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
                    title={item.label}
                    className={`flex items-center ${isCollapsed ? 'justify-center' : 'gap-3'} ${isCollapsed ? 'px-2 py-3' : 'px-4 py-3'} rounded-lg transition ${pathname === `/${locale}${item.href}`
                      ? 'bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400'
                      : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                      } ${isRTL && !isCollapsed ? 'flex-row-reverse text-right' : ''}`}
                    onClick={() => {
                      if (isMobileOpen) onMobileToggle()
                    }}
                  >
                    {item.icon}
                    {!isCollapsed && <span className="font-medium">{item.label}</span>}
                  </Link>
                )}
              </li>
            ))}
          </ul>
        </nav>
        {!isCollapsed && <SidebarSubscriptionBadge />}
        </div>
      </aside>

      {/* Overlay for mobile */}
      {isMobileOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-30 md:hidden"
          onClick={onMobileToggle}
        ></div>
      )}
    </>
  )
}
