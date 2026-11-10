'use client'

import { useLanguage } from '@/contexts/LanguageContext'

export function HowItWorks() {
  const { t, dir } = useLanguage()
  const steps = [t.howItWorks.step1, t.howItWorks.step2, t.howItWorks.step3]

  return (
    <section id="how-it-works" className="ledger-section" dir={dir}>
      <div className="poos-shell">
        <header className="ledger-section__head">
          <h2>{t.howItWorks.title}</h2>
          <p>{t.howItWorks.subtitle}</p>
        </header>
        <ol className="ledger-steps">
          {steps.map((step) => (
            <li className="ledger-step" key={step.number}>
              <span className="ledger-step__number">{step.number}</span>
              <h3>{step.title}</h3>
              <p>{step.description}</p>
            </li>
          ))}
        </ol>
      </div>
    </section>
  )
}
