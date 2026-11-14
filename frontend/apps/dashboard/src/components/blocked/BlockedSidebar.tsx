'use client'

import { useEffect } from 'react'
import { ChefHat, CreditCard, Lock, X } from 'lucide-react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { getLocaleFromPath } from '@/lib/translations'
import { useTenantStatus } from '@/hooks/useTenantStatus'
import { BlockedLogout } from './BlockedOverlay'

interface BlockedSidebarProps {
  isMobileOpen: boolean
  onMobileToggle: () => void
}

const messages = {
  en: {
    title: 'Access restricted', description: 'Review your billing details or contact your account administrator.',
    billing: 'Billing & subscription', close: 'Close navigation',
    items: ['Dashboard', 'Orders', 'Products', 'Staff', 'Settings'],
  },
  ar: {
    title: 'الوصول مقيّد', description: 'راجع تفاصيل الفوترة أو تواصل مع مسؤول حسابك.',
    billing: 'الفوترة والاشتراك', close: 'إغلاق القائمة',
    items: ['لوحة التحكم', 'الطلبات', 'المنتجات', 'الموظفون', 'الإعدادات'],
  },
}

function BlockedNotice({ locale, hasBilling, onNavigate }: {
  locale: 'en' | 'ar'; hasBilling: boolean; onNavigate: () => void
}) {
  const copy = messages[locale]
  return <div className="p-4 space-y-3">
    <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-red-900 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300">
      <p className="flex items-center gap-2 text-sm font-semibold"><Lock className="h-5 w-5 shrink-0" aria-hidden="true" />{copy.title}</p>
      <p className="mt-2 text-xs">{copy.description}</p>
    </div>
    {hasBilling && <Link href={`/${locale}/dashboard/settings/billing`} onClick={onNavigate}
      className="flex min-h-12 items-center gap-2 rounded-lg border border-gray-300 px-3 text-sm text-gray-700 dark:border-gray-600 dark:text-gray-200">
      <CreditCard className="h-4 w-4" aria-hidden="true" />{copy.billing}
    </Link>}
    <ul className="space-y-2" aria-label={copy.title}>
      {copy.items.map(item => <li key={item} aria-disabled="true" className="rounded-lg px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{item}</li>)}
    </ul>
  </div>
}

export function BlockedSidebar({ isMobileOpen, onMobileToggle }: BlockedSidebarProps) {
  const locale = getLocaleFromPath(usePathname())
  const { blockReason } = useTenantStatus()
  const copy = messages[locale]
  useEffect(() => {
    if (!isMobileOpen) return
    const closeOnEscape = (event: KeyboardEvent) => { if (event.key === 'Escape') onMobileToggle() }
    window.addEventListener('keydown', closeOnEscape)
    return () => window.removeEventListener('keydown', closeOnEscape)
  }, [isMobileOpen, onMobileToggle])
  return <>
    <aside dir={locale === 'ar' ? 'rtl' : 'ltr'} aria-label={copy.title}
      className={`fixed inset-y-0 ${locale === 'ar' ? 'right-0 border-l' : 'left-0 border-r'} z-40 w-64 max-w-[calc(100vw-2rem)] flex-col border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800 ${isMobileOpen ? 'flex' : 'hidden md:flex'}`}>
      <div className="flex min-h-16 items-center justify-between border-b border-gray-200 px-4 dark:border-gray-700">
        <ChefHat className="h-8 w-8 text-gray-500" aria-hidden="true" />
        <button type="button" onClick={onMobileToggle} aria-label={copy.close}
          className="flex h-12 w-12 items-center justify-center rounded-md text-gray-700 md:hidden dark:text-gray-200"><X className="h-5 w-5" aria-hidden="true" /></button>
      </div>
      <div className="min-h-0 flex-1 overflow-y-auto">
        <BlockedNotice locale={locale} hasBilling={blockReason === 'trial_expired' || blockReason === 'expired'} onNavigate={() => { if (isMobileOpen) onMobileToggle() }} />
      </div>
      <div className="border-t border-gray-200 p-4 dark:border-gray-700"><BlockedLogout /></div>
    </aside>
    {isMobileOpen && <button type="button" onClick={onMobileToggle} aria-label={copy.close}
      className="fixed inset-0 z-30 bg-black/40 md:hidden" />}
  </>
}
