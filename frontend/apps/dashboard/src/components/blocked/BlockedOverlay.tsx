'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useTranslations } from 'next-intl'
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

interface BlockedCopy {
  trialTitle: string
  expiredTitle: string
  expiredDescription: string
  disabledTitle: string
  disabledDescription: string
  restrictedTitle: string
  restrictedDescription: string
  billing: string
  billingNotice: string
  supportHelp: string
  contactSupport: string
  signOut: string
  signingOut: string
  signOutError: string
}

function blockContent(reason: string, copy: BlockedCopy) {
  if (reason === 'trial_expired') {
    return { Icon: Clock, title: copy.trialTitle, description: copy.expiredDescription, hasBilling: true }
  }
  if (reason === 'expired') {
    return { Icon: Clock, title: copy.expiredTitle, description: copy.expiredDescription, hasBilling: true }
  }
  if (reason === 'tenant_disabled') {
    return { Icon: Lock, title: copy.disabledTitle, description: copy.disabledDescription, hasBilling: false }
  }
  return { Icon: Lock, title: copy.restrictedTitle, description: copy.restrictedDescription, hasBilling: false }
}

function safeSupportEmail(email?: string): string | undefined {
  const normalized = email?.trim()
  return normalized && normalized.length <= 254 && /^[^\s@<>]+@[^\s@<>]+\.[^\s@<>]+$/.test(normalized)
    ? normalized
    : undefined
}

function BlockedActions({ copy, billingHref, hasBilling, supportEmail }: {
  copy: BlockedCopy
  billingHref: string
  hasBilling: boolean
  supportEmail?: string
}) {
  return <div className="mb-4 space-y-3">
    {hasBilling && <>
      <Button asChild className="h-12 w-full"><Link href={billingHref}>
        <CreditCard className="me-2 h-4 w-4" aria-hidden="true" />{copy.billing}
      </Link></Button>
      <p className="text-center text-xs text-gray-600 dark:text-gray-300">{copy.billingNotice}</p>
    </>}
    <p className="text-center text-sm text-gray-600 dark:text-gray-300">{copy.supportHelp}</p>
    {supportEmail && <Button asChild className="h-12 w-full" variant="outline">
      <a href={`mailto:${encodeURIComponent(supportEmail)}`}>{copy.contactSupport}: {supportEmail}</a>
    </Button>}
  </div>
}

function BlockedSignOut({ pathname, copy }: { pathname: string; copy: BlockedCopy }) {
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
      setError(copy.signOutError)
    } finally {
      setPending(false)
    }
  }
  return <>
    {error && <p role="alert" className="text-center text-sm text-red-700 dark:text-red-300">{error}</p>}
    <Button onClick={signOut} disabled={pending} aria-busy={pending} variant="outline" className="h-12 w-full">
      {pending ? copy.signingOut : copy.signOut}
    </Button>
  </>
}

function useBlockedCopy(): BlockedCopy {
  const t = useTranslations('subscription.blocked')
  return {
    trialTitle: t('trialExpired.title'),
    expiredTitle: t('subscriptionExpired.title'),
    expiredDescription: t('trialExpired.description'),
    disabledTitle: t('tenantDisabled.title'),
    disabledDescription: t('tenantDisabled.description'),
    restrictedTitle: t('restricted.title'),
    restrictedDescription: t('restricted.description'),
    billing: t('actions.billing'),
    billingNotice: t('actions.billingNotice'),
    supportHelp: t('actions.supportHelp'),
    contactSupport: t('actions.contactSupport'),
    signOut: t('actions.signOut'),
    signingOut: t('actions.signingOut'),
    signOutError: t('errors.signOut'),
  }
}

export function BlockedLogout() {
  const pathname = usePathname()
  return <BlockedSignOut pathname={pathname} copy={useBlockedCopy()} />
}

export function BlockedOverlay({ blockReason, supportEmail }: BlockedOverlayProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const copy = useBlockedCopy()
  const content = blockContent(blockReason, copy)
  const support = safeSupportEmail(supportEmail)
  return <div dir={locale === 'ar' ? 'rtl' : 'ltr'} className="flex min-h-[calc(100vh-4rem)] items-center justify-center bg-gray-50 dark:bg-gray-900">
    <div className="mx-4 w-full max-w-md rounded-lg bg-white p-8 shadow-2xl dark:bg-gray-800">
      <div className="mb-6 flex justify-center"><content.Icon className="h-16 w-16 text-orange-600 dark:text-orange-300" aria-hidden="true" /></div>
      <h1 className="mb-4 text-center text-2xl font-bold text-gray-900 dark:text-white">{content.title}</h1>
      <p className="mb-8 text-center text-gray-600 dark:text-gray-300">{content.description}</p>
      <BlockedActions copy={copy} billingHref={`/${locale}/dashboard/settings/billing`} hasBilling={content.hasBilling} supportEmail={support} />
      <div className="mb-4 border-t border-gray-200 dark:border-gray-700" />
      <BlockedSignOut pathname={pathname} copy={copy} />
    </div>
  </div>
}
