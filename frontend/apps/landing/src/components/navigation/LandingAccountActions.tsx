'use client'

import { useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { ChevronDown, LogOut } from 'lucide-react'
import { getConfig } from '@/lib/config'
import type { LandingSessionProfile, LandingSessionState } from '@/lib/auth-session'
import { getSiteCopy, type Locale } from '@/content/site'
import { TrialLink } from '@/components/marketing/TrialLink'
import { AccountEntryLinks } from './AccountEntryLinks'

interface LandingAccountActionsProps {
  locale: Locale
  session: LandingSessionState
  isSigningOut: boolean
  onRetry: () => void
  onSignOut: () => Promise<boolean>
}

function getInitials(name: string) {
  return name.trim().split(/\s+/).slice(0, 2).map(part => Array.from(part)[0] ?? '').join('').toLocaleUpperCase()
}

export function LandingAccountActions({ locale, session, isSigningOut, onRetry, onSignOut }: LandingAccountActionsProps) {
  const copy = getSiteCopy(locale)
  const [openForProfile, setOpenForProfile] = useState<LandingSessionProfile | null>(null)
  const rootRef = useRef<HTMLDivElement>(null)
  const triggerRef = useRef<HTMLButtonElement>(null)

  const isOpen = session.status === 'authenticated' && openForProfile === session.profile

  useEffect(() => {
    if (!isOpen) return
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key !== 'Escape') return
      event.preventDefault()
      setOpenForProfile(null)
      triggerRef.current?.focus()
    }
    const closeOutside = (event: MouseEvent) => {
      if (event.target instanceof Node && !rootRef.current?.contains(event.target)) setOpenForProfile(null)
    }
    document.addEventListener('keydown', closeOnEscape)
    document.addEventListener('mousedown', closeOutside)
    return () => {
      document.removeEventListener('keydown', closeOnEscape)
      document.removeEventListener('mousedown', closeOutside)
    }
  }, [isOpen])

  if (session.status === 'checking') {
    return <div className="sv-account-status" role="status" aria-live="polite" aria-busy="true">{copy.account.checking}</div>
  }

  if (session.status === 'temporarily-unavailable') {
    return <>
      <AccountEntryLinks locale={locale} />
      <div className="sv-desktop-trial"><TrialLink locale={locale} source="navigation" /></div>
      <div className="sv-account-unavailable" role="status">
        <span>{session.reason === 'sign-out' ? copy.account.signOutFailed : copy.account.unavailable}</span>
        <button className="sv-account-retry" type="button" onClick={onRetry}>{copy.account.retry}</button>
      </div>
    </>
  }

  if (session.status === 'guest') {
    return <>
      <AccountEntryLinks locale={locale} />
      <div className="sv-desktop-trial"><TrialLink locale={locale} source="navigation" /></div>
    </>
  }

  const { profile } = session
  const dashboardUrl = getConfig().dashboardUrl
  const dashboardPath = `${dashboardUrl}/${locale}/dashboard`
  const settingsPath = `${dashboardUrl}/${locale}/dashboard/settings/account`

  return <div className="sv-account-desktop" ref={rootRef}>
    <button
      className="sv-account-trigger"
      ref={triggerRef}
      type="button"
      aria-expanded={isOpen}
      aria-controls="sv-account-panel"
      aria-label={`${copy.account.menuLabelPrefix} ${profile.name}`}
      onClick={() => setOpenForProfile(isOpen ? null : profile)}>
      <span className="sv-account-avatar" aria-hidden="true">{getInitials(profile.name)}</span>
      <span className="sv-account-trigger__name"><bdi dir="auto">{profile.name}</bdi></span>
      <ChevronDown size={16} aria-hidden="true" />
    </button>
    <div className="sv-account-panel" id="sv-account-panel" hidden={!isOpen}>
      <div className="sv-account-panel__identity">
        <strong><bdi dir="auto">{profile.name}</bdi></strong>
        <span><bdi dir="auto">{profile.email}</bdi></span>
        {profile.restaurantName && <span className="sv-account-panel__restaurant"><small>{copy.account.restaurantContext}: </small><bdi dir="auto">{profile.restaurantName}</bdi></span>}
      </div>
      <nav aria-label={copy.account.actionsLabel}>
        <Link href={dashboardPath} prefetch={false} onClick={() => setOpenForProfile(null)}>{copy.account.dashboard}</Link>
        <Link href={settingsPath} prefetch={false} onClick={() => setOpenForProfile(null)}>{copy.account.settings}</Link>
      </nav>
      <button className="sv-account-panel__signout" type="button" disabled={isSigningOut} onClick={() => { void onSignOut() }}>
        <LogOut size={16} aria-hidden="true" />{isSigningOut ? copy.account.signingOut : copy.account.signOut}
      </button>
    </div>
  </div>
}
