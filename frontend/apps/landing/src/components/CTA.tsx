'use client'

import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { useLanguage } from '@/contexts/LanguageContext'

export function CTA() {
  const { t, language, dir } = useLanguage()
  const Arrow = dir === 'rtl' ? ArrowLeft : ArrowRight

  return (
    <section id="contact" className="ledger-cta" dir={dir}>
      <div className="poos-shell ledger-cta__inner">
        <div>
          <h2>{t.cta.title}</h2>
          <p className="poos-muted">{t.cta.subtitle}</p>
        </div>
        <Link className="poos-button" href={`/${language}/auth/register`}>
          {t.cta.button}<Arrow aria-hidden="true" size={18} />
        </Link>
      </div>
    </section>
  )
}
