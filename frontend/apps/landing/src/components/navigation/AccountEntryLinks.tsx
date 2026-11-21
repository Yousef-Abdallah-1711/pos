'use client'

import Link from 'next/link'
import type { Locale } from '@/content/site'
import { getSiteCopy } from '@/content/site'
import { localizedHref } from '@/content/navigation'

export function AccountEntryLinks({ locale, mobile = false, onNavigate }: { locale: Locale; mobile?: boolean; onNavigate?: () => void }) {
  const copy = getSiteCopy(locale)
  const className = mobile ? 'sv-mobile-navigation__link' : undefined

  return <Link className={mobile ? className : 'sv-login'} href={localizedHref(locale, 'login')} prefetch={false} onClick={onNavigate}>{copy.navigation.login}</Link>
}
