'use client'

import { useRef, useState } from 'react'
import Link from 'next/link'
import { ArrowUpRight, ChevronLeft, ChevronRight } from 'lucide-react'
import { useGSAP } from '@gsap/react'
import gsap from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import type { Locale } from '@/content/site'
import { getConversionCopy } from '@/content/conversionCopy'
import { getFeatures } from '@/content/features'
import { ProductCapture } from './ProductCapture'

gsap.registerPlugin(useGSAP, ScrollTrigger)

const controls = {
  en: { previous: 'Previous feature', next: 'Next feature', explore: 'Explore feature', sample: 'Actual product screen · synthetic data' },
  ar: { previous: 'الميزة السابقة', next: 'الميزة التالية', explore: 'اكتشف الميزة', sample: 'شاشة فعلية من المنتج · بيانات تجريبية' },
} as const

export function FeatureShowcase({ locale }: { locale: Locale }) {
  const root = useRef<HTMLElement>(null)
  const viewport = useRef<HTMLDivElement>(null)
  const track = useRef<HTMLDivElement>(null)
  const scrollTrigger = useRef<ScrollTrigger | null>(null)
  const [active, setActive] = useState(0)
  const features = getFeatures(locale)
  const c = getConversionCopy(locale)
  const ui = controls[locale]

  useGSAP(() => {
    if (!root.current || !viewport.current || !track.current) return
    const media = gsap.matchMedia()
    media.add('(min-width: 901px) and (prefers-reduced-motion: no-preference)', () => {
      const distance = () => Math.max(0, track.current!.scrollWidth - viewport.current!.clientWidth)
      const tween = gsap.to(track.current, {
        x: () => (locale === 'ar' ? 1 : -1) * distance(),
        ease: 'none',
        scrollTrigger: {
          trigger: root.current,
          start: 'top top+=78',
          end: () => `+=${distance()}`,
          pin: true,
          scrub: .8,
          invalidateOnRefresh: true,
          onUpdate: self => setActive(Math.min(features.length - 1, Math.round(self.progress * (features.length - 1)))),
        },
      })
      scrollTrigger.current = tween.scrollTrigger ?? null
      return () => { scrollTrigger.current = null; tween.kill() }
    })
    media.add('(max-width: 900px), (prefers-reduced-motion: reduce)', () => {
      gsap.set(track.current, { clearProps: 'transform' })
    })
    return () => media.revert()
  }, { scope: root, dependencies: [locale] })

  const goTo = (index: number) => {
    const next = Math.max(0, Math.min(features.length - 1, index))
    setActive(next)
    if (scrollTrigger.current) {
      const trigger = scrollTrigger.current
      window.scrollTo({ top: trigger.start + (trigger.end - trigger.start) * next / (features.length - 1), behavior: 'smooth' })
      return
    }
    track.current?.children[next]?.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'start' })
  }

  return <section id="capabilities" className="cv-section cv-showcase" ref={root} aria-labelledby="cv-capabilities-title">
    <div className="cv-wrap cv-showcase-heading">
      <div><p className="cv-eyebrow">{c.capabilitiesEyebrow}</p><h2 id="cv-capabilities-title">{c.capabilitiesTitle}</h2></div>
      <div><p>{c.capabilitiesBody}</p><div className="cv-showcase-controls">
        <button type="button" onClick={() => goTo(active - 1)} disabled={active === 0} aria-label={ui.previous}><ChevronLeft aria-hidden="true" /></button>
        <span aria-live="polite">{String(active + 1).padStart(2, '0')} / {String(features.length).padStart(2, '0')}</span>
        <button type="button" onClick={() => goTo(active + 1)} disabled={active === features.length - 1} aria-label={ui.next}><ChevronRight aria-hidden="true" /></button>
      </div></div>
    </div>
    <div className="cv-showcase-viewport" ref={viewport}>
      <div className="cv-showcase-track" ref={track}>
        {features.map((feature, index) => <article className="cv-feature-slide" key={feature.slug} id={`feature-slide-${feature.slug}`}>
          <div className="cv-feature-slide-copy">
            <p className="cv-index">{String(index + 1).padStart(2, '0')} / {String(features.length).padStart(2, '0')}</p>
            <p className="cv-feature-category">{feature.eyebrow}</p>
            <h3>{feature.headline}</h3>
            <p>{feature.description}</p>
            <Link href={`/${locale}/features/${feature.slug}`}>{ui.explore}<ArrowUpRight size={18} aria-hidden="true" /></Link>
          </div>
          <div className="cv-feature-slide-capture"><ProductCapture image={feature.image} alt={feature.title} locale={locale} /><small>{ui.sample}</small></div>
        </article>)}
      </div>
    </div>
  </section>
}
