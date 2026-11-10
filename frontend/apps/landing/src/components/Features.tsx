'use client'

import { useLanguage } from '@/contexts/LanguageContext'

const featureKeys = ['pos', 'kitchen', 'inventory', 'team', 'online'] as const

export function Features() {
  const { t, dir } = useLanguage()

  return (
    <section id="features" className="ledger-section" dir={dir}>
      <div className="poos-shell">
        <header className="ledger-section__head">
          <h2>{t.features.title}</h2>
          <p>{t.features.subtitle}</p>
        </header>
        <div className="ledger-entries">
          {featureKeys.map((key) => (
            <article className="ledger-entry" id={t.features[key].id} key={key}>
              <h3>{t.features[key].title}</h3>
              <p>{t.features[key].description}</p>
              <p className="ledger-entry__note">{t.features[key].note}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  )
}
