'use client'

import { useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { ArrowLeft, ArrowRight, Globe2, Menu, X } from 'lucide-react'
import { useLanguage } from '@/contexts/LanguageContext'
import { BrandMark } from '@/components/BrandMark'

export function Navbar() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const { language, setLanguage, t, dir } = useLanguage()
  const pathname = usePathname() ?? ''
  const locale = language
  const otherLocale = locale === 'en' ? 'ar' : 'en'
  const DirectionArrow = dir === 'rtl' ? ArrowLeft : ArrowRight
  const isLoginPage = pathname.endsWith('/auth/login')
  const isRegisterPage = pathname.endsWith('/auth/register')
  const isLandingPage = pathname === '/en' || pathname === '/ar'
  const menuButtonRef = useRef<HTMLButtonElement>(null)
  const mobileNavRef = useRef<HTMLDivElement>(null)
  const sectionHref = (sectionId: string) => isLandingPage ? `#${sectionId}` : `/${locale}#${sectionId}`
  const links = [
    { href: sectionHref('product'), label: t.nav.features },
    { href: sectionHref('how-it-works'), label: t.nav.howItWorks },
    { href: sectionHref('audience'), label: t.nav.audience },
    { href: sectionHref('faq'), label: t.nav.faq },
  ]

  const closeMenu = () => setMobileMenuOpen(false)

  useEffect(() => {
    if (!mobileMenuOpen) return

    const previousOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    mobileNavRef.current?.querySelector<HTMLAnchorElement>('a')?.focus()

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key !== 'Escape') return
      setMobileMenuOpen(false)
      menuButtonRef.current?.focus()
    }

    document.addEventListener('keydown', handleKeyDown)
    return () => {
      document.body.style.overflow = previousOverflow
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [mobileMenuOpen])

  return (
    <header className="poos-header" dir={dir}>
      <nav className="poos-header__inner poos-shell" aria-label={t.nav.mainNavigation}>
        <BrandMark href={`/${locale}`} label="POOS" />

        <div className="poos-header__links">
          {links.map((link) => (
            <a key={link.href} href={link.href} onClick={closeMenu}>{link.label}</a>
          ))}
        </div>

        <div className="poos-header__actions">
          <button
            type="button"
            className="poos-language"
            onClick={() => setLanguage(otherLocale)}
            aria-label={t.nav.switchLanguage}
          >
            <Globe2 aria-hidden="true" size={18} />
            <span>{otherLocale === 'ar' ? 'العربية' : 'English'}</span>
          </button>
          {!isLoginPage && <Link className="poos-header__login" href={`/${locale}/auth/login`}>{t.nav.login}</Link>}
          {!isRegisterPage && (
            <Link className="poos-button poos-header__start" href={`/${locale}/auth/register`}>
              {t.nav.getStarted}
              <DirectionArrow aria-hidden="true" size={17} />
            </Link>
          )}
          <button
            type="button"
            className="poos-menu-toggle"
            ref={menuButtonRef}
            onClick={() => setMobileMenuOpen((open) => !open)}
            aria-label={mobileMenuOpen ? t.nav.closeMenu : t.nav.openMenu}
            aria-expanded={mobileMenuOpen}
            aria-controls="poos-mobile-navigation"
          >
            {mobileMenuOpen ? <X aria-hidden="true" size={23} /> : <Menu aria-hidden="true" size={23} />}
          </button>
        </div>
      </nav>

      {mobileMenuOpen && (
        <div id="poos-mobile-navigation" ref={mobileNavRef} className="poos-mobile-nav" dir={dir}>
          <div className="poos-shell">
            {links.map((link) => (
              <a key={link.href} href={link.href} onClick={closeMenu}>{link.label}</a>
            ))}
            <div className="poos-mobile-nav__actions">
              <button type="button" className="poos-language" onClick={() => { setLanguage(otherLocale); closeMenu() }} aria-label={t.nav.switchLanguage}>
                <Globe2 aria-hidden="true" size={18} />
                <span>{otherLocale === 'ar' ? 'العربية' : 'English'}</span>
              </button>
              {!isLoginPage && <Link className="poos-button poos-button-secondary" href={`/${locale}/auth/login`} onClick={closeMenu}>{t.nav.login}</Link>}
              {!isRegisterPage && <Link className="poos-button" href={`/${locale}/auth/register`} onClick={closeMenu}>{t.nav.getStarted}</Link>}
            </div>
          </div>
        </div>
      )}
    </header>
  )
}
