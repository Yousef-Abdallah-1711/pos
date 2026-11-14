'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { Clock, CreditCard, Lock } from 'lucide-react'
import { Button } from '@/components/ui/button'
import api from '@/lib/api'
import { getDashboardLoginUrl } from '@/lib/authSession'
import { useAuthStore } from '@/stores/authStore'
import { getLocaleFromPath } from '@/lib/translations'

interface BlockedOverlayProps {
  blockReason: string
  supportEmail?: string
}

const messages = {
  en: {
    trialTitle: 'Your trial has expired', paidTitle: 'Your subscription has expired',
    expiredDescription: 'Review your account’s billing details and access options.',
    disabledTitle: 'Account access disabled', disabledDescription: 'Your account has been temporarily disabled.',
    restrictedTitle: 'Access restricted', restrictedDescription: 'Your account access is currently restricted.',
    billing: 'Review billing options', note: 'Opening billing does not authorize a charge.',
    help: 'Contact your account administrator if you need help restoring access.', contact: 'Contact support',
    signOut: 'Sign out', signingOut: 'Signing out…', logoutError: 'Sign out failed. Please try again.',
  },
  ar: {
    trialTitle: 'انتهت الفترة التجريبية', paidTitle: 'انتهت صلاحية الاشتراك',
    expiredDescription: 'راجع تفاصيل الفوترة وخيارات الوصول الخاصة بحسابك.',
    disabledTitle: 'تم تعطيل الوصول إلى الحساب', disabledDescription: 'تم تعطيل حسابك مؤقتًا.',
    restrictedTitle: 'الوصول مقيّد', restrictedDescription: 'الوصول إلى حسابك مقيّد حاليًا.',
    billing: 'مراجعة خيارات الفوترة', note: 'فتح صفحة الفوترة لا يمنح إذنًا بخصم أي مبلغ.',
    help: 'تواصل مع مسؤول حسابك إذا احتجت إلى مساعدة لاستعادة الوصول.', contact: 'التواصل مع الدعم',
    signOut: 'تسجيل الخروج', signingOut: 'جارٍ تسجيل الخروج…', logoutError: 'تعذّر تسجيل الخروج. حاول مرة أخرى.',
  },
}
type Copy = typeof messages.en

function blockContent(reason: string, copy: Copy) {
  const hasBilling = reason === 'trial_expired' || reason === 'expired'
  if (hasBilling) return { Icon: Clock, title: reason === 'expired' ? copy.paidTitle : copy.trialTitle, description: copy.expiredDescription, hasBilling }
  const disabled = reason === 'tenant_disabled'
  return { Icon: Lock, title: disabled ? copy.disabledTitle : copy.restrictedTitle, description: disabled ? copy.disabledDescription : copy.restrictedDescription, hasBilling }
}

function BlockedActions({ copy, billingHref, hasBilling, supportEmail }: {
  copy: Copy; billingHref: string; hasBilling: boolean; supportEmail?: string
}) {
  return <div className="space-y-3 mb-4">
    {hasBilling && <>
      <Button asChild className="h-12 w-full"><Link href={billingHref}>
        <CreditCard className="h-4 w-4 me-2" aria-hidden="true" />{copy.billing}
      </Link></Button>
      <p className="text-xs text-center text-gray-600 dark:text-gray-300">{copy.note}</p>
    </>}
    <p className="text-sm text-center text-gray-600 dark:text-gray-300">{copy.help}</p>
    {supportEmail && <Button asChild className="h-12 w-full" variant="outline">
      <a href={`mailto:${supportEmail}`}>{copy.contact}: {supportEmail}</a>
    </Button>}
  </div>
}

function BlockedSignOut({ pathname, copy }: { pathname: string; copy: Copy }) {
  const logout = useAuthStore(state => state.logout)
  const router = useRouter()
  const [pending, setPending] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const signOut = async () => {
    setPending(true)
    setError(null)
    try {
      await api.post('/auth/logout')
      logout()
      router.replace(getDashboardLoginUrl(pathname))
    } catch {
      setError(copy.logoutError)
    } finally {
      setPending(false)
    }
  }
  return <>
    {error && <p role="alert" className="text-sm text-center text-red-700 dark:text-red-300">{error}</p>}
    <Button onClick={signOut} disabled={pending} aria-busy={pending} variant="outline" className="h-12 w-full">
      {pending ? copy.signingOut : copy.signOut}
    </Button>
  </>
}

export function BlockedLogout() {
  const pathname = usePathname()
  return <BlockedSignOut pathname={pathname} copy={messages[getLocaleFromPath(pathname)]} />
}

export function BlockedOverlay({ blockReason, supportEmail }: BlockedOverlayProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const copy = messages[locale]
  const content = blockContent(blockReason, copy)
  return <div dir={locale === 'ar' ? 'rtl' : 'ltr'} className="flex items-center justify-center min-h-[calc(100vh-4rem)] bg-gray-50 dark:bg-gray-900">
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-2xl max-w-md w-full mx-4 p-8">
      <div className="flex justify-center mb-6"><content.Icon className="h-16 w-16 text-orange-600 dark:text-orange-300" aria-hidden="true" /></div>
      <h1 className="text-2xl font-bold text-center text-gray-900 dark:text-white mb-4">{content.title}</h1>
      <p className="text-center text-gray-600 dark:text-gray-300 mb-8">{content.description}</p>
      <BlockedActions copy={copy} billingHref={`/${locale}/dashboard/settings/billing`} hasBilling={content.hasBilling} supportEmail={supportEmail} />
      <div className="border-t border-gray-200 dark:border-gray-700 mb-4" />
      <BlockedSignOut pathname={pathname} copy={copy} />
    </div>
  </div>
}
