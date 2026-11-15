'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import type { MouseEvent } from 'react'
import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { Globe2, Menu, Moon, Sun, X } from 'lucide-react'
import { useTheme } from 'next-themes'
import { getSiteCopy, type Locale } from '@/content/site'
import { HEADER_LINKS, isNavigationRouteActive, localizedHref, SITE_ROUTES } from '@/content/navigation'
import { BrandMark } from './BrandMark'
import { TrialLink } from './marketing/TrialLink'
import { DesktopMegaMenu } from './navigation/DesktopMegaMenu'
import { MobileNavigation } from './navigation/MobileNavigation'

export function Navbar({ locale }: { locale: Locale }) {
  const { resolvedTheme, setTheme } = useTheme()
  const pathname = usePathname()
  const router = useRouter()
  const [mounted, setMounted] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)
  const triggerRef = useRef<HTMLButtonElement>(null)
  const copy = getSiteCopy(locale)
  const theme = mounted ? resolvedTheme : 'dark'
  const nextLocale = locale === 'en' ? 'ar' : 'en'
  const localizedPath = pathname.replace(/^\/(?:en|ar)(?=\/|$)/, `/${nextLocale}`)

  useEffect(() => setMounted(true), [])
  useEffect(() => {
    const desktop = window.matchMedia('(min-width: 1200px)')
    const closeOnDesktop = () => {
      if (desktop.matches) setMobileOpen(false)
    }
    closeOnDesktop()
    desktop.addEventListener('change', closeOnDesktop)
    return () => desktop.removeEventListener('change', closeOnDesktop)
  }, [])

  const closeMobile = useCallback(() => setMobileOpen(false), [])
  const changeLocale = useCallback(() => {
    closeMobile()
    router.push(`${localizedPath}${window.location.search}${window.location.hash}`)
  }, [closeMobile, localizedPath, router])
  const switchLocale = (event: MouseEvent<HTMLAnchorElement>) => {
    event.preventDefault()
    changeLocale()
  }

  return <header className="sv-header">
    <div className="sv-shell sv-nav">
      <BrandMark href={localizedHref(locale, 'home')} label="SERVIO" />
      <nav className="sv-desktop-links" aria-label={copy.navLabel}>
        <DesktopMegaMenu locale={locale} pathname={pathname} />
        {HEADER_LINKS.filter(item => item.route !== 'features').map(item => <Link
          key={item.id}
          href={localizedHref(locale, item.route)}
          aria-current={isNavigationRouteActive(pathname, SITE_ROUTES[item.route]) ? 'page' : undefined}>
          {copy.navigation[item.label]}
        </Link>)}
      </nav>
      <div className="sv-nav-actions">
        <Link className="sv-language" href={localizedPath} aria-label={locale === 'en' ? copy.switchToArabic : copy.switchToEnglish} onClick={switchLocale}>
          <Globe2 size={18} aria-hidden="true" /><span>{locale === 'en' ? 'العربية' : 'English'}</span>
        </Link>
        <button className="sv-theme-toggle" type="button" onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
          aria-label={theme === 'dark' ? copy.switchToLight : copy.switchToDark} title={theme === 'dark' ? copy.switchToLight : copy.switchToDark}>
          {theme === 'dark' ? <Sun size={18} aria-hidden="true" /> : <Moon size={18} aria-hidden="true" />}
        </button>
        <Link className="sv-login" href={localizedHref(locale, 'login')} prefetch={false}>{copy.navigation.login}</Link>
        <div className="sv-desktop-trial"><TrialLink locale={locale} source="navigation" /></div>
        <button className="sv-menu-button" ref={triggerRef} type="button" aria-controls="sv-mobile-navigation" aria-expanded={mobileOpen}
          aria-label={mobileOpen ? copy.closeMenu : copy.openMenu} onClick={() => setMobileOpen(value => !value)}>
          {mobileOpen ? <X aria-hidden="true" /> : <Menu aria-hidden="true" />}
        </button>
      </div>
    </div>
    <MobileNavigation
      locale={locale}
      pathname={pathname}
      open={mobileOpen}
      onClose={closeMobile}
      triggerRef={triggerRef}
      theme={theme ?? 'dark'}
      onToggleTheme={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
      onSwitchLocale={changeLocale}
    />
  </header>
}
