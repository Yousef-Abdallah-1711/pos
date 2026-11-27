'use client'

import { useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { ArrowUpRight, ChevronLeft, ChevronRight } from 'lucide-react'
import type { ScrollTrigger } from 'gsap/ScrollTrigger'
import type { Locale } from '@/content/site'
import { getConversionCopy } from '@/content/conversionCopy'
import { getFeatures } from '@/content/features'
import { ProductCapture } from './ProductCapture'
import { startDesktopMotion } from './desktopMotion'

const controls = {
  en: { previous: 'Previous feature', next: 'Next feature', explore: 'Explore feature', sample: 'Actual product screen · synthetic data' },
  ar: { previous: 'الميزة السابقة', next: 'الميزة التالية', explore: 'اكتشف الميزة', sample: 'شاشة فعلية من المنتج · بيانات تجريبية' },
} as const

export function FeatureShowcase({ locale }: { locale: Locale }) {
  const stage = useRef<HTMLDivElement>(null)
  const viewport = useRef<HTMLDivElement>(null)
  const track = useRef<HTMLDivElement>(null)
  const scrollTrigger = useRef<ScrollTrigger | null>(null)
  const [active, setActive] = useState(0)
  const features = getFeatures(locale)
  const c = getConversionCopy(locale)
  const ui = controls[locale]

  useEffect(() => startDesktopMotion(({ gsap, ScrollTrigger }) => {
    if (!stage.current || !viewport.current || !track.current) return
    const context = gsap.context(() => {
      const header = document.querySelector<HTMLElement>('.sv-header')
      const measureTop = () => {
        const top = Math.ceil((header?.getBoundingClientRect().height ?? 76) + 12)
        stage.current!.style.setProperty('--cv-showcase-top', `${top}px`)
        return top
      }
      measureTop()
      const distance = () => Math.max(0, track.current!.scrollWidth - viewport.current!.clientWidth)
      const slideOffsets = () => {
        const slides = Array.from(track.current!.children) as HTMLElement[]
        const first = slides[0].getBoundingClientRect()
        return slides.map(slide => {
          const rect = slide.getBoundingClientRect()
          return locale === 'ar' ? first.right - rect.right : rect.left - first.left
        })
      }
      let offsets = slideOffsets()
      let travel = distance()
      if (travel <= 0) return
      const tween = gsap.to(track.current, {
        x: () => (locale === 'ar' ? 1 : -1) * distance(),
        ease: 'none',
        scrollTrigger: {
          trigger: stage.current,
          start: () => `top top+=${measureTop()}`,
          end: () => `+=${Math.max(distance(), 1)}`,
          pin: stage.current,
          scrub: .45,
          invalidateOnRefresh: true,
          onRefresh: () => { travel = distance(); offsets = slideOffsets() },
          onUpdate: self => {
            const position = self.progress * travel
            const nearest = offsets.reduce((best, offset, index) => Math.abs(offset - position) < Math.abs(offsets[best] - position) ? index : best, 0)
            setActive(nearest)
          },
        },
      })
      scrollTrigger.current = tween.scrollTrigger ?? null
    }, stage.current)
    let active = true
    void document.fonts.ready.then(() => { if (active) ScrollTrigger.refresh() })
    return () => {
      active = false
      scrollTrigger.current = null
      context.revert()
      stage.current?.style.removeProperty('--cv-showcase-top')
    }
  }, '(min-width: 901px) and (min-height: 760px) and (prefers-reduced-motion: no-preference)'), [locale, features.length])

  const goTo = (index: number) => {
    const next = Math.max(0, Math.min(features.length - 1, index))
    setActive(next)
    if (scrollTrigger.current) {
      const trigger = scrollTrigger.current
      const slides = Array.from(track.current!.children) as HTMLElement[]
      const first = slides[0].getBoundingClientRect()
      const target = slides[next].getBoundingClientRect()
      const offset = locale === 'ar' ? first.right - target.right : target.left - first.left
      const distance = Math.max(1, track.current!.scrollWidth - viewport.current!.clientWidth)
      window.scrollTo({ top: trigger.start + (trigger.end - trigger.start) * Math.min(1, Math.max(0, offset / distance)), behavior: 'smooth' })
      return
    }
    track.current?.children[next]?.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'start' })
  }

  return <section id="capabilities" className="cv-section cv-showcase" aria-labelledby="cv-capabilities-title">
    <div className="cv-wrap cv-showcase-heading">
      <div><p className="cv-eyebrow">{c.capabilitiesEyebrow}</p><h2 id="cv-capabilities-title">{c.capabilitiesTitle}</h2></div>
      <div><p>{c.capabilitiesBody}</p></div>
    </div>
    <div className="cv-showcase-stage" ref={stage} onFocusCapture={event => {
      const slide = (event.target as HTMLElement).closest('.cv-feature-slide')
      if (slide && track.current) goTo(Array.from(track.current.children).indexOf(slide))
    }}>
      <div className="cv-showcase-stage-bar cv-wrap">
        <span className="cv-showcase-stage-label">{c.capabilitiesEyebrow}</span>
        <div className="cv-showcase-controls">
          <button type="button" onClick={() => goTo(active - 1)} disabled={active === 0} aria-label={ui.previous}><ChevronLeft aria-hidden="true" /></button>
          <span aria-live="polite">{String(active + 1).padStart(2, '0')} / {String(features.length).padStart(2, '0')}</span>
          <button type="button" onClick={() => goTo(active + 1)} disabled={active === features.length - 1} aria-label={ui.next}><ChevronRight aria-hidden="true" /></button>
        </div>
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
    </div>
  </section>
}
