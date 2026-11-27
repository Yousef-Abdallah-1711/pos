'use client'

import { useRef, useState } from 'react'
import Link from 'next/link'
import { ArrowUpRight } from 'lucide-react'
import { useGSAP } from '@gsap/react'
import gsap from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import type { Locale } from '@/content/site'
import { getFeatures } from '@/content/features'

gsap.registerPlugin(useGSAP, ScrollTrigger)

const copy = {
  en: { eyebrow: 'A connected workspace', title: 'The work belongs together.', body: 'These tools live in the same restaurant workspace. The lines show shared context, not automatic data transfer between every module.', explore: 'Explore this area', member: 'Workspace module', center: 'SERVIO workspace' },
  ar: { eyebrow: 'مساحة عمل مترابطة', title: 'أجزاء الشغل مكانها مع بعض.', body: 'الأدوات دي موجودة في مساحة عمل المطعم نفسها. الخطوط توضح انتماءها للمساحة، مش انتقال بيانات تلقائي بين كل جزء والتاني.', explore: 'اكتشف الجزء ده', member: 'جزء من مساحة العمل', center: 'مساحة عمل سيرفيو' },
} as const

export function BusinessEcosystem({ locale }: { locale: Locale }) {
  const root = useRef<HTMLElement>(null)
  const [active, setActive] = useState(0)
  const features = getFeatures(locale)
  const c = copy[locale]

  useGSAP(() => {
    if (!root.current) return
    const media = gsap.matchMedia()
    media.add('(min-width: 901px) and (prefers-reduced-motion: no-preference)', () => {
      const lines = root.current!.querySelectorAll('.cv-ecosystem-lines line')
      const tween = gsap.fromTo(lines,
        { strokeDasharray: 150, strokeDashoffset: 150 },
        { strokeDashoffset: 0, stagger: .08, ease: 'none', scrollTrigger: { trigger: root.current, start: 'top 70%', end: 'center 35%', scrub: .7 } })
      return () => tween.kill()
    })
    return () => media.revert()
  }, { scope: root, dependencies: [locale] })

  return <section className="cv-section cv-ecosystem" ref={root} aria-labelledby="cv-ecosystem-title">
    <div className="cv-wrap">
      <div className="cv-section-heading"><div><p className="cv-eyebrow">{c.eyebrow}</p><h2 id="cv-ecosystem-title">{c.title}</h2></div><p>{c.body}</p></div>
      <div className="cv-ecosystem-layout">
        <div className="cv-ecosystem-map">
          <svg className="cv-ecosystem-lines" viewBox="0 0 800 600" preserveAspectRatio="none" aria-hidden="true">
            <line x1="400" y1="300" x2="165" y2="95" /><line x1="400" y1="300" x2="400" y2="70" />
            <line x1="400" y1="300" x2="635" y2="95" /><line x1="400" y1="300" x2="150" y2="300" />
            <line x1="400" y1="300" x2="650" y2="300" /><line x1="400" y1="300" x2="230" y2="510" />
            <line x1="400" y1="300" x2="570" y2="510" />
          </svg>
          <div className="cv-ecosystem-center" aria-label={c.center}>SERVIO</div>
          {features.map((feature, index) => <button key={feature.slug} type="button" className={`cv-ecosystem-node cv-ecosystem-node--${index + 1}`} aria-pressed={active === index} onClick={() => setActive(index)}>{feature.title}</button>)}
        </div>
        <div className="cv-ecosystem-detail" aria-live="polite">
          <p className="cv-eyebrow">{c.member}</p>
          <h3>{features[active].title}</h3>
          <p>{features[active].description}</p>
          <Link href={`/${locale}/features/${features[active].slug}`}>{c.explore}<ArrowUpRight size={18} aria-hidden="true" /></Link>
        </div>
      </div>
    </div>
  </section>
}
