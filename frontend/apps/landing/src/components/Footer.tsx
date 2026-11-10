'use client'

import Link from 'next/link'
import { useLanguage } from '@/contexts/LanguageContext'
import { BrandMark } from '@/components/BrandMark'

export function Footer() {
  const { t, language, dir } = useLanguage()
  const locale = language
  const year = new Date().getFullYear()

  return (
    <footer className="poos-footer" dir={dir}>
      <div className="poos-shell poos-footer__inner">
        <div>
          <BrandMark href={`/${locale}`} label="POOS" inverse />
          <p className="poos-muted">{t.footer.tagline}</p>
        </div>
        <nav className="poos-footer__links" aria-label={t.nav.mainNavigation}>
          <Link href="#product">{t.footer.features}</Link>
          <Link href="#how-it-works">{t.nav.howItWorks}</Link>
          <Link href="#audience">{t.nav.audience}</Link>
          <Link href="#faq">{t.nav.faq}</Link>
          <Link href={`/${locale}/auth/login`}>{t.nav.login}</Link>
          <Link href={`/${locale}/auth/register`}>{t.nav.getStarted}</Link>
        </nav>
        <small className="poos-muted">© {year} POOS</small>
      </div>
    </footer>
  )
}
