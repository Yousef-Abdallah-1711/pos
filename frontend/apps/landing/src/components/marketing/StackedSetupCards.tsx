'use client'

import { useRef } from 'react'
import { useGSAP } from '@gsap/react'
import gsap from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import type { Locale } from '@/content/site'
import { getConversionCopy } from '@/content/conversionCopy'
import { TrialLink } from './TrialLink'

gsap.registerPlugin(useGSAP, ScrollTrigger)

export function StackedSetupCards({ locale }: { locale: Locale }) {
  const root = useRef<HTMLElement>(null)
  const c = getConversionCopy(locale)

  useGSAP(() => {
    if (!root.current) return
    const media = gsap.matchMedia()
    media.add('(min-width: 901px) and (prefers-reduced-motion: no-preference)', () => {
      const cards = gsap.utils.toArray<HTMLElement>('.cv-setup-card')
      const triggers = cards.slice(0, -1).map((card, index) => ScrollTrigger.create({
        trigger: cards[index + 1], start: 'top center', end: 'top 20%', scrub: .5,
        onUpdate: self => gsap.set(card, { scale: 1 - self.progress * .035, opacity: 1 - self.progress * .15 }),
      }))
      return () => triggers.forEach(trigger => trigger.kill())
    })
    return () => media.revert()
  }, { scope: root, dependencies: [locale] })

  return <section id="getting-started" className="cv-section cv-setup cv-setup-stacked" aria-labelledby="cv-setup-title" ref={root}>
    <div className="cv-wrap">
      <p className="cv-eyebrow">{c.setupEyebrow}</p><h2 id="cv-setup-title">{c.setupTitle}</h2>
      <ol className="cv-setup-card-list">
        {c.setupSteps.map((step, index) => <li className="cv-setup-card" key={step.title}>
          <span className="cv-index">{String(index + 1).padStart(2, '0')} / {String(c.setupSteps.length).padStart(2, '0')}</span>
          <div><h3>{step.title}</h3><p>{step.body}</p></div>
          <span className="cv-setup-card-mark" aria-hidden="true">SERVIO</span>
        </li>)}
      </ol>
      <TrialLink locale={locale} source="home-setup" label={c.createAccount} />
    </div>
  </section>
}
