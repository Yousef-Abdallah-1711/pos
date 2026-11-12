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
            <h1>{t.hero.title} <span>{t.hero.titleHighlight}</span></h1>
            <p>{t.hero.description}</p>
            <div className="ledger-spread__actions">
              <Link className="poos-button" href={`/${locale}/auth/register`}>
                {t.hero.cta}<Arrow aria-hidden="true" size={18} />
              </Link>
              <a className="poos-button poos-button-secondary" href="#how-it-works">{t.hero.demo}</a>
            </div>
          </div>

          <figure className="ledger-proof">
            <div className="ledger-proof__heading">
              <span className="ledger-proof__chrome-dots" aria-hidden="true"><i /><i /><i /></span>
              <span>{t.hero.captureLabel}</span>
              <span dir="ltr">POOS / {t.hero.captureSection}</span>
            </div>
            <div className="ledger-proof__body">
              <picture>
                <source
                  media="(max-width: 40rem)"
                  srcSet="/images/product/roles-mobile-430x700.png"
                  width="344"
                  height="560"
                />
                <Image
                  src="/images/product/roles-1440x900.png"
                  alt={t.hero.captureAlt}
                  width={1440}
                  height={900}
                  priority
                  loading="eager"
                  fetchPriority="high"
                  sizes="(max-width: 640px) 100vw, (max-width: 1200px) 52vw, 680px"
                />
              </picture>
            </div>
            <figcaption className="ledger-proof__caption">
              <span>{t.hero.captureCaption}</span>
              <span className="ledger-proof__caption-actions">
                <span className="ledger-proof__sample">{t.hero.sampleLabel}</span>
                <a href="/images/product/roles-1440x900.png" target="_blank" rel="noreferrer">
                  {t.hero.viewCapture}
                </a>
              </span>
            </figcaption>
          </figure>
        </div>

        <nav className="ledger-index" aria-label={t.hero.productAreas}>
          {t.hero.areas.map((area: { label: string; href: string }, index: number) => (
            <a href={area.href} key={`${area.href}:${area.label}`}>
              <span className="poos-primary" aria-hidden="true">{String(index + 1).padStart(2, '0')}</span>
              <span>{area.label}</span>
            </a>
          ))}
        </nav>
      </div>
    </section>
  )
}
