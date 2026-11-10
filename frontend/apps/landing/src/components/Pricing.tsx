'use client'

import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { useLanguage } from '@/contexts/LanguageContext'

export function Pricing() {
  const { t, language, dir } = useLanguage()
  const Arrow = dir === 'rtl' ? ArrowLeft : ArrowRight

  return (
    <section id="pricing" className="ledger-section" dir={dir}>
      <div className="poos-shell">
        <div className="pricing-note">
          <div>
            <h2>{t.pricing.title}</h2>
            <p>{t.pricing.subtitle}</p>
          </div>
          <Link className="poos-button" href={`/${language}/auth/register`}>
            {t.pricing.cta}<Arrow aria-hidden="true" size={18} />
          </Link>
        </div>
      </div>
    </section>
  )
}
