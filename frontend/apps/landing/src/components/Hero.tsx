'use client'

import Image from 'next/image'
import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { useLanguage } from '@/contexts/LanguageContext'

export function Hero() {
  const { t, language, dir } = useLanguage()
  const Arrow = dir === 'rtl' ? ArrowLeft : ArrowRight
  const locale = language

  return (
    <section className="landing-hero" dir={dir}>
      <div className="poos-shell">
        <div className="ledger-spread">
          <div className="ledger-spread__copy">
            <h1>{t.hero.title} <span className="poos-primary">{t.hero.titleHighlight}</span></h1>
            <p>{t.hero.description}</p>
            <div className="ledger-spread__actions">
              <Link className="poos-button" href={`/${locale}/auth/register`}>
                {t.hero.cta}<Arrow aria-hidden="true" size={18} />
              </Link>
              <a className="poos-button poos-button-secondary" href="#features">{t.hero.demo}</a>
            </div>
          </div>

          <figure className="ledger-proof">
            <div className="ledger-proof__heading">
              <span>{t.hero.captureLabel}</span>
              <span dir="ltr">POOS / {t.hero.captureSection}</span>
            </div>
            <div className="ledger-proof__body">
              <Image
                src="/images/poos-access-control-capture.png"
                alt={t.hero.captureAlt}
                width={1279}
                height={520}
                priority
                sizes="(max-width: 768px) 100vw, 56vw"
              />
            </div>
            <figcaption className="ledger-proof__caption">{t.hero.captureCaption}</figcaption>
          </figure>
        </div>

        <nav className="ledger-index" aria-label={t.hero.productAreas}>
          {t.hero.areas.map((area: { label: string; href: string }, index: number) => (
            <a href={area.href} key={area.href}>
              <span className="poos-primary" aria-hidden="true">{String(index + 1).padStart(2, '0')}</span>
              <span>{area.label}</span>
            </a>
          ))}
        </nav>
      </div>
    </section>
  )
}
