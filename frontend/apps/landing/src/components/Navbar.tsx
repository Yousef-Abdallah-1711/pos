'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, ArrowRight, Globe2, Menu, X } from 'lucide-react'
import { useLanguage } from '@/contexts/LanguageContext'
import { BrandMark } from '@/components/BrandMark'

export function Navbar() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const { language, setLanguage, t, dir } = useLanguage()
  const locale = language
  const otherLocale = locale === 'en' ? 'ar' : 'en'
  const DirectionArrow = dir === 'rtl' ? ArrowLeft : ArrowRight
  const links = [
    { href: '#features', label: t.nav.features },
    { href: '#how-it-works', label: t.nav.howItWorks },
    { href: '#pricing', label: t.nav.pricing },
    { href: '#faq', label: t.nav.faq },
  ]

  const closeMenu = () => setMobileMenuOpen(false)

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
          <Link className="poos-header__login" href={`/${locale}/auth/login`}>{t.nav.login}</Link>
          <Link className="poos-button poos-header__start" href={`/${locale}/auth/register`}>
            {t.nav.getStarted}
            <DirectionArrow aria-hidden="true" size={17} />
          </Link>
          <button
            type="button"
            className="poos-menu-toggle"
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
        <div id="poos-mobile-navigation" className="poos-mobile-nav" dir={dir}>
          <div className="poos-shell">
            {links.map((link) => (
              <a key={link.href} href={link.href} onClick={closeMenu}>{link.label}</a>
            ))}
            <div className="poos-mobile-nav__actions">
              <Link className="poos-button poos-button-secondary" href={`/${locale}/auth/login`} onClick={closeMenu}>{t.nav.login}</Link>
              <Link className="poos-button" href={`/${locale}/auth/register`} onClick={closeMenu}>{t.nav.getStarted}</Link>
            </div>
          </div>
        </div>
      )}
    </header>
  )
}
