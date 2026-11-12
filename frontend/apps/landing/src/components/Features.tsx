'use client'

import { useLanguage } from '@/contexts/LanguageContext'

export function Features() {
  const { t, dir } = useLanguage()

  return (
    <>
      <section id="product" className="product-ledger" dir={dir}>
        <div className="poos-shell">
          <header className="product-ledger__heading">
            <h2>{t.features.title}</h2>
            <p>{t.features.subtitle}</p>
          </header>
          <div className="product-ledger__spread">
            <article className="product-ledger__entry product-ledger__entry--front">
              <p>{t.features.frontLabel}</p>
              <h3>{t.features.frontTitle}</h3>
              <span>{t.features.frontNote}</span>
            </article>
            <div className="product-ledger__join" aria-label={t.features.connectionLabel}>
              <span className="product-ledger__join-mark" aria-hidden="true" />
              <strong>POOS</strong>
              <span>{t.features.connectionLabel}</span>
            </div>
            <article className="product-ledger__entry product-ledger__entry--back">
              <p>{t.features.backLabel}</p>
              <h3>{t.features.backTitle}</h3>
              <span>{t.features.backNote}</span>
            </article>
          </div>
        </div>
      </section>

      <section id="audience" className="audience-ledger" dir={dir}>
        <div className="poos-shell audience-ledger__inner">
          <h2>{t.audience.title}</h2>
          <div className="audience-ledger__entries">
            <p>{t.audience.independent}</p>
            <p>{t.audience.groups}</p>
          </div>
        </div>
      </section>
    </>
  )
}
