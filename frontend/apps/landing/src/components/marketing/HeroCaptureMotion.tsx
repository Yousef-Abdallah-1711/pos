'use client'

import { useRef, type ReactNode } from 'react'
import { useGSAP } from '@gsap/react'
import gsap from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'

gsap.registerPlugin(useGSAP, ScrollTrigger)

export function HeroCaptureMotion({ children }: { children: ReactNode }) {
  const root = useRef<HTMLDivElement>(null)

  useGSAP(() => {
    if (!root.current) return
    const media = gsap.matchMedia()
    media.add('(min-width: 901px) and (prefers-reduced-motion: no-preference)', () => {
      const capture = root.current?.querySelector('.cv-hero-capture')
      if (!capture) return
      const tween = gsap.fromTo(capture,
        { rotateX: 11, scale: .93, transformOrigin: 'center top' },
        { rotateX: 0, scale: 1, ease: 'none', scrollTrigger: {
          trigger: root.current,
          start: 'top bottom',
          end: 'bottom top+=120',
          scrub: .65,
          invalidateOnRefresh: true,
        } })
      return () => tween.kill()
    })
    const refresh = () => ScrollTrigger.refresh()
    void document.fonts.ready.then(refresh)
    window.addEventListener('load', refresh, { once: true })
    return () => {
      window.removeEventListener('load', refresh)
      media.revert()
    }
  }, { scope: root })

  return <div className="cv-hero-motion" ref={root}>{children}</div>
}
