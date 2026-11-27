'use client'

import { useEffect, useRef } from 'react'
import type { Locale } from '@/content/site'
import { ProductCapture } from './ProductCapture'
import { startDesktopMotion } from './desktopMotion'

const copy = {
  en: { backdrop: 'EVERY ORDER. CONNECTED.', title: 'One service day, one clear thread.', body: 'The people at the counter, in the kitchen and behind the business can follow the work in one restaurant workspace.', alt: 'Actual SERVIO point of sale screen' },
  ar: { backdrop: 'كل طلب في مسار واحد', title: 'يوم الشغل كله في مسار واضح.', body: 'الكاشير والمطبخ وإدارة المطعم يقدروا يتابعوا خطوات الشغل من مساحة عمل واحدة.', alt: 'شاشة فعلية لنقطة بيع سيرفيو' },
} as const

export function LayeredTypography({ locale }: { locale: Locale }) {
  const root = useRef<HTMLElement>(null)
  const c = copy[locale]

  useEffect(() => startDesktopMotion(({ gsap }) => {
    if (!root.current) return
    const context = gsap.context(() => {
      const background = root.current!.querySelector('.cv-layered-backdrop')
      const capture = root.current!.querySelector('.cv-layered-capture')
      if (!background || !capture) return
      const timeline = gsap.timeline({ scrollTrigger: { trigger: root.current, start: 'top bottom', end: 'bottom top', scrub: 1 } })
      timeline.fromTo(background, { xPercent: locale === 'ar' ? -8 : 8 }, { xPercent: locale === 'ar' ? 8 : -8, ease: 'none' }, 0)
      timeline.fromTo(capture, { yPercent: 9 }, { yPercent: -9, ease: 'none' }, 0)
    }, root.current)
    return () => context.revert()
  }), [locale])

  return <section className="cv-layered" ref={root} aria-labelledby="cv-layered-title">
    <span className="cv-layered-backdrop" aria-hidden="true">{c.backdrop}</span>
    <div className="cv-wrap cv-layered-inner">
      <div className="cv-layered-copy"><h2 id="cv-layered-title">{c.title}</h2><p>{c.body}</p></div>
      <div className="cv-layered-capture"><ProductCapture image="pos" alt={c.alt} locale={locale} /></div>
    </div>
  </section>
}
