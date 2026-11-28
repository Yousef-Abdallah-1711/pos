'use client'

import { useEffect, useRef, useState } from 'react'
import { ArrowDownRight, RotateCcw } from 'lucide-react'
import type { Locale } from '@/content/site'
import { getConversionCopy } from '@/content/conversionCopy'
import { ProductCapture } from './ProductCapture'
import { startDesktopMotion } from './desktopMotion'

const images = ['pos', 'dashboard', 'kitchen'] as const

export function OperationalBento({ locale }: { locale: Locale }) {
  const root = useRef<HTMLElement>(null)
  const [steps, setSteps] = useState([0, 0, 0])
  const c = getConversionCopy(locale)

  useEffect(() => startDesktopMotion(({ gsap, ScrollTrigger }) => {
    if (!root.current) return
    const context = gsap.context(() => {
      const cards = gsap.utils.toArray<HTMLElement>('.cv-bento-card')
      cards.forEach((card, index) => {
        const story = gsap.timeline({ paused: true })
          .fromTo(card, { y: 22, opacity: .82 }, { y: 0, opacity: 1, duration: .55, ease: 'power2.out' })
          .call(() => setSteps(current => current.map((step, i) => i === index ? 1 : step)), [], '+=.18')
          .call(() => setSteps(current => current.map((step, i) => i === index ? 2 : step)), [], '+=1.1')
        ScrollTrigger.create({ trigger: card, start: 'top 82%', once: true, onEnter: () => story.play() })
      })
    }, root.current)
    return () => context.revert()
  }), [])

  const setStep = (card: number, step: number) => setSteps(current => current.map((value, index) => index === card ? step : value))

  return <section className="cv-section cv-problem" ref={root} aria-labelledby="cv-problem-title">
    <div className="cv-wrap">
      <div className="cv-bento-intro">
        <div><p className="cv-eyebrow">{c.problemEyebrow}</p><h2 id="cv-problem-title">{c.problemTitle}</h2></div>
        <p>{c.problemBody}</p>
      </div>
      <div className="cv-bento-grid">
        {c.problems.map((item, index) => <article className={`cv-bento-card cv-bento-card--${index + 1}`} key={item.title}>
          <div className="cv-bento-card-head">
            <span className="cv-index">0{index + 1} / 03</span>
            <span className="cv-bento-card-kind">{item.kind}</span>
          </div>
          <div className="cv-bento-card-copy"><h3>{item.title}</h3><p>{item.body}</p></div>
          <div className="cv-bento-media">
            <ProductCapture image={images[index]} alt={item.alt} locale={locale} />
          </div>
          <div className="cv-bento-sample" data-step={steps[index]}>
            <span>{c.problemSample}</span>
            <strong>{c.problemSampleDetails[index][steps[index]]}</strong>
            <span className="cv-bento-sample-indicator" aria-hidden="true">{steps[index] === 2 ? '✓' : `0${steps[index] + 1}`}</span>
            {index === 0 && <div className="cv-bento-pos-example" aria-label={c.problemSample}>
              <div><small>{c.problemPos.menu}</small><span className="cv-bento-pos-item" data-selected={steps[index] > 0}>{c.problemPos.item}</span></div>
              <div><small>{c.problemPos.cart}</small><span>{steps[index] === 0 ? c.problemPos.empty : `${c.problemPos.item} · ${c.problemPos.quantity}: 1`}</span>
                {steps[index] > 0 && <span className="cv-bento-pos-submit" data-highlight={steps[index] === 2}>{c.problemPos.place} (1)</span>}
              </div>
            </div>}
          </div>
          <div className="cv-bento-story">
            <div className="cv-bento-story-state">
              <span>{c.problemIllustration} · {String(steps[index] + 1).padStart(2, '0')} / 03</span>
              <strong>{item.stages[steps[index]]}</strong>
            </div>
            <div className="cv-bento-story-controls" role="group" aria-label={item.title}>
              {item.stages.map((stage, stageIndex) => <button type="button" key={stage} aria-label={`${item.title}: ${stage}`} aria-pressed={steps[index] === stageIndex} onClick={() => setStep(index, stageIndex)} />)}
              <button type="button" className="cv-bento-replay" aria-label={c.problemReplay} onClick={() => setStep(index, 0)}><RotateCcw size={15} aria-hidden="true" /></button>
            </div>
          </div>
        </article>)}
      </div>
      <a className="cv-bento-bridge" href="#how-orders-move"><span>{c.problemBridge}</span><ArrowDownRight size={19} aria-hidden="true" /></a>
    </div>
  </section>
}
